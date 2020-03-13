/*
 * Copyright (c) 2019, 2020 Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package coherencerole

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-logr/logr"
	"github.com/go-test/deep"
	helmctl "github.com/operator-framework/operator-sdk/pkg/helm/controller"
	"github.com/operator-framework/operator-sdk/pkg/helm/release"
	"github.com/operator-framework/operator-sdk/pkg/helm/watches"
	"github.com/operator-framework/operator-sdk/pkg/k8sutil"
	coh "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	"github.com/oracle/coherence-operator/pkg/flags"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"os"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"strings"
	"sync"
	"time"
)

const (
	// The name of this controller. This is used in events, log messages, etc.
	controllerName = "coherencerole.controller"

	invalidRoleEventMessage      string = "invalid CoherenceRole '%s' cannot find parent CoherenceCluster '%s'"
	createMessage                string = "created CoherenceInternal '%s' from CoherenceRole '%s' successful"
	createFailedMessage          string = "create CoherenceInternal '%s' from CoherenceRole '%s' failed\n%s"
	updateMessage                string = "updated CoherenceInternal %s from CoherenceRole %s successful"
	updateFailedMessage          string = "update CoherenceInternal %s from CoherenceRole %s failed\n%s"
	scaleToZeroFailed            string = "scale of CoherenceRole %s to zero failed\n%s"
	failedToGetHelmValuesMessage string = "failed to get Helm values for CoherenceRole %s due to error\n%s"
	failedToGetParentCluster     string = "failed to get parent CoherenceCluster %s for CoherenceRole %s due to error\n%s"
	failedToReconcileRole        string = "failed to reconcile CoherenceRole %s due to error\n%s"
	failedToScaleRole            string = "failed to scale CoherenceRole %s from %d to %d due to error\n%s"

	eventReasonFailed  string = "failed"
	eventReasonCreated string = "SuccessfulCreate"
	eventReasonUpdated string = "SuccessfulUpdate"
	eventReasonScale   string = "Scaling"

	// The template used to create the CoherenceRole.Status.Selector
	selectorTemplate = "coherenceCluster=%s,coherenceRole=%s"

	statusHaRetryEnv = "STATUS_HA_RETRY"

	// The name of the Coherence container in the Coherence Pods
	CoherenceContainerName = "coherence"
	// The name of the Coherence Utils container in the Coherence Pods
	CoherenceUtilsContainerName = "coherence-k8s-utils"
)

var log = logf.Log.WithName(controllerName)

// Add creates a new CoherenceRole Controller and adds it to the Manager. The Manager will set fields on the Controller.
// and Start it when the Manager is Started.
func Add(mgr manager.Manager, opFlags *flags.CoherenceOperatorFlags) error {
	return add(mgr, newReconciler(mgr, opFlags))
}

// NewRoleReconciler returns a new reconcile.Reconciler.
func NewRoleReconciler(mgr manager.Manager, opFlags *flags.CoherenceOperatorFlags) *ReconcileCoherenceRole {
	return newReconciler(mgr, opFlags)
}

// newReconciler returns a new reconcile.Reconciler.
func newReconciler(mgr manager.Manager, opFlags *flags.CoherenceOperatorFlags) *ReconcileCoherenceRole {
	scheme := mgr.GetScheme()
	gvk := coh.GetCoherenceInternalGroupVersionKind(scheme)

	// Parse the StatusHA retry time from the
	retry := time.Minute
	s := os.Getenv(statusHaRetryEnv)
	if s != "" {
		d, err := time.ParseDuration(s)
		if err == nil {
			retry = d
		} else {
			fmt.Printf("The value of %s env-var is not a valid duration '%s' using default retry time", statusHaRetryEnv, s)
		}
	}

	return &ReconcileCoherenceRole{
		client:        mgr.GetClient(),
		scheme:        scheme,
		gvk:           gvk,
		events:        mgr.GetEventRecorderFor(controllerName),
		statusHARetry: retry,
		mgr:           mgr,
		resourceLocks: make(map[types.NamespacedName]bool),
		mutex:         sync.Mutex{},
		opFlags:       opFlags,
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler.
func add(mgr manager.Manager, r *ReconcileCoherenceRole) error {
	// Create a new controller
	c, err := controller.New(controllerName, mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource CoherenceRole
	err = c.Watch(&source.Kind{Type: &coh.CoherenceRole{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource - in this case we watch the StatefulSet created by the Helm chart
	err = c.Watch(&source.Kind{Type: &appsv1.StatefulSet{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    r.CreateEmptyHelmValues(nil),
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileCoherenceRole implements reconcile.Reconciler.
// If the reconcile.Reconciler API was to change then we'd get a compile error here.
var _ reconcile.Reconciler = &ReconcileCoherenceRole{}

// ReconcileCoherenceRole reconciles a CoherenceRole object and related CoherenceInternal
// Helm values resources.
type ReconcileCoherenceRole struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the api server
	client        client.Client
	scheme        *runtime.Scheme
	gvk           schema.GroupVersionKind
	events        record.EventRecorder
	statusHARetry time.Duration
	mgr           manager.Manager
	resourceLocks map[types.NamespacedName]bool
	mutex         sync.Mutex
	opFlags       *flags.CoherenceOperatorFlags
	initialized   bool
}

// Set the initialized flag for this controller.
func (r *ReconcileCoherenceRole) SetInitialized(i bool) {
	if r != nil {
		r.initialized = i
	}
}

// Attempt to lock the requested resource.
func (r *ReconcileCoherenceRole) lock(request reconcile.Request) bool {
	if r == nil {
		return false
	}
	r.mutex.Lock()
	defer r.mutex.Unlock()

	_, found := r.resourceLocks[request.NamespacedName]
	if found {
		log.Info("CoherenceRole " + request.Namespace + "/" + request.Name + " is locked")
		return false
	}

	r.resourceLocks[request.NamespacedName] = true
	log.Info("Acquired lock for CoherenceRole " + request.Namespace + "/" + request.Name)
	return true
}

// Unlock the requested resource
func (r *ReconcileCoherenceRole) unlock(request reconcile.Request) {
	if r != nil {
		r.mutex.Lock()
		defer r.mutex.Unlock()

		log.Info("Released lock for CoherenceRole " + request.Namespace + "/" + request.Name)
		delete(r.resourceLocks, request.NamespacedName)
	}
}

// Reconcile reads that state of a CoherenceRole object and makes changes based on the state read
// and what is in the CoherenceRole.Spec.
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileCoherenceRole) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	logger := log.WithValues("Namespace", request.Namespace, "Name", request.Name)

	if err := r.EnsureInitialized(logger); err != nil {
		return reconcile.Result{}, err
	}

	// Attempt to lock the requested resource. If the resource is locked then another
	// request for the same resource is already in progress so requeue this one.
	if ok := r.lock(request); !ok {
		logger.Info("CoherenceRole " + request.Namespace + "/" + request.Name + " is already locked, re-queuing request")
		return reconcile.Result{Requeue: true, RequeueAfter: 0}, nil
	}
	// Make sure that the request is unlocked when this method exits
	defer r.unlock(request)

	return r.reconcileInternal(request)
}

// Reconcile reads that state of a CoherenceRole object and makes changes based on the state read
// and what is in the CoherenceRole.Spec.
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileCoherenceRole) reconcileInternal(request reconcile.Request) (reconcile.Result, error) {
	logger := log.WithValues("Namespace", request.Namespace, "Name", request.Name)
	logger.Info("Reconciling CoherenceRole")

	// Fetch the CoherenceRole role
	role, found, err := r.getRole(request.Namespace, request.Name)
	if err != nil {
		// Error reading the object - requeue the request.
		// We can't call the error handler as we do not even have a role.
		// We log the error and do not requeue the request.
		logger.Error(err, "Error getting CoherenceRole to reconcile")
		return r.handleErrAndRequeue(err, nil, fmt.Sprintf(failedToReconcileRole, role.Name, err), logger)
	}

	if !found || role.GetDeletionTimestamp() != nil {
		logger.Info("CoherenceRole deleted")
		// Request object not found (could have been deleted after reconcile request) or this is a delete notification.
		// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
		// Ensure that the CoherenceInternal for this role is deleted.
		// It should be cleaned-up by k8s garbage collection but belt and braces as we've seen when
		// upgrading or reinstalling the Operator existing clusters are not always cleaned automatically
		r.deleteCoherenceInternal(request, logger)

		return reconcile.Result{Requeue: false}, nil
	}

	clusterName := role.GetCoherenceClusterName()

	// Fetch the owning CoherenceCluster
	cluster := &coh.CoherenceCluster{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Namespace: request.Namespace, Name: clusterName}, cluster)
	if err != nil {
		if errors.IsNotFound(err) {
			return r.handleErrAndFinish(nil, role, fmt.Sprintf(invalidRoleEventMessage, role.Name, clusterName), logger)
		}
		return r.handleErrAndRequeue(err, role, fmt.Sprintf(failedToGetParentCluster, clusterName, role.Name, err.Error()), logger)
	}

	// find the existing Helm values structure in k8s (this will be an unstructured.Unstructured)
	// it may not exist if this is a create request
	helmValues, err := r.GetExistingHelmValues(role)
	replicas := role.Spec.GetReplicas()

	switch {
	case replicas <= 0 && err == nil && helmValues.GetDeletionTimestamp() == nil:
		// Scaling down to zero so delete the helm values which will cause the StatefulSet to be deleted
		return r.scaleDownToZero(cluster, role, helmValues)
	case replicas <= 0 && err == nil && helmValues.GetDeletionTimestamp() != nil:
		// Helm values already deleted but not yet gone
		if err = r.updateStatus(role, nil, cluster); err != nil {
			// failed to update the CoherenceRole's status
			log.Error(err, "failed to update role status", "Namespace", role.Namespace, "Name", role.Name)
			return reconcile.Result{Requeue: true, RequeueAfter: time.Second * 5}, nil
		}
		return reconcile.Result{}, nil
	case replicas <= 0 && err != nil && errors.IsNotFound(err):
		// Helm values has been deleted
		if err = r.updateStatus(role, nil, cluster); err != nil {
			// failed to update the CoherenceRole's status
			log.Error(err, "failed to update role status", "Namespace", role.Namespace, "Name", role.Name)
			return reconcile.Result{Requeue: true, RequeueAfter: time.Second * 5}, nil
		}
		return reconcile.Result{}, nil
	case replicas > 0 && err != nil && errors.IsNotFound(err):
		// Helm values was not found so this is an insert of a new role
		return r.createRole(cluster, role)
	case err != nil:
		// the error is a real error
		return r.handleErrAndRequeue(err, role, fmt.Sprintf(failedToGetHelmValuesMessage, role.Name, err.Error()), logger)
	default:
		// The Helm values was found so this is an update
		return r.updateRole(cluster, role, helmValues)
	}
}

// createRole creates a new Helm values structure in k8s, which will in turn trigger a Helm install.
func (r *ReconcileCoherenceRole) getRole(namespace, name string) (*coh.CoherenceRole, bool, error) {
	role := &coh.CoherenceRole{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Namespace: namespace, Name: name}, role)

	switch {
	case err != nil && errors.IsNotFound(err):
		return role, false, nil
	case err != nil:
		return role, false, err
	default:
		return role, true, nil
	}
}

// createRole creates a new Helm values structure in k8s, which will in turn trigger a Helm install.
func (r *ReconcileCoherenceRole) createRole(cluster *coh.CoherenceCluster, role *coh.CoherenceRole) (reconcile.Result, error) {
	if role.Spec.GetReplicas() <= 0 {
		// nothing to do as the desired replica count is zero
		return reconcile.Result{}, nil
	}

	logger := log.WithValues("Namespace", role.Namespace, "Name", role.Name)
	logger.Info("Creating Coherence Role Helm values")

	// create the CoherenceInternal for the role
	ci := coh.NewCoherenceInternalSpec(cluster, role)
	// Ensure that the CoherenceInternal has images set
	// If the CoherenceInternalSpec does not have a Coherence or Coherence Utils images specified we set the defaults here.
	// This ensures that the image is fixed to either that specified in the cluster spec or to the current default
	// and means that the Helm controller does not do a rolling upgrade of the Pods if the Operator is upgraded.
	r.EnsureImages(ci, logger)

	// define a new Helm values map
	spec, err := coh.CoherenceInternalSpecAsMapFromSpec(ci)
	if err != nil {
		// this error would only occur if there was a json marshalling issue which would be unlikely
		return r.handleErrAndRequeue(err, nil, fmt.Sprintf(createFailedMessage, role.Name, role.Name, err), logger)
	}

	helmValues := r.CreateHelmValues(cluster, role, spec)
	labels := helmValues.GetLabels()
	labels[coh.CoherenceOperatorVersionLabel] = role.GetLabels()[coh.CoherenceOperatorVersionLabel]
	helmValues.SetLabels(labels)

	// Set this CoherenceRole instance as the owner and controller of the Helm values structure
	if err := controllerutil.SetControllerReference(role, helmValues, r.scheme); err != nil {
		return r.handleErrAndRequeue(err, nil, fmt.Sprintf(createFailedMessage, helmValues.GetName(), role.Name, err), logger)
	}

	// Clean-up previous Helm v3 state if any exists, we've seen orphaned state causing issues in testing
	hsList := corev1.SecretList{}
	err = r.client.List(context.TODO(), &hsList, client.InNamespace(role.Namespace))
	if err == nil {
		prefix := fmt.Sprintf("sh.helm.release.v1.%s.", role.Name)
		for _, hs := range hsList.Items {
			if strings.HasPrefix(hs.Name, prefix) {
				// Helm state exists so delete it
				logger.Info(fmt.Sprintf("Deleting existing Helm state for role %s in secret %s", role.Name, hs.Name))
				_ = r.client.Delete(context.TODO(), &hs)
			}
		}
	}

	// Create the CoherenceInternal resource in k8s which will be detected
	// by the Helm operator and trigger a Helm install

	d, _ := json.Marshal(helmValues)
	logger.V(2).Info("Creating CoherenceInternal\n------------\n" + string(d) + "\n------------\n")

	if err := r.client.Create(context.TODO(), helmValues); err != nil {
		return r.handleErrAndRequeue(err, nil, fmt.Sprintf(createFailedMessage, helmValues.GetName(), role.Name, err), logger)
	}

	// update this CoherenceRole's status
	role.Status.Status = coh.RoleStatusCreated
	role.Status.Replicas = role.Spec.GetReplicas()
	role.Status.Selector = fmt.Sprintf(selectorTemplate, cluster.Name, role.Spec.GetRoleName())
	err = r.client.Status().Update(context.TODO(), role)
	if err != nil {
		// failed to update the CoherenceRole's status
		// ToDo - handle this properly by re-queuing the request and then in the reconcile method properly handle setting status even if the role is in the desired state
		logger.Error(err, "failed to update role status")
	}

	// send a successful creation event
	msg := fmt.Sprintf(createMessage, helmValues.GetName(), role.Name)
	r.events.Event(role, corev1.EventTypeNormal, eventReasonCreated, msg)

	return reconcile.Result{Requeue: false}, nil
}

// updateRole updates an existing CoherenceInternal which will in turn trigger a Helm update.
func (r *ReconcileCoherenceRole) updateRole(cluster *coh.CoherenceCluster, role *coh.CoherenceRole, helmValues *unstructured.Unstructured) (reconcile.Result, error) {
	logger := log.WithValues("Namespace", role.Namespace, "Name", role.Name)
	logger.Info("Reconciling existing Coherence Role")

	clusterRole := cluster.GetRole(role.Spec.GetRoleName())
	clusterReplicas := clusterRole.GetReplicas()
	roleReplicas := role.Spec.GetReplicas()
	// Get the effective role - what the role spec should be according to the Cluster spec
	effectiveRole := clusterRole.DeepCopyWithDefaults(&cluster.Spec.CoherenceRoleSpec)
	effectiveRole.SetReplicas(effectiveRole.GetReplicas())

	if !reflect.DeepEqual(effectiveRole, &role.Spec) {
		// Role spec is not the same as the cluster's role spec - likely caused by a scale but could have
		// been caused by a direct update to the CoherenceRole, even though we really discourage that.

		if clusterReplicas == roleReplicas {
			// Something other than the Replicas has been changed so reset the role back to what is should be.
			diff := deep.Equal(effectiveRole, &role.Spec)
			logger.Info("CoherenceCluster role spec is different to CoherenceRole spec and will be reset to match the cluster - diff:\n" + strings.Join(diff, "\n"))
			effectiveRole.DeepCopyInto(&role.Spec)
		} else {
			// Update the cluster's Replicas count to match the role, which will cause this update to come around again.
			clusterRole.SetReplicas(roleReplicas)
			cluster.SetRole(clusterRole)
			logger.Info(fmt.Sprintf("Reconciling existing Coherence Role - updating cluster's role Replicas from %d to %d", clusterReplicas, clusterRole.GetReplicas()))
			err := r.client.Update(context.TODO(), cluster)
			if err != nil {
				return r.handleErrAndRequeue(err, nil, fmt.Sprintf(updateFailedMessage, helmValues.GetName(), role.Name, err), logger)
			}
		}
	}

	desiredReplicas := role.Spec.GetReplicas()

	// convert the unstructured data to a CoherenceInternal that we can deal with better
	existing, err := r.toCoherenceInternal(helmValues)
	if err != nil {
		return r.handleErrAndRequeue(err, nil, fmt.Sprintf(updateFailedMessage, helmValues.GetName(), role.Name, err), logger)
	}

	sts, err := r.findStatefulSet(role)
	if err != nil && !errors.IsNotFound(err) {
		return r.handleErrAndRequeue(err, role, fmt.Sprintf(updateFailedMessage, helmValues.GetName(), role.Name, err), logger)
	}

	if errors.IsNotFound(err) {
		// The StatefulSet is not found, it could be being created or recovered
		logger.Info(fmt.Sprintf("Re-queing update request. Could not find StatefulSet for CoherenceRole %s", role.Name))
		return reconcile.Result{Requeue: true, RequeueAfter: time.Second * 10}, nil
	}

	currentReplicas := existing.Spec.GetReplicas()
	desiredRole := r.CreateDesiredRole(cluster, role, &existing.Spec, sts)
	isUpgrade := r.isUpgrade(&existing.Spec, desiredRole)

	switch {
	case currentReplicas < desiredReplicas:
		logger.Info("Reconciling existing Coherence Role: case currentReplicas < desiredReplicas")
		// Scaling UP

		// if scaling up and upgrading then upgrade first and scale second
		// otherwise we'd have to upgrade all the scaled up members
		if isUpgrade {
			err := r.upgrade(role, helmValues, currentReplicas, desiredRole)
			if err == nil {
				// Requeue so that we then scale up after the upgrade.
				// We do things this way because the upgrade is still happening asynchronously
				// and could take time. By re-queuing the request it will come back around to
				// the reconcile method keep being re-queued until the cluster is once again
				// in a stable state when the scale up will then happen.
				return reconcile.Result{Requeue: true}, nil
			}
			return r.handleErrAndRequeue(err, nil, fmt.Sprintf(updateFailedMessage, helmValues.GetName(), role.Name, err), logger)
		}

		logger.Info(fmt.Sprintf("Request to scale up from %d to %d", currentReplicas, desiredReplicas))
		return r.scale(role, helmValues, existing, desiredReplicas, currentReplicas, sts)
	case currentReplicas > desiredReplicas:
		logger.Info("Reconciling existing Coherence Role: case currentReplicas > desiredReplicas")
		// Scaling DOWN

		// if scaling down and upgrading then scale down first and upgrade second
		// so that we do not have to upgrade the members we are scaling down
		logger.Info(fmt.Sprintf("Request to scale down from %d to %d", currentReplicas, desiredReplicas))
		result, err := r.scale(role, helmValues, existing, desiredReplicas, currentReplicas, sts)

		if err == nil && isUpgrade {
			// requeue the request so that we then upgrade
			// We do things this way because the scale down is still happening asynchronously
			// and could take time. By re-queuing the request it will come back around to
			// the reconcile method keep being re-queued until the cluster is once again
			// in a stable state when the upgrade will then happen.
			return reconcile.Result{Requeue: true}, nil
		}
		return result, err
	case isUpgrade:
		logger.Info("Reconciling existing Coherence Role: case isUpgrade")
		// no scaling, just a rolling upgrade
		if err := r.upgrade(role, helmValues, currentReplicas, desiredRole); err != nil {
			return r.handleErrAndRequeue(err, nil, fmt.Sprintf(updateFailedMessage, helmValues.GetName(), role.Name, err), logger)
		}
		return reconcile.Result{Requeue: false}, nil
	case sts != nil:
		logger.Info("Reconciling existing Coherence Role: case sts != nil")
		// nothing to do to update or scale
		// We probably arrived here due to a change in the StatefulSet for a role
		// In this case we can potentially update the role's status based on what changed in the StatefulSet
		err = r.updateStatus(role, sts, cluster)
		if err != nil {
			// failed to update the CoherenceRole's status
			log.Error(err, "failed to update role status", "Namespace", role.Namespace, "Name", role.Name)
			return reconcile.Result{Requeue: true, RequeueAfter: time.Second * 5}, nil
		}
	}

	logger.Info("Finished reconciling existing Coherence Role")
	return reconcile.Result{Requeue: false}, nil
}

// scaleDownToZero is called in response to the replica count of a role being set to zero.
func (r *ReconcileCoherenceRole) scaleDownToZero(cluster *coh.CoherenceCluster, role *coh.CoherenceRole, cohInt *unstructured.Unstructured) (reconcile.Result, error) {
	logger := log.WithValues("Namespace", role.Namespace, "Name", role.Name)
	logger.Info("Scaling existing Coherence Role to zero")

	// Delete the CoherenceInternal causing a Helm delete of the Pods
	if err := r.client.Delete(context.TODO(), cohInt); err != nil {
		return r.handleErrAndRequeue(err, role, fmt.Sprintf(scaleToZeroFailed, role.Name, err), logger)
	}

	// Update the role in the parent CoherenceCluster to have zero replicas.
	clusterRole := cluster.GetRole(role.Spec.GetRoleName())
	clusterRole.SetReplicas(0)
	cluster.SetRole(role.Spec)
	logger.Info(fmt.Sprintf("Updating replica count for role %s in parent cluster %s to zero", role.Name, cluster.Name))
	if err := r.client.Update(context.TODO(), cluster); err != nil {
		return r.handleErrAndRequeue(err, role, fmt.Sprintf(scaleToZeroFailed, role.Name, err), logger)
	}

	// send a successful update event
	msg := fmt.Sprintf(updateMessage, role.Name, role.Name)
	r.events.Event(role, corev1.EventTypeNormal, eventReasonUpdated, msg)

	return reconcile.Result{Requeue: false}, nil
}

// isUpgrade determines whether the current spec differs to the desired spec ignoring differences to the Replicas field.
func (r *ReconcileCoherenceRole) isUpgrade(current *coh.CoherenceInternalSpec, desired *coh.CoherenceInternalSpec) bool {
	clone := current.DeepCopy()
	clone.Replicas = desired.Replicas

	return !reflect.DeepEqual(clone, desired)
}

// upgrade triggers a rolling upgrade of the role
func (r *ReconcileCoherenceRole) upgrade(role *coh.CoherenceRole, existingRole *unstructured.Unstructured, replicas int32, desiredRole *coh.CoherenceInternalSpec) error {
	// Rolling upgrade
	reqLogger := log.WithValues("Namespace", role.Namespace, "Name", role.Name)
	reqLogger.Info("Rolling upgrade of existing Role")

	spec, err := coh.CoherenceInternalSpecAsMapFromSpec(desiredRole)
	if err != nil {
		return err
	}

	// update the CoherenceInternal, this should trigger an update of the Helm install
	desiredRole.Replicas = &replicas
	existingRole.Object["spec"] = spec

	if err = r.client.Update(context.TODO(), existingRole); err != nil {
		return err
	}

	// Update this CoherenceRole's status
	role.Status.Status = coh.RoleStatusRollingUpgrade
	if err = r.client.Update(context.TODO(), role); err != nil {
		reqLogger.Error(err, "failed to update Status")
	}

	// send a successful update event
	msg := fmt.Sprintf(updateMessage, role.Name, role.Name)
	r.events.Event(role, corev1.EventTypeNormal, eventReasonUpdated, msg)

	return nil
}

// Update the role's status based on the status of the StatefulSet.
func (r *ReconcileCoherenceRole) updateStatus(role *coh.CoherenceRole, sts *appsv1.StatefulSet, cluster *coh.CoherenceCluster) error {
	var err error

	log.Info("Updating role status", "Namespace", role.Namespace, "Name", role.Name, "Cluster", cluster.Name)

	if sts == nil {
		role.Status.CurrentReplicas = 0
		role.Status.ReadyReplicas = 0

		if err = r.client.Status().Update(context.TODO(), role); err != nil {
			// failed to update the CoherenceRole's status
			log.Error(err, "failed to update role status", "Namespace", role.Namespace, "Name", role.Name, "Cluster", cluster.Name)
			return err
		}
	} else if role.Status.CurrentReplicas != sts.Status.Replicas || role.Status.ReadyReplicas != sts.Status.ReadyReplicas {
		// Update this CoherenceRole's status
		role.Status.CurrentReplicas = sts.Status.CurrentReplicas
		role.Status.ReadyReplicas = sts.Status.ReadyReplicas

		if sts.Status.ReadyReplicas == role.Spec.GetReplicas() {
			role.Status.Status = coh.RoleStatusReady
		}

		if err = r.client.Status().Update(context.TODO(), role); err != nil {
			// failed to update the CoherenceRole's status
			log.Error(err, "failed to update role status", "Namespace", role.Namespace, "Name", role.Name, "Cluster", cluster.Name)
			return err
		}
	}

	// Update this role's status in the parent cluster
	// re-fetch the cluster in case it has been updated
	clusterToUpdate := &coh.CoherenceCluster{}
	if err := r.client.Get(context.TODO(), types.NamespacedName{Namespace: cluster.Namespace, Name: cluster.Name}, clusterToUpdate); err != nil {
		log.Error(err, "failed to get parent cluster to update status", "Namespace", cluster.Namespace, "Name", cluster.Name)
		return err
	}

	ready := role.Status.Status == coh.RoleStatusReady
	clusterToUpdate.SetRoleStatus(role.Spec.Role, ready, role.Status.ReadyReplicas, role.Status.Status)

	log.Info("Updating role's status in parent cluster", "Namespace", role.Namespace, "Name", role.Name, "Cluster", cluster.Name)

	if err = r.client.Status().Update(context.TODO(), clusterToUpdate); err != nil {
		// failed to update the CoherenceCluster's status
		log.Error(err, "failed to update role's parent cluster status", "Namespace", role.Namespace, "Name", role.Name, "Cluster", cluster.Name)
		return err
	}

	return err
}

// findStatefulSet finds the StatefulSet associated to the role.
func (r *ReconcileCoherenceRole) findStatefulSet(role *coh.CoherenceRole) (*appsv1.StatefulSet, error) {
	sts := &appsv1.StatefulSet{}
	if err := r.client.Get(context.TODO(), types.NamespacedName{Name: role.Name, Namespace: role.Namespace}, sts); err != nil {
		return nil, err
	}

	return sts, nil
}

// GetExistingHelmValues gets an existing unstructured Helm values from k8s for a given CoherenceRole
func (r *ReconcileCoherenceRole) GetExistingHelmValues(role *coh.CoherenceRole) (*unstructured.Unstructured, error) {
	cohInt := r.CreateEmptyHelmValues(role)
	err := r.client.Get(context.TODO(), types.NamespacedName{Namespace: role.Namespace, Name: role.Name}, cohInt)
	return cohInt, err
}

// CreateHelmValues creates an unstructured Helm values struct.
func (r *ReconcileCoherenceRole) CreateHelmValues(cluster *coh.CoherenceCluster, role *coh.CoherenceRole, spec map[string]interface{}) *unstructured.Unstructured {
	cohInternal := r.CreateEmptyHelmValues(role)
	cohInternal.Object["spec"] = spec

	// Set the labels for the CoherenceInternal
	labels := make(map[string]string)
	labels[coh.CoherenceClusterLabel] = cluster.Name
	labels[coh.CoherenceRoleLabel] = role.Spec.GetRoleName()
	cohInternal.SetLabels(labels)

	return cohInternal
}

// CreateEmptyHelmValues creates an empty (no Spec) unstructured Helm values.
func (r *ReconcileCoherenceRole) CreateEmptyHelmValues(role *coh.CoherenceRole) *unstructured.Unstructured {
	cohInternal := &unstructured.Unstructured{}

	cohInternal.SetGroupVersionKind(r.gvk)

	if role != nil {
		cohInternal.SetNamespace(role.Namespace)
		cohInternal.SetName(role.Name)
	}

	return cohInternal
}

// toCoherenceInternal converts an unstructured Helm values struct to a real CoherenceInternal struct.
func (r *ReconcileCoherenceRole) toCoherenceInternal(state *unstructured.Unstructured) (*coh.CoherenceInternal, error) {
	b, err := state.MarshalJSON()
	if err != nil {
		return nil, err
	}

	cohInternal := &coh.CoherenceInternal{}
	gvk := &schema.GroupVersionKind{
		Group:   r.gvk.Group,
		Kind:    r.gvk.Kind,
		Version: r.gvk.Version,
	}

	_, _, err = unstructured.UnstructuredJSONScheme.Decode(b, gvk, cohInternal)
	if err != nil {
		return nil, err
	}

	return cohInternal, nil
}

// handleErrAndRequeue is the common error handler
func (r *ReconcileCoherenceRole) handleErrAndRequeue(err error, role *coh.CoherenceRole, msg string, logger logr.Logger) (reconcile.Result, error) {
	return r.failed(err, role, msg, true, logger)
}

// handleErrAndFinish is the common error handler
func (r *ReconcileCoherenceRole) handleErrAndFinish(err error, role *coh.CoherenceRole, msg string, logger logr.Logger) (reconcile.Result, error) {
	return r.failed(err, role, msg, false, logger)
}

// failed is the common error handler
// ToDo: we need to be able to add some form of back-off so that failures are re-queued with a growing back-off time
func (r *ReconcileCoherenceRole) failed(err error, role *coh.CoherenceRole, msg string, requeue bool, logger logr.Logger) (reconcile.Result, error) {
	if err == nil {
		logger.V(0).Info(msg)
	} else {
		logger.Error(err, msg)
	}

	if role != nil {
		// update the status to failed.
		role.Status.Status = coh.RoleStatusFailed
		if e := r.client.Status().Update(context.TODO(), role); e != nil {
			// There isn't much we can do, we're already handling an error
			logger.Error(err, "failed to update role status")
		}

		// send a failure event
		r.events.Event(role, corev1.EventTypeNormal, eventReasonFailed, msg)
	}

	if requeue {
		return reconcile.Result{Requeue: true}, nil
	}
	return reconcile.Result{Requeue: false}, nil
}

// If the CoherenceInternalSpec does not have a Coherence image specified we set the default here.
// This ensures that the image is fixed to either that specified in the cluster spec or to the current default
// and means that the Helm controller does not upgrade the images if the Operator is upgraded.
func (r *ReconcileCoherenceRole) EnsureImages(ci *coh.CoherenceInternalSpec, logger logr.Logger) {
	coherenceImage := r.opFlags.GetCoherenceImage()
	if ci.EnsureCoherenceImage(coherenceImage) {
		logger.Info(fmt.Sprintf("Injected Coherence image name into role: '%s'", *coherenceImage))
	}

	utilsImage := r.opFlags.GetCoherenceUtilsImage()
	if ci.EnsureCoherenceUtilsImage(utilsImage) {
		logger.Info(fmt.Sprintf("Injected Coherence Utils image name into role: '%s'", *utilsImage))
	}
}

// Create the desired CoherenceInternalSpec for a given role.
func (r *ReconcileCoherenceRole) CreateDesiredRole(cluster *coh.CoherenceCluster, role *coh.CoherenceRole, existing *coh.CoherenceInternalSpec, sts *appsv1.StatefulSet) *coh.CoherenceInternalSpec {
	desiredRole := coh.NewCoherenceInternalSpec(cluster, role)

	coherenceImage := existing.GetCoherenceImage()

	if sts != nil && desiredRole.GetCoherenceImage() == nil {
		// if the desired Coherence image is still nil then this could be an update to a cluster
		// started with a much older Operator so we'll obtain the current image from the StatefulSet
		for _, c := range sts.Spec.Template.Spec.Containers {
			if c.Name == CoherenceContainerName {
				coherenceImage = &c.Image
			}
		}
	}

	if coherenceImage == nil {
		// If the Coherence image is still nil then use the default
		coherenceImage = r.opFlags.GetCoherenceImage()
	}

	utilsImage := existing.GetCoherenceUtilsImage()

	if sts != nil && utilsImage == nil {
		// if the desired Coherence Utils image is still nil then this could be an update to a cluster
		// started with a much older Operator so we'll obtain the current image from the StatefulSet
		for _, c := range sts.Spec.Template.Spec.InitContainers {
			if c.Name == CoherenceUtilsContainerName {
				utilsImage = &c.Image
			}
		}
	}

	if utilsImage == nil {
		// If the utils image is still nil then use the default
		utilsImage = r.opFlags.GetCoherenceUtilsImage()
	}

	desiredRole.EnsureCoherenceImage(coherenceImage)
	desiredRole.EnsureCoherenceUtilsImage(utilsImage)

	return desiredRole
}

func (r *ReconcileCoherenceRole) deleteCoherenceInternal(request reconcile.Request, logger logr.Logger) {
	ci := &unstructured.Unstructured{}
	ci.SetGroupVersionKind(r.gvk)
	ci.SetNamespace(request.Namespace)
	ci.SetName(request.Name)

	logger.Info(fmt.Sprintf("Ensuring CoherenceInternal '%s/%s' is deleted", request.Namespace, request.Name))
	_ = r.client.Delete(context.TODO(), ci)
}

func (r *ReconcileCoherenceRole) EnsureInitialized(logger logr.Logger) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.initialized {
		return nil
	}

	logger.Info("Initializing controller")

	// Reconcile all of the existing clusters first
	if err := r.reconcileExistingRoles(); err != nil {
		return err
	}

	r.initialized = true

	// Start the Helm controller
	logger.Info("Starting Helm controller")
	return r.setupHelm(r.mgr)
}

// Produces pseudo reconcile requests for all of the existing CoherenceRoles in the watch namespace
// We typically call this function first before the Helm operator has stated to ensure that everything
// is in the required state before the Helm operator gets a chance to think it has to update installs.
func (r *ReconcileCoherenceRole) reconcileExistingRoles() error {
	log.Info("Reconciling all existing CoherenceRoles")

	namespace, err := k8sutil.GetWatchNamespace()
	if err != nil {
		return err
	}

	logger := log.WithValues("Namespace", namespace)

	list := coh.CoherenceRoleList{}
	if err := r.client.List(context.TODO(), &list, client.InNamespace(namespace)); err != nil {
		return err
	}

	if len(list.Items) == 0 {
		logger.Info("Zero existing CoherenceRoles to reconcile")
	}

	for _, role := range list.Items {
		logger.Info("Reconciling existing role", "Name", role.GetName())
		request := reconcile.Request{
			NamespacedName: types.NamespacedName{
				Namespace: role.GetNamespace(),
				Name:      role.GetName(),
			},
		}

		_, err = r.reconcileInternal(request)
		if err != nil {
			log.Error(err, "Error reconciling existing role", "Name", role.GetName())
		}
	}

	return nil
}

func (r *ReconcileCoherenceRole) setupHelm(mgr manager.Manager) error {
	namespace, err := k8sutil.GetWatchNamespace()
	if err != nil {
		return err
	}

	// Setup Helm controller
	watchList, err := watches.Load(r.opFlags.WatchesFile)
	if err != nil {
		log.Error(err, "failed to load Helm watches")
		return err
	}

	fmt.Println(watchList)
	for _, w := range watchList {
		fmt.Println(w)
		err := helmctl.Add(mgr, helmctl.WatchOptions{
			Namespace:               namespace,
			GVK:                     w.GroupVersionKind,
			ManagerFactory:          release.NewManagerFactory(mgr, w.ChartDir),
			ReconcilePeriod:         r.opFlags.ReconcilePeriod,
			WatchDependentResources: w.WatchDependentResources,
		})
		if err != nil {
			log.Error(err, "failed to add Helm watche")
			return err
		}
	}

	return nil
}
