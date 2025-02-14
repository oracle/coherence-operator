/*
 * Copyright (c) 2020, 2025, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package statefulset

import (
	"context"
	coh "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/pkg/nodes"
	"github.com/oracle/coherence-operator/pkg/operator"
	"github.com/oracle/coherence-operator/pkg/probe"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"time"
)

type UpgradeStrategy interface {
	// IsOperatorManaged returns true if this strategy requires the Operator to manage the upgrade
	IsOperatorManaged() bool
	// RollingUpgrade performs the rolling upgrade
	// Parameters:
	// context.Context      - the Go context to use
	// *appsv1.StatefulSet  - a pointer to the StatefulSet to upgrade
	// string               - the name of the WKA service
	// kubernetes.Interface - the K8s client
	RollingUpgrade(context.Context, *appsv1.StatefulSet, string, kubernetes.Interface) (reconcile.Result, error)
}

func GetUpgradeStrategy(c coh.CoherenceResource, p probe.CoherenceProbe) UpgradeStrategy {
	spec, _ := c.GetStatefulSetSpec()
	if spec.RollingUpdateStrategy != nil {
		name := *spec.RollingUpdateStrategy
		if name == coh.UpgradeManual {
			return ManualUpgradeStrategy{}
		}
		if name == coh.UpgradeByNode {
			sp := spec.GetScalingProbe()
			return ByNodeUpgradeStrategy{
				cp:           p,
				scalingProbe: sp,
			}
		}
		if name == coh.UpgradeByNodeLabel {
			sp := spec.GetScalingProbe()
			if spec.RollingUpdateLabel == nil {
				return ByNodeUpgradeStrategy{
					cp:           p,
					scalingProbe: sp,
				}
			} else {
				return ByNodeLabelUpgradeStrategy{
					label:        *spec.RollingUpdateLabel,
					cp:           p,
					scalingProbe: sp,
				}
			}
		}
	}
	// default is by Pod
	return ByPodUpgradeStrategy{}
}

// ----- ByPodUpgradeStrategy ----------------------------------------------------------------------

var _ UpgradeStrategy = ByPodUpgradeStrategy{}

type ByPodUpgradeStrategy struct {
}

func (in ByPodUpgradeStrategy) RollingUpgrade(context.Context, *appsv1.StatefulSet, string, kubernetes.Interface) (reconcile.Result, error) {
	return reconcile.Result{}, nil
}

func (in ByPodUpgradeStrategy) IsOperatorManaged() bool {
	return false
}

// ----- ManualUpgradeStrategy ---------------------------------------------------------------------

var _ UpgradeStrategy = ManualUpgradeStrategy{}

type ManualUpgradeStrategy struct {
}

func (in ManualUpgradeStrategy) RollingUpgrade(context.Context, *appsv1.StatefulSet, string, kubernetes.Interface) (reconcile.Result, error) {
	return reconcile.Result{}, nil
}

func (in ManualUpgradeStrategy) IsOperatorManaged() bool {
	return false
}

// ----- ByNodeUpgradeStrategy ---------------------------------------------------------------------

var _ UpgradeStrategy = ByNodeUpgradeStrategy{}

type ByNodeUpgradeStrategy struct {
	cp           probe.CoherenceProbe
	scalingProbe *coh.Probe
}

func (in ByNodeUpgradeStrategy) RollingUpgrade(ctx context.Context, sts *appsv1.StatefulSet, svc string, c kubernetes.Interface) (reconcile.Result, error) {
	return rollingUpgrade(in.cp, in.scalingProbe, &PodNodeName{}, "NodeName", ctx, sts, svc, c)
}

func (in ByNodeUpgradeStrategy) IsOperatorManaged() bool {
	return true
}

// ----- ByNodeLabelUpgradeStrategy -----------------------------------------------------------------

var _ UpgradeStrategy = ByNodeLabelUpgradeStrategy{}

type ByNodeLabelUpgradeStrategy struct {
	cp           probe.CoherenceProbe
	scalingProbe *coh.Probe
	label        string
}

func (in ByNodeLabelUpgradeStrategy) RollingUpgrade(ctx context.Context, sts *appsv1.StatefulSet, svc string, c kubernetes.Interface) (reconcile.Result, error) {
	return rollingUpgrade(in.cp, in.scalingProbe, &PodNodeLabel{Label: in.label}, in.label, ctx, sts, svc, c)
}

func (in ByNodeLabelUpgradeStrategy) IsOperatorManaged() bool {
	return true
}

// ----- PodNodeIdSupplier -------------------------------------------------------------------------

type PodNodeIdSupplier interface {
	GetNodeId(context.Context, kubernetes.Interface, corev1.Pod) (string, error)
}

var _ PodNodeIdSupplier = &PodNodeName{}

type PodNodeName struct {
}

func (p *PodNodeName) GetNodeId(_ context.Context, _ kubernetes.Interface, pod corev1.Pod) (string, error) {
	return pod.Spec.NodeName, nil
}

var _ PodNodeIdSupplier = &PodNodeLabel{}

type PodNodeLabel struct {
	Label string
	cache map[string]string
}

func (p *PodNodeLabel) GetNodeId(ctx context.Context, c kubernetes.Interface, pod corev1.Pod) (string, error) {
	var err error
	var value string
	var found bool

	if p.cache == nil {
		p.cache = make(map[string]string)
	} else {
		value, found = p.cache[pod.Spec.NodeName]
		if found {
			return value, err
		}
	}

	value, err = nodes.GetExactLabelForNode(ctx, c, pod.Spec.NodeName, p.Label, log)
	if err != nil {
		return "", err
	}
	p.cache[pod.Spec.NodeName] = value
	return value, err
}

// ----- helper methods ----------------------------------------------------------------------------

func rollingUpgrade(cp probe.CoherenceProbe, scalingProbe *coh.Probe, fn PodNodeIdSupplier, idName string, ctx context.Context, sts *appsv1.StatefulSet, svc string, c kubernetes.Interface) (reconcile.Result, error) {
	var err error
	var replicas int32

	if sts.Spec.Replicas == nil {
		replicas = 1
		if sts.Status.ReadyReplicas == 0 {
			return reconcile.Result{}, nil
		}
	} else {
		replicas = *sts.Spec.Replicas
		if sts.Status.ReadyReplicas != *sts.Spec.Replicas {
			return reconcile.Result{}, nil
		}
	}

	if sts.Status.CurrentRevision == sts.Status.UpdateRevision {
		return reconcile.Result{}, nil
	}

	if !operator.IsNodeLookupEnabled() {
		log.Info("Cannot rolling upgrade of StatefulSet by Node, the Coherence Operator has Node lookup disabled", "Namespace", sts.Namespace, "Name", sts.Name)
		return reconcile.Result{}, nil
	}

	log.Info("Perform rolling upgrade of StatefulSet by Node", "Namespace", sts.Namespace, "Name", sts.Name)

	pods, err := cp.GetPodsForStatefulSet(ctx, sts)
	if err != nil {
		log.Error(err, "Error getting list of Pods for StatefulSet", "Namespace", sts.Namespace, "Name", sts.Name)
		return reconcile.Result{}, err
	}

	if len(pods.Items) == 0 {
		log.Info("Zero Pods found for StatefulSet", "Namespace", sts.Namespace, "Name", sts.Name)
		return reconcile.Result{}, err
	}

	if len(pods.Items) != int(replicas) {
		log.Info("Count of Pods found for StatefulSet does not match replicas", "Namespace", sts.Namespace, "Name", sts.Name, "Replicas", replicas, "Found", len(pods.Items))
		return reconcile.Result{}, err
	}

	revision := sts.Status.UpdateRevision

	podsToUpdate := corev1.PodList{}
	if len(pods.Items) > 1 {
		// we have multiple Pods
		podsById, allPodsById, err := groupPods(ctx, c, pods, revision, fn)
		if err != nil {
			return reconcile.Result{}, err
		}

		if len(podsById) == 0 {
			// nothing to do, all pods are at the required revision
			return reconcile.Result{}, nil
		}

		if len(allPodsById) == 1 {
			// There is only one Node, we cannot be NodeSafe so do not do anything
			id, err := fn.GetNodeId(ctx, c, pods.Items[0])
			if err != nil {
				return reconcile.Result{}, err
			}
			log.Info("All Pods have a single Node identifier and cannot be Safe, no Pods will be upgraded", "Namespace", sts.Namespace,
				"Name", sts.Name, "Replicas", len(pods.Items), "NodeId", idName, "IdValue", id)
			return reconcile.Result{}, nil
		}

		// create an array of Node identifiers
		var identifiers []string
		for k := range podsById {
			identifiers = append(identifiers, k)
		}

		// Create the list of Pods to be deleted (upgraded) by picking the first identifier
		node := identifiers[0]
		p := podsById[node]
		podsToUpdate.Items = p
	} else {
		// There is only on Pod and replicas == 1
		pod := pods.Items[0]
		podRevision := pod.Labels["controller-revision-hash"]
		if revision != podRevision {
			// The single Pod is not the required revision so upgrade it
			podsToUpdate.Items = append(podsToUpdate.Items, pod)
		}
	}

	if len(podsToUpdate.Items) > 0 {
		// We have Pods to be upgraded
		nodeId, _ := fn.GetNodeId(ctx, c, pods.Items[0])
		// Check Pods are "safe"
		if cp.ExecuteProbeForSubSetOfPods(ctx, sts, svc, scalingProbe, pods, podsToUpdate) {
			// delete the Pods
			log.Info("Upgrading all Pods for Node identifier", "Namespace", sts.Namespace, "Name", sts.Name, "NodeId", idName, "IdValue", nodeId, "Count", len(podsToUpdate.Items))
			err = deletePods(ctx, podsToUpdate, c)
		} else {
			log.Info("Pods failed Status HA check, upgrade is deferred for one minute", "Namespace", sts.Namespace, "Name", sts.Name, "NodeId", idName, "IdValue", nodeId)
			return reconcile.Result{Requeue: true, RequeueAfter: time.Minute}, nil
		}
	}

	// Even if we still have nodes to upgrade, we send s non-requeue result.
	// When the deleted Pods are rescheduled and become ready the StatefulSet status
	// will be updated, and we will end up back in this method
	return reconcile.Result{}, err
}

// deletePods will delete the pods in a pod list
func deletePods(ctx context.Context, pods corev1.PodList, c kubernetes.Interface) error {
	for _, pod := range pods.Items {
		log.Info("Attempting to delete Pod to trigger upgrade", "Namespace", pod.Namespace, "Name", pod.Name)
		if err := c.CoreV1().Pods(pod.Namespace).Delete(ctx, pod.Name, metav1.DeleteOptions{}); err != nil {
			log.Error(err, "Error deleting Pod", "Namespace", pod.Namespace, "Name", pod.Name)
			return err
		}
		log.Info("Deleted Pod", "Namespace", pod.Namespace, "Name", pod.Name)
	}
	return nil
}

// groupPods returns two maps of Pods by an identifier. The first is Pods with a specific controller revision, the second is all Pods
func groupPods(ctx context.Context, c kubernetes.Interface, pods corev1.PodList, revision string, fn PodNodeIdSupplier) (map[string][]corev1.Pod, map[string][]corev1.Pod, error) {
	allPodsById := make(map[string][]corev1.Pod)
	podsById := make(map[string][]corev1.Pod)
	for _, pod := range pods.Items {
		id, err := fn.GetNodeId(ctx, c, pod)
		if err != nil {
			return nil, nil, err
		}
		allPods, found := allPodsById[id]
		if !found {
			allPods = make([]corev1.Pod, 0)
		}
		allPods = append(allPods, pod)
		allPodsById[id] = allPods

		podRevision := pod.Labels["controller-revision-hash"]
		if revision != podRevision {
			ary, found := podsById[id]
			if !found {
				ary = make([]corev1.Pod, 0)
			}
			ary = append(ary, pod)
			podsById[id] = ary
		}
	}
	return podsById, allPodsById, nil
}
