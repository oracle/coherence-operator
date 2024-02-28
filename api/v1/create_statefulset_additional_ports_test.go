/*
 * Copyright (c) 2020, 2024, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1_test

import (
	. "github.com/onsi/gomega"
	coh "github.com/oracle/coherence-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/ptr"
	"testing"
)

func TestCreateStatefulSetWithPortsEmpty(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
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

	spec := coh.CoherenceResourceSpec{
		Ports: []coh.NamedPortSpec{
			{
				Name:     "test-port-one",
				Port:     9876,
				Protocol: &protocol,
				HostPort: int32Ptr(1234),
				HostIP:   stringPtr("10.10.1.0"),
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

	spec := coh.CoherenceResourceSpec{
		Ports: []coh.NamedPortSpec{
			{
				Name:     "test-port-one",
				Port:     9876,
				Protocol: &protocolOne,
				HostPort: int32Ptr(1234),
				HostIP:   stringPtr("10.10.1.0"),
			},
			{
				Name:     "test-port-two",
				Port:     5678,
				Protocol: &protocolTwo,
				HostPort: int32Ptr(7654),
				HostIP:   stringPtr("10.10.2.0"),
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

func TestCreateStatefulSetWithMetricsPortWhenNoPortValueSpecified(t *testing.T) {
	g := NewGomegaWithT(t)

	protocol := corev1.ProtocolTCP
	np := coh.NamedPortSpec{
		Name: coh.PortNameMetrics,
	}

	spec := coh.CoherenceResourceSpec{
		Ports: []coh.NamedPortSpec{np},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	addPorts(stsExpected, coh.ContainerNameCoherence, corev1.ContainerPort{
		Name:          coh.PortNameMetrics,
		ContainerPort: coh.DefaultMetricsPort,
		Protocol:      protocol,
	})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)

	svc := np.CreateService(deployment)
	g.Expect(svc).NotTo(BeNil())
	g.Expect(len(svc.Spec.Ports)).To(Equal(1))
	g.Expect(svc.Spec.Ports[0].Port).To(Equal(coh.DefaultMetricsPort))
	g.Expect(svc.Spec.Ports[0].TargetPort.IntVal).To(Equal(coh.DefaultMetricsPort))
}

func TestCreateStatefulSetWithMetricsPortWhenMetricsPortValueSpecified(t *testing.T) {
	g := NewGomegaWithT(t)

	protocol := corev1.ProtocolTCP
	np := coh.NamedPortSpec{
		Name: coh.PortNameMetrics,
	}

	spec := coh.CoherenceResourceSpec{
		Coherence: &coh.CoherenceSpec{
			Metrics: &coh.PortSpecWithSSL{
				Port: ptr.To(int32(1234)),
			},
		},
		Ports: []coh.NamedPortSpec{np},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	addPorts(stsExpected, coh.ContainerNameCoherence, corev1.ContainerPort{
		Name:          coh.PortNameMetrics,
		ContainerPort: 1234,
		Protocol:      protocol,
	})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)

	svc := np.CreateService(deployment)
	g.Expect(svc).NotTo(BeNil())
	g.Expect(len(svc.Spec.Ports)).To(Equal(1))
	g.Expect(svc.Spec.Ports[0].Port).To(Equal(int32(1234)))
	g.Expect(svc.Spec.Ports[0].TargetPort.IntVal).To(Equal(int32(1234)))
}

func TestCreateStatefulSetWithMetricsPortAndServicePortWhenNoPortValueSpecified(t *testing.T) {
	g := NewGomegaWithT(t)

	protocol := corev1.ProtocolTCP
	np := coh.NamedPortSpec{
		Name: coh.PortNameMetrics,
		Service: &coh.ServiceSpec{
			Port: ptr.To(int32(1234)),
		},
	}

	spec := coh.CoherenceResourceSpec{
		Ports: []coh.NamedPortSpec{np},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	addPorts(stsExpected, coh.ContainerNameCoherence, corev1.ContainerPort{
		Name:          coh.PortNameMetrics,
		ContainerPort: coh.DefaultMetricsPort,
		Protocol:      protocol,
	})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)

	svc := np.CreateService(deployment)
	g.Expect(svc).NotTo(BeNil())
	g.Expect(len(svc.Spec.Ports)).To(Equal(1))
	g.Expect(svc.Spec.Ports[0].Port).To(Equal(int32(1234)))
	g.Expect(svc.Spec.Ports[0].TargetPort.IntVal).To(Equal(coh.DefaultMetricsPort))
}

func TestCreateStatefulSetWithMetricsPortAndServicePortWhenMetricsPortValueSpecified(t *testing.T) {
	g := NewGomegaWithT(t)

	protocol := corev1.ProtocolTCP
	np := coh.NamedPortSpec{
		Name: coh.PortNameMetrics,
		Service: &coh.ServiceSpec{
			Port: ptr.To(int32(1234)),
		},
	}

	spec := coh.CoherenceResourceSpec{
		Coherence: &coh.CoherenceSpec{
			Metrics: &coh.PortSpecWithSSL{
				Port: ptr.To(int32(9876)),
			},
		},
		Ports: []coh.NamedPortSpec{np},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	addPorts(stsExpected, coh.ContainerNameCoherence, corev1.ContainerPort{
		Name:          coh.PortNameMetrics,
		ContainerPort: 9876,
		Protocol:      protocol,
	})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)

	svc := np.CreateService(deployment)
	g.Expect(svc).NotTo(BeNil())
	g.Expect(len(svc.Spec.Ports)).To(Equal(1))
	g.Expect(svc.Spec.Ports[0].Port).To(Equal(int32(1234)))
	g.Expect(svc.Spec.Ports[0].TargetPort.IntVal).To(Equal(int32(9876)))
}

func TestCreateStatefulSetWithManagementPortWhenNoPortValueSpecified(t *testing.T) {
	g := NewGomegaWithT(t)

	protocol := corev1.ProtocolTCP
	np := coh.NamedPortSpec{
		Name: coh.PortNameManagement,
	}

	spec := coh.CoherenceResourceSpec{
		Ports: []coh.NamedPortSpec{np},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	addPorts(stsExpected, coh.ContainerNameCoherence, corev1.ContainerPort{
		Name:          coh.PortNameManagement,
		ContainerPort: coh.DefaultManagementPort,
		Protocol:      protocol,
	})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)

	svc := np.CreateService(deployment)
	g.Expect(svc).NotTo(BeNil())
	g.Expect(len(svc.Spec.Ports)).To(Equal(1))
	g.Expect(svc.Spec.Ports[0].Port).To(Equal(coh.DefaultManagementPort))
	g.Expect(svc.Spec.Ports[0].TargetPort.IntVal).To(Equal(coh.DefaultManagementPort))
}

func TestCreateStatefulSetWithManagementPortWhenManagementPortValueSpecified(t *testing.T) {
	g := NewGomegaWithT(t)

	protocol := corev1.ProtocolTCP
	np := coh.NamedPortSpec{
		Name: coh.PortNameManagement,
	}

	spec := coh.CoherenceResourceSpec{
		Coherence: &coh.CoherenceSpec{
			Management: &coh.PortSpecWithSSL{
				Port: ptr.To(int32(1234)),
			},
		},
		Ports: []coh.NamedPortSpec{np},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	addPorts(stsExpected, coh.ContainerNameCoherence, corev1.ContainerPort{
		Name:          coh.PortNameManagement,
		ContainerPort: 1234,
		Protocol:      protocol,
	})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)

	svc := np.CreateService(deployment)
	g.Expect(svc).NotTo(BeNil())
	g.Expect(len(svc.Spec.Ports)).To(Equal(1))
	g.Expect(svc.Spec.Ports[0].Port).To(Equal(int32(1234)))
	g.Expect(svc.Spec.Ports[0].TargetPort.IntVal).To(Equal(int32(1234)))
}

func TestCreateStatefulSetWithManagementPortAndServicePortWhenNoPortValueSpecified(t *testing.T) {
	g := NewGomegaWithT(t)

	protocol := corev1.ProtocolTCP
	np := coh.NamedPortSpec{
		Name: coh.PortNameManagement,
		Service: &coh.ServiceSpec{
			Port: ptr.To(int32(1234)),
		},
	}

	spec := coh.CoherenceResourceSpec{
		Ports: []coh.NamedPortSpec{np},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	addPorts(stsExpected, coh.ContainerNameCoherence, corev1.ContainerPort{
		Name:          coh.PortNameManagement,
		ContainerPort: coh.DefaultManagementPort,
		Protocol:      protocol,
	})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)

	svc := np.CreateService(deployment)
	g.Expect(svc).NotTo(BeNil())
	g.Expect(len(svc.Spec.Ports)).To(Equal(1))
	g.Expect(svc.Spec.Ports[0].Port).To(Equal(int32(1234)))
	g.Expect(svc.Spec.Ports[0].TargetPort.IntVal).To(Equal(coh.DefaultManagementPort))
}

func TestCreateStatefulSetWithManagementPortAndServicePortWhenManagementPortValueSpecified(t *testing.T) {
	g := NewGomegaWithT(t)

	protocol := corev1.ProtocolTCP
	np := coh.NamedPortSpec{
		Name: coh.PortNameManagement,
		Service: &coh.ServiceSpec{
			Port: ptr.To(int32(1234)),
		},
	}

	spec := coh.CoherenceResourceSpec{
		Coherence: &coh.CoherenceSpec{
			Management: &coh.PortSpecWithSSL{
				Port: ptr.To(int32(9876)),
			},
		},
		Ports: []coh.NamedPortSpec{np},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	addPorts(stsExpected, coh.ContainerNameCoherence, corev1.ContainerPort{
		Name:          coh.PortNameManagement,
		ContainerPort: 9876,
		Protocol:      protocol,
	})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)

	svc := np.CreateService(deployment)
	g.Expect(svc).NotTo(BeNil())
	g.Expect(len(svc.Spec.Ports)).To(Equal(1))
	g.Expect(svc.Spec.Ports[0].Port).To(Equal(int32(1234)))
	g.Expect(svc.Spec.Ports[0].TargetPort.IntVal).To(Equal(int32(9876)))
}
