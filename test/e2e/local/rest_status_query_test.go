/*
 * Copyright (c) 2020, 2025, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */
package local

import (
	"fmt"
	. "github.com/onsi/gomega"
	coh "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	"net/http"
	"testing"
	"time"
)

func TestRestStatusQueryWithInvalidPath(t *testing.T) {
	// Make sure we defer clean-up when we're done!!
	testContext.CleanupAfterTest(t)
	g := NewGomegaWithT(t)

	url := "http://127.0.0.1:8000/status/foo"

	client := &http.Client{}

	println("Connecting with: ", url)
	resp, err := client.Get(url)
	if resp != nil {
		defer resp.Body.Close()
	}
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
}

func TestRestStatusQueryForUnknownDeployment(t *testing.T) {
	// Make sure we defer clean-up when we're done!!
	testContext.CleanupAfterTest(t)
	g := NewGomegaWithT(t)

	url := "http://127.0.0.1:8000/status/foo/bar"

	client := &http.Client{}

	println("Connecting with: ", url)
	resp, err := client.Get(url)
	if resp != nil {
		defer resp.Body.Close()
	}
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
}

func TestRestStatusQueryForDeployment(t *testing.T) {
	// Make sure we defer clean-up when we're done!!
	testContext.CleanupAfterTest(t)
	namespace := helper.GetTestNamespace()
	name := "minimal-cluster"
	g := NewGomegaWithT(t)

	// deploy the cluster
	helper.AssertDeployments(testContext, t, "deployment-minimal.yaml")

	// wait for Coherence resource to get to the ready state
	_, err := helper.WaitForCoherenceCondition(testContext, namespace, name, helper.StatusPhaseCondition(coh.ConditionTypeReady), 10*time.Second, 5*time.Minute)
	g.Expect(err).NotTo(HaveOccurred())

	url := fmt.Sprintf("http://127.0.0.1:8000/status/%s/minimal-cluster", namespace)

	client := &http.Client{}

	t.Logf("Testing status query URL is %s", url)
	resp, err := client.Get(url)
	if resp != nil {
		defer resp.Body.Close()
	}
	t.Logf("Received status response: %v err %v", resp, err)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(resp.StatusCode).To(Equal(http.StatusOK))
}
