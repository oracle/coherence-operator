/*
 * Copyright (c) 2019, 2020 Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package fakes

import (
	"encoding/json"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/operator-framework/operator-sdk/pkg/helm/client"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	kubev3 "helm.sh/helm/v3/pkg/kube"
	"helm.sh/helm/v3/pkg/release"
	storagev3 "helm.sh/helm/v3/pkg/storage"
	driverv3 "helm.sh/helm/v3/pkg/storage/driver"
	"io"
	k8serr "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/helm/pkg/chartutil"
	cpb "k8s.io/helm/pkg/proto/hapi/chart"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"strings"
)

var (
	cfg *action.Configuration
)

// A fake Helm install
type FakeHelm interface {
	// Perform a fake Operator helm install.
	FakeOperatorHelmInstall(mgr *FakeManager, namespace string, values helper.OperatorValues) (*HelmInstallResult, error)
}

// NewFakeHelm creates a FakeHelm from a manager, a ReconcileCoherenceCluster and a ReconcileCoherenceRole
func NewFakeHelm(mgr *FakeManager, clusterReconciler, roleReconciler reconcile.Reconciler, namespace string) (FakeHelm, error) {

	if cfg == nil {
		rcg, err := client.NewRESTClientGetter(mgr, namespace)
		if err != nil {
			return nil, err
		}

		storageBackend := storagev3.Init(driverv3.NewMemory())
		kubeClient := kubev3.New(rcg)

		cfg = &action.Configuration{
			RESTClientGetter: rcg,
			Releases:         storageBackend,
			KubeClient:       kubeClient,
			Log:              func(_ string, _ ...interface{}) {},
		}
	}

	fh := &fakeHelm{
		mgr:               mgr,
		clusterReconciler: clusterReconciler,
		roleReconciler:    roleReconciler,
		cfg:               cfg,
		namespace:         namespace,
	}

	return fh, nil
}

type fakeHelm struct {
	mgr               *FakeManager
	clusterReconciler reconcile.Reconciler
	roleReconciler    reconcile.Reconciler
	namespace         string
	cfg               *action.Configuration
}

func (f *fakeHelm) FakeOperatorHelmInstall(mgr *FakeManager, namespace string, values helper.OperatorValues) (*HelmInstallResult, error) {
	data, err := json.Marshal(values)
	if err != nil {
		return nil, err
	}

	valuesMap := make(map[string]interface{})
	err = json.Unmarshal(data, &valuesMap)
	if err != nil {
		return nil, err
	}

	chartDir, err := helper.FindOperatorHelmChartDir()
	if err != nil {
		return nil, err
	}

	return f.fakeHelmInstall(mgr, namespace, chartDir, valuesMap)
}

func (f *fakeHelm) fakeHelmInstall(mgr *FakeManager, namespace, chartDir string, values map[string]interface{}) (*HelmInstallResult, error) {
	var err error

	chart, err := loader.LoadDir(chartDir)
	if err != nil {
		return nil, err
	}

	install := action.NewInstall(f.cfg)
	install.DryRun = true
	install.Namespace = namespace
	install.ReleaseName = "operator"

	r, err := install.Run(chart, values)
	if err != nil {
		return nil, err
	}

	result, err := f.parseHelmManifest(mgr, r)
	if err != nil {
		return nil, err
	}

	return result, nil
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

func (f *fakeHelm) parseHelmManifest(mgr *FakeManager, release *release.Release) (*HelmInstallResult, error) {
	resources := make(map[schema.GroupVersionResource]map[string]runtime.Object)
	s := mgr.GetScheme()
	decoder := scheme.Codecs.UniversalDecoder()

	parts := strings.Split(release.Manifest, "\n---\n")
	list := make([]runtime.Object, len(parts))
	ownerRefs := make([]metav1.OwnerReference, 0)

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

			// remove owner references
			u.SetOwnerReferences(ownerRefs)
			data, err := yaml.Marshal(u.Object)
			if err != nil {
				return nil, err
			}

			o, err := s.New(u.GroupVersionKind())
			if err != nil {
				return nil, err
			}
			_, _, err = decoder.Decode(data, nil, o)
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
