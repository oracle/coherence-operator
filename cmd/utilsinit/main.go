/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package main

import (
	"fmt"
	"github.com/oracle/coherence-operator/pkg/utils"
	"os"
	"strings"
)

const (
	pathSep                 = string(os.PathSeparator)
	utilsDirEnv             = "UTIL_DIR"
	utilsDirDefault         = pathSep + "utils"
	filesDir                = pathSep + "files"
	snapshotDir             = pathSep + "snapshot"
	persistenceDir          = pathSep + "persistence"
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

	scriptsDir := utilDir + string(os.PathSeparator) + "scripts"
	confDir := utilDir + string(os.PathSeparator) + "conf"
	libDir := utilDir + string(os.PathSeparator) + "lib"

	fmt.Printf("Creating target directories under %s\n", utilDir)
	err = os.MkdirAll(scriptsDir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	err = os.MkdirAll(confDir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	err = os.MkdirAll(libDir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Copying files to %s\n", utilDir)

	fmt.Printf("Copying files/*.sh to %s\n", scriptsDir)
	err = utils.CopyDir(filesDir, scriptsDir, func(f string) bool { return strings.HasSuffix(f, ".sh") })
	if err != nil {
		panic(err)
	}

	fmt.Printf("Copying files/*.jar to %s\n", libDir)
	err = utils.CopyDir(filesDir, libDir, func(f string) bool { return strings.HasSuffix(f, ".jar") })
	if err != nil {
		panic(err)
	}

	cp := filesDir + string(os.PathSeparator) + "copy"
	_, err = os.Stat(cp)
	if err == nil {
		fmt.Println("Copying copy utility")
		err = utils.CopyFile(cp, utilDir+string(os.PathSeparator)+"copy")
		if err != nil {
			panic(err)
		}
	}

	opTest := filesDir + string(os.PathSeparator) + "op-test"
	_, err = os.Stat(opTest)
	if err == nil {
		fmt.Println("Copying op-test utility")
		err = utils.CopyFile(opTest, utilDir+string(os.PathSeparator)+"op-test")
		if err != nil {
			panic(err)
		}
	}

	dirNames := []string{snapshotDir, persistenceActiveDir, persistenceTrashDir, persistenceSnapshotsDir}
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
