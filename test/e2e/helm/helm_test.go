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
	cmd := exec.Command("helm", "install", "--namespace", ns, "--wait", "operator", chart)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	defer Cleanup(ns, "operator")

	err = cmd.Run()
	g.Expect(err).NotTo(HaveOccurred())
}

func Cleanup(namespace, name string) {
	cmd := exec.Command("helm", "uninstall", "--namespace", namespace, name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()
}
