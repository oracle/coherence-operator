/*
 * Copyright (c) 2020, 2025, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package local

import (
	"fmt"
	. "github.com/onsi/gomega"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	corev1 "k8s.io/api/core/v1"
	"net/http"
	"testing"
	"time"
)

func TestHelidonCdiCluster(t *testing.T) {
	testContext.CleanupAfterTest(t)
	_, pods := helper.AssertDeployments(testContext, t, "helidon-cluster.yaml")
	AssertHelidonEndpoint(t, pods)
}

func TestHelidonThreeCdiCluster(t *testing.T) {
	testContext.CleanupAfterTest(t)
	_, pods := helper.AssertDeployments(testContext, t, "helidon-cluster-3.yaml")
	AssertHelidonEndpoint(t, pods)
}

func TestHelidonTwoCdiCluster(t *testing.T) {
	testContext.CleanupAfterTest(t)
	_, pods := helper.AssertDeployments(testContext, t, "helidon-cluster-2.yaml")
	AssertHelidonEndpoint(t, pods)
}

// Assert that we can hit the Helidon web-app endpoint
func AssertHelidonEndpoint(t *testing.T, pods []corev1.Pod) {
	g := NewGomegaWithT(t)

	pf, ports, err := helper.StartPortForwarderForPodWithBackoff(&pods[0])
	g.Expect(err).NotTo(HaveOccurred())
	defer pf.Close()
	time.Sleep(10 * time.Second)
	client := &http.Client{}
	url := fmt.Sprintf("http://127.0.0.1:%d/ready", ports["web"])

	resp, err := client.Get(url)
	if resp != nil {
		defer resp.Body.Close()
	}
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(resp.StatusCode).To(Equal(http.StatusOK))
}
