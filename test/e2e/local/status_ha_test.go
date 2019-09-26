/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package local

import (
	goctx "context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/operator-framework/operator-sdk/pkg/log/zap"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	coh "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	"github.com/oracle/coherence-operator/pkg/controller/coherencerole"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/klog"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"testing"
	"time"

	. "github.com/onsi/gomega"
)

type StatusHATestCase struct {
	Cluster *coh.CoherenceCluster
	Name    string
}

func TestStatusHA(t *testing.T) {
	g := NewGomegaWithT(t)
	ns := helper.GetTestNamespace()

	logf.SetLogger(zap.Logger())

	flags := &flag.FlagSet{}
	klog.InitFlags(flags)
	_ = flags.Set("v", "4")

	clusterDefault, err := helper.NewCoherenceClusterFromYaml(ns, "status-ha-default.yaml")
	g.Expect(err).NotTo(HaveOccurred())
	clusterExec, err := helper.NewCoherenceClusterFromYaml(ns, "status-ha-exec.yaml")
	g.Expect(err).NotTo(HaveOccurred())
	clusterHttp, err := helper.NewCoherenceClusterFromYaml(ns, "status-ha-http.yaml")
	g.Expect(err).NotTo(HaveOccurred())
	clusterTcp, err := helper.NewCoherenceClusterFromYaml(ns, "status-ha-tcp.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	testCases := []StatusHATestCase{
		{Cluster: &clusterDefault, Name: "DefaultStatusHAHandler"},
		{Cluster: &clusterExec, Name: "ExecStatusHAHandler"},
		{Cluster: &clusterHttp, Name: "HttpStatusHAHandler"},
		{Cluster: &clusterTcp, Name: "TcpStatusHAHandler"},
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
	defer helper.DumpOperatorLogsAndCleanup(t, ctx)

	ns, err := ctx.GetNamespace()
	g.Expect(err).NotTo(HaveOccurred())

	d, _ := json.Marshal(tc.Cluster)
	fmt.Printf("StatusHA Test installing cluster:\n%s\n", string(d))

	err = f.Client.Create(goctx.TODO(), tc.Cluster, helper.DefaultCleanup(ctx))
	g.Expect(err).NotTo(HaveOccurred())

	g.Expect(len(tc.Cluster.Spec.Roles)).To(Equal(1))

	roleSpec := tc.Cluster.Spec.Roles[0]

	sts, err := helper.WaitForStatefulSetForRole(f.KubeClient, ns, tc.Cluster, roleSpec, time.Second*10, time.Minute*5, t)
	g.Expect(err).NotTo(HaveOccurred())

	// Get the list of Pods
	pods, err := helper.ListCoherencePodsForRole(f.KubeClient, ns, tc.Cluster.Name, roleSpec.GetRoleName())
	g.Expect(err).NotTo(HaveOccurred())

	// capture the Pod log in case we need it for debugging
	helper.DumpPodLog(f.KubeClient, &pods[0], t.Name(), t)

	pf, ports, err := helper.StartPortForwarderForPod(&pods[0])
	g.Expect(err).NotTo(HaveOccurred())
	defer pf.Close()

	role, err := helper.GetCoherenceRole(f, ns, roleSpec.GetFullRoleName(tc.Cluster))
	g.Expect(err).NotTo(HaveOccurred())

	ckr := coherencerole.ScaleableChecker{Client: f.Client.Client, Config: f.KubeConfig}
	ckr.SetGetPodHostName(func(pod corev1.Pod) string { return "127.0.0.1" })
	ckr.SetTranslatePort(func(name string, port int) int { return int(ports[name]) })
	ha := ckr.IsStatusHA(role, sts)
	g.Expect(ha).To(BeTrue())
}
