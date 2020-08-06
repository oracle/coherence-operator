/*
 * Copyright (c) 2020 Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package certification

import (
	"context"
	goctx "context"
	"encoding/json"
	"flag"
	"fmt"
	. "github.com/onsi/gomega"
	"github.com/operator-framework/operator-sdk/pkg/log/zap"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	coh "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	v1 "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	"io/ioutil"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/klog"
	"k8s.io/utils/pointer"
	"net/http"
	"os"
	"os/exec"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"strconv"
	"strings"
	"testing"
	"time"
)

type snapshotActionType int

const (
	canaryServiceName = "CanaryService"

	create snapshotActionType = iota
	recover
	delete
)

func TestCertifyMinimalSpec(t *testing.T) {
	g := NewGomegaWithT(t)

	f := framework.Global
	ctx := helper.CreateTestContext(t)
	defer helper.DumpOperatorLogsAndCleanup(t, ctx)

	ns := helper.GetTestNamespace()
	d := &v1.Coherence{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: ns,
			Name:      "certify-minimal",
		},
	}

	err := f.Client.Create(context.TODO(), d, helper.DefaultCleanup(ctx))
	g.Expect(err).NotTo(HaveOccurred())

	_, err = helper.WaitForStatefulSetForDeployment(f.KubeClient, ns, d, time.Second*10, time.Minute*5, t)
	g.Expect(err).NotTo(HaveOccurred())
}

func TestCertifyScaling(t *testing.T) {
	g := NewGomegaWithT(t)

	f := framework.Global
	ctx := helper.CreateTestContext(t)
	defer helper.DumpOperatorLogsAndCleanup(t, ctx)

	ns := helper.GetTestNamespace()
	d := &v1.Coherence{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: ns,
			Name:      "certify-scale",
		},
		Spec: v1.CoherenceResourceSpec{
			Replicas: pointer.Int32Ptr(1),
			ReadinessProbe: &v1.ReadinessProbeSpec{
				InitialDelaySeconds: pointer.Int32Ptr(10),
				PeriodSeconds:       pointer.Int32Ptr(10),
			},
		},
	}

	// Start with one replica
	err := f.Client.Create(context.TODO(), d, helper.DefaultCleanup(ctx))
	g.Expect(err).NotTo(HaveOccurred())
	_, err = helper.WaitForStatefulSetForDeployment(f.KubeClient, ns, d, time.Second*10, time.Minute*5, t)
	g.Expect(err).NotTo(HaveOccurred())

	// Scale Up to three
	err = scale(t, ns, d.Name, 3)
	g.Expect(err).NotTo(HaveOccurred())
	_, err = helper.WaitForStatefulSet(f.KubeClient, ns, d.Name, 3, time.Second*10, time.Minute*5, t)
	g.Expect(err).NotTo(HaveOccurred())

	// Scale down to one
	err = scale(t, ns, d.Name, 1)
	g.Expect(err).NotTo(HaveOccurred())
	_, err = helper.WaitForStatefulSet(f.KubeClient, ns, d.Name, 1, time.Second*10, time.Minute*5, t)
	g.Expect(err).NotTo(HaveOccurred())
}

func scale(t *testing.T, namespace, name string, replicas int32) error {
	cmd := exec.Command("kubectl", "-n", namespace, "scale", fmt.Sprintf("--replicas=%d", replicas), "coherence/"+name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	t.Log("Executing Scale Command: " + strings.Join(cmd.Args, " "))
	return cmd.Run()
}

func TestCertifyManagementDefaultPort(t *testing.T) {
	g := NewGomegaWithT(t)
	f := framework.Global

	ctx := helper.CreateTestContext(t)
	defer helper.DumpOperatorLogsAndCleanup(t, ctx)

	ns := helper.GetTestNamespace()

	d := &v1.Coherence{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: ns,
			Name:      "management-default",
		},
		Spec: coh.CoherenceResourceSpec{
			Coherence: &coh.CoherenceSpec{
				Management: &coh.PortSpecWithSSL{
					Enabled: pointer.BoolPtr(true),
				},
			},
			Ports: []v1.NamedPortSpec{
				{Name: coh.PortNameManagement},
			},
		},
	}

	err := f.Client.Create(context.TODO(), d, helper.DefaultCleanup(ctx))
	g.Expect(err).NotTo(HaveOccurred())

	_, err = helper.WaitForStatefulSetForDeployment(f.KubeClient, ns, d, time.Second*10, time.Minute*5, t)
	g.Expect(err).NotTo(HaveOccurred())

	// Get the deployment Pods
	pods, err := helper.ListCoherencePodsForDeployment(f.KubeClient, ns, "management-default")
	g.Expect(err).NotTo(HaveOccurred())

	// the default number of replicas is 3 so the first pod should be able to be used
	// Get only the first pod and add port forwarding
	var pod = pods[1]
	pf, ports, err := helper.StartPortForwarderForPod(&pod)
	defer pf.Close()

	println("Available ports:")
	for key, value := range ports {
		fmt.Println(key, value)
	}

	url := fmt.Sprintf("%s://127.0.0.1:%d/management/coherence/cluster", "http", ports[coh.PortNameManagement])

	var resp *http.Response
	client := &http.Client{}

	println("Connecting with: ", url)
	// try a max of 5 times
	for i := 0; i < 5; i++ {
		resp, err = client.Get(url)
		if err == nil || !true {
			break
		}
		time.Sleep(5 * time.Second)
	}
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(resp.StatusCode).To(Equal(http.StatusOK))

}
func TestCertifyManagementNonStandardPort(t *testing.T) {
	g := NewGomegaWithT(t)
	f := framework.Global

	ctx := helper.CreateTestContext(t)
	defer helper.DumpOperatorLogsAndCleanup(t, ctx)

	ns := helper.GetTestNamespace()

	d := &v1.Coherence{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: ns,
			Name:      "management-nondefault",
		},
		Spec: coh.CoherenceResourceSpec{
			Coherence: &coh.CoherenceSpec{
				Management: &coh.PortSpecWithSSL{
					Enabled: pointer.BoolPtr(true),
					Port:    pointer.Int32Ptr(30009),
				},
			},
			Ports: []v1.NamedPortSpec{
				{Name: coh.PortNameManagement,
					Port: 30009},
			},
		},
	}

	err := f.Client.Create(context.TODO(), d, helper.DefaultCleanup(ctx))
	g.Expect(err).NotTo(HaveOccurred())

	_, err = helper.WaitForStatefulSetForDeployment(f.KubeClient, ns, d, time.Second*10, time.Minute*5, t)
	g.Expect(err).NotTo(HaveOccurred())

	// Get the deployment Pods
	pods, err := helper.ListCoherencePodsForDeployment(f.KubeClient, ns, "management-nondefault")
	g.Expect(err).NotTo(HaveOccurred())

	// the default number of replicas is 3 so the first pod should be able to be used
	// Get only the first pod and add port forwarding
	var pod = pods[1]
	pf, ports, err := helper.StartPortForwarderForPod(&pod)
	defer pf.Close()

	println("Available ports:")
	for key, value := range ports {
		fmt.Println(key, value)
	}
	//url := fmt.Sprintf("%s://127.0.0.1:%d/metrics", "http", ports[coh.PortNameMetrics])
	url := fmt.Sprintf("%s://127.0.0.1:%d/management/coherence/cluster", "http", ports[coh.PortNameManagement])

	var resp *http.Response
	client := &http.Client{}

	println("Connecting with: ", url)
	// try a max of 5 times
	for i := 0; i < 5; i++ {
		resp, err = client.Get(url)
		if err == nil || !true {
			break
		}
		time.Sleep(5 * time.Second)
	}
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(resp.StatusCode).To(Equal(http.StatusOK))

}

func TestCertifyMetricsDefaultPort(t *testing.T) {
	g := NewGomegaWithT(t)
	f := framework.Global

	ctx := helper.CreateTestContext(t)
	defer helper.DumpOperatorLogsAndCleanup(t, ctx)

	ns := helper.GetTestNamespace()

	d := &v1.Coherence{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: ns,
			Name:      "metric-default",
		},
		Spec: coh.CoherenceResourceSpec{
			Coherence: &coh.CoherenceSpec{
				Metrics: &coh.PortSpecWithSSL{
					Enabled: pointer.BoolPtr(true),
				},
			},
			Ports: []v1.NamedPortSpec{
				{Name: coh.PortNameMetrics},
			},
		},
	}

	err := f.Client.Create(context.TODO(), d, helper.DefaultCleanup(ctx))
	g.Expect(err).NotTo(HaveOccurred())

	_, err = helper.WaitForStatefulSetForDeployment(f.KubeClient, ns, d, time.Second*10, time.Minute*5, t)
	g.Expect(err).NotTo(HaveOccurred())

	// Get the deployment Pods
	pods, err := helper.ListCoherencePodsForDeployment(f.KubeClient, ns, "metric-default")
	g.Expect(err).NotTo(HaveOccurred())

	// the default number of replicas is 3 so the first pod should be able to be used
	// Get only the first pod and add port forwarding
	var pod = pods[1]
	pf, ports, err := helper.StartPortForwarderForPod(&pod)
	defer pf.Close()

	fmt.Println("Available ports:")
	for key, value := range ports {
		fmt.Println(key, value)
	}

	url := fmt.Sprintf("%s://127.0.0.1:%d/metrics", "http", ports[coh.PortNameMetrics])

	var resp *http.Response
	client := &http.Client{}

	println("Connecting with: ", url)
	// try a max of 5 times
	for i := 0; i < 5; i++ {
		resp, err = client.Get(url)
		if err == nil || !true {
			break
		}
		time.Sleep(5 * time.Second)
	}
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(resp.StatusCode).To(Equal(http.StatusOK))

}

func TestCertifyMetricsNonStandardPort(t *testing.T) {
	g := NewGomegaWithT(t)
	f := framework.Global

	ctx := helper.CreateTestContext(t)
	defer helper.DumpOperatorLogsAndCleanup(t, ctx)

	ns := helper.GetTestNamespace()

	d := &v1.Coherence{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: ns,
			Name:      "metric-nondefault",
		},
		Spec: coh.CoherenceResourceSpec{
			Coherence: &coh.CoherenceSpec{
				Metrics: &coh.PortSpecWithSSL{
					Enabled: pointer.BoolPtr(true),
					Port:    pointer.Int32Ptr(9619),
				},
			},
			Ports: []v1.NamedPortSpec{
				{Name: coh.PortNameMetrics,
					Port: 9619},
			},
		},
	}

	err := f.Client.Create(context.TODO(), d, helper.DefaultCleanup(ctx))
	g.Expect(err).NotTo(HaveOccurred())

	_, err = helper.WaitForStatefulSetForDeployment(f.KubeClient, ns, d, time.Second*10, time.Minute*5, t)
	g.Expect(err).NotTo(HaveOccurred())

	// Get the deployment Pods
	pods, err := helper.ListCoherencePodsForDeployment(f.KubeClient, ns, "metric-nondefault")
	g.Expect(err).NotTo(HaveOccurred())

	// the default number of replicas is 3 so the first pod should be able to be used
	// Get only the first pod and add port forwarding
	var pod = pods[1]
	pf, ports, err := helper.StartPortForwarderForPod(&pod)
	defer pf.Close()

	fmt.Println("Available ports:")
	for key, value := range ports {
		fmt.Println(key, value)
	}
	//url := fmt.Sprintf("%s://127.0.0.1:%d/metrics", "http", ports[coh.PortNameMetrics])
	url := fmt.Sprintf("%s://127.0.0.1:%d/metrics", "http", ports[coh.PortNameMetrics])

	var resp *http.Response
	client := &http.Client{}

	println("Connecting with: ", url)
	// try a max of 5 times
	for i := 0; i < 5; i++ {
		resp, err = client.Get(url)
		if err == nil || !true {
			break
		}
		time.Sleep(5 * time.Second)
	}
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(resp.StatusCode).To(Equal(http.StatusOK))

}

// Deploy a Coherence resource with persistence enabled (this should enable active persistence).
// A PVC should be created for the StatefulSet. Create data in some caches, delete the deployment,
// re-deploy the deployment and assert that the data is recovered.
func TestActivePersistenceScaleUpAndDown(t *testing.T) {
	var yamlFile = "persistence-active-1.yaml"
	var pVolName = "persistence-volume"

	g := NewGomegaWithT(t)

	logf.SetLogger(zap.Logger())

	flags := &flag.FlagSet{}
	klog.InitFlags(flags)
	_ = flags.Set("v", "4")

	f := framework.Global
	ctx := helper.CreateTestContext(t)
	defer helper.DumpOperatorLogsAndCleanup(t, ctx)

	ns, err := ctx.GetWatchNamespace()
	g.Expect(err).NotTo(HaveOccurred())

	deployment, pods := ensurePods(g, ctx, yamlFile, ns, t)

	// check the pvc is created for the given volume
	if pVolName != "" {
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
		pvc := f.KubeClient.CoreV1().PersistentVolumeClaims(pvcName)
		g.Expect(pvc).ShouldNot(BeNil())
	}

	// create data in some caches
	err = helper.StartCanary(ns, deployment.GetName())
	g.Expect(err).NotTo(HaveOccurred())

	// Start with one replica
	_, err = helper.WaitForStatefulSetForDeployment(f.KubeClient, ns, &deployment, time.Second*10, time.Minute*5, t)
	g.Expect(err).NotTo(HaveOccurred())

	// Scale Up to three
	err = scale(t, ns, deployment.Name, 3)
	g.Expect(err).NotTo(HaveOccurred())
	_, err = helper.WaitForStatefulSet(f.KubeClient, ns, deployment.Name, 3, time.Second*10, time.Minute*5, t)
	g.Expect(err).NotTo(HaveOccurred())

	// Scale down to one
	err = scale(t, ns, deployment.Name, 1)
	g.Expect(err).NotTo(HaveOccurred())
	_, err = helper.WaitForStatefulSet(f.KubeClient, ns, deployment.Name, 1, time.Second*10, time.Minute*5, t)
	g.Expect(err).NotTo(HaveOccurred())

	// assert that the data is recovered
	err = helper.CheckCanary(ns, deployment.GetName())
	g.Expect(err).NotTo(HaveOccurred())

	// cleanup the data
	_ = helper.ClearCanary(ns, deployment.GetName())
}

// Deploy a Coherence resource with persistence enabled (this should enable active persistence).
// A PVC should be created for the StatefulSet. Create data in some caches, scale down to 0 the deployment,
// re-deploy the deployment and assert that the data is recovered.
func TestActivePersistenceScaleDownAndUp(t *testing.T) {
	var yamlFile = "persistence-active-3.yaml"
	var pVolName = "persistence-volume"

	g := NewGomegaWithT(t)

	logf.SetLogger(zap.Logger())

	flags := &flag.FlagSet{}
	klog.InitFlags(flags)
	_ = flags.Set("v", "4")

	f := framework.Global
	ctx := helper.CreateTestContext(t)
	defer helper.DumpOperatorLogsAndCleanup(t, ctx)

	ns, err := ctx.GetWatchNamespace()
	g.Expect(err).NotTo(HaveOccurred())

	deployment, pods := ensurePods(g, ctx, yamlFile, ns, t)

	// check the pvc is created for the given volume
	if pVolName != "" {
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
		pvc := f.KubeClient.CoreV1().PersistentVolumeClaims(pvcName)
		g.Expect(pvc).ShouldNot(BeNil())
	}

	// create data in some caches
	err = helper.StartCanary(ns, deployment.GetName())
	g.Expect(err).NotTo(HaveOccurred())

	// Start with three replicas
	_, err = helper.WaitForStatefulSetForDeployment(f.KubeClient, ns, &deployment, time.Second*10, time.Minute*5, t)
	g.Expect(err).NotTo(HaveOccurred())

	// Scale Down to One
	err = scale(t, ns, deployment.Name, 1)
	g.Expect(err).NotTo(HaveOccurred())
	_, err = helper.WaitForStatefulSet(f.KubeClient, ns, deployment.Name, 1, time.Second*10, time.Minute*5, t)
	g.Expect(err).NotTo(HaveOccurred())

	// Scale back up to Three
	err = scale(t, ns, deployment.Name, 3)
	g.Expect(err).NotTo(HaveOccurred())
	_, err = helper.WaitForStatefulSet(f.KubeClient, ns, deployment.Name, 3, time.Second*10, time.Minute*5, t)
	g.Expect(err).NotTo(HaveOccurred())

	// assert that the data is recovered
	err = helper.CheckCanary(ns, deployment.GetName())
	g.Expect(err).NotTo(HaveOccurred())

	// cleanup the data
	_ = helper.ClearCanary(ns, deployment.GetName())
}

// Deploy a Coherence resource with persistence enabled (this should enable active persistence).
// A PVC should be created for the StatefulSet. Create data in some caches, scale down to 0 the deployment,
// re-deploy the deployment and assert that the data is recovered.
func TestActivePersistenceTakeAndRestoreSnapshot(t *testing.T) {
	var yamlFile = "persistence-snapshot.yaml"
	var pVolName = "snapshot-volume"

	g := NewGomegaWithT(t)

	logf.SetLogger(zap.Logger())

	flags := &flag.FlagSet{}
	klog.InitFlags(flags)
	_ = flags.Set("v", "4")

	f := framework.Global
	ctx := helper.CreateTestContext(t)
	defer helper.DumpOperatorLogsAndCleanup(t, ctx)

	ns, err := ctx.GetWatchNamespace()
	g.Expect(err).NotTo(HaveOccurred())

	deployment, pods := ensurePods(g, ctx, yamlFile, ns, t)

	// check the pvc is created for the given volume
	if pVolName != "" {
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
		pvc := f.KubeClient.CoreV1().PersistentVolumeClaims(pvcName)
		g.Expect(pvc).ShouldNot(BeNil())
	}

	// create data in some caches
	err = helper.StartCanary(ns, deployment.GetName())
	g.Expect(err).NotTo(HaveOccurred())

	// take a snapshot
	err = processSnapshotRequest(pods[0], create)
	g.Expect(err).NotTo(HaveOccurred())

	defer processSnapshotRequest(pods[0], delete)

	localStorageRestartEnv := os.Getenv("LOCAL_STORAGE_RESTART")
	localStorageRestart, err := strconv.ParseBool(localStorageRestartEnv)
	if err != nil {
		localStorageRestart = false
	}
	// restart Coherence may be on a different instance, local storage will not work
	if !localStorageRestart {
		fmt.Println("Restarting...")
		// dump the pod logs before deleting
		helper.DumpPodsForTest(t, ctx)
		// delete the deployment
		err = helper.WaitForCoherenceCleanup(f, ns)
		g.Expect(err).NotTo(HaveOccurred())

		// re-deploy the deployment
		deployment, pods = ensurePods(g, ctx, yamlFile, ns, t)
	}

	// recover the snapshot
	err = processSnapshotRequest(pods[0], recover)
	g.Expect(err).NotTo(HaveOccurred())

	// assert that the data is recovered
	err = helper.CheckCanary(ns, deployment.GetName())
	g.Expect(err).NotTo(HaveOccurred())

	// cleanup the data
	_ = helper.ClearCanary(ns, deployment.GetName())
}

func ensurePods(g *GomegaWithT, ctx *framework.Context, yamlFile, ns string, t *testing.T) (v1.Coherence, []corev1.Pod) {
	f := framework.Global

	deployment, err := helper.NewSingleCoherenceFromYaml(ns, yamlFile)
	g.Expect(err).NotTo(HaveOccurred())

	d, _ := json.Marshal(deployment)
	fmt.Printf("Persistence Test installing deployment:\n%s\n", string(d))

	err = f.Client.Create(goctx.TODO(), &deployment, helper.DefaultCleanup(ctx))
	g.Expect(err).NotTo(HaveOccurred())

	_, err = helper.WaitForStatefulSetForDeployment(f.KubeClient, ns, &deployment, helper.RetryInterval, helper.Timeout, t)
	g.Expect(err).NotTo(HaveOccurred())

	// Get the list of Pods
	pods, err := helper.ListCoherencePodsForDeployment(f.KubeClient, ns, deployment.GetName())
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(len(pods)).Should(BeNumerically(">", 0))

	return deployment, pods
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
	if actionType == delete {
		httpMethod = "DELETE"
	}
	if actionType == recover {
		url = url + "/recover"
	}

	client := &http.Client{}
	var resp *http.Response
	var req *http.Request
	// try a max of 5 times
	for i := 0; i < 5; i++ {
		req, err = http.NewRequest(httpMethod, url, nil)
		if err == nil {
			resp, err = client.Do(req)
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
	err = wait.Poll(helper.RetryInterval, helper.Timeout, func() (done bool, err error) {
		url = fmt.Sprintf("http://127.0.0.1:%d/management/coherence/cluster/services/%s/persistence?fields=operationStatus",
			ports[v1.PortNameManagement], canaryServiceName)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Printf("Cannot create idle check request: %v\n", url)
			return false, err
		}
		resp, err = client.Do(req)
		if err != nil {
			fmt.Printf("Error in send idle check request: %v\n", url)
			return false, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Idle check request with incorrect status code: %v\n", resp.StatusCode)
			return false, err
		}

		bs, err := ioutil.ReadAll(resp.Body)
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
