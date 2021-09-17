/*
 * Copyright (c) 2019, 2021, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package local

import (
	goctx "context"
	coh "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/controllers/statefulset"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	corev1 "k8s.io/api/core/v1"
	"testing"
	"time"

	. "github.com/onsi/gomega"
)

type StatusHATestCase struct {
	Deployment *coh.Coherence
	Name       string
}

func TestStatusHA(t *testing.T) {
	// Make sure we defer clean-up when we're done!!
	testContext.CleanupAfterTest(t)

	g := NewGomegaWithT(t)
	ns := helper.GetTestNamespace()

	deploymentDefault, err := helper.NewSingleCoherenceFromYaml(ns, "status-ha-default.yaml")
	g.Expect(err).NotTo(HaveOccurred())
	deploymentExec, err := helper.NewSingleCoherenceFromYaml(ns, "status-ha-exec.yaml")
	g.Expect(err).NotTo(HaveOccurred())
	deploymentHTTP, err := helper.NewSingleCoherenceFromYaml(ns, "status-ha-http.yaml")
	g.Expect(err).NotTo(HaveOccurred())
	deploymentTCP, err := helper.NewSingleCoherenceFromYaml(ns, "status-ha-tcp.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	testCases := []StatusHATestCase{
		{Deployment: &deploymentDefault, Name: "DefaultStatusHAHandler"},
		{Deployment: &deploymentExec, Name: "ExecStatusHAHandler"},
		{Deployment: &deploymentHTTP, Name: "HttpStatusHAHandler"},
		{Deployment: &deploymentTCP, Name: "TcpStatusHAHandler"},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			assertStatusHA(t, tc)
		})
	}
}

func assertStatusHA(t *testing.T, tc StatusHATestCase) {
	// Make sure we defer clean-up when we're done!!
	testContext.CleanupAfterTest(t)

	g := NewGomegaWithT(t)
	ns := helper.GetTestNamespace()

	err := testContext.Client.Create(goctx.TODO(), tc.Deployment)
	g.Expect(err).NotTo(HaveOccurred())

	sts, err := helper.WaitForStatefulSetForDeployment(testContext, ns, tc.Deployment, time.Second*10, time.Minute*5)
	g.Expect(err).NotTo(HaveOccurred())

	// Get the list of Pods
	pods, err := helper.ListCoherencePodsForDeployment(testContext, ns, tc.Deployment.Name)
	g.Expect(err).NotTo(HaveOccurred())

	pf, ports, err := helper.StartPortForwarderForPod(&pods[0])
	g.Expect(err).NotTo(HaveOccurred())
	defer pf.Close()

	ckr := statefulset.CoherenceProbe{Client: testContext.Client, Config: testContext.Config}
	ckr.SetGetPodHostName(func(pod corev1.Pod) string { return "127.0.0.1" })
	ckr.SetTranslatePort(func(name string, port int) int { return int(ports[name]) })
	ha := ckr.IsStatusHA(testContext.Context, tc.Deployment, sts)
	g.Expect(ha).To(BeTrue())
}
