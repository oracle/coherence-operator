/*
 * Copyright (c) 2019, 2020, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package local

import (
	"context"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
	"time"

	. "github.com/onsi/gomega"
)

func TestMinimalCoherenceClusterWKA(t *testing.T) {
	assertWKA(t, "cluster-minimal.yaml")
}

func TestExculsionsWKA(t *testing.T) {
	assertWKA(t, "cluster-minimal.yaml")
}

func assertWKA(t *testing.T, yamlFile string) {
	g := NewGomegaWithT(t)
	f := framework.Global

	// Create the Operator SDK test context (this will deploy the Operator)
	ctx := helper.CreateTestContext(t)
	// Make sure we defer clean-up (uninstall the operator) when we're done
	defer helper.DumpOperatorLogsAndCleanup(t, ctx)

	// Get the test namespace
	namespace, err := ctx.GetNamespace()
	g.Expect(err).NotTo(HaveOccurred())

	cluster, err := helper.NewCoherenceClusterFromYaml(namespace, yamlFile)
	g.Expect(err).NotTo(HaveOccurred())

	// deploy the CoherenceCluster
	err = f.Client.Create(context.TODO(), &cluster, helper.DefaultCleanup(ctx))
	g.Expect(err).NotTo(HaveOccurred())

	roles := cluster.GetRoles()

	// Assert that a StatefulSet of the correct number or replicas is created for each role in the cluster
	for _, role := range roles {
		// Wait for the StatefulSet for the role to be ready - wait five minutes max
		sts, err := helper.WaitForStatefulSetForRole(f.KubeClient, namespace, &cluster, role, time.Second*10, time.Minute*5, t)
		g.Expect(err).NotTo(HaveOccurred())
		g.Expect(sts.Status.ReadyReplicas).To(Equal(role.GetReplicas()))
	}

	// Obtain the expected WKA list of Pod IP addresses
	var wkaPods []string
	for roleName, r := range cluster.GetRoles() {
		if r.Coherence.IsWKAMember() {
			pods, err := helper.ListCoherencePodsForRole(f.KubeClient, namespace, cluster.Name, roleName)
			g.Expect(err).NotTo(HaveOccurred())
			for _, pod := range pods {
				wkaPods = append(wkaPods, pod.Status.PodIP)
			}
		}
	}

	// Verify that the WKA service endpoints list has all of the required the Pod IP addresses.
	serviceName := cluster.GetWkaServiceName()
	ep, err := f.KubeClient.CoreV1().Endpoints(namespace).Get(serviceName, metav1.GetOptions{})
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(len(ep.Subsets)).NotTo(BeZero())

	subset := ep.Subsets[0]
	g.Expect(len(subset.Addresses)).To(Equal(len(wkaPods)))
	var actualWKA []string
	for _, address := range subset.Addresses {
		actualWKA = append(actualWKA, address.IP)
	}
	g.Expect(actualWKA).To(ConsistOf(wkaPods))
}
