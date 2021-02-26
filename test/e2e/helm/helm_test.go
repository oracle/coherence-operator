/*
 * Copyright (c) 2020, 2021 Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package helm

import (
	goctx "context"
	. "github.com/onsi/gomega"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"os"
	"os/exec"
	"testing"
)

func TestCreateWebhookCertSecretByDefault(t *testing.T) {
	g := NewGomegaWithT(t)
	result, err := helmInstall()
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(result).NotTo(BeNil())

	secret := &corev1.Secret{}
	err = result.Get("coherence-webhook-server-cert", secret)
	g.Expect(err).NotTo(HaveOccurred())
}

func TestCreateWebhookCertSecretWithName(t *testing.T) {
	g := NewGomegaWithT(t)
	result, err := helmInstall("--set", "webhookCertSecret=foo")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(result).NotTo(BeNil())

	secret := &corev1.Secret{}
	err = result.Get("foo", secret)
	g.Expect(err).NotTo(HaveOccurred())

	dep := &appsv1.Deployment{}
	err = result.Get("coherence-operator", dep)
	g.Expect(err).NotTo(HaveOccurred())

	c := findContainer("manager", dep)
	g.Expect(c).NotTo(BeNil())

	senv := findEnvVar("WEBHOOK_SECRET", c)
	g.Expect(senv).NotTo(BeNil())
	g.Expect(senv.Value).To(Equal("foo"))
}

func TestNotCreateWebhookCertSecretIfNotSelfSigned(t *testing.T) {
	g := NewGomegaWithT(t)
	result, err := helmInstall("--set", "webhookCertType=cert-manager")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(result).NotTo(BeNil())

	secret := &corev1.Secret{}
	err = result.Get("coherence-webhook-server-cert", secret)
	g.Expect(err).To(HaveOccurred())
}

func TestBasicHelmInstall(t *testing.T) {
	g := NewGomegaWithT(t)
	ns := helper.GetTestNamespace()

	t.Cleanup(func() {
		Cleanup(ns, "operator")
	})

	chart, err := helper.FindOperatorHelmChartDir()
	g.Expect(err).NotTo(HaveOccurred())

	cmd := exec.Command("helm", "install",
		"--set", "image="+helper.GetOperatorImage(),
		"--set", "defaultCoherenceUtilsImage="+helper.GetUtilsImage(),
		"--namespace", ns, "--wait", "operator", chart)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	AssertHelmInstall(t, cmd, g, ns)
}

func TestHelmInstallWithServiceAccountName(t *testing.T) {
	g := NewGomegaWithT(t)
	ns := helper.GetTestNamespace()

	t.Cleanup(func() {
		Cleanup(ns, "operator")
	})

	chart, err := helper.FindOperatorHelmChartDir()
	g.Expect(err).NotTo(HaveOccurred())

	cmd := exec.Command("helm", "install",
		"--set", "serviceAccountName=test-operator-account",
		"--set", "image="+helper.GetOperatorImage(),
		"--set", "defaultCoherenceUtilsImage="+helper.GetUtilsImage(),
		"--namespace", ns, "--wait", "operator", chart)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	AssertHelmInstall(t, cmd, g, ns)
}

func AssertHelmInstall(t *testing.T, cmd *exec.Cmd, g *GomegaWithT, ns string) {
	err := cmd.Run()
	g.Expect(err).NotTo(HaveOccurred())

	pods, err := helper.ListOperatorPods(testContext, ns)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(len(pods)).To(Equal(1))

	deployment, err := helper.NewSingleCoherenceFromYaml(ns, "coherence.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	defer testContext.Client.Delete(goctx.TODO(), &deployment)

	err = testContext.Client.Create(goctx.TODO(), &deployment)
	g.Expect(err).NotTo(HaveOccurred())

	_, err = helper.WaitForStatefulSetForDeployment(testContext, ns, &deployment, helper.RetryInterval, helper.Timeout)
	g.Expect(err).NotTo(HaveOccurred())
}

func Cleanup(namespace, name string) {
	cmd := exec.Command("helm", "uninstall", "--namespace", namespace, name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()
	_ = helper.WaitForDeleteOfPodsWithSelector(testContext, namespace, "control-plane=coherence", helper.RetryInterval, helper.Timeout)
}
