package net_testing

import (
	"context"
)

type clusterMemberSimulator struct {
	operatorHost string
	clusterHost  string
}

// _ is a simple variable to verify at compile time that clusterMemberSimulator implements ClusterMemberSimulator
var _ ClusterMemberSimulator = clusterMemberSimulator{}

// Run executes the Operator simulator test
func (in clusterMemberSimulator) Run(ctx context.Context) error {
	log.Info("Starting test", "Name", "Cluster Member Simulator")

	ports, err := getPorts()
	if err != nil {
		return err
	}

	tests := make(map[string]portTester)
	tests[TestPortOperatorRest] = simplePortTester{name: TestPortOperatorRest, host: in.operatorHost, port: ports[TestPortOperatorRest], protocol: "tcp"}
	tests[TestPortEcho] = simplePortTester{name: TestPortEcho, host: in.clusterHost, port: ports[TestPortEcho], protocol: "tcp"}
	tests[TestPortClusterPort] = simplePortTester{name: TestPortClusterPort, host: in.clusterHost, port: ports[TestPortClusterPort], protocol: "tcp"}
	tests[TestPortUnicastPort1] = simplePortTester{name: TestPortUnicastPort1, host: in.clusterHost, port: ports[TestPortUnicastPort1], protocol: "tcp"}
	tests[TestPortUnicastPort2] = simplePortTester{name: TestPortUnicastPort2, host: in.clusterHost, port: ports[TestPortUnicastPort2], protocol: "tcp"}
	tests[TestPortManagement] = simplePortTester{name: TestPortManagement, host: in.clusterHost, port: ports[TestPortManagement], protocol: "tcp"}
	tests[TestPortMetrics] = simplePortTester{name: TestPortMetrics, host: in.clusterHost, port: ports[TestPortMetrics], protocol: "tcp"}

	runTests(ctx, tests)
	return nil
}
