/*
 * Copyright (c) 2019, 2023, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package nettesting

import (
	"context"
	"fmt"
	"github.com/oracle/coherence-operator/pkg/operator"
	"github.com/pkg/errors"
	"net"
	"net/http"
	ctrl "sigs.k8s.io/controller-runtime"
	"time"
)

type serverRunner struct {
	servers map[string]WebServer
	running chan struct{}
}

// Run runs the network test server
func (in serverRunner) Run(ctx context.Context) error {
	var err error

	if in.servers == nil {
		in.servers = make(map[string]WebServer)
	}

	in.running = make(chan struct{})

	ports, err := getPorts()
	if err != nil {
		return err
	}

	for name, port := range ports {
		in.servers[name] = NewWebServer(name, port)
	}

	for name, server := range in.servers {
		log.Info("Starting HTTP server", "Name", name)
		if err = server.Start(ctx); err != nil {
			return err
		}
	}

	close(in.running)

	handler := ctrl.SetupSignalHandler()

	<-handler.Done()
	return nil
}

// WebServer is a simple test web server.
type WebServer interface {
	// GetAddress returns the address that this server is listening on.
	GetAddress() net.Addr
	// GetPort returns the port that this server is listening on.
	GetPort() int32
	// Close closes this server's listener
	Close() error
	// Start the server
	Start(ctx context.Context) error
	// Running closes when the server is running
	Running() <-chan struct{}
}

// NewWebServer creates a new WebServer on a specific port
func NewWebServer(name string, port int) WebServer {
	running := make(chan struct{})
	return simpleServer{name: name, port: port, running: running}
}

// simpleServer is an implementation of WebServer
type simpleServer struct {
	name       string
	port       int
	listener   net.Listener
	running    chan struct{}
	httpServer *http.Server
}

// _ is a simple variable to verify at compile time that simpleServer implements WebServer
var _ WebServer = simpleServer{}

func (in simpleServer) Start(_ context.Context) error {
	mux := http.NewServeMux()
	mux.Handle("/", handler{fn: in.requestHandler})

	address := fmt.Sprintf("%s:%d", operator.GetRestHost(), in.port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return errors.Wrap(err, "failed to start server")
	}

	in.listener = listener

	close(in.running)

	go func() {
		log.Info("Serving HTTP requests", "Name", in.name, "listenAddress", in.listener.Addr().String())
		panic(http.Serve(in.listener, mux))
	}()
	return nil
}

func (in simpleServer) Running() <-chan struct{} {
	return in.running
}

func (in simpleServer) GetName() string {
	return in.name
}

func (in simpleServer) GetAddress() net.Addr {
	return in.listener.Addr()
}

func (in simpleServer) GetPort() int32 {
	t, _ := net.ResolveTCPAddr(in.listener.Addr().Network(), in.listener.Addr().String())
	return int32(t.Port)
}

func (in simpleServer) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	in.httpServer.SetKeepAlivesEnabled(false)
	return in.httpServer.Shutdown(ctx)
}

func (in simpleServer) requestHandler(_ http.ResponseWriter, _ *http.Request) {
	log.Info("Received HTTP request", "Name", in.name)
}

// handler is a simple request handler
type handler struct {
	fn func(w http.ResponseWriter, r *http.Request)
}

// ServeHTTP handles the http request
func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.fn(w, r)
}
