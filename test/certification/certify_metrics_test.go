/*
 * Copyright (c) 2020, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 *
 */

package certification

import (
	"context"
	"fmt"
	. "github.com/onsi/gomega"
	v1 "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
	"net/http"
	"testing"
	"time"
)

func TestCertifyMetricsDefaultPort(t *testing.T) {
	// Ensure that everything is cleaned up after the test!
	testContext.CleanupAfterTest(t)
	g := NewGomegaWithT(t)

	ns := helper.GetTestNamespace()

	d := &v1.Coherence{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: ns,
			Name:      "metric-default",
		},
		Spec: v1.CoherenceResourceSpec{
			Coherence: &v1.CoherenceSpec{
				Metrics: &v1.PortSpecWithSSL{
					Enabled: pointer.BoolPtr(true),
				},
			},
			Ports: []v1.NamedPortSpec{
				{Name: v1.PortNameMetrics},
			},
		},
	}

	err := testContext.Client.Create(context.TODO(), d)
	g.Expect(err).NotTo(HaveOccurred())

	_, err = helper.WaitForStatefulSetForDeployment(testContext, ns, d, time.Second*10, time.Minute*5)
	g.Expect(err).NotTo(HaveOccurred())

	// Get the deployment Pods
	pods, err := helper.ListCoherencePodsForDeployment(testContext, ns, "metric-default")
	g.Expect(err).NotTo(HaveOccurred())

	// the default number of replicas is 3 so the first pod should be able to be used
	// Get only the first pod and add port forwarding
	var pod = pods[1]
	pf, ports, err := helper.StartPortForwarderForPod(&pod)
	g.Expect(err).NotTo(HaveOccurred())
	defer pf.Close()

	fmt.Println("Available ports:")
	for key, value := range ports {
		fmt.Println(key, value)
	}

	url := fmt.Sprintf("%s://127.0.0.1:%d/metrics", "http", ports[v1.PortNameMetrics])

	var resp *http.Response
	client := &http.Client{}

	println("Connecting with: ", url)
	// try a max of 5 times
	for i := 0; i < 5; i++ {
		resp, err = client.Get(url)
		if err == nil {
			break
		}
		time.Sleep(5 * time.Second)
	}
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(resp.StatusCode).To(Equal(http.StatusOK))

}

func TestCertifyMetricsNonStandardPort(t *testing.T) {
	// Ensure that everything is cleaned up after the test!
	testContext.CleanupAfterTest(t)
	g := NewGomegaWithT(t)

	ns := helper.GetTestNamespace()

	d := &v1.Coherence{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: ns,
			Name:      "metric-nondefault",
		},
		Spec: v1.CoherenceResourceSpec{
			Coherence: &v1.CoherenceSpec{
				Metrics: &v1.PortSpecWithSSL{
					Enabled: pointer.BoolPtr(true),
					Port:    pointer.Int32Ptr(9619),
				},
			},
			Ports: []v1.NamedPortSpec{
				{Name: v1.PortNameMetrics,
					Port: 9619},
			},
		},
	}

	err := testContext.Client.Create(context.TODO(), d)
	g.Expect(err).NotTo(HaveOccurred())

	_, err = helper.WaitForStatefulSetForDeployment(testContext, ns, d, time.Second*10, time.Minute*5)
	g.Expect(err).NotTo(HaveOccurred())

	// Get the deployment Pods
	pods, err := helper.ListCoherencePodsForDeployment(testContext, ns, "metric-nondefault")
	g.Expect(err).NotTo(HaveOccurred())

	// the default number of replicas is 3 so the first pod should be able to be used
	// Get only the first pod and add port forwarding
	var pod = pods[1]
	pf, ports, err := helper.StartPortForwarderForPod(&pod)
	g.Expect(err).NotTo(HaveOccurred())
	defer pf.Close()

	fmt.Println("Available ports:")
	for key, value := range ports {
		fmt.Println(key, value)
	}
	//url := fmt.Sprintf("%s://127.0.0.1:%d/metrics", "http", ports[coh.PortNameMetrics])
	url := fmt.Sprintf("%s://127.0.0.1:%d/metrics", "http", ports[v1.PortNameMetrics])

	var resp *http.Response
	client := &http.Client{}

	println("Connecting with: ", url)
	// try a max of 5 times
	for i := 0; i < 5; i++ {
		resp, err = client.Get(url)
		if err == nil {
			break
		}
		time.Sleep(5 * time.Second)
	}
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(resp.StatusCode).To(Equal(http.StatusOK))

}
