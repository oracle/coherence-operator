/*
 * Copyright (c) 2026, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package probe

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	coh "github.com/oracle/coherence-operator/api/v1"
	cohevents "github.com/oracle/coherence-operator/pkg/events"
	"github.com/oracle/coherence-operator/pkg/operator"
	"github.com/spf13/viper"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
			wantEvent:    "Warning ServiceSuspendFailed failed to suspend Coherence services in StatefulSet test-deployment",
			rejectEvent:  "ServiceSuspended",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tt.responseCode)
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

			if got := coherenceProbe.SuspendServices(context.Background(), deployment, sts); got != tt.wantStatus {
				t.Fatalf("SuspendServices() = %v, want %v", got, tt.wantStatus)
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
		})
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
