/*
 * Copyright (c) 2021, Oracle and/or its affiliates.
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

// initCommand creates the corba "init" sub-command
func initCommand() *cobra.Command {
	return &cobra.Command{
		Use:   v1.RunnerInit,
		Short: "Initialise a Coherence server",
		Long:  "Initialise a Coherence server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return initialise()
		},
	}
}

// initialise will initialise a Coherence Pod - typically this is run from an init-container
func initialise() error {
	var err error

	pathSep := string(os.PathSeparator)
	filesDir := pathSep + "files"
	loggingSrc := filesDir + pathSep + "logging"
	libSrc := filesDir + pathSep + "lib"
	snapshotDir := v1.VolumeMountPathSnapshots
	persistenceDir := v1.VolumeMountPathPersistence
	persistenceActiveDir := persistenceDir + pathSep + "active"
	persistenceTrashDir := persistenceDir + pathSep + "trash"
	persistenceSnapshotsDir := persistenceDir + pathSep + "snapshots"

	fmt.Println("Starting container initialisation")

	utilDir := os.Getenv(v1.EnvVarCohUtilDir)
	if utilDir == "" {
		utilDir = v1.VolumeMountPathUtils
	}

	loggingDir := utilDir + pathSep + "logging"

	libDir := os.Getenv(v1.EnvVarCohUtilLibDir)
	if libDir == "" {
		libDir = utilDir + pathSep + "lib"
	}

	fmt.Printf("Creating target directories under %s\n", utilDir)
	err = os.MkdirAll(loggingDir, os.ModePerm)
	if err != nil {
		return err
	}

	err = os.MkdirAll(libDir, os.ModePerm)
	if err != nil {
		return err
	}

	fmt.Printf("Copying files to %s\n", utilDir)
	fmt.Printf("Copying %s to %s\n", loggingSrc, loggingDir)
	err = utils.CopyDir(loggingSrc, loggingDir, func(f string) bool { return true })
	if err != nil {
		return err
	}

	fmt.Printf("Copying %s to %s\n", libSrc, libDir)
	err = utils.CopyDir(libSrc, libDir, func(f string) bool { return true })
	if err != nil {
		return err
	}

	cp := filesDir + pathSep + "copy"
	_, err = os.Stat(cp)
	if err == nil {
		fmt.Println("Copying copy utility")
		err = utils.CopyFile(cp, utilDir+pathSep+"copy")
		if err != nil {
			return err
		}
	}

	run := filesDir + pathSep + "runner"
	_, err = os.Stat(run)
	if err == nil {
		fmt.Println("Copying runner utility")
		err = utils.CopyFile(run, utilDir+pathSep+"runner")
		if err != nil {
			return err
		}
	}

	var dirNames []string

	_, err = os.Stat(persistenceDir)
	if err == nil {
		// if "/persistence" exists then we'll create the sub-directories
		dirNames = append(dirNames, persistenceActiveDir, persistenceTrashDir, persistenceSnapshotsDir)
	}

	_, err = os.Stat(snapshotDir)
	if err == nil {
		// if "/snapshot" exists then we'll create the cluster snapshot directory
		clusterName := os.Getenv(v1.EnvVarCohClusterName)
		if clusterName != "" {
			snapshotClusterDir := pathSep + "snapshot" + pathSep + clusterName
			dirNames = append(dirNames, snapshotClusterDir)
		}
	}

	for _, dirName := range dirNames {
		fmt.Printf("Creating directory %s\n", dirName)
		err = os.MkdirAll(dirName, os.ModePerm)
		if err != nil {
			return err
		}
		info, err := os.Stat(dirName)
		if err != nil {
			return err
		}
		if info.Mode().Perm() != os.ModePerm {
			err = os.Chmod(dirName, os.ModePerm)
			if err != nil {
				return err
			}
		}
	}

	fmt.Println("Finished container initialisation")
	return nil
}
