/*
 * Copyright (c) 2020, 2022, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package prometheus

import (
	"encoding/json"
	"fmt"
	. "github.com/onsi/gomega"
	coh "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	monitoring "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	client "github.com/prometheus-operator/prometheus-operator/pkg/client/versioned/typed/monitoring/v1"
	"io/ioutil"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
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

	promClient, err := client.NewForConfig(testContext.Config)
	g.Expect(err).NotTo(HaveOccurred())

	AssertPrometheus(t, "prometheus-test.yaml", promPod, promClient)
}

func AssertPrometheus(t *testing.T, yamlFile string, promPod corev1.Pod, promClient *client.MonitoringV1Client) {
	g := NewGomegaWithT(t)

	ShouldGetPrometheusConfig(t, promPod)

	// Deploy the Coherence cluster
	deployments, cohPods := helper.AssertDeployments(testContext, t, yamlFile)
	deployment := deployments["test"]

	err := ShouldEventuallyHaveServiceMonitor(t, deployment.Namespace, "test-metrics", promClient, 10*time.Second, 5*time.Minute)
	g.Expect(err).NotTo(HaveOccurred())

	// Wait for Coherence metrics to appear in Prometheus
	ShouldEventuallySeeClusterMetrics(t, promPod, cohPods)

	// Ensure that we can see the deployments size metric
	ShouldGetClusterSizeMetric(t, promPod)

	// Ensure we can update the Coherence deployment and cause the ServiceMonitor to be updated
	ShouldPatchServiceMonitor(t, deployment, promClient)
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
	r, err := PrometheusAPIRequest(pod, "/api/v1/status/config")
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

func ShouldPatchServiceMonitor(t *testing.T, deployment coh.Coherence, promClient *client.MonitoringV1Client) {
	g := NewGomegaWithT(t)

	current := &coh.Coherence{}
	err := testContext.Client.Get(testContext.Context, types.NamespacedName{Namespace: deployment.Namespace, Name: deployment.Name}, current)
	g.Expect(err).NotTo(HaveOccurred())

	// update the ServiceMonitor interval to cause an update
	current.Spec.Ports[0].ServiceMonitor.Interval = "10s"
	err = testContext.Client.Update(testContext.Context, current)
	g.Expect(err).NotTo(HaveOccurred())

	err = ShouldEventuallyHaveServiceMonitorWithState(t, deployment.Namespace, "test-metrics", hasInterval, promClient, 10*time.Second, 5*time.Minute)
	g.Expect(err).NotTo(HaveOccurred())

}

func ShouldEventuallyHaveServiceMonitor(t *testing.T, namespace, name string, promClient *client.MonitoringV1Client, retryInterval, timeout time.Duration) error {
	return ShouldEventuallyHaveServiceMonitorWithState(t, namespace, name, alwaysTrue, promClient, retryInterval, timeout)
}

type ServiceMonitorPredicate func(*testing.T, *monitoring.ServiceMonitor) bool

func alwaysTrue(*testing.T, *monitoring.ServiceMonitor) bool {
	return true
}

func hasInterval(t *testing.T, sm *monitoring.ServiceMonitor) bool {
	if len(sm.Spec.Endpoints) > 0 && sm.Spec.Endpoints[0].Interval == "10s" {
		return true
	}
	t.Logf("Waiting for availability of ServiceMonitor resource %s - with endpoint interval of 10s", sm.Name)
	return false
}

func ShouldEventuallyHaveServiceMonitorWithState(t *testing.T, namespace, name string, predicate ServiceMonitorPredicate, promClient *client.MonitoringV1Client, retryInterval, timeout time.Duration) error {
	var sm *monitoring.ServiceMonitor

	err := wait.Poll(retryInterval, timeout, func() (done bool, err error) {
		sm, err = promClient.ServiceMonitors(namespace).Get(testContext.Context, name, v1.GetOptions{})
		if err != nil {
			if apierrors.IsNotFound(err) {
				t.Logf("Waiting for availability of ServiceMonitor resource %s - NotFound", name)
				return false, nil
			}
			t.Logf("Waiting for availability of ServiceMonitor resource %s - %s", name, err.Error())
			return false, nil
		}
		if predicate(t, sm) {
			return true, nil
		}
		t.Logf("Waiting for availability of ServiceMonitor resource %s - %s to match predicate", name, err.Error())
		return false, nil
	})

	return err
}

func PrometheusQuery(t *testing.T, pod corev1.Pod, query string, result interface{}) error {
	r, err := PrometheusAPIRequest(pod, "/api/v1/query?query="+query)
	if err != nil {
		return err
	}

	if r.Status != "success" {
		return fmt.Errorf("prometheus returned a non-success status '%s' Data='%s'", r.Status, string(r.Data))
	}
	t.Logf("Query: /api/v1/query?query=%s Result: status=%s %s", query, r.Status, string(r.Data))
	return r.GetData(result)
}

func PrometheusAPIRequest(pod corev1.Pod, path string) (*PrometheusAPIResult, error) {
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

	result := &PrometheusAPIResult{}
	err = json.Unmarshal(data, result)
	return result, err
}

type PrometheusAPIResult struct {
	Status string          `json:"status"`
	Data   json.RawMessage `json:"data"`
}

func (r *PrometheusAPIResult) GetData(v interface{}) error {
	if r == nil {
		return fmt.Errorf("called on a nil PrometheusAPIResult")
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
