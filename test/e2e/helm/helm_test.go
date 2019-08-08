package helm_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/oracle/coherence-operator/test/e2e/helper/matchers"

	"github.com/oracle/coherence-operator/test/e2e/helper"

	coreV1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"time"
)

// This file is an example of how to write a test for installing the Operator Helm chart.
var _ = Describe("Operator Helm Chart", func() {
	var hm *helper.HelmReleaseManager
	var err error

	// Each "When" section is a single test...
	When("installing Helm chart with empty values", func() {

		// The JustBefore function is where the Helm install happens
		JustBeforeEach(func() {
			// Create the values to use
			values := make(map[string]interface{})

			// Create a HelmReleaseManager with a release name and values
			hm, err = HelmHelper.NewHelmReleaseManager("foo", &values)
			Expect(err).ToNot(HaveOccurred())

			// Install the chart
			_, err = hm.InstallRelease()
			Expect(err).ToNot(HaveOccurred())
		})

		// The JustAfter function will ensure the chart is uninstalled
		JustAfterEach(func() {
			// ensure that the chart is uninstalled
			_, err := hm.UninstallRelease()
			Expect(err).ToNot(HaveOccurred())
		})

		// There should be ONLY ONE It section in this When section to do the assertions.
		// If there were multiple IT sections the chart would be installed and uninstalled before and after each It
		// because that is how Ginkgo works.
		It("should deploy the Operator", func() {
			ns := HelmHelper.Namespace
			client := HelmHelper.KubeClient

			pods, err := helper.WaitForOperatorPods(HelmHelper.KubeClient, ns, time.Second*10, time.Minute*5)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(pods)).To(Equal(1))

			pod, err := client.CoreV1().Pods(ns).Get(pods[0].Name, metav1.GetOptions{})
			Expect(err).ToNot(HaveOccurred())

			err = helper.WaitForPodReady(client, pod.Namespace, pod.Name, time.Second*10, time.Minute*5)
			Expect(err).ToNot(HaveOccurred())

			container := pod.Spec.Containers[0]
			Expect(container.Name).To(Equal("coherence-operator"))
			Expect(container.Env).To(HaveEnvVar(coreV1.EnvVar{Name: "OPERATOR_NAME", Value: "coherence-operator"}))
		})
	})

})
