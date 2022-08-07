/*
 * Copyright (c) 2020, 2022, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package controllers

import (
	"context"
	"fmt"
	"github.com/go-logr/logr"
	"github.com/oracle/coherence-operator/controllers/predicates"
	"github.com/oracle/coherence-operator/controllers/reconciler"
	"github.com/oracle/coherence-operator/controllers/secret"
	"github.com/oracle/coherence-operator/controllers/servicemonitor"
	"github.com/oracle/coherence-operator/controllers/statefulset"
	"github.com/oracle/coherence-operator/pkg/operator"
	"github.com/oracle/coherence-operator/pkg/rest"
	"github.com/oracle/coherence-operator/pkg/utils"
	"github.com/pkg/errors"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/spf13/viper"
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
	"strconv"

	coh "github.com/oracle/coherence-operator/api/v1"
)

const (
	controllerName = "controllers.Coherence"

	reconcileFailedMessage       string = "failed to reconcile Coherence resource '%s' in namespace '%s'\n%s"
	createResourcesFailedMessage string = "create resources for Coherence resource '%s' in namespace '%s' failed\n%s"
)

// CoherenceReconciler reconciles a Coherence object
type CoherenceReconciler struct {
	client.Client
	reconciler.CommonReconciler
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

// +kubebuilder:rbac:groups=coherence.oracle.com,resources=coherence;coherence/finalizers;coherence/status,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=pods;pods/exec;services;endpoints;events;configmaps;secrets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=monitoring.coreos.com,resources=servicemonitors,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=cert-manager.io,resources=certificates;issuers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;update;patch;delete

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
		// Error reading the current deployment state from k8s.
		return reconcile.Result{}, err
	}

	// Check whether this is a deletion
	deleteTime := deployment.GetDeletionTimestamp()
	if deleteTime != nil {
		// Check whether finalization needs to be run
		if utils.StringArrayContains(deployment.GetFinalizers(), coh.CoherenceFinalizer) {
			log.Info("Coherence resource deleted at " + deleteTime.String() + ", running finalizer")
			// Run finalization logic.
			// If the finalization logic fails, don't remove the finalizer so
			// that we can retry during the next reconciliation.
			if err := in.finalizeDeployment(ctx, deployment); err != nil {
				log.Error(err, "Failed to remove finalizer")
				return ctrl.Result{Requeue: true}, nil
			}
			// Remove the finalizer. Once all finalizers have been
			// removed, the object will be deleted.
			controllerutil.RemoveFinalizer(deployment, coh.CoherenceFinalizer)
			err := in.GetClient().Update(ctx, deployment)
			if err != nil {
				return ctrl.Result{}, err
			}
		} else {
			log.Info("Coherence resource deleted at " + deleteTime.String() + ", finalizer already removed")
		}

		return ctrl.Result{}, nil
	}

	// The request is an add or update

	// Ensure the hash label is present (it should have been added by the web-hook, so this should be a no-op).
	// The hash may not have been added if the Coherence resource was added/modified when the Operator was uninstalled.
	if hashApplied, err := in.ensureHashApplied(ctx, deployment); hashApplied || err != nil {
		return ctrl.Result{Requeue: true}, err
	}

	// Add finalizer for this CR if required (it should have been added by the web-hook but may not have been if the
	// Coherence resource was added when the Operator was uninstalled)
	if finalizerApplied, err := in.ensureFinalizerApplied(ctx, deployment); finalizerApplied || err != nil {
		return ctrl.Result{Requeue: true}, err
	}

	// ensure that the deployment has an initial status
	if deployment.Status.Phase == "" {
		err := in.UpdateDeploymentStatusPhase(ctx, request.NamespacedName, coh.ConditionTypeInitialized)
		if err != nil {
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
			log.Info("Added default replicas to Coherence resource, re-queuing request", "Replicas", strconv.Itoa(int(replicas)))
			return reconcile.Result{}, err
		}
	}

	// ensure that the Operator configuration Secret exists
	if err = EnsureOperatorSecret(ctx, request.Namespace, in.GetClient(), in.Log); err != nil {
		err = errors.Wrap(err, "ensuring Operator configuration secret")
		return in.HandleErrAndRequeue(ctx, err, nil, fmt.Sprintf(reconcileFailedMessage, request.Name, request.Namespace, err), in.Log)
	}

	// ensure that the state store exists
	storage, err := utils.NewStorage(request.NamespacedName, in.GetManager())
	if err != nil {
		err = errors.Wrap(err, "obtaining desired state store")
		return in.HandleErrAndRequeue(ctx, err, nil, fmt.Sprintf(reconcileFailedMessage, request.Name, request.Namespace, err), in.Log)
	}

	hash := deployment.GetLabels()[coh.LabelCoherenceHash]
	var desiredResources coh.Resources

	storeHash, found := storage.GetHash()
	if !found || storeHash != hash || deployment.Status.Phase != coh.ConditionTypeReady {
		// Storage state was saved with no hash or a different hash so is not in the desired state
		// or the Coherence resource is not in the Ready state
		// Create the desired resources the deployment
		if desiredResources, err = deployment.Spec.CreateKubernetesResources(deployment); err != nil {
			return in.HandleErrAndRequeue(ctx, err, nil, fmt.Sprintf(createResourcesFailedMessage, request.Name, request.Namespace, err), in.Log)
		}

		if found {

			// The "storeHash" is not "", so it must have been processed by the Operator (could have been a previous version).
			// There was a bug prior to 3.2.8 where the hash was calculated at the wrong point in the defaulting web-hook,
			// so the "currentHash" may be wrong, and hence differ from the recalculated "hash".
			if deployment.IsBeforeVersion("3.2.8") {
				// the AnnotationOperatorVersion annotation was added in the 3.2.8 web-hook, so if it is missing
				// the Coherence resource was added or updated prior to 3.2.8
				// In this case we just ignore the difference in hash.
				// There is an edge case where the Coherence resource could have legitimately been updated whilst
				// the Operator and web-hooks were uninstalled. In that case we would ignore the update until another
				// update is made. The simplest way for the customer to work around this is to add the
				// AnnotationOperatorVersion annotation with some value, which will then be overwritten by the web-hook
				// and the Coherence resource will be correctly processes.
				desiredResources = storage.GetLatest()
				log.Info("Ignoring hash difference for pre-3.2.8 resource", "hash", hash, "store", storeHash)
			}
		}
	} else {
		// storage state was saved with the current hash so is already in the desired state
		desiredResources = storage.GetLatest()
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
	if applied, err := in.ensureVersionAnnotationApplied(ctx, deployment); applied || err != nil {
		return ctrl.Result{Requeue: true}, err
	}

	// process the secondary resources in the order they should be created
	var failures []Failure
	for _, rec := range in.reconcilers {
		log.Info("Reconciling Coherence resource secondary resources", "controller", rec.GetControllerName())
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
		if err = in.UpdateDeploymentStatusPhase(ctx, request.NamespacedName, coh.ConditionTypeStopped); err != nil {
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

func (in *CoherenceReconciler) SetupWithManager(mgr ctrl.Manager) error {
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
		secret.NewSecretReconciler(mgr),
		reconciler.NewServiceReconciler(mgr),
		servicemonitor.NewServiceMonitorReconciler(mgr),
		statefulset.NewStatefulSetReconciler(mgr),
	}

	in.reconcilers = reconcilers
	in.SetCommonReconciler(controllerName, mgr)
	in.SetPatchType(types.MergePatchType)

	// Watch for changes to secondary resources
	for _, sub := range reconcilers {
		if err := in.watchSecondaryResource(sub); err != nil {
			return err
		}
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&coh.Coherence{}).
		Named("coherence").
		Complete(in)
}

// Watch the resources to be reconciled
func (in *CoherenceReconciler) watchSecondaryResource(s reconciler.SecondaryResourceReconciler) error {
	if !s.CanWatch() {
		// this reconciler does not do watches
		return nil
	}

	// Create a new controller
	c, err := controller.New(s.GetControllerName(), s.GetManager(), controller.Options{Reconciler: s.GetReconciler()})
	if err != nil {
		return err
	}

	src := &source.Kind{Type: s.GetTemplate()}
	h := &handler.EnqueueRequestForOwner{IsController: true, OwnerType: &coh.Coherence{}}
	p := predicates.SecondaryPredicate{}
	if err := c.Watch(src, h, p); err != nil {
		return err
	}

	return nil
}

func (in *CoherenceReconciler) GetReconciler() reconcile.Reconciler { return in }

// ensureHashApplied ensures that the hash label is present in the Coherence resource, patching it if required
func (in *CoherenceReconciler) ensureHashApplied(ctx context.Context, c *coh.Coherence) (bool, error) {
	currentHash := ""
	labels := c.GetLabels()
	if len(labels) > 0 {
		currentHash = labels[coh.LabelCoherenceHash]
	}

	// Re-fetch the Coherence resource to ensure we have the most recent copy
	latest := &coh.Coherence{}
	c.DeepCopyInto(latest)
	hash, _ := coh.EnsureHashLabel(latest)

	if currentHash != hash {
		if c.IsBeforeVersion("3.2.8") {
			// Before 3.2.8 there was a bug calculating the has in the defaulting web-hook
			// This would cause the hashes to be different here, when in fact they should not be
			// If the Coherence resource being processes has no version annotation, or a version
			// prior to 3.2.8 then we return as if the hashes matched
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

	version := operator.GetVersion()

	if currentVersion != version {
		// make a copy of the Coherence resource to use in the three-way patch
		latest := c.DeepCopy()
		latest.AddAnnotation(coh.AnnotationOperatorVersion, version)

		callback := func() {
			in.Log.Info(fmt.Sprintf("Applied %s annotation", coh.AnnotationOperatorVersion))
		}

		applied, err := in.ThreeWayPatchWithCallback(ctx, c.Name, c, c, latest, callback)
		if err != nil {
			return false, errors.Wrapf(err, "failed to update Coherence resource %s/%s with version annotation", c.Namespace, c.Name)
		}
		return applied, nil
	}
	return false, nil
}

func (in *CoherenceReconciler) ensureFinalizerApplied(ctx context.Context, c *coh.Coherence) (bool, error) {
	if utils.StringArrayDoesNotContain(c.GetFinalizers(), coh.CoherenceFinalizer) {
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
			probe := statefulset.CoherenceProbe{
				Client: in.GetClient(),
				Config: in.GetManager().GetConfig(),
			}
			if probe.SuspendServices(ctx, c, sts) == statefulset.ServiceSuspendFailed {
				return fmt.Errorf("failed to suspend services")
			}
		}
	}
	return nil
}

// EnsureOperatorSecret ensures that the Operator configuration secret exists in the namespace.
func EnsureOperatorSecret(ctx context.Context, namespace string, c client.Client, log logr.Logger) error {
	s := &coreV1.Secret{}

	err := c.Get(ctx, types.NamespacedName{Name: coh.OperatorConfigName, Namespace: namespace}, s)
	if err != nil && !apierrors.IsNotFound(err) {
		return err
	}

	restHostAndPort := rest.GetServerHostAndPort()

	s.SetNamespace(namespace)
	s.SetName(coh.OperatorConfigName)

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
