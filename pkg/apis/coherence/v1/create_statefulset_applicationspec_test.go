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

func TestCreateStatefulSetWithApplicationType(t *testing.T) {
	// Create a spec with an ApplicationSpec with an application type
	spec := coh.CoherenceDeploymentSpec{
		Application: &coh.ApplicationSpec{
			Type: stringPtr("foo"),
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	// Add the expected environment variables
	addEnvVars(stsExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: "APP_TYPE", Value: "foo"})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithApplicationMain(t *testing.T) {
	// Create a spec with an ApplicationSpec with a main
	mainClass := "com.tangosol.net.CacheFactory"
	spec := coh.CoherenceDeploymentSpec{
		Application: &coh.ApplicationSpec{
			Main: stringPtr(mainClass),
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	// Add the expected environment variables
	addEnvVars(stsExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: "COH_MAIN_CLASS", Value: mainClass})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithApplicationMainArgs(t *testing.T) {
	// Create a spec with an ApplicationSpec with a main
	spec := coh.CoherenceDeploymentSpec{
		Application: &coh.ApplicationSpec{
			Args: []string{"arg1", "arg2"},
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	// Add the expected environment variables
	addEnvVars(stsExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: coh.EnvVarAppMainArgs, Value: "arg1 arg2"})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithApplicationEmptyMainArgs(t *testing.T) {
	// Create a spec with an ApplicationSpec with a main
	spec := coh.CoherenceDeploymentSpec{
		Application: &coh.ApplicationSpec{
			Args: []string{},
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithWorkingDirectory(t *testing.T) {
	// Create a spec with an ApplicationSpec with an application directory
	dir := "/home/foo/app"
	spec := coh.CoherenceDeploymentSpec{
		Application: &coh.ApplicationSpec{
			WorkingDir: &dir,
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	addEnvVars(stsExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: coh.EnvVarCohAppDir, Value: dir})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}
