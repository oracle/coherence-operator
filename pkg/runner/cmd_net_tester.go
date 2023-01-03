/*
 * Copyright (c) 2022, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package runner

import (
	"context"
	"github.com/oracle/coherence-operator/pkg/net_testing"
	"github.com/spf13/cobra"
)

const (
	// CommandNetTest is the argument to launch a network test.
	CommandNetTest = "net-test"

	// CommandNetTestServer is the argument to launch a network test server.
	CommandNetTestServer = "server"

	// CommandNetTestOperator is the argument to launch a network test Operator simulation.
	CommandNetTestOperator = "operator"

	// CommandNetTestCluster is the argument to launch a network test Coherence cluster member simulation.
	CommandNetTestCluster = "cluster"

	// ArgOperatorHostName is the argument to use to specify the Operator host name to connect to
	ArgOperatorHostName = "operator-host"

	// ArgClusterHostName is the argument to use to specify the Operator host name to connect to
	ArgClusterHostName = "cluster-host"
)

// networkTestCommand creates the corba net-test sub-command
func networkTestCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   CommandNetTest,
		Short: "Network test",
		Long:  "Run a network communication test",
		Args:  cobra.ArbitraryArgs,
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	rootCmd.AddCommand(networkTestServerCommand())
	rootCmd.AddCommand(networkTestOperatorCommand())
	rootCmd.AddCommand(networkTestClusterCommand())

	return rootCmd
}

// networkTestServerCommand creates the network test server sub-command
func networkTestServerCommand() *cobra.Command {
	return &cobra.Command{
		Use:   CommandNetTestServer,
		Short: "Network test server",
		Long:  "Starts a network communication test server",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return netTestServer(cmd)
		},
	}
}

// netTestOperator runs a network test server
func netTestServer(cmd *cobra.Command) error {
	test := net_testing.NewServerRunner()
	if err := test.Run(context.Background()); err != nil {
		return err
	}
	return nil
}

// networkTestServerCommand creates the network test operator simulator sub-command
func networkTestOperatorCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   CommandNetTestOperator,
		Short: "Network test Operator simulator",
		Long:  "Run a network communication test Operator simulator",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return netTestOperator(cmd)
		},
	}

	flagSet := cmd.Flags()
	flagSet.String(ArgClusterHostName, "127.0.0.1", "The host name of the Coherence cluster simulator")

	return cmd
}

// netTestOperator runs a network test Operator simulator
func netTestOperator(cmd *cobra.Command) error {
	flagSet := cmd.Flags()

	hostName, err := flagSet.GetString(ArgClusterHostName)
	if err != nil {
		return err
	}

	test := net_testing.NewOperatorSimulatorRunner(hostName)
	if err := test.Run(context.Background()); err != nil {
		return err
	}
	return nil
}

// networkTestClusterCommand creates the network test Coherence cluster member simulator sub-command
func networkTestClusterCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   CommandNetTestCluster,
		Short: "Network test Operator simulator",
		Long:  "Run a network communication test Operator simulator",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return netTestCluster(cmd)
		},
	}

	flagSet := cmd.Flags()
	flagSet.String(ArgOperatorHostName, "127.0.0.1", "The host name of the Coherence Operator simulator")
	flagSet.String(ArgClusterHostName, "127.0.0.1", "The host name of the Coherence cluster member simulator")

	return cmd
}

// netTestCluster runs a network test Coherence cluster member simulator
func netTestCluster(cmd *cobra.Command) error {

	flagSet := cmd.Flags()

	operatorHost, err := flagSet.GetString(ArgOperatorHostName)
	if err != nil {
		return err
	}

	clusterHost, err := flagSet.GetString(ArgClusterHostName)
	if err != nil {
		return err
	}

	test := net_testing.NewClusterMemberRunner(operatorHost, clusterHost)
	if err := test.Run(context.Background()); err != nil {
		return err
	}
	return nil
}
