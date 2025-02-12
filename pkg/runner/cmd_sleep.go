/*
 * Copyright (c) 2020, 2025, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package runner

import (
	v1 "github.com/oracle/coherence-operator/api/v1"
	"github.com/spf13/cobra"
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
			return run(cmd, func(details *RunDetails, _ *cobra.Command) {
				sleep(details, args)
			})
		},
	}
}

func sleep(details *RunDetails, args []string) {
	details.Command = CommandSleep
	details.AppType = AppTypeJava
	details.MainClass = "com.oracle.coherence.k8s.Sleep"
	details.MainArgs = args
	details.UseOperatorHealth = true
	details.addArg("-Dcoherence.distributed.localstorage=false")
	details.setenv(v1.EnvVarCohRole, "sleep")
	details.unsetenv(v1.EnvVarJvmMemoryHeap)
}
