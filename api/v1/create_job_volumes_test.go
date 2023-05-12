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

func TestCreateJobWithEmptyVolumes(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		Volumes: []corev1.Volume{},
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	stsExpected := createMinimalExpectedJob(deployment)

	// assert that the Job is as expected
	assertJobCreation(t, deployment, stsExpected)
}

func TestCreateJobWithOneVolume(t *testing.T) {

	volumeOne := corev1.Volume{
		Name: "volume-one",
		VolumeSource: corev1.VolumeSource{
			HostPath: &corev1.HostPathVolumeSource{
				Path: "/volumes/one",
			},
		},
	}

	spec := coh.CoherenceResourceSpec{
		Volumes: []corev1.Volume{volumeOne},
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	stsExpected := createMinimalExpectedJob(deployment)
	stsExpected.Spec.Template.Spec.Volumes = append(stsExpected.Spec.Template.Spec.Volumes, volumeOne)

	// assert that the Job is as expected
	assertJobCreation(t, deployment, stsExpected)
}

func TestCreateJobWithTwoVolumes(t *testing.T) {

	volumeOne := corev1.Volume{
		Name: "volume-one",
		VolumeSource: corev1.VolumeSource{
			HostPath: &corev1.HostPathVolumeSource{
				Path: "/volumes/one",
			},
		},
	}

	volumeTwo := corev1.Volume{
		Name: "volume-two",
		VolumeSource: corev1.VolumeSource{
			HostPath: &corev1.HostPathVolumeSource{
				Path: "/volumes/two",
			},
		},
	}

	spec := coh.CoherenceResourceSpec{
		Volumes: []corev1.Volume{volumeOne, volumeTwo},
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	stsExpected := createMinimalExpectedJob(deployment)
	stsExpected.Spec.Template.Spec.Volumes = append(stsExpected.Spec.Template.Spec.Volumes, volumeOne, volumeTwo)

	// assert that the Job is as expected
	assertJobCreation(t, deployment, stsExpected)
}
