/*
 * Copyright (c) 2019, 2020 Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package helper

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/operator-framework/operator-sdk/pkg/helm/release"
	"github.com/oracle/coherence-operator/api/v1"
	"github.com/pborman/uuid"
	rel "helm.sh/helm/v3/pkg/release"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/utils/pointer"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// A helper for managing Helm charts.
// This helper uses an internal Helm and Tiller API and does
// not require Tiller to be installed on the k8s cluster.
type HelmHelper struct {
	Manager        manager.Manager
	KubeClient     kubernetes.Interface
	Namespace      string
	managerFactory release.ManagerFactory
	cleanup        []func()
}

// ReleaseManager manages a Helm release. It can install, update, reconcile,
// and uninstall a release.
type HelmReleaseManager struct {
	Manager release.Manager
}

// Obtain a new HelmHelper for managing the specified Helm chart.
// The helper will use the currently configured Kubernetes config.
func NewHelmHelper(chartDir string) (*HelmHelper, error) {
	cfg, _, err := GetKubeconfigAndNamespace("")
	if err != nil {
		return nil, fmt.Errorf("error (1): %v", err)
	}

	namespace := GetTestNamespace()

	mgr, err := createManager(cfg, namespace)
	if err != nil {
		return nil, fmt.Errorf("error (2): %v", err)
	}

	err = v1.AddToScheme(mgr.GetScheme())
	if err != nil {
		return nil, fmt.Errorf("error (3): %v", err)
	}

	f := release.NewManagerFactory(mgr, chartDir)

	kubeclient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("error (4): %v", err)
	}

	// Ensure that the namespace exists
	_, err = kubeclient.CoreV1().Namespaces().Get(context.TODO(), namespace, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			_, err = kubeclient.CoreV1().Namespaces().Create(context.TODO(), &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{Namespace: namespace, Name: namespace},
			}, metav1.CreateOptions{})
			if err != nil {
				return nil, fmt.Errorf("error (5): %v", err)
			}
		} else {
			return nil, fmt.Errorf("error (6): %v", err)
		}
	}

	helper := &HelmHelper{
		Manager:        mgr,
		Namespace:      namespace,
		managerFactory: f,
		KubeClient:     kubeclient,
	}

	return helper, nil
}

// Obtain a new HelmHelper for managing the Coherence Operator Helm chart.
func NewOperatorChartHelper() (*HelmHelper, error) {
	chart, err := FindOperatorHelmChartDir()
	if err != nil {
		return nil, err
	}
	return NewOperatorChartHelperForChart(chart)
}

// Obtain a new HelmHelper for managing the specified Coherence Operator Helm chart.
func NewOperatorChartHelperForChart(chart string) (*HelmHelper, error) {
	f, err := os.Stat(chart)
	if err != nil {
		return nil, err
	}

	if !f.IsDir() {
		return nil, errors.New("Chart location is not a directory: " + chart)
	}

	return NewHelmHelper(chart)
}

// Obtain a new manager for managing a specific Operator Helm release with a release name and values.
func (h *HelmHelper) NewOperatorHelmReleaseManager(releaseName string, values *OperatorValues) (*HelmReleaseManager, error) {
	if values == nil {
		values = &OperatorValues{}
	}

	if values.ImagePullSecrets == nil {
		values.ImagePullSecrets = GetImagePullSecrets()
	}

	values.FullnameOverride = pointer.StringPtr(releaseName)

	data, err := json.Marshal(values)
	if err != nil {
		return nil, err
	}

	m := make(map[string]interface{})
	err = json.Unmarshal(data, &m)
	if err != nil {
		return nil, err
	}

	y, err := yaml.Marshal(values)
	fmt.Printf("Installing Operator chart release %s with values.yaml\n%s\n", releaseName, string(y))

	return h.NewHelmReleaseManager(releaseName, m)
}

// Obtain a new manager for managing a specific Helm release with a release name and values.
func (h *HelmHelper) NewHelmReleaseManager(releaseName string, values map[string]interface{}) (*HelmReleaseManager, error) {
	u := &unstructured.Unstructured{}
	u.SetGroupVersionKind(schema.GroupVersionKind{Group: "coherence.oracle.com", Version: "v1", Kind: "Operator"})
	u.SetNamespace(h.Namespace)
	u.Object["spec"] = values

	uid := uuid.Parse(uuid.New()).String()
	u.SetUID(types.UID(uid))
	u.SetName(releaseName)

	empty := make(map[string]string)
	m, err := h.managerFactory.NewManager(u, empty)
	if err != nil {
		return nil, err
	}

	err = m.Sync(context.TODO())
	if err != nil {
		return nil, err
	}

	return &HelmReleaseManager{Manager: m}, nil
}

// Clean-up the specified Helm release if it has been installed.
// If the Manager UninstallRelease method returns an error it will just be logged.
func (h *HelmHelper) Cleanup(m *HelmReleaseManager) {
	_, err := m.UninstallRelease()
	if err != nil {
		fmt.Printf("Error uninstalling Helm release '%s' due to %s", m.ReleaseName(), err.Error())
	}
}

// ReleaseName returns the name of the release.
func (h *HelmReleaseManager) ReleaseName() string {
	if h == nil {
		return ""
	}
	return h.Manager.ReleaseName()
}

// IsInstalled returns true if the release has been installed.
func (h *HelmReleaseManager) IsInstalled() bool {
	if h == nil {
		return false
	}
	return h.Manager.IsInstalled()
}

// InstallRelease performs an install of the chart.
func (h *HelmReleaseManager) InstallRelease() (*rel.Release, error) {
	if h == nil {
		return nil, errors.New("InstallRelease called on a nil HelmReleaseManager")
	}
	return h.Manager.InstallRelease(context.TODO())
}

// UpdateRelease performs an update of the release.
func (h *HelmReleaseManager) UpdateRelease() (*rel.Release, *rel.Release, error) {
	if h == nil {
		return nil, nil, errors.New("UpdateRelease called on a nil HelmReleaseManager")
	}
	return h.Manager.UpdateRelease(context.TODO())
}

// ReconcileRelease creates or patches resources as necessary to match the
// deployed release's manifest.
func (h *HelmReleaseManager) ReconcileRelease() (*rel.Release, error) {
	if h == nil {
		return nil, errors.New("ReconcileRelease called on a nil HelmReleaseManager")
	}
	return h.Manager.ReconcileRelease(context.TODO())
}

// UninstallRelease performs an uninstall of the release.
func (h *HelmReleaseManager) UninstallRelease() (*rel.Release, error) {
	if h == nil {
		return nil, errors.New("UninstallRelease called on a nil HelmReleaseManager")
	}
	return h.Manager.UninstallRelease(context.TODO())
}

func createManager(cfg *rest.Config, namespace string) (manager.Manager, error) {
	mgr, err := manager.New(cfg, manager.Options{
		Namespace:          namespace,
		MapperProvider:     apiutil.NewDiscoveryRESTMapper,
		LeaderElection:     false,
		MetricsBindAddress: "0",
	})
	if err != nil {
		return nil, err
	}

	return mgr, nil
}
