/*
 * Copyright (c) 2019, 2020, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package local

import (
	"context"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	coh "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"testing"
	"time"

	. "github.com/onsi/gomega"
)

// ----- tests --------------------------------------------------------------

func TestMinimalCoherenceCluster(t *testing.T) {
	assertCluster(t, "cluster-minimal.yaml", map[string]int32{coh.DefaultRoleName: coh.DefaultReplicas})
}

func TestNoRoleOneReplica(t *testing.T) {
	assertCluster(t, "cluster-no-role-one-replica.yaml", map[string]int32{coh.DefaultRoleName: 1})
}

func TestOneRoleDefaultReplicas(t *testing.T) {
	assertCluster(t, "cluster-one-role-default-replica.yaml", map[string]int32{"data": coh.DefaultReplicas})
}

func TestOneRoleOneReplicas(t *testing.T) {
	assertCluster(t, "cluster-one-role-one-replica.yaml", map[string]int32{coh.DefaultRoleName: 1})
}

func TestOneRoleTwoReplicas(t *testing.T) {
	assertCluster(t, "cluster-one-role-two-replica.yaml", map[string]int32{"data": 2})
}

func TestTwoRolesDefaultReplicas(t *testing.T) {
	assertCluster(t, "cluster-two-roles-default-replica.yaml", map[string]int32{"data": coh.DefaultReplicas, "proxy": coh.DefaultReplicas})
}

func TestTwoRolesOneReplicas(t *testing.T) {
	assertCluster(t, "cluster-two-roles-one-replica.yaml", map[string]int32{"data": 1, "proxy": 1})
}

func TestTwoRolesTwoReplicas(t *testing.T) {
	assertCluster(t, "cluster-two-roles-different-replica.yaml", map[string]int32{"data": 2, "proxy": 1})
}

func TestStartQuorumDependentRoleReadySingleRole(t *testing.T) {
	g := NewGomegaWithT(t)

	cluster, pods := assertCluster(t, "cluster-ready-quorum.yaml", map[string]int32{"data": 2, "test": 1})
	ready := helper.GetLastPodReadyTime(pods, "data")
	scheduled := helper.GetFirstPodScheduledTime(pods, "test")

	g.Expect(scheduled.Before(&ready)).To(BeFalse())

	dataStatus := cluster.GetRoleStatus("data")
	testStatus := cluster.GetRoleStatus("test")
	dataReady := dataStatus.GetCondition(coh.RoleStatusReady)
	testCreated := testStatus.GetCondition(coh.RoleStatusCreated)

	g.Expect(testCreated.LastTransitionTime.Before(&dataReady.LastTransitionTime)).To(BeFalse())
}

func TestStartQuorumDependentRoleReadyTwoRoles(t *testing.T) {
	g := NewGomegaWithT(t)

	cluster, pods := assertCluster(t, "cluster-ready-quorum-two-roles.yaml", map[string]int32{"data": 2, "proxy": 2, "test": 1})
	readyData := helper.GetLastPodReadyTime(pods, "data")
	readyProxy := helper.GetLastPodReadyTime(pods, "proxy")
	scheduled := helper.GetFirstPodScheduledTime(pods, "test")

	g.Expect(scheduled.Before(&readyData)).To(BeFalse())
	g.Expect(scheduled.Before(&readyProxy)).To(BeFalse())

	dataStatus := cluster.GetRoleStatus("data")
	dataReady := dataStatus.GetCondition(coh.RoleStatusReady)
	proxyStatus := cluster.GetRoleStatus("proxy")
	proxyReady := proxyStatus.GetCondition(coh.RoleStatusReady)
	testStatus := cluster.GetRoleStatus("test")
	testCreated := testStatus.GetCondition(coh.RoleStatusCreated)

	g.Expect(testCreated.LastTransitionTime.Before(&dataReady.LastTransitionTime)).To(BeFalse())
	g.Expect(testCreated.LastTransitionTime.Before(&proxyReady.LastTransitionTime)).To(BeFalse())
}

func TestStartQuorumDependentRoleReadyChained(t *testing.T) {
	g := NewGomegaWithT(t)

	cluster, pods := assertCluster(t, "cluster-ready-quorum-chained.yaml", map[string]int32{"data": 2, "proxy": 2, "test": 1})

	readyData := helper.GetLastPodReadyTime(pods, "data")
	scheduledProxy := helper.GetFirstPodScheduledTime(pods, "proxy")

	g.Expect(scheduledProxy.Before(&readyData)).To(BeFalse())

	dataStatus := cluster.GetRoleStatus("data")
	proxyStatus := cluster.GetRoleStatus("proxy")
	dataReady := dataStatus.GetCondition(coh.RoleStatusReady)
	proxyCreated := proxyStatus.GetCondition(coh.RoleStatusCreated)

	g.Expect(proxyCreated.LastTransitionTime.Before(&dataReady.LastTransitionTime)).To(BeFalse())

	readyProxy := helper.GetLastPodReadyTime(pods, "proxy")
	scheduledTest := helper.GetFirstPodScheduledTime(pods, "test")

	g.Expect(scheduledTest.Before(&readyProxy)).To(BeFalse())

	proxyReady := proxyStatus.GetCondition(coh.RoleStatusReady)
	testStatus := cluster.GetRoleStatus("test")
	testCreated := testStatus.GetCondition(coh.RoleStatusCreated)

	g.Expect(testCreated.LastTransitionTime.Before(&proxyReady.LastTransitionTime)).To(BeFalse())
}

func TestStartQuorumDependentRoleOnePodReadySingleRole(t *testing.T) {
	g := NewGomegaWithT(t)

	_, pods := assertCluster(t, "cluster-ready-quorum-one-pod.yaml", map[string]int32{"data": 2, "test": 1})
	ready := helper.GetFirstPodReadyTime(pods, "data")
	scheduled := helper.GetFirstPodScheduledTime(pods, "test")

	g.Expect(scheduled.Before(&ready)).To(BeFalse())
}

func TestStartQuorumDependentRoleOnePodReadyTwoRoles(t *testing.T) {
	g := NewGomegaWithT(t)

	_, pods := assertCluster(t, "cluster-ready-quorum-two-roles-one-pod.yaml", map[string]int32{"data": 2, "proxy": 2, "test": 1})
	readyData := helper.GetFirstPodReadyTime(pods, "data")
	readyProxy := helper.GetFirstPodReadyTime(pods, "proxy")
	scheduled := helper.GetFirstPodScheduledTime(pods, "test")

	g.Expect(scheduled.Before(&readyData)).To(BeFalse())
	g.Expect(scheduled.Before(&readyProxy)).To(BeFalse())
}

// ----- helpers ------------------------------------------------------------

// Test that a cluster can be created using the specified yaml.
func assertCluster(t *testing.T, yamlFile string, expectedRoles map[string]int32) (*coh.CoherenceCluster, []corev1.Pod) {
	// initialise Gomega so we can use matchers
	g := NewGomegaWithT(t)
	f := framework.Global

	// work out the total expected roles and cluster size
	totalRoles := 0
	clusterSize := 0
	for _, size := range expectedRoles {
		clusterSize = clusterSize + int(size)
		if size > 0 {
			totalRoles = totalRoles + 1
		}
	}

	// Create the Operator SDK test context (this will deploy the Operator)
	ctx := helper.CreateTestContext(t)
	// Make sure we defer clean-up (uninstall the operator) when we're done
	defer helper.DumpOperatorLogsAndCleanup(t, ctx)

	// Get the test namespace
	namespace, err := ctx.GetNamespace()
	g.Expect(err).NotTo(HaveOccurred())

	cluster, err := helper.NewCoherenceClusterFromYaml(namespace, yamlFile)

	// verify the cluster size is expected
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(cluster.GetClusterSize()).To(Equal(clusterSize))

	// deploy the CoherenceCluster
	err = f.Client.Create(context.TODO(), &cluster, helper.DefaultCleanup(ctx))
	g.Expect(err).NotTo(HaveOccurred())

	roles := cluster.GetRoles()

	// Assert that a CoherenceRole is created for each role in the cluster
	for _, role := range roles {
		roleName := role.GetFullRoleName(&cluster)
		// Wait for a CoherenceRole to be created
		role, err := helper.WaitForCoherenceRole(f, namespace, roleName, time.Second*10, time.Minute*2, t)
		g.Expect(err).NotTo(HaveOccurred())

		expectedReplicas, found := expectedRoles[role.Spec.GetRoleName()]
		g.Expect(found).To(BeTrue(), "Found Role with unexpected name '"+roleName+"'")
		g.Expect(role.Spec.GetReplicas()).To(Equal(expectedReplicas))
	}

	// Assert that a StatefulSet of the correct number or replicas is created for each role in the cluster
	for _, role := range roles {
		// Wait for the StatefulSet for the role to be ready - wait five minutes max
		sts, err := helper.WaitForStatefulSetForRole(f.KubeClient, namespace, &cluster, role, time.Second*10, time.Minute*5, t)
		g.Expect(err).NotTo(HaveOccurred())
		g.Expect(sts.Status.ReadyReplicas).To(Equal(role.GetReplicas()))
	}

	// Get all of the Pods in the cluster
	pods, err := helper.ListCoherencePodsForCluster(f.KubeClient, namespace, cluster.Name)
	g.Expect(err).NotTo(HaveOccurred())

	// assert that the correct number of Pods is returned
	g.Expect(len(pods)).To(Equal(clusterSize))

	// Verify that the WKA service has the same number of endpoints as the cluster size.
	serviceName := cluster.GetWkaServiceName()
	ep, err := f.KubeClient.CoreV1().Endpoints(namespace).Get(serviceName, metav1.GetOptions{})
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(len(ep.Subsets)).NotTo(BeZero())

	subset := ep.Subsets[0]
	g.Expect(len(subset.Addresses)).To(Equal(clusterSize))

	opts := client.ObjectKey{Namespace: namespace, Name: cluster.Name}
	err = f.Client.Get(context.TODO(), opts, &cluster)
	g.Expect(err).NotTo(HaveOccurred())

	return &cluster, pods
}
