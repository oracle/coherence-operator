/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	coherence "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Testing PersistentStorageSpec struct", func() {

	Context("Copying a PersistentStorageSpec using DeepCopyWithDefaults", func() {
		var (
			original *coherence.PersistentStorageSpec
			defaults *coherence.PersistentStorageSpec
			clone    *coherence.PersistentStorageSpec

			block      = corev1.PersistentVolumeBlock
			filesystem = corev1.PersistentVolumeFilesystem
			pvcOne     = &corev1.PersistentVolumeClaimSpec{
				AccessModes: []corev1.PersistentVolumeAccessMode{"ReadWriteOnce"},
				Resources: corev1.ResourceRequirements{
					Requests: map[corev1.ResourceName]resource.Quantity{"storage": resource.MustParse("2Gi")},
				},
				StorageClassName: stringPtr("sc1"),
				DataSource:       &corev1.TypedLocalObjectReference{Name: "pvc1", Kind: "PersistentVolumeClaim"},
				VolumeMode:       &block,
				VolumeName:       "name1",
				Selector:         &metav1.LabelSelector{MatchLabels: map[string]string{"component": "coh1"}},
			}
			pvcTwo = &corev1.PersistentVolumeClaimSpec{
				AccessModes: []corev1.PersistentVolumeAccessMode{"ReadWriteOnce"},
				Resources: corev1.ResourceRequirements{
					Requests: map[corev1.ResourceName]resource.Quantity{"storage": resource.MustParse("3Gi")},
				},
				StorageClassName: stringPtr("sc2"),
				DataSource:       &corev1.TypedLocalObjectReference{Name: "pvc2", Kind: "PersistentVolumeClaim"},
				VolumeMode:       &filesystem,
				VolumeName:       "name2",
				Selector:         &metav1.LabelSelector{MatchLabels: map[string]string{"component": "coh2"}},
			}

			volumeOne = &corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}}
			volumeTwo = &corev1.VolumeSource{NFS: &corev1.NFSVolumeSource{Server: "10.100.100.200", Path: "/"}}
		)

		JustBeforeEach(func() {
			clone = original.DeepCopyWithDefaults(defaults)
		})

		When("original and defaults are nil", func() {
			BeforeEach(func() {
				original = nil
				defaults = nil
			})

			It("the copy should be nil", func() {
				Expect(clone).Should(BeNil())
			})
		})

		When("defaults is nil", func() {
			BeforeEach(func() {
				original = &coherence.PersistentStorageSpec{
					Enabled:               boolPtr(true),
					PersistentVolumeClaim: pvcOne,
					Volume:                volumeOne,
				}

				defaults = nil
			})

			It("should copy the original Enabled", func() {
				Expect(*clone.Enabled).To(Equal(*original.Enabled))
			})

			It("should copy the original PersistentVolumeClaim", func() {
				Expect(*clone.PersistentVolumeClaim).To(Equal(*original.PersistentVolumeClaim))
			})

			It("should copy the original Volume", func() {
				Expect(*clone.Volume).To(Equal(*original.Volume))
			})
		})

		When("original is nil", func() {
			BeforeEach(func() {
				defaults = &coherence.PersistentStorageSpec{
					Enabled:               boolPtr(true),
					PersistentVolumeClaim: pvcOne,
					Volume:                volumeOne,
				}

				original = nil
			})

			It("should copy the defaults Enabled", func() {
				Expect(*clone.Enabled).To(Equal(*defaults.Enabled))
			})

			It("should copy the defaults PersistentVolumeClaim", func() {
				Expect(*clone.PersistentVolumeClaim).To(Equal(*defaults.PersistentVolumeClaim))
			})

			It("should copy the defaults Volume", func() {
				Expect(*clone.Volume).To(Equal(*defaults.Volume))
			})
		})

		When("all original fields are set", func() {
			BeforeEach(func() {
				original = &coherence.PersistentStorageSpec{
					Enabled:               boolPtr(true),
					PersistentVolumeClaim: pvcOne,
					Volume:                volumeOne,
				}

				defaults = &coherence.PersistentStorageSpec{
					Enabled:               boolPtr(false),
					PersistentVolumeClaim: pvcTwo,
					Volume:                volumeTwo,
				}
			})

			It("should copy the original Enabled", func() {
				Expect(*clone.Enabled).To(Equal(*original.Enabled))
			})

			It("should copy the original PersistentVolumeClaim", func() {
				Expect(*clone.PersistentVolumeClaim).To(Equal(*original.PersistentVolumeClaim))
			})

			It("should copy the original Volume", func() {
				Expect(*clone.Volume).To(Equal(*original.Volume))
			})
		})

		When("original Enabled is nil", func() {
			BeforeEach(func() {
				original = &coherence.PersistentStorageSpec{
					Enabled:               nil,
					PersistentVolumeClaim: pvcOne,
					Volume:                volumeOne,
				}

				defaults = &coherence.PersistentStorageSpec{
					Enabled:               boolPtr(true),
					PersistentVolumeClaim: pvcTwo,
					Volume:                volumeTwo,
				}
			})

			It("should copy the defaults Enabled", func() {
				Expect(*clone.Enabled).To(Equal(*defaults.Enabled))
			})

			It("should copy the original PersistentVolumeClaim", func() {
				Expect(*clone.PersistentVolumeClaim).To(Equal(*original.PersistentVolumeClaim))
			})

			It("should copy the original Volume", func() {
				Expect(*clone.Volume).To(Equal(*original.Volume))
			})
		})

		When("original PersistentVolumeClaim is nil", func() {
			BeforeEach(func() {
				original = &coherence.PersistentStorageSpec{
					Enabled:               boolPtr(true),
					PersistentVolumeClaim: nil,
					Volume:                volumeOne,
				}

				defaults = &coherence.PersistentStorageSpec{
					Enabled:               boolPtr(false),
					PersistentVolumeClaim: pvcTwo,
					Volume:                volumeTwo,
				}
			})

			It("should copy the original Enabled", func() {
				Expect(*clone.Enabled).To(Equal(*original.Enabled))
			})

			It("should copy the defaults PersistentVolumeClaim", func() {
				Expect(*clone.PersistentVolumeClaim).To(Equal(*defaults.PersistentVolumeClaim))
			})

			It("should copy the original Volume", func() {
				Expect(*clone.Volume).To(Equal(*original.Volume))
			})
		})

		When("original Volume is nil", func() {
			BeforeEach(func() {
				original = &coherence.PersistentStorageSpec{
					Enabled:               boolPtr(true),
					PersistentVolumeClaim: pvcOne,
					Volume:                nil,
				}

				defaults = &coherence.PersistentStorageSpec{
					Enabled:               boolPtr(false),
					PersistentVolumeClaim: pvcTwo,
					Volume:                volumeTwo,
				}
			})

			It("should copy the original Enabled", func() {
				Expect(*clone.Enabled).To(Equal(*original.Enabled))
			})

			It("should copy the original PersistentVolumeClaim", func() {
				Expect(*clone.PersistentVolumeClaim).To(Equal(*original.PersistentVolumeClaim))
			})

			It("should copy the defaults Volume", func() {
				Expect(*clone.Volume).To(Equal(*defaults.Volume))
			})
		})
	})
})
