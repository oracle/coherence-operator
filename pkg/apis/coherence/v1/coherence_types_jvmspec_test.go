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

var _ = Describe("Testing JVM JvmDebugSpec struct", func() {

	Context("Copying an JvmDebugSpec using DeepCopyWithDefaults", func() {
		var original *coherence.JvmDebugSpec
		var defaults *coherence.JvmDebugSpec
		var clone *coherence.JvmDebugSpec

		var jvmOne = &coherence.JvmDebugSpec{
			Enabled: boolPtr(true),
			Suspend: boolPtr(true),
			Attach:  stringPtr("10.10.100.10:1234"),
			Port:    int32Ptr(8080),
		}

		var jvmTwo = &coherence.JvmDebugSpec{
			Enabled: boolPtr(false),
			Suspend: boolPtr(false),
			Attach:  stringPtr("10.10.100.99:9876"),
			Port:    int32Ptr(9090),
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
				original = jvmOne
				defaults = nil
			})

			It("should copy the original", func() {
				Expect(clone).To(Equal(original))
			})
		})

		When("defaults is empty", func() {
			BeforeEach(func() {
				original = jvmOne
				defaults = &coherence.JvmDebugSpec{}
			})

			It("should copy the original", func() {
				Expect(clone).To(Equal(original))
			})
		})

		When("original is nil", func() {
			BeforeEach(func() {
				original = nil
				defaults = jvmTwo
			})

			It("should copy the defaults", func() {
				Expect(clone).To(Equal(defaults))
			})
		})

		When("original is empty", func() {
			BeforeEach(func() {
				original = &coherence.JvmDebugSpec{}
				defaults = jvmTwo
			})

			It("should copy the defaults", func() {
				Expect(clone).To(Equal(defaults))
			})
		})

		When("original Enabled is nil", func() {
			BeforeEach(func() {
				original = jvmOne.DeepCopy()
				original.Enabled = nil
				defaults = jvmTwo.DeepCopy()
			})

			It("should copy the defaults Enabled", func() {
				expected := original.DeepCopy()
				expected.Enabled = defaults.Enabled
				Expect(clone).To(Equal(expected))
			})
		})

		When("original Port is nil", func() {
			BeforeEach(func() {
				original = jvmOne.DeepCopy()
				original.Port = nil
				defaults = jvmTwo.DeepCopy()
			})

			It("should copy the defaults Port", func() {
				expected := original.DeepCopy()
				expected.Port = defaults.Port
				Expect(clone).To(Equal(expected))
			})
		})

		When("original Suspend is nil", func() {
			BeforeEach(func() {
				original = jvmOne.DeepCopy()
				original.Suspend = nil
				defaults = jvmTwo.DeepCopy()
			})

			It("should copy the defaults Suspend", func() {
				expected := original.DeepCopy()
				expected.Suspend = defaults.Suspend
				Expect(clone).To(Equal(expected))
			})
		})

		When("original Attach is nil", func() {
			BeforeEach(func() {
				original = jvmOne.DeepCopy()
				original.Attach = nil
				defaults = jvmTwo.DeepCopy()
			})

			It("should copy the defaults Attach", func() {
				expected := original.DeepCopy()
				expected.Attach = defaults.Attach
				Expect(clone).To(Equal(expected))
			})
		})
	})
})

