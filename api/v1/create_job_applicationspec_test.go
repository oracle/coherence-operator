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

func TestCreateJobWithApplicationType(t *testing.T) {
	// Create a spec with an ApplicationSpec with an application type
	spec := coh.CoherenceResourceSpec{
		Application: &coh.ApplicationSpec{
			Type: stringPtr("foo"),
		},
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	jobExpected := createMinimalExpectedJob(deployment)
	// Add the expected environment variables
	addEnvVarsToAllJobContainers(jobExpected, corev1.EnvVar{Name: coh.EnvVarAppType, Value: "foo"})

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}

func TestCreateJobWithApplicationMain(t *testing.T) {
	// Create a spec with an ApplicationSpec with a main
	mainClass := "com.tangosol.net.CacheFactory"
	spec := coh.CoherenceResourceSpec{
		Application: &coh.ApplicationSpec{
			Main: stringPtr(mainClass),
		},
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	jobExpected := createMinimalExpectedJob(deployment)
	// Add the expected environment variables
	addEnvVarsToAllJobContainers(jobExpected, corev1.EnvVar{Name: "COHERENCE_OPERATOR_MAIN_CLASS", Value: mainClass})

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}

func TestCreateJobWithApplicationMainArgs(t *testing.T) {
	// Create a spec with an ApplicationSpec with a main
	spec := coh.CoherenceResourceSpec{
		Application: &coh.ApplicationSpec{
			Args: []string{"arg1", "arg2"},
		},
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	jobExpected := createMinimalExpectedJob(deployment)
	// Add the expected environment variables
	addEnvVarsToAllJobContainers(jobExpected, corev1.EnvVar{Name: coh.EnvVarAppMainArgs, Value: "arg1 arg2"})

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}

func TestCreateJobWithApplicationMainArgsEmpty(t *testing.T) {
	// Create a spec with an ApplicationSpec with a main
	spec := coh.CoherenceResourceSpec{
		Application: &coh.ApplicationSpec{
			Args: []string{},
		},
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	jobExpected := createMinimalExpectedJob(deployment)

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}

func TestCreateJobWithWorkingDirectory(t *testing.T) {
	// Create a spec with an ApplicationSpec with an application directory
	dir := "/home/foo/app"
	spec := coh.CoherenceResourceSpec{
		Application: &coh.ApplicationSpec{
			WorkingDir: &dir,
		},
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	jobExpected := createMinimalExpectedJob(deployment)
	addEnvVarsToAllJobContainers(jobExpected, corev1.EnvVar{Name: coh.EnvVarCohAppDir, Value: dir})

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}
