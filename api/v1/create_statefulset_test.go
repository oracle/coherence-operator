/*
 * Copyright (c) 2020, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1_test

import (
	coh "github.com/oracle/coherence-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"
	"testing"
)

func TestCreateStatefulSetFromMinimalRoleSpec(t *testing.T) {
	// Create minimal spec spec
	spec := coh.CoherenceResourceSpec{}
	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithName(t *testing.T) {
	// create a spec with a name
	spec := coh.CoherenceResourceSpec{
		Role: "data",
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)

	// The sts name should be the full spec name
	stsExpected.Name = deployment.Name

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithReplicas(t *testing.T) {
	// create a spec with a name
	spec := coh.CoherenceResourceSpec{
		Replicas: pointer.Int32Ptr(50),
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithEnvVars(t *testing.T) {
	// create a spec with environment variables
	ev := []corev1.EnvVar{
		{Name: "FOO", Value: "FOO_VAL"},
	}
	spec := coh.CoherenceResourceSpec{
		Env: ev,
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	// add the expected environment variables
	addEnvVars(stsExpected, coh.ContainerNameCoherence, ev...)

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithEmptyEnvVars(t *testing.T) {
	// create a spec with empty environment variables
	spec := coh.CoherenceResourceSpec{
		Env: []corev1.EnvVar{},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithHealthPort(t *testing.T) {
	// create a spec with a custom health port
	spec := coh.CoherenceResourceSpec{
		HealthPort: int32Ptr(210),
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	stsExpected.Spec.Template.Spec.Containers[0].ReadinessProbe.HTTPGet.Port = intstr.FromInt(210)
	stsExpected.Spec.Template.Spec.Containers[0].LivenessProbe.HTTPGet.Port = intstr.FromInt(210)

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithLabels(t *testing.T) {
	// create a spec with empty environment variables
	labels := make(map[string]string)
	labels["foo"] = "foo-label"
	labels["bar"] = "bar-label"

	spec := coh.CoherenceResourceSpec{
		Labels: labels,
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	for k, v := range labels {
		stsExpected.Spec.Template.Labels[k] = v
	}

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithAnnotations(t *testing.T) {
	// create a spec with empty environment variables
	annotations := make(map[string]string)
	annotations["foo"] = "foo-annotation"
	annotations["bar"] = "bar-annotation"

	spec := coh.CoherenceResourceSpec{
		Annotations: annotations,
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	if stsExpected.Spec.Template.Annotations == nil {
		stsExpected.Spec.Template.Annotations = make(map[string]string)
	}
	for k, v := range annotations {
		stsExpected.Spec.Template.Annotations[k] = v
	}

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithResources(t *testing.T) {
	res := corev1.ResourceRequirements{
		Limits: map[corev1.ResourceName]resource.Quantity{
			corev1.ResourceCPU: resource.MustParse("8"),
		},
		Requests: map[corev1.ResourceName]resource.Quantity{
			corev1.ResourceCPU: resource.MustParse("4"),
		},
	}

	spec := coh.CoherenceResourceSpec{
		Resources: &res,
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	stsExpected.Spec.Template.Spec.Containers[0].Resources = res

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithAffinity(t *testing.T) {
	// Create a test affinity spec
	sel := metav1.LabelSelector{
		MatchLabels: map[string]string{"Foo": "Bar"},
	}
	affinity := corev1.Affinity{
		PodAffinity: &corev1.PodAffinity{
			RequiredDuringSchedulingIgnoredDuringExecution: []corev1.PodAffinityTerm{
				{
					LabelSelector: &sel,
				},
			},
		},
	}

	// Create the spec with the affinity spec
	spec := coh.CoherenceResourceSpec{
		Affinity: &affinity,
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	stsExpected.Spec.Template.Spec.Affinity = &affinity

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithNodeSelector(t *testing.T) {
	selector := make(map[string]string)
	selector["foo"] = "foo-label"
	selector["bar"] = "bar-label"

	// Create the spec with the node selector
	spec := coh.CoherenceResourceSpec{
		NodeSelector: selector,
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	stsExpected.Spec.Template.Spec.NodeSelector = selector

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithTolerations(t *testing.T) {
	tolerations := []corev1.Toleration{
		{
			Key:      "Foo",
			Operator: corev1.TolerationOpEqual,
			Value:    "Bar",
		},
	}

	spec := coh.CoherenceResourceSpec{
		Tolerations: tolerations,
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	stsExpected.Spec.Template.Spec.Tolerations = tolerations

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithSecurityContext(t *testing.T) {
	ctx := corev1.PodSecurityContext{
		RunAsUser:    pointer.Int64Ptr(1000),
		RunAsNonRoot: boolPtr(true),
	}

	spec := coh.CoherenceResourceSpec{
		SecurityContext: &ctx,
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	stsExpected.Spec.Template.Spec.SecurityContext = &ctx

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithShareProcessNamespaceFalse(t *testing.T) {
	spec := coh.CoherenceResourceSpec{
		ShareProcessNamespace: boolPtr(false),
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	stsExpected.Spec.Template.Spec.ShareProcessNamespace = boolPtr(false)

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithShareProcessNamespaceTrue(t *testing.T) {
	spec := coh.CoherenceResourceSpec{
		ShareProcessNamespace: boolPtr(true),
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	stsExpected.Spec.Template.Spec.ShareProcessNamespace = boolPtr(true)

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithHostIPCFalse(t *testing.T) {
	spec := coh.CoherenceResourceSpec{
		HostIPC: boolPtr(false),
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	stsExpected.Spec.Template.Spec.HostIPC = false

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithHostIPCNamespaceTrue(t *testing.T) {
	spec := coh.CoherenceResourceSpec{
		HostIPC: boolPtr(true),
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	stsExpected.Spec.Template.Spec.HostIPC = true

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}
