/*
 * Copyright (c) 2020, 2024, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1_test

import (
	"github.com/go-test/deep"
	. "github.com/onsi/gomega"
	coh "github.com/oracle/coherence-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	"testing"
)

func TestCreateHeadlessServiceForMinimalDeployment(t *testing.T) {
	// Create the test deployment
	deployment := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "test-ns",
			Name:      "test",
		},
	}

	// create the expected WKA service
	labels := deployment.CreateCommonLabels()
	labels[coh.LabelCoherenceCluster] = "test"
	labels[coh.LabelComponent] = coh.LabelComponentCoherenceHeadless

	// The selector for the service (match all Pods with the same cluster label)
	selector := make(map[string]string)
	selector[coh.LabelCoherenceDeployment] = "test"
	selector[coh.LabelCoherenceCluster] = "test"
	selector[coh.LabelCoherenceRole] = "test"
	selector[coh.LabelComponent] = coh.LabelComponentCoherencePod

	expected := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "test-ns",
			Name:      "test-sts",
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			ClusterIP:                corev1.ClusterIPNone,
			PublishNotReadyAddresses: true,
			Ports:                    getDefaultServicePorts(),
			Selector:                 selector,
		},
	}

	// assert that the Services are as expected
	assertHeadlessService(t, deployment, expected)
}

func TestCreateHeadlessServiceWithSingleIPFamily(t *testing.T) {
	// Create the test deployment
	deployment := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "test-ns",
			Name:      "test",
		},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			HeadlessServiceIpFamilies: []corev1.IPFamily{
				corev1.IPv6Protocol,
			},
		},
	}

	// create the expected WKA service
	labels := deployment.CreateCommonLabels()
	labels[coh.LabelCoherenceCluster] = "test"
	labels[coh.LabelComponent] = coh.LabelComponentCoherenceHeadless

	// The selector for the service (match all Pods with the same cluster label)
	selector := make(map[string]string)
	selector[coh.LabelCoherenceDeployment] = "test"
	selector[coh.LabelCoherenceCluster] = "test"
	selector[coh.LabelCoherenceRole] = "test"
	selector[coh.LabelComponent] = coh.LabelComponentCoherencePod

	expected := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "test-ns",
			Name:      "test-sts",
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			ClusterIP:                corev1.ClusterIPNone,
			PublishNotReadyAddresses: true,
			Ports:                    getDefaultServicePorts(),
			Selector:                 selector,
			IPFamilyPolicy:           ptr.To(corev1.IPFamilyPolicySingleStack),
			IPFamilies:               []corev1.IPFamily{corev1.IPv6Protocol},
		},
	}

	// assert that the Services are as expected
	assertHeadlessService(t, deployment, expected)
}

func TestCreateHeadlessServiceWithDualStackIPFamily(t *testing.T) {
	// Create the test deployment
	deployment := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "test-ns",
			Name:      "test",
		},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			HeadlessServiceIpFamilies: []corev1.IPFamily{
				corev1.IPv4Protocol,
				corev1.IPv6Protocol,
			},
		},
	}

	// create the expected WKA service
	labels := deployment.CreateCommonLabels()
	labels[coh.LabelCoherenceCluster] = "test"
	labels[coh.LabelComponent] = coh.LabelComponentCoherenceHeadless

	// The selector for the service (match all Pods with the same cluster label)
	selector := make(map[string]string)
	selector[coh.LabelCoherenceDeployment] = "test"
	selector[coh.LabelCoherenceCluster] = "test"
	selector[coh.LabelCoherenceRole] = "test"
	selector[coh.LabelComponent] = coh.LabelComponentCoherencePod

	expected := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "test-ns",
			Name:      "test-sts",
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			ClusterIP:                corev1.ClusterIPNone,
			PublishNotReadyAddresses: true,
			Ports:                    getDefaultServicePorts(),
			Selector:                 selector,
			IPFamilyPolicy:           ptr.To(corev1.IPFamilyPolicyPreferDualStack),
			IPFamilies:               []corev1.IPFamily{corev1.IPv4Protocol, corev1.IPv6Protocol},
		},
	}

	// assert that the Services are as expected
	assertHeadlessService(t, deployment, expected)
}

func assertHeadlessService(t *testing.T, deployment *coh.Coherence, expected *corev1.Service) {
	g := NewGomegaWithT(t)

	resActual := deployment.Spec.CreateHeadlessService(deployment)
	resExpected := coh.Resource{
		Kind: coh.ResourceTypeService,
		Name: expected.GetName(),
		Spec: expected,
	}

	diffs := deep.Equal(resActual, resExpected)
	g.Expect(diffs).To(BeNil())
}
