/*
 * Copyright (c) 2020, 2025, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package runner

import (
	"fmt"
	v1 "github.com/oracle/coherence-operator/api/v1"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os"
)

// configCommand creates the corba "config" sub-command
func configCommand(env map[string]string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   v1.RunnerConfig,
		Short: "Create the Operator JVM args files a Coherence server",
		Long:  "Create the Operator JVM args files a Coherence server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return maybeRun(cmd, createArgsFile)
		},
	}

	utilDir, found := env[v1.EnvVarCohUtilDir]
	if !found || utilDir == "" {
		utilDir = v1.VolumeMountPathUtils
	}

	flagSet := cmd.Flags()
	flagSet.String(ArgUtilsDir, utilDir, "The utils files root directory")

	return cmd
}

// config will create the JVM args files for a Coherence Pod - typically this is run from an init-container
func createArgsFile(details *RunDetails, _ *cobra.Command) (bool, error) {
	var err error

	populateServerDetails(details)
	err = configureCommand(details)
	if err != nil {
		return false, errors.Wrap(err, "failed to configure server command")
	}

	cpFile := fmt.Sprintf("%s%c%s", details.UtilsDir, os.PathSeparator, v1.OperatorClasspathFile)
	classpath := details.getClasspath()
	err = os.WriteFile(cpFile, []byte(classpath), os.ModePerm)
	fmt.Printf("Created class path file %s\n", cpFile)
	fmt.Println("--------------------")
	fmt.Println(classpath)
	fmt.Println("--------------------")
	if err != nil {
		return false, errors.Wrap(err, "failed to write coherence classpath file")
	}

	args := details.Args

	argFileName := fmt.Sprintf("%s%c%s", details.UtilsDir, os.PathSeparator, v1.OperatorJvmArgsFile)
	argFile, err := os.Create(argFileName)
	if err != nil {
		return false, errors.Wrap(err, "failed to create coherence jvm args file")
	}
	defer argFile.Close()
	for _, arg := range args {
		_, err = argFile.WriteString(arg + "\n")
		if err != nil {
			return false, errors.Wrap(err, "failed to write coherence jvm args file")
		}
	}

	err = os.Chmod(argFileName, os.ModePerm)
	if err != nil {
		return false, errors.Wrap(err, "failed to set file-mode on coherence jvm args file")
	}
	fmt.Printf("Created JVM args file %s\n", argFileName)
	fmt.Println("--------------------")
	for _, arg := range args {
		fmt.Println(arg)
	}
	fmt.Println("--------------------")

	return false, err
}
