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
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"testing"
)

// Test that a deployment works using the minimal valid yaml for a Coherence
func TestMinimalJob(t *testing.T) {
	// Make sure we defer clean-up when we're done!!
	testContext.CleanupAfterTest(t)
	g := NewWithT(t)

	deployments, _ := helper.AssertDeployments(testContext, t, "deployment-job-minimal.yaml")

	data, ok := deployments["minimal-job"]
	g.Expect(ok).To(BeTrue(), "did not find expected 'minimal-job' deployment")

	hasFinalizer := controllerutil.ContainsFinalizer(&data, coh.CoherenceFinalizer)
	g.Expect(hasFinalizer).To(BeTrue())
}
