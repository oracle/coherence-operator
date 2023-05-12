/*
 * Copyright (c) 2020, 2023, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1_test

import (
	coh "github.com/oracle/coherence-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	"testing"
)

func TestCreateStatefulSetWithEmptyVolumeClaimTemplates(t *testing.T) {

	spec := coh.CoherenceStatefulSetResourceSpec{
		VolumeClaimTemplates: []coh.PersistentVolumeClaim{},
	}

	// Create the test deployment
	deployment := createTestCoherenceDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithOneVolumeClaimTemplate(t *testing.T) {

	volumeOne := coh.PersistentVolumeClaim{
		Metadata: coh.PersistentVolumeClaimObjectMeta{
			Name: "PVCOne",
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			VolumeName:       "PVOne",
			StorageClassName: stringPtr("TestStorage"),
		},
	}

	spec := coh.CoherenceStatefulSetResourceSpec{
		VolumeClaimTemplates: []coh.PersistentVolumeClaim{volumeOne},
	}

	// Create the test deployment
	deployment := createTestCoherenceDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	stsExpected.Spec.VolumeClaimTemplates = append(stsExpected.Spec.VolumeClaimTemplates, volumeOne.ToPVC())

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithTwoVolumeClaimTemplates(t *testing.T) {

	volumeOne := coh.PersistentVolumeClaim{
		Metadata: coh.PersistentVolumeClaimObjectMeta{
			Name: "PVCOne",
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			VolumeName:       "PVOne",
			StorageClassName: stringPtr("TestStorage"),
		},
	}

	volumeTwo := coh.PersistentVolumeClaim{
		Metadata: coh.PersistentVolumeClaimObjectMeta{
			Name: "PVCTwo",
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			VolumeName:       "PVTwo",
			StorageClassName: stringPtr("TestStorage"),
		},
	}

	spec := coh.CoherenceStatefulSetResourceSpec{
		VolumeClaimTemplates: []coh.PersistentVolumeClaim{volumeOne, volumeTwo},
	}

	// Create the test deployment
	deployment := createTestCoherenceDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	stsExpected.Spec.VolumeClaimTemplates = append(stsExpected.Spec.VolumeClaimTemplates, volumeOne.ToPVC(), volumeTwo.ToPVC())

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}
