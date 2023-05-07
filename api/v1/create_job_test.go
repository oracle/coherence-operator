/*
 * Copyright (c) 2019, 2023, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1_test

import (
	coh "github.com/oracle/coherence-operator/api/v1"
	"k8s.io/utils/pointer"
	"testing"
)

func TestCreateJobFromMinimalRoleSpec(t *testing.T) {
	// Create minimal spec spec
	spec := coh.CoherenceResourceSpec{
		RunAsJob: pointer.Bool(true),
	}
	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected Job
	expected := createMinimalExpectedJob(deployment)
	// assert that the Job is as expected
	assertJobCreation(t, deployment, expected)
}

func TestCreateJobWithReplicas(t *testing.T) {
	// Create minimal spec spec
	spec := coh.CoherenceResourceSpec{
		Replicas: pointer.Int32(19),
		RunAsJob: pointer.Bool(true),
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)

	// Create expected Job
	expected := createMinimalExpectedJob(deployment)
	expected.Spec.Parallelism = pointer.Int32(19)

	// assert that the Job is as expected
	assertJobCreation(t, deployment, expected)
}

func TestCreateJobWithReplicasAndCompletions(t *testing.T) {
	// Create minimal spec spec
	spec := coh.CoherenceResourceSpec{
		Replicas: pointer.Int32(19),
		RunAsJob: pointer.Bool(true),
		JobSpec: &coh.CoherenceJob{
			Completions: pointer.Int32(21),
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)

	// Create expected Job
	expected := createMinimalExpectedJob(deployment)
	expected.Spec.Parallelism = pointer.Int32(19)
	expected.Spec.Completions = pointer.Int32(21)

	// assert that the Job is as expected
	assertJobCreation(t, deployment, expected)
}

func TestCreateJobWithReplicasAndSyncedCompletions(t *testing.T) {
	// Create minimal spec spec
	spec := coh.CoherenceResourceSpec{
		Replicas: pointer.Int32(19),
		RunAsJob: pointer.Bool(true),
		JobSpec: &coh.CoherenceJob{
			SyncCompletionsToReplicas: pointer.Bool(true),
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)

	// Create expected Job
	expected := createMinimalExpectedJob(deployment)
	expected.Spec.Parallelism = pointer.Int32(19)
	expected.Spec.Completions = pointer.Int32(19)

	// assert that the Job is as expected
	assertJobCreation(t, deployment, expected)
}

func TestCreateJobWithReplicasAndSyncedCompletionsOverride(t *testing.T) {
	// Create minimal spec spec
	spec := coh.CoherenceResourceSpec{
		Replicas: pointer.Int32(19),
		RunAsJob: pointer.Bool(true),
		JobSpec: &coh.CoherenceJob{
			Completions:               pointer.Int32(21),
			SyncCompletionsToReplicas: pointer.Bool(true),
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)

	// Create expected Job
	expected := createMinimalExpectedJob(deployment)
	expected.Spec.Parallelism = pointer.Int32(19)
	expected.Spec.Completions = pointer.Int32(19)

	// assert that the Job is as expected
	assertJobCreation(t, deployment, expected)
}
