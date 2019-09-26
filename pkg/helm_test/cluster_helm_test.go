/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package helm_test

import (
	. "github.com/onsi/gomega"
	cohv1 "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"testing"
)

/*
 * These tests verify the various scenarios for setting Coherence cache configuration
 * in a CoherenceCluster.
 */

func TestClusterFromMinimalYaml(t *testing.T) {
	g := NewGomegaWithT(t)

	// Use the specified yaml files to create a CoherenceCluster and trigger a fake end-to-end
	// reconcile to obtain the resources that would have been created by the Helm operator.
	result, cluster, err := CreateCluster("minimal-cluster.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	g.Expect(result.Size()).To(Equal(3))

	// should have created config map
	name := cluster.GetFullRoleName(cohv1.DefaultRoleName) + "-scripts"
	cm := corev1.ConfigMap{}
	err = result.Get(name, &cm)
	g.Expect(err).NotTo(HaveOccurred())

	// should have created one headless Service
	name = cluster.GetFullRoleName(cohv1.DefaultRoleName) + "-headless"
	svc := corev1.Service{}
	err = result.Get(name, &svc)
	g.Expect(err).NotTo(HaveOccurred())

	// Should have one StatefulSet
	list := appsv1.StatefulSetList{}
	err = result.List(&list)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(len(list.Items)).To(Equal(1))

	// Obtain the StatefulSet that Helm would have created
	sts, err := findStatefulSet(result, cluster, cohv1.DefaultRoleName)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(*sts.Spec.Replicas).To(Equal(cohv1.DefaultReplicas))

}

func TestClusterImplicitRoleOneReplica(t *testing.T) {
	g := NewGomegaWithT(t)

	// Use the specified yaml files to create a CoherenceCluster and trigger a fake end-to-end
	// reconcile to obtain the resources that would have been created by the Helm operator.
	result, cluster, err := CreateCluster("cluster-test-implicit-role-one-replica.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	g.Expect(result.Size()).To(Equal(3))

	// should have created config map
	name := cluster.GetFullRoleName(cohv1.DefaultRoleName) + "-scripts"
	cm := corev1.ConfigMap{}
	err = result.Get(name, &cm)
	g.Expect(err).NotTo(HaveOccurred())

	// should have created one headless Service
	name = cluster.GetFullRoleName(cohv1.DefaultRoleName) + "-headless"
	svc := corev1.Service{}
	err = result.Get(name, &svc)
	g.Expect(err).NotTo(HaveOccurred())

	// Should have one StatefulSet
	list := appsv1.StatefulSetList{}
	err = result.List(&list)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(len(list.Items)).To(Equal(1))

	// Obtain the StatefulSet that Helm would have created
	sts, err := findStatefulSet(result, cluster, cohv1.DefaultRoleName)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(*sts.Spec.Replicas).To(Equal(int32(1)))
}

func TestClusterExplicitRoles(t *testing.T) {
	g := NewGomegaWithT(t)

	// Use the specified yaml files to create a CoherenceCluster and trigger a fake end-to-end
	// reconcile to obtain the resources that would have been created by the Helm operator.
	result, cluster, err := CreateCluster("cluster-test-explicit-roles.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	g.Expect(result.Size()).To(Equal(6))

	// should have created config map for data role
	name := cluster.GetFullRoleName("data") + "-scripts"
	cm := corev1.ConfigMap{}
	err = result.Get(name, &cm)
	g.Expect(err).NotTo(HaveOccurred())

	// should have created config map for proxy role
	name = cluster.GetFullRoleName("proxy") + "-scripts"
	cm = corev1.ConfigMap{}
	err = result.Get(name, &cm)
	g.Expect(err).NotTo(HaveOccurred())

	// should have created headless Service for data role
	name = cluster.GetFullRoleName("data") + "-headless"
	svc := corev1.Service{}
	err = result.Get(name, &svc)
	g.Expect(err).NotTo(HaveOccurred())

	// should have created headless Service for data role
	name = cluster.GetFullRoleName("proxy") + "-headless"
	svc = corev1.Service{}
	err = result.Get(name, &svc)
	g.Expect(err).NotTo(HaveOccurred())

	// Should have two StatefulSets
	list := appsv1.StatefulSetList{}
	err = result.List(&list)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(len(list.Items)).To(Equal(2))

	// Obtain the StatefulSet that Helm would have created for the data role
	sts, err := findStatefulSet(result, cluster, "data")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(*sts.Spec.Replicas).To(Equal(int32(1)))

	// Obtain the StatefulSet that Helm would have created for the data role
	sts, err = findStatefulSet(result, cluster, "proxy")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(*sts.Spec.Replicas).To(Equal(int32(2)))
}

func TestClusterExplicitRolesWithDefaultReplicas(t *testing.T) {
	g := NewGomegaWithT(t)

	// Use the specified yaml files to create a CoherenceCluster and trigger a fake end-to-end
	// reconcile to obtain the resources that would have been created by the Helm operator.
	result, cluster, err := CreateCluster("cluster-test-explicit-roles-with-default.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	g.Expect(result.Size()).To(Equal(6))

	// should have created config map for data role
	name := cluster.GetFullRoleName("data") + "-scripts"
	cm := corev1.ConfigMap{}
	err = result.Get(name, &cm)
	g.Expect(err).NotTo(HaveOccurred())

	// should have created config map for proxy role
	name = cluster.GetFullRoleName("proxy") + "-scripts"
	cm = corev1.ConfigMap{}
	err = result.Get(name, &cm)
	g.Expect(err).NotTo(HaveOccurred())

	// should have created headless Service for data role
	name = cluster.GetFullRoleName("data") + "-headless"
	svc := corev1.Service{}
	err = result.Get(name, &svc)
	g.Expect(err).NotTo(HaveOccurred())

	// should have created headless Service for data role
	name = cluster.GetFullRoleName("proxy") + "-headless"
	svc = corev1.Service{}
	err = result.Get(name, &svc)
	g.Expect(err).NotTo(HaveOccurred())

	// Should have two StatefulSets
	list := appsv1.StatefulSetList{}
	err = result.List(&list)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(len(list.Items)).To(Equal(2))

	// Obtain the StatefulSet that Helm would have created for the data role
	sts, err := findStatefulSet(result, cluster, "data")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(*sts.Spec.Replicas).To(Equal(int32(4)))

	// Obtain the StatefulSet that Helm would have created for the data role
	sts, err = findStatefulSet(result, cluster, "proxy")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(*sts.Spec.Replicas).To(Equal(int32(2)))
}
