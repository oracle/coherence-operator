/*
 * Copyright (c) 2019, 2023, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package nettesting

import (
	"context"
	"fmt"
	coh "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/pkg/operator"
	"net"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"strconv"
	"time"
)

const (
	// TestPortEcho is the name of the echo port
	TestPortEcho = "Echo"
	// TestPortHealth is the name of the health port
	TestPortHealth = "Health"
	// TestPortWebHook is the name of the web-hook port
	TestPortWebHook = "WebHook"
	// TestPortClusterPort is the name of the Coherence cluster port
	TestPortClusterPort = "ClusterPort"
	// TestPortUnicastPort1 is the name of the first Coherence unicast port
	TestPortUnicastPort1 = "UnicastPort1"
	// TestPortUnicastPort2 is the name of the second Coherence unicast port
	TestPortUnicastPort2 = "UnicastPort2"
	// TestPortMetrics is the name of the Coherence metrics port
	TestPortMetrics = "Metrics"
	// TestPortManagement is the name of the Coherence management port
	TestPortManagement = "Management"
	// TestPortOperatorRest is the name of the Operator REST port
	TestPortOperatorRest = "OperatorRest"
)

var (
	log = ctrl.Log.WithName("net-test")
)

// ServerRunner manages a server that listens on various ports
// to allow network connectivity tests to those ports.
type ServerRunner interface {
	Run(ctx context.Context) error
}

// ClientSimulator runs a test simulating connectivity to a port.
type ClientSimulator interface {
	Run(ctx context.Context) error
}

// NewServerRunner create a new ServerRunner
func NewServerRunner() ServerRunner {
	servers := make(map[string]WebServer)
	running := make(chan struct{})
	return serverRunner{servers: servers, running: running}
}

// NewOperatorSimulatorRunner create a new ServerRunner
func NewOperatorSimulatorRunner(host string) ClientSimulator {
	return operatorSimulator{host: host}
}

// NewClusterMemberRunner create a new ClusterMemberSimulator
func NewClusterMemberRunner(operatorHost, clusterHost string) ClientSimulator {
	return clusterMemberSimulator{operatorHost: operatorHost, clusterHost: clusterHost}
}

func NewWebHookClientRunner(operatorHost string) ClientSimulator {
	return webHookClientSimulator{operatorHost: operatorHost}
}

func NewSimpleClientRunner(host string, port int, protocol string) ClientSimulator {
	if protocol == "" {
		protocol = "tcp"
	}
	return simpleClientSimulator{host: host, port: port, protocol: protocol}
}

type portTester interface {
	testPort(ctx context.Context)
}

type simplePortTester struct {
	name     string
	host     string
	port     int
	protocol string
}

func (in simplePortTester) testPort(_ context.Context) {
	var err error

	log.Info("Testing connectivity", "Host", in.host, "PortName", in.name, "Port", in.port)

	con, err := net.DialTimeout(in.protocol, fmt.Sprintf("%s:%d", in.host, in.port), time.Second*10)
	if err != nil {
		log.Info("Testing connectivity FAILED", "Host", in.host, "PortName", in.name, "Port", in.port, "Error", err.Error())
	} else {
		log.Info("Testing connectivity PASSED", "Host", in.host, "PortName", in.name, "Port", in.port)
		_ = con.Close()
	}
}

func runTests(ctx context.Context, tests map[string]portTester) {
	for _, test := range tests {
		test.testPort(ctx)
	}
}

func getPorts() (map[string]int, error) {
	var err error

	ports := make(map[string]int)
	ports[TestPortEcho] = 7

	if ports[TestPortHealth], err = findPort(coh.EnvVarCohHealthPort, int(coh.DefaultHealthPort)); err != nil {
		return nil, err
	}

	if ports[TestPortWebHook], err = findPort("WEBHOOK_PORT", 443); err != nil {
		return nil, err
	}

	if ports[TestPortClusterPort], err = findPort("COHERENCE_CLUSTERPORT", 7574); err != nil {
		return nil, err
	}

	if ports[TestPortUnicastPort1], err = findPort(coh.EnvVarCoherenceLocalPort, int(coh.DefaultUnicastPort)); err != nil {
		return nil, err
	}

	if ports[TestPortUnicastPort2], err = findPort(coh.EnvVarCoherenceLocalPortAdjust, int(coh.DefaultUnicastPortAdjust)); err != nil {
		return nil, err
	}

	if ports[TestPortMetrics], err = findPort(coh.EnvVarCohMetricsPrefix+coh.EnvVarCohPortSuffix, int(coh.DefaultMetricsPort)); err != nil {
		return nil, err
	}

	if ports[TestPortManagement], err = findPort(coh.EnvVarCohMgmtPrefix+coh.EnvVarCohPortSuffix, int(coh.DefaultManagementPort)); err != nil {
		return nil, err
	}

	if ports[TestPortOperatorRest], err = findPort("OPERATOR_REST_PORT", int(operator.GetRestPort())); err != nil {
		return nil, err
	}

	return ports, nil
}

func findPort(portEnvVar string, defaultPort int) (int, error) {
	var err error
	port := defaultPort
	if portEnvVar != "" {
		if p, found := os.LookupEnv(portEnvVar); found {
			if port, err = strconv.Atoi(p); err != nil {
				return 0, err
			}
		}
	}
	return port, nil
}

type simpleClientSimulator struct {
	host     string
	port     int
	protocol string
}

// _ is a simple variable to verify at compile time that simpleClientSimulator implements ClientSimulator
var _ ClientSimulator = simpleClientSimulator{}

// Run executes the simple client test
func (in simpleClientSimulator) Run(ctx context.Context) error {
	log.Info("Starting test", "Name", "Simple Client")

	name := fmt.Sprintf("%s-%d", in.host, in.port)
	tests := make(map[string]portTester)
	tests[name] = simplePortTester{name: name, host: in.host, port: in.port, protocol: in.protocol}

	runTests(ctx, tests)
	return nil
}
