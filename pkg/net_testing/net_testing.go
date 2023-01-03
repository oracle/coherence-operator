package net_testing

import (
	"context"
	"fmt"
	coh "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/pkg/operator"
	"net"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"strconv"
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

// OperatorSimulator runs a test simulating the connectivity
// required by the Coherence Operator.
type OperatorSimulator interface {
	Run(ctx context.Context) error
}

// ClusterMemberSimulator runs a test simulating the connectivity
// required by the Coherence cluster member.
type ClusterMemberSimulator interface {
	Run(ctx context.Context) error
}

// NewServerRunner create a new ServerRunner
func NewServerRunner() ServerRunner {
	servers := make(map[string]WebServer)
	running := make(chan struct{})
	return serverRunner{servers: servers, running: running}
}

// NewOperatorSimulatorRunner create a new ServerRunner
func NewOperatorSimulatorRunner(host string) OperatorSimulator {
	return operatorSimulator{host: host}
}

// NewClusterMemberRunner create a new ClusterMemberSimulator
func NewClusterMemberRunner(operatorHost, clusterHost string) ClusterMemberSimulator {
	return clusterMemberSimulator{operatorHost: operatorHost, clusterHost: clusterHost}
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

func (in simplePortTester) testPort(ctx context.Context) {
	var err error

	log.Info("Testing connectivity", "PortName", in.name, "Port", in.port)

	con, err := net.Dial(in.protocol, fmt.Sprintf("%s:%d", in.host, in.port))
	if err != nil {
		log.Info("Testing connectivity FAILED", "PortName", in.name, "Port", in.port, "Error", err.Error())
	} else {
		log.Info("Testing connectivity PASSED", "PortName", in.name, "Port", in.port)
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
