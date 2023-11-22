/*
 * Copyright (c) 2023, Oracle and/or its affiliates.
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
	// CommandSleep is the argument to sleep for a number of seconds.
	CommandSleep = "sleep"
)

// queryPlusCommand creates the corba "sleep" sub-command
func sleepCommand() *cobra.Command {
	return &cobra.Command{
		Use:   CommandSleep,
		Short: "Sleep for a number of seconds",
		Long:  "Sleep for a number of seconds",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(cmd, sleep)
		},
	}
}

func sleep(details *RunDetails, _ *cobra.Command) {
	details.Command = CommandSleep
	details.AppType = AppTypeJava
	details.MainClass = "com.oracle.coherence.k8s.Sleep"
	if len(os.Args) > 2 {
		details.MainArgs = os.Args[2:]
	}
	details.UseOperatorHealth = true
	details.addArg("-Dcoherence.distributed.localstorage=false")
	details.setenv(v1.EnvVarCohRole, "sleep")
	details.unsetenv(v1.EnvVarJvmMemoryHeap)
}
