/*
 * Copyright (c) 2021, 2023, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package compatibility

import (
	"fmt"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
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

	// Ensure that the Operator has been deployed to the test namespace
	namespace := helper.GetTestNamespace()
	pods, err := helper.ListOperatorPods(testContext, namespace)
	if err != nil {
		fmt.Printf("Error looking for Operator Pods in namespace %s : %+v", namespace, err)
		os.Exit(1)
	}
	if len(pods) == 0 {
		fmt.Printf("Cannot find any Operator Pods in namespace %s. "+
			"This test suite requires an Operator is already deployed", namespace)
		os.Exit(1)
	}

	fmt.Printf("Waiting for Operator Pod %s to be ready in namespace %s.", pods[0].Name, namespace)
	err = helper.WaitForPodReady(testContext, namespace, pods[0].Name, 10*time.Second, 5*time.Minute)
	if err != nil {
		fmt.Printf("Failed waiting for Operator Pod %s to be ready in namespace %s.", pods[0].Name, namespace)
		os.Exit(1)
	}

	exitCode := m.Run()
	testContext.Logf("Tests completed with return code %d", exitCode)
	testContext.Close()
	os.Exit(exitCode)
}

func scale(t *testing.T, namespace, name string, replicas int32) error {
	cmd := exec.Command("kubectl", "-n", namespace, "scale", fmt.Sprintf("--replicas=%d", replicas), "coherence/"+name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	t.Log("Executing Scale Command: " + strings.Join(cmd.Args, " "))
	return cmd.Run()
}
