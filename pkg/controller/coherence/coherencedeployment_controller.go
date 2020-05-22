/*
 * Copyright (c) 2020, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package coherence

import (
	"fmt"
	monitoringv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/go-test/deep"
	coh "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	"github.com/oracle/coherence-operator/pkg/controller/reconciler"
	"github.com/oracle/coherence-operator/pkg/controller/servicemonitor"
	"github.com/oracle/coherence-operator/pkg/controller/statefulset"
	"github.com/oracle/coherence-operator/pkg/flags"
	"github.com/oracle/coherence-operator/pkg/operator"
	"github.com/oracle/coherence-operator/pkg/utils"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

const (
	controllerName = "controller_coherence"

	reconcileFailedMessage       string = "failed to reconcile Coherence resource '%s' in namespace '%s'\n%s"
	createResourcesFailedMessage string = "create resources for Coherence resource '%s' in namesapce '%s' failed\n%s"
)

// Add creates a new Coherence resource Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, NewReconciler(mgr))
}

// NewReconciler returns a new ReconcileCoherence
func NewReconciler(mgr manager.Manager) *ReconcileCoherence {
	return NewReconcilerWithFlags(mgr, flags.GetOperatorFlags())
}

// NewReconciler creates a new ReconcileCoherence
func NewReconcilerWithFlags(mgr manager.Manager, opFlags *flags.CoherenceOperatorFlags) *ReconcileCoherence {
	r := &ReconcileCoherence{opFlags: opFlags}

	gv := schema.GroupVersion{
		Group:   coh.ServiceMonitorGroup,
		Version: coh.ServiceMonitorVersion,
	}
	mgr.GetScheme().AddKnownTypes(gv, &monitoringv1.ServiceMonitor{}, &monitoringv1.ServiceMonitorList{})

	// Create the sub-resource reconcilers IN THE ORDER THAT RESOURCES MUST BE CREATED.
	// This is important to ensure, for example, that a ConfigMap is created before the
	// StatefulSet that uses it.
	reconcilers := []reconciler.SecondaryResourceReconciler{
		reconciler.NewConfigMapReconciler(mgr),
		reconciler.NewSecretReconciler(mgr),
		reconciler.NewServiceReconciler(mgr),
		servicemonitor.NewServiceMonitorReconciler(mgr),
		statefulset.NewStatefulSetReconciler(mgr),
	}

	r.reconcilers = reconcilers

	r.SetCommonReconciler(controllerName, mgr)
	return r
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r *ReconcileCoherence) error {
	// Create a new controller
	c, err := controller.New("coherence-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Coherence
	err = c.Watch(&source.Kind{Type: &coh.Coherence{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resources
	for _, sub := range r.reconcilers {
		if err = watchSecondaryResource(sub); err != nil {
			return err
		}
	}

	return nil
}

// Watch the resources to be reconciled
func watchSecondaryResource(r reconciler.SecondaryResourceReconciler) error {
	if !r.CanWatch() {
		// this reconciler does not do watches
		return nil
	}

	// Create a new controller
	c, err := controller.New(r.GetControllerName(), r.GetManager(), controller.Options{Reconciler: r.GetReconciler()})
	if err != nil {
		return err
	}

	src := &source.Kind{Type: r.GetTemplate()}
	h := &handler.EnqueueRequestForOwner{IsController: true, OwnerType: &coh.Coherence{}}
	if err := c.Watch(src, h); err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileCoherence implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileCoherence{}

// ReconcileCoherence reconciles a Coherence resource object
type ReconcileCoherence struct {
	reconciler.CommonReconciler
	opFlags     *flags.CoherenceOperatorFlags
	reconcilers []reconciler.SecondaryResourceReconciler
}

func (in *ReconcileCoherence) SetPatchType(pt types.PatchType) {
	for _, r := range in.reconcilers {
		r.SetPatchType(pt)
	}
	in.CommonReconciler.SetPatchType(pt)
}

// Reconcile reads that state of the cluster for a Coherence resource object and makes changes based on the state read
// and what is in the Coherence.Spec
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (in *ReconcileCoherence) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	logger := in.GetLog().WithValues("Namespace", request.Namespace, "Name", request.Name)
	logger.Info("Reconciling Coherence resource")

	// Attempt to lock the requested resource. If the resource is locked then another
	// request for the same resource is already in progress so requeue this one.
	if ok := in.Lock(request); !ok {
		logger.Info("Coherence resource " + request.Namespace + "/" + request.Name + " is already locked, re-queuing request")
		return reconcile.Result{Requeue: true, RequeueAfter: 0}, nil
	}
	// Make sure that the request is unlocked when this method exits
	defer in.Unlock(request)

	// Fetch the Coherence resource instance
	deployment, found, err := in.MaybeFindDeployment(request.Namespace, request.Name)
	if err != nil {
		// Error reading the current deployment state from k8s.
		return reconcile.Result{}, err
	}

	if !found || deployment.GetDeletionTimestamp() != nil {
		logger.Info("Coherence resource deleted")
		return reconcile.Result{}, nil
	}

	// ensure that the deployment has an initial status
	if deployment.Status.Phase == "" {
		err := in.UpdateDeploymentStatusPhase(request.NamespacedName, coh.ConditionTypeInitialized)
		if err != nil {
			return reconcile.Result{}, err
		}
	}

	// Check whether the deployment has a replica count specified
	// Ideally we'd do this with a validating/defaulting web-hook but maybe in a later version.
	if deployment.Spec.Replicas == nil {
		// No replica count so we patch the deployment to have the default replicas value.
		// The reason we do this is because the kubectl scale command will fail if the replicas
		// field is not set to a non-nil value.
		patch := &coh.Coherence{}
		deployment.DeepCopyInto(patch)
		replicas := deployment.GetReplicas()
		patch.Spec.Replicas = &replicas
		_, err = in.ThreeWayPatch(deployment.Name, deployment, deployment, patch)
		if err != nil {
			return reconcile.Result{}, err
		}
	}

	// ensure that the Operator configuration Secret exists
	if err = operator.EnsureOperatorSecret(request.Namespace, in.GetClient(), logger); err != nil {
		err = errors.Wrap(err, "ensuring Operator configuration secret")
		return in.HandleErrAndRequeue(err, nil, fmt.Sprintf(reconcileFailedMessage, request.Name, request.Namespace, err), logger)
	}

	// ensure that the state store exists
	storage, err := utils.NewStorage(request.NamespacedName, in.GetManager())
	if err != nil {
		err = errors.Wrap(err, "obtaining desired state store")
		return in.HandleErrAndRequeue(err, nil, fmt.Sprintf(reconcileFailedMessage, request.Name, request.Namespace, err), logger)
	}

	// obtain the original resources from the state store
	originalResources := storage.GetLatest()
	// Create the desired resources the deployment
	desiredResources, err := deployment.Spec.CreateKubernetesResources(deployment, flags.GetOperatorFlags())
	if err != nil {
		return in.HandleErrAndRequeue(err, nil, fmt.Sprintf(createResourcesFailedMessage, request.Name, request.Namespace, err), logger)
	}

	// compare the original with the desired to determine whether there is anything to update
	var isDiff bool
	switch {
	case originalResources.Items == nil && desiredResources.Items == nil:
		isDiff = false
	case originalResources.Items != nil && desiredResources.Items == nil:
		isDiff = true
	case originalResources.Items == nil && desiredResources.Items != nil:
		isDiff = true
	default:
		diff := deep.Equal(originalResources.Items, desiredResources.Items)
		isDiff = len(diff) > 0
	}

	// create the result
	result := reconcile.Result{Requeue: false}

	if isDiff {
		logger.Info("Reconciling Coherence resource secondary resources")
		// make the deployment the owner of all of the secondary resources about to be reconciled
		if deployment != nil {
			if err := desiredResources.SetController(deployment, in.GetManager().GetScheme()); err != nil {
				return reconcile.Result{}, err
			}
		}

		// update the store to have the desired state as the latest state.
		if err = storage.Store(desiredResources, deployment); err != nil {
			err = errors.Wrap(err, "storing latest state in state store")
			return reconcile.Result{}, err
		}

		// process the secondary resources in the order they should be created
		for _, rec := range in.reconcilers {
			r, err := rec.ReconcileResources(request, deployment, storage)
			if err != nil {
				return reconcile.Result{}, err
			}
			result.Requeue = result.Requeue || r.Requeue
		}
	} else {
		// original and desired are identical so there is nothing else to do
		logger.Info("Reconciled secondary resources for deployment, nothing to update")
	}

	// if replica count is zero update the status to Stopped
	if deployment.GetReplicas() == 0 {
		if err = in.UpdateDeploymentStatusPhase(request.NamespacedName, coh.ConditionTypeStopped); err != nil {
			err = errors.Wrap(err, "error updating deployment status")
		}
	}

	logger.Info(fmt.Sprintf("Finished reconciling Coherence resource. Result='%v'", result))
	return result, err
}

func (in *ReconcileCoherence) GetReconciler() reconcile.Reconciler { return in }
