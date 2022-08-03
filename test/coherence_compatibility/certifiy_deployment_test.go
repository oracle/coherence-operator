/*
 * Copyright (c) 2021, 2022 Oracle and/or its affiliates.
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
		Spec: v1.CoherenceResourceSpec{
			Image: &imageName,
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
