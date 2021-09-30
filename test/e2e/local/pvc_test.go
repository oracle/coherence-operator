/*
 * Copyright (c) 2019, 2021, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package local

import (
	. "github.com/onsi/gomega"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestAdditionalVolumeClaimTemplates(t *testing.T) {
	g := NewGomegaWithT(t)
	testContext.CleanupAfterTest(t)

	ns := helper.GetTestNamespace()
	pvcName := "data-volume-storage-with-pvc-0"

	defer deletePVC(ns, pvcName)

	helper.AssertDeployments(testContext, t, "pvc.yaml")

	g.Expect(ns).NotTo(BeEmpty())
	pvc, err := testContext.KubeClient.CoreV1().PersistentVolumeClaims(ns).
		Get(testContext.Context, pvcName, metav1.GetOptions{})
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(pvc).NotTo(BeNil())
	g.Expect(pvc.Status.Phase).To(Equal(corev1.ClaimBound))

	labels := pvc.GetLabels()
	g.Expect(labels["coherence.oracle.com/test1"]).To(Equal("pvc-test-1"))
	g.Expect(labels["coherence.oracle.com/test2"]).To(Equal("pvc-test-2"))

	annotations := pvc.GetAnnotations()
	g.Expect(annotations["test-key1"]).To(Equal("test-value-1"))
	g.Expect(annotations["test-key2"]).To(Equal("test-value-2"))
}

func deletePVC(namespace, name string) {
	_ = testContext.KubeClient.CoreV1().PersistentVolumeClaims(namespace).
		Delete(testContext.Context, name, metav1.DeleteOptions{})
}
