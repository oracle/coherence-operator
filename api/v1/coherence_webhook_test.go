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
