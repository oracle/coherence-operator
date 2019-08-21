// The rest package provides a ReST server for the Coherence Operator.
package rest

import (
	"context"
	"fmt"
	"github.com/oracle/coherence-operator/pkg/flags"
	"github.com/oracle/coherence-operator/pkg/net"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"strings"
)

const (
	zoneLabel = "failure-domain.beta.kubernetes.io/zone"
)

// The logger to use to log messages
var log = logf.Log.WithName("rest-server")

// The k8s client
var k8sClient client.Client

// StartRestServer starts a ReST server to server Coherence Operator requests,
// for example node zone information.
func StartRestServer(mgr manager.Manager, host string, port int32) {
	k8sClient = mgr.GetClient()

	http.HandleFunc("/zone/", getZoneForNode)

	go func() {
		address := fmt.Sprintf("%s:%d", host, port)
		log.Info("http requests serving on " + address)
		err := http.ListenAndServe(address, nil)
		if err != nil {
			panic("ListenAndServe: " + err.Error())
		}
	}()
}

// GetHostAndPort returns the address and port that this endpoint can be reached on by external processes.
func GetHostAndPort() string {
	f := flags.GetOperatorFlags()

	var service string
	var port int32

	if f.ServiceName != "" {
		// use the service name if it was specifically set
		service = f.ServiceName
	} else if f.RestHost != "0.0.0.0" {
		// if no service name was set but ReST is bound to a specific address then use that
		service = f.RestHost
	} else {
		// ReST is bound to 0.0.0.0 so use any of our local addresses.
		// This does not guarantee we're reachable but would be OK in local testing
		ip, err := net.GetLocalAddress()
		if err == nil && ip != nil {
			service = fmt.Sprint(ip.String())
		}
	}

	if f.ServicePort != -1 {
		port = f.ServicePort
	} else {
		port = f.RestPort
	}

	return fmt.Sprintf("%s:%d", service, port)
}

// getZoneForNode is a GET request that returns the zone label for a k8s node
// or a 404 response if the node does not exist.
func getZoneForNode(w http.ResponseWriter, r *http.Request) {
	node := &corev1.Node{}

	pos := strings.LastIndex(r.URL.Path, "/")
	name := r.URL.Path[1+pos:]
	err := k8sClient.Get(context.TODO(), client.ObjectKey{Namespace: "", Name: name}, node)

	if err != nil {
		log.Error(err, "Error getting node "+name+" from k8s")
		if errors.IsNotFound(err) {
			w.WriteHeader(404)
			fmt.Println(w, "Node "+name+" does not exist")
		} else {
			w.WriteHeader(500)
			fmt.Println(w, "An error occurred getting node information from k8s")
		}
	}

	if _, err = fmt.Fprint(w, node.Labels[zoneLabel]); err != nil {
		log.Error(err, "Error writing zone response for node "+name)
	}
}
