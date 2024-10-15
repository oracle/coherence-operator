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
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/ptr"
	"testing"
)

func TestCreateWKAServiceForMinimalDeployment(t *testing.T) {
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
	assertWKAService(t, deployment, expected)
}

func TestCreateWKAServiceWithAppLabel(t *testing.T) {
	// Create the test deployment
	deployment := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "test-ns",
			Name:      "test",
		},
		Spec: coh.CoherenceStatefulSetResourceSpec{
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
	assertWKAService(t, deployment, expected)
}

func TestCreateWKAServiceWithVersionLabel(t *testing.T) {
	// Create the test deployment
	deployment := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "test-ns",
			Name:      "test",
		},
		Spec: coh.CoherenceStatefulSetResourceSpec{
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
	assertWKAService(t, deployment, expected)
}

func TestCreateWKAServiceForDeploymentWithClusterName(t *testing.T) {
	// Create the test deployment
	deployment := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "test-ns",
			Name:      "test",
		},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			Cluster: ptr.To("test-cluster"),
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
	assertWKAService(t, deployment, expected)
}

func TestCreateWKAServiceForDeploymentWithAdditionalLabels(t *testing.T) {
	extraLabels := make(map[string]string)
	extraLabels["one"] = "label-one"
	extraLabels["two"] = "label-two"

	// Create the test deployment
	deployment := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "test-ns",
			Name:      "test",
		},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Coherence: &coh.CoherenceSpec{
					WKA: &coh.CoherenceWKASpec{
						Labels: extraLabels,
					},
				},
			},
			Cluster: ptr.To("test-cluster"),
		},
	}

	// create the expected WKA service
	labels := deployment.CreateCommonLabels()
	labels[coh.LabelCoherenceCluster] = "test-cluster"
	labels[coh.LabelComponent] = coh.LabelComponentWKA
	labels["one"] = "label-one"
	labels["two"] = "label-two"

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
	assertWKAService(t, deployment, expected)
}

func TestCreateWKAServiceForDeploymentWithAdditionalAnnotations(t *testing.T) {
	extraAnnotations := make(map[string]string)
	extraAnnotations["one"] = "label-one"
	extraAnnotations["two"] = "label-two"

	// Create the test deployment
	deployment := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "test-ns",
			Name:      "test",
		},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Coherence: &coh.CoherenceSpec{
					WKA: &coh.CoherenceWKASpec{
						Annotations: extraAnnotations,
					},
				},
			},
			Cluster: ptr.To("test-cluster"),
		},
	}

	// create the expected WKA service
	labels := deployment.CreateCommonLabels()
	labels[coh.LabelCoherenceCluster] = "test-cluster"
	labels[coh.LabelComponent] = coh.LabelComponentWKA

	ann := make(map[string]string)
	ann["service.alpha.kubernetes.io/tolerate-unready-endpoints"] = "true"
	ann["one"] = "label-one"
	ann["two"] = "label-two"

	// The selector for the service (match all Pods with the same cluster label)
	selector := make(map[string]string)
	selector[coh.LabelCoherenceCluster] = "test-cluster"
	selector[coh.LabelComponent] = coh.LabelComponentCoherencePod
	selector[coh.LabelCoherenceWKAMember] = "true"

	expected := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:   "test-ns",
			Name:        "test-wka",
			Labels:      labels,
			Annotations: ann,
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
	assertWKAService(t, deployment, expected)
}

func TestCreateWKAServiceWithIPFamily(t *testing.T) {
	// Create the test deployment
	deployment := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "test-ns",
			Name:      "test",
		},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Coherence: &coh.CoherenceSpec{
					WKA: &coh.CoherenceWKASpec{
						IPFamily: ptr.To(corev1.IPv4Protocol),
					},
				},
			},
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
			IPFamilyPolicy:           ptr.To(corev1.IPFamilyPolicySingleStack),
			IPFamilies:               []corev1.IPFamily{corev1.IPv4Protocol},
		},
	}

	// assert that the Services are as expected
	assertWKAService(t, deployment, expected)
}

func assertWKAService(t *testing.T, deployment *coh.Coherence, expected *corev1.Service) {
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

func getDefaultServicePorts() []corev1.ServicePort {
	return []corev1.ServicePort{
		{
			Name:        coh.PortNameCoherence,
			Protocol:    corev1.ProtocolTCP,
			AppProtocol: ptr.To(coh.AppProtocolTcp),
			Port:        7,
			TargetPort:  intstr.FromInt32(7),
		},
		{
			Name:        coh.PortNameCoherenceLocal,
			Protocol:    corev1.ProtocolTCP,
			AppProtocol: ptr.To(coh.AppProtocolTcp),
			Port:        coh.DefaultUnicastPort,
			TargetPort:  intstr.FromString(coh.PortNameCoherenceLocal),
		},
		{
			Name:        coh.PortNameCoherenceCluster,
			Protocol:    corev1.ProtocolTCP,
			AppProtocol: ptr.To(coh.AppProtocolTcp),
			Port:        coh.DefaultClusterPort,
			TargetPort:  intstr.FromString(coh.PortNameCoherenceCluster),
		},
		{
			Name:        coh.PortNameHealth,
			Protocol:    corev1.ProtocolTCP,
			AppProtocol: ptr.To(coh.AppProtocolHttp),
			Port:        coh.DefaultHealthPort,
			TargetPort:  intstr.FromString(coh.PortNameHealth),
		},
	}
}
