/*
 * Copyright (c) 2020, 2025, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package runner

import (
	v1 "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/pkg/runner/run_details"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	// CommandConsole is the argument to launch a Coherence console.
	CommandConsole = "console"
)

// consoleCommand creates the cobra sub-command to run a Coherence CacheFactory console.
func consoleCommand(v *viper.Viper) *cobra.Command {
	cmd := &cobra.Command{
		Use:   CommandConsole,
		Short: "Start a Coherence interactive console",
		Long:  "Starts a Coherence interactive console",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(cmd, func(details *run_details.RunDetails, cmd *cobra.Command) {
				console(details, args, v)
			})
		},
	}
	addEnvVarFlag(cmd)
	addJvmArgFlag(cmd)
	setupFlags(cmd, v)
	return cmd
}

// Configure the runner to run a Coherence CacheFactory console
func console(details *run_details.RunDetails, args []string, v *viper.Viper) {
	details.Command = CommandConsole
	loadConfigFiles(details)

	if details.IsSpringBoot() {
		details.AddSystemPropertyArg(v1.SysPropSpringLoaderMain, v1.ConsoleMain)
	} else {
		details.AppType = v1.AppTypeJava
		details.MainClass = v1.ConsoleMain
	}

	details.AddSystemPropertyArg(v1.SysPropCoherenceRole, CommandConsole)
	details.AddSystemPropertyArg(v1.SysPropCoherenceDistributedLocalStorage, "false")
	details.AddSystemPropertyArg(v1.SysPropCoherenceLocalPortAdjust, "true")
	details.AddSystemPropertyArg(v1.SysPropCoherenceManagementHttp, "none")
	details.AddSystemPropertyArg(v1.SysPropCoherenceManagementHttpPort, "0")
	details.AddSystemPropertyArg(v1.SysPropCoherenceMetricsHttpEnabled, "false")
	details.AddSystemPropertyArg(v1.SysPropCoherenceMetricsHttpPort, "0")
	details.AddSystemPropertyArg(v1.SysPropOperatorHealthEnabled, "false")
	details.AddSystemPropertyArg(v1.SysPropCoherenceHealthHttpPort, "0")
	details.AddSystemPropertyArg(v1.SysPropCoherenceGrpcEnabled, "false")
	details.AddDiagnosticOption("-XX:NativeMemoryTracking=off")
	details.MainArgs = args
}
