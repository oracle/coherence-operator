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

func TestCreateStatefulSetWithPortsEmpty(t *testing.T) {

	spec := coh.CoherenceDeploymentSpec{
		Ports: []coh.NamedPortSpec{},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithPortsWithOneAdditionalPort(t *testing.T) {

	protocol := corev1.ProtocolUDP

	spec := coh.CoherenceDeploymentSpec{
		Ports: []coh.NamedPortSpec{
			{
				Name: "test-port-one",
				PortSpec: coh.PortSpec{
					Port:     9876,
					Protocol: &protocol,
					HostPort: int32Ptr(1234),
					HostIP:   stringPtr("10.10.1.0"),
				},
			},
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	addPorts(stsExpected, coh.ContainerNameCoherence, corev1.ContainerPort{
		Name:          "test-port-one",
		ContainerPort: 9876,
		HostPort:      1234,
		Protocol:      protocol,
		HostIP:        "10.10.1.0",
	})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithPortsWithTwoAdditionalPorts(t *testing.T) {

	protocolOne := corev1.ProtocolUDP
	protocolTwo := corev1.ProtocolSCTP

	spec := coh.CoherenceDeploymentSpec{
		Ports: []coh.NamedPortSpec{
			{
				Name: "test-port-one",
				PortSpec: coh.PortSpec{
					Port:     9876,
					Protocol: &protocolOne,
					HostPort: int32Ptr(1234),
					HostIP:   stringPtr("10.10.1.0"),
				},
			},
			{
				Name: "test-port-two",
				PortSpec: coh.PortSpec{
					Port:     5678,
					Protocol: &protocolTwo,
					HostPort: int32Ptr(7654),
					HostIP:   stringPtr("10.10.2.0"),
				},
			},
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	addPorts(stsExpected, coh.ContainerNameCoherence,
		corev1.ContainerPort{
			Name:          "test-port-one",
			ContainerPort: 9876,
			HostPort:      1234,
			Protocol:      protocolOne,
			HostIP:        "10.10.1.0",
		},
		corev1.ContainerPort{
			Name:          "test-port-two",
			ContainerPort: 5678,
			HostPort:      7654,
			Protocol:      protocolTwo,
			HostIP:        "10.10.2.0",
		})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}
