/*
 * Copyright (c) 2020, 2025, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package controllers

import (
	"context"
	"fmt"
	"github.com/go-logr/logr"
	coh "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/controllers/predicates"
	"github.com/oracle/coherence-operator/controllers/reconciler"
	"github.com/oracle/coherence-operator/controllers/secret"
	"github.com/oracle/coherence-operator/controllers/servicemonitor"
	"github.com/oracle/coherence-operator/controllers/statefulset"
	"github.com/oracle/coherence-operator/pkg/clients"
	"github.com/oracle/coherence-operator/pkg/events"
	"github.com/oracle/coherence-operator/pkg/operator"
	"github.com/oracle/coherence-operator/pkg/probe"
	"github.com/oracle/coherence-operator/pkg/rest"
	"github.com/oracle/coherence-operator/pkg/utils"
	"github.com/pkg/errors"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/spf13/viper"
	coreV1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	"strconv"
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
	ClientSet   clients.ClientSet
	Log         logr.Logger
	Scheme      *runtime.Scheme
	reconcilers []reconciler.SecondaryResourceReconciler
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
// +kubebuilder:rbac:groups=cert-manager.io,resources=certificates;issuers,verbs=get;list;watch;create;update;patch;delete
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
		return reconcile.Result{Requeue: true, RequeueAfter: 0}, nil
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
		// returning an error will requeue the event so we will try again ToDo: probably need som backoff here
		return reconcile.Result{}, errors.Wrap(err, "getting Coherence resource")
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
			if err := in.finalizeDeployment(ctx, deployment); err != nil {
				msg := fmt.Sprintf("failed to finalize Coherence resource, %s", err.Error())
				in.GetEventRecorder().Event(deployment, coreV1.EventTypeWarning, reconciler.EventReasonDeleted, msg)
				log.Error(err, "Failed to remove finalizer")
				return ctrl.Result{Requeue: true}, nil
			}
			// Remove the finalizer. Once all finalizers have been
			// removed, the object will be deleted.
			controllerutil.RemoveFinalizer(deployment, coh.CoherenceFinalizer)
			err := in.GetClient().Update(ctx, deployment)
			if err != nil {
				if apierrors.IsNotFound(err) {
					log.Info("Failed to remove the finalizer from the Coherence resource, it looks like it had already been deleted")
					return ctrl.Result{}, nil
				}
				msg := fmt.Sprintf("failed to remove finalizers from Coherence resource, %s", err.Error())
				in.GetEventRecorder().Event(deployment, coreV1.EventTypeWarning, reconciler.EventReasonDeleted, msg)
				return ctrl.Result{}, errors.Wrap(err, "trying to remove finalizer from Coherence resource")
			}
		} else {
			log.Info("Coherence resource deleted at " + deleteTime.String() + ", finalizer already removed")
		}
		// nothing else to do
		return ctrl.Result{}, nil
	}

	// The request is an add or update

	// Ensure the hash label is present (it should have been added by the web-hook, so this should be a no-op).
	// The hash may not have been added if the Coherence resource was added/modified when the Operator was uninstalled.
	if hashApplied, err := in.ensureHashApplied(ctx, deployment); hashApplied || err != nil {
		return ctrl.Result{Requeue: true}, err
	}

	if deployment.Spec.AllowUnsafeDelete != nil && *deployment.Spec.AllowUnsafeDelete {
		if controllerutil.ContainsFinalizer(deployment, coh.CoherenceFinalizer) {
			err := in.ensureFinalizerRemoved(ctx, deployment)
			if err != nil && !apierrors.IsNotFound(err) {
				return ctrl.Result{Requeue: true}, errors.Wrap(err, "failed to remove finalizer")
			}
			log.Info("Removed Finalizer from Coherence resource as AllowUnsafeDelete has been set to true")
		} else {
			log.Info("Finalizer not added to Coherence resource as AllowUnsafeDelete has been set to true")
		}
	} else {
		// Add finalizer for this CR if required (it should have been added by the web-hook but may not have been if the
		// Coherence resource was added when the Operator was uninstalled)
		if finalizerApplied, err := in.ensureFinalizerApplied(ctx, deployment); finalizerApplied || err != nil {
			var msg string
			if err != nil {
				msg = fmt.Sprintf("failed to add finalizers to Coherence resource, %s", err.Error())
				in.GetEventRecorder().Event(deployment, coreV1.EventTypeWarning, reconciler.EventReasonFailed, msg)
			} else {
				in.GetEventRecorder().Event(deployment, coreV1.EventTypeNormal, reconciler.EventReasonUpdated, "added finalizer")
			}
			// we need to requeue as we have updated the Coherence resource
			return ctrl.Result{Requeue: true}, err
		}
	}

	// ensure that the deployment has an initial status
	if deployment.Status.Phase == "" {
		err := in.UpdateCoherenceStatusPhase(ctx, request.NamespacedName, coh.ConditionTypeInitialized)
		if err != nil {
			// failed to set the status
			return reconcile.Result{}, err
		}
	}

	// Check whether the deployment has a replica count specified
	// Ideally we'd do this with a validating/defaulting web-hook but maybe in a later version.
	if deployment.Spec.Replicas == nil {
		// No replica count, so we patch the deployment to have the default replicas value.
		// The reason we do this, is because the kubectl scale command will fail if the replicas
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
	if err = in.ensureOperatorSecret(ctx, deployment, in.GetClient(), in.Log); err != nil {
		err = errors.Wrap(err, "ensuring Operator configuration secret")
		return in.HandleErrAndRequeue(ctx, err, nil, fmt.Sprintf(reconcileFailedMessage, request.Name, request.Namespace, err), in.Log)
	}

	// ensure that the state store exists
	storage, err := utils.NewStorage(request.NamespacedName, in.GetManager())
	if err != nil {
		err = errors.Wrap(err, "obtaining desired state store")
		in.GetEventRecorder().Event(deployment, coreV1.EventTypeWarning, reconciler.EventReasonFailed,
			fmt.Sprintf("failed to obtain state store: %s", err.Error()))
		return in.HandleErrAndRequeue(ctx, err, nil, fmt.Sprintf(reconcileFailedMessage, request.Name, request.Namespace, err), in.Log)
	}

	hash := deployment.GetLabels()[coh.LabelCoherenceHash]
	storeHash, _ := storage.GetHash()
	var desiredResources coh.Resources

	desiredResources, err = getDesiredResources(deployment, storage, log)
	if err != nil {
		return in.HandleErrAndRequeue(ctx, err, nil, fmt.Sprintf(createResourcesFailedMessage, request.Name, request.Namespace, err), in.Log)
	}

	// create the result
	result := ctrl.Result{Requeue: false}

	log.Info("Reconciling Coherence resource secondary resources", "hash", hash, "store", storeHash)

	// make the deployment the owner of all the secondary resources about to be reconciled
	if err := desiredResources.SetController(deployment, in.GetManager().GetScheme()); err != nil {
		return reconcile.Result{}, err
	}

	// set the hash on all the secondary resources to match the deployment's hash
	desiredResources.SetHashLabels(hash)

	// update the store to have the desired state as the latest state.
	if err = storage.Store(desiredResources, deployment); err != nil {
		err = errors.Wrap(err, "storing latest state in state store")
		return reconcile.Result{}, err
	}

	// Ensure the version annotation is present (it should have been added by the web-hook, so this should be a no-op).
	// The hash may not have been added if the Coherence resource was added/modified when the Operator was uninstalled.
	applied, err := in.ensureVersionAnnotationApplied(ctx, deployment)
	if err != nil {
		// We failed to update the Coherence resource
		in.GetEventRecorder().Eventf(deployment, coreV1.EventTypeWarning, reconciler.EventReasonFailed,
			fmt.Sprintf("failed to ensure version annotation applied: %s", err.Error()))
		return ctrl.Result{}, err
	}
	if applied {
		// We updated the Coherence resource, so requeue the event
		in.GetEventRecorder().Eventf(deployment, coreV1.EventTypeNormal, reconciler.EventReasonUpdated,
			fmt.Sprintf("applied version annotation"))
		return ctrl.Result{Requeue: true}, nil
	}

	// process the secondary resources in the order they should be created
	var failures []Failure
	for _, rec := range in.reconcilers {
		r, err := rec.ReconcileAllResourceOfKind(ctx, request, deployment, storage)
		if err != nil {
			failures = append(failures, Failure{Name: rec.GetControllerName(), Error: err})
		}
		result.Requeue = result.Requeue || r.Requeue
	}

	if len(failures) > 0 {
		// one or more reconcilers failed:
		for _, failure := range failures {
			log.Error(failure.Error, "Secondary Reconciler failed", "Reconciler", failure.Name)
		}
		return reconcile.Result{}, fmt.Errorf("one or more secondary resource reconcilers failed to reconcile")
	}

	// if replica count is zero update the status to Stopped
	if deployment.GetReplicas() == 0 {
		if err = in.UpdateCoherenceStatusPhase(ctx, request.NamespacedName, coh.ConditionTypeStopped); err != nil {
			return result, errors.Wrap(err, "error updating deployment status")
		}
	}

	// Update the Status with the hash
	if err = in.UpdateDeploymentStatusHash(ctx, request.NamespacedName, hash); err != nil {
		return result, errors.Wrap(err, "error updating deployment status hash")
	}

	log.Info("Finished reconciling Coherence resource", "requeue", result.Requeue, "after", result.RequeueAfter.String())
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
	in.SetPatchType(types.MergePatchType)

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

// ensureHashApplied ensures that the hash label is present in the Coherence resource, patching it if required
// Returns true if the hash label was applied to the Coherence resource, or false if the label was already present
func (in *CoherenceReconciler) ensureHashApplied(ctx context.Context, c *coh.Coherence) (bool, error) {
	currentHash := ""
	labels := c.GetLabels()
	if len(labels) > 0 {
		currentHash = labels[coh.LabelCoherenceHash]
	}

	// Re-fetch the Coherence resource to ensure we have the most recent copy
	latest := c.DeepCopy()
	hash, _ := coh.EnsureCoherenceHashLabel(latest)

	if currentHash != hash {
		if c.IsBeforeVersion("3.4.2") {
			// Before 3.4.2 there was a bug calculating the hash in the defaulting web-hook
			// This would cause the hashes to be different here, when in fact they should not be
			// If the Coherence resource being processes has no version annotation, or a version
			// prior to 3.4.2 then we return as if the hashes matched
			if labels == nil {
				labels = make(map[string]string)
			}
			labels[coh.LabelCoherenceHash] = hash
			c.SetLabels(labels)
			return false, nil
		}
		callback := func() {
			in.Log.Info(fmt.Sprintf("Applied %s label", coh.LabelCoherenceHash), "newHash", hash, "currentHash", currentHash)
		}

		applied, err := in.ThreeWayPatchWithCallback(ctx, c.Name, c, c, latest, callback)
		if err != nil {
			return false, errors.Wrapf(err, "failed to update Coherence resource %s/%s with hash", c.Namespace, c.Name)
		}
		return applied, nil
	}
	return false, nil
}

// ensureVersionAnnotationApplied ensures that the version annotation is present in the Coherence resource, patching it if required
func (in *CoherenceReconciler) ensureVersionAnnotationApplied(ctx context.Context, c *coh.Coherence) (bool, error) {
	currentVersion, _ := c.GetVersionAnnotation()
	operatorVersion := operator.GetVersion()

	if currentVersion != operatorVersion {
		// make a copy of the Coherence resource to use in the three-way patch
		latest := c.DeepCopy()
		latest.AddAnnotation(coh.AnnotationOperatorVersion, operatorVersion)

		callback := func() {
			in.Log.Info(fmt.Sprintf("Applied %s annotation", coh.AnnotationOperatorVersion), "value", operatorVersion)
		}

		applied, err := in.ThreeWayPatchWithCallback(ctx, c.Name, c, c, latest, callback)
		if err != nil {
			return false, errors.Wrapf(err, "failed to update Coherence resource %s/%s with operatorVersion annotation", c.Namespace, c.Name)
		}
		return applied, nil
	}
	return false, nil
}

// ensureFinalizerApplied ensures the finalizer is applied to the Coherence resource
func (in *CoherenceReconciler) ensureFinalizerApplied(ctx context.Context, c *coh.Coherence) (bool, error) {
	if !controllerutil.ContainsFinalizer(c, coh.CoherenceFinalizer) {
		// Re-fetch the Coherence resource to ensure we have the most recent copy
		latest := &coh.Coherence{}
		c.DeepCopyInto(latest)
		controllerutil.AddFinalizer(latest, coh.CoherenceFinalizer)

		callback := func() {
			in.Log.Info("Added finalizer to Coherence resource", "Namespace", c.Namespace, "Name", c.Name, "Finalizer", coh.CoherenceFinalizer)
		}

		// Perform a three-way patch to apply the finalizer
		applied, err := in.ThreeWayPatchWithCallback(ctx, c.Name, c, c, latest, callback)
		if err != nil {
			return false, errors.Wrapf(err, "failed to update Coherence resource %s/%s with finalizer", c.Namespace, c.Name)
		}
		return applied, nil
	}
	return false, nil
}

// ensureFinalizerApplied ensures the finalizer is removed from the Coherence resource
func (in *CoherenceReconciler) ensureFinalizerRemoved(ctx context.Context, c *coh.Coherence) error {
	if controllerutil.RemoveFinalizer(c, coh.CoherenceFinalizer) {
		err := in.GetClient().Update(ctx, c)
		if err != nil {
			in.Log.Info("Failed to remove the finalizer from the Coherence resource, it looks like it had already been deleted")
			return err
		}
	}
	return nil
}

// finalizeDeployment performs any required finalizer tasks for the Coherence resource
func (in *CoherenceReconciler) finalizeDeployment(ctx context.Context, c *coh.Coherence) error {
	// determine whether we can skip service suspension
	if viper.GetBool(operator.FlagSkipServiceSuspend) {
		in.Log.Info("Skipping suspension of Coherence services in deployment " + c.Name +
			operator.FlagSkipServiceSuspend + " is set to true")
		return nil
	}
	if !c.Spec.IsSuspendServicesOnShutdown() {
		in.Log.Info("Skipping suspension of Coherence services in deployment " + c.Name +
			" Spec.SuspendServicesOnShutdown is set to false")
		return nil
	}
	if c.GetReplicas() == 0 {
		in.Log.Info("Skipping suspension of Coherence services in deployment " + c.Name +
			" Spec.Replicas is zero")
		return nil
	}

	in.Log.Info("Finalizing Coherence resource", "Namespace", c.Namespace, "Name", c.Name)
	// Get the StatefulSet
	sts, stsExists, err := in.MaybeFindStatefulSet(ctx, c.Namespace, c.Name)
	if err != nil {
		return errors.Wrapf(err, "getting StatefulSet %s/%s", c.Namespace, c.Name)
	}
	if stsExists {
		if sts.Status.ReadyReplicas == 0 {
			in.Log.Info("Skipping suspension of Coherence services in deployment " + c.Name + " - No Pods are ready")
		} else {
			// Do service suspension...
			p := probe.CoherenceProbe{
				Client:        in.GetClient(),
				Config:        in.GetManager().GetConfig(),
				EventRecorder: events.NewOwnedEventRecorder(c, in.GetEventRecorder()),
			}
			if p.SuspendServices(ctx, c, sts) == probe.ServiceSuspendFailed {
				return fmt.Errorf("failed to suspend services")
			}
		}
	}
	return nil
}

// ensureOperatorSecret ensures that the Operator configuration secret exists in the namespace.
func (in *CoherenceReconciler) ensureOperatorSecret(ctx context.Context, deployment *coh.Coherence, c client.Client, log logr.Logger) error {
	namespace := deployment.Namespace
	s := &coreV1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:        coh.OperatorConfigName,
			Namespace:   namespace,
			Labels:      deployment.CreateGlobalLabels(),
			Annotations: deployment.CreateGlobalAnnotations(),
		},
	}

	err := c.Get(ctx, types.NamespacedName{Name: coh.OperatorConfigName, Namespace: namespace}, s)
	if err != nil && !apierrors.IsNotFound(err) {
		return err
	}

	restHostAndPort := rest.GetServerHostAndPort()

	oldValue := s.Data[coh.OperatorConfigKeyHost]
	if oldValue == nil || string(oldValue) != restHostAndPort {
		// data is different so create/update

		if s.StringData == nil {
			s.StringData = make(map[string]string)
		}
		s.StringData[coh.OperatorConfigKeyHost] = restHostAndPort

		log.Info("Operator configuration updated", "Key", coh.OperatorConfigKeyHost, "OldValue", string(oldValue), "NewValue", restHostAndPort)
		if apierrors.IsNotFound(err) {
			// for some reason we're getting here even if the secret exists so delete it!!
			_ = c.Delete(ctx, s)
			log.Info("Creating configuration secret " + coh.OperatorConfigName)
			err = c.Create(ctx, s)
		} else {
			log.Info("Updating configuration secret " + coh.OperatorConfigName)
			err = c.Update(ctx, s)
		}
	}

	return err
}

// watchSecondaryResource registers the secondary resource reconcilers to watch the resources to be reconciled
func watchSecondaryResource(mgr ctrl.Manager, s reconciler.SecondaryResourceReconciler, owner coh.CoherenceResource) error {
	var err error
	if !s.CanWatch() {
		// this reconciler does not do watches
		return nil
	}

	// Create a new controller
	c, err := controller.New(s.GetControllerName(), s.GetManager(), controller.Options{Reconciler: s.GetReconciler()})
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