var _ = Describe("Testing JVMGarbageCollectorSpec struct", func() {

	Context("Copying an JVMGarbageCollectorSpec using DeepCopyWithDefaults", func() {
		var original *coherence.JvmGarbageCollectorSpec
		var defaults *coherence.JvmGarbageCollectorSpec
		var clone *coherence.JvmGarbageCollectorSpec

		var jvmOne = &coherence.JvmGarbageCollectorSpec{
			Collector: stringPtr("foo"),
			Args:      []string{"argOne", "argTwo"},
			Logging:   boolPtr(true),
		}

		var jvmTwo = &coherence.JvmGarbageCollectorSpec{
			Collector: stringPtr("bar"),
			Args:      []string{"argThree", "argFour"},
			Logging:   boolPtr(false),
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
				original = jvmOne
				defaults = nil
			})

			It("should copy the original", func() {
				Expect(clone).To(Equal(original))
			})
		})

		When("defaults is empty", func() {
			BeforeEach(func() {
				original = jvmOne
				defaults = &coherence.JvmGarbageCollectorSpec{}
			})

			It("should copy the original", func() {
				Expect(clone).To(Equal(original))
			})
		})

		When("original is nil", func() {
			BeforeEach(func() {
				original = nil
				defaults = jvmTwo
			})

			It("should copy the defaults", func() {
				Expect(clone).To(Equal(defaults))
			})
		})

		When("original is empty", func() {
			BeforeEach(func() {
				original = &coherence.JvmGarbageCollectorSpec{}
				defaults = jvmTwo
			})

			It("should copy the defaults", func() {
				Expect(clone).To(Equal(defaults))
			})
		})

		When("original Collector is nil", func() {
			BeforeEach(func() {
				original = jvmOne.DeepCopy()
				original.Collector = nil
				defaults = jvmTwo.DeepCopy()
			})

			It("should copy the defaults Collector", func() {
				expected := original.DeepCopy()
				expected.Collector = defaults.Collector
				Expect(clone).To(Equal(expected))
			})
		})

		When("original Logging is nil", func() {
			BeforeEach(func() {
				original = jvmOne.DeepCopy()
				original.Logging = nil
				defaults = jvmTwo.DeepCopy()
			})

			It("should copy the defaults Logging", func() {
				expected := original.DeepCopy()
				expected.Logging = defaults.Logging
				Expect(clone).To(Equal(expected))
			})
		})

		When("original Args is nil", func() {
			BeforeEach(func() {
				original = jvmOne.DeepCopy()
				original.Args = nil
				defaults = jvmTwo.DeepCopy()
			})

			It("should copy the defaults Args", func() {
				expected := original.DeepCopy()
				expected.Args = defaults.Args
				Expect(clone).To(Equal(expected))
			})
		})

		When("original Args is empty", func() {
			BeforeEach(func() {
				original = jvmOne.DeepCopy()
				original.Args = []string{}
				defaults = jvmTwo.DeepCopy()
			})

			It("should have empty Args", func() {
				expected := original.DeepCopy()
				expected.Args = []string{}
				Expect(clone).To(Equal(expected))
			})
		})
	})
})

