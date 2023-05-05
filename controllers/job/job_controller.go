/*
 * Copyright (c) 2019, 2023, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package job

import (
	"context"
	"fmt"
	"github.com/go-logr/logr"
	coh "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/controllers/reconciler"
	"github.com/oracle/coherence-operator/pkg/utils"
	"github.com/pkg/errors"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"time"
)

const (
	// The name of this controller. This is used in events, log messages, etc.
	controllerName = "controllers.Job"

	CreateMessage        string = "successfully created Job for Coherence resource '%s'"
	FailedToScaleMessage string = "failed to scale Coherence resource %s from %d to %d due to error\n%s"
	FailedToPatchMessage string = "failed to patch Coherence resource %s due to error\n%s"

	EventReasonScale string = "Scaling"
)

// blank assignment to verify that ReconcileServiceMonitor implements reconcile.Reconciler.
// If the reconcile.Reconciler API was to change then we'd get a compile error here.
var _ reconcile.Reconciler = &ReconcileJob{}

// NewJobReconciler returns a new Job reconciler.
func NewJobReconciler(mgr manager.Manager) reconciler.SecondaryResourceReconciler {

	r := &ReconcileJob{
		ReconcileSecondaryResource: reconciler.ReconcileSecondaryResource{
			Kind:     coh.ResourceTypeJob,
			Template: &batchv1.Job{},
		},
	}

	r.SetCommonReconciler(controllerName, mgr)
	return r
}

// ReconcileJob is a reconciler for Jobs.
type ReconcileJob struct {
	reconciler.ReconcileSecondaryResource
}

func (in *ReconcileJob) GetReconciler() reconcile.Reconciler { return in }

func (in *ReconcileJob) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	logger := in.GetLog().WithValues("Namespace", request.Namespace, "Name", request.Name, "Kind", "Job")
	logger.Info("Starting reconcile")

	// Attempt to lock the requested resource. If the resource is locked then another
	// request for the same resource is already in progress so requeue this one.
	if ok := in.Lock(request); !ok {
		return reconcile.Result{Requeue: true, RequeueAfter: 0}, nil
	}
	// Make sure that the request is unlocked when this method exits
	defer in.Unlock(request)

	storage, err := utils.NewStorage(request.NamespacedName, in.GetManager())
	if err != nil {
		return reconcile.Result{}, err
	}

	result, err := in.ReconcileAllResourceOfKind(ctx, request, nil, storage)
	logger.Info("Completed reconcile")
	return result, err
}

// ReconcileAllResourceOfKind will process the specified reconcile request for the specified deployment.
// The previous state being reconciled can be obtained from the storage parameter.
func (in *ReconcileJob) ReconcileAllResourceOfKind(ctx context.Context, request reconcile.Request, deployment *coh.Coherence, storage utils.Storage) (reconcile.Result, error) {
	logger := in.GetLog().WithValues("Namespace", request.Namespace, "Name", request.Name)
	logger.Info("Reconciling Job for deployment")

	result := reconcile.Result{}

	// Fetch the Job's current state
	stsCurrent, stsExists, err := in.MaybeFindJob(ctx, request.Namespace, request.Name)
	if err != nil {
		logger.Info("Finished reconciling Job for deployment. Error getting Job", "error", err.Error())
		return result, errors.Wrapf(err, "getting Job %s/%s", request.Namespace, request.Name)
	}

	if stsExists && stsCurrent.GetDeletionTimestamp() != nil {
		logger.Info("Finished reconciling Job. The Job is being deleted")
		// The Job exists but is being deleted
		return result, nil
	}

	if stsExists && deployment == nil {
		// find the owning Coherence resource
		if deployment, err = in.FindOwningCoherenceResource(ctx, stsCurrent); err != nil {
			logger.Info("Finished reconciling Job. Error finding parent Coherence resource", "error", err.Error())
			return reconcile.Result{}, err
		}
	}

	switch {
	case deployment == nil || deployment.GetReplicas() == 0:
		// The Coherence resource does not exist, or it exists and is scaling down to zero replicas
		if stsExists {
			// The Job does exist though, so it needs to be deleted.
			if deployment != nil {
				// If we get here, we must be scaling down to zero as the Coherence resource exists
				// If the Coherence resource did not exist then service suspension already happened
				// when the Coherence resource was deleted.
				logger.Info("Scaling down to zero")
				err = in.UpdateDeploymentStatusActionsState(ctx, request.NamespacedName, false)
				// TODO: what to do with error?
				if err != nil {
					logger.Info("Error updating deployment status", "error", err.Error())
				}
			}
			// delete the Job
			err = in.Delete(ctx, request.Namespace, request.Name, logger)
		} else {
			// The Job and parent resource has been deleted so no more to do
			_, err = in.UpdateDeploymentStatus(ctx, request)
			return reconcile.Result{}, err
		}
	case !stsExists:
		// Job does not exist but deployment does so create the Job (checking any start quorum)
		result, err = in.createJob(ctx, deployment, storage, logger)
	default:
		// Both Job and deployment exists so this is maybe an update
		result, err = in.updateJob(ctx, deployment, stsCurrent, storage, logger)
	}

	if err != nil {
		logger.Info("Finished reconciling Job with error", "error", err.Error())
		return result, err
	}

	logger.Info("Finished reconciling Job for deployment")
	return result, nil
}

func (in *ReconcileJob) createJob(ctx context.Context, deployment *coh.Coherence, storage utils.Storage, logger logr.Logger) (reconcile.Result, error) {
	logger.Info("Creating Job")

	ok, reason := in.CanCreate(ctx, deployment)

	if !ok {
		// start quorum not met, send event and update deployment status
		in.GetEventRecorder().Event(deployment, corev1.EventTypeNormal, "Waiting", reason)
		_ = in.UpdateDeploymentStatusCondition(ctx, deployment.GetNamespacedName(), coh.Condition{
			Type:    coh.ConditionTypeWaiting,
			Status:  corev1.ConditionTrue,
			Reason:  "StatusQuorum",
			Message: reason,
		})
		return reconcile.Result{Requeue: true, RequeueAfter: time.Second * 30}, nil
	}

	err := in.Create(ctx, deployment.Name, storage, logger)
	if err == nil {
		// ensure that the deployment has a Created status
		err := in.UpdateDeploymentStatusPhase(ctx, deployment.GetNamespacedName(), coh.ConditionTypeCreated)
		if err != nil {
			return reconcile.Result{}, errors.Wrap(err, "updating deployment status")
		}

		// send a successful creation event
		msg := fmt.Sprintf(CreateMessage, deployment.Name)
		in.GetEventRecorder().Event(deployment, corev1.EventTypeNormal, reconciler.EventReasonCreated, msg)
	}

	logger.Info("Created Job")
	return reconcile.Result{}, err
}

func (in *ReconcileJob) updateJob(ctx context.Context, deployment *coh.Coherence, current *batchv1.Job, storage utils.Storage, logger logr.Logger) (reconcile.Result, error) {
	logger.Info("Updating job")

	var err error
	result := reconcile.Result{}

	// get the desired resource state from the store
	resource, found := storage.GetLatest().GetResource(coh.ResourceTypeJob, current.Name)
	if !found {
		// Desired state not found requeue and the request should sort itself out next time around
		logger.Info("Cannot locate desired state for Job, possibly a deletion, re-queuing request")
		return reconcile.Result{Requeue: true}, nil
	}
	if resource.IsDelete() {
		// we should never get here, requeue and the request should sort itself out next time around
		logger.Info("In update path for Job, but is a deletion - re-queuing request")
		return reconcile.Result{Requeue: true}, nil
	}

	desired := resource.Spec.(*batchv1.Job)
	desiredReplicas := in.getReplicas(desired)
	currentReplicas := in.getReplicas(current)

	switch {
	case currentReplicas < desiredReplicas:
		// scale up - if also updating we do the rolling upgrade first followed by the
		// scale up so that we do not do a rolling upgrade of the bigger scaled up cluster

		// try the patch first
		result, err = in.patchJob(ctx, deployment, current, desired, storage, logger)
		if err == nil && !result.Requeue {
			// there was nothing else to patch, so we can do the scale up
			result, err = in.scale(ctx, deployment, current, currentReplicas, desiredReplicas)
		}
	case currentReplicas > desiredReplicas:
		// scale down - if also updating we scale down followed by a rolling upgrade so that
		// we do the rolling upgrade on the smaller scaled down cluster.

		// do the scale down
		_, err = in.scale(ctx, deployment, current, currentReplicas, desiredReplicas)
		// requeue the request so that we do any upgrade next time around
		result.Requeue = true
	default:
		// just an update
		result, err = in.patchJob(ctx, deployment, current, desired, storage, logger)
	}

	return result, err
}

// A nil safe method to get the number of replicas for a Job
func (in *ReconcileJob) getReplicas(sts *batchv1.Job) int32 {
	if sts == nil || sts.Spec.Parallelism == nil {
		return 1
	}
	return *sts.Spec.Parallelism
}

func (in *ReconcileJob) scale(ctx context.Context, deployment *coh.Coherence, job *batchv1.Job, current, desired int32) (reconcile.Result, error) {
	// if the Job is not stable we cannot scale (e.g. it might already be in the middle of a rolling upgrade)
	logger := in.GetLog().WithValues("Namespace", deployment.Name, "Name", deployment.Name)
	logger.Info("Scaling Job", "Current", current, "Desired", desired)

	// ensure that the deployment has a Scaling status
	if err := in.UpdateDeploymentStatusPhase(ctx, deployment.GetNamespacedName(), coh.ConditionTypeScaling); err != nil {
		logger.Error(err, "Error updating deployment status to Scaling")
	}

	return in.parallelScale(ctx, deployment, job, desired)
}

// parallelScale will scale the Job by the required amount in one request.
func (in *ReconcileJob) parallelScale(ctx context.Context, deployment *coh.Coherence, sts *batchv1.Job, replicas int32) (reconcile.Result, error) {
	logger := in.GetLog().WithValues("Namespace", deployment.Name, "Name", deployment.Name)
	logger.Info("Scaling Job", "Replicas", replicas)

	events := in.GetEventRecorder()

	// Update this Coherence resource's status
	deployment.Status.Phase = coh.ConditionTypeScaling
	deployment.Status.Replicas = replicas

	if err := in.UpdateDeploymentStatusPhase(ctx, deployment.GetNamespacedName(), coh.ConditionTypeScaling); err != nil {
		logger.Error(err, "Error updating deployment status to Scaling")
	}

	// Create the desired state
	stsDesired := &batchv1.Job{}
	sts.DeepCopyInto(stsDesired)
	stsDesired.Spec.Parallelism = &replicas

	// ThreeWayPatch theJob to trigger it to scale
	_, err := in.ThreeWayPatch(ctx, sts.Name, sts, sts, stsDesired)
	if err != nil {
		// send a failed scale event
		msg := fmt.Sprintf("failed to scale Job %s from %d to %d", sts.Name, in.getReplicas(sts), replicas)
		events.Event(deployment, corev1.EventTypeNormal, EventReasonScale, msg)
		return reconcile.Result{}, errors.Wrap(err, msg)
	}

	// send a successful scale event
	msg := fmt.Sprintf("scaled Job %s from %d to %d", sts.Name, in.getReplicas(sts), replicas)
	events.Event(deployment, corev1.EventTypeNormal, EventReasonScale, msg)

	return reconcile.Result{}, nil
}

// Patch the Job if required, returning a bool to indicate whether a patch was applied.
func (in *ReconcileJob) patchJob(ctx context.Context, deployment *coh.Coherence, current, desired *batchv1.Job, storage utils.Storage, logger logr.Logger) (reconcile.Result, error) {
	hashMatches := in.HashLabelsMatch(current, storage)
	resource, _ := storage.GetPrevious().GetResource(coh.ResourceTypeJob, current.GetName())
	original := &batchv1.Job{}

	if resource.IsPresent() {
		err := resource.As(original)
		if err != nil {
			return in.HandleErrAndRequeue(ctx, err, deployment, fmt.Sprintf(FailedToPatchMessage, deployment.Name, err.Error()), logger)
		}
	} else {
		// there was no previous
		original = desired
	}

	errorList := coh.ValidateJobUpdate(desired, original)
	if len(errorList) > 0 {
		msg := fmt.Sprintf("upddates to the statefuleset would have been invalid, the update will not be re-queued: %v", errorList)
		events := in.GetEventRecorder()
		events.Event(deployment, corev1.EventTypeWarning, reconciler.EventReasonUpdated, msg)
		return reconcile.Result{Requeue: false}, fmt.Errorf(msg)
	}

	// We NEVER change the replicas or Status in an update.
	// Replicas is handled by scaling, so we always set the desired replicas to match the current replicas
	desired.Spec.Parallelism = current.Spec.Parallelism
	original.Spec.Parallelism = current.Spec.Parallelism

	// We NEVER patch finalizers
	original.ObjectMeta.Finalizers = current.ObjectMeta.Finalizers
	desired.ObjectMeta.Finalizers = current.ObjectMeta.Finalizers

	// We need to ensure we do not create a patch due to differences in
	// Job Status, so we blank out the status fields
	desired.Status = batchv1.JobStatus{}
	current.Status = batchv1.JobStatus{}
	original.Status = batchv1.JobStatus{}

	desiredPodSpec := desired.Spec.Template
	currentPodSpec := desired.Spec.Template
	originalPodSpec := desired.Spec.Template

	// ensure we do not patch any fields that may be set by a previous version of the Operator
	// as this will cause a rolling update of the Pods, typically these are fields where
	// the Operator sets defaults, and we changed the default behaviour
	in.BlankContainerFields(deployment, &desiredPodSpec)
	in.BlankContainerFields(deployment, &currentPodSpec)
	in.BlankContainerFields(deployment, &originalPodSpec)

	// Sort the environment variables, so we do not patch on just a re-ordering of env vars
	in.SortEnvForAllContainers(&desiredPodSpec)
	in.SortEnvForAllContainers(&currentPodSpec)
	in.SortEnvForAllContainers(&originalPodSpec)

	// ensure the Coherence image is present so that we do not patch on a Coherence resource
	// from pre-3.1.x that does not have images set
	if deployment.Spec.Image == nil {
		cohImage := in.GetCoherenceImage(&desiredPodSpec)
		in.SetCoherenceImage(&originalPodSpec, cohImage)
		in.SetCoherenceImage(&currentPodSpec, cohImage)
	}

	// ensure the Operator image is present so that we do not patch on a Coherence resource
	// from pre-3.1.x that does not have images set
	if deployment.Spec.CoherenceUtils == nil || deployment.Spec.CoherenceUtils.Image == nil {
		operatorImage := in.GetOperatorImage(&desiredPodSpec)
		in.SetOperatorImage(&originalPodSpec, operatorImage)
		in.SetOperatorImage(&currentPodSpec, operatorImage)
	}

	// a callback function that the 3-way patch method will call just before it applies a patch
	// if there is any patch to apply, this will check StatusHA if required and update the deployment status
	callback := func() {
		// ensure that the deployment has an "Upgrading" status
		if err := in.UpdateDeploymentStatusPhase(ctx, deployment.GetNamespacedName(), coh.ConditionTypeRollingUpgrade); err != nil {
			logger.Error(err, "Error updating deployment status to Upgrading")
		}
	}

	// fix the CreationTimestamp so that it is not in the patch
	desired.SetCreationTimestamp(current.GetCreationTimestamp())
	// create the patch to see whether there is anything to update
	patch, data, err := in.CreateThreeWayPatch(current.GetName(), original, desired, current, reconciler.PatchIgnore)
	if err != nil {
		return reconcile.Result{}, errors.Wrapf(err, "failed to create patch for Job/%s", current.GetName())
	}

	if patch == nil {
		// nothing to patch so just return
		return reconcile.Result{}, nil
	}

	logger.Info("Updating Job")

	// now apply the patch
	patched, err := in.ApplyThreeWayPatchWithCallback(ctx, current.GetName(), current, patch, data, callback)

	// log the result of patching
	switch {
	case err != nil:
		logger.Info("Error patching Job " + err.Error())
		return in.HandleErrAndRequeue(ctx, err, deployment, fmt.Sprintf(FailedToPatchMessage, deployment.Name, err.Error()), logger)
	case !patched:
		return reconcile.Result{Requeue: true}, nil
	}

	if patched && hashMatches {
		logger.Info("Patch applied to job even though hashes matched (possible external update)")
	}

	return reconcile.Result{}, nil
}
