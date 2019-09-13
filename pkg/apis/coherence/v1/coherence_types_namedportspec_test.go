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
)

var _ = Describe("Testing NamedPortSpec struct", func() {

	Context("Copying a NamedPortSpec using DeepCopyWithDefaults", func() {
		var original *coherence.NamedPortSpec
		var defaults *coherence.NamedPortSpec
		var clone *coherence.NamedPortSpec
		var expected *coherence.NamedPortSpec

		NewPortSpecOne := func() *coherence.NamedPortSpec {
			return &coherence.NamedPortSpec{
				Name: "foo",
				PortSpec: coherence.PortSpec{
					Port: 8000,
				},
			}
		}

		NewPortSpecTwo := func() *coherence.NamedPortSpec {
			return &coherence.NamedPortSpec{
				Name: "bar",
				PortSpec: coherence.PortSpec{
					Port: 9000,
				},
			}
		}

		ValidateResult := func() {
			It("should have correct Name", func() {
				Expect(clone.Name).To(Equal(expected.Name))
			})

			It("should have correct Protocol", func() {
				Expect(clone.Protocol).To(Equal(expected.Protocol))
			})

			It("should have correct PortSpec", func() {
				Expect(clone.PortSpec).To(Equal(expected.PortSpec))
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

		When("defaults is nil copy should match original", func() {
			BeforeEach(func() {
				original = NewPortSpecOne()
				defaults = nil
				expected = original
			})

			ValidateResult()
		})

		When("original is nil copy should match defaults", func() {
			BeforeEach(func() {
				defaults = NewPortSpecOne()
				original = nil
				expected = defaults
			})

			ValidateResult()
		})

		When("all original fields are set copy should match original", func() {
			BeforeEach(func() {
				original = NewPortSpecOne()
				defaults = NewPortSpecTwo()
				expected = original
			})

			ValidateResult()
		})

		When("original Name is blank copy should have defaults name", func() {
			BeforeEach(func() {
				original = NewPortSpecOne()
				original.Name = ""
				defaults = NewPortSpecTwo()

				expected = NewPortSpecOne()
				expected.Name = defaults.Name
			})

			ValidateResult()
		})

		Context("Merging []NamedPortSpec", func() {
			var primary []coherence.NamedPortSpec
			var secondary []coherence.NamedPortSpec
			var merged []coherence.NamedPortSpec

			var portOne = coherence.NamedPortSpec{
				Name: "One",
				PortSpec: coherence.PortSpec{
					Port:     7000,
					Protocol: stringPtr("TCP"),
				},
			}

			var portTwo = coherence.NamedPortSpec{
				Name: "Two",
				PortSpec: coherence.PortSpec{
					Port:     8000,
					Protocol: stringPtr("UDP"),
				},
			}

			var portThree = coherence.NamedPortSpec{
				Name: "Three",
				PortSpec: coherence.PortSpec{
					Port:     9000,
					Protocol: stringPtr("ABC"),
				},
			}

			JustBeforeEach(func() {
				merged = coherence.MergeNamedPortSpecs(primary, secondary)
			})

			When("primary and secondary slices are nil", func() {
				BeforeEach(func() {
					primary = nil
					secondary = nil
				})

				It("the result should be nil", func() {
					Expect(merged).To(BeNil())
				})
			})

			When("primary slice is not nil and the secondary slice is nil", func() {
				BeforeEach(func() {
					primary = []coherence.NamedPortSpec{portOne, portTwo, portThree}
					secondary = nil
				})

				It("the result should be the primary slice", func() {
					Expect(merged).To(Equal(primary))
				})
			})

			When("primary slice is not nil and the secondary slice is empty", func() {
				BeforeEach(func() {
					primary = []coherence.NamedPortSpec{portOne, portTwo, portThree}
					secondary = []coherence.NamedPortSpec{}
				})

				It("the result should be the primary slice", func() {
					Expect(merged).To(Equal(primary))
				})
			})

			When("primary slice is nil and the secondary slice is not nil", func() {
				BeforeEach(func() {
					primary = nil
					secondary = []coherence.NamedPortSpec{portOne, portTwo, portThree}
				})

				It("the result should be the secondary slice", func() {
					Expect(merged).To(Equal(secondary))
				})
			})

			When("primary slice is empty and the secondary slice is not nil", func() {
				BeforeEach(func() {
					primary = []coherence.NamedPortSpec{}
					secondary = []coherence.NamedPortSpec{portOne, portTwo, portThree}
				})

				It("the result should be the secondary slice", func() {
					Expect(merged).To(Equal(secondary))
				})
			})

			When("primary slice is populated and the secondary slice is populated", func() {
				BeforeEach(func() {
					primary = []coherence.NamedPortSpec{portOne, portTwo}
					secondary = []coherence.NamedPortSpec{portThree}
				})

				It("the result should contain the correct number of ports", func() {
					Expect(len(merged)).To(Equal(3))
				})

				It("the result should contain portOne at position 0", func() {
					Expect(merged[0]).To(Equal(portOne))
				})

				It("the result should contain portTwo at position 1", func() {
					Expect(merged[1]).To(Equal(portTwo))
				})

				It("the result should contain portThree at position 2", func() {
					Expect(merged[2]).To(Equal(portThree))
				})
			})

			When("primary slice is populated and the secondary slice is populated with matching ports", func() {
				var p1 = coherence.NamedPortSpec{
					Name: "Foo",
					PortSpec: coherence.PortSpec{
						Port: 7000,
					},
				}

				var p2 = coherence.NamedPortSpec{
					Name: "Foo",
					PortSpec: coherence.PortSpec{
						Protocol: stringPtr("TCP"),
					},
				}

				var pm = coherence.NamedPortSpec{
					Name: "Foo",
					PortSpec: coherence.PortSpec{
						Port:     7000,
						Protocol: stringPtr("TCP"),
					},
				}

				BeforeEach(func() {
					primary = []coherence.NamedPortSpec{portOne, p1}
					secondary = []coherence.NamedPortSpec{portTwo, p2}
				})

				It("the result should contain the correct number of ports", func() {
					Expect(len(merged)).To(Equal(3))
				})

				It("the result should contain portOne at position 0", func() {
					Expect(merged[0]).To(Equal(portOne))
				})

				It("the result should contain the merged port at position 1", func() {
					Expect(merged[1]).To(Equal(pm))
				})

				It("the result should contain portThree at position 2", func() {
					Expect(merged[2]).To(Equal(portTwo))
				})
			})
		})
	})
})
