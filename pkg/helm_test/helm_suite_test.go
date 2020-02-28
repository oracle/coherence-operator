/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

// The helm_test package contains tests that take a CoherenceCluster and
// pass it through the operator controllers to verify that the resulting
// kubernetes resources generated by the Helm install are correct.
package helm_test

import (
	cohv1 "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	"github.com/oracle/coherence-operator/pkg/controller/coherencecluster"
	"github.com/oracle/coherence-operator/pkg/controller/coherencerole"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	stubs "github.com/oracle/coherence-operator/pkg/fakes"
	k8serr "k8s.io/apimachinery/pkg/api/errors"
)

// ----- helpers ------------------------------------------------------------

const (
	// The name fo the Coherence container in the StatefulSet
	coherenceContainer   = "coherence"
	applicationContainer = "application"
	fluentdContainer     = "fluentd"
	fluentdImage         = "fluent/fluentd-kubernetes-daemonset:v1.3.3-debian-elasticsearch-1.3"
)

// Use the specified yaml files to create a CoherenceCluster and trigger a fake end-to-end
// reconcile to obtain the resources that would have been created by the Helm operator.
func CreateCluster(yamlFile string) (*stubs.HelmInstallResult, *cohv1.CoherenceCluster, error) {
	namespace := "test-namespace"
	cluster, err := helper.NewCoherenceClusterFromYaml(namespace, yamlFile)
	if err != nil {
		return nil, nil, err
	}

	mgr, err := stubs.NewFakeManager()
	if err != nil {
		return nil, nil, err
	}

	cr := coherencecluster.NewClusterReconciler(mgr)
	rr := coherencerole.NewRoleReconciler(mgr)
	helm, err := stubs.NewFakeHelm(mgr, cr, rr, namespace)
	if err != nil {
		return nil, nil, err
	}

	r, err := helm.HelmInstallFromCoherenceCluster(&cluster)

	return r, &cluster, err
}

// Shared function to find a ConfigMap in a Helm result.
func findConfigMap(result *stubs.HelmInstallResult, name string) (corev1.ConfigMap, error) {
	cm := corev1.ConfigMap{}
	err := result.Get(name, &cm)
	return cm, err
}

// Shared function to find a StatefulSet in a Helm result.
func findStatefulSet(result *stubs.HelmInstallResult, cluster *cohv1.CoherenceCluster, roleName string) (appsv1.StatefulSet, error) {
	name := cluster.GetFullRoleName(roleName)
	sts := appsv1.StatefulSet{}
	err := result.Get(name, &sts)
	return sts, err
}

// Shared function to find a specific container in a StatefulSet spec for a role in a cluster
func findContainerForRole(result *stubs.HelmInstallResult, cluster *cohv1.CoherenceCluster, roleName string, containerName string) (corev1.Container, error) {
	sts, err := findStatefulSet(result, cluster, roleName)
	if err != nil {
		return corev1.Container{}, err
	}
	return findContainer(sts, containerName)
}

// Shared function to find a specific container in a StatefulSet spec
func findContainer(sts appsv1.StatefulSet, name string) (corev1.Container, error) {
	for _, c := range sts.Spec.Template.Spec.Containers {
		if c.Name == name {
			return c, nil
		}
	}

	return corev1.Container{}, k8serr.NewNotFound(schema.GroupResource{Group: "k8s.io/api/core/v1", Resource: "Container"}, name)
}

// Shared function to find a specific init-container in a StatefulSet spec
func findInitContainer(sts appsv1.StatefulSet, name string) (corev1.Container, error) {
	for _, c := range sts.Spec.Template.Spec.InitContainers {
		if c.Name == name {
			return c, nil
		}
	}

	return corev1.Container{}, k8serr.NewNotFound(schema.GroupResource{Group: "k8s.io/api/core/v1", Resource: "Container"}, name)
}

// Shared function to find a specific init-container in a StatefulSet spec for a role in a cluster
func findInitContainerForRole(result *stubs.HelmInstallResult, cluster *cohv1.CoherenceCluster, roleName string, containerName string) (corev1.Container, error) {
	sts, err := findStatefulSet(result, cluster, roleName)
	if err != nil {
		return corev1.Container{}, err
	}
	return findInitContainer(sts, containerName)
}

// Shared function to find a specific volume mount in a Container spec
func findVolumeMount(container corev1.Container, name string) (corev1.VolumeMount, error) {
	for _, v := range container.VolumeMounts {
		if v.Name == name {
			return v, nil
		}
	}

	return corev1.VolumeMount{}, k8serr.NewNotFound(schema.GroupResource{Group: "k8s.io/api/core/v1", Resource: "VolumeMount"}, name)
}

// Shared function to find a specific volume in a StatefulSet spec
func findVolume(sts appsv1.StatefulSet, name string) (corev1.Volume, error) {
	for _, v := range sts.Spec.Template.Spec.Volumes {
		if v.Name == name {
			return v, nil
		}
	}

	return corev1.Volume{}, k8serr.NewNotFound(schema.GroupResource{Group: "k8s.io/api/core/v1", Resource: "Volume"}, name)
}

// Shared function to find a specific VolumeClaimTemplate in a StatefulSet spec
func findPersistentVolumeClaim(sts appsv1.StatefulSet, name string) (corev1.PersistentVolumeClaim, error) {
	for _, v := range sts.Spec.VolumeClaimTemplates {
		if v.Name == name {
			return v, nil
		}
	}

	return corev1.PersistentVolumeClaim{}, k8serr.NewNotFound(schema.GroupResource{Group: "k8s.io/api/core/v1", Resource: "PersistentVolumeClaim"}, name)
}
