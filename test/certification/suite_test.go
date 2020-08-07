/*
 * Copyright (c) 2020 Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package certification

import (
	"fmt"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	"os"
	"testing"
)

var testContext helper.TestContext

// The entry point for the test suite
func TestMain(m *testing.M) {
	var err error

	// Create a new TestContext - DO NOT start any controllers.
	if testContext, err = helper.NewContext(false); err != nil {
		fmt.Printf("Error: %+v", err)
		os.Exit(1)
	}

	exitCode := m.Run()
	testContext.Logf("Tests completed with return code %d", exitCode)
	testContext.Close()
	os.Exit(exitCode)
}
