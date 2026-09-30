/*
 * Copyright (c) 2026, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package remote

import (
	"fmt"
	"testing"
	"time"

	. "github.com/onsi/gomega"
	coh "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/test/e2e/helper"
)

func TestJobWithReadyHTTPAction(t *testing.T) {
	testContext.CleanupAfterTest(t)
	g := NewWithT(t)

	name := "job-with-ready-http"
	jobs, _ := helper.AssertCoherenceJobs(testContext, t, "job-with-ready-http-action.yaml")

	job, ok := jobs[name]
	g.Expect(ok).To(BeTrue(), fmt.Sprintf("did not find expected '%s' deployment", name))

	condition := remoteJobProbesExecuted{count: int(job.GetReplicas())}
	_, err := helper.WaitForCoherenceJobCondition(testContext, job.Namespace, job.Name, condition, time.Second*10, time.Minute*5)
	g.Expect(err).NotTo(HaveOccurred())
}

type remoteJobProbesExecuted struct {
	count int
}

func (in remoteJobProbesExecuted) Test(d coh.CoherenceResource) bool {
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

func (in remoteJobProbesExecuted) String() string {
	return fmt.Sprintf("Job HTTP ready actions executed on %d pods", in.count)
}
