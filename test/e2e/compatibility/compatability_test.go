/*
 * Copyright (c) 2020, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package compatibility_test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-test/deep"
	. "github.com/onsi/gomega"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	coh "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/pointer"
	"os"
	"testing"
	"time"
)

func TestCompatibility(t *testing.T) {
	helper.AssumeRunningCompatibilityTests(t)
	versions := helper.GetCompatibleOperatorVersions()
	for _, version := range versions {
		name := fmt.Sprintf("%s", version)
		t.Run(name, func(t *testing.T) {
			assertCompatibilityForVersion(t, version)
		})
	}
}

func assertCompatibilityForVersion(t *testing.T, prevVersion string) {
	var err error

	g := NewGomegaWithT(t)
	f := framework.Global
	namespace := f.OperatorNamespace

	values := helper.OperatorValues{}

	chart, err := helper.FindPreviousOperatorHelmChartDir(prevVersion)
	g.Expect(err).NotTo(HaveOccurred())
	t.Logf("Running compatibility test against previous chart %s", chart)

	_, err = os.Stat(chart)
	if err != nil {
		t.Skipf("Skipping compatibility test. Cannot locate previous version chart at %s", chart)
	}

	ctx := helper.CreateTestContext(t)
	// Make sure we defer clean-up (uninstall the operator and Coherence cluster) when we're done
	cleaner := Cleanup{t: t, ctx: ctx}
	defer cleaner.Run()

	// Create the previous version helper.HelmHelper
	hhPrev, err := helper.NewOperatorChartHelperForChart(chart)
	g.Expect(err).ToNot(HaveOccurred())

	cl := hhPrev.KubeClient

	// Create a previous version HelmReleaseManager with a release name and values
	rmPrev, err := hhPrev.NewOperatorHelmReleaseManager("prev-operator", &values)
	g.Expect(err).ToNot(HaveOccurred())
	defer CleanupHelm(t, rmPrev, hhPrev)
	cleaner.rm = rmPrev
	cleaner.hh = hhPrev

	// Install the previous Operator chart
	t.Logf("Installing previous Operator version %s", prevVersion)
	_, err = rmPrev.InstallRelease()
	g.Expect(err).ToNot(HaveOccurred())

	// The chart is installed but the Pod(s) may not exist yet so wait for it...
	// (we wait a maximum of 5 minutes, retrying every 10 seconds)
	pods, err := helper.WaitForOperatorPods(cl, namespace, time.Second*10, time.Minute*5)
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(len(pods)).To(Equal(1))

	// Deploy the Coherence cluster using the previous operator
	t.Logf("Installing Coherence cluster")
	cluster, err := DeployCoherenceCluster(t, ctx, namespace, "coherence.yaml")
	g.Expect(err).ToNot(HaveOccurred())

	// Get the cluster StatefulSet before Operator Upgrade
	role := cluster.GetFirstRole()
	stsBefore, err := helper.WaitForStatefulSetForDeployment(cl, namespace, &cluster, role, time.Second*10, time.Minute*5, t)
	g.Expect(err).ToNot(HaveOccurred())

	dir := t.Name() + "/Before"
	helper.DumpOperatorLog(framework.Global.KubeClient, namespace, dir, t)
	helper.DumpState(namespace, dir, t)

	// Upgrade to the current Operator - we do this by running cleanup to remove the previous operator and then install the new one
	t.Logf("Removing previous operator version %s", prevVersion)
	_, err = rmPrev.UninstallRelease()
	g.Expect(err).ToNot(HaveOccurred())

	// Wait for the Operator Pod to be removed
	t.Log("Waiting for removal of previous operator...")
	err = helper.WaitForOperatorDeletion(cl, namespace, time.Second*10, time.Minute*5, t)
	g.Expect(err).ToNot(HaveOccurred())

	// Create a current version HelmReleaseManager with a release name and values
	version := helper.GetOperatorVersion()

	// Create the current version helper.HelmHelper
	hhCurr, err := helper.NewOperatorChartHelper()
	g.Expect(err).ToNot(HaveOccurred())

	t.Logf("Installing current version of Operator %s", version)
	rmCurr, err := hhCurr.NewOperatorHelmReleaseManager("current-operator", &values)
	g.Expect(err).ToNot(HaveOccurred())
	cleaner.rm = rmCurr
	cleaner.hh = hhCurr

	_, err = rmCurr.InstallRelease()
	g.Expect(err).ToNot(HaveOccurred())

	// The chart is installed but the Pod may not exist yet so wait for it...
	// (we wait a maximum of 5 minutes, retrying every 10 seconds)
	pods, err = helper.WaitForOperatorPods(hhCurr.KubeClient, namespace, time.Second*10, time.Minute*5)
	d, err := json.Marshal(pods[0])
	g.Expect(err).ToNot(HaveOccurred())
	t.Logf("JSON for new Operator Pod version %s:\n%s", version, string(d))
	image := helper.GetOperatorImage()
	g.Expect(pods[0].Spec.Containers[0].Image).To(Equal(image))

	// wait for one minute and ensure that the Coherence Cluster did not restart
	time.Sleep(time.Minute * 1)

	// Get the cluster StatefulSet after Operator Upgrade
	stsAfter, err := helper.WaitForStatefulSetForDeployment(cl, namespace, &cluster, role, time.Second*10, time.Minute*5, t)
	g.Expect(err).ToNot(HaveOccurred())

	// Assert that the StatefulSet has not been changed by the upgrade (i.e. its generation is unchanged)
	diffs := deep.Equal(stsBefore, stsAfter)
	t.Logf("Difference between StatefulSet for v%s and v%s:\n", prevVersion, version)
	for _, diff := range diffs {
		t.Log(diff)
	}
	g.Expect(stsAfter.Generation).To(Equal(stsBefore.Generation))

	// Ensure that everything is still linked, i.e. the cluster to the role to the StatefulSet
	// by upgrading the Coherence cluster. We're just going to add a label to the Pods

	t.Log("Re-fetching Coherence cluster to update with Pod labels")

	// re-fetch the Coherence cluster as it might have changed
	err = f.Client.Get(context.TODO(), types.NamespacedName{Namespace: namespace, Name: cluster.Name}, &cluster)
	g.Expect(err).ToNot(HaveOccurred())
	// add the labels
	labels := cluster.Spec.CoherenceRoleSpec.Labels
	if labels == nil {
		labels = make(map[string]string)
	}
	labels["foo"] = "bar"
	cluster.Spec.CoherenceRoleSpec.Labels = labels

	// Update the Coherence cluster in k8s
	t.Log("Updating Coherence cluster")
	err = f.Client.Update(context.TODO(), &cluster)
	g.Expect(err).ToNot(HaveOccurred())

	// wait for at all Pods with the new label, this will verify that the update worked and the Coherence cluster is still good
	t.Log("Waiting for Pods to update with new labels")
	_, err = helper.WaitForPodsWithLabel(cl, namespace, "foo=bar", int(role.GetReplicas()), time.Second*10, time.Minute*5)
	g.Expect(err).ToNot(HaveOccurred())
	// Get the cluster StatefulSet after cluster update
	_, err = helper.WaitForStatefulSetForDeployment(cl, namespace, &cluster, role, time.Second*10, time.Minute*5, t)
	g.Expect(err).ToNot(HaveOccurred())

	// Delete the CoherenceCluster
	t.Logf("Re-fetching Coherence cluster %s/%s to delete it", namespace, cluster.Name)
	err = f.Client.Get(context.TODO(), types.NamespacedName{Namespace: namespace, Name: cluster.Name}, &cluster)
	g.Expect(err).ToNot(HaveOccurred())
	cr := coh.CoherenceRole{}
	cr.SetNamespace(namespace)
	cr.SetName(cluster.GetFullRoleName(role.GetRoleName()))
	t.Logf("Fetching CoherenceRole %s/%s", namespace, cr.GetName())
	err = f.Client.Get(context.TODO(), types.NamespacedName{Namespace: namespace, Name: cr.GetName()}, &cr)
	g.Expect(err).ToNot(HaveOccurred())
	sts := appsv1.StatefulSet{}
	err = f.Client.Get(context.TODO(), types.NamespacedName{Namespace: namespace, Name: cr.GetName()}, &sts)
	g.Expect(err).ToNot(HaveOccurred())

	t.Logf("Deleting Coherence cluster %s/%s", namespace, cluster.Name)
	err = f.Client.Delete(context.TODO(), &cluster)
	g.Expect(err).ToNot(HaveOccurred())

	t.Logf("Waiting for CoherenceCluster %s/%s to be removed", namespace, cluster.Name)
	err = helper.WaitForDeletion(f, namespace, cluster.Name, &cluster, time.Second*10, time.Minute*5, t)
	g.Expect(err).ToNot(HaveOccurred())

	t.Logf("Waiting for CoherenceRole %s/%s to be removed", namespace, cr.Name)
	err = helper.WaitForDeletion(f, namespace, cr.Name, &cr, time.Second*10, time.Minute*5, t)
	g.Expect(err).ToNot(HaveOccurred())

	t.Logf("Waiting for StatefulSet %s/%s to be removed", namespace, sts.GetName())
	err = helper.WaitForDeletion(f, namespace, sts.GetName(), &sts, time.Second*10, time.Minute*5, t)
	g.Expect(err).ToNot(HaveOccurred())

	// Wait for the updated cluster Pods to be deleted
	t.Log("Waiting for Coherence cluster Pods to be removed")
	selector := fmt.Sprintf("coherenceCluster=%s", cluster.Name)
	err = helper.WaitForDeleteOfPodsWithSelector(cl, namespace, selector, time.Second*10, time.Minute*5, t)
	g.Expect(err).ToNot(HaveOccurred())
}
