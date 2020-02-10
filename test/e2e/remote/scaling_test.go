/*
 * Copyright (c) 2019, 2020 Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package remote

import (
	goctx "context"
	"fmt"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	cohv1 "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	"golang.org/x/net/context"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/testing_frameworks/integration"
	"strings"
	"testing"
	"time"

	. "github.com/onsi/gomega"
)

// Test scaling up and down with different policies.
// This test is an example of using sub-tests to run the test with different test cases.
func TestScaling(t *testing.T) {
	testCases := []struct {
		testName string
		start    int32
		end      int32
		policy   cohv1.ScalingPolicy
	}{
		{"UpParallelScaling", 1, 3, cohv1.ParallelScaling},
		{"UpParallelUpSafeDownScaling", 1, 3, cohv1.ParallelUpSafeDownScaling},
		{"UpSafeScaling", 1, 3, cohv1.SafeScaling},
		{"DownParallelScaling", 3, 1, cohv1.ParallelScaling},
		{"DownParallelUpSafeDownScaling", 3, 1, cohv1.ParallelUpSafeDownScaling},
		{"DownSafeScaling", 3, 1, cohv1.SafeScaling},
	}

	for id, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			assertScale(t, id, tc.policy, tc.start, tc.end, roleScaler)
		})
	}
}

// Test that a role can be scaled up using the kubectl scale command
func TestScalingUpWithKubectl(t *testing.T) {
	assertScale(t, 10, cohv1.ParallelUpSafeDownScaling, 1, 3, kubeCtlRoleScaler)
}

// Test that a role can be scaled down using the kubectl scale command
func TestScalingDownWithKubectl(t *testing.T) {
	assertScale(t, 20, cohv1.ParallelUpSafeDownScaling, 3, 1, kubeCtlRoleScaler)
}

// If a role is scaled down to zero it should be deleted and just its parent CoherenceCluster should remain.
// This test scales down by directly updating the replica count in the role to zero.
func TestScaleDownToZero(t *testing.T) {
	assertScaleDownToZero(t, 30, roleScaler)
}

// If a role is scaled down to zero it should be deleted and just its parent CoherenceCluster should remain.
// This test scales down using the "kubectl scale --relicas=0" command
func TestScaleDownToZeroUsingKubectl(t *testing.T) {
	assertScaleDownToZero(t, 40, kubeCtlRoleScaler)
}

// ----- helper methods ------------------------------------------------

// ScaleFunction is a function that can scale a role up or down
type ScaleFunction func(t *testing.T, role *cohv1.CoherenceRole, replicas int32) error

// A scaler function that scales a role by directly updating it to have a set number of replicas
var roleScaler = func(t *testing.T, role *cohv1.CoherenceRole, replicas int32) error {
	role.Spec.SetReplicas(replicas)
	f := framework.Global
	return f.Client.Update(goctx.TODO(), role)
}

// A scaler function that scales a role using the kubectl scale command
var kubeCtlRoleScaler = func(t *testing.T, role *cohv1.CoherenceRole, replicas int32) error {
	versionArg := "--resource-version=" + role.ResourceVersion
	replicasArg := fmt.Sprintf("--replicas=%d", replicas)
	roleArg := "coherencerole/" + role.GetName()
	kubectl := integration.KubeCtl{}
	args := []string{"-n", role.GetNamespace(), "scale", replicasArg, versionArg, roleArg}

	t.Logf("Executing kubectl %s", strings.Join(args, " "))

	stdout, stderr, err := kubectl.Run(args...)
	o, _ := ioutil.ReadAll(stdout)
	t.Logf("kubectl scale stdout:\n%s\n", string(o))
	e, _ := ioutil.ReadAll(stderr)
	t.Logf("kubectl scale stderr:\n%s\n", string(e))
	return err
}

// Assert that a cluster can be created and scaled using the specified policy.
func assertScale(t *testing.T, id int, policy cohv1.ScalingPolicy, replicasStart, replicasScale int32, scaler ScaleFunction) {
	g := NewGomegaWithT(t)

	ctx := helper.CreateTestContext(t)
	defer helper.DumpOperatorLogsAndCleanup(t, ctx)

	namespace, err := ctx.GetNamespace()
	g.Expect(err).NotTo(HaveOccurred())

	cluster, err := helper.NewCoherenceClusterFromYaml(namespace, "scaling-test.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	//Give the cluster a unique name based on the test name
	cluster.SetName(fmt.Sprintf("%s-%d", cluster.GetName(), id))

	// Get the role and update it's replica count and scaling policy
	roleSpec := cluster.GetFirstRole()
	roleSpec.SetReplicas(replicasStart)

	if roleSpec.Scaling == nil {
		roleSpec.Scaling = &cohv1.ScalingSpec{}
	}
	roleSpec.Scaling.Policy = &policy

	// NOTE: we MUST set the role back into the role array because in the cluster
	// because in Go (unlike some other languages) we seem to have a COPY of what
	// is in the role array.
	cluster.Spec.Roles[0] = roleSpec

	clusterName := cluster.GetName()
	roleName := roleSpec.GetRoleName()

	// Do the canary test unless parallel scaling down
	doCanary := replicasStart < replicasScale || policy != cohv1.ParallelScaling

	f := framework.Global

	installSimpleCluster(t, ctx, cluster)

	if doCanary {
		t.Log("Initialising canary cache")
		err = helper.StartCanary(namespace, clusterName, roleName)
		g.Expect(err).NotTo(HaveOccurred())
	}

	role, err := helper.GetCoherenceRole(f, namespace, clusterName+"-"+roleName)
	g.Expect(err).NotTo(HaveOccurred())

	err = scaler(t, role, replicasScale)
	g.Expect(err).NotTo(HaveOccurred())

	assertRoleEventuallyInDesiredState(t, cluster, role.Spec, replicasScale)

	if doCanary {
		t.Log("Checking canary cache")
		err = helper.CheckCanary(namespace, clusterName, roleName)
		g.Expect(err).NotTo(HaveOccurred())
	}

	// The parent CoherenceCluster should have the correct replica count for the role
	cl := cohv1.CoherenceCluster{}
	err = f.Client.Get(context.TODO(), types.NamespacedName{Namespace: namespace, Name: cluster.Name}, &cl)
	g.Expect(err).NotTo(HaveOccurred())
	r := cl.GetRole(roleName)
	g.Expect(r.GetReplicas()).To(Equal(replicasScale))
}

func assertScaleDownToZero(t *testing.T, uid int, scaler ScaleFunction) {
	const (
		zero int32 = 0
	)

	g := NewGomegaWithT(t)

	ctx := helper.CreateTestContext(t)

	namespace, err := ctx.GetNamespace()
	g.Expect(err).NotTo(HaveOccurred())

	cluster, err := helper.NewCoherenceClusterFromYaml(namespace, "scaling-to-zero-test.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	defer cleanup(t, namespace, cluster.Name, ctx)

	//Give the cluster a unique name based on the test name
	cluster.SetName(fmt.Sprintf("%s-%d", cluster.GetName(), uid))

	// Get the role and update it's replica count and scaling policy
	roleSpec := cluster.GetFirstRole()
	roleFullName := cluster.GetFullRoleName(roleSpec.GetRoleName())

	installSimpleCluster(t, ctx, cluster)

	f := framework.Global
	role, err := helper.GetCoherenceRole(f, namespace, roleFullName)
	g.Expect(err).NotTo(HaveOccurred())

	// Scale the role down to zero
	err = scaler(t, role, zero)
	g.Expect(err).NotTo(HaveOccurred())

	// Wait for deletion of the CoherenceInternal
	u := helper.NewUnstructuredCoherenceInternal()
	err = helper.WaitForDeletion(f, namespace, roleFullName, &u, time.Second*5, time.Minute*5)
	g.Expect(err).NotTo(HaveOccurred())

	// The CoherenceCluster should still exist
	cl := cohv1.CoherenceCluster{}
	err = f.Client.Get(context.TODO(), types.NamespacedName{Namespace: namespace, Name: cluster.Name}, &cl)
	g.Expect(err).NotTo(HaveOccurred())
	// The replica count for the role spec in the cluster should be zero
	r := cl.GetRole(roleSpec.GetRoleName())
	g.Expect(r.GetReplicas()).To(Equal(zero))

	// wait for the role to match the condition
	fullName := r.GetFullRoleName(&cluster)
	condition := helper.ReplicasRoleCondition(0)
	_, err = helper.WaitForCoherenceRoleCondition(f, cluster.Namespace, fullName, condition, time.Second*10, time.Minute*5, t)
	g.Expect(err).NotTo(HaveOccurred())
}
