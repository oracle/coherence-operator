/*
 * Copyright (c) 2020, 2022, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package main

import (
	"fmt"
	"github.com/oracle/coherence-operator/pkg/operator"
	"github.com/oracle/coherence-operator/pkg/runner"
	"os"
	"runtime"
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
	// Branch is the runner Git branch injected by the Go linker at build time.
	Branch string
	// Date is the runner build date injected by the Go linker at build time.
	Date string
	// Author is the username of the account at build time
	Author string
)

func main() {
	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))

	printVersion()
	operator.SetVersion(Version)

	if _, err := runner.Execute(); err != nil {
		logf.Log.WithName("runner").Error(err, "Unexpected error while executing command")
		os.Exit(1)
	}
}

func printVersion() {
	log.Info(fmt.Sprintf("Operator Version: %s", Version))
	log.Info(fmt.Sprintf("Operator Build Date: %s", Date))
	log.Info(fmt.Sprintf("Operator Built By: %s", Author))
	log.Info(fmt.Sprintf("Operator Git Commit: %s (%s)", Commit, Branch))
	log.Info(fmt.Sprintf("Go Version: %s", runtime.Version()))
	log.Info(fmt.Sprintf("Go OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH))
}
