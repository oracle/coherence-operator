/*
 * Copyright (c) 2021, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package remote

import (
	"context"
	. "github.com/onsi/gomega"
	coh "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	"time"

	"testing"
)

func TestRollingUpgrade(t *testing.T) {
	// Make sure we defer clean-up when we're done!!
	testContext.CleanupAfterTest(t)

	g := NewGomegaWithT(t)

	namespace := helper.GetTestNamespace()

	t.Log("Deploying initial version of Coherence cluster")
	// Do the initial deployment
	deployments, _ := helper.AssertDeployments(testContext, t, "rolling-upgrade.yaml")
	// Get the expected single deployment from the returned map
	deployment, ok := deployments["rolling-cluster"]
	g.Expect(ok).To(BeTrue())

	// Get the latest state for the deployment
	upgrade := coh.Coherence{}
	err := testContext.Client.Get(context.TODO(), deployment.GetNamespacedName(), &upgrade)
	g.Expect(err).NotTo(HaveOccurred())

	// Upgrade the version label and JVM Heap
	updatedHeap := "500m"
	t.Log("Deploying updated version of Coherence cluster")
	upgrade.Spec.Labels["version"] = "two"
	upgrade.Spec.JVM.Memory.HeapSize = &updatedHeap
	err = testContext.Client.Update(context.TODO(), &upgrade)
	g.Expect(err).NotTo(HaveOccurred())

	// wait for the expected updated Pods
	t.Log("Waiting for all Pods to be updated")
	pods, err := helper.WaitForPodsWithLabel(testContext, namespace, "version=two", 3, time.Second*10, time.Minute*5)
	g.Expect(err).NotTo(HaveOccurred())

	for _, pod := range pods {
		for _, ev := range pod.Spec.Containers[0].Env {
			if ev.Name == coh.EnvVarJvmMemoryHeap {
				g.Expect(ev.Value).To(Equal(updatedHeap), "Expected heap incorrect for Pod "+pod.Name)
			}
		}
	}
}

// Test the scenario where a cluster is updated that causes a scaled down to 1 and a
// rolling upgrade.
// The Operator will first scale down, then it will upgrade the single member that is left.
// Without persistence this causes data loss, but in this case persistence is active
// so there should be no data loss after the Pod is upgraded.
func TestScalingWithRollingUpgrade(t *testing.T) {
	// Ensure that everything is cleaned up after the test!
	testContext.CleanupAfterTest(t)

	g := NewGomegaWithT(t)

	namespace := helper.GetTestNamespace()

	deployment, err := helper.NewSingleCoherenceFromYaml(namespace, "rolling-upgrade-with-persistence.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	name := deployment.Name

	// Start with three replicas
	deployment.SetReplicas(3)

	// Install the deployment
	installSimpleDeployment(t, deployment)

	// Load the canary data
	err = helper.StartCanary(testContext, namespace, name)
	g.Expect(err).NotTo(HaveOccurred())

	// Set the replicas to one and add a label so that we also do a rolling upgrade
	installed, err := helper.GetCoherence(testContext, namespace, name)
	g.Expect(err).NotTo(HaveOccurred())
	installed.SetReplicas(1)
	installed.Spec.Labels = make(map[string]string)
	installed.Spec.Labels["updated"] = "true"

	// trigger the update
	err = testContext.Client.Update(context.TODO(), installed)
	g.Expect(err).NotTo(HaveOccurred())

	// Wait for the remaining Pod to be updated
	_, err = helper.WaitForPodsWithLabel(testContext, namespace, "updated=true", 1, time.Second*5, time.Minute*5)
	g.Expect(err).NotTo(HaveOccurred())

	// Wait for ready
	_, err = helper.WaitForStatefulSet(testContext, namespace, name, 1, time.Second*5, time.Minute*5)
	g.Expect(err).NotTo(HaveOccurred())

	// Check the canary data
	err = helper.CheckCanary(testContext, namespace, name)
	g.Expect(err).NotTo(HaveOccurred())
}
