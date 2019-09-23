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
	"testing"
)

func TestApplicationSpecDeepCopyWithDefaults(t *testing.T) {

}

var _ = Describe("Testing ApplicationSpec struct", func() {

	Context("Copying an ApplicationSpec using DeepCopyWithDefaults", func() {
		var original *coherence.ApplicationSpec
		var defaults *coherence.ApplicationSpec
		var clone *coherence.ApplicationSpec

		var always = corev1.PullAlways
		var never = corev1.PullNever

		var appOne = &coherence.ApplicationSpec{
			Type:      stringPtr("java"),
			MainClass: stringPtr("TestMainOne"),
			Args:      []string{},
			ImageSpec: coherence.ImageSpec{
				Image:           stringPtr("app:1.0"),
				ImagePullPolicy: &always,
			},
			LibDir:    stringPtr("/test/libOne"),
			ConfigDir: stringPtr("/test/confOne"),
		}

		var appTwo = &coherence.ApplicationSpec{
			Type:      stringPtr("node"),
			MainClass: stringPtr("TestMainTwo"),
			Args:      []string{},
			ImageSpec: coherence.ImageSpec{
				Image:           stringPtr("app:2.0"),
				ImagePullPolicy: &never,
			},
			LibDir:    stringPtr("/test/libTwo"),
			ConfigDir: stringPtr("/test/confTwo"),
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
				original = appOne
				defaults = nil
			})

			It("should copy the original Args", func() {
				Expect(clone.Args).To(Equal(original.Args))
			})

			It("should copy the original Type", func() {
				Expect(*clone.Type).To(Equal(*original.Type))
			})

			It("should copy the original Type", func() {
				Expect(*clone.LibDir).To(Equal(*original.LibDir))
			})

			It("should copy the original Type", func() {
				Expect(*clone.ConfigDir).To(Equal(*original.ConfigDir))
			})

			It("should copy the original Image", func() {
				Expect(*clone.Image).To(Equal(*original.Image))
			})

			It("should copy the original ImagePullPolicy", func() {
				Expect(*clone.ImagePullPolicy).To(Equal(*original.ImagePullPolicy))
			})
		})

		When("defaults is empty", func() {
			BeforeEach(func() {
				original = appOne
				defaults = &coherence.ApplicationSpec{}
			})

			It("should copy the original Args", func() {
				Expect(clone.Args).To(Equal(original.Args))
			})

			It("should copy the original Type", func() {
				Expect(*clone.Type).To(Equal(*original.Type))
			})

			It("should copy the original Type", func() {
				Expect(*clone.LibDir).To(Equal(*original.LibDir))
			})

			It("should copy the original Type", func() {
				Expect(*clone.ConfigDir).To(Equal(*original.ConfigDir))
			})

			It("should copy the original Image", func() {
				Expect(*clone.Image).To(Equal(*original.Image))
			})

			It("should copy the original ImagePullPolicy", func() {
				Expect(*clone.ImagePullPolicy).To(Equal(*original.ImagePullPolicy))
			})
		})

		When("original is nil", func() {
			BeforeEach(func() {
				original = nil
				defaults = appTwo
			})

			It("should copy the defaults Args", func() {
				Expect(clone.Args).To(Equal(defaults.Args))
			})

			It("should copy the defaults Type", func() {
				Expect(*clone.Type).To(Equal(*defaults.Type))
			})

			It("should copy the defaults Type", func() {
				Expect(*clone.LibDir).To(Equal(*defaults.LibDir))
			})

			It("should copy the defaults Type", func() {
				Expect(*clone.ConfigDir).To(Equal(*defaults.ConfigDir))
			})

			It("should copy the defaults Image", func() {
				Expect(*clone.Image).To(Equal(*defaults.Image))
			})

			It("should copy the defaults ImagePullPolicy", func() {
				Expect(*clone.ImagePullPolicy).To(Equal(*defaults.ImagePullPolicy))
			})
		})

		When("original is empty", func() {
			BeforeEach(func() {
				original = &coherence.ApplicationSpec{}
				defaults = appTwo
			})

			It("should copy the defaults Args", func() {
				Expect(clone.Args).To(Equal(defaults.Args))
			})

			It("should copy the defaults Type", func() {
				Expect(*clone.Type).To(Equal(*defaults.Type))
			})

			It("should copy the defaults LibDir", func() {
				Expect(*clone.LibDir).To(Equal(*defaults.LibDir))
			})

			It("should copy the defaults ConfigDir", func() {
				Expect(*clone.ConfigDir).To(Equal(*defaults.ConfigDir))
			})

			It("should copy the defaults Image", func() {
				Expect(*clone.Image).To(Equal(*defaults.Image))
			})

			It("should copy the defaults ImagePullPolicy", func() {
				Expect(*clone.ImagePullPolicy).To(Equal(*defaults.ImagePullPolicy))
			})
		})

		When("original Args is nil", func() {
			BeforeEach(func() {
				original = appOne.DeepCopy()
				original.Args = nil
				defaults = appTwo.DeepCopy()
				defaults.Args = []string{"one", "two"}
			})

			It("should copy the defaults Args", func() {
				expected := original.DeepCopy()
				expected.Args = defaults.Args
				Expect(clone).To(Equal(expected))
			})
		})

		When("original Args is empty", func() {
			BeforeEach(func() {
				original = appOne.DeepCopy()
				original.Args = []string{}
				defaults = appTwo.DeepCopy()
				defaults.Args = []string{"one", "two"}
			})

			It("should copy the original Args", func() {
				Expect(clone).To(Equal(original))
			})
		})

		When("defaults Args is nil", func() {
			BeforeEach(func() {
				original = appOne.DeepCopy()
				original.Args = []string{"one", "two"}
				defaults = appTwo.DeepCopy()
				defaults.Args = nil
			})

			It("should copy the original Args", func() {
				expected := original.DeepCopy()
				Expect(clone).To(Equal(expected))
			})
		})

		When("defaults Args is empty", func() {
			BeforeEach(func() {
				original = appOne.DeepCopy()
				original.Args = []string{"one", "two"}
				defaults = appTwo.DeepCopy()
				defaults.Args = []string{}
			})

			It("should copy the original Args", func() {
				expected := original.DeepCopy()
				Expect(clone).To(Equal(expected))
			})
		})

		When("original and defaults Args is populated", func() {
			BeforeEach(func() {
				original = appOne.DeepCopy()
				original.Args = []string{"one", "two"}
				defaults = appTwo.DeepCopy()
				defaults.Args = []string{"three", "four"}
			})

			It("should copy the original Args", func() {
				expected := original.DeepCopy()
				Expect(clone).To(Equal(expected))
			})
		})
	})
})
