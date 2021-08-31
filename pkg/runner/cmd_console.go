/*
 * Copyright (c) 2021, Oracle and/or its affiliates.
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
	// CommandConsole is the argument to launch a Coherence console.
	CommandConsole = "console"
)

// consoleCommand creates the cobra sub-command to run a Coherence CacheFactory console.
func consoleCommand() *cobra.Command {
	return &cobra.Command{
		Use:   CommandConsole,
		Short: "Start a Coherence interactive console",
		Long:  "Starts a Coherence interactive console",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(cmd, console)
		},
	}
}

// Configure the runner to run a Coherence CacheFactory console
func console(details *RunDetails, _ *cobra.Command) {
	details.Command = CommandConsole
	details.AppType = AppTypeJava
	details.MainClass = "com.tangosol.net.CacheFactory"
	details.addArg("-Dcoherence.distributed.localstorage=false")
	details.setenv(v1.EnvVarCohRole, "console")
	details.unsetenv(v1.EnvVarJvmMemoryHeap)
	if len(os.Args) > 2 {
		details.MainArgs = os.Args[2:]
	}
}
