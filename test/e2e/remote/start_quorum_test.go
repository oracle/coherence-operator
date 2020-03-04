/*
 * Copyright (c) 2019, 2020 Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package remote

import (
	"context"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	corev1 "k8s.io/api/core/v1"
	"testing"
	"time"

	. "github.com/onsi/gomega"
)

func TestStartQuorumDependentRoleReady(t *testing.T) {
	assertCluster(t, "cluster-ready-quorum.yaml")
}

func assertCluster(t *testing.T, yamlFile string) {

	g := NewGomegaWithT(t)

	ctx := helper.CreateTestContext(t)
	namespace, err := ctx.GetNamespace()
	g.Expect(err).NotTo(HaveOccurred())

	startTime := time.Now()
	cluster, err := helper.NewCoherenceClusterFromYaml(namespace, yamlFile)
	g.Expect(err).NotTo(HaveOccurred())

	defer cleanup(t, namespace, cluster.Name, ctx)

	f := framework.Global

	err = f.Client.Create(context.TODO(), &cluster, helper.DefaultCleanup(ctx))
	g.Expect(err).NotTo(HaveOccurred())

	roleTimes := PodTimes{Scheduled: startTime, Ready: startTime}
	testTimes := PodTimes{Scheduled: startTime, Ready: startTime}

	// Assert roles start (assert in the order specified by the roles array, which is the expected start order)
	for _, r := range cluster.Spec.Roles {
		fullName := r.GetFullRoleName(&cluster)
		_, err := helper.WaitForStatefulSet(f.KubeClient, cluster.Namespace, fullName, r.GetReplicas(), time.Second*10, time.Minute*2, t)
		g.Expect(err).NotTo(HaveOccurred())

		pods, err := helper.ListCoherencePodsForCluster(f.KubeClient, namespace, cluster.Name)
		g.Expect(err).NotTo(HaveOccurred())

		t := PodTimes{Scheduled: startTime, Ready: startTime}

		for _, pod := range pods {
			for _, c := range pod.Status.Conditions {
				if c.Type == corev1.PodReady {
					if c.LastTransitionTime.After(t.Ready) {
						t.Ready = c.LastTransitionTime.Time
					}
				}
				if c.Type == corev1.PodScheduled {
					if c.LastTransitionTime.After(t.Scheduled) {
						t.Scheduled = c.LastTransitionTime.Time
					}
				}
			}
		}

		if r.Role == "test" {
			testTimes = t
		} else {
			if t.Scheduled.After(roleTimes.Scheduled) {
				roleTimes.Scheduled = t.Scheduled
			}
			if t.Ready.After(roleTimes.Ready) {
				roleTimes.Ready = t.Ready
			}
		}
	}

	t.Logf("Test Times ready=" + testTimes.Ready.String() + " scheduled=" + testTimes.Scheduled.String())
	t.Logf("Role Times ready=" + roleTimes.Ready.String() + " scheduled=" + roleTimes.Scheduled.String())
	jitter := time.Second * 1
	// add some jitter to the Ready time to round it up a second
	if testTimes.Scheduled.Before(roleTimes.Ready.Add(jitter)) {
		t.Fatalf("Expected test role scheduled time " + testTimes.Scheduled.String() + " to be after (or equal to) dependent pods ready time ")
	}
}

type PodTimes struct {
	Scheduled time.Time
	Ready     time.Time
}
