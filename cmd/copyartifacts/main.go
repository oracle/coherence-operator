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
	"path/filepath"
	"strings"
)

func main() {
	appDir := os.Getenv("APP_DIR")
	extAppDir := os.Getenv("EXTERNAL_APP_DIR")
	libDir := os.Getenv("LIB_DIR")
	extLibDir := os.Getenv("EXTERNAL_LIB_DIR")
	confDir := os.Getenv("CONF_DIR")
	extConfDir := os.Getenv("EXTERNAL_CONF_DIR")

	fmt.Printf("Lib directory is: '%s'\n", libDir)
	fmt.Printf("External lib directory is: '%s'\n", extLibDir)
	fmt.Printf("Config directory is: '%s'\n", confDir)
	fmt.Printf("External Config directory is: '%s'\n", extConfDir)

	_ = os.MkdirAll(extAppDir, os.ModePerm)
	_ = os.MkdirAll(extLibDir, os.ModePerm)
	_ = os.MkdirAll(extLibDir, os.ModePerm)

	_, err := os.Stat(appDir)
	if err == nil {
		err = utils.CopyDir(appDir, extAppDir, utils.AlwaysFilter())
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Printf("App directory '%s' does not exist - no files to copy\n", appDir)
	}

	extLibUnderApp := strings.HasPrefix(filepath.Dir(extLibDir), extAppDir)
	extConfUnderApp := strings.HasPrefix(filepath.Dir(extConfDir), extAppDir)
	libUnderApp := strings.HasPrefix(filepath.Dir(libDir), appDir)
	confUnderApp := strings.HasPrefix(filepath.Dir(confDir), appDir)

	if extLibUnderApp && libUnderApp {
		fmt.Printf("Lib directory '%s' is under App directory '%s' - no files to copy\n", libDir, appDir)
	} else {
		_, err = os.Stat(libDir)
		if err == nil {
			err = utils.CopyDir(libDir, extLibDir, utils.AlwaysFilter())
			if err != nil {
				panic(err)
			}
		} else {
			fmt.Printf("Lib directory '%s' does not exist - no files to copy\n", libDir)
		}
	}

	if extConfUnderApp && confUnderApp {
		fmt.Printf("Config directory '%s' is under App directory '%s' - no files to copy\n", libDir, appDir)
	} else {
		_, err = os.Stat(confDir)
		if err == nil {
			err = utils.CopyDir(confDir, extConfDir, utils.AlwaysFilter())
			if err != nil {
				panic(err)
			}
		} else {
			fmt.Printf("Config directory '%s' does not exist - no files to copy\n", confDir)
		}
	}
}
