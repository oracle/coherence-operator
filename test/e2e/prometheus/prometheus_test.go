/*
 * Copyright (c) 2020, 2026, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package prometheus

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	. "github.com/onsi/gomega"
	coh "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	monitoring "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	client "github.com/prometheus-operator/prometheus-operator/pkg/client/versioned/typed/monitoring/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"
)

const clusterSizeMetricNameRegex = "coherence_cluster_size|vendor:coherence_cluster_size|vendor_Coherence_Cluster_Size"

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

	// Wait for Prometheus to scrape every Coherence Pod successfully
	ShouldEventuallyHaveHealthyTargets(t, promPod, cohPods)

	// Ensure that we can see the deployments size metric
	ShouldEventuallyGetClusterSizeMetric(t, promPod, cohPods, time.Time{})

	// Ensure we can update the Coherence deployment and cause the ServiceMonitor to be updated
	ShouldPatchServiceMonitor(t, deployment, promClient)
	ShouldChangeServiceMonitorLabelDrop(t, deployment, promClient)
	metricsNotBefore := time.Now()
	ShouldEventuallyHaveHealthyTargets(t, promPod, cohPods)
	ShouldEventuallyGetClusterSizeMetric(t, promPod, cohPods, metricsNotBefore)
}

func IsPrometheusInstalled() (bool, corev1.Pod, error) {
	promNamespace := helper.GetPrometheusNamespace()
	promPods, err := helper.ListPodsWithLabelSelector(testContext, promNamespace, "app.kubernetes.io/name=prometheus")
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

// Ensure that Prometheus can successfully scrape every Coherence Pod.
func ShouldEventuallyHaveHealthyTargets(t *testing.T, promPod corev1.Pod, cohPods []corev1.Pod) {
	g := NewGomegaWithT(t)
	g.Expect(cohPods).NotTo(BeEmpty())

	namespace := cohPods[0].Namespace
	for _, pod := range cohPods {
		g.Expect(pod.Namespace).To(Equal(namespace))
	}

	t.Logf("Waiting for healthy Prometheus targets for all Coherence Pods")

	err := wait.PollUntilContextTimeout(testContext.Context, time.Second*20, time.Minute*2, true, func(context.Context) (done bool, err error) {
		result := PrometheusVector{}
		err = PrometheusQuery(t, promPod, fmt.Sprintf("up{namespace=%q}", namespace), &result)
		if err != nil {
			return false, err
		}

		m := make(map[string]bool)
		for _, pod := range cohPods {
			m[pod.Name] = false
		}

		for _, v := range result.Result {
			name := v.Labels["pod"]
			if _, ok := m[name]; !ok {
				continue
			}
			value, err := v.GetValue()
			if err != nil {
				return false, err
			}
			if value == 1 {
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

// Ensure that Prometheus has a cluster size metric from every Coherence Pod.
// When notBefore is set, the source sample must have been scraped after that time.
func ShouldEventuallyGetClusterSizeMetric(t *testing.T, promPod corev1.Pod, cohPods []corev1.Pod, notBefore time.Time) {
	g := NewGomegaWithT(t)
	g.Expect(cohPods).NotTo(BeEmpty())

	namespace := cohPods[0].Namespace
	expectedPods := make(map[string]bool)
	for _, pod := range cohPods {
		g.Expect(pod.Namespace).To(Equal(namespace))
		expectedPods[pod.Name] = false
	}

	query := fmt.Sprintf("timestamp({__name__=~%q,namespace=%q})", clusterSizeMetricNameRegex, namespace)
	var notBeforeUnix float64
	if !notBefore.IsZero() {
		notBeforeUnix = float64(notBefore.UnixNano()) / float64(time.Second)
	}
	err := wait.PollUntilContextTimeout(testContext.Context, 5*time.Second, 2*time.Minute, true, func(context.Context) (bool, error) {
		metrics := PrometheusVector{}
		if err := PrometheusQuery(t, promPod, query, &metrics); err != nil {
			return false, err
		}

		found := make(map[string]bool, len(expectedPods))
		for podName := range expectedPods {
			found[podName] = false
		}
		for _, metric := range metrics.Result {
			podName := metric.Labels["pod"]
			if _, ok := found[podName]; !ok {
				continue
			}
			sampleTimestamp, err := metric.GetValue()
			if err != nil {
				return false, err
			}
			if notBefore.IsZero() || sampleTimestamp >= notBeforeUnix {
				found[podName] = true
			}
		}

		for _, pod := range cohPods {
			if !found[pod.Name] {
				return false, nil
			}
		}
		return true, nil
	})

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

// Exercise removal of the last target rule and restoration through the parent CR.
func ShouldChangeServiceMonitorLabelDrop(t *testing.T, deployment coh.Coherence, promClient *client.MonitoringV1Client) {
	g := NewGomegaWithT(t)
	// Bypass the manager cache for read-after-write and no-roll verification.
	reader := testContext.Manager.GetAPIReader()
	key := types.NamespacedName{Namespace: deployment.Namespace, Name: deployment.Name}
	broad := monitoring.RelabelConfig{Action: "labeldrop", Regex: "(endpoint|instance|job|service)"}
	narrow := monitoring.RelabelConfig{Action: "labeldrop", Regex: "(endpoint|service)"}
	disabled := false
	states := []struct {
		name             string
		flag             *bool
		custom, expected []monitoring.RelabelConfig
	}{
		{name: "baseline", expected: []monitoring.RelabelConfig{broad}},
		{name: "remove last rule", flag: &disabled},
		{name: "customer configuration", flag: &disabled, custom: []monitoring.RelabelConfig{narrow}, expected: []monitoring.RelabelConfig{narrow}},
		{name: "restore default", custom: []monitoring.RelabelConfig{narrow}, expected: []monitoring.RelabelConfig{narrow, broad}},
	}
	type podState struct {
		UID      types.UID
		Restarts map[string]int32
	}
	var baseline map[string]podState
	var revision string
	for i, state := range states {
		t.Logf("ServiceMonitor label drop: %s", state.name)
		current := &coh.Coherence{}
		g.Expect(reader.Get(testContext.Context, key, current)).To(Succeed())
		portIndex := -1
		for index, port := range current.Spec.Ports {
			if port.Name == "metrics" {
				portIndex = index
				break
			}
		}
		g.Expect(portIndex).To(BeNumerically(">=", 0))
		if i > 0 {
			current.Spec.Ports[portIndex].ServiceMonitor.UseDefaultLabelDrop = state.flag
			current.Spec.Ports[portIndex].ServiceMonitor.Relabelings = state.custom
			g.Expect(testContext.Client.Update(testContext.Context, current)).To(Succeed())
		}
		current = &coh.Coherence{}
		g.Expect(reader.Get(testContext.Context, key, current)).To(Succeed())
		// Read-back proves the new optional field survived schema admission/storage.
		g.Expect(current.Spec.Ports[portIndex].ServiceMonitor.UseDefaultLabelDrop).To(Equal(state.flag))
		g.Expect(current.Spec.Ports[portIndex].ServiceMonitor.Relabelings).To(Equal(state.custom))
		expectedHash := current.GetGenerationString()
		sts := &appsv1.StatefulSet{}
		err := wait.PollUntilContextTimeout(testContext.Context, time.Second, 5*time.Minute, true, func(ctx context.Context) (bool, error) {
			sm, err := promClient.ServiceMonitors(key.Namespace).Get(ctx, "test-metrics", v1.GetOptions{})
			if err != nil {
				return false, err
			}
			if len(sm.Spec.Endpoints) != 1 || sm.Spec.Endpoints[0].Port != "metrics" {
				return false, nil
			}
			rules := sm.Spec.Endpoints[0].RelabelConfigs
			if len(rules) != len(state.expected) {
				return false, nil
			}
			for index := range rules {
				if !reflect.DeepEqual(rules[index], state.expected[index]) {
					return false, nil
				}
			}
			if err = reader.Get(ctx, key, sts); err != nil {
				return false, err
			}
			return sts.Labels[coh.LabelCoherenceHash] == expectedHash && sts.Status.ObservedGeneration >= sts.Generation, nil
		})
		g.Expect(err).NotTo(HaveOccurred(), state.name)
		pods := &corev1.PodList{}
		g.Expect(reader.List(testContext.Context, pods, runtimeclient.InNamespace(key.Namespace), runtimeclient.MatchingLabels(sts.Spec.Selector.MatchLabels))).To(Succeed())
		snapshot := map[string]podState{}
		for _, pod := range pods.Items {
			restarts := map[string]int32{}
			for _, container := range pod.Status.ContainerStatuses {
				restarts["container/"+container.Name] = container.RestartCount
			}
			for _, container := range pod.Status.InitContainerStatuses {
				restarts["init/"+container.Name] = container.RestartCount
			}
			snapshot[pod.Name] = podState{UID: pod.UID, Restarts: restarts}
		}
		g.Expect(snapshot).NotTo(BeEmpty())
		if i == 0 {
			baseline = snapshot
			revision = sts.Status.UpdateRevision
			g.Expect(revision).NotTo(BeEmpty())
		} else {
			g.Expect(snapshot).To(Equal(baseline), state.name)
			g.Expect(sts.Status.UpdateRevision).To(Equal(revision), state.name)
		}
	}
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

	t.Logf("Waiting for ServiceMonitor resource %s/%s to be available", namespace, name)

	err := wait.PollUntilContextTimeout(context.Background(), retryInterval, timeout, true, func(context.Context) (done bool, err error) {
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
		if err == nil {
			t.Logf("Waiting for availability of ServiceMonitor resource %s to match predicate", name)
		} else {
			t.Logf("Waiting for availability of ServiceMonitor resource %s to match predicate. error: %s", name, err.Error())
		}
		return false, nil
	})

	return err
}

func PrometheusQuery(t *testing.T, pod corev1.Pod, query string, result interface{}) error {
	r, err := PrometheusAPIRequest(pod, "/api/v1/query?query="+url.QueryEscape(query))
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
	pf, ports, err := helper.StartPortForwarderForPodWithBackoff(&pod)
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

	url := fmt.Sprintf("http://%s:%d%s%s", pf.Hostname, ports["web"], sep, path)
	resp, err := http.Get(url)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(resp.Body)
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

func (m *PrometheusMetric) GetValue() (float64, error) {
	if m == nil || len(m.Value) < 2 {
		return 0, fmt.Errorf("prometheus metric has no sample value")
	}

	switch value := m.Value[1].(type) {
	case string:
		return strconv.ParseFloat(value, 64)
	case float64:
		return value, nil
	default:
		return 0, fmt.Errorf("unexpected Prometheus sample value type %T", value)
	}
}

func (m *PrometheusMetric) GetName() string {
	if m == nil {
		return ""
	}
	return m.Labels["__name__"]
}
