package v1_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	coherence "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
)

var _ = Describe("Testing ServiceSpec struct", func() {

	Context("Copying a ServiceSpec using DeepCopyWithDefaults", func() {
		var original *coherence.ServiceSpec
		var defaults *coherence.ServiceSpec
		var clone *coherence.ServiceSpec
		var expected *coherence.ServiceSpec

		NewServiceSpecOne := func() *coherence.ServiceSpec {
			return &coherence.ServiceSpec{
				Enabled:        boolPtr(true),
				Type:           stringPtr("LoadBalancerIP"),
				Domain:         stringPtr("cluster.local"),
				LoadBalancerIP: stringPtr("10.10.10.20"),
				Annotations:    map[string]string{ "foo": "1"},
				ExternalPort:   int32Ptr(9099),
			}
		}

		NewServiceSpecTwo := func() *coherence.ServiceSpec {
			return &coherence.ServiceSpec{
				Enabled:        boolPtr(true),
				Type:           stringPtr("ClusterIP"),
				Domain:         stringPtr("cluster.local2"),
				LoadBalancerIP: stringPtr("10.10.10.21"),
				Annotations:    map[string]string{ "foo": "2"},
				ExternalPort:   int32Ptr(9098),
			}
		}

		ValidateResult := func() {
			It("should have correct Enabled", func() {
				Expect(*clone.Enabled).To(Equal(*expected.Enabled))
			})

			It("should have correct Type", func() {
				Expect(*clone.Type).To(Equal(*expected.Type))
			})

			It("should have correct Domain", func() {
				Expect(*clone.Domain).To(Equal(*expected.Domain))
			})

			It("should have correct LoadBalancerIP", func() {
				Expect(*clone.LoadBalancerIP).To(Equal(*expected.LoadBalancerIP))
			})

			It("should have correct Annotations", func() {
				Expect(clone.Annotations).To(Equal(expected.Annotations))
			})

			It("should have correct ExternalPort", func() {
				Expect(*clone.ExternalPort).To(Equal(*expected.ExternalPort))
			})
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
				original = NewServiceSpecOne()
				defaults = nil
				expected = original
			})

			ValidateResult()
		})

		When("original is nil", func() {
			BeforeEach(func() {
				defaults = NewServiceSpecOne()
				original = nil
				expected = defaults
			})

			ValidateResult()
		})

		When("all original fields are set", func() {
			BeforeEach(func() {
				original = NewServiceSpecOne()
				defaults = NewServiceSpecTwo()
				expected = original
			})

			ValidateResult()
		})

		When("original Enabled is nil", func() {
			BeforeEach(func() {
				original = NewServiceSpecOne()
				original.Enabled = nil
				defaults = NewServiceSpecTwo()

				expected = NewServiceSpecOne()
				expected.Enabled = defaults.Enabled
			})

			ValidateResult()
		})

		When("original Type is nil", func() {
			BeforeEach(func() {
				original = NewServiceSpecOne()
				original.Type = nil
				defaults = NewServiceSpecTwo()

				expected = NewServiceSpecOne()
				expected.Type = defaults.Type
			})

			ValidateResult()
		})

		When("original Domain is nil", func() {
			BeforeEach(func() {
				original = NewServiceSpecOne()
				original.Domain = nil
				defaults = NewServiceSpecTwo()

				expected = NewServiceSpecOne()
				expected.Domain = defaults.Domain
			})

			ValidateResult()
		})
		When("original LoadBalancerIP is nil", func() {
			BeforeEach(func() {
				original = NewServiceSpecOne()
				original.LoadBalancerIP = nil
				defaults = NewServiceSpecTwo()

				expected = NewServiceSpecOne()
				expected.LoadBalancerIP = defaults.LoadBalancerIP
			})

			ValidateResult()
		})

		When("original Annotations is nil", func() {
			BeforeEach(func() {
				original = NewServiceSpecOne()
				original.Annotations = nil
				defaults = NewServiceSpecTwo()

				expected = NewServiceSpecOne()
				expected.Annotations = defaults.Annotations
			})

			ValidateResult()
		})

		When("original ExternalPort is nil", func() {
			BeforeEach(func() {
				original = NewServiceSpecOne()
				original.ExternalPort = nil
				defaults = NewServiceSpecTwo()

				expected = NewServiceSpecOne()
				expected.ExternalPort = defaults.ExternalPort
			})

			ValidateResult()
		})
	})
})
