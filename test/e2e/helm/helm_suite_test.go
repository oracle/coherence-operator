/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package helm_test

import (
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	"testing"
)

var HelmHelper *helper.HelmHelper

// This is the Ginkgo test suite entry point. In here we configure the
// HelmHelper that can then be used by the rest of the tests in the suite
func TestCoherenceRoleControler(t *testing.T) {
	RegisterFailHandler(Fail)

	// Create a helper.HelmHelper
	h, err := helper.NewOperatorChartHelper()
	if err != nil {
		t.Fatal(err)
	}

	HelmHelper = h

	// Make Ginkgo run the rest of the test suite
	junitReporter := reporters.NewJUnitReporter("test-report.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "Coherence Operator Helm Suite", []Reporter{junitReporter})
}
