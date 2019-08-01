package coherencerole

import (
	"context"
	"fmt"
	coh "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/kubernetes/pkg/apis/core"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"time"
)

// scale will scale a role up or down
func (r *ReconcileCoherenceRole) scale(role *coh.CoherenceRole, cohInternal *unstructured.Unstructured, existing *coh.CoherenceInternal, desired int32, current int32, sts *appsv1.StatefulSet) (reconcile.Result, error) {
	var policy coh.ScalingPolicy

	if role.Spec.ScalingPolicy == nil {
		// the scaling policy is not set the look at the storage enabled flag
		if role.Spec.StorageEnabled != nil || *role.Spec.StorageEnabled {
			// storage enabled is either not set or is true so do safe scaling
			policy = coh.SafeScaling
		} else {
			// storage enabled is false so do parallel scaling
			policy = coh.ParallelScaling
		}
	} else {
		// scaling policy is set so use it
		policy = *role.Spec.ScalingPolicy
	}

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
	if r.isStatusHA(role, sts) {
		var replicas int32

		if desired > current {
			replicas = current + 1
		} else {
			replicas = current - 1
		}

		// use the parallel method to just scale by one
		_, err := r.parallelScale(role, cohInternal, existing, replicas, current)
		if err == nil {
			if replicas == desired {
				// we're at the desired size so finished scaling
				return reconcile.Result{}, nil
			} else {
				// scaled by one but not yet at the desired size - requeue request after one minute
				return reconcile.Result{Requeue: true, RequeueAfter: time.Minute}, nil
			}
		} else {
			// failed
			return reconcile.Result{}, err
		}
	}

	// Not StatusHA - wait one minute
	return reconcile.Result{Requeue: true, RequeueAfter: time.Minute}, nil
}

// parallelScale will scale the role by the required amount in one request.
func (r *ReconcileCoherenceRole) parallelScale(role *coh.CoherenceRole, cohInternal *unstructured.Unstructured, existing *coh.CoherenceInternal, desired int32, current int32) (reconcile.Result, error) {
	// update the CoherenceInternal, this should trigger an update of the Helm install to scale the StatefulSet
	existing.Spec.ClusterSize = desired
	cohInternal.Object["spec"] = existing.Spec
	err := r.client.Update(context.TODO(), cohInternal)
	if err != nil {
		return reconcile.Result{}, err
	}

	// Update this CoherenceRole's status
	role.Status.Status = coh.RoleStatusScaling
	role.Status.Replicas = desired
	err = r.client.Status().Update(context.TODO(), role)

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

// isStatusHA will return true if the cluster represented by the role is StatusHA.
func (r *ReconcileCoherenceRole) isStatusHA(role *coh.CoherenceRole, sts *appsv1.StatefulSet) bool {
	list := core.PodList{}
	opts := client.ListOptions{}
	opts.InNamespace(role.Namespace)
	opts.MatchingLabels(sts.Spec.Selector.MatchLabels)

	err := r.client.List(context.TODO(), &opts, &list)
	if err != nil {
		return false
	}

	if len(list.Items) == 0 {
		return false
	}

	for _, pod := range list.Items {
		fmt.Println(pod.Name)
	}

	return true
}
