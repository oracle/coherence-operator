/*
 * Copyright (c) 2020, 2025, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package local

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	coh "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/ptr"
	"net/http"
	"os"
	"testing"
	"time"

	. "github.com/onsi/gomega"
)

// A struct used to define a management test case.
type ManagementTestCase struct {
	Deployment    *coh.Coherence
	Name          string
	KeyFile       string
	CertFile      string
	CaCertFile    string
	ShouldSucceed bool
}

// TestManagement is a go test that uses sub-tests (test cases) to basically run the
// same test with different parameters. In this case different Coherence
// configurations with management over rest configured with and without SSL.
func TestManagementOverRest(t *testing.T) {
	helper.SkipIfCoherenceVersionLessThan(t, 12, 2, 1, 4)

	// Make sure we defer clean-up when we're done!!
	testContext.CleanupAfterTest(t)
	// initialise Gomega so we can use matchers
	g := NewWithT(t)

	// Get the test namespace
	namespace := helper.GetTestNamespace()

	// Get the test SSL information (secret name etc.)
	_, ssl, err := helper.GetTestSslSecret()
	g.Expect(err).NotTo(HaveOccurred())

	// require 2-way auth
	ssl.RequireClientCert = ptr.To(true)

	// load the test Coherence from a yaml files
	deploymentWithoutSSL, err := helper.NewSingleCoherenceFromYaml(namespace, "management-test.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	// load the test Coherence resource that used a distroless JIB image from a yaml files
	deploymentJib, err := helper.NewSingleCoherenceFromYaml(namespace, "management-jib-test.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	j, _ := json.Marshal(deploymentWithoutSSL)
	fmt.Printf("Deployment: %s\n", string(j))

	// Copy deploymentWithoutSSL and configure it to use SSL at the Spec level in all deployments
	deploymentSSL := &coh.Coherence{}
	deploymentWithoutSSL.DeepCopyInto(deploymentSSL)

	// Set the SSL settings
	deploymentSSL.Name = "management-ssl"
	deploymentSSL.Spec.Coherence.Management.Enabled = ptr.To(true)
	deploymentSSL.Spec.Coherence.Management.SSL = ssl

	j, _ = json.Marshal(deploymentWithoutSSL)
	fmt.Printf("Deployment: %s\n", string(j))
	j, _ = json.Marshal(deploymentSSL)
	fmt.Printf("DeploymentSSL: %s\n", string(j))

	// Create the test cases
	testCases := []ManagementTestCase{
		{&deploymentWithoutSSL, "PlainHTTP", "", "", "", true},
		{&deploymentJib, "JIB", "", "", "", true},
		{deploymentSSL, "WithSSL", "groot.key", "groot.crt", "guardians-ca.crt", true},
		{deploymentSSL, "ClientHasBadKey", "yondu.key", "groot.crt", "guardians-ca.crt", false},
		{deploymentSSL, "BadCert", "groot.key", "yondu.crt", "guardians-ca.crt", false},
		{deploymentSSL, "BadCaCert", "groot.key", "groot.crt", "ravagers-ca.crt", false},
	}

	// Run the test cases...
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			testManagementOverRest(t, tc)
		})
	}
}

// This is the actual test method that creates the Coherence resource, waits for it to start
// and then asserts that management over rest can be retrieved from the endpoints for the
// Deployment Pods using SSL or not depending on the configuration.
func testManagementOverRest(t *testing.T, tc ManagementTestCase) {
	g := NewWithT(t)
	ns := helper.GetTestNamespace()

	// deploy the Coherence resource
	deployment := tc.Deployment.DeepCopy()
	err := testContext.Client.Create(context.TODO(), deployment)
	g.Expect(err).NotTo(HaveOccurred())

	// defer clean-up so that we remove the deployment after this test case is finished
	defer cleanupManagement(t, deployment, ns)

	assertManagementOverRest(t, tc)
}

// assert management for a test case
func assertManagementOverRest(t *testing.T, tc ManagementTestCase) {
	g := NewWithT(t)
	ns := helper.GetTestNamespace()

	deployment := tc.Deployment
	replicas := deployment.GetReplicas()

	// Wait for the StatefulSet for the deployment to be ready - wait five minutes max
	sts, err := helper.WaitForStatefulSetForDeployment(testContext, ns, deployment, time.Second*10, time.Minute*5)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(sts.Status.ReadyReplicas).To(Equal(replicas))

	// determine whether the deployment is using SSL
	isSSL := deployment.Spec.Coherence.Management.IsSSLEnabled()

	fmt.Printf("Deployment %s SSL Enabled = %t\n", deployment.Name, isSSL)

	// Get the deployment Pods
	pods, err := helper.ListCoherencePodsForDeployment(testContext, ns, deployment.Name)
	g.Expect(err).NotTo(HaveOccurred())

	// For each Pod test whether we can connect to management over rest
	for _, pod := range pods {
		if isSSL {
			err = requestManagementWithSSL(pod, tc, tc.ShouldSucceed)
			if tc.ShouldSucceed {
				g.Expect(err).NotTo(HaveOccurred())
			} else {
				g.Expect(err).To(HaveOccurred())
			}
		} else {
			err = requestManagementWithoutSSL(pod, true)
			g.Expect(err).NotTo(HaveOccurred())
		}
	}
}

// test connecting to an SSL enabled Pod.
func requestManagementWithSSL(pod corev1.Pod, tc ManagementTestCase, retry bool) error {
	certDir, err := helper.FindTestCertsDir()
	if err != nil {
		return err
	}

	keyFile := certDir + string(os.PathSeparator) + tc.KeyFile
	certFile := certDir + string(os.PathSeparator) + tc.CertFile
	caCertFile := certDir + string(os.PathSeparator) + tc.CaCertFile

	caCert, err := os.ReadFile(caCertFile)
	if err != nil {
		return err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return err
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:      caCertPool,
				Certificates: []tls.Certificate{cert},
			},
		},
	}

	return assertManagementRequest(pod, client, "https", retry)
}

// test connecting to a plain http Pod.
func requestManagementWithoutSSL(pod corev1.Pod, retry bool) error {
	client := &http.Client{}
	return assertManagementRequest(pod, client, "http", retry)
}

// make a management over rest request.
func assertManagementRequest(pod corev1.Pod, client *http.Client, protocol string, retry bool) error {
	pf, ports, err := helper.StartPortForwarderForPod(&pod)
	if err != nil {
		return err
	}

	defer pf.Close()

	url := fmt.Sprintf("%s://127.0.0.1:%d/management/coherence/cluster", protocol, ports[coh.PortNameManagement])

	var resp *http.Response

	// try a max of 5 times
	for i := 0; i < 5; i++ {
		resp, err = client.Get(url)
		if err == nil || !retry {
			break
		}
		time.Sleep(5 * time.Second)
	}

	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request returned non-200 status %d", resp.StatusCode)
	}

	return nil
}

func cleanupManagement(t *testing.T, deployment *coh.Coherence, ns string) {
	helper.DumpState(testContext, ns, t.Name())

	err := testContext.Client.Delete(context.TODO(), deployment)
	if err != nil {
		t.Log(err)
	}
	err = helper.WaitForCoherenceCleanup(testContext, ns)
	if err != nil {
		t.Log(err)
	}
}
