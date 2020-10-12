/*
 * Copyright (c) 2020, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package statefulset

import (
	"context"
	"fmt"
	coh "github.com/oracle/coherence-operator/api/v1"
	mgmt "github.com/oracle/coherence-operator/pkg/management"
	"github.com/oracle/coherence-operator/pkg/operator"
	"github.com/spf13/viper"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/rest"
	"k8s.io/kubernetes/pkg/probe"
	httpprobe "k8s.io/kubernetes/pkg/probe/http"
	tcprobe "k8s.io/kubernetes/pkg/probe/tcp"
	"net/http"
	"net/url"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strconv"
	"strings"
)

type CoherenceProbe struct {
	Client         client.Client
	Config         *rest.Config
	getPodHostName func(pod corev1.Pod) string
	translatePort  func(name string, port int) int
}

func (in *CoherenceProbe) SetGetPodHostName(fn func(pod corev1.Pod) string) {
	if in == nil {
		return
	}
	in.getPodHostName = fn
}

func (in *CoherenceProbe) GetPodHostName(pod corev1.Pod) string {
	if in.getPodHostName == nil {
		return pod.Status.PodIP
	}
	return in.getPodHostName(pod)
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
func (in *CoherenceProbe) IsStatusHA(deployment *coh.Coherence, sts *appsv1.StatefulSet) bool {
	log.Info("Checking StatefulSet "+sts.Name+" for StatusHA",
		"Namespace", deployment.Namespace, "Name", deployment.Name)
	p := deployment.Spec.GetScalingProbe()
	return in.ExecuteProbe(deployment, sts, p)
}

// SuspendServices will request services be suspended in the Coherence cluster.
// This is called prior to stopping a StatefulSet to then have a graceful shutdown.
// The number of Pods matching the StatefulSet selector must match the StatefulSet replica count
// ALl Pods must be in the ready state
// All Pods must pass the StatusHA check
func (in *CoherenceProbe) SuspendServices(deployment *coh.Coherence, sts *appsv1.StatefulSet) bool {
	if viper.GetBool(operator.FlagSkipServiceSuspend) {
		log.Info("Skipping suspension of Coherence services in StatefulSet "+sts.Name+
			operator.FlagSkipServiceSuspend+" is set to true",
			"Namespace", deployment.Namespace, "Name", deployment.Name)
		return true
	}

	// check whether the Coherence deployment supports service suspension
	if deployment.Annotations[coh.ANNOTATION_FEATURE_SUSPEND] != "true" {
		log.Info("Skipping suspension of Coherence services in StatefulSet "+sts.Name+
			coh.ANNOTATION_FEATURE_SUSPEND+" annotation is missing or not set to true",
			"Namespace", deployment.Namespace, "Name", deployment.Name)
		return true
	}

	if deployment.Spec.SuspendServicesOnShutdown != nil && !*deployment.Spec.SuspendServicesOnShutdown {
		log.Info("Skipping suspension of Coherence services in StatefulSet "+sts.Name+
			" spec.SuspendServicesOnShutdown is set to false",
			"Namespace", deployment.Namespace, "Name", deployment.Name)
		return true
	}
	log.Info("Suspending Coherence services in StatefulSet "+sts.Name,
		"Namespace", deployment.Namespace, "Name", deployment.Name)
	p := deployment.Spec.GetSuspendProbe()
	return in.ExecuteProbe(deployment, sts, p)
}

func (in *CoherenceProbe) ExecuteProbe(deployment *coh.Coherence, sts *appsv1.StatefulSet, probe *coh.Probe) bool {
	logger := log.WithValues("Namespace", deployment.Namespace, "Name", deployment.Name)
	list := corev1.PodList{}

	labels := client.MatchingLabels{}
	for k, v := range sts.Spec.Selector.MatchLabels {
		labels[k] = v
	}

	err := in.Client.List(context.TODO(), &list, client.InNamespace(deployment.Namespace), labels)
	if err != nil {
		log.Error(err, "Error getting list of Pods for StatefulSet "+sts.Name)
		return false
	}

	// All Pods must be in the Running Phase
	for _, pod := range list.Items {
		if !in.IsPodReady(pod) {
			logger.Info("Cannot execute probe, one or more Pods is not in a ready state")
			return false
		}
	}

	count := int32(len(list.Items))
	switch {
	case count == 0:
		logger.Info("Cannot find any Pods for StatefulSet " + sts.Name)
		return true
	case sts.Spec.Replicas == nil && count != 1:
		logger.Info(fmt.Sprintf("Pod count of %d does not yet match StatefulSet replica count: 1", count))
		return false
	case sts.Spec.Replicas != nil && count != *sts.Spec.Replicas:
		logger.Info(fmt.Sprintf("Pod count of %d does not yet match StatefulSet replica count: %d", count, *sts.Spec.Replicas))
		return false
	}

	for _, pod := range list.Items {
		if pod.Status.Phase == "Running" {
			if log.Enabled() {
				log.Info("Using pod " + pod.Name + " to execute probe")
			}

			ha, err := in.RunProbe(pod, probe)
			if err == nil {
				log.Info(fmt.Sprintf("Executed probe using pod %s result=%t", pod.Name, ha))
				return ha
			}
			log.Info(fmt.Sprintf("Execute probe using pod %s (%t) error %s", pod.Name, ha, err.Error()))
		} else {
			log.Info("Skipping execute probe for pod " + pod.Name + " as Pod status not in running phase")
		}
	}

	return false
}

// Determine whether the specified Pods are in the Ready state.
func (in *CoherenceProbe) IsPodReady(pod corev1.Pod) bool {
	if pod.DeletionTimestamp != nil || pod.Status.Phase != corev1.PodRunning {
		return false
	}

	for _, c := range pod.Status.Conditions {
		if c.Type == corev1.PodReady && c.Status == corev1.ConditionTrue {
			return true
		}
	}
	return false
}

func (in *CoherenceProbe) RunProbe(pod corev1.Pod, handler *coh.Probe) (bool, error) {
	switch {
	case handler.Exec != nil:
		return in.ProbeUsingExec(pod, handler)
	case handler.HTTPGet != nil:
		return in.ProbeUsingHTTP(pod, handler)
	case handler.TCPSocket != nil:
		return in.ProbeUsingTCP(pod, handler)
	default:
		return true, nil
	}
}

func (in *CoherenceProbe) ProbeUsingExec(pod corev1.Pod, handler *coh.Probe) (bool, error) {
	req := &mgmt.ExecRequest{
		Pod:       pod.Name,
		Container: coh.ContainerNameCoherence,
		Namespace: pod.Namespace,
		Command:   handler.Exec.Command,
		Arg:       []string{},
		Timeout:   handler.GetTimeout(),
	}

	exitCode, _, _, err := mgmt.PodExec(req, in.Config)

	log.Info(fmt.Sprintf("Exec Probe: '%s' result=%d error=%s", strings.Join(handler.Exec.Command, ", "), exitCode, err))

	if err != nil {
		return false, err
	}

	return exitCode == 0, nil
}

func (in *CoherenceProbe) ProbeUsingHTTP(pod corev1.Pod, handler *coh.Probe) (bool, error) {
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

	log.Info(fmt.Sprintf("HTTP Probe URL: %s result=%s msg=%s error=%s", u.String(), result, s, err))

	return result == probe.Success, err
}

func (in *CoherenceProbe) ProbeUsingTCP(pod corev1.Pod, handler *coh.Probe) (bool, error) {
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

	log.Info(fmt.Sprintf("TCP Probe: %s:%d result=%s error=%s", host, port, result, err))

	return result == probe.Success, err
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
