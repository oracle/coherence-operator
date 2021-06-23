/*
 * Copyright (c) 2020, 2021, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package controllers

import (
	"context"
	"fmt"
	monitoringv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/go-logr/logr"
	"github.com/go-test/deep"
	"github.com/oracle/coherence-operator/controllers/predicates"
	"github.com/oracle/coherence-operator/controllers/reconciler"
	"github.com/oracle/coherence-operator/controllers/servicemonitor"
	"github.com/oracle/coherence-operator/controllers/statefulset"
	"github.com/oracle/coherence-operator/pkg/operator"
	"github.com/oracle/coherence-operator/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
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

	coh "github.com/oracle/coherence-operator/api/v1"
)

const (
	controllerName = "controller_coherence"

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

// blank assignment to verify that CoherenceReconciler implements reconcile.Reconciler
// There will be a compile time error here if this breaks
var _ reconcile.Reconciler = &CoherenceReconciler{}

// +kubebuilder:rbac:groups=coherence.oracle.com,resources=coherence,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=coherence.oracle.com,resources=coherence/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=pods;pods/exec;services;endpoints;events;configmaps;secrets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=monitoring.coreos.com,resources=servicemonitors,verbs=get;list;watch;create;update;patch;delete

func (in *CoherenceReconciler) Reconcile(ctx context.Context, request ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	log := in.Log.WithValues("coherence", request.NamespacedName)

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
	err := in.GetClient().Get(ctx, types.NamespacedName{Namespace: request.Namespace, Name: request.Name}, deployment)
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
	if deployment.GetDeletionTimestamp() != nil {
		// Check whether finalization needs to be run
		if utils.StringArrayContains(deployment.GetFinalizers(), coh.Finalizer) {
			// Run finalization logic.
			// If the finalization logic fails, don't remove the finalizer so
			// that we can retry during the next reconciliation.
			if err := in.finalizeDeployment(ctx, deployment); err != nil {
				log.Error(err, "Failed to remove finalizer")
				return ctrl.Result{Requeue: true}, nil
			}
			// Remove the finalizer. Once all finalizers have been
			// removed, the object will be deleted.
			controllerutil.RemoveFinalizer(deployment, coh.Finalizer)
			err := in.GetClient().Update(ctx, deployment)
			if err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	// The request is an add or update

	// Add finalizer for this CR if required
	if utils.StringArrayDoesNotContain(deployment.GetFinalizers(), coh.Finalizer) {
		// Adding the finalizer causes an update so the request will come around again
		if err := in.addFinalizer(ctx, deployment); err != nil {
			return ctrl.Result{}, err
		}
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
		// No replica count so we patch the deployment to have the default replicas value.
		// The reason we do this is because the kubectl scale command will fail if the replicas
		// field is not set to a non-nil value.
		patch := &coh.Coherence{}
		deployment.DeepCopyInto(patch)
		replicas := deployment.GetReplicas()
		patch.Spec.Replicas = &replicas
		_, err = in.ThreeWayPatch(ctx, deployment.Name, deployment, deployment, patch)
		if err != nil {
			return reconcile.Result{}, err
		}
	}

	// ensure that the Operator configuration Secret exists
	if err = coh.EnsureOperatorSecret(ctx, request.Namespace, in.GetClient(), in.Log); err != nil {
		err = errors.Wrap(err, "ensuring Operator configuration secret")
		return in.HandleErrAndRequeue(ctx, err, nil, fmt.Sprintf(reconcileFailedMessage, request.Name, request.Namespace, err), in.Log)
	}

	// ensure that the state store exists
	storage, err := utils.NewStorage(request.NamespacedName, in.GetManager())
	if err != nil {
		err = errors.Wrap(err, "obtaining desired state store")
		return in.HandleErrAndRequeue(ctx, err, nil, fmt.Sprintf(reconcileFailedMessage, request.Name, request.Namespace, err), in.Log)
	}

	// obtain the original resources from the state store
	originalResources := storage.GetLatest()
	// Create the desired resources the deployment
	desiredResources, err := deployment.Spec.CreateKubernetesResources(deployment)
	if err != nil {
		return in.HandleErrAndRequeue(ctx, err, nil, fmt.Sprintf(createResourcesFailedMessage, request.Name, request.Namespace, err), in.Log)
	}

	// compare the original with the desired to determine whether there is anything to update
	var isDiff bool
	switch {
	case originalResources.Items == nil && desiredResources.Items == nil:
		isDiff = false
	case originalResources.Items != nil && desiredResources.Items == nil:
		isDiff = true
	case originalResources.Items == nil && desiredResources.Items != nil:
		isDiff = true
	default:
		diff := deep.Equal(originalResources.Items, desiredResources.Items)
		isDiff = len(diff) > 0
	}

	// create the result
	result := ctrl.Result{Requeue: false}

	if isDiff {
		in.Log.Info("Reconciling Coherence resource secondary resources")
		// make the deployment the owner of all of the secondary resources about to be reconciled
		if err := desiredResources.SetController(deployment, in.GetManager().GetScheme()); err != nil {
			return reconcile.Result{}, err
		}

		// update the store to have the desired state as the latest state.
		if err = storage.Store(desiredResources, deployment); err != nil {
			err = errors.Wrap(err, "storing latest state in state store")
			return reconcile.Result{}, err
		}

		// process the secondary resources in the order they should be created
		for _, rec := range in.reconcilers {
			r, err := rec.ReconcileResources(ctx, request, deployment, storage)
			if err != nil {
				return reconcile.Result{}, err
			}
			result.Requeue = result.Requeue || r.Requeue
		}
	} else {
		// original and desired are identical so there is nothing else to do
		in.Log.Info("Reconciled secondary resources for deployment, nothing to update")
	}

	// if replica count is zero update the status to Stopped
	if deployment.GetReplicas() == 0 {
		if err = in.UpdateDeploymentStatusPhase(ctx, request.NamespacedName, coh.ConditionTypeStopped); err != nil {
			err = errors.Wrap(err, "error updating deployment status")
		}
	}

	in.Log.Info(fmt.Sprintf("Finished reconciling Coherence resource. Result='%v'", result))
	return result, err
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
		reconciler.NewSecretReconciler(mgr),
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

func (in *CoherenceReconciler) addFinalizer(ctx context.Context, c *coh.Coherence) error {
	in.Log.Info("Adding Finalizer to Coherence resource", "Namespace", c.Namespace, "Name", c.Name)
	controllerutil.AddFinalizer(c, coh.Finalizer)

	// Update CR
	err := in.GetClient().Update(ctx, c)
	if err != nil {
		in.Log.Error(err, "Failed to update Coherence resource with finalizer",
			"Namespace", c.Namespace, "Name", c.Name)
		return err
	}
	return nil
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
