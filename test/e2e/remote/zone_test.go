/*
 * Copyright (c) 2019, 2020, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package remote

import (
	goctx "context"
	"fmt"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	coh "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/pkg/flags"
	"github.com/oracle/coherence-operator/pkg/management"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	"golang.org/x/net/context"
	"net/http"
	"testing"
	"time"

	. "github.com/onsi/gomega"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Verify that a Coherence resource deployed by the Operator has the correct site value
// set from the Node's failure domain zone.
func TestSiteLabel(t *testing.T) {
	// This test uses Management over ReST to verify the site
	helper.SkipIfCoherenceVersionLessThan(t, 12, 2, 1, 4)

	fn := func(member management.MemberData) string {
		return member.SiteName
	}

	dfn := func(namespace string) string {
		return fmt.Sprintf("zone-zone-test-sts.%s.svc.cluster.local", namespace)
	}

	assertLabel(t, "zone", flags.DefaultSiteLabel, fn, dfn)
}

// Verify that a Coherence resource deployed by the Operator has the correct rack value
// set from the Node's failure domain region.
func TestRackLabel(t *testing.T) {
	// This test uses Management over ReST to verify the rack
	helper.SkipIfCoherenceVersionLessThan(t, 12, 2, 1, 4)

	fn := func(member management.MemberData) string {
		return member.RackName
	}

	dfn := func(namespace string) string {
		return "n/a"
	}

	assertLabel(t, "rack", flags.DefaultRackLabel, fn, dfn)
}

func assertLabel(t *testing.T, name string, label string, fn func(management.MemberData) string, dfn func(string) string) {
	g := NewGomegaWithT(t)
	f := framework.Global

	ctx := helper.CreateTestContext(t)
	defer helper.DumpOperatorLogs(t)

	namespace, err := ctx.GetWatchNamespace()
	g.Expect(err).NotTo(HaveOccurred())

	// load the test Coherence resource from a yaml files
	deployment, err := helper.NewSingleCoherenceFromYaml(namespace, "zone-test.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	deployment.SetName(name + "-zone-test")

	// deploy to k8s
	err = f.Client.Create(goctx.TODO(), &deployment, helper.DefaultCleanup(ctx))
	g.Expect(err).NotTo(HaveOccurred())

	replicas := deployment.GetReplicas()

	// Wait for the StatefulSet for the deployment to be ready - wait five minutes max
	sts, err := helper.WaitForStatefulSetForDeployment(f.KubeClient, namespace, &deployment, time.Second*10, time.Minute*5)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(sts.Status.ReadyReplicas).To(Equal(replicas))

	// Get the list of Pods
	pods, err := helper.ListCoherencePodsForDeployment(f.KubeClient, namespace, deployment.Name)
	g.Expect(err).NotTo(HaveOccurred())

	// capture the Pod log in case we need it for debugging
	helper.DumpPodLog(f.KubeClient, &pods[0], t.Name(), t)

	// Port forward to the first Pod
	pf, ports, err := helper.StartPortForwarderForPod(&pods[0])
	g.Expect(err).NotTo(HaveOccurred())
	defer pf.Close()

	// Do a Management over ReST query for the deployment members
	cl := &http.Client{}
	members, _, err := management.GetMembers(cl, "127.0.0.1", ports[coh.PortNameManagement])
	g.Expect(err).NotTo(HaveOccurred())

	// assert that the site for each member matches the Node's zone label
	for _, member := range members.Items {
		g.Expect(member.MachineName).NotTo(BeEmpty())
		// The member's machine name is the k8s Node name
		node, err := f.KubeClient.CoreV1().Nodes().Get(context.TODO(), member.MachineName, metav1.GetOptions{})
		g.Expect(err).NotTo(HaveOccurred())
		zone := node.GetLabels()[label]

		actual := fn(member)
		if zone != "" {
			t.Logf("Expecting label value to be: %s", zone)
			g.Expect(actual).To(Equal(zone))
		} else {
			// when running locally (for example in Docker on MacOS) the node might not
			// have a zone unless one has been explicitly set by the developer.
			t.Logf("Expecting label value to be: %s", dfn(namespace))
			g.Expect(actual).To(Equal(dfn(namespace)))
		}
	}
}
