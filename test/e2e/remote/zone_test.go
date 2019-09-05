/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package remote

import (
	goctx "context"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	"github.com/oracle/coherence-operator/pkg/flags"
	"github.com/oracle/coherence-operator/pkg/management"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	"net/http"
	"testing"
	"time"

	. "github.com/onsi/gomega"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Verify that a CoherenceCluster deployed by the Operator has the correct site value
// set from the Node's failure domain zone.
func TestSiteLabel(t *testing.T) {
	// This test uses Management over ReST to verify the site
	helper.SkipIfCoherenceVersionLessThan(t, 12, 2, 1, 4)

	fn := func(member management.MemberData) string {
		return member.SiteName
	}

	assertLabel(t, flags.DefaultSiteLabel, fn)
}

// Verify that a CoherenceCluster deployed by the Operator has the correct rack value
// set from the Node's failure domain region.
func TestRackLabel(t *testing.T) {
	// This test uses Management over ReST to verify the rack
	helper.SkipIfCoherenceVersionLessThan(t, 12, 2, 1, 4)

	fn := func(member management.MemberData) string {
		return member.RackName
	}

	assertLabel(t, flags.DefaultRackLabel, fn)
}

func assertLabel(t *testing.T, label string, fn func(management.MemberData) string) {
	g := NewGomegaWithT(t)

	f := framework.Global

	ctx := helper.CreateTestContext(t)
	defer helper.DumpOperatorLogsAndCleanup(t, ctx)

	namespace, err := ctx.GetNamespace()
	g.Expect(err).NotTo(HaveOccurred())

	// load the test CoherenceCluster from a yaml files
	cluster, err := helper.NewCoherenceClusterFromYaml(namespace, "zone-test.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	// deploy the CoherenceCluster
	err = f.Client.Create(goctx.TODO(), &cluster, helper.DefaultCleanup(ctx))
	g.Expect(err).NotTo(HaveOccurred())

	role := cluster.Spec.Roles[0]
	replicas := role.GetReplicas()

	// Wait for the StatefulSet for the role to be ready - wait five minutes max
	sts, err := helper.WaitForStatefulSetForRole(f.KubeClient, namespace, &cluster, role, time.Second*10, time.Minute*5, t)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(sts.Status.ReadyReplicas).To(Equal(replicas))

	// Get the list of Pods
	pods, err := helper.ListCoherencePodsForRole(f.KubeClient, namespace, cluster.Name, role.GetRoleName())
	g.Expect(err).NotTo(HaveOccurred())

	// capture the Pod log in case we need it for debugging
	helper.DumpPodLog(f.KubeClient, &pods[0], t.Name(), t)

	// Port forward to the first Pod
	pf, ports, err := helper.StartPortForwarderForPod(&pods[0])
	g.Expect(err).NotTo(HaveOccurred())
	defer pf.Close()

	// Do a Management over ReST query for the cluster members
	cl := &http.Client{}
	members, _, err := management.GetMembers(cl, "127.0.0.1", ports["mgmt-port"])
	g.Expect(err).NotTo(HaveOccurred())

	// assert that the site for each member matches the Node's zone label
	for _, member := range members.Items {
		g.Expect(member.MachineName).NotTo(BeEmpty())
		// The member's machine name is the k8s Node name
		node, err := f.KubeClient.CoreV1().Nodes().Get(member.MachineName, metav1.GetOptions{})
		g.Expect(err).NotTo(HaveOccurred())
		zone := node.GetLabels()[label]

		actual := fn(member)

		if zone != "" {
			g.Expect(actual).To(Equal(zone))
		} else {
			// when running locally (for example in Docker on MacOS) the node might not
			// have a zone unless one has been explicitly set by the developer.
			g.Expect(actual).To(Equal("n/a"))
		}
	}
}
