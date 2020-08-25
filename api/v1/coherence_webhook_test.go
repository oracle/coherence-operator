/*
 * Copyright (c) 2020, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 *
 */

package v1_test

import (
	. "github.com/onsi/gomega"
	coh "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/pkg/operator"
	"github.com/spf13/viper"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
	"testing"
)

func TestDefaultReplicasIsSet(t *testing.T) {
	g := NewGomegaWithT(t)

	c := coh.Coherence{}
	c.Default()
	g.Expect(c.Spec.Replicas).NotTo(BeNil())
	g.Expect(*c.Spec.Replicas).To(Equal(coh.DefaultReplicas))
}

func TestDefaultReplicasIsNotOverriddenWhenAlreadySet(t *testing.T) {
	g := NewGomegaWithT(t)

	var replicas int32 = 19
	c := coh.Coherence{
		Spec: coh.CoherenceResourceSpec{
			Replicas: &replicas,
		},
	}
	c.Default()
	g.Expect(c.Spec.Replicas).NotTo(BeNil())
	g.Expect(*c.Spec.Replicas).To(Equal(replicas))
}

func TestCoherenceImageIsSet(t *testing.T) {
	g := NewGomegaWithT(t)

	viper.Set(operator.FlagCoherenceImage, "foo")

	c := coh.Coherence{}
	c.Default()
	g.Expect(c.Spec.Image).NotTo(BeNil())
	g.Expect(*c.Spec.Image).To(Equal("foo"))
}

func TestCoherenceImageIsNotOverriddenWhenAlreadySet(t *testing.T) {
	g := NewGomegaWithT(t)

	viper.Set(operator.FlagCoherenceImage, "foo")
	image := "bar"
	c := coh.Coherence{
		Spec: coh.CoherenceResourceSpec{
			Image: &image,
		},
	}

	c.Default()
	g.Expect(c.Spec.Image).NotTo(BeNil())
	g.Expect(*c.Spec.Image).To(Equal(image))
}

func TestUtilsImageIsSet(t *testing.T) {
	g := NewGomegaWithT(t)

	viper.Set(operator.FlagUtilsImage, "foo")

	c := coh.Coherence{}
	c.Default()
	g.Expect(c.Spec.CoherenceUtils).NotTo(BeNil())
	g.Expect(c.Spec.CoherenceUtils.Image).NotTo(BeNil())
	g.Expect(*c.Spec.CoherenceUtils.Image).To(Equal("foo"))
}

func TestUtilsImageIsNotOverriddenWhenAlreadySet(t *testing.T) {
	g := NewGomegaWithT(t)

	viper.Set(operator.FlagUtilsImage, "foo")
	image := "bar"
	c := coh.Coherence{
		Spec: coh.CoherenceResourceSpec{
			CoherenceUtils: &coh.ImageSpec{
				Image: &image,
			},
		},
	}

	c.Default()
	g.Expect(c.Spec.CoherenceUtils).NotTo(BeNil())
	g.Expect(c.Spec.CoherenceUtils.Image).NotTo(BeNil())
	g.Expect(*c.Spec.CoherenceUtils.Image).To(Equal(image))
}

func TestPersistenceModeChangeNotAllowed(t *testing.T) {
	g := NewGomegaWithT(t)

	cm := "on-demand"
	current := &coh.Coherence{
		Spec: coh.CoherenceResourceSpec{
			Coherence: &coh.CoherenceSpec{
				Persistence: &coh.PersistenceSpec{
					Mode:                  &cm,
					PersistentStorageSpec: coh.PersistentStorageSpec{},
					Snapshots:             &coh.PersistentStorageSpec{
						PersistentVolumeClaim: &corev1.PersistentVolumeClaimSpec{},
						Volume:                &corev1.VolumeSource{},
					},
				},
			},
		},
	}

	pm := "active"
	prev := &coh.Coherence{
		Spec: coh.CoherenceResourceSpec{
			Coherence: &coh.CoherenceSpec{
				Persistence: &coh.PersistenceSpec{
					Mode:                  &pm,
					PersistentStorageSpec: coh.PersistentStorageSpec{},
					Snapshots:             &coh.PersistentStorageSpec{
						PersistentVolumeClaim: &corev1.PersistentVolumeClaimSpec{},
						Volume:                &corev1.VolumeSource{},
					},
				},
			},
		},
	}

	err := current.ValidateUpdate(prev)
	g.Expect(err).To(HaveOccurred())
}

