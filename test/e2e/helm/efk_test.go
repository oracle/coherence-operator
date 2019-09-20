/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package helm_test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v6/esapi"
	. "github.com/onsi/gomega"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	cohv1 "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	"github.com/oracle/coherence-operator/pkg/fakes"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"net/http"
	"os"
	"testing"

	es6 "github.com/elastic/go-elasticsearch/v6"
	"github.com/oracle/coherence-operator/test/e2e/helper"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"time"
)

// Test installing the Operator with EFK enabled.
func TestOperatorWithEFK(t *testing.T) {
	helmHelper, err := helper.NewOperatorChartHelper()
	if err != nil {
		t.Fatal(err)
	}

	// Create the Operator SDK test context (this will deploy the Operator)
	ctx := helper.CreateTestContext(t)
	// Make sure we defer clean-up (uninstall the operator and Coherence cluster) when we're done
	defer helper.DumpOperatorLogsAndCleanup(t, ctx)

	// Create the values to use
	values := helper.OperatorValues{
		InstallEFK: true,
	}

	assertEFK(t, values, "op1", ctx, helmHelper)
}

// Test installing the Operator with an external EFK stack.
func TestOperatorWithExternalEFK(t *testing.T) {
	helmHelper, err := helper.NewOperatorChartHelper()
	if err != nil {
		t.Fatal(err)
	}

	// Create the Operator SDK test context (this will deploy the Operator)
	ctx := helper.CreateTestContext(t)
	// Make sure we defer clean-up (uninstall the operator and Coherence cluster) when we're done
	defer helper.DumpOperatorLogsAndCleanup(t, ctx)

	// install the External EFK stack
	esHost := installExternalEFK(t, ctx, false)

	// Create the values to use with install EFK disabled and the ES endpoint set to the ES Pod
	values := helper.OperatorValues{
		InstallEFK:            false,
		ElasticsearchEndpoint: &helper.ElasticsearchEndpointSpec{Host: &esHost},
	}

	assertEFK(t, values, "op2", ctx, helmHelper)
}

// Test installing the Operator with an external EFK stack.
func TestOperatorWithExternalEFKAndMonitoringSecret(t *testing.T) {
	helmHelper, err := helper.NewOperatorChartHelper()
	if err != nil {
		t.Fatal(err)
	}

	// Create the Operator SDK test context (this will deploy the Operator)
	ctx := helper.CreateTestContext(t)
	// Make sure we defer clean-up (uninstall the operator and Coherence cluster) when we're done
	defer helper.DumpOperatorLogsAndCleanup(t, ctx)

	// install the External EFK stack and create the monitoring secret
	_ = installExternalEFK(t, ctx, true)

	// Create the values to use with install EFK disabled and No ES endpoint set.
	// The Coherence clusters should get the ES endpoint from the pre-created secret
	values := helper.OperatorValues{
		InstallEFK: false,
	}

	assertEFK(t, values, "op3", ctx, helmHelper)
}

