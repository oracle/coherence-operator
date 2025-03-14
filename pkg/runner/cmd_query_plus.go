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
			return run(cmd, func(details *run_details.RunDetails, cmd *cobra.Command) {
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
func queryPlus(details *run_details.RunDetails, args []string, v *viper.Viper) {
	details.Command = CommandQueryPlus
	loadConfigFiles(details)

	if details.IsSpringBoot() {
		details.AddArg("-Dloader.main=" + v1.QueryPlusMain)
	} else {
		details.AppType = v1.AppTypeJava
		details.MainClass = v1.QueryPlusMain
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
