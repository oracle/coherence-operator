/*
 * Copyright (c) 2020, 2025, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package runner

import (
	v1 "github.com/oracle/coherence-operator/api/v1"
	"github.com/spf13/cobra"
	"os"
)

const (
	// CommandJShell is the argument to launch a JShell console.
	CommandJShell = "jshell"
)

// queryPlusCommand creates the corba "jshell" sub-command
func jShellCommand() *cobra.Command {
	return &cobra.Command{
		Use:   CommandJShell,
		Short: "Start a Coherence interactive JShell console",
		Long:  "Starts a Coherence interactive JShell console",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(cmd, jShell)
		},
	}
}

// Configure the runner to run a Coherence JShell console
func jShell(details *RunDetails, _ *cobra.Command) {
	details.Command = CommandQueryPlus
	details.AppType = AppTypeJava
	details.MainClass = "jdk.internal.jshell.tool.JShellToolProvider"
	if len(os.Args) > 2 {
		details.MainArgs = os.Args[2:]
	}
	details.addArg("-Dcoherence.distributed.localstorage=false")
	details.setenv(v1.EnvVarCohRole, "jshell")
	details.unsetenv(v1.EnvVarJvmMemoryHeap)
	details.unsetenv(v1.EnvVarCoherenceLocalPortAdjust)
	details.unsetenv(v1.EnvVarCohMgmtPrefix + v1.EnvVarCohEnabledSuffix)
	details.unsetenv(v1.EnvVarCohMetricsPrefix + v1.EnvVarCohEnabledSuffix)
}
