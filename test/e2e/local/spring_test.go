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
	"os"
	"testing"
)

func TestStartSpringCluster(t *testing.T) {
	testContext.CleanupAfterTest(t)
	_, pods := helper.AssertDeployments(testContext, t, "spring-cluster.yaml")
	AssertSpringEndpoint(t, pods)
}

func TestStartSpringFatJarCluster(t *testing.T) {
	testContext.CleanupAfterTest(t)
	_, pods := helper.AssertDeployments(testContext, t, "spring-fat-jar-cluster.yaml")
	AssertSpringEndpoint(t, pods)
}

func TestStartSpringBuildpacksCluster(t *testing.T) {
	skip := os.Getenv("SKIP_SPRING_CNBP")
	if skip == "true" {
		return
	}
	testContext.CleanupAfterTest(t)
	_, pods := helper.AssertDeployments(testContext, t, "spring-buildpack-cluster.yaml")
	AssertSpringEndpoint(t, pods)
}

func TestStartSpringTwoCluster(t *testing.T) {
	testContext.CleanupAfterTest(t)
	_, pods := helper.AssertDeployments(testContext, t, "spring-cluster-2.yaml")
	AssertSpringEndpoint(t, pods)
}

func TestStartSpringTwoFatJarCluster(t *testing.T) {
	testContext.CleanupAfterTest(t)
	_, pods := helper.AssertDeployments(testContext, t, "spring-fat-jar-cluster-2.yaml")
	AssertSpringEndpoint(t, pods)
}

func TestStartSpringTwoBuildpacksCluster(t *testing.T) {
	skip := os.Getenv("SKIP_SPRING_CNBP")
	if skip == "true" {
		return
	}
	testContext.CleanupAfterTest(t)
	_, pods := helper.AssertDeployments(testContext, t, "spring-buildpack-cluster-2.yaml")
	AssertSpringEndpoint(t, pods)
}

// Assert that we can hit the Spring Boot web-app endpoint
func AssertSpringEndpoint(t *testing.T, pods []corev1.Pod) {
	g := NewGomegaWithT(t)

	pf, ports, err := helper.StartPortForwarderForPodWithBackoff(&pods[0])
	g.Expect(err).NotTo(HaveOccurred())
	defer pf.Close()

	client := &http.Client{}
	url := fmt.Sprintf("http://%s:%d/", pf.Hostname, ports["web"])

	resp, err := client.Get(url)
	if resp != nil {
		defer resp.Body.Close()
	}
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(resp.StatusCode).To(Equal(http.StatusOK))
}
