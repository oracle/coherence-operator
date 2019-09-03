/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package fakes

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/operator-framework/operator-sdk/pkg/helm/client"
	"github.com/operator-framework/operator-sdk/pkg/helm/engine"
	cohv1 "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	"github.com/pborman/uuid"
	"io"
	k8serr "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	apitypes "k8s.io/apimachinery/pkg/types"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/helm/pkg/chartutil"
	helmengine "k8s.io/helm/pkg/engine"
	"k8s.io/helm/pkg/kube"
	cpb "k8s.io/helm/pkg/proto/hapi/chart"
	helm "k8s.io/helm/pkg/proto/hapi/services"
	"k8s.io/helm/pkg/storage"
	"k8s.io/helm/pkg/storage/driver"
	"k8s.io/helm/pkg/tiller"
	"k8s.io/helm/pkg/tiller/environment"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"strings"
)

// A fake Helm install
type FakeHelm interface {
	// HelmInstallFromCoherenceCluster takes a CoherenceCluster and passes it through
	// the operator reconciler chain to create a Helm install and returns the result of
	// the install. The result will contain all of the resources created by the Helm install.
	HelmInstallFromCoherenceCluster(cluster *cohv1.CoherenceCluster) (*HelmInstallResult, error)
	// Perform a fake Operator helm install.
	FakeOperatorHelmInstall(mgr *FakeManager, namespace string, values helper.OperatorValues) (*HelmInstallResult, error)
}

// NewFakeHelm creates a FakeHelm from a manager, a ReconcileCoherenceCluster and a ReconcileCoherenceRole
func NewFakeHelm(mgr *FakeManager, clusterReconciler, roleReconciler reconcile.Reconciler) FakeHelm {
	return &fakeHelm{mgr, clusterReconciler, roleReconciler}
}

type fakeHelm struct {
	mgr               *FakeManager
	clusterReconciler reconcile.Reconciler
	roleReconciler    reconcile.Reconciler
}

func (f *fakeHelm) HelmInstallFromCoherenceCluster(cluster *cohv1.CoherenceCluster) (*HelmInstallResult, error) {
	if cluster == nil {
		return &HelmInstallResult{}, nil
	}

	_ = f.mgr.GetClient().Create(context.TODO(), cluster)

	clusterRequest := reconcile.Request{
		NamespacedName: apitypes.NamespacedName{
			Namespace: cluster.Namespace,
			Name:      cluster.Name,
		},
	}

	_, err := f.clusterReconciler.Reconcile(clusterRequest)
	if err != nil {
		return nil, err
	}

	list := f.mgr.GetCoherenceRoles(cluster.GetNamespace())

	var result *HelmInstallResult

	for _, role := range list.Items {
		roleRequest := reconcile.Request{
			NamespacedName: apitypes.NamespacedName{
				Namespace: role.Namespace,
				Name:      role.Name,
			},
		}

		_, err := f.roleReconciler.Reconcile(roleRequest)
		if err != nil {
			return nil, err
		}

		values := f.mgr.AssertCoherenceInternalExists(role.Namespace, role.Name)

		r, err := f.FakeHCoherenceHelmInstall(f.mgr, cluster.GetNamespace(), values)
		if err != nil {
			return nil, err
		}

		if result == nil {
			result = r
		} else {
			result.Merge(r)
		}
	}

	return result, nil
}

func (f *fakeHelm) FakeHCoherenceHelmInstall(mgr *FakeManager, namespace string, values *unstructured.Unstructured) (*HelmInstallResult, error) {
	data, err := yaml.Marshal(values.Object["spec"])
	if err != nil {
		return nil, err
	}

	chartDir, err := helper.FindCoherenceHelmChartDir()
	if err != nil {
		return nil, err
	}

	return f.fakeHelmInstall(mgr, namespace, chartDir, data)
}

func (f *fakeHelm) FakeOperatorHelmInstall(mgr *FakeManager, namespace string, values helper.OperatorValues) (*HelmInstallResult, error) {
	data, err := values.ToYaml()
	if err != nil {
		return nil, err
	}

	chartDir, err := helper.FindOperatorHelmChartDir()
	if err != nil {
		return nil, err
	}

	return f.fakeHelmInstall(mgr, namespace, chartDir, data)
}

