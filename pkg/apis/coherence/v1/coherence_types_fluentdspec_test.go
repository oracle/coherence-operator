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

var _ = Describe("Testing FluentdSpec struct", func() {

	Context("Copying an FluentdSpec using DeepCopyWithDefaults", func() {
		var (
			original *coherence.FluentdSpec
			defaults *coherence.FluentdSpec
			clone    *coherence.FluentdSpec

			always = v1.PullAlways
			never  = v1.PullNever

			specOne = coherence.ImageSpec{
				Image:           stringPtr("foo.1.0"),
				ImagePullPolicy: &always,
			}
			specTwo = coherence.ImageSpec{
				Image:           stringPtr("foo.2.0"),
				ImagePullPolicy: &never,
			}

			fdOne = &coherence.FluentdSpec{
				ImageSpec:  specOne,
				Enabled:    boolPtr(true),
				ConfigFile: stringPtr("one.yaml"),
				Tag:        stringPtr("tag-one"),
			}

			fdTwo = &coherence.FluentdSpec{
				ImageSpec:  specTwo,
				Enabled:    boolPtr(false),
				ConfigFile: stringPtr("two.yaml"),
				Tag:        stringPtr("tag-two"),
			}
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
				original = fdOne
				defaults = nil
			})

			It("should copy the original Image field", func() {
				Expect(*clone.Image).To(Equal(*original.Image))
			})

			It("should copy the original ImagePullPolicy field", func() {
				Expect(*clone.ImagePullPolicy).To(Equal(*original.ImagePullPolicy))
			})

			It("should copy the original Enabled field", func() {
				Expect(*clone.Enabled).To(Equal(*original.Enabled))
			})

			It("should copy the original ConfigFile field", func() {
				Expect(*clone.ConfigFile).To(Equal(*original.ConfigFile))
			})

			It("should copy the original Tag field", func() {
				Expect(*clone.Tag).To(Equal(*original.Tag))
			})
		})

		When("original is nil", func() {
			BeforeEach(func() {
				defaults = fdOne
				original = nil
			})

			It("should copy the defaults", func() {
				Expect(*clone).To(Equal(*defaults))
			})
		})

		When("all fields in the original are set", func() {
			BeforeEach(func() {
				original = fdOne
				defaults = fdTwo
			})

			It("should copy the original", func() {
				Expect(*clone).To(Equal(*original))
			})
		})

		When("the original Enabled field is nil", func() {
			BeforeEach(func() {
				original = fdOne.DeepCopy()
				original.Enabled = nil
				defaults = fdTwo
			})

			It("should copy the defaults Enabled field", func() {
				expected := original.DeepCopy()
				expected.Enabled = defaults.Enabled

				Expect(*clone).To(Equal(*expected))
			})
		})

		When("the original ConfigFile field is nil", func() {
			BeforeEach(func() {
				original = fdOne.DeepCopy()
				original.ConfigFile = nil
				defaults = fdTwo
			})

			It("should copy the defaults ConfigFile field", func() {
				expected := original.DeepCopy()
				expected.ConfigFile = defaults.ConfigFile

				Expect(*clone).To(Equal(*expected))
			})
		})

		When("the original Tag field is nil", func() {
			BeforeEach(func() {
				original = fdOne.DeepCopy()
				original.Tag = nil
				defaults = fdTwo
			})

			It("should copy the defaults Tag field", func() {
				expected := original.DeepCopy()
				expected.Tag = defaults.Tag

				Expect(*clone).To(Equal(*expected))
			})
		})

		When("the original ImageSpec field is empty ImageSpec struct", func() {
			BeforeEach(func() {
				original = fdOne.DeepCopy()
				original.ImageSpec = coherence.ImageSpec{}
				defaults = fdTwo
			})

			It("should copy the defaults Image field", func() {
				expected := original.DeepCopy()
				expected.ImageSpec = defaults.ImageSpec

				Expect(*clone).To(Equal(*expected))
			})
		})

		When("the original ImageSpec.Image field is nil", func() {
			BeforeEach(func() {
				original = fdOne.DeepCopy()
				original.Image = nil
				defaults = fdTwo
			})

			It("should copy the defaults Image field", func() {
				expected := original.DeepCopy()
				expected.Image = defaults.Image
				Expect(*clone).To(Equal(*expected))
			})
		})

		When("the original ImageSpec.ImagePullPolicy field is nil", func() {
			BeforeEach(func() {
				BeforeEach(func() {
					original = fdOne.DeepCopy()
					original.ConfigFile = nil
					defaults = fdTwo
				})

				It("should copy the defaults ImagePullPolicy field", func() {
					expected := original.DeepCopy()
					expected.ImagePullPolicy = defaults.ImagePullPolicy
					Expect(*clone).To(Equal(*expected))
				})
			})
		})
	})
})
