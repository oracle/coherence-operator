/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package helm_test

import (
	"encoding/json"
	"fmt"
	. "github.com/onsi/gomega"
	"io/ioutil"
	corev1 "k8s.io/api/core/v1"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/rest"
	"k8s.io/utils/pointer"
	"net/http"
	"strings"
	"testing"

	"github.com/oracle/coherence-operator/test/e2e/helper"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"time"
)

func TestOperatorWithPrometheus(t *testing.T) {
	helmHelper, err := helper.NewOperatorChartHelper()
	if err != nil {
		t.Fatal(err)
	}

	// Create the Operator SDK test context (this will deploy the Operator)
	ctx := helper.CreateTestContext(t)
	// Make sure we defer clean-up (uninstall the operator and Coherence cluster) when we're done
	defer helper.DumpOperatorLogsAndCleanup(t, ctx)

	hasCRDs, err := HasPrometheusCRDs(helmHelper.Manager.GetConfig())
	fmt.Printf("Check for Prometheus CRDs - found=%t\n", hasCRDs)

	// Create the values to use to install the operator with Prometheus but without Grafana
	values := helper.OperatorValues{
		Prometheusoperator: &helper.PrometheusOperatorSpec{
			Enabled: pointer.BoolPtr(true),
			PrometheusOperator: &helper.PrometheusOp{
				CreateCustomResource: pointer.BoolPtr(!hasCRDs),
			},
			Prometheus: &helper.Prometheus{
				PrometheusSpec: &helper.PrometheusSpec{ScrapeInterval: pointer.StringPtr("5s")},
			},
			Grafana: &helper.Grafana{
				Enabled: pointer.BoolPtr(false),
			},
		},
	}

	g := NewGomegaWithT(t)
	namespace := helmHelper.Namespace
	client := helmHelper.KubeClient

	// Create a HelmReleaseManager with a release name and values
	hm, err := helmHelper.NewOperatorHelmReleaseManager("op", &values)
	g.Expect(err).ToNot(HaveOccurred())

	// Defer cleanup (helm delete) to make sure it happens when this method exits
	defer CleanupHelm(t, hm, helmHelper)

	// Install the chart
	_, err = hm.InstallRelease()
	g.Expect(err).ToNot(HaveOccurred())

	// Wait for the Prometheus Pod(s)
	promPods, err := helper.WaitForPodsWithLabel(client, namespace, "app=prometheus", 1, time.Second*5, time.Minute*1)

	if err != nil {
		fmt.Printf("Found zero Prometheus Pods. Pods in namespace %s:\n", namespace)
		lst, _ := client.CoreV1().Pods(namespace).List(metav1.ListOptions{})
		if len(lst.Items) == 0 {
			fmt.Println("Found zero pods in namespace")
		} else {
			for _, p := range lst.Items {
				fmt.Printf("Pod: %s labels %v\n", p.Name, p.Labels)
			}
		}
		t.Fatal("Did not find any Prometheus Pods")
	}

	g.Expect(err).ToNot(HaveOccurred())

	// Get one of the Prometheus Pods - there should only be one anyway
	promPod := promPods[0]

	// The chart is installed but the Prometheus Pod we have may not be ready yet so wait for it...
	// (we wait a maximum of 5 minutes, retrying every 5 seconds)
	err = helper.WaitForPodReady(client, promPod.Namespace, promPod.Name, time.Second*5, time.Minute*5)
	g.Expect(err).ToNot(HaveOccurred())

	// Deploy the Coherence cluster with the metrics port exposed on a service.
	// We need to do this because the default Prometheus install uses a ServiceMonitor
	cluster, err := DeployCoherenceCluster(t, ctx, namespace, "coherence-with-metrics.yaml")
	g.Expect(err).ToNot(HaveOccurred())

	// ensure that we can hit the Prometheus API
	ShouldGetPrometheusConfig(t, promPod)

	// ensure that Coherence metrics eventually appear from all of the Coherence Pods
	cohPods, err := helper.ListCoherencePodsForCluster(client, namespace, cluster.Name)
	ShouldEventuallySeeClusterMetrics(t, promPod, cohPods)

	// Ensure that we can see the cluster size metric
	ShouldGetClusterSizeMetric(t, promPod)
}

