/*
 * Copyright (c) 2020, 2025, Oracle and/or its affiliates.
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
	jobExpected.Spec.Template.Spec.Containers[0].WorkingDir = dir
	jobExpected.Spec.Template.Spec.InitContainers[1].WorkingDir = dir
	addEnvVarsToAllJobContainers(jobExpected, corev1.EnvVar{Name: coh.EnvVarCohAppDir, Value: dir})

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}

func TestCreateJobUseImageEntryPoint(t *testing.T) {
	spec := coh.CoherenceResourceSpec{
		Application: &coh.ApplicationSpec{
			UseImageEntryPoint: ptr.To(true),
		},
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	jobExpected := createMinimalExpectedJob(deployment)
	jobExpected.Spec.Template.Spec.Containers[0].Command = nil
	jobExpected.Spec.Template.Spec.Containers[0].Args = nil
	jobExpected.Spec.Template.Spec.Containers[0].Env = append(jobExpected.Spec.Template.Spec.Containers[0].Env,
		corev1.EnvVar{Name: coh.EnvVarJdkOptions, Value: "@/coherence-operator/utils/coherence-entrypoint-args.txt"})

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}

func TestCreateJobUseImageEntryPointAndUseJdkOptsFalse(t *testing.T) {
	spec := coh.CoherenceResourceSpec{
		Application: &coh.ApplicationSpec{
			UseImageEntryPoint: ptr.To(true),
			UseJdkJavaOptions:  ptr.To(false),
		},
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	jobExpected := createMinimalExpectedJob(deployment)
	jobExpected.Spec.Template.Spec.Containers[0].Command = nil
	jobExpected.Spec.Template.Spec.Containers[0].Args = nil

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}

func TestCreateJobUseImageEntryPointWithExistingJdkOpts(t *testing.T) {
	spec := coh.CoherenceResourceSpec{
		Application: &coh.ApplicationSpec{
			UseImageEntryPoint: ptr.To(true),
		},
		Env: []corev1.EnvVar{
			{
				Name:  coh.EnvVarJdkOptions,
				Value: "-Dfoo=bar",
			},
		},
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	jobExpected := createMinimalExpectedJob(deployment)
	jobExpected.Spec.Template.Spec.Containers[0].Command = nil
	jobExpected.Spec.Template.Spec.Containers[0].Args = nil
	jobExpected.Spec.Template.Spec.Containers[0].Env = append(jobExpected.Spec.Template.Spec.Containers[0].Env,
		corev1.EnvVar{Name: coh.EnvVarJdkOptions, Value: "-Dfoo=bar @/coherence-operator/utils/coherence-entrypoint-args.txt"})
	jobExpected.Spec.Template.Spec.InitContainers[0].Env = append(jobExpected.Spec.Template.Spec.InitContainers[0].Env,
		corev1.EnvVar{Name: coh.EnvVarJdkOptions, Value: "-Dfoo=bar"})
	jobExpected.Spec.Template.Spec.InitContainers[1].Env = append(jobExpected.Spec.Template.Spec.InitContainers[1].Env,
		corev1.EnvVar{Name: coh.EnvVarJdkOptions, Value: "-Dfoo=bar"})

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}

func TestCreateJobWithApplicationEntryPoint(t *testing.T) {
	spec := coh.CoherenceResourceSpec{
		Application: &coh.ApplicationSpec{
			EntryPoint: []string{"foo", "bar"},
		},
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	jobExpected := createMinimalExpectedJob(deployment)
	jobExpected.Spec.Template.Spec.Containers[0].Command = []string{"foo", "bar"}

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}

func TestCreateJobUseImageEntryPointWithAltJavaOptEnvVar(t *testing.T) {
	spec := coh.CoherenceResourceSpec{
		Application: &coh.ApplicationSpec{
			UseImageEntryPoint:      ptr.To(true),
			AlternateJdkJavaOptions: ptr.To("ALT_OPTS"),
		},
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	jobExpected := createMinimalExpectedJob(deployment)
	jobExpected.Spec.Template.Spec.Containers[0].Command = nil
	jobExpected.Spec.Template.Spec.Containers[0].Args = nil
	jobExpected.Spec.Template.Spec.Containers[0].Env = append(jobExpected.Spec.Template.Spec.Containers[0].Env,
		corev1.EnvVar{Name: coh.EnvVarJdkOptions, Value: "@/coherence-operator/utils/coherence-entrypoint-args.txt"},
		corev1.EnvVar{Name: "ALT_OPTS", Value: "@/coherence-operator/utils/coherence-entrypoint-args.txt"})
	jobExpected.Spec.Template.Spec.InitContainers[0].Env = append(jobExpected.Spec.Template.Spec.InitContainers[0].Env,
		corev1.EnvVar{Name: "ALT_OPTS", Value: "@/coherence-operator/utils/coherence-entrypoint-args.txt"})
	jobExpected.Spec.Template.Spec.InitContainers[1].Env = append(jobExpected.Spec.Template.Spec.InitContainers[1].Env,
		corev1.EnvVar{Name: "ALT_OPTS", Value: "@/coherence-operator/utils/coherence-entrypoint-args.txt"})

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}

func TestCreateJobUseImageEntryPointAndUseJdkOptsFalseWithAltJavaOptEnvVar(t *testing.T) {
	spec := coh.CoherenceResourceSpec{
		Application: &coh.ApplicationSpec{
			UseImageEntryPoint:      ptr.To(true),
			UseJdkJavaOptions:       ptr.To(false),
			AlternateJdkJavaOptions: ptr.To("ALT_OPTS"),
		},
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	jobExpected := createMinimalExpectedJob(deployment)
	jobExpected.Spec.Template.Spec.Containers[0].Command = nil
	jobExpected.Spec.Template.Spec.Containers[0].Args = nil
	jobExpected.Spec.Template.Spec.Containers[0].Env = append(jobExpected.Spec.Template.Spec.Containers[0].Env,
		corev1.EnvVar{Name: "ALT_OPTS", Value: "@/coherence-operator/utils/coherence-entrypoint-args.txt"})
	jobExpected.Spec.Template.Spec.InitContainers[0].Env = append(jobExpected.Spec.Template.Spec.InitContainers[0].Env,
		corev1.EnvVar{Name: "ALT_OPTS", Value: "@/coherence-operator/utils/coherence-entrypoint-args.txt"})
	jobExpected.Spec.Template.Spec.InitContainers[1].Env = append(jobExpected.Spec.Template.Spec.InitContainers[1].Env,
		corev1.EnvVar{Name: "ALT_OPTS", Value: "@/coherence-operator/utils/coherence-entrypoint-args.txt"})

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}
