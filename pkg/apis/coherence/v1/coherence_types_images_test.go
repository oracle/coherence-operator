package v1_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	coherence "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
)

var _ = Describe("Testing Images struct", func() {

	Context("Copying an Images using DeepCopyWithDefaults", func() {
		var (
			specOne   = &coherence.ImageSpec{Image: stringPtr("foo.1.0")}
			specTwo   = &coherence.ImageSpec{Image: stringPtr("foo.2.0")}
			specThree = &coherence.ImageSpec{Image: stringPtr("foo.3.0")}
			specFour  = &coherence.ImageSpec{Image: stringPtr("foo.4.0")}

			userOne   = &coherence.UserArtifactsImageSpec{ConfigDir: stringPtr("/conf")}
			userTwo   = &coherence.UserArtifactsImageSpec{ConfigDir: stringPtr("/conf")}
			fluentOne = &coherence.FluentdImageSpec{ImageSpec: coherence.ImageSpec{Image: stringPtr("fluent:1.0")}}
			fluentTwo = &coherence.FluentdImageSpec{ImageSpec: coherence.ImageSpec{Image: stringPtr("fluent:2.0")}}

			imagesOne = &coherence.Images{
				Coherence:      specOne,
				CoherenceUtils: specTwo,
				UserArtifacts:  userOne,
				Fluentd:        fluentOne,
			}

			imagesTwo = &coherence.Images{
				Coherence:      specThree,
				CoherenceUtils: specFour,
				UserArtifacts:  userTwo,
				Fluentd:        fluentTwo,
			}

			original *coherence.Images
			defaults *coherence.Images
			clone    *coherence.Images
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
				original = imagesOne.DeepCopy()
				defaults = nil
			})

			It("the clone should equal the original", func() {
				Expect(clone).To(Equal(original))
			})
		})

		When("original is nil", func() {
			BeforeEach(func() {
				original = nil
				defaults = imagesTwo.DeepCopy()
			})

			It("the clone should equal the defaults", func() {
				Expect(clone).To(Equal(defaults))
			})
		})

		When("all fields in the original are set", func() {
			BeforeEach(func() {
				original = imagesOne.DeepCopy()
				defaults = imagesTwo.DeepCopy()
			})

			It("the clone should equal the original", func() {
				Expect(clone).To(Equal(original))
			})
		})

		When("the original Coherence field is nil", func() {
			BeforeEach(func() {
				original = imagesOne.DeepCopy()
				defaults = imagesTwo.DeepCopy()

				original.Coherence = nil
			})

			It("the clone should equal the original with the Coherence field from the defaults", func() {
				expected := original.DeepCopy()
				expected.Coherence = defaults.Coherence

				Expect(clone).To(Equal(expected))
			})
		})

		When("the original CoherenceUtils field is nil", func() {
			BeforeEach(func() {
				original = imagesOne.DeepCopy()
				defaults = imagesTwo.DeepCopy()

				original.CoherenceUtils = nil
			})

			It("the clone should equal the original with the CoherenceUtils field from the defaults", func() {
				expected := original.DeepCopy()
				expected.CoherenceUtils = defaults.CoherenceUtils

				Expect(clone).To(Equal(expected))
			})
		})

		When("the original UserArtifacts field is nil", func() {
			BeforeEach(func() {
				original = imagesOne.DeepCopy()
				defaults = imagesTwo.DeepCopy()

				original.UserArtifacts = nil
			})

			It("the clone should equal the original with the UserArtifacts field from the defaults", func() {
				expected := original.DeepCopy()
				expected.UserArtifacts = defaults.UserArtifacts

				Expect(clone).To(Equal(expected))
			})
		})

		When("the original Fluentd field is nil", func() {
			BeforeEach(func() {
				original = imagesOne.DeepCopy()
				defaults = imagesTwo.DeepCopy()

				original.Fluentd = nil
			})

			It("the clone should equal the original with the Fluentd field from the defaults", func() {
				expected := original.DeepCopy()
				expected.Fluentd = defaults.Fluentd

				Expect(clone).To(Equal(expected))
			})
		})
	})
})