// Ensure that the Prometheus status/config endpoint can be accessed.
func ShouldGetPrometheusConfig(t *testing.T, pod corev1.Pod) {
	g := NewGomegaWithT(t)
	r, err := PrometheusApiRequest(pod, "/api/v1/status/config")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(r.Status).To(Equal("success"))
}

// Ensure that Coherence metrics for all Coherence Pods eventually appear.
func ShouldEventuallySeeClusterMetrics(t *testing.T, promPod corev1.Pod, cohPods []corev1.Pod) {
	g := NewGomegaWithT(t)

	err := wait.Poll(time.Second*5, time.Minute*5, func() (done bool, err error) {
		result := PrometheusVector{}
		err = PrometheusQuery(promPod, "up", &result)
		if err != nil {
			return false, err
		}

		m := make(map[string]bool)
		for _, pod := range cohPods {
			m[pod.Name] = false
		}

		for _, v := range result.Result {
			if v.Labels["job"] == "coherence-service-metrics" {
				name := v.Labels["pod"]
				m[name] = true
			}
		}

		for _, pod := range cohPods {
			if m[pod.Name] == false {
				return false, nil
			}
		}
		return true, nil
	})

	g.Expect(err).NotTo(HaveOccurred())
}

// Ensure that we can get the cluster size metric
func ShouldGetClusterSizeMetric(t *testing.T, pod corev1.Pod) {
	g := NewGomegaWithT(t)

	metrics := PrometheusVector{}
	err := PrometheusQuery(pod, "vendor:coherence_cluster_size", &metrics)
	g.Expect(err).NotTo(HaveOccurred())
}

func PrometheusQuery(pod corev1.Pod, query string, result interface{}) error {
	r, err := PrometheusApiRequest(pod, "/api/v1/query?query="+query)
	if err != nil {
		return err
	}

	if r.Status != "success" {
		return fmt.Errorf("prometheus returned a non-success status '%s' Data='%s'", r.Status, string(r.Data))
	}

	return r.GetData(result)
}

func PrometheusApiRequest(pod corev1.Pod, path string) (*PrometheusApiResult, error) {
	// Start the port forwarder for the Pod
	pf, ports, err := helper.StartPortForwarderForPod(&pod)
	if err != nil {
		return nil, err
	}
	// Defer closing the port forwarder to ensure we clean up
	defer pf.Close()

	var sep string
	if strings.HasPrefix(path, "/") {
		sep = ""
	} else {
		sep = "/"
	}

	url := fmt.Sprintf("http://127.0.0.1:%d%s%s", ports["web"], sep, path)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	result := &PrometheusApiResult{}
	err = json.Unmarshal(data, result)
	return result, err
}

func HasPrometheusCRDs(cfg *rest.Config) (bool, error) {
	promCrds := []string{
		"alertmanagers.monitoring.coreos.com",
		"prometheuses.monitoring.coreos.com",
		"prometheusrules.monitoring.coreos.com",
		"servicemonitors.monitoring.coreos.com",
	}

	cl, err := apiextensionsclient.NewForConfig(cfg)
	if err != nil {
		return false, err
	}

	count := 0
	for _, crd := range promCrds {
		crds := cl.ApiextensionsV1beta1().CustomResourceDefinitions()
		_, err := crds.Get(crd, metav1.GetOptions{})
		if err != nil && !errors.IsNotFound(err) {
			return false, err
		}

		if err == nil {
			count++
		}
	}

	numCRDs := len(promCrds)

	switch count {
	case 0:
		return false, nil
	case numCRDs:
		return true, nil
	default:
		return false, fmt.Errorf("only %d of %d Prometheus CRDs are installed - either all %d must be installed or none must be installed", count, numCRDs, numCRDs)
	}
}

type PrometheusApiResult struct {
	Status string          `json:"status"`
	Data   json.RawMessage `json:"data"`
}

func (r *PrometheusApiResult) GetData(v interface{}) error {
	if r == nil {
		return fmt.Errorf("called on a nil PrometheusApiResult")
	}

	return json.Unmarshal(r.Data, v)
}

type PrometheusVector struct {
	ResultType string             `json:"resultType"`
	Result     []PrometheusMetric `json:"result"`
}

type PrometheusMetric struct {
	Labels map[string]string `json:"metric"`
	Value  []interface{}     `json:"value"`
}

func (m *PrometheusMetric) GetName() string {
	if m == nil {
		return ""
	}
	return m.Labels["__name__"]
}
