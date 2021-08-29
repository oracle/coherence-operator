/*
 * Copyright (c) 2021, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package local

import (
	. "github.com/onsi/gomega"
	coh "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/pkg/runner"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	"testing"
)

func TestStatusForExistingCuster(t *testing.T) {
	g := NewGomegaWithT(t)
	testContext.CleanupAfterTest(t)

	helper.AssertDeployments(testContext, t, "deployment-minimal.yaml")

	ns := helper.GetTestNamespace()

	var env map[string]string

	// Should be Ready
	args := []string{"status", "--operator-url", "http://localhost:8000", "--namespace", ns, "--name", "minimal-cluster", "--timeout", "1m"}
	_, err := runner.ExecuteWithArgs(env, args)
	g.Expect(err).NotTo(HaveOccurred())

	// Should not be Scaling
	args = []string{"status", "--operator-url", "http://localhost:8000", "--namespace", ns, "--name", "minimal-cluster", "--condition", string(coh.ConditionTypeScaling), "--timeout", "1m"}
	_, err = runner.ExecuteWithArgs(env, args)
	g.Expect(err).To(HaveOccurred())
}

func TestStatusForNonExistentCluster(t *testing.T) {
	g := NewGomegaWithT(t)
	testContext.CleanupAfterTest(t)

	ns := helper.GetTestNamespace()

	var env map[string]string
	args := []string{"status", "--operator-url", "http://localhost:8000", "--namespace", ns, "--name", "foo", "--timeout", "30s"}

	// should not be found
	_, err := runner.ExecuteWithArgs(env, args)
	g.Expect(err).To(HaveOccurred())
}
