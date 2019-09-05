/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package helm_test

import (
	"fmt"
	. "github.com/onsi/gomega"
	. "github.com/oracle/coherence-operator/test/e2e/helper/matchers"
	"net/http"
	"testing"

	"github.com/oracle/coherence-operator/test/e2e/helper"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"time"
)

// Test installing the Operator with the default values.
func TestBasicHelmInstall(t *testing.T) {
	g := NewGomegaWithT(t)

	// Create a helper.HelmHelper
	helmHelper, err := helper.NewOperatorChartHelper()
	if err != nil {
		t.Fatal(err)
	}

	// Create the values to use (in this case empty)
	values := helper.OperatorValues{}

	// If we wanted to load a YAML file into the values then we can use
	// the values.LoadFromYaml method with a file name relative to this
	// test files location
	//err = values.LoadFromYaml("test.yaml")
	//g.Expect(err).ToNot(HaveOccurred())

	// Create a HelmReleaseManager with a release name and values
	hm, err := helmHelper.NewOperatorHelmReleaseManager("operator", &values)
	g.Expect(err).ToNot(HaveOccurred())

	// Defer cleanup (helm delete) to make sure it happens when this method exits
	defer CleanupHelm(t, hm, helmHelper)

	// Install the chart
	_, err = hm.InstallRelease()
	g.Expect(err).ToNot(HaveOccurred())

	// The chart is installed but the Pod(s) may not exist yet so wait for it...
	// (we wait a maximum of 5 minutes, retrying every 10 seconds)
	pods, err := helper.WaitForOperatorPods(helmHelper.KubeClient, helmHelper.Namespace, time.Second*10, time.Minute*5)
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(len(pods)).To(Equal(1))

	// The Pod(s) exist so get one of them using the k8s client from the helper
	// which is in the HelmHelper.KubeClient var configured in the suite .go file.
	pod, err := helmHelper.KubeClient.CoreV1().Pods(helmHelper.Namespace).Get(pods[0].Name, metav1.GetOptions{})
	g.Expect(err).ToNot(HaveOccurred())

	// The chart is installed but the Pod we have may not be ready yet so wait for it...
	// (we wait a maximum of 5 minutes, retrying every 10 seconds)
	err = helper.WaitForPodReady(helmHelper.KubeClient, pod.Namespace, pod.Name, time.Second*10, time.Minute*5)
	g.Expect(err).ToNot(HaveOccurred())

	// Assert some things
	container := pod.Spec.Containers[0]
	g.Expect(container.Name).To(Equal("coherence-operator"))
	g.Expect(container.Env).To(HaveEnvVar(corev1.EnvVar{Name: "OPERATOR_NAME", Value: "coherence-operator"}))

	// Obtain a PortForwarder for the Pod - this will forward all of the ports defined in the Pod spec.
	// The values returned are a PortForwarder, a map of port name to local port and any error
	fwd, ports, err := helper.StartPortForwarderForPod(pod)
	g.Expect(err).ToNot(HaveOccurred())

	// Defer closing the PortForwarder so we clean-up properly
	defer fwd.Close()

	// The ReST port in the Operator container spec is named "rest"
	restPort := ports["rest"]

	// Do a GET on the Zone endpoint for the Pod's NodeName, we should get no error and a 200 response
	resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d/site/%s", restPort, pod.Spec.NodeName))
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(resp.StatusCode).To(Equal(200))
}
