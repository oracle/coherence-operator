/*
 * Copyright (c) 2020, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1_test

import (
	coh "github.com/oracle/coherence-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	"testing"
)

func TestCreateStatefulSetWithContainersEmpty(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		SideCars: []corev1.Container{},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithOneExtraContainer(t *testing.T) {
	c := corev1.Container{
		Name:  "one",
		Image: "image-one:1.0",
	}
	spec := coh.CoherenceResourceSpec{
		SideCars: []corev1.Container{c},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)

	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)

	// Create expected container
	conExpected := corev1.Container{
		Name:         "one",
		Image:        "image-one:1.0",
		Env:          deployment.Spec.CreateCommonEnv(deployment),
		VolumeMounts: deployment.Spec.CreateCommonVolumeMounts(),
	}

	stsExpected.Spec.Template.Spec.Containers = append(stsExpected.Spec.Template.Spec.Containers, conExpected)

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithOneExtraContainerWithOverriddenEnvVar(t *testing.T) {
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
	deployment := createTestDeployment(spec)

	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)

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

	stsExpected.Spec.Template.Spec.Containers = append(stsExpected.Spec.Template.Spec.Containers, conExpected)

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithTwoExtraContainers(t *testing.T) {
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
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)

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

	stsExpected.Spec.Template.Spec.Containers = append(stsExpected.Spec.Template.Spec.Containers, conExpected1, conExpected2)

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithInitContainersEmpty(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		InitContainers: []corev1.Container{},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithOneExtraInitContainer(t *testing.T) {
	c := corev1.Container{
		Name:  "one",
		Image: "image-one:1.0",
	}
	spec := coh.CoherenceResourceSpec{
		InitContainers: []corev1.Container{c},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)

	// Create expected container
	conExpected := corev1.Container{
		Name:         "one",
		Image:        "image-one:1.0",
		Env:          deployment.Spec.CreateCommonEnv(deployment),
		VolumeMounts: deployment.Spec.CreateCommonVolumeMounts(),
	}

	stsExpected.Spec.Template.Spec.InitContainers = append(stsExpected.Spec.Template.Spec.InitContainers, conExpected)

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithTwoExtraInitContainers(t *testing.T) {
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
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)

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

	stsExpected.Spec.Template.Spec.InitContainers = append(stsExpected.Spec.Template.Spec.InitContainers, conExpected1, conExpected2)

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithExtraContainerAndVolume(t *testing.T) {
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
	deployment := createTestDeployment(spec)

	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	stsExpected.Spec.Template.Spec.Volumes = append(stsExpected.Spec.Template.Spec.Volumes, vol)
	stsExpected.Spec.Template.Spec.Containers[0].VolumeMounts = append(stsExpected.Spec.Template.Spec.Containers[0].VolumeMounts, mount)

	// Create expected container
	conExpected := corev1.Container{
		Name:         "one",
		Image:        "image-one:1.0",
		Env:          deployment.Spec.CreateCommonEnv(deployment),
		VolumeMounts: append(deployment.Spec.CreateCommonVolumeMounts(), mount),
	}

	stsExpected.Spec.Template.Spec.Containers = append(stsExpected.Spec.Template.Spec.Containers, conExpected)

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}
