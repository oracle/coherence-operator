package v1_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	coherence "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
)

var _ = Describe("Testing ReadinessProbeSpec struct", func() {

	Context("Copying a ReadinessProbeSpec using DeepCopyWithDefaults", func() {
		var original *coherence.ReadinessProbeSpec
		var defaults *coherence.ReadinessProbeSpec
		var clone *coherence.ReadinessProbeSpec

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
				original = &coherence.ReadinessProbeSpec{
					InitialDelaySeconds: int32Pointer(10),
					TimeoutSeconds:      int32Pointer(20),
					PeriodSeconds:       int32Pointer(30),
					SuccessThreshold:    int32Pointer(40),
					FailureThreshold:    int32Pointer(50),
				}

				defaults = nil
			})

			It("should copy the original InitialDelaySeconds", func() {
				Expect(*clone.InitialDelaySeconds).To(Equal(*original.InitialDelaySeconds))
			})

			It("should copy the original TimeoutSeconds", func() {
				Expect(*clone.TimeoutSeconds).To(Equal(*original.TimeoutSeconds))
			})

			It("should copy the original PeriodSeconds", func() {
				Expect(*clone.PeriodSeconds).To(Equal(*original.PeriodSeconds))
			})

			It("should copy the original SuccessThreshold", func() {
				Expect(*clone.SuccessThreshold).To(Equal(*original.SuccessThreshold))
			})

			It("should copy the original FailureThreshold", func() {
				Expect(*clone.FailureThreshold).To(Equal(*original.FailureThreshold))
			})
		})

		When("original is nil", func() {
			BeforeEach(func() {
				defaults = &coherence.ReadinessProbeSpec{
					InitialDelaySeconds: int32Pointer(10),
					TimeoutSeconds:      int32Pointer(20),
					PeriodSeconds:       int32Pointer(30),
					SuccessThreshold:    int32Pointer(40),
					FailureThreshold:    int32Pointer(50),
				}

				original = nil
			})

			It("should copy the defaults InitialDelaySeconds", func() {
				Expect(*clone.InitialDelaySeconds).To(Equal(*defaults.InitialDelaySeconds))
			})

			It("should copy the defaults TimeoutSeconds", func() {
				Expect(*clone.TimeoutSeconds).To(Equal(*defaults.TimeoutSeconds))
			})

			It("should copy the defaults PeriodSeconds", func() {
				Expect(*clone.PeriodSeconds).To(Equal(*defaults.PeriodSeconds))
			})

			It("should copy the defaults SuccessThreshold", func() {
				Expect(*clone.SuccessThreshold).To(Equal(*defaults.SuccessThreshold))
			})

			It("should copy the defaults FailureThreshold", func() {
				Expect(*clone.FailureThreshold).To(Equal(*defaults.FailureThreshold))
			})
		})

		When("all original fields are set", func() {
			BeforeEach(func() {
				original = &coherence.ReadinessProbeSpec{
					InitialDelaySeconds: int32Pointer(10),
					TimeoutSeconds:      int32Pointer(20),
					PeriodSeconds:       int32Pointer(30),
					SuccessThreshold:    int32Pointer(40),
					FailureThreshold:    int32Pointer(50),
				}

				defaults = &coherence.ReadinessProbeSpec{
					InitialDelaySeconds: int32Pointer(100),
					TimeoutSeconds:      int32Pointer(200),
					PeriodSeconds:       int32Pointer(300),
					SuccessThreshold:    int32Pointer(400),
					FailureThreshold:    int32Pointer(500),
				}
			})

			It("should copy the original InitialDelaySeconds", func() {
				Expect(*clone.InitialDelaySeconds).To(Equal(*original.InitialDelaySeconds))
			})

			It("should copy the defaults TimeoutSeconds", func() {
				Expect(*clone.TimeoutSeconds).To(Equal(*original.TimeoutSeconds))
			})

			It("should copy the original PeriodSeconds", func() {
				Expect(*clone.PeriodSeconds).To(Equal(*original.PeriodSeconds))
			})

			It("should copy the original SuccessThreshold", func() {
				Expect(*clone.SuccessThreshold).To(Equal(*original.SuccessThreshold))
			})

			It("should copy the original FailureThreshold", func() {
				Expect(*clone.FailureThreshold).To(Equal(*original.FailureThreshold))
			})
		})

		When("original InitialDelaySeconds is nil", func() {
			BeforeEach(func() {
				original = &coherence.ReadinessProbeSpec{
					InitialDelaySeconds: nil,
					TimeoutSeconds:      int32Pointer(20),
					PeriodSeconds:       int32Pointer(30),
					SuccessThreshold:    int32Pointer(40),
					FailureThreshold:    int32Pointer(50),
				}

				defaults = &coherence.ReadinessProbeSpec{
					InitialDelaySeconds: int32Pointer(100),
					TimeoutSeconds:      int32Pointer(200),
					PeriodSeconds:       int32Pointer(300),
					SuccessThreshold:    int32Pointer(400),
					FailureThreshold:    int32Pointer(500),
				}
			})

			It("should copy the defaults InitialDelaySeconds", func() {
				Expect(*clone.InitialDelaySeconds).To(Equal(*defaults.InitialDelaySeconds))
			})

			It("should copy the original TimeoutSeconds", func() {
				Expect(*clone.TimeoutSeconds).To(Equal(*original.TimeoutSeconds))
			})

			It("should copy the original PeriodSeconds", func() {
				Expect(*clone.PeriodSeconds).To(Equal(*original.PeriodSeconds))
			})

			It("should copy the original SuccessThreshold", func() {
				Expect(*clone.SuccessThreshold).To(Equal(*original.SuccessThreshold))
			})

			It("should copy the original FailureThreshold", func() {
				Expect(*clone.FailureThreshold).To(Equal(*original.FailureThreshold))
			})
		})

		When("original TimeoutSeconds is nil", func() {
			BeforeEach(func() {
				original = &coherence.ReadinessProbeSpec{
					InitialDelaySeconds: int32Pointer(10),
					TimeoutSeconds:      nil,
					PeriodSeconds:       int32Pointer(30),
					SuccessThreshold:    int32Pointer(40),
					FailureThreshold:    int32Pointer(50),
				}

				defaults = &coherence.ReadinessProbeSpec{
					InitialDelaySeconds: int32Pointer(100),
					TimeoutSeconds:      int32Pointer(200),
					PeriodSeconds:       int32Pointer(300),
					SuccessThreshold:    int32Pointer(400),
					FailureThreshold:    int32Pointer(500),
				}
			})

			It("should copy the original InitialDelaySeconds", func() {
				Expect(*clone.InitialDelaySeconds).To(Equal(*original.InitialDelaySeconds))
			})

			It("should copy the defaults TimeoutSeconds", func() {
				Expect(*clone.TimeoutSeconds).To(Equal(*defaults.TimeoutSeconds))
			})

			It("should copy the original PeriodSeconds", func() {
				Expect(*clone.PeriodSeconds).To(Equal(*original.PeriodSeconds))
			})

			It("should copy the original SuccessThreshold", func() {
				Expect(*clone.SuccessThreshold).To(Equal(*original.SuccessThreshold))
			})

			It("should copy the original FailureThreshold", func() {
				Expect(*clone.FailureThreshold).To(Equal(*original.FailureThreshold))
			})
		})

		When("original PeriodSeconds is nil", func() {
			BeforeEach(func() {
				original = &coherence.ReadinessProbeSpec{
					InitialDelaySeconds: int32Pointer(10),
					TimeoutSeconds:      int32Pointer(20),
					PeriodSeconds:       nil,
					SuccessThreshold:    int32Pointer(40),
					FailureThreshold:    int32Pointer(50),
				}

				defaults = &coherence.ReadinessProbeSpec{
					InitialDelaySeconds: int32Pointer(100),
					TimeoutSeconds:      int32Pointer(200),
					PeriodSeconds:       int32Pointer(300),
					SuccessThreshold:    int32Pointer(400),
					FailureThreshold:    int32Pointer(500),
				}
			})

			It("should copy the original InitialDelaySeconds", func() {
				Expect(*clone.InitialDelaySeconds).To(Equal(*original.InitialDelaySeconds))
			})

			It("should copy the original TimeoutSeconds", func() {
				Expect(*clone.TimeoutSeconds).To(Equal(*original.TimeoutSeconds))
			})

			It("should copy the defaults PeriodSeconds", func() {
				Expect(*clone.PeriodSeconds).To(Equal(*defaults.PeriodSeconds))
			})

			It("should copy the original SuccessThreshold", func() {
				Expect(*clone.SuccessThreshold).To(Equal(*original.SuccessThreshold))
			})

			It("should copy the original FailureThreshold", func() {
				Expect(*clone.FailureThreshold).To(Equal(*original.FailureThreshold))
			})
		})

		When("original SuccessThreshold is nil", func() {
			BeforeEach(func() {
				original = &coherence.ReadinessProbeSpec{
					InitialDelaySeconds: int32Pointer(10),
					TimeoutSeconds:      int32Pointer(20),
					PeriodSeconds:       int32Pointer(30),
					SuccessThreshold:    nil,
					FailureThreshold:    int32Pointer(50),
				}

				defaults = &coherence.ReadinessProbeSpec{
					InitialDelaySeconds: int32Pointer(100),
					TimeoutSeconds:      int32Pointer(200),
					PeriodSeconds:       int32Pointer(300),
					SuccessThreshold:    int32Pointer(400),
					FailureThreshold:    int32Pointer(500),
				}
			})

			It("should copy the original InitialDelaySeconds", func() {
				Expect(*clone.InitialDelaySeconds).To(Equal(*original.InitialDelaySeconds))
			})

			It("should copy the original TimeoutSeconds", func() {
				Expect(*clone.TimeoutSeconds).To(Equal(*original.TimeoutSeconds))
			})

			It("should copy the original PeriodSeconds", func() {
				Expect(*clone.PeriodSeconds).To(Equal(*original.PeriodSeconds))
			})

			It("should copy the defaults SuccessThreshold", func() {
				Expect(*clone.SuccessThreshold).To(Equal(*defaults.SuccessThreshold))
			})

			It("should copy the original FailureThreshold", func() {
				Expect(*clone.FailureThreshold).To(Equal(*original.FailureThreshold))
			})
		})

		When("original FailureThreshold is nil", func() {
			BeforeEach(func() {
				original = &coherence.ReadinessProbeSpec{
					InitialDelaySeconds: int32Pointer(10),
					TimeoutSeconds:      int32Pointer(20),
					PeriodSeconds:       int32Pointer(30),
					SuccessThreshold:    int32Pointer(40),
					FailureThreshold:    nil,
				}

				defaults = &coherence.ReadinessProbeSpec{
					InitialDelaySeconds: int32Pointer(100),
					TimeoutSeconds:      int32Pointer(200),
					PeriodSeconds:       int32Pointer(300),
					SuccessThreshold:    int32Pointer(400),
					FailureThreshold:    int32Pointer(500),
				}
			})

			It("should copy the original InitialDelaySeconds", func() {
				Expect(*clone.InitialDelaySeconds).To(Equal(*original.InitialDelaySeconds))
			})

			It("should copy the defaults TimeoutSeconds", func() {
				Expect(*clone.TimeoutSeconds).To(Equal(*original.TimeoutSeconds))
			})

			It("should copy the original PeriodSeconds", func() {
				Expect(*clone.PeriodSeconds).To(Equal(*original.PeriodSeconds))
			})

			It("should copy the original SuccessThreshold", func() {
				Expect(*clone.SuccessThreshold).To(Equal(*original.SuccessThreshold))
			})

			It("should copy the defaults FailureThreshold", func() {
				Expect(*clone.FailureThreshold).To(Equal(*defaults.FailureThreshold))
			})
		})
	})

})
