/*
 * Copyright (c) 2026, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */
package errorhandling

import (
	"context"
	"fmt"
	"testing"
	"time"

	coh "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/pkg/probe"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/events"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestSuspensionRecoveryPreservesDeletionPolicy(t *testing.T) {
	failures := map[string]error{
		"legacy":       fmt.Errorf("failed to suspend services"),
		"unauthorized": &probe.SuspensionError{Cause: &probe.RequestError{StatusCode: 401}},
		"forbidden":    &probe.SuspensionError{Cause: &probe.RequestError{StatusCode: 403}},
		"tls":          &probe.SuspensionError{Cause: &probe.RequestError{Kind: "transport or TLS"}},
		"secret":       &probe.SuspensionError{Cause: &probe.RequestError{Kind: "Secret unavailable"}},
	}
	for name, failure := range failures {
		for _, deleting := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/deleting=%t", name, deleting), func(t *testing.T) {
				scheme := runtime.NewScheme()
				if err := coh.AddToScheme(scheme); err != nil {
					t.Fatal(err)
				}
				c := &coh.Coherence{ObjectMeta: metav1.ObjectMeta{Name: "test", Namespace: "test", Finalizers: []string{coh.CoherenceFinalizer}}}
				if deleting {
					now := metav1.Now()
					c.DeletionTimestamp = &now
				}
				kube := fake.NewClientBuilder().WithScheme(scheme).WithStatusSubresource(c).WithObjects(c).Build()
				handler := &ErrorHandler{Client: kube, EventRecorder: events.NewFakeRecorder(20)}
				result, err := handler.HandleError(context.Background(), fmt.Errorf("wrapped: %w", failure), c, "suspension failed")
				if err != nil || result.RequeueAfter != 10*time.Second {
					t.Fatalf("expected existing recovery retry, got %+v, %v", result, err)
				}
				actual := &coh.Coherence{}
				if err := kube.Get(context.Background(), client.ObjectKeyFromObject(c), actual); err != nil {
					t.Fatal(err)
				}
				_, bypass := actual.Annotations["coherence.oracle.com/finalizer-bypass"]
				if bypass != deleting {
					t.Fatalf("deletion bypass=%t, deleting=%t", bypass, deleting)
				}
				if actual.Annotations[AnnotationErrorCount] != "1" {
					t.Fatalf("existing error tracking was skipped: %v", actual.Annotations)
				}
				if len(actual.Finalizers) != 1 {
					t.Fatal("recovery should annotate the resource for finalization, not remove its finalizer directly")
				}
			})
		}
	}
}
