/*
 * Copyright (c) 2020, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1_test

import (
	coh "github.com/oracle/coherence-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	"testing"
)

func TestCreateStatefulSetWithNetworkSpec(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		Network: &coh.NetworkSpec{},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithNetworkSpecWithDNSConfigNameServers(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		Network: &coh.NetworkSpec{
			DNSConfig: &coh.PodDNSConfig{
				Nameservers: []string{"one", "two"},
			},
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	stsExpected.Spec.Template.Spec.DNSConfig = &corev1.PodDNSConfig{
		Nameservers: []string{"one", "two"},
	}

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithNetworkSpecWithDNSConfigEmptyNameServers(t *testing.T) {

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
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithNetworkSpecWithDNSConfigSearches(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		Network: &coh.NetworkSpec{
			DNSConfig: &coh.PodDNSConfig{
				Searches: []string{"one", "two"},
			},
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	stsExpected.Spec.Template.Spec.DNSConfig = &corev1.PodDNSConfig{
		Searches: []string{"one", "two"},
	}

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithNetworkSpecWithDNSConfigEmptySearches(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		Network: &coh.NetworkSpec{
			DNSConfig: &coh.PodDNSConfig{
				Searches: []string{},
				Options:  nil,
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

func TestCreateStatefulSetWithNetworkSpecWithDNSConfigOptions(t *testing.T) {

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
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	stsExpected.Spec.Template.Spec.DNSConfig = &corev1.PodDNSConfig{
		Options: []corev1.PodDNSConfigOption{
			{
				Name:  "Foo",
				Value: stringPtr("Bar"),
			},
		},
	}

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithNetworkSpecWithDNSConfigEmptyOptions(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		Network: &coh.NetworkSpec{
			DNSConfig: &coh.PodDNSConfig{
				Options: []corev1.PodDNSConfigOption{},
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

func TestCreateStatefulSetWithNetworkSpecWithDNSPolicy(t *testing.T) {

	policy := corev1.DNSClusterFirst

	spec := coh.CoherenceResourceSpec{
		Network: &coh.NetworkSpec{
			DNSPolicy: &policy,
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	stsExpected.Spec.Template.Spec.DNSPolicy = policy

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithNetworkSpecWithHostAliases(t *testing.T) {

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
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	stsExpected.Spec.Template.Spec.HostAliases = aliases

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithNetworkSpecWithHostNetworkFalse(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		Network: &coh.NetworkSpec{
			HostNetwork: boolPtr(false),
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	stsExpected.Spec.Template.Spec.HostNetwork = false

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithNetworkSpecWithHostNetworkTrue(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		Network: &coh.NetworkSpec{
			HostNetwork: boolPtr(true),
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	stsExpected.Spec.Template.Spec.HostNetwork = true

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithNetworkSpecWithHostname(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		Network: &coh.NetworkSpec{
			Hostname: stringPtr("foo.com"),
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	stsExpected.Spec.Template.Spec.Hostname = "foo.com"

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}
