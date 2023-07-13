/*
 * Copyright (c) 2023, Oracle and/or its affiliates.
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
		CoherenceUtils: &coh.ImageSpec{},
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	jobExpected := createMinimalExpectedJob(deployment)

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}

func TestCreateJobWithCoherenceUtilsWithImage(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		CoherenceUtils: &coh.ImageSpec{
			Image: stringPtr("utils:1.0"),
		},
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	jobExpected := createMinimalExpectedJob(deployment)
	// Set the expected Operator image name
	jobExpected.Spec.Template.Spec.InitContainers[0].Image = "utils:1.0"

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}

func TestCreateJobWithCoherenceUtilsWithImagePullPolicy(t *testing.T) {
	policy := corev1.PullAlways
	spec := coh.CoherenceResourceSpec{
		CoherenceUtils: &coh.ImageSpec{
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
