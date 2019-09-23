/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1_test

import (
	"errors"
	"fmt"
	"github.com/ghodss/yaml"
	v1 "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// Returns a pointer to an int32
func int32Ptr(x int32) *int32 {
	return &x
}

// Returns a pointer to an int32
func boolPtr(x bool) *bool {
	return &x
}

// Returns a pointer to a string
func stringPtr(x string) *string {
	return &x
}

func LoadCoherenceRoleFromCoherenceClusterYamlFile(file string) v1.CoherenceRoleSpec {
	cluster, err := LoadCoherenceClusterFromYamlFile(file)
	if err != nil {
		fmt.Println(err)
		return v1.CoherenceRoleSpec{}
	}
	return cluster.Spec.CoherenceRoleSpec
}

func LoadCoherenceClusterFromYamlFile(file string) (*v1.CoherenceCluster, error) {
	if file == "" {
		return nil, errors.New("missing file name")
	}

	actualFile, err := findActualFile(file)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(actualFile)
	if err != nil {
		return nil, errors.New("Failed to read file " + actualFile + " caused by " + err.Error())
	}

	// expand any ${env-var} references in the yaml file
	s := os.ExpandEnv(string(data))

	cluster := &v1.CoherenceCluster{}
	err = yaml.Unmarshal([]byte(s), cluster)
	if err != nil {
		return nil, errors.New("Failed to parse yaml file " + actualFile + " caused by " + err.Error())
	}

	return cluster, nil
}

func findActualFile(file string) (string, error) {
	_, err := os.Stat(file)
	if err == nil {
		return file, nil
	}

	// files does not exist
	if !strings.HasPrefix(file, "/") {
		// the file does not exist and is not absolute so try relative to a location
		// in the call stack by walking up the stack and trying each location.
		i := 0
		for {
			_, caller, _, ok := runtime.Caller(i)
			if ok {
				dir := filepath.Dir(caller)
				f := dir + string(os.PathSeparator) + file
				_, e := os.Stat(f)
				if e == nil {
					return f, nil
				}
			} else {
				// no more call stack
				break
			}
			i = i + 1
		}
	}

	return "", err
}
