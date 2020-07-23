/*
 * Copyright (c) 2019, 2020, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package local

import (
	"context"
	. "github.com/onsi/gomega"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	coh "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	"testing"
	"time"
)

func TestRollingUpgrade(t *testing.T) {
	ctx := helper.CreateTestContext(t)
	defer helper.DumpOperatorLogsAndCleanup(t, ctx)

	g := NewGomegaWithT(t)
	f := framework.Global

	namespace, err := ctx.GetWatchNamespace()
	g.Expect(err).NotTo(HaveOccurred())

	t.Log("Deploying initial version of Coherence cluster")
	// Do the initial deployment
	deployments, _ := AssertDeploymentsWithContext(t, ctx, "rolling-upgrade.yaml")
	// Get the expected single deployment from the returned map
	deployment, ok := deployments["rolling-cluster"]
	g.Expect(ok).To(BeTrue())

	// Get the latest state for the deployment
	upgrade := coh.Coherence{}
	err = f.Client.Get(context.TODO(), deployment.GetNamespacedName(), &upgrade)
	g.Expect(err).NotTo(HaveOccurred())

	// Upgrade the version label and JVM Heap
	updatedHeap := "500m"
	t.Log("Deploying updated version of Coherence cluster")
	upgrade.Spec.Labels["version"] = "two"
	upgrade.Spec.JVM.Memory.HeapSize = &updatedHeap
	err = f.Client.Update(context.TODO(), &upgrade)
	g.Expect(err).NotTo(HaveOccurred())

	// wait for the expected updated Pods
	t.Log("Waiting for all Pods to be updated")
	pods, err := helper.WaitForPodsWithLabel(f.KubeClient, namespace, "version=two", 3, time.Second*10, time.Minute*5)
	g.Expect(err).NotTo(HaveOccurred())

	for _, pod := range pods {
		for _, ev := range pod.Spec.Containers[0].Env {
			if ev.Name == coh.EnvVarJvmMemoryHeap {
				g.Expect(ev.Value).To(Equal(updatedHeap), "Expected heap incorrect for Pod "+pod.Name)
			}
		}
	}
}