func TestPersistenceModeChangeAllowedIfReplicasIsZero(t *testing.T) {
	g := NewGomegaWithT(t)

	cm := "on-demand"
	current := &coh.Coherence{
		Spec: coh.CoherenceResourceSpec{
			Replicas: int32Ptr(0),
			Coherence: &coh.CoherenceSpec{
				Persistence: &coh.PersistenceSpec{
					Mode:                  &cm,
					PersistentStorageSpec: coh.PersistentStorageSpec{},
					Snapshots:             &coh.PersistentStorageSpec{
						PersistentVolumeClaim: &corev1.PersistentVolumeClaimSpec{},
						Volume:                &corev1.VolumeSource{},
					},
				},
			},
		},
	}

	pm := "active"
	prev := &coh.Coherence{
		Spec: coh.CoherenceResourceSpec{
			Coherence: &coh.CoherenceSpec{
				Persistence: &coh.PersistenceSpec{
					Mode:                  &pm,
					PersistentStorageSpec: coh.PersistentStorageSpec{},
					Snapshots:             &coh.PersistentStorageSpec{
						PersistentVolumeClaim: &corev1.PersistentVolumeClaimSpec{},
						Volume:                &corev1.VolumeSource{},
					},
				},
			},
		},
	}

	err := current.ValidateUpdate(prev)
	g.Expect(err).NotTo(HaveOccurred())
}

func TestPersistenceModeChangeAllowedIfPreviousReplicasIsZero(t *testing.T) {
	g := NewGomegaWithT(t)

	cm := "on-demand"
	current := &coh.Coherence{
		Spec: coh.CoherenceResourceSpec{
			Coherence: &coh.CoherenceSpec{
				Persistence: &coh.PersistenceSpec{
					Mode:                  &cm,
					PersistentStorageSpec: coh.PersistentStorageSpec{},
					Snapshots:             &coh.PersistentStorageSpec{
						PersistentVolumeClaim: &corev1.PersistentVolumeClaimSpec{},
						Volume:                &corev1.VolumeSource{},
					},
				},
			},
		},
	}

	pm := "active"
	prev := &coh.Coherence{
		Spec: coh.CoherenceResourceSpec{
			Replicas: int32Ptr(0),
			Coherence: &coh.CoherenceSpec{
				Persistence: &coh.PersistenceSpec{
					Mode:                  &pm,
					PersistentStorageSpec: coh.PersistentStorageSpec{},
					Snapshots:             &coh.PersistentStorageSpec{
						PersistentVolumeClaim: &corev1.PersistentVolumeClaimSpec{},
						Volume:                &corev1.VolumeSource{},
					},
				},
			},
		},
	}

	err := current.ValidateUpdate(prev)
	g.Expect(err).NotTo(HaveOccurred())
}

func TestPersistenceVolumeChangeNotAllowed(t *testing.T) {
	g := NewGomegaWithT(t)

	cm := "on-demand"
	current := &coh.Coherence{
		Spec: coh.CoherenceResourceSpec{
			Coherence: &coh.CoherenceSpec{
				Persistence: &coh.PersistenceSpec{
					Mode:                  &cm,
					PersistentStorageSpec: coh.PersistentStorageSpec{
						PersistentVolumeClaim: &corev1.PersistentVolumeClaimSpec{
							VolumeName: "foo",
						},
						Volume:                &corev1.VolumeSource{},
					},
					Snapshots:             &coh.PersistentStorageSpec{
						PersistentVolumeClaim: &corev1.PersistentVolumeClaimSpec{},
						Volume:                &corev1.VolumeSource{},
					},
				},
			},
		},
	}

	pm := "active"
	prev := &coh.Coherence{
		Spec: coh.CoherenceResourceSpec{
			Coherence: &coh.CoherenceSpec{
				Persistence: &coh.PersistenceSpec{
					Mode:                  &pm,
					PersistentStorageSpec: coh.PersistentStorageSpec{
						PersistentVolumeClaim: &corev1.PersistentVolumeClaimSpec{
							VolumeName: "bar",
						},
						Volume:                &corev1.VolumeSource{},
					},
					Snapshots:             &coh.PersistentStorageSpec{
						PersistentVolumeClaim: &corev1.PersistentVolumeClaimSpec{},
						Volume:                &corev1.VolumeSource{},
					},
				},
			},
		},
	}

	err := current.ValidateUpdate(prev)
	g.Expect(err).To(HaveOccurred())
}

