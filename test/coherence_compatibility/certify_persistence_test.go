/*
 * Copyright (c) 2021, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 *
 */

package compatibility

import (
	goctx "context"
	"encoding/json"
	"fmt"
	. "github.com/onsi/gomega"
	v1 "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	corev1 "k8s.io/api/core/v1"
	"testing"
	"time"
)

// Deploy a Coherence resource with persistence enabled (this should enable active persistence).
// A PVC should be created for the StatefulSet. Create data in some caches, delete the deployment,
// re-deploy the deployment and assert that the data is recovered.
func TestActivePersistenceScaleUpAndDown(t *testing.T) {
	var yamlFile = "persistence-active-1.yaml"
	var pVolName = "persistence-volume"

	// Ensure that everything is cleaned up after the test!
	testContext.CleanupAfterTest(t)
	g := NewGomegaWithT(t)

	ns := helper.GetTestNamespace()
	deployment, pods := ensurePods(g, yamlFile, ns)

	// check the pvc is created for the given volume
	pvcName := ""
	for _, vol := range pods[0].Spec.Volumes {
		if vol.Name == pVolName {
			if vol.PersistentVolumeClaim != nil {
				pvcName = vol.PersistentVolumeClaim.ClaimName
			}
			break
		}
	}

	// check the pvc is created
	g.Expect(pvcName).NotTo(Equal(""))
	pvc := testContext.KubeClient.CoreV1().PersistentVolumeClaims(pvcName)
	g.Expect(pvc).ShouldNot(BeNil())

	// create data in some caches
	err := helper.StartCanary(testContext, ns, deployment.GetName())
	g.Expect(err).NotTo(HaveOccurred())

	// Start with one replica
	//_, err = helper.WaitForStatefulSetForDeployment(f.KubeClient, ns, &deployment, time.Second*10, time.Minute*5, t)
	//g.Expect(err).NotTo(HaveOccurred())

	// Scale Up to three
	err = scale(t, ns, deployment.Name, 3)
	g.Expect(err).NotTo(HaveOccurred())
	_, err = helper.WaitForStatefulSet(testContext, ns, deployment.Name, 3, time.Second*10, time.Minute*5)
	g.Expect(err).NotTo(HaveOccurred())

	// Scale down to one
	err = scale(t, ns, deployment.Name, 1)
	g.Expect(err).NotTo(HaveOccurred())
	_, err = helper.WaitForStatefulSet(testContext, ns, deployment.Name, 1, time.Second*10, time.Minute*5)
	g.Expect(err).NotTo(HaveOccurred())

	// assert that the data is recovered
	err = helper.CheckCanaryEventuallyGood(testContext, ns, deployment.GetName())
	g.Expect(err).NotTo(HaveOccurred())

	// cleanup the data
	_ = helper.ClearCanary(testContext, ns, deployment.GetName())
}

// Deploy a Coherence resource with persistence enabled (this should enable active persistence).
// A PVC should be created for the StatefulSet. Create data in some caches, scale down to 0 the deployment,
// re-deploy the deployment and assert that the data is recovered.
func TestActivePersistenceScaleDownAndUp(t *testing.T) {
	var yamlFile = "persistence-active-3.yaml"
	var pVolName = "persistence-volume"

	// Ensure that everything is cleaned up after the test!
	testContext.CleanupAfterTest(t)
	g := NewGomegaWithT(t)

	ns := helper.GetTestNamespace()

	deployment, pods := ensurePods(g, yamlFile, ns)

	// check the pvc is created for the given volume
	pvcName := ""
	for _, vol := range pods[0].Spec.Volumes {
		if vol.Name == pVolName {
			if vol.PersistentVolumeClaim != nil {
				pvcName = vol.PersistentVolumeClaim.ClaimName
			}
			break
		}
	}

	// check the pvc is created
	g.Expect(pvcName).NotTo(Equal(""))
	pvc := testContext.KubeClient.CoreV1().PersistentVolumeClaims(pvcName)
	g.Expect(pvc).ShouldNot(BeNil())

	// create data in some caches
	err := helper.StartCanary(testContext, ns, deployment.GetName())
	g.Expect(err).NotTo(HaveOccurred())

	// Start with three replicas
	//_, err = helper.WaitForStatefulSetForDeployment(f.KubeClient, ns, &deployment, time.Second*10, time.Minute*5, t)
	//g.Expect(err).NotTo(HaveOccurred())

	// Scale Down to One
	err = scale(t, ns, deployment.Name, 1)
	g.Expect(err).NotTo(HaveOccurred())
	_, err = helper.WaitForStatefulSet(testContext, ns, deployment.Name, 1, time.Second*10, time.Minute*5)
	g.Expect(err).NotTo(HaveOccurred())

	// Scale back up to Three
	err = scale(t, ns, deployment.Name, 3)
	g.Expect(err).NotTo(HaveOccurred())
	_, err = helper.WaitForStatefulSet(testContext, ns, deployment.Name, 3, time.Second*10, time.Minute*5)
	g.Expect(err).NotTo(HaveOccurred())

	// assert that the data is recovered
	err = helper.CheckCanaryEventuallyGood(testContext, ns, deployment.GetName())
	g.Expect(err).NotTo(HaveOccurred())

	// cleanup the data
	_ = helper.ClearCanary(testContext, ns, deployment.GetName())
}

func ensurePods(g *GomegaWithT, yamlFile, ns string) (v1.Coherence, []corev1.Pod) {
	deployment, err := helper.NewSingleCoherenceFromYaml(ns, yamlFile)
	g.Expect(err).NotTo(HaveOccurred())

	d, _ := json.Marshal(deployment)
	fmt.Printf("Persistence Test installing deployment:\n%s\n", string(d))

	err = testContext.Client.Create(goctx.TODO(), &deployment)
	g.Expect(err).NotTo(HaveOccurred())

	_, err = helper.WaitForStatefulSetForDeployment(testContext, ns, &deployment, helper.RetryInterval, helper.Timeout)
	g.Expect(err).NotTo(HaveOccurred())

	// Get the list of Pods
	pods, err := helper.ListCoherencePodsForDeployment(testContext, ns, deployment.GetName())
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(len(pods)).Should(BeNumerically(">", 0))

	return deployment, pods
}
