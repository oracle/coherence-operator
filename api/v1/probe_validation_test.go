/*
 * Copyright (c) 2026, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */
package v1

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func TestBasicAuthRequiresTargetPodAddress(t *testing.T) {
	selector := func(key string) corev1.SecretKeySelector {
		return corev1.SecretKeySelector{LocalObjectReference: corev1.LocalObjectReference{Name: "auth"}, Key: key}
	}
	probe := &Probe{HTTP: &HTTPAction{
		HTTPEndpoint: HTTPEndpoint{Port: intstr.FromInt(6676), Scheme: corev1.URISchemeHTTPS},
		BasicAuth:    &HTTPBasicAuth{Username: selector("username"), Password: selector("password")},
	}}
	if err := probe.ValidateAction(); err != nil {
		t.Fatalf("Pod-targeted Basic authentication was rejected: %v", err)
	}

	probe.HTTP.TLS = &HTTPClientTLS{CASecret: selector("ca.crt"), ServerName: "health.example"}
	if err := probe.ValidateAction(); err != nil {
		t.Fatalf("TLS server name was treated as a connection host: %v", err)
	}

	probe.HTTP.Host = "health.example"
	if err := probe.ValidateAction(); err == nil || err.Error() != "basic authentication requires the target Pod address; host must be empty" {
		t.Fatalf("custom host with Basic authentication was not rejected: %v", err)
	}

	probe.HTTP.BasicAuth = nil
	if err := probe.ValidateAction(); err != nil {
		t.Fatalf("unauthenticated custom host was rejected: %v", err)
	}
}
