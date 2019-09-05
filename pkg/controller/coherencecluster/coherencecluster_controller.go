/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

// Package coherencecluster contains the Coherence Operator controller for the CoherenceCluster crd
package coherencecluster

import (
	"context"
	"fmt"
	"github.com/go-logr/logr"
	"github.com/go-test/deep"
	coherence "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/tools/record"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"strings"
	"sync"
)

const (
	// The name of this controller. This is used in events, log messages, etc.
	controllerName string = "coherencecluster.controller"

	createEventMessage       string = "created CoherenceRole '%s' in CoherenceCluster '%s' successful"
	createEventFailedMessage string = "create CoherenceRole '%s' in CoherenceCluster '%s' failed\n%s"
	updateEventMessage       string = "updated CoherenceRole %s in CoherenceCluster %s successful"
	updateFailedEventMessage string = "update CoherenceRole %s in CoherenceCluster %s failed\n%s"
	deleteEventMessage       string = "deleted CoherenceRole %s in CoherenceCluster %s successful"
	deleteFailedEventMessage string = "delete CoherenceRole %s in CoherenceCluster %s failed\n%s"

	eventReasonCreated      string = "SuccessfulCreate"
	eventReasonFailedCreate string = "FailedCreate"
	eventReasonUpdated      string = "SuccessfulUpdate"
	eventReasonFailedUpdate string = "FailedUpdate"
	eventReasonDeleted      string = "SuccessfulDelete"
	eventReasonFailedDelete string = "FailedDelete"
)

var log = logf.Log.WithName(controllerName)

// Add creates a new CoherenceCluster Controller and adds it to the Manager.
// The Manager will set fields on the Controller.
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// NewClusterReconciler returns a new reconcile.Reconciler.
func NewClusterReconciler(mgr manager.Manager) reconcile.Reconciler {
	return newReconciler(mgr)
}

