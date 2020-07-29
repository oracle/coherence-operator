/*
 * Copyright (c) 2019, 2020 Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package remote

import (
	goctx "context"
	"fmt"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	cohv1 "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	"golang.org/x/net/context"
	"io/ioutil"
	appsv1 "k8s.io/api/apps/v1"
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
			assertScale(t, id, tc.policy, tc.start, tc.end, deploymentScaler)
		})
	}
}

// Test that a deployment can be scaled up using the kubectl scale command
func TestScalingUpWithKubectl(t *testing.T) {
	assertScale(t, 10, cohv1.ParallelUpSafeDownScaling, 1, 3, kubeCtlScaler)
}

// Test that a deployment can be scaled down using the kubectl scale command
func TestScalingDownWithKubectl(t *testing.T) {
	assertScale(t, 20, cohv1.ParallelUpSafeDownScaling, 3, 1, kubeCtlScaler)
}

// If a deployment is scaled down to zero it should be deleted and just its parent Coherence resource should remain.
// This test scales down by directly updating the replica count in the deployment to zero.
func TestScaleDownToZero(t *testing.T) {
	assertScaleDownToZero(t, 30, deploymentScaler)
}

// If a deployment is scaled down to zero it should be deleted and just its parent Coherence resource should remain.
// This test scales down using the "kubectl scale --relicas=0" command
func TestScaleDownToZeroUsingKubectl(t *testing.T) {
	assertScaleDownToZero(t, 40, kubeCtlScaler)
}

// ----- helper methods ------------------------------------------------

// ScaleFunction is a function that can scale a deployment up or down
type ScaleFunction func(t *testing.T, d *cohv1.Coherence, replicas int32) error

// A scaler function that scales a deployment by directly updating it to have a set number of replicas
var deploymentScaler = func(t *testing.T, d *cohv1.Coherence, replicas int32) error {
	f := framework.Global
	current, err := helper.GetCoherence(f, d.Namespace, d.Name)
	if err != nil {
		return err
	}
	current.Spec.SetReplicas(replicas)
	t.Logf("Scaling %s to %d", current.Name, replicas)
	return f.Client.Update(goctx.TODO(), current)
}

// A scaler function that scales a deployment using the kubectl scale command
var kubeCtlScaler = func(t *testing.T, d *cohv1.Coherence, replicas int32) error {
	f := framework.Global
	current, err := helper.GetCoherence(f, d.Namespace, d.Name)
	if err != nil {
		return err
	}

	versionArg := "--resource-version=" + current.ResourceVersion
	replicasArg := fmt.Sprintf("--replicas=%d", replicas)
	deploymentArg := "coherence/" + current.GetName()
	kubectl := integration.KubeCtl{}
	args := []string{"-n", current.GetNamespace(), "scale", replicasArg, versionArg, deploymentArg}

	t.Logf("Executing kubectl %s", strings.Join(args, " "))

	stdout, stderr, err := kubectl.Run(args...)
	o, _ := ioutil.ReadAll(stdout)
	t.Logf("kubectl scale stdout:\n%s\n", string(o))
	e, _ := ioutil.ReadAll(stderr)
	t.Logf("kubectl scale stderr:\n%s\n", string(e))
	return err
}

// Assert that a deployment can be created and scaled using the specified policy.
func assertScale(t *testing.T, id int, policy cohv1.ScalingPolicy, replicasStart, replicasScale int32, scaler ScaleFunction) {
	g := NewGomegaWithT(t)

	t.Log("assertScale() - Starting...")
	ctx := helper.CreateTestContext(t)
	defer helper.DumpOperatorLogs(t)

	namespace, err := ctx.GetWatchNamespace()
	g.Expect(err).NotTo(HaveOccurred())

	deployment, err := helper.NewSingleCoherenceFromYaml(namespace, "scaling-test.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	//Give the deployment a unique name based on the test name
	deployment.SetName(fmt.Sprintf("%s-%d", deployment.GetName(), id))

	// update the replica count and scaling policy
	deployment.SetReplicas(replicasStart)

	if deployment.Spec.Scaling == nil {
		deployment.Spec.Scaling = &cohv1.ScalingSpec{}
	}
	deployment.Spec.Scaling.Policy = &policy

	// Do the canary test unless parallel scaling down
	doCanary := replicasStart < replicasScale || policy != cohv1.ParallelScaling

	t.Logf("assertScale() - doCanary=%t", doCanary)
	t.Log("assertScale() - Installing Coherence deployment...")
	installSimpleDeployment(t, ctx, deployment)
	t.Log("assertScale() - Installed Coherence deployment")

	if doCanary {
		t.Log("Initialising canary cache")
		err = helper.StartCanary(testContext, namespace, deployment.Name)
		g.Expect(err).NotTo(HaveOccurred())
	}

	// Get the current deployment state so that we can scale it
	err = scaler(t, &deployment, replicasScale)
	g.Expect(err).NotTo(HaveOccurred())

	assertDeploymentEventuallyInDesiredState(t, deployment, replicasScale)

	if doCanary {
		t.Log("Checking canary cache")
		err = helper.CheckCanary(testContext, namespace, deployment.Name)
		g.Expect(err).NotTo(HaveOccurred())
	}
}

func assertScaleDownToZero(t *testing.T, uid int, scaler ScaleFunction) {
	const (
		zero int32 = 0
	)

	g := NewGomegaWithT(t)
	ctx := helper.CreateTestContext(t)

	namespace, err := ctx.GetWatchNamespace()
	g.Expect(err).NotTo(HaveOccurred())

	deployment, err := helper.NewSingleCoherenceFromYaml(namespace, "scaling-to-zero-test.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	defer cleanup(t, ctx, deployment.GetNamespacedName())

	//Give the deployment a unique name based on the test name
	deployment.SetName(fmt.Sprintf("%s-%d", deployment.GetName(), uid))

	installSimpleDeployment(t, ctx, deployment)

	f := framework.Global

	// Scale the deployment down to zero
	err = scaler(t, &deployment, zero)
	g.Expect(err).NotTo(HaveOccurred())

	// Wait for deletion of the StatefulSet
	sts := appsv1.StatefulSet{}
	err = helper.WaitForDeletion(f, namespace, deployment.Name, &sts, time.Second*5, time.Minute*5)
	g.Expect(err).NotTo(HaveOccurred())

	// The Coherence resource should still exist
	updated := cohv1.Coherence{}
	err = f.Client.Get(context.TODO(), deployment.GetNamespacedName(), &updated)
	g.Expect(err).NotTo(HaveOccurred())
	// The replica count for the deployment spec in the deployment should be zero
	g.Expect(updated.GetReplicas()).To(Equal(zero))

	// wait for the deployment to match the condition
	condition := helper.ReplicaCountCondition(0)
	_, err = helper.WaitForCoherenceCondition(f, namespace, deployment.Name, condition, time.Second*10, time.Minute*5)
	g.Expect(err).NotTo(HaveOccurred())
}
