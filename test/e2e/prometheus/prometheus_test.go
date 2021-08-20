/*
 * Copyright (c) 2020, 2021, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package prometheus

import (
	"encoding/json"
	"fmt"
	. "github.com/onsi/gomega"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	"io/ioutil"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestPrometheus(t *testing.T) {
	// Ensure that everything is cleaned up after the test!
	testContext.CleanupAfterTest(t)

	g := NewGomegaWithT(t)
	ok, promPod, err := IsPrometheusInstalled()
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(ok).To(BeTrue(), "Cannot find any Prometheus Pods - this test requires Prometheus to have been installed")

	AssertPrometheus(t, "prometheus-test.yaml", promPod)
}

func AssertPrometheus(t *testing.T, yamlFile string, promPod corev1.Pod) {
	ShouldGetPrometheusConfig(t, promPod)

	// Deploy the Coherence cluster
	_, cohPods := helper.AssertDeployments(testContext, t, yamlFile)

	// Wait for Coherence metrics to appear in Prometheus
	ShouldEventuallySeeClusterMetrics(t, promPod, cohPods)

	// Ensure that we can see the deployments size metric
	ShouldGetClusterSizeMetric(t, promPod)
}

func IsPrometheusInstalled() (bool, corev1.Pod, error) {
	promNamespace := helper.GetPrometheusNamespace()
	promPods, err := helper.ListPodsWithLabelSelector(testContext, promNamespace, "app=prometheus")
	if err != nil || len(promPods) == 0 {
		return false, corev1.Pod{}, err
	}
	return true, promPods[0], nil
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

	err := wait.Poll(time.Second*20, time.Minute*15, func() (done bool, err error) {
		result := PrometheusVector{}
		err = PrometheusQuery(t, promPod, "up", &result)
		if err != nil {
			return false, err
		}

		var namespace string
		m := make(map[string]bool)
		for _, pod := range cohPods {
			namespace = pod.Namespace
			m[pod.Name] = false
		}

		for _, v := range result.Result {
			if v.Labels["namespace"] == namespace {
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
	err := PrometheusQuery(t, pod, "vendor:coherence_cluster_size", &metrics)
	g.Expect(err).NotTo(HaveOccurred())
}

func PrometheusQuery(t *testing.T, pod corev1.Pod, query string, result interface{}) error {
	r, err := PrometheusApiRequest(pod, "/api/v1/query?query="+query)
	if err != nil {
		return err
	}

	if r.Status != "success" {
		return fmt.Errorf("prometheus returned a non-success status '%s' Data='%s'", r.Status, string(r.Data))
	} else {
		t.Logf("Query: /api/v1/query?query=%s Result: status=%s %s", query, r.Status, string(r.Data))
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
