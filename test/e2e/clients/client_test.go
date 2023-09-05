/*
 * Copyright (c) 2021, 2023, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package clients

import (
	. "github.com/onsi/gomega"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
	"testing"
	"time"
)

// TestSimpleClients is the parent test function for simple client tests
func TestSimpleClients(t *testing.T) {
	g := NewGomegaWithT(t)
	testContext.CleanupAfterTest(t)

	// Start the Coherence cluster
	cluster, err := DeployTestCluster(testContext, t, "storage.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	// Create all the child test cases.
	testCases := []ClientTestCase{
		// Simple Extend test direct connection
		{ClientType: ClientTypeExtend, Name: "ExtendInternalDirect", Cluster: cluster, Test: simpleClientTest},
		// Simple Extend test name-service connection
		{ClientType: ClientTypeExtend, Name: "ExtendInternalNS", Cluster: cluster, Test: simpleClientTest, CacheConfig: "test-cache-config-ns.xml"},
		// Simple gRPC test
		{ClientType: ClientTypeGrpc, Name: "GrpcInternal", Cluster: cluster, Test: simpleClientTest},
	}

	// Execute all the child test cases.
	Execute(t, testCases)
}

// simpleClientTest runs a simple Extend or gRPC client test where the client runs inside k8s
// The test client is run in a K8s Job, which will just run and complete with an exit code of success or failure.
func simpleClientTest(t *testing.T, tc ClientTestCase) {
	g := NewGomegaWithT(t)

	ns := helper.GetTestClientNamespace()

	// Use the lower-case test name as the Job name, escaping any / to a dash
	jobName := strings.ToLower(strings.ReplaceAll(t.Name(), "/", "-"))

	// ensure we dump logs and other information when we exit this method
	defer helper.DumpState(testContext, ns, t.Name())

	// Get the k8s API Jobs client
	client := testContext.KubeClient.BatchV1().Jobs(ns)

	// ensure we delete the Job when the test finishes
	t.Cleanup(func() {
		_ = helper.DeleteJob(testContext, ns, jobName)
	})

	// ensure there is no Job left from a previous test
	_ = helper.DeleteJob(testContext, ns, jobName)

	// Create the Job to run the client
	job := CreateClientJob(ns, jobName, tc)
	_, err := client.Create(testContext.Context, job, metav1.CreateOptions{})
	g.Expect(err).NotTo(HaveOccurred())

	// All Pods for the Job will have this label so that we can find them
	label := "job-name=" + jobName

	// Wait for Pods to appear
	t.Logf("Looking for Pods with label '%s' in namespace %s\n", label, ns)
	pods, err := helper.WaitForPodsWithLabel(testContext, ns, label, 1, 10*time.Second, 5*time.Minute)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(pods).NotTo(BeEmpty())

	// Wait for the Job to complete
	err = helper.WaitForJobCompletion(testContext, ns, pods[0].Name, 10*time.Second, 5*time.Minute)
	g.Expect(err).NotTo(HaveOccurred())
}

/*
// simpleExternalClientTest runs a simple Extend or gRPC client test where the client runs outside k8s
// The test client is run in a container on the local host, which will just run and complete with an exit code of success or failure.
func simpleExternalClientTest(t *testing.T, tc ClientTestCase) {
	g := NewGomegaWithT(t)
	host, port, err := findIngress(tc, string(tc.ClientType))
	g.Expect(err).NotTo(HaveOccurred())

	cfg := tc.CacheConfig
	if cfg == "" {
		cfg = "test-cache-config.xml"
	}

	image := helper.GetClientImage()

	var hostEnvVar string
	var portEnvVar string

	switch tc.ClientType {
	case ClientTypeExtend:
		hostEnvVar = "COHERENCE_EXTEND_ADDRESS"
		portEnvVar = "COHERENCE_EXTEND_PORT"
		break
	case ClientTypeGrpc:
		hostEnvVar = "COHERENCE_GRPC_CHANNELS_DEFAULT_HOST"
		portEnvVar = "COHERENCE_GRPC_CHANNELS_DEFAULT_PORT"
		break
	}

	t.Logf("docker run -it --rm -e CLIENT_TYPE=%s -e COHERENCE_CACHECONFIG=%s -e COHERENCE_DISTRIBUTED_LOCALSTORAGE=false -e %s=%s -e %s=%d %s", string(tc.ClientType), cfg, hostEnvVar, host, portEnvVar, port, image)
	time.Sleep(20*time.Second)
}

func findIngress(tc ClientTestCase, portName string) (string, int32, error) {
	ingress, found := tc.Cluster.ServiceIngress[portName]
	if !found {
		return "", -1, fmt.Errorf("could not find ingress configuration for port %s", portName)
	}

	for _, i := range ingress {
		if i.IP != "" {
			for _, p := range i.Ports {
				if p.Error == nil {
					return "127.0.0.1", p.Port, nil
				}
			}
		}
	}

	return "", -1, fmt.Errorf("could not find a working ingress configuration for port %s", portName)
}
*/
