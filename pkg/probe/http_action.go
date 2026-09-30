/*
 * Copyright (c) 2026, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */
package probe

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"syscall"

	coh "github.com/oracle/coherence-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// RequestError preserves failure classification without exposing credentials or response bodies.
// Failures are terminal within a probe attempt, including ambiguous transport failures.
type RequestError struct {
	Kind       string
	StatusCode int
	Cause      error
}

func (e *RequestError) Error() string {
	if e.StatusCode != 0 {
		return fmt.Sprintf("HTTP probe rejected: status %d", e.StatusCode)
	}
	return "HTTP probe failed: " + e.diagnostic()
}
func (e *RequestError) Unwrap() error { return e.Cause }

// diagnostic returns only controlled text. The wrapped cause remains available
// through Unwrap, but may contain request, credential, Secret or response data.
func (e *RequestError) diagnostic() string {
	switch e.Kind {
	case "transport timeout":
		return "transport timeout"
	case "transport", "transport or TLS":
		return transportDiagnostic(e.Cause)
	case "configuration":
		return "configuration validation failure"
	case "Secret client unavailable", "Secret unavailable", "Secret key missing or empty",
		"port configuration", "endpoint configuration", "request configuration",
		"invalid CA bundle", "invalid Basic username", "rejected response", "uncertain response",
		"unknown active Pod connection", "invalid active Pod connection",
		"managed security requires the generated suspend action", "missing active Pod connection",
		"unsupported action":
		return e.Kind
	default:
		return "request failure"
	}
}

func transportDiagnostic(err error) string {
	if errors.Is(err, context.DeadlineExceeded) {
		return "transport timeout"
	}
	var timeout interface{ Timeout() bool }
	if errors.As(err, &timeout) && timeout.Timeout() {
		return "transport timeout"
	}
	var unknownAuthority x509.UnknownAuthorityError
	if errors.As(err, &unknownAuthority) {
		return "TLS certificate signed by an unknown or untrusted CA"
	}
	var hostnameError x509.HostnameError
	if errors.As(err, &hostnameError) {
		return "TLS certificate hostname mismatch"
	}
	if errors.Is(err, syscall.ECONNREFUSED) {
		return "connection refused"
	}
	var netError *net.OpError
	if errors.As(err, &netError) {
		return "transport failure"
	}
	return "transport or TLS failure"
}

// SuspensionCauseKind identifies safe suspension diagnostics that are not HTTP
// request failures.
type SuspensionCauseKind string

const (
	SuspensionCausePodNotReady          SuspensionCauseKind = "PodNotReady"
	SuspensionCausePodCountMismatch     SuspensionCauseKind = "PodCountMismatch"
	SuspensionCauseProbeUnavailable     SuspensionCauseKind = "ProbeUnavailable"
	SuspensionCauseWorkloadLookupFailed SuspensionCauseKind = "WorkloadLookupFailed"
)

// SuspensionCauseError carries structured, sanitized context while preserving
// an underlying error for errors.Is/errors.As callers.
type SuspensionCauseError struct {
	Kind         SuspensionCauseKind
	PodName      string
	PodPhase     string
	ActualPods   int32
	ExpectedPods int32
	Cause        error
}

func (e *SuspensionCauseError) Error() string {
	switch e.Kind {
	case SuspensionCausePodNotReady:
		return fmt.Sprintf("probe precondition failed: Pod %s is not ready (%s)", e.PodName, safePodPhase(e.PodPhase))
	case SuspensionCausePodCountMismatch:
		return fmt.Sprintf("probe precondition failed: found %d Pods, expected %d", e.ActualPods, e.ExpectedPods)
	case SuspensionCauseProbeUnavailable:
		return "probe precondition failed: no running Pod could execute the probe"
	case SuspensionCauseWorkloadLookupFailed:
		return "workload state lookup failed"
	default:
		return "suspension failure"
	}
}

func (e *SuspensionCauseError) Unwrap() error { return e.Cause }

func safePodPhase(phase string) string {
	switch phase {
	case "Pending", "Running", "Succeeded", "Failed", "Unknown", "Terminating":
		return phase
	default:
		return "unknown state"
	}
}

