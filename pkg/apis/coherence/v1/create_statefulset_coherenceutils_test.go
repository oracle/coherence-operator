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

func TestCreateStatefulSetWithCoherenceUtilsEmpty(t *testing.T) {

	spec := coh.CoherenceDeploymentSpec{
		CoherenceUtils: &coh.ImageSpec{},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithCoherenceUtilsWithImage(t *testing.T) {

	spec := coh.CoherenceDeploymentSpec{
		CoherenceUtils: &coh.ImageSpec{
			Image: stringPtr("utils:1.0"),
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	// Set the expected utils image name
	stsExpected.Spec.Template.Spec.InitContainers[0].Image = "utils:1.0"

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithCoherenceUtilsWithImagePullPolicy(t *testing.T) {
	policy := corev1.PullAlways
	spec := coh.CoherenceDeploymentSpec{
		CoherenceUtils: &coh.ImageSpec{
			ImagePullPolicy: &policy,
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	// Set the expected utils image pull policy
	stsExpected.Spec.Template.Spec.InitContainers[0].ImagePullPolicy = policy

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}
