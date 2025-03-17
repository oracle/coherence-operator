/*
 * Copyright (c) 2020, 2025, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package compatibility

import (
	"context"
	. "github.com/onsi/gomega"
	v1 "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	"testing"
	"time"
)

func TestCoherenceCompatibilityMinimalSpec(t *testing.T) {
	// Ensure that everything is cleaned up after the test!
	testContext.CleanupAfterTest(t)

	g := NewGomegaWithT(t)

	imageName := helper.GetCoherenceCompatibilityImage()

	ns := helper.GetTestNamespace()
	d := &v1.Coherence{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: ns,
			Name:      "certify-coherence",
		},
		Spec: v1.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: v1.CoherenceResourceSpec{
				Image: &imageName,
			},
		},
	}

	// ensure the imagePullSecrets are correctly injected
	secrets := helper.GetImagePullSecrets()
	d.Spec.ImagePullSecrets = append(d.Spec.ImagePullSecrets, secrets...)

	err := testContext.Client.Create(context.TODO(), d)
	g.Expect(err).NotTo(HaveOccurred())

	_, err = helper.WaitForStatefulSetForDeployment(testContext, ns, d, time.Second*10, time.Minute*5)
	g.Expect(err).NotTo(HaveOccurred())
}

func TestCoherenceCompatibilityScaling(t *testing.T) {
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

	if helper.GetTestCoherenceIsJava8() {
		d.Spec.JVM = &v1.JVMSpec{
			Java8: ptr.To(true),
		}
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

// This test is not really using Java8 for images using Coherence 14.1.1-2206
// and above. The test just sets the Java8 field to true to test running with
// the legacy container entry point.
func TestCoherenceCompatibilityJava8(t *testing.T) {
	// Ensure that everything is cleaned up after the test!
	testContext.CleanupAfterTest(t)

	g := NewGomegaWithT(t)

	ns := helper.GetTestNamespace()
	d := &v1.Coherence{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: ns,
			Name:      "certify-coherence",
		},
		Spec: v1.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: v1.CoherenceResourceSpec{
				JVM: &v1.JVMSpec{
					Java8: ptr.To(true),
				},
			},
		},
	}

	// ensure the imagePullSecrets are correctly injected
	secrets := helper.GetImagePullSecrets()
	d.Spec.ImagePullSecrets = append(d.Spec.ImagePullSecrets, secrets...)

	err := testContext.Client.Create(context.TODO(), d)
	g.Expect(err).NotTo(HaveOccurred())

	_, err = helper.WaitForStatefulSetForDeployment(testContext, ns, d, time.Second*10, time.Minute*5)
	g.Expect(err).NotTo(HaveOccurred())
}
