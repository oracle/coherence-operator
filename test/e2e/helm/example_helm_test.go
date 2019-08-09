package helm_test

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/oracle/coherence-operator/test/e2e/helper/matchers"
	"net/http"

	"github.com/oracle/coherence-operator/test/e2e/helper"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"time"
)

// This file is an example of how to write a test for installing the Operator Helm chart
// using the Ginkgo test framework.
var _ = Describe("Operator Helm Chart", func() {
	var hm *helper.HelmReleaseManager
	var err error

	// Normally in Ginkgo each "When" section is a set of related tests but for testing
	// the Helm chart install we make it just a single test.
	// To write multiple related tests in a single Go file put each test into its own When section
	// inside a single Describe section.
	When("installing Helm chart with empty values", func() {

		// The JustBefore function is where the Helm install happens
		JustBeforeEach(func() {
			// Create the values to use
			values := helper.OperatorValues{}

			// If we wanted to load a YAML file into the values then we can use
			// the values.LoadFromYaml method with a file name relative to this
			// test files location
			//err = values.LoadFromYaml("test.yaml")

			// Create a HelmReleaseManager with a release name and values
			hm, err = HelmHelper.NewOperatorHelmReleaseManager("foo", &values)
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
		// because Ginkgo runs the JustBeforeEach and BeforeEach before every IT and runs JustAfterEach and AfterEach
		// after every It.
		It("should deploy the Operator", func() {
			// The chart is installed but the Pod(s) may not exist yet so wait for it...
			// (we wait a maximum of 5 minutes, retrying every 10 seconds)
			pods, err := helper.WaitForOperatorPods(HelmHelper.KubeClient, HelmHelper.Namespace, time.Second*10, time.Minute*5)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(pods)).To(Equal(1))

			// The Pod(s) exist so get one of them using the k8s client from the helper
			// which is in the HelmHelper.KubeClient var configured in the suite .go file.
			pod, err := HelmHelper.KubeClient.CoreV1().Pods(HelmHelper.Namespace).Get(pods[0].Name, metav1.GetOptions{})
			Expect(err).ToNot(HaveOccurred())

			// The chart is installed but the Pod we have may not be ready yet so wait for it...
			// (we wait a maximum of 5 minutes, retrying every 10 seconds)
			err = helper.WaitForPodReady(HelmHelper.KubeClient, pod.Namespace, pod.Name, time.Second*10, time.Minute*5)
			Expect(err).ToNot(HaveOccurred())

			// Assert some things
			container := pod.Spec.Containers[0]
			Expect(container.Name).To(Equal("coherence-operator"))
			Expect(container.Env).To(HaveEnvVar(corev1.EnvVar{Name: "OPERATOR_NAME", Value: "coherence-operator"}))

			// Obtain a PortForwarder for the Pod - this will forward all of the ports defined in the Pod spec.
			// The values returned are a PortForwarder, a map of port name to local port and any error
			fwd, ports, err := helper.StartPortForwarderForPod(pod)
			Expect(err).ToNot(HaveOccurred())

			// Defer closing the PortForwarder so we clean-up properly
			defer fwd.Close()

			// The ReST port in the Operator container spec is named "rest"
			restPort := ports["rest"]

			// Do a GET on the Zone endpoint for the Pod's NodeName, we should get no error and a 200 response
			resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d/zone/%s", restPort, pod.Spec.NodeName))
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(200))
		})
	})

})
