/*
 * Copyright (c) 2020, 2021, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package controller_test

import (
	. "github.com/onsi/gomega"
	coh "github.com/oracle/coherence-operator/api/v1"
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
