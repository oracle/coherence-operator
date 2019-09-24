/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package local

import (
	"context"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	coherence "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	"github.com/oracle/coherence-operator/pkg/management"
	"github.com/oracle/coherence-operator/test/e2e/helper"

	"io/ioutil"
	"net/http"

	corev1 "k8s.io/api/core/v1"
	"testing"
	"time"

	. "github.com/onsi/gomega"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	loggingConfigMapName = "logging-config-map"
)

/*
Deploy a CoherenceCluster with the log level changed in the "spec" section.
*/
func TestCoherenceLogLevelSpecSection(t *testing.T) {
	var (
		clusterName = "mycluster"
		roleName    = coherence.DefaultRoleName
		replicas    = coherence.DefaultReplicas
	)
	g := NewGomegaWithT(t)

	ctx := helper.CreateTestContext(t)
	defer helper.DumpOperatorLogsAndCleanup(t, ctx)

	namespace, err := ctx.GetNamespace()
	g.Expect(err).NotTo(HaveOccurred())

	assertClusterCreated(t, ctx, namespace, "spec_log_level.yaml",
		map[string]int32{coherence.DefaultRoleName: replicas})
	assertLogConfigOnMembers(g, namespace, clusterName, roleName, int(replicas), 9)
}

/*
Deploy a CoherenceCluster with the log level changed in the "roles" section.
*/
func TestCoherenceLogLevelRolesSection(t *testing.T) {
	var (
		clusterName       = "mycluster"
		roleName          = "storage"
		replicas    int32 = 1
	)
	g := NewGomegaWithT(t)

	ctx := helper.CreateTestContext(t)
	defer helper.DumpOperatorLogsAndCleanup(t, ctx)

	namespace, err := ctx.GetNamespace()
	g.Expect(err).NotTo(HaveOccurred())

	assertClusterCreated(t, ctx, namespace, "role_log_level.yaml",
		map[string]int32{coherence.DefaultRoleName: replicas})
	assertLogConfigOnMembers(g, namespace, clusterName, roleName, int(replicas), 9)
}

/*
Deploy a CoherenceCluster with the logging config set to a custom file in the
"spec" section.
*/
func TestCoherenceLogConfigFileSpecSection(t *testing.T) {
	var (
		clusterName       = "mycluster"
		roleName          = coherence.DefaultRoleName
		replicas    int32 = coherence.DefaultReplicas
	)
	g := NewGomegaWithT(t)

	ctx := helper.CreateTestContext(t)
	defer helper.DumpOperatorLogsAndCleanup(t, ctx)

	namespace, err := ctx.GetNamespace()
	g.Expect(err).NotTo(HaveOccurred())

	assertClusterCreated(t, ctx, namespace, "spec_log_configfile.yaml",
		map[string]int32{coherence.DefaultRoleName: replicas})
	assertLogConfigOnMembers(g, namespace, clusterName, roleName, int(replicas), 5)
}

/*
Deploy a CoherenceCluster with the logging config set to a custom file in the
"roles" section.
*/
func TestCoherenceLogConfigFileRolesSection(t *testing.T) {
	var (
		clusterName       = "mycluster"
		roleName          = "storage"
		replicas    int32 = 1
	)
	g := NewGomegaWithT(t)

	ctx := helper.CreateTestContext(t)
	defer helper.DumpOperatorLogsAndCleanup(t, ctx)

	namespace, err := ctx.GetNamespace()
	g.Expect(err).NotTo(HaveOccurred())

	assertClusterCreated(t, ctx, namespace, "role_log_configfile.yaml",
		map[string]int32{coherence.DefaultRoleName: replicas})
	assertLogConfigOnMembers(g, namespace, clusterName, roleName, int(replicas), 5)
}

/*
Create a ConfigMap containing the custom logging configuration.
Deploy a CoherenceCluster with the logging config set to the ConfigMap file in the
"spec" section.
*/
func TestCoherenceLogConfigMapSpecSection(t *testing.T) {
	var (
		clusterName       = "mycluster"
		roleName          = coherence.DefaultRoleName
		replicas    int32 = coherence.DefaultReplicas
	)
	g := NewGomegaWithT(t)

	ctx := helper.CreateTestContext(t)
	defer helper.DumpOperatorLogsAndCleanup(t, ctx)

	namespace, err := ctx.GetNamespace()
	g.Expect(err).NotTo(HaveOccurred())

	createLoggingConfigMap(g, namespace)

	assertClusterCreated(t, ctx, namespace, "spec_log_configmap.yaml",
		map[string]int32{coherence.DefaultRoleName: replicas})
	assertLogConfigOnMembers(g, namespace, clusterName, roleName, int(replicas), 5)

	deleteLoggingConfigMap(namespace)
}

