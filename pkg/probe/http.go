/*
 * Copyright (c) 2020, 2025, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package probe

import (
	"crypto/tls"
	"fmt"
	"github.com/oracle/coherence-operator/pkg/operator"
	"io"
	"k8s.io/apimachinery/pkg/util/net"
	"net/http"
	"net/url"
	"time"
)

// NewHTTPProbe creates Probe that will skip TLS verification while probing.
func NewHTTPProbe() HTTPProbe {
	tlsConfig := &tls.Config{InsecureSkipVerify: true}
	return NewHTTPProbeWithTLSConfig(tlsConfig)
}

// NewHTTPProbeWithTLSConfig takes tls config as parameter.
func NewHTTPProbeWithTLSConfig(config *tls.Config) HTTPProbe {
	transport := net.SetTransportDefaults(&http.Transport{TLSClientConfig: config, DisableKeepAlives: true})
	return httpProbe{transport}
}

// HTTPProbe is an interface that defines the Probe function for doing HTTP readiness/liveness checks.
type HTTPProbe interface {
	Probe(url *url.URL, headers http.Header, timeout time.Duration) (Result, string, error)
}

type httpProbe struct {
	transport *http.Transport
}

// Probe returns a ProbeRunner capable of running an HTTP check.
func (pr httpProbe) Probe(url *url.URL, headers http.Header, timeout time.Duration) (Result, string, error) {
	return DoHTTPProbe(url, headers, &http.Client{Timeout: timeout, Transport: pr.transport})
}

// GetHTTPInterface is an interface for making HTTP requests, that returns a response and error.
type GetHTTPInterface interface {
	Do(req *http.Request) (*http.Response, error)
}

// DoHTTPProbe checks if a GET request to the url succeeds.
// If the HTTP response code is successful (i.e. 400 > code >= 200), it returns Success.
// If the HTTP response code is unsuccessful or HTTP communication fails, it returns Failure.
// This is exported because some other packages may want to do direct HTTP probes.
func DoHTTPProbe(url *url.URL, headers http.Header, client GetHTTPInterface) (Result, string, error) {
	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		// Convert errors into failures to catch timeouts.
		return Failure, err.Error(), nil
	}
	if _, ok := headers["User-Agent"]; !ok {
		if headers == nil {
			headers = http.Header{}
		}
		// explicitly set User-Agent, so it's not set to default Go value
		headers.Set("User-Agent", fmt.Sprintf("coherence-operator/%s", operator.GetVersion()))
	}
	req.Header = headers
	if headers.Get("Host") != "" {
		req.Host = headers.Get("Host")
	}
	log.Info("Executing HTTP Probe", "URL", url.String(), "Headers", headers)
	res, err := client.Do(req)
	if err != nil {
		// Convert errors into failures to catch timeouts.
		return Failure, err.Error(), nil
	}
	defer closeBody(res)
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return Failure, "", err
	}
	body := string(b)
	if res.StatusCode >= http.StatusOK && res.StatusCode < http.StatusBadRequest {
		return Success, body, nil
	}
	return Failure, fmt.Sprintf("HTTP probe failed with statuscode: %d", res.StatusCode), nil
}

func closeBody(res *http.Response) {
	// close the response body, ignoring any errors
	_ = res.Body.Close()
}
