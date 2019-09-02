package local

import (
	"context"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	coh "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	mgmt "github.com/oracle/coherence-operator/pkg/management"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	"net/http"
	"testing"
	"time"

	. "github.com/onsi/gomega"
)

// ----- tests --------------------------------------------------------------

func TestMinimalCoherenceCluster(t *testing.T) {
	assertCluster(t, "cluster-minimal.yaml", map[string]int32{coh.DefaultRoleName: coh.DefaultReplicas})
}

func TestNoRoleOneReplica(t *testing.T) {
	assertCluster(t, "cluster-no-role-one-replica.yaml", map[string]int32{coh.DefaultRoleName: 1})
}

func TestOneRoleDefaultReplicas(t *testing.T) {
	assertCluster(t, "cluster-one-role-default-replica.yaml", map[string]int32{"data": coh.DefaultReplicas})
}

func TestOneRoleOneReplicas(t *testing.T) {
	assertCluster(t, "cluster-one-role-one-replica.yaml", map[string]int32{coh.DefaultRoleName: 1})
}

func TestOneRoleTwoReplicas(t *testing.T) {
	assertCluster(t, "cluster-one-role-two-replica.yaml", map[string]int32{"data": 2})
}

func TestTwoRolesDefaultReplicas(t *testing.T) {
	assertCluster(t, "cluster-two-roles-default-replica.yaml", map[string]int32{"data": coh.DefaultReplicas, "proxy": coh.DefaultReplicas})
}

func TestTwoRolesOneReplicas(t *testing.T) {
	assertCluster(t, "cluster-two-roles-one-replica.yaml", map[string]int32{"data": 1, "proxy": 1})
}

func TestTwoRolesTwoReplicas(t *testing.T) {
	assertCluster(t, "cluster-two-roles-different-replica.yaml", map[string]int32{"data": 2, "proxy": 1})
}

// ----- helpers ------------------------------------------------------------

// Test that a cluster can be created using the specified yaml.
func assertCluster(t *testing.T, yamlFile string, expectedRoles map[string]int32) {
	// initialise Gomega so we can use matchers
	g := NewGomegaWithT(t)
	f := framework.Global

	// work out the total expected roles and cluster size
	totalRoles := 0
	clusterSize := 0
	for _, size := range expectedRoles {
		clusterSize = clusterSize + int(size)
		if size > 0 {
			totalRoles = totalRoles + 1
		}
	}

	// Create the Operator SDK test context (this will deploy the Operator)
	ctx := helper.CreateTestContext(t)
	// Make sure we defer clean-up (uninstall the operator) when we're done
	defer helper.DumpOperatorLogsAndCleanup(t, ctx)

	// Get the test namespace
	namespace, err := ctx.GetNamespace()
	g.Expect(err).NotTo(HaveOccurred())

	cluster, err := helper.NewCoherenceClusterFromYaml(namespace, yamlFile)

	// verify the cluster size is expected
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(cluster.GetClusterSize()).To(Equal(clusterSize))

	// deploy the CoherenceCluster
	err = f.Client.Create(context.TODO(), &cluster, helper.DefaultCleanup(ctx))
	g.Expect(err).NotTo(HaveOccurred())

	roles := cluster.GetRoles()

	// Assert that a CoherenceRole is created for each role in the cluster
	for _, role := range roles {
		roleName := role.GetFullRoleName(&cluster)
		// Wait for a CoherenceRole to be created
		role, err := helper.WaitForCoherenceRole(f, namespace, roleName, time.Second*10, time.Minute*2, t)
		g.Expect(err).NotTo(HaveOccurred())

		expectedReplicas, found := expectedRoles[role.Spec.GetRoleName()]
		g.Expect(found).To(BeTrue(), "Found Role with unexpected name '"+roleName+"'")
		g.Expect(role.Spec.GetReplicas()).To(Equal(expectedReplicas))
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
	g.Expect(len(pods)).To(Equal(clusterSize))

	// Start a port-forwarder that will forward ALL ports on a Pod (the first pod in the list)
	pf, ports, err := helper.StartPortForwarderForPod(&pods[0])
	g.Expect(err).NotTo(HaveOccurred())

	// ensure the port-forwarder is closed when this method exits
	defer pf.Close()

	// Do a Management over ReST query to get the cluster size
	clusterData, status, err := mgmt.GetCluster(&http.Client{}, "127.0.0.1", ports[mgmt.PortName])
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(status).To(Equal(http.StatusOK))
	g.Expect(clusterData.ClusterSize).To(Equal(clusterSize))
}
