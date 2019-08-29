package remote

import (
	goctx "context"
	"fmt"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	cohv1 "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	"io/ioutil"
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

	scaler := func(t *testing.T, role *cohv1.CoherenceRole, replicas int32) error {
		role.Spec.SetReplicas(replicas)
		f := framework.Global
		return f.Client.Update(goctx.TODO(), role)
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			assertScale(t, tc.policy, tc.start, tc.end, scaler)
		})
	}
}

func TestScalingWithKubectl(t *testing.T) {

	scaler := func(t *testing.T, role *cohv1.CoherenceRole, replicas int32) error {
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

	assertScale(t, cohv1.ParallelUpSafeDownScaling, 1, 3, scaler)
}

// ----- helper methods ------------------------------------------------

type ScaleFunction func(t *testing.T, role *cohv1.CoherenceRole, replicas int32) error

// Assert that a cluster can be created and scaled using the specified policy.
func assertScale(t *testing.T, policy cohv1.ScalingPolicy, replicasStart, replicasScale int32, scaler ScaleFunction) {
	var (
		clusterName = "test-cluster"
		roleName    = "one"
	)
	g := NewGomegaWithT(t)

	ctx := helper.CreateTestContext(t)
	defer helper.DumpOperatorLogsAndCleanup(t, ctx)

	namespace, err := ctx.GetNamespace()
	g.Expect(err).NotTo(HaveOccurred())

	cluster, err := helper.NewCoherenceClusterFromYaml(namespace, "scaling-test.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	// Get the role and update it's replica count and scaling policy
	roleSpec := cluster.Spec.Roles[0]
	roleSpec.SetReplicas(replicasStart)
	roleSpec.ScalingPolicy = &policy
	// NOTE: we MUST set the role back into the role array because in the cluster
	// because in Go (unlike some other languages) we seem to have a COPY of what
	// is in the role array.
	cluster.Spec.Roles[0] = roleSpec

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
}

// installSimpleCluster installs a cluster and asserts that the underlying StatefulSet resources reach the correct state.
func installSimpleCluster(t *testing.T, ctx *framework.TestCtx, cluster cohv1.CoherenceCluster) {
	g := NewGomegaWithT(t)

	f := framework.Global

	err := f.Client.Create(goctx.TODO(), &cluster, helper.DefaultCleanup(ctx))
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
	_, err := helper.WaitForCoherenceRoleCondition(f, cluster.Namespace, fullName, condition, time.Second*10, time.Minute*5, t)
	g.Expect(err).NotTo(HaveOccurred())

	t.Logf("Asserting StatefulSet %s exists with %d replicas\n", fullName, replicas)

	// wait for the StatefulSet to have three replicas
	sts, err := helper.WaitForStatefulSet(f.KubeClient, cluster.Namespace, fullName, replicas, time.Second*10, time.Minute*5, t)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(sts.Status.ReadyReplicas).To(Equal(replicas))
}
