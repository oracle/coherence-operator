/*
 * Copyright (c) 2020, 2024, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package secret

import (
	"context"
	coh "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/controllers/reconciler"
	"github.com/oracle/coherence-operator/pkg/clients"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	// The name of this controller. This is used in events, log messages, etc.
	controllerName = "controllers.Secret"
)

// blank assignment to verify that ReconcileServiceMonitor implements reconcile.Reconciler.
// If the reconcile.Reconciler API was to change then we'd get a compile error here.
var _ reconcile.Reconciler = &ReconcileSecret{}

// NewSecretReconciler returns a new Secret reconciler.
func NewSecretReconciler(mgr manager.Manager, cs clients.ClientSet) reconciler.SecondaryResourceReconciler {
	r := &ReconcileSecret{
		SimpleReconciler: reconciler.SimpleReconciler{
			ReconcileSecondaryResource: reconciler.ReconcileSecondaryResource{
				Kind:      coh.ResourceTypeSecret,
				Template:  &corev1.Secret{},
				SkipWatch: true,
			},
		},
	}

	r.SetCommonReconciler(controllerName, mgr, cs)
	return r
}

// ReconcileSecret is a reconciler for Secrets.
type ReconcileSecret struct {
	reconciler.SimpleReconciler
}

// Reconcile reads that state of the secondary resource for a deployment and makes changes based on the
// state read and the desired state based on the parent Coherence resource.
func (in *ReconcileSecret) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	// We need to see whether this is a storage Secret and if so we do not reconcile it
	s, found, err := in.FindResource(ctx, request.Namespace, request.Name)
	if err == nil && found {
		if _, ok := s.GetLabels()[coh.LabelCoherenceStore]; ok {
			// this is the storage secret so we can skip it
			return reconcile.Result{}, nil
		}
	}

	logger := in.GetLog().WithValues("Namespace", request.Namespace, "Name", request.Name, "Kind", in.Kind.Name())
	logger.Info("Starting reconcile")

	// Attempt to lock the requested resource. If the resource is locked then another
	// request for the same resource is already in progress so requeue this one.
	if ok := in.Lock(request); !ok {
		logger.Info("Completed reconcile. Already locked, re-queuing")
		return reconcile.Result{Requeue: true, RequeueAfter: 0}, nil
	}
	// Make sure that the request is unlocked when this method exits
	defer in.Unlock(request)

	err = in.ReconcileSingleResource(ctx, request.Namespace, request.Name, nil, nil, logger)
	logger.Info("Completed reconcile")
	return reconcile.Result{}, err
}
