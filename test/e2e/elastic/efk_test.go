/*
 * Copyright (c) 2020, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package elastic

import (
	"context"
	"fmt"
	es "github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	. "github.com/onsi/gomega"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"testing"
	"time"
)

func TestElasticSearch(t *testing.T) {
	// Ensure that everything is cleaned up after the test!
	testContext.CleanupAfterTest(t)

	g := NewGomegaWithT(t)

	esPod, kPod := AssertElasticsearchInstalled(testContext, t)

	// Create the ConfigMap with the Fluentd config
	cm := &corev1.ConfigMap{}
	err := helper.LoadFromYamlFile("efk-configmap.yaml", cm)
	g.Expect(err).ToNot(HaveOccurred())

	cm.SetNamespace(helper.GetTestNamespace())

	err = testContext.Client.Create(context.TODO(), cm)
	g.Expect(err).ToNot(HaveOccurred())

	// Create an Elasticsearch client
	cl := ESClient{Pod: esPod}
	// Test that the client connects
	shouldConnectToES(t, cl)

	// Deploy the Coherence cluster
	_, pods := helper.AssertDeployments(testContext, t, "efk-test.yaml")
	assertAllHaveFluentdContainers(t, pods)

	// It can take a while for things to start to appear in Elasticsearch so wait...
	err = cl.WaitForCoherenceIndices(t, time.Second*5, time.Minute*5)
	g.Expect(err).ToNot(HaveOccurred())

	assertHaveLogsFromAllCoherencePods(t, cl, pods)
	assertCoherenceClusterIndexInKibana(t, kPod)
}

// Assert that it is possible to connect to Elasticsearch
func shouldConnectToES(t *testing.T, cl ESClient) {
	g := NewGomegaWithT(t)
	_, _ = AssertElasticsearchInstalled(testContext, t)

	res, err := cl.Query(func(es *es.Client) (*esapi.Response, error) {
		return es.Info()
	})

	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(res.IsError()).To(BeFalse(), fmt.Sprintf("Error response from ES %s", res.String()))
}

// Assert that all of the Pods have a Fluentd container
func assertAllHaveFluentdContainers(t *testing.T, pods []corev1.Pod) {
	g := NewGomegaWithT(t)

	for _, pod := range pods {
		found := false
		for _, c := range pod.Spec.Containers {
			if c.Name == "fluentd" {
				found = true
				break
			}
		}
		g.Expect(found).To(BeTrue(), fmt.Sprintf("Pod %s does not have a Fluentd container", pod.Name))
	}
}

// Assert that there are log messages in Elasticsearch for each Coherence Cluster member Pod
func assertHaveLogsFromAllCoherencePods(t *testing.T, cl ESClient, pods []corev1.Pod) {
	g := NewGomegaWithT(t)

	// An ES search query template to find log messages with a specific "host" field value
	esQuery := `{"query": {"match": {"host": "%s"}}}`

	// Iterate over the Pods in the cluster and make sure that there are log messages from each one.
	// The "host" field in the log message will be the Pod name.
	for _, pod := range pods {
		t.Logf("Looking for ES logs for Pod %s", pod.Name)
		fn := NewESSearchFunction("coherence-cluster-*", fmt.Sprintf(esQuery, pod.Name))

		// Do the check in a loop as it can take a while for logs to appear for the Pods
		err := wait.Poll(time.Second*5, time.Minute*2, func() (done bool, err error) {
			result := ESSearchResult{}
			err = cl.QueryAndParse(&result, fn)
			if err != nil {
				return false, err
			}
			return result.Size() != 0, nil
		})

		g.Expect(err).ToNot(HaveOccurred(), "Did not find logs in Elasticsearch for Pod "+pod.Name)
	}
}
