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
	// CommandJShell is the argument to launch a JShell console.
	CommandJShell = "jshell"
)

// queryPlusCommand creates the corba "jshell" sub-command
func jShellCommand(v *viper.Viper) *cobra.Command {
	cmd := &cobra.Command{
		Use:   CommandJShell,
		Short: "Start a Coherence interactive JShell console",
		Long:  "Starts a Coherence interactive JShell console",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(cmd, func(details *run_details.RunDetails, cmd *cobra.Command) {
				jShell(details, args, v)
			})
		},
	}
	addEnvVarFlag(cmd)
	addJvmArgFlag(cmd)
	setupFlags(cmd, v)
	return cmd
}

// Configure the runner to run a Coherence JShell console
func jShell(details *run_details.RunDetails, args []string, v *viper.Viper) {
	details.AppType = v1.AppTypeJShell
	details.Command = CommandJShell
	loadConfigFiles(details)

	details.AddSystemPropertyArg(v1.SysPropCoherenceRole, CommandJShell)
	details.AddSystemPropertyArg(v1.SysPropCoherenceDistributedLocalStorage, "false")
	details.AddSystemPropertyArg(v1.SysPropCoherenceLocalPortAdjust, "true")
	details.AddSystemPropertyArg(v1.SysPropCoherenceManagementHttp, "none")
	details.AddSystemPropertyArg(v1.SysPropCoherenceManagementHttpPort, "0")
	details.AddSystemPropertyArg(v1.SysPropCoherenceMetricsHttpEnabled, "false")
	details.AddSystemPropertyArg(v1.SysPropCoherenceMetricsHttpPort, "0")
	details.AddSystemPropertyArg(v1.SysPropOperatorHealthEnabled, "false")
	details.AddSystemPropertyArg(v1.SysPropCoherenceHealthHttpPort, "0")
	details.AddSystemPropertyArg(v1.SysPropCoherenceGrpcEnabled, "false")
	details.AddArg("-XX:NativeMemoryTracking=off")
	details.AddArg("-XshowSettings:none")
	details.AddArg("-XX:-PrintCommandLineFlags")
	details.AddArg("-XX:-PrintFlagsFinal")
	details.MainArgs = args
}
