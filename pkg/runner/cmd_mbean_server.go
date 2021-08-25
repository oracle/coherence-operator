/*
 * Copyright (c) 2021, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package runner

import (
	v1 "github.com/oracle/coherence-operator/api/v1"
	"github.com/spf13/cobra"
)

const (
	// CommandMBeanServer is the argument to launch a MBean server.
	CommandMBeanServer = "mbeanserver"
)

// mbeanServerCommand creates the cobra sub-command to run a Coherence MBean server.
func mbeanServerCommand() *cobra.Command {
	return &cobra.Command{
		Use:   CommandMBeanServer,
		Short: "Start a Coherence MBean server",
		Long:  "Starts a Coherence MBean server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(cmd, mbeanServer)
		},
	}
}

// Configure the runner to run a JMXMP MBean server
func mbeanServer(details *RunDetails, _ *cobra.Command) {
	details.Command = CommandMBeanServer
	details.AppType = AppTypeJava
	details.addClasspath(details.UtilsDir + "/lib/*")
	details.MainClass = "com.oracle.coherence.k8s.JmxmpServer"
	details.MainArgs = []string{}
	details.setenv(v1.EnvVarJvmJmxmpEnabled, "true")
	details.setenv(v1.EnvVarCohRole, "MBeanServer")
	details.addArg("-Dcoherence.distributed.localstorage=false")
	details.addArg("-Dcoherence.management=all")
	details.addArg("-Dcoherence.management.remote=true")
	details.addArg("-Dcom.sun.management.jmxremote.ssl=false")
	details.addArg("-Dcom.sun.management.jmxremote.authenticate=false")
}
