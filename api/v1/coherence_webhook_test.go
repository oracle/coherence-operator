/*
 * Copyright (c) 2020, 2025, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1_test

import (
	"context"
	. "github.com/onsi/gomega"
	coh "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/pkg/operator"
	"github.com/spf13/viper"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/ptr"
	"testing"
	"time"
)

func TestDefaultReplicasIsSet(t *testing.T) {
	g := NewGomegaWithT(t)

	c := coh.Coherence{}
	err := c.Default(context.Background(), &c)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(c.Spec.CoherenceResourceSpec.Replicas).NotTo(BeNil())
	g.Expect(*c.Spec.CoherenceResourceSpec.Replicas).To(Equal(coh.DefaultReplicas))
}

func TestAddVersionAnnotation(t *testing.T) {
	g := NewGomegaWithT(t)

	c := coh.Coherence{}
	err := c.Default(context.Background(), &c)
	g.Expect(err).NotTo(HaveOccurred())
	an := c.GetAnnotations()
	g.Expect(an).NotTo(BeNil())
	g.Expect(an[coh.AnnotationOperatorVersion]).To(Equal(operator.GetVersion()))
	g.Expect(*c.Spec.CoherenceResourceSpec.Replicas).To(Equal(coh.DefaultReplicas))
}

func TestShouldAddFinalizer(t *testing.T) {
	g := NewGomegaWithT(t)
	c := coh.Coherence{}
	err := c.Default(context.Background(), &c)
	g.Expect(err).NotTo(HaveOccurred())
	finalizers := c.GetFinalizers()
	g.Expect(len(finalizers)).To(Equal(1))
	g.Expect(finalizers).To(ContainElement(coh.CoherenceFinalizer))
}

func TestShouldNotAddFinalizerAgainIfPresent(t *testing.T) {
	g := NewGomegaWithT(t)
	c := coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{
			Name:       "foo",
			Finalizers: []string{coh.CoherenceFinalizer},
		},
	}
	err := c.Default(context.Background(), &c)
	g.Expect(err).NotTo(HaveOccurred())
	finalizers := c.GetFinalizers()
	g.Expect(len(finalizers)).To(Equal(1))
	g.Expect(finalizers).To(ContainElement(coh.CoherenceFinalizer))
}

func TestShouldNotRemoveFinalizersAlreadyPresent(t *testing.T) {
	g := NewGomegaWithT(t)
	c := coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{
			Name:       "foo",
			Finalizers: []string{"foo", "bar"},
		},
	}
	err := c.Default(context.Background(), &c)
	g.Expect(err).NotTo(HaveOccurred())
	finalizers := c.GetFinalizers()
	g.Expect(len(finalizers)).To(Equal(3))
	g.Expect(finalizers).To(ContainElement("foo"))
	g.Expect(finalizers).To(ContainElement("bar"))
	g.Expect(finalizers).To(ContainElement(coh.CoherenceFinalizer))
}

func TestNoNotAddFinalizerToDeletedResource(t *testing.T) {
	g := NewGomegaWithT(t)

	dt := &metav1.Time{
		Time: time.Now(),
	}

	c := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{
			Name:              "foo",
			DeletionTimestamp: dt,
		},
		Spec: coh.CoherenceStatefulSetResourceSpec{},
	}

	err := c.Default(context.Background(), c)
	g.Expect(err).NotTo(HaveOccurred())
	finalizers := c.GetFinalizers()
	g.Expect(finalizers).To(BeNil())
}

func TestDefaultReplicasIsNotOverriddenWhenAlreadySet(t *testing.T) {
	g := NewGomegaWithT(t)

	var replicas int32 = 19
	c := coh.Coherence{
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Replicas: &replicas,
			},
		},
	}
	err := c.Default(context.Background(), &c)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(c.Spec.Replicas).NotTo(BeNil())
	g.Expect(*c.Spec.Replicas).To(Equal(replicas))
}

func TestCoherenceLocalPortIsSet(t *testing.T) {
	g := NewGomegaWithT(t)

	c := coh.Coherence{}
	err := c.Default(context.Background(), &c)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(c.Spec.Coherence).NotTo(BeNil())
	g.Expect(*c.Spec.Coherence.LocalPort).To(Equal(coh.DefaultUnicastPort))
}

func TestCoherenceLocalPortIsNotOverridden(t *testing.T) {
	g := NewGomegaWithT(t)

	var port int32 = 1234

	c := coh.Coherence{
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Coherence: &coh.CoherenceSpec{
					LocalPort: int32Ptr(port),
				},
			},
		},
	}
	err := c.Default(context.Background(), &c)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(c.Spec.Coherence).NotTo(BeNil())
	g.Expect(*c.Spec.Coherence.LocalPort).To(Equal(port))
}

func TestCoherenceLocalPortIsNotSetOnUpdate(t *testing.T) {
	g := NewGomegaWithT(t)

	c := coh.Coherence{}
	c.Status.Phase = coh.ConditionTypeReady
	err := c.Default(context.Background(), &c)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(c.Spec.Coherence).To(BeNil())
}

func TestCoherenceLocalPortAdjustIsSet(t *testing.T) {
	g := NewGomegaWithT(t)

	lpa := intstr.FromInt32(coh.DefaultUnicastPortAdjust)
	c := coh.Coherence{}
	err := c.Default(context.Background(), &c)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(c.Spec.Coherence).NotTo(BeNil())
	g.Expect(*c.Spec.Coherence.LocalPortAdjust).To(Equal(lpa))
}

func TestCoherenceLocalPortAdjustIsNotOverridden(t *testing.T) {
	g := NewGomegaWithT(t)

	lpa := intstr.FromInt32(9876)
	c := coh.Coherence{
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Coherence: &coh.CoherenceSpec{
					LocalPortAdjust: &lpa,
				},
			},
		},
	}
	err := c.Default(context.Background(), &c)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(c.Spec.Coherence).NotTo(BeNil())
	g.Expect(*c.Spec.Coherence.LocalPortAdjust).To(Equal(lpa))
}

func TestCoherenceLocalPortAdjustIsNotSetOnUpdate(t *testing.T) {
	g := NewGomegaWithT(t)

	c := coh.Coherence{}
	c.Status.Phase = coh.ConditionTypeReady
	err := c.Default(context.Background(), &c)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(c.Spec.Coherence).To(BeNil())
}

func TestCoherenceImageIsSet(t *testing.T) {
	g := NewGomegaWithT(t)

	viper.Set(operator.FlagCoherenceImage, "foo")

	c := coh.Coherence{}
	err := c.Default(context.Background(), &c)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(c.Spec.Image).NotTo(BeNil())
	g.Expect(*c.Spec.Image).To(Equal("foo"))
}

func TestCoherenceImageIsNotOverriddenWhenAlreadySet(t *testing.T) {
	g := NewGomegaWithT(t)

	viper.Set(operator.FlagCoherenceImage, "foo")
	image := "bar"
	c := coh.Coherence{
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Image: &image,
			},
		},
	}

	err := c.Default(context.Background(), &c)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(c.Spec.Image).NotTo(BeNil())
	g.Expect(*c.Spec.Image).To(Equal(image))
}

func TestPersistenceModeChangeNotAllowed(t *testing.T) {
	g := NewGomegaWithT(t)

	cm := "on-demand"
	current := &coh.Coherence{
		Spec: coh.CoherenceStatefulSetResourceSpec{
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
	prev := &coh.Coherence{
		Spec: coh.CoherenceStatefulSetResourceSpec{
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

	_, err := current.ValidateUpdate(context.Background(), current, prev)
	g.Expect(err).To(HaveOccurred())
}

func TestPersistenceModeChangeAllowedIfReplicasIsZero(t *testing.T) {
	g := NewGomegaWithT(t)

	cm := "on-demand"
	current := &coh.Coherence{
		Spec: coh.CoherenceStatefulSetResourceSpec{
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
	prev := &coh.Coherence{
		Spec: coh.CoherenceStatefulSetResourceSpec{
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

	_, err := current.ValidateUpdate(context.Background(), current, prev)
	g.Expect(err).NotTo(HaveOccurred())
}

func TestPersistenceModeChangeAllowedIfPreviousReplicasIsZero(t *testing.T) {
	g := NewGomegaWithT(t)

	cm := "on-demand"
	current := &coh.Coherence{
		Spec: coh.CoherenceStatefulSetResourceSpec{
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
	prev := &coh.Coherence{
		Spec: coh.CoherenceStatefulSetResourceSpec{
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

	_, err := current.ValidateUpdate(context.Background(), current, prev)
	g.Expect(err).NotTo(HaveOccurred())
}

func TestPersistenceVolumeChangeNotAllowed(t *testing.T) {
	g := NewGomegaWithT(t)

	cm := "on-demand"
	current := &coh.Coherence{
		Spec: coh.CoherenceStatefulSetResourceSpec{
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
	prev := &coh.Coherence{
		Spec: coh.CoherenceStatefulSetResourceSpec{
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

	_, err := current.ValidateUpdate(context.Background(), current, prev)
	g.Expect(err).To(HaveOccurred())
}

func TestValidateCreateReplicasWhenReplicasIsNil(t *testing.T) {
	g := NewGomegaWithT(t)

	current := &coh.Coherence{
		Spec: coh.CoherenceStatefulSetResourceSpec{},
	}

	_, err := current.ValidateUpdate(context.Background(), current, current)
	g.Expect(err).NotTo(HaveOccurred())
}

func TestValidateCreateReplicasWhenReplicasIsPositive(t *testing.T) {
	g := NewGomegaWithT(t)

	current := &coh.Coherence{
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Replicas: ptr.To(int32(19)),
			},
		},
	}

	_, err := current.ValidateUpdate(context.Background(), current, current)
	g.Expect(err).NotTo(HaveOccurred())
}

func TestValidateCreateReplicasWhenReplicasIsZero(t *testing.T) {
	g := NewGomegaWithT(t)

	current := &coh.Coherence{
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Replicas: ptr.To(int32(19)),
			},
		},
	}

	_, err := current.ValidateUpdate(context.Background(), current, current)
	g.Expect(err).NotTo(HaveOccurred())
}

func TestValidateCreateReplicasWhenReplicasIsInvalid(t *testing.T) {
	g := NewGomegaWithT(t)

	current := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "foo"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Replicas: ptr.To(int32(-1)),
			},
		},
	}

	_, err := current.ValidateUpdate(context.Background(), current, current)
	g.Expect(err).To(HaveOccurred())
}

func TestValidateUpdateReplicasWhenReplicasIsNil(t *testing.T) {
	g := NewGomegaWithT(t)

	current := &coh.Coherence{
		Spec: coh.CoherenceStatefulSetResourceSpec{},
	}

	prev := &coh.Coherence{
		Spec: coh.CoherenceStatefulSetResourceSpec{},
	}

	_, err := current.ValidateUpdate(context.Background(), current, prev)
	g.Expect(err).NotTo(HaveOccurred())
}

func TestValidateUpdateReplicasWhenReplicasIsPositive(t *testing.T) {
	g := NewGomegaWithT(t)

	current := &coh.Coherence{
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Replicas: ptr.To(int32(19)),
			},
		},
	}

	prev := &coh.Coherence{
		Spec: coh.CoherenceStatefulSetResourceSpec{},
	}

	_, err := current.ValidateUpdate(context.Background(), current, prev)
	g.Expect(err).NotTo(HaveOccurred())
}

func TestValidateUpdateReplicasWhenReplicasIsZero(t *testing.T) {
	g := NewGomegaWithT(t)

	current := &coh.Coherence{
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Replicas: ptr.To(int32(19)),
			},
		},
	}

	prev := &coh.Coherence{
		Spec: coh.CoherenceStatefulSetResourceSpec{},
	}

	_, err := current.ValidateUpdate(context.Background(), current, prev)
	g.Expect(err).NotTo(HaveOccurred())
}

func TestValidateUpdateReplicasWhenReplicasIsInvalid(t *testing.T) {
	g := NewGomegaWithT(t)

	current := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "foo"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Replicas: ptr.To(int32(-1)),
			},
		},
	}

	prev := &coh.Coherence{
		Spec: coh.CoherenceStatefulSetResourceSpec{},
	}

	_, err := current.ValidateUpdate(context.Background(), current, prev)
	g.Expect(err).To(HaveOccurred())
}

func TestValidateVolumeClaimUpdateWhenVolumeClaimsNil(t *testing.T) {
	g := NewGomegaWithT(t)

	current := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "foo"},
		Spec:       coh.CoherenceStatefulSetResourceSpec{},
	}

	prev := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "foo"},
		Spec:       coh.CoherenceStatefulSetResourceSpec{},
	}

	_, err := current.ValidateUpdate(context.Background(), current, prev)
	g.Expect(err).NotTo(HaveOccurred())
}

func TestValidateVolumeClaimUpdateWhenVolumeClaimsNilAndEmpty(t *testing.T) {
	g := NewGomegaWithT(t)

	current := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "foo"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			VolumeClaimTemplates: []coh.PersistentVolumeClaim{},
		},
	}

	prev := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "foo"},
		Spec:       coh.CoherenceStatefulSetResourceSpec{},
	}

	_, err := current.ValidateUpdate(context.Background(), current, prev)
	g.Expect(err).NotTo(HaveOccurred())
}

func TestValidateVolumeClaimUpdateWhenVolumeClaimsAdded(t *testing.T) {
	g := NewGomegaWithT(t)

	current := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "foo"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			VolumeClaimTemplates: []coh.PersistentVolumeClaim{
				{
					Metadata: coh.PersistentVolumeClaimObjectMeta{Name: "foo"},
					Spec:     corev1.PersistentVolumeClaimSpec{},
				},
			},
		},
	}

	prev := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "foo"},
		Spec:       coh.CoherenceStatefulSetResourceSpec{},
	}

	_, err := current.ValidateUpdate(context.Background(), current, prev)
	g.Expect(err).To(HaveOccurred())
}

func TestValidateVolumeClaimUpdateWhenVolumeClaimsRemoved(t *testing.T) {
	g := NewGomegaWithT(t)

	current := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "foo"},
		Spec:       coh.CoherenceStatefulSetResourceSpec{},
	}

	prev := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "foo"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			VolumeClaimTemplates: []coh.PersistentVolumeClaim{
				{
					Metadata: coh.PersistentVolumeClaimObjectMeta{Name: "foo"},
					Spec:     corev1.PersistentVolumeClaimSpec{},
				},
			},
		},
	}

	_, err := current.ValidateUpdate(context.Background(), current, prev)
	g.Expect(err).To(HaveOccurred())
}

func TestValidateNodePortsOnCreateWithValidPorts(t *testing.T) {
	g := NewGomegaWithT(t)

	current := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "foo"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
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

	_, err := current.ValidateUpdate(context.Background(), current, current)
	g.Expect(err).NotTo(HaveOccurred())
}

func TestValidateNodePortsOnCreateWithInvalidPort(t *testing.T) {
	g := NewGomegaWithT(t)

	current := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "foo"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
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

	_, err := current.ValidateUpdate(context.Background(), current, current)
	g.Expect(err).To(HaveOccurred())
}

func TestValidateNodePortsOnUpdateWithValidPorts(t *testing.T) {
	g := NewGomegaWithT(t)

	current := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "foo"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
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

	prev := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "foo"},
		Spec:       coh.CoherenceStatefulSetResourceSpec{},
	}

	_, err := current.ValidateUpdate(context.Background(), current, prev)
	g.Expect(err).NotTo(HaveOccurred())
}

func TestValidateNodePortsOnUpdateWithInvalidPort(t *testing.T) {
	g := NewGomegaWithT(t)

	current := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "foo"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
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

	prev := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "foo"},
		Spec:       coh.CoherenceStatefulSetResourceSpec{},
	}

	_, err := current.ValidateUpdate(context.Background(), prev, current)
	g.Expect(err).To(HaveOccurred())
}
