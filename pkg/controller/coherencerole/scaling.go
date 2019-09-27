/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package coherencerole

import (
	"context"
	"fmt"
	coh "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	mgmt "github.com/oracle/coherence-operator/pkg/management"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/rest"
	"k8s.io/kubernetes/pkg/probe"
	httpprobe "k8s.io/kubernetes/pkg/probe/http"
	tcprobe "k8s.io/kubernetes/pkg/probe/tcp"
	"net/http"
	"net/url"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"strconv"
	"strings"
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
		}
		return r.safeScale(role, cohInternal, existing, desired, current, sts)
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

	checker := ScalableChecker{Client: r.client, Config: r.mgr.GetConfig()}
	ha := current == 1 || checker.IsStatusHA(role, sts)

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
			}
			// scaled by one but not yet at the desired size - requeue request after one minute
			return reconcile.Result{Requeue: true, RequeueAfter: time.Minute}, nil
		}
		// failed
		return r.handleErrAndRequeue(err, role, fmt.Sprintf(failedToScaleRole, role.Name, current, replicas, err.Error()), logger)
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

	existing.Spec.Replicas = &desired
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

type ScalableChecker struct {
	Client         client.Client
	Config         *rest.Config
	getPodHostName func(pod corev1.Pod) string
	translatePort  func(name string, port int) int
}

func (in *ScalableChecker) SetGetPodHostName(fn func(pod corev1.Pod) string) {
	if in == nil {
		return
	}
	in.getPodHostName = fn
}

func (in *ScalableChecker) GetPodHostName(pod corev1.Pod) string {
	if in.getPodHostName == nil {
		return pod.Status.PodIP
	}
	return in.getPodHostName(pod)
}

func (in *ScalableChecker) SetTranslatePort(fn func(name string, port int) int) {
	if in == nil {
		return
	}
	in.translatePort = fn
}

func (in *ScalableChecker) TranslatePort(name string, port int) int {
	if in.translatePort == nil {
		return port
	}
	return in.translatePort(name, port)
}

// IsStatusHA will return true if the cluster represented by the role is StatusHA.
func (in *ScalableChecker) IsStatusHA(role *coh.CoherenceRole, sts *appsv1.StatefulSet) bool {
	list := corev1.PodList{}
	opts := client.ListOptions{}
	opts.InNamespace(role.Namespace)
	opts.MatchingLabels(sts.Spec.Selector.MatchLabels)

	if log.Enabled() {
		log.Info("Checking StatefulSet "+sts.Name+" for StatusHA", "Namespace", role.Name, "Name", role.Name)
	}

	err := in.Client.List(context.TODO(), &opts, &list)
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

	scalingProbe := role.Spec.GetScalingProbe()

	for _, pod := range list.Items {
		if pod.Status.Phase == "Running" {
			if log.Enabled() {
				log.Info("Checking pod " + pod.Name + " for StatusHA")
			}

			ha, err := in.CanScale(pod, scalingProbe)
			if err == nil {
				log.Info(fmt.Sprintf("Checked pod %s for StatusHA (%t)", pod.Name, ha))
				return ha
			}
			log.Info(fmt.Sprintf("Checked pod %s for StatusHA (%t) error %s", pod.Name, ha, err.Error()))
		} else {
			log.Info("Skipping StatusHA checking for pod " + pod.Name + " as Pod status not in running phase")
		}
	}

	return false
}

// Determine whether a role allowed to scale using the configured probe.
func (in *ScalableChecker) CanScale(pod corev1.Pod, handler *coh.ScalingProbe) (bool, error) {
	switch {
	case handler.Exec != nil:
		return in.ExecIsPodStatusHA(pod, handler)
	case handler.HTTPGet != nil:
		return in.HTTPIsPodStatusHA(pod, handler)
	case handler.TCPSocket != nil:
		return in.TCPIsPodStatusHA(pod, handler)
	default:
		return true, nil
	}
}

