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
)

func main() {
	libDir := os.Getenv("LIB_DIR")
	extLibDir := os.Getenv("EXTERNAL_LIB_DIR")
	confDir := os.Getenv("CONF_DIR")
	extConfDir := os.Getenv("EXTERNAL_CONF_DIR")

	fmt.Printf("Lib directory is: '%s'\n", libDir)
	fmt.Printf("External lib directory is: '%s'\n", extLibDir)
	fmt.Printf("Config directory is: '%s'\n", confDir)
	fmt.Printf("External Config directory is: '%s'\n", extConfDir)

	_ = os.MkdirAll(extLibDir, os.ModePerm)
	_ = os.MkdirAll(extLibDir, os.ModePerm)

	_, err := os.Stat(libDir)
	if err == nil {
		err = utils.CopyDir(libDir, extLibDir, utils.AlwaysFilter())
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Printf("Lib directory '%s' does not exist - no files to copy\n", libDir)
	}

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
