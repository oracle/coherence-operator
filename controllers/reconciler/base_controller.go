/*
 * Copyright (c) 2020, 2025, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package reconciler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-logr/logr"
	coh "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/pkg/clients"
	"github.com/oracle/coherence-operator/pkg/operator"
	"github.com/oracle/coherence-operator/pkg/utils"
	"github.com/pkg/errors"
	"golang.org/x/mod/semver"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
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
	"sort"
	"strings"
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
	// EventReasonReconciling is the reason description for an reconciling event.
	EventReasonReconciling string = "Reconciling"
	// EventReasonScaling is the reason description for an scaling event.
	EventReasonScaling string = "Scaling"
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
	GetClientSet() clients.ClientSet
	GetEventRecorder() record.EventRecorder
	GetLog() logr.Logger
	GetReconciler() reconcile.Reconciler
	SetPatchType(patchType types.PatchType)
}

// CommonReconciler is a base controller structure.
type CommonReconciler struct {
	name      string
	mgr       manager.Manager
	clientSet clients.ClientSet
	locks     map[types.NamespacedName]bool
	mutex     *sync.Mutex
	logger    logr.Logger
	patchType types.PatchType
}

func (in *CommonReconciler) GetControllerName() string       { return in.name }
func (in *CommonReconciler) GetManager() manager.Manager     { return in.mgr }
func (in *CommonReconciler) GetClient() client.Client        { return in.mgr.GetClient() }
func (in *CommonReconciler) GetClientSet() clients.ClientSet { return in.clientSet }
func (in *CommonReconciler) GetMutex() *sync.Mutex           { return in.mutex }
func (in *CommonReconciler) GetPatchType() types.PatchType   { return in.patchType }
func (in *CommonReconciler) SetPatchType(pt types.PatchType) { in.patchType = pt }
func (in *CommonReconciler) GetEventRecorder() record.EventRecorder {
	return in.mgr.GetEventRecorderFor(in.name)
}
func (in *CommonReconciler) GetLog() logr.Logger {
	return in.logger
}

func (in *CommonReconciler) SetCommonReconciler(name string, mgr manager.Manager, cs clients.ClientSet) {
	in.name = name
	in.mgr = mgr
	in.clientSet = cs
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

// UpdateCoherenceStatusPhase updates the Coherence resource's status.
func (in *CommonReconciler) UpdateCoherenceStatusPhase(ctx context.Context, key types.NamespacedName, phase coh.ConditionType) error {
	return in.UpdateCoherenceStatusCondition(ctx, key, coh.Condition{Type: phase, Status: corev1.ConditionTrue})
}

// UpdateCoherenceJobStatusPhase updates the CoherenceJob resource's status.
func (in *CommonReconciler) UpdateCoherenceJobStatusPhase(ctx context.Context, key types.NamespacedName, phase coh.ConditionType) error {
	return in.UpdateCoherenceJobStatusCondition(ctx, key, coh.Condition{Type: phase, Status: corev1.ConditionTrue})
}

// UpdateCoherenceStatusCondition updates the Coherence resource's status.
func (in *CommonReconciler) UpdateCoherenceStatusCondition(ctx context.Context, key types.NamespacedName, c coh.Condition) error {
	return in.updateDeploymentStatusCondition(ctx, key, c, &coh.Coherence{})
}

// UpdateCoherenceJobStatusCondition updates the CoherenceJob resource's status.
func (in *CommonReconciler) UpdateCoherenceJobStatusCondition(ctx context.Context, key types.NamespacedName, c coh.Condition) error {
	return in.updateDeploymentStatusCondition(ctx, key, c, &coh.CoherenceJob{})
}

// UpdateDeploymentStatusCondition updates the Coherence resource's status.
func (in *CommonReconciler) updateDeploymentStatusCondition(ctx context.Context, key types.NamespacedName, c coh.Condition, deployment coh.CoherenceResource) error {
	var err error
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
		updated := deployment.DeepCopyResource()
		status := updated.GetStatus()
		if status.SetCondition(deployment, c) {
			patch, err := in.CreateTwoWayPatchOfType(types.MergePatchType, deployment.GetName(), updated, deployment)
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
			updated.Status.Version = operator.GetVersion()
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

// MaybeFindJob finds the required Job, returning the StatefulSet and a flag indicating whether it was found.
func (in *CommonReconciler) MaybeFindJob(ctx context.Context, namespace, name string) (*batchv1.Job, bool, error) {
	job := &batchv1.Job{}
	err := in.GetClient().Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, job)
	switch {
	case err != nil && apierrors.IsNotFound(err):
		return nil, false, nil
	case err != nil:
		return job, false, err
	default:
		return job, true, nil
	}
}

// UpdateDeploymentStatusActionsState updates the Coherence resource's status ActionsExecuted flag.
func (in *CommonReconciler) UpdateDeploymentStatusActionsState(ctx context.Context, key types.NamespacedName, actionExecuted bool) error {
	deployment := &coh.Coherence{}
	err := in.GetClient().Get(ctx, key, deployment)
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
		if deployment.Status.ActionsExecuted != actionExecuted {
			updated := deployment.DeepCopy()
			updated.Status.ActionsExecuted = actionExecuted
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

// IsVersionAnnotationEqualOrBefore returns true if the specified object
// has a version annotation with a version the same as ot before the
// specified version or has no version annotation.
func (in *CommonReconciler) IsVersionAnnotationEqualOrBefore(m metav1.Object, version string) bool {
	a := m.GetAnnotations()
	if a != nil {
		v, found := a[coh.AnnotationOperatorVersion]
		if found && v != "" {
			if v == version {
				return true
			}
			if version[0] != 'v' {
				version = "v" + version
			}
			if v[0] != 'v' {
				v = "v" + v
			}
			return semver.Compare(v, version) <= 0
		}
	}
	return true
}

// CanCreate determines whether any specified start quorum has been met.
func (in *CommonReconciler) CanCreate(ctx context.Context, deployment coh.CoherenceResource) (bool, string) {
	spec := deployment.GetSpec()
	if len(spec.StartQuorum) == 0 {
		// there is no start quorum
		return true, ""
	}

	logger := in.GetLog().WithValues("Namespace", deployment.GetNamespace(), "Name", deployment.GetName())
	logger.Info("Checking deployment start quorum")

	var quorum []string

	for _, q := range spec.StartQuorum {
		if q.Deployment == "" {
			// this start-quorum does not have a dependency name so skip it
			continue
		}
		// work out which Namespace to look for the dependency in
		var namespace string
		if q.Namespace == "" {
			// start-quorum does not specify a namespace so use the same one as the deployment
			namespace = deployment.GetNamespace()
		} else {
			// start-quorum does specify a namespace so use it
			namespace = q.Namespace
		}
		dep, found, err := in.MaybeFindDeployment(ctx, namespace, q.Deployment)
		switch {
		case err != nil:
			// cannot create due to an error looking up the deployment
			quorum = append(quorum, fmt.Sprintf("error finding deployment '%s' - %s", q.Deployment, err.Error()))
		case !found:
			// cannot create as the deployment does not yet exist
			quorum = append(quorum, fmt.Sprintf("deployment '%s/%s' does not exist", namespace, q.Deployment))
		case found && q.PodCount > 0 && dep.Status.ReadyReplicas < q.PodCount:
			// deployment exists and quorum requires a specific number of ready Pods
			quorum = append(quorum, fmt.Sprintf("role '%s/%s' to have %d ready Pods (ready=%d)", namespace, q.Deployment, q.PodCount, dep.Status.ReadyReplicas))
		case found && dep.Status.Phase != coh.ConditionTypeReady:
			// deployment exists and quorum requires all pods ready
			quorum = append(quorum, fmt.Sprintf("deployment '%s' is not ready", q.Deployment))
		}
	}

	if len(quorum) > 0 {
		reason := "Waiting for start quorum to be met: \"" + strings.Join(quorum, "\" and \"") + "\""
		logger.Info(reason)
		return false, reason
	}
	return true, ""
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

	in.GetLog().V(2).Info(fmt.Sprintf("Created patch for %s/%s\n%s", kind, name, string(data)))

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

	in.GetLog().WithValues().Info(fmt.Sprintf("Patching %s/%s", kind, name), "Patch", string(data))
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
	in.GetLog().Info(fmt.Sprintf("Created patch for %s/%s", kind, name), "Patch", string(data))

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
func (in *CommonReconciler) HandleErrAndRequeue(ctx context.Context, err error, deployment coh.CoherenceResource, msg string, logger logr.Logger) (reconcile.Result, error) {
	return in.Failed(ctx, err, deployment, msg, true, logger)
}

// HandleErrAndFinish is the common error handler
func (in *CommonReconciler) HandleErrAndFinish(ctx context.Context, err error, deployment *coh.Coherence, msg string, logger logr.Logger) (reconcile.Result, error) {
	return in.Failed(ctx, err, deployment, msg, false, logger)
}

// Failed is a common error handler
// ToDo: we need to be able to add some form of back-off so that failures are re-queued with a growing back-off time
func (in *CommonReconciler) Failed(ctx context.Context, err error, deployment coh.CoherenceResource, msg string, requeue bool, logger logr.Logger) (reconcile.Result, error) {
	if err == nil {
		logger.V(0).Info(msg)
	} else {
		logger.V(0).Info(msg + ": " + err.Error())
	}

	if deployment != nil {
		// update the status to failed.
		status := deployment.GetStatus()
		status.Phase = coh.ConditionTypeFailed
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
func (in *CommonReconciler) FindOwningCoherenceResource(ctx context.Context, o client.Object) (coh.CoherenceResource, error) {
	if o != nil {
		for _, ref := range o.GetOwnerReferences() {
			if ref.Kind == coh.ResourceTypeCoherence.Name() {
				return in.FindDeployment(ctx, o.GetNamespace(), ref.Name)
			}
			if ref.Kind == coh.ResourceTypeCoherenceJob.Name() {
				return in.FindCoherenceJob(ctx, o.GetNamespace(), ref.Name)
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

// FindCoherenceJob finds the CoherenceJob resource.
func (in *CommonReconciler) FindCoherenceJob(ctx context.Context, namespace, name string) (*coh.CoherenceJob, error) {
	deployment, _, err := in.MaybeFindCoherenceJob(ctx, namespace, name)
	return deployment, err
}

// MaybeFindCoherenceJob possibly finds a CoherenceJob resource.
func (in *CommonReconciler) MaybeFindCoherenceJob(ctx context.Context, namespace, name string) (*coh.CoherenceJob, bool, error) {
	deployment := &coh.CoherenceJob{}
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

// BlankContainerFields blanks out any fields that we do not want to include in the patch
// Typically these are fields where we changed the default behaviour in the newer Operator versions
func (in *CommonReconciler) BlankContainerFields(deployment coh.CoherenceResource, template *corev1.PodTemplateSpec) {
	spec := deployment.GetSpec()
	if spec.Affinity == nil {
		// affinity not set by user so do not diff on it
		template.Spec.Affinity = nil
	}
	in.BlankOperatorInitContainerFields(template)
	in.BlankCoherenceContainerFields(template)
}

// BlankOperatorInitContainerFields blanks out any fields that may have been set by a previous Operator version.
// DO NOT blank out anything that the user has control over as they may have
// updated them, so we need to include them in the patch
func (in *CommonReconciler) BlankOperatorInitContainerFields(template *corev1.PodTemplateSpec) {
	for i := range template.Spec.InitContainers {
		c := template.Spec.InitContainers[i]
		if c.Name == coh.ContainerNameOperatorInit {
			// This is the Operator init-container
			// blank out the container command field
			c.Command = []string{}
			// set the updated init-container back into the StatefulSet
			template.Spec.InitContainers[i] = c
		}
	}
}

// BlankCoherenceContainerFields blanks out any fields that may have been set by a previous Operator version.
// DO NOT blank out anything that the user has control over as they may have
// updated them, so we need to include them in the patch
func (in *CommonReconciler) BlankCoherenceContainerFields(template *corev1.PodTemplateSpec) {
	for i := range template.Spec.Containers {
		c := template.Spec.Containers[i]
		if c.Name == coh.ContainerNameCoherence {
			// This is the Coherence Container
			// blank out the container command field
			c.Command = []string{}
			// set the updated container back into the StatefulSet
			template.Spec.Containers[i] = c
		}
	}
}

// SortEnvForAllContainers sorts the environment variable slice for all containers.
func (in *CommonReconciler) SortEnvForAllContainers(template *corev1.PodTemplateSpec) {
	for i := range template.Spec.InitContainers {
		c := template.Spec.InitContainers[i]
		in.SortEnvForContainer(&c)
		template.Spec.InitContainers[i] = c
	}
	for i := range template.Spec.Containers {
		c := template.Spec.Containers[i]
		in.SortEnvForContainer(&c)
		template.Spec.Containers[i] = c
	}
}

// SortEnvForContainer sorts the environment variable slice for a container.
func (in *CommonReconciler) SortEnvForContainer(c *corev1.Container) {
	sort.Slice(c.Env, func(i, j int) bool {
		return c.Env[i].Name < c.Env[j].Name
	})
}

// GetOperatorImage gets the Operator image name from the init container.
func (in *CommonReconciler) GetOperatorImage(template *corev1.PodTemplateSpec) string {
	for i := range template.Spec.InitContainers {
		c := template.Spec.InitContainers[i]
		if c.Name == coh.ContainerNameOperatorInit {
			return c.Image
		}
	}
	return ""
}

// SetOperatorImage sets the Operator image name in the init container.
func (in *CommonReconciler) SetOperatorImage(template *corev1.PodTemplateSpec, image string) {
	for i := range template.Spec.InitContainers {
		c := template.Spec.InitContainers[i]
		if c.Name == coh.ContainerNameOperatorInit {
			c.Image = image
			template.Spec.InitContainers[i] = c
		}
	}
}

// GetCoherenceImage gets the Coherence image name from the coherence container.
func (in *CommonReconciler) GetCoherenceImage(template *corev1.PodTemplateSpec) string {
	for i := range template.Spec.Containers {
		c := template.Spec.Containers[i]
		if c.Name == coh.ContainerNameCoherence {
			return c.Image
		}
	}
	return ""
}

// SetCoherenceImage sets the Coherence image name in the coherence container.
func (in *CommonReconciler) SetCoherenceImage(template *corev1.PodTemplateSpec, image string) {
	for i := range template.Spec.Containers {
		c := template.Spec.Containers[i]
		if c.Name == coh.ContainerNameCoherence {
			c.Image = image
			template.Spec.Containers[i] = c
		}
	}
}

// ----- SecondaryResourceReconciler ----------------------------------------------

// SecondaryResourceReconciler is a reconciler for sub-resources.
type SecondaryResourceReconciler interface {
	BaseReconciler
	GetTemplate() client.Object
	ReconcileAllResourceOfKind(context.Context, reconcile.Request, coh.CoherenceResource, utils.Storage) (reconcile.Result, error)
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
func (in *ReconcileSecondaryResource) ReconcileAllResourceOfKind(ctx context.Context, request reconcile.Request, deployment coh.CoherenceResource, storage utils.Storage) (reconcile.Result, error) {
	logger := in.GetLog().WithValues("Namespace", request.Namespace, "Name", request.Name, "Kind", in.Kind.Name())

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
	return reconcile.Result{}, nil
}

// HashLabelsMatch determines whether the Coherence Hash label on the specified Object matches the hash on the storage.
func (in *ReconcileSecondaryResource) HashLabelsMatch(o metav1.Object, storage utils.Storage) bool {
	storageHash, storageHashFound := storage.GetHash()
	objectHash, objectHashFound := o.GetLabels()[coh.LabelCoherenceHash]
	in.logger.Info("Checking hash", "ObjectName", o.GetName(), "ObjectKind", in.Kind.Name(), "StorageHash", storageHash, "StorageHashFound", storageHashFound, "ObjectHash", objectHash, "ObjectHashFound", objectHashFound)
	return storageHashFound == objectHashFound && storageHash == objectHash
}

// ReconcileSingleResource reconciles a specific resource.
func (in *ReconcileSecondaryResource) ReconcileSingleResource(ctx context.Context, namespace, name string, owner coh.CoherenceResource, storage utils.Storage, logger logr.Logger) error {
	logger = logger.WithValues("Resource", name)
	logger.Info(fmt.Sprintf("Reconciling %v", in.Kind))

	// Fetch the resource's current state
	resource, exists, err := in.FindResource(ctx, namespace, name)
	if err != nil {
		// Error reading the object - requeue the request.
		// We can't call the error handler as we do not even have an owning Coherence resource.
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
		logger.Info(fmt.Sprintf("Finished reconciling %v", in.Kind))
		return nil
	}

	if storage == nil && owner != nil {
		if storage, err = utils.NewStorage(owner.GetNamespacedName(), in.GetManager()); err != nil {
			return err
		}
	}

	switch {
	case owner == nil:
		if exists {
			// The owning Coherence resource does not exist but the resource still does,
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
		err = in.Update(ctx, name, resource, storage, logger)
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
	logger.Info("Deleting")
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
	hashMatches := in.HashLabelsMatch(current, storage)
	if hashMatches {
		logger.Info(fmt.Sprintf("Nothing to update for %v, hashes match", in.Kind))
		return nil
	}

	original, _ := storage.GetPrevious().GetResource(in.Kind, name)
	desired, found := storage.GetLatest().GetResource(in.Kind, name)
	if !found {
		logger.Info(fmt.Sprintf("Cannot update %v as latest state not present in store", in.Kind))
		return nil
	}

	patched, err := in.ThreeWayPatch(ctx, name, current, original.Spec, desired.Spec)
	if patched && hashMatches {
		logger.Info(fmt.Sprintf("Patch applied to %v", in.Kind))
	}
	return err
}

func (in *ReconcileSecondaryResource) FindResource(ctx context.Context, namespace, name string) (client.Object, bool, error) {
	object := in.NewFromTemplate(namespace, name)
	err := in.GetClient().Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, object)
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
