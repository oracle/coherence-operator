/*
 * Copyright (c) 2020, 2023, Oracle and/or its affiliates.
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

var expectedJobResources = 5

func TestCreateResourcesForMinimalJobDeployment(t *testing.T) {
	g := NewGomegaWithT(t)

	// Create the test deployment (a minimal deployment configuration)
	deployment := coh.CoherenceJob{
		ObjectMeta: metav1.ObjectMeta{
			Name: "operator-test-job",
		},
	}

	resources, mgr := ReconcileJob(t, deployment)

	// Verify the expected k8s events
	AssertJobCreationEvent(t, deployment.Name, mgr)
	AssertNoRemainingEvents(t, mgr)

	// Should have created the correct number of resources
	g.Expect(len(resources)).To(Equal(expectedJobResources))

	// Resource 0 is the CoherenceJob resource
	c, err := toCoherenceJob(mgr, resources[0])
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

	// Resource 4 = Job
	job, err := toJob(mgr, resources[4])
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(job.GetName()).To(Equal(deployment.Name))
	// Job should have default replica count
	g.Expect(job.Spec.Parallelism).NotTo(BeNil())
	g.Expect(*job.Spec.Parallelism).To(Equal(coh.DefaultJobReplicas))
}

func TestShouldNotAddFinalizer(t *testing.T) {
	g := NewGomegaWithT(t)

	// Create the test deployment (a minimal deployment configuration)
	deployment := coh.CoherenceJob{
		ObjectMeta: metav1.ObjectMeta{
			Name: "operator-test-job",
		},
	}

	resources, mgr := ReconcileJob(t, deployment)
	// Should have created resources
	g.Expect(len(resources)).NotTo(BeZero())
	// Resource 0 = CoherenceJob resource
	c, err := toCoherenceJob(mgr, resources[0])
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(c.GetName()).To(Equal(deployment.GetName()))
	// the finalizer should not be present
	g.Expect(c.GetFinalizers()).NotTo(ContainElement(coh.CoherenceFinalizer))
}

func TestCreateJobResourcesDeploymentNotInWKA(t *testing.T) {
	g := NewGomegaWithT(t)

	deployment := coh.CoherenceJob{
		ObjectMeta: metav1.ObjectMeta{
			Name: "operator-test",
		},
		Spec: coh.CoherenceJobResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Coherence: &coh.CoherenceSpec{
					ExcludeFromWKA: pointer.Bool(true),
				},
			},
		},
	}

	resources, mgr := ReconcileJob(t, deployment)

	// Verify the expected k8s events
	AssertJobCreationEvent(t, deployment.Name, mgr)
	AssertNoRemainingEvents(t, mgr)

	// Should have created the correct number of resources
	g.Expect(len(resources)).To(Equal(expectedJobResources))

	// Resource 0 is the CoherenceJob resource
	c, err := toCoherenceJob(mgr, resources[0])
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

	// Resource 4 = Job
	job, err := toJob(mgr, resources[4])
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(job.GetName()).To(Equal(deployment.Name))
	// Job should have default replica count
	g.Expect(job.Spec.Parallelism).NotTo(BeNil())
	g.Expect(*job.Spec.Parallelism).To(Equal(coh.DefaultJobReplicas))
}

func TestCreateJobResourcesDeploymentWithExistingWKA(t *testing.T) {
	g := NewGomegaWithT(t)

	deployment := coh.CoherenceJob{
		ObjectMeta: metav1.ObjectMeta{
			Name: "operator-test",
		},
		Spec: coh.CoherenceJobResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Coherence: &coh.CoherenceSpec{
					WKA: &coh.CoherenceWKASpec{
						Deployment: "foo",
						Namespace:  "",
					},
				},
			},
		},
	}

	resources, mgr := ReconcileJob(t, deployment)

	// Verify the expected k8s events
	AssertJobCreationEvent(t, deployment.Name, mgr)
	AssertNoRemainingEvents(t, mgr)

	// Should have created the correct number of resources
	g.Expect(len(resources)).To(Equal(expectedJobResources - 1))

	// Resource 0 is the CoherenceJob resource
	c, err := toCoherenceJob(mgr, resources[0])
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

	// Resource 3 = Job
	job, err := toJob(mgr, resources[3])
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(job.GetName()).To(Equal(deployment.Name))
	// Job should have default replica count
	g.Expect(job.Spec.Parallelism).NotTo(BeNil())
	g.Expect(*job.Spec.Parallelism).To(Equal(coh.DefaultJobReplicas))
}

func TestCreateJobResourcesForDeploymentWithReplicaCount(t *testing.T) {
	g := NewGomegaWithT(t)

	deployment := coh.CoherenceJob{
		ObjectMeta: metav1.ObjectMeta{
			Name: "operator-test",
		},
		Spec: coh.CoherenceJobResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Replicas: pointer.Int32(5),
			},
		},
	}

	// run the reconciler
	resources, mgr := ReconcileJob(t, deployment)

	// Should have created the correct number of resources
	g.Expect(len(resources)).To(Equal(expectedJobResources))

	// Resource 4 = Job
	job, err := toJob(mgr, resources[4])
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(job.GetName()).To(Equal(deployment.Name))
	// Job should have default replica count
	g.Expect(job.Spec.Parallelism).NotTo(BeNil())
	g.Expect(*job.Spec.Parallelism).To(Equal(int32(5)))
}

func TestCreateJobResourcesForDeploymentWithClusterName(t *testing.T) {
	g := NewGomegaWithT(t)

	// Create the test deployment with a clusterName
	clusterName := "test-cluster"
	deployment := coh.CoherenceJob{
		ObjectMeta: metav1.ObjectMeta{
			Name: "operator-test",
		},
		Spec: coh.CoherenceJobResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Cluster: &clusterName,
			},
		},
	}

	resources, mgr := ReconcileJob(t, deployment)

	// Verify the expected k8s events
	AssertJobCreationEvent(t, deployment.Name, mgr)
	AssertNoRemainingEvents(t, mgr)

	// Should have created the correct number of resources
	g.Expect(len(resources)).To(Equal(expectedJobResources))

	// Resource 3 = WKA Service
	wka, err := toService(mgr, resources[3])
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(wka.GetName()).To(Equal(deployment.GetWkaServiceName()))
	// WKA service should have the correct selector
	g.Expect(len(wka.Spec.Selector)).To(Equal(3))
	g.Expect(wka.Spec.Selector).To(HaveKeyWithValue(coh.LabelCoherenceCluster, clusterName))
	g.Expect(wka.Spec.Selector).To(HaveKeyWithValue(coh.LabelCoherenceWKAMember, "true"))
	g.Expect(wka.Spec.Selector).To(HaveKeyWithValue(coh.LabelComponent, coh.LabelComponentCoherencePod))

	// Resource 4 = Job
	job, err := toJob(mgr, resources[4])
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(job.GetName()).To(Equal(deployment.Name))
	// Job should have correct cluster name env-var
	container, found := FindContainerInJob(coh.ContainerNameCoherence, job)
	g.Expect(found).To(BeTrue())
	g.Expect(container.Env).To(matchers.HaveEnvVar(corev1.EnvVar{Name: coh.EnvVarCohClusterName, Value: clusterName}))
}

func TestCreateJobResourcesForDeploymentWithHealthPort(t *testing.T) {
	g := NewGomegaWithT(t)

	// Create the test deployment with a clusterName
	var health = 19
	deployment := coh.CoherenceJob{
		ObjectMeta: metav1.ObjectMeta{
			Name: "operator-test",
		},
		Spec: coh.CoherenceJobResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				HealthPort: pointer.Int32(int32(health)),
			},
		},
	}

	_, mgr := ReconcileJob(t, deployment)

	// Get the Job
	job, err := mgr.Client.GetJob(deployment.Namespace, deployment.Name)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(job.GetName()).To(Equal(deployment.Name))

	container, found := FindContainerInJob(coh.ContainerNameCoherence, job)
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
