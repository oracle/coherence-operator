/*
 * Copyright (c) 2026, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */
package probe

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"syscall"
	"testing"

	coh "github.com/oracle/coherence-operator/api/v1"
	cohevents "github.com/oracle/coherence-operator/pkg/events"
	"github.com/oracle/coherence-operator/pkg/operator"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/tools/events"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestSanitizedRequestAndSuspensionErrors(t *testing.T) {
	sensitive := errors.New("Authorization: Basic dXNlcjpwYXNzd29yZA==; username=operator; password=private; Secret=auth; response body=private")
	tests := []struct {
		name string
		err  *RequestError
		want string
	}{
		{name: "rejected status", err: &RequestError{StatusCode: http.StatusForbidden, Cause: sensitive}, want: "HTTP probe rejected: status 403"},
		{name: "timeout", err: &RequestError{Kind: "transport or TLS", Cause: context.DeadlineExceeded}, want: "HTTP probe failed: transport timeout"},
		{name: "connection refused", err: &RequestError{Kind: "transport or TLS", Cause: &net.OpError{Op: "dial", Err: syscall.ECONNREFUSED}}, want: "HTTP probe failed: connection refused"},
		{name: "unknown CA", err: &RequestError{Kind: "transport or TLS", Cause: x509.UnknownAuthorityError{}}, want: "HTTP probe failed: TLS certificate signed by an unknown or untrusted CA"},
		{name: "hostname mismatch", err: &RequestError{Kind: "transport or TLS", Cause: x509.HostnameError{Certificate: &x509.Certificate{}, Host: "private.example"}}, want: "HTTP probe failed: TLS certificate hostname mismatch"},
		{name: "general transport", err: &RequestError{Kind: "transport or TLS", Cause: sensitive}, want: "HTTP probe failed: transport or TLS failure"},
		{name: "Secret unavailable", err: &RequestError{Kind: "Secret unavailable", Cause: sensitive}, want: "HTTP probe failed: Secret unavailable"},
		{name: "validation", err: &RequestError{Kind: "configuration", Cause: sensitive}, want: "HTTP probe failed: configuration validation failure"},
		{name: "unknown kind", err: &RequestError{Kind: sensitive.Error(), Cause: sensitive}, want: "HTTP probe failed: request failure"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Fatalf("RequestError.Error() = %q, want %q", got, tt.want)
			}
			wrapped := &SuspensionError{Cause: fmt.Errorf("wrapped: %w", tt.err)}
			if got := wrapped.Error(); got != "required service suspension failed: "+tt.want {
				t.Fatalf("SuspensionError.Error() = %q", got)
			}
			if strings.Contains(wrapped.Error(), "private") || strings.Contains(wrapped.Error(), "Authorization") || strings.Contains(wrapped.Error(), "dXNlcj") {
				t.Fatalf("sensitive data escaped sanitizer: %q", wrapped.Error())
			}
			if !errors.Is(wrapped, sensitive) && tt.err.Cause == sensitive {
				t.Fatalf("errors.Is() lost the wrapped cause: %v", wrapped)
			}
			var requestError *RequestError
			if !errors.As(wrapped, &requestError) || requestError != tt.err {
				t.Fatalf("errors.As() lost RequestError: %v", wrapped)
			}
		})
	}
}

func TestSanitizedSuspensionCauseErrors(t *testing.T) {
	sensitive := errors.New("Secret auth contained password private")
	tests := []struct {
		name  string
		cause *SuspensionCauseError
		want  string
	}{
		{
			name:  "Pod not ready",
			cause: &SuspensionCauseError{Kind: SuspensionCausePodNotReady, PodName: "storage-0", PodPhase: "Running"},
			want:  "probe precondition failed: Pod storage-0 is not ready (Running)",
		},
		{
			name:  "Pod phase is sanitized",
			cause: &SuspensionCauseError{Kind: SuspensionCausePodNotReady, PodName: "storage-0", PodPhase: sensitive.Error()},
			want:  "probe precondition failed: Pod storage-0 is not ready (unknown state)",
		},
		{
			name:  "Pod count mismatch",
			cause: &SuspensionCauseError{Kind: SuspensionCausePodCountMismatch, ActualPods: 1, ExpectedPods: 3},
			want:  "probe precondition failed: found 1 Pods, expected 3",
		},
		{
			name:  "workload lookup",
			cause: &SuspensionCauseError{Kind: SuspensionCauseWorkloadLookupFailed, Cause: sensitive},
			want:  "workload state lookup failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wrapped := &SuspensionError{Cause: tt.cause}
			want := "required service suspension failed: " + tt.want
			if got := wrapped.Error(); got != want {
				t.Fatalf("SuspensionError.Error() = %q, want %q", got, want)
			}
			if strings.Contains(wrapped.Error(), "private") || strings.Contains(wrapped.Error(), "password") {
				t.Fatalf("sensitive data escaped sanitizer: %q", wrapped.Error())
			}
			if tt.cause.Cause == sensitive && !errors.Is(wrapped, sensitive) {
				t.Fatalf("errors.Is() lost the wrapped cause: %v", wrapped)
			}
			var causeError *SuspensionCauseError
			if !errors.As(wrapped, &causeError) || causeError != tt.cause {
				t.Fatalf("errors.As() lost SuspensionCauseError: %v", wrapped)
			}
		})
	}
}

