/*
 * Copyright (c) 2026, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package events

import (
	"fmt"
	"testing"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type recordedEvent struct {
	regarding runtime.Object
	related   runtime.Object
	eventType string
	reason    string
	action    string
	note      string
}

type actionRecorder struct {
	event recordedEvent
}

func (in *actionRecorder) Eventf(regarding runtime.Object, related runtime.Object, eventType, reason, action, note string, args ...interface{}) {
	in.event = recordedEvent{
		regarding: regarding,
		related:   related,
		eventType: eventType,
		reason:    reason,
		action:    action,
		note:      fmt.Sprintf(note, args...),
	}
}

func TestOwnedEventRecorderUsesNonEmptyAction(t *testing.T) {
	owner := &corev1.Pod{}

	tests := []struct {
		name      string
		record    func(*OwnedEventRecorder)
		eventType string
		reason    string
		note      string
	}{
		{
			name: "event",
			record: func(recorder *OwnedEventRecorder) {
				recorder.Event(corev1.EventTypeNormal, "Created", "created resource")
			},
			eventType: corev1.EventTypeNormal,
			reason:    "Created",
			note:      "created resource",
		},
		{
			name: "event with literal percent",
			record: func(recorder *OwnedEventRecorder) {
				recorder.Event(corev1.EventTypeNormal, "Progressing", "progress 100%")
			},
			eventType: corev1.EventTypeNormal,
			reason:    "Progressing",
			note:      "progress 100%",
		},
		{
			name: "formatted event",
			record: func(recorder *OwnedEventRecorder) {
				recorder.Eventf(corev1.EventTypeWarning, "Failed", "failed resource %s", "test")
			},
			eventType: corev1.EventTypeWarning,
			reason:    "Failed",
			note:      "failed resource test",
		},
		{
			name: "formatted event with percent verb",
			record: func(recorder *OwnedEventRecorder) {
				recorder.Eventf(corev1.EventTypeNormal, "Progressing", "progress %d%%", 100)
			},
			eventType: corev1.EventTypeNormal,
			reason:    "Progressing",
			note:      "progress 100%",
		},
		{
			name: "formatted event with percent in argument",
			record: func(recorder *OwnedEventRecorder) {
				recorder.Eventf(corev1.EventTypeWarning, "Failed", "failed: %s", "disk 90%")
			},
			eventType: corev1.EventTypeWarning,
			reason:    "Failed",
			note:      "failed: disk 90%",
		},
		{
			name: "formatted info",
			record: func(recorder *OwnedEventRecorder) {
				recorder.Infof("ServiceSuspended", "suspended StatefulSet %s", "test")
			},
			eventType: corev1.EventTypeNormal,
			reason:    "ServiceSuspended",
			note:      "suspended StatefulSet test",
		},
		{
			name: "formatted warning",
			record: func(recorder *OwnedEventRecorder) {
				recorder.Warnf("ServiceSuspendFailed", "failed StatefulSet %s", "test")
			},
			eventType: corev1.EventTypeWarning,
			reason:    "ServiceSuspendFailed",
			note:      "failed StatefulSet test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target := &actionRecorder{}
			recorder := NewOwnedEventRecorder(owner, target)

			tt.record(&recorder)

			if target.event.regarding != owner {
				t.Fatalf("regarding = %v, want owner %v", target.event.regarding, owner)
			}
			if target.event.related != nil {
				t.Fatalf("related = %v, want nil", target.event.related)
			}
			if target.event.eventType != tt.eventType {
				t.Errorf("event type = %q, want %q", target.event.eventType, tt.eventType)
			}
			if target.event.reason != tt.reason {
				t.Errorf("reason = %q, want %q", target.event.reason, tt.reason)
			}
			if target.event.action != "Operator" {
				t.Errorf("action = %q, want %q", target.event.action, "Operator")
			}
			if target.event.note != tt.note {
				t.Errorf("note = %q, want %q", target.event.note, tt.note)
			}
		})
	}
}
