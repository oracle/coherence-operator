/*
 * Copyright (c) 2021, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package statefulset

import (
	"net"
	"strconv"
	"time"
)

// NewTCPProbe creates TCPProbe.
func NewTCPProbe() TCPProbe {
	return tcpProbe{}
}

// TCPProbe is an interface that defines the Probe function for doing TCP readiness/liveness checks.
type TCPProbe interface {
	Probe(host string, port int, timeout time.Duration) (Result, string, error)
}

type tcpProbe struct{}

// Probe returns a ProbeRunner capable of running an TCP check.
func (pr tcpProbe) Probe(host string, port int, timeout time.Duration) (Result, string, error) {
	return DoTCPProbe(net.JoinHostPort(host, strconv.Itoa(port)), timeout)
}

// DoTCPProbe checks that a TCP socket to the address can be opened.
// If the socket can be opened, it returns Success
// If the socket fails to open, it returns Failure.
// This is exported because some other packages may want to do direct TCP probes.
func DoTCPProbe(addr string, timeout time.Duration) (Result, string, error) {
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		// Convert errors to failures to handle timeouts.
		return Failure, err.Error(), nil
	}
	err = conn.Close()
	if err != nil {
		log.Error(err, "Unexpected error closing TCP probe socket")
	}
	return Success, "", nil
}
