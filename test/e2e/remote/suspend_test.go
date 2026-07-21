/*
 * Copyright (c) 2020, 2026, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package remote

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	. "github.com/onsi/gomega"
	cohv1 "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/pkg/utils"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	testFinalizer           = "coherence.oracle.com/test"
	runIPv6FinalizerTestEnv = "RUN_IPV6_FINALIZER_TEST"
)

func TestSuspendServices(t *testing.T) {
	g := NewGomegaWithT(t)

	ns := helper.GetTestNamespace()
	c, err := helper.NewSingleCoherenceFromYamlWithSuffix(ns, "suspend-test.yaml", "-suspend")
	g.Expect(err).NotTo(HaveOccurred())
	runSuspendFinalizerTest(t, c, false, nil)
}

func TestSuspendServicesWithIPv6PodIP(t *testing.T) {
	if os.Getenv(runIPv6FinalizerTestEnv) != "true" {
		t.Skipf("set %s=true to run the IPv6-primary dual-stack test", runIPv6FinalizerTestEnv)
	}

	g := NewGomegaWithT(t)

	ns := helper.GetTestNamespace()
	c, err := helper.NewSingleCoherenceFromYamlWithSuffix(ns, "suspend-test.yaml", "-suspend-ipv6")
	g.Expect(err).NotTo(HaveOccurred())

	c.Spec.Replicas = ptr.To(int32(1))
	c.Spec.HeadlessServiceIpFamilies = []corev1.IPFamily{corev1.IPv6Protocol}
	c.Spec.SuspendProbe = nil
	if c.Spec.Coherence == nil {
		c.Spec.Coherence = &cohv1.CoherenceSpec{}
	}
	if c.Spec.Coherence.WKA == nil {
		c.Spec.Coherence.WKA = &cohv1.CoherenceWKASpec{}
	}
	c.Spec.Coherence.WKA.IPFamily = ptr.To(corev1.IPv6Protocol)

	runSuspendFinalizerTest(t, c, true, validateIPv6PrimaryPod)
}

func runSuspendFinalizerTest(t *testing.T, c cohv1.Coherence, requirePersistencePVC bool, podValidator func(*testing.T, []corev1.Pod)) {
	t.Helper()
	// Ensure that everything is cleaned up after the test!
	testContext.CleanupAfterTest(t)
	ctx := context.Background()
	g := NewGomegaWithT(t)

	if requirePersistencePVC {
		requireDefaultStorageClass(t)
	}

	err := testContext.Client.Create(ctx, &c)
	g.Expect(err).NotTo(HaveOccurred())

	if requirePersistencePVC {
		pvc, pvcErr := waitForPersistencePVC(&c)
		g.Expect(pvcErr).NotTo(HaveOccurred())
		t.Logf("PersistentVolumeClaim %s/%s is Bound", pvc.Namespace, pvc.Name)
	}

	assertDeploymentEventuallyInDesiredState(t, c, c.GetReplicas())

	// get the StatefulSet for the deployment
	sts, err := testContext.KubeClient.AppsV1().StatefulSets(c.Namespace).Get(ctx, c.Name, metav1.GetOptions{})
	g.Expect(err).NotTo(HaveOccurred())
	pods, err := helper.ListCoherencePodsForDeployment(testContext, c.Namespace, c.Name)
	g.Expect(err).NotTo(HaveOccurred())
	if podValidator != nil {
		podValidator(t, pods)
	}

	err = addTestFinalizer(&c)
	g.Expect(err).NotTo(HaveOccurred())
	defer removeAllFinalizersLoggingErrors(t, &c)

	// Delete the deployment which should cause services to be suspended
	// The deployment will not be deleted yet as we still have the test finalizer in place
	err = testContext.Client.Delete(ctx, &c)
	g.Expect(err).NotTo(HaveOccurred())
	// The Operator should run its finalizer and suspend services
	err = waitForFinalizerTasks(c.GetNamespacedName())
	g.Expect(err).NotTo(HaveOccurred())

	// assert that the cache service is suspended
	svc, err := ManagementOverRestRequest(&c, "/management/coherence/cluster/services/PartitionedCache")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(svc["quorumStatus"]).To(ContainElement("Suspended"))

	err = waitForServiceSuspendedEvent(c.UID, c.Namespace, sts.Name)
	g.Expect(err).NotTo(HaveOccurred())

	// remove the test finalizer which should then let everything be deleted
	err = removeAllFinalizers(&c)
	g.Expect(err).NotTo(HaveOccurred())
	// the StatefulSet should eventually be deleted
	err = helper.WaitForDelete(testContext, sts)
	g.Expect(err).NotTo(HaveOccurred())
}

func validateIPv6PrimaryPod(t *testing.T, pods []corev1.Pod) {
	t.Helper()
	if len(pods) != 1 {
		t.Fatalf("expected exactly one workload Pod, found %d: %v", len(pods), podAddressSummary(pods))
	}

	pod := pods[0]
	addresses := podIPStrings(pod)
	primary := net.ParseIP(pod.Status.PodIP)
	if primary == nil || primary.To4() != nil {
		t.Fatalf("expected Pod %s primary IP to be valid IPv6, primary=%q podIPs=%v", pod.Name, pod.Status.PodIP, addresses)
	}

	hasIPv4 := false
	hasIPv6 := false
	for _, podIP := range pod.Status.PodIPs {
		parsed := net.ParseIP(podIP.IP)
		if parsed == nil {
			t.Fatalf("Pod %s has invalid address %q, primary=%q podIPs=%v", pod.Name, podIP.IP, pod.Status.PodIP, addresses)
		}
		if parsed.To4() == nil {
			hasIPv6 = true
		} else {
			hasIPv4 = true
		}
	}
	if !hasIPv4 || !hasIPv6 {
		t.Fatalf("expected Pod %s to have IPv4 and IPv6 addresses, primary=%q podIPs=%v", pod.Name, pod.Status.PodIP, addresses)
	}
}

func podAddressSummary(pods []corev1.Pod) []string {
	summary := make([]string, 0, len(pods))
	for _, pod := range pods {
		summary = append(summary, fmt.Sprintf("%s primary=%q podIPs=%v", pod.Name, pod.Status.PodIP, podIPStrings(pod)))
	}
	return summary
}

func podIPStrings(pod corev1.Pod) []string {
	addresses := make([]string, 0, len(pod.Status.PodIPs))
	for _, podIP := range pod.Status.PodIPs {
		addresses = append(addresses, podIP.IP)
	}
	return addresses
}

func requireDefaultStorageClass(t *testing.T) {
	t.Helper()
	g := NewGomegaWithT(t)

	classes, err := testContext.KubeClient.StorageV1().StorageClasses().List(testContext.Context, metav1.ListOptions{})
	g.Expect(err).NotTo(HaveOccurred())
	var defaults []string
	for _, class := range classes.Items {
		annotations := class.GetAnnotations()
		if strings.EqualFold(annotations["storageclass.kubernetes.io/is-default-class"], "true") ||
			strings.EqualFold(annotations["storageclass.beta.kubernetes.io/is-default-class"], "true") {
			defaults = append(defaults, class.Name)
		}
	}
	g.Expect(defaults).NotTo(BeEmpty(), "expected a default StorageClass, found %s", mustJSON(classes.Items))
}

func waitForPersistencePVC(c *cohv1.Coherence) (*corev1.PersistentVolumeClaim, error) {
	selector := fmt.Sprintf("%s=%s,%s=%s",
		cohv1.LabelComponent, cohv1.LabelComponentPVC,
		cohv1.LabelCoherenceDeployment, c.Name)
	var claims []corev1.PersistentVolumeClaim

	err := wait.PollUntilContextTimeout(testContext.Context, time.Second, 3*time.Minute, true, func(ctx context.Context) (bool, error) {
		list, err := testContext.KubeClient.CoreV1().PersistentVolumeClaims(c.Namespace).List(ctx, metav1.ListOptions{LabelSelector: selector})
		if err != nil {
			return false, err
		}
		claims = list.Items
		if len(claims) > 1 {
			return false, fmt.Errorf("expected one persistence PVC matching %q, found %s", selector, mustJSON(claims))
		}
		if len(claims) == 0 {
			testContext.Logf("Waiting for persistence PVC matching %q", selector)
			return false, nil
		}
		if claims[0].Status.Phase != corev1.ClaimBound {
			testContext.Logf("Waiting for persistence PVC %s/%s to be Bound - phase %s",
				claims[0].Namespace, claims[0].Name, claims[0].Status.Phase)
			return false, nil
		}
		return true, nil
	})
	if err != nil {
		classes, classesErr := testContext.KubeClient.StorageV1().StorageClasses().List(testContext.Context, metav1.ListOptions{})
		events, eventsErr := testContext.KubeClient.CoreV1().Events(c.Namespace).List(testContext.Context, metav1.ListOptions{})
		var classItems interface{} = "<unavailable>"
		if classes != nil {
			classItems = classes.Items
		}
		var eventItems interface{} = "<unavailable>"
		if events != nil {
			eventItems = events.Items
		}
		return nil, fmt.Errorf("waiting for persistence PVC matching %q: %w; PVCs=%s; StorageClasses=%s (error=%v); Events=%s (error=%v)",
			selector, err, mustJSON(claims), mustJSON(classItems), classesErr, mustJSON(eventItems), eventsErr)
	}
	return claims[0].DeepCopy(), nil
}

func waitForServiceSuspendedEvent(uid types.UID, namespace, statefulSetName string) error {
	var matchingEvents []string
	err := wait.PollUntilContextTimeout(testContext.Context, time.Second, 5*time.Minute, true, func(ctx context.Context) (bool, error) {
		list, err := testContext.KubeClient.EventsV1().Events(namespace).List(ctx, metav1.ListOptions{})
		if err != nil {
			return false, err
		}

		matchingEvents = matchingEvents[:0]
		foundSuccess := false
		foundFailure := false
		for _, event := range list.Items {
			if event.Regarding.UID != uid {
				continue
			}
			matchingEvents = append(matchingEvents, fmt.Sprintf("type=%s reason=%s note=%q", event.Type, event.Reason, event.Note))
			if event.Reason == "ServiceSuspendFailed" {
				foundFailure = true
			}
			if event.Type == corev1.EventTypeNormal && event.Reason == "ServiceSuspended" && strings.Contains(event.Note, statefulSetName) {
				foundSuccess = true
			}
		}
		if foundFailure {
			return false, fmt.Errorf("found ServiceSuspendFailed event for Coherence UID %s: %v", uid, matchingEvents)
		}
		return foundSuccess, nil
	})
	if err != nil {
		return fmt.Errorf("waiting for normal ServiceSuspended event naming StatefulSet %s for Coherence UID %s: %w; matching events: %v",
			statefulSetName, uid, err, matchingEvents)
	}
	return nil
}

func mustJSON(value interface{}) string {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Sprintf("<json error: %v>", err)
	}
	return string(data)
}

func TestNotSuspendServicesWhenSuspendDisabled(t *testing.T) {
	// Ensure that everything is cleaned up after the test!
	testContext.CleanupAfterTest(t)
	ctx := context.Background()
	g := NewGomegaWithT(t)

	ns := helper.GetTestNamespace()
	c, err := helper.NewSingleCoherenceFromYamlWithSuffix(ns, "suspend-test.yaml", "-suspend-disable")
	g.Expect(err).NotTo(HaveOccurred())

	// Set the flag to NOT suspend on shutdown
	c.Spec.SuspendServicesOnShutdown = ptr.To(false)

	installSimpleDeployment(t, c)

	// get the StatefulSet for the deployment
	sts, err := testContext.KubeClient.AppsV1().StatefulSets(ns).Get(ctx, c.Name, metav1.GetOptions{})
	g.Expect(err).NotTo(HaveOccurred())

	err = addTestFinalizer(&c)
	g.Expect(err).NotTo(HaveOccurred())

	// Delete the deployment which should cause services to be suspended
	// The deployment will not be deleted yet as we still have the test finalizer in place
	err = testContext.Client.Delete(ctx, &c)
	g.Expect(err).NotTo(HaveOccurred())
	// The Operator should run its finalizer and suspend services
	err = waitForFinalizerTasks(c.GetNamespacedName())
	g.Expect(err).NotTo(HaveOccurred())

	// assert that the cache service is suspended
	svc, err := ManagementOverRestRequest(&c, "/management/coherence/cluster/services/PartitionedCache")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(svc["quorumStatus"]).NotTo(ContainElement("Suspended"))

	// remove the test finalizer which should then let everything be deleted
	err = removeAllFinalizers(&c)
	g.Expect(err).NotTo(HaveOccurred())
	// the StatefulSet should eventually be deleted
	err = helper.WaitForDelete(testContext, sts)
	g.Expect(err).NotTo(HaveOccurred())
}

func TestSuspendServicesOnScaleDownToZero(t *testing.T) {
	// Ensure that everything is cleaned up after the test!
	testContext.CleanupAfterTest(t)
	ctx := context.Background()
	g := NewGomegaWithT(t)

	ns := helper.GetTestNamespace()
	c, err := helper.NewSingleCoherenceFromYamlWithSuffix(ns, "suspend-test.yaml", "-scale-zero")
	g.Expect(err).NotTo(HaveOccurred())

	installSimpleDeployment(t, c)

	err = addTestFinalizer(&c)
	g.Expect(err).NotTo(HaveOccurred())

	// Add a finalizer to the StatefulSet to stop it being deleted
	sts := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      c.Name,
			Namespace: c.Namespace,
		},
	}
	err = addTestFinalizer(sts)
	g.Expect(err).NotTo(HaveOccurred())
	// ensure we remove the finalizer
	defer removeAllFinalizersLoggingErrors(t, sts)

	// re-fetch the latest Coherence state and scale down to zero, which should cause services to be suspended
	err = testContext.Client.Get(ctx, c.GetNamespacedName(), &c)
	g.Expect(err).NotTo(HaveOccurred())
	patch := client.RawPatch(types.MergePatchType, []byte(`{"spec":{"replicas":0}}`))
	err = testContext.Client.Patch(testContext.Context, &c, patch)
	g.Expect(err).NotTo(HaveOccurred())

	// The Operator should suspend services and delete the StatefulSet causing its deletion timestamp to be set
	// As we added a finalizer to the StatefulSet it will not actually get deleted yet
	err = waitForStatefulSetDeletionTimestamp(c.GetNamespacedName())
	g.Expect(err).NotTo(HaveOccurred())

	// assert that the cache service is suspended
	svc, err := ManagementOverRestRequest(&c, "/management/coherence/cluster/services/PartitionedCache")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(svc["quorumStatus"]).To(ContainElement("Suspended"))

	// remove the test finalizer from the StatefulSet and Coherence deployment which should then let everything be deleted
	err = removeAllFinalizers(sts)
	g.Expect(err).NotTo(HaveOccurred())
	err = removeAllFinalizers(&c)
	g.Expect(err).NotTo(HaveOccurred())
	// the StatefulSet should eventually be deleted
	err = helper.WaitForDelete(testContext, sts)
	g.Expect(err).NotTo(HaveOccurred())
}

func TestNotSuspendServicesOnScaleDownToZeroIfSuspendDisabled(t *testing.T) {
	// Ensure that everything is cleaned up after the test!
	testContext.CleanupAfterTest(t)
	ctx := context.Background()
	g := NewGomegaWithT(t)

	ns := helper.GetTestNamespace()
	c, err := helper.NewSingleCoherenceFromYamlWithSuffix(ns, "suspend-test.yaml", "-disabled-scale-zero")
	g.Expect(err).NotTo(HaveOccurred())

	// Set the flag to NOT suspend on shutdown
	c.Spec.SuspendServicesOnShutdown = ptr.To(false)

	installSimpleDeployment(t, c)

	err = addTestFinalizer(&c)
	g.Expect(err).NotTo(HaveOccurred())

	// Add a finalizer to the StatefulSet to stop it being deleted
	sts := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      c.Name,
			Namespace: c.Namespace,
		},
	}
	err = addTestFinalizer(sts)
	g.Expect(err).NotTo(HaveOccurred())
	// ensure we remove the finalizer
	defer removeAllFinalizersLoggingErrors(t, sts)

	// re-fetch the latest Coherence state and scale down to zero, which should cause services to be suspended
	err = testContext.Client.Get(ctx, c.GetNamespacedName(), &c)
	g.Expect(err).NotTo(HaveOccurred())
	patch := client.RawPatch(types.MergePatchType, []byte(`{"spec":{"replicas":0}}`))
	err = testContext.Client.Patch(testContext.Context, &c, patch)
	g.Expect(err).NotTo(HaveOccurred())

	// The Operator should suspend services and delete the StatefulSet causing its deletion timestamp to be set
	// As we added a finalizer to the StatefulSet it will not actually get deleted yet
	err = waitForStatefulSetDeletionTimestamp(c.GetNamespacedName())
	g.Expect(err).NotTo(HaveOccurred())

	// assert that the cache service is suspended
	svc, err := ManagementOverRestRequest(&c, "/management/coherence/cluster/services/PartitionedCache")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(svc["quorumStatus"]).NotTo(ContainElement("Suspended"))

	// remove the test finalizer from the StatefulSet and Coherence deployment which should then let everything be deleted
	err = removeAllFinalizers(sts)
	g.Expect(err).NotTo(HaveOccurred())
	err = removeAllFinalizers(&c)
	g.Expect(err).NotTo(HaveOccurred())
	// the StatefulSet should eventually be deleted
	err = helper.WaitForDelete(testContext, sts)
	g.Expect(err).NotTo(HaveOccurred())
}

func TestNotSuspendServicesInMultipleDeployments(t *testing.T) {
	// Ensure that everything is cleaned up after the test!
	testContext.CleanupAfterTest(t)
	ctx := context.Background()
	g := NewGomegaWithT(t)

	ns := helper.GetTestNamespace()
	clusterName := "test-cluster"
	cOne, err := helper.NewSingleCoherenceFromYaml(ns, "suspend-test.yaml")
	g.Expect(err).NotTo(HaveOccurred())
	cTwo := cohv1.Coherence{}
	cOne.DeepCopyInto(&cTwo)
	cOne.SetName("test-one")
	cOne.Spec.Cluster = &clusterName
	cTwo.SetName("test-two")
	cTwo.Spec.Cluster = &clusterName

	// install deployment one
	installSimpleDeployment(t, cOne)
	// install deployment two
	installSimpleDeployment(t, cTwo)

	// assert that cluster size is correct
	size := cOne.GetReplicas() + cTwo.GetReplicas()
	data, err := ManagementOverRestRequest(&cOne, "/management/coherence/cluster")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(data["clusterSize"]).To(BeEquivalentTo(size))

	// Delete deployment two, which should cause services to be suspended
	err = testContext.Client.Delete(ctx, &cTwo)
	g.Expect(err).NotTo(HaveOccurred())
	// wait for deployment two to be deleted
	err = helper.WaitForDelete(testContext, &cTwo)
	g.Expect(err).NotTo(HaveOccurred())

	// assert that the cache service is NOT suspended
	svc, err := ManagementOverRestRequest(&cOne, "/management/coherence/cluster/services/PartitionedCache")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(svc["quorumStatus"]).NotTo(ContainElement("Suspended"))
}

func waitForFinalizerTasks(n types.NamespacedName) error {
	ctx := context.Background()
	// wait for the Operator finalizer to be removed which signals that the Operator finalization
	// is complete and services should be suspended.
	c := &cohv1.Coherence{}
	return wait.PollUntilContextTimeout(context.Background(), time.Second, 5*time.Minute, true, func(context.Context) (done bool, err error) {
		if err := testContext.Client.Get(ctx, n, c); err != nil {
			return false, err
		}
		return utils.StringArrayDoesNotContain(c.GetFinalizers(), cohv1.CoherenceFinalizer), nil
	})
}

func waitForStatefulSetDeletionTimestamp(n types.NamespacedName) error {
	ctx := context.Background()
	sts := &appsv1.StatefulSet{}
	return wait.PollUntilContextTimeout(context.Background(), time.Second, 5*time.Minute, true, func(context.Context) (done bool, err error) {
		if err := testContext.Client.Get(ctx, n, sts); err != nil {
			return false, err
		}
		return sts.GetDeletionTimestamp() != nil, nil
	})
}

func addTestFinalizer(o client.Object) error {
	ctx := context.Background()
	k := helper.ObjectKey(o)
	if err := testContext.Client.Get(ctx, k, o); err != nil {
		return err
	}
	s := `{"metadata":{"finalizers":[`
	for _, f := range o.GetFinalizers() {
		s += fmt.Sprintf(`"%s",`, f)
	}
	s += fmt.Sprintf(`"%s"]}}`, testFinalizer)

	patch := client.RawPatch(types.MergePatchType, []byte(s))
	return testContext.Client.Patch(ctx, o, patch)
}

func removeAllFinalizers(o client.Object) error {
	ctx := context.Background()
	k := helper.ObjectKey(o)
	o.DeepCopyObject()
	if err := testContext.Client.Get(ctx, k, o); err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	patch := client.RawPatch(types.MergePatchType, []byte(`{"metadata":{"finalizers":[]}}`))
	err := testContext.Client.Patch(ctx, o, patch)
	if errors.IsNotFound(err) {
		return nil
	}
	return err
}

func removeAllFinalizersLoggingErrors(t *testing.T, o client.Object) {
	if err := removeAllFinalizers(o); err != nil {
		t.Logf("Error removing finalizer from %s/%s - %s", o.GetNamespace(), o.GetName(), err.Error())
	}
}

func ManagementOverRestRequest(c *cohv1.Coherence, path string) (map[string]interface{}, error) {
	pods, err := helper.ListCoherencePodsForDeployment(testContext, c.Namespace, c.Name)
	if err != nil {
		return nil, err
	}
	if len(pods) == 0 {
		return nil, fmt.Errorf("could not find any Pods for Coherence deployment %s", c.Name)
	}

	pf, ports, err := helper.StartPortForwarderForPodWithBackoff(&pods[0])
	if err != nil {
		return nil, err
	}
	defer pf.Close()

	var sep string
	if path[0] == '/' {
		sep = ""
	} else {
		sep = "/"
	}

	url := fmt.Sprintf("http://%s:%d%s%s", pf.Hostname, ports[cohv1.PortNameManagement], sep, path)
	var resp *http.Response

	// try a max of 5 times
	cl := &http.Client{}
	for i := 0; i < 5; i++ {
		resp, err = cl.Get(url)
		if err == nil {
			break
		}
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request returned non-200 status %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	m := make(map[string]interface{})
	err = json.Unmarshal(data, &m)
	return m, err
}
