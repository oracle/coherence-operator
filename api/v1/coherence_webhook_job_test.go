/*
 * Copyright (c) 2020, 2024, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1_test

import (
	. "github.com/onsi/gomega"
	coh "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/pkg/operator"
	"github.com/spf13/viper"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/ptr"
	"testing"
)

func TestJobDefaultReplicasIsSet(t *testing.T) {
	g := NewGomegaWithT(t)

	c := coh.CoherenceJob{}
	c.Default()
	g.Expect(c.Spec.Replicas).NotTo(BeNil())
	g.Expect(*c.Spec.Replicas).To(Equal(coh.DefaultJobReplicas))
}

func TestJobAddVersionAnnotation(t *testing.T) {
	g := NewGomegaWithT(t)

	c := coh.CoherenceJob{}
	c.Default()
	g.Expect(c.Annotations).NotTo(BeNil())
	g.Expect(c.Annotations[coh.AnnotationOperatorVersion]).To(Equal(operator.GetVersion()))
	g.Expect(c.Spec).NotTo(BeNil())
	replicas := c.Spec.CoherenceResourceSpec.Replicas
	g.Expect(*replicas).To(Equal(coh.DefaultJobReplicas))
}

func TestJobShouldAddFinalizer(t *testing.T) {
	g := NewGomegaWithT(t)
	c := coh.CoherenceJob{}
	c.Default()
	finalizers := c.GetFinalizers()
	g.Expect(len(finalizers)).To(Equal(0))
}

func TestJobShouldNotAddFinalizerAgainIfPresent(t *testing.T) {
	g := NewGomegaWithT(t)
	c := coh.CoherenceJob{
		ObjectMeta: metav1.ObjectMeta{
			Name:       "foo",
			Finalizers: []string{coh.CoherenceFinalizer},
		},
	}
	c.Default()
	finalizers := c.GetFinalizers()
	g.Expect(len(finalizers)).To(Equal(1))
	g.Expect(finalizers).To(ContainElement(coh.CoherenceFinalizer))
}

func TestJobShouldNotRemoveFinalizersAlreadyPresent(t *testing.T) {
	g := NewGomegaWithT(t)
	c := coh.CoherenceJob{
		ObjectMeta: metav1.ObjectMeta{
			Name:       "foo",
			Finalizers: []string{"foo", "bar"},
		},
	}
	c.Default()
	finalizers := c.GetFinalizers()
	g.Expect(len(finalizers)).To(Equal(2))
	g.Expect(finalizers).To(ContainElement("foo"))
	g.Expect(finalizers).To(ContainElement("bar"))
}

func TestJobDefaultReplicasIsNotOverriddenWhenAlreadySet(t *testing.T) {
	g := NewGomegaWithT(t)

	var replicas int32 = 19
	c := coh.CoherenceJob{
		Spec: coh.CoherenceJobResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Replicas: &replicas,
			},
		},
	}
	c.Default()
	g.Expect(c.Spec.Replicas).NotTo(BeNil())
	g.Expect(*c.Spec.Replicas).To(Equal(replicas))
}

func TestJobCoherenceLocalPortIsSet(t *testing.T) {
	g := NewGomegaWithT(t)

	c := coh.CoherenceJob{}
	c.Default()
	g.Expect(c.Spec.Coherence).NotTo(BeNil())
	g.Expect(*c.Spec.Coherence.LocalPort).To(Equal(coh.DefaultUnicastPort))
}

func TestJobCoherenceLocalPortIsNotOverridden(t *testing.T) {
	g := NewGomegaWithT(t)

	var port int32 = 1234

	c := coh.CoherenceJob{
		Spec: coh.CoherenceJobResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Coherence: &coh.CoherenceSpec{
					LocalPort: int32Ptr(port),
				},
			},
		},
	}
	c.Default()
	g.Expect(c.Spec.Coherence).NotTo(BeNil())
	g.Expect(*c.Spec.Coherence.LocalPort).To(Equal(port))
}

func TestJobCoherenceLocalPortIsNotSetOnUpdate(t *testing.T) {
	g := NewGomegaWithT(t)

	c := coh.CoherenceJob{}
	c.Status.Phase = coh.ConditionTypeReady
	c.Default()
	g.Expect(c.Spec.Coherence).NotTo(BeNil())
	g.Expect(c.Spec.Coherence.LocalPort).To(BeNil())
	g.Expect(c.Spec.Coherence.LocalPortAdjust).To(BeNil())
}

func TestJobCoherenceLocalPortAdjustIsSet(t *testing.T) {
	g := NewGomegaWithT(t)

	lpa := intstr.FromInt32(coh.DefaultUnicastPortAdjust)
	c := coh.CoherenceJob{}
	c.Default()
	g.Expect(c.Spec.Coherence).NotTo(BeNil())
	g.Expect(*c.Spec.CoherenceResourceSpec.Coherence.LocalPortAdjust).To(Equal(lpa))
}

func TestJobCoherenceLocalPortAdjustIsNotOverridden(t *testing.T) {
	g := NewGomegaWithT(t)

	lpa := intstr.FromInt32(9876)
	c := coh.CoherenceJob{
		Spec: coh.CoherenceJobResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Coherence: &coh.CoherenceSpec{
					LocalPortAdjust: &lpa,
				},
			},
		},
	}
	c.Default()
	g.Expect(c.Spec.Coherence).NotTo(BeNil())
	g.Expect(*c.Spec.Coherence.LocalPortAdjust).To(Equal(lpa))
}

func TestJobCoherenceLocalPortAdjustIsNotSetOnUpdate(t *testing.T) {
	g := NewGomegaWithT(t)

	c := coh.CoherenceJob{}
	c.Status.Phase = coh.ConditionTypeReady
	c.Default()
	g.Expect(c.Spec.Coherence).NotTo(BeNil())
	g.Expect(c.Spec.Coherence.LocalPortAdjust).To(BeNil())
}

func TestJobCoherenceImageIsSet(t *testing.T) {
	g := NewGomegaWithT(t)

	viper.Set(operator.FlagCoherenceImage, "foo")

	c := coh.CoherenceJob{}
	c.Default()
	g.Expect(c.Spec.Image).NotTo(BeNil())
	g.Expect(*c.Spec.Image).To(Equal("foo"))
}

func TestJobCoherenceImageIsNotOverriddenWhenAlreadySet(t *testing.T) {
	g := NewGomegaWithT(t)

	viper.Set(operator.FlagCoherenceImage, "foo")
	image := "bar"
	c := coh.CoherenceJob{
		Spec: coh.CoherenceJobResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Image: &image,
			},
		},
	}

	c.Default()
	g.Expect(c.Spec.Image).NotTo(BeNil())
	g.Expect(*c.Spec.Image).To(Equal(image))
}

func TestJobUtilsImageIsSet(t *testing.T) {
	g := NewGomegaWithT(t)

	viper.Set(operator.FlagOperatorImage, "foo")

	c := coh.CoherenceJob{}
	c.Default()
	g.Expect(c.Spec.CoherenceUtils).NotTo(BeNil())
	g.Expect(c.Spec.CoherenceUtils.Image).NotTo(BeNil())
	g.Expect(*c.Spec.CoherenceUtils.Image).To(Equal("foo"))
}

func TestJobUtilsImageIsNotOverriddenWhenAlreadySet(t *testing.T) {
	g := NewGomegaWithT(t)

	viper.Set(operator.FlagOperatorImage, "foo")
	image := "bar"
	c := coh.CoherenceJob{
		Spec: coh.CoherenceJobResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				CoherenceUtils: &coh.ImageSpec{
					Image: &image,
				},
			},
		},
	}

	c.Default()
	g.Expect(c.Spec.CoherenceUtils).NotTo(BeNil())
	g.Expect(c.Spec.CoherenceUtils.Image).NotTo(BeNil())
	g.Expect(*c.Spec.CoherenceUtils.Image).To(Equal(image))
}

func TestJobPersistenceModeChangeNotAllowed(t *testing.T) {
	g := NewGomegaWithT(t)

	cm := "on-demand"
	current := &coh.CoherenceJob{
		Spec: coh.CoherenceJobResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Coherence: &coh.CoherenceSpec{
					Persistence: &coh.PersistenceSpec{
						Mode:                  &cm,
						PersistentStorageSpec: coh.PersistentStorageSpec{},
						Snapshots: &coh.PersistentStorageSpec{
							PersistentVolumeClaim: &corev1.PersistentVolumeClaimSpec{},
							Volume:                &corev1.VolumeSource{},
						},
					},
				},
			},
		},
	}

	pm := "active"
	prev := &coh.CoherenceJob{
		Spec: coh.CoherenceJobResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Coherence: &coh.CoherenceSpec{
					Persistence: &coh.PersistenceSpec{
						Mode:                  &pm,
						PersistentStorageSpec: coh.PersistentStorageSpec{},
						Snapshots: &coh.PersistentStorageSpec{
							PersistentVolumeClaim: &corev1.PersistentVolumeClaimSpec{},
							Volume:                &corev1.VolumeSource{},
						},
					},
				},
			},
		},
	}

	_, err := current.ValidateUpdate(prev)
	g.Expect(err).To(HaveOccurred())
}

func TestJobPersistenceModeChangeAllowedIfReplicasIsZero(t *testing.T) {
	g := NewGomegaWithT(t)

	cm := "on-demand"
	current := &coh.CoherenceJob{
		Spec: coh.CoherenceJobResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Replicas: int32Ptr(0),
				Coherence: &coh.CoherenceSpec{
					Persistence: &coh.PersistenceSpec{
						Mode:                  &cm,
						PersistentStorageSpec: coh.PersistentStorageSpec{},
						Snapshots: &coh.PersistentStorageSpec{
							PersistentVolumeClaim: &corev1.PersistentVolumeClaimSpec{},
							Volume:                &corev1.VolumeSource{},
						},
					},
				},
			},
		},
	}

	pm := "active"
	prev := &coh.CoherenceJob{
		Spec: coh.CoherenceJobResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Coherence: &coh.CoherenceSpec{
					Persistence: &coh.PersistenceSpec{
						Mode:                  &pm,
						PersistentStorageSpec: coh.PersistentStorageSpec{},
						Snapshots: &coh.PersistentStorageSpec{
							PersistentVolumeClaim: &corev1.PersistentVolumeClaimSpec{},
							Volume:                &corev1.VolumeSource{},
						},
					},
				},
			},
		},
	}

	_, err := current.ValidateUpdate(prev)
	g.Expect(err).NotTo(HaveOccurred())
}

func TestJobPersistenceModeChangeAllowedIfPreviousReplicasIsZero(t *testing.T) {
	g := NewGomegaWithT(t)

	cm := "on-demand"
	current := &coh.CoherenceJob{
		Spec: coh.CoherenceJobResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Coherence: &coh.CoherenceSpec{
					Persistence: &coh.PersistenceSpec{
						Mode:                  &cm,
						PersistentStorageSpec: coh.PersistentStorageSpec{},
						Snapshots: &coh.PersistentStorageSpec{
							PersistentVolumeClaim: &corev1.PersistentVolumeClaimSpec{},
							Volume:                &corev1.VolumeSource{},
						},
					},
				},
			},
		},
	}

	pm := "active"
	prev := &coh.CoherenceJob{
		Spec: coh.CoherenceJobResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Replicas: int32Ptr(0),
				Coherence: &coh.CoherenceSpec{
					Persistence: &coh.PersistenceSpec{
						Mode:                  &pm,
						PersistentStorageSpec: coh.PersistentStorageSpec{},
						Snapshots: &coh.PersistentStorageSpec{
							PersistentVolumeClaim: &corev1.PersistentVolumeClaimSpec{},
							Volume:                &corev1.VolumeSource{},
						},
					},
				},
			},
		},
	}

	_, err := current.ValidateUpdate(prev)
	g.Expect(err).NotTo(HaveOccurred())
}

func TestJobPersistenceVolumeChangeNotAllowed(t *testing.T) {
	g := NewGomegaWithT(t)

	cm := "on-demand"
	current := &coh.CoherenceJob{
		Spec: coh.CoherenceJobResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Coherence: &coh.CoherenceSpec{
					Persistence: &coh.PersistenceSpec{
						Mode: &cm,
						PersistentStorageSpec: coh.PersistentStorageSpec{
							PersistentVolumeClaim: &corev1.PersistentVolumeClaimSpec{
								VolumeName: "foo",
							},
							Volume: &corev1.VolumeSource{},
						},
						Snapshots: &coh.PersistentStorageSpec{
							PersistentVolumeClaim: &corev1.PersistentVolumeClaimSpec{},
							Volume:                &corev1.VolumeSource{},
						},
					},
				},
			},
		},
	}

	pm := "active"
	prev := &coh.CoherenceJob{
		Spec: coh.CoherenceJobResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Coherence: &coh.CoherenceSpec{
					Persistence: &coh.PersistenceSpec{
						Mode: &pm,
						PersistentStorageSpec: coh.PersistentStorageSpec{
							PersistentVolumeClaim: &corev1.PersistentVolumeClaimSpec{
								VolumeName: "bar",
							},
							Volume: &corev1.VolumeSource{},
						},
						Snapshots: &coh.PersistentStorageSpec{
							PersistentVolumeClaim: &corev1.PersistentVolumeClaimSpec{},
							Volume:                &corev1.VolumeSource{},
						},
					},
				},
			},
		},
	}

	_, err := current.ValidateUpdate(prev)
	g.Expect(err).To(HaveOccurred())
}

func TestJobValidateCreateReplicasWhenReplicasIsNil(t *testing.T) {
	g := NewGomegaWithT(t)

	current := &coh.CoherenceJob{
		Spec: coh.CoherenceJobResourceSpec{},
	}

	_, err := current.ValidateCreate()
	g.Expect(err).NotTo(HaveOccurred())
}

func TestJobValidateCreateReplicasWhenReplicasIsPositive(t *testing.T) {
	g := NewGomegaWithT(t)

	current := &coh.CoherenceJob{
		Spec: coh.CoherenceJobResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Replicas: ptr.To(int32(19)),
			},
		},
	}

	_, err := current.ValidateCreate()
	g.Expect(err).NotTo(HaveOccurred())
}

func TestJobValidateCreateReplicasWhenReplicasIsZero(t *testing.T) {
	g := NewGomegaWithT(t)

	current := &coh.CoherenceJob{
		Spec: coh.CoherenceJobResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Replicas: ptr.To(int32(19)),
			},
		},
	}

	_, err := current.ValidateCreate()
	g.Expect(err).NotTo(HaveOccurred())
}

func TestJobValidateCreateReplicasWhenReplicasIsInvalid(t *testing.T) {
	g := NewGomegaWithT(t)

	current := &coh.CoherenceJob{
		ObjectMeta: metav1.ObjectMeta{Name: "foo"},
		Spec: coh.CoherenceJobResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Replicas: ptr.To(int32(-1)),
			},
		},
	}

	_, err := current.ValidateCreate()
	g.Expect(err).To(HaveOccurred())
}

func TestJobValidateUpdateReplicasWhenReplicasIsNil(t *testing.T) {
	g := NewGomegaWithT(t)

	current := &coh.CoherenceJob{
		Spec: coh.CoherenceJobResourceSpec{},
	}

	prev := &coh.CoherenceJob{
		Spec: coh.CoherenceJobResourceSpec{},
	}

	_, err := current.ValidateUpdate(prev)
	g.Expect(err).NotTo(HaveOccurred())
}

func TestJobValidateUpdateReplicasWhenReplicasIsPositive(t *testing.T) {
	g := NewGomegaWithT(t)

	current := &coh.CoherenceJob{
		Spec: coh.CoherenceJobResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Replicas: ptr.To(int32(9)),
			},
		},
	}

	prev := &coh.CoherenceJob{
		Spec: coh.CoherenceJobResourceSpec{},
	}

	_, err := current.ValidateUpdate(prev)
	g.Expect(err).NotTo(HaveOccurred())
}

func TestJobValidateUpdateReplicasWhenReplicasIsZero(t *testing.T) {
	g := NewGomegaWithT(t)

	current := &coh.CoherenceJob{
		Spec: coh.CoherenceJobResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Replicas: ptr.To(int32(19)),
			},
		},
	}

	prev := &coh.CoherenceJob{
		Spec: coh.CoherenceJobResourceSpec{},
	}

	_, err := current.ValidateUpdate(prev)
	g.Expect(err).NotTo(HaveOccurred())
}

func TestJobValidateUpdateReplicasWhenReplicasIsInvalid(t *testing.T) {
	g := NewGomegaWithT(t)

	current := &coh.CoherenceJob{
		ObjectMeta: metav1.ObjectMeta{Name: "foo"},
		Spec: coh.CoherenceJobResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Replicas: ptr.To(int32(-1)),
			},
		},
	}

	prev := &coh.CoherenceJob{
		Spec: coh.CoherenceJobResourceSpec{},
	}

	_, err := current.ValidateUpdate(prev)
	g.Expect(err).To(HaveOccurred())
}

func TestJobValidateVolumeClaimUpdateWhenVolumeClaimsNil(t *testing.T) {
	g := NewGomegaWithT(t)

	current := &coh.CoherenceJob{
		ObjectMeta: metav1.ObjectMeta{Name: "foo"},
		Spec:       coh.CoherenceJobResourceSpec{},
	}

	prev := &coh.CoherenceJob{
		ObjectMeta: metav1.ObjectMeta{Name: "foo"},
		Spec:       coh.CoherenceJobResourceSpec{},
	}

	_, err := current.ValidateUpdate(prev)
	g.Expect(err).NotTo(HaveOccurred())
}

func TestJobValidateNodePortsOnCreateWithValidPorts(t *testing.T) {
	g := NewGomegaWithT(t)

	current := &coh.CoherenceJob{
		ObjectMeta: metav1.ObjectMeta{Name: "foo"},
		Spec: coh.CoherenceJobResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Ports: []coh.NamedPortSpec{
					{
						Name:     "p1",
						Port:     1234,
						NodePort: ptr.To(int32(30000)),
					},
					{
						Name:     "p2",
						Port:     1235,
						NodePort: ptr.To(int32(32767)),
					},
				},
			},
		},
	}

	_, err := current.ValidateCreate()
	g.Expect(err).NotTo(HaveOccurred())
}

func TestJobValidateNodePortsOnCreateWithInvalidPort(t *testing.T) {
	g := NewGomegaWithT(t)

	current := &coh.CoherenceJob{
		ObjectMeta: metav1.ObjectMeta{Name: "foo"},
		Spec: coh.CoherenceJobResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Ports: []coh.NamedPortSpec{
					{
						Name:     "p1",
						Port:     1234,
						NodePort: ptr.To(int32(30000)),
					},
					{
						Name:     "p2",
						Port:     1235,
						NodePort: ptr.To(int32(32767)),
					},
					{
						Name:     "p3",
						Port:     1235,
						NodePort: ptr.To(int32(1234)),
					},
				},
			},
		},
	}

	_, err := current.ValidateCreate()
	g.Expect(err).To(HaveOccurred())
}

func TestJobValidateNodePortsOnUpdateWithValidPorts(t *testing.T) {
	g := NewGomegaWithT(t)

	current := &coh.CoherenceJob{
		ObjectMeta: metav1.ObjectMeta{Name: "foo"},
		Spec: coh.CoherenceJobResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Ports: []coh.NamedPortSpec{
					{
						Name:     "p1",
						Port:     1234,
						NodePort: ptr.To(int32(30000)),
					},
					{
						Name:     "p2",
						Port:     1235,
						NodePort: ptr.To(int32(32767)),
					},
				},
			},
		},
	}

	prev := &coh.CoherenceJob{
		ObjectMeta: metav1.ObjectMeta{Name: "foo"},
		Spec:       coh.CoherenceJobResourceSpec{},
	}

	_, err := current.ValidateUpdate(prev)
	g.Expect(err).NotTo(HaveOccurred())
}

func TestJobValidateNodePortsOnUpdateWithInvalidPort(t *testing.T) {
	g := NewGomegaWithT(t)

	current := &coh.CoherenceJob{
		ObjectMeta: metav1.ObjectMeta{Name: "foo"},
		Spec: coh.CoherenceJobResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Ports: []coh.NamedPortSpec{
					{
						Name:     "p1",
						Port:     1234,
						NodePort: ptr.To(int32(30000)),
					},
					{
						Name:     "p2",
						Port:     1235,
						NodePort: ptr.To(int32(32767)),
					},
					{
						Name:     "p3",
						Port:     1235,
						NodePort: ptr.To(int32(1234)),
					},
				},
			},
		},
	}

	prev := &coh.CoherenceJob{
		ObjectMeta: metav1.ObjectMeta{Name: "foo"},
		Spec:       coh.CoherenceJobResourceSpec{},
	}

	_, err := current.ValidateUpdate(prev)
	g.Expect(err).To(HaveOccurred())
}
