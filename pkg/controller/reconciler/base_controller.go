/*
 * Copyright (c) 2020, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package reconciler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-logr/logr"
	"github.com/operator-framework/operator-sdk/pkg/status"
	coh "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
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
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"strings"
	"sync"
)

const (
	// If the patch json is this we can skip the patch
	patchIgnore = "{\"metadata\":{\"creationTimestamp\":null},\"status\":{\"replicas\":0}}"

	EventReasonCreated string = "Created"
	EventReasonUpdated string = "Updated"
	EventReasonFailed  string = "Failed"
	EventReasonDeleted string = "Deleted"
)

var (
	// the common lock mutex used by the reconcilers.
	commonMutex = &sync.Mutex{}
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

// A base controller structure.
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

	patchEnv := os.Getenv("USE_STRATEGIC_PATCH")
	if strings.ToLower(patchEnv) == "true" {
		in.patchType = types.StrategicMergePatchType
	} else {
		in.patchType = types.MergePatchType
	}
}

// Attempt to lock the requested resource.
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

// Unlock the requested resource
func (in *CommonReconciler) Unlock(request reconcile.Request) {
	if in != nil {
		in.mutex.Lock()
		defer in.mutex.Unlock()

		in.logger.V(2).Info(fmt.Sprintf("Released lock for %s/%s", request.Namespace, request.Name))
		delete(in.locks, request.NamespacedName)
	}
}

// Update the Coherence resource's status
func (in *CommonReconciler) UpdateDeploymentStatus(request reconcile.Request) error {
	var err error
	var sts *appsv1.StatefulSet
	sts, _, err = in.MaybeFindStatefulSet(request.Namespace, request.Name)
	if err != nil {
		// an error occurred
		err = errors.Wrapf(err, "getting StatefulSet %s", request.Name)
		return err
	}

	deployment := &coh.Coherence{}
	err = in.GetClient().Get(context.TODO(), request.NamespacedName, deployment)
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
		if updated.Status.Update(deployment, &sts.Status) {
			patch, err := in.CreateTwoWayPatch(deployment.Name, updated, deployment)
			if err != nil {
				return errors.Wrap(err, "creating Coherence resource status patch")
			}
			if patch != nil {
				err = in.GetClient().Status().Patch(context.TODO(), deployment, patch)
				if err != nil {
					return errors.Wrap(err, "updating Coherence resource status")
				}
			}
		}
	}
	return err
}

// Update the Coherence resource's status
func (in *CommonReconciler) UpdateDeploymentStatusPhase(key types.NamespacedName, phase status.ConditionType) error {
	return in.UpdateDeploymentStatusCondition(key, status.Condition{Type: phase, Status: corev1.ConditionTrue})
}

// Update the Coherence resource's status
func (in *CommonReconciler) UpdateDeploymentStatusCondition(key types.NamespacedName, c status.Condition) error {
	var err error
	deployment := &coh.Coherence{}
	err = in.GetClient().Get(context.TODO(), key, deployment)
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
			patch, err := in.CreateTwoWayPatch(deployment.Name, updated, deployment)
			if err != nil {
				return errors.Wrap(err, "creating Coherence resource status patch")
			}
			if patch != nil {
				err = in.GetClient().Status().Patch(context.TODO(), deployment, patch)
				if err != nil {
					return errors.Wrap(err, "updating Coherence resource status")
				}
			}
		}
	}
	return err
}

// Find the required StatefulSet, returning the StatefulSet and a flag indicating whether it was found.
func (in *CommonReconciler) MaybeFindStatefulSet(namespace, name string) (*appsv1.StatefulSet, bool, error) {
	sts := &appsv1.StatefulSet{}
	err := in.GetClient().Get(context.TODO(), types.NamespacedName{Namespace: namespace, Name: name}, sts)
	switch {
	case err != nil && apierrors.IsNotFound(err):
		return sts, false, nil
	case err != nil:
		return sts, false, err
	default:
		return sts, true, nil
	}
}

// Perform a two-way merge patch on the resource.
func (in *CommonReconciler) TwoWayPatch(name string, current, desired runtime.Object) error {
	patch, err := in.CreateTwoWayPatch(name, desired, current, patchIgnore)
	if err != nil {
		kind := current.GetObjectKind().GroupVersionKind().Kind
		return errors.Wrapf(err, "failed to create patch for %s/%s", kind, name)
	}

	if patch == nil {
		// nothing to patch so just return
		return nil
	}

	err = in.GetManager().GetClient().Patch(context.TODO(), current, patch)
	if err != nil {
		kind := current.GetObjectKind().GroupVersionKind().Kind
		return errors.Wrapf(err, "cannot patch  %s/%s", kind, name)
	}

	return nil
}

// Create a two-way patch between the original state, the current state and the desired state of a k8s resource.
func (in *CommonReconciler) CreateTwoWayPatch(name string, desired, current runtime.Object, ignore ...string) (client.Patch, error) {
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

	return client.RawPatch(in.patchType, data), nil
}

// Perform a three-way merge patch on the resource returning true if a patch was required otherwise false.
func (in *CommonReconciler) ThreeWayPatch(name string, current, original, desired runtime.Object) (bool, error) {
	return in.ThreeWayPatchWithCallback(name, current, original, desired, nil)
}

// Perform a three-way merge patch on the resource returning true if a patch was required otherwise false.
func (in *CommonReconciler) ThreeWayPatchWithCallback(name string, current, original, desired runtime.Object, callback func()) (bool, error) {
	kind := current.GetObjectKind().GroupVersionKind().Kind
	// fix the CreationTimestamp so that it is not in the patch
	desired.(metav1.Object).SetCreationTimestamp(current.(metav1.Object).GetCreationTimestamp())
	// create the patch
	patch, err := in.CreateThreeWayPatch(name, original, desired, current, patchIgnore)
	if err != nil {
		return false, errors.Wrapf(err, "failed to create patch for %s/%s", kind, name)
	}

	if patch == nil {
		// nothing to patch so just return
		return false, nil
	}

	// execute any callback
	if callback != nil {
		callback()
	}

	in.GetLog().WithValues().Info(fmt.Sprintf("Patching %s/%s", kind, name))
	err = in.GetManager().GetClient().Patch(context.TODO(), current, patch)
	if err != nil {
		return false, errors.Wrapf(err, "cannot patch  %s/%s", kind, name)
	}

	return true, nil
}

// Create a three-way patch between the original state, the current state and the desired state of a k8s resource.
func (in *CommonReconciler) CreateThreeWayPatch(name string, original, desired, current runtime.Object, ignore ...string) (client.Patch, error) {
	data, err := in.CreateThreeWayPatchData(name, original, desired, current, ignore...)
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
	in.GetLog().Info(fmt.Sprintf("Created patch for %s/%s %s", kind, name, string(data)))

	return client.RawPatch(in.patchType, data), nil
}

// Create a three-way patch between the original state, the current state and the desired state of a k8s resource.
func (in *CommonReconciler) CreateThreeWayPatchData(name string, original, desired, current runtime.Object, ignore ...string) ([]byte, error) {
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
func (in *CommonReconciler) HandleErrAndRequeue(err error, deployment *coh.Coherence, msg string, logger logr.Logger) (reconcile.Result, error) {
	return in.Failed(err, deployment, msg, true, logger)
}

// handleErrAndFinish is the common error handler
func (in *CommonReconciler) HandleErrAndFinish(err error, deployment *coh.Coherence, msg string, logger logr.Logger) (reconcile.Result, error) {
	return in.Failed(err, deployment, msg, false, logger)
}

// failed is the common error handler
// ToDo: we need to be able to add some form of back-off so that failures are re-queued with a growing back-off time
func (in *CommonReconciler) Failed(err error, deployment *coh.Coherence, msg string, requeue bool, logger logr.Logger) (reconcile.Result, error) {
	if err == nil {
		logger.V(0).Info(msg)
	} else {
		logger.Error(err, msg)
	}

	if deployment != nil {
		// update the status to failed.
		deployment.Status.Phase = coh.ConditionTypeFailed
		if e := in.GetClient().Status().Update(context.TODO(), deployment); e != nil {
			// There isn't much we can do, we're already handling an error
			logger.Error(err, "failed to update deployment status")
		}

		// send a failure event
		in.GetEventRecorder().Event(deployment, corev1.EventTypeNormal, "failed", msg)
	}

	if requeue {
		return reconcile.Result{Requeue: true}, nil
	}
	return reconcile.Result{Requeue: false}, nil
}

// Find the Coherence resource for a given request.
func (in *CommonReconciler) FindDeployment(request reconcile.Request) (*coh.Coherence, error) {
	deployment, _, err := in.MaybeFindDeployment(request.Namespace, request.Name)
	return deployment, err
}

// Maybe find a Coherence resource.
func (in *CommonReconciler) MaybeFindDeployment(namespace, name string) (*coh.Coherence, bool, error) {
	deployment := &coh.Coherence{}
	err := in.GetClient().Get(context.TODO(), types.NamespacedName{Namespace: namespace, Name: name}, deployment)

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

// A reconciler for sub-resource.
type SecondaryResourceReconciler interface {
	BaseReconciler
	GetTemplate() runtime.Object
	ReconcileResources(reconcile.Request, *coh.Coherence, utils.Storage) (reconcile.Result, error)
	CanWatch() bool
}

// ----- ReconcileSecondaryResource ----------------------------------------------

// ReconcileSecondaryResource reconciles a secondary resource for a Coherence resource
type ReconcileSecondaryResource struct {
	CommonReconciler
	Template  runtime.Object
	Kind      coh.ResourceType
	SkipWatch bool
}

func (in *ReconcileSecondaryResource) GetTemplate() runtime.Object { return in.Template }
func (in *ReconcileSecondaryResource) CanWatch() bool              { return !in.SkipWatch }

// ReconcileResources reconciles the state of the desired resources for the reconciler
func (in *ReconcileSecondaryResource) ReconcileResources(request reconcile.Request, deployment *coh.Coherence, storage utils.Storage) (reconcile.Result, error) {
	logger := in.GetLog().WithValues("Namespace", request.Namespace, "Name", request.Name)

	var err error
	diff := storage.GetLatest().DiffForKind(in.Kind, storage.GetPrevious())
	for _, res := range diff {
		if res.IsDelete() {
			if err = in.Delete(request.Namespace, res.Name, logger); err != nil {
				logger.Info(fmt.Sprintf("Finished reconciling all %s for deployment with error: %s", in.Kind, err.Error()))
				return reconcile.Result{}, err
			}
		} else {
			if err = in.ReconcileResource(request.Namespace, res.Name, deployment, storage); err != nil {
				logger.Info(fmt.Sprintf("Finished reconciling all %s for deployment with error: %s", in.Kind, err.Error()))
				return reconcile.Result{}, err
			}
		}
	}
	return reconcile.Result{}, nil
}

func (in *ReconcileSecondaryResource) ReconcileResource(namespace, name string, deployment *coh.Coherence, storage utils.Storage) error {
	logger := in.GetLog().WithValues("Namespace", namespace, "Name", name)

	// Fetch the resource's current state
	resource, exists, err := in.FindResource(namespace, name)
	if err != nil {
		// Error reading the object - requeue the request.
		// We can't call the error handler as we do not even have a deployment.
		// We log the error and do not requeue the request.
		return errors.Wrapf(err, "getting %s %s/%s", in.Kind, namespace, name)
	}

	if exists && resource.GetDeletionTimestamp() != nil {
		// The resource exists but is being deleted
		exists = false
	}

	switch {
	case deployment == nil || deployment.GetReplicas() == 0:
		if exists {
			// The deployment does not exist (or is scaled down to zero) but the resource still does,
			// ensure that the resource is deleted.
			// This should not actually be required as everything is owned by the deployment
			// and there should be a cascaded delete by k8s so it's belt and braces.
			err = in.Delete(namespace, name, logger)
		}
	case !exists:
		// Resource does not exist but deployment does so create it
		err = in.Create(name, storage, logger)
	default:
		// Both the resource and deployment exists so this is maybe an update
		err = in.Update(name, resource.(runtime.Object), storage, logger)
	}
	return err
}

// Delete the resource
func (in *ReconcileSecondaryResource) NewFromTemplate(namespace, name string) runtime.Object {
	// create a new resource from copying the empty template
	resource := in.Template.DeepCopyObject()
	// set the resource's namespace and name
	metaObj := resource.(metav1.Object)
	metaObj.SetNamespace(namespace)
	metaObj.SetName(name)
	return resource
}

// Create the specified resource
func (in *ReconcileSecondaryResource) Create(name string, storage utils.Storage, logger logr.Logger) error {
	logger.Info(fmt.Sprintf("Creating %v/%s for deployment", in.Kind, name))
	// Get the resource state
	resource, found := storage.GetLatest().GetResource(in.Kind, name)
	if !found {
		return fmt.Errorf("cannot create %v/%s for deployment as latest state not present in store", in.Kind, name)
	}
	// create the resource
	if err := in.GetClient().Create(context.TODO(), resource.Spec); err != nil {
		logger.Info(fmt.Sprintf("Failed creating %v for deployment - %s", in.Kind, err.Error()))
		return errors.Wrapf(err, "failed to create %v/%s", in.Kind, name)
	}
	return nil
}

// Delete the resource
func (in *ReconcileSecondaryResource) Delete(namespace, name string, logger logr.Logger) error {
	logger.Info(fmt.Sprintf("Deleting %v/%s for deployment", in.Kind, name))
	// create a new resource from copying the empty template
	resource := in.NewFromTemplate(namespace, name)

	// Do the delete
	err := in.GetClient().Delete(context.TODO(), resource, client.PropagationPolicy(metav1.DeletePropagationBackground))
	// check for an error (we ignore not-found as this means k8s has already done the delete for us)
	if err != nil && !apierrors.IsNotFound(err) {
		err = errors.Wrapf(err, "failed to delete %s/%s", in.Kind, name)
		logger.Error(err, fmt.Sprintf("Failed to delete %s/%s", in.Kind, name))
	}
	return nil
}

// Maybe update the resource if the current state differs from the desired state.
func (in *ReconcileSecondaryResource) Update(name string, current runtime.Object, storage utils.Storage, logger logr.Logger) error {
	original, _ := storage.GetPrevious().GetResource(in.Kind, name)
	desired, found := storage.GetLatest().GetResource(in.Kind, name)
	if !found {
		return fmt.Errorf("cannot update %s/%s as latest state not present in store", in.Kind, name)
	}

	_, err := in.ThreeWayPatch(name, current, original.Spec, desired.Spec)
	return err
}

func (in *ReconcileSecondaryResource) FindResource(namespace, name string) (metav1.Object, bool, error) {
	object := in.NewFromTemplate(namespace, name).(metav1.Object)
	err := in.GetClient().Get(context.TODO(), types.NamespacedName{Namespace: namespace, Name: name}, object.(runtime.Object))
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
