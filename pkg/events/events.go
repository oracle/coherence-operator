/*
 * Copyright (c) 2020, 2025, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package events

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
)

// OwnedEventRecorder is a wrapper around an EventRecorder that
// keeps a reference to the object associated with the events
type OwnedEventRecorder struct {
	owner    runtime.Object
	recorder record.EventRecorder
}

// NewOwnedEventRecorder creates a new OwnedEventRecorder
func NewOwnedEventRecorder(owner runtime.Object, recorder record.EventRecorder) OwnedEventRecorder {
	return OwnedEventRecorder{
		owner:    owner,
		recorder: recorder,
	}
}

// Event constructs an event from the given information and puts it in the queue for sending.
// 'eventType' of this event, and can be one of Normal, Warning. New types could be added in future
// 'reason' is the reason this event is generated. 'reason' should be short and unique; it
// should be in UpperCamelCase format (starting with a capital letter). "reason" will be used
// to automate handling of events, so imagine people writing switch statements to handle them.
// You want to make that easy.
// 'message' is intended to be human-readable.
//
// The resulting event will be created in the same namespace as the reference object.
func (in *OwnedEventRecorder) Event(eventType, reason, message string) {
	if in != nil && in.owner != nil && in.recorder != nil {
		in.recorder.Event(in.owner, eventType, reason, message)
	}
}

// Warn is just like Event, and sends a Warning event.
func (in *OwnedEventRecorder) Warn(reason, message string) {
	in.Event(corev1.EventTypeWarning, reason, message)
}

// Info is just like Event, and sends a Normal event.
func (in *OwnedEventRecorder) Info(reason, message string) {
	in.Event(corev1.EventTypeNormal, reason, message)
}

// Eventf is just like Event, but with Sprintf for the message field.
func (in *OwnedEventRecorder) Eventf(eventType, reason, messageFmt string, args ...interface{}) {
	if in != nil && in.owner != nil && in.recorder != nil {
		in.recorder.Eventf(in.owner, eventType, reason, messageFmt, args...)
	}
}

// Warnf is just like Eventf, and sends a Warning event.
func (in *OwnedEventRecorder) Warnf(reason, messageFmt string, args ...interface{}) {
	in.Eventf(corev1.EventTypeWarning, reason, messageFmt, args...)
}

// Infof is just like Eventf, and sends a Normal event.
func (in *OwnedEventRecorder) Infof(reason, messageFmt string, args ...interface{}) {
	in.Eventf(corev1.EventTypeNormal, reason, messageFmt, args...)
}

// AnnotatedEventf is just like eventf, but with annotations attached
func (in *OwnedEventRecorder) AnnotatedEventf(annotations map[string]string, eventType, reason, messageFmt string, args ...interface{}) {
	if in != nil && in.owner != nil && in.recorder != nil {
		in.recorder.AnnotatedEventf(in.owner, annotations, eventType, reason, messageFmt, args...)
	}
}

// AnnotatedWarnf is just like AnnotatedEventf, and sends a Warning event.
func (in *OwnedEventRecorder) AnnotatedWarnf(annotations map[string]string, reason, messageFmt string, args ...interface{}) {
	in.AnnotatedEventf(annotations, corev1.EventTypeWarning, reason, messageFmt, args...)
}

// AnnotatedInfof is just like AnnotatedEventf, and sends a Normal event.
func (in *OwnedEventRecorder) AnnotatedInfof(annotations map[string]string, reason, messageFmt string, args ...interface{}) {
	in.AnnotatedEventf(annotations, corev1.EventTypeNormal, reason, messageFmt, args...)
}
