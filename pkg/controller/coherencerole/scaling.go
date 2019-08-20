package coherencerole

import (
	"context"
	"fmt"
	coh "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	mgmt "github.com/oracle/coherence-operator/pkg/management"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"time"
)

// scale will scale a role up or down
func (r *ReconcileCoherenceRole) scale(role *coh.CoherenceRole, cohInternal *unstructured.Unstructured, existing *coh.CoherenceInternal, desired int32, current int32, sts *appsv1.StatefulSet) (reconcile.Result, error) {
	policy := role.Spec.GetEffectiveScalingPolicy()

	switch policy {
	case coh.SafeScaling:
		return r.safeScale(role, cohInternal, existing, desired, current, sts)
	case coh.ParallelScaling:
		return r.parallelScale(role, cohInternal, existing, desired, current)
	case coh.ParallelUpSafeDownScaling:
		if desired > current {
			return r.parallelScale(role, cohInternal, existing, desired, current)
		} else {
			return r.safeScale(role, cohInternal, existing, desired, current, sts)
		}
	default:
		// shouldn't get here, but better safe than sorry
		return r.safeScale(role, cohInternal, existing, desired, current, sts)
	}
}

// safeScale will scale a role up or down by one and requeue the request.
func (r *ReconcileCoherenceRole) safeScale(role *coh.CoherenceRole, cohInternal *unstructured.Unstructured, existing *coh.CoherenceInternal, desired int32, current int32, sts *appsv1.StatefulSet) (reconcile.Result, error) {
	logger := log.WithValues("Namespace", role.Name, "Name", role.Name)

	if sts.Status.ReadyReplicas != current {
		logger.Info(fmt.Sprintf("Role %s is not StatusHA - re-queing scaling request. Stateful set ready replicas is %d", role.Name, sts.Status.ReadyReplicas))
	}

	ha := current == 1 || r.IsStatusHA(role, sts)

	if ha {
		var replicas int32

		if desired > current {
			replicas = current + 1
		} else {
			replicas = current - 1
		}

		logger.Info(fmt.Sprintf("Role %s is StatusHA, safely scaling from %d to %d (final desired replicas %d)", role.Name, current, replicas, desired))

		// use the parallel method to just scale by one
		_, err := r.parallelScale(role, cohInternal, existing, replicas, current)
		if err == nil {
			if replicas == desired {
				// we're at the desired size so finished scaling
				return reconcile.Result{Requeue: false}, nil
			} else {
				// scaled by one but not yet at the desired size - requeue request after one minute
				return reconcile.Result{Requeue: true, RequeueAfter: time.Minute}, nil
			}
		} else {
			// failed
			return r.handleErrAndRequeue(err, role, fmt.Sprintf(failedToScaleRole, role.Name, current, replicas, err.Error()), logger)
		}
	}

	// Not StatusHA - wait one minute
	logger.Info(fmt.Sprintf("Role %s is not StatusHA - re-queing scaling request", role.Name))
	return reconcile.Result{Requeue: true, RequeueAfter: r.statusHARetry}, nil
}

// parallelScale will scale the role by the required amount in one request.
func (r *ReconcileCoherenceRole) parallelScale(role *coh.CoherenceRole, cohInternal *unstructured.Unstructured, existing *coh.CoherenceInternal, desired int32, current int32) (reconcile.Result, error) {
	// update the CoherenceInternal, this should trigger an update of the Helm install to scale the StatefulSet

	// Update this CoherenceRole's status
	role.Status.Status = coh.RoleStatusScaling
	role.Status.Replicas = desired
	err := r.client.Status().Update(context.TODO(), role)
	if err != nil {
		// failed to update the CoherenceRole's status
		// ToDo - handle this properly by re-queuing the request and then in the reconcile method properly handle setting status even if the role is in the desired state
		log.Error(err, "failed to update role status")
	}

	existing.Spec.ClusterSize = desired
	cohInternal.Object["spec"] = existing.Spec
	err = r.client.Update(context.TODO(), cohInternal)
	if err != nil {
		// send a failed scale event
		msg := fmt.Sprintf("failed to scale Helm install %s in CoherenceRole %s from %d to %d", role.Name, role.Name, current, desired)
		r.events.Event(role, corev1.EventTypeNormal, "SuccessfulScale", msg)

		return reconcile.Result{}, err
	}

	// send a successful scale event
	msg := fmt.Sprintf("scaled Helm install %s in CoherenceRole %s from %d to %d", role.Name, role.Name, current, desired)
	r.events.Event(role, corev1.EventTypeNormal, eventReasonScale, msg)

	return reconcile.Result{}, nil
}

// IsStatusHA will return true if the cluster represented by the role is StatusHA.
func (r *ReconcileCoherenceRole) IsStatusHA(role *coh.CoherenceRole, sts *appsv1.StatefulSet) bool {
	list := corev1.PodList{}
	opts := client.ListOptions{}
	opts.InNamespace(role.Namespace)
	opts.MatchingLabels(sts.Spec.Selector.MatchLabels)

	if log.Enabled() {
		log.Info("Checking StatefulSet "+sts.Name+" for StatusHA", "Namespace", role.Name, "Name", role.Name)
	}

	err := r.client.List(context.TODO(), &opts, &list)
	if err != nil {
		log.Error(err, "Error getting list of Pods for StatefulSet "+sts.Name)
		return false
	}

	if len(list.Items) == 0 {
		if log.Enabled() {
			log.Info("Cannot find any Pods for StatefulSet " + sts.Name + " - assuming StatusHA is true")
		}
		return true
	}

	for _, pod := range list.Items {
		if pod.Status.Phase == "Running" {
			ip := pod.Status.PodIP
			ha, err := IsPodStatusHA(ip)
			if log.Enabled() {
				log.Info("Checking pod " + pod.Name + " for StatusHA")
			}
			if err == nil {
				return ha
			}
		} else {
			log.Info("Skipping StatusHA checking for pod " + pod.Name + " as Pod status not in running phase")
		}
	}

	return false
}

// Determine whether a Pod's Coherence Services are StatusHA.
func IsPodStatusHA(podIP string) (bool, error) {
	cl := &http.Client{}

	services, _, err := mgmt.GetServices(cl, podIP, 30000)
	if err != nil {
		if log.Enabled() {
			log.Info("Error querying services from podIP " + podIP + "\n" + err.Error())
		}
		return false, err
	}

	for _, service := range services.Items {
		if service.Type == "DistributedCache" {
			part, rc, err := mgmt.GetPartitionAssignment(cl, podIP, 30000, service.Name)
			if err == nil {
				if rc == http.StatusOK {
					// we must have more than one service member and backups > 0 to event think about being HA.
					if part.BackupCount > 0 && part.ServiceNodeCount > 1 {
						if part.HAStatusCode <= 1 || part.RemainingDistributionCount != 0 {
							// we're not HA
							return false, nil
						}
					}
				}
			} else {
				log.Info("Error accessing podIP " + podIP + "\n" + err.Error())
			}
		}
	}

	return true, nil
}
