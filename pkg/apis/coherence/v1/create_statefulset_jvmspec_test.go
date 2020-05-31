/*
 * Copyright (c) 2020, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1_test

import (
	coh "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	corev1 "k8s.io/api/core/v1"
	"testing"
)

func TestCreateStatefulSetWithEmptyJvmSpec(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		JVM: &coh.JVMSpec{},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithJvmSpecWithEmptyArgs(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		JVM: &coh.JVMSpec{
			Args: []string{},
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithJvmSpecWithArgs(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		JVM: &coh.JVMSpec{
			Args: []string{"argOne", "argTwo"},
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	addEnvVars(stsExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: coh.EnvVarJvmArgs, Value: "argOne argTwo"})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithJvmSpecWithEmptyClasspath(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		JVM: &coh.JVMSpec{
			Classpath: []string{},
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithJvmSpecWithClasspath(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		JVM: &coh.JVMSpec{
			Classpath: []string{"/foo", "/bar"},
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	addEnvVars(stsExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: coh.EnvVarJvmExtraClasspath, Value: "/foo:/bar"})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithJvmSpecWithUseContainerLimitsTrue(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		JVM: &coh.JVMSpec{
			UseContainerLimits: boolPtr(true),
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	addEnvVars(stsExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: coh.EnvVarJvmUseContainerLimits, Value: "true"})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithJvmSpecWithUseContainerLimitsFalse(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		JVM: &coh.JVMSpec{
			UseContainerLimits: boolPtr(false),
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	addEnvVars(stsExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: coh.EnvVarJvmUseContainerLimits, Value: "false"})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithJvmSpecWithDebugEnabledFalse(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		JVM: &coh.JVMSpec{
			Debug: &coh.JvmDebugSpec{
				Enabled: boolPtr(false),
				Suspend: boolPtr(true),
				Attach:  stringPtr("10.10.10.10:5001"),
				Port:    int32Ptr(1234),
			},
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithJvmSpecWithDebugEnabledTrueSuspendTrue(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		JVM: &coh.JVMSpec{
			Debug: &coh.JvmDebugSpec{
				Enabled: boolPtr(true),
				Suspend: boolPtr(true),
				Attach:  stringPtr("10.10.10.10:5001"),
				Port:    int32Ptr(1234),
			},
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	addEnvVars(stsExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: "JVM_DEBUG_ENABLED", Value: "true"})
	addEnvVars(stsExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: "JVM_DEBUG_SUSPEND", Value: "true"})
	addEnvVars(stsExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: "JVM_DEBUG_ATTACH", Value: "10.10.10.10:5001"})
	addEnvVars(stsExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: "JVM_DEBUG_PORT", Value: "1234"})
	// add the expected debug port
	addPorts(stsExpected, coh.ContainerNameCoherence, corev1.ContainerPort{
		Name:          coh.PortNameDebug,
		ContainerPort: 1234,
	})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithJvmSpecWithDebugEnabledTrueSuspendFalse(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		JVM: &coh.JVMSpec{
			Debug: &coh.JvmDebugSpec{
				Enabled: boolPtr(true),
				Suspend: boolPtr(false),
				Attach:  stringPtr("10.10.10.10:5001"),
				Port:    int32Ptr(1234),
			},
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	addEnvVars(stsExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: "JVM_DEBUG_ENABLED", Value: "true"})
	addEnvVars(stsExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: "JVM_DEBUG_ATTACH", Value: "10.10.10.10:5001"})
	addEnvVars(stsExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: "JVM_DEBUG_PORT", Value: "1234"})
	// add the expected debug port
	addPorts(stsExpected, coh.ContainerNameCoherence, corev1.ContainerPort{
		Name:          coh.PortNameDebug,
		ContainerPort: 1234,
	})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithJvmSpecWithGarbageCollector(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		JVM: &coh.JVMSpec{
			Gc: &coh.JvmGarbageCollectorSpec{
				Collector: stringPtr("G1"),
			},
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	addEnvVars(stsExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: "JVM_GC_COLLECTOR", Value: "G1"})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithJvmSpecWithGarbageCollectorArgs(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		JVM: &coh.JVMSpec{
			Gc: &coh.JvmGarbageCollectorSpec{
				Args:    []string{"-XX:GC-ArgOne", "-XX:GC-ArgTwo"},
				Logging: nil,
			},
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	addEnvVars(stsExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: "JVM_GC_ARGS", Value: "-XX:GC-ArgOne -XX:GC-ArgTwo"})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithJvmSpecWithGarbageCollectorLoggingFalse(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		JVM: &coh.JVMSpec{
			Gc: &coh.JvmGarbageCollectorSpec{
				Logging: boolPtr(false),
			},
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	addEnvVars(stsExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: "JVM_GC_LOGGING", Value: "false"})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithJvmSpecWithGarbageCollectorLoggingTrue(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		JVM: &coh.JVMSpec{
			Gc: &coh.JvmGarbageCollectorSpec{
				Logging: boolPtr(true),
			},
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	addEnvVars(stsExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: "JVM_GC_LOGGING", Value: "true"})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithJvmSpecWithDiagnosticsVolume(t *testing.T) {

	hostPath := &corev1.HostPathVolumeSource{Path: "/home/root/debug"}
	spec := coh.CoherenceResourceSpec{
		JVM: &coh.JVMSpec{
			DiagnosticsVolume: &corev1.VolumeSource{
				HostPath: hostPath,
			},
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet with the specified JVM diagnostic volume
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	coh.ReplaceVolume(stsExpected, corev1.Volume{
		Name: coh.VolumeNameJVM,
		VolumeSource: corev1.VolumeSource{
			HostPath: hostPath,
		},
	})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithJvmSpecWithMemorySettings(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		JVM: &coh.JVMSpec{
			Memory: &coh.JvmMemorySpec{
				HeapSize:             stringPtr("5g"),
				StackSize:            stringPtr("500m"),
				MetaspaceSize:        stringPtr("1g"),
				DirectMemorySize:     stringPtr("4g"),
				NativeMemoryTracking: stringPtr("detail"),
			},
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	addEnvVars(stsExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: "JVM_HEAP_SIZE", Value: "5g"})
	addEnvVars(stsExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: "JVM_STACK_SIZE", Value: "500m"})
	addEnvVars(stsExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: "JVM_METASPACE_SIZE", Value: "1g"})
	addEnvVars(stsExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: "JVM_DIRECT_MEMORY_SIZE", Value: "4g"})
	addEnvVars(stsExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: "JVM_NATIVE_MEMORY_TRACKING", Value: "detail"})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithJvmSpecWithExitOnOomTrue(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		JVM: &coh.JVMSpec{
			Memory: &coh.JvmMemorySpec{
				OnOutOfMemory: &coh.JvmOutOfMemorySpec{
					Exit:     boolPtr(true),
					HeapDump: nil,
				},
			},
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	addEnvVars(stsExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: "JVM_OOM_EXIT", Value: "true"})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithJvmSpecWithExitOnOomFalse(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		JVM: &coh.JVMSpec{
			Memory: &coh.JvmMemorySpec{
				OnOutOfMemory: &coh.JvmOutOfMemorySpec{
					Exit:     boolPtr(false),
					HeapDump: nil,
				},
			},
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	addEnvVars(stsExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: "JVM_OOM_EXIT", Value: "false"})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithJvmSpecWithHeapDumpOnOomTrue(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		JVM: &coh.JVMSpec{
			Memory: &coh.JvmMemorySpec{
				OnOutOfMemory: &coh.JvmOutOfMemorySpec{
					HeapDump: boolPtr(true),
				},
			},
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	addEnvVars(stsExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: "JVM_OOM_HEAP_DUMP", Value: "true"})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithJvmSpecWithHeapDumpOnOomFalse(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		JVM: &coh.JVMSpec{
			Memory: &coh.JvmMemorySpec{
				OnOutOfMemory: &coh.JvmOutOfMemorySpec{
					HeapDump: boolPtr(false),
				},
			},
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	addEnvVars(stsExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: "JVM_OOM_HEAP_DUMP", Value: "false"})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithJvmSpecWithJmxmpEnabledTrue(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		JVM: &coh.JVMSpec{
			Jmxmp: &coh.JvmJmxmpSpec{
				Enabled: boolPtr(true),
			},
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	addEnvVars(stsExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: "JVM_JMXMP_ENABLED", Value: "true"})
	addEnvVars(stsExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: "JVM_JMXMP_PORT", Value: "9099"})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithJvmSpecWithJmxmpEnabledFalse(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		JVM: &coh.JVMSpec{
			Jmxmp: &coh.JvmJmxmpSpec{
				Enabled: boolPtr(false),
			},
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	addEnvVars(stsExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: "JVM_JMXMP_ENABLED", Value: "false"})
	addEnvVars(stsExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: "JVM_JMXMP_PORT", Value: "9099"})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithJvmSpecWithJmxmpEnabledWithPort(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		JVM: &coh.JVMSpec{
			Jmxmp: &coh.JvmJmxmpSpec{
				Enabled: boolPtr(true),
				Port:    int32Ptr(1234),
			},
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	addEnvVars(stsExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: "JVM_JMXMP_ENABLED", Value: "true"})
	addEnvVars(stsExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: "JVM_JMXMP_PORT", Value: "1234"})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}
