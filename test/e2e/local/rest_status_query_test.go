/*
 * Copyright (c) 2021, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */
package local

import (
	"fmt"
	. "github.com/onsi/gomega"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	"net/http"
	"testing"
)

func TestRestStatusQueryWithInvalidPath(t *testing.T) {
	// Make sure we defer clean-up when we're done!!
	testContext.CleanupAfterTest(t)
	g := NewGomegaWithT(t)

	url := "http://127.0.0.1:8000/status/foo"

	client := &http.Client{}

	println("Connecting with: ", url)
	resp, err := client.Get(url)
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
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
}

func TestRestStatusQueryForDeployment(t *testing.T) {
	// Make sure we defer clean-up when we're done!!
	testContext.CleanupAfterTest(t)
	g := NewGomegaWithT(t)

	helper.AssertDeployments(testContext, t, "deployment-minimal.yaml")

	namespace := helper.GetTestNamespace()
	url := fmt.Sprintf("http://127.0.0.1:8000/status/%s/minimal-cluster", namespace)

	client := &http.Client{}

	t.Logf("Testing status query URL is %s", url)
	resp, err := client.Get(url)
	t.Logf("Received status response: %v err %v", resp, err)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(resp.StatusCode).To(Equal(http.StatusOK))
}
