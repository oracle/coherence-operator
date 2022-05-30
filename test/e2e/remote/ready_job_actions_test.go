/*
 * Copyright (c) 2021, 2022, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package remote

import (
	. "github.com/onsi/gomega"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	"testing"
	"time"
)

func TestExecuteJobActionWhenReady(t *testing.T) {
	g := NewWithT(t)
	testContext.CleanupAfterTest(t)

	deployment, _ := helper.AssertDeployments(testContext, t, "deployment-execute-job-action.yaml")
	data, ok := deployment["action-test"]
	g.Expect(ok).To(BeTrue(), "did not find expected 'action-test' deployment")

	_, err := helper.WaitForDeploymentReady(testContext, data.Namespace, data.Name, time.Second*5, time.Minute*5)
	g.Expect(err).NotTo(HaveOccurred())

	// find action job
	jobs, err := helper.WaitForJobsWithLabel(testContext, data.Namespace, "test=actions", 1, time.Second*5, time.Minute*2)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(jobs).To(HaveLen(1))
	firstJob := jobs[0]

	// find action pod
	pods, err := helper.WaitForPodsWithLabel(testContext, data.Namespace, "job-name="+jobs[0].Name, 1, time.Second*5, time.Minute*2)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(pods).To(HaveLen(1))

	// scale to 3
	err = deploymentScaler(t, &data, 3)
	g.Expect(err).NotTo(HaveOccurred())
	assertDeploymentEventuallyInDesiredState(t, data, 3)

	// no new action jobs
	jobs, err = helper.WaitForJobsWithLabel(testContext, data.Namespace, "test=actions", 2, time.Second*5, time.Minute*1)
	g.Expect(err).To(HaveOccurred())
	g.Expect(jobs).To(HaveLen(1))
	g.Expect(jobs[0].Name).Should(Equal(firstJob.Name))

	// scale to 0
	err = deploymentScaler(t, &data, 0)
	g.Expect(err).NotTo(HaveOccurred())

	// Wait for deletion of the StatefulSet
	sts := appsv1.StatefulSet{}
	err = helper.WaitForDeletion(testContext, data.Namespace, data.Name, &sts, time.Second*5, time.Minute*5)
	g.Expect(err).NotTo(HaveOccurred())

	// no new action jobs
	jobs, err = helper.WaitForJobsWithLabel(testContext, data.Namespace, "test=actions", 2, time.Second*5, time.Minute*2)
	g.Expect(err).To(HaveOccurred())
	g.Expect(jobs).To(HaveLen(1)) // COMPLETED one
	g.Expect(jobs[0].Name).Should(Equal(firstJob.Name))

	// scale to 2
	err = deploymentScaler(t, &data, 2)
	g.Expect(err).NotTo(HaveOccurred())
	assertDeploymentEventuallyInDesiredState(t, data, 2)

	// find action jobs
	jobs, err = helper.WaitForJobsWithLabel(testContext, data.Namespace, "test=actions", 2, time.Second*5, time.Minute*2)
	g.Expect(err).NotTo(HaveOccurred())
	count := len(jobs)
	g.Expect(count).To(Equal(2))
	g.Expect([]string{jobs[0].Name, jobs[1].Name}).Should(ContainElement(firstJob.Name))
	var secondJob batchv1.Job
	if jobs[0].Name == firstJob.Name {
		secondJob = jobs[1]
	} else {
		secondJob = jobs[0]
	}

	// find action pod
	pods, err = helper.WaitForPodsWithLabel(testContext, data.Namespace, "job-name="+secondJob.Name, 1, time.Second*5, time.Minute*2)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(pods).To(HaveLen(1))
}
