/*
 * Copyright (c) 2020, 2021, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package compatibility

import (
	"context"
	. "github.com/onsi/gomega"
	cohv1 "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/test/e2e/helper"
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

	//// delete the previous Operator version
	//t.Logf("Unnstalling previous Operator version: %s\n", version)
	//err = UninstallOperator(ns, name)
	//g.Expect(err).NotTo(HaveOccurred())
	//
	//// wait for Operator Pod to be deleted
	//err = helper.WaitForOperatorCleanup(testContext, ns)
	//g.Expect(err).NotTo(HaveOccurred())

	// Upgrade to this version
	t.Logf("Helm upgrade to current Operator version\n")
	UpgradeToCurrentVersion(g, ns, name)

	// wait a few minutes to allow the new Operator to reconcile the existing Coherence cluster
	t.Logf("Upgraded to current Operator version - waiting for reconcile...\n")
	time.Sleep(2 * time.Minute)

	// Get the current state fo the StatefulSet
	stsAfter := &appsv1.StatefulSet{}
	err = testContext.Client.Get(context.TODO(), types.NamespacedName{Namespace: ns, Name: stsBefore.Name}, stsAfter)
	g.Expect(err).NotTo(HaveOccurred())

	// assert that the StatefulSet has not been updated
	g.Expect(stsAfter.Generation).To(Equal(stsBefore.Generation))

	// scale up to make sure that the Operator can still manage the Coherence cluster
	cmd := exec.Command("kubectl", "-n", ns, "scale", "coherence", d.Name, "--replicas=3")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	g.Expect(err).NotTo(HaveOccurred())
	_ = assertDeploymentEventuallyInDesiredState(t, d, 3)
}

func InstallPreviousVersion(g *GomegaWithT, ns, name, version, selector string) {
	cmd := exec.Command("helm", "install", "--version", version,
		"--namespace", ns, name, "coherence/coherence-operator")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	g.Expect(err).NotTo(HaveOccurred())

	pods, err := helper.WaitForPodsWithSelector(testContext, ns, selector, time.Second*10, time.Minute*5)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(len(pods)).To(Equal(1))
	err = helper.WaitForPodReady(testContext.KubeClient, ns, pods[0].Name, time.Second*10, time.Minute*5)
	g.Expect(err).NotTo(HaveOccurred())
}

func UpgradeToCurrentVersion(g *GomegaWithT, ns, name string) {
	chart, err := helper.FindOperatorHelmChartDir()
	g.Expect(err).NotTo(HaveOccurred())

	cmd := exec.Command("helm", "upgrade",
		"--set", "image="+helper.GetOperatorImage(),
		"--set", "defaultCoherenceUtilsImage="+helper.GetUtilsImage(),
		"--namespace", ns, "--wait", name, chart)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	g.Expect(err).NotTo(HaveOccurred())

	version := os.Getenv("VERSION")
	selector := "app.kubernetes.io/version=" + version

	pods, err := helper.WaitForPodsWithSelector(testContext, ns, selector, time.Second*10, time.Minute*5)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(len(pods)).To(Equal(1))
	err = helper.WaitForPodReady(testContext.KubeClient, ns, pods[0].Name, time.Second*10, time.Minute*5)
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