var _ = Describe("Testing JvmMemorySpec struct", func() {

	Context("Copying an JvmMemorySpec using DeepCopyWithDefaults", func() {
		var original *coherence.JvmMemorySpec
		var defaults *coherence.JvmMemorySpec
		var clone *coherence.JvmMemorySpec

		var jvmOne = &coherence.JvmMemorySpec{
			HeapSize:             stringPtr("10g"),
			StackSize:            stringPtr("50m"),
			MetaspaceSize:        stringPtr("256m"),
			NativeMemoryTracking: stringPtr("foo"),
			OnOutOfMemory: &coherence.JvmOutOfMemorySpec{
				Exit:     boolPtr(true),
				HeapDump: boolPtr(true),
			},
		}

		var jvmTwo = &coherence.JvmMemorySpec{
			HeapSize:             stringPtr("5g"),
			StackSize:            stringPtr("5m"),
			MetaspaceSize:        stringPtr("25m"),
			NativeMemoryTracking: stringPtr("bar"),
			OnOutOfMemory: &coherence.JvmOutOfMemorySpec{
				Exit:     boolPtr(false),
				HeapDump: boolPtr(false),
			},
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
				original = jvmOne
				defaults = nil
			})

			It("should copy the original", func() {
				Expect(clone).To(Equal(original))
			})
		})

		When("defaults is empty", func() {
			BeforeEach(func() {
				original = jvmOne
				defaults = &coherence.JvmMemorySpec{}
			})

			It("should copy the original", func() {
				Expect(clone).To(Equal(original))
			})
		})

		When("original is nil", func() {
			BeforeEach(func() {
				original = nil
				defaults = jvmTwo
			})

			It("should copy the defaults", func() {
				Expect(clone).To(Equal(defaults))
			})
		})

		When("original is empty", func() {
			BeforeEach(func() {
				original = &coherence.JvmMemorySpec{}
				defaults = jvmTwo
			})

			It("should copy the defaults", func() {
				Expect(clone).To(Equal(defaults))
			})
		})

		When("original HeapSize is nil", func() {
			BeforeEach(func() {
				original = jvmOne.DeepCopy()
				original.HeapSize = nil
				defaults = jvmTwo.DeepCopy()
			})

			It("should copy the defaults HeapSize", func() {
				expected := original.DeepCopy()
				expected.HeapSize = defaults.HeapSize
				Expect(clone).To(Equal(expected))
			})
		})

		When("original MetaspaceSize is nil", func() {
			BeforeEach(func() {
				original = jvmOne.DeepCopy()
				original.MetaspaceSize = nil
				defaults = jvmTwo.DeepCopy()
			})

			It("should copy the defaults MetaspaceSize", func() {
				expected := original.DeepCopy()
				expected.MetaspaceSize = defaults.MetaspaceSize
				Expect(clone).To(Equal(expected))
			})
		})

		When("original StackSize is nil", func() {
			BeforeEach(func() {
				original = jvmOne.DeepCopy()
				original.StackSize = nil
				defaults = jvmTwo.DeepCopy()
			})

			It("should copy the defaults StackSize", func() {
				expected := original.DeepCopy()
				expected.StackSize = defaults.StackSize
				Expect(clone).To(Equal(expected))
			})
		})

		When("original OnOutOfMemory is nil", func() {
			BeforeEach(func() {
				original = jvmOne.DeepCopy()
				original.OnOutOfMemory = nil
				defaults = jvmTwo.DeepCopy()
			})

			It("should copy the defaults OnOutOfMemory", func() {
				expected := original.DeepCopy()
				expected.OnOutOfMemory = defaults.OnOutOfMemory
				Expect(clone).To(Equal(expected))
			})
		})
	})
})

var _ = Describe("Testing JvmOutOfMemorySpec struct", func() {

	Context("Copying an JvmOutOfMemorySpec using DeepCopyWithDefaults", func() {
		var original *coherence.JvmOutOfMemorySpec
		var defaults *coherence.JvmOutOfMemorySpec
		var clone *coherence.JvmOutOfMemorySpec

		var jvmOne = &coherence.JvmOutOfMemorySpec{
			Exit:     boolPtr(true),
			HeapDump: boolPtr(true),
		}

		var jvmTwo = &coherence.JvmOutOfMemorySpec{
			Exit:     boolPtr(false),
			HeapDump: boolPtr(false),
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
				original = jvmOne
				defaults = nil
			})

			It("should copy the original", func() {
				Expect(clone).To(Equal(original))
			})
		})

		When("defaults is empty", func() {
			BeforeEach(func() {
				original = jvmOne
				defaults = &coherence.JvmOutOfMemorySpec{}
			})

			It("should copy the original", func() {
				Expect(clone).To(Equal(original))
			})
		})

		When("original is nil", func() {
			BeforeEach(func() {
				original = nil
				defaults = jvmTwo
			})

			It("should copy the defaults", func() {
				Expect(clone).To(Equal(defaults))
			})
		})

		When("original is empty", func() {
			BeforeEach(func() {
				original = &coherence.JvmOutOfMemorySpec{}
				defaults = jvmTwo
			})

			It("should copy the defaults", func() {
				Expect(clone).To(Equal(defaults))
			})
		})

		When("original Exit is nil", func() {
			BeforeEach(func() {
				original = jvmOne.DeepCopy()
				original.Exit = nil
				defaults = jvmTwo.DeepCopy()
			})

			It("should copy the defaults Exit", func() {
				expected := original.DeepCopy()
				expected.Exit = defaults.Exit
				Expect(clone).To(Equal(expected))
			})
		})

		When("original HeapDump is nil", func() {
			BeforeEach(func() {
				original = jvmOne.DeepCopy()
				original.HeapDump = nil
				defaults = jvmTwo.DeepCopy()
			})

			It("should copy the defaults HeapDump", func() {
				expected := original.DeepCopy()
				expected.HeapDump = defaults.HeapDump
				Expect(clone).To(Equal(expected))
			})
		})
	})
})

