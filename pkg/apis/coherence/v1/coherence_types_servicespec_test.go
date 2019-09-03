/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	coherence "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	corev1 "k8s.io/api/core/v1"
)

var _ = Describe("Testing ServiceSpec struct", func() {

	Context("Copying a ServiceSpec using DeepCopyWithDefaults", func() {
		var original *coherence.ServiceSpec
		var defaults *coherence.ServiceSpec
		var clone *coherence.ServiceSpec
		var expected *coherence.ServiceSpec

		var clusterIP = corev1.ServiceTypeClusterIP
		var loadBalancer = corev1.ServiceTypeLoadBalancer

		NewServiceSpecOne := func() *coherence.ServiceSpec {
			return &coherence.ServiceSpec{
				Enabled:        boolPtr(true),
				Type:           &clusterIP,
				Port:           int32Ptr(80),
				LoadBalancerIP: stringPtr("10.10.10.20"),
				Annotations:    map[string]string{"foo": "1"},
			}
		}

		NewServiceSpecTwo := func() *coherence.ServiceSpec {
			return &coherence.ServiceSpec{
				Enabled:        boolPtr(true),
				Type:           &loadBalancer,
				Port:           int32Ptr(8080),
				LoadBalancerIP: stringPtr("10.10.10.21"),
				Annotations:    map[string]string{"foo": "2"},
			}
		}

		ValidateResult := func() {
			It("should have correct Enabled", func() {
				Expect(*clone.Enabled).To(Equal(*expected.Enabled))
			})

			It("should have correct Type", func() {
				Expect(*clone.Type).To(Equal(*expected.Type))
			})

			It("should have correct Port", func() {
				Expect(*clone.Port).To(Equal(*expected.Port))
			})

			It("should have correct LoadBalancerIP", func() {
				Expect(*clone.LoadBalancerIP).To(Equal(*expected.LoadBalancerIP))
			})

			It("should have correct Annotations", func() {
				Expect(clone.Annotations).To(Equal(expected.Annotations))
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

		When("original Port is nil", func() {
			BeforeEach(func() {
				original = NewServiceSpecOne()
				original.Port = nil
				defaults = NewServiceSpecTwo()

				expected = NewServiceSpecOne()
				expected.Port = defaults.Port
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
	})
})
