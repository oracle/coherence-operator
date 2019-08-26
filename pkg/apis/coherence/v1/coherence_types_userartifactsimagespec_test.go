package v1_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	coherence "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	corev1 "k8s.io/api/core/v1"
)

var _ = Describe("Testing UserArtifactsImageSpec struct", func() {

	Context("Copying an UserArtifactsImageSpec using DeepCopyWithDefaults", func() {
		var original *coherence.UserArtifactsImageSpec
		var defaults *coherence.UserArtifactsImageSpec
		var clone *coherence.UserArtifactsImageSpec
		var always = corev1.PullAlways
		var never = corev1.PullNever

		var userArtifactsOne = &coherence.UserArtifactsImageSpec{
			ImageSpec: coherence.ImageSpec{
				Image:           stringPtr("one:1.0"),
				ImagePullPolicy: &always,
			},
			LibDir:    stringPtr("/one/lib"),
			ConfigDir: stringPtr("/one/cfg"),
		}

		var userArtifactsTwo = &coherence.UserArtifactsImageSpec{
			ImageSpec: coherence.ImageSpec{
				Image:           stringPtr("two:1.0"),
				ImagePullPolicy: &never,
			},
			LibDir:    stringPtr("/two/lib"),
			ConfigDir: stringPtr("/two/cfg"),
		}

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
				original = userArtifactsOne
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
				original = nil
				defaults = userArtifactsTwo
			})

			It("copy should equal defaults", func() {
				Expect(clone).To(Equal(defaults))
			})
		})

		When("all fields in the original are set", func() {
			BeforeEach(func() {
				original = userArtifactsOne
				defaults = userArtifactsTwo
			})

			It("copy should equal original", func() {
				Expect(clone).To(Equal(original))
			})
		})

		When("the original Image is nil", func() {
			BeforeEach(func() {
				original = userArtifactsOne.DeepCopy()
				original.Image = nil
				defaults = userArtifactsTwo
			})

			It("copy should equal original with the defaults Image", func() {
				expected := original.DeepCopy()
				expected.Image = defaults.Image
				Expect(clone).To(Equal(expected))
			})
		})

		When("the original ImagePullPolicy is nil", func() {
			BeforeEach(func() {
				original = userArtifactsOne.DeepCopy()
				original.ImagePullPolicy = nil
				defaults = userArtifactsTwo
			})

			It("copy should equal original with the defaults ImagePullPolicy", func() {
				expected := original.DeepCopy()
				expected.ImagePullPolicy = defaults.ImagePullPolicy
				Expect(clone).To(Equal(expected))
			})
		})

		When("the original LibDir is nil", func() {
			BeforeEach(func() {
				original = userArtifactsOne.DeepCopy()
				original.LibDir = nil
				defaults = userArtifactsTwo
			})

			It("copy should equal original with the defaults LibDir", func() {
				expected := original.DeepCopy()
				expected.LibDir = defaults.LibDir
				Expect(clone).To(Equal(expected))
			})
		})

		When("the original ConfigDir is nil", func() {
			BeforeEach(func() {
				original = userArtifactsOne.DeepCopy()
				original.ConfigDir = nil
				defaults = userArtifactsTwo
			})

			It("copy should equal original with the defaults ConfigDir", func() {
				expected := original.DeepCopy()
				expected.ConfigDir = defaults.ConfigDir
				Expect(clone).To(Equal(expected))
			})
		})
	})
})
