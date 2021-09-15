/*
 * Copyright (c) 2019, 2020, 2021 Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package remote

import (
	goctx "context"
	"fmt"
	coh "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/pkg/management"
	"github.com/oracle/coherence-operator/pkg/operator"
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
	// This test uses Management over REST to verify the site
	helper.SkipIfCoherenceVersionLessThan(t, 12, 2, 1, 4)
	// Ensure that everything is cleaned up after the test!
	testContext.CleanupAfterTest(t)

	fn := func(member management.MemberData) string {
		return member.SiteName
	}

	dfn := func(namespace string) string {
		return fmt.Sprintf("zone-zone-test-sts.%s.svc.cluster.local", namespace)
	}

	assertLabel(t, "zone", operator.DefaultSiteLabels, fn, dfn)
}

// Verify that a Coherence resource deployed by the Operator has the correct rack value
// set from the Node's failure domain region.
func TestRackLabel(t *testing.T) {
	// This test uses Management over REST to verify the rack
	helper.SkipIfCoherenceVersionLessThan(t, 12, 2, 1, 4)
	// Ensure that everything is cleaned up after the test!
	testContext.CleanupAfterTest(t)

	fn := func(member management.MemberData) string {
		return member.RackName
	}

	dfn := func(namespace string) string {
		return "n/a"
	}

	assertLabel(t, "rack", operator.DefaultRackLabels, fn, dfn)
}

func assertLabel(t *testing.T, name string, labels []string, fn func(management.MemberData) string, dfn func(string) string) {
	g := NewGomegaWithT(t)
	namespace := helper.GetTestNamespace()

	// load the test Coherence resource from a yaml files
	deployment, err := helper.NewSingleCoherenceFromYaml(namespace, "zone-test.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	deployment.SetName(name + "-zone-test")

	// deploy to k8s
	err = testContext.Client.Create(goctx.TODO(), &deployment)
	g.Expect(err).NotTo(HaveOccurred())

	replicas := deployment.GetReplicas()

	// Wait for the StatefulSet for the deployment to be ready - wait five minutes max
	sts, err := helper.WaitForStatefulSetForDeployment(testContext, namespace, &deployment, time.Second*10, time.Minute*5)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(sts.Status.ReadyReplicas).To(Equal(replicas))

	// Get the list of Pods
	pods, err := helper.ListCoherencePodsForDeployment(testContext, namespace, deployment.Name)
	g.Expect(err).NotTo(HaveOccurred())

	// capture the Pod log in case we need it for debugging
	helper.DumpPodLog(testContext, &pods[0], t.Name())

	// Port forward to the first Pod
	pf, ports, err := helper.StartPortForwarderForPod(&pods[0])
	g.Expect(err).NotTo(HaveOccurred())
	defer pf.Close()

	// Do a Management over REST query for the deployment members
	cl := &http.Client{}
	members, _, err := management.GetMembers(cl, "127.0.0.1", ports[coh.PortNameManagement])
	g.Expect(err).NotTo(HaveOccurred())

	// assert that the site or rack for each member matches the Node's zone label
	for _, member := range members.Items {
		g.Expect(member.MachineName).NotTo(BeEmpty())
		// The member's machine name is the k8s Node name
		node, err := testContext.KubeClient.CoreV1().Nodes().Get(context.TODO(), member.MachineName, metav1.GetOptions{})
		g.Expect(err).NotTo(HaveOccurred())

		actual := fn(member)
		expected := ""

		// find the first non-blank label as that is the label that should have been used
		for _, label := range labels {
			labelValue := node.GetLabels()[label]
			t.Logf("Node %s label '%s' value '%s'", node.Name, label, labelValue)
			if labelValue != "" {
				expected = labelValue
				break
			}
		}

		if expected == "" {
			// when running locally (for example in Docker on MacOS) the node might not
			// have labels unless they have been explicitly set by the developer.
			expected = dfn(namespace)
			t.Logf("Node %s found zero non-blank labels, using default expected '%s'", node.Name, expected)
		}

		t.Logf("Expecting label value to be: %s", expected)
		g.Expect(actual).To(Equal(expected))
	}
}
