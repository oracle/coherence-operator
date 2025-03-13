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

func TestCreateJobWithContainersEmpty(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		SideCars: []corev1.Container{},
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	jobExpected := createMinimalExpectedJob(deployment)

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}

func TestCreateJobWithOneExtraContainer(t *testing.T) {
	c := corev1.Container{
		Name:  "one",
		Image: "image-one:1.0",
	}
	spec := coh.CoherenceResourceSpec{
		SideCars: []corev1.Container{c},
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)

	// Create expected Job
	jobExpected := createMinimalExpectedJob(deployment)

	// Create expected container
	conExpected := corev1.Container{
		Name:         "one",
		Image:        "image-one:1.0",
		Env:          deployment.Spec.CreateCommonEnv(deployment),
		VolumeMounts: deployment.Spec.CreateCommonVolumeMounts(),
	}

	jobExpected.Spec.Template.Spec.Containers = append(jobExpected.Spec.Template.Spec.Containers, conExpected)

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}

func TestCreateJobWithOneExtraContainerWithOverriddenEnvVar(t *testing.T) {
	c := corev1.Container{
		Name:  "one",
		Image: "image-one:1.0",
		Env: []corev1.EnvVar{
			{Name: "foo", Value: "bar"},
			{Name: coh.EnvVarCohRole, Value: "overridden"},
		},
	}
	spec := coh.CoherenceResourceSpec{
		SideCars: []corev1.Container{c},
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)

	// Create expected Job
	jobExpected := createMinimalExpectedJob(deployment)

	// Create expected container
	env := append(deployment.Spec.CreateCommonEnv(deployment), corev1.EnvVar{Name: "foo", Value: "bar"})
	for i, e := range env {
		if e.Name == coh.EnvVarCohRole {
			env[i] = corev1.EnvVar{Name: coh.EnvVarCohRole, Value: "overridden"}
		}
	}

	conExpected := corev1.Container{
		Name:         "one",
		Image:        "image-one:1.0",
		Env:          env,
		VolumeMounts: deployment.Spec.CreateCommonVolumeMounts(),
	}

	jobExpected.Spec.Template.Spec.Containers = append(jobExpected.Spec.Template.Spec.Containers, conExpected)

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}

func TestCreateJobWithTwoExtraContainers(t *testing.T) {
	c1 := corev1.Container{
		Name:  "one",
		Image: "image-one:1.0",
	}
	c2 := corev1.Container{
		Name:  "two",
		Image: "image-two:1.0",
	}
	spec := coh.CoherenceResourceSpec{
		SideCars: []corev1.Container{c1, c2},
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	jobExpected := createMinimalExpectedJob(deployment)

	// Create expected container1
	conExpected1 := corev1.Container{
		Name:         "one",
		Image:        "image-one:1.0",
		Env:          deployment.Spec.CreateCommonEnv(deployment),
		VolumeMounts: deployment.Spec.CreateCommonVolumeMounts(),
	}
	// Create expected container2
	conExpected2 := corev1.Container{
		Name:         "two",
		Image:        "image-two:1.0",
		Env:          deployment.Spec.CreateCommonEnv(deployment),
		VolumeMounts: deployment.Spec.CreateCommonVolumeMounts(),
	}

	jobExpected.Spec.Template.Spec.Containers = append(jobExpected.Spec.Template.Spec.Containers, conExpected1, conExpected2)

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}

func TestCreateJobWithInitContainersEmpty(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		InitContainers: []corev1.Container{},
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	jobExpected := createMinimalExpectedJob(deployment)

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}

func TestCreateJobWithOneExtraInitContainer(t *testing.T) {
	c := corev1.Container{
		Name:  "one",
		Image: "image-one:1.0",
	}
	spec := coh.CoherenceResourceSpec{
		InitContainers: []corev1.Container{c},
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	jobExpected := createMinimalExpectedJob(deployment)

	// Create expected container
	conExpected := corev1.Container{
		Name:         "one",
		Image:        "image-one:1.0",
		Env:          deployment.Spec.CreateCommonEnv(deployment),
		VolumeMounts: deployment.Spec.CreateCommonVolumeMounts(),
	}

	jobExpected.Spec.Template.Spec.InitContainers = append(jobExpected.Spec.Template.Spec.InitContainers, conExpected)

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}

func TestCreateJobWithTwoExtraInitContainers(t *testing.T) {
	c1 := corev1.Container{
		Name:  "one",
		Image: "image-one:1.0",
	}
	c2 := corev1.Container{
		Name:  "two",
		Image: "image-two:1.0",
	}
	spec := coh.CoherenceResourceSpec{
		InitContainers: []corev1.Container{c1, c2},
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	jobExpected := createMinimalExpectedJob(deployment)

	// Create expected container1
	conExpected1 := corev1.Container{
		Name:         "one",
		Image:        "image-one:1.0",
		Env:          deployment.Spec.CreateCommonEnv(deployment),
		VolumeMounts: deployment.Spec.CreateCommonVolumeMounts(),
	}
	// Create expected container2
	conExpected2 := corev1.Container{
		Name:         "two",
		Image:        "image-two:1.0",
		Env:          deployment.Spec.CreateCommonEnv(deployment),
		VolumeMounts: deployment.Spec.CreateCommonVolumeMounts(),
	}

	jobExpected.Spec.Template.Spec.InitContainers = append(jobExpected.Spec.Template.Spec.InitContainers, conExpected1, conExpected2)

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}

func TestCreateJobWithExtraContainerAndVolume(t *testing.T) {
	c := corev1.Container{
		Name:  "one",
		Image: "image-one:1.0",
	}
	mount := corev1.VolumeMount{
		Name:      "logs",
		MountPath: "/logs",
	}
	vol := corev1.Volume{
		Name: "logs",
		VolumeSource: corev1.VolumeSource{
			EmptyDir: &corev1.EmptyDirVolumeSource{},
		},
	}
	spec := coh.CoherenceResourceSpec{
		SideCars:     []corev1.Container{c},
		VolumeMounts: []corev1.VolumeMount{mount},
		Volumes:      []corev1.Volume{vol},
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)

	// Create expected Job
	jobExpected := createMinimalExpectedJob(deployment)
	jobExpected.Spec.Template.Spec.Volumes = append(jobExpected.Spec.Template.Spec.Volumes, vol)
	jobExpected.Spec.Template.Spec.Containers[0].VolumeMounts = append(jobExpected.Spec.Template.Spec.Containers[0].VolumeMounts, mount)
	jobExpected.Spec.Template.Spec.InitContainers[0].VolumeMounts = append(jobExpected.Spec.Template.Spec.InitContainers[0].VolumeMounts, mount)
	jobExpected.Spec.Template.Spec.InitContainers[1].VolumeMounts = append(jobExpected.Spec.Template.Spec.InitContainers[1].VolumeMounts, mount)

	// Create expected container
	conExpected := corev1.Container{
		Name:         "one",
		Image:        "image-one:1.0",
		Env:          deployment.Spec.CreateCommonEnv(deployment),
		VolumeMounts: append(deployment.Spec.CreateCommonVolumeMounts(), mount),
	}

	jobExpected.Spec.Template.Spec.Containers = append(jobExpected.Spec.Template.Spec.Containers, conExpected)

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}
