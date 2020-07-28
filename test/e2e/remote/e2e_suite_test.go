/*
 * Copyright (c) 2019, 2020 Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package remote

import (
	"context"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	cohv1 "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	"k8s.io/apimachinery/pkg/types"
	"testing"
	"time"

	. "github.com/onsi/gomega"
)

// Operator SDK test suite entry point
func TestMain(m *testing.M) {
	framework.MainEntry(m)
}

func cleanup(t *testing.T, ctx *framework.Context, names ...types.NamespacedName) {
	helper.DumpOperatorLogs(t)
	for _, name := range names {
		deleteDeployment(name.Namespace, name.Name)
	}
	ctx.Cleanup()
}

// deleteDeployment deletes a deployment.
func deleteDeployment(namespace, name string) {
	deployment := cohv1.Coherence{}
	f := framework.Global

	err := f.Client.Get(context.TODO(), types.NamespacedName{Namespace: namespace, Name: name}, &deployment)

	if err == nil {
		_ = f.Client.Delete(context.TODO(), &deployment)
	}
}

// installSimpleDeployment installs a deployment and asserts that the underlying StatefulSet resources reach the correct state.
func installSimpleDeployment(t *testing.T, ctx *framework.Context, d cohv1.Coherence) {
	g := NewGomegaWithT(t)
	f := framework.Global
	err := f.Client.Create(context.TODO(), &d, helper.DefaultCleanup(ctx))
	g.Expect(err).NotTo(HaveOccurred())
	assertDeploymentEventuallyInDesiredState(t, d, d.GetReplicas())
}

// assertDeploymentEventuallyInDesiredState asserts that a Coherence resource exists and has the correct spec and that the
// underlying StatefulSet exists with the correct status and ready replicas.
func assertDeploymentEventuallyInDesiredState(t *testing.T, d cohv1.Coherence, replicas int32) {
	g := NewGomegaWithT(t)
	f := framework.Global

	t.Logf("Asserting Coherence resource %s exists with %d replicas\n", d.Name, replicas)

	// create a DeploymentStateCondition that checks a deployment's replica count
	condition := helper.ReplicaCountCondition(replicas)

	// wait for the deployment to match the condition
	_, err := helper.WaitForCoherenceCondition(f, d.Namespace, d.Name, condition, time.Second*10, time.Minute*5, t)
	g.Expect(err).NotTo(HaveOccurred())

	t.Logf("Asserting StatefulSet %s exists with %d replicas\n", d.Name, replicas)

	// wait for the StatefulSet to have the required ready replicas
	sts, err := helper.WaitForStatefulSet(f.KubeClient, d.Namespace, d.Name, replicas, time.Second*10, time.Minute*5)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(sts.Status.ReadyReplicas).To(Equal(replicas))

	t.Logf("Asserting StatefulSet %s exist with %d replicas - Done!\n", d.Name, replicas)
}
