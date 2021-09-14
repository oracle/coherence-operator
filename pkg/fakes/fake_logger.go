/*
 * Copyright (c) 2019, 2021 Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package fakes

import (
	"github.com/go-logr/logr"
	"testing"
)

// TestLogger is a logr.Logger that prints through a testing.T object.
type TestLogger struct {
	T *testing.T
}

var _ logr.Logger = TestLogger{}

func (log TestLogger) Info(msg string, keysAndValues ...interface{}) {
	log.T.Logf("%s %v", msg, keysAndValues)
}

func (TestLogger) Enabled() bool {
	return false
}

func (log TestLogger) Error(err error, msg string, args ...interface{}) {
	log.T.Logf("%s: %v -- %v", msg, err, args)
}

func (log TestLogger) V(v int) logr.Logger {
	return log
}

func (log TestLogger) WithName(_ string) logr.Logger {
	return log
}

func (log TestLogger) WithValues(_ ...interface{}) logr.Logger {
	return log
}
