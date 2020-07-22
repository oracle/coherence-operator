/*
 * Copyright (c) 2020 Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package certification

import (
	"context"
	"fmt"
	. "github.com/onsi/gomega"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	coh "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	v1 "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
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

