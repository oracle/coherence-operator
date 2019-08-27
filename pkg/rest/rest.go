// The rest package provides a ReST server for the Coherence Operator.
package rest

import (
	"fmt"
	"github.com/oracle/coherence-operator/pkg/flags"
	"github.com/oracle/coherence-operator/pkg/net"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"net/http"
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
var k8sClient *kubernetes.Clientset

// StartRestServer starts a ReST server to server Coherence Operator requests,
// for example node zone information.
func StartRestServer(mgr manager.Manager, host string, port int32) {
	k8sClient = kubernetes.NewForConfigOrDie(mgr.GetConfig())

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

// getZoneForNode is a GET request that returns the zone label for a k8s node.
func getZoneForNode(w http.ResponseWriter, r *http.Request) {
	var zone string

	pos := strings.LastIndex(r.URL.Path, "/")
	name := r.URL.Path[1+pos:]
	node, err := k8sClient.CoreV1().Nodes().Get(name, metav1.GetOptions{})

	if err == nil {
		zone = node.Labels[zoneLabel]
	} else {
		log.Error(err, "Error getting node "+name+" from k8s")
		zone = ""
	}

	w.WriteHeader(200)
	if _, err = fmt.Fprint(w, zone); err != nil {
		log.Error(err, "Error writing zone response for node "+name)
	}
}
