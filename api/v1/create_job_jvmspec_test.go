/*
 * Copyright (c) 2020, 2025, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1_test

import (
	"testing"

	coh "github.com/oracle/coherence-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
)

func TestCreateJobWithEmptyJvmSpec(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		JVM: &coh.JVMSpec{},
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	jobExpected := createMinimalExpectedJob(deployment)

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}

func TestCreateJobWithJvmSpecWithEmptyArgs(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		JVM: &coh.JVMSpec{
			Args: []string{},
		},
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	jobExpected := createMinimalExpectedJob(deployment)

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}

func TestCreateJobWithJvmSpecWithArgs(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		JVM: &coh.JVMSpec{
			Args: []string{"argOne", "argTwo"},
		},
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	jobExpected := createMinimalExpectedJob(deployment)
	addEnvVarsToAllJobContainers(jobExpected, corev1.EnvVar{Name: coh.EnvVarJvmArgs, Value: "argOne argTwo"})

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}

func TestCreateJobWithJvmSpecWithEmptyClasspath(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		JVM: &coh.JVMSpec{
			Classpath: []string{},
		},
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	jobExpected := createMinimalExpectedJob(deployment)

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}

func TestCreateJobWithJvmSpecWithClasspath(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		JVM: &coh.JVMSpec{
			Classpath: []string{"/foo", "/bar"},
		},
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	jobExpected := createMinimalExpectedJob(deployment)
	addEnvVarsToAllJobContainers(jobExpected, corev1.EnvVar{Name: coh.EnvVarJvmExtraClasspath, Value: "/foo:/bar"})

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}

func TestCreateJobWithJvmSpecWithUseContainerLimitsTrue(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		JVM: &coh.JVMSpec{
			UseContainerLimits: boolPtr(true),
		},
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	jobExpected := createMinimalExpectedJob(deployment)
	addEnvVarsToAllJobContainers(jobExpected, corev1.EnvVar{Name: coh.EnvVarJvmUseContainerLimits, Value: "true"})

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}

func TestCreateJobWithJvmSpecWithUseContainerLimitsFalse(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		JVM: &coh.JVMSpec{
			UseContainerLimits: boolPtr(false),
		},
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	jobExpected := createMinimalExpectedJob(deployment)
	addEnvVarsToAllJobContainers(jobExpected, corev1.EnvVar{Name: coh.EnvVarJvmUseContainerLimits, Value: "false"})

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}

func TestCreateJobWithJvmSpecWithDebugEnabledFalse(t *testing.T) {

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
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	jobExpected := createMinimalExpectedJob(deployment)

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}

func TestCreateJobWithJvmSpecWithDebugEnabledTrueSuspendTrue(t *testing.T) {

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
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	jobExpected := createMinimalExpectedJob(deployment)
	addEnvVarsToAllJobContainers(jobExpected, corev1.EnvVar{Name: "JVM_DEBUG_ENABLED", Value: "true"})
	addEnvVarsToAllJobContainers(jobExpected, corev1.EnvVar{Name: "JVM_DEBUG_SUSPEND", Value: "true"})
	addEnvVarsToAllJobContainers(jobExpected, corev1.EnvVar{Name: "JVM_DEBUG_ATTACH", Value: "10.10.10.10:5001"})
	addEnvVarsToAllJobContainers(jobExpected, corev1.EnvVar{Name: "JVM_DEBUG_PORT", Value: "1234"})
	// add the expected debug port
	addPortsForJob(jobExpected, coh.ContainerNameCoherence, corev1.ContainerPort{
		Name:          coh.PortNameDebug,
		ContainerPort: 1234,
	})

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}

func TestCreateJobWithJvmSpecWithDebugEnabledTrueSuspendFalse(t *testing.T) {

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
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	jobExpected := createMinimalExpectedJob(deployment)
	addEnvVarsToAllJobContainers(jobExpected, corev1.EnvVar{Name: "JVM_DEBUG_ENABLED", Value: "true"})
	addEnvVarsToAllJobContainers(jobExpected, corev1.EnvVar{Name: "JVM_DEBUG_ATTACH", Value: "10.10.10.10:5001"})
	addEnvVarsToAllJobContainers(jobExpected, corev1.EnvVar{Name: "JVM_DEBUG_PORT", Value: "1234"})
	// add the expected debug port
	addPortsForJob(jobExpected, coh.ContainerNameCoherence, corev1.ContainerPort{
		Name:          coh.PortNameDebug,
		ContainerPort: 1234,
	})

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}

func TestCreateJobWithJvmSpecWithGarbageCollector(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		JVM: &coh.JVMSpec{
			Gc: &coh.JvmGarbageCollectorSpec{
				Collector: stringPtr("G1"),
			},
		},
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	jobExpected := createMinimalExpectedJob(deployment)
	addEnvVarsToAllJobContainers(jobExpected, corev1.EnvVar{Name: "JVM_GC_COLLECTOR", Value: "G1"})

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}

func TestCreateJobWithJvmSpecWithGarbageCollectorArgs(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		JVM: &coh.JVMSpec{
			Gc: &coh.JvmGarbageCollectorSpec{
				Args: []string{"-XX:GC-ArgOne", "-XX:GC-ArgTwo"},
			},
		},
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	jobExpected := createMinimalExpectedJob(deployment)
	addEnvVarsToAllJobContainers(jobExpected, corev1.EnvVar{Name: "JVM_GC_ARGS", Value: "-XX:GC-ArgOne -XX:GC-ArgTwo"})

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}

func TestCreateJobWithJvmSpecWithDiagnosticsVolume(t *testing.T) {

	hostPath := &corev1.HostPathVolumeSource{Path: "/home/root/debug"}
	spec := coh.CoherenceResourceSpec{
		JVM: &coh.JVMSpec{
			DiagnosticsVolume: &corev1.VolumeSource{
				HostPath: hostPath,
			},
		},
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job with the specified JVM diagnostic volume
	jobExpected := createMinimalExpectedJob(deployment)
	coh.ReplaceVolumeInJob(jobExpected, corev1.Volume{
		Name: coh.VolumeNameJVM,
		VolumeSource: corev1.VolumeSource{
			HostPath: hostPath,
		},
	})

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}

func TestCreateJobWithJvmSpecWithMemorySettings(t *testing.T) {

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
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	jobExpected := createMinimalExpectedJob(deployment)
	addEnvVarsToAllJobContainers(jobExpected, corev1.EnvVar{Name: "JVM_HEAP_SIZE", Value: "5g"})
	addEnvVarsToAllJobContainers(jobExpected, corev1.EnvVar{Name: "JVM_STACK_SIZE", Value: "500m"})
	addEnvVarsToAllJobContainers(jobExpected, corev1.EnvVar{Name: "JVM_METASPACE_SIZE", Value: "1g"})
	addEnvVarsToAllJobContainers(jobExpected, corev1.EnvVar{Name: "JVM_DIRECT_MEMORY_SIZE", Value: "4g"})
	addEnvVarsToAllJobContainers(jobExpected, corev1.EnvVar{Name: "JVM_NATIVE_MEMORY_TRACKING", Value: "detail"})

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}

func TestCreateJobWithJvmSpecWithExitOnOomTrue(t *testing.T) {

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
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	jobExpected := createMinimalExpectedJob(deployment)
	addEnvVarsToAllJobContainers(jobExpected, corev1.EnvVar{Name: "JVM_OOM_EXIT", Value: "true"})

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}

func TestCreateJobWithJvmSpecWithExitOnOomFalse(t *testing.T) {

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
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	jobExpected := createMinimalExpectedJob(deployment)
	addEnvVarsToAllJobContainers(jobExpected, corev1.EnvVar{Name: "JVM_OOM_EXIT", Value: "false"})

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}

func TestCreateJobWithJvmSpecWithHeapDumpOnOomTrue(t *testing.T) {

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
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	jobExpected := createMinimalExpectedJob(deployment)
	addEnvVarsToAllJobContainers(jobExpected, corev1.EnvVar{Name: "JVM_OOM_HEAP_DUMP", Value: "true"})

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}

func TestCreateJobWithJvmSpecWithHeapDumpOnOomFalse(t *testing.T) {

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
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	jobExpected := createMinimalExpectedJob(deployment)
	addEnvVarsToAllJobContainers(jobExpected, corev1.EnvVar{Name: "JVM_OOM_HEAP_DUMP", Value: "false"})

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}
