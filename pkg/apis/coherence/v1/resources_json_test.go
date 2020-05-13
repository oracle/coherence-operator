/*
 * Copyright (c) 2020, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1_test

import (
	"encoding/json"
	"github.com/go-test/deep"
	"github.com/oracle/coherence-operator/pkg/apis"
	v1 "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubectl/pkg/scheme"
	"testing"

	. "github.com/onsi/gomega"
)

func TestSerializeResources(t *testing.T) {
	om := metav1.ObjectMeta{
		Namespace: "operator-test",
		Name:      "foo",
	}

	_ = apis.AddToScheme(scheme.Scheme)

	resources := []v1.Resource{
		{
			Kind: v1.ResourceTypeDeployment,
			Name: "foo",
			Spec: &v1.CoherenceDeployment{ObjectMeta: om},
		},
		{
			Kind: v1.ResourceTypeConfigMap,
			Name: "foo",
			Spec: &corev1.ConfigMap{ObjectMeta: om},
		},
		{
			Kind: v1.ResourceTypeSecret,
			Name: "foo",
			Spec: &corev1.Secret{ObjectMeta: om},
		},
		{
			Kind: v1.ResourceTypeService,
			Name: "foo",
			Spec: &corev1.Service{ObjectMeta: om},
		},
		{
			Kind: v1.ResourceTypeService,
			Name: "foo",
			Spec: &appsv1.StatefulSet{ObjectMeta: om},
		},
	}

	for _, resource := range resources {
		t.Run(resource.Kind.Name(), func(t *testing.T) {
			AssertResourcesRoundTrip(t, v1.Resources{Version: 1, Items: []v1.Resource{resource}})
		})
	}
}

func TestSerializeMultipleResources(t *testing.T) {
	om := metav1.ObjectMeta{
		Namespace: "operator-test",
		Name:      "foo",
	}

	_ = apis.AddToScheme(scheme.Scheme)

	resources := []v1.Resource{
		{
			Kind: v1.ResourceTypeDeployment,
			Name: "foo",
			Spec: &v1.CoherenceDeployment{ObjectMeta: om},
		},
		{
			Kind: v1.ResourceTypeConfigMap,
			Name: "foo",
			Spec: &corev1.ConfigMap{ObjectMeta: om},
		},
		{
			Kind: v1.ResourceTypeSecret,
			Name: "foo",
			Spec: &corev1.Secret{ObjectMeta: om},
		},
		{
			Kind: v1.ResourceTypeService,
			Name: "foo",
			Spec: &corev1.Service{ObjectMeta: om},
		},
		{
			Kind: v1.ResourceTypeService,
			Name: "foo",
			Spec: &appsv1.StatefulSet{ObjectMeta: om},
		},
	}

	AssertResourcesRoundTrip(t, v1.Resources{Version: 1, Items: resources})
}

func AssertResourcesRoundTrip(t *testing.T, in v1.Resources) {
	g := NewGomegaWithT(t)
	result := v1.Resources{}

	in.EnsureGVK(scheme.Scheme)

	b, err := json.Marshal(in)
	g.Expect(err).NotTo(HaveOccurred())

	err = json.Unmarshal(b, &result)
	g.Expect(err).NotTo(HaveOccurred())

	diff := deep.Equal(result, result)
	g.Expect(diff).To(BeEmpty())
}
