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

func TestCreateStatefulSetWithConfigMapVolumesEmpty(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		ConfigMapVolumes: []coh.ConfigMapVolumeSpec{},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithConfigMapVolume(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		ConfigMapVolumes: []coh.ConfigMapVolumeSpec{
			{
				Name:      "test-config",
				MountPath: "/apps/config",
			},
		},
	}

	vm := corev1.VolumeMount{
		Name:      "test-config",
		MountPath: "/apps/config",
	}

	vol := corev1.Volume{
		Name: "test-config",
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: "test-config",
				},
			},
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	stsExpected.Spec.Template.Spec.Volumes = append(stsExpected.Spec.Template.Spec.Volumes, vol)
	for i, c := range stsExpected.Spec.Template.Spec.InitContainers {
		c.VolumeMounts = append(c.VolumeMounts, vm)
		stsExpected.Spec.Template.Spec.InitContainers[i] = c
	}
	for i, c := range stsExpected.Spec.Template.Spec.Containers {
		c.VolumeMounts = append(c.VolumeMounts, vm)
		stsExpected.Spec.Template.Spec.Containers[i] = c
	}

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}
