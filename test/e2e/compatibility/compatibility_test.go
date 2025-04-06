/*
 * Copyright (c) 2020, 2025, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package compatibility

import (
	"context"
	"fmt"
	. "github.com/onsi/gomega"
	cohv1 "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/test/e2e/helm"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	"golang.org/x/mod/semver"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/types"
	"os"
	"os/exec"
	"testing"
	"time"
)

func TestCompatibility(t *testing.T) {
	g := NewGomegaWithT(t)

	ns := helper.GetTestNamespace()
	name := "operator"
	defer CleanupBlind(t, ns, name)

	version := os.Getenv("COMPATIBLE_VERSION")
	g.Expect(version).NotTo(BeEmpty(), "COMPATIBLE_VERSION environment variable has not been set")
	selector := os.Getenv("COMPATIBLE_SELECTOR")
	g.Expect(selector).NotTo(BeEmpty(), "COMPATIBLE_SELECTOR environment variable has not been set")

	// Install Previous version
	t.Logf("Helm install previous Operator version: %s\n", version)
	InstallPreviousVersion(g, ns, name, version, selector)

	// Install a Coherence deployment
	d, err := helper.NewSingleCoherenceFromYaml(ns, "coherence.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	err = nil
	for i := 0; i < 10; i++ {
		err = testContext.Client.Create(context.TODO(), &d)
		if err == nil {
			break
		}
		t.Logf("Coherence cluster install failed, will retry in 5 seconds: %s", err.Error())
		time.Sleep(5 * time.Second)
	}

	g.Expect(err).NotTo(HaveOccurred())
	stsBefore := assertDeploymentEventuallyInDesiredState(t, d, d.GetReplicas())

	// dump the before upgrade state
	dir := fmt.Sprintf("%s-%s-before", t.Name(), version)
	helper.DumpState(testContext, ns, dir)

	if semver.Compare("v"+version, "v3.5.0") < 0 {
		// upgrading from a pre-3.5.0 version so we need to patch the CRDs
		helm.PatchAllCRDs(t, g, name, ns)
	}

	// Upgrade to this version
	UpgradeToCurrentVersion(t, g, ns, name)

	// wait a few minutes to allow the new Operator to reconcile the existing Coherence cluster
	// usually this would be quick, but on a slow build machine it could be a few minutes
	t.Logf("Upgraded to current Operator version - waiting for reconcile...\n")
	time.Sleep(5 * time.Minute)

	// Get the current state of the StatefulSet
	stsAfter := &appsv1.StatefulSet{}
	err = testContext.Client.Get(context.TODO(), types.NamespacedName{Namespace: ns, Name: stsBefore.Name}, stsAfter)
	g.Expect(err).NotTo(HaveOccurred())

	// dump the after upgrade state
	dir = fmt.Sprintf("%s-%s-after", t.Name(), version)
	helper.DumpState(testContext, ns, dir)

	// assert that the StatefulSet has not been updated
	g.Expect(stsAfter.Generation).To(Equal(stsBefore.Generation))

	// scale up to make sure that the Operator can still manage the Coherence cluster
	n := fmt.Sprintf("coherence/%s", d.Name)
	t.Logf("Scaling coherence resource %s in namespace %s to 3 replicas\n", n, ns)
	cmd := exec.Command("kubectl", "-n", ns, "scale", n, "--replicas=3")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	g.Expect(err).NotTo(HaveOccurred())
	_ = assertDeploymentEventuallyInDesiredState(t, d, 3)
}

func InstallPreviousVersion(g *GomegaWithT, ns, name, version, selector string) {
	chartDir, err := helper.FindOperatorTestHelmChartDir()
	g.Expect(err).NotTo(HaveOccurred())

	prevDir := chartDir + string(os.PathSeparator) + version

	err = os.RemoveAll(prevDir)
	g.Expect(err).NotTo(HaveOccurred())

	err = os.MkdirAll(prevDir, os.ModePerm)
	g.Expect(err).NotTo(HaveOccurred())

	cmd := exec.Command("helm", "fetch", "--version", version,
		"--untar", "--untardir", prevDir, "coherence/coherence-operator")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	g.Expect(err).NotTo(HaveOccurred())

	valuesFile := prevDir + string(os.PathSeparator) + "coherence-operator" + string(os.PathSeparator) + "values.yaml"

	values := helper.OperatorValues{}
	err = values.LoadFromYaml(valuesFile)
	g.Expect(err).NotTo(HaveOccurred())

	cmd = exec.Command("helm", "install", "--version", version,
		"--namespace", ns, name, "coherence/coherence-operator")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	g.Expect(err).NotTo(HaveOccurred())

	replicas := values.GetReplicas(1)
	pods, err := helper.WaitForPodsWithSelectorAndReplicas(testContext, ns, selector, replicas, time.Second*10, time.Minute*5)
	g.Expect(err).NotTo(HaveOccurred())
	err = helper.WaitForPodReady(testContext, ns, pods[0].Name, time.Second*10, time.Minute*5)
	g.Expect(err).NotTo(HaveOccurred())
}

func UpgradeToCurrentVersion(t *testing.T, g *GomegaWithT, ns, name string) {
	t.Logf("Helm upgrade to current Operator version\n")
	chart, err := helper.FindOperatorHelmChartDir()
	g.Expect(err).NotTo(HaveOccurred())

	cmd := exec.Command("helm", "upgrade",
		"--set", "image="+helper.GetOperatorImage(),
		"--set", "defaultCoherenceUtilsImage="+helper.GetOperatorImage(),
		"--namespace", ns, "--wait", name, chart)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	t.Logf("Helm upgrade to current Operator version - executing Helm upgrade\n")
	err = cmd.Run()
	g.Expect(err).NotTo(HaveOccurred())
	t.Logf("Helm upgrade to current Operator version - Helm upgrade successful\n")

	version := os.Getenv("VERSION")
	selector := "app.kubernetes.io/version=" + version

	t.Logf("Helm upgrade to current Operator version - Waiting for pods in namespace %s with selector \"%s\"\n", ns, selector)
	pods, err := helper.WaitForPodsWithSelectorAndReplicas(testContext, ns, selector, 3, time.Second*10, time.Minute*5)
	g.Expect(err).NotTo(HaveOccurred())
	t.Logf("Helm upgrade to current Operator version - Waiting for Pods %s in namespace %s to be ready\n", pods[0].Name, ns)
	err = helper.WaitForPodReady(testContext, ns, pods[0].Name, time.Second*10, time.Minute*5)
	g.Expect(err).NotTo(HaveOccurred())
}

func CleanupBlind(t *testing.T, namespace, name string) {
	_ = Cleanup(t, namespace, name)
}

func Cleanup(t *testing.T, namespace, name string) error {
	helper.DumpOperatorLogs(t, testContext)
	_ = helper.WaitForCoherenceCleanup(testContext, namespace)
	return UninstallOperator(namespace, name)
}

func UninstallOperator(namespace, name string) error {
	cmd := exec.Command("helm", "delete", "--namespace", namespace, name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// assertDeploymentEventuallyInDesiredState asserts that a Coherence resource exists and has the correct spec and that the
// underlying StatefulSet exists with the correct status and ready replicas.
func assertDeploymentEventuallyInDesiredState(t *testing.T, d cohv1.Coherence, replicas int32) *appsv1.StatefulSet {
	g := NewGomegaWithT(t)

	testContext.Logf("Asserting Coherence resource %s exists with %d replicas", d.Name, replicas)

	// create a DeploymentStateCondition that checks a deployment's replica count
	condition := helper.ReplicaCountCondition(replicas)

	// wait for the deployment to match the condition
	_, err := helper.WaitForCoherenceCondition(testContext, d.Namespace, d.Name, condition, time.Second*10, time.Minute*5)
	g.Expect(err).NotTo(HaveOccurred())

	testContext.Logf("Asserting StatefulSet %s exists with %d replicas", d.Name, replicas)

	// wait for the StatefulSet to have the required ready replicas
	sts, err := helper.WaitForStatefulSet(testContext, d.Namespace, d.Name, replicas, time.Second*10, time.Minute*5)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(sts.Status.ReadyReplicas).To(Equal(replicas))

	testContext.Logf("Asserting StatefulSet %s exist with %d replicas - Done!", d.Name, replicas)

	return sts
}
