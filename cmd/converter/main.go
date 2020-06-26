/*
 * Copyright (c) 2020 Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package main

import (
	"fmt"
	"github.com/oracle/coherence-operator/pkg/apis/coherence/legacy"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		f := os.Args[1]
		err := legacy.Convert(f, os.Stdout)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}
}
