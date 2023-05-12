/*
 * Copyright (c) 2023, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1_test

import (
	coh "github.com/oracle/coherence-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	"testing"
)

func TestCreateJobWithNetworkSpec(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		Network: &coh.NetworkSpec{},
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	stsExpected := createMinimalExpectedJob(deployment)

	// assert that the Job is as expected
	assertJobCreation(t, deployment, stsExpected)
}

func TestCreateJobWithNetworkSpecWithDNSConfigNameServers(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		Network: &coh.NetworkSpec{
			DNSConfig: &coh.PodDNSConfig{
				Nameservers: []string{"one", "two"},
			},
		},
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	stsExpected := createMinimalExpectedJob(deployment)
	stsExpected.Spec.Template.Spec.DNSConfig = &corev1.PodDNSConfig{
		Nameservers: []string{"one", "two"},
	}

	// assert that the Job is as expected
	assertJobCreation(t, deployment, stsExpected)
}

func TestCreateJobWithNetworkSpecWithDNSConfigEmptyNameServers(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		Network: &coh.NetworkSpec{
			DNSConfig: &coh.PodDNSConfig{
				Nameservers: []string{},
				Searches:    nil,
				Options:     nil,
			},
		},
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	stsExpected := createMinimalExpectedJob(deployment)

	// assert that the Job is as expected
	assertJobCreation(t, deployment, stsExpected)
}

func TestCreateJobWithNetworkSpecWithDNSConfigSearches(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		Network: &coh.NetworkSpec{
			DNSConfig: &coh.PodDNSConfig{
				Searches: []string{"one", "two"},
			},
		},
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	stsExpected := createMinimalExpectedJob(deployment)
	stsExpected.Spec.Template.Spec.DNSConfig = &corev1.PodDNSConfig{
		Searches: []string{"one", "two"},
	}

	// assert that the Job is as expected
	assertJobCreation(t, deployment, stsExpected)
}

func TestCreateJobWithNetworkSpecWithDNSConfigEmptySearches(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		Network: &coh.NetworkSpec{
			DNSConfig: &coh.PodDNSConfig{
				Searches: []string{},
				Options:  nil,
			},
		},
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	stsExpected := createMinimalExpectedJob(deployment)

	// assert that the Job is as expected
	assertJobCreation(t, deployment, stsExpected)
}

func TestCreateJobWithNetworkSpecWithDNSConfigOptions(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		Network: &coh.NetworkSpec{
			DNSConfig: &coh.PodDNSConfig{
				Options: []corev1.PodDNSConfigOption{
					{
						Name:  "Foo",
						Value: stringPtr("Bar"),
					},
				},
			},
		},
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	stsExpected := createMinimalExpectedJob(deployment)
	stsExpected.Spec.Template.Spec.DNSConfig = &corev1.PodDNSConfig{
		Options: []corev1.PodDNSConfigOption{
			{
				Name:  "Foo",
				Value: stringPtr("Bar"),
			},
		},
	}

	// assert that the Job is as expected
	assertJobCreation(t, deployment, stsExpected)
}

func TestCreateJobWithNetworkSpecWithDNSConfigEmptyOptions(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		Network: &coh.NetworkSpec{
			DNSConfig: &coh.PodDNSConfig{
				Options: []corev1.PodDNSConfigOption{},
			},
		},
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	stsExpected := createMinimalExpectedJob(deployment)

	// assert that the Job is as expected
	assertJobCreation(t, deployment, stsExpected)
}

func TestCreateJobWithNetworkSpecWithDNSPolicy(t *testing.T) {

	policy := corev1.DNSClusterFirst

	spec := coh.CoherenceResourceSpec{
		Network: &coh.NetworkSpec{
			DNSPolicy: &policy,
		},
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	stsExpected := createMinimalExpectedJob(deployment)
	stsExpected.Spec.Template.Spec.DNSPolicy = policy

	// assert that the Job is as expected
	assertJobCreation(t, deployment, stsExpected)
}

func TestCreateJobWithNetworkSpecWithHostAliases(t *testing.T) {

	aliases := []corev1.HostAlias{
		{
			IP:        "10.10.10.10",
			Hostnames: []string{"foo.com"},
		},
	}

	spec := coh.CoherenceResourceSpec{
		Network: &coh.NetworkSpec{
			HostAliases: aliases,
		},
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	stsExpected := createMinimalExpectedJob(deployment)
	stsExpected.Spec.Template.Spec.HostAliases = aliases

	// assert that the Job is as expected
	assertJobCreation(t, deployment, stsExpected)
}

func TestCreateJobWithNetworkSpecWithHostNetworkFalse(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		Network: &coh.NetworkSpec{
			HostNetwork: boolPtr(false),
		},
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	stsExpected := createMinimalExpectedJob(deployment)
	stsExpected.Spec.Template.Spec.HostNetwork = false

	// assert that the Job is as expected
	assertJobCreation(t, deployment, stsExpected)
}

func TestCreateJobWithNetworkSpecWithHostNetworkTrue(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		Network: &coh.NetworkSpec{
			HostNetwork: boolPtr(true),
		},
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	stsExpected := createMinimalExpectedJob(deployment)
	stsExpected.Spec.Template.Spec.HostNetwork = true

	// assert that the Job is as expected
	assertJobCreation(t, deployment, stsExpected)
}

func TestCreateJobWithNetworkSpecWithHostname(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		Network: &coh.NetworkSpec{
			Hostname: stringPtr("foo.com"),
		},
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	stsExpected := createMinimalExpectedJob(deployment)
	stsExpected.Spec.Template.Spec.Hostname = "foo.com"

	// assert that the Job is as expected
	assertJobCreation(t, deployment, stsExpected)
}

func TestCreateJobWithNetworkSpecWithSetHostnameAsFQDNTrue(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		Network: &coh.NetworkSpec{
			SetHostnameAsFQDN: boolPtr(true),
		},
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	stsExpected := createMinimalExpectedJob(deployment)
	stsExpected.Spec.Template.Spec.SetHostnameAsFQDN = boolPtr(true)

	// assert that the Job is as expected
	assertJobCreation(t, deployment, stsExpected)
}

func TestCreateJobWithNetworkSpecWithSetHostnameAsFQDNFalse(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		Network: &coh.NetworkSpec{
			SetHostnameAsFQDN: boolPtr(false),
		},
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	stsExpected := createMinimalExpectedJob(deployment)
	stsExpected.Spec.Template.Spec.SetHostnameAsFQDN = boolPtr(false)

	// assert that the Job is as expected
	assertJobCreation(t, deployment, stsExpected)
}

func TestCreateJobWithNetworkSpecWithSubdomain(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		Network: &coh.NetworkSpec{
			Subdomain: stringPtr("foo"),
		},
	}

	// Create the test deployment
	deployment := createTestCoherenceJob(spec)
	// Create expected Job
	stsExpected := createMinimalExpectedJob(deployment)
	stsExpected.Spec.Template.Spec.Subdomain = "foo"

	// assert that the Job is as expected
	assertJobCreation(t, deployment, stsExpected)
}
