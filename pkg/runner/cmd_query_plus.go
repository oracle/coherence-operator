/*
 * Copyright (c) 2020, 2025, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package runner

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	// CommandQueryPlus is the argument to launch a QueryPlus console.
	CommandQueryPlus = "queryplus"
)

// queryPlusCommand creates the corba "queryplus" sub-command
func queryPlusCommand(v *viper.Viper) *cobra.Command {
	cmd := &cobra.Command{
		Use:   CommandQueryPlus,
		Short: "Start a Coherence interactive QueryPlus console",
		Long:  "Starts a Coherence interactive QueryPlus console",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(cmd, func(details *RunDetails, cmd *cobra.Command) {
				queryPlus(details, args, v)
			})
		},
	}
	addEnvVarFlag(cmd)
	addJvmArgFlag(cmd)
	setupFlags(cmd, v)
	return cmd
}

// Configure the runner to run a Coherence Query Plus console
func queryPlus(details *RunDetails, args []string, v *viper.Viper) {
	details.Command = CommandQueryPlus
	loadConfigFiles(details)

	if details.IsSpringBoot() {
		details.addArg("-Dloader.main=" + QueryPlusMain)
	} else {
		details.AppType = AppTypeJava
		details.MainClass = QueryPlusMain
	}

	details.addArg("-Dcoherence.role=console")
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
