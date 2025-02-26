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
	// CommandSleep is the argument to sleep for a number of seconds.
	CommandSleep = "sleep"
)

// queryPlusCommand creates the corba "sleep" sub-command
func sleepCommand(v *viper.Viper) *cobra.Command {
	cmd := &cobra.Command{
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
	addEnvVarFlag(cmd)
	addJvmArgFlag(cmd)
	setupFlags(cmd, v)
	return cmd
}

func sleep(details *RunDetails, args []string) {
	details.Command = CommandSleep
	details.AppType = AppTypeJava
	details.MainClass = "com.oracle.coherence.k8s.Sleep"
	details.MainArgs = args
	details.UseOperatorHealth = true
	details.addArg("-Dcoherence.distributed.localstorage=false")
	details.addArg("-Dcoherence.localport.adjust=true")
	details.addArg("-Dcoherence.management.http.enabled=false")
	details.addArg("-Dcoherence.metrics.http.enabled=false")
	details.setenv(v1.EnvVarCohRole, "sleep")
}
