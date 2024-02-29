/*
 * Copyright (c) 2020, 2024, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1_test

import (
	coh "github.com/oracle/coherence-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/ptr"
	"testing"
)

func TestCreateJobFromMinimalRoleSpec(t *testing.T) {
	// Create minimal spec spec
	spec := coh.CoherenceResourceSpec{}
	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	expected := createMinimalExpectedJob(deployment)
	// assert that the Job is as expected
	assertJobCreation(t, deployment, expected)
}

func TestCreateJobWithReplicas(t *testing.T) {
	// Create minimal spec spec
	spec := coh.CoherenceResourceSpec{
		Replicas: ptr.To(int32(19)),
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)

	// Create expected Job
	expected := createMinimalExpectedJob(deployment)
	expected.Spec.Parallelism = ptr.To(int32(19))

	// assert that the Job is as expected
	assertJobCreation(t, deployment, expected)
}

func TestCreateJobWithReplicasAndCompletions(t *testing.T) {
	// Create minimal spec spec
	spec := coh.CoherenceJobResourceSpec{
		CoherenceResourceSpec: coh.CoherenceResourceSpec{
			Replicas: ptr.To(int32(19)),
		},
		Completions: ptr.To(int32(21)),
	}

	// Create the test deployment
	deployment := createTestCoherenceJobDeployment(spec)

	// Create expected Job
	expected := createMinimalExpectedJob(deployment)
	expected.Spec.Parallelism = ptr.To(int32(19))
	expected.Spec.Completions = ptr.To(int32(21))

	// assert that the Job is as expected
	assertJobCreation(t, deployment, expected)
}

func TestCreateJobWithReplicasAndSyncedCompletions(t *testing.T) {
	// Create minimal spec spec
	spec := coh.CoherenceJobResourceSpec{
		CoherenceResourceSpec: coh.CoherenceResourceSpec{
			Replicas: ptr.To(int32(19)),
		},
		SyncCompletionsToReplicas: ptr.To(true),
	}

	// Create the test deployment
	deployment := createTestCoherenceJobDeployment(spec)

	// Create expected Job
	expected := createMinimalExpectedJob(deployment)
	expected.Spec.Parallelism = ptr.To(int32(19))
	expected.Spec.Completions = ptr.To(int32(19))

	// assert that the Job is as expected
	assertJobCreation(t, deployment, expected)
}

func TestCreateJobWithReplicasAndSyncedCompletionsOverride(t *testing.T) {
	// Create minimal spec spec
	spec := coh.CoherenceJobResourceSpec{
		CoherenceResourceSpec: coh.CoherenceResourceSpec{
			Replicas: ptr.To(int32(19)),
		},
		Completions:               ptr.To(int32(21)),
		SyncCompletionsToReplicas: ptr.To(true),
	}

	// Create the test deployment
	deployment := createTestCoherenceJobDeployment(spec)

	// Create expected Job
	expected := createMinimalExpectedJob(deployment)
	expected.Spec.Parallelism = ptr.To(int32(19))
	expected.Spec.Completions = ptr.To(int32(19))

	// assert that the Job is as expected
	assertJobCreation(t, deployment, expected)
}

func TestCreateJobWithEnvVarsFrom(t *testing.T) {
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

	spec := coh.CoherenceJobResourceSpec{
		CoherenceResourceSpec: coh.CoherenceResourceSpec{
			Env: []corev1.EnvVar{},
		},
		EnvFrom: from,
	}

	// Create the test deployment
	deployment := createTestCoherenceJobDeployment(spec)
	// Create expected StatefulSet
	expected := createMinimalExpectedJob(deployment)

	addEnvVarsFromToJob(expected, coh.ContainerNameCoherence, from...)

	// assert that the StatefulSet is as expected
	assertJobCreation(t, deployment, expected)
}
