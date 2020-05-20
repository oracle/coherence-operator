/*
 * Copyright (c) 2019, 2020, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package local

import (
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	coh "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	"github.com/oracle/coherence-operator/pkg/management"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	"io/ioutil"
	"net/http"

	corev1 "k8s.io/api/core/v1"
	"testing"

	. "github.com/onsi/gomega"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	loggingConfigMapName = "logging-config-map"
)

var defaultReplicas = int(coh.DefaultReplicas)

func TestDeploymentWithCoherenceLogLevel(t *testing.T) {
	g := NewGomegaWithT(t)

	ctx := helper.CreateTestContext(t)
	defer helper.DumpOperatorLogsAndCleanup(t, ctx)

	namespace, err := ctx.GetWatchNamespace()
	g.Expect(err).NotTo(HaveOccurred())

	AssertDeploymentsWithContext(t, ctx, "deployment-with-log-level.yaml")
	assertLogConfigOnMembers(g, namespace, "test", defaultReplicas, 9)
}

func TestDeploymentWithLoggingConfigFile(t *testing.T) {
	g := NewGomegaWithT(t)

	ctx := helper.CreateTestContext(t)
	defer helper.DumpOperatorLogsAndCleanup(t, ctx)

	namespace, err := ctx.GetWatchNamespace()
	g.Expect(err).NotTo(HaveOccurred())

	AssertDeploymentsWithContext(t, ctx, "deployment-with-log-configfile.yaml")
	assertLogConfigOnMembers(g, namespace, "test", defaultReplicas, 5)
}

// ----- helpers ------------------------------------------------------------

// Delete logging-config-map from the given namespace
func deleteLoggingConfigMap(namespace string) {
	_ = framework.Global.KubeClient.CoreV1().ConfigMaps(namespace).Delete(loggingConfigMapName, &metav1.DeleteOptions{})
}

// Create required logging-config-map in a test namespace.
func createLoggingConfigMap(g *GomegaWithT, namespace string) {
	f := framework.Global
	logConfigMap := &corev1.ConfigMap{}
	logConfigMap.SetNamespace(namespace)
	logConfigMap.SetName(loggingConfigMapName)

	logConfigMap.Data = make(map[string]string)
	rootDir, err := helper.FindProjectRootDir()
	g.Expect(err).NotTo(HaveOccurred())

	var cmData []byte
	cmData, err = ioutil.ReadFile(rootDir + "/test/e2e/local/logging.properties")
	g.Expect(err).NotTo(HaveOccurred())
	logConfigMap.Data["logging.properties"] = string(cmData)

	existingCM, _ := f.KubeClient.CoreV1().ConfigMaps(namespace).Get(loggingConfigMapName,
		metav1.GetOptions{})
	if existingCM != nil {
		deleteLoggingConfigMap(namespace)
	}

	_, err = f.KubeClient.CoreV1().ConfigMaps(namespace).Create(logConfigMap)
	g.Expect(err).NotTo(HaveOccurred())
}

// Verify log configuration on each member of a cluster which is created using specified yaml file.
func assertLogConfigOnMembers(g *GomegaWithT, namespace, deployment string, replicas, logLevel int) {
	ok, _ := helper.IsCoherenceVersionAtLeast(12, 2, 1, 4)
	if ok {
		f := framework.Global

		// Get the list of Pods
		pods, err := helper.ListCoherencePodsForDeployment(f.KubeClient, namespace, deployment)
		g.Expect(err).NotTo(HaveOccurred())
		g.Expect(len(pods)).NotTo(BeZero())

		// Port forward to the first Pod
		pf, ports, err := helper.StartPortForwarderForPod(&pods[0])
		g.Expect(err).NotTo(HaveOccurred())
		defer pf.Close()

		cl := &http.Client{}
		members, _, err := management.GetMembers(cl, "127.0.0.1", ports[coh.PortNameManagement])
		g.Expect(err).NotTo(HaveOccurred())
		g.Expect(len(members.Items)).To(Equal(replicas))

		// assert that each member has required logging configuration.
		for _, member := range members.Items {
			g.Expect(member.MachineName).NotTo(BeEmpty())
			g.Expect(member.LoggingLevel).To(Equal(logLevel))
		}
	}
}
