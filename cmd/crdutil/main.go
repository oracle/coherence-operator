/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package main

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	"io/ioutil"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"os"
)

// A utility to set-up the CRDs in the Helm Chart at build time
func main() {
	pathSep := string(os.PathSeparator)
	chartDir, err := helper.FindOperatorHelmChartDir()
	fmt.Println("chart dir is " + chartDir)
	templateDir := chartDir + pathSep + "templates"
	fmt.Println("template dir is " + chartDir)
	panicIfErr(err)

	err = os.MkdirAll(templateDir, os.ModePerm)
	panicIfErr(err)

	crdDir, err := helper.FindCrdDir()
	panicIfErr(err)

	crds := []string{
		"coherence_v1_coherencecluster_crd.yaml",
		"coherence_v1_coherenceinternal_crd.yaml",
		"coherence_v1_coherencerole_crd.yaml",
	}

	for _, crdName := range crds {
		fmt.Println("crdutil: processing CRD " + crdName)
		crd := v1beta1.CustomResourceDefinition{}
		fileName := crdDir + pathSep + crdName
		fmt.Println("crdutil: readin CRD " + crdName + " from " + fileName)
		d, err := ioutil.ReadFile(fileName)
		panicIfErr(err)
		fmt.Println("crdutil: unmarshal CRD " + crdName)
		err = yaml.Unmarshal(d, &crd)
		panicIfErr(err)

		fmt.Println("crdutil: adding annotations to CRD " + crdName)
		var ann map[string]string
		if crd.Annotations == nil {
			ann = make(map[string]string)
		} else {
			ann = crd.Annotations
		}

		ann["helm.sh/hook"] = "crd-install"
		ann["helm.sh/hook-delete-policy"] = "before-hook-creation"
		crd.Annotations = ann

		fmt.Println("crdutil: marshalling CRD " + crdName)
		d, err = yaml.Marshal(crd)
		panicIfErr(err)

		templateName := templateDir + pathSep + crdName
		fmt.Println("crdutil: writing CRD " + crdName + " to " + templateName)
		out, err := os.Create(templateName)
		panicIfErr(err)
		_, err = out.WriteString("{{- if and .Release.IsInstall .Values.createCustomResource -}}\n")
		panicIfErr(err)
		_, err = out.Write(d)
		panicIfErr(err)
		_, err = out.WriteString("{{- end }}\n")
		panicIfErr(err)
		err = out.Close()
		panicIfErr(err)
		fmt.Println("crdutil: finished CRD " + crdName)
	}
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}
