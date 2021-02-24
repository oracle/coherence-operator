/*
 * Copyright (c) 2020, 2021, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package servicemonitor

import (
	"context"
	"fmt"
	monitoringv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	monclientv1 "github.com/coreos/prometheus-operator/pkg/client/versioned/typed/monitoring/v1"
	"github.com/go-logr/logr"
	coh "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/controllers/reconciler"
	"github.com/oracle/coherence-operator/pkg/utils"
	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/discovery"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	// The name of this controller. This is used in events, log messages, etc.
	controllerName = "controllers.ServiceMonitor"
)

// blank assignment to verify that ReconcileServiceMonitor implements reconcile.Reconciler.
// If the reconcile.Reconciler API was to change then we'd get a compile error here.
var _ reconcile.Reconciler = &ReconcileServiceMonitor{}

// NewServiceReconciler returns a new Service reconciler.
func NewServiceMonitorReconciler(mgr manager.Manager) reconciler.SecondaryResourceReconciler {
	r := &ReconcileServiceMonitor{
		ReconcileSecondaryResource: reconciler.ReconcileSecondaryResource{
			Kind:      coh.ResourceTypeServiceMonitor,
			Template:  &monitoringv1.ServiceMonitor{},
			SkipWatch: true,
		},
		monClient: monclientv1.NewForConfigOrDie(mgr.GetConfig()),
	}

	r.SetCommonReconciler(controllerName, mgr)
	return r
}

type ReconcileServiceMonitor struct {
	reconciler.ReconcileSecondaryResource
	monClient *monclientv1.MonitoringV1Client
}

func (in *ReconcileServiceMonitor) GetReconciler() reconcile.Reconciler { return in }

// Reconcile reads that state of the ServiceMonitors for a deployment and makes changes based on the
// state read and the desired state based on the parent Coherence resource.
func (in *ReconcileServiceMonitor) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	// Attempt to lock the requested resource. If the resource is locked then another
	// request for the same resource is already in progress so requeue this one.
	if ok := in.Lock(request); !ok {
		return reconcile.Result{Requeue: true, RequeueAfter: 0}, nil
	}
	// Make sure that the request is unlocked when this method exits
	defer in.Unlock(request)

	// Obtain the parent Coherence resource
	deployment, err := in.FindDeployment(request)
	if err != nil {
		return reconcile.Result{}, err
	}

	storage, err := utils.NewStorage(request.NamespacedName, in.GetManager())
	if err != nil {
		return reconcile.Result{}, err
	}

	return in.ReconcileResources(request, deployment, storage)
}

// ReconcileResources reconciles the state of the desired ServiceMonitors for the reconciler
func (in *ReconcileServiceMonitor) ReconcileResources(req reconcile.Request, d *coh.Coherence, store utils.Storage) (reconcile.Result, error) {
	logger := in.GetLog().WithValues("Namespace", req.Namespace, "Name", req.Name)

	var err error
	diff := store.GetLatest().DiffForKind(in.Kind, store.GetPrevious())
	for _, res := range diff {
		if res.IsDelete() {
			if err = in.monClient.ServiceMonitors(req.Namespace).Delete(context.TODO(), res.Name, metav1.DeleteOptions{}); err != nil {
				logger.Info(fmt.Sprintf("Finished reconciling all %s for d with error: %s", in.Kind, err.Error()))
				return reconcile.Result{}, err
			}
		} else {
			if err = in.ReconcileResource(req.Namespace, res.Name, d, store); err != nil {
				logger.Info(fmt.Sprintf("Finished reconciling all %s for d with error: %s", in.Kind, err.Error()))
				return reconcile.Result{}, err
			}
		}
	}
	return reconcile.Result{}, nil
}

func (in *ReconcileServiceMonitor) ReconcileResource(namespace, name string, deployment *coh.Coherence, storage utils.Storage) error {
	logger := in.GetLog().WithValues("Namespace", namespace, "Name", name)

	// See whether it is even possible to handle Prometheus ServiceMonitor resources
	if !in.hasServiceMonitor() {
		logger.Info(fmt.Sprintf("Cannot reconcile ServiceMonitor %s as the ServiceMonitor CR is not installed", name))
		return nil
	}

	// Fetch the sm's current state
	var exists bool
	sm, err := in.monClient.ServiceMonitors(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	switch {
	case err != nil && apierrors.IsNotFound(err):
		// we can ignore not found errors
		err = nil
		exists = false
	case err != nil:
		// Error reading the object - requeue the request.
		// We can't call the error handler as we do not even have a deployment.
		// We log the error and do not requeue the request.
		return errors.Wrapf(err, "getting ServiceMonitor %s/%s", namespace, name)
	default:
		// sm was found
		exists = true
	}

	if exists && sm.GetDeletionTimestamp() != nil {
		// The Service exists but is being deleted
		exists = false
	}

	switch {
	case deployment == nil || deployment.GetReplicas() == 0:
		if exists {
			// The deployment does not exist (or is scaled down to zero) but the sm still does,
			// ensure that the sm is deleted.
			// This should not actually be required as everything is owned by the deployment
			// and there should be a cascaded delete by k8s so it's belt and braces.
			err = in.monClient.ServiceMonitors(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
		}
	case !exists:
		// ServiceMonitor does not exist but deployment does so create it
		err = in.CreateServiceMonitor(namespace, name, storage, logger)
		if err != nil {
			err = errors.Wrapf(err, "Failed to create ServiceMonitor %s/%s", namespace, name)
		}
	default:
		// Both the sm and deployment exists so this is maybe an update
		err = in.Update(name, sm, storage, logger)
	}
	return err
}

// Create the ServiceMonitor
func (in *ReconcileServiceMonitor) CreateServiceMonitor(namespace, name string, storage utils.Storage, logger logr.Logger) error {
	logger.Info(fmt.Sprintf("Creating ServiceMonitor/%s for deployment", name))
	// Get the ServiceMonitor desired state
	resource, found := storage.GetLatest().GetResource(in.Kind, name)
	if !found {
		return fmt.Errorf("cannot create ServiceMonitor %s/%s for deployment as latest state not present in store", namespace, name)
	}
	// create the ServiceMonitor
	sm := resource.Spec.(*monitoringv1.ServiceMonitor)
	_, err := in.monClient.ServiceMonitors(namespace).Create(context.TODO(), sm, metav1.CreateOptions{})
	if err != nil {
		logger.Info(fmt.Sprintf("Failed creating ServiceMonitor %s/%s - %s", namespace, name, err.Error()))
		return errors.Wrapf(err, "failed to create ServiceMonitor %s/%s", namespace, name)
	}
	return nil
}

// Maybe update the ServiceMonitor if the current state differs from the desired state.
func (in *ReconcileServiceMonitor) UpdateServiceMonitor(namespace, name string, current runtime.Object, storage utils.Storage) error {
	original, _ := storage.GetPrevious().GetResource(in.Kind, name)
	desired, found := storage.GetLatest().GetResource(in.Kind, name)
	if !found {
		return fmt.Errorf("cannot update ServiceMonitor %s/%s as latest state not present in store", namespace, name)
	}

	// fix the CreationTimestamp so that it is not in the patch
	desired.Spec.(metav1.Object).SetCreationTimestamp(current.(metav1.Object).GetCreationTimestamp())
	// create the patch
	data, err := in.CreateThreeWayPatchData(original.Spec, desired.Spec, current)
	if err != nil {
		return errors.Wrapf(err, "failed to create patch for ServiceMonitor %s/%s", namespace, name)
	}

	// check whether the patch counts as no-patch
	if string(data) == "{}" {
		// empty patch
		return nil
	}

	in.GetLog().WithValues().Info(fmt.Sprintf("Patching ServiceMonitor %s/%s", namespace, name))
	_, err = in.monClient.ServiceMonitors(namespace).Patch(context.TODO(), name, in.GetPatchType(), data, metav1.PatchOptions{})
	if err != nil {
		return errors.Wrapf(err, "cannot patch ServiceMonitor %s/%s", namespace, name)
	}

	return nil
}

// hasServiceMonitor checks if the Prometheus ServiceMonitor CRD is registered in the cluster.
func (in *ReconcileServiceMonitor) hasServiceMonitor() bool {
	dc := discovery.NewDiscoveryClientForConfigOrDie(in.GetManager().GetConfig())
	apiVersion := coh.ServiceMonitorGroupVersion
	kind := coh.ServiceMonitorKind
	ok, err := ResourceExists(dc, apiVersion, kind)
	if err != nil {
		in.GetLog().Error(err, "error checking for Prometheus ServiceMonitor CRD")
		return false
	}
	return ok
}

// ResourceExists returns true if the given resource kind exists
// in the given api groupversion
func ResourceExists(dc discovery.DiscoveryInterface, apiGroupVersion, kind string) (bool, error) {

	_, apiLists, err := dc.ServerGroupsAndResources()
	if err != nil {
		return false, err
	}
	for _, apiList := range apiLists {
		if apiList.GroupVersion == apiGroupVersion {
			for _, r := range apiList.APIResources {
				if r.Kind == kind {
					return true, nil
				}
			}
		}
	}
	return false, nil
}
