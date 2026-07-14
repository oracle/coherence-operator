/*
 * Copyright (c) 2026, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package rest

import (
	"context"
	"net"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/oracle/coherence-operator/pkg/operator"
	"github.com/spf13/viper"
)

func TestDefaultWildcardAcceptsLoopbackConnections(t *testing.T) {
	setRestConfiguration(t, operator.DefaultRestHost, 0, "", -1)
	s := startTestServer(t)
	port := strconv.Itoa(int(s.GetPort()))

	assertCanConnect(t, net.JoinHostPort("127.0.0.1", port))
	if supportsIPv6Loopback() {
		assertCanConnect(t, net.JoinHostPort("::1", port))
	}
}

func TestExplicitRestHosts(t *testing.T) {
	tests := []struct {
		name     string
		host     string
		dialHost string
		network  string
	}{
		{name: "IPv4", host: "127.0.0.1", dialHost: "127.0.0.1", network: "tcp4"},
		{name: "IPv6", host: "::1", dialHost: "::1", network: "tcp6"},
		{name: "bracketed IPv6", host: "[::1]", dialHost: "::1", network: "tcp6"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.network == "tcp6" && !supportsIPv6Loopback() {
				t.Skip("IPv6 loopback is not available")
			}

			setRestConfiguration(t, tt.host, 0, "", -1)
			s := startTestServer(t)
			assertCanConnect(t, net.JoinHostPort(tt.dialHost, strconv.Itoa(int(s.GetPort()))))
		})
	}
}

func TestGetHostAndPort(t *testing.T) {
	t.Run("explicit IPv4", func(t *testing.T) {
		setRestConfiguration(t, "127.0.0.1", 8001, "", -1)
		s := &server{}
		if actual := s.GetHostAndPort(); actual != "127.0.0.1:8001" {
			t.Fatalf("GetHostAndPort() = %q, want %q", actual, "127.0.0.1:8001")
		}
	})

	t.Run("explicit IPv6", func(t *testing.T) {
		setRestConfiguration(t, "::1", 8002, "", -1)
		s := &server{}
		if actual := s.GetHostAndPort(); actual != "[::1]:8002" {
			t.Fatalf("GetHostAndPort() = %q, want %q", actual, "[::1]:8002")
		}
	})

	t.Run("bracketed explicit IPv6", func(t *testing.T) {
		setRestConfiguration(t, "[::1]", 8003, "", -1)
		s := &server{}
		if actual := s.GetHostAndPort(); actual != "[::1]:8003" {
			t.Fatalf("GetHostAndPort() = %q, want %q", actual, "[::1]:8003")
		}
	})

	t.Run("service", func(t *testing.T) {
		setRestConfiguration(t, operator.DefaultRestHost, 8000, "operator-rest", 9000)
		operator.GetViper().Set(operator.FlagOperatorNamespace, "testing")
		s := &server{}
		if actual := s.GetHostAndPort(); actual != "operator-rest.testing.svc:9000" {
			t.Fatalf("GetHostAndPort() = %q, want %q", actual, "operator-rest.testing.svc:9000")
		}
	})

	t.Run("wildcard", func(t *testing.T) {
		setRestConfiguration(t, operator.DefaultRestHost, 0, "", -1)
		s := startTestServer(t)
		actual := s.GetHostAndPort()
		host, port, err := net.SplitHostPort(actual)
		if err != nil {
			t.Fatalf("GetHostAndPort() returned invalid address %q: %v", actual, err)
		}
		if host == "" || isWildcardHost(host) {
			t.Fatalf("GetHostAndPort() returned unusable wildcard host %q", host)
		}
		if port != strconv.Itoa(int(s.GetPort())) {
			t.Fatalf("GetHostAndPort() port = %q, want %d", port, s.GetPort())
		}
	})
}

func startTestServer(t *testing.T) *server {
	t.Helper()

	s := &server{
		running: make(chan struct{}),
		endpoints: map[string]func(http.ResponseWriter, *http.Request){
			"/": func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusNoContent)
			},
		},
	}
	if err := s.Start(context.Background()); err != nil {
		t.Fatalf("failed to start REST server: %v", err)
	}
	t.Cleanup(func() {
		if err := s.Close(); err != nil {
			t.Errorf("failed to close REST server: %v", err)
		}
	})
	return s
}

func setRestConfiguration(t *testing.T, host string, port int32, serviceName string, servicePort int32) {
	t.Helper()

	previous := operator.GetViper()
	v := viper.New()
	v.Set(operator.FlagRestHost, host)
	v.Set(operator.FlagRestPort, port)
	v.Set(operator.FlagServiceName, serviceName)
	v.Set(operator.FlagServicePort, servicePort)
	operator.SetViper(v)
	t.Cleanup(func() {
		operator.SetViper(previous)
	})
}

func assertCanConnect(t *testing.T, address string) {
	t.Helper()

	conn, err := net.DialTimeout("tcp", address, 5*time.Second)
	if err != nil {
		t.Fatalf("failed to connect to %s: %v", address, err)
	}
	if err = conn.Close(); err != nil {
		t.Fatalf("failed to close connection to %s: %v", address, err)
	}
}

func supportsIPv6Loopback() bool {
	listener, err := net.Listen("tcp6", "[::1]:0")
	if err != nil {
		return false
	}
	_ = listener.Close()
	return true
}
