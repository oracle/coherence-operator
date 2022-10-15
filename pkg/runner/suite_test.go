/*
 * Copyright (c) 2020, 2022, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package runner

import (
	"fmt"
	v1 "github.com/oracle/coherence-operator/api/v1"
	"os"
	"testing"
)

var TestAppDir string

// The entry point for the test suite
func TestMain(m *testing.M) {
	// create a temporary folder to represent the app directory
	dir, err := os.CreateTemp("", "operator-tests")
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	TestAppDir = dir.Name()
	err = os.Setenv(v1.EnvVarCohAppDir, TestAppDir)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	err = os.MkdirAll(TestAppDir+string(os.PathSeparator)+"resources", os.ModePerm)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	err = os.MkdirAll(TestAppDir+string(os.PathSeparator)+"classes", os.ModePerm)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	libs := TestAppDir + string(os.PathSeparator) + "libs"
	err = os.MkdirAll(libs, os.ModePerm)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	_, err = os.Create(libs + string(os.PathSeparator) + "foo1.jar")
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	_, err = os.Create(libs + string(os.PathSeparator) + "bar1.JAR")
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	_, err = os.Create(libs + string(os.PathSeparator) + "bar1.txt")
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	classpath := TestAppDir + string(os.PathSeparator) + "classpath"
	err = os.MkdirAll(classpath, os.ModePerm)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	_, err = os.Create(classpath + string(os.PathSeparator) + "foo2.jar")
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	_, err = os.Create(classpath + string(os.PathSeparator) + "bar2.JAR")
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	_, err = os.Create(classpath + string(os.PathSeparator) + "bar2.txt")
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	exitCode := m.Run()

	_ = os.RemoveAll(TestAppDir)
	os.Exit(exitCode)
}
