/*
 * Copyright (c) 2020, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package controller_test

import (
	"fmt"
	. "github.com/onsi/gomega"
	coh "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/test/e2e/helper/matchers"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
	"testing"
)

func TestCreateResourcesForMinimalDeployment(t *testing.T) {
	g := NewGomegaWithT(t)

	// Create the test deployment (a minimal deployment configuration)
	deployment := coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{
			Name: "operator-test",
		},
	}

	resources, mgr := Reconcile(t, deployment)

	// Verify the expected k8s events
	AssertStatefulSetCreationEvent(t, deployment.Name, mgr)
	AssertNoRemainingEvents(t, mgr)

	// Should have created the correct number of resources
	g.Expect(len(resources)).To(Equal(6))

	// Resource 0 = Deployment
	c, err := toCoherence(mgr, resources[0])
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(c.GetName()).To(Equal(deployment.GetName()))

	// Resource 1 = Operator config Secret
	opCfg, err := toSecret(mgr, resources[1])
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(opCfg.GetName()).To(Equal(coh.OperatorConfigName))

	// Resource 2 = deployment storage Secret
	store, err := toSecret(mgr, resources[2])
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(store.GetName()).To(Equal(deployment.Name))

	// Resource 3 = WKA Service
	wka, err := toService(mgr, resources[3])
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(wka.GetName()).To(Equal(deployment.GetWkaServiceName()))

	// Resource 4 = Service for StatefulSet
	ss, err := toService(mgr, resources[4])
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(ss.GetName()).To(Equal(deployment.GetHeadlessServiceName()))

	// Resource 5 = StatefulSet
	sts, err := toStatefulSet(mgr, resources[5])
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(sts.GetName()).To(Equal(deployment.Name))
	// StatefulSet should have default replica count
	g.Expect(sts.Spec.Replicas).NotTo(BeNil())
	g.Expect(*sts.Spec.Replicas).To(Equal(coh.DefaultReplicas))
}

func TestCreateResourcesDeploymentNotInWKA(t *testing.T) {
	g := NewGomegaWithT(t)

	deployment := coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{
			Name: "operator-test",
		},
		Spec: coh.CoherenceResourceSpec{
			Coherence: &coh.CoherenceSpec{
				ExcludeFromWKA: pointer.BoolPtr(true),
			},
		},
	}

	resources, mgr := Reconcile(t, deployment)

	// Verify the expected k8s events
	AssertStatefulSetCreationEvent(t, deployment.Name, mgr)
	AssertNoRemainingEvents(t, mgr)

	// Should have created the correct number of resources
	g.Expect(len(resources)).To(Equal(6))

	// Resource 0 = Deployment
	c, err := toCoherence(mgr, resources[0])
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(c.GetName()).To(Equal(deployment.GetName()))

	// Resource 1 = Operator config Secret
	opCfg, err := toSecret(mgr, resources[1])
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(opCfg.GetName()).To(Equal(coh.OperatorConfigName))

	// Resource 2 = deployment storage Secret
	store, err := toSecret(mgr, resources[2])
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(store.GetName()).To(Equal(deployment.Name))

	// Resource 3 = WKA Service
	wka, err := toService(mgr, resources[3])
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(wka.GetName()).To(Equal(deployment.GetWkaServiceName()))

	// Resource 4 = Service for StatefulSet
	ss, err := toService(mgr, resources[4])
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(ss.GetName()).To(Equal(deployment.GetHeadlessServiceName()))

	// Resource 5 = StatefulSet
	sts, err := toStatefulSet(mgr, resources[5])
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(sts.GetName()).To(Equal(deployment.Name))
	// StatefulSet should have default replica count
	g.Expect(sts.Spec.Replicas).NotTo(BeNil())
	g.Expect(*sts.Spec.Replicas).To(Equal(coh.DefaultReplicas))
}

func TestCreateResourcesDeploymentWithExistingWKA(t *testing.T) {
	g := NewGomegaWithT(t)

	deployment := coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{
			Name: "operator-test",
		},
		Spec: coh.CoherenceResourceSpec{
			Coherence: &coh.CoherenceSpec{
				WKA: &coh.CoherenceWKASpec{
					Deployment: "foo",
					Namespace:  "",
				},
			},
		},
	}

	resources, mgr := Reconcile(t, deployment)

	// Verify the expected k8s events
	AssertStatefulSetCreationEvent(t, deployment.Name, mgr)
	AssertNoRemainingEvents(t, mgr)

	// Should have created the correct number of resources
	g.Expect(len(resources)).To(Equal(5))

	// Resource 0 = Deployment
	c, err := toCoherence(mgr, resources[0])
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(c.GetName()).To(Equal(deployment.GetName()))

	// Resource 1 = Operator config Secret
	opCfg, err := toSecret(mgr, resources[1])
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(opCfg.GetName()).To(Equal(coh.OperatorConfigName))

	// Resource 2 = deployment storage Secret
	store, err := toSecret(mgr, resources[2])
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(store.GetName()).To(Equal(deployment.Name))

	// Resource 3 = Service for StatefulSet
	ss, err := toService(mgr, resources[3])
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(ss.GetName()).To(Equal(deployment.GetHeadlessServiceName()))

	// Resource 4 = StatefulSet
	sts, err := toStatefulSet(mgr, resources[4])
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(sts.GetName()).To(Equal(deployment.Name))
	// StatefulSet should have default replica count
	g.Expect(sts.Spec.Replicas).NotTo(BeNil())
	g.Expect(*sts.Spec.Replicas).To(Equal(coh.DefaultReplicas))
}

func TestCreateResourcesForDeploymentWithReplicaCount(t *testing.T) {
	g := NewGomegaWithT(t)

	deployment := coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{
			Name: "operator-test",
		},
		Spec: coh.CoherenceResourceSpec{
			Replicas: pointer.Int32Ptr(5),
		},
	}

	// run the reconciler
	resources, mgr := Reconcile(t, deployment)

	// Should have created the correct number of resources
	g.Expect(len(resources)).To(Equal(6))

	// Resource 5 = StatefulSet
	sts, err := toStatefulSet(mgr, resources[5])
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(sts.GetName()).To(Equal(deployment.Name))
	// StatefulSet should have default replica count
	g.Expect(sts.Spec.Replicas).NotTo(BeNil())
	g.Expect(*sts.Spec.Replicas).To(Equal(int32(5)))
}

func TestCreateResourcesForDeploymentWithClusterName(t *testing.T) {
	g := NewGomegaWithT(t)

	// Create the test deployment with a clusterName
	clusterName := "test-cluster"
	deployment := coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{
			Name: "operator-test",
		},
		Spec: coh.CoherenceResourceSpec{
			Cluster: &clusterName,
		},
	}

	resources, mgr := Reconcile(t, deployment)

	// Verify the expected k8s events
	AssertStatefulSetCreationEvent(t, deployment.Name, mgr)
	AssertNoRemainingEvents(t, mgr)

	// Should have created the correct number of resources
	g.Expect(len(resources)).To(Equal(6))

	// Resource 3 = WKA Service
	wka, err := toService(mgr, resources[3])
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(wka.GetName()).To(Equal(deployment.GetWkaServiceName()))
	// WKA service should have the correct selector
	g.Expect(len(wka.Spec.Selector)).To(Equal(3))
	g.Expect(wka.Spec.Selector).To(HaveKeyWithValue(coh.LabelCoherenceCluster, clusterName))
	g.Expect(wka.Spec.Selector).To(HaveKeyWithValue(coh.LabelCoherenceWKAMember, "true"))
	g.Expect(wka.Spec.Selector).To(HaveKeyWithValue(coh.LabelComponent, coh.LabelComponentCoherencePod))

	// Resource 5 = StatefulSet
	sts, err := toStatefulSet(mgr, resources[5])
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(sts.GetName()).To(Equal(deployment.Name))
	// StatefulSet should have correct cluster name env-var
	container, found := FindContainer(coh.ContainerNameCoherence, sts)
	g.Expect(found).To(BeTrue())
	g.Expect(container.Env).To(matchers.HaveEnvVar(corev1.EnvVar{Name: coh.EnvVarCohClusterName, Value: clusterName}))
}

func TestCreateResourcesForDeploymentWithHealthPort(t *testing.T) {
	g := NewGomegaWithT(t)

	// Create the test deployment with a clusterName
	var health = 19
	deployment := coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{
			Name: "operator-test",
		},
		Spec: coh.CoherenceResourceSpec{
			HealthPort: pointer.Int32Ptr(int32(health)),
		},
	}

	_, mgr := Reconcile(t, deployment)

	// Get the StatefulSet
	sts, err := mgr.Client.GetStatefulSet(deployment.Namespace, deployment.Name)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(sts.GetName()).To(Equal(deployment.Name))

	container, found := FindContainer(coh.ContainerNameCoherence, sts)
	g.Expect(found).To(BeTrue())
	// Coherence container should have correct health port
	port, found := FindContainerPort(container, coh.PortNameHealth)
	g.Expect(found).To(BeTrue())
	g.Expect(port.ContainerPort).To(Equal(int32(health)))
	// Coherence container should have correct health env-var
	g.Expect(container.Env).To(matchers.HaveEnvVar(corev1.EnvVar{Name: coh.EnvVarCohHealthPort, Value: fmt.Sprintf("%d", health)}))
	// Coherence container readiness probe should use correct port
	g.Expect(container.ReadinessProbe.HTTPGet.Port.IntValue()).To(Equal(health))
	// Coherence container liveness probe should use correct port
	g.Expect(container.LivenessProbe.HTTPGet.Port.IntValue()).To(Equal(health))
}
