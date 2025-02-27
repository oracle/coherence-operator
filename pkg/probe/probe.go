/*
 * Copyright (c) 2020, 2025, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package probe

import (
	"context"
	"fmt"
	coh "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/pkg/events"
	mgmt "github.com/oracle/coherence-operator/pkg/management"
	"github.com/oracle/coherence-operator/pkg/operator"
	"github.com/spf13/viper"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/rest"
	"net/http"
	"net/url"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"strconv"
	"strings"
)

// Result is a string used to handle the results for probing container readiness/liveness
type Result string

const (
	// Success Result
	Success Result = "success"
	// Failure Result
	Failure Result = "failure"
	// Unknown Result
	Unknown Result = "unknown"
)

var log = logf.Log.WithName("Probe")

type CoherenceProbe struct {
	Client         client.Client
	Config         *rest.Config
	EventRecorder  events.OwnedEventRecorder
	getPodHostName func(pod corev1.Pod) string
	translatePort  func(name string, port int) int
}

func (in *CoherenceProbe) SetGetPodHostName(fn func(pod corev1.Pod) string) {
	if in == nil {
		return
	}
	in.getPodHostName = fn
}

func (in *CoherenceProbe) GetPodIpOrHostName(pod corev1.Pod) string {
	hostName, found := pod.Labels[operator.LabelTestHostName]
	if found {
		return hostName
	}
	if in.getPodHostName == nil {
		if pod.Status.PodIP == "" {
			return pod.Spec.Hostname + "." + pod.Spec.Subdomain + "." + pod.GetNamespace() + ".svc"
		}
		return pod.Status.PodIP
	}
	return in.getPodHostName(pod)
}

func (in *CoherenceProbe) GetPodHostName(pod corev1.Pod, svc string) string {
	hostName, found := pod.Labels[operator.LabelTestHostName]
	if found {
		return hostName
	}
	if in.getPodHostName != nil {
		return in.getPodHostName(pod)
	}
	return fmt.Sprintf("%s.%s.%s", pod.Name, svc, pod.Namespace)
}

func (in *CoherenceProbe) SetTranslatePort(fn func(name string, port int) int) {
	if in == nil {
		return
	}
	in.translatePort = fn
}

func (in *CoherenceProbe) TranslatePort(name string, port int) int {
	if in.translatePort == nil {
		return port
	}
	return in.translatePort(name, port)
}

// IsStatusHA will return true if the deployment represented by the deployment is StatusHA.
// The number of Pods matching the StatefulSet selector must match the StatefulSet replica count
// ALl Pods must be in the ready state
// All Pods must pass the StatusHA check
func (in *CoherenceProbe) IsStatusHA(ctx context.Context, deployment coh.CoherenceResource, sts *appsv1.StatefulSet) bool {
	log.Info("Checking StatefulSet "+sts.Name+" for StatusHA",
		"Namespace", deployment.GetNamespace(), "Name", deployment.GetName())

	spec, found := deployment.GetStatefulSetSpec()
	if found {
		p := spec.GetScalingProbe()
		return in.ExecuteProbe(ctx, sts, deployment.GetWkaServiceName(), p)
	}
	return true
}

type ServiceSuspendStatus int

const ( // iota is reset to 0
	ServiceSuspendSkipped    ServiceSuspendStatus = iota // == 0
	ServiceSuspendSuccessful ServiceSuspendStatus = iota // == 1
	ServiceSuspendFailed     ServiceSuspendStatus = iota // == 2
)

// SuspendServices will request services be suspended in the Coherence cluster.
// This is called prior to stopping a StatefulSet to then have a graceful shutdown.
// The number of Pods matching the StatefulSet selector must match the StatefulSet replica count
// ALl Pods must be in the ready state
// All Pods must pass the StatusHA check
func (in *CoherenceProbe) SuspendServices(ctx context.Context, deployment coh.CoherenceResource, sts *appsv1.StatefulSet) ServiceSuspendStatus {
	ns := deployment.GetNamespace()
	name := deployment.GetName()

	if deployment.GetType() != coh.CoherenceTypeStatefulSet {
		return ServiceSuspendSkipped
	}

	c := deployment.(*coh.Coherence)
	stsSpec, _ := c.GetStatefulSetSpec()

	if viper.GetBool(operator.FlagSkipServiceSuspend) {
		msg := fmt.Sprintf("Skipping suspension of Coherence services in StatefulSet %s, flag %s is set to true",
			sts.Name, operator.FlagSkipServiceSuspend)
		log.Info(msg, "Namespace", ns, "Name", name)
		in.EventRecorder.Warn("ServiceSuspendSkipped", msg)
		return ServiceSuspendSkipped
	}

	if !stsSpec.IsSuspendServicesOnShutdown() {
		msg := fmt.Sprintf("Skipping suspension of Coherence services in StatefulSet %s, spec.SuspendServicesOnShutdown is set to false", sts.Name)
		log.Info(msg, "Namespace", ns, "Name", name)
		in.EventRecorder.Warn("ServiceSuspendSkipped", msg)
		return ServiceSuspendSkipped
	}

	log.Info("Suspending Coherence services in StatefulSet "+sts.Name, "Namespace", ns, "Name", name)
	if in.ExecuteProbe(ctx, sts, deployment.GetWkaServiceName(), stsSpec.GetSuspendProbe()) {
		in.EventRecorder.Warnf("ServiceSuspendFailed", "failed to suspend Coherence services in StatefulSet %s", sts.Name)
		return ServiceSuspendSuccessful
	}
	in.EventRecorder.Infof("ServiceSuspended", "suspended Coherence services in StatefulSet %s", sts.Name)
	return ServiceSuspendFailed
}

func (in *CoherenceProbe) GetPodsForStatefulSet(ctx context.Context, sts *appsv1.StatefulSet) (corev1.PodList, error) {
	pods := corev1.PodList{}
	labels := client.MatchingLabels{}
	for k, v := range sts.Spec.Selector.MatchLabels {
		labels[k] = v
	}
	err := in.Client.List(ctx, &pods, client.InNamespace(sts.GetNamespace()), labels)
	return pods, err
}

func (in *CoherenceProbe) ExecuteProbe(ctx context.Context, sts *appsv1.StatefulSet, svc string, probe *coh.Probe) bool {
	pods, err := in.GetPodsForStatefulSet(ctx, sts)
	if err != nil {
		log.Error(err, "Error getting list of Pods for StatefulSet "+sts.Name)
		in.EventRecorder.Infof("CheckStatusHA", "Failed to get pods for StatefulSet %s: %s", sts.Name, err.Error())
		return false
	}
	return in.ExecuteProbeForSubSetOfPods(ctx, sts, svc, probe, pods, pods)
}

func (in *CoherenceProbe) ExecuteProbeForSubSetOfPods(ctx context.Context, sts *appsv1.StatefulSet, svc string, probe *coh.Probe, stsPods, pods corev1.PodList) bool {
	logger := log.WithValues("Namespace", sts.GetNamespace(), "Name", sts.GetName())

	// All Pods must be in the Running Phase
	for _, pod := range stsPods.Items {
		if ready, phase := in.IsPodReady(pod); !ready {
			msg := fmt.Sprintf("Cannot execute probe, one or more Pods is not in a ready state - %s (%v) ", pod.Name, phase)
			logger.Info(msg)
			in.EventRecorder.Warn("CheckStatusHA", msg)
			return false
		}
	}

	count := int32(len(stsPods.Items))
	switch {
	case count == 0:
		msg := fmt.Sprintf("Skipping StatusHA check, no Pods found in StatefulSet %s", sts.Name)
		logger.Info(msg)
		in.EventRecorder.Info("CheckStatusHA", msg)
		return true
	case sts.Spec.Replicas == nil && count != 1:
		msg := fmt.Sprintf("Pod count of %d does not yet match StatefulSet replica count: 1", count)
		logger.Info(msg)
		in.EventRecorder.Info("CheckStatusHA", msg)
		return false
	case sts.Spec.Replicas != nil && count != *sts.Spec.Replicas:
		msg := fmt.Sprintf("Pod count of %d does not yet match StatefulSet replica count: %d", count, *sts.Spec.Replicas)
		in.EventRecorder.Info("CheckStatusHA", msg)
		logger.Info(msg)
		return false
	}

	for _, pod := range pods.Items {
		if pod.Status.Phase == "Running" {
			if log.Enabled() {
				log.Info("Using pod " + pod.Name + " to execute probe")
			}

			ha, err := in.RunProbe(ctx, pod, svc, probe)
			if err == nil {
				msg := fmt.Sprintf("Executed probe using pod %s result=%t", pod.Name, ha)
				log.Info(msg)
				in.EventRecorder.Info("CheckStatusHA", msg)
				return ha
			}
			msg := fmt.Sprintf("Execute probe using pod %s (%t) error %s", pod.Name, ha, err.Error())
			in.EventRecorder.Warn("CheckStatusHA", msg)
			log.Info(msg)
		} else {
			msg := fmt.Sprintf("Skipping execute probe for pod %s as Pod status not in running phase", pod.Name)
			log.Info(msg)
			in.EventRecorder.Warn("CheckStatusHA", msg)
		}
	}

	return false
}

// IsPodReady determines whether the specified Pods are in the Ready state.
func (in *CoherenceProbe) IsPodReady(pod corev1.Pod) (bool, string) {
	if pod.DeletionTimestamp != nil || pod.Status.Phase != corev1.PodRunning {
		return false, "Terminating"
	}

	for _, c := range pod.Status.Conditions {
		if c.Type == corev1.PodReady && c.Status == corev1.ConditionTrue {
			return true, string(pod.Status.Phase)
		}
	}
	return false, string(pod.Status.Phase)
}

func (in *CoherenceProbe) RunProbe(ctx context.Context, pod corev1.Pod, svc string, handler *coh.Probe) (bool, error) {
	switch {
	case handler.Exec != nil:
		return in.ProbeUsingExec(ctx, pod, handler)
	case handler.HTTPGet != nil:
		return in.ProbeUsingHTTP(pod, svc, handler)
	case handler.TCPSocket != nil:
		return in.ProbeUsingTCP(pod, handler)
	default:
		return true, nil
	}
}

func (in *CoherenceProbe) ProbeUsingExec(ctx context.Context, pod corev1.Pod, handler *coh.Probe) (bool, error) {
	req := &mgmt.ExecRequest{
		Pod:       pod.Name,
		Container: coh.ContainerNameCoherence,
		Namespace: pod.Namespace,
		Command:   handler.Exec.Command,
		Arg:       []string{},
		Timeout:   handler.GetTimeout(),
	}

	exitCode, _, _, err := mgmt.PodExec(ctx, req, in.Config)

	log.Info(fmt.Sprintf("Exec Probe: '%s' result=%d error=%s", strings.Join(handler.Exec.Command, ", "), exitCode, err))

	if err != nil {
		return false, err
	}

	return exitCode == 0, nil
}

func (in *CoherenceProbe) ProbeUsingHTTP(pod corev1.Pod, svc string, handler *coh.Probe) (bool, error) {
	var (
		scheme   corev1.URIScheme
		hostOrIP string
		port     int
		path     string
	)

	action := handler.HTTPGet

	if action.Scheme == "" {
		scheme = corev1.URISchemeHTTP
	} else {
		scheme = action.Scheme
	}

	if action.Host == "" {
		hostOrIP = in.GetPodIpOrHostName(pod)
	} else {
		hostOrIP = action.Host
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

	u, err := url.Parse(fmt.Sprintf("%s://%s:%d/%s", scheme, hostOrIP, port, path))
	if err != nil {
		return false, err
	}

	header := http.Header{}
	header.Set("Host", in.GetPodHostName(pod, svc))

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

	p := NewHTTPProbe()
	result, s, err := p.Probe(u, header, handler.GetTimeout())

	log.Info("Executed HTTP Probe", "URL", u, "Result", fmt.Sprintf("%v", result), "Msg", s, "Error", err)

	return result == Success, err
}

func (in *CoherenceProbe) ProbeUsingTCP(pod corev1.Pod, handler *coh.Probe) (bool, error) {
	var (
		host string
		port int
	)

	action := handler.TCPSocket

	if action.Host == "" {
		host = in.GetPodIpOrHostName(pod)
	} else {
		host = action.Host
	}

	port, err := in.findPort(pod, action.Port)
	if err != nil {
		return false, err
	}

	p := NewTCPProbe()
	result, _, err := p.Probe(host, port, handler.GetTimeout())

	log.Info(fmt.Sprintf("TCP Probe: %s:%d result=%s error=%s", host, port, result, err))

	return result == Success, err
}

func (in *CoherenceProbe) findPort(pod corev1.Pod, port intstr.IntOrString) (int, error) {
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

func (in *CoherenceProbe) findPortInPod(pod corev1.Pod, name string) (int, error) {
	if name == coh.PortNameHealth {
		p, found := pod.Labels[operator.LabelTestHealthPort]
		if found {
			return strconv.Atoi(p)
		}
	}

	for _, container := range pod.Spec.Containers {
		if container.Name == coh.ContainerNameCoherence {
			return in.findPortInContainer(pod, container, name)
		}
	}

	return -1, fmt.Errorf("cannot find coherence container in Pod '%s'", pod.Name)
}

func (in *CoherenceProbe) findPortInContainer(pod corev1.Pod, container corev1.Container, name string) (int, error) {
	for _, port := range container.Ports {
		if port.Name == name {
			p := in.TranslatePort(port.Name, int(port.ContainerPort))
			return p, nil
		}
	}

	return -1, fmt.Errorf("cannot find port '%s' in coherence container in Pod '%s'", name, pod.Name)
}
