/*
 * Copyright (c) 2019, 2020, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package local

import (
	"context"
	"fmt"
	. "github.com/onsi/gomega"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	coh "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	appsv1 "k8s.io/api/apps/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"testing"
	"time"
)

// ----- tests --------------------------------------------------------------

// Test that a deployment works using the minimal valid yaml for a Coherence
func TestMinimalDeployment(t *testing.T) {
	AssertDeployments(t, "deployment-minimal.yaml")
}

// Test that a deployment works with a replica count of 1
func TestDeploymentWithOneReplica(t *testing.T) {
	AssertDeployments(t, "deployment-one-replica.yaml")
}

// Test that a deployment works using the a yaml file containing two Coherence
// specs that have the same cluster name.
func TestTwoDeploymentsOneCluster(t *testing.T) {
	AssertDeployments(t, "deployment-multi.yaml")
}

// Test that two deployments with dependencies start in the correct order
func TestStartQuorumRequireAllPodsReady(t *testing.T) {
	g := NewGomegaWithT(t)

	// Start the two deployments
	deployments, pods := AssertDeployments(t, "deployment-with-start-quorum.yaml")
	data, ok := deployments["data"]
	g.Expect(ok).To(BeTrue(), "did not find expected 'data' deployment")
	test, ok := deployments["test"]
	g.Expect(ok).To(BeTrue(), "did not find expected 'test' deployment")

	g.Expect(data.Status.Phase).To(Equal(coh.ConditionTypeReady))
	g.Expect(test.Status.Phase).To(Equal(coh.ConditionTypeReady))

	ready := data.Status.Conditions.GetCondition(coh.ConditionTypeReady)
	g.Expect(ready).NotTo(BeNil())
	created := test.Status.Conditions.GetCondition(coh.ConditionTypeCreated)
	g.Expect(created).NotTo(BeNil())
	// created time should not be before ready time
	g.Expect(created.LastTransitionTime.Time.Before(ready.LastTransitionTime.Time)).To(BeFalse())

	// earliest test Pod scheduled should not be before last data Pod ready
	dataPodReady := helper.GetLastPodReadyTime(pods, "data")
	testPodScheduled := helper.GetFirstPodScheduledTime(pods, "test")
	g.Expect(testPodScheduled.Before(&dataPodReady)).To(BeFalse())
}

// Test that two deployments with dependency on single Pod ready start in the correct order
func TestStartQuorumRequireOnePodReady(t *testing.T) {
	g := NewGomegaWithT(t)

	// Start the two deployments
	deployments, pods := AssertDeployments(t, "deployment-with-start-quorum-one-pod.yaml")
	data, ok := deployments["data"]
	g.Expect(ok).To(BeTrue(), "did not find expected 'data' deployment")
	test, ok := deployments["test"]
	g.Expect(ok).To(BeTrue(), "did not find expected 'test' deployment")

	g.Expect(data.Status.Phase).To(Equal(coh.ConditionTypeReady))
	g.Expect(test.Status.Phase).To(Equal(coh.ConditionTypeReady))

	// Get the time the first data Pod was ready
	dataPodReady := helper.GetFirstPodReadyTime(pods, "data")

	created := test.Status.Conditions.GetCondition(coh.ConditionTypeCreated)
	g.Expect(created).NotTo(BeNil())
	// created time should not be before first data Pod ready time
	g.Expect(created.LastTransitionTime.Time.Before(dataPodReady.Time)).To(BeFalse(),
		fmt.Sprintf("Expected test created %s after data ready %s", created.LastTransitionTime.Time.String(), dataPodReady.Time.String()))

	// earliest test Pod scheduled should not be before last data Pod ready
	testPodScheduled := helper.GetFirstPodScheduledTime(pods, "test")
	g.Expect(testPodScheduled.Before(&dataPodReady)).To(BeFalse())
}

func TestTwoDeploymentsOneClusterWithWKAExclusion(t *testing.T) {
	AssertDeployments(t, "deployment-with-wka-exclusion.yaml")
}

// Test that a cluster can be created with zero replicas.
func TestDeploymentWithZeroReplicas(t *testing.T) {
	// initialise Gomega so we can use matchers
	g := NewGomegaWithT(t)
	f := framework.Global

	// Create the Operator SDK test context (this will deploy the Operator)
	ctx := helper.CreateTestContext(t)
	// Make sure we defer clean-up (uninstall the operator) when we're done
	defer helper.DumpOperatorLogsAndCleanup(t, ctx)

	// Get the test namespace
	namespace, err := ctx.GetWatchNamespace()
	g.Expect(err).NotTo(HaveOccurred())

	deployments, err := helper.NewCoherenceFromYaml(namespace, "deployment-with-zero-replicas.yaml")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(len(deployments)).To(Equal(1))
	deployment := deployments[0]

	// deploy the Coherence
	err = f.Client.Create(context.TODO(), &deployment, helper.DefaultCleanup(ctx))
	g.Expect(err).NotTo(HaveOccurred())

	// The deployment should eventually be in the Stopped phase
	condition := helper.StatusPhaseCondition(coh.ConditionTypeStopped)
	_, err = helper.WaitForCoherenceCondition(f, namespace, deployment.Name, condition, time.Second, time.Minute*5, t)
	g.Expect(err).NotTo(HaveOccurred())

	// There should be no StatefulSet
	sts := &appsv1.StatefulSet{}
	err = f.Client.Get(context.TODO(), deployment.GetNamespacedName(), sts)
	g.Expect(err).To(HaveOccurred())
	g.Expect(apierrors.IsNotFound(err)).To(BeTrue())
}
