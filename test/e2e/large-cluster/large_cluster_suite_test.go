/*
 * Copyright (c) 2020, 2024, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package large_cluster

import (
	"fmt"
	. "github.com/onsi/gomega"
	cohv1 "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"net/http"
	"os"
	"testing"
	"time"
)

var testContext helper.TestContext
var nodeList *corev1.NodeList
var nodeMap map[string]corev1.Node
var clusterCount = 0
var statusHA = true

var zones = [...]string{"zone-1", "zone-2", "zone-3"}
var faultDomains = [...]string{"fd-1", "fd-2"}

// The entry point for the test suite
func TestMain(m *testing.M) {
	var err error

	helper.EnsureTestEnvVars()

	// Create a new TestContext
	if testContext, err = helper.NewContext(true); err != nil {
		fmt.Printf("Error: %+v", err)
		os.Exit(1)
	}

	testContext.RestEndpoints["/ha"] = isStatusHA

	if err = testContext.Start(); err != nil {
		fmt.Printf("Error: %+v", err)
		os.Exit(1)
	}

	nodeMap = make(map[string]corev1.Node)
	nodeList, err = testContext.KubeClient.CoreV1().Nodes().List(testContext.Context, metav1.ListOptions{})
	if err != nil {
		fmt.Printf("Failed to get Node list from K8s. %+v", err)
		os.Exit(1)
	}

	zoneCount := len(zones)
	fdCount := len(faultDomains)
	for i, node := range nodeList.Items {
		nodeMap[node.Name] = node
		zone := zones[i%zoneCount]
		fd := faultDomains[(i/zoneCount)%fdCount]
		labels := node.Labels
		labels["failure-domain.beta.kubernetes.io/region"] = "region-1"
		labels["failure-domain.beta.kubernetes.io/zone"] = zone
		labels["topology.kubernetes.io/region"] = "region-1"
		labels["topology.kubernetes.io/zone"] = zone
		labels["oci.oraclecloud.com/fault-domain"] = fd
		node.Labels = labels
		updated, err := testContext.KubeClient.CoreV1().Nodes().Update(testContext.Context, &node, metav1.UpdateOptions{})
		if err != nil {
			fmt.Printf("Failed to label Node. %+v", err)
			os.Exit(1)
		}
		nodeList.Items[i] = *updated
	}

	exitCode := m.Run()
	testContext.Logf("Tests completed with return code %d", exitCode)
	testContext.Close()
	os.Exit(exitCode)
}

func SetStatusHA(ha bool) {
	statusHA = ha
}

func isStatusHA(w http.ResponseWriter, r *http.Request) {
	if statusHA {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
	_, _ = fmt.Fprint(w, "")
}

func GenerateClusterName() string {
	clusterCount++
	return fmt.Sprintf("cluster-%d", clusterCount)
}

func GetTestContext() helper.TestContext {
	return testContext
}
func GetRestPort() int32 {
	return testContext.RestServer.GetPort()
}

// installSimpleDeployment installs a deployment and asserts that the underlying
// StatefulSet resources reach the correct state.
func installSimpleDeployment(t *testing.T, d cohv1.Coherence) (cohv1.Coherence, appsv1.StatefulSet) {
	g := NewGomegaWithT(t)
	helper.AddLoopbackTestHostnameLabel(&d)
	err := testContext.Client.Create(testContext.Context, &d)
	g.Expect(err).NotTo(HaveOccurred())
	return assertDeploymentEventuallyInDesiredState(t, d, d.GetReplicas())
}

// assertDeploymentEventuallyInDesiredState asserts that a Coherence resource exists and has the correct spec and that the
// underlying StatefulSet exists with the correct status and ready replicas.
func assertDeploymentEventuallyInDesiredState(t *testing.T, d cohv1.Coherence, replicas int32) (cohv1.Coherence, appsv1.StatefulSet) {
	g := NewGomegaWithT(t)

	testContext.Logf("Asserting Coherence resource %s exists with %d replicas", d.Name, replicas)

	// create a DeploymentStateCondition that checks a deployment's replica count
	condition := helper.ReplicaCountCondition(replicas)

	// wait for the deployment to match the condition
	_, err := helper.WaitForCoherenceCondition(testContext, d.Namespace, d.Name, condition, time.Second*10, time.Minute*5)
	g.Expect(err).NotTo(HaveOccurred())

	testContext.Logf("Asserting StatefulSet %s exists with %d replicas", d.Name, replicas)

	// wait for the StatefulSet to have the required ready replicas
	sts, err := helper.WaitForStatefulSet(testContext, d.Namespace, d.Name, replicas, time.Second*10, time.Minute*5)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(sts.Status.ReadyReplicas).To(Equal(replicas))

	testContext.Logf("Asserting StatefulSet %s exist with %d replicas - Done!", d.Name, replicas)

	err = testContext.Client.Get(testContext.Context, types.NamespacedName{Namespace: d.Namespace, Name: d.Name}, &d)
	g.Expect(err).NotTo(HaveOccurred())
	err = testContext.Client.Get(testContext.Context, types.NamespacedName{Namespace: d.Namespace, Name: d.Name}, sts)
	g.Expect(err).NotTo(HaveOccurred())

	return d, *sts
}
