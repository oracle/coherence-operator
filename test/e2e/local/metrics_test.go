/*
 * Copyright (c) 2019, 2020, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package local

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	coh "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	"io/ioutil"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/pointer"
	"net/http"
	"os"
	"testing"
	"time"

	. "github.com/onsi/gomega"
)

// A struct used to define a metrics test case.
type MetricsTestCase struct {
	Deployment    *coh.Coherence
	Name          string
	Ctx           *framework.Context
	KeyFile       string
	CertFile      string
	CaCertFile    string
	ShouldSucceed bool
}

// TestMetrics is a go test that uses sub-tests (test cases) to basically run the
// same test with different parameters. In this case different Coherence resource
// configurations with metrics configured with and without SSL.
func TestMetrics(t *testing.T) {
	helper.SkipIfCoherenceVersionLessThan(t, 12, 2, 1, 4)

	// initialise Gomega so we can use matchers
	g := NewGomegaWithT(t)

	// Create the Operator SDK test context (this will deploy the Operator)
	ctx := helper.CreateTestContext(t)
	// Make sure we defer clean-up (uninstall the operator) when we're done
	defer helper.DumpOperatorLogs(t)

	// Get the test namespace
	namespace, err := ctx.GetWatchNamespace()
	g.Expect(err).NotTo(HaveOccurred())

	// Get the test SSL information (secret name etc.)
	_, ssl, err := helper.GetTestSslSecret()
	g.Expect(err).NotTo(HaveOccurred())

	// require 2-way auth
	ssl.RequireClientCert = pointer.BoolPtr(true)

	// load the test Coherence resource from a yaml files
	deploymentWithoutSSL, err := helper.NewSingleCoherenceFromYaml(namespace, "metrics-test.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	// load the test Coherence resource that used a distroless JIB image from a yaml files
	deploymentJib, err := helper.NewSingleCoherenceFromYaml(namespace, "metrics-jib-test.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	// Copy deploymentWithoutSSL and configure it to use SSL at the Spec level in all deployments
	deploymentSSL := &coh.Coherence{}
	deploymentWithoutSSL.DeepCopyInto(deploymentSSL)

	// Set the SSL settings
	deploymentSSL.Name = "metrics-ssl"
	deploymentSSL.Spec.Coherence.Metrics.Enabled = pointer.BoolPtr(true)
	deploymentSSL.Spec.Coherence.Metrics.SSL = ssl

	// Create the test cases
	testCases := []MetricsTestCase{
		{&deploymentWithoutSSL, "PlainHTTP", ctx, "", "", "", true},
		{&deploymentJib, "JIB", ctx, "", "", "", true},
		{deploymentSSL, "WithSSL", ctx, "groot.key", "groot.cert", "guardians-ca.crt", true},
		{deploymentSSL, "ClientHasBadKey", ctx, "yondu.key", "groot.cert", "guardians-ca.crt", false},
		{deploymentSSL, "BadCert", ctx, "groot.key", "yondu.cert", "guardians-ca.crt", false},
		{deploymentSSL, "BadCaCert", ctx, "groot.key", "groot.cert", "ravagers-ca.crt", false},
	}

	// Run the test cases...
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			testClusterMetrics(t, tc)
		})
	}
}

// This is the actual test method that creates the Coherence resource, waits for it to start
// and then asserts that metrics can be retrieved from the endpoints for the Deployment Pods
// using SSL or not depending on the configuration.
func testClusterMetrics(t *testing.T, tc MetricsTestCase) {
	f := framework.Global
	g := NewGomegaWithT(t)

	ns, err := tc.Ctx.GetWatchNamespace()
	g.Expect(err).NotTo(HaveOccurred())

	// deploy the Coherence resource
	deployment := tc.Deployment.DeepCopy()
	err = f.Client.Create(context.TODO(), deployment, helper.DefaultCleanup(tc.Ctx))
	g.Expect(err).NotTo(HaveOccurred())

	// defer clean-up so that we remove the deployment after this test case is finished
	defer cleanupMetrics(t, deployment, ns)

	assertMetrics(t, tc)
}

// assert metrics for a test case
func assertMetrics(t *testing.T, tc MetricsTestCase) {
	f := framework.Global
	g := NewGomegaWithT(t)
	ns := tc.Deployment.GetNamespace()

	replicas := tc.Deployment.GetReplicas()

	// Wait for the StatefulSet for the deployment to be ready - wait five minutes max
	sts, err := helper.WaitForStatefulSetForDeployment(f.KubeClient, ns, tc.Deployment, time.Second*10, time.Minute*5, t)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(sts.Status.ReadyReplicas).To(Equal(replicas))

	// determine whether the deployment is using SSL
	isSSL := tc.Deployment.Spec.Coherence.Metrics.IsSSLEnabled()

	// Get the deployment Pods
	pods, err := helper.ListCoherencePodsForDeployment(f.KubeClient, ns, tc.Deployment.GetName())
	g.Expect(err).NotTo(HaveOccurred())

	// For each Pod test whether we can connect to metrics
	for _, pod := range pods {
		if isSSL {
			err = requestMetricsWithSSL(pod, tc, tc.ShouldSucceed)
			if tc.ShouldSucceed {
				g.Expect(err).NotTo(HaveOccurred())
			} else {
				g.Expect(err).To(HaveOccurred())
			}
		} else {
			err = requestMetricsWithoutSSL(pod, true)
			g.Expect(err).NotTo(HaveOccurred())
		}
	}
}

// test connecting to an SSL enabled Pod.
func requestMetricsWithSSL(pod corev1.Pod, tc MetricsTestCase, retry bool) error {
	certDir, err := helper.FindTestCertsDir()
	if err != nil {
		return err
	}

	keyFile := certDir + string(os.PathSeparator) + tc.KeyFile
	certFile := certDir + string(os.PathSeparator) + tc.CertFile
	caCertFile := certDir + string(os.PathSeparator) + tc.CaCertFile

	caCert, err := ioutil.ReadFile(caCertFile)
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

	return assertMetricsRequest(pod, client, "https", retry)
}

// test connecting to a plain http Pod.
func requestMetricsWithoutSSL(pod corev1.Pod, retry bool) error {
	client := &http.Client{}
	return assertMetricsRequest(pod, client, "http", retry)
}

// make a metrics request.
func assertMetricsRequest(pod corev1.Pod, client *http.Client, protocol string, retry bool) error {
	pf, ports, err := helper.StartPortForwarderForPod(&pod)
	if err != nil {
		return err
	}

	defer pf.Close()

	url := fmt.Sprintf("%s://127.0.0.1:%d/metrics", protocol, ports[coh.PortNameMetrics])

	var resp *http.Response

	// try a max of 5 times
	for i := 0; i < 5; i++ {
		resp, err = client.Get(url)
		if err == nil || !retry {
			break
		}
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request returned non-200 status %d", resp.StatusCode)
	}

	return nil
}

func cleanupMetrics(t *testing.T, deployment *coh.Coherence, ns string) {
	helper.DumpState(ns, t.Name(), t)

	f := framework.Global
	err := f.Client.Delete(context.TODO(), deployment)
	if err != nil {
		t.Log(err)
	}
	err = helper.WaitForCoherenceCleanup(f, ns)
	if err != nil {
		t.Log(err)
	}
}
