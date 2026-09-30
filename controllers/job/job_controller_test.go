/*
 * Copyright (c) 2026, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package job

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/go-logr/logr"
	coh "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/pkg/clients"
	"github.com/oracle/coherence-operator/pkg/operator"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/rest"
	k8sevents "k8s.io/client-go/tools/events"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/recorder"
)

type probeTestManager struct {
	manager.Manager
	client   client.Client
	recorder recorder.EventRecorder
}

func (m *probeTestManager) GetClient() client.Client {
	return m.client
}

func (m *probeTestManager) GetConfig() *rest.Config {
	return &rest.Config{}
}

func (m *probeTestManager) GetEventRecorder(string) recorder.EventRecorder {
	return m.recorder
}

func TestMaybeExecuteProbeRecordsHTTPResult(t *testing.T) {
	for _, tc := range []struct {
		name       string
		statusCode int
		success    bool
		errorText  string
	}{
		{name: "success", statusCode: http.StatusOK, success: true},
		{name: "failure", statusCode: http.StatusServiceUnavailable, errorText: "HTTP probe rejected: status 503"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tc.statusCode)
			}))
			defer server.Close()

			host, port := serverAddress(t, server.URL)
			reconciler, job, deployment, _ := readyActionTestObjects(t, host, coh.Probe{ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{Path: "/test", Port: intstr.FromInt(port)},
			}})

			statuses, err := reconciler.maybeExecuteProbe(context.Background(), job, deployment, logr.Discard())
			if err != nil {
				t.Fatal(err)
			}
			if len(statuses) != 1 || statuses[0].Success == nil || *statuses[0].Success != tc.success {
				t.Fatalf("statuses=%+v", statuses)
			}
			if tc.errorText == "" {
				if statuses[0].Error != nil {
					t.Fatalf("unexpected error: %s", *statuses[0].Error)
				}
			} else if statuses[0].Error == nil || *statuses[0].Error != tc.errorText {
				t.Fatalf("error=%v, want %q", statuses[0].Error, tc.errorText)
			}
		})
	}
}

func TestMaybeExecuteProbeRecordsTCPFailure(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	host, value, err := net.SplitHostPort(listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	port, err := strconv.Atoi(value)
	if err != nil {
		t.Fatal(err)
	}
	if err = listener.Close(); err != nil {
		t.Fatal(err)
	}

	reconciler, job, deployment, _ := readyActionTestObjects(t, host, coh.Probe{ProbeHandler: corev1.ProbeHandler{
		TCPSocket: &corev1.TCPSocketAction{Port: intstr.FromInt(port)},
	}})
	statuses, err := reconciler.maybeExecuteProbe(context.Background(), job, deployment, logr.Discard())
	if err != nil {
		t.Fatal(err)
	}
	if len(statuses) != 1 || statuses[0].Success == nil || *statuses[0].Success {
		t.Fatalf("statuses=%+v", statuses)
	}
	const expected = "probe returned an unsuccessful result"
	if statuses[0].Error == nil || *statuses[0].Error != expected {
		t.Fatalf("error=%v, want %q", statuses[0].Error, expected)
	}
}

func TestMaybeExecuteLegacyNoOpProbeWarns(t *testing.T) {
	for _, tt := range []struct {
		name   string
		probe  coh.Probe
		reason string
	}{
		{name: "empty handler", probe: coh.Probe{}, reason: "the probe handler is empty"},
		{name: "gRPC handler", probe: coh.Probe{ProbeHandler: corev1.ProbeHandler{GRPC: &corev1.GRPCAction{Port: 1408}}}, reason: "gRPC probe actions are not supported"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			reconciler, job, deployment, eventRecorder := readyActionTestObjects(t, "127.0.0.1", tt.probe)

			statuses, err := reconciler.maybeExecuteProbe(context.Background(), job, deployment, logr.Discard())
			if err != nil {
				t.Fatal(err)
			}
			if len(statuses) != 1 || statuses[0].Success == nil || !*statuses[0].Success || statuses[0].Error != nil {
				t.Fatalf("statuses=%+v", statuses)
			}

			select {
			case event := <-eventRecorder.Events:
				if !strings.Contains(event, "Warning ProbeActionNoOp") || !strings.Contains(event, tt.reason) {
					t.Fatalf("event=%q, want ProbeActionNoOp containing %q", event, tt.reason)
				}
			default:
				t.Fatal("expected ProbeActionNoOp warning event")
			}
		})
	}
}

func TestShouldExecuteProbeOncePerReadyTransition(t *testing.T) {
	reconciler := &ReconcileJob{}
	readyTime := metav1.NewTime(time.Now())
	condition := &corev1.PodCondition{Type: corev1.PodReady, Status: corev1.ConditionTrue, LastTransitionTime: readyTime}

	if !reconciler.shouldExecuteProbe(coh.CoherenceJobProbeStatus{}, condition) {
		t.Fatal("first probe attempt was skipped")
	}
	if reconciler.shouldExecuteProbe(coh.CoherenceJobProbeStatus{LastReadyTime: &readyTime, Success: ptr.To(false)}, condition) {
		t.Fatal("failed probe was repeated without a new ready transition")
	}
	if reconciler.shouldExecuteProbe(coh.CoherenceJobProbeStatus{LastReadyTime: &readyTime, Success: ptr.To(true)}, condition) {
		t.Fatal("successful probe was repeated without a new ready transition")
	}

	later := metav1.NewTime(readyTime.Add(time.Second))
	condition.LastTransitionTime = later
	if !reconciler.shouldExecuteProbe(coh.CoherenceJobProbeStatus{LastReadyTime: &readyTime, Success: ptr.To(true)}, condition) {
		t.Fatal("probe was not executed after a new ready transition")
	}
}

func readyActionTestObjects(t *testing.T, host string, action coh.Probe) (*ReconcileJob, *batchv1.Job, *coh.CoherenceJob, *k8sevents.FakeRecorder) {
	t.Helper()
	selector := map[string]string{"job": "ready-action"}
	readyTime := metav1.NewTime(time.Now())
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "ready-action-0",
			Namespace: "test",
			Labels: map[string]string{
				"job":                      selector["job"],
				operator.LabelTestHostName: host,
			},
		},
		Status: corev1.PodStatus{Conditions: []corev1.PodCondition{{
			Type: corev1.PodReady, Status: corev1.ConditionTrue, LastTransitionTime: readyTime,
		}}},
	}
	eventRecorder := k8sevents.NewFakeRecorder(10)
	manager := &probeTestManager{client: fake.NewClientBuilder().WithObjects(pod).Build(), recorder: eventRecorder}
	reconciler := &ReconcileJob{}
	reconciler.SetCommonReconciler("test.Job", manager, clients.ClientSet{})

	deployment := &coh.CoherenceJob{
		ObjectMeta: metav1.ObjectMeta{Name: "ready-action", Namespace: "test"},
		Spec: coh.CoherenceJobResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{Replicas: ptr.To(int32(1))},
			Cluster:               "test",
			ReadyAction:           &coh.CoherenceJobProbe{Probe: action},
		},
	}
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{Name: "ready-action", Namespace: "test"},
		Spec:       batchv1.JobSpec{Selector: &metav1.LabelSelector{MatchLabels: selector}},
		Status:     batchv1.JobStatus{Ready: ptr.To(int32(1))},
	}
	return reconciler, job, deployment, eventRecorder
}

func serverAddress(t *testing.T, rawURL string) (string, int) {
	t.Helper()
	u, err := url.Parse(rawURL)
	if err != nil {
		t.Fatal(err)
	}
	host, value, err := net.SplitHostPort(u.Host)
	if err != nil {
		t.Fatal(err)
	}
	port, err := strconv.Atoi(value)
	if err != nil {
		t.Fatal(err)
	}
	return host, port
}
