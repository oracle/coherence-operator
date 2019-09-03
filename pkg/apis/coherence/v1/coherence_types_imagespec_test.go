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
	v1 "k8s.io/api/core/v1"
)

var _ = Describe("Testing ImageSpec struct", func() {

	Context("Copying an ImageSpec using DeepCopyWithDefaults", func() {
		var original *coherence.ImageSpec
		var defaults *coherence.ImageSpec
		var clone *coherence.ImageSpec

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
				always := v1.PullAlways

				original = &coherence.ImageSpec{
					Image:           stringPtr("foo:1.0"),
					ImagePullPolicy: &always,
				}

				defaults = nil
			})

			It("should copy the original Image", func() {
				Expect(*clone.Image).To(Equal(*original.Image))
			})

			It("should copy the original ImagePullPolicy", func() {
				Expect(*clone.ImagePullPolicy).To(Equal(*original.ImagePullPolicy))
			})
		})

		When("defaults is empty", func() {
			BeforeEach(func() {
				always := v1.PullAlways

				original = &coherence.ImageSpec{
					Image:           stringPtr("foo:1.0"),
					ImagePullPolicy: &always,
				}

				defaults = &coherence.ImageSpec{}
			})

			It("should copy the original Image", func() {
				Expect(*clone.Image).To(Equal(*original.Image))
			})

			It("should copy the original ImagePullPolicy", func() {
				Expect(*clone.ImagePullPolicy).To(Equal(*original.ImagePullPolicy))
			})
		})

		When("original is nil", func() {
			BeforeEach(func() {
				always := v1.PullAlways

				defaults = &coherence.ImageSpec{
					Image:           stringPtr("foo:1.0"),
					ImagePullPolicy: &always,
				}

				original = nil
			})

			It("should copy the defaults Image", func() {
				Expect(*clone.Image).To(Equal(*defaults.Image))
			})

			It("should copy the defaults ImagePullPolicy", func() {
				Expect(*clone.ImagePullPolicy).To(Equal(*defaults.ImagePullPolicy))
			})
		})

		When("all fields in the original are set", func() {
			BeforeEach(func() {
				always := v1.PullAlways
				never := v1.PullNever

				original = &coherence.ImageSpec{
					Image:           stringPtr("foo:1.0"),
					ImagePullPolicy: &always,
				}

				defaults = &coherence.ImageSpec{
					Image:           stringPtr("foo:2.0"),
					ImagePullPolicy: &never,
				}
			})

			It("should copy the original Image", func() {
				Expect(*clone.Image).To(Equal(*original.Image))
			})

			It("should copy the original ImagePullPolicy", func() {
				Expect(*clone.ImagePullPolicy).To(Equal(*original.ImagePullPolicy))
			})
		})

		When("the original Image is nil", func() {
			BeforeEach(func() {
				always := v1.PullAlways
				never := v1.PullNever

				original = &coherence.ImageSpec{
					Image:           nil,
					ImagePullPolicy: &always,
				}

				defaults = &coherence.ImageSpec{
					Image:           stringPtr("foo:2.0"),
					ImagePullPolicy: &never,
				}
			})

			It("should copy the defaults Image", func() {
				Expect(*clone.Image).To(Equal(*defaults.Image))
			})

			It("should copy the original ImagePullPolicy", func() {
				Expect(*clone.ImagePullPolicy).To(Equal(*original.ImagePullPolicy))
			})
		})

		When("the original ImagePullPolicy is nil", func() {
			BeforeEach(func() {
				never := v1.PullNever

				original = &coherence.ImageSpec{
					Image:           stringPtr("foo:1.0"),
					ImagePullPolicy: nil,
				}

				defaults = &coherence.ImageSpec{
					Image:           stringPtr("foo:2.0"),
					ImagePullPolicy: &never,
				}
			})

			It("should copy the original Image", func() {
				Expect(*clone.Image).To(Equal(*original.Image))
			})

			It("should copy the defaults ImagePullPolicy", func() {
				Expect(*clone.ImagePullPolicy).To(Equal(*defaults.ImagePullPolicy))
			})
		})
	})
})
