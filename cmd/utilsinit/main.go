/*
 * Copyright (c) 2019, 2020, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package main

import (
	"fmt"
	v1 "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	"github.com/oracle/coherence-operator/pkg/utils"
	"os"
)

const (
	pathSep                 = string(os.PathSeparator)
	utilsDirEnv             = "UTIL_DIR"
	clusterEnv              = "COH_CLUSTER_NAME"
	filesDir                = pathSep + "files"
	configSrc               = filesDir + pathSep + "config"
	loggingSrc              = filesDir + pathSep + "logging"
	libSrc                  = filesDir + pathSep + "lib"
	utilsDirDefault         = v1.VolumeMountPathUtils
	snapshotDir             = v1.VolumeMountPathSnapshots
	persistenceDir          = v1.VolumeMountPathPersistence
	persistenceActiveDir    = persistenceDir + pathSep + "active"
	persistenceTrashDir     = persistenceDir + pathSep + "trash"
	persistenceSnapshotsDir = persistenceDir + pathSep + "snapshots"
)

func main() {
	var err error

	fmt.Println("Starting container initialisation")

	utilDir := os.Getenv(utilsDirEnv)
	if utilDir == "" {
		utilDir = utilsDirDefault
	}

	configDir := utilDir + pathSep + "config"
	loggingDir := utilDir + pathSep + "logging"
	libDir := utilDir + pathSep + "lib"

	fmt.Printf("Creating target directories under %s\n", utilDir)
	err = os.MkdirAll(configDir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	err = os.MkdirAll(loggingDir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	err = os.MkdirAll(libDir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Copying files to %s\n", utilDir)
	fmt.Printf("Copying %s to %s\n", configSrc, configDir)
	err = utils.CopyDir(configSrc, configDir, func(f string) bool { return true })
	if err != nil {
		panic(err)
	}

	fmt.Printf("Copying %s to %s\n", loggingSrc, loggingDir)
	err = utils.CopyDir(loggingSrc, loggingDir, func(f string) bool { return true })
	if err != nil {
		panic(err)
	}

	fmt.Printf("Copying %s to %s\n", libSrc, libDir)
	err = utils.CopyDir(libSrc, libDir, func(f string) bool { return true })
	if err != nil {
		panic(err)
	}

	cp := filesDir + pathSep + "copy"
	_, err = os.Stat(cp)
	if err == nil {
		fmt.Println("Copying copy utility")
		err = utils.CopyFile(cp, utilDir+pathSep+"copy")
		if err != nil {
			panic(err)
		}
	}

	run := filesDir + pathSep + "runner"
	_, err = os.Stat(run)
	if err == nil {
		fmt.Println("Copying runner utility")
		err = utils.CopyFile(run, utilDir+pathSep+"runner")
		if err != nil {
			panic(err)
		}
	}

	opTest := filesDir + pathSep + "op-test"
	_, err = os.Stat(opTest)
	if err == nil {
		fmt.Println("Copying op-test utility")
		err = utils.CopyFile(opTest, utilDir+pathSep+"op-test")
		if err != nil {
			panic(err)
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
		clusterName := os.Getenv(clusterEnv)
		if clusterName != "" {
			snapshotClusterDir := pathSep + "snapshot" + pathSep + clusterName
			dirNames = append(dirNames, snapshotClusterDir)
		}
	}

	for _, dirName := range dirNames {
		fmt.Printf("Creating directory %s\n", dirName)
		err = os.MkdirAll(dirName, os.ModePerm)
		if err != nil {
			panic(err)
		}
		info, err := os.Stat(dirName)
		if err != nil {
			panic(err)
		}
		if info.Mode().Perm() != os.ModePerm {
			err = os.Chmod(dirName, os.ModePerm)
			if err != nil {
				panic(err)
			}
		}
	}

	fmt.Println("Finished container initialisation")
}
