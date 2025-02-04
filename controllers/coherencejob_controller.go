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
	"github.com/oracle/coherence-operator/controllers/job"
	"github.com/oracle/coherence-operator/controllers/reconciler"
	"github.com/oracle/coherence-operator/controllers/secret"
	"github.com/oracle/coherence-operator/controllers/servicemonitor"
	"github.com/oracle/coherence-operator/pkg/clients"
	"github.com/oracle/coherence-operator/pkg/operator"
	"github.com/oracle/coherence-operator/pkg/rest"
	"github.com/oracle/coherence-operator/pkg/utils"
	"github.com/pkg/errors"
	coreV1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"strconv"

	coh "github.com/oracle/coherence-operator/api/v1"
)

// CoherenceJobReconciler reconciles a CoherenceJob object
type CoherenceJobReconciler struct {
	client.Client
	reconciler.CommonReconciler
	ClientSet   clients.ClientSet
	Log         logr.Logger
	Scheme      *runtime.Scheme
	reconcilers []reconciler.SecondaryResourceReconciler
}

var (
	jobControllerName = "controllers.CoherenceJob"
)

// blank assignment to verify that CoherenceJobReconciler implements reconcile.Reconciler
// There will be a compile-time error here if this breaks
var _ reconcile.Reconciler = &CoherenceJobReconciler{}

func (in *CoherenceJobReconciler) Reconcile(ctx context.Context, request ctrl.Request) (ctrl.Result, error) {
	deployment := &coh.CoherenceJob{}
	return in.ReconcileDeployment(ctx, request, deployment)
}

func (in *CoherenceJobReconciler) ReconcileDeployment(ctx context.Context, request ctrl.Request, deployment *coh.CoherenceJob) (ctrl.Result, error) {
	var err error

	log := in.Log.WithValues("namespace", request.Namespace, "name", request.Name)

	log.Info("Reconciling CoherenceJob resource")

	// Attempt to lock the requested resource. If the resource is locked then another
	// request for the same resource is already in progress so requeue this one.
	if ok := in.Lock(request); !ok {
		log.Info("CoherenceJob resource " + request.Namespace + "/" + request.Name + " is already locked, requeue request")
		return reconcile.Result{Requeue: true, RequeueAfter: 0}, nil
	}
	// Make sure that the request is unlocked when this method exits
	defer in.Unlock(request)

	// Fetch the CoherenceJob resource instance
	err = in.GetClient().Get(ctx, types.NamespacedName{Namespace: request.Namespace, Name: request.Name}, deployment)
	if err != nil {
		if apierrors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected.
			// Return and don't requeue
			log.Info("CoherenceJob resource not found. Ignoring request since object must be deleted.")
			return ctrl.Result{}, nil
		}
		// Error reading the current deployment state from k8s.
		return reconcile.Result{}, errors.Wrap(err, "getting CoherenceJob resource")
	}

	// Check whether this is a deletion
	deleteTime := deployment.GetDeletionTimestamp()
	if deleteTime != nil {
		// Check whether finalization needs to be run
		if controllerutil.ContainsFinalizer(deployment, coh.CoherenceFinalizer) {
			log.Info("CoherenceJob resource deleted at " + deleteTime.String() + ", running finalizer")
			// Remove the finalizer. Once all finalizers have been
			// removed, the object will be deleted.
			controllerutil.RemoveFinalizer(deployment, coh.CoherenceFinalizer)
			err := in.GetClient().Update(ctx, deployment)
			if err != nil {
				if apierrors.IsNotFound(err) {
					log.Info("Failed to remove the finalizer from the CoherenceJob resource, it looks like it had already been deleted")
					return ctrl.Result{}, nil
				}
				return ctrl.Result{}, errors.Wrap(err, "trying to remove finalizer from CoherenceJob resource")
			}
		} else {
			log.Info("CoherenceJob resource deleted at " + deleteTime.String() + ", finalizer already removed")
		}

		return ctrl.Result{}, nil
	}

	// The request is an add or update

	// Ensure the hash label is present (it should have been added by the web-hook, so this should be a no-op).
	// The hash may not have been added if the CoherenceJob resource was added/modified when the Operator was uninstalled.
	if hashApplied, err := in.ensureHashApplied(ctx, deployment); hashApplied || err != nil {
		return ctrl.Result{Requeue: true}, err
	}

	spec, _ := deployment.GetJobResourceSpec()

	// ensure that the deployment has an initial status
	status := deployment.GetStatus()
	if status.Phase == "" {
		err := in.UpdateCoherenceJobStatusPhase(ctx, request.NamespacedName, coh.ConditionTypeInitialized)
		if err != nil {
			return reconcile.Result{}, err
		}
	}

	// Check whether the deployment has a replica count specified
	// Ideally we'd do this with a validating/defaulting web-hook but maybe in a later version.
	//spec = deployment.GetSpec()
	if spec.Replicas == nil {
		// No replica count, so we patch the deployment to have the default replicas value.
		// The reason we do this, is because the kubectl scale command will fail if the replicas
		// field is not set to a non-nil value.
		patch := deployment.DeepCopyResource()
		replicas := deployment.GetReplicas()
		patch.SetReplicas(replicas)
		_, err = in.ThreeWayPatch(ctx, deployment.GetName(), deployment, deployment, patch)
		if err != nil {
			log.Info("Added default replicas to CoherenceJob resource, re-queuing request", "Replicas", strconv.Itoa(int(replicas)))
			return reconcile.Result{}, err
		}
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
		return in.HandleErrAndRequeue(ctx, err, nil, fmt.Sprintf(reconcileFailedMessage, request.Name, request.Namespace, err), in.Log)
	}

	hash := deployment.GetLabels()[coh.LabelCoherenceHash]
	storeHash, _ := storage.GetHash()
	var desiredResources coh.Resources

	desiredResources, err = checkJobHash(deployment, storage, log)
	if err != nil {
		return in.HandleErrAndRequeue(ctx, err, nil, fmt.Sprintf(createResourcesFailedMessage, request.Name, request.Namespace, err), in.Log)
	}

	// create the result
	result := ctrl.Result{Requeue: false}

	log.Info("Reconciling CoherenceJob resource secondary resources", "hash", hash, "store", storeHash)

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
	// The hash may not have been added if the CoherenceJob resource was added/modified when the Operator was uninstalled.
	if applied, err := in.ensureVersionAnnotationApplied(ctx, deployment); applied || err != nil {
		return ctrl.Result{Requeue: true}, err
	}

	// process the secondary resources in the order they should be created
	var failures []Failure
	for _, rec := range in.reconcilers {
		log.Info("Reconciling CoherenceJob resource secondary resources", "controller", rec.GetControllerName())
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
		if err = in.UpdateCoherenceJobStatusPhase(ctx, request.NamespacedName, coh.ConditionTypeStopped); err != nil {
			return result, errors.Wrap(err, "error updating deployment status")
		}
	}

	// Update the Status with the hash
	if err = in.UpdateDeploymentStatusHash(ctx, request.NamespacedName, hash); err != nil {
		return result, errors.Wrap(err, "error updating deployment status hash")
	}

	log.Info("Finished reconciling CoherenceJob resource", "requeue", result.Requeue, "after", result.RequeueAfter.String())
	return result, nil
}

