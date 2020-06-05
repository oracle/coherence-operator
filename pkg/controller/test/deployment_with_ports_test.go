/*
 * Copyright (c) 2020, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package controller_test

import (
	. "github.com/onsi/gomega"
	coh "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"
	"testing"
)

func TestCreateDeploymentWithPort(t *testing.T) {
	g := NewGomegaWithT(t)

	tcp := corev1.ProtocolUDP

	deployment := coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{
			Name: "operator-test",
		},
		Spec: coh.CoherenceResourceSpec{
			Ports: []coh.NamedPortSpec{
				{
					Name:     "extend",
					Port:     20000,
					Protocol: &tcp,
					NodePort: pointer.Int32Ptr(80),
					HostPort: pointer.Int32Ptr(8080),
					HostIP:   pointer.StringPtr("10.10.10.1"),
				},
			},
		},
	}

	// run the reconciler
	resources, mgr := Reconcile(t, deployment)

	// Should have created the correct number of resources
	g.Expect(len(resources)).To(Equal(7))

	// Resource 5 = Service for the port
	svc2, err := toService(mgr, resources[5])
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(svc2.GetName()).To(Equal(deployment.Name + "-extend"))
	g.Expect(len(svc2.Spec.Ports)).To(Equal(1))
	g.Expect(svc2.Spec.Ports[0]).To(Equal(corev1.ServicePort{
		Name:       "extend",
		Protocol:   tcp,
		Port:       20000,
		TargetPort: intstr.FromInt(20000),
		NodePort:   80,
	}))

	// Resource 6 = StatefulSet
	sts, err := toStatefulSet(mgr, resources[6])
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(sts.GetName()).To(Equal(deployment.Name))

	container, found := FindContainer(coh.ContainerNameCoherence, sts)
	g.Expect(found).To(BeTrue())
	port, found := FindContainerPort(container, "extend")
	g.Expect(found).To(BeTrue())
	g.Expect(port).To(Equal(corev1.ContainerPort{
		Name:          "extend",
		HostPort:      8080,
		ContainerPort: 20000,
		Protocol:      tcp,
		HostIP:        "10.10.10.1",
	}))
}

func TestUpdateDeploymentWithPort(t *testing.T) {
	g := NewGomegaWithT(t)

	tcp := corev1.ProtocolUDP

	original := coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{
			Name: "operator-test",
		},
		Spec: coh.CoherenceResourceSpec{
			Ports: []coh.NamedPortSpec{
				{
					Name:     "extend",
					Port:     20000,
					Protocol: &tcp,
					NodePort: pointer.Int32Ptr(80),
					HostPort: pointer.Int32Ptr(8080),
					HostIP:   pointer.StringPtr("10.10.10.1"),
				},
			},
		},
	}

	updated := coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{
			Name: "operator-test",
		},
		Spec: coh.CoherenceResourceSpec{
			Ports: []coh.NamedPortSpec{
				{
					Name:     "extend",
					Port:     30000,
					Protocol: &tcp,
					NodePort: pointer.Int32Ptr(88),
					HostPort: pointer.Int32Ptr(8888),
					HostIP:   pointer.StringPtr("10.10.10.2"),
				},
			},
		},
	}

	// run the reconciler
	mgr := ReconcileAndUpdate(t, original, updated)

	// Get the updated Service for the port
	svc, err := mgr.Client.GetService(original.Namespace, original.Name+"-extend")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(svc.GetName()).To(Equal(original.Name + "-extend"))
	g.Expect(len(svc.Spec.Ports)).To(Equal(1))
	g.Expect(svc.Spec.Ports[0]).To(Equal(corev1.ServicePort{
		Name:       "extend",
		Protocol:   tcp,
		Port:       30000,
		TargetPort: intstr.FromInt(30000),
		NodePort:   88,
	}))

	// Get the updated StatefulSet
	sts, err := mgr.Client.GetStatefulSet(original.Namespace, original.Name)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(sts.GetName()).To(Equal(original.Name))

	container, found := FindContainer(coh.ContainerNameCoherence, sts)
	g.Expect(found).To(BeTrue())
	port, found := FindContainerPort(container, "extend")
	g.Expect(found).To(BeTrue())
	g.Expect(port).To(Equal(corev1.ContainerPort{
		Name:          "extend",
		HostPort:      8888,
		ContainerPort: 30000,
		Protocol:      tcp,
		HostIP:        "10.10.10.2",
	}))
}

func TestUpdateDeploymentWithAdditionalPort(t *testing.T) {
	g := NewGomegaWithT(t)

	tcp := corev1.ProtocolUDP

	original := coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{
			Name: "operator-test",
		},
		Spec: coh.CoherenceResourceSpec{
			Ports: []coh.NamedPortSpec{
				{
					Name:     "management",
					Port:     10000,
					Protocol: &tcp,
					NodePort: pointer.Int32Ptr(80),
					HostPort: pointer.Int32Ptr(8080),
					HostIP:   pointer.StringPtr("10.10.10.1"),
				},
			},
		},
	}

	updated := coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{
			Name: "operator-test",
		},
		Spec: coh.CoherenceResourceSpec{
			Ports: []coh.NamedPortSpec{
				{
					Name:     coh.PortNameManagement,
					Port:     10000,
					Protocol: &tcp,
					NodePort: pointer.Int32Ptr(80),
					HostPort: pointer.Int32Ptr(8080),
					HostIP:   pointer.StringPtr("10.10.10.1"),
				},
				{
					Name:     coh.PortNameMetrics,
					Port:     20000,
					Protocol: &tcp,
					NodePort: pointer.Int32Ptr(88),
					HostPort: pointer.Int32Ptr(8888),
					HostIP:   pointer.StringPtr("10.10.10.2"),
				},
			},
		},
	}

	// run the reconciler
	mgr := ReconcileAndUpdate(t, original, updated)

	// Get the Service for the original port
	svc, err := mgr.Client.GetService(original.Namespace, original.Name+"-management")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(svc.GetName()).To(Equal(original.Name + "-management"))
	g.Expect(len(svc.Spec.Ports)).To(Equal(1))
	g.Expect(svc.Spec.Ports[0]).To(Equal(corev1.ServicePort{
		Name:       "management",
		Protocol:   tcp,
		Port:       10000,
		TargetPort: intstr.FromInt(10000),
		NodePort:   80,
	}))

	// Get the Service for the added port
	svcAdded, err := mgr.Client.GetService(original.Namespace, original.Name+"-metrics")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(svcAdded.GetName()).To(Equal(original.Name + "-metrics"))
	g.Expect(len(svcAdded.Spec.Ports)).To(Equal(1))
	g.Expect(svcAdded.Spec.Ports[0]).To(Equal(corev1.ServicePort{
		Name:       coh.PortNameMetrics,
		Protocol:   tcp,
		Port:       20000,
		TargetPort: intstr.FromInt(20000),
		NodePort:   88,
	}))

	// Get the updated StatefulSet
	sts, err := mgr.Client.GetStatefulSet(original.Namespace, original.Name)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(sts.GetName()).To(Equal(original.Name))

	container, found := FindContainer(coh.ContainerNameCoherence, sts)
	g.Expect(found).To(BeTrue())
	// The StatefulSet should have the original port
	managementPort, found := FindContainerPort(container, coh.PortNameManagement)
	g.Expect(found).To(BeTrue())
	g.Expect(managementPort).To(Equal(corev1.ContainerPort{
		Name:          "management",
		HostPort:      8080,
		ContainerPort: 10000,
		Protocol:      tcp,
		HostIP:        "10.10.10.1",
	}))
	// The StatefulSet should have the added port
	metricsPort, found := FindContainerPort(container, coh.PortNameMetrics)
	g.Expect(found).To(BeTrue())
	g.Expect(metricsPort).To(Equal(corev1.ContainerPort{
		Name:          coh.PortNameMetrics,
		HostPort:      8888,
		ContainerPort: 20000,
		Protocol:      tcp,
		HostIP:        "10.10.10.2",
	}))
}