// Assert that the Operator Helm install works and a Coherence Cluster's logs appear in Elasticsearch and Kibana
func assertEFK(t *testing.T, values helper.OperatorValues, releaseName string, ctx *framework.TestCtx, helmHelper *helper.HelmHelper) {
	g := NewGomegaWithT(t)

	namespace := helmHelper.Namespace
	client := helmHelper.KubeClient

	// Create a HelmReleaseManager with a release name and values
	hm, err := helmHelper.NewOperatorHelmReleaseManager(releaseName, &values)
	g.Expect(err).ToNot(HaveOccurred())

	// Defer cleanup (helm delete) to make sure it happens when this method exits
	defer CleanupHelm(t, hm, helmHelper)

	// Install the chart
	_, err = hm.InstallRelease()
	g.Expect(err).ToNot(HaveOccurred())

	// The chart is installed but the Pod(s) may not exist yet so wait for it...
	// (we wait a maximum of 5 minutes, retrying every 5 seconds)
	oPods, err := helper.WaitForOperatorPods(client, namespace, time.Second*5, time.Minute*5)
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(len(oPods)).To(Equal(1))

	// Find the Elasticsearch Pod(s)
	esPods, err := helper.ListPodsWithLabelSelector(client, namespace, "component=elasticsearch")
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(len(esPods)).To(Equal(1))

	// The Elasticsearch Pod(s) exist so get one of them using the k8s client from the helper
	// which is in the client var configured in the suite .go file.
	esPod, err := client.CoreV1().Pods(namespace).Get(esPods[0].Name, metav1.GetOptions{})
	g.Expect(err).ToNot(HaveOccurred())

	// Create an Elasticsearch client
	cl := helper.ElasticSearchClient{Pod: esPod}

	// The chart is installed but the ES Pod we have may not be ready yet so wait for it...
	// (we wait a maximum of 5 minutes, retrying every 5 seconds)
	err = helper.WaitForPodReady(client, esPod.Namespace, esPod.Name, time.Second*5, time.Minute*5)
	g.Expect(err).ToNot(HaveOccurred())

	// Find the kibana Pod(s)
	kPods, err := helper.ListPodsWithLabelSelector(client, namespace, "component=kibana")
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(len(kPods)).To(Equal(1))

	// The Kibana Pod(s) exist so get one of them using the k8s client from the helper
	// which is in the client var configured in the suite .go file.
	kPod, err := client.CoreV1().Pods(namespace).Get(kPods[0].Name, metav1.GetOptions{})
	g.Expect(err).ToNot(HaveOccurred())

	// The chart is installed but the Kibana Pod we have may not be ready yet so wait for it...
	// (we wait a maximum of 5 minutes, retrying every 5 seconds)
	err = helper.WaitForPodReady(client, kPod.Namespace, kPod.Name, time.Second*5, time.Minute*5)
	g.Expect(err).ToNot(HaveOccurred())

	// Deploy the Coherence cluster
	cluster, err := DeployCoherenceCluster(t, ctx, namespace, "coherence-with-fluentd.yaml")
	g.Expect(err).ToNot(HaveOccurred())

	// It can take a while for things to start to appear in Elasticsearch so wait...
	err = cl.WaitForCoherenceIndices(time.Second*5, time.Minute*10, t)
	g.Expect(err).ToNot(HaveOccurred())

	// The rest of the tests are executed as sub-tests.
	// This allows us to run a number of tests and see which fail rather
	// than having one big test method that fails at the first bad assertion
	testCases := []struct {
		name string
		fn   func(t *testing.T)
	}{
		{"ShouldConnectToES", func(t *testing.T) { ShouldConnectToES(t, cl) }},
		{"CoherencePodsShouldHaveFluentdContainer", func(t *testing.T) { CoherencePodsShouldHaveFluentdContainer(t, cluster, helmHelper) }},
		{"ShouldHaveCoherenceClusterIndices", func(t *testing.T) { ShouldHaveCoherenceClusterIndices(t, cl) }},
		{"ShouldHaveLogsFromAllCoherencePods", func(t *testing.T) { ShouldHaveLogsFromAllCoherencePods(t, cl, cluster, helmHelper) }},
		{"ShouldHaveCoherenceClusterIndexInKibana", func(t *testing.T) { ShouldHaveCoherenceClusterIndexInKibana(t, kPod) }},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.fn(t)
		})
	}
}

// With Fluentd enabled the Coherence Cluster Pods should have a sidecar container running fluentd
func CoherencePodsShouldHaveFluentdContainer(t *testing.T, cluster cohv1.CoherenceCluster, helm *helper.HelmHelper) {
	g := NewGomegaWithT(t)

	pods, err := helper.ListCoherencePodsForCluster(helm.KubeClient, helm.Namespace, cluster.Name)
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(len(pods)).ToNot(BeZero())

	for _, pod := range pods {
		containers := make(map[string]corev1.Container)
		for _, c := range pod.Spec.Containers {
			containers[c.Name] = c
		}

		_, ok := containers["fluentd"]
		g.Expect(ok).To(BeTrue(), "Pod "+pod.Name+" does not have a fluentd container")
	}
}

