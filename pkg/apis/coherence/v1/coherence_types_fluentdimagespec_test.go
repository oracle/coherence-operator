package v1_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	coherence "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	v1 "k8s.io/api/core/v1"
)

var _ = Describe("Testing FluentdImageSpec struct", func() {

	Context("Copying an FluentdImageSpec using DeepCopyWithDefaults", func() {
		var (
			original *coherence.FluentdImageSpec
			defaults *coherence.FluentdImageSpec
			clone    *coherence.FluentdImageSpec

			appOne = &coherence.FluentdApplicationSpec{ConfigFile: stringPointer("one.yaml")}
			appTwo = &coherence.FluentdApplicationSpec{ConfigFile: stringPointer("two.yaml")}

			always = v1.PullAlways
			never  = v1.PullNever

			specOne = coherence.ImageSpec{
				Image:           stringPointer("foo.1.0"),
				ImagePullPolicy: &always,
			}
			specTwo = coherence.ImageSpec{
				Image:           stringPointer("foo.2.0"),
				ImagePullPolicy: &never,
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
				original = &coherence.FluentdImageSpec{
					ImageSpec:   specOne,
					Application: appOne,
				}

				defaults = nil
			})

			It("should copy the original Image field", func() {
				Expect(*clone.Image).To(Equal(*original.Image))
			})

			It("should copy the original ImagePullPolicy field", func() {
				Expect(*clone.ImagePullPolicy).To(Equal(*original.ImagePullPolicy))
			})

			It("should copy the original Application field", func() {
				Expect(*clone.Application).To(Equal(*original.Application))
			})
		})

		When("original is nil", func() {
			BeforeEach(func() {
				defaults = &coherence.FluentdImageSpec{
					ImageSpec:   specTwo,
					Application: appTwo,
				}

				original = nil
			})

			It("should copy the defaults Image field", func() {
				Expect(*clone.Image).To(Equal(*defaults.Image))
			})

			It("should copy the defaults ImagePullPolicy field", func() {
				Expect(*clone.ImageSpec.ImagePullPolicy).To(Equal(*defaults.ImageSpec.ImagePullPolicy))
			})

			It("should copy the defaults Application field", func() {
				Expect(*clone.Application).To(Equal(*defaults.Application))
			})
		})

		When("all fields in the original are set", func() {
			BeforeEach(func() {
				original = &coherence.FluentdImageSpec{
					ImageSpec:   specOne,
					Application: appOne,
				}

				defaults = &coherence.FluentdImageSpec{
					ImageSpec:   specTwo,
					Application: appTwo,
				}
			})

			It("should copy the original Image field", func() {
				Expect(clone.ImageSpec).To(Equal(original.ImageSpec))
			})

			It("should copy the original ImagePullPolicy field", func() {
				Expect(*clone.ImagePullPolicy).To(Equal(*original.ImagePullPolicy))
			})

			It("should copy the original Application field", func() {
				Expect(*clone.Application).To(Equal(*original.Application))
			})
		})

		When("the original ImageSpec field is empty ImageSpec struct", func() {
			BeforeEach(func() {
				original = &coherence.FluentdImageSpec{
					ImageSpec:   coherence.ImageSpec{},
					Application: appOne,
				}

				defaults = &coherence.FluentdImageSpec{
					ImageSpec:   specTwo,
					Application: appTwo,
				}
			})

			It("should copy the defaults Image field", func() {
				Expect(*clone.Image).To(Equal(*defaults.Image))
			})

			It("should copy the defaults ImagePullPolicy field", func() {
				Expect(*clone.ImagePullPolicy).To(Equal(*defaults.ImagePullPolicy))
			})

			It("should copy the original Application field", func() {
				Expect(*clone.Application).To(Equal(*original.Application))
			})
		})

		When("the original ImageSpec.Image field is nil", func() {
			BeforeEach(func() {
				original = &coherence.FluentdImageSpec{
					ImageSpec: coherence.ImageSpec{
						Image:           nil,
						ImagePullPolicy: &always,
					},
					Application: appOne,
				}

				defaults = &coherence.FluentdImageSpec{
					ImageSpec:   specTwo,
					Application: appTwo,
				}
			})

			It("should copy the defaults Image field", func() {
				Expect(*clone.Image).To(Equal(*defaults.Image))
			})

			It("should copy the defaults ImagePullPolicy field", func() {
				Expect(*clone.ImagePullPolicy).To(Equal(*original.ImagePullPolicy))
			})

			It("should copy the original Application field", func() {
				Expect(*clone.Application).To(Equal(*original.Application))
			})
		})

		When("the original ImageSpec.ImagePullPolicy field is nil", func() {
			BeforeEach(func() {
				original = &coherence.FluentdImageSpec{
					ImageSpec: coherence.ImageSpec{
						Image:           stringPointer("foo:1.0"),
						ImagePullPolicy: nil,
					},
					Application: appOne,
				}

				defaults = &coherence.FluentdImageSpec{
					ImageSpec:   specTwo,
					Application: appTwo,
				}
			})

			It("should copy the defaults Image field", func() {
				Expect(*clone.Image).To(Equal(*original.Image))
			})

			It("should copy the defaults ImagePullPolicy field", func() {
				Expect(*clone.ImagePullPolicy).To(Equal(*defaults.ImagePullPolicy))
			})

			It("should copy the original Application field", func() {
				Expect(*clone.Application).To(Equal(*original.Application))
			})
		})

		When("the original Application field is nil", func() {
			BeforeEach(func() {
				original = &coherence.FluentdImageSpec{
					ImageSpec:   specOne,
					Application: nil,
				}

				defaults = &coherence.FluentdImageSpec{
					ImageSpec:   specTwo,
					Application: appTwo,
				}
			})

			It("should copy the original Image field", func() {
				Expect(*clone.Image).To(Equal(*original.Image))
			})

			It("should copy the original ImagePullPolicy field", func() {
				Expect(*clone.ImagePullPolicy).To(Equal(*original.ImagePullPolicy))
			})

			It("should copy the original Application field", func() {
				Expect(*clone.Application).To(Equal(*defaults.Application))
			})
		})
	})
})
