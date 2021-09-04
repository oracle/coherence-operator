/*
 * Copyright (c) 2020, 2021, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package clients

import (
	"fmt"
	"github.com/oracle/coherence-operator/pkg/operator"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	"github.com/spf13/viper"
	"os"
	"testing"
)

var testContext helper.TestContext

// The entry point for the test suite
func TestMain(m *testing.M) {
	var err error

	helper.EnsureTestEnvVars()

	if testContext, err = helper.NewContext(true); err != nil {
		fmt.Printf("Error: %+v", err)
		os.Exit(1)
	}

	// We never do service suspension in local tests as the external Operator cannot see the Pods directly
	viper.Set(operator.FlagSkipServiceSuspend, true)

	exitCode := m.Run()
	testContext.Logf("Tests completed with return code %d", exitCode)
	testContext.Close()
	os.Exit(exitCode)
}