func (in *CoherenceProbe) secretValue(ctx context.Context, namespace string, ref corev1.SecretKeySelector) ([]byte, error) {
	if in.Client == nil {
		return nil, &RequestError{Kind: "Secret client unavailable"}
	}
	secret := &corev1.Secret{}
	if err := in.Client.Get(ctx, client.ObjectKey{Namespace: namespace, Name: ref.Name}, secret); err != nil {
		return nil, &RequestError{Kind: "Secret unavailable", Cause: err}
	}
	value, ok := secret.Data[ref.Key]
	if !ok || len(value) == 0 {
		return nil, &RequestError{Kind: "Secret key missing or empty"}
	}
	return value, nil
}

// ProbeUsingHTTPAction executes a verified, context-bound HTTP action without redirects.
func (in *CoherenceProbe) ProbeUsingHTTPAction(ctx context.Context, pod corev1.Pod, svc string, handler *coh.Probe) (bool, error) {
	if err := handler.ValidateAction(); err != nil {
		return false, &RequestError{Kind: "configuration", Cause: err}
	}
	ctx, cancel := context.WithTimeout(ctx, handler.GetTimeout())
	defer cancel()
	a := handler.HTTP
	scheme := a.Scheme
	if scheme == "" {
		scheme = corev1.URISchemeHTTP
	}
	host := a.Host
	if host == "" {
		host = in.GetPodIpOrHostName(pod)
	}
	port, err := in.findPort(pod, a.Port)
	if err != nil {
		return false, &RequestError{Kind: "port configuration", Cause: err}
	}
	u, err := buildHTTPProbeURL(scheme, host, port, a.Path)
	if err != nil {
		return false, &RequestError{Kind: "endpoint configuration", Cause: err}
	}
	method := a.Method
	if method == "" {
		method = http.MethodGet
	}
	req, err := http.NewRequestWithContext(ctx, method, u.String(), nil)
	if err != nil {
		return false, &RequestError{Kind: "request configuration", Cause: err}
	}
	req.Host = in.GetPodHostName(pod, svc)
	for _, h := range a.HTTPHeaders {
		if strings.EqualFold(h.Name, "Host") {
			req.Host = h.Value
		} else {
			req.Header.Add(h.Name, h.Value)
		}
	}
	config := &tls.Config{MinVersion: tls.VersionTLS12}
	if a.TLS != nil {
		pem, err := in.secretValue(ctx, pod.Namespace, a.TLS.CASecret)
		if err != nil {
			return false, err
		}
		config.RootCAs = x509.NewCertPool()
		if !config.RootCAs.AppendCertsFromPEM(pem) {
			return false, &RequestError{Kind: "invalid CA bundle"}
		}
		config.ServerName = a.TLS.ServerName
	}
	if a.BasicAuth != nil {
		username, err := in.secretValue(ctx, pod.Namespace, a.BasicAuth.Username)
		if err != nil {
			return false, err
		}
		password, err := in.secretValue(ctx, pod.Namespace, a.BasicAuth.Password)
		if err != nil {
			return false, err
		}
		if strings.ContainsAny(string(username), ":\r\n") {
			return false, &RequestError{Kind: "invalid Basic username"}
		}
		req.SetBasicAuth(string(username), string(password))
	}
	transport := &http.Transport{TLSClientConfig: config, DisableKeepAlives: true}
	defer transport.CloseIdleConnections()
	httpClient := &http.Client{Transport: transport, CheckRedirect: func(_ *http.Request, _ []*http.Request) error { return http.ErrUseLastResponse }}
	res, err := httpClient.Do(req)
	if err != nil {
		return false, &RequestError{Kind: "transport or TLS", Cause: err}
	}
	defer func() {
		_ = res.Body.Close()
	}()
	if res.StatusCode != http.StatusOK {
		return false, &RequestError{Kind: "rejected response", StatusCode: res.StatusCode}
	}
	if _, err = io.Copy(io.Discard, io.LimitReader(res.Body, 64*1024)); err != nil {
		return false, &RequestError{Kind: "uncertain response", Cause: err}
	}
	return true, nil
}

// SuspensionError preserves the cause for diagnostics and the existing lifecycle recovery policy.
type SuspensionError struct{ Cause error }

func (e *SuspensionError) Error() string {
	if e.Cause == nil {
		return "required service suspension failed"
	}
	var r *RequestError
	if errors.As(e.Cause, &r) {
		return "required service suspension failed: " + r.Error()
	}
	var diagnostic *SuspensionCauseError
	if errors.As(e.Cause, &diagnostic) {
		return "required service suspension failed: " + diagnostic.Error()
	}
	return "required service suspension failed: suspension failure"
}
func (e *SuspensionError) Unwrap() error { return e.Cause }
