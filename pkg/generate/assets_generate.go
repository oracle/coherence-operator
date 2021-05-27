/*
 * Copyright (c) 2020, 2021, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package main

import (
	"github.com/shurcooL/vfsgen"
	"log"
	"net/http"
)

var assets http.FileSystem = http.Dir("build/_output/assets")

func main() {
	err := vfsgen.Generate(assets, vfsgen.Options{
		Filename:     "pkg/data/zz_generated_assets.go",
		PackageName:  "data",
		VariableName: "Assets",
	})
	if err != nil {
		log.Fatalln(err)
	}
}
