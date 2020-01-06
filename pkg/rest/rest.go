/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

// The rest package provides a ReST server for the Coherence Operator.
package rest

import (
	"fmt"
	"github.com/oracle/coherence-operator/pkg/flags"
	onet "github.com/oracle/coherence-operator/pkg/net"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8s "k8s.io/client-go/kubernetes"
	"net"
	"net/http"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"strings"
)

// The logger to use to log messages
var log = logf.Log.WithName("rest-server")

type handler struct {
	fn func(w http.ResponseWriter, r *http.Request)
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.fn(w, r)
}

type Server interface {
	// GetAddress returns the address that this server is listening on.
	GetAddress() net.Addr
	// GetAddress returns the address that this server is listening on.
	GetPort() int32
	// Close closes this server's listener
	Close() error
	// GetHostAndPort returns the address that the ReST server should be reached on by external processes
	GetHostAndPort(*flags.CoherenceOperatorFlags) string
}

// StartRestServer starts a ReST server to server Coherence Operator requests,
// for example node zone information.
func StartRestServer(m manager.Manager, cf *flags.CoherenceOperatorFlags) (Server, error) {
	address := fmt.Sprintf("%s:%d", cf.RestHost, cf.RestPort)

	client, err := k8s.NewForConfig(m.GetConfig())
	if err != nil {
		return nil, err
	}

	s := server{cohFlags: cf, client: client}

	mux := http.NewServeMux()
	mux.Handle("/site/", handler{fn: s.getSiteLabelForNode})
	mux.Handle("/rack/", handler{fn: s.getRackLabelForNode})

	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}

	s.listener = listener

	go func() {
		log.Info("Serving ReST requests on " + s.listener.Addr().String())
		panic(http.Serve(s.listener, mux))
	}()

	return s, nil
}

type server struct {
	cohFlags *flags.CoherenceOperatorFlags
	listener net.Listener
	client   *k8s.Clientset
}

func (s server) GetAddress() net.Addr {
	return s.listener.Addr()
}

func (s server) GetPort() int32 {
	t, _ := net.ResolveTCPAddr(s.listener.Addr().Network(), s.listener.Addr().String())
	return int32(t.Port)
}

func (s server) Close() error {
	return s.listener.Close()
}

// GetHostAndPort returns the address and port that this endpoint can be reached on by external processes.
func (s server) GetHostAndPort(cof *flags.CoherenceOperatorFlags) string {
	f := flags.GetOperatorFlags()

	var service string
	var port int32

	switch {
	case f.ServiceName != "":
		// use the service name if it was specifically set
		service = f.ServiceName
	case f.RestHost != "0.0.0.0":
		// if no service name was set but ReST is bound to a specific address then use that
		service = f.RestHost
	default:
		// ReST is bound to 0.0.0.0 so use any of our local addresses.
		// This does not guarantee we're reachable but would be OK in local testing
		ip, err := onet.GetLocalAddress()
		if err == nil && ip != nil {
			service = fmt.Sprint(ip.String())
		}
	}

	switch {
	case f.ServicePort != -1:
		port = f.ServicePort
	case f.RestPort > 0:
		port = f.RestPort
	default:
		port = s.GetPort()
	}

	return fmt.Sprintf("%s:%d", service, port)
}

// getSiteLabelForNode is a GET request that returns the node label on a k8s node to use for a Coherence site value.
func (s server) getSiteLabelForNode(w http.ResponseWriter, r *http.Request) {
	s.getLabelForNode(s.cohFlags.SiteLabel, w, r)
}

// getRackLabelForNode is a GET request that returns the node label on a k8s node to use for a Coherence rack value.
func (s server) getRackLabelForNode(w http.ResponseWriter, r *http.Request) {
	s.getLabelForNode(s.cohFlags.RackLabel, w, r)
}

// getRackLabelForNode is a GET request that returns the node label on a k8s node to use for a Coherence rack value.
func (s server) getLabelForNode(label string, w http.ResponseWriter, r *http.Request) {
	var value string
	pos := strings.LastIndex(r.URL.Path, "/")
	name := r.URL.Path[1+pos:]

	log.Info(fmt.Sprintf("Querying for node name='%s' URL: %s", name, r.URL.Path))

	node, err := s.client.CoreV1().Nodes().Get(name, metav1.GetOptions{})

	if err == nil {
		value = node.Labels[label]
	} else {
		log.Error(err, "Error getting node "+name+" from k8s")
		value = ""
	}

	w.WriteHeader(200)
	if _, err = fmt.Fprint(w, value); err != nil {
		log.Error(err, "Error writing value response for node "+name)
	}
}