func TestValidateCreateReplicasWhenReplicasIsNil(t *testing.T) {
	g := NewGomegaWithT(t)

	current := &coh.Coherence{
		Spec: coh.CoherenceResourceSpec{},
	}

	err := current.ValidateCreate()
	g.Expect(err).NotTo(HaveOccurred())
}

func TestValidateCreateReplicasWhenReplicasIsPositive(t *testing.T) {
	g := NewGomegaWithT(t)

	current := &coh.Coherence{
		Spec: coh.CoherenceResourceSpec{
			Replicas: pointer.Int32Ptr(19),
		},
	}

	err := current.ValidateCreate()
	g.Expect(err).NotTo(HaveOccurred())
}

func TestValidateCreateReplicasWhenReplicasIsZero(t *testing.T) {
	g := NewGomegaWithT(t)

	current := &coh.Coherence{
		Spec: coh.CoherenceResourceSpec{
			Replicas: pointer.Int32Ptr(19),
		},
	}

	err := current.ValidateCreate()
	g.Expect(err).NotTo(HaveOccurred())
}

func TestValidateCreateReplicasWhenReplicasIsInvalid(t *testing.T) {
	g := NewGomegaWithT(t)

	current := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "foo"},
		Spec: coh.CoherenceResourceSpec{
			Replicas: pointer.Int32Ptr(-1),
		},
	}

	err := current.ValidateCreate()
	g.Expect(err).To(HaveOccurred())
}

func TestValidateUpdateReplicasWhenReplicasIsNil(t *testing.T) {
	g := NewGomegaWithT(t)

	current := &coh.Coherence{
		Spec: coh.CoherenceResourceSpec{},
	}

	prev := &coh.Coherence{
		Spec: coh.CoherenceResourceSpec{},
	}

	err := current.ValidateUpdate(prev)
	g.Expect(err).NotTo(HaveOccurred())
}

func TestValidateUpdateReplicasWhenReplicasIsPositive(t *testing.T) {
	g := NewGomegaWithT(t)

	current := &coh.Coherence{
		Spec: coh.CoherenceResourceSpec{
			Replicas: pointer.Int32Ptr(19),
		},
	}

	prev := &coh.Coherence{
		Spec: coh.CoherenceResourceSpec{},
	}

	err := current.ValidateUpdate(prev)
	g.Expect(err).NotTo(HaveOccurred())
}

func TestValidateUpdateReplicasWhenReplicasIsZero(t *testing.T) {
	g := NewGomegaWithT(t)

	current := &coh.Coherence{
		Spec: coh.CoherenceResourceSpec{
			Replicas: pointer.Int32Ptr(19),
		},
	}

	prev := &coh.Coherence{
		Spec: coh.CoherenceResourceSpec{},
	}

	err := current.ValidateUpdate(prev)
	g.Expect(err).NotTo(HaveOccurred())
}

func TestValidateUpdateReplicasWhenReplicasIsInvalid(t *testing.T) {
	g := NewGomegaWithT(t)

	current := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "foo"},
		Spec: coh.CoherenceResourceSpec{
			Replicas: pointer.Int32Ptr(-1),
		},
	}

	prev := &coh.Coherence{
		Spec: coh.CoherenceResourceSpec{},
	}

	err := current.ValidateUpdate(prev)
	g.Expect(err).To(HaveOccurred())
}
