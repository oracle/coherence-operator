/*
 * Copyright (c) 2020, 2025, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package controllers

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-logr/logr"
	coh "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/controllers/errorhandling"
	"github.com/oracle/coherence-operator/controllers/finalizer"
	"github.com/oracle/coherence-operator/controllers/predicates"
	"github.com/oracle/coherence-operator/controllers/reconciler"
	"github.com/oracle/coherence-operator/controllers/resources"
	"github.com/oracle/coherence-operator/controllers/secret"
	"github.com/oracle/coherence-operator/controllers/servicemonitor"
	"github.com/oracle/coherence-operator/controllers/statefulset"
	"github.com/oracle/coherence-operator/controllers/status"
	"github.com/oracle/coherence-operator/pkg/clients"
	"github.com/oracle/coherence-operator/pkg/operator"
	"github.com/oracle/coherence-operator/pkg/utils"
	"github.com/pkg/errors"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	coreV1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

const (
	// The name of this controller
	controllerName = "controllers.Coherence"

	// The error message template to use to indicate a reconcile failure.
	reconcileFailedMessage string = "failed to reconcile Coherence resource '%s' in namespace '%s'\n%s"

	// The error message template to use to indicate a resource creation failure.
	createResourcesFailedMessage string = "create resources for Coherence resource '%s' in namespace '%s' failed\n%s"
)

// CoherenceReconciler reconciles a Coherence resource
type CoherenceReconciler struct {
	client.Client
	reconciler.CommonReconciler
	ClientSet        clients.ClientSet
	Log              logr.Logger
	Scheme           *runtime.Scheme
	reconcilers      []reconciler.SecondaryResourceReconciler
	finalizerManager *finalizer.FinalizerManager
	statusManager    *status.StatusManager
	resourcesManager *resources.OperatorSecretManager
}

// Failure is a simple holder for a named error
type Failure struct {
	Name  string
	Error error
}

// blank assignment to verify that CoherenceReconciler implements reconcile.Reconciler
// There will be a compile-time error here if this breaks
var _ reconcile.Reconciler = &CoherenceReconciler{}

// +kubebuilder:rbac:groups=coherence.oracle.com,resources=coherence;coherencejob;coherence/finalizers;coherencejob/finalizers;coherence/status;coherencejob/status,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=pods;pods/exec;services;endpoints;events;configmaps;secrets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=monitoring.coreos.com,resources=servicemonitors,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;update;patch;delete

// Reconcile performs a full reconciliation for the Coherence resource referred to by the Request.
// The Controller will requeue the Request to be processed again if an error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (in *CoherenceReconciler) Reconcile(ctx context.Context, request ctrl.Request) (ctrl.Result, error) {
	var err error

	log := in.Log.WithValues("namespace", request.Namespace, "name", request.Name)

	log.Info("Reconciling Coherence resource")

	// Attempt to lock the requested resource. If the resource is locked then another
	// request for the same resource is already in progress so requeue this one.
	if ok := in.Lock(request); !ok {
		log.Info("Coherence resource " + request.Namespace + "/" + request.Name + " is already locked, requeue request")
		return reconcile.Result{RequeueAfter: time.Minute}, nil
	}
	// Make sure that the request is unlocked when this method exits
	defer in.Unlock(request)

	// Fetch the Coherence resource instance
	deployment := &coh.Coherence{}
	err = in.GetClient().Get(ctx, types.NamespacedName{Namespace: request.Namespace, Name: request.Name}, deployment)
	if err != nil {
		if apierrors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected.
			// Return and don't requeue
			log.Info("Coherence resource not found. Ignoring request since object must be deleted.")
			return ctrl.Result{}, nil
		}
		// else... error reading the current deployment state from k8s.
		msg := fmt.Sprintf("failed to find Coherence resource, %s", err.Error())
		in.GetEventRecorder().Event(deployment, coreV1.EventTypeWarning, reconciler.EventReasonFailed, msg)
		// returning an error will requeue the event so we will try again
		wrappedErr := errorhandling.NewGetResourceError(request.Name, request.Namespace, err)
		return reconcile.Result{}, wrappedErr
	}

	// Check whether this is a deletion
	deleteTime := deployment.GetDeletionTimestamp()
	if deleteTime != nil {
		// Check whether finalization needs to be run
		if controllerutil.ContainsFinalizer(deployment, coh.CoherenceFinalizer) {
			log.Info("Coherence resource deleted at " + deleteTime.String() + ", running finalizer")
			// Run finalization logic.
			// If the finalization logic fails, don't remove the finalizer so
			// that we can retry during the next reconciliation.
			in.GetEventRecorder().Event(deployment, coreV1.EventTypeNormal, reconciler.EventReasonDeleted, "running finalizers")
			if err := in.finalizerManager.FinalizeDeployment(ctx, deployment, in.MaybeFindStatefulSet); err != nil {
				msg := fmt.Sprintf("failed to finalize Coherence resource, %s", err.Error())
				in.GetEventRecorder().Event(deployment, coreV1.EventTypeWarning, reconciler.EventReasonDeleted, msg)
				log.Error(err, "Failed to remove finalizer")
				return ctrl.Result{RequeueAfter: time.Minute}, nil
			}
			// Remove the finalizer. Once all finalizers have been
			// removed, the object will be deleted.
			err := in.finalizerManager.EnsureFinalizerRemoved(ctx, deployment)
			if err != nil {
				if apierrors.IsNotFound(err) {
					log.Info("Failed to remove the finalizer from the Coherence resource, it looks like it had already been deleted")
					return ctrl.Result{}, nil
				}
				msg := fmt.Sprintf("failed to remove finalizers from Coherence resource, %s", err.Error())
				in.GetEventRecorder().Event(deployment, coreV1.EventTypeWarning, reconciler.EventReasonDeleted, msg)
				wrappedErr := errorhandling.NewOperationError("remove_finalizer", err).
					WithContext("resource", deployment.GetName()).
					WithContext("namespace", deployment.GetNamespace())
				return ctrl.Result{}, wrappedErr
			}
		} else {
			log.Info("Coherence resource deleted at " + deleteTime.String() + ", finalizer already removed")
		}
		// nothing else to do
		return ctrl.Result{}, nil
	}

	// This is an add request or update request

	if deployment.Spec.AllowUnsafeDelete != nil && *deployment.Spec.AllowUnsafeDelete {
		if controllerutil.ContainsFinalizer(deployment, coh.CoherenceFinalizer) {
			err := in.finalizerManager.EnsureFinalizerRemoved(ctx, deployment)
			if err != nil && !apierrors.IsNotFound(err) {
				return ctrl.Result{}, errorhandling.NewOperationError("remove_finalizer", err).
					WithContext("resource", deployment.GetName()).
					WithContext("namespace", deployment.GetNamespace()).
					WithContext("reason", "allow_unsafe_delete")
			}
			log.Info("Removed Finalizer from Coherence resource as AllowUnsafeDelete has been set to true")
		} else {
			log.Info("Finalizer not added to Coherence resource as AllowUnsafeDelete has been set to true")
		}
	} else {
		// Add a finalizer for this CR if required
		if finalizerApplied, err := in.finalizerManager.EnsureFinalizerApplied(ctx, deployment); finalizerApplied || err != nil {
			var msg string
			if err != nil {
				msg = fmt.Sprintf("failed to add finalizers to Coherence resource, %s", err.Error())
				in.GetEventRecorder().Event(deployment, coreV1.EventTypeWarning, reconciler.EventReasonFailed, msg)
			} else {
				in.GetEventRecorder().Event(deployment, coreV1.EventTypeNormal, reconciler.EventReasonUpdated, "added finalizer")
			}
			// we need to requeue as we have updated the Coherence resource
			return ctrl.Result{RequeueAfter: time.Minute}, err
		}
	}

	// ensure that the deployment has an initial status
	if deployment.Status.Phase == "" {
		err := in.statusManager.UpdateCoherenceStatusPhase(ctx, request.NamespacedName, coh.ConditionTypeInitialized)
		if err != nil {
			// failed to set the status
			return reconcile.Result{}, err
		}
	}

	// Check whether the deployment has a replica count specified
	if deployment.Spec.Replicas == nil {
		// No replica count, so we patch the deployment to have the default replicas value.
		// The reason we do this is the kubectl scale command will fail if the replicas
		// field is not set to a non-nil value.
		patch := &coh.Coherence{}
		deployment.DeepCopyInto(patch)
		replicas := deployment.GetReplicas()
		patch.Spec.Replicas = &replicas
		_, err = in.ThreeWayPatch(ctx, deployment.Name, deployment, deployment, patch)
		if err != nil {
			in.GetEventRecorder().Event(deployment, coreV1.EventTypeWarning, reconciler.EventReasonFailed,
				fmt.Sprintf("failed to add default replicas to Coherence resource, %s", err.Error()))
			return reconcile.Result{}, errors.Wrap(err, "failed to add default replicas to Coherence resource")
		}
		msg := "Added default replicas to Coherence resource, re-queuing request"
		log.Info(msg, "Replicas", strconv.Itoa(int(replicas)))
		in.GetEventRecorder().Event(deployment, coreV1.EventTypeNormal, reconciler.EventReasonUpdated, msg)
		return reconcile.Result{}, err
	}

	// ensure that the Operator configuration Secret exists
	if err = in.resourcesManager.EnsureOperatorSecret(ctx, deployment); err != nil {
		err = errorhandling.NewResourceError("ensure", "operator_secret", deployment.GetNamespace(), err)
		return in.HandleErrAndRequeue(ctx, err, deployment, fmt.Sprintf(reconcileFailedMessage, request.Name, request.Namespace, err), in.Log)
	}

	// ensure that the state store exists
	storage, err := utils.NewStorage(request.NamespacedName, in.GetManager(), in.GetPatcher())
	if err != nil {
		err = errorhandling.NewOperationError("obtain_state_store", err).
			WithContext("resource", deployment.GetName()).
			WithContext("namespace", deployment.GetNamespace())
		in.GetEventRecorder().Event(deployment, coreV1.EventTypeWarning, reconciler.EventReasonFailed,
			fmt.Sprintf("failed to obtain state store: %s", err.Error()))
		return in.HandleErrAndRequeue(ctx, err, deployment, fmt.Sprintf(reconcileFailedMessage, request.Name, request.Namespace, err), in.Log)
	}

	// create the result
	result := ctrl.Result{}

	hash := deployment.GetGenerationString()
	storeHash, _ := storage.GetHash()
	var desiredResources coh.Resources

	if hash == storeHash && deployment.IsBeforeOrSameVersion("3.4.3") {
		deployment.UpdateStatusVersion(operator.GetVersion())
		if err = storage.ResetHash(ctx, deployment); err != nil {
			return result, errors.Wrap(err, "error updating storage status hash")
		}
		hashNew := deployment.GetGenerationString()
		if err = in.statusManager.UpdateDeploymentStatusHash(ctx, request.NamespacedName, hashNew); err != nil {
			return result, errors.Wrap(err, "error updating deployment status hash")
		}
		log.Info("Updated pre-3.5.0 Coherence resource status hash", "From", hash, "To", hashNew)
		return result, nil
	}

	desiredResources, err = getDesiredResources(deployment, storage, log)
	if err != nil {
		err = errorhandling.NewOperationError("get_desired_resources", err).
			WithContext("resource", deployment.GetName()).
			WithContext("namespace", deployment.GetNamespace())
		return in.HandleErrAndRequeue(ctx, err, deployment, fmt.Sprintf(createResourcesFailedMessage, request.Name, request.Namespace, err), in.Log)
	}

	log.Info("Reconciling Coherence resource secondary resources", "hash", hash, "store", storeHash)

	// make the deployment the owner of all the secondary resources about to be reconciled
	if err := desiredResources.SetController(deployment, in.GetManager().GetScheme()); err != nil {
		err = errorhandling.NewOperationError("set_controller", err).
			WithContext("resource", deployment.GetName()).
			WithContext("namespace", deployment.GetNamespace())
		return reconcile.Result{}, err
	}

	// set the hash on all the secondary resources to match the deployment's hash
	desiredResources.SetHashLabelAndAnnotations(hash)

	// update the store to have the desired state as the latest state.
	if err = storage.Store(ctx, desiredResources, deployment); err != nil {
		err = errorhandling.NewOperationError("store_state", err).
			WithContext("resource", deployment.GetName()).
			WithContext("namespace", deployment.GetNamespace())
		return reconcile.Result{}, err
	}

	// Ensure the version is present
	deployment.UpdateStatusVersion(operator.GetVersion())

	// process the secondary resources in the order they should be created
	var failures []Failure
	for _, rec := range in.reconcilers {
		r, err := rec.ReconcileAllResourceOfKind(ctx, request, deployment, storage)
		if err != nil {
			failures = append(failures, Failure{Name: rec.GetControllerName(), Error: err})
			result.RequeueAfter = time.Minute
		} else if r.RequeueAfter > 0 && (result.RequeueAfter <= 0 || r.RequeueAfter < result.RequeueAfter) {
			result.RequeueAfter = r.RequeueAfter
		}
	}

	if len(failures) > 0 {
		// one or more reconcilers failed:
		for _, failure := range failures {
			log.Error(failure.Error, "Secondary Reconciler failed", "Reconciler", failure.Name)
		}

		// Create a composite error with context
		err = errorhandling.NewOperationError("reconcile_secondary_resources", nil).
			WithContext("resource", deployment.GetName()).
			WithContext("namespace", deployment.GetNamespace()).
			WithContext("failed_reconcilers", fmt.Sprintf("%d", len(failures)))

		// Add the first failure as the underlying error
		if len(failures) > 0 {
			err.(*errorhandling.OperationError).Err = failures[0].Error
		}

		return reconcile.Result{}, err
	}

	// if replica count is zero update the status to Stopped
	if deployment.GetReplicas() == 0 {
		if err = in.statusManager.UpdateCoherenceStatusPhase(ctx, request.NamespacedName, coh.ConditionTypeStopped); err != nil {
			return result, errorhandling.NewOperationError("update_status", err).
				WithContext("resource", deployment.GetName()).
				WithContext("namespace", deployment.GetNamespace()).
				WithContext("status", string(coh.ConditionTypeStopped))
		}
	}

	// Update the Status with the hash
	if err = in.statusManager.UpdateDeploymentStatusHash(ctx, request.NamespacedName, hash); err != nil {
		return result, errorhandling.NewOperationError("update_status_hash", err).
			WithContext("resource", deployment.GetName()).
			WithContext("namespace", deployment.GetNamespace()).
			WithContext("hash", hash)
	}

	log.Info("Finished reconciling Coherence resource", "RequeueAfter", result.RequeueAfter)
	return result, nil
}

func (in *CoherenceReconciler) SetupWithManager(mgr ctrl.Manager, cs clients.ClientSet) error {
	SetupMonitoringResources(mgr)

	// Create the sub-resource reconcilers IN THE ORDER THAT RESOURCES MUST BE CREATED.
	// This is important to ensure, for example, that a ConfigMap is created before the
	// StatefulSet that uses it.
	reconcilers := []reconciler.SecondaryResourceReconciler{
		reconciler.NewConfigMapReconciler(mgr, cs),
		secret.NewSecretReconciler(mgr, cs),
		reconciler.NewServiceReconciler(mgr, cs),
		servicemonitor.NewServiceMonitorReconciler(mgr, cs),
		statefulset.NewStatefulSetReconciler(mgr, cs),
	}

	in.reconcilers = reconcilers
	in.SetCommonReconciler(controllerName, mgr, cs)
	in.GetPatcher().SetPatchType(types.MergePatchType)

	// Initialize the manager fields
	in.finalizerManager = &finalizer.FinalizerManager{
		Client:        mgr.GetClient(),
		Log:           in.Log.WithName("finalizer"),
		EventRecorder: in.GetEventRecorder(),
		Patcher:       in.GetPatcher(),
	}

	in.statusManager = &status.StatusManager{
		Client:  mgr.GetClient(),
		Log:     in.Log.WithName("status"),
		Patcher: in.GetPatcher(),
	}

	in.resourcesManager = &resources.OperatorSecretManager{
		Client: mgr.GetClient(),
		Log:    in.Log.WithName("resources"),
	}

	template := &coh.Coherence{}

	// Watch for changes to secondary resources
	for _, sub := range reconcilers {
		if err := watchSecondaryResource(mgr, sub, template); err != nil {
			return err
		}
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(template).
		Named("coherence").
		WithOptions(controller.Options{MaxConcurrentReconciles: 1}).
		Complete(in)
}

// GetReconciler returns this reconciler.
func (in *CoherenceReconciler) GetReconciler() reconcile.Reconciler { return in }

// watchSecondaryResource registers the secondary resource reconcilers to watch the resources to be reconciled
func watchSecondaryResource(mgr ctrl.Manager, s reconciler.SecondaryResourceReconciler, owner coh.CoherenceResource) error {
	var err error
	if !s.CanWatch() {
		// this reconciler does not do watches
		return nil
	}

	// Create a new controller
	opts := controller.Options{Reconciler: s.GetReconciler(), MaxConcurrentReconciles: 1}
	c, err := controller.New(s.GetControllerName(), s.GetManager(), opts)
	if err != nil {
		return err
	}

	h := handler.EnqueueRequestForOwner(mgr.GetScheme(), mgr.GetRESTMapper(), owner)
	p := predicates.SecondaryPredicate{}
	src := source.Kind(mgr.GetCache(), s.GetTemplate(), h, p)
	err = c.Watch(src)
	return err
}

// SetupMonitoringResources ensures the Prometheus types are registered with the manager.
func SetupMonitoringResources(mgr ctrl.Manager) {
	gv := schema.GroupVersion{
		Group:   coh.ServiceMonitorGroup,
		Version: coh.ServiceMonitorVersion,
	}
	mgr.GetScheme().AddKnownTypes(gv, &monitoringv1.ServiceMonitor{}, &monitoringv1.ServiceMonitorList{})
}
