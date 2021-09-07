/*
 * Copyright (c) 2021, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package clients

import "testing"

const (
	defaultCacheConfig = "test-cache-config.xml"
)

// ClientType represents a type of text client
type ClientType string

var (
	// ClientTypeExtend is a ClientType representing Extend client tests
	ClientTypeExtend ClientType = "extend"
	// ClientTypeGrpc is a ClientType representing gRPC client tests
	ClientTypeGrpc ClientType = "grpc"
)

// ClientTestCase represents the data required to run a client test case.
type ClientTestCase struct {
	// ClientType is the type of test client (e.g. extend or grpc).
	ClientType ClientType
	// Name is the name of the test.
	Name string
	// Cluster is the test Coherence cluster.
	Cluster *CoherenceCluster
	// CacheConfig is the name of the cache configuration file the client can connect to.
	CacheConfig string
	// Test is the test method to run
	Test TestExecution
}

type TestExecution func(*testing.T, ClientTestCase)

// Execute runs all the child test cases.
func Execute(parentTest *testing.T, testCases []ClientTestCase) {
	for _, tc := range testCases {
		parentTest.Run(tc.Name, func(t *testing.T) {
			tc.Test(t, tc)
		})
	}
}
