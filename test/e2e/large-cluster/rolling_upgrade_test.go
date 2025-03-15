/*
 * Copyright (c) 2020, 2025, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package large_cluster

import (
	"fmt"
	. "github.com/onsi/gomega"
	coh "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/pkg/operator"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	"slices"
	"testing"
	"time"
)

func TestUpgradeByNodeWithOnePodPerNode(t *testing.T) {
	testContext.CleanupAfterTest(t)

	ns := helper.GetTestNamespace()

	labels := map[string]string{}
	labels[operator.LabelTestHostName] = "127.0.0.1"
	labels[operator.LabelTestHealthPort] = fmt.Sprintf("%d", GetRestPort())

	replicas := 3

	c, _ := installSimpleDeployment(t, coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: ns,
			Name:      GenerateClusterName(),
		},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Replicas: ptr.To(int32(replicas)),
				Labels:   labels,
				ReadinessProbe: &coh.ReadinessProbeSpec{
					InitialDelaySeconds: ptr.To(int32(10)),
					PeriodSeconds:       ptr.To(int32(10)),
					FailureThreshold:    ptr.To(int32(20)),
				},
			},
			RollingUpdateStrategy: ptr.To(coh.UpgradeByNode),
		},
	})

	nodeNameGetter := func(pod corev1.Pod) string {
		return pod.Spec.NodeName
	}

	DoUpgradeTest(t, c, nodeNameGetter)
}

func TestUpgradeByNodeWithTwoPodsPerNode(t *testing.T) {
	testContext.CleanupAfterTest(t)

	ns := helper.GetTestNamespace()

	labels := map[string]string{}
	labels[operator.LabelTestHostName] = "127.0.0.1"
	labels[operator.LabelTestHealthPort] = fmt.Sprintf("%d", GetRestPort())

	replicas := len(nodeList.Items) * 2

	c, _ := installSimpleDeployment(t, coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: ns,
			Name:      GenerateClusterName(),
		},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Replicas: ptr.To(int32(replicas)),
				Labels:   labels,
				ReadinessProbe: &coh.ReadinessProbeSpec{
					InitialDelaySeconds: ptr.To(int32(10)),
					PeriodSeconds:       ptr.To(int32(10)),
					FailureThreshold:    ptr.To(int32(20)),
				},
			},
			RollingUpdateStrategy: ptr.To(coh.UpgradeByNode),
		},
	})

	nodeNameGetter := func(pod corev1.Pod) string {
		return pod.Spec.NodeName
	}

	DoUpgradeTest(t, c, nodeNameGetter)
}

func TestUpgradeByNodeLabelTwoPodsPerNode(t *testing.T) {
	testContext.CleanupAfterTest(t)

	ns := helper.GetTestNamespace()

	labels := map[string]string{}
	labels[operator.LabelTestHostName] = "127.0.0.1"
	labels[operator.LabelTestHealthPort] = fmt.Sprintf("%d", GetRestPort())

	replicas := len(nodeList.Items) * 2

	zoneLabel := "topology.kubernetes.io/zone"

	c, _ := installSimpleDeployment(t, coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: ns,
			Name:      GenerateClusterName(),
		},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Replicas: ptr.To(int32(replicas)),
				Labels:   labels,
				ReadinessProbe: &coh.ReadinessProbeSpec{
					InitialDelaySeconds: ptr.To(int32(10)),
					PeriodSeconds:       ptr.To(int32(10)),
					FailureThreshold:    ptr.To(int32(20)),
				},
			},
			RollingUpdateStrategy: ptr.To(coh.UpgradeByNodeLabel),
			RollingUpdateLabel:    &zoneLabel,
		},
	})

	nodeZoneGetter := func(pod corev1.Pod) string {
		node, ok := nodeMap[pod.Spec.NodeName]
		if ok {
			return node.Labels[zoneLabel]
		}
		return ""
	}

	DoUpgradeTest(t, c, nodeZoneGetter)
}

func DoUpgradeTest(t *testing.T, c coh.Coherence, idFunction func(corev1.Pod) string) {
	var err error
	g := NewWithT(t)

	replicas := int(c.GetReplicas())
	ns := c.Namespace

	before, err := helper.ListCoherencePodsForDeployment(testContext, c.Namespace, c.Name)
	g.Expect(err).To(BeNil())
	g.Expect(len(before)).To(Equal(replicas))

	if c.Spec.Labels == nil {
		c.Spec.Labels = make(map[string]string)
	}
	c.Spec.Labels["test"] = "one"

	err = testContext.Client.Update(testContext.Context, &c)
	g.Expect(err).NotTo(HaveOccurred())

	_, err = helper.WaitForPodsWithLabel(testContext, ns, "test=one", int(replicas), time.Second*5, time.Minute*10)
	g.Expect(err).NotTo(HaveOccurred())

	// Wait for the final Pod to be ready, i.e. all Pods in ready state
	_, err = helper.WaitForStatefulSetPodCondition(testContext, c.Namespace, c.Name, int32(replicas), corev1.PodReady, time.Second*10, time.Minute*5)
	g.Expect(err).To(BeNil())

	after, err := helper.ListCoherencePodsForDeployment(testContext, c.Namespace, c.Name)
	g.Expect(err).To(BeNil())
	g.Expect(len(after)).To(Equal(replicas))

	AssertPodRestartTimes(t, before, after, idFunction)
}

func AssertPodRestartTimes(t *testing.T, before, after []corev1.Pod, idFunction func(corev1.Pod) string) {
	g := NewWithT(t)

	// get the list of nodes the before pods were running on
	nodeNames := make([]string, 0)
	m := SplitPodsById(before, idFunction)
	for name := range m {
		nodeNames = append(nodeNames, name)
	}

	// dump data to logs for analysis on failure
	for _, pods := range m {
		for _, pod := range pods {
			ap, ok := FindPod(pod.Name, after)
			g.Expect(ok).To(BeTrue())
			t.Logf(">>>> Pod %s beforeNode %s beforeId %s afterNode %s afterId %s scheduled %v ready %v",
				pod.Name, pod.Spec.NodeName, idFunction(pod), ap.Spec.NodeName, idFunction(ap), GetPodScheduledTime(ap), GetPodReadyTime(ap))
		}
	}

	// for each after node get min scheduled and max ready and store by before node
	mapScheduled := make(map[string]time.Time)
	mapReady := make(map[string]time.Time)

	for _, pod := range after {
		scheduled := GetPodScheduledTime(pod)
		ready := GetPodReadyTime(pod)

		bp, ok := FindPod(pod.Name, before)
		g.Expect(ok).To(BeTrue())
		n := idFunction(bp)

		s, found := mapScheduled[n]
		if !found || s.After(scheduled) {
			mapScheduled[n] = scheduled
		}

		r, found := mapReady[n]
		if !found || r.Before(ready) {
			mapReady[n] = ready
		}
	}

	// make sure none of the ranges overlap
	for name, scheduled := range mapScheduled {
		ready := mapReady[name]
		for _, otherName := range nodeNames {
			if otherName != name {
				scheduledOther := mapScheduled[otherName]
				readyOther := mapReady[otherName]
				// compare scheduled time
				i := scheduled.Compare(scheduledOther)
				// must not be equal
				g.Expect(i).NotTo(BeZero())
				// ready must be the same comparison as scheduled,
				// i.e. if scheduled is before then ready must be before
				g.Expect(ready.Compare(readyOther)).To(Equal(i), fmt.Sprintf("node %s scheduled and ready overlap with node %s", name, otherName))
			}
		}
	}
}

func SplitPodsById(pods []corev1.Pod, idFunction func(corev1.Pod) string) map[string][]corev1.Pod {
	m := make(map[string][]corev1.Pod)
	for _, pod := range pods {
		nodeName := idFunction(pod)
		podsForNode, found := m[nodeName]
		if !found {
			podsForNode = make([]corev1.Pod, 0)
		}
		podsForNode = append(podsForNode, pod)
		m[nodeName] = podsForNode
	}
	return m
}

func SortPodsByScheduledTime(pods []corev1.Pod) {
	sorter := func(a, b corev1.Pod) int {
		aTime := GetPodScheduledTime(a)
		bTime := GetPodScheduledTime(b)
		return aTime.Compare(bTime)
	}
	slices.SortFunc(pods, sorter)
}

func GetPodScheduledTime(pod corev1.Pod) time.Time {
	for _, c := range pod.Status.Conditions {
		if c.Type == corev1.PodScheduled && c.Status == corev1.ConditionTrue {
			return c.LastTransitionTime.Time
		}
	}
	return time.Time{}
}

func GetPodReadyTime(pod corev1.Pod) time.Time {
	for _, c := range pod.Status.Conditions {
		if c.Type == corev1.PodReady && c.Status == corev1.ConditionTrue {
			return c.LastTransitionTime.Time
		}
	}
	return time.Time{}
}

func FindPod(name string, pods []corev1.Pod) (corev1.Pod, bool) {
	for _, pod := range pods {
		if pod.Name == name {
			return pod, true
		}
	}
	return corev1.Pod{}, false
}
