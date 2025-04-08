/*
 * Copyright (c) 2020, 2025, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package remote

import (
	"context"
	"encoding/json"
	"fmt"
	. "github.com/onsi/gomega"
	cohv1 "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/pkg/utils"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	"io"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/utils/ptr"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"testing"
	"time"
)

const testFinalizer = "coherence.oracle.com/test"

func TestSuspendServices(t *testing.T) {
	// Ensure that everything is cleaned up after the test!
	testContext.CleanupAfterTest(t)
	ctx := context.Background()
	g := NewGomegaWithT(t)

	ns := helper.GetTestNamespace()
	c, err := helper.NewSingleCoherenceFromYamlWithSuffix(ns, "suspend-test.yaml", "-suspend")
	g.Expect(err).NotTo(HaveOccurred())

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
	g.Expect(svc["quorumStatus"]).To(ContainElement("Suspended"))

	// remove the test finalizer which should then let everything be deleted
	err = removeAllFinalizers(&c)
	g.Expect(err).NotTo(HaveOccurred())
	// the StatefulSet should eventually be deleted
	err = helper.WaitForDelete(testContext, sts)
	g.Expect(err).NotTo(HaveOccurred())
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
	controllerutil.AddFinalizer(o, testFinalizer)
	return testContext.Client.Update(ctx, o)
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
