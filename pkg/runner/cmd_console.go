/*
 * Copyright (c) 2020, 2025, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package runner

import (
	v1 "github.com/oracle/coherence-operator/api/v1"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
)

const (
	// CommandConsole is the argument to launch a Coherence console.
	CommandConsole = "console"
)

// consoleCommand creates the cobra sub-command to run a Coherence CacheFactory console.
func consoleCommand(v *viper.Viper) *cobra.Command {
	return &cobra.Command{
		Use:   CommandConsole,
		Short: "Start a Coherence interactive console",
		Long:  "Starts a Coherence interactive console",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(cmd, func(details *RunDetails, _ *cobra.Command) {
				console(details, args, v)
			})
		},
	}
}

// Configure the runner to run a Coherence CacheFactory console
func console(details *RunDetails, args []string, v *viper.Viper) {
	app := strings.ToLower(v.GetString(v1.EnvVarAppType))
	if app == AppTypeSpring {
		details.AppType = AppTypeSpring
		details.MainClass = SpringBootMain
		details.addArg("-Dloader.main=" + ConsoleMain)
	} else {
		details.AppType = AppTypeJava
		details.MainClass = ConsoleMain
	}
	details.Command = CommandConsole
	details.addArg("-Dcoherence.distributed.localstorage=false")
	details.setenv(v1.EnvVarCohRole, "console")
	details.unsetenv(v1.EnvVarJvmMemoryHeap)
	details.unsetenv(v1.EnvVarCoherenceLocalPortAdjust)
	details.unsetenv(v1.EnvVarCohMgmtPrefix + v1.EnvVarCohEnabledSuffix)
	details.unsetenv(v1.EnvVarCohMetricsPrefix + v1.EnvVarCohEnabledSuffix)
	details.MainArgs = args
}
