/*
 * Copyright (c) 2019, 2023, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package nettesting

import "context"

type webHookClientSimulator struct {
	operatorHost string
}

// _ is a simple variable to verify at compile time that clusterMemberSimulator implements ClientSimulator
var _ ClientSimulator = webHookClientSimulator{}

// Run executes the Operator simulator test
func (in webHookClientSimulator) Run(ctx context.Context) error {
	log.Info("Starting test", "Name", "Web-Hook Client")

	ports, err := getPorts()
	if err != nil {
		return err
	}

	tests := make(map[string]portTester)
	tests[TestPortWebHook] = simplePortTester{name: TestPortWebHook, host: in.operatorHost, port: ports[TestPortWebHook], protocol: "tcp"}

	runTests(ctx, tests)
	return nil
}