// newReconciler returns a new reconcile.Reconciler.
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileCoherenceCluster{
		client:        mgr.GetClient(),
		scheme:        mgr.GetScheme(),
		events:        mgr.GetRecorder(controllerName),
		resourceLocks: make(map[types.NamespacedName]bool),
		mutex:         sync.Mutex{},
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler.
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New(controllerName, mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to CoherenceCluster resources
	err = c.Watch(&source.Kind{Type: &coherence.CoherenceCluster{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileCoherenceCluster implements reconcile.Reconciler.
var _ reconcile.Reconciler = &ReconcileCoherenceCluster{}

// ReconcileCoherenceCluster reconciles a CoherenceCluster object
type ReconcileCoherenceCluster struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the api server
	client        client.Client
	scheme        *runtime.Scheme
	events        record.EventRecorder
	resourceLocks map[types.NamespacedName]bool
	mutex         sync.Mutex
}

// Attempt to lock the requested resource.
func (r *ReconcileCoherenceCluster) lock(request reconcile.Request) bool {
	if r == nil {
		return false
	}
	r.mutex.Lock()
	defer r.mutex.Unlock()

	_, found := r.resourceLocks[request.NamespacedName]
	if found {
		log.Info("CoherenceCluster " + request.Namespace + "/" + request.Name + " is locked")
		return false
	}

	r.resourceLocks[request.NamespacedName] = true
	log.Info("Acquired lock for CoherenceCluster " + request.Namespace + "/" + request.Name)
	return true
}

// Unlock the requested resource
func (r *ReconcileCoherenceCluster) unlock(request reconcile.Request) {
	if r != nil {
		r.mutex.Lock()
		defer r.mutex.Unlock()

		log.Info("Released lock for CoherenceCluster " + request.Namespace + "/" + request.Name)
		delete(r.resourceLocks, request.NamespacedName)
	}
}

// Reconcile reads that state of a CoherenceCluster object and makes changes based on the state read
// and what is in the CoherenceCluster.Spec.
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileCoherenceCluster) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	logger := log.WithValues("Namespace", request.Namespace, "Cluster.Name", request.Name)
	logger.Info("Reconciling CoherenceCluster")

	// Attempt to lock the requested resource. If the resource is locked then another
	// request for the same resource is already in progress so requeue this one.
	if ok := r.lock(request); !ok {
		logger.Info("CoherenceCluster " + request.Namespace + "/" + request.Name + " is already locked, re-queuing request")
		return reconcile.Result{Requeue: true, RequeueAfter: 0}, nil
	}
	// Make sure that the request is unlocked when this method exits
	defer r.unlock(request)

	// Fetch the CoherenceCluster instance
	cluster := &coherence.CoherenceCluster{}
	err := r.client.Get(context.TODO(), request.NamespacedName, cluster)
	if err != nil {
		if errors.IsNotFound(err) {
			logger.Info("CoherenceCluster '" + request.Name + "' deleted")
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	clusterName := cluster.GetName()

	existingRoles := make(map[string]coherence.CoherenceRole)
	if err = r.findExistingRoles(clusterName, cluster.Namespace, existingRoles); err != nil {
		return reconcile.Result{}, err
	}

	// get the desired role specs from the cluster
	desiredRoles, desiredRoleNames := r.getDesiredRoles(cluster)

	// remove any existing roles that are not in the desired spec
	for name, role := range existingRoles {
		if _, found := desiredRoles[name]; !found {
			err = r.deleteRole(params{cluster: cluster, existingRole: role, reqLogger: logger})
			if err != nil {
				return reconcile.Result{}, err
			}
		}
	}

	// Remove any existing roles where the replica count in the desired spec is zero
	for _, role := range desiredRoles {
		if role.GetReplicas() == 0 {
			existingRole, found := existingRoles[role.GetRoleName()]
			if found {
				err = r.deleteRole(params{cluster: cluster, existingRole: existingRole, reqLogger: logger})
				if err != nil {
					return reconcile.Result{}, err
				}
			}
		}
	}

	// Process the inserts and updates in the order they are specified in the cluster spec
	for _, roleName := range desiredRoleNames {
		role := desiredRoles[roleName]
		if role.GetReplicas() > 0 {
			// Check whether this CoherenceRole already exists
			existingRole, found := existingRoles[role.GetRoleName()]

			parameters := params{
				request:      request,
				cluster:      cluster,
				desiredRole:  role,
				existingRole: existingRole,
				reqLogger:    logger,
			}

			if found {
				// this is a request to update a cluster role
				result, err := r.updateRole(parameters)
				if err != nil || result.Requeue {
					return result, err
				}
			} else if err := r.createRole(parameters); err != nil {
				// this is a request for a new cluster role
				return reconcile.Result{}, err
			}
		}
	}

	// we're done so do not requeue
	return reconcile.Result{}, nil
}

// params contains the parameters to the insertRole and updateRole methods in a struct.
// This makes the method signatures a little more compact.
type params struct {
	request      reconcile.Request
	cluster      *coherence.CoherenceCluster
	desiredRole  coherence.CoherenceRoleSpec
	existingRole coherence.CoherenceRole
	reqLogger    logr.Logger
}

// createRole create a new cluster role.
func (r *ReconcileCoherenceCluster) createRole(p params) error {
	if p.desiredRole.GetReplicas() <= 0 {
		// should not get here but do nothing the desired replica count is zero
		return nil
	}

	fullName := p.desiredRole.GetFullRoleName(p.cluster)

	logger := p.reqLogger.WithValues("Role", fullName)
	logger.Info("Creating CoherenceRole")

	// make sure that the WKA service exists
	if err := r.ensureWkaService(p.cluster); err != nil {
		return err
	}

	// define a new CoherenceRole structure
	role := &coherence.CoherenceRole{}
	role.SetNamespace(p.request.Namespace)
	role.SetName(fullName)
	role.Spec = *p.desiredRole.DeepCopyWithDefaults(&p.cluster.Spec.CoherenceRoleSpec)

	labels := make(map[string]string)
	labels[coherence.CoherenceClusterLabel] = p.cluster.Name
	labels[coherence.CoherenceRoleLabel] = p.desiredRole.GetRoleName()
	role.SetLabels(labels)

	// Set CoherenceCluster instance as the owner and controller of the CoherenceRole structure
	if err := controllerutil.SetControllerReference(p.cluster, role, r.scheme); err != nil {
		return err
	}

	// Create the CoherenceRole resource in k8s which will be detected by the role controller
	if err := r.client.Create(context.TODO(), role); err != nil {
		msg := fmt.Sprintf(createEventFailedMessage, role.Name, p.cluster.Name, err.Error())
		r.events.Event(p.cluster, v1.EventTypeNormal, eventReasonFailedCreate, msg)
		return err
	}

	// send a successful creation event
	msg := fmt.Sprintf(createEventMessage, role.Name, p.cluster.Name)
	r.events.Event(p.cluster, v1.EventTypeNormal, eventReasonCreated, msg)

	return nil
}

// updateRole updates an existing cluster role.
func (r *ReconcileCoherenceCluster) updateRole(p params) (reconcile.Result, error) {
	logger := p.reqLogger.WithValues("Role", p.existingRole.GetName())

	if reflect.DeepEqual(p.existingRole.Spec, p.desiredRole) {
		// nothing to do
		logger.Info("Existing Role is at the desired spec")
		return reconcile.Result{}, nil
	}

	diff := deep.Equal(p.existingRole.Spec, p.desiredRole)
	logger.Info("Updating CoherenceRole - diff\n    " + strings.Join(diff, "\n    "+"\n"))

	// Create the CoherenceRole resource in k8s which will be detected by the role controller
	p.existingRole.Spec = p.desiredRole
	err := r.client.Update(context.TODO(), &p.existingRole)

	if err == nil {
		// send a successful update event
		msg := fmt.Sprintf(updateEventMessage, p.existingRole.Name, p.cluster.Name)
		r.events.Event(p.cluster, v1.EventTypeNormal, eventReasonUpdated, msg)
	} else {
		// send a failed update event
		msg := fmt.Sprintf(updateFailedEventMessage, p.existingRole.Name, p.cluster.Name, err.Error())
		r.events.Event(p.cluster, v1.EventTypeNormal, eventReasonFailedUpdate, msg)
	}

	return reconcile.Result{}, err
}

// deleteRole deletes an existing cluster role.
// This will ultimately trigger un-deployment of the related cluster members.
func (r *ReconcileCoherenceCluster) deleteRole(p params) error {
	logger := p.reqLogger.WithValues("Role", p.existingRole.GetName())
	logger.Info("Deleting existing Role")

	err := r.client.Delete(context.TODO(), &p.existingRole)
	if err != nil {
		msg := fmt.Sprintf(deleteFailedEventMessage, p.existingRole.Name, p.cluster.Name, err.Error())
		r.events.Event(p.cluster, v1.EventTypeNormal, eventReasonFailedDelete, msg)
		return err
	}

	// send a successful deletion event
	msg := fmt.Sprintf(deleteEventMessage, p.existingRole.Name, p.cluster.Name)
	r.events.Event(p.cluster, v1.EventTypeNormal, eventReasonDeleted, msg)

	return nil
}

// ensureWkaService ensures that the headless service used for WKA exists for the specified cluster.
// A service will be created if one does not exist
func (r *ReconcileCoherenceCluster) ensureWkaService(cluster *coherence.CoherenceCluster) error {

	name := types.NamespacedName{
		Namespace: cluster.Namespace,
		Name:      cluster.GetWkaServiceName(),
	}

	err := r.client.Get(context.TODO(), name, &v1.Service{})

	if err != nil && errors.IsNotFound(err) {
		reqLogger := log.WithValues("Namespace", cluster.Namespace, "Cluster.Name", cluster.Name)
		reqLogger.Info("Creating WKA service '" + name.Name + "'")

		service := &v1.Service{}

		service.Namespace = name.Namespace
		service.Name = name.Name

		service.Annotations = make(map[string]string)
		service.Annotations["service.alpha.kubernetes.io/tolerate-unready-endpoints"] = "true"

		service.Labels = make(map[string]string)
		service.Labels[coherence.CoherenceClusterLabel] = cluster.Name
		service.Labels[coherence.CoherenceComponentLabel] = "coherenceWkaService"

		service.Spec = v1.ServiceSpec{}
		service.Spec.ClusterIP = v1.ClusterIPNone
		service.Spec.Ports = make([]v1.ServicePort, 1)

		service.Spec.Ports[0] = v1.ServicePort{
			Name:       "coherence-extend",
			Protocol:   v1.ProtocolTCP,
			Port:       30000,
			TargetPort: intstr.FromString("mgmt-port"),
		}

		service.Spec.Selector = make(map[string]string)
		service.Spec.Selector[coherence.CoherenceClusterLabel] = cluster.Name
		service.Spec.Selector["component"] = "coherencePod"

		// Set CoherenceCluster instance as the owner and controller of the service structure
		if err := controllerutil.SetControllerReference(cluster, service, r.scheme); err != nil {
			return err
		}

		return r.client.Create(context.TODO(), service)
	}

	return err
}

// getDesiredRoles returns a map with all of the desired roles from the cluster and a slice of role names in the order they
// were specified in the cluster.
// If the cluster has no roles then the default role will be used as the single role spec for the cluster.
func (r *ReconcileCoherenceCluster) getDesiredRoles(cluster *coherence.CoherenceCluster) (map[string]coherence.CoherenceRoleSpec, []string) {
	defaultSpec := cluster.Spec.CoherenceRoleSpec
	if len(cluster.Spec.Roles) == 0 {
		return map[string]coherence.CoherenceRoleSpec{defaultSpec.GetRoleName(): defaultSpec}, []string{defaultSpec.GetRoleName()}
	}

	m := make(map[string]coherence.CoherenceRoleSpec)
	names := make([]string, len(cluster.Spec.Roles))
	index := 0
	for _, role := range cluster.Spec.Roles {
		clone := role.DeepCopyWithDefaults(&defaultSpec)
		names[index] = role.GetRoleName()
		m[names[index]] = *clone
		index++
	}
	return m, names
}

// findExistingRoles populates a map with all of the existing (deployed) cluster roles for the cluster name.
func (r *ReconcileCoherenceCluster) findExistingRoles(clusterName string, namespace string, roles map[string]coherence.CoherenceRole) error {
	list := &coherence.CoherenceRoleList{}

	opts := client.ListOptions{
		Namespace: namespace,
	}

	err := opts.SetLabelSelector(coherence.CoherenceClusterLabel + "=" + clusterName)
	if err != nil {
		return err
	}

	err = r.client.List(context.TODO(), &opts, list)
	if err != nil {
		return err
	}

	for _, role := range list.Items {
		roles[role.Spec.GetRoleName()] = role
	}

	return nil
}
