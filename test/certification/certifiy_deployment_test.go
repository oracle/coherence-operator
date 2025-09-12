/*
 * Copyright (c) 2020, 2025, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package certification

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	. "github.com/onsi/gomega"
	v1 "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
)

func TestCertifyMinimalSpec(t *testing.T) {
	// Ensure that everything is cleaned up after the test!
	testContext.CleanupAfterTest(t)

	g := NewGomegaWithT(t)

	ns := helper.GetTestClusterNamespace()
	d := &v1.Coherence{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: ns,
			Name:      "certify-minimal",
		},
	}

	err := testContext.Client.Create(context.TODO(), d)
	g.Expect(err).NotTo(HaveOccurred())

	_, err = helper.WaitForStatefulSetForDeployment(testContext, ns, d, time.Second*10, time.Minute*5)
	g.Expect(err).NotTo(HaveOccurred())
}

func TestCertifyScaling(t *testing.T) {
	// Ensure that everything is cleaned up after the test!
	testContext.CleanupAfterTest(t)
	g := NewGomegaWithT(t)

	ns := helper.GetTestClusterNamespace()
	d := &v1.Coherence{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: ns,
			Name:      "certify-scale",
		},
		Spec: v1.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: v1.CoherenceResourceSpec{
				Replicas: ptr.To(int32(1)),
				ReadinessProbe: &v1.ReadinessProbeSpec{
					InitialDelaySeconds: ptr.To(int32(10)),
					PeriodSeconds:       ptr.To(int32(10)),
				},
			},
		},
	}

	// Start with one replica
	err := testContext.Client.Create(context.TODO(), d)
	g.Expect(err).NotTo(HaveOccurred())
	_, err = helper.WaitForStatefulSetForDeployment(testContext, ns, d, time.Second*10, time.Minute*5)
	g.Expect(err).NotTo(HaveOccurred())

	// Scale Up to three
	err = scale(t, ns, d.Name, 3)
	g.Expect(err).NotTo(HaveOccurred())
	_, err = helper.WaitForStatefulSet(testContext, ns, d.Name, 3, time.Second*10, time.Minute*5)
	g.Expect(err).NotTo(HaveOccurred())

	// Scale down to one
	err = scale(t, ns, d.Name, 1)
	g.Expect(err).NotTo(HaveOccurred())
	_, err = helper.WaitForStatefulSet(testContext, ns, d.Name, 1, time.Second*10, time.Minute*5)
	g.Expect(err).NotTo(HaveOccurred())
}

// Test the scenario where we create a Coherence cluster without a "replicas" field, which will default to three Pods.
// Then scale up the cluster to four.
// Then apply an update using the same Coherence resource with no replicas field.
// After the update is applied, the cluster should still be four and not revert to three.
func TestCertifyScalingWithUpdate(t *testing.T) {
	// Ensure that everything is cleaned up after the test!
	testContext.CleanupAfterTest(t)
	g := NewGomegaWithT(t)

	ns := helper.GetTestClusterNamespace()

	// the name of the cluster from scale-with-update-one.yaml and scale-with-update-two.yaml
	name := "certify-scale-update"

	// Start with the default three replicas
	err := apply(t, ns, "scale-with-update-one.yaml")
	g.Expect(err).NotTo(HaveOccurred())
	_, err = helper.WaitForPodsWithLabel(testContext, ns, "one=testOne", 3, time.Second*10, time.Minute*10)
	g.Expect(err).NotTo(HaveOccurred())

	// Scale Up to four
	err = scale(t, ns, name, 4)
	g.Expect(err).NotTo(HaveOccurred())
	_, err = helper.WaitForStatefulSet(testContext, ns, name, 4, time.Second*10, time.Minute*5)
	g.Expect(err).NotTo(HaveOccurred())

	// apply the update
	err = apply(t, ns, "scale-with-update-two.yaml")
	g.Expect(err).NotTo(HaveOccurred())
	// There should eventually be four Pods with the additional label
	_, err = helper.WaitForPodsWithLabel(testContext, ns, "two=testTwo", 4, time.Second*10, time.Minute*10)
	g.Expect(err).NotTo(HaveOccurred())
}

// Test the scenario where we create a Coherence cluster with a "replicas" field.
// Then scale up the cluster to four.
// Then apply an update using the same Coherence resource with no replicas field.
// After the update is applied, the cluster should still be four and not revert to three.
func TestCertifyScalingClusterWithReplicasThenUpdate(t *testing.T) {
	// Ensure that everything is cleaned up after the test!
	testContext.CleanupAfterTest(t)
	g := NewGomegaWithT(t)

	ns := helper.GetTestClusterNamespace()

	// the name of the cluster from scale-with-update-one.yaml and scale-with-update-two.yaml
	name := "certify-scale-update"

	// Start with the two replicas
	err := apply(t, ns, "scale-with-update-with-replicas.yaml")
	g.Expect(err).NotTo(HaveOccurred())
	_, err = helper.WaitForPodsWithLabel(testContext, ns, "one=testOne", 2, time.Second*10, time.Minute*10)
	g.Expect(err).NotTo(HaveOccurred())

	// Scale Up to four
	err = scale(t, ns, name, 4)
	g.Expect(err).NotTo(HaveOccurred())
	_, err = helper.WaitForStatefulSet(testContext, ns, name, 4, time.Second*10, time.Minute*5)
	g.Expect(err).NotTo(HaveOccurred())

	// apply the update
	err = apply(t, ns, "scale-with-update-two.yaml")
	g.Expect(err).NotTo(HaveOccurred())
	// There should eventually be four Pods with the additional label
	_, err = helper.WaitForPodsWithLabel(testContext, ns, "two=testTwo", 4, time.Second*10, time.Minute*10)
	g.Expect(err).NotTo(HaveOccurred())
}

func scale(t *testing.T, namespace, name string, replicas int32) error {
	cmd := exec.Command("kubectl", "-n", namespace, "scale", fmt.Sprintf("--replicas=%d", replicas), "coherence/"+name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	t.Log("Executing Scale Command: " + strings.Join(cmd.Args, " "))
	return cmd.Run()
}

func apply(t *testing.T, namespace, fileName string) error {
	actualFile, err := helper.FindActualFile(fileName)
	if err != nil {
		return err
	}
	cmd := exec.Command("kubectl", "-n", namespace, "apply", "-f", actualFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	t.Log("Executing Kubectl Command: " + strings.Join(cmd.Args, " "))
	return cmd.Run()
}
