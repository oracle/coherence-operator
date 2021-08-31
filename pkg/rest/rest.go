/*
 * Copyright (c) 2019, 2020, 2021 Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

// Package rest provides a REST server for the Coherence Operator.
package rest

import (
	"context"
	"fmt"
	v1 "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/pkg/clients"
	onet "github.com/oracle/coherence-operator/pkg/net"
	"github.com/oracle/coherence-operator/pkg/operator"
	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	k8s "k8s.io/client-go/kubernetes"
	"net"
	"net/http"
	ctrl "sigs.k8s.io/controller-runtime"
	"strconv"
	"strings"
	"time"
)

// The logger to use to log messages
var (
	log = ctrl.Log.WithName("rest-server")
	svr *server
)

type handler struct {
	fn func(w http.ResponseWriter, r *http.Request)
}

// ServeHTTP handles the http request
func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.fn(w, r)
}

// Server is the Operator REST server.
type Server interface {
	// GetAddress returns the address that this server is listening on.
	GetAddress() net.Addr
	// GetPort returns the port that this server is listening on.
	GetPort() int32
	// Close closes this server's listener
	Close() error
	// GetHostAndPort returns the address that the ReST server should be reached on by external processes
	GetHostAndPort() string
	// Start the REST server
	Start(ctx context.Context) error
	// SetupWithManager will configure the server to run when the manager starts
	SetupWithManager(mgr ctrl.Manager) error
	// Running closes when the server is running
	Running() <-chan struct{}
}

// GetServerHostAndPort obtains the host and port that the REST server is listening on of empty string if the server
// is not started.
func GetServerHostAndPort() string {
	if svr == nil {
		return ""
	}
	return svr.GetHostAndPort()
}

// NewServer will create a new REST server
func NewServer(c clients.ClientSet) Server {
	running := make(chan struct{})
	if svr == nil {
		svr = &server{
			client:  c.KubeClient,
			running: running,
		}
	}
	return svr
}

type server struct {
	listener   net.Listener
	client     k8s.Interface
	mgr        ctrl.Manager
	ctx        context.Context
	running    chan struct{}
	httpServer *http.Server
}

// SetupWithManager configures this server from the specified Manager.
func (s server) SetupWithManager(mgr ctrl.Manager) error {
	s.mgr = mgr
	return mgr.Add(s)
}

func (s server) Running() <-chan struct{} {
	return s.running
}

// Start starts this server
func (s server) Start(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.Handle("/site/", handler{fn: s.getSiteLabelForNode})
	mux.Handle("/rack/", handler{fn: s.getRackLabelForNode})
	mux.Handle("/status/", handler{fn: s.getCoherenceStatus})

	address := fmt.Sprintf("%s:%d", operator.GetRestHost(), operator.GetRestPort())
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return errors.Wrap(err, "failed to start REST server")
	}

	s.listener = listener

	close(s.running)

	go func() {
		log.Info("Serving REST requests", "listenAddress", s.listener.Addr().String())
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
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	s.httpServer.SetKeepAlivesEnabled(false)
	return s.httpServer.Shutdown(ctx)
}

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
	labelUsed := "<None>"

	path := r.URL.Path
	// strip off any trailing slash
	if last := len(path) - 1; last >= 0 && path[last] == '/' {
		path = path[:last]
	}

	pos := strings.LastIndex(path, "/")
	name := r.URL.Path[1+pos:]

	node, err := s.client.CoreV1().Nodes().Get(context.TODO(), name, metav1.GetOptions{})

	if err == nil {
		var ok bool
		for _, label := range labels {
			if value, ok = node.Labels[label]; ok && value != "" {
				labelUsed = label
				break
			}
		}
	} else {
		if apierrors.IsNotFound(err) {
			log.Info("GET query for node labels - NotFound", "node", name, "label", labelUsed, "value", value, "remoteAddress", r.RemoteAddr)
		} else {
			log.Error(err, "GET query for node labels - Error", "node", name, "label", labelUsed, "value", value, "remoteAddress", r.RemoteAddr)
		}
		value = ""
		labelUsed = ""
	}

	w.WriteHeader(http.StatusOK)
	if _, err = fmt.Fprint(w, value); err != nil {
		log.Error(err, "Error writing value response for node "+name)
	} else {
		log.Info("GET query for node labels", "node", name, "label", labelUsed, "value", value, "remoteAddress", r.RemoteAddr)
	}
}

// getCoherenceStatus is a GET request that returns the status of a Coherence deployment.
// The namespace and name of the deployment are extracted from the request path.
// For example, a path of /status/foo/bar would check the status of Coherence resource "bar"
// in namespace "foo".
// By default, the request checks that the deployment has a status of Ready.
// It is possible to pass in a different status using the ?phase=<expected-phase> query parameter.
func (s server) getCoherenceStatus(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	// strip off any trailing slash
	if last := len(path) - 1; last >= 0 && path[last] == '/' {
		path = path[:last]
	}
	// strip off any leading slash
	if path[0] == '/' {
		path = path[1:]
	}

	segments := strings.Split(path, "/")
	if len(segments) != 3 {
		log.Info("GET query for Coherence deployment - invalid path", "remoteAddress", r.RemoteAddr, "path", r.URL.Path)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, "Invalid status request. Required path is /status/<namespace>/<name>")
		return
	}

	coh := v1.Coherence{}
	err := s.mgr.GetClient().Get(s.ctx, types.NamespacedName{
		Namespace: segments[1],
		Name:      segments[2],
	}, &coh)

	phase := r.URL.Query().Get("phase")
	if phase == "" {
		phase = string(v1.ConditionTypeReady)
	}

	if err != nil {
		if apierrors.IsNotFound(err) {
			w.WriteHeader(http.StatusNotFound)
			log.Info("GET status query for Coherence deployment - NotFound", "namespace", segments[1], "name", segments[2], "remoteAddress", r.RemoteAddr)
			_, _ = fmt.Fprintf(w, `{"Namespace": "%s", "Name": "%s", "Required": "%s", "Actual": "NotFound"}`, segments[1], segments[2], phase)
		} else {
			log.Error(err, "GET status query for Coherence deployment - Error", "namespace", segments[1], "name", segments[2], "remoteAddress", r.RemoteAddr)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintf(w, `{"Namespace": "%s", "Name": "%s", "Required": "%s", "Actual": "Error", "Cause": "%s"}`, segments[1], segments[2], phase, err.Error())
		}
		return
	}

	actual := string(coh.Status.Phase)
	match := strings.EqualFold(phase, actual)

	var status int
	if match {
		status = http.StatusOK
	} else {
		status = http.StatusBadRequest
	}

	log.Info("GET query for Coherence deployment status", "code", strconv.Itoa(status), "required", phase, "actual", actual, "namespace", segments[1], "name", segments[2], "remoteAddress", r.RemoteAddr)
	w.WriteHeader(status)
	_, _ = fmt.Fprintf(w, `{"Namespace": "%s", "Name": "%s", "Required": "%s", "Actual": "%s"}`, segments[1], segments[2], phase, actual)
}
