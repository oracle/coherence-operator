/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package remote

import (
	"context"
	cohv1 "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	"k8s.io/apimachinery/pkg/types"
	"testing"

	f "github.com/operator-framework/operator-sdk/pkg/test"
)

// Operator SDK test suite entry point
func TestMain(m *testing.M) {
	f.MainEntry(m)
}

// deleteCluster deletes a cluster.
func deleteCluster(namespace, name string) {
	cluster := cohv1.CoherenceCluster{}
	f := f.Global

	err := f.Client.Get(context.TODO(), types.NamespacedName{Namespace: namespace, Name: name}, &cluster)

	if err == nil {
		_ = f.Client.Delete(context.TODO(), &cluster)
	}
}
