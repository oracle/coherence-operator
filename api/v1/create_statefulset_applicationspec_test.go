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
	stsExpected.Spec.Template.Spec.Containers[0].WorkingDir = dir
	stsExpected.Spec.Template.Spec.InitContainers[1].WorkingDir = dir
	addEnvVarsToAll(stsExpected, corev1.EnvVar{Name: coh.EnvVarCohAppDir, Value: dir})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetUseImageEntryPoint(t *testing.T) {
	spec := coh.CoherenceResourceSpec{
		Application: &coh.ApplicationSpec{
			UseImageEntryPoint: ptr.To(true),
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	stsExpected.Spec.Template.Spec.Containers[0].Command = nil
	stsExpected.Spec.Template.Spec.Containers[0].Args = nil
	stsExpected.Spec.Template.Spec.Containers[0].Env = append(stsExpected.Spec.Template.Spec.Containers[0].Env,
		corev1.EnvVar{Name: coh.EnvVarJdkOptions, Value: "@/coherence-operator/utils/coherence-entrypoint-args.txt"})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetUseImageEntryPointAndUseJdkOptsFalse(t *testing.T) {
	spec := coh.CoherenceResourceSpec{
		Application: &coh.ApplicationSpec{
			UseImageEntryPoint: ptr.To(true),
			UseJdkJavaOptions:  ptr.To(false),
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	stsExpected.Spec.Template.Spec.Containers[0].Command = nil
	stsExpected.Spec.Template.Spec.Containers[0].Args = nil

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetUseImageEntryPointWithExistingJdkOpts(t *testing.T) {
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
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	stsExpected.Spec.Template.Spec.Containers[0].Command = nil
	stsExpected.Spec.Template.Spec.Containers[0].Args = nil
	stsExpected.Spec.Template.Spec.Containers[0].Env = append(stsExpected.Spec.Template.Spec.Containers[0].Env,
		corev1.EnvVar{Name: coh.EnvVarJdkOptions, Value: "-Dfoo=bar @/coherence-operator/utils/coherence-entrypoint-args.txt"})
	stsExpected.Spec.Template.Spec.InitContainers[0].Env = append(stsExpected.Spec.Template.Spec.InitContainers[0].Env,
		corev1.EnvVar{Name: coh.EnvVarJdkOptions, Value: "-Dfoo=bar"})
	stsExpected.Spec.Template.Spec.InitContainers[1].Env = append(stsExpected.Spec.Template.Spec.InitContainers[1].Env,
		corev1.EnvVar{Name: coh.EnvVarJdkOptions, Value: "-Dfoo=bar"})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithApplicationEntryPoint(t *testing.T) {
	spec := coh.CoherenceResourceSpec{
		Application: &coh.ApplicationSpec{
			EntryPoint: []string{"foo", "bar"},
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	stsExpected.Spec.Template.Spec.Containers[0].Command = []string{"foo", "bar"}

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetUseImageEntryPointWithAltJavaOptEnvVar(t *testing.T) {
	spec := coh.CoherenceResourceSpec{
		Application: &coh.ApplicationSpec{
			UseImageEntryPoint:      ptr.To(true),
			AlternateJdkJavaOptions: ptr.To("ALT_OPTS"),
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	stsExpected.Spec.Template.Spec.Containers[0].Command = nil
	stsExpected.Spec.Template.Spec.Containers[0].Args = nil
	stsExpected.Spec.Template.Spec.Containers[0].Env = append(stsExpected.Spec.Template.Spec.Containers[0].Env,
		corev1.EnvVar{Name: coh.EnvVarJdkOptions, Value: "@/coherence-operator/utils/coherence-entrypoint-args.txt"},
		corev1.EnvVar{Name: "ALT_OPTS", Value: "@/coherence-operator/utils/coherence-entrypoint-args.txt"})
	stsExpected.Spec.Template.Spec.InitContainers[0].Env = append(stsExpected.Spec.Template.Spec.InitContainers[0].Env,
		corev1.EnvVar{Name: "ALT_OPTS", Value: "@/coherence-operator/utils/coherence-entrypoint-args.txt"})
	stsExpected.Spec.Template.Spec.InitContainers[1].Env = append(stsExpected.Spec.Template.Spec.InitContainers[1].Env,
		corev1.EnvVar{Name: "ALT_OPTS", Value: "@/coherence-operator/utils/coherence-entrypoint-args.txt"})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetUseImageEntryPointAndUseJdkOptsFalseWithAltJavaOptEnvVar(t *testing.T) {
	spec := coh.CoherenceResourceSpec{
		Application: &coh.ApplicationSpec{
			UseImageEntryPoint:      ptr.To(true),
			UseJdkJavaOptions:       ptr.To(false),
			AlternateJdkJavaOptions: ptr.To("ALT_OPTS"),
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	stsExpected.Spec.Template.Spec.Containers[0].Command = nil
	stsExpected.Spec.Template.Spec.Containers[0].Args = nil
	stsExpected.Spec.Template.Spec.Containers[0].Env = append(stsExpected.Spec.Template.Spec.Containers[0].Env,
		corev1.EnvVar{Name: "ALT_OPTS", Value: "@/coherence-operator/utils/coherence-entrypoint-args.txt"})
	stsExpected.Spec.Template.Spec.InitContainers[0].Env = append(stsExpected.Spec.Template.Spec.InitContainers[0].Env,
		corev1.EnvVar{Name: "ALT_OPTS", Value: "@/coherence-operator/utils/coherence-entrypoint-args.txt"})
	stsExpected.Spec.Template.Spec.InitContainers[1].Env = append(stsExpected.Spec.Template.Spec.InitContainers[1].Env,
		corev1.EnvVar{Name: "ALT_OPTS", Value: "@/coherence-operator/utils/coherence-entrypoint-args.txt"})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}