func (in *ScalableChecker) ExecIsPodStatusHA(pod corev1.Pod, handler *coh.ScalingProbe) (bool, error) {
	req := &mgmt.ExecRequest{
		Pod:       pod.Name,
		Container: "coherence",
		Namespace: pod.Namespace,
		Command:   handler.Exec.Command,
		Arg:       []string{},
		Timeout:   handler.GetTimeout(),
	}

	exitCode, _, _, err := mgmt.PodExec(req, in.Config)

	log.Info(fmt.Sprintf("StatusHA check Exec: '%s' result=%d error=%s", strings.Join(handler.Exec.Command, ", "), exitCode, err))

	if err != nil {
		return false, err
	}

	return exitCode == 0, nil
}

func (in *ScalableChecker) HTTPIsPodStatusHA(pod corev1.Pod, handler *coh.ScalingProbe) (bool, error) {
	var (
		scheme corev1.URIScheme
		host   string
		port   int
		path   string
	)

	action := handler.HTTPGet

	if action.Scheme == "" {
		scheme = corev1.URISchemeHTTP
	} else {
		scheme = action.Scheme
	}

	if action.Host == "" {
		host = in.GetPodHostName(pod)
	} else {
		host = action.Host
	}

	port, err := in.findPort(pod, action.Port)
	if err != nil {
		return false, err
	}

	if strings.HasPrefix(action.Path, "/") {
		path = action.Path[1:]
	} else {
		path = action.Path
	}

	u, err := url.Parse(fmt.Sprintf("%s://%s:%d/%s", scheme, host, port, path))
	if err != nil {
		return false, err
	}

	header := http.Header{}
	if action.HTTPHeaders != nil {
		for _, h := range action.HTTPHeaders {
			hh, found := header[h.Name]
			if found {
				header[h.Name] = append(hh, h.Value)
			} else {
				header[h.Name] = []string{h.Value}
			}
		}
	}

	p := httpprobe.New()
	result, s, err := p.Probe(u, header, handler.GetTimeout())

	log.Info(fmt.Sprintf("StatusHA check URL: %s result=%s msg=%s error=%s", u.String(), result, s, err))

	return result == probe.Success, err
}

func (in *ScalableChecker) TCPIsPodStatusHA(pod corev1.Pod, handler *coh.ScalingProbe) (bool, error) {
	var (
		host string
		port int
	)

	action := handler.TCPSocket

	if action.Host == "" {
		host = in.GetPodHostName(pod)
	} else {
		host = action.Host
	}

	port, err := in.findPort(pod, action.Port)
	if err != nil {
		return false, err
	}

	p := tcprobe.New()
	result, _, err := p.Probe(host, port, handler.GetTimeout())

	log.Info(fmt.Sprintf("StatusHA check TCP: %s:%d result=%s error=%s", host, port, result, err))

	return result == probe.Success, err
}

func (in *ScalableChecker) findPort(pod corev1.Pod, port intstr.IntOrString) (int, error) {
	if port.Type == intstr.Int {
		return port.IntValue(), nil
	}

	s := port.String()
	i, err := strconv.Atoi(s)
	if err == nil {
		// string is an int
		return i, nil
	}
	// string is a port name
	return in.findPortInPod(pod, s)
}

func (in *ScalableChecker) findPortInPod(pod corev1.Pod, name string) (int, error) {
	for _, container := range pod.Spec.Containers {
		if container.Name == "coherence" {
			return in.findPortInContainer(pod, container, name)
		}
	}

	return -1, fmt.Errorf("cannot find coherence container in Pod '%s'", pod.Name)
}

func (in *ScalableChecker) findPortInContainer(pod corev1.Pod, container corev1.Container, name string) (int, error) {
	for _, port := range container.Ports {
		if port.Name == name {
			p := in.TranslatePort(port.Name, int(port.ContainerPort))
			return p, nil
		}
	}

	return -1, fmt.Errorf("cannot find port '%s' in coherence container in Pod '%s'", name, pod.Name)
}
