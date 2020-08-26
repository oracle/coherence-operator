/*
 * Copyright (c) 2020, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1_test

import (
	"fmt"
	"github.com/go-test/deep"
	. "github.com/onsi/gomega"
	coh "github.com/oracle/coherence-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sort"
	"testing"
)

func TestCreateServicesWithAdditionalPortsEmpty(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		Ports: []coh.NamedPortSpec{},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)

	// assert that the Services are as expected
	assertService(t, deployment)
}

func TestCreateServicesWithPortsWithOneAdditionalPortWithServiceEnabledFalse(t *testing.T) {

	protocol := corev1.ProtocolUDP

	spec := coh.CoherenceResourceSpec{
		Ports: []coh.NamedPortSpec{
			{
				Name:     "test-port-one",
				Port:     9876,
				Protocol: &protocol,
				NodePort: int32Ptr(2020),
				HostPort: int32Ptr(1234),
				HostIP:   stringPtr("10.10.1.0"),
				Service: &coh.ServiceSpec{
					Enabled: boolPtr(false),
				},
			},
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)

	// assert that the Services are as expected
	assertService(t, deployment)
}

func TestCreateServicesWithPortsWithOneAdditionalPort(t *testing.T) {

	protocol := corev1.ProtocolUDP

	spec := coh.CoherenceResourceSpec{
		Ports: []coh.NamedPortSpec{
			{
				Name:     "test-port-one",
				Port:     9876,
				Protocol: &protocol,
				NodePort: int32Ptr(2020),
				HostPort: int32Ptr(1234),
				HostIP:   stringPtr("10.10.1.0"),
				Service: &coh.ServiceSpec{
					Enabled: boolPtr(true),
				},
			},
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)

	// Create the expected labels
	labels := deployment.CreateCommonLabels()
	labels[coh.LabelComponent] = coh.LabelComponentPortService
	labels[coh.LabelPort] = "test-port-one"

	// Create the expected service selector labels
	selectorLabels := deployment.CreateCommonLabels()
	selectorLabels[coh.LabelComponent] = coh.LabelComponentCoherencePod

	// Create expected Service
	svcExpected := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:   fmt.Sprintf("%s-test-port-one", deployment.Name),
			Labels: labels,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       "test-port-one",
					Protocol:   protocol,
					Port:       9876,
					TargetPort: intstr.FromInt(9876),
					NodePort:   2020,
				},
			},
			Selector: selectorLabels,
		},
	}

	// assert that the Services are as expected
	assertService(t, deployment, &svcExpected)
}

func TestCreateServicesWithPortsWithOneAdditionalPortWithServiceName(t *testing.T) {

	protocol := corev1.ProtocolUDP

	spec := coh.CoherenceResourceSpec{
		Ports: []coh.NamedPortSpec{
			{
				Name:     "test-port-one",
				Port:     9876,
				Protocol: &protocol,
				NodePort: int32Ptr(2020),
				HostPort: int32Ptr(1234),
				HostIP:   stringPtr("10.10.1.0"),
				Service: &coh.ServiceSpec{
					Enabled: boolPtr(true),
					Name:    stringPtr("test-service"),
				},
			},
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)

	// Create the expected labels
	labels := deployment.CreateCommonLabels()
	labels[coh.LabelComponent] = coh.LabelComponentPortService
	labels[coh.LabelPort] = "test-port-one"

	// Create the expected service selector labels
	selectorLabels := deployment.CreateCommonLabels()
	selectorLabels[coh.LabelComponent] = coh.LabelComponentCoherencePod

	// Create expected Service
	svcExpected := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "test-service",
			Labels: labels,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       "test-port-one",
					Protocol:   protocol,
					Port:       9876,
					TargetPort: intstr.FromInt(9876),
					NodePort:   2020,
				},
			},
			Selector: selectorLabels,
		},
	}

	// assert that the Services are as expected
	assertService(t, deployment, &svcExpected)
}

func TestCreateServicesWithPortsWithOneAdditionalPortWithServicePort(t *testing.T) {

	protocol := corev1.ProtocolUDP

	spec := coh.CoherenceResourceSpec{
		Ports: []coh.NamedPortSpec{
			{
				Name:     "test-port-one",
				Port:     9876,
				Protocol: &protocol,
				NodePort: int32Ptr(2020),
				HostPort: int32Ptr(1234),
				HostIP:   stringPtr("10.10.1.0"),
				Service: &coh.ServiceSpec{
					Enabled: boolPtr(true),
					Port:    int32Ptr(80),
				},
			},
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)

	// Create the expected labels
	labels := deployment.CreateCommonLabels()
	labels[coh.LabelComponent] = coh.LabelComponentPortService
	labels[coh.LabelPort] = "test-port-one"

	// Create the expected service selector labels
	selectorLabels := deployment.CreateCommonLabels()
	selectorLabels[coh.LabelComponent] = coh.LabelComponentCoherencePod

	// Create expected Service
	svcExpected := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:   fmt.Sprintf("%s-test-port-one", deployment.Name),
			Labels: labels,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       "test-port-one",
					Protocol:   protocol,
					Port:       80,
					TargetPort: intstr.FromInt(9876),
					NodePort:   2020,
				},
			},
			Selector: selectorLabels,
		},
	}

	// assert that the Services are as expected
	assertService(t, deployment, &svcExpected)
}

func TestCreateServicesWithPortsWithOneAdditionalPortWithServiceFields(t *testing.T) {

	protocol := corev1.ProtocolUDP
	svcType := corev1.ServiceTypeNodePort
	trafficPolicy := corev1.ServiceExternalTrafficPolicyTypeLocal
	ipFamily := corev1.IPv4Protocol
	affinity := corev1.ServiceAffinityNone
	cfg := corev1.SessionAffinityConfig{
		ClientIP: &corev1.ClientIPConfig{
			TimeoutSeconds: int32Ptr(1000),
		},
	}

	spec := coh.CoherenceResourceSpec{
		Ports: []coh.NamedPortSpec{
			{
				Name:     "test-port-one",
				Port:     9876,
				Protocol: &protocol,
				NodePort: int32Ptr(2020),
				HostPort: int32Ptr(1234),
				HostIP:   stringPtr("10.10.1.0"),
				Service: &coh.ServiceSpec{
					Enabled:                  boolPtr(true),
					Type:                     &svcType,
					ClusterIP:                stringPtr("192.168.1.30"),
					ExternalIPs:              []string{"10.10.10.99", "10.10.10.100"},
					LoadBalancerIP:           stringPtr("10.99.0.0"),
					LoadBalancerSourceRanges: []string{"10.10.10.0", "10.10.10.255"},
					ExternalName:             stringPtr("test-external-name"),
					HealthCheckNodePort:      int32Ptr(1000),
					PublishNotReadyAddresses: boolPtr(true),
					ExternalTrafficPolicy:    &trafficPolicy,
					SessionAffinity:          &affinity,
					SessionAffinityConfig:    &cfg,
					IPFamily:                 &ipFamily,
				},
			},
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)

	// Create the expected labels
	labels := deployment.CreateCommonLabels()
	labels[coh.LabelComponent] = coh.LabelComponentPortService
	labels[coh.LabelPort] = "test-port-one"

	// Create the expected service selector labels
	selectorLabels := deployment.CreateCommonLabels()
	selectorLabels[coh.LabelComponent] = coh.LabelComponentCoherencePod

	// Create expected Service
	svcExpected := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:   fmt.Sprintf("%s-test-port-one", deployment.Name),
			Labels: labels,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       "test-port-one",
					Protocol:   protocol,
					Port:       9876,
					TargetPort: intstr.FromInt(9876),
					NodePort:   2020,
				},
			},
			Selector:                 selectorLabels,
			Type:                     svcType,
			ClusterIP:                "192.168.1.30",
			ExternalIPs:              []string{"10.10.10.99", "10.10.10.100"},
			LoadBalancerIP:           "10.99.0.0",
			LoadBalancerSourceRanges: []string{"10.10.10.0", "10.10.10.255"},
			ExternalName:             "test-external-name",
			HealthCheckNodePort:      1000,
			PublishNotReadyAddresses: true,
			ExternalTrafficPolicy:    trafficPolicy,
			SessionAffinity:          affinity,
			SessionAffinityConfig:    &cfg,
			IPFamily:                 &ipFamily,
		},
	}

	// assert that the Services are as expected
	assertService(t, deployment, &svcExpected)
}

func TestCreateServicesWithPortsWithOneAdditionalPortWithServiceLabels(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		Ports: []coh.NamedPortSpec{
			{
				Name: "test-port-one",
				Port: 9876,
				Service: &coh.ServiceSpec{
					Enabled: boolPtr(true),
					Labels:  map[string]string{"LabelOne": "One", "LabelTwo": "Two"},
				},
			},
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)

	// Create the expected labels
	labels := deployment.CreateCommonLabels()
	labels[coh.LabelComponent] = coh.LabelComponentPortService
	labels[coh.LabelPort] = "test-port-one"
	labels["LabelOne"] = "One"
	labels["LabelTwo"] = "Two"

	// Create the expected service selector labels
	selectorLabels := deployment.CreateCommonLabels()
	selectorLabels[coh.LabelComponent] = coh.LabelComponentCoherencePod

	// Create expected Service
	svcExpected := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:   fmt.Sprintf("%s-test-port-one", deployment.Name),
			Labels: labels,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       "test-port-one",
					Port:       9876,
					TargetPort: intstr.FromInt(9876),
					Protocol:   corev1.ProtocolTCP,
				},
			},
			Selector: selectorLabels,
		},
	}

	// assert that the Services are as expected
	assertService(t, deployment, &svcExpected)
}

func TestCreateServicesWithPortsWithOneAdditionalPortWithServiceAnnotations(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		Ports: []coh.NamedPortSpec{
			{
				Name: "test-port-one",
				Port: 9876,
				Service: &coh.ServiceSpec{
					Enabled:     boolPtr(true),
					Annotations: map[string]string{"AnnOne": "One", "AnnTwo": "Two"},
				},
			},
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)

	// Create the expected labels
	labels := deployment.CreateCommonLabels()
	labels[coh.LabelComponent] = coh.LabelComponentPortService
	labels[coh.LabelPort] = "test-port-one"

	// Create the expected service selector labels
	selectorLabels := deployment.CreateCommonLabels()
	selectorLabels[coh.LabelComponent] = coh.LabelComponentCoherencePod

	// Create expected Service
	svcExpected := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        fmt.Sprintf("%s-test-port-one", deployment.Name),
			Labels:      labels,
			Annotations: map[string]string{"AnnOne": "One", "AnnTwo": "Two"},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       "test-port-one",
					Port:       9876,
					TargetPort: intstr.FromInt(9876),
					Protocol:   corev1.ProtocolTCP,
				},
			},
			Selector: selectorLabels,
		},
	}

	// assert that the Services are as expected
	assertService(t, deployment, &svcExpected)
}

func TestCreateServicesWithPortsWithTwoAdditionalPorts(t *testing.T) {

	protocolOne := corev1.ProtocolUDP
	protocolTwo := corev1.ProtocolSCTP

	spec := coh.CoherenceResourceSpec{
		Ports: []coh.NamedPortSpec{
			{
				Name:     "test-port-one",
				Port:     9876,
				Protocol: &protocolOne,
			},
			{
				Name:     "test-port-two",
				Port:     5678,
				Protocol: &protocolTwo,
			},
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)

	// Create the expected labelsOne
	labelsOne := deployment.CreateCommonLabels()
	labelsOne[coh.LabelComponent] = coh.LabelComponentPortService
	labelsOne[coh.LabelPort] = "test-port-one"

	// Create the expected labelsOne
	labelsTwo := deployment.CreateCommonLabels()
	labelsTwo[coh.LabelComponent] = coh.LabelComponentPortService
	labelsTwo[coh.LabelPort] = "test-port-two"

	// Create the expected service selector labelsOne
	selectorLabels := deployment.CreateCommonLabels()
	selectorLabels[coh.LabelComponent] = coh.LabelComponentCoherencePod

	// Create expected first Service
	svcExpectedOne := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:   fmt.Sprintf("%s-test-port-one", deployment.Name),
			Labels: labelsOne,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       "test-port-one",
					Port:       9876,
					TargetPort: intstr.FromInt(9876),
					Protocol:   protocolOne,
				},
			},
			Selector: selectorLabels,
		},
	}

	// Create expected second Service
	svcExpectedTwo := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:   fmt.Sprintf("%s-test-port-two", deployment.Name),
			Labels: labelsTwo,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       "test-port-two",
					Port:       5678,
					TargetPort: intstr.FromInt(5678),
					Protocol:   protocolTwo,
				},
			},
			Selector: selectorLabels,
		},
	}

	// assert that the Services are as expected
	assertService(t, deployment, &svcExpectedOne, &svcExpectedTwo)
}

func assertService(t *testing.T, deployment *coh.Coherence, servicesExpected ...metav1.Object) {
	g := NewGomegaWithT(t)

	res := deployment.Spec.CreateServicesForPort(deployment)

	// Sort the expected services
	sort.SliceStable(servicesExpected, func(i, j int) bool {
		return servicesExpected[i].GetName() < servicesExpected[j].GetName()
	})

	// Sort the actual services
	sort.SliceStable(res, func(i, j int) bool {
		return res[i].Name < res[j].Name
	})

	var resExpected []coh.Resource
	for _, r := range servicesExpected {
		resExpected = append(resExpected, coh.Resource{
			Kind: coh.ResourceTypeService,
			Name: r.GetName(),
			Spec: r.(runtime.Object),
		})
	}
	diffs := deep.Equal(res, resExpected)
	g.Expect(diffs).To(BeNil())
}
