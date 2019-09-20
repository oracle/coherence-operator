/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	es6 "github.com/elastic/go-elasticsearch/v6"
	"github.com/elastic/go-elasticsearch/v6/esapi"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"strings"
	"time"
)

// The Coherence Cluster Index Pattern used by Kibana
const KibanaIndexPatternCoherenceCluster = "6abb1220-3feb-11e9-a9a3-4b1c09db6e6a"

// An elastic search client associated to an Elasticsearch Pod
type ElasticSearchClient struct {
	Pod *corev1.Pod
}

// A function that takes an Elasticsearch client and executes an API call
type ESFunction func(es *es6.Client) (*esapi.Response, error)

func NewESSearchFunction(index, query string) ESFunction {
	return func(es *es6.Client) (*esapi.Response, error) {
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
func (c *ElasticSearchClient) WaitForCoherenceIndices(retryInterval, timeout time.Duration, logger Logger) error {

	// A query function to retrieve ES indices
	fn := func(es *es6.Client) (*esapi.Response, error) {
		return es.Cat.Indices(es.Cat.Indices.WithPretty())
	}

	var res *esapi.Response

	err := wait.Poll(retryInterval, timeout, func() (done bool, err error) {
		res, err = c.Query(fn)
		if err != nil {
			logger.Logf("Waiting for Coherence indices in Elasticsearch - %s\n", err.Error())
			return false, err
		}
		ok := strings.Contains(strings.ToLower(res.String()), "coherence")
		return ok, nil
	})

	if err != nil {
		var s string
		if res == nil {
			s = "response is nil"
		} else {
			s = res.String()
		}
		logger.Logf("Error waiting for Coherence indices in Elasticsearch - error: '%s' response:\n%s\n", err.Error(), s)
	}

	return err
}

// Perform an Elasticsearch API call
func (c *ElasticSearchClient) Query(fn ESFunction) (*esapi.Response, error) {
	if c == nil {
		return nil, fmt.Errorf("this ElasticSearchClient is nil")
	}

	fwd, ports, err := StartPortForwarderForPod(c.Pod)
	if err != nil {
		return nil, err
	}

	// Defer closing the PortForwarder so we clean-up properly
	defer fwd.Close()

	// The ReST port in the Elasticsearch container spec is named "rest"
	port := ports["rest"]

	esHost := fmt.Sprintf("http://127.0.0.1:%d", port)
	cl, err := es6.NewClient(es6.Config{Addresses: []string{esHost}})
	if err != nil {
		return nil, err
	}

	return fn(cl)
}

// Perform an Elasticsearch API call and parse the response
func (c *ElasticSearchClient) QueryAndParse(o interface{}, fn ESFunction) error {
	if c == nil {
		return fmt.Errorf("this ElasticSearchClient is nil")
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

// SearchResult represents the result of an Elasticsearch search operation
type ESSearchResult struct {
	Took         uint64          `json:"took"`
	TimedOut     bool            `json:"timed_out"`
	Shards       ESShard         `json:"_shards"`
	Hits         ESResultHits    `json:"hits"`
	Aggregations json.RawMessage `json:"aggregations"`
}

// Return the number of hits in the search results
func (r *ESSearchResult) Size() int {
	if r == nil {
		return 0
	}
	return len(r.Hits.Hits)
}

// An Elasticsearch shard
type ESShard struct {
	Total      int `json:"total"`
	Successful int `json:"successful"`
	Failed     int `json:"failed"`
}

// ESResultHits represents the hits in the result of an Elasticsearch search
type ESResultHits struct {
	Total    int           `json:"total"`
	MaxScore float32       `json:"max_score"`
	Hits     []ESSearchHit `json:"hits"`
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

// Obtain the raw json of the Hit Source as a map
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
