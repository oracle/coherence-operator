/*
 * Copyright (c) 2020, 2024, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1_test

import (
	"fmt"
	coh "github.com/oracle/coherence-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/ptr"
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
		Replicas: ptr.To(int32(50)),
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithRackLabel(t *testing.T) {
	spec := coh.CoherenceResourceSpec{
		RackLabel: ptr.To("coherence.oracle.com/test"),
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	url := fmt.Sprintf("%s?nodeLabel=%s", coh.OperatorRackURL, "coherence.oracle.com/test")
	addEnvVars(stsExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: coh.EnvVarCohRack, Value: url})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithSiteLabel(t *testing.T) {
	spec := coh.CoherenceResourceSpec{
		SiteLabel: ptr.To("coherence.oracle.com/test"),
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	url := fmt.Sprintf("%s?nodeLabel=%s", coh.OperatorSiteURL, "coherence.oracle.com/test")
	addEnvVars(stsExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: coh.EnvVarCohSite, Value: url})

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

func TestCreateStatefulSetWithEnvVarsFrom(t *testing.T) {
	cm := corev1.ConfigMapEnvSource{
		LocalObjectReference: corev1.LocalObjectReference{
			Name: "test-vars",
		},
		Optional: ptr.To(true),
	}

	var from []corev1.EnvFromSource
	from = append(from, corev1.EnvFromSource{
		Prefix:       "foo",
		ConfigMapRef: &cm,
	})

	spec := coh.CoherenceStatefulSetResourceSpec{
		CoherenceResourceSpec: coh.CoherenceResourceSpec{
			Env: []corev1.EnvVar{},
		},
		EnvFrom: from,
	}

	// Create the test deployment
	deployment := createTestCoherenceDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)

	addEnvVarsFrom(stsExpected, coh.ContainerNameCoherence, from...)

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
	stsExpected.Spec.Template.Spec.Containers[0].ReadinessProbe.HTTPGet.Port = intstr.FromInt32(210)
	stsExpected.Spec.Template.Spec.Containers[0].LivenessProbe.HTTPGet.Port = intstr.FromInt32(210)

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

func TestCreateStatefulSetWithAnnotationsFromCoherenceResource(t *testing.T) {
	// create a spec with empty environment variables
	annotations := make(map[string]string)
	annotations["key1"] = "value1"
	annotations["key2"] = "value2"

	spec := coh.CoherenceResourceSpec{}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	deployment.SetAnnotations(annotations)

	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	if stsExpected.Annotations == nil {
		stsExpected.Annotations = make(map[string]string)
	}
	for k, v := range annotations {
		stsExpected.Annotations[k] = v
	}

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithAnnotationsOverriddenFromCoherenceResource(t *testing.T) {
	// create a spec with empty environment variables
	annotationsOne := make(map[string]string)
	annotationsOne["key1"] = "value1"
	annotationsOne["key2"] = "value2"

	annotationsTwo := make(map[string]string)
	annotationsTwo["key3"] = "value3"
	annotationsTwo["key4"] = "value4"

	spec := coh.CoherenceStatefulSetResourceSpec{
		StatefulSetAnnotations: annotationsTwo,
	}

	// Create the test deployment
	deployment := createTestCoherenceDeployment(spec)
	deployment.SetAnnotations(annotationsOne)

	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	if stsExpected.Annotations == nil {
		stsExpected.Annotations = make(map[string]string)
	}
	for k, v := range annotationsTwo {
		stsExpected.Annotations[k] = v
	}

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithPodAnnotations(t *testing.T) {
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
		RunAsUser:    ptr.To(int64(1000)),
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

func TestCreateStatefulSetWithContainerSecurityContext(t *testing.T) {
	ctx := corev1.SecurityContext{
		Capabilities: &corev1.Capabilities{
			Add: []corev1.Capability{"foo"},
		},
		RunAsUser:    ptr.To(int64(1000)),
		RunAsNonRoot: boolPtr(true),
	}

	spec := coh.CoherenceResourceSpec{
		ContainerSecurityContext: &ctx,
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	// Add the expected security context to both the init-container and the Coherence container
	stsExpected.Spec.Template.Spec.InitContainers[0].SecurityContext = &ctx
	stsExpected.Spec.Template.Spec.Containers[0].SecurityContext = &ctx

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

func TestCreateStatefulSetWithAppLabel(t *testing.T) {
	spec := coh.CoherenceResourceSpec{
		AppLabel: stringPtr("foo"),
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	stsExpected.Labels["app"] = "foo"
	stsExpected.Spec.Template.Labels["app"] = "foo"

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithActiveDeadlineSeconds(t *testing.T) {
	spec := coh.CoherenceResourceSpec{
		ActiveDeadlineSeconds: ptr.To(int64(19)),
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	stsExpected.Spec.Template.Spec.ActiveDeadlineSeconds = ptr.To(int64(19))

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithEnableServiceLinksFalse(t *testing.T) {
	spec := coh.CoherenceResourceSpec{
		EnableServiceLinks: ptr.To(false),
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	stsExpected.Spec.Template.Spec.EnableServiceLinks = ptr.To(false)

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithEnableServiceLinksTrue(t *testing.T) {
	spec := coh.CoherenceResourceSpec{
		EnableServiceLinks: ptr.To(true),
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	stsExpected.Spec.Template.Spec.EnableServiceLinks = ptr.To(true)

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithPreemptionPolicy(t *testing.T) {
	policy := corev1.PreemptNever
	spec := coh.CoherenceResourceSpec{
		PreemptionPolicy: &policy,
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	stsExpected.Spec.Template.Spec.PreemptionPolicy = &policy

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithPriorityClassName(t *testing.T) {
	spec := coh.CoherenceResourceSpec{
		PriorityClassName: stringPtr("foo"),
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	stsExpected.Spec.Template.Spec.PriorityClassName = "foo"

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithRestartPolicy(t *testing.T) {
	policy := corev1.RestartPolicyOnFailure
	spec := coh.CoherenceResourceSpec{
		RestartPolicy: &policy,
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	stsExpected.Spec.Template.Spec.RestartPolicy = policy

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithRuntimeClassName(t *testing.T) {
	spec := coh.CoherenceResourceSpec{
		RuntimeClassName: stringPtr("foo"),
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	stsExpected.Spec.Template.Spec.RuntimeClassName = stringPtr("foo")

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithSchedulerName(t *testing.T) {
	spec := coh.CoherenceResourceSpec{
		SchedulerName: stringPtr("foo"),
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	stsExpected.Spec.Template.Spec.SchedulerName = "foo"

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithTopologySpreadConstraintsEmpty(t *testing.T) {
	spec := coh.CoherenceResourceSpec{
		TopologySpreadConstraints: []corev1.TopologySpreadConstraint{},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	stsExpected.Spec.Template.Spec.TopologySpreadConstraints = []corev1.TopologySpreadConstraint{}

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithTopologySpreadConstraints(t *testing.T) {
	selector := make(map[string]string)
	selector["foo"] = "bar"

	constraint := corev1.TopologySpreadConstraint{
		MaxSkew:           19,
		TopologyKey:       "foo",
		WhenUnsatisfiable: corev1.DoNotSchedule,
		LabelSelector: &metav1.LabelSelector{
			MatchLabels: selector,
		},
		MinDomains: ptr.To(int32(2)),
	}

	spec := coh.CoherenceResourceSpec{
		TopologySpreadConstraints: []corev1.TopologySpreadConstraint{constraint},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	stsExpected.Spec.Template.Spec.TopologySpreadConstraints = []corev1.TopologySpreadConstraint{constraint}

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithGlobalLabels(t *testing.T) {
	m := make(map[string]string)
	m["one"] = "value-one"
	m["two"] = "value-two"

	spec := coh.CoherenceStatefulSetResourceSpec{
		Global: &coh.GlobalSpec{
			Labels: m,
		},
	}

	// Create the test deployment
	deployment := createTestCoherenceDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	labelsExpected := stsExpected.Labels
	labelsExpected["one"] = "value-one"
	labelsExpected["two"] = "value-two"

	podLabelsExpected := stsExpected.Spec.Template.Labels
	podLabelsExpected["one"] = "value-one"
	podLabelsExpected["two"] = "value-two"

	// assert that the Job is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithGlobalAnnotations(t *testing.T) {
	m := make(map[string]string)
	m["one"] = "value-one"
	m["two"] = "value-two"

	spec := coh.CoherenceStatefulSetResourceSpec{
		Global: &coh.GlobalSpec{
			Annotations: m,
		},
	}

	// Create the test deployment
	deployment := createTestCoherenceDeployment(spec)
	// Create expected Job
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	annExpected := stsExpected.Annotations
	if annExpected == nil {
		annExpected = make(map[string]string)
	}
	annExpected["one"] = "value-one"
	annExpected["two"] = "value-two"
	stsExpected.Annotations = annExpected

	podAnnExpected := stsExpected.Spec.Template.Annotations
	if podAnnExpected == nil {
		podAnnExpected = make(map[string]string)
	}
	podAnnExpected["one"] = "value-one"
	podAnnExpected["two"] = "value-two"
	stsExpected.Spec.Template.Annotations = podAnnExpected

	// assert that the Job is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}