var _ = Describe("Testing JVMSpec struct", func() {

	Context("Copying an JVMSpec using DeepCopyWithDefaults", func() {
		var original *coherence.JVMSpec
		var defaults *coherence.JVMSpec
		var clone *coherence.JVMSpec

		var jvmOne = &coherence.JVMSpec{
			Debug: &coherence.JvmDebugSpec{
				Enabled: boolPtr(true),
				Suspend: boolPtr(true),
				Attach:  stringPtr("10.1.2.3:1234"),
				Port:    int32Ptr(8080),
			},
			UseContainerLimits: boolPtr(true),
			FlightRecorder:     boolPtr(true),
			Gc: &coherence.JvmGarbageCollectorSpec{
				Collector: stringPtr("G1"),
				Args:      []string{"one", "two"},
				Logging:   boolPtr(true),
			},
			DiagnosticsVolume: &corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: "/foo",
				},
			},
			Memory: &coherence.JvmMemorySpec{
				HeapSize:             stringPtr("10g"),
				StackSize:            stringPtr("500m"),
				MetaspaceSize:        stringPtr("256mg"),
				NativeMemoryTracking: stringPtr("off"),
				OnOutOfMemory: &coherence.JvmOutOfMemorySpec{
					Exit:     boolPtr(true),
					HeapDump: boolPtr(true),
				},
			},
		}

		var jvmTwo = &coherence.JVMSpec{
			Debug: &coherence.JvmDebugSpec{
				Enabled: boolPtr(false),
				Suspend: boolPtr(false),
				Attach:  stringPtr("10.1.2.99:9876"),
				Port:    int32Ptr(9090),
			},
			UseContainerLimits: boolPtr(false),
			FlightRecorder:     boolPtr(false),
			Gc: &coherence.JvmGarbageCollectorSpec{
				Collector: stringPtr("Parallel"),
				Args:      []string{"three", "four"},
				Logging:   boolPtr(false),
			},
			DiagnosticsVolume: &corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: "/bar",
				},
			},
			Memory: &coherence.JvmMemorySpec{
				HeapSize:             stringPtr("90g"),
				StackSize:            stringPtr("900m"),
				MetaspaceSize:        stringPtr("956mg"),
				NativeMemoryTracking: stringPtr("all"),
				OnOutOfMemory: &coherence.JvmOutOfMemorySpec{
					Exit:     boolPtr(false),
					HeapDump: boolPtr(false),
				},
			},
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
				original = jvmOne
				defaults = nil
			})

			It("should copy the original", func() {
				Expect(clone).To(Equal(original))
			})
		})

		When("defaults is empty", func() {
			BeforeEach(func() {
				original = jvmOne
				defaults = &coherence.JVMSpec{}
			})

			It("should copy the original", func() {
				Expect(clone).To(Equal(original))
			})
		})

		When("original is nil", func() {
			BeforeEach(func() {
				original = nil
				defaults = jvmTwo
			})

			It("should copy the defaults", func() {
				Expect(clone).To(Equal(defaults))
			})
		})

		When("original is empty", func() {
			BeforeEach(func() {
				original = &coherence.JVMSpec{}
				defaults = jvmTwo
			})

			It("should copy the defaults", func() {
				Expect(clone).To(Equal(defaults))
			})
		})

		When("original Debug is nil", func() {
			BeforeEach(func() {
				original = jvmOne.DeepCopy()
				original.Debug = nil
				defaults = jvmTwo.DeepCopy()
			})

			It("should copy the defaults Debug", func() {
				expected := original.DeepCopy()
				expected.Debug = defaults.Debug
				Expect(clone).To(Equal(expected))
			})
		})

		When("original DiagnosticsVolume is nil", func() {
			BeforeEach(func() {
				original = jvmOne.DeepCopy()
				original.DiagnosticsVolume = nil
				defaults = jvmTwo.DeepCopy()
			})

			It("should copy the defaults DiagnosticsVolume", func() {
				expected := original.DeepCopy()
				expected.DiagnosticsVolume = defaults.DiagnosticsVolume
				Expect(clone).To(Equal(expected))
			})
		})

		When("original FlightRecorder is nil", func() {
			BeforeEach(func() {
				original = jvmOne.DeepCopy()
				original.FlightRecorder = nil
				defaults = jvmTwo.DeepCopy()
			})

			It("should copy the defaults FlightRecorder", func() {
				expected := original.DeepCopy()
				expected.FlightRecorder = defaults.FlightRecorder
				Expect(clone).To(Equal(expected))
			})
		})

		When("original GC is nil", func() {
			BeforeEach(func() {
				original = jvmOne.DeepCopy()
				original.Gc = nil
				defaults = jvmTwo.DeepCopy()
			})

			It("should copy the defaults GC", func() {
				expected := original.DeepCopy()
				expected.Gc = defaults.Gc
				Expect(clone).To(Equal(expected))
			})
		})

		When("original Memory is nil", func() {
			BeforeEach(func() {
				original = jvmOne.DeepCopy()
				original.Memory = nil
				defaults = jvmTwo.DeepCopy()
			})

			It("should copy the defaults Memory", func() {
				expected := original.DeepCopy()
				expected.Memory = defaults.Memory
				Expect(clone).To(Equal(expected))
			})
		})

		When("original UseContainerLimits is nil", func() {
			BeforeEach(func() {
				original = jvmOne.DeepCopy()
				original.UseContainerLimits = nil
				defaults = jvmTwo.DeepCopy()
			})

			It("should copy the defaults UseContainerLimits", func() {
				expected := original.DeepCopy()
				expected.UseContainerLimits = defaults.UseContainerLimits
				Expect(clone).To(Equal(expected))
			})
		})

		When("original Args is nil", func() {
			BeforeEach(func() {
				original = jvmOne.DeepCopy()
				original.Args = nil
				defaults = jvmTwo.DeepCopy()
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
				original = jvmOne.DeepCopy()
				original.Args = []string{}
				defaults = jvmTwo.DeepCopy()
				defaults.Args = []string{"one", "two"}
			})

			It("should merge the original and defaults Args", func() {
				expected := original.DeepCopy()
				expected.Args = []string{"one", "two"}
				Expect(clone).To(Equal(expected))
			})
		})

		When("defaults Args is nil", func() {
			BeforeEach(func() {
				original = jvmOne.DeepCopy()
				original.Args = []string{"one", "two"}
				defaults = jvmTwo.DeepCopy()
				defaults.Args = nil
			})

			It("should copy the original Args", func() {
				expected := original.DeepCopy()
				Expect(clone).To(Equal(expected))
			})
		})

		When("defaults Args is empty", func() {
			BeforeEach(func() {
				original = jvmOne.DeepCopy()
				original.Args = []string{"one", "two"}
				defaults = jvmTwo.DeepCopy()
				defaults.Args = []string{}
			})

			It("should copy the original Args", func() {
				expected := original.DeepCopy()
				Expect(clone).To(Equal(expected))
			})
		})

		When("original and defaults Args is populated", func() {
			BeforeEach(func() {
				original = jvmOne.DeepCopy()
				original.Args = []string{"one", "two"}
				defaults = jvmTwo.DeepCopy()
				defaults.Args = []string{"three", "four"}
			})

			It("should merge the original and default Args, default first", func() {
				expected := original.DeepCopy()
				expected.Args = []string{"three", "four", "one", "two"}
				Expect(clone).To(Equal(expected))
			})
		})
	})
})
