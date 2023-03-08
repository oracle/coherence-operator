/*
 * Copyright (c) 2021, 2023, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package local

import (
	"context"
	"fmt"
	. "github.com/onsi/gomega"
	cohv1 "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/pointer"
	"testing"
	"time"
)

// Test that a deployment gets deleted
func TestDeleteDeployment(t *testing.T) {
	// Make sure we defer clean-up when we're done!!
	testContext.CleanupAfterTest(t)
	g := NewGomegaWithT(t)

	namespace := helper.GetTestNamespace()

	deployment, err := helper.NewSingleCoherenceFromYaml(namespace, "deployment-minimal.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	name := deployment.Name

	err = testContext.Client.Create(context.TODO(), &deployment)
	g.Expect(err).NotTo(HaveOccurred())

	// create a DeploymentStateCondition that checks a deployment's replica count
	condition := helper.ReplicaCountCondition(deployment.GetReplicas())

	// wait for the deployment to match the condition
	_, err = helper.WaitForCoherenceCondition(testContext, namespace, name, condition, time.Second*10, time.Minute*5)
	g.Expect(err).NotTo(HaveOccurred())

	err = testContext.Client.Delete(context.TODO(), &deployment)
	g.Expect(err).NotTo(HaveOccurred())

	err = helper.WaitForDelete(testContext, &deployment)
	g.Expect(err).NotTo(HaveOccurred())
}

// Test that a deployment with zero ready pods gets deleted
func TestDeleteDeploymentWithZeroReadyPods(t *testing.T) {
	// Make sure we defer clean-up when we're done!!
	testContext.CleanupAfterTest(t)
	g := NewGomegaWithT(t)

	namespace := helper.GetTestNamespace()

	deployment, err := helper.NewSingleCoherenceFromYaml(namespace, "deployment-minimal.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	name := deployment.Name

	// set the image to an invalid name so that Pods never start
	deployment.Spec.Image = pointer.String("invalid-image:1.0.0")

	err = testContext.Client.Create(context.TODO(), &deployment)
	g.Expect(err).NotTo(HaveOccurred())

	// wait for the StatefulSet to appear - it will have zero ready replicas
	_, err = helper.WaitForStatefulSet(testContext, namespace, name, 0, time.Second*10, time.Minute*5)
	g.Expect(err).NotTo(HaveOccurred())

	// Delete the Coherence deployment
	err = testContext.Client.Delete(context.TODO(), &deployment)
	g.Expect(err).NotTo(HaveOccurred())

	err = helper.WaitForDelete(testContext, &deployment)
	g.Expect(err).NotTo(HaveOccurred())
}

// Test that a deployment where one or more Pods in a Coherence resource cannot be created due to a lack of resources
func TestDeleteDeploymentWithAllPendingPods(t *testing.T) {
	// Make sure we defer clean-up when we're done!!
	testContext.CleanupAfterTest(t)
	g := NewGomegaWithT(t)

	namespace := helper.GetTestNamespace()

	deployment, err := helper.NewSingleCoherenceFromYaml(namespace, "deployment-no-resources.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	err = testContext.Client.Create(context.TODO(), &deployment)
	g.Expect(err).NotTo(HaveOccurred())

	// wait for all pods to be in pending state
	var fieldSelector = fmt.Sprintf("%s=%s", "status.phase", corev1.PodPending)
	labelSelector := fmt.Sprintf("%s=%s", cohv1.LabelComponent, cohv1.LabelComponentCoherencePod)
	_, err = helper.WaitForPodsWithLabelAndField(testContext, namespace, labelSelector, fieldSelector, int(deployment.GetReplicas()), time.Second*10, time.Minute*2)
	g.Expect(err).NotTo(HaveOccurred())

	// Delete the Coherence deployment
	err = testContext.Client.Delete(context.TODO(), &deployment)
	g.Expect(err).NotTo(HaveOccurred())

	err = helper.WaitForDelete(testContext, &deployment)
	g.Expect(err).NotTo(HaveOccurred())
}
