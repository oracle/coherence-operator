/*
 * Copyright (c) 2026, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package probe

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	coh "github.com/oracle/coherence-operator/api/v1"
	cohevents "github.com/oracle/coherence-operator/pkg/events"
	"github.com/oracle/coherence-operator/pkg/operator"
	"github.com/spf13/viper"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/tools/events"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestBuildHTTPProbeURL(t *testing.T) {
	tests := []struct {
		name   string
		scheme corev1.URIScheme
		host   string
		port   int
		path   string
		want   string
	}{
		{
			name:   "IPv4 address",
			scheme: corev1.URISchemeHTTP,
			host:   "10.0.0.8",
			port:   6676,
			path:   "/suspend",
			want:   "http://10.0.0.8:6676/suspend",
		},
		{
			name:   "DNS hostname with HTTPS",
			scheme: corev1.URISchemeHTTPS,
			host:   "storage-0.storage.test.svc",
			port:   6676,
			path:   "ready",
			want:   "https://storage-0.storage.test.svc:6676/ready",
		},
		{
			name:   "IPv6 address",
			scheme: corev1.URISchemeHTTP,
			host:   "fd02:0:0:6::8e35",
			port:   6676,
			path:   "/suspend",
			want:   "http://[fd02:0:0:6::8e35]:6676/suspend",
		},
		{
			name:   "path query",
			scheme: corev1.URISchemeHTTP,
			host:   "127.0.0.1",
			port:   8080,
			path:   "/health?verbose=true",
			want:   "http://127.0.0.1:8080/health?verbose=true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := buildHTTPProbeURL(tt.scheme, tt.host, tt.port, tt.path)
			if err != nil {
				t.Fatalf("buildHTTPProbeURL() returned an error: %v", err)
			}
			if got := u.String(); got != tt.want {
				t.Fatalf("buildHTTPProbeURL() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestSuspendServicesReportsResult(t *testing.T) {
	originalSkipServiceSuspend := viper.GetBool(operator.FlagSkipServiceSuspend)
	viper.Set(operator.FlagSkipServiceSuspend, false)
	t.Cleanup(func() {
		viper.Set(operator.FlagSkipServiceSuspend, originalSkipServiceSuspend)
	})

	tests := []struct {
		name         string
		responseCode int
		wantStatus   ServiceSuspendStatus
		wantEvent    string
		wantError    string
		rejectEvent  string
	}{
		{
			name:         "successful suspend",
			responseCode: http.StatusOK,
			wantStatus:   ServiceSuspendSuccessful,
			wantEvent:    "Normal ServiceSuspended suspended Coherence services in StatefulSet test-deployment",
			rejectEvent:  "ServiceSuspendFailed",
		},
		{
			name:         "failed suspend",
			responseCode: http.StatusInternalServerError,
			wantStatus:   ServiceSuspendFailed,
			wantEvent:    "Warning ServiceSuspendFailed failed to suspend Coherence services in StatefulSet test-deployment: required service suspension failed: HTTP probe rejected: status 500",
			wantError:    "required service suspension failed: HTTP probe rejected: status 500",
			rejectEvent:  "ServiceSuspended",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			const responseBody = "private response body containing Authorization: Basic dXNlcjpwYXNzd29yZA=="
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPut {
					t.Errorf("generated suspension used %s", r.Method)
				}
				w.WriteHeader(tt.responseCode)
				_, _ = w.Write([]byte(responseBody))
			}))
			defer server.Close()

			host, port, err := net.SplitHostPort(strings.TrimPrefix(server.URL, "http://"))
			if err != nil {
				t.Fatalf("failed to split test server address: %v", err)
			}

			labels := map[string]string{
				"app":                        "test-deployment",
				operator.LabelTestHostName:   host,
				operator.LabelTestHealthPort: port,
			}
			pod := &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-deployment-0",
					Namespace: "test-namespace",
					Labels:    labels,
				},
				Status: corev1.PodStatus{
					Phase: corev1.PodRunning,
					Conditions: []corev1.PodCondition{
						{
							Type:   corev1.PodReady,
							Status: corev1.ConditionTrue,
						},
					},
				},
			}
			sts := &appsv1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-deployment",
					Namespace: "test-namespace",
				},
				Spec: appsv1.StatefulSetSpec{
					Replicas: ptr.To(int32(1)),
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{"app": "test-deployment"},
					},
				},
			}
			deployment := &coh.Coherence{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-deployment",
					Namespace: "test-namespace",
				},
			}

			recorder := events.NewFakeRecorder(10)
			coherenceProbe := CoherenceProbe{
				Client:        fake.NewClientBuilder().WithObjects(pod).Build(),
				EventRecorder: cohevents.NewOwnedEventRecorder(deployment, recorder),
			}

			got, suspendErr := coherenceProbe.SuspendServicesWithError(context.Background(), deployment, sts)
			if got != tt.wantStatus {
				t.Fatalf("SuspendServicesWithError() = %v, want %v", got, tt.wantStatus)
			}
			if tt.wantError == "" {
				if suspendErr != nil {
					t.Fatalf("SuspendServicesWithError() returned %v", suspendErr)
				}
			} else if suspendErr == nil || suspendErr.Error() != tt.wantError {
				t.Fatalf("SuspendServicesWithError() error = %v, want %q", suspendErr, tt.wantError)
			}

			var recordedEvents []string
			for len(recorder.Events) > 0 {
				recordedEvents = append(recordedEvents, <-recorder.Events)
			}
			if !containsEvent(recordedEvents, tt.wantEvent) {
				t.Errorf("events %q do not contain %q", recordedEvents, tt.wantEvent)
			}
			if containsEventWithReason(recordedEvents, tt.rejectEvent) {
				t.Errorf("events %q unexpectedly contain reason %q", recordedEvents, tt.rejectEvent)
			}
			if containsEventWithReason(recordedEvents, "ProbeActionNoOp") {
				t.Errorf("events %q unexpectedly report a no-op", recordedEvents)
			}
			if suspendErr != nil && strings.Contains(suspendErr.Error(), responseBody) {
				t.Fatalf("response body escaped in error: %v", suspendErr)
			}
			for _, event := range recordedEvents {
				if strings.Contains(event, responseBody) || strings.Contains(event, "dXNlcj") {
					t.Fatalf("response body or credential escaped in event: %q", event)
				}
			}
		})
	}
}

func TestSuspendServicesReportsLegacyNoOpAsSkipped(t *testing.T) {
	originalSkipServiceSuspend := viper.GetBool(operator.FlagSkipServiceSuspend)
	viper.Set(operator.FlagSkipServiceSuspend, false)
	t.Cleanup(func() {
		viper.Set(operator.FlagSkipServiceSuspend, originalSkipServiceSuspend)
	})

	for _, tt := range []struct {
		name   string
		probe  *coh.Probe
		reason string
	}{
		{name: "empty handler", probe: &coh.Probe{}, reason: "the probe handler is empty"},
		{name: "gRPC handler", probe: &coh.Probe{ProbeHandler: corev1.ProbeHandler{GRPC: &corev1.GRPCAction{Port: 1408}}}, reason: "gRPC probe actions are not supported"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			labels := map[string]string{"app": "test-deployment"}
			pod := &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{Name: "test-deployment-0", Namespace: "test-namespace", Labels: labels},
				Status: corev1.PodStatus{
					Phase:      corev1.PodRunning,
					Conditions: []corev1.PodCondition{{Type: corev1.PodReady, Status: corev1.ConditionTrue}},
				},
			}
			sts := &appsv1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{Name: "test-deployment", Namespace: "test-namespace"},
				Spec: appsv1.StatefulSetSpec{
					Replicas: ptr.To(int32(1)),
					Selector: &metav1.LabelSelector{MatchLabels: labels},
				},
			}
			deployment := &coh.Coherence{
				ObjectMeta: metav1.ObjectMeta{Name: "test-deployment", Namespace: "test-namespace"},
				Spec:       coh.CoherenceStatefulSetResourceSpec{SuspendProbe: tt.probe},
			}
			eventRecorder := events.NewFakeRecorder(10)
			coherenceProbe := CoherenceProbe{
				Client:        fake.NewClientBuilder().WithObjects(pod).Build(),
				EventRecorder: cohevents.NewOwnedEventRecorder(deployment, eventRecorder),
			}

			status, err := coherenceProbe.SuspendServicesWithError(context.Background(), deployment, sts)
			if status != ServiceSuspendSkipped || err != nil {
				t.Fatalf("status=%v error=%v, want skipped without error", status, err)
			}

			var recordedEvents []string
			for len(eventRecorder.Events) > 0 {
				recordedEvents = append(recordedEvents, <-eventRecorder.Events)
			}
			if !containsEventWithReason(recordedEvents, "ProbeActionNoOp") {
				t.Fatalf("events %q do not contain ProbeActionNoOp", recordedEvents)
			}
			if containsEventWithReason(recordedEvents, "ServiceSuspended") {
				t.Fatalf("events %q incorrectly report successful suspension", recordedEvents)
			}
			if !containsEventText(recordedEvents, "Configured suspend probe", tt.reason, "no request or command was sent") {
				t.Fatalf("events %q do not explain the no-op", recordedEvents)
			}
		})
	}
}

func TestSuspendServicesReportsPreconditionCause(t *testing.T) {
	originalSkipServiceSuspend := viper.GetBool(operator.FlagSkipServiceSuspend)
	viper.Set(operator.FlagSkipServiceSuspend, false)
	t.Cleanup(func() {
		viper.Set(operator.FlagSkipServiceSuspend, originalSkipServiceSuspend)
	})

	tests := []struct {
		name      string
		podReady  bool
		replicas  int32
		wantCause string
	}{
		{
			name:      "Pod not ready",
			podReady:  false,
			replicas:  1,
			wantCause: "probe precondition failed: Pod test-deployment-0 is not ready (Running)",
		},
		{
			name:      "Pod count mismatch",
			podReady:  true,
			replicas:  2,
			wantCause: "probe precondition failed: found 1 Pods, expected 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conditions := []corev1.PodCondition{}
			if tt.podReady {
				conditions = append(conditions, corev1.PodCondition{Type: corev1.PodReady, Status: corev1.ConditionTrue})
			}
			labels := map[string]string{"app": "test-deployment"}
			pod := &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{Name: "test-deployment-0", Namespace: "test-namespace", Labels: labels},
				Status:     corev1.PodStatus{Phase: corev1.PodRunning, Conditions: conditions},
			}
			sts := &appsv1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{Name: "test-deployment", Namespace: "test-namespace"},
				Spec: appsv1.StatefulSetSpec{
					Replicas: ptr.To(tt.replicas),
					Selector: &metav1.LabelSelector{MatchLabels: labels},
				},
			}
			deployment := &coh.Coherence{ObjectMeta: metav1.ObjectMeta{Name: "test-deployment", Namespace: "test-namespace"}}
			recorder := events.NewFakeRecorder(10)
			coherenceProbe := CoherenceProbe{
				Client:        fake.NewClientBuilder().WithObjects(pod).Build(),
				EventRecorder: cohevents.NewOwnedEventRecorder(deployment, recorder),
			}

			status, err := coherenceProbe.SuspendServicesWithError(context.Background(), deployment, sts)
			wantError := "required service suspension failed: " + tt.wantCause
			if status != ServiceSuspendFailed || err == nil || err.Error() != wantError {
				t.Fatalf("status=%v error=%v, want %v and %q", status, err, ServiceSuspendFailed, wantError)
			}
			wantEvent := "Warning ServiceSuspendFailed failed to suspend Coherence services in StatefulSet test-deployment: " + wantError
			var recordedEvents []string
			for len(recorder.Events) > 0 {
				recordedEvents = append(recordedEvents, <-recorder.Events)
			}
			if !containsEvent(recordedEvents, wantEvent) {
				t.Fatalf("events %q do not contain %q", recordedEvents, wantEvent)
			}
		})
	}
}

func TestStatusHARejectedResponseDoesNotFallThrough(t *testing.T) {
	for _, status := range []int{http.StatusBadRequest, http.StatusServiceUnavailable} {
		t.Run(strconv.Itoa(status), func(t *testing.T) {
			secondCalls := 0
			first := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/ha" {
					t.Errorf("probe path = %q, want /ha", r.URL.Path)
				}
				w.WriteHeader(status)
			}))
			defer first.Close()
			second := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				secondCalls++
				w.WriteHeader(http.StatusOK)
			}))
			defer second.Close()

			pods := corev1.PodList{Items: []corev1.Pod{
				readyPodForProbeServer(t, "storage-0", first.URL),
				readyPodForProbeServer(t, "storage-1", second.URL),
			}}
			sts := &appsv1.StatefulSet{Spec: appsv1.StatefulSetSpec{Replicas: ptr.To(int32(len(pods.Items)))}}
			deployment := &coh.Coherence{}
			coherenceProbe := CoherenceProbe{
				EventRecorder: cohevents.NewOwnedEventRecorder(deployment, events.NewFakeRecorder(10)),
			}

			ok, err := coherenceProbe.executeProbeForPods(
				context.Background(), sts, "test", deployment.Spec.GetDefaultScalingProbe(), pods, pods, probeExecutionOptions{})
			var requestError *RequestError
			if ok || !errors.As(err, &requestError) || requestError.StatusCode != status {
				t.Fatalf("ok=%v err=%v, want rejected status %d", ok, err, status)
			}
			if secondCalls != 0 {
				t.Fatalf("second Pod was probed %d times after status %d", secondCalls, status)
			}
		})
	}
}

func TestStatusHATransportFailureDoesNotFallThrough(t *testing.T) {
	first := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	firstURL := first.URL
	first.Close()

	secondCalls := 0
	second := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		secondCalls++
		w.WriteHeader(http.StatusOK)
	}))
	defer second.Close()

	pods := corev1.PodList{Items: []corev1.Pod{
		readyPodForProbeServer(t, "storage-0", firstURL),
		readyPodForProbeServer(t, "storage-1", second.URL),
	}}
	sts := &appsv1.StatefulSet{Spec: appsv1.StatefulSetSpec{Replicas: ptr.To(int32(len(pods.Items)))}}
	deployment := &coh.Coherence{}
	coherenceProbe := CoherenceProbe{
		EventRecorder: cohevents.NewOwnedEventRecorder(deployment, events.NewFakeRecorder(10)),
	}

	ok, err := coherenceProbe.executeProbeForPods(
		context.Background(), sts, "test", deployment.Spec.GetDefaultScalingProbe(), pods, pods, probeExecutionOptions{})
	var requestError *RequestError
	if ok || !errors.As(err, &requestError) || requestError.Kind != "transport" {
		t.Fatalf("ok=%v err=%v, want authoritative transport failure", ok, err)
	}
	if secondCalls != 0 {
		t.Fatalf("second Pod was probed %d times after transport failure", secondCalls)
	}
}

func TestResolveScalingProbePreservesSource(t *testing.T) {
	for _, spec := range []*coh.CoherenceStatefulSetResourceSpec{
		nil,
		{},
		{Scaling: &coh.ScalingSpec{}},
	} {
		resolved := ResolveScalingProbe(spec)
		if !resolved.Generated || resolved.Probe == nil || resolved.Probe.HTTPGet == nil || resolved.Probe.HTTPGet.Path != "/ha" {
			t.Fatalf("default scaling probe was not identified as generated: %#v", resolved)
		}
	}

	custom := &coh.Probe{ProbeHandler: corev1.ProbeHandler{HTTPGet: &corev1.HTTPGetAction{
		Path:   "/ha",
		Port:   intstr.FromString(coh.PortNameHealth),
		Scheme: corev1.URISchemeHTTP,
	}}}
	resolved := ResolveScalingProbe(&coh.CoherenceStatefulSetResourceSpec{Scaling: &coh.ScalingSpec{Probe: custom}})
	if resolved.Generated || resolved.Probe != custom {
		t.Fatalf("explicit scaling probe lost its source or identity: %#v", resolved)
	}
}

func TestStatusHALegacyNoOpWarnsAndSucceeds(t *testing.T) {
	for _, tt := range []struct {
		name   string
		probe  *coh.Probe
		reason string
	}{
		{name: "empty handler", probe: &coh.Probe{}, reason: "the probe handler is empty"},
		{name: "gRPC handler", probe: &coh.Probe{ProbeHandler: corev1.ProbeHandler{GRPC: &corev1.GRPCAction{Port: 1408}}}, reason: "gRPC probe actions are not supported"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			labels := map[string]string{"app": "storage"}
			pod := &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{Name: "storage-0", Namespace: "test", Labels: labels},
				Status: corev1.PodStatus{
					Phase:      corev1.PodRunning,
					Conditions: []corev1.PodCondition{{Type: corev1.PodReady, Status: corev1.ConditionTrue}},
				},
			}
			sts := &appsv1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{Name: "storage", Namespace: "test"},
				Spec: appsv1.StatefulSetSpec{
					Replicas: ptr.To(int32(1)),
					Selector: &metav1.LabelSelector{MatchLabels: labels},
				},
			}
			deployment := &coh.Coherence{ObjectMeta: metav1.ObjectMeta{Name: "storage", Namespace: "test"}}
			eventRecorder := events.NewFakeRecorder(10)
			coherenceProbe := CoherenceProbe{
				Client:        fake.NewClientBuilder().WithObjects(pod).Build(),
				EventRecorder: cohevents.NewOwnedEventRecorder(deployment, eventRecorder),
			}

			if !coherenceProbe.ExecuteScalingProbe(context.Background(), sts, "storage-wka", ResolvedScalingProbe{Probe: tt.probe}) {
				t.Fatal("legacy no-op StatusHA probe did not retain successful compatibility behavior")
			}

			var recordedEvents []string
			for len(eventRecorder.Events) > 0 {
				recordedEvents = append(recordedEvents, <-eventRecorder.Events)
			}
			if !containsEventWithReason(recordedEvents, "ProbeActionNoOp") ||
				!containsEventText(recordedEvents, "Configured StatusHA probe", tt.reason, "no request or command was sent") {
				t.Fatalf("events %q do not explain the StatusHA no-op", recordedEvents)
			}
		})
	}
}

func TestLegacyNoOpClassification(t *testing.T) {
	for _, tt := range []struct {
		name     string
		probe    *coh.Probe
		wantNoOp bool
	}{
		{name: "nil", probe: nil},
		{name: "empty", probe: &coh.Probe{}, wantNoOp: true},
		{name: "gRPC", probe: &coh.Probe{ProbeHandler: corev1.ProbeHandler{GRPC: &corev1.GRPCAction{Port: 1408}}}, wantNoOp: true},
		{name: "HTTP", probe: &coh.Probe{HTTP: &coh.HTTPAction{}}},
		{name: "HTTP GET", probe: &coh.Probe{ProbeHandler: corev1.ProbeHandler{HTTPGet: &corev1.HTTPGetAction{}}}},
		{name: "exec", probe: &coh.Probe{ProbeHandler: corev1.ProbeHandler{Exec: &corev1.ExecAction{}}}},
		{name: "TCP", probe: &coh.Probe{ProbeHandler: corev1.ProbeHandler{TCPSocket: &corev1.TCPSocketAction{}}}},
		{name: "exec with gRPC", probe: &coh.Probe{ProbeHandler: corev1.ProbeHandler{Exec: &corev1.ExecAction{}, GRPC: &corev1.GRPCAction{Port: 1408}}}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, got := legacyNoOpReason(tt.probe)
			if got != tt.wantNoOp {
				t.Fatalf("legacyNoOpReason() no-op=%t, want %t", got, tt.wantNoOp)
			}
		})
	}
}

func TestStatusHAUsesHTTPSOnlyForGeneratedProbe(t *testing.T) {
	t.Run("generated secure probe uses HTTPS", func(t *testing.T) {
		server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/ha" {
				t.Errorf("probe path = %q, want /ha", r.URL.Path)
			}
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		assertStatusHAForServer(t, server.URL, nil, true)
	})

	t.Run("generated non-secure probe keeps HTTP", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/ha" {
				t.Errorf("probe path = %q, want /ha", r.URL.Path)
			}
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		assertStatusHAForServer(t, server.URL, nil, false)
	})

	for _, port := range []string{"named health port", "numeric application port"} {
		t.Run("explicit HTTP /ha on "+port+" remains HTTP", func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/ha" {
					t.Errorf("probe path = %q, want /ha", r.URL.Path)
				}
				w.WriteHeader(http.StatusOK)
			}))
			defer server.Close()

			_, portText, err := net.SplitHostPort(strings.TrimPrefix(server.URL, "http://"))
			if err != nil {
				t.Fatal(err)
			}
			probePort := intstr.FromString(coh.PortNameHealth)
			if port == "numeric application port" {
				portNumber, err := strconv.Atoi(portText)
				if err != nil {
					t.Fatal(err)
				}
				probePort = intstr.FromInt(portNumber)
			}
			custom := &coh.Probe{ProbeHandler: corev1.ProbeHandler{HTTPGet: &corev1.HTTPGetAction{
				Path:   "/ha",
				Port:   probePort,
				Scheme: corev1.URISchemeHTTP,
			}}}
			assertStatusHAForServer(t, server.URL, custom, true)
		})
	}
}

func assertStatusHAForServer(t *testing.T, address string, custom *coh.Probe, secure bool) {
	t.Helper()
	pod := readyPodForProbeServer(t, "storage-0", address)
	pod.Namespace = "test"
	pod.Labels["app"] = "storage"
	if secure {
		pod.Annotations = map[string]string{coh.HealthMutatorAnnotation: secureHealthMutatorAnnotation(t)}
	}
	sts := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{Name: "storage", Namespace: "test"},
		Spec: appsv1.StatefulSetSpec{
			Replicas: ptr.To(int32(1)),
			Selector: &metav1.LabelSelector{MatchLabels: map[string]string{"app": "storage"}},
		},
	}
	deployment := &coh.Coherence{ObjectMeta: metav1.ObjectMeta{Name: "storage", Namespace: "test"}}
	if custom != nil {
		deployment.Spec.Scaling = &coh.ScalingSpec{Probe: custom}
	}
	coherenceProbe := CoherenceProbe{
		Client:        fake.NewClientBuilder().WithObjects(&pod).Build(),
		EventRecorder: cohevents.NewOwnedEventRecorder(deployment, events.NewFakeRecorder(10)),
	}
	if !coherenceProbe.IsStatusHA(context.Background(), deployment, sts) {
		t.Fatalf("StatusHA failed for secure=%t custom=%#v address=%s", secure, custom, address)
	}
}

func secureHealthMutatorAnnotation(t *testing.T) string {
	t.Helper()
	selector := corev1.SecretKeySelector{LocalObjectReference: corev1.LocalObjectReference{Name: "health"}, Key: "value"}
	connection := &coh.HealthMutatorConnection{
		Secure:    true,
		BasicAuth: &coh.HTTPBasicAuth{Username: selector, Password: selector},
		TLS:       &coh.HTTPClientTLS{CASecret: selector, ServerName: "health.example"},
		ServerTLS: &coh.SSLSpec{
			Enabled:              ptr.To(true),
			KeyStore:             ptr.To("server.jks"),
			KeyStorePasswordFile: ptr.To("store-password"),
			KeyPasswordFile:      ptr.To("key-password"),
		},
	}
	if err := connection.Validate(); err != nil {
		t.Fatalf("invalid secure Health test connection: %v", err)
	}
	data, err := json.Marshal(connection)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func readyPodForProbeServer(t *testing.T, name, address string) corev1.Pod {
	t.Helper()
	host, port, err := net.SplitHostPort(strings.TrimPrefix(strings.TrimPrefix(address, "https://"), "http://"))
	if err != nil {
		t.Fatalf("failed to split test server address: %v", err)
	}
	return corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				operator.LabelTestHostName:   host,
				operator.LabelTestHealthPort: port,
			},
		},
		Status: corev1.PodStatus{
			Phase: corev1.PodRunning,
			Conditions: []corev1.PodCondition{{
				Type: corev1.PodReady, Status: corev1.ConditionTrue,
			}},
		},
	}
}

func containsEvent(recordedEvents []string, expected string) bool {
	for _, event := range recordedEvents {
		if event == expected {
			return true
		}
	}
	return false
}

func containsEventWithReason(recordedEvents []string, reason string) bool {
	for _, event := range recordedEvents {
		fields := strings.Fields(event)
		if len(fields) > 1 && fields[1] == reason {
			return true
		}
	}
	return false
}

func containsEventText(recordedEvents []string, fragments ...string) bool {
	for _, event := range recordedEvents {
		matches := true
		for _, fragment := range fragments {
			if !strings.Contains(event, fragment) {
				matches = false
				break
			}
		}
		if matches {
			return true
		}
	}
	return false
}
