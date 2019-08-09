package v1_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	coherence "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
)

var _ = Describe("Testing CoherenceServiceSpec struct", func() {

	Context("Copying a CoherenceServiceSpec using DeepCopyWithDefaults", func() {
		var original *coherence.CoherenceServiceSpec
		var defaults *coherence.CoherenceServiceSpec
		var clone *coherence.CoherenceServiceSpec
		var expected *coherence.CoherenceServiceSpec

		NewCoherenceServiceSpecOne := func() *coherence.CoherenceServiceSpec {
			return &coherence.CoherenceServiceSpec{
				ServiceSpec:        coherence.ServiceSpec{
					Enabled:        boolPtr(true),
					Type:           stringPtr("LoadBalancerIP"),
					Domain:         stringPtr("cluster.local"),
					LoadBalancerIP: stringPtr("10.10.10.20"),
					Annotations:    map[string]string{ "foo": "1"},
					ExternalPort:   int32Ptr(20000),
				},
				ManagementHttpPort: int32Ptr(30000),
				MetricsHttpPort:    int32Ptr(9612),
			}
		}

		NewCoherenceServiceSpecTwo := func() *coherence.CoherenceServiceSpec {
			return &coherence.CoherenceServiceSpec{
				ServiceSpec:        coherence.ServiceSpec{
					Enabled:      boolPtr(true),
					Type:         stringPtr("ClusterIP"),
					Domain:       stringPtr("cluster.local2"),
					LoadBalancerIP: stringPtr("10.10.10.21"),
					Annotations:  map[string]string{ "foo": "2"},
					ExternalPort: int32Ptr(20001),
				},
				ManagementHttpPort: int32Ptr(30001),
				MetricsHttpPort:    int32Ptr(9613),
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

			It("should have correct ManagementHttpPort", func() {
				Expect(*clone.ManagementHttpPort).To(Equal(*expected.ManagementHttpPort))
			})

			It("should have correct MetricsHttpPort", func() {
				Expect(*clone.MetricsHttpPort).To(Equal(*expected.MetricsHttpPort))
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
				original = NewCoherenceServiceSpecOne()
				defaults = nil
				expected = original
			})

			ValidateResult()
		})

		When("original is nil", func() {
			BeforeEach(func() {
				defaults = NewCoherenceServiceSpecOne()
				original = nil
				expected = defaults
			})

			ValidateResult()
		})

		When("all original fields are set", func() {
			BeforeEach(func() {
				original = NewCoherenceServiceSpecOne()
				defaults = NewCoherenceServiceSpecTwo()
				expected = original
			})

			ValidateResult()
		})

		When("original Enabled is nil", func() {
			BeforeEach(func() {
				original = NewCoherenceServiceSpecOne()
				original.Enabled = nil
				defaults = NewCoherenceServiceSpecTwo()

				expected = NewCoherenceServiceSpecOne()
				expected.Enabled = defaults.Enabled
			})

			ValidateResult()
		})

		When("original Type is nil", func() {
			BeforeEach(func() {
				original = NewCoherenceServiceSpecOne()
				original.Type = nil
				defaults = NewCoherenceServiceSpecTwo()

				expected = NewCoherenceServiceSpecOne()
				expected.Type = defaults.Type
			})

			ValidateResult()
		})

		When("original Domain is nil", func() {
			BeforeEach(func() {
				original = NewCoherenceServiceSpecOne()
				original.Domain = nil
				defaults = NewCoherenceServiceSpecTwo()

				expected = NewCoherenceServiceSpecOne()
				expected.Domain = defaults.Domain
			})

			ValidateResult()
		})
		When("original LoadBalancerIP is nil", func() {
			BeforeEach(func() {
				original = NewCoherenceServiceSpecOne()
				original.LoadBalancerIP = nil
				defaults = NewCoherenceServiceSpecTwo()

				expected = NewCoherenceServiceSpecOne()
				expected.LoadBalancerIP = defaults.LoadBalancerIP
			})

			ValidateResult()
		})

		When("original Annotations is nil", func() {
			BeforeEach(func() {
				original = NewCoherenceServiceSpecOne()
				original.Annotations = nil
				defaults = NewCoherenceServiceSpecTwo()

				expected = NewCoherenceServiceSpecOne()
				expected.Annotations = defaults.Annotations
			})

			ValidateResult()
		})

		When("original ExternalPort is nil", func() {
			BeforeEach(func() {
				original = NewCoherenceServiceSpecOne()
				original.ExternalPort = nil
				defaults = NewCoherenceServiceSpecTwo()

				expected = NewCoherenceServiceSpecOne()
				expected.ExternalPort = defaults.ExternalPort
			})

			ValidateResult()
		})

		When("original ManagementHttpPort is nil", func() {
			BeforeEach(func() {
				original = NewCoherenceServiceSpecOne()
				original.ManagementHttpPort = nil
				defaults = NewCoherenceServiceSpecTwo()

				expected = NewCoherenceServiceSpecOne()
				expected.ManagementHttpPort = defaults.ManagementHttpPort
			})

			ValidateResult()
		})

		When("original MetricsHttpPort is nil", func() {
			BeforeEach(func() {
				original = NewCoherenceServiceSpecOne()
				original.MetricsHttpPort = nil
				defaults = NewCoherenceServiceSpecTwo()

				expected = NewCoherenceServiceSpecOne()
				expected.MetricsHttpPort = defaults.MetricsHttpPort
			})

			ValidateResult()
		})
	})
})
