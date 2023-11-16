/*
 * Copyright (c) 2023, Oracle and/or its affiliates.
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
	"testing"
)

func TestCreateWKAServiceForMinimalJonDeployment(t *testing.T) {
	// Create the test deployment
	deployment := &coh.CoherenceJob{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "test-ns",
			Name:      "test",
		},
	}

	// create the expected WKA service
	labels := deployment.CreateCommonLabels()
	labels[coh.LabelCoherenceCluster] = "test"
	labels[coh.LabelComponent] = coh.LabelComponentWKA

	// The selector for the service (match all Pods with the same cluster label)
	selector := make(map[string]string)
	selector[coh.LabelCoherenceCluster] = "test"
	selector[coh.LabelComponent] = coh.LabelComponentCoherencePod
	selector[coh.LabelCoherenceWKAMember] = "true"

	expected := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "test-ns",
			Name:      "test-wka",
			Labels:    labels,
			Annotations: map[string]string{
				"service.alpha.kubernetes.io/tolerate-unready-endpoints": "true",
			},
		},
		Spec: corev1.ServiceSpec{
			ClusterIP: corev1.ClusterIPNone,
			// Pods must be part of the WKA service even if not ready
			PublishNotReadyAddresses: true,
			Ports:                    getDefaultServicePorts(),
			Selector:                 selector,
		},
	}

	// assert that the Services are as expected
	assertWKAServiceForJob(t, deployment, expected)
}

func TestCreateWKAServiceForJobWithAppLabel(t *testing.T) {
	// Create the test deployment
	deployment := &coh.CoherenceJob{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "test-ns",
			Name:      "test",
		},
		Spec: coh.CoherenceJobResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				AppLabel: stringPtr("foo"),
			},
		},
	}

	// create the expected WKA service
	labels := deployment.CreateCommonLabels()
	labels[coh.LabelCoherenceCluster] = "test"
	labels[coh.LabelComponent] = coh.LabelComponentWKA
	labels[coh.LabelApp] = "foo"

	// The selector for the service (match all Pods with the same cluster label)
	selector := make(map[string]string)
	selector[coh.LabelCoherenceCluster] = "test"
	selector[coh.LabelComponent] = coh.LabelComponentCoherencePod
	selector[coh.LabelCoherenceWKAMember] = "true"

	expected := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "test-ns",
			Name:      "test-wka",
			Labels:    labels,
			Annotations: map[string]string{
				"service.alpha.kubernetes.io/tolerate-unready-endpoints": "true",
			},
		},
		Spec: corev1.ServiceSpec{
			ClusterIP: corev1.ClusterIPNone,
			// Pods must be part of the WKA service even if not ready
			PublishNotReadyAddresses: true,
			Ports:                    getDefaultServicePorts(),
			Selector:                 selector,
		},
	}

	// assert that the Services are as expected
	assertWKAServiceForJob(t, deployment, expected)
}

func TestCreateWKAServiceForJobWithVersionLabel(t *testing.T) {
	// Create the test deployment
	deployment := &coh.CoherenceJob{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "test-ns",
			Name:      "test",
		},
		Spec: coh.CoherenceJobResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				VersionLabel: stringPtr("v1.0.0"),
			},
		},
	}

	// create the expected WKA service
	labels := deployment.CreateCommonLabels()
	labels[coh.LabelCoherenceCluster] = "test"
	labels[coh.LabelComponent] = coh.LabelComponentWKA
	labels[coh.LabelVersion] = "v1.0.0"

	// The selector for the service (match all Pods with the same cluster label)
	selector := make(map[string]string)
	selector[coh.LabelCoherenceCluster] = "test"
	selector[coh.LabelComponent] = coh.LabelComponentCoherencePod
	selector[coh.LabelCoherenceWKAMember] = "true"

	expected := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "test-ns",
			Name:      "test-wka",
			Labels:    labels,
			Annotations: map[string]string{
				"service.alpha.kubernetes.io/tolerate-unready-endpoints": "true",
			},
		},
		Spec: corev1.ServiceSpec{
			ClusterIP: corev1.ClusterIPNone,
			// Pods must be part of the WKA service even if not ready
			PublishNotReadyAddresses: true,
			Ports:                    getDefaultServicePorts(),
			Selector:                 selector,
		},
	}

	// assert that the Services are as expected
	assertWKAServiceForJob(t, deployment, expected)
}

func TestCreateWKAServiceForJobWithClusterName(t *testing.T) {
	// Create the test deployment
	deployment := &coh.CoherenceJob{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "test-ns",
			Name:      "test",
		},
		Spec: coh.CoherenceJobResourceSpec{
			Cluster: "test-cluster",
		},
	}

	// create the expected WKA service
	labels := deployment.CreateCommonLabels()
	labels[coh.LabelCoherenceCluster] = "test-cluster"
	labels[coh.LabelComponent] = coh.LabelComponentWKA

	// The selector for the service (match all Pods with the same cluster label)
	selector := make(map[string]string)
	selector[coh.LabelCoherenceCluster] = "test-cluster"
	selector[coh.LabelComponent] = coh.LabelComponentCoherencePod
	selector[coh.LabelCoherenceWKAMember] = "true"

	expected := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "test-ns",
			Name:      "test-wka",
			Labels:    labels,
			Annotations: map[string]string{
				"service.alpha.kubernetes.io/tolerate-unready-endpoints": "true",
			},
		},
		Spec: corev1.ServiceSpec{
			ClusterIP: corev1.ClusterIPNone,
			// Pods must be part of the WKA service even if not ready
			PublishNotReadyAddresses: true,
			Ports:                    getDefaultServicePorts(),
			Selector:                 selector,
		},
	}

	// assert that the Services are as expected
	assertWKAServiceForJob(t, deployment, expected)
}

func assertWKAServiceForJob(t *testing.T, deployment *coh.CoherenceJob, expected *corev1.Service) {
	g := NewGomegaWithT(t)

	resActual := deployment.Spec.CreateWKAService(deployment)
	resExpected := coh.Resource{
		Kind: coh.ResourceTypeService,
		Name: expected.GetName(),
		Spec: expected,
	}

	diffs := deep.Equal(resActual, resExpected)
	g.Expect(diffs).To(BeNil())
}