// Assert that it is possible to connect to Elasticsearch
func ShouldConnectToES(t *testing.T, esClient helper.ElasticSearchClient) {
	g := NewGomegaWithT(t)

	res, err := esClient.Query(func(es *es6.Client) (*esapi.Response, error) {
		return es.Info()
	})

	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(res.IsError()).To(BeFalse(), fmt.Sprintf("Error response from ES %s", res.String()))
}

// Assert that the Coherence cluster indices exist in Elasticsearch
func ShouldHaveCoherenceClusterIndices(t *testing.T, esClient helper.ElasticSearchClient) {
	g := NewGomegaWithT(t)

	res, err := esClient.Query(func(es *es6.Client) (*esapi.Response, error) {
		return es.Indices.Get([]string{"coherence-cluster-*"})
	})

	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(res.IsError()).To(BeFalse(), fmt.Sprintf("Error response from ES %s", res.String()))

	m := make(map[string]interface{})
	err = json.NewDecoder(res.Body).Decode(&m)
	g.Expect(len(m)).NotTo(BeZero())
}

// Assert that there are log messages in Elasticsearch for each Coherence Cluster member Pod
func ShouldHaveLogsFromAllCoherencePods(t *testing.T, esClient helper.ElasticSearchClient, cluster cohv1.CoherenceCluster, helm *helper.HelmHelper) {
	g := NewGomegaWithT(t)

	pods, err := helper.ListCoherencePodsForCluster(helm.KubeClient, helm.Namespace, cluster.Name)
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(len(pods)).ToNot(BeZero())

	// An ES search query template to find log messages with a specific "host" field value
	esQuery := `{"query": {"match": {"host": "%s"}}}`

	// Iterate over the Pods in the cluster and make sure that there are log messages from each one.
	// The "host" field in the log message will be the Pod name.
	for _, pod := range pods {
		fn := helper.NewESSearchFunction("coherence-cluster-*", fmt.Sprintf(esQuery, pod.Name))

		// Do the check in a loop as it can take a while for logs to appear for the Pods
		err = wait.Poll(time.Second*5, time.Minute*2, func() (done bool, err error) {
			result := helper.ESSearchResult{}
			err = esClient.QueryAndParse(&result, fn)
			if err != nil {
				return false, err
			}
			return result.Size() != 0, nil
		})

		g.Expect(err).ToNot(HaveOccurred(), "Did not find logs in Elasticsearch for Pod "+pod.Name)
	}
}

// Assert that the Coherence Cluster index pattern is present in Kibana
func ShouldHaveCoherenceClusterIndexInKibana(t *testing.T, kibana *corev1.Pod) {
	g := NewGomegaWithT(t)

	// start a port forwarder to the Kibana Pod
	fwd, ports, err := helper.StartPortForwarderForPod(kibana)
	g.Expect(err).ToNot(HaveOccurred())

	// Defer closing the PortForwarder so we clean-up properly
	defer fwd.Close()

	// The ReST port in the Kibana container spec is named "kibana"
	port := ports["kibana"]

	// Query Kibana for the Coherence Cluster index pattern
	url := fmt.Sprintf("http://127.0.0.1:%d//api/saved_objects/index-pattern/%s", port, helper.KibanaIndexPatternCoherenceCluster)

	// Do the query in a loop as it might take time to appear
	err = wait.Poll(time.Second*5, time.Minute*2, func() (done bool, err error) {
		res, err := http.Get(url)
		if err != nil {
			return false, err
		}
		return res.StatusCode == http.StatusOK, nil
	})
}

