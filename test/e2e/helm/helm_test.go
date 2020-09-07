/*
 * Copyright (c) 2020 Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package helm

import (
	. "github.com/onsi/gomega"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	"os"
	"os/exec"
	"testing"
)

func TestBasicHelmInstall(t *testing.T) {
	g := NewGomegaWithT(t)

	chart, err := helper.FindOperatorHelmChartDir()
	g.Expect(err).NotTo(HaveOccurred())

	ns := helper.GetTestNamespace()
	cmd := exec.Command("helm", "install",
		"--set", "image="+helper.GetOperatorImage(),
		"--set", "defaultCoherenceUtilsImage="+helper.GetUtilsImage(),
		"--namespace", ns, "--wait", "operator", chart)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	defer Cleanup(ns, "operator")

	err = cmd.Run()
	g.Expect(err).NotTo(HaveOccurred())

	pods, err := helper.ListOperatorPods(testContext, ns)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(len(pods)).To(Equal(1))
}

func Cleanup(namespace, name string) {
	cmd := exec.Command("helm", "uninstall", "--namespace", namespace, name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()
}
