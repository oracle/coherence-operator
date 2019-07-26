// Package coherencecluster contains the Coherence Operator controller for the CoherenceCluster crd
package coherencecluster

import (
	"context"
	"fmt"
	"github.com/go-logr/logr"
	coherence "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	"github.com/oracle/coherence-operator/pkg/flags"
	v1 "k8s.io/api/core/v1"
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
)

// The name of this controller. This is used in events, log messages, etc.
const controllerName = "coherencecluster-controller"

var log = logf.Log.WithName(controllerName)

// Add creates a new CoherenceCluster Controller and adds it to the Manager.
// The Manager will set fields on the Controller.
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler.
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileCoherenceCluster{
		client: mgr.GetClient(),
		scheme: mgr.GetScheme(),
		events: mgr.GetRecorder(controllerName),
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler.
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New(controllerName, mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource CoherenceCluster
	err = c.Watch(&source.Kind{Type: &coherence.CoherenceCluster{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource CoherenceRole and requeue the owner CoherenceCluster
	err = c.Watch(&source.Kind{Type: &coherence.CoherenceRole{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &coherence.CoherenceCluster{},
	})
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
	client client.Client
	scheme *runtime.Scheme
	events record.EventRecorder
	flags  *flags.CoherenceOperatorFlags
}

// Reconcile reads that state of a CoherenceCluster object and makes changes based on the state read
// and what is in the CoherenceCluster.Spec.
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileCoherenceCluster) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Namespace", request.Namespace, "Cluster.Name", request.Name)
	reqLogger.Info("Reconciling CoherenceCluster")

	// Fetch the CoherenceCluster instance
	cluster := &coherence.CoherenceCluster{}
	err := r.client.Get(context.TODO(), request.NamespacedName, cluster)
	if err != nil {
		if errors.IsNotFound(err) {
			reqLogger.Info("CoherenceCluster '" + request.Name + "' deleted")
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

	// remove any existing roles that are not in the desired spec
	newRoleNames := listRoles(cluster)
	for name, role := range existingRoles {
		if newRoleNames[name] == false {
			err = r.deleteRole(params{existingRole: role, reqLogger: reqLogger})
			if err != nil {
				return reconcile.Result{}, err
			}
		}
	}

	// Remove any existing roles where the replica count in the desired spec is zero
	for _, role := range cluster.Spec.Roles {
		if role.GetReplicas() == 0 {
			existingRole, found := existingRoles[role.RoleName]
			if found {
				err = r.deleteRole(params{existingRole: existingRole, reqLogger: reqLogger})
				if err != nil {
					return reconcile.Result{}, err
				}
			}
		}
	}

	// Process the inserts and updates
	for _, role := range cluster.Spec.Roles {
		// Check whether this CoherenceInternal already exists
		existingRole, found := existingRoles[role.RoleName]

		if found {
			// this is a request to update a cluster role

			parameters := params{
				request:      request,
				cluster:      cluster,
				desiredRole:  role,
				existingRole: existingRole,
				reqLogger:    reqLogger,
			}

			result, err := r.updateRole(parameters)
			if err != nil || result.Requeue {
				return result, err
			}
		} else {
			// this is a request for a new cluster role

			// make sure that the WKA service exists
			if err := r.ensureWkaService(cluster); err != nil {
				return reconcile.Result{}, err
			}

			parameters := params{
				request:      request,
				cluster:      cluster,
				desiredRole:  role,
				existingRole: existingRole,
				reqLogger:    reqLogger,
			}

			if err := r.createRole(parameters); err != nil {
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
		// nothing to do as the desired replica count is zero
		return nil
	}

	logger := p.reqLogger.WithValues("Role", p.existingRole.GetName())
	logger.Info("Creating new Role")

	// define a new CoherenceRole structure
	role := &coherence.CoherenceRole{}
	role.SetNamespace(p.request.Namespace)
	role.SetName(p.request.Name + "-" + p.desiredRole.RoleName)
	role.Spec = p.desiredRole

	labels := make(map[string]string)
	labels["coherenceCluster"] = p.cluster.Name
	labels["coherenceRole"] = p.desiredRole.RoleName
	role.SetLabels(labels)

	// Set CoherenceCluster instance as the owner and controller of the CoherenceRole structure
	if err := controllerutil.SetControllerReference(p.cluster, role, r.scheme); err != nil {
		return err
	}

	// Create the CoherenceRole resource in k8s which will be detected by the role controller
	if err := r.client.Create(context.TODO(), role); err != nil {
		return err
	}

	// send a successful creation event
	msg := fmt.Sprintf("create CoherenceRole %s in CoherenceCluster %s successful", role.Name, p.cluster.Name)
	r.events.Event(p.cluster, v1.EventTypeNormal, "SuccessfulCreate", msg)

	return nil
}

// updateRole updates an existing cluster role.
func (r *ReconcileCoherenceCluster) updateRole(p params) (reconcile.Result, error) {
	logger := p.reqLogger.WithValues("Role", p.existingRole.GetName())
	logger.Info("Update existing Role")

	if reflect.DeepEqual(p.existingRole.Spec, p.desiredRole) {
		// nothing to do
		logger.Info("Existing Role is at the desired spec")
		return reconcile.Result{}, nil
	}

	// Create the CoherenceRole resource in k8s which will be detected by the role controller
	p.existingRole.Spec = p.desiredRole
	err := r.client.Update(context.TODO(), &p.existingRole)

	if err == nil {
		// send a successful update event
		msg := fmt.Sprintf("update CoherenceRole %s in CoherenceCluster %s successful", p.existingRole.Name, p.cluster.Name)
		r.events.Event(p.cluster, v1.EventTypeNormal, "SuccessfulUpdate", msg)
	} else {
		// send a failed update event
		msg := fmt.Sprintf("update CoherenceRole %s in CoherenceCluster %s failed\n%s", p.existingRole.Name, p.cluster.Name, err.Error())
		r.events.Event(p.cluster, v1.EventTypeNormal, "FailedUpdate", msg)
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
		return err
	}

	// send a successful deletion event
	msg := fmt.Sprintf("delete CoherenceRole %s in CoherenceCluster %s successful", p.existingRole.Name, p.cluster.Name)
	r.events.Event(p.cluster, v1.EventTypeNormal, "SuccessfulDelete", msg)

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
			Port:       20000,
			TargetPort: intstr.FromString("extend-port"),
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

// listRoles creates a map of all of the role names in the specified CoherenceCluster.
func listRoles(cluster *coherence.CoherenceCluster) map[string]bool {
	m := make(map[string]bool)
	for _, role := range cluster.Spec.Roles {
		m[role.RoleName] = true
	}
	return m
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
		roles[role.Spec.RoleName] = role
	}

	return nil
}
