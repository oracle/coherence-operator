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
	"github.com/oracle/coherence-operator/test/e2e/helper"
	"k8s.io/apimachinery/pkg/types"
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
	defer helper.DumpOperatorLogsAndCleanup(t, ctx)

	// Create the previous version helper.HelmHelper
	helmHelperPrevious, err := helper.NewOperatorChartHelperForChart(chart)
	if err != nil {
		t.Fatal(err)
	}

	// Create the current version helper.HelmHelper
	helmHelperCurrent, err := helper.NewOperatorChartHelper()
	if err != nil {
		t.Fatal(err)
	}

	namespace := helmHelperCurrent.Namespace
	cl := helmHelperCurrent.KubeClient

	// Create a previous version HelmReleaseManager with a release name and values
	rmPrevious, err := helmHelperPrevious.NewOperatorHelmReleaseManager("prev-operator", &values)
	g.Expect(err).ToNot(HaveOccurred())
	defer CleanupHelm(t, rmPrevious, helmHelperPrevious)

	// Create a current version HelmReleaseManager with a release name and values
	rmCurrent, err := helmHelperCurrent.NewOperatorHelmReleaseManager("current-operator", &values)
	g.Expect(err).ToNot(HaveOccurred())
	defer CleanupHelm(t, rmCurrent, helmHelperCurrent)

	// Delete the CRDs so that the previous version Operator installs the previous version CRDs
	t.Logf("Removing CRDs")
	err = helper.UninstallCrds(t)
	g.Expect(err).NotTo(HaveOccurred())

	// Install the previous Operator chart
	t.Logf("Installing previous Operator version %s", prevVersion)
	_, err = rmPrevious.InstallRelease()
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
	stsBefore, err := helper.WaitForStatefulSetForRole(cl, namespace, &cluster, cluster.GetFirstRole(), time.Second*10, time.Minute*5, t)
	g.Expect(err).ToNot(HaveOccurred())

	dir := t.Name() + "/Before"
	helper.DumpOperatorLog(framework.Global.KubeClient, namespace, dir, t)
	helper.DumpState(namespace, dir, t)

	// Upgrade to the current Operator - we do this by running cleanup to remove the previous operator and then install the new one
	t.Logf("Removing previous operator version %s", prevVersion)
	CleanupHelm(t, rmPrevious, helmHelperPrevious)
	t.Logf("Installing current operator")
	_, err = rmCurrent.InstallRelease()
	g.Expect(err).ToNot(HaveOccurred())

	// The chart is installed but the Pod may not exist yet so wait for it...
	// (we wait a maximum of 5 minutes, retrying every 10 seconds)
	pods, err = helper.WaitForOperatorPods(helmHelperCurrent.KubeClient, helmHelperCurrent.Namespace, time.Second*10, time.Minute*5)
	d, err := json.Marshal(pods[0])
	g.Expect(err).ToNot(HaveOccurred())
	version := helper.GetOperatorVersion()
	t.Logf("JSON for new Operator Pod version %s:\n%s", version, string(d))
	image := helper.GetOperatorImage()
	g.Expect(pods[0].Spec.Containers[0].Image).To(Equal(image))

	// wait for one minute and ensure that the Coherence Cluster did not restart
	time.Sleep(time.Minute * 1)

	// Get the cluster StatefulSet after Operator Upgrade
	stsAfter, err := helper.WaitForStatefulSetForRole(cl, namespace, &cluster, cluster.GetFirstRole(), time.Second*10, time.Minute*5, t)
	g.Expect(err).ToNot(HaveOccurred())

	// Assert that the StatefulSet has not been changed by the upgrade (i.e. its generation is unchanged)
	diffs := deep.Equal(stsBefore, stsAfter)
	t.Logf("Difference between StatefulSet for v%s and v%s:\n", prevVersion, version)
	for _, diff := range diffs {
		t.Log(diff)
	}
	g.Expect(stsAfter.Generation).To(Equal(stsBefore.Generation))

	// Ensure that everything is still linked, i.e. the cluster to the role to the CoherenceInternal
	// by upgrading the Coherence cluster. We're just going to add a label to the Pods

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
	err = f.Client.Update(context.TODO(), &cluster)
	g.Expect(err).ToNot(HaveOccurred())

	// wait for at least one Pod with the new label, this will verify that the update worked and the Coherence cluster is still good
	_, err = helper.WaitForPodsWithLabel(cl, namespace, "foo=bar", 1, time.Second*10, time.Minute*5)
	g.Expect(err).ToNot(HaveOccurred())
}
