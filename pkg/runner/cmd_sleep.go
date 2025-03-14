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
			return run(cmd, func(details *RunDetails, cmd *cobra.Command) {
				sleep(details, args, v)
			})
		},
	}
	addEnvVarFlag(cmd)
	addJvmArgFlag(cmd)
	setupFlags(cmd, v)
	return cmd
}

func sleep(details *RunDetails, args []string, v *viper.Viper) {
	app := strings.ToLower(v.GetString(v1.EnvVarAppType))
	if app == v1.AppTypeSpring2 {
		details.AppType = v1.AppTypeSpring2
		details.MainClass = v1.SpringBootMain2
		details.addArg("-Dloader.main=com.oracle.coherence.k8s.Sleep")
	} else {
		details.AppType = v1.AppTypeJava
		details.MainClass = "com.oracle.coherence.k8s.Sleep"
	}
	details.Command = CommandSleep
	details.MainArgs = args
	details.UseOperatorHealth = true
	details.addArg("-Dcoherence.distributed.localstorage=false")
	details.addArg("-Dcoherence.localport.adjust=true")
	details.addArg("-Dcoherence.management.http=none")
	details.addArg("-Dcoherence.management.http.port=0")
	details.addArg("-Dcoherence.metrics.http.enabled=false")
	details.addArg("-Dcoherence.metrics.http.port=0")
	details.addArg("-Dcoherence.operator.health.enabled=false")
	details.addArg("-Dcoherence.grpc.enabled=false")
	details.setenv(v1.EnvVarJvmMemoryNativeTracking, "off")
	details.setenv(v1.EnvVarCohRole, "sleep")
	details.setenv(v1.EnvVarCohHealthPort, "0")
	details.MainArgs = args
}
