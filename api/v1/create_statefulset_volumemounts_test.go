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

func TestCreateStatefulSetWithEmptyVolumeMounts(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		VolumeMounts: []corev1.VolumeMount{},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithOneVolumeMount(t *testing.T) {

	mountOne := corev1.VolumeMount{
		Name:      "volume-one",
		ReadOnly:  true,
		MountPath: "/home/root/one",
	}

	spec := coh.CoherenceResourceSpec{
		VolumeMounts: []corev1.VolumeMount{mountOne},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	stsExpected.Spec.Template.Spec.Containers[0].VolumeMounts = append(stsExpected.Spec.Template.Spec.Containers[0].VolumeMounts, mountOne)
	stsExpected.Spec.Template.Spec.InitContainers[0].VolumeMounts = append(stsExpected.Spec.Template.Spec.InitContainers[0].VolumeMounts, mountOne)
	stsExpected.Spec.Template.Spec.InitContainers[1].VolumeMounts = append(stsExpected.Spec.Template.Spec.InitContainers[1].VolumeMounts, mountOne)

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}
func TestCreateStatefulSetWithTwoVolumeMounts(t *testing.T) {

	mountOne := corev1.VolumeMount{
		Name:      "volume-one",
		ReadOnly:  true,
		MountPath: "/home/root/one",
	}

	mountTwo := corev1.VolumeMount{
		Name:      "volume-two",
		ReadOnly:  true,
		MountPath: "/home/root/two",
	}

	spec := coh.CoherenceResourceSpec{
		VolumeMounts: []corev1.VolumeMount{mountOne, mountTwo},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	stsExpected.Spec.Template.Spec.Containers[0].VolumeMounts = append(stsExpected.Spec.Template.Spec.Containers[0].VolumeMounts, mountOne, mountTwo)
	stsExpected.Spec.Template.Spec.InitContainers[0].VolumeMounts = append(stsExpected.Spec.Template.Spec.InitContainers[0].VolumeMounts, mountOne, mountTwo)
	stsExpected.Spec.Template.Spec.InitContainers[1].VolumeMounts = append(stsExpected.Spec.Template.Spec.InitContainers[1].VolumeMounts, mountOne, mountTwo)

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}
