/*
 * Copyright (c) 2019, 2023, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1_test

import (
	. "github.com/onsi/gomega"
	coh "github.com/oracle/coherence-operator/api/v1"
	"testing"
)

func TestCreateResourcesFromMinimalSpec(t *testing.T) {
	g := NewGomegaWithT(t)

	// Create minimal spec spec
	spec := coh.CoherenceStatefulSetResourceSpec{}
	// Create the test deployment
	deployment := createTestCoherenceDeployment(spec)

	// Create expected Job
	expected := createMinimalExpectedStatefulSet(deployment)
	// assert that the resources is as expected
	res := assertResourceCreation(t, deployment)
	sts := res.GetResourcesOfKind(coh.ResourceTypeStatefulSet)
	g.Expect(len(sts)).To(Equal(1))
	assertStatefulSet(t, sts[0], expected)
}
