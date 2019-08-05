package v1_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	coherence "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
)

var _ = Describe("Testing LoggingSpec struct", func() {

	Context("Copying a LoggingSpec using DeepCopyWithDefaults", func() {
		var original *coherence.LoggingSpec
		var defaults *coherence.LoggingSpec
		var clone *coherence.LoggingSpec

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
				original = &coherence.LoggingSpec{
					Level:         int32Ptr(9),
					ConfigFile:    stringPtr("logging.properties"),
					ConfigMapName: stringPtr("loggingMap"),
				}

				defaults = nil
			})

			It("should copy the original Level", func() {
				Expect(*clone.Level).To(Equal(*original.Level))
			})

			It("should copy the original ConfigFile", func() {
				Expect(*clone.ConfigFile).To(Equal(*original.ConfigFile))
			})

			It("should copy the original ConfigMapName", func() {
				Expect(*clone.ConfigMapName).To(Equal(*original.ConfigMapName))
			})
		})

		When("original is nil", func() {
			BeforeEach(func() {
				defaults = &coherence.LoggingSpec{
					Level:         int32Ptr(9),
					ConfigFile:    stringPtr("logging.properties"),
					ConfigMapName: stringPtr("loggingMap"),
				}

				original = nil
			})

			It("should copy the original Level", func() {
				Expect(*clone.Level).To(Equal(*defaults.Level))
			})

			It("should copy the original ConfigFile", func() {
				Expect(*clone.ConfigFile).To(Equal(*defaults.ConfigFile))
			})

			It("should copy the original ConfigMapName", func() {
				Expect(*clone.ConfigMapName).To(Equal(*defaults.ConfigMapName))
			})
		})

		When("all original fields are set", func() {
			BeforeEach(func() {
				original = &coherence.LoggingSpec{
					Level:         int32Ptr(9),
					ConfigFile:    stringPtr("logging.properties"),
					ConfigMapName: stringPtr("loggingMap"),
				}

				defaults = &coherence.LoggingSpec{
					Level:         int32Ptr(7),
					ConfigFile:    stringPtr("logging2.properties"),
					ConfigMapName: stringPtr("loggingMap2"),
				}
			})

			It("should copy the original Level", func() {
				Expect(*clone.Level).To(Equal(*original.Level))
			})

			It("should copy the original ConfigFile", func() {
				Expect(*clone.ConfigFile).To(Equal(*original.ConfigFile))
			})

			It("should copy the original ConfigMapName", func() {
				Expect(*clone.ConfigMapName).To(Equal(*original.ConfigMapName))
			})
		})

		When("original Level is nil", func() {
			BeforeEach(func() {
				original = &coherence.LoggingSpec{
					Level:         nil,
					ConfigFile:    stringPtr("logging.properties"),
					ConfigMapName: stringPtr("loggingMap"),
				}

				defaults = &coherence.LoggingSpec{
					Level:         int32Ptr(7),
					ConfigFile:    stringPtr("logging2.properties"),
					ConfigMapName: stringPtr("loggingMap2"),
				}
			})

			It("should copy the original Level", func() {
				Expect(*clone.Level).To(Equal(*defaults.Level))
			})

			It("should copy the original ConfigFile", func() {
				Expect(*clone.ConfigFile).To(Equal(*original.ConfigFile))
			})

			It("should copy the original ConfigMapName", func() {
				Expect(*clone.ConfigMapName).To(Equal(*original.ConfigMapName))
			})
		})

		When("original ConfigFile is nil", func() {
			BeforeEach(func() {
				original = &coherence.LoggingSpec{
					Level:         int32Ptr(9),
					ConfigFile:    nil,
					ConfigMapName: stringPtr("loggingMap"),
				}

				defaults = &coherence.LoggingSpec{
					Level:         int32Ptr(7),
					ConfigFile:    stringPtr("logging2.properties"),
					ConfigMapName: stringPtr("loggingMap2"),
				}
			})

			It("should copy the original Level", func() {
				Expect(*clone.Level).To(Equal(*original.Level))
			})

			It("should copy the original ConfigFile", func() {
				Expect(*clone.ConfigFile).To(Equal(*defaults.ConfigFile))
			})

			It("should copy the original ConfigMapName", func() {
				Expect(*clone.ConfigMapName).To(Equal(*original.ConfigMapName))
			})
		})

		When("original ConfigMapName is nil", func() {
			BeforeEach(func() {
				original = &coherence.LoggingSpec{
					Level:         int32Ptr(9),
					ConfigFile:    stringPtr("logging.properties"),
					ConfigMapName: nil,
				}

				defaults = &coherence.LoggingSpec{
					Level:         int32Ptr(7),
					ConfigFile:    stringPtr("logging2.properties"),
					ConfigMapName: stringPtr("loggingMap2"),
				}
			})

			It("should copy the original Level", func() {
				Expect(*clone.Level).To(Equal(*original.Level))
			})

			It("should copy the original ConfigFile", func() {
				Expect(*clone.ConfigFile).To(Equal(*original.ConfigFile))
			})

			It("should copy the original ConfigMapName", func() {
				Expect(*clone.ConfigMapName).To(Equal(*defaults.ConfigMapName))
			})
		})
	})
})
