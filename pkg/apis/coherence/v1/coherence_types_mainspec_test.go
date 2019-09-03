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

var _ = Describe("Testing MainSpec struct", func() {

	Context("Copying a MainSpec using DeepCopyWithDefaults", func() {
		var original *coherence.MainSpec
		var defaults *coherence.MainSpec
		var clone *coherence.MainSpec

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
				original = &coherence.MainSpec{
					Class:     stringPtr("com.tangosol.net.DefaultCacheServer"),
					Arguments: stringPtr("-Dcoherence.localhost=192.168.0.301"),
				}

				defaults = nil
			})

			It("should copy the original Class", func() {
				Expect(*clone.Class).To(Equal(*original.Class))
			})

			It("should copy the original Arguements", func() {
				Expect(*clone.Arguments).To(Equal(*original.Arguments))
			})
		})

		When("original is nil", func() {
			BeforeEach(func() {
				defaults = &coherence.MainSpec{
					Class:     stringPtr("com.tangosol.net.DefaultCacheServer"),
					Arguments: stringPtr("-Dcoherence.localhost=192.168.0.301"),
				}

				original = nil
			})

			It("should copy the defaults Class", func() {
				Expect(*clone.Class).To(Equal(*defaults.Class))
			})

			It("should copy the defaults Arguements", func() {
				Expect(*clone.Arguments).To(Equal(*defaults.Arguments))
			})
		})

		When("all original fields are set", func() {
			BeforeEach(func() {
				original = &coherence.MainSpec{
					Class:     stringPtr("com.tangosol.net.DefaultCacheServer"),
					Arguments: stringPtr("-Dcoherence.localhost=192.168.0.301"),
				}

				defaults = &coherence.MainSpec{
					Class:     stringPtr("com.tangosol.net.DefaultCacheServer2"),
					Arguments: stringPtr("-Dcoherence.localhost=192.168.0.302"),
				}
			})

			It("should copy the original Class", func() {
				Expect(*clone.Class).To(Equal(*original.Class))
			})

			It("should copy the original Arguements", func() {
				Expect(*clone.Arguments).To(Equal(*original.Arguments))
			})
		})

		When("original Class is nil", func() {
			BeforeEach(func() {
				original = &coherence.MainSpec{
					Class:     nil,
					Arguments: stringPtr("-Dcoherence.localhost=192.168.0.301"),
				}

				defaults = &coherence.MainSpec{
					Class:     stringPtr("com.tangosol.net.DefaultCacheServer2"),
					Arguments: stringPtr("-Dcoherence.localhost=192.168.0.302"),
				}
			})

			It("should copy the defaults Class", func() {
				Expect(*clone.Class).To(Equal(*defaults.Class))
			})

			It("should copy the original Arguements", func() {
				Expect(*clone.Arguments).To(Equal(*original.Arguments))
			})
		})

		When("original Arguments is nil", func() {
			BeforeEach(func() {
				original = &coherence.MainSpec{
					Class:     stringPtr("com.tangosol.net.DefaultCacheServer"),
					Arguments: nil,
				}

				defaults = &coherence.MainSpec{
					Class:     stringPtr("com.tangosol.net.DefaultCacheServer2"),
					Arguments: stringPtr("-Dcoherence.localhost=192.168.0.302"),
				}
			})

			It("should copy the original Class", func() {
				Expect(*clone.Class).To(Equal(*original.Class))
			})

			It("should copy the defaults Arguements", func() {
				Expect(*clone.Arguments).To(Equal(*defaults.Arguments))
			})
		})
	})
})
