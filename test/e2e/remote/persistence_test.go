/*
 * Copyright (c) 2020, 2025, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package remote

import (
	goctx "context"
	"encoding/json"
	"fmt"
	. "github.com/onsi/gomega"
	v1 "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	"golang.org/x/net/context"
	"io"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"
)

type snapshotActionType int

const (
	canaryServiceName = "CanaryService"

	Create snapshotActionType = iota
	Recover
	Delete
)

// Deploy a Coherence resource with persistence enabled (this should enable active persistence).
// A PVC should be created for the StatefulSet. Create data in some caches, delete the deployment,
// re-deploy the deployment and assert that the data is recovered.
func TestActivePersistence(t *testing.T) {
	// Make sure we defer clean-up when we're done!!
	testContext.CleanupAfterTest(t)
	assertPersistence("persistence-active.yaml", "persistence-volume", false, false, true, t)
}

// Deploy a Coherence resource with the minimal default configuration. Persistence will be on-demand.
// Put data in a cache, take a snapshot, delete the data, recover the snapshot,
// assert that the data is recovered.
func TestOnDemandPersistence(t *testing.T) {
	// Make sure we defer clean-up when we're done!!
	testContext.CleanupAfterTest(t)
	assertPersistence("persistence-on-demand.yaml", "", true, true, false, t)
}

// Deploy a Coherence resource with snapshot enabled. Persistence will be on-demand,
// a PVC will be created for the StatefulSet to use for snapshots. Put data in a cache, take a snapshot,
// delete the deployment, re-deploy the deployment, recover the snapshot, assert that the data is recovered.
func TestSnapshotPersistence(t *testing.T) {
	// Make sure we defer clean-up when we're done!!
	testContext.CleanupAfterTest(t)
	assertPersistence("persistence-snapshot.yaml", "snapshot-volume", true, false, true, t)
}

// Deploy a Coherence resource with both persistence and snapshot configured. Persistence will be active,
// a PVC will be created for the StatefulSet. Put data in a cache, take a snapshot,
// delete the deployment, re-deploy the deployment, recover the snapshot, assert that the data is recovered.
func TestSnapshotWithActivePersistence(t *testing.T) {
	// Make sure we defer clean-up when we're done!!
	testContext.CleanupAfterTest(t)
	assertPersistence("persistence-active-snapshot.yaml", "snapshot-volume", true, true, true, t)
}

// Deploy a Coherence resource with both persistence and snapshot configured and a securityContext so
// that the container is not running as the root user.
// Persistence will be active, a PVC will be created for the StatefulSet. Put data in a cache, take a snapshot,
// delete the deployment, re-deploy the deployment, recover the snapshot, assert that the data is recovered.
func TestSnapshotWithActivePersistenceWithSecurityContext(t *testing.T) {
	// Make sure we defer clean-up when we're done!!
	testContext.CleanupAfterTest(t)
	assertPersistence("persistence-active-snapshot-security.yaml", "snapshot-volume", true, true, true, t)
}

func assertPersistence(yamlFile, pVolName string, isSnapshot, isClearCanary, isRestart bool, t *testing.T) {
	g := NewGomegaWithT(t)
	ns := helper.GetTestNamespace()

	deployment, pods := ensurePods(g, yamlFile, ns)

	// check the pvc is created for the given volume
	if pVolName != "" {
		fmt.Printf("Checking existence of PVC %s\n", pVolName)
		pvcName := ""
		for _, vol := range pods[0].Spec.Volumes {
			if vol.Name == pVolName {
				if vol.PersistentVolumeClaim != nil {
					pvcName = vol.PersistentVolumeClaim.ClaimName
				}
				break
			}
		}

		// check the pvc is created
		g.Expect(pvcName).NotTo(Equal(""))
		pvc := testContext.KubeClient.CoreV1().PersistentVolumeClaims(pvcName)
		g.Expect(pvc).ShouldNot(BeNil())
	}

	// create data in some caches
	fmt.Println("Initialise Canary Cache")
	err := helper.StartCanary(testContext, ns, deployment.GetName())
	g.Expect(err).NotTo(HaveOccurred())

	if isSnapshot {
		// take a snapshot
		fmt.Println("Creating Snapshot")
		err = processSnapshotRequest(pods[0], Create)
		g.Expect(err).NotTo(HaveOccurred())

		defer processSnapshotRequestBlind(pods[0], Delete)
	}

	if isClearCanary {
		fmt.Println("Clearing Canary Cache")
		// delete the data
		err = helper.ClearCanary(testContext, ns, deployment.GetName())
		g.Expect(err).NotTo(HaveOccurred())
	}

	localStorageRestartEnv := os.Getenv("LOCAL_STORAGE_RESTART")
	localStorageRestart, err := strconv.ParseBool(localStorageRestartEnv)
	if err != nil {
		localStorageRestart = false
	}
	// restart Coherence may be on a different instance, local storage will not work
	if isRestart && !localStorageRestart {
		// dump the pod logs before deleting
		helper.DumpPodsForTest(testContext, t)
		// delete the deployment
		fmt.Println("Deleting Coherence deployment")
		err = helper.WaitForCoherenceCleanup(testContext, ns)
		g.Expect(err).NotTo(HaveOccurred())

		// re-deploy the deployment
		fmt.Println("Re-starting Coherence deployment")
		deployment, pods = ensurePods(g, yamlFile, ns)
	}

	if isSnapshot {
		// recover the snapshot
		fmt.Println("Recovering Snapshot")
		err = processSnapshotRequest(pods[0], Recover)
		g.Expect(err).NotTo(HaveOccurred())
	}

	// assert that the data is recovered
	fmt.Println("Checking Canary cache")
	err = helper.CheckCanaryEventuallyGood(testContext, ns, deployment.GetName())
	g.Expect(err).NotTo(HaveOccurred())

	// cleanup the data
	_ = helper.ClearCanary(testContext, ns, deployment.GetName())
}

func ensurePods(g *GomegaWithT, yamlFile, ns string) (v1.Coherence, []corev1.Pod) {
	deployment, err := helper.NewSingleCoherenceFromYaml(ns, yamlFile)
	g.Expect(err).NotTo(HaveOccurred())

	d, _ := json.Marshal(deployment)
	fmt.Printf("Persistence Test installing deployment:\n%s\n", string(d))

	err = testContext.Client.Create(goctx.TODO(), &deployment)
	g.Expect(err).NotTo(HaveOccurred())

	_, err = helper.WaitForStatefulSetForDeployment(testContext, ns, &deployment, helper.RetryInterval, helper.Timeout)
	g.Expect(err).NotTo(HaveOccurred())

	// Get the list of Pods
	pods, err := helper.ListCoherencePodsForDeployment(testContext, ns, deployment.GetName())
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(len(pods)).Should(BeNumerically(">", 0))

	return deployment, pods
}

func processSnapshotRequestBlind(pod corev1.Pod, actionType snapshotActionType) {
	err := processSnapshotRequest(pod, actionType)
	if err != nil {
		fmt.Printf("Failed to process snapshot request (type=%d) due to %s\n", actionType, err.Error())
	}
}

func processSnapshotRequest(pod corev1.Pod, actionType snapshotActionType) error {
	pf, ports, err := helper.StartPortForwarderForPod(&pod)
	if err != nil {
		return err
	}

	defer pf.Close()

	url := fmt.Sprintf("http://127.0.0.1:%d/management/coherence/cluster/services/%s/persistence/snapshots/snapshotOne",
		ports[v1.PortNameManagement], canaryServiceName)
	httpMethod := "POST"
	if actionType == Delete {
		httpMethod = "DELETE"
	}
	if actionType == Recover {
		url += "/recover"
	}

	client := &http.Client{}
	var resp *http.Response
	var req *http.Request
	// try a max of 5 times
	for i := 0; i < 5; i++ {
		req, err = http.NewRequest(httpMethod, url, nil)
		if err == nil {
			resp, err = client.Do(req)
			if resp != nil {
				_ = resp.Body.Close()
			}
			if err == nil {
				break
			}
		}
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("snapshot request returned non-200 status %d", resp.StatusCode)
	}

	// wait for idle
	err = wait.PollUntilContextTimeout(context.Background(), helper.RetryInterval, helper.Timeout, true, func(context.Context) (done bool, err error) {
		url = fmt.Sprintf("http://127.0.0.1:%d/management/coherence/cluster/services/%s/persistence?fields=operationStatus",
			ports[v1.PortNameManagement], canaryServiceName)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Printf("Cannot create idle check request: %v\n", url)
			return false, err
		}
		resp, err = client.Do(req)
		if resp != nil {
			defer resp.Body.Close()
		}
		if err != nil {
			fmt.Printf("Error in send idle check request: %v\n", url)
			return false, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Idle check request with incorrect status code: %v\n", resp.StatusCode)
			return false, err
		}

		bs, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Print("Error in reading idle check response")
			return false, err
		}

		var data map[string]interface{}
		if err = json.Unmarshal(bs, &data); err != nil {
			fmt.Print("Error in unmarshalling idle check response")
			return false, nil
		}
		opStatus := data["operationStatus"]
		fmt.Printf("Persistence opStatus: %v\n", opStatus)
		if opStatus == "Idle" {
			return true, nil
		}
		return false, nil
	})

	return err
}
