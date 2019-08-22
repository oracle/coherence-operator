package v1_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	coherence "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	corev1 "k8s.io/api/core/v1"
)

var _ = Describe("Testing JMXSpec struct", func() {

	Context("Copying a JMXSpec using DeepCopyWithDefaults", func() {
		var original *coherence.JMXSpec
		var defaults *coherence.JMXSpec
		var clone *coherence.JMXSpec
		var expected *coherence.JMXSpec
		var clusterIP = corev1.ServiceTypeClusterIP
		var loadBalancer = corev1.ServiceTypeLoadBalancer

		NewJMXSpecOne := func() *coherence.JMXSpec {
			return &coherence.JMXSpec{
				Enabled:  boolPtr(true),
				Replicas: int32Ptr(3),
				MaxHeap:  stringPtr("2Gi"),
				Service: &coherence.ServiceSpec{
					Type:           &clusterIP,
					LoadBalancerIP: stringPtr("10.10.10.20"),
					Annotations:    map[string]string{"foo": "1"},
					Port:           int32Ptr(9099),
				},
			}
		}

		NewJMXSpecTwo := func() *coherence.JMXSpec {
			return &coherence.JMXSpec{
				Enabled:  boolPtr(true),
				Replicas: int32Ptr(6),
				MaxHeap:  stringPtr("3Gi"),
				Service: &coherence.ServiceSpec{
					Type:           &loadBalancer,
					LoadBalancerIP: stringPtr("10.10.10.21"),
					Annotations:    map[string]string{"foo": "2"},
					Port:           int32Ptr(9098),
				},
			}
		}

		ValidateResult := func() {
			It("should have correct Enabled", func() {
				Expect(*clone.Enabled).To(Equal(*expected.Enabled))
			})

			It("should have correct Replicas", func() {
				Expect(*clone.Replicas).To(Equal(*expected.Replicas))
			})

			It("should have correct MaxHeap", func() {
				Expect(*clone.MaxHeap).To(Equal(*expected.MaxHeap))
			})

			It("should have correct Service", func() {
				Expect(*clone.Service).To(Equal(*expected.Service))
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
				original = NewJMXSpecOne()
				defaults = nil
				expected = original
			})

			ValidateResult()
		})

		When("original is nil", func() {
			BeforeEach(func() {
				defaults = NewJMXSpecOne()
				original = nil
				expected = defaults
			})

			ValidateResult()
		})

		When("all original fields are set", func() {
			BeforeEach(func() {
				original = NewJMXSpecOne()
				defaults = NewJMXSpecTwo()
				expected = original
			})

			ValidateResult()
		})

		When("original Enabled is nil", func() {
			BeforeEach(func() {
				original = NewJMXSpecOne()
				original.Enabled = nil
				defaults = NewJMXSpecTwo()

				expected = NewJMXSpecOne()
				expected.Enabled = defaults.Enabled
			})

			ValidateResult()
		})

		When("original Replicas is nil", func() {
			BeforeEach(func() {
				original = NewJMXSpecOne()
				original.Replicas = nil
				defaults = NewJMXSpecTwo()

				expected = NewJMXSpecOne()
				expected.Replicas = defaults.Replicas
			})

			ValidateResult()
		})

		When("original MaxHeap is nil", func() {
			BeforeEach(func() {
				original = NewJMXSpecOne()
				original.MaxHeap = nil
				defaults = NewJMXSpecTwo()

				expected = NewJMXSpecOne()
				expected.MaxHeap = defaults.MaxHeap
			})

			ValidateResult()
		})
		When("original Service is nil", func() {
			BeforeEach(func() {
				original = NewJMXSpecOne()
				original.Service = nil
				defaults = NewJMXSpecTwo()

				expected = NewJMXSpecOne()
				expected.Service = defaults.Service
			})

			ValidateResult()
		})
	})
})
