/*
 * Copyright (c) 2020, 2022, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package helm

import (
	goctx "context"
	"fmt"
	. "github.com/onsi/gomega"
	coh "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/pkg/operator"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
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

func TestNotCreateWebhookCertSecretIfManualCertManager(t *testing.T) {
	g := NewGomegaWithT(t)
	result, err := helmInstall("--set", "webhookCertType=manual")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(result).NotTo(BeNil())

	secret := &corev1.Secret{}
	err = result.Get("coherence-webhook-server-cert", secret)
	g.Expect(err).To(HaveOccurred())
}

func TestBasicHelmInstall(t *testing.T) {
	g := NewGomegaWithT(t)
	cmd, err := createHelmCommand()
	g.Expect(err).NotTo(HaveOccurred())
	AssertHelmInstallWithSubTest(t, "basic", cmd, g, AssertThreeReplicas)
}

func TestSetReplicas(t *testing.T) {
	g := NewGomegaWithT(t)
	cmd, err := createHelmCommand("--set", "replicas=1")
	g.Expect(err).NotTo(HaveOccurred())
	AssertHelmInstallWithSubTest(t, "basic", cmd, g, AssertSingleReplica)
}

func TestSetReResources(t *testing.T) {
	g := NewGomegaWithT(t)
	cmd, err := createHelmCommand("--set", "replicas=1", "--set", "resources.requests.cpu=250m", "--set", "resources.requests.memory=64Mi", "--set", "resources.limits.cpu=512m", "--set", "resources.limits.memory=128Mi")
	g.Expect(err).NotTo(HaveOccurred())
	AssertHelmInstallWithSubTest(t, "basic", cmd, g, AssertResources)
}

func AssertResources() error {
	ns := helper.GetTestNamespace()
	pods, err := helper.ListOperatorPods(testContext, ns)
	if err != nil {
		return err
	}
	if len(pods) == 0 {
		return fmt.Errorf("expected at least one Coherence Operator Pod but found zero")
	}
	resources := pods[0].Spec.Containers[0].Resources

	var qty *resource.Quantity

	qty = resources.Requests.Name("cpu", resource.BinarySI)
	if qty == nil {
		return fmt.Errorf("expected a cpu requests quantity")
	}
	if qty.String() != "250m" {
		return fmt.Errorf("expected a cpu requests of 250m")
	}

	qty = resources.Requests.Name("memory", resource.BinarySI)
	if qty == nil {
		return fmt.Errorf("expected a memory requests quantity")
	}
	if qty.String() != "64Mi" {
		return fmt.Errorf("expected a memory requests of 64Mi")
	}

	qty = resources.Limits.Name("cpu", resource.BinarySI)
	if qty == nil {
		return fmt.Errorf("expected a cpu limits quantity")
	}
	if qty.String() != "512m" {
		return fmt.Errorf("expected a cpu limit of 512m")
	}

	qty = resources.Limits.Name("memory", resource.BinarySI)
	if qty == nil {
		return fmt.Errorf("expected a memory limits quantity")
	}
	if qty.String() != "128Mi" {
		return fmt.Errorf("expected a memory limit of 128Mi")
	}

	return nil
}

func TestHelmInstallWithServiceAccountName(t *testing.T) {
	g := NewGomegaWithT(t)
	cmd, err := createHelmCommand("--set", "serviceAccountName=test-operator-account")
	g.Expect(err).NotTo(HaveOccurred())
	AssertHelmInstallWithSubTest(t, "account", cmd, g, emptySubTest)
}

func TestHelmInstallWithoutClusterRoles(t *testing.T) {
	g := NewGomegaWithT(t)
	cmd, err := createHelmCommand("--set", "clusterRoles=false")
	g.Expect(err).NotTo(HaveOccurred())
	AssertHelmInstallWithSubTest(t, "no-roles", cmd, g, AssertNoClusterRoles)
}

func TestHelmInstallWithoutClusterRolesWithNodeRole(t *testing.T) {
	g := NewGomegaWithT(t)
	cmd, err := createHelmCommand("--set", "clusterRoles=false", "--set", "nodeRoles=true")
	g.Expect(err).NotTo(HaveOccurred())
	AssertHelmInstallWithSubTest(t, "node-role", cmd, g, AssertOnlyNodeClusterRoles)
}

func createHelmCommand(args ...string) (*exec.Cmd, error) {
	ns := helper.GetTestNamespace()

	chart, err := helper.FindOperatorHelmChartDir()
	if err != nil {
		return nil, err
	}

	argList := []string{"install",
		"--set", "image=" + helper.GetOperatorImage(),
		"--set", "defaultCoherenceUtilsImage=" + helper.GetOperatorImage()}

	argList = append(argList, args...)

	argList = append(argList, args...)
	argList = append(argList, "--namespace", ns, "--wait", "operator", chart)

	cmd := exec.Command("helm", argList...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd, nil
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

func AssertSingleReplica() error {
	ns := helper.GetTestNamespace()
	pods, err := helper.ListOperatorPods(testContext, ns)
	if err != nil {
		return err
	}
	count := len(pods)
	if count != 1 {
		return fmt.Errorf("expected a single Coherence Operator Pod but found %d", count)
	}
	return nil
}

func AssertThreeReplicas() error {
	ns := helper.GetTestNamespace()
	pods, err := helper.ListOperatorPods(testContext, ns)
	if err != nil {
		return err
	}
	count := len(pods)
	if count != 3 {
		return fmt.Errorf("expected three Coherence Operator Pods but found %d", count)
	}
	return nil
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
				return fmt.Errorf("no Coherence ClusterRole should exist but found %s", cr.Name)
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
				return fmt.Errorf("no Coherence ClusterRoleBinding should exist but found %s", crb.Name)
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

func AssertHelmInstallWithSubTest(t *testing.T, id string, cmd *exec.Cmd, g *GomegaWithT, test SubTest) {
	ns := helper.GetTestNamespace()

	t.Cleanup(func() {
		if t.Failed() {
			helper.DumpOperatorLogs(t, testContext)
		}
		Cleanup(ns, "operator")
	})

	t.Logf("Asserting Helm install. Removing Webhooks")
	err := RemoveWebHook()
	g.Expect(err).NotTo(HaveOccurred())
	t.Logf("Asserting Helm install. Removing RBAC")
	err = RemoveRBAC()
	g.Expect(err).NotTo(HaveOccurred())

	t.Logf("Asserting Helm install. Performing Helm install")
	err = cmd.Run()
	g.Expect(err).NotTo(HaveOccurred())

	t.Logf("Asserting Helm install. Ensure Operator Pod is ready")
	pods, err := helper.ListOperatorPods(testContext, ns)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(len(pods)).NotTo(Equal(0))

	pod := pods[0]
	err = helper.WaitForPodReady(testContext.KubeClient, pod.Namespace, pod.Name, 10*time.Second, 5*time.Minute)
	g.Expect(err).NotTo(HaveOccurred())

	t.Logf("Asserting Helm install. Deploying Coherence resource")
	deployment, err := helper.NewSingleCoherenceFromYaml(ns, "coherence.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	name := deployment.GetName()
	deployment.SetName(name + "-" + id)

	defer deleteCoherence(t, &deployment)

	err = testContext.Client.Create(goctx.TODO(), &deployment)
	g.Expect(err).NotTo(HaveOccurred())

	_, err = helper.WaitForStatefulSetForDeployment(testContext, ns, &deployment, helper.RetryInterval, helper.Timeout)
	g.Expect(err).NotTo(HaveOccurred())

	err = test()
	g.Expect(err).NotTo(HaveOccurred())
}

func deleteCoherence(t *testing.T, d *coh.Coherence) {
	if err := testContext.Client.Delete(goctx.TODO(), d); err != nil {
		t.Logf("Error deleting Coherence deployment %s - %s", d.GetNamespacedName(), err.Error())
	}
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
	// DefaultValidatingWebhookName
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

// RemoveRBAC removes all the RBAC rules that the Operator install creates to
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
