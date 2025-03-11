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

func TestCreateStatefulSetWithApplicationType(t *testing.T) {
	// Create a spec with an ApplicationSpec with an application type
	spec := coh.CoherenceResourceSpec{
		Application: &coh.ApplicationSpec{
			Type: stringPtr("foo"),
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	// Add the expected environment variables
	addEnvVarsToAll(stsExpected, corev1.EnvVar{Name: coh.EnvVarAppType, Value: "foo"})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithApplicationMain(t *testing.T) {
	// Create a spec with an ApplicationSpec with a main
	mainClass := "com.tangosol.net.CacheFactory"
	spec := coh.CoherenceResourceSpec{
		Application: &coh.ApplicationSpec{
			Main: stringPtr(mainClass),
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	// Add the expected environment variables
	addEnvVarsToAll(stsExpected, corev1.EnvVar{Name: "COHERENCE_OPERATOR_MAIN_CLASS", Value: mainClass})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithApplicationMainArgs(t *testing.T) {
	// Create a spec with an ApplicationSpec with a main
	spec := coh.CoherenceResourceSpec{
		Application: &coh.ApplicationSpec{
			Args: []string{"arg1", "arg2"},
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	// Add the expected environment variables
	addEnvVarsToAll(stsExpected, corev1.EnvVar{Name: coh.EnvVarAppMainArgs, Value: "arg1 arg2"})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithApplicationMainArgsEmpty(t *testing.T) {
	// Create a spec with an ApplicationSpec with a main
	spec := coh.CoherenceResourceSpec{
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
	spec := coh.CoherenceResourceSpec{
		Application: &coh.ApplicationSpec{
			WorkingDir: &dir,
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	addEnvVarsToAll(stsExpected, corev1.EnvVar{Name: coh.EnvVarCohAppDir, Value: dir})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}