func (in *CoherenceJobReconciler) SetupWithManager(mgr ctrl.Manager, cs clients.ClientSet) error {
	SetupMonitoringResources(mgr)

	// Create the sub-resource reconcilers IN THE ORDER THAT RESOURCES MUST BE CREATED.
	// This is important to ensure, for example, that a ConfigMap is created before the
	// StatefulSet that uses it.
	reconcilers := []reconciler.SecondaryResourceReconciler{
		reconciler.NewNamedConfigMapReconciler(mgr, cs, "controllers.JobConfigMap"),
		secret.NewNamedSecretReconciler(mgr, cs, "controllers.JobSecret"),
		reconciler.NewNamedServiceReconciler(mgr, cs, "controllers.JobService"),
		servicemonitor.NewNamedServiceMonitorReconciler(mgr, cs, "controllers.JobServiceMonitor"),
		job.NewJobReconciler(mgr, cs),
	}

	in.reconcilers = reconcilers
	in.SetCommonReconciler(jobControllerName, mgr, cs)
	in.SetPatchType(types.MergePatchType)

	template := &coh.CoherenceJob{}

	// Watch for changes to secondary resources
	for _, sub := range reconcilers {
		if err := watchSecondaryResource(mgr, sub, template); err != nil {
			return err
		}
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(template).
		Named("coherencejob").
		Complete(in)
}

func (in *CoherenceJobReconciler) GetReconciler() reconcile.Reconciler { return in }

// ensureHashApplied ensures that the hash label is present in the CoherenceJob resource, patching it if required
func (in *CoherenceJobReconciler) ensureHashApplied(ctx context.Context, c *coh.CoherenceJob) (bool, error) {
	currentHash := ""
	labels := c.GetLabels()
	if len(labels) > 0 {
		currentHash = labels[coh.LabelCoherenceHash]
	}

	// Re-fetch the CoherenceJob resource to ensure we have the most recent copy
	latest := c.DeepCopy()
	hash, _ := coh.EnsureJobHashLabel(latest)

	if currentHash != hash {
		if c.IsBeforeVersion("3.4.2") {
			// Before 3.4.2 there was a bug calculating the has in the defaulting web-hook
			// This would cause the hashes to be different here, when in fact they should not be
			// If the CoherenceJob resource being processes has no version annotation, or a version
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

		applied, err := in.ThreeWayPatchWithCallback(ctx, c.GetName(), c, c, latest, callback)
		if err != nil {
			return false, errors.Wrapf(err, "failed to update CoherenceJob resource %s/%s with hash", c.GetNamespace(), c.GetName())
		}
		return applied, nil
	}
	return false, nil
}

// ensureVersionAnnotationApplied ensures that the version annotation is present in the CoherenceJob resource, patching it if required
func (in *CoherenceJobReconciler) ensureVersionAnnotationApplied(ctx context.Context, c coh.CoherenceResource) (bool, error) {
	currentVersion, _ := c.GetVersionAnnotation()
	operatorVersion := operator.GetVersion()

	if currentVersion != operatorVersion {
		// make a copy of the CoherenceJob resource to use in the three-way patch
		latest := c.DeepCopyResource()
		latest.AddAnnotation(coh.AnnotationOperatorVersion, operatorVersion)

		callback := func() {
			in.Log.Info(fmt.Sprintf("Applied %s annotation", coh.AnnotationOperatorVersion), "value", operatorVersion)
		}

		applied, err := in.ThreeWayPatchWithCallback(ctx, c.GetName(), c, c, latest, callback)
		if err != nil {
			return false, errors.Wrapf(err, "failed to update CoherenceJob resource %s/%s with operatorVersion annotation", c.GetNamespace(), c.GetName())
		}
		return applied, nil
	}
	return false, nil
}

// ensureOperatorSecret ensures that the Operator configuration secret exists in the namespace.
func (in *CoherenceJobReconciler) ensureOperatorSecret(ctx context.Context, deployment *coh.CoherenceJob, c client.Client, log logr.Logger) error {
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
