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
	// CommandQueryPlus is the argument to launch a QueryPlus console.
	CommandQueryPlus = "queryplus"
)

// queryPlusCommand creates the corba "queryplus" sub-command
func queryPlusCommand() *cobra.Command {
	return &cobra.Command{
		Use:   CommandQueryPlus,
		Short: "Start a Coherence interactive QueryPlus console",
		Long:  "Starts a Coherence interactive QueryPlus console",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(cmd, queryPlus)
		},
	}
}

// Configure the runner to run a Coherence Query Plus console
func queryPlus(details *RunDetails, _ *cobra.Command) {
	details.Command = CommandQueryPlus
	details.AppType = AppTypeJava
	details.MainClass = "com.tangosol.coherence.dslquery.QueryPlus"
	if len(os.Args) > 2 {
		details.MainArgs = os.Args[2:]
	}
	details.addArg("-Dcoherence.distributed.localstorage=false")
	details.setenv(v1.EnvVarCohRole, "queryPlus")
	details.unsetenv(v1.EnvVarJvmMemoryHeap)
}
