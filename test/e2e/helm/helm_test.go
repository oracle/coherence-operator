/*
 * Copyright (c) 2020, 2021 Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package helm

import (
	goctx "context"
	"fmt"
	. "github.com/onsi/gomega"
	"github.com/oracle/coherence-operator/pkg/operator"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"os/exec"
	"strings"
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

	AssertHelmInstall("basic", cmd, g, ns)
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

	AssertHelmInstall("account", cmd, g, ns)
}

func TestHelmInstallWithoutClusterRoles(t *testing.T) {
	g := NewGomegaWithT(t)
	ns := helper.GetTestNamespace()

	t.Cleanup(func() {
		Cleanup(ns, "operator")
	})

	chart, err := helper.FindOperatorHelmChartDir()
	g.Expect(err).NotTo(HaveOccurred())

	cmd := exec.Command("helm", "install",
		"--set", "clusterRoles=false",
		"--set", "image="+helper.GetOperatorImage(),
		"--set", "defaultCoherenceUtilsImage="+helper.GetUtilsImage(),
		"--namespace", ns, "--wait", "operator", chart)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	AssertHelmInstallWithSubTest("no-roles", cmd, g, ns, AssertNoClusterRoles)
}

func TestHelmInstallWithoutClusterRolesWithNodeRole(t *testing.T) {
	g := NewGomegaWithT(t)
	ns := helper.GetTestNamespace()

	t.Cleanup(func() {
		Cleanup(ns, "operator")
	})

	chart, err := helper.FindOperatorHelmChartDir()
	g.Expect(err).NotTo(HaveOccurred())

	cmd := exec.Command("helm", "install",
		"--set", "clusterRoles=false",
		"--set", "nodeRoles=true",
		"--set", "image="+helper.GetOperatorImage(),
		"--set", "defaultCoherenceUtilsImage="+helper.GetUtilsImage(),
		"--namespace", ns, "--wait", "operator", chart)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	AssertHelmInstallWithSubTest("node-role", cmd, g, ns, AssertOnlyNodeClusterRoles)
}

type SubTest func() error

var emptySubTest = func() error {
	return nil
}

func AssertNoClusterRoles() error {
	return AssertRBAC(false)
}

func AssertOnlyNodeClusterRoles() error {
	return AssertRBAC(true)
}

func AssertRBAC(allowNode bool) error {
	rbacClient := testContext.KubeClient.RbacV1()

	crList, err := rbacClient.ClusterRoles().List(goctx.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, cr := range crList.Items {
		if strings.HasPrefix(strings.ToLower(cr.Name), "coherence-operator") {
			if !allowNode || cr.Name != "coherence-operator-node-viewer" {
				return fmt.Errorf("no Coherence ClusterRole shoudl exist but found %s", cr.Name)
			}
		}
	}

	crbList, err := rbacClient.ClusterRoleBindings().List(goctx.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, crb := range crbList.Items {
		if strings.HasPrefix(strings.ToLower(crb.Name), "coherence-operator") {
			if !allowNode || crb.Name != "coherence-operator-node-viewer" {
				return fmt.Errorf("no Coherence ClusterRoleBinding shoudl exist but found %s", crb.Name)
			}
		}
	}

	ns := helper.GetTestNamespace()
	_, err = rbacClient.Roles(ns).Get(goctx.TODO(), "coherence-operator", metav1.GetOptions{})
	if err != nil {
		return err
	}

	_, err = rbacClient.RoleBindings(ns).Get(goctx.TODO(), "coherence-operator", metav1.GetOptions{})
	if err != nil {
		return err
	}

	return nil
}

func AssertHelmInstall(id string, cmd *exec.Cmd, g *GomegaWithT, ns string) {
	AssertHelmInstallWithSubTest(id, cmd, g, ns, emptySubTest)
}

func AssertHelmInstallWithSubTest(id string, cmd *exec.Cmd, g *GomegaWithT, ns string, test SubTest) {
	err := RemoveWebHook()
	g.Expect(err).NotTo(HaveOccurred())
	err = RemoveRBAC()
	g.Expect(err).NotTo(HaveOccurred())

	err = cmd.Run()
	g.Expect(err).NotTo(HaveOccurred())

	pods, err := helper.ListOperatorPods(testContext, ns)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(len(pods)).To(Equal(1))

	deployment, err := helper.NewSingleCoherenceFromYaml(ns, "coherence.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	name := deployment.GetName()
	deployment.SetName(name + "-" + id)

	defer testContext.Client.Delete(goctx.TODO(), &deployment)

	err = testContext.Client.Create(goctx.TODO(), &deployment)
	g.Expect(err).NotTo(HaveOccurred())

	_, err = helper.WaitForStatefulSetForDeployment(testContext, ns, &deployment, helper.RetryInterval, helper.Timeout)
	g.Expect(err).NotTo(HaveOccurred())

	err = test()
	g.Expect(err).NotTo(HaveOccurred())
}

func Cleanup(namespace, name string) {
	_ = helper.WaitForCoherenceCleanup(testContext, namespace)

	cmd := exec.Command("helm", "uninstall", "--namespace", namespace, name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()
	_ = helper.WaitForDeleteOfPodsWithSelector(testContext, namespace, "control-plane=coherence", helper.RetryInterval, helper.Timeout)
}

// Remove the web-hooks that the Operator install creates to
// ensure that nothing is left from a previous test.
func RemoveWebHook() error {
	//DefaultValidatingWebhookName
	client := testContext.KubeClient.AdmissionregistrationV1()

	err := client.MutatingWebhookConfigurations().Delete(goctx.TODO(), operator.DefaultMutatingWebhookName, metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return err
	}

	err = client.ValidatingWebhookConfigurations().Delete(goctx.TODO(), operator.DefaultValidatingWebhookName, metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return err
	}

	return nil
}

// Remove all of the RBAC rules that the Operator install creates to
// ensure that nothing is left from a previous test.
func RemoveRBAC() error {
	var err error
	rbacClient := testContext.KubeClient.RbacV1()
	clusterRolesClient := rbacClient.ClusterRoles()
	clusterRoleBindingsClient := rbacClient.ClusterRoleBindings()

	crList, err := clusterRolesClient.List(goctx.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, cr := range crList.Items {
		if strings.HasPrefix(strings.ToLower(cr.Name), "coherence-operator") {
			if err := clusterRolesClient.Delete(goctx.TODO(), cr.Name, metav1.DeleteOptions{}); err != nil {
				return err
			}
		}
	}

	crbList, err := clusterRoleBindingsClient.List(goctx.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, crb := range crbList.Items {
		if strings.HasPrefix(strings.ToLower(crb.Name), "coherence-operator") {
			if err := clusterRoleBindingsClient.Delete(goctx.TODO(), crb.Name, metav1.DeleteOptions{}); err != nil {
				return err
			}
		}
	}

	ns := helper.GetTestNamespace()
	rolesClient := rbacClient.Roles(ns)
	roleBindingsClient := rbacClient.RoleBindings(ns)

	if role, err := rolesClient.Get(goctx.TODO(), "coherence-operator", metav1.GetOptions{}); err == nil {
		if err := rolesClient.Delete(goctx.TODO(), role.Name, metav1.DeleteOptions{}); err != nil {
			return err
		}
	}

	if roleBinding, err := roleBindingsClient.Get(goctx.TODO(), "coherence-operator", metav1.GetOptions{}); err == nil {
		if err := roleBindingsClient.Delete(goctx.TODO(), roleBinding.Name, metav1.DeleteOptions{}); err != nil {
			return err
		}
	}

	return nil
}
