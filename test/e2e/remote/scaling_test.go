package remote

import (
	goctx "context"
	"fmt"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	coherence "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
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
	doCanary := replicasStart < replicasScale || policy != coherence.ParallelScaling

	f := framework.Global

	installSimpleCluster(t, ctx, cluster)

	if doCanary {
		t.Log("Initialising canary cache")
		err = startCanary(namespace, clusterName, roleName)
		g.Expect(err).NotTo(HaveOccurred())
	}

	role, err := helper.GetCoherenceRole(f, namespace, clusterName+"-"+roleName)
	g.Expect(err).NotTo(HaveOccurred())

	role.Spec.SetReplicas(replicasScale)
	err = f.Client.Update(goctx.TODO(), role)
	g.Expect(err).NotTo(HaveOccurred())

	assertRoleEventuallyInDesiredState(t, cluster, role.Spec)

	if doCanary {
		t.Log("Checking canary cache")
		err = checkCanary(namespace, clusterName, roleName)
		g.Expect(err).NotTo(HaveOccurred())
	}
}

// installSimpleCluster installs a cluster and asserts that the underlying StatefulSet resources reach the correct state.
func installSimpleCluster(t *testing.T, ctx *framework.TestCtx, cluster coherence.CoherenceCluster) {
	g := NewGomegaWithT(t)

	f := framework.Global

	err := f.Client.Create(goctx.TODO(), &cluster, helper.DefaultCleanup(ctx))
	g.Expect(err).NotTo(HaveOccurred())

	if len(cluster.Spec.Roles) > 0 {
		for _, r := range cluster.Spec.Roles {
			assertRoleEventuallyInDesiredState(t, cluster, r)
		}
	} else {
		assertRoleEventuallyInDesiredState(t, cluster, cluster.Spec.CoherenceRoleSpec)
	}
}

// assertRoleEventuallyInDesiredState asserts that a CoherenceRole exists and has the correct spec and that the
// underlying StatefulSet exists with the correct status and ready replicas.
func assertRoleEventuallyInDesiredState(t *testing.T, cluster coherence.CoherenceCluster, r coherence.CoherenceRoleSpec) {
	g := NewGomegaWithT(t)
	f := framework.Global
	fullName := r.GetFullRoleName(&cluster)

	t.Logf("Asserting CoherenceRole %s exists\n", fullName)

	role, err := helper.WaitForCoherenceRole(f, cluster.Namespace, fullName, helper.RetryInterval, helper.Timeout, t)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(role.Spec.GetRoleName()).To(Equal(r.GetRoleName()))
	g.Expect(role.Spec.GetReplicas()).To(Equal(r.GetReplicas()))

	replicas := r.GetReplicas()

	t.Logf("Asserting StatefulSet %s exists with %d replicas\n", fullName, replicas)

	sts, err := helper.WaitForStatefulSetForRole(f.KubeClient, cluster.Namespace, &cluster, role.Spec, time.Second*10, time.Minute*5, t)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(sts.Status.ReadyReplicas).To(Equal(replicas))
}

// Initialise the canary test in the role being scaled.
func startCanary(namespace, clusterName, roleName string) error {
	return canary(namespace, clusterName, roleName, "canaryStart")
}

// Invoke the canary test in the role being scaled.
func checkCanary(namespace, clusterName, roleName string) error {
	return canary(namespace, clusterName, roleName, "canaryCheck")
}

// Make a canary ReST PUT call to Pod zero of the role.
func canary(namespace, clusterName, roleName, endpoint string) error {
	podName := fmt.Sprintf("%s-%s-0", clusterName, roleName)
	f := framework.Global

	pod, err := f.KubeClient.CoreV1().Pods(namespace).Get(podName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	forwarder, ports, err := helper.StartPortForwarderForPod(pod)
	if err != nil {
		return err
	}

	defer forwarder.Close()

	url := fmt.Sprintf("http://127.0.0.1:%d/%s", ports["rest"], endpoint)
	client := &http.Client{}
	request, err := http.NewRequest(http.MethodPut, url, strings.NewReader(""))
	if err != nil {
		return err
	}

	request.ContentLength = 0
	resp, err := client.Do(request)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("expected http response %d but received %d from '%s'", http.StatusOK, resp.StatusCode, url)
	}

	return nil
}