func actionForServer(t *testing.T, address string) *coh.Probe {
	t.Helper()
	host, port, err := net.SplitHostPort(strings.TrimPrefix(strings.TrimPrefix(address, "https://"), "http://"))
	if err != nil {
		t.Fatal(err)
	}
	n, _ := strconv.Atoi(port)
	return &coh.Probe{HTTP: &coh.HTTPAction{Method: "PUT", HTTPEndpoint: coh.HTTPEndpoint{Host: host, Port: intstr.FromInt(n), Path: "/suspend"}}}
}

func TestHTTPActionDispatchAndRejections(t *testing.T) {
	for _, code := range []int{200, 201, 302, 401, 403, 404, 405, 500} {
		t.Run(strconv.Itoa(code), func(t *testing.T) {
			calls := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls++
				if r.Method != "PUT" || r.URL.Path != "/suspend" {
					t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
				}
				w.Header().Set("Location", "/suspend")
				w.WriteHeader(code)
			}))
			defer server.Close()
			p := &CoherenceProbe{}
			ok, err := p.RunProbe(context.Background(), corev1.Pod{}, "test", actionForServer(t, server.URL))
			if ok != (code == 200) || (err == nil) != (code == 200) || calls != 1 {
				t.Fatalf("ok=%v err=%v calls=%d", ok, err, calls)
			}
			if err != nil {
				var failure *RequestError
				if !errors.As(err, &failure) || failure.StatusCode != code {
					t.Fatalf("lost status: %v", err)
				}
			}
		})
	}
}

func TestHTTPActionVerifiedTLSAndCredentials(t *testing.T) {
	calls := 0
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		user, password, ok := r.BasicAuth()
		if !ok || user != "operator" || password != "private:password" {
			w.WriteHeader(401)
			return
		}
		w.WriteHeader(200)
	}))
	defer server.Close()
	selector := func(key string) corev1.SecretKeySelector {
		return corev1.SecretKeySelector{LocalObjectReference: corev1.LocalObjectReference{Name: "auth"}, Key: key}
	}
	secret := &corev1.Secret{ObjectMeta: metav1.ObjectMeta{Name: "auth", Namespace: "test"}, Data: map[string][]byte{
		"user": []byte("operator"), "password": []byte("private:password"), "ca": pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: server.Certificate().Raw}),
	}}
	p := &CoherenceProbe{Client: fake.NewClientBuilder().WithObjects(secret).Build()}
	pod := corev1.Pod{ObjectMeta: metav1.ObjectMeta{Namespace: "test"}}
	action := actionForServer(t, server.URL)
	action.HTTP.Scheme = corev1.URISchemeHTTPS
	action.HTTP.BasicAuth = &coh.HTTPBasicAuth{Username: selector("user"), Password: selector("password")}
	action.HTTP.TLS = &coh.HTTPClientTLS{CASecret: selector("ca")}
	if ok, err := p.RunProbe(context.Background(), pod, "test", action); ok || err == nil || calls != 0 {
		t.Fatalf("authenticated custom host was not rejected before dispatch: ok=%v err=%v calls=%d", ok, err, calls)
	}
	pod.Status.PodIP = action.HTTP.Host
	action.HTTP.Host = ""
	if ok, err := p.RunProbe(context.Background(), pod, "test", action); !ok || err != nil {
		t.Fatalf("valid request failed: %v", err)
	}
	action.HTTP.TLS.ServerName = "wrong.example"
	if ok, err := p.RunProbe(context.Background(), pod, "test", action); ok || err == nil || err.Error() != "HTTP probe failed: TLS certificate hostname mismatch" {
		t.Fatalf("wrong TLS name: ok=%v err=%v", ok, err)
	}
	action.HTTP.TLS.ServerName = ""
	action.HTTP.BasicAuth.Password.Key = "missing"
	if ok, err := p.RunProbe(context.Background(), pod, "test", action); ok || err == nil || err.Error() != "HTTP probe failed: Secret key missing or empty" || strings.Contains(err.Error(), "private:password") {
		t.Fatalf("missing key: %v %v", ok, err)
	}
	action.HTTP.BasicAuth.Password.Name = "missing-auth"
	action.HTTP.BasicAuth.Password.Key = "password"
	if ok, err := p.RunProbe(context.Background(), pod, "test", action); ok || err == nil || err.Error() != "HTTP probe failed: Secret unavailable" {
		t.Fatalf("missing Secret: %v %v", ok, err)
	}
	action.HTTP.Scheme = corev1.URISchemeHTTP
	if ok, err := p.RunProbe(context.Background(), pod, "test", action); ok || err == nil {
		t.Fatal("accepted plaintext credentials")
	}
	if calls != 1 {
		t.Fatalf("invalid configurations reached server: %d", calls)
	}
}

