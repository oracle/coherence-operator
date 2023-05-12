/*
 * Copyright (c) 2023, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1_test

import (
	coh "github.com/oracle/coherence-operator/api/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
	"testing"
)

func TestCreateJobWithPersistenceEmpty(t *testing.T) {
	spec := coh.CoherenceResourceSpec{
		Coherence: &coh.CoherenceSpec{
			Persistence: &coh.PersistenceSpec{},
		},
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	jobExpected := createMinimalExpectedJob(deployment)

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}

func TestCreateJobWithPersistenceModeOnDemand(t *testing.T) {
	spec := coh.CoherenceResourceSpec{
		Coherence: &coh.CoherenceSpec{
			Persistence: &coh.PersistenceSpec{
				Mode: pointer.String("on-demand"),
			},
		},
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	jobExpected := createMinimalExpectedJob(deployment)
	addEnvVarsToJob(jobExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: coh.EnvVarCohPersistenceMode, Value: "on-demand"})

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}

func TestCreateJobWithPersistenceModeActive(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		Coherence: &coh.CoherenceSpec{
			Persistence: &coh.PersistenceSpec{
				Mode: pointer.String("active"),
			},
		},
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	jobExpected := createMinimalExpectedJob(deployment)
	addEnvVarsToJob(jobExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: coh.EnvVarCohPersistenceMode, Value: "active"})

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}

func TestCreateJobWithPersistenceModeActiveAsync(t *testing.T) {
	spec := coh.CoherenceResourceSpec{
		Coherence: &coh.CoherenceSpec{
			Persistence: &coh.PersistenceSpec{
				Mode: pointer.String("active-async"),
			},
		},
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	jobExpected := createMinimalExpectedJob(deployment)
	addEnvVarsToJob(jobExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: coh.EnvVarCohPersistenceMode, Value: "active-async"})

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}

func TestCreateJobWithPersistenceVolume(t *testing.T) {
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

	jobExpected, deployment := createResourcesForJob(spec)

	// add the expected volume to the Pod
	jobExpected.Spec.Template.Spec.Volumes = append(jobExpected.Spec.Template.Spec.Volumes,
		corev1.Volume{Name: coh.VolumeNamePersistence, VolumeSource: vol})

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}

func TestCreateJobWithPersistencePVC(t *testing.T) {
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

	jobExpected, deployment := createResourcesForJob(spec)

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}

func TestCreateJobWithPersistenceVolumeAndPVC(t *testing.T) {
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

	jobExpected, deployment := createResourcesForJob(spec)

	// add the expected volume to the Pod - the PVC should not be set
	jobExpected.Spec.Template.Spec.Volumes = append(jobExpected.Spec.Template.Spec.Volumes,
		corev1.Volume{Name: coh.VolumeNamePersistence, VolumeSource: vol})

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}

func TestCreateJobWithPersistenceSnapshotVolume(t *testing.T) {
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
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	jobExpected := createMinimalExpectedJob(deployment)

	// Add the expected environment variables
	addEnvVarsToJob(jobExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: coh.EnvVarCohSnapshotDir, Value: coh.VolumeMountPathSnapshots})
	addEnvVarsToJob(jobExpected, coh.ContainerNameOperatorInit, corev1.EnvVar{Name: coh.EnvVarCohSnapshotDir, Value: coh.VolumeMountPathSnapshots})

	// add the expected volume mount to the Operator init-container
	jobExpected.Spec.Template.Spec.InitContainers[0].VolumeMounts = append(jobExpected.Spec.Template.Spec.InitContainers[0].VolumeMounts, corev1.VolumeMount{
		Name:      coh.VolumeNameSnapshots,
		MountPath: coh.VolumeMountPathSnapshots,
	})

	// add the expected volume mount to the coherence container
	jobExpected.Spec.Template.Spec.Containers[0].VolumeMounts = append(jobExpected.Spec.Template.Spec.Containers[0].VolumeMounts, corev1.VolumeMount{
		Name:      coh.VolumeNameSnapshots,
		MountPath: coh.VolumeMountPathSnapshots,
	})

	// add the expected volume to the Pod
	jobExpected.Spec.Template.Spec.Volumes = append(jobExpected.Spec.Template.Spec.Volumes,
		corev1.Volume{Name: coh.VolumeNameSnapshots, VolumeSource: vol})

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}

func TestCreateJobWithPersistenceSnapshotPVC(t *testing.T) {
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
	deployment := createTestCoherenceJob(spec)

	// Create expected Job
	jobExpected := createMinimalExpectedJob(deployment)

	// Add the expected environment variables
	addEnvVarsToJob(jobExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: coh.EnvVarCohSnapshotDir, Value: coh.VolumeMountPathSnapshots})
	addEnvVarsToJob(jobExpected, coh.ContainerNameOperatorInit, corev1.EnvVar{Name: coh.EnvVarCohSnapshotDir, Value: coh.VolumeMountPathSnapshots})

	// add the expected volume mount to the Operator init-container
	jobExpected.Spec.Template.Spec.InitContainers[0].VolumeMounts = append(jobExpected.Spec.Template.Spec.InitContainers[0].VolumeMounts, corev1.VolumeMount{
		Name:      coh.VolumeNameSnapshots,
		MountPath: coh.VolumeMountPathSnapshots,
	})

	// add the expected volume mount to the coherence container
	jobExpected.Spec.Template.Spec.Containers[0].VolumeMounts = append(jobExpected.Spec.Template.Spec.Containers[0].VolumeMounts, corev1.VolumeMount{
		Name:      coh.VolumeNameSnapshots,
		MountPath: coh.VolumeMountPathSnapshots,
	})

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}

func TestCreateJobWithPersistenceSnapshotVolumeAndPVC(t *testing.T) {
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
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	jobExpected := createMinimalExpectedJob(deployment)

	// Add the expected environment variables
	addEnvVarsToJob(jobExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: coh.EnvVarCohSnapshotDir, Value: coh.VolumeMountPathSnapshots})
	addEnvVarsToJob(jobExpected, coh.ContainerNameOperatorInit, corev1.EnvVar{Name: coh.EnvVarCohSnapshotDir, Value: coh.VolumeMountPathSnapshots})

	// add the expected volume mount to the Operator init-container
	jobExpected.Spec.Template.Spec.InitContainers[0].VolumeMounts = append(jobExpected.Spec.Template.Spec.InitContainers[0].VolumeMounts, corev1.VolumeMount{
		Name:      coh.VolumeNameSnapshots,
		MountPath: coh.VolumeMountPathSnapshots,
	})

	// add the expected volume mount to the coherence container
	jobExpected.Spec.Template.Spec.Containers[0].VolumeMounts = append(jobExpected.Spec.Template.Spec.Containers[0].VolumeMounts, corev1.VolumeMount{
		Name:      coh.VolumeNameSnapshots,
		MountPath: coh.VolumeMountPathSnapshots,
	})

	// add the expected volume to the Pod - the PVC should not be used
	jobExpected.Spec.Template.Spec.Volumes = append(jobExpected.Spec.Template.Spec.Volumes,
		corev1.Volume{Name: coh.VolumeNameSnapshots, VolumeSource: vol})

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}

func TestCreateJobWithPersistenceAndSnapshotVolume(t *testing.T) {
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
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	jobExpected := createMinimalExpectedJob(deployment)

	// Add the expected environment variables
	addEnvVarsToJob(jobExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: coh.EnvVarCohPersistenceDir, Value: coh.VolumeMountPathPersistence})
	addEnvVarsToJob(jobExpected, coh.ContainerNameOperatorInit, corev1.EnvVar{Name: coh.EnvVarCohPersistenceDir, Value: coh.VolumeMountPathPersistence})
	addEnvVarsToJob(jobExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: coh.EnvVarCohSnapshotDir, Value: coh.VolumeMountPathSnapshots})
	addEnvVarsToJob(jobExpected, coh.ContainerNameOperatorInit, corev1.EnvVar{Name: coh.EnvVarCohSnapshotDir, Value: coh.VolumeMountPathSnapshots})

	// add the expected volume mount to the Operator init-container
	jobExpected.Spec.Template.Spec.InitContainers[0].VolumeMounts = append(jobExpected.Spec.Template.Spec.InitContainers[0].VolumeMounts,
		corev1.VolumeMount{Name: coh.VolumeNamePersistence, MountPath: coh.VolumeMountPathPersistence},
		corev1.VolumeMount{Name: coh.VolumeNameSnapshots, MountPath: coh.VolumeMountPathSnapshots})

	// add the expected volume mount to the coherence container
	jobExpected.Spec.Template.Spec.Containers[0].VolumeMounts = append(jobExpected.Spec.Template.Spec.Containers[0].VolumeMounts,
		corev1.VolumeMount{Name: coh.VolumeNamePersistence, MountPath: coh.VolumeMountPathPersistence},
		corev1.VolumeMount{Name: coh.VolumeNameSnapshots, MountPath: coh.VolumeMountPathSnapshots})

	// add the expected volume to the Pod
	jobExpected.Spec.Template.Spec.Volumes = append(jobExpected.Spec.Template.Spec.Volumes,
		corev1.Volume{Name: coh.VolumeNamePersistence, VolumeSource: volOne},
		corev1.Volume{Name: coh.VolumeNameSnapshots, VolumeSource: volTwo})

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}

func TestCreateJobWithPersistenceAndSnapshotPVC(t *testing.T) {
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
	deployment := createTestCoherenceJob(spec)

	// Create expected Job
	jobExpected := createMinimalExpectedJob(deployment)

	// Add the expected environment variables
	addEnvVarsToJob(jobExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: coh.EnvVarCohPersistenceDir, Value: coh.VolumeMountPathPersistence})
	addEnvVarsToJob(jobExpected, coh.ContainerNameOperatorInit, corev1.EnvVar{Name: coh.EnvVarCohPersistenceDir, Value: coh.VolumeMountPathPersistence})
	addEnvVarsToJob(jobExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: coh.EnvVarCohSnapshotDir, Value: coh.VolumeMountPathSnapshots})
	addEnvVarsToJob(jobExpected, coh.ContainerNameOperatorInit, corev1.EnvVar{Name: coh.EnvVarCohSnapshotDir, Value: coh.VolumeMountPathSnapshots})

	// add the expected volume mount to the Operator init-container
	jobExpected.Spec.Template.Spec.InitContainers[0].VolumeMounts = append(jobExpected.Spec.Template.Spec.InitContainers[0].VolumeMounts,
		corev1.VolumeMount{Name: coh.VolumeNamePersistence, MountPath: coh.VolumeMountPathPersistence},
		corev1.VolumeMount{Name: coh.VolumeNameSnapshots, MountPath: coh.VolumeMountPathSnapshots})

	// add the expected volume mount to the coherence container
	jobExpected.Spec.Template.Spec.Containers[0].VolumeMounts = append(jobExpected.Spec.Template.Spec.Containers[0].VolumeMounts,
		corev1.VolumeMount{Name: coh.VolumeNamePersistence, MountPath: coh.VolumeMountPathPersistence},
		corev1.VolumeMount{Name: coh.VolumeNameSnapshots, MountPath: coh.VolumeMountPathSnapshots})

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}

func createResourcesForJob(spec coh.CoherenceResourceSpec) (*batchv1.Job, *coh.CoherenceJob) {
	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	jobExpected := createMinimalExpectedJob(deployment)

	// Add the expected environment variables
	addEnvVarsToJob(jobExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: coh.EnvVarCohPersistenceDir, Value: coh.VolumeMountPathPersistence})
	addEnvVarsToJob(jobExpected, coh.ContainerNameOperatorInit, corev1.EnvVar{Name: coh.EnvVarCohPersistenceDir, Value: coh.VolumeMountPathPersistence})

	// add the expected volume mount to the Operator init-container
	jobExpected.Spec.Template.Spec.InitContainers[0].VolumeMounts = append(jobExpected.Spec.Template.Spec.InitContainers[0].VolumeMounts, corev1.VolumeMount{
		Name:      coh.VolumeNamePersistence,
		MountPath: coh.VolumeMountPathPersistence,
	})

	// add the expected volume mount to the coherence container
	jobExpected.Spec.Template.Spec.Containers[0].VolumeMounts = append(jobExpected.Spec.Template.Spec.Containers[0].VolumeMounts, corev1.VolumeMount{
		Name:      coh.VolumeNamePersistence,
		MountPath: coh.VolumeMountPathPersistence,
	})

	return jobExpected, deployment
}
