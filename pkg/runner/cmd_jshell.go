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
			return run(cmd, func(details *RunDetails, cmd *cobra.Command) {
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
func jShell(details *RunDetails, args []string, v *viper.Viper) {
	details.AppType = v1.AppTypeJShell
	details.Command = CommandJShell
	loadConfigFiles(details)

	details.addArg("-Dcoherence.role=jshell")
	details.addArg("-Dcoherence.distributed.localstorage=false")
	details.addArg("-Dcoherence.localport.adjust=true")
	details.addArg("-Dcoherence.management.http=none")
	details.addArg("-Dcoherence.management.http.port=0")
	details.addArg("-Dcoherence.metrics.http.enabled=false")
	details.addArg("-Dcoherence.metrics.http.port=0")
	details.addArg("-Dcoherence.operator.health.enabled=false")
	details.addArg("-Dcoherence.health.http.port=0")
	details.addArg("-Dcoherence.grpc.enabled=false")
	details.addArg("-XX:NativeMemoryTracking=summary")
	details.MainArgs = args
}
