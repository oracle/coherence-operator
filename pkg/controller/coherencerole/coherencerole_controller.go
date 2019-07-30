// Package coherencerole contains the Coherence Operator controller for the CoherenceRole crd
package coherencerole

import (
	"context"
	"fmt"
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	coherence "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// The name of this controller. This is used in events, log messages, etc.
const controllerName = "coherencerole-controller"

// The template used to create the CoherenceRole.Status.Selector
const selectorTemplate = "coherenceCluster=%s,coherenceRole=%s"

var log = logf.Log.WithName(controllerName)

// Add creates a new CoherenceRole Controller and adds it to the Manager. The Manager will set fields on the Controller.
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler.
func newReconciler(mgr manager.Manager) *ReconcileCoherenceRole {
	scheme := mgr.GetScheme()
	gvk := coherence.GetCoherenceInternalGroupVersionKind(scheme)

	return &ReconcileCoherenceRole{
		client: mgr.GetClient(),
		scheme: scheme,
		gvk:    gvk,
		events: mgr.GetRecorder(controllerName),
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
	err = c.Watch(&source.Kind{Type: &coherence.CoherenceRole{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource - in this case we watch the StatefulSet created by the Helm chart
	err = c.Watch(&source.Kind{Type: &appsv1.StatefulSet{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    r.createEmptyCoherenceInternal(nil),
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileCoherenceRole implements reconcile.Reconciler.
var _ reconcile.Reconciler = &ReconcileCoherenceRole{}

// ReconcileCoherenceRole reconciles a CoherenceRole object
type ReconcileCoherenceRole struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the api server
	client client.Client
	scheme *runtime.Scheme
	gvk    schema.GroupVersionKind
	events record.EventRecorder
}

// Reconcile reads that state of a CoherenceRole object and makes changes based on the state read
// and what is in the CoherenceRole.Spec.
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileCoherenceRole) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling CoherenceRole")

	// Fetch the CoherenceRole role
	role := &coherence.CoherenceRole{}
	err := r.client.Get(context.TODO(), request.NamespacedName, role)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	clusterName := role.GetCoherenceClusterName()

	// Fetch the owning CoherenceCluster
	cluster := &coherence.CoherenceCluster{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Namespace: request.Namespace, Name: clusterName}, cluster)
	if err != nil {
		if errors.IsNotFound(err) {
			reqLogger.Error(err, "A CoherenceRole must have an associated CoherenceCluster and should not be created outside of a CoherenceCluster.")

			// update the status to failed.
			role.Status.Status = coherence.RoleStatusFailed
			_ = r.client.Status().Update(context.TODO(), role)

			// send a successful creation event
			msg := fmt.Sprintf("Invalid CoherenceRole '%s' cannot find parent CoherenceCluster '%s'", role.Name, clusterName)
			r.events.Event(role, corev1.EventTypeNormal, "Failed", msg)

			return reconcile.Result{Requeue: false}, nil
		} else {
			return reconcile.Result{}, err
		}
	}

	// find the existing CoherenceInternal
	existingRole, err := r.getCoherenceInternal(role)
	if err != nil && !errors.IsNotFound(err) {
		return reconcile.Result{}, err
	}

	p := params{
		request:     request,
		cluster:     cluster,
		role:        role,
		cohInternal: existingRole,
		reqLogger:   reqLogger,
	}

	if err != nil && errors.IsNotFound(err) {
		// this is an insert of a new role
		return r.createRole(p)
	} else {
		return r.updateRole(p)
	}
}

// createRole creates a new CoherenceInternal which will in turn trigger a Helm install.
func (r *ReconcileCoherenceRole) createRole(p params) (reconcile.Result, error) {
	if p.role.Spec.GetReplicas() <= 0 {
		// nothing to do as the desired replica count is zero
		return reconcile.Result{}, nil
	}

	logger := p.reqLogger.WithValues("Role", p.cohInternal.GetName())
	logger.Info("Creating CoherenceInternal")

	// define a new CoherenceInternal structure
	spec, err := coherence.NewCoherenceInternalSpecAsMap(p.cluster, p.role)
	if err != nil {
		return reconcile.Result{}, err
	}
	cohInternal := r.createEmptyCoherenceInternal(p.role)
	cohInternal.Object["spec"] = spec

	// Set the labels for the CoherenceInternal
	labels := make(map[string]string)
	labels[coherence.CoherenceClusterLabel] = p.cluster.Name
	labels[coherence.CoherenceRoleLabel] = p.role.Spec.GetRoleName()
	cohInternal.SetLabels(labels)

	// Set CoherenceCluster instance as the owner and controller of the CoherenceInternal structure
	if err := controllerutil.SetControllerReference(p.cluster, cohInternal, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	p.role.Status.Status = coherence.RoleStatusCreated
	p.role.Status.Replicas = p.role.Spec.GetReplicas()
	p.role.Status.Selector = fmt.Sprintf(selectorTemplate, p.cluster.Name, p.role.Spec.GetRoleName())
	_ = r.client.Status().Update(context.TODO(), p.role)

	// Create the CoherenceInternal resource in k8s which will be detected
	// by the Helm operator and trigger a Helm install
	if err := r.client.Create(context.TODO(), cohInternal); err != nil {
		return reconcile.Result{}, err
	}

	// send a successful creation event
	msg := fmt.Sprintf("create Helm install '%s' in CoherenceRole '%s' successful", cohInternal.GetName(), p.role.Name)
	r.events.Event(p.role, corev1.EventTypeNormal, "SuccessfulCreate", msg)

	return reconcile.Result{}, nil
}

// updateRole updates an existing CoherenceInternal which will in turn trigger a Helm update.
func (r *ReconcileCoherenceRole) updateRole(p params) (reconcile.Result, error) {
	logger := p.reqLogger.WithValues("Role", p.cohInternal.GetName())
	logger.Info("Reconciling existing CoherenceRole")

	clusterRole := p.cluster.GetRole(p.role.Spec.GetRoleName())
	if !reflect.DeepEqual(clusterRole, p.role.Spec) {
		// role spec is not the same as the cluster's role spec - likely caused by a scale
		// update the cluster and the update will come around again
		logger.Info("CoherenceCluster role spec is different to this spec - updating CoherenceCluster '" + p.cluster.Name + "'")
		p.cluster.SetRole(p.role.Spec)
		err := r.client.Update(context.TODO(), p.cluster)
		if err != nil {
			return reconcile.Result{}, err
		}
	}

	// convert the unstructured data to a CoherenceInternal that we can deal with better
	existing, err := r.toCoherenceInternal(p.cohInternal)
	if err != nil {
		return reconcile.Result{}, err
	}

	currentReplicas := existing.Spec.ClusterSize
	desiredReplicas := p.role.Spec.GetReplicas()
	desiredRole := coherence.NewCoherenceInternalSpec(p.cluster, p.role)
	isUpgrade := !reflect.DeepEqual(existing.Spec, desiredRole)

	sts, err := r.findStatefulSet(p.role)
	if err != nil {
		logger.Info("Could not get StatefulSet")
	}

	if !isUpgrade && currentReplicas == desiredReplicas {
		// nothing to do
		logger.Info("Existing CoherenceRole is at the desired spec")
		return reconcile.Result{}, nil
	}

	if currentReplicas < desiredReplicas {
		// Scaling UP
		logger.Info(fmt.Sprintf("Scaling up existing Role from %d to %d", currentReplicas, desiredReplicas))
		// if scaling up and upgrading then upgrade first and scale second

		// update the CoherenceInternal, this should trigger an update of the Helm install to scale the StatefulSet
		existing.Spec.ClusterSize = desiredReplicas
		p.cohInternal.Object["spec"] = existing.Spec
		err = r.client.Update(context.TODO(), p.cohInternal)
		if err != nil {
			return reconcile.Result{}, err
		}

		// Update this CoherenceRole's status
		p.role.Status.Status = coherence.RoleStatusScaling
		p.role.Status.Replicas = desiredReplicas
		_ = r.client.Status().Update(context.TODO(), p.role)

		// send a successful scale event
		msg := fmt.Sprintf("scaled Helm install %s in CoherenceRole %s from %d to %d", p.role.Name, p.role.Name, currentReplicas, desiredReplicas)
		r.events.Event(p.role, corev1.EventTypeNormal, "SuccessfulScale", msg)

	} else if currentReplicas > desiredReplicas {
		// Scaling DOWN
		logger.Info(fmt.Sprintf("Scaling down existing Role from %d to %d", currentReplicas, desiredReplicas))
		// if scaling down and upgrading then scale down first and upgrade second

		// update the CoherenceInternal, this should trigger an update of the Helm install to scale the StatefulSet
		existing.Spec.ClusterSize = desiredReplicas
		p.cohInternal.Object["spec"] = existing.Spec
		err = r.client.Update(context.TODO(), p.cohInternal)
		if err != nil {
			return reconcile.Result{}, err
		}

		// Update this CoherenceRole's status
		p.role.Status.Status = coherence.RoleStatusScaling
		p.role.Status.Replicas = desiredReplicas
		err = r.client.Status().Update(context.TODO(), p.role)
		if err != nil {
			logger.Error(err, "Failed to update Status")
		}

		// send a successful scale event
		msg := fmt.Sprintf("scaled Helm install %s in CoherenceRole %s from %d to %d", p.role.Name, p.role.Name, currentReplicas, desiredReplicas)
		r.events.Event(p.role, corev1.EventTypeNormal, "SuccessfulScale", msg)

	} else if isUpgrade {
		// Rolling upgrade
		logger.Info("Rolling upgrade of existing Role")

		// update the CoherenceInternal, this should trigger an update of the Helm install
		p.cohInternal.Object["spec"] = desiredRole
		err = r.client.Status().Update(context.TODO(), p.cohInternal)
		if err != nil {
			return reconcile.Result{}, err
		}

		// Update this CoherenceRole's status
		p.role.Status.Status = coherence.RoleStatusRollingUpgrade
		p.role.Status.Replicas = desiredReplicas
		err = r.client.Update(context.TODO(), p.role)
		if err != nil {
			logger.Error(err, "Failed to update Status")
		}

		// send a successful scale event
		msg := fmt.Sprintf("Upgraded Helm install %s in CoherenceRole %s", p.role.Name, p.role.Name)
		r.events.Event(p.cluster, corev1.EventTypeNormal, "SuccessfulUpgrade", msg)
	} else if sts != nil {
		// nothing to do to update or scale - update our status if the StatefulSet has changed

		if p.role.Status.CurrentReplicas != sts.Status.Replicas || p.role.Status.ReadyReplicas != sts.Status.ReadyReplicas {
			// Update this CoherenceRole's status
			p.role.Status.CurrentReplicas = sts.Status.Replicas
			p.role.Status.ReadyReplicas = sts.Status.ReadyReplicas
			if sts.Status.ReadyReplicas == desiredReplicas {
				p.role.Status.Status = coherence.RoleStatusReady
			}
			_ = r.client.Status().Update(context.TODO(), p.role)
		}
	}

	return reconcile.Result{}, nil
}

// findStatefulSet finds the StatefulSet associated to the role.
func (r *ReconcileCoherenceRole) findStatefulSet(role *coherence.CoherenceRole) (*appsv1.StatefulSet, error) {
	sts := &appsv1.StatefulSet{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: role.Name, Namespace: role.Namespace}, sts)

	if err != nil {
		return nil, err
	}

	return sts, nil
}

// GetCoherenceInternal gets the unstructured CoherenceInternal from k8s for a given CoherenceRole
func (r *ReconcileCoherenceRole) getCoherenceInternal(role *coherence.CoherenceRole) (*unstructured.Unstructured, error) {
	cohInt := r.createEmptyCoherenceInternal(role)
	err := r.client.Get(context.TODO(), types.NamespacedName{Namespace: role.Namespace, Name: role.Name}, cohInt)
	return cohInt, err
}

// createEmptyCoherenceInternal creates an empty (no Spec) unstructured CoherenceInternal.
func (r *ReconcileCoherenceRole) createEmptyCoherenceInternal(role *coherence.CoherenceRole) *unstructured.Unstructured {
	cohInternal := &unstructured.Unstructured{}

	cohInternal.SetGroupVersionKind(r.gvk)

	if role != nil {
		cohInternal.SetNamespace(role.Namespace)
		cohInternal.SetName(role.Name)
	}

	return cohInternal
}

// toCoherenceInternal converts an unstructured CoherenceInternal to a real CoherenceInternal struct.
func (r *ReconcileCoherenceRole) toCoherenceInternal(state *unstructured.Unstructured) (*coherence.CoherenceInternal, error) {
	b, err := state.MarshalJSON()
	if err != nil {
		return nil, err
	}

	cohInternal := &coherence.CoherenceInternal{}
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

// params is the parameters to the insertRole and updateRole methods in a struct.
// This makes the method signatures a little more compact
type params struct {
	request     reconcile.Request
	cluster     *coherence.CoherenceCluster
	role        *coherence.CoherenceRole
	cohInternal *unstructured.Unstructured
	reqLogger   logr.Logger
}
