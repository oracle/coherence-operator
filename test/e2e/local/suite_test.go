/*
 * Copyright (c) 2020 Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package local

import (
	"github.com/oracle/coherence-operator/test/e2e/helper"
	"os"
	"testing"
)

var testContext helper.TestContext

// The entry point for the test suite
func TestMain(m *testing.M) {
	var err error
	
	if testContext, err = helper.NewContext(true); err != nil {
		panic(err)
	}
	exitCode := m.Run()
	testContext.Logf("Tests completed with return code %d", exitCode)
	testContext.Close()
	os.Exit(exitCode)
}
