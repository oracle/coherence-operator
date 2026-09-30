/*
 * Copyright (c) 2026, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package controllers

import (
	"context"
	"encoding/json"
	"testing"

	coh "github.com/oracle/coherence-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

type podListCountingClient struct {
	client.Client
	podListCalls int
}

func (c *podListCountingClient) List(ctx context.Context, list client.ObjectList, opts ...client.ListOption) error {
	if _, ok := list.(*corev1.PodList); ok {
		c.podListCalls++
	}
	return c.Client.List(ctx, list, opts...)
}

func TestHealthSecretRequestsListsPodsOnce(t *testing.T) {
	const (
		namespace  = "test"
		secretName = "health-secret"
	)
	scheme := runtime.NewScheme()
	if err := corev1.AddToScheme(scheme); err != nil {
		t.Fatal(err)
	}
	if err := coh.AddToScheme(scheme); err != nil {
		t.Fatal(err)
	}

	connection := func(name string) *coh.HealthMutatorConnection {
		selector := corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{Name: name},
			Key:                  "value",
		}
		return &coh.HealthMutatorConnection{
			BasicAuth: &coh.HTTPBasicAuth{Username: selector, Password: selector},
		}
	}
	resource := func(name string, desired *coh.HealthMutatorConnection) *coh.Coherence {
		return &coh.Coherence{
			ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace},
			Spec:       coh.CoherenceStatefulSetResourceSpec{HealthMutatorConnection: desired},
		}
	}
	pod := func(name string, owner *coh.Coherence, live *coh.HealthMutatorConnection) *corev1.Pod {
		data, err := json.Marshal(live)
		if err != nil {
			t.Fatal(err)
		}
		return &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:        name,
				Namespace:   owner.Namespace,
				Labels:      owner.Spec.CreatePodSelectorLabels(owner),
				Annotations: map[string]string{coh.HealthMutatorAnnotation: string(data)},
			},
		}
	}

	desired := resource("desired", connection(secretName))
	liveOnly := resource("live-only", nil)
	unrelated := resource("unrelated", nil)
	otherNamespace := resource("other-namespace", connection(secretName))
	otherNamespace.Namespace = "other"

	objects := []client.Object{
		desired,
		liveOnly,
		unrelated,
		otherNamespace,
		pod("desired-0", desired, connection(secretName)),
		pod("live-only-0", liveOnly, connection(secretName)),
		pod("unrelated-0", unrelated, connection("other-secret")),
	}
	kube := &podListCountingClient{Client: fake.NewClientBuilder().WithScheme(scheme).WithObjects(objects...).Build()}
	requests := mapHealthSecretRequests(context.Background(), kube, &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{Name: secretName, Namespace: namespace},
	})

	if kube.podListCalls != 1 {
		t.Fatalf("expected one Pod list, got %d", kube.podListCalls)
	}
	if len(requests) != 2 {
		t.Fatalf("expected two requests, got %v", requests)
	}
	requested := map[string]int{}
	for _, request := range requests {
		requested[request.Name]++
	}
	if requested[desired.Name] != 1 {
		t.Errorf("expected one request for desired reference, got %d", requested[desired.Name])
	}
	if requested[liveOnly.Name] != 1 {
		t.Errorf("expected one request for live Pod reference, got %d", requested[liveOnly.Name])
	}
	if requested[unrelated.Name] != 0 || requested[otherNamespace.Name] != 0 {
		t.Errorf("unexpected requests: %v", requested)
	}
}
