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
		details.AddArg("-Dloader.main=" + v1.ConsoleMain)
	} else {
		details.AppType = v1.AppTypeJava
		details.MainClass = v1.ConsoleMain
	}

	details.AddArg("-Dcoherence.role=console")
	details.AddArg("-Dcoherence.distributed.localstorage=false")
	details.AddArg("-Dcoherence.localport.adjust=true")
	details.AddArg("-Dcoherence.management.http=none")
	details.AddArg("-Dcoherence.management.http.port=0")
	details.AddArg("-Dcoherence.metrics.http.enabled=false")
	details.AddArg("-Dcoherence.metrics.http.port=0")
	details.AddArg("-Dcoherence.operator.health.enabled=false")
	details.AddArg("-Dcoherence.health.http.port=0")
	details.AddArg("-Dcoherence.grpc.enabled=false")
	details.AddDiagnosticOption("-XX:NativeMemoryTracking=off")
	details.MainArgs = args
}
