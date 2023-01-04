/*
 * Copyright (c) 2019, 2023, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package runner

import (
	"context"
	"github.com/oracle/coherence-operator/pkg/nettesting"
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

	// CommandNetTestHook is the argument to launch a network test that connects to the Operator web-hook port.
	CommandNetTestHook = "hook"

	// CommandNetTestClient is the argument to launch a simple network test client.
	CommandNetTestClient = "client"

	// ArgOperatorHostName is the argument to use to specify the Operator host name to connect to
	ArgOperatorHostName = "operator-host"

	// ArgClusterHostName is the argument to use to specify the Operator host name to connect to
	ArgClusterHostName = "cluster-host"

	// ArgHostName is the argument to use to specify the host name to connect to
	ArgHostName = "host"

	// ArgPort is the argument to use to specify the port to connect to
	ArgPort = "port"

	// ArgProtocol is the argument to use to specify the network protocol to use
	ArgProtocol = "protocol"
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
	rootCmd.AddCommand(networkWebHookClientCommand())
	rootCmd.AddCommand(networkSimpleClientCommand())

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
func netTestServer(_ *cobra.Command) error {
	test := nettesting.NewServerRunner()
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

	test := nettesting.NewOperatorSimulatorRunner(hostName)
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

	test := nettesting.NewClusterMemberRunner(operatorHost, clusterHost)
	if err := test.Run(context.Background()); err != nil {
		return err
	}
	return nil
}

// networkTestClusterCommand creates the network test Coherence cluster member simulator sub-command
func networkWebHookClientCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   CommandNetTestHook,
		Short: "Network test web-hook client",
		Long:  "Run a network communication test web-hook client",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return netTestWebHook(cmd)
		},
	}

	flagSet := cmd.Flags()
	flagSet.String(ArgOperatorHostName, "127.0.0.1", "The host name of the Coherence Operator simulator")

	return cmd
}

// netTestWebHook runs a web-hook connectivity test
func netTestWebHook(cmd *cobra.Command) error {

	flagSet := cmd.Flags()

	operatorHost, err := flagSet.GetString(ArgOperatorHostName)
	if err != nil {
		return err
	}

	test := nettesting.NewWebHookClientRunner(operatorHost)
	if err := test.Run(context.Background()); err != nil {
		return err
	}
	return nil
}

// networkTestClusterCommand creates the network test Coherence cluster member simulator sub-command
func networkSimpleClientCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   CommandNetTestClient,
		Short: "Network test client",
		Long:  "Run a network communication test client for a single host and port",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return netTestClient(cmd)
		},
	}

	flagSet := cmd.Flags()
	flagSet.String(ArgHostName, "127.0.0.1", "The host name of the server to connect to")
	flagSet.Int(ArgPort, -1, "The port to connect to")
	flagSet.String(ArgProtocol, "tcp", "The host name of the server to connect to")

	return cmd
}

// netTestClient runs a connectivity test client
func netTestClient(cmd *cobra.Command) error {

	flagSet := cmd.Flags()

	host, err := flagSet.GetString(ArgHostName)
	if err != nil {
		return err
	}

	port, err := flagSet.GetInt(ArgPort)
	if err != nil {
		return err
	}

	protocol, err := flagSet.GetString(ArgProtocol)
	if err != nil {
		return err
	}

	test := nettesting.NewSimpleClientRunner(host, port, protocol)
	if err := test.Run(context.Background()); err != nil {
		return err
	}
	return nil
}
