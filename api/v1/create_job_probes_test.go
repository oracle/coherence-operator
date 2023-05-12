/*
 * Copyright (c) 2023, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1_test

import (
	coh "github.com/oracle/coherence-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"testing"
)

func TestCreateJobWithEmptyReadinessProbeSpec(t *testing.T) {
	probe := coh.ReadinessProbeSpec{}

	spec := coh.CoherenceResourceSpec{
		ReadinessProbe: &probe,
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)

	// Create expected probe
	probeExpected := spec.UpdateDefaultReadinessProbeAction(spec.CreateDefaultReadinessProbe())
	// Create expected Job
	stsExpected := createMinimalExpectedJob(deployment)
	stsExpected.Spec.Template.Spec.Containers[0].ReadinessProbe = probeExpected

	// assert that the Job is as expected
	assertJobCreation(t, deployment, stsExpected)
}

func TestCreateJobWithReadinessProbeSpec(t *testing.T) {
	probe := coh.ReadinessProbeSpec{
		InitialDelaySeconds: int32Ptr(10),
		TimeoutSeconds:      int32Ptr(20),
		PeriodSeconds:       int32Ptr(30),
		SuccessThreshold:    int32Ptr(40),
		FailureThreshold:    int32Ptr(50),
	}

	spec := coh.CoherenceResourceSpec{
		ReadinessProbe: &probe,
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	stsExpected := createMinimalExpectedJob(deployment)
	stsExpected.Spec.Template.Spec.Containers[0].ReadinessProbe = &corev1.Probe{
		ProbeHandler: corev1.ProbeHandler{
			Exec: nil,
			HTTPGet: &corev1.HTTPGetAction{
				Path:   coh.DefaultReadinessPath,
				Port:   intstr.FromInt(int(coh.DefaultHealthPort)),
				Scheme: "HTTP",
			},
			TCPSocket: nil,
		},
		InitialDelaySeconds: 10,
		TimeoutSeconds:      20,
		PeriodSeconds:       30,
		SuccessThreshold:    40,
		FailureThreshold:    50,
	}

	// assert that the Job is as expected
	assertJobCreation(t, deployment, stsExpected)
}

func TestCreateJobWithReadinessProbeSpecWithHttpGet(t *testing.T) {
	handler := &corev1.HTTPGetAction{
		Path: "/test/ready",
		Port: intstr.FromInt(1234),
	}

	probe := coh.ReadinessProbeSpec{
		ProbeHandler: coh.ProbeHandler{
			Exec:      nil,
			HTTPGet:   handler,
			TCPSocket: nil,
		},
		InitialDelaySeconds: int32Ptr(10),
		TimeoutSeconds:      int32Ptr(20),
		PeriodSeconds:       int32Ptr(30),
		SuccessThreshold:    int32Ptr(40),
		FailureThreshold:    int32Ptr(50),
	}

	spec := coh.CoherenceResourceSpec{
		ReadinessProbe: &probe,
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	stsExpected := createMinimalExpectedJob(deployment)
	stsExpected.Spec.Template.Spec.Containers[0].ReadinessProbe = &corev1.Probe{
		ProbeHandler: corev1.ProbeHandler{
			HTTPGet: handler,
		},
		InitialDelaySeconds: 10,
		TimeoutSeconds:      20,
		PeriodSeconds:       30,
		SuccessThreshold:    40,
		FailureThreshold:    50,
	}

	// assert that the Job is as expected
	assertJobCreation(t, deployment, stsExpected)
}

func TestCreateJobWithReadinessProbeSpecWithTCPSocket(t *testing.T) {
	handler := &corev1.TCPSocketAction{
		Port: intstr.FromInt(1234),
		Host: "foo.com",
	}

	probe := coh.ReadinessProbeSpec{
		ProbeHandler: coh.ProbeHandler{
			TCPSocket: handler,
		},
		InitialDelaySeconds: int32Ptr(10),
		TimeoutSeconds:      int32Ptr(20),
		PeriodSeconds:       int32Ptr(30),
		SuccessThreshold:    int32Ptr(40),
		FailureThreshold:    int32Ptr(50),
	}

	spec := coh.CoherenceResourceSpec{
		ReadinessProbe: &probe,
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	stsExpected := createMinimalExpectedJob(deployment)
	stsExpected.Spec.Template.Spec.Containers[0].ReadinessProbe = &corev1.Probe{
		ProbeHandler: corev1.ProbeHandler{
			TCPSocket: handler,
		},
		InitialDelaySeconds: 10,
		TimeoutSeconds:      20,
		PeriodSeconds:       30,
		SuccessThreshold:    40,
		FailureThreshold:    50,
	}

	// assert that the Job is as expected
	assertJobCreation(t, deployment, stsExpected)
}

func TestCreateJobWithReadinessProbeSpecWithExec(t *testing.T) {
	handler := &corev1.ExecAction{
		Command: []string{"exec", "something"},
	}

	probe := coh.ReadinessProbeSpec{
		ProbeHandler: coh.ProbeHandler{
			Exec: handler,
		},
		InitialDelaySeconds: int32Ptr(10),
		TimeoutSeconds:      int32Ptr(20),
		PeriodSeconds:       int32Ptr(30),
		SuccessThreshold:    int32Ptr(40),
		FailureThreshold:    int32Ptr(50),
	}

	spec := coh.CoherenceResourceSpec{
		ReadinessProbe: &probe,
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	stsExpected := createMinimalExpectedJob(deployment)
	stsExpected.Spec.Template.Spec.Containers[0].ReadinessProbe = &corev1.Probe{
		ProbeHandler: corev1.ProbeHandler{
			Exec: handler,
		},
		InitialDelaySeconds: 10,
		TimeoutSeconds:      20,
		PeriodSeconds:       30,
		SuccessThreshold:    40,
		FailureThreshold:    50,
	}

	// assert that the Job is as expected
	assertJobCreation(t, deployment, stsExpected)
}

func TestCreateJobWithEmptyLivenessProbeSpec(t *testing.T) {
	probe := coh.ReadinessProbeSpec{}

	spec := coh.CoherenceResourceSpec{
		LivenessProbe: &probe,
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)

	// Create expected probe
	probeExpected := spec.UpdateDefaultLivenessProbeAction(spec.CreateDefaultLivenessProbe())
	// Create expected Job
	stsExpected := createMinimalExpectedJob(deployment)
	stsExpected.Spec.Template.Spec.Containers[0].LivenessProbe = probeExpected

	// assert that the Job is as expected
	assertJobCreation(t, deployment, stsExpected)
}

func TestCreateJobWithLivenessProbeSpec(t *testing.T) {

	probe := coh.ReadinessProbeSpec{
		InitialDelaySeconds: int32Ptr(10),
		TimeoutSeconds:      int32Ptr(20),
		PeriodSeconds:       int32Ptr(30),
		SuccessThreshold:    int32Ptr(40),
		FailureThreshold:    int32Ptr(50),
	}

	spec := coh.CoherenceResourceSpec{
		LivenessProbe: &probe,
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	stsExpected := createMinimalExpectedJob(deployment)
	stsExpected.Spec.Template.Spec.Containers[0].LivenessProbe = &corev1.Probe{
		ProbeHandler: corev1.ProbeHandler{
			Exec: nil,
			HTTPGet: &corev1.HTTPGetAction{
				Path:   coh.DefaultLivenessPath,
				Port:   intstr.FromInt(int(coh.DefaultHealthPort)),
				Scheme: "HTTP",
			},
			TCPSocket: nil,
		},
		InitialDelaySeconds: 10,
		TimeoutSeconds:      20,
		PeriodSeconds:       30,
		SuccessThreshold:    40,
		FailureThreshold:    50,
	}

	// assert that the Job is as expected
	assertJobCreation(t, deployment, stsExpected)
}

func TestCreateJobWithLivenessProbeSpecWithHttpGet(t *testing.T) {

	handler := &corev1.HTTPGetAction{
		Path: "/test/ready",
		Port: intstr.FromInt(1234),
	}

	probe := coh.ReadinessProbeSpec{
		ProbeHandler: coh.ProbeHandler{
			Exec:      nil,
			HTTPGet:   handler,
			TCPSocket: nil,
		},
		InitialDelaySeconds: int32Ptr(10),
		TimeoutSeconds:      int32Ptr(20),
		PeriodSeconds:       int32Ptr(30),
		SuccessThreshold:    int32Ptr(40),
		FailureThreshold:    int32Ptr(50),
	}

	spec := coh.CoherenceResourceSpec{
		LivenessProbe: &probe,
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	stsExpected := createMinimalExpectedJob(deployment)
	stsExpected.Spec.Template.Spec.Containers[0].LivenessProbe = &corev1.Probe{
		ProbeHandler: corev1.ProbeHandler{
			HTTPGet: handler,
		},
		InitialDelaySeconds: 10,
		TimeoutSeconds:      20,
		PeriodSeconds:       30,
		SuccessThreshold:    40,
		FailureThreshold:    50,
	}

	// assert that the Job is as expected
	assertJobCreation(t, deployment, stsExpected)
}

func TestCreateJobWithLivenessProbeSpecWithTCPSocket(t *testing.T) {

	handler := &corev1.TCPSocketAction{
		Port: intstr.FromInt(1234),
		Host: "foo.com",
	}

	probe := coh.ReadinessProbeSpec{
		ProbeHandler: coh.ProbeHandler{
			TCPSocket: handler,
		},
		InitialDelaySeconds: int32Ptr(10),
		TimeoutSeconds:      int32Ptr(20),
		PeriodSeconds:       int32Ptr(30),
		SuccessThreshold:    int32Ptr(40),
		FailureThreshold:    int32Ptr(50),
	}

	spec := coh.CoherenceResourceSpec{
		LivenessProbe: &probe,
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	stsExpected := createMinimalExpectedJob(deployment)
	stsExpected.Spec.Template.Spec.Containers[0].LivenessProbe = &corev1.Probe{
		ProbeHandler: corev1.ProbeHandler{
			TCPSocket: handler,
		},
		InitialDelaySeconds: 10,
		TimeoutSeconds:      20,
		PeriodSeconds:       30,
		SuccessThreshold:    40,
		FailureThreshold:    50,
	}

	// assert that the Job is as expected
	assertJobCreation(t, deployment, stsExpected)
}

func TestCreateJobWithLivenessProbeSpecWithExec(t *testing.T) {

	handler := &corev1.ExecAction{
		Command: []string{"exec", "something"},
	}

	probe := coh.ReadinessProbeSpec{
		ProbeHandler: coh.ProbeHandler{
			Exec: handler,
		},
		InitialDelaySeconds: int32Ptr(10),
		TimeoutSeconds:      int32Ptr(20),
		PeriodSeconds:       int32Ptr(30),
		SuccessThreshold:    int32Ptr(40),
		FailureThreshold:    int32Ptr(50),
	}

	spec := coh.CoherenceResourceSpec{
		LivenessProbe: &probe,
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	stsExpected := createMinimalExpectedJob(deployment)
	stsExpected.Spec.Template.Spec.Containers[0].LivenessProbe = &corev1.Probe{
		ProbeHandler: corev1.ProbeHandler{
			Exec: handler,
		},
		InitialDelaySeconds: 10,
		TimeoutSeconds:      20,
		PeriodSeconds:       30,
		SuccessThreshold:    40,
		FailureThreshold:    50,
	}

	// assert that the Job is as expected
	assertJobCreation(t, deployment, stsExpected)
}

func TestCreateJobWithNilStartupProbeSpec(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		StartupProbe: nil,
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)

	// Create expected Job
	stsExpected := createMinimalExpectedJob(deployment)

	// assert that the Job is as expected
	assertJobCreation(t, deployment, stsExpected)
}

func TestCreateJobWithEmptyStartupProbeSpec(t *testing.T) {

	probe := coh.ReadinessProbeSpec{}

	spec := coh.CoherenceResourceSpec{
		StartupProbe: &probe,
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)

	// Create expected probe
	probeExpected := spec.UpdateDefaultLivenessProbeAction(spec.CreateDefaultLivenessProbe())

	// Create expected Job
	stsExpected := createMinimalExpectedJob(deployment)
	stsExpected.Spec.Template.Spec.Containers[0].StartupProbe = probeExpected

	// assert that the Job is as expected
	assertJobCreation(t, deployment, stsExpected)
}

func TestCreateJobWithStartupProbeSpec(t *testing.T) {

	probe := coh.ReadinessProbeSpec{
		InitialDelaySeconds: int32Ptr(10),
		TimeoutSeconds:      int32Ptr(20),
		PeriodSeconds:       int32Ptr(30),
		SuccessThreshold:    int32Ptr(40),
		FailureThreshold:    int32Ptr(50),
	}

	spec := coh.CoherenceResourceSpec{
		StartupProbe: &probe,
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	stsExpected := createMinimalExpectedJob(deployment)
	stsExpected.Spec.Template.Spec.Containers[0].StartupProbe = &corev1.Probe{
		ProbeHandler: corev1.ProbeHandler{
			Exec: nil,
			HTTPGet: &corev1.HTTPGetAction{
				Path:   coh.DefaultLivenessPath,
				Port:   intstr.FromInt(int(coh.DefaultHealthPort)),
				Scheme: "HTTP",
			},
			TCPSocket: nil,
		},
		InitialDelaySeconds: 10,
		TimeoutSeconds:      20,
		PeriodSeconds:       30,
		SuccessThreshold:    40,
		FailureThreshold:    50,
	}

	// assert that the Job is as expected
	assertJobCreation(t, deployment, stsExpected)
}

func TestCreateJobWithStartupProbeSpecWithHttpGet(t *testing.T) {

	handler := &corev1.HTTPGetAction{
		Path: "/test/ready",
		Port: intstr.FromInt(1234),
	}

	probe := coh.ReadinessProbeSpec{
		ProbeHandler: coh.ProbeHandler{
			Exec:      nil,
			HTTPGet:   handler,
			TCPSocket: nil,
		},
		InitialDelaySeconds: int32Ptr(10),
		TimeoutSeconds:      int32Ptr(20),
		PeriodSeconds:       int32Ptr(30),
		SuccessThreshold:    int32Ptr(40),
		FailureThreshold:    int32Ptr(50),
	}

	spec := coh.CoherenceResourceSpec{
		StartupProbe: &probe,
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	stsExpected := createMinimalExpectedJob(deployment)
	stsExpected.Spec.Template.Spec.Containers[0].StartupProbe = &corev1.Probe{
		ProbeHandler: corev1.ProbeHandler{
			HTTPGet: handler,
		},
		InitialDelaySeconds: 10,
		TimeoutSeconds:      20,
		PeriodSeconds:       30,
		SuccessThreshold:    40,
		FailureThreshold:    50,
	}

	// assert that the Job is as expected
	assertJobCreation(t, deployment, stsExpected)
}

func TestCreateJobWithStartupProbeSpecWithTCPSocket(t *testing.T) {

	handler := &corev1.TCPSocketAction{
		Port: intstr.FromInt(1234),
		Host: "foo.com",
	}

	probe := coh.ReadinessProbeSpec{
		ProbeHandler: coh.ProbeHandler{
			TCPSocket: handler,
		},
		InitialDelaySeconds: int32Ptr(10),
		TimeoutSeconds:      int32Ptr(20),
		PeriodSeconds:       int32Ptr(30),
		SuccessThreshold:    int32Ptr(40),
		FailureThreshold:    int32Ptr(50),
	}

	spec := coh.CoherenceResourceSpec{
		StartupProbe: &probe,
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	stsExpected := createMinimalExpectedJob(deployment)
	stsExpected.Spec.Template.Spec.Containers[0].StartupProbe = &corev1.Probe{
		ProbeHandler: corev1.ProbeHandler{
			TCPSocket: handler,
		},
		InitialDelaySeconds: 10,
		TimeoutSeconds:      20,
		PeriodSeconds:       30,
		SuccessThreshold:    40,
		FailureThreshold:    50,
	}

	// assert that the Job is as expected
	assertJobCreation(t, deployment, stsExpected)
}

func TestCreateJobWithStartupProbeSpecWithExec(t *testing.T) {

	handler := &corev1.ExecAction{
		Command: []string{"exec", "something"},
	}

	probe := coh.ReadinessProbeSpec{
		ProbeHandler: coh.ProbeHandler{
			Exec: handler,
		},
		InitialDelaySeconds: int32Ptr(10),
		TimeoutSeconds:      int32Ptr(20),
		PeriodSeconds:       int32Ptr(30),
		SuccessThreshold:    int32Ptr(40),
		FailureThreshold:    int32Ptr(50),
	}

	spec := coh.CoherenceResourceSpec{
		StartupProbe: &probe,
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	stsExpected := createMinimalExpectedJob(deployment)
	stsExpected.Spec.Template.Spec.Containers[0].StartupProbe = &corev1.Probe{
		ProbeHandler: corev1.ProbeHandler{
			Exec: handler,
		},
		InitialDelaySeconds: 10,
		TimeoutSeconds:      20,
		PeriodSeconds:       30,
		SuccessThreshold:    40,
		FailureThreshold:    50,
	}

	// assert that the Job is as expected
	assertJobCreation(t, deployment, stsExpected)
}