func (f *fakeHelm) fakeHelmInstall(mgr *FakeManager, namespace, chartDir string, values []byte) (*HelmInstallResult, error) {
	storageBackend := storage.Init(driver.NewMemory())

	cfg, _, err := helper.GetKubeconfigAndNamespace("")
	if err != nil {
		return nil, err
	}

	mgrReal, err := manager.New(cfg, manager.Options{
		Namespace:      namespace,
		MapperProvider: apiutil.NewDiscoveryRESTMapper,
		LeaderElection: false,
	})
	if err != nil {
		return nil, err
	}

	tillerKubeClient, err := client.NewFromManager(mgrReal)
	if err != nil {
		return nil, err
	}

	u := &unstructured.Unstructured{}
	u.SetGroupVersionKind(schema.GroupVersionKind{Group: "coherence.oracle.com", Version: "v1", Kind: "Operator"})
	u.SetNamespace(namespace)

	uid := uuid.Parse(uuid.New()).String()
	u.SetUID(types.UID(uid))
	u.SetName("test")

	releaseServer, err := f.getReleaseServer(u, storageBackend, tillerKubeClient)
	if err != nil {
		return nil, err
	}

	chart, err := f.loadChart(chartDir)
	if err != nil {
		return nil, err
	}

	chart.Dependencies = make([]*cpb.Chart, 0)

	config := &cpb.Config{Raw: string(values)}

	dryRunReq := &helm.InstallReleaseRequest{
		Name:         "operator",
		Chart:        chart,
		Values:       config,
		DryRun:       true,
		Namespace:    namespace,
		DisableHooks: true,
	}

	response, err := releaseServer.InstallRelease(context.TODO(), dryRunReq)
	if err != nil {
		return nil, err
	}

	result, err := f.parseHelmManifest(mgr, response)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (f *fakeHelm) getReleaseServer(cr *unstructured.Unstructured, storageBackend *storage.Storage, tillerKubeClient *kube.Client) (*tiller.ReleaseServer, error) {
	controllerRef := metav1.NewControllerRef(cr, cr.GroupVersionKind())
	ownerRefs := []metav1.OwnerReference{
		*controllerRef,
	}
	baseEngine := helmengine.New()
	e := engine.NewOwnerRefEngine(baseEngine, ownerRefs)
	var ey environment.EngineYard = map[string]environment.Engine{
		environment.GoTplEngine: e,
	}
	env := &environment.Environment{
		EngineYard: ey,
		Releases:   storageBackend,
		KubeClient: tillerKubeClient,
	}
	kubeconfig, err := tillerKubeClient.ToRESTConfig()
	if err != nil {
		return nil, err
	}
	cs, err := clientset.NewForConfig(kubeconfig)
	if err != nil {
		return nil, err
	}

	return tiller.NewReleaseServer(env, cs, false), nil
}

func (f *fakeHelm) loadChart(chartDir string) (*cpb.Chart, error) {
	// chart is mutated by the call to processRequirements,
	// so we need to reload it from disk every time.
	chart, err := chartutil.LoadDir(chartDir)
	if err != nil {
		return nil, fmt.Errorf("failed to load chart: %s", err)
	}

	return chart, nil
}

func (f *fakeHelm) parseHelmManifest(mgr *FakeManager, response *helm.InstallReleaseResponse) (*HelmInstallResult, error) {
	resources := make(map[schema.GroupVersionResource]map[string]runtime.Object)
	s := mgr.GetScheme()
	decoder := scheme.Codecs.UniversalDecoder()

	parts := strings.Split(response.Release.Manifest, "\n---\n")
	list := make([]runtime.Object, len(parts))
	index := 0
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			u := unstructured.Unstructured{}
			err := yaml.Unmarshal([]byte(trimmed), &u)
			if err != nil {
				return nil, err
			}
			gvr, _ := meta.UnsafeGuessKindToResource(u.GroupVersionKind())

			o, err := s.New(u.GroupVersionKind())
			if err != nil {
				return nil, err
			}
			_, _, err = decoder.Decode([]byte(trimmed), nil, o)
			if err != nil {
				return nil, err
			}

			m, ok := resources[gvr]
			if !ok {
				m = make(map[string]runtime.Object)
			}
			list[index] = o
			index++
			m[u.GetName()] = o
			resources[gvr] = m
		}
	}

	ordered := list[0:index]
	return &HelmInstallResult{resources: resources, ordered: ordered, mgr: mgr, decoder: decoder}, nil
}

