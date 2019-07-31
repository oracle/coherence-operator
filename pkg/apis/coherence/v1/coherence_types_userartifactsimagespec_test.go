package v1_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	coherence "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
)

var _ = Describe("Testing UserArtifactsImageSpec struct", func() {

	Context("Copying an UserArtifactsImageSpec using DeepCopyWithDefaults", func() {
		var original *coherence.UserArtifactsImageSpec
		var defaults *coherence.UserArtifactsImageSpec
		var clone *coherence.UserArtifactsImageSpec

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
				original = &coherence.UserArtifactsImageSpec{
					LibDir:    stringPtr("/lib"),
					ConfigDir: stringPtr("/conf"),
				}

				defaults = nil
			})

			It("should copy the original LibDir", func() {
				Expect(*clone.LibDir).To(Equal(*original.LibDir))
			})

			It("should copy the original ConfigDir", func() {
				Expect(*clone.ConfigDir).To(Equal(*original.ConfigDir))
			})
		})

		When("original is nil", func() {
			BeforeEach(func() {
				defaults = &coherence.UserArtifactsImageSpec{
					LibDir:    stringPtr("/lib"),
					ConfigDir: stringPtr("/conf"),
				}

				original = nil
			})

			It("should copy the defaults LibDir", func() {
				Expect(*clone.LibDir).To(Equal(*defaults.LibDir))
			})

			It("should copy the defaults ConfigDir", func() {
				Expect(*clone.ConfigDir).To(Equal(*defaults.ConfigDir))
			})
		})

		When("all fields in the original are set", func() {
			BeforeEach(func() {
				original = &coherence.UserArtifactsImageSpec{
					LibDir:    stringPtr("/lib"),
					ConfigDir: stringPtr("/conf"),
				}

				defaults = &coherence.UserArtifactsImageSpec{
					LibDir:    stringPtr("/libDefault"),
					ConfigDir: stringPtr("/confDefault"),
				}
			})

			It("should copy the original LibDir", func() {
				Expect(*clone.LibDir).To(Equal(*original.LibDir))
			})

			It("should copy the original ConfigDir", func() {
				Expect(*clone.ConfigDir).To(Equal(*original.ConfigDir))
			})
		})

		When("the original Image is nil", func() {
			BeforeEach(func() {
				original = &coherence.UserArtifactsImageSpec{
					LibDir:    nil,
					ConfigDir: stringPtr("/conf"),
				}

				defaults = &coherence.UserArtifactsImageSpec{
					LibDir:    stringPtr("/libDefault"),
					ConfigDir: stringPtr("/confDefault"),
				}
			})

			It("should copy the defaults LibDir", func() {
				Expect(*clone.LibDir).To(Equal(*defaults.LibDir))
			})

			It("should copy the original ConfigDir", func() {
				Expect(*clone.ConfigDir).To(Equal(*original.ConfigDir))
			})
		})

		When("the original ConfigDir is nil", func() {
			BeforeEach(func() {
				original = &coherence.UserArtifactsImageSpec{
					LibDir:    stringPtr("/lib"),
					ConfigDir: nil,
				}

				defaults = &coherence.UserArtifactsImageSpec{
					LibDir:    stringPtr("/libDefault"),
					ConfigDir: stringPtr("/confDefault"),
				}
			})

			It("should copy the original LibDir", func() {
				Expect(*clone.LibDir).To(Equal(*original.LibDir))
			})

			It("should copy the defaults ConfigDir", func() {
				Expect(*clone.ConfigDir).To(Equal(*defaults.ConfigDir))
			})
		})
	})
})