// ----- helper methods -------------------------------------------------

// Install an external EFK stack.
// We do this by doing a fake Operator Helm install and then pulling out
// the bits we need for the ELK stack and creating them using the k8s client.
// Returns the endpoint to use to contact Elasticsearch
func installExternalEFK(t *testing.T, ctx *framework.TestCtx, includeSecret bool) string {
	f := framework.Global
	g := NewGomegaWithT(t)

	namespace := helper.GetTestNamespace()

	// We use the fake Operator Helm install with EFK enabled to obtain
	// the yaml to use to install the EFK stack.
	mgrFake := fakes.NewFakeManager()
	values := helper.OperatorValues{}

	chartDir, err := helper.FindOperatorHelmChartDir()
	g.Expect(err).ToNot(HaveOccurred())

	err = values.LoadFromYaml(chartDir + string(os.PathSeparator) + "values.yaml")
	g.Expect(err).ToNot(HaveOccurred())

	values.InstallEFK = true

	helm := fakes.NewFakeHelm(mgrFake, nil, nil)
	result, err := helm.FakeOperatorHelmInstall(mgrFake, namespace, values)
	g.Expect(err).ToNot(HaveOccurred())

	// Find the bits we need from the Helm install result
	esdp := &appsv1.Deployment{}
	essvc := &corev1.Service{}
	kbdp := &appsv1.Deployment{}
	kbsvc := &corev1.Service{}
	sccm := &corev1.ConfigMap{}
	kbdbcm := &corev1.ConfigMap{}

	err = result.Get("elasticsearch", esdp)
	g.Expect(err).ToNot(HaveOccurred())
	err = result.Get("elasticsearch", essvc)
	g.Expect(err).ToNot(HaveOccurred())
	err = result.Get("kibana", kbdp)
	g.Expect(err).ToNot(HaveOccurred())
	err = result.Get("kibana", kbsvc)
	g.Expect(err).ToNot(HaveOccurred())
	err = result.Get("operator-coherence-operator-importscript", sccm)
	g.Expect(err).ToNot(HaveOccurred())
	err = result.Get("operator-coherence-operator-dashboards", kbdbcm)
	g.Expect(err).ToNot(HaveOccurred())

	esdp.SetNamespace(namespace)
	essvc.SetNamespace(namespace)
	kbdp.SetNamespace(namespace)
	kbsvc.SetNamespace(namespace)
	sccm.SetNamespace(namespace)
	kbdbcm.SetNamespace(namespace)

	// deploy the EFK stack
	err = f.Client.Create(context.TODO(), sccm, helper.DefaultCleanup(ctx))
	g.Expect(err).ToNot(HaveOccurred())
	err = f.Client.Create(context.TODO(), kbdbcm, helper.DefaultCleanup(ctx))
	g.Expect(err).ToNot(HaveOccurred())
	err = f.Client.Create(context.TODO(), essvc, helper.DefaultCleanup(ctx))
	g.Expect(err).ToNot(HaveOccurred())
	err = f.Client.Create(context.TODO(), kbsvc, helper.DefaultCleanup(ctx))
	g.Expect(err).ToNot(HaveOccurred())
	err = f.Client.Create(context.TODO(), esdp, helper.DefaultCleanup(ctx))
	g.Expect(err).ToNot(HaveOccurred())
	err = f.Client.Create(context.TODO(), kbdp, helper.DefaultCleanup(ctx))
	g.Expect(err).ToNot(HaveOccurred())

	if includeSecret {
		sec := &corev1.Secret{}
		err = result.Get("coherence-monitoring-config", sec)
		g.Expect(err).ToNot(HaveOccurred())

		sec.SetNamespace(namespace)

		err = f.Client.Create(context.TODO(), sec, helper.DefaultCleanup(ctx))
		g.Expect(err).ToNot(HaveOccurred())
	}

	return essvc.GetName()
}
