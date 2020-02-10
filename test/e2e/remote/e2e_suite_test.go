/*
 * Copyright (c) 2019, 2020 Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package remote

import (
	"context"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	cohv1 "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	"k8s.io/apimachinery/pkg/types"
	"testing"
	"time"

	. "github.com/onsi/gomega"
	f "github.com/operator-framework/operator-sdk/pkg/test"
)

// Operator SDK test suite entry point
func TestMain(m *testing.M) {
	f.MainEntry(m)
}

func cleanup(t *testing.T, namespace, clusterName string, ctx *framework.TestCtx) {
	helper.DumpOperatorLogs(t, ctx)
	deleteCluster(namespace, clusterName)
	ctx.Cleanup()
}

// deleteCluster deletes a cluster.
func deleteCluster(namespace, name string) {
	cluster := cohv1.CoherenceCluster{}
	f := f.Global

	err := f.Client.Get(context.TODO(), types.NamespacedName{Namespace: namespace, Name: name}, &cluster)

	if err == nil {
		_ = f.Client.Delete(context.TODO(), &cluster)
	}
}

// installSimpleCluster installs a cluster and asserts that the underlying StatefulSet resources reach the correct state.
func installSimpleCluster(t *testing.T, ctx *framework.TestCtx, cluster cohv1.CoherenceCluster) {
	g := NewGomegaWithT(t)

	f := framework.Global

	err := f.Client.Create(context.TODO(), &cluster, helper.DefaultCleanup(ctx))
	g.Expect(err).NotTo(HaveOccurred())

	if len(cluster.Spec.Roles) > 0 {
		for _, r := range cluster.Spec.Roles {
			assertRoleEventuallyInDesiredState(t, cluster, r, r.GetReplicas())
		}
	} else {
		r := cluster.Spec.CoherenceRoleSpec
		assertRoleEventuallyInDesiredState(t, cluster, r, r.GetReplicas())
	}
}

// assertRoleEventuallyInDesiredState asserts that a CoherenceRole exists and has the correct spec and that the
// underlying StatefulSet exists with the correct status and ready replicas.
func assertRoleEventuallyInDesiredState(t *testing.T, cluster cohv1.CoherenceCluster, r cohv1.CoherenceRoleSpec, replicas int32) {
	g := NewGomegaWithT(t)
	f := framework.Global
	fullName := r.GetFullRoleName(&cluster)

	t.Logf("Asserting CoherenceRole %s exists\n", fullName)

	t.Logf("Asserting CoherenceRole %s exists with %d replicas\n", fullName, replicas)

	// create a RoleStateCondition that checks a role's replica count
	condition := helper.ReplicasRoleCondition(replicas)

	// wait for the role to match the condition
	_, err := helper.WaitForCoherenceRoleCondition(f, cluster.Namespace, fullName, condition, time.Second*10, time.Minute*10, t)
	g.Expect(err).NotTo(HaveOccurred())

	t.Logf("Asserting StatefulSet %s exists with %d replicas\n", fullName, replicas)

	// wait for the StatefulSet to have the required ready replicas
	sts, err := helper.WaitForStatefulSet(f.KubeClient, cluster.Namespace, fullName, replicas, time.Second*10, time.Minute*10, t)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(sts.Status.ReadyReplicas).To(Equal(replicas))

	t.Logf("Asserting StatefulSet %s exist with %d replicas - Done!\n", fullName, replicas)
}
