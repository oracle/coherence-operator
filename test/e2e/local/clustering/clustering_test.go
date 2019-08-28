// The clustering package contains functional tests related to Coherence clustering.
package clustering

import (
	"context"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	mgmt "github.com/oracle/coherence-operator/pkg/management"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	"net/http"
	"testing"
	"time"

	. "github.com/onsi/gomega"
)

// TestClustering verifies that different CoherenceCluster configurations form a cluster.
func TestMinimalCoherenceCluster(t *testing.T) {
	// initialise Gomega so we can use matchers
	g := NewGomegaWithT(t)
	f := framework.Global

	// Create the Operator SDK test context (this will deploy the Operator)
	ctx := helper.CreateTestContext(t)
	// Make sure we defer clean-up (uninstall the operator) when we're done
	defer helper.DumpOperatorLogsAndCleanup(t, ctx)

	// Get the test namespace
	namespace, err := ctx.GetNamespace()
	g.Expect(err).NotTo(HaveOccurred())

	cluster, err := helper.NewCoherenceClusterFromYaml(namespace, "clustering-minimal.yaml")

	// deploy the CoherenceCluster
	err = f.Client.Create(context.TODO(), &cluster, helper.DefaultCleanup(ctx))
	g.Expect(err).NotTo(HaveOccurred())

	roles := cluster.GetRoles()

	// Assert that a CoherenceRole is created for each role in the cluster
	for _, role := range roles {
		roleName := role.GetFullRoleName(&cluster)
		// Wait for a CoherenceRole to be created - we expect one for a minimal CoherenceCluster
		_, err := helper.WaitForCoherenceRole(f, namespace, roleName, time.Second*10, time.Minute*2, t)
		g.Expect(err).NotTo(HaveOccurred())

	}

	// Assert that a StatefulSet of the correct number or replicas is created for each role in the cluster
	for _, role := range roles {
		// Wait for the StatefulSet for the role to be ready - wait five minutes max
		sts, err := helper.WaitForStatefulSetForRole(f.KubeClient, namespace, &cluster, role, time.Second*10, time.Minute*5, t)
		g.Expect(err).NotTo(HaveOccurred())
		g.Expect(sts.Status.ReadyReplicas).To(Equal(role.GetReplicas()))
	}

	// Get all of the Pods in the cluster
	pods, err := helper.ListCoherencePodsForCluster(f.KubeClient, namespace, cluster.Name)
	g.Expect(err).NotTo(HaveOccurred())

	// assert that the correct number of Pods is returned
	size := cluster.GetClusterSize()
	g.Expect(len(pods)).To(Equal(size))

	// Start a port-forwarder that will forward ALL ports on a Pod (the first pod in the list)
	pf, ports, err := helper.StartPortForwarderForPod(&pods[0])
	g.Expect(err).NotTo(HaveOccurred())

	// ensure the port-forwarder is closed when this method exits
	defer pf.Close()

	// Do a Management over ReST query to get the cluster size
	clusterData, status, err := mgmt.GetCluster(&http.Client{}, "127.0.0.1", ports[mgmt.PORT_NAME])
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(status).To(Equal(http.StatusOK))
	g.Expect(clusterData.ClusterSize).To(Equal(size))
}
