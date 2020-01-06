/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package remote

import (
	"context"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	cohv1 "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	"k8s.io/apimachinery/pkg/types"
	"testing"

	f "github.com/operator-framework/operator-sdk/pkg/test"
)

// Operator SDK test suite entry point
func TestMain(m *testing.M) {
	f.MainEntry(m)
}

func cleanup(t *testing.T, namespace, clusterName string, ctx *framework.TestCtx) {
	helper.DumpOperatorLogs(t, ctx)
	deleteCluster(namespace, clusterName)
	ctx.Cleanup()
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
