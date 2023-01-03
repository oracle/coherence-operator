package net_testing

import (
	"context"
	"github.com/oracle/coherence-operator/pkg/clients"
	ctrl "sigs.k8s.io/controller-runtime"
)

type operatorSimulator struct {
	host string
}

// _ is a simple variable to verify at compile time that operatorSimulator implements OperatorSimulator
var _ OperatorSimulator = operatorSimulator{}

// Run executes the Operator simulator test
func (o operatorSimulator) Run(ctx context.Context) error {
	log.Info("Starting test", "Name", "Operator Simulator")

	ports, err := getPorts()
	if err != nil {
		return err
	}

	tests := make(map[string]portTester)
	tests["K8s"] = k8sApiTester{}
	tests[TestPortHealth] = simplePortTester{name: TestPortHealth, host: o.host, port: ports[TestPortHealth], protocol: "tcp"}

	runTests(ctx, tests)
	return nil
}

type k8sApiTester struct {
}

func (in k8sApiTester) testPort(ctx context.Context) {
	log.Info("Testing connectivity", "PortName", "K8s API Server")

	cfg, err := ctrl.GetConfig()
	if err != nil {
		log.Info("Testing connectivity failed", "PortName", "K8s API Server", "Error", err.Error())
		return
	}

	cs, err := clients.NewForConfig(cfg)
	if err != nil {
		log.Info("Testing connectivity failed", "PortName", "K8s API Server", "Error", err.Error())
		return
	}

	version, err := cs.DiscoveryClient.ServerVersion()
	if err != nil {
		log.Info("Testing connectivity failed", "PortName", "K8s API Server", "Error", err.Error())
		return
	}

	log.Info("Testing connectivity passed", "PortName", "K8s API Server", "Version", version)
}
