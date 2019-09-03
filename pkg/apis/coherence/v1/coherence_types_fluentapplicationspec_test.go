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
)

var _ = Describe("Testing FluentdApplicationSpec struct", func() {

	Context("Copying an FluentdApplicationSpec using DeepCopyWithDefaults", func() {
		var original *coherence.FluentdApplicationSpec
		var defaults *coherence.FluentdApplicationSpec
		var clone *coherence.FluentdApplicationSpec

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
				original = &coherence.FluentdApplicationSpec{
					ConfigFile: stringPtr("config.yaml"),
					Tag:        stringPtr("abc-123"),
				}

				defaults = nil
			})

			It("should copy the original ConfigFile", func() {
				Expect(*clone.ConfigFile).To(Equal(*original.ConfigFile))
			})

			It("should copy the original Tag", func() {
				Expect(*clone.Tag).To(Equal(*original.Tag))
			})
		})

		When("original is nil", func() {
			BeforeEach(func() {
				defaults = &coherence.FluentdApplicationSpec{
					ConfigFile: stringPtr("config.yaml"),
					Tag:        stringPtr("abc-123"),
				}

				original = nil
			})

			It("should copy the defaults ConfigFile", func() {
				Expect(*clone.ConfigFile).To(Equal(*defaults.ConfigFile))
			})

			It("should copy the defaults Tag", func() {
				Expect(*clone.Tag).To(Equal(*defaults.Tag))
			})
		})

		When("all fields in the original are set", func() {
			BeforeEach(func() {
				original = &coherence.FluentdApplicationSpec{
					ConfigFile: stringPtr("config.yaml"),
					Tag:        stringPtr("abc-123"),
				}

				defaults = &coherence.FluentdApplicationSpec{
					ConfigFile: stringPtr("default-config.yaml"),
					Tag:        stringPtr("def-456"),
				}
			})

			It("should copy the original ConfigFile", func() {
				Expect(*clone.ConfigFile).To(Equal(*original.ConfigFile))
			})

			It("should copy the original Tag", func() {
				Expect(*clone.Tag).To(Equal(*original.Tag))
			})
		})

		When("the original Image is nil", func() {
			BeforeEach(func() {
				original = &coherence.FluentdApplicationSpec{
					ConfigFile: nil,
					Tag:        stringPtr("abc-123"),
				}

				defaults = &coherence.FluentdApplicationSpec{
					ConfigFile: stringPtr("default-config.yaml"),
					Tag:        stringPtr("def-456"),
				}
			})

			It("should copy the defaults ConfigFile", func() {
				Expect(*clone.ConfigFile).To(Equal(*defaults.ConfigFile))
			})

			It("should copy the original Tag", func() {
				Expect(*clone.Tag).To(Equal(*original.Tag))
			})
		})

		When("the original Tag is nil", func() {
			BeforeEach(func() {
				original = &coherence.FluentdApplicationSpec{
					ConfigFile: stringPtr("config.yaml"),
					Tag:        nil,
				}

				defaults = &coherence.FluentdApplicationSpec{
					ConfigFile: stringPtr("default-config.yaml"),
					Tag:        stringPtr("def-456"),
				}
			})

			It("should copy the original ConfigFile", func() {
				Expect(*clone.ConfigFile).To(Equal(*original.ConfigFile))
			})

			It("should copy the defaults Tag", func() {
				Expect(*clone.Tag).To(Equal(*defaults.Tag))
			})
		})
	})
})
