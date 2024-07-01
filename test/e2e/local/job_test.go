/*
 * Copyright (c) 2020, 2024, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package local

import (
	"fmt"
	. "github.com/onsi/gomega"
	coh "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/ptr"
	"testing"
	"time"
)

// Test that a minimal CoherenceJob works
func TestMinimalJob(t *testing.T) {
	// Make sure we defer clean-up when we're done!!
	testContext.CleanupAfterTest(t)
	g := NewWithT(t)

	jobs, _ := helper.AssertCoherenceJobs(testContext, t, "job-minimal.yaml")

	_, ok := jobs["minimal-job"]
	g.Expect(ok).To(BeTrue(), "did not find expected 'minimal-job' deployment")
}

func TestJobWithSingleSuccessfulReplica(t *testing.T) {
	// Make sure we defer clean-up when we're done!!
	testContext.CleanupAfterTest(t)
	g := NewWithT(t)

	ns := helper.GetTestNamespace()
	name := "job-single-successful"

	pods := deployJob(t, ns, name, 1, 0)
	g.Expect(len(pods)).To(Equal(1))

	condition := helper.JobSucceededCondition(1)
	_, err := helper.WaitForCoherenceJobCondition(testContext, ns, name, condition, time.Second*10, time.Minute*5)
	g.Expect(err).NotTo(HaveOccurred())

	condition = helper.StatusPhaseCondition(coh.ConditionTypeCompleted)
	_, err = helper.WaitForCoherenceJobCondition(testContext, ns, name, condition, time.Second*10, time.Minute*5)
	g.Expect(err).NotTo(HaveOccurred())
}

func TestJobWithMultipleSuccessfulReplicas(t *testing.T) {
	// Make sure we defer clean-up when we're done!!
	testContext.CleanupAfterTest(t)
	g := NewWithT(t)

	ns := helper.GetTestNamespace()
	name := "job-multiple-successful"

	replicas := 3
	pods := deployJob(t, ns, name, int32(replicas), 0)
	g.Expect(len(pods)).To(Equal(replicas))

	condition := helper.JobSucceededCondition(3)
	_, err := helper.WaitForCoherenceJobCondition(testContext, ns, name, condition, time.Second*10, time.Minute*5)
	g.Expect(err).NotTo(HaveOccurred())

	condition = helper.StatusPhaseCondition(coh.ConditionTypeCompleted)
	_, err = helper.WaitForCoherenceJobCondition(testContext, ns, name, condition, time.Second*10, time.Minute*5)
	g.Expect(err).NotTo(HaveOccurred())
}

func TestJobWithSingleFailedReplica(t *testing.T) {
	// Make sure we defer clean-up when we're done!!
	testContext.CleanupAfterTest(t)
	g := NewWithT(t)

	ns := helper.GetTestNamespace()
	name := "job-single-failed"

	pods := deployJob(t, ns, name, 1, 1)
	g.Expect(len(pods)).To(Equal(1))

	condition := helper.JobFailedCondition(1)
	_, err := helper.WaitForCoherenceJobCondition(testContext, ns, name, condition, time.Second*10, time.Minute*5)
	g.Expect(err).NotTo(HaveOccurred())
}

func TestJobWithReadyAction(t *testing.T) {
	// Make sure we defer clean-up when we're done!!
	testContext.CleanupAfterTest(t)
	g := NewWithT(t)

	name := "job-with-ready"

	jobs, _ := helper.AssertCoherenceJobs(testContext, t, "job-with-ready-action.yaml")

	job, ok := jobs[name]
	g.Expect(ok).To(BeTrue(), fmt.Sprintf("did not find expected '%s' deployment", name))

	condition := jobProbesExecuted{count: int(job.GetReplicas())}
	_, err := helper.WaitForCoherenceJobCondition(testContext, job.Namespace, job.Name, condition, time.Second*10, time.Minute*5)
	g.Expect(err).NotTo(HaveOccurred())
}

func deployJob(t *testing.T, ns, name string, replicas, exitCode int32) []corev1.Pod {
	g := NewWithT(t)

	t.Logf("Deploying CoherenceJob %s in namespace %s", name, ns)

	jobs, err := helper.NewCoherenceJobFromYaml(ns, "job-test.yaml")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(len(jobs)).To(Equal(1))

	jobs[0].Name = name
	jobs[0].Spec.Replicas = ptr.To(replicas)

	ev := corev1.EnvVar{
		Name:  "APPLICATION_EXIT_CODE",
		Value: fmt.Sprintf("%d", exitCode),
	}
	jobs[0].Spec.Env = append(jobs[0].Spec.Env, ev)

	m, pods := helper.AssertCoherenceJobsSpec(testContext, t, jobs)

	_, ok := m[name]
	g.Expect(ok).To(BeTrue(), fmt.Sprintf("did not find expected '%s' deployment", name))

	t.Logf("Deployed CoherenceJob %s in namespace %s with %d pods", name, ns, len(pods))

	return pods
}

type jobProbesExecuted struct {
	count int
}

func (in jobProbesExecuted) Test(d coh.CoherenceResource) bool {
	status := d.GetStatus()
	if len(status.JobProbes) == 0 {
		return false
	}

	success := 0
	for _, s := range status.JobProbes {
		if s.Success != nil && *s.Success {
			success++
		}
	}
	return success == in.count
}

func (in jobProbesExecuted) String() string {
	return fmt.Sprintf("Job ready probes executed on %d pods", in.count)
}
