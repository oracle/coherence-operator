/*
 * Copyright (c) 2020, 2021, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package reconciler

import (
	"context"
	coh "github.com/oracle/coherence-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// blank assignment to verify that SimpleReconciler implements reconcile.Reconciler.
// If the reconcile.Reconciler API was to change then we'd get a compile error here.
var _ reconcile.Reconciler = &SimpleReconciler{}

func NewConfigMapReconciler(mgr manager.Manager) SecondaryResourceReconciler {
	return NewSimpleReconciler(mgr, "controllers.ConfigMap", coh.ResourceTypeConfigMap, &corev1.ConfigMap{})
}

func NewServiceReconciler(mgr manager.Manager) SecondaryResourceReconciler {
	return NewSimpleReconciler(mgr, "controllers.Service", coh.ResourceTypeService, &corev1.Service{})
}

// NewSimpleReconciler returns a new SimpleReconciler.
func NewSimpleReconciler(mgr manager.Manager, name string, kind coh.ResourceType, template client.Object) SecondaryResourceReconciler {
	r := &SimpleReconciler{
		ReconcileSecondaryResource: ReconcileSecondaryResource{
			Kind:     kind,
			Template: template,
		},
	}

	r.SetCommonReconciler(name, mgr)
	return r
}

type SimpleReconciler struct {
	ReconcileSecondaryResource
}

func (in *SimpleReconciler) GetReconciler() reconcile.Reconciler { return in }

// Reconcile reads that state of the secondary resource for a deployment and makes changes based on the
// state read and the desired state based on the parent Coherence resource.
func (in *SimpleReconciler) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
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

	err := in.ReconcileSingleResource(ctx, request.Namespace, request.Name, nil, nil, true, logger)
	logger.Info("Completed reconcile")
	return reconcile.Result{}, err
}
