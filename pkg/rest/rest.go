/*
 * Copyright (c) 2019, 2020, 2021 Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

// The rest package provides a ReST server for the Coherence Operator.
package rest

import (
	"context"
	"fmt"
	"github.com/oracle/coherence-operator/pkg/clients"
	onet "github.com/oracle/coherence-operator/pkg/net"
	"github.com/oracle/coherence-operator/pkg/operator"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8s "k8s.io/client-go/kubernetes"
	"net"
	"net/http"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"strings"
)

// The logger to use to log messages
var (
	log   = logf.Log.WithName("rest-server")
	svr   *server
)

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
	GetHostAndPort() string
	// Start the REST server
	Start(stop <-chan struct{}) error
	SetupWithManager(mgr ctrl.Manager) error
}

// Obtain the host and port that the REST server is listening on of empty string if the server
// is not started.
func GetServerHostAndPort() string {
	if svr == nil {
		return ""
	}
	return svr.GetHostAndPort()
}

func NewServer(c clients.ClientSet) Server {
	if svr == nil {
		svr = &server{
			client: c.KubeClient,
		}
	}
	return svr
}

type server struct {
	listener net.Listener
	client   k8s.Interface
}

func (s server) SetupWithManager(mgr ctrl.Manager) error {
	return mgr.Add(s)
}

func (s server) Start(stop <-chan struct{}) error {
	mux := http.NewServeMux()
	mux.Handle("/site/", handler{fn: s.getSiteLabelForNode})
	mux.Handle("/rack/", handler{fn: s.getRackLabelForNode})

	address := fmt.Sprintf("%s:%d", operator.GetRestHost(), operator.GetRestPort())
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return errors.Wrap(err, "failed to start REST server")
	}

	s.listener = listener

	go func() {
		log.Info("Serving REST requests on " + s.listener.Addr().String())
		panic(http.Serve(s.listener, mux))
	}()
	return nil
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
func (s server) GetHostAndPort() string {
	var service string
	var port int32

	restHost := operator.GetRestHost()
	serviceName := operator.GetRestServiceName()

	switch {
	case serviceName != "":
		// use the service name if it was specifically set
		service = serviceName
	case restHost != "0.0.0.0":
		// if no service name was set but REST is bound to a specific address then use that
		service = restHost
	default:
		// REST is bound to 0.0.0.0 so use any of our local addresses.
		// This does not guarantee we're reachable but would be OK in local testing
		ip, err := onet.GetLocalAddress()
		if err == nil && ip != nil {
			service = fmt.Sprint(ip.String())
		}
	}

	restPort := operator.GetRestPort()
	servicePort := operator.GetRestServicePort()

	switch {
	case servicePort != -1:
		port = servicePort
	case restPort > 0:
		port = restPort
	default:
		port = s.GetPort()
	}

	return fmt.Sprintf("%s:%d", service, port)
}

// getSiteLabelForNode is a GET request that returns the node label on a k8s node to use for a Coherence site value.
func (s server) getSiteLabelForNode(w http.ResponseWriter, r *http.Request) {
	s.getLabelForNode(operator.GetSiteLabel(), w, r)
}

// getRackLabelForNode is a GET request that returns the node label on a k8s node to use for a Coherence rack value.
func (s server) getRackLabelForNode(w http.ResponseWriter, r *http.Request) {
	s.getLabelForNode(operator.GetRackLabel(), w, r)
}

// getRackLabelForNode is a GET request that returns the node label on a k8s node to use for a Coherence rack value.
func (s server) getLabelForNode(labels []string, w http.ResponseWriter, r *http.Request) {
	var value string
	pos := strings.LastIndex(r.URL.Path, "/")
	name := r.URL.Path[1+pos:]

	log.Info(fmt.Sprintf("Querying for node name='%s' URL: %s", name, r.URL.Path))

	node, err := s.client.CoreV1().Nodes().Get(context.TODO(), name, metav1.GetOptions{})

	if err == nil {
		var ok bool
		for _, label := range labels {
			if value, ok = node.Labels[label]; ok && value != "" {
				break
			}
		}
	} else {
		log.Error(err, "Error getting node "+name+" from k8s")
		value = ""
	}

	w.WriteHeader(200)
	if _, err = fmt.Fprint(w, value); err != nil {
		log.Error(err, "Error writing value response for node "+name)
	}
}