func TestLegacyEmptyAndGRPCProbesRemainSuccessfulNoOps(t *testing.T) {
	p := &CoherenceProbe{}
	for _, action := range []*coh.Probe{{}, {ProbeHandler: corev1.ProbeHandler{GRPC: &corev1.GRPCAction{Port: 1234}}}} {
		if ok, err := p.RunProbe(context.Background(), corev1.Pod{}, "test", action); !ok || err != nil {
			t.Fatalf("legacy no-op action failed: %v %v", ok, err)
		}
	}
	if ok, err := p.RunProbe(context.Background(), corev1.Pod{}, "test", nil); ok || err == nil {
		t.Fatalf("nil action succeeded: %v %v", ok, err)
	}
}

func TestSupportedLegacyProbeTakesPrecedenceOverGRPC(t *testing.T) {
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		calls++
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	host, portText, err := net.SplitHostPort(strings.TrimPrefix(server.URL, "http://"))
	if err != nil {
		t.Fatal(err)
	}
	port, err := strconv.Atoi(portText)
	if err != nil {
		t.Fatal(err)
	}
	owner := &coh.Coherence{}
	eventRecorder := events.NewFakeRecorder(10)
	p := &CoherenceProbe{EventRecorder: cohevents.NewOwnedEventRecorder(owner, eventRecorder)}
	action := &coh.Probe{ProbeHandler: corev1.ProbeHandler{
		HTTPGet: &corev1.HTTPGetAction{Host: host, Port: intstr.FromInt(port), Path: "/"},
		GRPC:    &corev1.GRPCAction{Port: 1408},
	}}

	if ok, err := p.RunProbe(context.Background(), corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "storage-0"}}, "test", action); !ok || err != nil {
		t.Fatalf("supported HTTP GET action did not take precedence: ok=%t err=%v", ok, err)
	}
	if calls != 1 {
		t.Fatalf("HTTP GET action calls=%d, want 1", calls)
	}
	for len(eventRecorder.Events) > 0 {
		if event := <-eventRecorder.Events; strings.Contains(event, "ProbeActionNoOp") {
			t.Fatalf("supported mixed probe emitted no-op warning: %q", event)
		}
	}
}

func TestNewHTTPProbeCannotBeCombinedWithLegacyAction(t *testing.T) {
	p := &CoherenceProbe{}
	action := &coh.Probe{
		HTTP:         &coh.HTTPAction{},
		ProbeHandler: corev1.ProbeHandler{Exec: &corev1.ExecAction{Command: []string{"true"}}},
	}
	if ok, err := p.RunProbe(context.Background(), corev1.Pod{}, "test", action); ok || err == nil {
		t.Fatalf("mixed new and legacy actions succeeded: %v %v", ok, err)
	}
}

func TestRejectedSuspensionDoesNotContactSecondPod(t *testing.T) {
	for _, code := range []int{401, 403, 404, 405, 302} {
		t.Run(strconv.Itoa(code), func(t *testing.T) {
			calls := 0
			first := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(code) }))
			defer first.Close()
			second := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls++; w.WriteHeader(200) }))
			defer second.Close()
			pods := corev1.PodList{}
			for i, address := range []string{first.URL, second.URL} {
				action := actionForServer(t, address)
				pods.Items = append(pods.Items, corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: strconv.Itoa(i), Labels: map[string]string{operator.LabelTestHostName: action.HTTP.Host, operator.LabelTestHealthPort: action.HTTP.Port.String()}}, Status: corev1.PodStatus{Phase: corev1.PodRunning, Conditions: []corev1.PodCondition{{Type: corev1.PodReady, Status: corev1.ConditionTrue}}}})
			}
			sts := &appsv1.StatefulSet{Spec: appsv1.StatefulSetSpec{Replicas: ptr.To(int32(2))}}
			deployment := &coh.Coherence{}
			p := &CoherenceProbe{EventRecorder: cohevents.NewOwnedEventRecorder(deployment, events.NewFakeRecorder(20))}
			ok, err := p.executeProbeForPods(context.Background(), sts, "test", deployment.Spec.GetDefaultSuspendProbe(), pods, pods, probeExecutionOptions{suspension: true})
			var rejected *RequestError
			if ok || !errors.As(err, &rejected) || rejected.StatusCode != code || calls != 0 {
				t.Fatalf("ok=%v err=%v second calls=%d", ok, err, calls)
			}
		})
	}
}
