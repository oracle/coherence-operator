/*
 * Copyright (c) 2020, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1_test

import (
	coh "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	corev1 "k8s.io/api/core/v1"
	"testing"
)

func TestCreateStatefulSetWithCoherenceSpecEmpty(t *testing.T) {

	spec := coh.CoherenceDeploymentSpec{
		Coherence: &coh.CoherenceSpec{},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithImage(t *testing.T) {

	spec := coh.CoherenceDeploymentSpec{
		Image: stringPtr("coherence:1.0"),
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	stsExpected.Spec.Template.Spec.Containers[0].Image = "coherence:1.0"

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithImagePullPolicy(t *testing.T) {
	policy := corev1.PullAlways
	spec := coh.CoherenceDeploymentSpec{
		ImagePullPolicy: &policy,
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	stsExpected.Spec.Template.Spec.Containers[0].ImagePullPolicy = policy

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithCoherenceSpecWithStorageEnabledTrue(t *testing.T) {

	spec := coh.CoherenceDeploymentSpec{
		Coherence: &coh.CoherenceSpec{
			StorageEnabled: boolPtr(true),
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	addEnvVars(stsExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: coh.EnvVarCohStorage, Value: "true"})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithCoherenceSpecWithStorageEnabledFalse(t *testing.T) {

	spec := coh.CoherenceDeploymentSpec{
		Coherence: &coh.CoherenceSpec{
			StorageEnabled: boolPtr(false),
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	addEnvVars(stsExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: coh.EnvVarCohStorage, Value: "false"})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithCoherenceSpecWithCacheConfig(t *testing.T) {

	spec := coh.CoherenceDeploymentSpec{
		Coherence: &coh.CoherenceSpec{
			CacheConfig: stringPtr("test-config.xml"),
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	addEnvVars(stsExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: coh.EnvVarCohCacheConfig, Value: "test-config.xml"})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithCoherenceSpecWithOverrideConfig(t *testing.T) {

	spec := coh.CoherenceDeploymentSpec{
		Coherence: &coh.CoherenceSpec{
			OverrideConfig: stringPtr("test-override.xml"),
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	addEnvVars(stsExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: coh.EnvVarCohOverride, Value: "test-override.xml"})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithCoherenceSpecWithLogLevel(t *testing.T) {

	spec := coh.CoherenceDeploymentSpec{
		Coherence: &coh.CoherenceSpec{
			LogLevel: int32Ptr(9),
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	addEnvVars(stsExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: coh.EnvVarCohLogLevel, Value: "9"})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithCoherenceSpecWithExcludeFromWKATrue(t *testing.T) {

	spec := coh.CoherenceDeploymentSpec{
		Coherence: &coh.CoherenceSpec{
			ExcludeFromWKA: boolPtr(true),
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	stsExpected.Spec.Template.Labels[coh.LabelCoherenceWKAMember] = "false"

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithCoherenceSpecWithExcludeFromWKAFalse(t *testing.T) {

	spec := coh.CoherenceDeploymentSpec{
		Coherence: &coh.CoherenceSpec{
			ExcludeFromWKA: boolPtr(false),
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	stsExpected.Spec.Template.Labels[coh.LabelCoherenceWKAMember] = "true"

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}
