/*
 * Copyright (c) 2026, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */
package controllers_test

import (
	"context"
	"strings"
	"testing"

	"github.com/go-logr/logr"
	coh "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/controllers"
	"github.com/oracle/coherence-operator/pkg/clients"
	"github.com/oracle/coherence-operator/pkg/fakes"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	kubescheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
)

func TestInvalidHealthMutatorConfigurationReportsFailure(t *testing.T) {
	if err := coh.AddToScheme(kubescheme.Scheme); err != nil {
		t.Fatal(err)
	}

	selector := corev1.SecretKeySelector{
		LocalObjectReference: corev1.LocalObjectReference{Name: "health-secret"},
		Key:                  "value",
	}
	validConnection := &coh.HealthMutatorConnection{
		Secure:    true,
		BasicAuth: &coh.HTTPBasicAuth{Username: selector, Password: selector},
		TLS:       &coh.HTTPClientTLS{CASecret: selector, ServerName: "health.example"},
		ServerTLS: &coh.SSLSpec{
			Enabled:              ptr.To(true),
			KeyStore:             ptr.To("server.p12"),
			KeyStorePasswordFile: ptr.To("store-password"),
			KeyPasswordFile:      ptr.To("key-password"),
		},
	}

	tests := []struct {
		name        string
		connection  *coh.HealthMutatorConnection
		probe       *coh.Probe
		application *coh.ApplicationSpec
		message     string
	}{
		{
			name:       "incomplete secure connection",
			connection: &coh.HealthMutatorConnection{Secure: true},
			message:    "secure health mutators require credentials",
		},
		{
			name:       "managed connection with custom suspend probe",
			connection: validConnection,
			probe: &coh.Probe{HTTP: &coh.HTTPAction{
				Method:       "PUT",
				HTTPEndpoint: coh.HTTPEndpoint{Path: "/suspend", Port: intstr.FromInt(8089)},
			}},
			message: "managed secure health mode excludes custom suspend probes",
		},
		{
			name:       "managed connection without JVM argument delivery",
			connection: validConnection,
			application: &coh.ApplicationSpec{
				UseImageEntryPoint: ptr.To(true),
				UseJdkJavaOptions:  ptr.To(false),
			},
			message: "requires application.useJdkJavaOptions=true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource := &coh.Coherence{
				ObjectMeta: metav1.ObjectMeta{Name: "test", Namespace: "test"},
			}
			resource.Spec.HealthMutatorConnection = tt.connection
			resource.Spec.SuspendProbe = tt.probe
			resource.Spec.Application = tt.application

			kube := fakes.NewFakeClient(resource)
			events := fakes.NewFakeEventRecorder(10)
			manager := &fakes.FakeManager{Scheme: kubescheme.Scheme, Client: kube, Events: events}
			reconciler := &controllers.CoherenceReconciler{Client: kube, Log: logr.Discard()}
			reconciler.SetCommonReconciler("test.Coherence", manager, clients.ClientSet{})

			result, err := reconciler.Reconcile(context.Background(), ctrl.Request{NamespacedName: types.NamespacedName{Namespace: resource.Namespace, Name: resource.Name}})
			if err != nil {
				t.Fatalf("reconcile returned an error: %v", err)
			}
			if result != (ctrl.Result{}) {
				t.Fatalf("invalid permanent configuration was requeued: %+v", result)
			}

			actual := &coh.Coherence{}
			if err := kube.Get(context.Background(), types.NamespacedName{Namespace: resource.Namespace, Name: resource.Name}, actual); err != nil {
				t.Fatal(err)
			}
			if actual.Status.Phase != coh.ConditionTypeFailed {
				t.Fatalf("phase = %q, want %q", actual.Status.Phase, coh.ConditionTypeFailed)
			}
			condition := actual.Status.Conditions.GetCondition(coh.ConditionTypeFailed)
			if condition == nil || condition.Status != corev1.ConditionTrue || condition.Reason != "InvalidHealthMutatorConfiguration" || !strings.Contains(condition.Message, tt.message) {
				t.Fatalf("unexpected Failed condition: %+v", condition)
			}

			event, found := manager.NextEvent()
			if !found {
				t.Fatal("missing invalid configuration event")
			}
			if event.Type != corev1.EventTypeWarning || event.Reason != "InvalidHealthMutatorConfiguration" || !strings.Contains(event.Message, tt.message) {
				t.Fatalf("unexpected event: %+v", event)
			}
			if _, found := manager.NextEvent(); found {
				t.Fatal("unexpected additional event")
			}

			result, err = reconciler.Reconcile(context.Background(), ctrl.Request{NamespacedName: types.NamespacedName{Namespace: resource.Namespace, Name: resource.Name}})
			if err != nil || result != (ctrl.Result{}) {
				t.Fatalf("unchanged invalid configuration did not remain quiet: result=%+v err=%v", result, err)
			}
			if _, found := manager.NextEvent(); found {
				t.Fatal("unchanged invalid configuration emitted another event")
			}
		})
	}
}