type HelmInstallResult struct {
	resources map[schema.GroupVersionResource]map[string]runtime.Object
	ordered   []runtime.Object
	mgr       *FakeManager
	decoder   runtime.Decoder
}

type HelmInstallResultFilter func(runtime.Object) bool

func (h *HelmInstallResult) ToString(filter HelmInstallResultFilter, w io.Writer) error {
	var sep = []byte("\n---\n")

	for _, res := range h.ordered {
		if filter == nil || filter(res) {
			_, err := w.Write(sep)
			if err != nil {
				return err
			}

			d, err := yaml.Marshal(res)
			if err != nil {
				return err
			}
			_, err = w.Write(d)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (h *HelmInstallResult) Size() int {
	if h == nil {
		return 0
	}
	return len(h.ordered)
}

func (h *HelmInstallResult) Merge(other *HelmInstallResult) {
	if h == nil || other == nil {
		return
	}

	if h.resources == nil {
		h.resources = make(map[schema.GroupVersionResource]map[string]runtime.Object)
	}

	if other.resources != nil {
		for k, v := range other.resources {
			m, ok := h.resources[k]
			if !ok {
				m = make(map[string]runtime.Object)
			}
			for km, vm := range v {
				m[km] = vm
			}
			h.resources[k] = m
		}
	}

	if h.ordered == nil {
		h.ordered = []runtime.Object{}
	}

	h.ordered = append(h.ordered, other.ordered...)
}

func (h *HelmInstallResult) Get(name string, o runtime.Object) error {
	if h == nil {
		return fmt.Errorf("resource '%s' not found", name)
	}

	gvr, err := h.getGVRFromObject(o, h.mgr.GetScheme())
	if err != nil {
		return err
	}

	if h.resources == nil {
		return k8serr.NewNotFound(gvr.GroupResource(), name)
	}

	m, ok := h.resources[gvr]
	if !ok {
		return k8serr.NewNotFound(gvr.GroupResource(), name)
	}

	item, ok := m[name]
	if !ok {
		return k8serr.NewNotFound(gvr.GroupResource(), name)
	}

	j, err := json.Marshal(item)
	if err != nil {
		return err
	}

	_, _, err = h.decoder.Decode(j, nil, o)
	if err != nil {
		return err
	}

	return nil
}

func (h *HelmInstallResult) List(list runtime.Object) error {
	if h == nil || h.resources == nil {
		return nil
	}

	gvk, err := getGVKFromList(list, h.mgr.GetScheme())
	if err != nil {
		return err
	}

	gvr, _ := meta.UnsafeGuessKindToResource(gvk)
	m, ok := h.resources[gvr]
	if ok {
		items := make([]runtime.Object, len(m))
		i := 0
		for _, o := range m {
			items[i] = o.DeepCopyObject()
			i++
		}

		if err := meta.SetList(list, items); err != nil {
			return err
		}
	}

	return nil
}

func (h *HelmInstallResult) getGVRFromObject(obj runtime.Object, scheme *runtime.Scheme) (schema.GroupVersionResource, error) {
	gvk, err := apiutil.GVKForObject(obj, scheme)
	if err != nil {
		return schema.GroupVersionResource{}, err
	}
	gvr, _ := meta.UnsafeGuessKindToResource(gvk)
	return gvr, nil
}

func getGVKFromList(list runtime.Object, scheme *runtime.Scheme) (schema.GroupVersionKind, error) {
	gvk, err := apiutil.GVKForObject(list, scheme)
	if err != nil {
		return schema.GroupVersionKind{}, err
	}

	if gvk.Kind == "List" {
		return schema.GroupVersionKind{}, fmt.Errorf("cannot derive GVK for generic List type %T (kind %q)", list, gvk)
	}

	if !strings.HasSuffix(gvk.Kind, "List") {
		return schema.GroupVersionKind{}, fmt.Errorf("non-list type %T (kind %q) passed as output", list, gvk)
	}
	// we need the non-list GVK, so chop off the "List" from the end of the kind
	gvk.Kind = gvk.Kind[:len(gvk.Kind)-4]
	return gvk, nil
}
