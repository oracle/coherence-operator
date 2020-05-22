/*
 * Copyright (c) 2020, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package compatibility_test

import (
	"context"
	"fmt"
	. "github.com/onsi/gomega"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	v1 "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	framework.MainEntry(m)
}

func CleanupHelm(t *testing.T, hm *helper.HelmReleaseManager, hh *helper.HelmHelper) {
	// ensure that the chart is uninstalled
	_, err := hm.UninstallRelease()
	if err != nil {
		fmt.Println("Failed to uninstall helm release " + err.Error())
	}

	// Wait for the Operator Pods to die as terminating Pods can mess up the next test that runs
	// if it is too quick after this test.
	err = helper.WaitForOperatorCleanup(hh.KubeClient, hh.Namespace, t)
	if err != nil {
		fmt.Println("Failed waiting for Operator clean-up " + err.Error())
	}
}

func DeployCoherenceCluster(t *testing.T, ctx *framework.Context, namespace, yamlFile string) (v1.CoherenceCluster, error) {
	g := NewGomegaWithT(t)
	f := framework.Global

	cluster, err := helper.NewCoherenceFromYaml(namespace, yamlFile)
	g.Expect(err).NotTo(HaveOccurred())

	// deploy the CoherenceCluster
	err = f.Client.Create(context.TODO(), &cluster, helper.DefaultCleanup(ctx))
	if err != nil {
		return cluster, err
	}

	// Wait for the StatefulSet(s)
	roles := cluster.GetRoles()
	for _, role := range roles {
		_, err = helper.WaitForStatefulSetForDeployment(f.KubeClient, namespace, &cluster, role, time.Second*5, time.Minute*5, t)
		if err != nil {
			return cluster, err
		}
	}

	return cluster, nil
}

type Cleanup struct {
	t   *testing.T
	ctx *framework.Context
	rm  *helper.HelmReleaseManager
	hh  *helper.HelmHelper
}

func (in Cleanup) Run() {
	if in.t != nil {
		ns := helper.GetTestNamespace()
		helper.DumpState(ns, in.t.Name(), in.t)
		helper.DumpOperatorLogsAndCleanup(in.t, in.ctx)

		if in.rm != nil && in.hh != nil {
			CleanupHelm(in.t, in.rm, in.hh)
		}
	}
}
