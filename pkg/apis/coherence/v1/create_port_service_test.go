/*
 * Copyright (c) 2019, 2020, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1_test

import (
	"github.com/go-test/deep"
	. "github.com/onsi/gomega"
	coh "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"testing"
)

func TestNamedPortSpec_CreateServiceWithMinimalFields(t *testing.T) {
	g := NewGomegaWithT(t)
	c := coh.CoherenceDeployment{}
	c.Name = "test-deployment"
	c.Spec.Role = "storage"

	np := coh.NamedPortSpec{
		Name: "foo",
		PortSpec: coh.PortSpec{
			Port: 19,
		},
	}

	labels := c.CreateCommonLabels()
	labels[coh.LabelComponent] = coh.LabelComponentPortService
	labels[coh.LabelPort] = np.Name

	selector := c.CreateCommonLabels()
	selector[coh.LabelComponent] = coh.LabelComponentCoherencePod

	expected := corev1.Service{
		ObjectMeta: v1.ObjectMeta{
			Name:   "test-deployment-foo",
			Labels: labels,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       "foo",
					Protocol:   corev1.ProtocolTCP,
					Port:       19,
					TargetPort: intstr.FromString("foo"),
					NodePort:   0,
				},
			},
			Selector: selector,
		},
	}

	svc := np.CreateService(&c)
	g.Expect(svc).NotTo(BeNil())
	g.Expect(deep.Equal(*svc, expected)).To(BeNil())
}

func TestNamedPortSpec_CreateServiceWithProtocol(t *testing.T) {
	g := NewGomegaWithT(t)
	d := coh.CoherenceDeployment{}
	d.Name = "test-deployment"
	d.Spec.Role = "storage"

	udp := corev1.ProtocolUDP

	np := coh.NamedPortSpec{
		Name: "foo",
		PortSpec: coh.PortSpec{
			Port:     19,
			Protocol: &udp,
		},
	}

	svc := np.CreateService(&d)
	g.Expect(svc).NotTo(BeNil())
	g.Expect(svc.Spec.Ports[0].Protocol).To(Equal(udp))
}

func TestNamedPortSpec_CreateServiceWithNodePort(t *testing.T) {
	g := NewGomegaWithT(t)
	d := coh.CoherenceDeployment{}
	d.Name = "test-deployment"
	d.Spec.Role = "storage"

	np := coh.NamedPortSpec{
		Name: "foo",
		PortSpec: coh.PortSpec{
			Port:     19,
			NodePort: int32Ptr(6676),
		},
	}

	svc := np.CreateService(&d)
	g.Expect(svc).NotTo(BeNil())
	g.Expect(svc.Spec.Ports[0].NodePort).To(Equal(int32(6676)))
}

func TestNamedPortSpec_CreateServiceWithService(t *testing.T) {
	g := NewGomegaWithT(t)
	d := coh.CoherenceDeployment{}
	d.Name = "test-deployment"
	d.Spec.Role = "storage"

	tp := corev1.ServiceTypeClusterIP
	ipf := corev1.IPv4Protocol
	etpt := corev1.ServiceExternalTrafficPolicyTypeLocal
	sa := corev1.ServiceAffinityClientIP
	sac := corev1.SessionAffinityConfig{
		ClientIP: &corev1.ClientIPConfig{TimeoutSeconds: int32Ptr(9876)},
	}

	np := coh.NamedPortSpec{
		Name: "foo",
		PortSpec: coh.PortSpec{
			Port: 19,
			Service: &coh.ServiceSpec{
				Name:                     stringPtr("bar"),
				Port:                     int32Ptr(99),
				Type:                     &tp,
				ClusterIP:                stringPtr("10.10.10.99"),
				ExternalIPs:              []string{"192.164.1.99", "192.164.1.100"},
				LoadBalancerIP:           stringPtr("10.10.10.10"),
				Labels:                   nil,
				Annotations:              nil,
				SessionAffinity:          &sa,
				LoadBalancerSourceRanges: []string{"A", "B"},
				ExternalName:             stringPtr("ext-bar"),
				ExternalTrafficPolicy:    &etpt,
				HealthCheckNodePort:      int32Ptr(1234),
				PublishNotReadyAddresses: boolPtr(true),
				SessionAffinityConfig:    &sac,
				IPFamily:                 &ipf,
			},
		},
	}

	labels := d.CreateCommonLabels()
	labels[coh.LabelComponent] = coh.LabelComponentPortService
	labels[coh.LabelPort] = np.Name

	selector := d.CreateCommonLabels()
	selector[coh.LabelComponent] = coh.LabelComponentCoherencePod

	expected := corev1.Service{
		ObjectMeta: v1.ObjectMeta{
			Name:   "bar",
			Labels: labels,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       "foo",
					Protocol:   corev1.ProtocolTCP,
					Port:       99,
					TargetPort: intstr.FromString("foo"),
				},
			},
			Selector:                 selector,
			ClusterIP:                "10.10.10.99",
			Type:                     tp,
			ExternalIPs:              []string{"192.164.1.99", "192.164.1.100"},
			SessionAffinity:          sa,
			LoadBalancerIP:           "10.10.10.10",
			LoadBalancerSourceRanges: []string{"A", "B"},
			ExternalName:             "ext-bar",
			ExternalTrafficPolicy:    etpt,
			HealthCheckNodePort:      1234,
			PublishNotReadyAddresses: true,
			SessionAffinityConfig:    &sac,
			IPFamily:                 &ipf,
		},
	}

	svc := np.CreateService(&d)
	g.Expect(svc).NotTo(BeNil())
	g.Expect(deep.Equal(*svc, expected)).To(BeNil())
}
