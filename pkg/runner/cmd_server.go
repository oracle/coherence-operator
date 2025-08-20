/*
 * Copyright (c) 2020, 2025, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package runner

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	v1 "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/pkg/runner/run_details"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	// CommandServer is the argument to launch a server.
	CommandServer = "server"
)

// serverCommand creates the corba "server" sub-command
func serverCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   CommandServer,
		Short: "Start a Coherence server",
		Long:  "Starts a Coherence server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(cmd, server)
		},
	}

	return cmd
}

// Configure the runner to run a Coherence Server
func server(details *run_details.RunDetails, _ *cobra.Command) {
	details.Command = CommandServer
	loadConfigFiles(details)

	ma, found := details.LookupEnv(v1.EnvVarAppMainArgs)
	if found {
		if ma != "" {
			for _, arg := range strings.Split(ma, " ") {
				details.MainArgs = append(details.MainArgs, details.ExpandEnv(arg))
			}
		}
	}
}

func loadConfigFiles(details *run_details.RunDetails) {
	var err error

	if err = loadClassPathFile(details); err != nil {
		fmt.Printf("Error loading class path file %v\n", err)
	}

	if err = loadJvmArgsFile(details); err != nil {
		fmt.Printf("Error loading class path file %v\n", err)
	}

	if details.IsSpringBoot() {
		if err = loadSpringBootArgsFile(details); err != nil {
			fmt.Printf("Error loading main class file %v\n", err)
		}
	}

	if err = loadMainClassFile(details); err != nil {
		fmt.Printf("Error loading main class file %v\n", err)
	}
}

func loadClassPathFile(details *run_details.RunDetails) error {
	file := fmt.Sprintf(v1.FileNamePattern, details.UtilsDir, os.PathSeparator, v1.OperatorClasspathFile)
	data, err := os.ReadFile(file)
	if err != nil {
		return errors.Wrapf(err, "error reading %s", file)
	}
	details.Classpath = string(data)
	return nil
}

func loadJvmArgsFile(details *run_details.RunDetails) error {
	file := fmt.Sprintf(v1.FileNamePattern, details.UtilsDir, os.PathSeparator, v1.OperatorJvmArgsFile)
	lines, err := readLines(file)
	if err != nil {
		return errors.Wrapf(err, "error reading %s", file)
	}
	details.AddArgs(lines...)
	return nil
}

func loadMainClassFile(details *run_details.RunDetails) error {
	dir := details.GetenvOrDefault(v1.EnvVarCohUtilDir, details.UtilsDir)
	file := fmt.Sprintf(v1.FileNamePattern, dir, os.PathSeparator, v1.OperatorMainClassFile)
	lines, err := readLines(file)
	if err != nil {
		return errors.Wrapf(err, "error reading %s", file)
	}
	if len(lines) > 0 {
		details.MainClass = lines[0]
		if len(lines) > 1 {
			details.MainArgs = append(details.MainArgs, lines[1:]...)
		}
	}
	return nil
}

func loadSpringBootArgsFile(details *run_details.RunDetails) error {
	file := fmt.Sprintf(v1.FileNamePattern, details.UtilsDir, os.PathSeparator, v1.OperatorSpringBootArgsFile)
	lines, err := readLines(file)
	if err != nil {
		return errors.Wrapf(err, "error reading %s", file)
	}
	details.AddArgs(lines...)
	return nil
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
