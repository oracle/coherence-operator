/*
 * Copyright (c) 2020, 2025, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package runner

import (
	"bytes"
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
			return maybeRun(cmd, createsFiles)
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

// createsFiles will create the various config files
func createsFiles(details *RunDetails, _ *cobra.Command) (bool, error) {
	populateServerDetails(details)
	if err := createClassPathFile(details); err != nil {
		return false, err
	}
	if err := createArgsFile(details); err != nil {
		return false, err
	}
	if err := createCliConfig(details); err != nil {
		return false, err
	}
	return false, nil
}

// createClassPathFile will create the class path files for a Coherence Pod - typically this is run from an init-container
func createClassPathFile(details *RunDetails) error {
	err := configureCommand(details)
	if err != nil {
		return errors.Wrap(err, "failed to configure server command")
	}

	cpFile := fmt.Sprintf("%s%c%s", details.UtilsDir, os.PathSeparator, v1.OperatorClasspathFile)
	classpath := details.getClasspath()
	err = os.WriteFile(cpFile, []byte(classpath), os.ModePerm)
	fmt.Printf("Created class path file %s\n", cpFile)
	fmt.Println("--------------------")
	fmt.Println(classpath)
	fmt.Println("--------------------")
	if err != nil {
		return errors.Wrap(err, "failed to write coherence classpath file")
	}
	return nil
}

// createArgsFile will create the JVM args files for a Coherence Pod - typically this is run from an init-container
func createArgsFile(details *RunDetails) error {
	var err error

	args := details.Args
	argFileName := fmt.Sprintf("%s%c%s", details.UtilsDir, os.PathSeparator, v1.OperatorJvmArgsFile)

	var buffer bytes.Buffer
	for _, arg := range args {
		buffer.WriteString(arg + "\n")
	}
	if err := os.WriteFile(argFileName, buffer.Bytes(), os.ModePerm); err != nil {
		return errors.Wrap(err, "failed to write JVM args file "+argFileName)
	}

	fmt.Printf("Created JVM args file %s\n", argFileName)
	fmt.Println("--------------------")
	fmt.Println(buffer.String())
	fmt.Println("--------------------")

	return err
}

func createCliConfig(details *RunDetails) error {
	home := details.getenvOrDefault(v1.EnvVarCohCtlHome, details.UtilsDir)
	fileName := fmt.Sprintf("%s%c%s", home, os.PathSeparator, "cohctl.yaml")

	cluster := details.Getenv(v1.EnvVarCohClusterName)
	port := details.Getenv(v1.EnvVarCohMgmtPrefix + v1.EnvVarCohPortSuffix)
	if port == "" {
		port = fmt.Sprintf("%d", v1.DefaultManagementPort)
	}
	protocol := details.Getenv(v1.EnvVarCohCliProtocol)
	if protocol == "" {
		protocol = "http"
	}

	var buffer bytes.Buffer
	buffer.WriteString("clusters:\n")
	buffer.WriteString("    - name: default\n")
	buffer.WriteString("      discoverytype: manual\n")
	buffer.WriteString("      connectiontype: " + protocol + "\n")
	buffer.WriteString("      connectionurl: " + protocol + "://127.0.0.1:" + port + "/management/coherence/cluster\n")
	buffer.WriteString("      nameservicediscovery: \"\"\n")
	buffer.WriteString("      clusterversion: \"\"\n")
	buffer.WriteString("      clustername: \"" + cluster + "\"\n")
	buffer.WriteString("      clustertype: Standalone\n")
	buffer.WriteString("      manuallycreated: true\n")
	buffer.WriteString("      baseclasspath: \"\"\n")
	buffer.WriteString("      additionalclasspath: \"\"\n")
	buffer.WriteString("      arguments: \"\"\n")
	buffer.WriteString("      managementport: " + port + "\n")
	buffer.WriteString("      persistencemode: \"\"\n")
	buffer.WriteString("      loggingdestination: \"\"\n")
	buffer.WriteString("      managementavailable: false\n")
	buffer.WriteString("color: \"on\"\n")
	buffer.WriteString("currentcontext: default\n")
	buffer.WriteString("debug: false\n")
	buffer.WriteString("defaultbytesformat: m\n")
	buffer.WriteString("ignoreinvalidcerts: false\n")
	buffer.WriteString("requesttimeout: 30\n")
	if err := os.WriteFile(fileName, buffer.Bytes(), os.ModePerm); err != nil {
		return errors.Wrap(err, "failed to write coherence CLI config file "+fileName)
	}

	fmt.Printf("Created CLI config file %s\n", fileName)
	fmt.Println("--------------------")
	fmt.Println(buffer.String())
	fmt.Println("--------------------")

	return nil
}
