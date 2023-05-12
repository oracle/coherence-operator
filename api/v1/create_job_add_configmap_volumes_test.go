/*
 * Copyright (c) 2023, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1_test

import (
	coh "github.com/oracle/coherence-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	"testing"
)

func TestCreateJobWithConfigMapVolumesEmpty(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		ConfigMapVolumes: []coh.ConfigMapVolumeSpec{},
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	jobExpected := createMinimalExpectedJob(deployment)

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}

func TestCreateJobWithConfigMapVolume(t *testing.T) {

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
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	jobExpected := createMinimalExpectedJob(deployment)
	jobExpected.Spec.Template.Spec.Volumes = append(jobExpected.Spec.Template.Spec.Volumes, vol)
	for i, c := range jobExpected.Spec.Template.Spec.InitContainers {
		c.VolumeMounts = append(c.VolumeMounts, vm)
		jobExpected.Spec.Template.Spec.InitContainers[i] = c
	}
	for i, c := range jobExpected.Spec.Template.Spec.Containers {
		c.VolumeMounts = append(c.VolumeMounts, vm)
		jobExpected.Spec.Template.Spec.Containers[i] = c
	}

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}
