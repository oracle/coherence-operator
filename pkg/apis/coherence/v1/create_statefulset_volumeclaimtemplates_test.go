/*
 * Copyright (c) 2020, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1_test

import (
	coh "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestCreateStatefulSetWithEmptyVolumeClaimTemplates(t *testing.T) {

	spec := coh.CoherenceDeploymentSpec{
		VolumeClaimTemplates: []corev1.PersistentVolumeClaim{},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithOneVolumeClaimTemplate(t *testing.T) {

	volumeOne := corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name: "PVCOne",
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			VolumeName:       "PVOne",
			StorageClassName: stringPtr("TestStorage"),
		},
	}

	spec := coh.CoherenceDeploymentSpec{
		VolumeClaimTemplates: []corev1.PersistentVolumeClaim{volumeOne},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	stsExpected.Spec.VolumeClaimTemplates = append(stsExpected.Spec.VolumeClaimTemplates, volumeOne)

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithTwoVolumeClaimTemplates(t *testing.T) {

	volumeOne := corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name: "PVCOne",
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			VolumeName:       "PVOne",
			StorageClassName: stringPtr("TestStorage"),
		},
	}

	volumeTwo := corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name: "PVCTwo",
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			VolumeName:       "PVTwo",
			StorageClassName: stringPtr("TestStorage"),
		},
	}

	spec := coh.CoherenceDeploymentSpec{
		VolumeClaimTemplates: []corev1.PersistentVolumeClaim{volumeOne, volumeTwo},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	stsExpected.Spec.VolumeClaimTemplates = append(stsExpected.Spec.VolumeClaimTemplates, volumeOne, volumeTwo)

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}
