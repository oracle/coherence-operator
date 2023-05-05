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
	// assert that the StatefulSet is as expected
	assertJobCreation(t, deployment, expected)
}
