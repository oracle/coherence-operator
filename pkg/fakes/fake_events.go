/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package fakes

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// FakeEventRecorder is used as a fake during tests. It is thread safe. It is usable
// when created manually and not by NewFakeRecorder, however all events may be
// thrown away in this case.
type FakeEventRecorder struct {
	Events chan FakeEvent
}

type FakeEvent struct {
	Owner       runtime.Object
	Type        string
	Reason      string
	Message     string
	Timestamp   metav1.Time
	Annotations map[string]string
}

func (f *FakeEventRecorder) Event(object runtime.Object, eventtype, reason, message string) {
	if f.Events != nil {
		f.Events <- FakeEvent{Owner: object, Type: eventtype, Reason: reason, Message: message}
	}
}

func (f *FakeEventRecorder) Eventf(object runtime.Object, eventtype, reason, messageFmt string, args ...interface{}) {
	if f.Events != nil {
		f.Events <- FakeEvent{Owner: object, Type: eventtype, Reason: reason, Message: fmt.Sprintf(messageFmt, args...)}
	}
}

func (f *FakeEventRecorder) PastEventf(object runtime.Object, timestamp metav1.Time, eventtype, reason, messageFmt string, args ...interface{}) {
	f.Events <- FakeEvent{Owner: object, Type: eventtype, Reason: reason, Message: fmt.Sprintf(messageFmt, args...), Timestamp: timestamp}
}

func (f *FakeEventRecorder) AnnotatedEventf(object runtime.Object, annotations map[string]string, eventtype, reason, messageFmt string, args ...interface{}) {
	var clone map[string]string
	if annotations != nil {
		clone := make(map[string]string)
		for k, v := range annotations {
			clone[k] = v
		}
	}
	f.Events <- FakeEvent{Owner: object, Type: eventtype, Reason: reason, Message: fmt.Sprintf(messageFmt, args...), Annotations: clone}
}

// NewFakeEventRecorder creates new fake event recorder with event channel with
// buffer of given size.
func NewFakeEventRecorder(bufferSize int) *FakeEventRecorder {
	return &FakeEventRecorder{
		Events: make(chan FakeEvent, bufferSize),
	}
}
