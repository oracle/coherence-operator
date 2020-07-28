/*
 * Copyright (c) 2020, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */
package local


import (
	"context"
	. "github.com/onsi/gomega"
	coh "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"testing"
	"time"
)

// Test that a cluster can be created using the specified yaml.
func AssertDeployments(t *testing.T, yamlFile string) (map[string]coh.Coherence, []corev1.Pod) {
	// Make sure we defer clean-up (uninstall the operator) when we're done
	defer helper.DumpOperatorLogs(t, testContext)
	return AssertDeploymentsWithContext(t, yamlFile)
}

// Test that a cluster can be created using the specified yaml.
func AssertDeploymentsWithContext(t *testing.T, yamlFile string) (map[string]coh.Coherence, []corev1.Pod) {
	// initialise Gomega so we can use matchers
	g := NewGomegaWithT(t)

	// Get the test namespace
	namespace := helper.GetTestNamespace()

	deployments, err := helper.NewCoherenceFromYaml(namespace, yamlFile)
	g.Expect(err).NotTo(HaveOccurred())

	// we must have at least one deployment
	g.Expect(len(deployments)).NotTo(BeZero())

	// assert all deployments have the same cluster name
	clusterName := deployments[0].GetCoherenceClusterName()
	for _, d := range deployments {
		g.Expect(d.GetCoherenceClusterName()).To(Equal(clusterName))
	}

	// work out the expected cluster size
	expectedClusterSize := 0
	expectedWkaSize := 0
	for _, d := range deployments {
		t.Logf("Deployment %s has replica count %d", d.Name, d.GetReplicas())
		replicas := int(d.GetReplicas())
		expectedClusterSize += replicas
		if d.Spec.Coherence.IsWKAMember() {
			expectedWkaSize += replicas
		}
	}
	t.Logf("Expected cluster size is %d", expectedClusterSize)

	for _, d := range deployments {
		t.Logf("Deploying %s", d.Name)
		// deploy the Coherence resource
		err = testContext.Client.Create(context.TODO(), &d)
		g.Expect(err).NotTo(HaveOccurred())
	}

	// Assert that a StatefulSet of the correct number or replicas is created for each roleSpec in the cluster
	for _, d := range deployments {
		t.Logf("Waiting for StatefulSet for deployment %s", d.Name)
		// Wait for the StatefulSet for the roleSpec to be ready - wait five minutes max
		sts, err := helper.WaitForStatefulSet(testContext, namespace, d.Name, d.GetReplicas(), time.Second*10, time.Minute*5)
		g.Expect(err).NotTo(HaveOccurred())
		g.Expect(sts.Status.ReadyReplicas).To(Equal(d.GetReplicas()))
		t.Logf("Have StatefulSet for deployment %s", d.Name)
	}

	// Get all of the Pods in the cluster
	t.Logf("Getting all Pods for cluster '%s'", clusterName)
	pods, err := helper.ListCoherencePodsForCluster(testContext, namespace, clusterName)
	g.Expect(err).NotTo(HaveOccurred())
	t.Logf("Found %d Pods for cluster '%s'", len(pods), clusterName)

	// assert that the correct number of Pods is returned
	g.Expect(len(pods)).To(Equal(expectedClusterSize))

	// Verify that the WKA service has the same number of endpoints as the cluster size.
	serviceName := deployments[0].GetWkaServiceName()
	
	ep, err := testContext.KubeClient.CoreV1().Endpoints(namespace).Get(context.TODO(), serviceName, metav1.GetOptions{})
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(len(ep.Subsets)).NotTo(BeZero())

	subset := ep.Subsets[0]
	g.Expect(len(subset.Addresses)).To(Equal(expectedWkaSize))

	m := make(map[string]coh.Coherence)
	for _, d := range deployments {
		opts := client.ObjectKey{Namespace: namespace, Name: d.Name}
		dpl := coh.Coherence{}
		err = testContext.Client.Get(context.TODO(), opts, &dpl)
		g.Expect(err).NotTo(HaveOccurred())
		m[dpl.Name] = dpl
	}

	// Obtain the expected WKA list of Pod IP addresses
	var wkaPods []string
	for _, d := range deployments {
		if d.Spec.Coherence.IsWKAMember() {
			pods, err := helper.ListCoherencePodsForDeployment(testContext, d.Namespace, d.Name)
			g.Expect(err).NotTo(HaveOccurred())
			for _, pod := range pods {
				wkaPods = append(wkaPods, pod.Status.PodIP)
			}
		}
	}

	// Verify that the WKA service endpoints list for each deployment has all of the required the Pod IP addresses.
	for _, d := range deployments {
		serviceName := d.GetWkaServiceName()
		ep, err = testContext.KubeClient.CoreV1().Endpoints(namespace).Get(context.TODO(), serviceName, metav1.GetOptions{})
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

	return m, pods
}
