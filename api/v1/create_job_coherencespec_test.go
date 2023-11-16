/*
 * Copyright (c) 2023, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1_test

import (
	. "github.com/onsi/gomega"
	coh "github.com/oracle/coherence-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/intstr"
	"testing"
)

func TestCreateJobWithCoherenceSpecEmpty(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		Coherence: &coh.CoherenceSpec{},
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	jobExpected := createMinimalExpectedJob(deployment)

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}

func TestCreateJobWithImage(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		Image: stringPtr("coherence:1.0"),
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	jobExpected := createMinimalExpectedJob(deployment)
	jobExpected.Spec.Template.Spec.Containers[0].Image = "coherence:1.0"

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}

func TestCreateJobWithImagePullPolicy(t *testing.T) {
	policy := corev1.PullAlways
	spec := coh.CoherenceResourceSpec{
		ImagePullPolicy: &policy,
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	jobExpected := createMinimalExpectedJob(deployment)
	jobExpected.Spec.Template.Spec.Containers[0].ImagePullPolicy = policy

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}

func TestCreateJobWithCoherenceSpecWithStorageEnabledTrue(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		Coherence: &coh.CoherenceSpec{
			StorageEnabled: boolPtr(true),
		},
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	jobExpected := createMinimalExpectedJob(deployment)
	addEnvVarsToJob(jobExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: coh.EnvVarCohStorage, Value: "true"})

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}

func TestCreateJobWithCoherenceLocalPort(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		Coherence: &coh.CoherenceSpec{
			LocalPort: int32Ptr(1234),
		},
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	jobExpected := createMinimalExpectedJob(deployment)
	addEnvVarsToJob(jobExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: coh.EnvVarCoherenceLocalPort, Value: "1234"})

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}

func TestCreateJobWithCoherenceLocalPortAdjustTrue(t *testing.T) {
	lpa := intstr.FromString("true")
	spec := coh.CoherenceResourceSpec{
		Coherence: &coh.CoherenceSpec{
			LocalPortAdjust: &lpa,
		},
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	jobExpected := createMinimalExpectedJob(deployment)
	addEnvVarsToJob(jobExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: coh.EnvVarCoherenceLocalPortAdjust, Value: "true"})

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}

func TestCreateJobWithCoherenceLocalPortAdjustFalse(t *testing.T) {
	lpa := intstr.FromString("false")
	spec := coh.CoherenceResourceSpec{
		Coherence: &coh.CoherenceSpec{
			LocalPortAdjust: &lpa,
		},
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	jobExpected := createMinimalExpectedJob(deployment)
	addEnvVarsToJob(jobExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: coh.EnvVarCoherenceLocalPortAdjust, Value: "false"})

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}

func TestCreateJobWithCoherenceLocalPortAdjust(t *testing.T) {
	lpa := intstr.FromInt32(9876)
	spec := coh.CoherenceResourceSpec{
		Coherence: &coh.CoherenceSpec{
			LocalPortAdjust: &lpa,
		},
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	jobExpected := createMinimalExpectedJob(deployment)
	addEnvVarsToJob(jobExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: coh.EnvVarCoherenceLocalPortAdjust, Value: "9876"})

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}

func TestCreateJobWithCoherenceSpecWithStorageEnabledFalse(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		Coherence: &coh.CoherenceSpec{
			StorageEnabled: boolPtr(false),
		},
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	jobExpected := createMinimalExpectedJob(deployment)
	addEnvVarsToJob(jobExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: coh.EnvVarCohStorage, Value: "false"})

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}

func TestCreateJobWithCoherenceSpecWithCacheConfig(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		Coherence: &coh.CoherenceSpec{
			CacheConfig: stringPtr("test-config.xml"),
		},
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	jobExpected := createMinimalExpectedJob(deployment)
	addEnvVarsToJob(jobExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: coh.EnvVarCohCacheConfig, Value: "test-config.xml"})

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}

func TestCreateJobWithCoherenceSpecWithOverrideConfig(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		Coherence: &coh.CoherenceSpec{
			OverrideConfig: stringPtr("test-override.xml"),
		},
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	jobExpected := createMinimalExpectedJob(deployment)
	addEnvVarsToJob(jobExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: coh.EnvVarCohOverride, Value: "test-override.xml"})

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}

func TestCreateJobWithCoherenceSpecWithLogLevel(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		Coherence: &coh.CoherenceSpec{
			LogLevel: int32Ptr(9),
		},
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	jobExpected := createMinimalExpectedJob(deployment)
	addEnvVarsToJob(jobExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: coh.EnvVarCohLogLevel, Value: "9"})

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}

func TestCreateJobWithCoherenceSpecWithExcludeFromWKATrue(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		Coherence: &coh.CoherenceSpec{
			ExcludeFromWKA: boolPtr(true),
		},
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	jobExpected := createMinimalExpectedJob(deployment)
	jobExpected.Spec.Template.Labels[coh.LabelCoherenceWKAMember] = "false"

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}

func TestCreateJobWithCoherenceSpecWithExcludeFromWKAFalse(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		Coherence: &coh.CoherenceSpec{
			ExcludeFromWKA: boolPtr(false),
		},
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	jobExpected := createMinimalExpectedJob(deployment)
	jobExpected.Spec.Template.Labels[coh.LabelCoherenceWKAMember] = "true"

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}

func TestCreateJobWithCoherenceSpecWithTracingRatio(t *testing.T) {

	ratio := resource.MustParse("0.0005")
	spec := coh.CoherenceResourceSpec{
		Coherence: &coh.CoherenceSpec{
			Tracing: &coh.CoherenceTracingSpec{
				Ratio: &ratio,
			},
		},
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	jobExpected := createMinimalExpectedJob(deployment)
	addEnvVarsToJob(jobExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: coh.EnvVarCohTracingRatio, Value: "500u"})

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}

func TestCreateJobWithCoherenceSpecWithIpMonitorDefault(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		Coherence: &coh.CoherenceSpec{
			EnableIPMonitor: nil,
		},
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	jobExpected := createMinimalExpectedJob(deployment)

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}

func TestCreateJobWithCoherenceSpecWithIpMonitorDisabled(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		Coherence: &coh.CoherenceSpec{
			EnableIPMonitor: boolPtr(false),
		},
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	jobExpected := createMinimalExpectedJob(deployment)

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}

func TestCreateJobWithCoherenceSpecWithIpMonitorEnabled(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		Coherence: &coh.CoherenceSpec{
			EnableIPMonitor: boolPtr(true),
		},
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	jobExpected := createMinimalExpectedJob(deployment)
	addEnvVarsToJob(jobExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: coh.EnvVarEnableIPMonitor, Value: "TRUE"})

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}

func TestCreateJobWithCoherenceSpecWithWkaSameNamespace(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		Coherence: &coh.CoherenceSpec{
			WKA: &coh.CoherenceWKASpec{
				Deployment: "storage",
				Namespace:  "",
			},
		},
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	jobExpected := createMinimalExpectedJob(deployment)
	addEnvVarsToJob(jobExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: coh.EnvVarCohWka, Value: deployment.GetWKA()})
	jobExpected.Spec.Template.Labels[coh.LabelCoherenceWKAMember] = "false"

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}

func TestCreateJobWithCoherenceSpecWithWkaDifferentNamespace(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		Coherence: &coh.CoherenceSpec{
			WKA: &coh.CoherenceWKASpec{
				Addresses:  []string{},
				Deployment: "storage",
				Namespace:  "data",
			},
		},
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	jobExpected := createMinimalExpectedJob(deployment)
	expectedWka := deployment.GetWKA()
	addEnvVarsToJob(jobExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: coh.EnvVarCohWka, Value: expectedWka})
	jobExpected.Spec.Template.Labels[coh.LabelCoherenceWKAMember] = "false"

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}

func TestCreateJobWithCoherenceSpecWithWkaAddress(t *testing.T) {
	g := NewGomegaWithT(t)

	spec := coh.CoherenceResourceSpec{
		Coherence: &coh.CoherenceSpec{
			WKA: &coh.CoherenceWKASpec{
				Addresses:  []string{"storage.foo.bar.local"},
				Namespace:  "data",
				Deployment: "bad-cluster",
			},
		},
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	jobExpected := createMinimalExpectedJob(deployment)
	expectedWka := "storage.foo.bar.local"
	g.Expect(deployment.GetWKA()).To(Equal(expectedWka))
	addEnvVarsToJob(jobExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: coh.EnvVarCohWka, Value: expectedWka})
	jobExpected.Spec.Template.Labels[coh.LabelCoherenceWKAMember] = "false"

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}

func TestCreateJobWithCoherenceSpecWithMultipleWkaAddresses(t *testing.T) {
	g := NewGomegaWithT(t)

	spec := coh.CoherenceResourceSpec{
		Coherence: &coh.CoherenceSpec{
			WKA: &coh.CoherenceWKASpec{
				Addresses:  []string{"storage.one.bar.local", "storage.two.bar.local"},
				Namespace:  "data",
				Deployment: "bad-cluster",
			},
		},
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	jobExpected := createMinimalExpectedJob(deployment)
	expectedWka := "storage.one.bar.local,storage.two.bar.local"
	g.Expect(deployment.GetWKA()).To(Equal(expectedWka))
	addEnvVarsToJob(jobExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: coh.EnvVarCohWka, Value: expectedWka})
	jobExpected.Spec.Template.Labels[coh.LabelCoherenceWKAMember] = "false"

	// assert that the Job is as expected
	assertJobCreation(t, deployment, jobExpected)
}
