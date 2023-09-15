/*
 * Copyright (c) 2023, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1_test

import (
	. "github.com/onsi/gomega"
	coh "github.com/oracle/coherence-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	"testing"
)

// Tests for image name validation

func TestCoherenceWithNoImageNames(t *testing.T) {
	g := NewGomegaWithT(t)

	c := coh.Coherence{}
	_, err := c.ValidateCreate()
	g.Expect(err).NotTo(HaveOccurred())
}

func TestCoherenceCreateWithValidImageName(t *testing.T) {
	g := NewGomegaWithT(t)

	c := coh.Coherence{
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Image: stringPtr("test/coherence:1.0"),
			},
		},
	}
	_, err := c.ValidateCreate()
	g.Expect(err).NotTo(HaveOccurred())
}

func TestCoherenceCreateWithInvalidImageName(t *testing.T) {
	g := NewGomegaWithT(t)

	c := coh.Coherence{
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Image: stringPtr("test/bad image name:1.0"),
			},
		},
	}
	_, err := c.ValidateCreate()
	g.Expect(err).To(HaveOccurred())
}

func TestCoherenceCreateWithImageNameWithTrailingSpace(t *testing.T) {
	g := NewGomegaWithT(t)

	c := coh.Coherence{
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Image: stringPtr("test/coherence:1.0 "),
			},
		},
	}
	_, err := c.ValidateCreate()
	g.Expect(err).To(HaveOccurred())
}

func TestCoherenceCreateWithValidOperatorImageName(t *testing.T) {
	g := NewGomegaWithT(t)

	c := coh.Coherence{
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				CoherenceUtils: &coh.ImageSpec{
					Image: stringPtr("test/coherence:1.0"),
				},
			},
		},
	}
	_, err := c.ValidateCreate()
	g.Expect(err).NotTo(HaveOccurred())
}

func TestCoherenceCreateWithInvalidOperatorImageName(t *testing.T) {
	g := NewGomegaWithT(t)

	c := coh.Coherence{
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				CoherenceUtils: &coh.ImageSpec{
					Image: stringPtr("test/bad image name:1.0"),
				},
			},
		},
	}
	_, err := c.ValidateCreate()
	g.Expect(err).To(HaveOccurred())
}

func TestCoherenceUpdateWithInvalidImageName(t *testing.T) {
	g := NewGomegaWithT(t)

	c := coh.Coherence{
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Image: stringPtr("test/bad image name:1.0"),
			},
		},
	}
	_, err := c.ValidateUpdate(&c)
	g.Expect(err).To(HaveOccurred())
}

func TestCoherenceCreateWithValidInitContainerImageName(t *testing.T) {
	g := NewGomegaWithT(t)

	c := coh.Coherence{
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				InitContainers: []corev1.Container{
					{
						Name:  "side-one",
						Image: "test/coherence:1.0",
					},
				},
			},
		},
	}
	_, err := c.ValidateCreate()
	g.Expect(err).NotTo(HaveOccurred())
}

func TestCoherenceCreateWithInvalidInitContainerImageName(t *testing.T) {
	g := NewGomegaWithT(t)

	c := coh.Coherence{
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				InitContainers: []corev1.Container{
					{
						Name:  "side-one",
						Image: "test/bad image name:1.0",
					},
				},
			},
		},
	}
	_, err := c.ValidateCreate()
	g.Expect(err).To(HaveOccurred())
}

func TestCoherenceCreateWithValidSidecarImageName(t *testing.T) {
	g := NewGomegaWithT(t)

	c := coh.Coherence{
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				SideCars: []corev1.Container{
					{
						Name:  "side-one",
						Image: "test/coherence:1.0",
					},
				},
			},
		},
	}
	_, err := c.ValidateCreate()
	g.Expect(err).NotTo(HaveOccurred())
}

func TestCoherenceCreateWithInvalidSidecarImageName(t *testing.T) {
	g := NewGomegaWithT(t)

	c := coh.Coherence{
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				SideCars: []corev1.Container{
					{
						Name:  "side-one",
						Image: "test/bad image name:1.0",
					},
				},
			},
		},
	}
	_, err := c.ValidateCreate()
	g.Expect(err).To(HaveOccurred())
}

func TestJobWithNoImageNames(t *testing.T) {
	g := NewGomegaWithT(t)

	c := coh.CoherenceJob{}
	_, err := c.ValidateCreate()
	g.Expect(err).NotTo(HaveOccurred())
}

func TestJobCreateWithInvalidImageName(t *testing.T) {
	g := NewGomegaWithT(t)

	c := coh.CoherenceJob{
		Spec: coh.CoherenceJobResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Image: stringPtr("test/bad image name:1.0"),
			},
		},
	}
	_, err := c.ValidateCreate()
	g.Expect(err).To(HaveOccurred())
}

func TestJobCreateWithValidImageDigest(t *testing.T) {
	g := NewGomegaWithT(t)

	c := coh.CoherenceJob{
		Spec: coh.CoherenceJobResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Image: stringPtr("ghcr.io@sha256:f8a592ee6d31c02feea037c269a87564ae666f91480d1d6be24ff9dd1675c7d0"),
			},
		},
	}
	_, err := c.ValidateCreate()
	g.Expect(err).NotTo(HaveOccurred())
}

func TestJobCreateWithInvalidImageDigest(t *testing.T) {
	g := NewGomegaWithT(t)

	c := coh.CoherenceJob{
		Spec: coh.CoherenceJobResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Image: stringPtr("test@sha256:1234"),
			},
		},
	}
	_, err := c.ValidateCreate()
	g.Expect(err).To(HaveOccurred())
}

func TestJobUpdateWithInvalidImageName(t *testing.T) {
	g := NewGomegaWithT(t)

	c := coh.CoherenceJob{
		Spec: coh.CoherenceJobResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Image: stringPtr("test/bad image name:1.0"),
			},
		},
	}
	_, err := c.ValidateUpdate(&c)
	g.Expect(err).To(HaveOccurred())
}
