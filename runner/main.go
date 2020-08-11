/*
 * Copyright (c) 2020 Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package main

import (
	"fmt"
	"github.com/oracle/coherence-operator/pkg/runner"
	"os"
	"strings"
)

var (
	// BuildInfo is a pipe delimited string of build information injected by the Go linker at build time.
	BuildInfo string
)

func main() {
	err := runner.Run(os.Args, envToMap())
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func envToMap() map[string]string {
	env := make(map[string]string)
	for _, ev := range os.Environ() {
		key, val := getKeyVal(ev)
		env[key] = val
	}

	env[runner.EnvVarBuildInfo] = BuildInfo
	if BuildInfo != "" {
		parts := strings.Split(BuildInfo, "|")

		if len(parts) > 0 {
			env[runner.EnvVarVersion] = parts[0]
		}

		if len(parts) > 1 {
			env[runner.EnvVarGitCommit] = parts[1]
		}

		if len(parts) > 2 {
			env[runner.EnvVarBuildDate] = strings.Replace(parts[2], ".", " ", -1)
		}
	}

	return env
}

func getKeyVal(ev string) (string, string) {
	values := strings.Split(ev, "=")
	switch len(values) {
	case 1:
		return values[0], ""
	case 2:
		return values[0], values[1]
	default:
		return values[0], strings.Join(values[1:], "=")
	}
}
