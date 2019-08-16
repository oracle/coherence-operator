package remote

import (
	goctx "context"
	"fmt"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	coherence "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
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

	testImage := os.Getenv("TEST_USER_IMAGE")
	if testImage == "" {
		t.Fatal("The TEST_USER_IMAGE environment variable must point to a valid Coherence test image")
	}

	// Cannot execute this test using a local operator because safe scaling requires
	// the operator to make ReST calls to Pods which it can only do when properly
	// deployed into the k8s cluster.
	g.Expect(framework.Global.LocalOperator).To(BeFalse())

	ctx := helper.CreateTestContext(t)
	defer cleanup(t, ctx)

	namespace, err := ctx.GetNamespace()
	g.Expect(err).NotTo(HaveOccurred())

	artifacts := &coherence.UserArtifactsImageSpec{ImageSpec: coherence.ImageSpec{Image: &testImage}}
	config := "test-cache-config.xml"
	main := "com.oracle.coherence.k8s.testing.RestServer"

	roleOne := coherence.CoherenceRoleSpec{
		Role:          roleName,
		Replicas:      &replicasStart,
		ScalingPolicy: &policy,
		Images:        &coherence.Images{UserArtifacts: artifacts},
		Ports:         map[string]int32{"rest": 8080},
		CacheConfig:   &config,
		Main:          &coherence.MainSpec{Class: &main},
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

func cleanup(t *testing.T, ctx *framework.TestCtx) {
	namespace, err := ctx.GetNamespace()
	if err == nil {
		helper.DumpOperatorLog(namespace, t.Name(), t)
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
