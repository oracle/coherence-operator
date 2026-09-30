/*
 * Copyright (c) 2026, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */
package finalizer

import (
	"context"
	"errors"
	coh "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/pkg/probe"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/events"
	"k8s.io/utils/ptr"
	"net"
	"net/http"
	"net/http/httptest"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"strconv"
	"testing"
)

func TestDesiredZeroPreservesHistoricalFinalizationSkip(t *testing.T) {
	c := &coh.Coherence{ObjectMeta: metav1.ObjectMeta{Name: "test", Namespace: "test"}}
	c.Spec.Replicas = ptr.To(int32(0))
	manager := &FinalizerManager{EventRecorder: events.NewFakeRecorder(20)}
	lookup := func(context.Context, string, string) (*appsv1.StatefulSet, bool, error) {
		t.Fatal("desired-zero finalization unexpectedly inspected workload state")
		return nil, false, nil
	}
	if err := manager.FinalizeDeployment(context.Background(), c, lookup); err != nil {
		t.Fatalf("historical desired-zero skip failed: %v", err)
	}
}

func TestNoReadyPodsPreservesHistoricalFinalizationSkip(t *testing.T) {
	c := &coh.Coherence{ObjectMeta: metav1.ObjectMeta{Name: "test", Namespace: "test"}}
	sts := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{Name: "test", Namespace: "test"},
		Spec:       appsv1.StatefulSetSpec{Replicas: ptr.To(int32(1))},
		Status:     appsv1.StatefulSetStatus{Replicas: 1, ReadyReplicas: 0},
	}
	manager := &FinalizerManager{EventRecorder: events.NewFakeRecorder(20)}
	lookup := func(context.Context, string, string) (*appsv1.StatefulSet, bool, error) { return sts, true, nil }
	if err := manager.FinalizeDeployment(context.Background(), c, lookup); err != nil {
		t.Fatalf("historical no-ready-Pods skip failed: %v", err)
	}
}

func TestDeletionBypassAnnotationPreservesPresenceSemantics(t *testing.T) {
	for _, value := range []string{"true", "false", ""} {
		t.Run("value="+value, func(t *testing.T) {
			c := &coh.Coherence{ObjectMeta: metav1.ObjectMeta{Name: "test", Namespace: "test", Annotations: map[string]string{"coherence.oracle.com/finalizer-bypass": value}}}
			manager := &FinalizerManager{EventRecorder: events.NewFakeRecorder(10)}
			lookup := func(context.Context, string, string) (*appsv1.StatefulSet, bool, error) {
				t.Fatal("existing bypass should skip workload lookup and suspension")
				return nil, false, nil
			}
			if err := manager.FinalizeDeployment(context.Background(), c, lookup); err != nil {
				t.Fatalf("existing bypass was rejected: %v", err)
			}
		})
	}
}

func TestDeletionErrorCountPolicyAfterRejectedSuspension(t *testing.T) {
	for _, status := range []int{http.StatusUnauthorized, http.StatusForbidden} {
		for _, count := range []string{"", "invalid", "999999999999999999999999", "3", "4", "99"} {
			t.Run(strconv.Itoa(status)+"/count="+count, func(t *testing.T) {
				requests := 0
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					requests++
					if r.Method != http.MethodPut || r.URL.Path != "/suspend" {
						t.Errorf("unexpected suspension request: %s %s", r.Method, r.URL.Path)
					}
					w.WriteHeader(status)
				}))
				defer server.Close()
				host, portText, err := net.SplitHostPort(server.Listener.Addr().String())
				if err != nil {
					t.Fatal(err)
				}
				port, err := strconv.Atoi(portText)
				if err != nil {
					t.Fatal(err)
				}
				c := &coh.Coherence{ObjectMeta: metav1.ObjectMeta{Name: "test", Namespace: "test", Annotations: map[string]string{"coherence.oracle.com/error-count": count}}}
				c.Spec.Replicas = ptr.To(int32(1))
				pod := &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{Name: "test-0", Namespace: "test", Labels: map[string]string{"app": "test"}},
					Spec:       corev1.PodSpec{Containers: []corev1.Container{{Name: coh.ContainerNameCoherence, Ports: []corev1.ContainerPort{{Name: coh.PortNameHealth, ContainerPort: int32(port)}}}}},
					Status:     corev1.PodStatus{Phase: corev1.PodRunning, PodIP: host, Conditions: []corev1.PodCondition{{Type: corev1.PodReady, Status: corev1.ConditionTrue}}},
				}
				sts := &appsv1.StatefulSet{ObjectMeta: metav1.ObjectMeta{Name: "test", Namespace: "test"}, Spec: appsv1.StatefulSetSpec{Replicas: ptr.To(int32(1)), Selector: &metav1.LabelSelector{MatchLabels: map[string]string{"app": "test"}}}, Status: appsv1.StatefulSetStatus{ReadyReplicas: 1, Replicas: 1}}
				manager := &FinalizerManager{Client: fake.NewClientBuilder().WithObjects(pod).Build(), EventRecorder: events.NewFakeRecorder(20)}
				lookup := func(context.Context, string, string) (*appsv1.StatefulSet, bool, error) { return sts, true, nil }
				err = manager.FinalizeDeployment(context.Background(), c, lookup)
				if count == "4" || count == "99" {
					if err != nil {
						t.Fatalf("existing error-count policy should permit deletion: %v", err)
					}
				} else {
					var requestErr *probe.RequestError
					if !errors.As(err, &requestErr) || requestErr.StatusCode != status {
						t.Fatalf("expected preserved HTTP %d failure, got %v", status, err)
					}
				}
				if requests != 1 {
					t.Fatalf("expected one PUT and no fallback, got %d requests", requests)
				}
			})
		}
	}
}
