/*
 * Copyright (c) 2020, 2021, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package reconciler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-logr/logr"
	"github.com/operator-framework/operator-lib/status"
	coh "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/pkg/utils"
	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sync"
)

//goland:noinspection GoUnusedConst
const (
	// PatchIgnore - If the patch json is this we can skip the patch.
	PatchIgnore = "{\"metadata\":{\"creationTimestamp\":null},\"status\":{\"replicas\":0}}"

	// EventReasonCreated is the reason description for a created event.
	EventReasonCreated string = "Created"
	// EventReasonUpdated is the reason description for an updated event.
	EventReasonUpdated string = "Updated"
	// EventReasonFailed is the reason description for a failed event.
	EventReasonFailed string = "Failed"
	// EventReasonDeleted is the reason description for a deleted event.
	EventReasonDeleted string = "Deleted"
)

var (
	// commonMutex is the common lock mutex used by the reconcilers.
	commonMutex = &sync.Mutex{}
	// commonLocks is the map of locked namespace names.
	commonLocks = make(map[types.NamespacedName]bool)
)

// ----- CommonReconciler -----------------------------------------------------

type BaseReconciler interface {
	GetControllerName() string
	GetManager() manager.Manager
	GetClient() client.Client
	GetEventRecorder() record.EventRecorder
	GetLog() logr.Logger
	GetReconciler() reconcile.Reconciler
	SetPatchType(patchType types.PatchType)
}

// CommonReconciler is a base controller structure.
type CommonReconciler struct {
	name      string
	mgr       manager.Manager
	locks     map[types.NamespacedName]bool
	mutex     *sync.Mutex
	logger    logr.Logger
	patchType types.PatchType
}

func (in *CommonReconciler) GetControllerName() string       { return in.name }
func (in *CommonReconciler) GetManager() manager.Manager     { return in.mgr }
func (in *CommonReconciler) GetClient() client.Client        { return in.mgr.GetClient() }
func (in *CommonReconciler) GetMutex() *sync.Mutex           { return in.mutex }
func (in *CommonReconciler) GetPatchType() types.PatchType   { return in.patchType }
func (in *CommonReconciler) SetPatchType(pt types.PatchType) { in.patchType = pt }
func (in *CommonReconciler) GetEventRecorder() record.EventRecorder {
	return in.mgr.GetEventRecorderFor(in.name)
}
func (in *CommonReconciler) GetLog() logr.Logger {
	if in.logger == nil {
		in.logger = logf.Log.WithName(in.name)
	}
	return in.logger
}

func (in *CommonReconciler) SetCommonReconciler(name string, mgr manager.Manager) {
	in.name = name
	in.mgr = mgr
	in.mutex = commonMutex
	in.logger = logf.Log.WithName(name)
	in.patchType = types.StrategicMergePatchType
}

// Lock attempts to lock the requested resource.
func (in *CommonReconciler) Lock(request reconcile.Request) bool {
	if in == nil {
		return false
	}
	in.mutex.Lock()
	defer in.mutex.Unlock()

	if in.locks == nil {
		in.locks = commonLocks
	}

	_, found := in.locks[request.NamespacedName]
	if found {
		in.logger.V(2).Info("Resource " + request.Namespace + "/" + request.Name + " is locked")
		return false
	}

	in.locks[request.NamespacedName] = true
	in.logger.V(2).Info(fmt.Sprintf("Acquired lock for %s/%s", request.Namespace, request.Name))
	return true
}

// Unlock unlocks the requested resource
func (in *CommonReconciler) Unlock(request reconcile.Request) {
	if in != nil {
		in.mutex.Lock()
		defer in.mutex.Unlock()

		in.logger.V(2).Info(fmt.Sprintf("Released lock for %s/%s", request.Namespace, request.Name))
		delete(in.locks, request.NamespacedName)
	}
}

// UpdateDeploymentStatus updates the Coherence resource's status.
func (in *CommonReconciler) UpdateDeploymentStatus(ctx context.Context, request reconcile.Request) (*coh.Coherence, error) {
	var err error
	var sts *appsv1.StatefulSet
	sts, _, err = in.MaybeFindStatefulSet(ctx, request.Namespace, request.Name)
	if err != nil {
		// an error occurred
		err = errors.Wrapf(err, "getting StatefulSet %s", request.Name)
		return nil, err
	}

	deployment := &coh.Coherence{}
	err = in.GetClient().Get(ctx, request.NamespacedName, deployment)
	switch {
	case err != nil && apierrors.IsNotFound(err):
		// deployment not found - possibly deleted
		err = nil
	case err != nil:
		// an error occurred
		err = errors.Wrapf(err, "getting deployment %s", request.Name)
	case deployment.GetDeletionTimestamp() != nil:
		// deployment is being deleted
		err = nil
	default:
		updated := deployment.DeepCopy()
		var stsStatus *appsv1.StatefulSetStatus
		if sts == nil {
			stsStatus = nil
		} else {
			stsStatus = &sts.Status
		}
		if updated.Status.Update(deployment, stsStatus) {
			err = in.GetClient().Status().Update(ctx, updated)
		}
	}
	return deployment, err
}

// UpdateDeploymentStatusPhase updates the Coherence resource's status.
func (in *CommonReconciler) UpdateDeploymentStatusPhase(ctx context.Context, key types.NamespacedName, phase status.ConditionType) error {
	return in.UpdateDeploymentStatusCondition(ctx, key, status.Condition{Type: phase, Status: corev1.ConditionTrue})
}

// UpdateDeploymentStatusCondition updates the Coherence resource's status.
func (in *CommonReconciler) UpdateDeploymentStatusCondition(ctx context.Context, key types.NamespacedName, c status.Condition) error {
	var err error
	deployment := &coh.Coherence{}
	err = in.GetClient().Get(ctx, key, deployment)
	switch {
	case err != nil && apierrors.IsNotFound(err):
		// deployment not found - possibly deleted
		err = nil
	case err != nil:
		// an error occurred
		err = errors.Wrapf(err, "getting deployment %s", key.Name)
	case deployment.GetDeletionTimestamp() != nil:
		// deployment is being deleted
		err = nil
	default:
		updated := deployment.DeepCopy()
		if updated.Status.SetCondition(deployment, c) {
			patch, err := in.CreateTwoWayPatchOfType(types.MergePatchType, deployment.Name, updated, deployment)
			if err != nil {
				return errors.Wrap(err, "creating Coherence resource status patch")
			}
			if patch != nil {
				err = in.GetClient().Status().Patch(ctx, deployment, patch)
				if err != nil {
					return errors.Wrap(err, "updating Coherence resource status")
				}
			}
		}
	}
	return err
}

// UpdateDeploymentStatusHash updates the Coherence resource's status hash.
func (in *CommonReconciler) UpdateDeploymentStatusHash(ctx context.Context, key types.NamespacedName, hash string) error {
	var err error
	deployment := &coh.Coherence{}
	err = in.GetClient().Get(ctx, key, deployment)
	switch {
	case err != nil && apierrors.IsNotFound(err):
		// deployment not found - possibly deleted
		err = nil
	case err != nil:
		// an error occurred
		err = errors.Wrapf(err, "getting deployment %s", key.Name)
	case deployment.GetDeletionTimestamp() != nil:
		// deployment is being deleted
		err = nil
	default:
		if deployment.Status.Hash != hash {
			updated := deployment.DeepCopy()
			updated.Status.Hash = hash
			patch, err := in.CreateTwoWayPatchOfType(types.MergePatchType, deployment.Name, updated, deployment)
			if err != nil {
				return errors.Wrap(err, "creating Coherence resource status patch")
			}
			if patch != nil {
				err = in.GetClient().Status().Patch(ctx, deployment, patch)
				if err != nil {
					return errors.Wrap(err, "updating Coherence resource status")
				}
			}
		}
	}
	return err
}

// MaybeFindStatefulSet finds the required StatefulSet, returning the StatefulSet and a flag indicating whether it was found.
func (in *CommonReconciler) MaybeFindStatefulSet(ctx context.Context, namespace, name string) (*appsv1.StatefulSet, bool, error) {
	sts := &appsv1.StatefulSet{}
	err := in.GetClient().Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, sts)
	switch {
	case err != nil && apierrors.IsNotFound(err):
		return nil, false, nil
	case err != nil:
		return sts, false, err
	default:
		return sts, true, nil
	}
}

// TwoWayPatch performs a two-way merge patch on the resource.
func (in *CommonReconciler) TwoWayPatch(ctx context.Context, name string, current, desired client.Object) (bool, error) {
	patch, err := in.CreateTwoWayPatch(name, desired, current, PatchIgnore)
	if err != nil {
		kind := current.GetObjectKind().GroupVersionKind().Kind
		return false, errors.Wrapf(err, "failed to create patch for %s/%s", kind, name)
	}

	if patch == nil {
		// nothing to patch so just return
		return false, nil
	}

	err = in.GetManager().GetClient().Patch(ctx, current, patch)
	if err != nil {
		kind := current.GetObjectKind().GroupVersionKind().Kind
		return false, errors.Wrapf(err, "cannot patch  %s/%s", kind, name)
	}

	return true, nil
}

// CreateTwoWayPatch creates a two-way patch between the original state, the current state and the desired state of a k8s resource.
func (in *CommonReconciler) CreateTwoWayPatch(name string, desired, current runtime.Object, ignore ...string) (client.Patch, error) {
	return in.CreateTwoWayPatchOfType(in.patchType, name, desired, current, ignore...)
}

// CreateTwoWayPatchOfType creates a two-way patch between the current state and the desired state of a k8s resource.
func (in *CommonReconciler) CreateTwoWayPatchOfType(patchType types.PatchType, name string, desired, current runtime.Object, ignore ...string) (client.Patch, error) {
	currentData, err := json.Marshal(current)
	if err != nil {
		return nil, errors.Wrap(err, "serializing current configuration")
	}
	desiredData, err := json.Marshal(desired)
	if err != nil {
		return nil, errors.Wrap(err, "serializing desired configuration")
	}

	// Get a versioned object
	versionedObject := in.asVersioned(desired)

	patchMeta, err := strategicpatch.NewPatchMetaFromStruct(versionedObject)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create patch metadata from object")
	}

	data, err := strategicpatch.CreateTwoWayMergePatchUsingLookupPatchMeta(currentData, desiredData, patchMeta)
	if err != nil {
		return nil, errors.Wrap(err, "creating three-way patch")
	}

	// check whether the patch counts as no-patch
	ignore = append(ignore, "{}")
	for _, s := range ignore {
		if string(data) == s {
			// empty patch
			return nil, err
		}
	}

	// log the patch
	kind := current.GetObjectKind().GroupVersionKind().Kind
	logger := in.GetLog()
	if logger != nil {
		in.GetLog().V(2).Info(fmt.Sprintf("Created patch for %s/%s\n%s", kind, name, string(data)))
	}

	return client.RawPatch(patchType, data), nil
}

// ThreeWayPatch performs a three-way merge patch on the resource returning true if a patch was required otherwise false.
func (in *CommonReconciler) ThreeWayPatch(ctx context.Context, name string, current, original, desired client.Object) (bool, error) {
	return in.ThreeWayPatchWithCallback(ctx, name, current, original, desired, nil)
}

// ThreeWayPatchWithCallback performs a three-way merge patch on the resource returning true if a patch was required otherwise false.
func (in *CommonReconciler) ThreeWayPatchWithCallback(ctx context.Context, name string, current, original, desired client.Object, callback func()) (bool, error) {
	kind := current.GetObjectKind().GroupVersionKind().Kind
	// fix the CreationTimestamp so that it is not in the patch
	desired.(metav1.Object).SetCreationTimestamp(current.(metav1.Object).GetCreationTimestamp())
	// create the patch
	patch, data, err := in.CreateThreeWayPatch(name, original, desired, current, PatchIgnore)
	if err != nil {
		return false, errors.Wrapf(err, "failed to create patch for %s/%s", kind, name)
	}

	if patch == nil {
		// nothing to patch so just return
		return false, nil
	}

	return in.ApplyThreeWayPatchWithCallback(ctx, name, current, patch, data, callback)
}

// ApplyThreeWayPatchWithCallback performs a three-way merge patch on the resource returning true if a patch was required otherwise false.
func (in *CommonReconciler) ApplyThreeWayPatchWithCallback(ctx context.Context, name string, current client.Object, patch client.Patch, data []byte, callback func()) (bool, error) {
	kind := current.GetObjectKind().GroupVersionKind().Kind

	// execute any callback
	if callback != nil {
		callback()
	}

	in.GetLog().WithValues().Info(fmt.Sprintf("Patching %s/%s with %s", kind, name, string(data)))
	err := in.GetManager().GetClient().Patch(ctx, current, patch)
	if err != nil {
		return false, errors.Wrapf(err, "failed to patch  %s/%s with %s", kind, name, string(data))
	}

	return true, nil
}

// CreateThreeWayPatch creates a three-way patch between the original state, the current state and the desired state of a k8s resource.
func (in *CommonReconciler) CreateThreeWayPatch(name string, original, desired, current runtime.Object, ignore ...string) (client.Patch, []byte, error) {
	data, err := in.CreateThreeWayPatchData(original, desired, current)
	if err != nil {
		return nil, data, errors.Wrap(err, "creating three-way patch")
	}

	// check whether the patch counts as no-patch
	ignore = append(ignore, "{}")
	for _, s := range ignore {
		if string(data) == s {
			// empty patch
			return nil, data, err
		}
	}

	// log the patch
	kind := current.GetObjectKind().GroupVersionKind().Kind
	in.GetLog().Info(fmt.Sprintf("Created patch for %s/%s %s", kind, name, string(data)))

	return client.RawPatch(in.patchType, data), data, nil
}

// CreateThreeWayPatchData creates a three-way patch between the original state, the current state and the desired state of a k8s resource.
func (in *CommonReconciler) CreateThreeWayPatchData(original, desired, current runtime.Object) ([]byte, error) {
	originalData, err := json.Marshal(original)
	if err != nil {
		return nil, errors.Wrap(err, "serializing original configuration")
	}
	currentData, err := json.Marshal(current)
	if err != nil {
		return nil, errors.Wrap(err, "serializing current configuration")
	}
	desiredData, err := json.Marshal(desired)
	if err != nil {
		return nil, errors.Wrap(err, "serializing desired configuration")
	}

	// Get a versioned object
	versionedObject := in.asVersioned(desired)

	patchMeta, err := strategicpatch.NewPatchMetaFromStruct(versionedObject)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create patch metadata from object")
	}

	return strategicpatch.CreateThreeWayMergePatch(originalData, desiredData, currentData, patchMeta, true)
}

// asVersioned converts the given object into a runtime.Template with the correct group and version set.
func (in *CommonReconciler) asVersioned(obj runtime.Object) runtime.Object {
	var gv = runtime.GroupVersioner(schema.GroupVersions(scheme.Scheme.PrioritizedVersionsAllGroups()))
	if obj, err := runtime.ObjectConvertor(scheme.Scheme).ConvertToVersion(obj, gv); err == nil {
		return obj
	}
	return obj
}

// HandleErrAndRequeue is the common error handler
func (in *CommonReconciler) HandleErrAndRequeue(ctx context.Context, err error, deployment *coh.Coherence, msg string, logger logr.Logger) (reconcile.Result, error) {
	return in.Failed(ctx, err, deployment, msg, true, logger)
}

// HandleErrAndFinish is the common error handler
func (in *CommonReconciler) HandleErrAndFinish(ctx context.Context, err error, deployment *coh.Coherence, msg string, logger logr.Logger) (reconcile.Result, error) {
	return in.Failed(ctx, err, deployment, msg, false, logger)
}

// Failed is a common error handler
// ToDo: we need to be able to add some form of back-off so that failures are requeued with a growing back-off time
func (in *CommonReconciler) Failed(ctx context.Context, err error, deployment *coh.Coherence, msg string, requeue bool, logger logr.Logger) (reconcile.Result, error) {
	if err == nil {
		logger.V(0).Info(msg)
	} else {
		logger.V(0).Info(msg + ": " + err.Error())
	}

	if deployment != nil {
		// update the status to failed.
		deployment.Status.Phase = coh.ConditionTypeFailed
		if e := in.GetClient().Status().Update(ctx, deployment); e != nil {
			// There isn't much we can do, we're already handling an error
			logger.V(0).Info("failed to update deployment status due to: " + e.Error())
		}

		// send a failure event
		in.GetEventRecorder().Event(deployment, corev1.EventTypeNormal, "failed", msg)
	}

	if requeue {
		return reconcile.Result{Requeue: true}, nil
	}
	return reconcile.Result{Requeue: false}, nil
}

// FindOwningCoherenceResource finds the owning Coherence resource.
func (in *CommonReconciler) FindOwningCoherenceResource(ctx context.Context, o client.Object) (*coh.Coherence, error) {
	if o != nil {
		for _, ref := range o.GetOwnerReferences() {
			if ref.Kind == coh.ResourceTypeCoherence.Name() {
				return in.FindDeployment(ctx, o.GetNamespace(), ref.Name)
			}
		}
	}
	return nil, nil
}

// FindDeployment finds the Coherence resource.
func (in *CommonReconciler) FindDeployment(ctx context.Context, namespace, name string) (*coh.Coherence, error) {
	deployment, _, err := in.MaybeFindDeployment(ctx, namespace, name)
	return deployment, err
}

// MaybeFindDeployment possibly finds a Coherence resource.
func (in *CommonReconciler) MaybeFindDeployment(ctx context.Context, namespace, name string) (*coh.Coherence, bool, error) {
	deployment := &coh.Coherence{}
	err := in.GetClient().Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, deployment)

	switch {
	case err != nil && apierrors.IsNotFound(err):
		// the deployment does not exist
		return nil, false, nil
	case err != nil:
		// an error occurred
		return deployment, false, err
	default:
		// the deployment exists
		return deployment, true, nil
	}
}

// ----- SecondaryResourceReconciler ----------------------------------------------

// SecondaryResourceReconciler is a reconciler for sub-resources.
type SecondaryResourceReconciler interface {
	BaseReconciler
	GetTemplate() client.Object
	ReconcileAllResourceOfKind(context.Context, reconcile.Request, *coh.Coherence, utils.Storage) (reconcile.Result, error)
	CanWatch() bool
}

// ----- ReconcileSecondaryResource ----------------------------------------------

// ReconcileSecondaryResource reconciles secondary resources of a specific Kind for a specific Coherence resource
type ReconcileSecondaryResource struct {
	CommonReconciler
	Template  client.Object
	Kind      coh.ResourceType
	SkipWatch bool
}

func (in *ReconcileSecondaryResource) GetTemplate() client.Object { return in.Template }
func (in *ReconcileSecondaryResource) CanWatch() bool             { return !in.SkipWatch }

// ReconcileAllResourceOfKind reconciles the state of all the desired resources of the specified Kind for the reconciler
func (in *ReconcileSecondaryResource) ReconcileAllResourceOfKind(ctx context.Context, request reconcile.Request, deployment *coh.Coherence, storage utils.Storage) (reconcile.Result, error) {
	logger := in.GetLog().WithValues("Namespace", request.Namespace, "Name", request.Name, "Kind", in.Kind.Name())
	logger.Info(fmt.Sprintf("Reconciling all %v", in.Kind))

	var err error
	resources := storage.GetLatest().GetResourcesOfKind(in.Kind)
	for _, res := range resources {
		if res.IsDelete() {
			if err = in.Delete(ctx, request.Namespace, res.Name, logger); err != nil {
				logger.Info(fmt.Sprintf("Finished reconciling all %v with error", in.Kind), "error", err.Error())
				return reconcile.Result{}, err
			}
		} else {
			if err = in.ReconcileSingleResource(ctx, request.Namespace, res.Name, deployment, storage, logger); err != nil {
				logger.Info(fmt.Sprintf("Finished reconciling all %v with error", in.Kind), "error", err.Error())
				return reconcile.Result{}, err
			}
		}
	}
	logger.Info(fmt.Sprintf("Finished reconciling all %v", in.Kind))
	return reconcile.Result{}, nil
}

// HashLabelsMatch determines whether the Coherence Hash label on the specified Object matches the hash on the storage.
func (in *ReconcileSecondaryResource) HashLabelsMatch(o metav1.Object, storage utils.Storage) bool {
	storageHash, storageHashFound := storage.GetHash()
	objectHash, objectHashFound := o.GetLabels()[coh.LabelCoherenceHash]
	return storageHashFound == objectHashFound && storageHash == objectHash
}

// ReconcileSingleResource reconciles a specific resource.
func (in *ReconcileSecondaryResource) ReconcileSingleResource(ctx context.Context, namespace, name string, owner *coh.Coherence, storage utils.Storage, logger logr.Logger) error {
	logger = logger.WithValues("Resource", name)
	logger.Info(fmt.Sprintf("Reconciling single %v", in.Kind))

	// Fetch the resource's current state
	resource, exists, err := in.FindResource(ctx, namespace, name)
	if err != nil {
		// Error reading the object - requeue the request.
		// We can't call the error handler as we do not even have a owning Coherence resource.
		// We log the error and do not requeue the request.
		return errors.Wrapf(err, "getting %s %s/%s", in.Kind, namespace, name)
	}

	if exists && resource.GetDeletionTimestamp() != nil {
		// The resource exists but is being deleted
		exists = false
	}

	if owner == nil {
		// try to find the owning Coherence resource
		if owner, err = in.FindOwningCoherenceResource(ctx, resource); err != nil {
			return err
		}
	}

	if owner != nil && in.Kind.Name() == coh.ResourceTypeSecret.Name() && name == owner.GetName() {
		// this a reconcile event for the storage secret, we can ignore it
		logger.Info(fmt.Sprintf("Finished reconciling single %v", in.Kind))
		return nil
	}

	if storage == nil && owner != nil {
		if storage, err = utils.NewStorage(owner.GetNamespacedName(), in.GetManager()); err != nil {
			return err
		}
	}

	switch {
	case owner == nil || owner.GetReplicas() == 0:
		if exists {
			// The owning Coherence resource does not exist (or is scaled down to zero) but the resource still does,
			// ensure that the resource is deleted.
			// This should not actually be required as everything is owned by the owning Coherence resource
			// and there should be a cascaded delete by k8s, so it's belt and braces.
			err = in.Delete(ctx, namespace, name, logger)
		}
	case !exists:
		// Resource does not exist but owning Coherence resource does so create it
		err = in.Create(ctx, name, storage, logger)
	default:
		// Both the resource and owning Coherence resource exist so this is maybe an update
		err = in.Update(ctx, name, resource.(client.Object), storage, logger)
	}

	logger.Info(fmt.Sprintf("Finished reconciling single %v", in.Kind))
	return err
}

// NewFromTemplate creates a resource from a template resource.
func (in *ReconcileSecondaryResource) NewFromTemplate(namespace, name string) client.Object {
	// create a new resource from copying the empty template
	resource := in.Template.DeepCopyObject()
	// set the resource's namespace and name
	metaObj := resource.(metav1.Object)
	metaObj.SetNamespace(namespace)
	metaObj.SetName(name)
	return resource.(client.Object)
}

// Create the specified resource
func (in *ReconcileSecondaryResource) Create(ctx context.Context, name string, storage utils.Storage, logger logr.Logger) error {
	logger.Info(fmt.Sprintf("Creating %v", in.Kind))
	// Get the resource state
	resource, found := storage.GetLatest().GetResource(in.Kind, name)
	if !found {
		logger.Info(fmt.Sprintf("Cannot create %v as latest state not present in store", in.Kind))
		return nil
	}
	// create the resource
	if err := in.GetClient().Create(ctx, resource.Spec); err != nil {
		return errors.Wrapf(err, "failed to create %v/%s", in.Kind, name)
	}
	return nil
}

// Delete the resource
func (in *ReconcileSecondaryResource) Delete(ctx context.Context, namespace, name string, logger logr.Logger) error {
	logger.Info("Deleting StatefulSet")
	// create a new resource from copying the empty template
	resource := in.NewFromTemplate(namespace, name)

	// Do the deletion
	err := in.GetClient().Delete(ctx, resource, client.PropagationPolicy(metav1.DeletePropagationBackground))
	// check for an error (we ignore not-found as this means k8s has already done the delete for us)
	if err != nil && !apierrors.IsNotFound(err) {
		err = errors.Wrapf(err, "failed to delete %s/%s", in.Kind, name)
		logger.Error(err, fmt.Sprintf("Failed to delete %s/%s", in.Kind, name))
	}
	return nil
}

// Update possibly updates the resource if the current state differs from the desired state.
func (in *ReconcileSecondaryResource) Update(ctx context.Context, name string, current client.Object, storage utils.Storage, logger logr.Logger) error {
	logger.Info(fmt.Sprintf("Updating %v", in.Kind))

	hashMatches := in.HashLabelsMatch(current, storage)
	original, _ := storage.GetPrevious().GetResource(in.Kind, name)
	desired, found := storage.GetLatest().GetResource(in.Kind, name)
	if !found {
		logger.Info(fmt.Sprintf("Cannot update %v as latest state not present in store", in.Kind))
		return nil
	}

	patched, err := in.ThreeWayPatch(ctx, name, current, original.Spec, desired.Spec)
	if patched && hashMatches {
		logger.Info(fmt.Sprintf("Patch applied to %v even though hashes matched (possible external update)", in.Kind))
	}
	return err
}

func (in *ReconcileSecondaryResource) FindResource(ctx context.Context, namespace, name string) (client.Object, bool, error) {
	object := in.NewFromTemplate(namespace, name)
	err := in.GetClient().Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, object.(client.Object))
	var found bool
	switch {
	case err != nil && apierrors.IsNotFound(err):
		// we can ignore not found errors
		err = nil
		found = false
	case err != nil:
		// some error other than not-found occurred
		found = false
	default:
		// resource was found
		found = true
	}
	return object, found, err
}
