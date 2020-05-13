/*
 * Copyright (c) 2019, 2020, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package helm_test

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

func CleanupHelm(t *testing.T, hm *helper.HelmReleaseManager, helmHelper *helper.HelmHelper) {
	helper.DumpState(helmHelper.Namespace, t.Name(), t)

	// ensure that the chart is uninstalled
	_, err := hm.UninstallRelease()
	if err != nil {
		fmt.Println("Failed to uninstall helm release " + err.Error())
	}

	// Wait for the Operator Pods to die as terminating Pods can mess up the next test that runs
	// if it is too quick after this test.
	err = helper.WaitForOperatorCleanup(helmHelper.KubeClient, helmHelper.Namespace, t)
	if err != nil {
		fmt.Println("Failed waiting for Operator clean-up " + err.Error())
	}
}

func DeployCoherence(t *testing.T, ctx *framework.Context, namespace, yamlFile string) ([]v1.CoherenceDeployment, error) {
	g := NewGomegaWithT(t)
	f := framework.Global

	deployments, err := helper.NewCoherenceDeploymentFromYaml(namespace, yamlFile)
	g.Expect(err).NotTo(HaveOccurred())

	for _, d := range deployments {
		err = f.Client.Create(context.TODO(), &d, helper.DefaultCleanup(ctx))
		if err != nil {
			return deployments, err
		}
	}

	// Wait for the StatefulSet(s)
	for _, d := range deployments {
		_, err = helper.WaitForStatefulSetForDeployment(f.KubeClient, namespace, &d, time.Second*5, time.Minute*5, t)
		if err != nil {
			return deployments, err
		}
	}

	return deployments, nil
}