/*
Create a ConfigMap containing the custom logging configuration.
Deploy a CoherenceCluster with the logging config set to the ConfigMap file in the
"roles" section.
*/
func TestCoherenceLogConfigMapRolesSection(t *testing.T) {
	var (
		clusterName       = "mycluster"
		roleName          = "storage"
		replicas    int32 = 1
	)
	g := NewGomegaWithT(t)

	ctx := helper.CreateTestContext(t)
	defer helper.DumpOperatorLogsAndCleanup(t, ctx)

	namespace, err := ctx.GetNamespace()
	g.Expect(err).NotTo(HaveOccurred())

	createLoggingConfigMap(g, namespace)

	assertClusterCreated(t, ctx, namespace, "role_log_configmap.yaml",
		map[string]int32{coherence.DefaultRoleName: replicas})
	assertLogConfigOnMembers(g, namespace, clusterName, roleName, int(replicas), 5)

	deleteLoggingConfigMap(namespace)
}

// ----- helpers ------------------------------------------------------------

// Verify the cluster is created using given yamlFile.
func assertClusterCreated(t *testing.T, ctx *framework.TestCtx, namespace, yamlFile string, expectedRoles map[string]int32) {
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

	// Verify that the WKA service has the same number of endpoints as the cluster size.
	serviceName := cluster.GetWkaServiceName()
	ep, err := f.KubeClient.CoreV1().Endpoints(namespace).Get(serviceName, metav1.GetOptions{})
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(len(ep.Subsets)).NotTo(BeZero())

	subset := ep.Subsets[0]
	g.Expect(len(subset.Addresses)).To(Equal(clusterSize))
}

// Delete logging-config-map from the given namespace
func deleteLoggingConfigMap(namespace string) {
	_ = framework.Global.KubeClient.CoreV1().ConfigMaps(namespace).Delete(loggingConfigMapName, &metav1.DeleteOptions{})
}

// Create required logging-config-map in a test namespace.
func createLoggingConfigMap(g *GomegaWithT, namespace string) {
	f := framework.Global
	logConfigMap := &corev1.ConfigMap{}
	logConfigMap.SetNamespace(namespace)
	logConfigMap.SetName(loggingConfigMapName)

	logConfigMap.Data = make(map[string]string)
	rootDir, err := helper.FindProjectRootDir()
	g.Expect(err).NotTo(HaveOccurred())

	var cmData []byte
	cmData, err = ioutil.ReadFile(rootDir + "/test/e2e/local/logging.properties")
	g.Expect(err).NotTo(HaveOccurred())
	logConfigMap.Data["logging.properties"] = string(cmData)

	existingCM, _ := f.KubeClient.CoreV1().ConfigMaps(namespace).Get(loggingConfigMapName,
		metav1.GetOptions{})
	if existingCM != nil {
		deleteLoggingConfigMap(namespace)
	}

	_, err = f.KubeClient.CoreV1().ConfigMaps(namespace).Create(logConfigMap)
	g.Expect(err).NotTo(HaveOccurred())
}

// Verify log configuration on each member of a cluster which is created using specified yaml file.
func assertLogConfigOnMembers(g *GomegaWithT, namespace, clusterName, roleName string, replicas, logLevel int) {
	ok, _ := helper.IsCoherenceVersionAtLeast(12, 2, 1, 4)
	if ok {
		f := framework.Global

		// Get the list of Pods
		pods, err := helper.ListCoherencePodsForRole(f.KubeClient, namespace, clusterName, roleName)
		g.Expect(err).NotTo(HaveOccurred())

		// Port forward to the first Pod
		pf, ports, err := helper.StartPortForwarderForPod(&pods[0])
		g.Expect(err).NotTo(HaveOccurred())
		defer pf.Close()

		cl := &http.Client{}
		members, _, err := management.GetMembers(cl, "127.0.0.1", ports[management.PortName])
		g.Expect(err).NotTo(HaveOccurred())
		g.Expect(len(members.Items)).To(Equal(replicas))

		// assert that each member has required logging configuration.
		for _, member := range members.Items {
			g.Expect(member.MachineName).NotTo(BeEmpty())
			g.Expect(member.LoggingLevel).To(Equal(logLevel))
		}
	}
}
