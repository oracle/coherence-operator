/*
 * Copyright (c) 2020, 2025, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1_test

import (
	coh "github.com/oracle/coherence-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	"testing"
)

func TestCreateJobWithCoherenceUtilsEmpty(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		CoherenceUtils: &coh.CoherenceUtilsSpec{},
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	jobExpected := createMinimalExpectedJob(deployment)

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}

func TestCreateJobWithCoherenceUtilsWithImagePullPolicy(t *testing.T) {
	policy := corev1.PullAlways
	spec := coh.CoherenceResourceSpec{
		CoherenceUtils: &coh.CoherenceUtilsSpec{
			ImagePullPolicy: &policy,
		},
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	jobExpected := createMinimalExpectedJob(deployment)
	// Set the expected Operator image pull policy
	jobExpected.Spec.Template.Spec.InitContainers[0].ImagePullPolicy = policy

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}
