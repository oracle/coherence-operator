package remote

import (
	goctx "context"
	"fmt"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	coherence "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	"io"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"strings"
	"testing"
	"time"

	. "github.com/onsi/gomega"
)

// Test scaling up and down with different policies.
// This test is an example of using sub-tests to run the test with different test cases.
func TestScaling(t *testing.T) {
	testCases := []struct {
		start  int32
		end    int32
		policy coherence.ScalingPolicy
	}{
		{1, 3, coherence.ParallelScaling},           // scale up
		{1, 3, coherence.ParallelUpSafeDownScaling}, // scale up
		{1, 3, coherence.SafeScaling},               // scale up
		{3, 1, coherence.ParallelScaling},           // scale down
		{3, 1, coherence.ParallelUpSafeDownScaling}, // scale down
		{3, 1, coherence.SafeScaling},               // scale down
	}

	for _, tc := range testCases {
		var dir string
		if tc.start > tc.end {
			dir = "Down"
		} else {
			dir = "Up"
		}

		name := fmt.Sprintf("%s from %d to %d with policy %s", dir, tc.start, tc.end, tc.policy)

		t.Run(name, func(t *testing.T) {
			assertScale(t, tc.policy, tc.start, tc.end)
		})
	}
}

// ----- helper methods ------------------------------------------------

// Assert that a cluster can be created and scaled using the specified policy.
func assertScale(t *testing.T, policy coherence.ScalingPolicy, replicasStart, replicasScale int32) {
	var (
		clusterName = "test-cluster"
		roleName    = "one"
	)
	g := NewGomegaWithT(t)

	// Cannot execute this test using a local operator because safe scaling requires
	// the operator to make ReST calls to Pods which it can only do when properly
	// deployed into the k8s cluster.
	g.Expect(framework.Global.LocalOperator).To(BeFalse())

	ctx := helper.CreateTestContext(t)
	defer cleanup(t, ctx)

	namespace, err := ctx.GetNamespace()
	g.Expect(err).NotTo(HaveOccurred())

	roleOne := coherence.CoherenceRoleSpec{
		Role:          roleName,
		Replicas:      &replicasStart,
		ScalingPolicy: &policy,
	}

	cluster := coherence.CoherenceCluster{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      clusterName,
		},
		Spec: coherence.CoherenceClusterSpec{
			CoherenceRoleSpec: coherence.CoherenceRoleSpec{
				ReadinessProbe: helper.Readiness,
			},
			Roles: []coherence.CoherenceRoleSpec{roleOne},
		},
	}

	installSimpleCluster(t, ctx, cluster)

	f := framework.Global
	role, err := helper.GetCoherenceRole(f, namespace, clusterName+"-"+roleName)
	g.Expect(err).NotTo(HaveOccurred())

	role.Spec.SetReplicas(replicasScale)
	err = f.Client.Update(goctx.TODO(), role)
	g.Expect(err).NotTo(HaveOccurred())

	assertRole(t, cluster, role.Spec)

}

func cleanup(t *testing.T, ctx *framework.TestCtx) {
	namespace, err := ctx.GetNamespace()
	if err == nil {
		dumpOperatorLog(t, namespace)
	}
	ctx.Cleanup()
}

// installSimpleCluster installs a cluster and asserts that the underlying StatefulSet resources reach the correct state.
func installSimpleCluster(t *testing.T, ctx *framework.TestCtx, cluster coherence.CoherenceCluster) {
	g := NewGomegaWithT(t)

	f := framework.Global

	err := f.Client.Create(goctx.TODO(), &cluster, helper.DefaultCleanup(ctx))
	g.Expect(err).NotTo(HaveOccurred())

	t.Logf("Installing CoherenceCluster %s:\n%#v\n", cluster.Name, cluster)

	if len(cluster.Spec.Roles) > 0 {
		for _, r := range cluster.Spec.Roles {
			assertRole(t, cluster, r)
		}
	} else {
		assertRole(t, cluster, cluster.Spec.CoherenceRoleSpec)
	}
}

// assertRole asserts that a CoherenceRole exists and has the correct spec and that the underlying StatefulSet
// exists with the correct status and ready replicas.
func assertRole(t *testing.T, cluster coherence.CoherenceCluster, r coherence.CoherenceRoleSpec) {
	g := NewGomegaWithT(t)
	f := framework.Global
	fullName := r.GetFullRoleName(&cluster)

	t.Logf("Asserting CoherenceRole %s exists\n", fullName)

	role, err := helper.WaitForCoherenceRole(t, f, cluster.Namespace, fullName, helper.RetryInterval, helper.Timeout)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(role.Spec.GetRoleName()).To(Equal(r.GetRoleName()))
	g.Expect(role.Spec.GetReplicas()).To(Equal(r.GetReplicas()))

	replicas := r.GetReplicas()

	t.Logf("Asserting StatefulSet %s exists with %d replicas\n", fullName, replicas)

	sts, err := helper.WaitForStatefulSet(t, f.KubeClient, cluster.Namespace, fullName, replicas, time.Second*10, time.Minute*5)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(sts.Status.ReadyReplicas).To(Equal(replicas))
}

// Dump the Operator Pod log to a file.
func dumpOperatorLog(t *testing.T, namespace string) {
	f := framework.Global

	t.Log("Dumping Operator log for test " + t.Name())

	logs := os.Getenv("TEST_LOGS")
	if logs == "" {
		t.Log("Cannot capture Operator logs as log folder env var TEST_LOGS is not set")
		return
	}

	list, err := f.KubeClient.CoreV1().Pods(namespace).List(metav1.ListOptions{LabelSelector: "name=coherence-operator"})
	if err == nil {
		if len(list.Items) > 0 {
			logOpts := corev1.PodLogOptions{}
			pod := list.Items[0]

			res := f.KubeClient.CoreV1().Pods(namespace).GetLogs(pod.Name, &logOpts)
			s, err := res.Stream()
			if err == nil {
				name := logs + "/" + strings.ReplaceAll(t.Name(), "/", "_")
				err = os.MkdirAll(name, os.ModePerm)
				if err == nil {
					out, err := os.Create(name + "/operator.log")
					if err == nil {
						_, err = io.Copy(out, s)
					}
				}
			}
		} else {
			t.Log("Could not capture Operator Pod log. No Pods found.")
		}
	}

	if err != nil {
		t.Logf("Could not capture Operator Pod log due to error: %s\n", err.Error())
	}
}
