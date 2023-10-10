/*
 * Copyright (c) 2021, 2023, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package runner

import (
	"fmt"
	v1 "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/pkg/utils"
	"github.com/spf13/cobra"
	"os"
)

const (
	// ArgCommand is an additional command to run after initialisation
	ArgCommand = "cmd"
	// ArgRootDir is the root directory to initialise and directories files in
	ArgRootDir = "root"
	// ArgUtilsDir is the utils directory name
	ArgUtilsDir = "utils"
	// ArgPersistenceDir is root the persistence directory name
	ArgPersistenceDir = "persistence"
	// ArgSnapshotsDir is root the snapshots directory name
	ArgSnapshotsDir = "snapshots"
)

// EnvFunction is a function that returns an environment variable for a given name.
type EnvFunction func(string) string

// initCommand creates the corba "init" sub-command
func initCommand(env map[string]string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   v1.RunnerInit,
		Short: "Initialise a Coherence server",
		Long:  "Initialise a Coherence server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return maybeRun(cmd, initialise)
		},
	}

	utilDir, found := env[v1.EnvVarCohUtilDir]
	if !found || utilDir == "" {
		utilDir = v1.VolumeMountPathUtils
	}

	persistenceDir := v1.VolumeMountPathPersistence
	snapshotDir := v1.VolumeMountPathSnapshots

	flagSet := cmd.Flags()
	flagSet.StringSlice(ArgCommand, nil, "An additional command to run after initialisation")
	flagSet.String(ArgRootDir, "", "The root directory to use to initialise files and directories in")
	flagSet.String(ArgUtilsDir, utilDir, "The utils files root directory")
	flagSet.String(ArgPersistenceDir, persistenceDir, "The root persistence directory")
	flagSet.String(ArgSnapshotsDir, snapshotDir, "The root snapshots directory")

	return cmd
}

// initialise will initialise a Coherence Pod - typically this is run from an init-container
func initialise(details *RunDetails, cmd *cobra.Command) (bool, error) {
	return initialiseWithEnv(cmd, details.Getenv)
}

// initialise will initialise a Coherence Pod - typically this is run from an init-container
func initialiseWithEnv(cmd *cobra.Command, getEnv EnvFunction) (bool, error) {
	var err error

	flagSet := cmd.Flags()

	rootDir, err := flagSet.GetString(ArgRootDir)
	if err != nil {
		return false, err
	}

	pathSep := string(os.PathSeparator)
	filesDir := rootDir + pathSep + "files"
	loggingSrc := filesDir + pathSep + "logging"
	libSrc := filesDir + pathSep + "lib"

	persistenceDir, err := flagSet.GetString(ArgPersistenceDir)
	if err != nil {
		return false, err
	}
	persistenceActiveDir := persistenceDir + pathSep + "active"
	persistenceTrashDir := persistenceDir + pathSep + "trash"
	persistenceSnapshotsDir := persistenceDir + pathSep + "snapshots"

	snapshotDir, err := flagSet.GetString(ArgSnapshotsDir)
	if err != nil {
		return false, err
	}

	fmt.Println("Starting container initialisation")

	utilDir, err := flagSet.GetString(ArgUtilsDir)
	if err != nil {
		return false, err
	}

	loggingDir := utilDir + pathSep + "logging"

	libDir := getEnv(v1.EnvVarCohUtilLibDir)
	if libDir == "" {
		libDir = utilDir + pathSep + "lib"
	}

	fmt.Printf("Creating target directories under %s\n", utilDir)
	err = os.MkdirAll(loggingDir, os.ModePerm)
	if err != nil {
		return false, err
	}

	err = os.MkdirAll(libDir, os.ModePerm)
	if err != nil {
		return false, err
	}

	fmt.Printf("Copying files to %s\n", utilDir)
	fmt.Printf("Copying %s to %s\n", loggingSrc, loggingDir)
	err = utils.CopyDir(loggingSrc, loggingDir, func(f string) bool { return true })
	if err != nil {
		return false, err
	}

	fmt.Printf("Copying %s to %s\n", libSrc, libDir)
	err = utils.CopyDir(libSrc, libDir, func(f string) bool { return true })
	if err != nil {
		return false, err
	}

	cp := filesDir + pathSep + "copy"
	_, err = os.Stat(cp)
	if err == nil {
		fmt.Println("Copying copy utility")
		err = utils.CopyFile(cp, utilDir+pathSep+"copy")
		if err != nil {
			return false, err
		}
	}

	run := filesDir + pathSep + "runner"
	_, err = os.Stat(run)
	if err == nil {
		fmt.Println("Copying runner utility")
		err = utils.CopyFile(run, utilDir+pathSep+"runner")
		if err != nil {
			return false, err
		}
	}

	cohctl := filesDir + pathSep + "cohctl"
	if _, err := os.Stat(cohctl); err == nil {
		fmt.Printf("Copying cohctl utility to \"%s%scohctl\"\n", utilDir, pathSep)
		err = utils.CopyFile(cohctl, utilDir+pathSep+"cohctl")
		if err != nil {
			fmt.Printf("Failed to copy cohctl utility to \"%s%scohctl\" - %s\n", utilDir, pathSep, err.Error())
		}
	}

	var dirNames []string

	_, err = os.Stat(persistenceDir)
	if err == nil {
		// if "/persistence" exists then we'll create the subdirectories
		dirNames = append(dirNames, persistenceActiveDir, persistenceTrashDir, persistenceSnapshotsDir)
	}

	_, err = os.Stat(snapshotDir)
	if err == nil {
		// if "/snapshot" exists then we'll create the cluster snapshot directory
		clusterName := getEnv(v1.EnvVarCohClusterName)
		if clusterName != "" {
			snapshotClusterDir := pathSep + "snapshot" + pathSep + clusterName
			dirNames = append(dirNames, snapshotClusterDir)
		}
	}

	for _, dirName := range dirNames {
		fmt.Printf("Creating directory %s\n", dirName)
		err = os.MkdirAll(dirName, os.ModePerm)
		if err != nil {
			return false, err
		}
		info, err := os.Stat(dirName)
		if err != nil {
			return false, err
		}
		if info.Mode().Perm() != os.ModePerm {
			err = os.Chmod(dirName, os.ModePerm)
			if err != nil {
				return false, err
			}
		}
	}

	fmt.Println("Finished container initialisation")

	c, err := flagSet.GetStringSlice(ArgCommand)
	if err != nil {
		return false, err
	}
	if len(c) != 0 {
		fmt.Printf("Running post initialisation command: %s\n", c)
		_, err = ExecuteWithArgs(nil, c)
		return true, err
	}

	return false, err
}
