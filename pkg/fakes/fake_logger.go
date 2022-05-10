/*
 * Copyright (c) 2019, 2022 Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package fakes

import (
	"github.com/go-logr/logr"
	"testing"
)

// TestLogSink is a logr.LogSink that prints through a testing.T object.
type TestLogSink struct {
	T *testing.T
}

var _ logr.LogSink = TestLogSink{}

func (sink TestLogSink) Init(info logr.RuntimeInfo) {
}

func (sink TestLogSink) Enabled(level int) bool {
	return false
}

func (sink TestLogSink) Info(level int, msg string, _ ...interface{}) {
	sink.T.Logf("%s", msg)
}

func (sink TestLogSink) Error(err error, msg string, keysAndValues ...interface{}) {
	sink.T.Logf("%s: %v -- %v", msg, err, keysAndValues)
}

func (sink TestLogSink) WithValues(keysAndValues ...interface{}) logr.LogSink {
	return sink
}

func (sink TestLogSink) WithName(name string) logr.LogSink {
	return sink
}
