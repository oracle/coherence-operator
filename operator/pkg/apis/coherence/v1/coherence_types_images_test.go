package v1_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	coherence "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
)

var _ = Describe("Testing Images struct", func() {

	Context("Copying an Images using DeepCopyWithDefaults", func() {
		var (
			original *coherence.Images
			defaults *coherence.Images
			clone    *coherence.Images

			specOne   = &coherence.ImageSpec{Image: stringPointer("foo.1.0")}
			specTwo   = &coherence.ImageSpec{Image: stringPointer("foo.2.0")}
			specThree = &coherence.ImageSpec{Image: stringPointer("foo.3.0")}
			specFour  = &coherence.ImageSpec{Image: stringPointer("foo.4.0")}
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
				original = &coherence.Images{
					Coherence:      specOne,
					CoherenceUtils: specTwo,
				}

				defaults = nil
			})

			It("should copy the original Coherence field", func() {
				Expect(*clone.Coherence).To(Equal(*original.Coherence))
			})

			It("should copy the original CoherenceUtils field", func() {
				Expect(*clone.CoherenceUtils).To(Equal(*original.CoherenceUtils))
			})
		})

		When("original is nil", func() {
			BeforeEach(func() {
				defaults = &coherence.Images{
					Coherence:      specThree,
					CoherenceUtils: specFour,
				}

				original = nil
			})

			It("should copy the defaults Coherence field", func() {
				Expect(*clone.Coherence).To(Equal(*defaults.Coherence))
			})

			It("should copy the defaults CoherenceUtils field", func() {
				Expect(*clone.CoherenceUtils).To(Equal(*defaults.CoherenceUtils))
			})
		})

		When("all fields in the original are set", func() {
			BeforeEach(func() {
				original = &coherence.Images{
					Coherence:      specOne,
					CoherenceUtils: specTwo,
				}

				defaults = &coherence.Images{
					Coherence:      specThree,
					CoherenceUtils: specFour,
				}
			})

			It("should copy the original Coherence field", func() {
				Expect(*clone.Coherence).To(Equal(*original.Coherence))
			})

			It("should copy the original CoherenceUtils field", func() {
				Expect(*clone.CoherenceUtils).To(Equal(*original.CoherenceUtils))
			})
		})

		When("the original Coherence field is nil", func() {
			BeforeEach(func() {
				original = &coherence.Images{
					Coherence:      nil,
					CoherenceUtils: specTwo,
				}

				defaults = &coherence.Images{
					Coherence:      specThree,
					CoherenceUtils: specFour,
				}
			})

			It("should copy the defaults Coherence field", func() {
				Expect(*clone.Coherence).To(Equal(*defaults.Coherence))
			})

			It("should copy the original CoherenceUtils field", func() {
				Expect(*clone.CoherenceUtils).To(Equal(*original.CoherenceUtils))
			})
		})

		When("the original CoherenceUtils field is nil", func() {
			BeforeEach(func() {
				original = &coherence.Images{
					Coherence:      specOne,
					CoherenceUtils: nil,
				}

				defaults = &coherence.Images{
					Coherence:      specThree,
					CoherenceUtils: specFour,
				}
			})

			It("should copy the original Coherence field", func() {
				Expect(*clone.Coherence).To(Equal(*original.Coherence))
			})

			It("should copy the defaults CoherenceUtils field", func() {
				Expect(*clone.CoherenceUtils).To(Equal(*defaults.CoherenceUtils))
			})
		})
	})
})
