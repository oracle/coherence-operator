/*
 * Copyright (c) 2019, 2023, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package local

import (
	. "github.com/onsi/gomega"
	coh "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	"os/exec"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"testing"
)

// Test that the Coherence CLI can be executed in a Pod
func TestCoherenceCLI(t *testing.T) {
	// Make sure we defer clean-up when we're done!!
	testContext.CleanupAfterTest(t)
	g := NewWithT(t)

	deployments, _ := helper.AssertDeployments(testContext, t, "deployment-cli.yaml")

	data, ok := deployments["storage"]
	g.Expect(ok).To(BeTrue(), "did not find expected 'storage' deployment")

	hasFinalizer := controllerutil.ContainsFinalizer(&data, coh.CoherenceFinalizer)
	g.Expect(hasFinalizer).To(BeTrue())

	_, err := exec.Command("kubectl", "-n", data.Namespace, "exec", "storage-0",
		"-c", "coherence", "--", "/coherence-operator/utils/cohctl", "get", "members").CombinedOutput()
	g.Expect(err).NotTo(HaveOccurred())
}
