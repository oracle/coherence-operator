/*
 * Copyright (c) 2019, 2023, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package nettesting

import (
	"context"
	"github.com/oracle/coherence-operator/pkg/clients"
	ctrl "sigs.k8s.io/controller-runtime"
)

type operatorSimulator struct {
	host string
}

// _ is a simple variable to verify at compile time that operatorSimulator implements ClientSimulator
var _ ClientSimulator = operatorSimulator{}

// Run executes the Operator simulator test
func (o operatorSimulator) Run(ctx context.Context) error {
	log.Info("Starting test", "Name", "Operator Simulator")

	ports, err := getPorts()
	if err != nil {
		return err
	}

	tests := make(map[string]portTester)
	tests["K8s"] = k8sAPITester{}
	tests[TestPortHealth] = simplePortTester{name: TestPortHealth, host: o.host, port: ports[TestPortHealth], protocol: "tcp"}

	runTests(ctx, tests)
	return nil
}

type k8sAPITester struct {
}

func (in k8sAPITester) testPort(_ context.Context) {
	log.Info("Testing connectivity", "PortName", "K8s API Server")

	cfg, err := ctrl.GetConfig()
	if err != nil {
		log.Info("Testing connectivity FAILED", "PortName", "K8s API Server", "Error", err.Error())
		return
	}

	cs, err := clients.NewForConfig(cfg)
	if err != nil {
		log.Info("Testing connectivity FAILED", "PortName", "K8s API Server", "Error", err.Error())
		return
	}

	version, err := cs.DiscoveryClient.ServerVersion()
	if err != nil {
		log.Info("Testing connectivity FAILED", "PortName", "K8s API Server", "Error", err.Error())
		return
	}

	log.Info("Testing connectivity PASSED", "PortName", "K8s API Server", "Version", version)
}
