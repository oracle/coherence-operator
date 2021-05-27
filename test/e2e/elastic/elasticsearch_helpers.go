/*
 * Copyright (c) 2020, 2021, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package elastic

import (
	"bytes"
	"encoding/json"
	"fmt"
	es "github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	. "github.com/onsi/gomega"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

// AssertElasticsearchInstalled determines whether Elasticsearch and Kibana are installed returning
// the installed Elasticsearch and Kibana Pods
func AssertElasticsearchInstalled(ctx helper.TestContext, t *testing.T) (*corev1.Pod, *corev1.Pod) {
	g := NewGomegaWithT(t)

	var esPod *corev1.Pod
	var kPod *corev1.Pod

	namespace := helper.GetTestNamespace()

	// Find the Elasticsearch Pods
	pods, err := helper.ListPodsWithLabelSelector(ctx, namespace, "app=elasticsearch-master")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(len(pods)).NotTo(BeZero(), "No Elasticsearch Pods found")
	esPod = &pods[0]

	// Find the Kibana Pods
	pods, err = helper.ListPodsWithLabelSelector(ctx, namespace, "app=kibana")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(len(pods)).NotTo(BeZero(), "No Kibana Pods found")
	kPod = &pods[0]

	return esPod, kPod
}

// Assert that the Coherence Cluster index pattern is present in Kibana
func assertCoherenceClusterIndexInKibana(t *testing.T, kPod *corev1.Pod) {
	g := NewGomegaWithT(t)

	// start a port forwarder to the Kibana Pod
	fwd, ports, err := helper.StartPortForwarderForPod(kPod)
	g.Expect(err).ToNot(HaveOccurred())

	// Defer closing the PortForwarder so we clean-up properly
	defer fwd.Close()

	// The ReST port in the Kibana container spec is named "kPod"
	port := ports["kPod"]

	indexPattern := os.Getenv("KIBANA_INDEX_PATTERN")
	if indexPattern == "" {
		indexPattern = "6abb1220-3feb-11e9-a9a3-4b1c09db6e6a"
	}

	// Query Kibana for the Coherence Cluster index pattern
	url := fmt.Sprintf("http://127.0.0.1:%d//api/saved_objects/index-pattern/%s", port, indexPattern)

	// Do the query in a loop as it might take time to appear
	err = wait.Poll(time.Second*5, time.Minute*2, func() (done bool, err error) {
		res, err2 := http.Get(url)
		if err2 != nil {
			return false, err2
		}
		return res.StatusCode == http.StatusOK, nil
	})
}

// ESClient is an elastic search client associated to an Elasticsearch Pod
type ESClient struct {
	Pod *corev1.Pod
}

// ESFunction is a function that takes an Elasticsearch client and executes an API call
type ESFunction func(es *es.Client) (*esapi.Response, error)

func NewESSearchFunction(index, query string) ESFunction {
	return func(es *es.Client) (*esapi.Response, error) {
		var buf bytes.Buffer
		buf.Write([]byte(query))
		return es.Search(
			es.Search.WithIndex(index),
			es.Search.WithBody(&buf),
			es.Search.WithTrackTotalHits(true),
			es.Search.WithPretty())
	}
}

// WaitForCoherenceIndices waits for a Coherence indices to appear in Elasticsearch.
func (c *ESClient) WaitForCoherenceIndices(t *testing.T, retryInterval, timeout time.Duration) error {

	// A query function to retrieve ES indices
	fn := func(es *es.Client) (*esapi.Response, error) {
		return es.Cat.Indices(es.Cat.Indices.WithPretty())
	}

	var res *esapi.Response

	err := wait.Poll(retryInterval, timeout, func() (done bool, err error) {
		res, err = c.Query(fn)
		if err != nil {
			t.Logf("Waiting for Coherence indices in Elasticsearch - %s\n", err.Error())
			return false, err
		}
		ok := strings.Contains(strings.ToLower(res.String()), "coherence-cluster-")
		return ok, nil
	})

	if err != nil {
		var s string
		if res == nil {
			s = "response is nil"
		} else {
			s = res.String()
		}
		t.Logf("Error waiting for Coherence indices in Elasticsearch - error: '%s' response:\n%s\n", err.Error(), s)
	}

	return err
}

// Query performs an Elasticsearch API call
func (c *ESClient) Query(fn ESFunction) (*esapi.Response, error) {
	if c == nil {
		return nil, fmt.Errorf("this ESClient is nil")
	}

	fwd, ports, err := helper.StartPortForwarderForPod(c.Pod)
	if err != nil {
		return nil, err
	}

	// Defer closing the PortForwarder so we clean-up properly
	defer fwd.Close()

	// The ReST port in the Elasticsearch container spec is named "rest"
	port := ports["http"]

	esHost := fmt.Sprintf("http://127.0.0.1:%d", port)
	cl, err := es.NewClient(es.Config{Addresses: []string{esHost}})
	if err != nil {
		return nil, err
	}

	return fn(cl)
}

// QueryAndParse performs an Elasticsearch API call and parse the response
func (c *ESClient) QueryAndParse(o interface{}, fn ESFunction) error {
	if c == nil {
		return fmt.Errorf("this ESClient is nil")
	}

	res, err := c.Query(fn)
	if err != nil {
		return err
	}

	if res.IsError() {
		return fmt.Errorf("error response from Elasticsearch: %s", res.String())
	}

	return json.NewDecoder(res.Body).Decode(o)
}

// ESSearchResult represents the result of an Elasticsearch search operation
type ESSearchResult struct {
	Took         uint64          `json:"took"`
	TimedOut     bool            `json:"timed_out"`
	Shards       ESShard         `json:"_shards"`
	Hits         ESResultHits    `json:"hits"`
	Aggregations json.RawMessage `json:"aggregations"`
}

// Size returns the number of hits in the search results
func (r *ESSearchResult) Size() int {
	if r == nil {
		return 0
	}
	return len(r.Hits.Hits)
}

// ESShard is an Elasticsearch shard
type ESShard struct {
	Total      int `json:"total"`
	Successful int `json:"successful"`
	Failed     int `json:"failed"`
}

// ESResultHits represents the hits in the result of an Elasticsearch search
type ESResultHits struct {
	Total    ESSearchTotal `json:"total"`
	MaxScore float32       `json:"max_score"`
	Hits     []ESSearchHit `json:"hits"`
}

type ESSearchTotal struct {
	Value    int    `json:"value"`
	Relation string `json:"relation"`
}

// ESSearchHit is an individual Elasticsearch search hit
type ESSearchHit struct {
	Index     string              `json:"_index"`
	Type      string              `json:"_type"`
	ID        string              `json:"_id"`
	Score     float32             `json:"_score"`
	Source    json.RawMessage     `json:"_source"`
	Highlight map[string][]string `json:"highlight,omitempty"`
}

// GetSource returns the raw json of the Hit Source as a map
func (h *ESSearchHit) GetSource() (*map[string]interface{}, error) {
	m := make(map[string]interface{})

	if h != nil {
		err := json.Unmarshal(h.Source, &m)
		if err != nil {
			return &m, err
		}
	}

	return &m, nil
}
