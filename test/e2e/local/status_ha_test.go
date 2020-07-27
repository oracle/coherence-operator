/*
 * Copyright (c) 2019, 2020 Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package local

import (
	goctx "context"
	"flag"
	"github.com/operator-framework/operator-sdk/pkg/log/zap"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	coh "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/controllers/statefulset"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/klog"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"testing"
	"time"

	. "github.com/onsi/gomega"
)

type StatusHATestCase struct {
	Deployment *coh.Coherence
	Name       string
}

func TestStatusHA(t *testing.T) {
	g := NewGomegaWithT(t)
	ns := helper.GetTestNamespace()

	logf.SetLogger(zap.Logger())

	flags := &flag.FlagSet{}
	klog.InitFlags(flags)
	_ = flags.Set("v", "4")

	deploymentDefault, err := helper.NewSingleCoherenceFromYaml(ns, "status-ha-default.yaml")
	g.Expect(err).NotTo(HaveOccurred())
	deploymentExec, err := helper.NewSingleCoherenceFromYaml(ns, "status-ha-exec.yaml")
	g.Expect(err).NotTo(HaveOccurred())
	deploymentHttp, err := helper.NewSingleCoherenceFromYaml(ns, "status-ha-http.yaml")
	g.Expect(err).NotTo(HaveOccurred())
	deploymentTcp, err := helper.NewSingleCoherenceFromYaml(ns, "status-ha-tcp.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	testCases := []StatusHATestCase{
		{Deployment: &deploymentDefault, Name: "DefaultStatusHAHandler"},
		{Deployment: &deploymentExec, Name: "ExecStatusHAHandler"},
		{Deployment: &deploymentHttp, Name: "HttpStatusHAHandler"},
		{Deployment: &deploymentTcp, Name: "TcpStatusHAHandler"},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			assertStatusHA(t, tc)
		})
	}
}

func assertStatusHA(t *testing.T, tc StatusHATestCase) {
	g := NewGomegaWithT(t)
	f := framework.Global
	ctx := helper.CreateTestContext(t)
	defer helper.DumpOperatorLogs(t)

	ns, err := ctx.GetWatchNamespace()
	g.Expect(err).NotTo(HaveOccurred())

	err = f.Client.Create(goctx.TODO(), tc.Deployment, helper.DefaultCleanup(ctx))
	g.Expect(err).NotTo(HaveOccurred())

	sts, err := helper.WaitForStatefulSetForDeployment(f.KubeClient, ns, tc.Deployment, time.Second*10, time.Minute*5, t)
	g.Expect(err).NotTo(HaveOccurred())

	// Get the list of Pods
	pods, err := helper.ListCoherencePodsForDeployment(f.KubeClient, ns, tc.Deployment.Name)
	g.Expect(err).NotTo(HaveOccurred())

	// capture the Pod log in case we need it for debugging
	helper.DumpPodLog(f.KubeClient, &pods[0], t.Name(), t)

	pf, ports, err := helper.StartPortForwarderForPod(&pods[0])
	g.Expect(err).NotTo(HaveOccurred())
	defer pf.Close()

	ckr := statefulset.ScalableChecker{Client: f.Client.Client, Config: f.KubeConfig}
	ckr.SetGetPodHostName(func(pod corev1.Pod) string { return "127.0.0.1" })
	ckr.SetTranslatePort(func(name string, port int) int { return int(ports[name]) })
	ha := ckr.IsStatusHA(tc.Deployment, sts)
	g.Expect(ha).To(BeTrue())
}
