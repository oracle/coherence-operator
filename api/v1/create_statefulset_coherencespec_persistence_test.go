/*
 * Copyright (c) 2020, 2022, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1_test

import (
	coh "github.com/oracle/coherence-operator/api/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
	"testing"
)

func TestCreateStatefulSetWithPersistenceEmpty(t *testing.T) {
	spec := coh.CoherenceResourceSpec{
		Coherence: &coh.CoherenceSpec{
			Persistence: &coh.PersistenceSpec{},
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithPersistenceModeOnDemand(t *testing.T) {
	spec := coh.CoherenceResourceSpec{
		Coherence: &coh.CoherenceSpec{
			Persistence: &coh.PersistenceSpec{
				Mode: pointer.StringPtr("on-demand"),
			},
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	addEnvVars(stsExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: coh.EnvVarCohPersistenceMode, Value: "on-demand"})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithPersistenceModeActive(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		Coherence: &coh.CoherenceSpec{
			Persistence: &coh.PersistenceSpec{
				Mode: pointer.StringPtr("active"),
			},
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	addEnvVars(stsExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: coh.EnvVarCohPersistenceMode, Value: "active"})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithPersistenceModeActiveAsync(t *testing.T) {
	spec := coh.CoherenceResourceSpec{
		Coherence: &coh.CoherenceSpec{
			Persistence: &coh.PersistenceSpec{
				Mode: pointer.StringPtr("active-async"),
			},
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	addEnvVars(stsExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: coh.EnvVarCohPersistenceMode, Value: "active-async"})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithPersistenceVolume(t *testing.T) {
	vol := corev1.VolumeSource{
		HostPath: &corev1.HostPathVolumeSource{Path: "/data/persistence"},
	}

	spec := coh.CoherenceResourceSpec{
		Coherence: &coh.CoherenceSpec{
			Persistence: &coh.PersistenceSpec{
				PersistentStorageSpec: coh.PersistentStorageSpec{
					Volume: &vol,
				},
			},
		},
	}

	stsExpected, deployment := createResources(spec)

	// add the expected volume to the Pod
	stsExpected.Spec.Template.Spec.Volumes = append(stsExpected.Spec.Template.Spec.Volumes,
		corev1.Volume{Name: coh.VolumeNamePersistence, VolumeSource: vol})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithPersistencePVC(t *testing.T) {
	mode := corev1.PersistentVolumeFilesystem
	pvc := corev1.PersistentVolumeClaimSpec{
		AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
		Selector: &metav1.LabelSelector{
			MatchLabels: map[string]string{"Foo": "Bar"},
		},
		VolumeName:       "test-volume",
		StorageClassName: stringPtr("Foo"),
		VolumeMode:       &mode,
	}

	spec := coh.CoherenceResourceSpec{
		Coherence: &coh.CoherenceSpec{
			Persistence: &coh.PersistenceSpec{
				PersistentStorageSpec: coh.PersistentStorageSpec{
					PersistentVolumeClaim: &pvc,
				},
			},
		},
	}

	stsExpected, deployment := createResources(spec)

	// add the expected PVC
	labels := deployment.CreateCommonLabels()
	labels[coh.LabelComponent] = coh.LabelComponentPVC
	stsExpected.Spec.VolumeClaimTemplates = append(stsExpected.Spec.VolumeClaimTemplates, corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:   coh.VolumeNamePersistence,
			Labels: labels,
		},
		Spec: pvc,
	})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithPersistenceVolumeAndPVC(t *testing.T) {
	mode := corev1.PersistentVolumeFilesystem
	vol := corev1.VolumeSource{
		HostPath: &corev1.HostPathVolumeSource{Path: "/data/persistence"},
	}
	pvc := corev1.PersistentVolumeClaimSpec{
		AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
		Selector: &metav1.LabelSelector{
			MatchLabels: map[string]string{"Foo": "Bar"},
		},
		VolumeName:       "test-volume",
		StorageClassName: stringPtr("Foo"),
		VolumeMode:       &mode,
	}

	spec := coh.CoherenceResourceSpec{
		Coherence: &coh.CoherenceSpec{
			Persistence: &coh.PersistenceSpec{
				PersistentStorageSpec: coh.PersistentStorageSpec{
					Volume:                &vol,
					PersistentVolumeClaim: &pvc,
				},
			},
		},
	}

	stsExpected, deployment := createResources(spec)

	// add the expected volume to the Pod - the PVC should not be set
	stsExpected.Spec.Template.Spec.Volumes = append(stsExpected.Spec.Template.Spec.Volumes,
		corev1.Volume{Name: coh.VolumeNamePersistence, VolumeSource: vol})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithPersistenceSnapshotVolume(t *testing.T) {
	vol := corev1.VolumeSource{
		HostPath: &corev1.HostPathVolumeSource{Path: "/data/snapshots"},
	}

	spec := coh.CoherenceResourceSpec{
		Coherence: &coh.CoherenceSpec{
			Persistence: &coh.PersistenceSpec{
				Snapshots: &coh.PersistentStorageSpec{
					Volume: &vol,
				},
			},
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)

	// Add the expected environment variables
	addEnvVars(stsExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: coh.EnvVarCohSnapshotDir, Value: coh.VolumeMountPathSnapshots})
	addEnvVars(stsExpected, coh.ContainerNameOperatorInit, corev1.EnvVar{Name: coh.EnvVarCohSnapshotDir, Value: coh.VolumeMountPathSnapshots})

	// add the expected volume mount to the Operator init-container
	stsExpected.Spec.Template.Spec.InitContainers[0].VolumeMounts = append(stsExpected.Spec.Template.Spec.InitContainers[0].VolumeMounts, corev1.VolumeMount{
		Name:      coh.VolumeNameSnapshots,
		MountPath: coh.VolumeMountPathSnapshots,
	})

	// add the expected volume mount to the coherence container
	stsExpected.Spec.Template.Spec.Containers[0].VolumeMounts = append(stsExpected.Spec.Template.Spec.Containers[0].VolumeMounts, corev1.VolumeMount{
		Name:      coh.VolumeNameSnapshots,
		MountPath: coh.VolumeMountPathSnapshots,
	})

	// add the expected volume to the Pod
	stsExpected.Spec.Template.Spec.Volumes = append(stsExpected.Spec.Template.Spec.Volumes,
		corev1.Volume{Name: coh.VolumeNameSnapshots, VolumeSource: vol})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithPersistenceSnapshotPVC(t *testing.T) {
	mode := corev1.PersistentVolumeFilesystem
	pvc := corev1.PersistentVolumeClaimSpec{
		AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
		Selector: &metav1.LabelSelector{
			MatchLabels: map[string]string{"Foo": "Bar"},
		},
		VolumeName:       "test-volume",
		StorageClassName: stringPtr("Foo"),
		VolumeMode:       &mode,
	}

	spec := coh.CoherenceResourceSpec{
		Coherence: &coh.CoherenceSpec{
			Persistence: &coh.PersistenceSpec{
				Snapshots: &coh.PersistentStorageSpec{
					PersistentVolumeClaim: &pvc,
				},
			},
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)

	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)

	// Add the expected environment variables
	addEnvVars(stsExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: coh.EnvVarCohSnapshotDir, Value: coh.VolumeMountPathSnapshots})
	addEnvVars(stsExpected, coh.ContainerNameOperatorInit, corev1.EnvVar{Name: coh.EnvVarCohSnapshotDir, Value: coh.VolumeMountPathSnapshots})

	// add the expected volume mount to the Operator init-container
	stsExpected.Spec.Template.Spec.InitContainers[0].VolumeMounts = append(stsExpected.Spec.Template.Spec.InitContainers[0].VolumeMounts, corev1.VolumeMount{
		Name:      coh.VolumeNameSnapshots,
		MountPath: coh.VolumeMountPathSnapshots,
	})

	// add the expected volume mount to the coherence container
	stsExpected.Spec.Template.Spec.Containers[0].VolumeMounts = append(stsExpected.Spec.Template.Spec.Containers[0].VolumeMounts, corev1.VolumeMount{
		Name:      coh.VolumeNameSnapshots,
		MountPath: coh.VolumeMountPathSnapshots,
	})

	// add the expected snapshot PVC
	labels := deployment.CreateCommonLabels()
	labels[coh.LabelComponent] = coh.LabelComponentPVC
	stsExpected.Spec.VolumeClaimTemplates = append(stsExpected.Spec.VolumeClaimTemplates, corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:   coh.VolumeNameSnapshots,
			Labels: labels,
		},
		Spec: pvc,
	})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithPersistenceSnapshotVolumeAndPVC(t *testing.T) {
	vol := corev1.VolumeSource{
		HostPath: &corev1.HostPathVolumeSource{Path: "/data/snapshots"},
	}
	mode := corev1.PersistentVolumeFilesystem
	pvc := corev1.PersistentVolumeClaimSpec{
		AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
		Selector: &metav1.LabelSelector{
			MatchLabels: map[string]string{"Foo": "Bar"},
		},
		VolumeName:       "test-volume",
		StorageClassName: stringPtr("Foo"),
		VolumeMode:       &mode,
	}

	spec := coh.CoherenceResourceSpec{
		Coherence: &coh.CoherenceSpec{
			Persistence: &coh.PersistenceSpec{
				Snapshots: &coh.PersistentStorageSpec{
					Volume:                &vol,
					PersistentVolumeClaim: &pvc,
				},
			},
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)

	// Add the expected environment variables
	addEnvVars(stsExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: coh.EnvVarCohSnapshotDir, Value: coh.VolumeMountPathSnapshots})
	addEnvVars(stsExpected, coh.ContainerNameOperatorInit, corev1.EnvVar{Name: coh.EnvVarCohSnapshotDir, Value: coh.VolumeMountPathSnapshots})

	// add the expected volume mount to the Operator init-container
	stsExpected.Spec.Template.Spec.InitContainers[0].VolumeMounts = append(stsExpected.Spec.Template.Spec.InitContainers[0].VolumeMounts, corev1.VolumeMount{
		Name:      coh.VolumeNameSnapshots,
		MountPath: coh.VolumeMountPathSnapshots,
	})

	// add the expected volume mount to the coherence container
	stsExpected.Spec.Template.Spec.Containers[0].VolumeMounts = append(stsExpected.Spec.Template.Spec.Containers[0].VolumeMounts, corev1.VolumeMount{
		Name:      coh.VolumeNameSnapshots,
		MountPath: coh.VolumeMountPathSnapshots,
	})

	// add the expected volume to the Pod - the PVC should not be used
	stsExpected.Spec.Template.Spec.Volumes = append(stsExpected.Spec.Template.Spec.Volumes,
		corev1.Volume{Name: coh.VolumeNameSnapshots, VolumeSource: vol})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithPersistenceAndSnapshotVolume(t *testing.T) {
	volOne := corev1.VolumeSource{
		HostPath: &corev1.HostPathVolumeSource{Path: "/data/active"},
	}
	volTwo := corev1.VolumeSource{
		HostPath: &corev1.HostPathVolumeSource{Path: "/data/snapshots"},
	}

	spec := coh.CoherenceResourceSpec{
		Coherence: &coh.CoherenceSpec{
			Persistence: &coh.PersistenceSpec{
				PersistentStorageSpec: coh.PersistentStorageSpec{
					Volume: &volOne,
				},
				Snapshots: &coh.PersistentStorageSpec{
					Volume: &volTwo,
				},
			},
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)

	// Add the expected environment variables
	addEnvVars(stsExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: coh.EnvVarCohPersistenceDir, Value: coh.VolumeMountPathPersistence})
	addEnvVars(stsExpected, coh.ContainerNameOperatorInit, corev1.EnvVar{Name: coh.EnvVarCohPersistenceDir, Value: coh.VolumeMountPathPersistence})
	addEnvVars(stsExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: coh.EnvVarCohSnapshotDir, Value: coh.VolumeMountPathSnapshots})
	addEnvVars(stsExpected, coh.ContainerNameOperatorInit, corev1.EnvVar{Name: coh.EnvVarCohSnapshotDir, Value: coh.VolumeMountPathSnapshots})

	// add the expected volume mount to the Operator init-container
	stsExpected.Spec.Template.Spec.InitContainers[0].VolumeMounts = append(stsExpected.Spec.Template.Spec.InitContainers[0].VolumeMounts,
		corev1.VolumeMount{Name: coh.VolumeNamePersistence, MountPath: coh.VolumeMountPathPersistence},
		corev1.VolumeMount{Name: coh.VolumeNameSnapshots, MountPath: coh.VolumeMountPathSnapshots})

	// add the expected volume mount to the coherence container
	stsExpected.Spec.Template.Spec.Containers[0].VolumeMounts = append(stsExpected.Spec.Template.Spec.Containers[0].VolumeMounts,
		corev1.VolumeMount{Name: coh.VolumeNamePersistence, MountPath: coh.VolumeMountPathPersistence},
		corev1.VolumeMount{Name: coh.VolumeNameSnapshots, MountPath: coh.VolumeMountPathSnapshots})

	// add the expected volume to the Pod
	stsExpected.Spec.Template.Spec.Volumes = append(stsExpected.Spec.Template.Spec.Volumes,
		corev1.Volume{Name: coh.VolumeNamePersistence, VolumeSource: volOne},
		corev1.Volume{Name: coh.VolumeNameSnapshots, VolumeSource: volTwo})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithPersistenceAndSnapshotPVC(t *testing.T) {
	mode := corev1.PersistentVolumeFilesystem
	pvcOne := corev1.PersistentVolumeClaimSpec{
		AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
		Selector: &metav1.LabelSelector{
			MatchLabels: map[string]string{"Foo": "Bar"},
		},
		VolumeName:       "test-volume-one",
		StorageClassName: stringPtr("Foo"),
		VolumeMode:       &mode,
	}
	pvcTwo := corev1.PersistentVolumeClaimSpec{
		AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
		Selector: &metav1.LabelSelector{
			MatchLabels: map[string]string{"Foo": "Bar"},
		},
		VolumeName:       "test-volume-two",
		StorageClassName: stringPtr("Foo"),
		VolumeMode:       &mode,
	}

	spec := coh.CoherenceResourceSpec{
		Coherence: &coh.CoherenceSpec{
			Persistence: &coh.PersistenceSpec{
				PersistentStorageSpec: coh.PersistentStorageSpec{
					PersistentVolumeClaim: &pvcOne,
				},
				Snapshots: &coh.PersistentStorageSpec{
					PersistentVolumeClaim: &pvcTwo,
				},
			},
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)

	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)

	// Add the expected environment variables
	addEnvVars(stsExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: coh.EnvVarCohPersistenceDir, Value: coh.VolumeMountPathPersistence})
	addEnvVars(stsExpected, coh.ContainerNameOperatorInit, corev1.EnvVar{Name: coh.EnvVarCohPersistenceDir, Value: coh.VolumeMountPathPersistence})
	addEnvVars(stsExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: coh.EnvVarCohSnapshotDir, Value: coh.VolumeMountPathSnapshots})
	addEnvVars(stsExpected, coh.ContainerNameOperatorInit, corev1.EnvVar{Name: coh.EnvVarCohSnapshotDir, Value: coh.VolumeMountPathSnapshots})

	// add the expected volume mount to the Operator init-container
	stsExpected.Spec.Template.Spec.InitContainers[0].VolumeMounts = append(stsExpected.Spec.Template.Spec.InitContainers[0].VolumeMounts,
		corev1.VolumeMount{Name: coh.VolumeNamePersistence, MountPath: coh.VolumeMountPathPersistence},
		corev1.VolumeMount{Name: coh.VolumeNameSnapshots, MountPath: coh.VolumeMountPathSnapshots})

	// add the expected volume mount to the coherence container
	stsExpected.Spec.Template.Spec.Containers[0].VolumeMounts = append(stsExpected.Spec.Template.Spec.Containers[0].VolumeMounts,
		corev1.VolumeMount{Name: coh.VolumeNamePersistence, MountPath: coh.VolumeMountPathPersistence},
		corev1.VolumeMount{Name: coh.VolumeNameSnapshots, MountPath: coh.VolumeMountPathSnapshots})

	// add the expected PVCs
	labels := deployment.CreateCommonLabels()
	labels[coh.LabelComponent] = coh.LabelComponentPVC
	stsExpected.Spec.VolumeClaimTemplates = append(stsExpected.Spec.VolumeClaimTemplates, corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:   coh.VolumeNamePersistence,
			Labels: labels,
		},
		Spec: pvcOne,
	})

	stsExpected.Spec.VolumeClaimTemplates = append(stsExpected.Spec.VolumeClaimTemplates, corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:   coh.VolumeNameSnapshots,
			Labels: labels,
		},
		Spec: pvcTwo,
	})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func createResources(spec coh.CoherenceResourceSpec) (*appsv1.StatefulSet, *coh.Coherence) {
	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)

	// Add the expected environment variables
	addEnvVars(stsExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: coh.EnvVarCohPersistenceDir, Value: coh.VolumeMountPathPersistence})
	addEnvVars(stsExpected, coh.ContainerNameOperatorInit, corev1.EnvVar{Name: coh.EnvVarCohPersistenceDir, Value: coh.VolumeMountPathPersistence})

	// add the expected volume mount to the Operator init-container
	stsExpected.Spec.Template.Spec.InitContainers[0].VolumeMounts = append(stsExpected.Spec.Template.Spec.InitContainers[0].VolumeMounts, corev1.VolumeMount{
		Name:      coh.VolumeNamePersistence,
		MountPath: coh.VolumeMountPathPersistence,
	})

	// add the expected volume mount to the coherence container
	stsExpected.Spec.Template.Spec.Containers[0].VolumeMounts = append(stsExpected.Spec.Template.Spec.Containers[0].VolumeMounts, corev1.VolumeMount{
		Name:      coh.VolumeNamePersistence,
		MountPath: coh.VolumeMountPathPersistence,
	})

	return stsExpected, deployment
}
