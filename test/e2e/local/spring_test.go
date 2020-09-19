/*
 *  Copyright (c) 2020, Oracle and/or its affiliates.
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
)

func TestStartSpringCluster(t *testing.T) {
	testContext.CleanupAfterTest(t)
	_, pods := helper.AssertDeployments(testContext, t, "spring-cluster.yaml")
	AssertEndpoint(t, pods)
}

func TestStartSpringFatJarCluster(t *testing.T) {
	testContext.CleanupAfterTest(t)
	_, pods := helper.AssertDeployments(testContext, t, "spring-fat-jar-cluster.yaml")
	AssertEndpoint(t, pods)
}

func TestStartSpringBuildpacksCluster(t *testing.T) {
	testContext.CleanupAfterTest(t)
	_, pods := helper.AssertDeployments(testContext, t, "spring-buildpack-cluster.yaml")
	AssertEndpoint(t, pods)
}

// Assert that we can hit the Spring Boot web-app endpoint
func AssertEndpoint(t *testing.T, pods []corev1.Pod) {
	g := NewGomegaWithT(t)

	pf, ports, err := helper.StartPortForwarderForPod(&pods[0])
	g.Expect(err).NotTo(HaveOccurred())
	defer pf.Close()

	client := &http.Client{}
	url := fmt.Sprintf("http://127.0.0.1:%d/", ports["web"])

	resp, err := client.Get(url)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(resp.StatusCode).To(Equal(http.StatusOK))
}
