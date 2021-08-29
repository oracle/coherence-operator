/*
 * Copyright (c) 2020, 2021, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package main

import (
	"github.com/oracle/coherence-operator/pkg/runner"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var (
	log = ctrl.Log.WithName("runner")

	// Version is the runner version injected by the Go linker at build time.
	Version string
	// Commit is the runner Git commit injected by the Go linker at build time.
	Commit string
	// Date is the runner build date injected by the Go linker at build time.
	Date string
)

func main() {
	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))

	log.Info("Coherence Operator runner version information", "Version", Version, "Commit", Commit, "BuildDate", Date)
	if _, err := runner.Execute(); err != nil {
		logf.Log.WithName("runner").Error(err, "Unexpected error while executing command")
		os.Exit(1)
	}
}
