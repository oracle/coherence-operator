/*
 * Copyright (c) 2020, 2025, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

// Package helper contains various helpers for use in end-to-end testing.
package helper

import (
	goctx "context"
	"encoding/json"
	"fmt"
	. "github.com/onsi/gomega"
	coh "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/pkg/operator"
	"golang.org/x/net/context"
	"io"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/ptr"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"strings"
	"testing"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"time"
)

const operatorPodSelector = "control-plane=coherence"

var (
	RetryInterval = time.Second * 5
	Timeout       = time.Minute * 3
)

// WaitForStatefulSetForDeployment waits for a StatefulSet to be created for the specified deployment.
func WaitForStatefulSetForDeployment(ctx TestContext, namespace string, deployment *coh.Coherence, retryInterval, timeout time.Duration) (*appsv1.StatefulSet, error) {
	return WaitForStatefulSet(ctx, namespace, deployment.Name, deployment.Spec.GetReplicas(), retryInterval, timeout)
}

// WaitForDeploymentReady waits for a Coherence deployment to be ready.
func WaitForDeploymentReady(ctx TestContext, namespace, name string, retryInterval, timeout time.Duration) (*coh.Coherence, error) {
	var d = &coh.Coherence{}
	var key = types.NamespacedName{Namespace: namespace, Name: name}

	err := wait.PollUntilContextTimeout(ctx.Context, retryInterval, timeout, true, func(context.Context) (done bool, err error) {
		err = ctx.Client.Get(ctx.Context, key, d)
		if err != nil {
			if apierrors.IsNotFound(err) {
				ctx.Logf("Waiting for availability of Coherence deployment %s/%s - NotFound", namespace, name)
				return false, nil
			}
			ctx.Logf("Waiting for availability of Coherence deployment %s/%s - %s", namespace, name, err.Error())
			return false, err
		}

		if d.Status.Phase == coh.ConditionTypeReady {
			return true, nil
		}

		ctx.Logf("Waiting for Coherence deployment %s/%s to be Ready - %s (%d/%d)",
			namespace, name, d.Status.Phase, d.Status.ReadyReplicas, d.Status.Replicas)
		return false, nil
	})

	if err != nil && d != nil {
		d, _ := json.MarshalIndent(d, "", "    ")
		ctx.Logf("Error waiting for StatefulSet%s", string(d))
	}
	return d, err
}

// WaitForStatefulSet waits for a StatefulSet to be created with the specified number of replicas.
func WaitForStatefulSet(ctx TestContext, namespace, stsName string, replicas int32, retryInterval, timeout time.Duration) (*appsv1.StatefulSet, error) {
	var sts *appsv1.StatefulSet

	err := wait.PollUntilContextTimeout(ctx.Context, retryInterval, timeout, true, func(context.Context) (done bool, err error) {
		sts, err = ctx.KubeClient.AppsV1().StatefulSets(namespace).Get(ctx.Context, stsName, metav1.GetOptions{})
		if err != nil {
			if apierrors.IsNotFound(err) {
				ctx.Logf("Waiting for availability of StatefulSet %s - NotFound", stsName)
				return false, nil
			}
			ctx.Logf("Waiting for availability of %s StatefulSet - %s", stsName, err.Error())
			return false, err
		}

		if sts.Status.ReadyReplicas == replicas {
			ctx.Logf("StatefulSet %s replicas = (%d/%d)", stsName, sts.Status.ReadyReplicas, replicas)
			return true, nil
		}
		ctx.Logf("Waiting for full availability of StatefulSet %s (%d/%d)", stsName, sts.Status.ReadyReplicas, replicas)
		return false, nil
	})

	if err != nil && sts != nil {
		d, _ := json.MarshalIndent(sts, "", "    ")
		ctx.Logf("Error waiting for StatefulSet %s", string(d))
	}
	return sts, err
}

// WaitForStatefulSetPodCondition waits for all Pods in a StatefulSet to have the specified condition.
func WaitForStatefulSetPodCondition(ctx TestContext, namespace, stsName string, replicas int32, c corev1.PodConditionType, retryInterval, timeout time.Duration) (*appsv1.StatefulSet, error) {
	sts, err := WaitForStatefulSet(ctx, namespace, stsName, replicas, retryInterval, timeout)
	if err != nil {
		return nil, err
	}

	err = wait.PollUntilContextTimeout(ctx.Context, retryInterval, timeout, true, func(context.Context) (done bool, err error) {
		pods, err := ListPodsForStatefulSet(ctx, sts)
		if err != nil {
			if apierrors.IsNotFound(err) {
				ctx.Logf("Waiting for Pods in StatefulSet %s to reach condition %s - NotFound", stsName, c)
				return false, nil
			}
			ctx.Logf("Waiting for Pods in StatefulSet %s to reach condition %s - %s", stsName, c, err.Error())
			return false, err
		}

		ready := true
		count := 0
		for _, pod := range pods.Items {
			for _, cond := range pod.Status.Conditions {
				if cond.Type == c {
					if cond.Status == corev1.ConditionTrue {
						count++
					} else {
						ready = false
					}
				}
			}
		}

		ctx.Logf("Waiting for Pods in StatefulSet %s to reach condition %s = (%d/%d)", stsName, c, count, len(pods.Items))
		return ready, nil
	})

	return sts, err
}

// WaitForJob waits for a Job to be created with the specified number of replicas.
func WaitForJob(ctx TestContext, namespace, stsName string, replicas int32, retryInterval, timeout time.Duration) (*batchv1.Job, error) {
	var job *batchv1.Job

	err := wait.PollUntilContextTimeout(ctx.Context, retryInterval, timeout, true, func(context.Context) (done bool, err error) {
		job, err = ctx.KubeClient.BatchV1().Jobs(namespace).Get(ctx.Context, stsName, metav1.GetOptions{})
		if err != nil {
			if apierrors.IsNotFound(err) {
				ctx.Logf("Waiting for availability of Job %s - NotFound", stsName)
				return false, nil
			}
			ctx.Logf("Waiting for availability of %s Job - %s", stsName, err.Error())
			return false, err
		}

		var ready int32
		readyPtr := job.Status.Ready
		if readyPtr != nil {
			ready = *readyPtr
		} else {
			ready = job.Status.Succeeded
			if job.Status.Ready != nil {
				ready += job.Status.Active
			}
		}

		if ready == replicas {
			ctx.Logf("Job %s replicas = (%d/%d)", stsName, ready, replicas)
			return true, nil
		}
		ctx.Logf("Waiting for full availability of Job %s (%d/%d)", stsName, ready, replicas)
		return false, nil
	})

	if err != nil && job != nil {
		d, _ := json.MarshalIndent(job, "", "    ")
		ctx.Logf("Error waiting for Job %s", string(d))
	}
	return job, err
}

func WaitForCoherenceJobCondition(ctx TestContext, namespace, name string, condition DeploymentStateCondition, retryInterval, timeout time.Duration) (*coh.CoherenceJob, error) {
	var job *coh.CoherenceJob

	err := wait.PollUntilContextTimeout(ctx.Context, retryInterval, timeout, true, func(context.Context) (done bool, err error) {
		job, err = GetCoherenceJob(ctx, namespace, name)
		if err != nil {
			if apierrors.IsNotFound(err) {
				ctx.Logf("Waiting for availability of CoherenceJob resource %s - NotFound", name)
				return false, nil
			}
			ctx.Logf("Waiting for availability of CoherenceJob resource %s - %s", name, err.Error())
			return false, nil
		}
		valid := true
		if condition != nil {
			valid = condition.Test(job)
			if !valid {
				ctx.Logf("Waiting for CoherenceJob resource %s to meet condition '%s'", name, condition.String())
			}
		}
		return valid, nil
	})

	return job, err
}

// WaitForJobCompletion waits for a specified k8s Job to complete.
func WaitForJobCompletion(ctx TestContext, namespace, name string, retryInterval, timeout time.Duration) error {
	k8s := ctx.KubeClient
	err := wait.PollUntilContextTimeout(ctx.Context, retryInterval, timeout, true, func(context.Context) (done bool, err error) {
		p, err := k8s.CoreV1().Pods(namespace).Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		if len(p.Status.ContainerStatuses) > 0 {
			ready := true
			for _, s := range p.Status.ContainerStatuses {
				if s.State.Terminated == nil || s.State.Terminated.ExitCode != 0 {
					ready = false
					break
				}
			}
			return ready, nil
		}
		return false, nil
	})

	return err
}

// DeploymentStateCondition is a function that takes a deployment and determines whether it meets a condition
type DeploymentStateCondition interface {
	Test(coh.CoherenceResource) bool
	String() string
}

// An always true DeploymentStateCondition
type alwaysCondition struct{}

func (a alwaysCondition) Test(coh.CoherenceResource) bool {
	return true
}

func (a alwaysCondition) String() string {
	return "true"
}

type replicaCountCondition struct {
	replicas int32
}

func (in replicaCountCondition) Test(d coh.CoherenceResource) bool {
	return d.GetStatus().ReadyReplicas == in.replicas
}

func (in replicaCountCondition) String() string {
	return fmt.Sprintf("replicas == %d", in.replicas)
}

func ReplicaCountCondition(replicas int32) DeploymentStateCondition {
	return replicaCountCondition{replicas: replicas}
}

type phaseCondition struct {
	phase coh.ConditionType
}

func (in phaseCondition) Test(d coh.CoherenceResource) bool {
	return d.GetStatus().Phase == in.phase
}

func (in phaseCondition) String() string {
	return fmt.Sprintf("phase == %s", in.phase)
}

func StatusPhaseCondition(phase coh.ConditionType) DeploymentStateCondition {
	return phaseCondition{phase: phase}
}

type jobCompletedCondition struct {
	count int32
}

func (in jobCompletedCondition) Test(d coh.CoherenceResource) bool {
	status := d.GetStatus()
	return (status.Succeeded + status.Failed) == in.count
}

func (in jobCompletedCondition) String() string {
	return fmt.Sprintf("completed count == %d", in.count)
}

func JobCompletedCondition(count int32) DeploymentStateCondition {
	return jobCompletedCondition{count: count}
}

type jobSucceededCondition struct {
	count int32
}

func (in jobSucceededCondition) Test(d coh.CoherenceResource) bool {
	status := d.GetStatus()
	return status.Succeeded == in.count
}

func (in jobSucceededCondition) String() string {
	return fmt.Sprintf("succeeded count == %d", in.count)
}

func JobSucceededCondition(count int32) DeploymentStateCondition {
	return jobSucceededCondition{count: count}
}

type jobFailedCondition struct {
	count int32
}

func (in jobFailedCondition) Test(d coh.CoherenceResource) bool {
	status := d.GetStatus()
	return status.Failed == in.count
}

func (in jobFailedCondition) String() string {
	return fmt.Sprintf("failed count == %d", in.count)
}

func JobFailedCondition(count int32) DeploymentStateCondition {
	return jobFailedCondition{count: count}
}

// WaitForCoherence waits for a Coherence resource to be created.
//
//goland:noinspection GoUnusedExportedFunction
func WaitForCoherence(ctx TestContext, namespace, name string, retryInterval, timeout time.Duration) (*coh.Coherence, error) {
	return WaitForCoherenceCondition(ctx, namespace, name, alwaysCondition{}, retryInterval, timeout)
}

// WaitForCoherenceCondition waits for a Coherence resource to be created.
func WaitForCoherenceCondition(testCtx TestContext, namespace, name string, condition DeploymentStateCondition, retryInterval, timeout time.Duration) (*coh.Coherence, error) {
	var deployment *coh.Coherence

	ctx, _ := context.WithTimeout(testCtx.Context, timeout)

	err := wait.PollUntilContextTimeout(ctx, retryInterval, timeout, true, func(context.Context) (done bool, err error) {
		deployment, err = GetCoherence(testCtx, namespace, name)
		if err != nil {
			if apierrors.IsNotFound(err) {
				testCtx.Logf("Waiting for availability of Coherence resource %s - NotFound", name)
				return false, nil
			}
			testCtx.Logf("Waiting for availability of Coherence resource %s - %s", name, err.Error())
			return false, nil
		}
		valid := true
		if condition != nil {
			valid = condition.Test(deployment)
			if !valid {
				testCtx.Logf("Waiting for Coherence resource %s to meet condition '%s'", name, condition.String())
			}
		}
		return valid, nil
	})

	return deployment, err
}

// GetCoherence gets the specified Coherence resource
func GetCoherence(ctx TestContext, namespace, name string) (*coh.Coherence, error) {
	opts := client.ObjectKey{Namespace: namespace, Name: name}
	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
		},
	}

	err := ctx.Client.Get(ctx.Context, opts, d)

	return d, err
}

// GetCoherenceJob gets the specified CoherenceJob resource
func GetCoherenceJob(ctx TestContext, namespace, name string) (*coh.CoherenceJob, error) {
	opts := client.ObjectKey{Namespace: namespace, Name: name}
	d := &coh.CoherenceJob{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
		},
	}

	err := ctx.Client.Get(ctx.Context, opts, d)

	return d, err
}

// WaitForOperatorPods waits for a Coherence Operator Pods to be created.
func WaitForOperatorPods(ctx TestContext, namespace string, retryInterval, timeout time.Duration) ([]corev1.Pod, error) {
	return WaitForPodsWithSelector(ctx, namespace, operatorPodSelector, retryInterval, timeout)
}

func DeleteJob(ctx TestContext, namespace, jobName string) error {
	cl := ctx.KubeClient.BatchV1().Jobs(namespace)
	if err := cl.Delete(ctx.Context, jobName, metav1.DeleteOptions{PropagationPolicy: ptr.To(metav1.DeletePropagationBackground)}); err != nil && !errors.IsNotFound(err) {
		return err
	}
	pods, err := ListPodsWithLabelSelector(ctx, namespace, "job-name="+jobName)
	if err != nil {
		return err
	}

	for i := range pods {
		_ = ctx.Client.Delete(ctx.Context, &pods[i])
	}

	return nil
}

// WaitForPodsWithSelector waits for a Coherence Operator Pods to be created.
func WaitForPodsWithSelector(ctx TestContext, namespace, selector string, retryInterval, timeout time.Duration) ([]corev1.Pod, error) {
	var pods []corev1.Pod

	err := wait.PollUntilContextTimeout(ctx.Context, retryInterval, timeout, true, func(context.Context) (done bool, err error) {
		pods, err = ListPodsWithLabelSelector(ctx, namespace, selector)
		if err != nil {
			return false, err
		}
		return len(pods) > 0, nil
	})
	return pods, err
}

// WaitForPodsWithSelectorAndReplicas waits for Pods to be created.
func WaitForPodsWithSelectorAndReplicas(ctx TestContext, namespace, selector string, replicas int, retryInterval, timeout time.Duration) ([]corev1.Pod, error) {
	var pods []corev1.Pod

	err := wait.PollUntilContextTimeout(ctx.Context, retryInterval, timeout, true, func(context.Context) (done bool, err error) {
		pods, err = ListPodsWithLabelSelector(ctx, namespace, selector)
		if err != nil {
			return false, err
		}
		return len(pods) == replicas, nil
	})
	return pods, err
}

// WaitForOperatorDeletion waits for deletion of the Operator Pods.
//
//goland:noinspection GoUnusedExportedFunction
func WaitForOperatorDeletion(ctx TestContext, namespace string, retryInterval, timeout time.Duration) error {
	return WaitForDeleteOfPodsWithSelector(ctx, namespace, operatorPodSelector, retryInterval, timeout)
}

// WaitForDeleteOfPodsWithSelector waits for Pods to be deleted.
func WaitForDeleteOfPodsWithSelector(ctx TestContext, namespace, selector string, retryInterval, timeout time.Duration) error {
	ctx.Logf("Waiting for Pods in namespace %s with selector '%s' to be deleted", namespace, selector)
	var pods []corev1.Pod

	err := wait.PollUntilContextTimeout(ctx.Context, retryInterval, timeout, true, func(context.Context) (done bool, err error) {
		ctx.Logf("List Pods in namespace %s with selector '%s'", namespace, selector)
		pods, err = ListPodsWithLabelSelector(ctx, namespace, selector)
		if err != nil {
			ctx.Logf("Error listing Pods in namespace %s with selector '%s' - %s", namespace, selector, err.Error())
			return false, err
		}
		ctx.Logf("Found %d Pods in namespace %s with selector '%s'", len(pods), namespace, selector)
		return len(pods) == 0, nil
	})

	ctx.Logf("Finished waiting for Pods in namespace %s with selector '%s' to be deleted. Error=%s", namespace, selector, err)
	return err
}

// WaitForDeletion waits for deletion of the specified resource.
func WaitForDeletion(ctx TestContext, namespace, name string, resource client.Object, retryInterval, timeout time.Duration) error {
	gvk, _ := apiutil.GVKForObject(resource, ctx.Manager.GetScheme())
	ctx.Logf("Waiting for deletion of %v %s/%s", gvk, namespace, name)

	key := types.NamespacedName{Namespace: namespace, Name: name}

	err := wait.PollUntilContextTimeout(ctx.Context, retryInterval, timeout, true, func(context.Context) (done bool, err error) {
		err = ctx.Client.Get(ctx.Context, key, resource)
		switch {
		case err != nil && errors.IsNotFound(err):
			return true, nil
		case err != nil && !errors.IsNotFound(err):
			ctx.Logf("Waiting for deletion of %v %s/%s - Error=%s", gvk, namespace, name, err)
			return false, err
		default:
			ctx.Logf("Still waiting for deletion of %v %s/%s", gvk, namespace, name)
			return false, nil
		}
	})

	ctx.Logf("Finished waiting for deletion of %s %s/%s - Error=%s", gvk, namespace, name, err)
	return err
}

// ListOperatorPods lists the Operator Pods that exist - this is Pods with the label "name=coh-operator"
func ListOperatorPods(ctx TestContext, namespace string) ([]corev1.Pod, error) {
	return ListPodsWithLabelSelector(ctx, namespace, operatorPodSelector)
}

// ListCoherencePodsForCluster lists the Coherence Cluster Pods that exist for a cluster - this is Pods with the label "coherenceCluster=<cluster>"
func ListCoherencePodsForCluster(ctx TestContext, namespace, cluster string) ([]corev1.Pod, error) {
	return ListPodsWithLabelSelector(ctx, namespace, fmt.Sprintf("%s=%s", coh.LabelCoherenceCluster, cluster))
}

// WaitForPodsWithLabel waits for at least the required number of Pods matching the specified labels selector to be created.
func WaitForPodsWithLabel(ctx TestContext, namespace, selector string, count int, retryInterval, timeout time.Duration) ([]corev1.Pod, error) {
	var pods []corev1.Pod

	err := wait.PollUntilContextTimeout(ctx.Context, retryInterval, timeout, true, func(context.Context) (done bool, err error) {
		pods, err = ListPodsWithLabelSelector(ctx, namespace, selector)
		if err != nil {
			ctx.Logf("Waiting for at least %d Pods with label selector '%s' - failed due to %s", count, selector, err.Error())
			return false, err
		}
		found := len(pods) >= count
		if !found {
			ctx.Logf("Waiting for at least %d Pods with label selector '%s' - found %d", count, selector, len(pods))
		}
		return found, nil
	})

	return pods, err
}

// WaitForJobsWithLabel waits for at least the required number of Jobs matching the specified labels selector to be created.
func WaitForJobsWithLabel(ctx TestContext, namespace, selector string, count int, retryInterval, timeout time.Duration) ([]batchv1.Job, error) {
	var jobs []batchv1.Job

	err := wait.PollUntilContextTimeout(ctx.Context, retryInterval, timeout, true, func(context.Context) (done bool, err error) {
		jobs, err = ListJobsWithLabelSelector(ctx, namespace, selector)
		if err != nil {
			ctx.Logf("Waiting for at least %d Jobs with label selector '%s' - failed due to %s", count, selector, err.Error())
			return false, err
		}
		found := len(jobs) >= count
		if !found {
			ctx.Logf("Waiting for at least %d Jobs with label selector '%s' - found %d", count, selector, len(jobs))
		}
		return found, nil
	})

	return jobs, err
}

// WaitForPodsWithLabelAndField waits for at least the required number of pending Pods
func WaitForPodsWithLabelAndField(ctx TestContext, namespace, labelSelector, fieldSelector string, count int, retryInterval, timeout time.Duration) ([]corev1.Pod, error) {
	var pods []corev1.Pod

	err := wait.PollUntilContextTimeout(ctx.Context, retryInterval, timeout, true, func(context.Context) (done bool, err error) {
		pods, err = ListPodsWithLabelAndFieldSelector(ctx, namespace, labelSelector, fieldSelector)
		if err != nil {
			ctx.Logf("Waiting for %d Pods with label selector '%s' and field selector '%s' - failed due to %s", count, labelSelector, fieldSelector, err.Error())
			return false, err
		}

		found := len(pods) >= count
		if !found {
			ctx.Logf("Waiting for %d Pods with label selector '%s' and field selector '%s' - found %d", count, labelSelector, fieldSelector, len(pods))
			return found, nil
		}

		return found, nil
	})

	return pods, err
}

// ListCoherencePodsForDeployment lists the Pods that exist for a deployment - this is Pods with the label "coherenceDeployment=<deployment>"
func ListCoherencePodsForDeployment(ctx TestContext, namespace, deployment string) ([]corev1.Pod, error) {
	selector := fmt.Sprintf("%s=%s", coh.LabelCoherenceDeployment, deployment)
	return ListPodsWithLabelSelector(ctx, namespace, selector)
}

// ListCoherencePods lists the all Coherence deployment Pods in a namespace
func ListCoherencePods(ctx TestContext, namespace string) ([]corev1.Pod, error) {
	selector := fmt.Sprintf("%s=%s", coh.LabelComponent, coh.LabelComponentCoherencePod)
	return ListPodsWithLabelSelector(ctx, namespace, selector)
}

// ListPodsWithLabelSelector lists the Coherence Cluster Pods that exist for a given label selector.
func ListPodsWithLabelSelector(ctx TestContext, namespace, selector string) ([]corev1.Pod, error) {
	opts := metav1.ListOptions{LabelSelector: selector}

	list, err := ctx.KubeClient.CoreV1().Pods(namespace).List(ctx.Context, opts)
	if err != nil {
		return []corev1.Pod{}, err
	}

	return list.Items, nil
}

// ListJobsWithLabelSelector lists the Coherence Cluster Jobs that exist for a given label selector.
func ListJobsWithLabelSelector(ctx TestContext, namespace, selector string) ([]batchv1.Job, error) {
	opts := metav1.ListOptions{LabelSelector: selector}

	list, err := ctx.KubeClient.BatchV1().Jobs(namespace).List(ctx.Context, opts)
	if err != nil {
		return []batchv1.Job{}, err
	}

	return list.Items, nil
}

// ListPodsWithLabelAndFieldSelector lists the Coherence Cluster Pods that exist for a given label and field selectors.
func ListPodsWithLabelAndFieldSelector(ctx TestContext, namespace, labelSelector, fieldSelector string) ([]corev1.Pod, error) {
	opts := metav1.ListOptions{LabelSelector: labelSelector, FieldSelector: fieldSelector}

	list, err := ctx.KubeClient.CoreV1().Pods(namespace).List(ctx.Context, opts)
	if err != nil {
		return []corev1.Pod{}, err
	}

	return list.Items, nil
}

func ListPodsForStatefulSet(ctx TestContext, sts *appsv1.StatefulSet) (corev1.PodList, error) {
	pods := corev1.PodList{}
	var replicas int
	if sts.Spec.Replicas == nil {
		replicas = 1
	} else {
		replicas = int(*sts.Spec.Replicas)
	}

	name := types.NamespacedName{Namespace: sts.Namespace}
	for i := 0; i < replicas; i++ {
		name.Name = fmt.Sprintf("%s-%d", sts.Name, i)
		pod := corev1.Pod{}
		err := ctx.Client.Get(ctx.Context, name, &pod)
		if err != nil {
			if apierrors.IsNotFound(err) {
				t := metav1.Now()
				pod.Namespace = name.Namespace
				pod.Name = name.Name
				pod.DeletionTimestamp = &t
			} else {
				return pods, err
			}
		}
		pods.Items = append(pods.Items, pod)
	}
	return pods, nil
}

// WaitForPodReady waits for a Pods to be ready.
//
//goland:noinspection GoUnusedExportedFunction
func WaitForPodReady(ctx TestContext, namespace, name string, retryInterval, timeout time.Duration) error {
	err := wait.PollUntilContextTimeout(ctx.Context, retryInterval, timeout, true, func(context.Context) (done bool, err error) {
		k8s := ctx.KubeClient
		p, err := k8s.CoreV1().Pods(namespace).Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		if len(p.Status.ContainerStatuses) > 0 {
			ready := true
			for _, s := range p.Status.ContainerStatuses {
				if !s.Ready {
					ready = false
					break
				}
			}
			return ready, nil
		}
		return false, nil
	})

	return err
}

// WaitForCoherenceCleanup waits until there are no Coherence resources left in the test namespace.
// The default clean-up hooks only wait for deletion of resources directly created via the test client
func WaitForCoherenceCleanup(ctx TestContext, namespace string) error {
	ctx.Logf("Waiting for clean-up of Coherence resources in namespace %s", namespace)

	err := waitForCoherenceJobCleanup(ctx, namespace)
	if err != nil {
		return err
	}

	list := &coh.CoherenceList{}
	err = ctx.Client.List(goctx.TODO(), list, client.InNamespace(namespace))
	if err != nil {
		return err
	}

	// Do a plain delete first
	for i := range list.Items {
		item := list.Items[i]
		ctx.Logf("Deleting Coherence resource %s in namespace %s", item.Name, item.Namespace)
		err = ctx.Client.Delete(goctx.TODO(), &item)
		if err != nil {
			ctx.Logf("Error deleting Coherence resource %s - %s", item.Name, err.Error())
		}
	}

	// Obtain any remaining Coherence resources
	err = ctx.Client.List(goctx.TODO(), list, client.InNamespace(namespace))
	if err != nil {
		return err
	}

	// Delete all the Coherence resources - patching out any finalizer
	patch := client.RawPatch(types.MergePatchType, []byte(`{"metadata":{"finalizers":[]}}`))
	for i := range list.Items {
		item := list.Items[i]
		ctx.Logf("Patching Coherence resource %s in namespace %s to remove finalizers", item.Name, item.Namespace)
		if err := ctx.Client.Patch(ctx.Context, &item, patch); err != nil {
			ctx.Logf("error patching Coherence %s: %+v", item.Name, err)
		}
		ctx.Logf("Deleting Coherence resource %s in namespace %s", item.Name, item.Namespace)
		err = ctx.Client.Delete(goctx.TODO(), &item)
		if err != nil {
			ctx.Logf("Error deleting Coherence resource %s - %s", item.Name, err.Error())
		}
	}

	// Wait for removal of the Coherence resources
	err = wait.PollUntilContextTimeout(ctx.Context, RetryInterval, Timeout, true, func(context.Context) (done bool, err error) {
		err = ctx.Client.List(goctx.TODO(), list, client.InNamespace(namespace))
		if err == nil || isNoResources(err) || errors.IsNotFound(err) {
			if len(list.Items) > 0 {
				ctx.Logf("Waiting for deletion of %d Coherence resources", len(list.Items))
				return false, nil
			}
			return true, nil
		}
		ctx.Logf("Error waiting for deletion of Coherence resources: %s\n%+v", err.Error(), err)
		return false, nil
	})

	if err == nil {
		// wait for all Coherence Pods to be deleted
		err = wait.PollUntilContextTimeout(ctx.Context, RetryInterval, Timeout, true, func(context.Context) (done bool, err error) {
			list, err := ListCoherencePods(ctx, namespace)
			if err == nil {
				if len(list) > 0 {
					ctx.Logf("Waiting for deletion of %d Coherence Pods", len(list))
					return false, nil
				}
				return true, nil
			}
			ctx.Logf("Error waiting for deletion of Coherence Pods: %s", err.Error())
			return false, nil
		})
	}

	return err
}

// waitForCoherenceJobCleanup waits until there are no CoherenceJob resources left in the test namespace.
// The default clean-up hooks only wait for deletion of resources directly created via the test client
func waitForCoherenceJobCleanup(ctx TestContext, namespace string) error {
	ctx.Logf("Waiting for clean-up of CoherenceJob resources in namespace %s", namespace)

	list := &coh.CoherenceJobList{}
	err := ctx.Client.List(goctx.TODO(), list, client.InNamespace(namespace))
	if err != nil {
		return err
	}

	// Do a plain delete first
	for i := range list.Items {
		item := list.Items[i]
		ctx.Logf("Deleting CoherenceJob resource %s in namespace %s", item.Name, item.Namespace)
		err = ctx.Client.Delete(goctx.TODO(), &item)
		if err != nil {
			ctx.Logf("Error deleting CoherenceJob resource %s - %s", item.Name, err.Error())
		}
	}

	// Obtain any remaining CoherenceJob resources
	err = ctx.Client.List(goctx.TODO(), list, client.InNamespace(namespace))
	if err != nil {
		return err
	}

	// Delete all the CoherenceJob resources - patching out any finalizer
	patch := client.RawPatch(types.MergePatchType, []byte(`{"metadata":{"finalizers":[]}}`))
	for i := range list.Items {
		item := list.Items[i]
		ctx.Logf("Patching CoherenceJob resource %s in namespace %s to remove finalizers", item.Name, item.Namespace)
		if err := ctx.Client.Patch(ctx.Context, &item, patch); err != nil {
			ctx.Logf("error patching CoherenceJob %s: %+v", item.Name, err)
		}
		ctx.Logf("Deleting CoherenceJob resource %s in namespace %s", item.Name, item.Namespace)
		err = ctx.Client.Delete(goctx.TODO(), &item)
		if err != nil {
			ctx.Logf("Error deleting CoherenceJob resource %s - %s", item.Name, err.Error())
		}
	}

	// Wait for removal of the CoherenceJob resources
	err = wait.PollUntilContextTimeout(ctx.Context, RetryInterval, Timeout, true, func(context.Context) (done bool, err error) {
		err = ctx.Client.List(goctx.TODO(), list, client.InNamespace(namespace))
		if err == nil || isNoResources(err) || errors.IsNotFound(err) {
			if len(list.Items) > 0 {
				ctx.Logf("Waiting for deletion of %d CoherenceJob resources", len(list.Items))
				return false, nil
			}
			return true, nil
		}
		ctx.Logf("Error waiting for deletion of CoherenceJob resources: %s\n%+v", err.Error(), err)
		return false, nil
	})

	return err
}

func isNoResources(err error) bool {
	return err != nil && strings.HasPrefix(err.Error(), "no matches for kind")
}

// WaitForOperatorCleanup waits until there are no Operator Pods in the test namespace.
func WaitForOperatorCleanup(ctx TestContext, namespace string) error {
	ctx.Logf("Waiting for deletion of Coherence Operator Pods")
	// wait for all Operator Pods to be deleted
	err := wait.PollUntilContextTimeout(ctx.Context, RetryInterval, Timeout, true, func(context.Context) (done bool, err error) {
		list, err := ListOperatorPods(ctx, namespace)
		if err == nil {
			if len(list) > 0 {
				ctx.Logf("Waiting for deletion of %d Coherence Operator Pods", len(list))
				return false, nil
			}
			return true, nil
		}
		ctx.Logf("Error waiting for deletion of Coherence Operator Pods: %s", err.Error())
		return false, nil
	})

	ctx.Logger.Info("Coherence Operator Pods deleted")
	return err
}

// DumpOperatorLog dumps the Operator Pod log to a file.
func DumpOperatorLog(ctx TestContext, namespace, directory string) {
	list, err := ctx.KubeClient.CoreV1().Pods(namespace).List(ctx.Context, metav1.ListOptions{LabelSelector: "name=coherence-operator"})
	if err == nil {
		if len(list.Items) > 0 {
			pod := list.Items[0]
			DumpPodLog(ctx, &pod, directory)
		} else {
			ctx.Logger.Info("Could not capture Operator Pod log. No Pods found.")
		}
	}

	if err != nil {
		ctx.Logf("Could not capture Operator Pod log due to error: %s", err.Error())
	}
}

// DumpPodLog dumps the Pod log to a file.
func DumpPodLog(ctx TestContext, pod *corev1.Pod, directory string) {
	logs, err := FindTestLogsDir()
	if err != nil {
		ctx.Logger.Info("cannot capture logs due to " + err.Error())
		return
	}

	ctx.Logger.Info("Capturing Pod logs for " + pod.Name)

	pathSep := string(os.PathSeparator)
	name := logs + pathSep + directory
	err = os.MkdirAll(name, os.ModePerm)
	if err != nil {
		ctx.Logger.Info("cannot capture logs for Pod " + pod.Name + " due to " + err.Error())
	}

	for _, container := range pod.Spec.InitContainers {
		DumpContainerLogs(ctx, container, pod, name)
	}
	for _, container := range pod.Spec.Containers {
		DumpContainerLogs(ctx, container, pod, name)
	}
}

// DumpContainerLogs dumps the logs for a container
func DumpContainerLogs(ctx TestContext, container corev1.Container, pod *corev1.Pod, directory string) {
	var err error
	pathSep := string(os.PathSeparator)
	res := ctx.KubeClient.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, &corev1.PodLogOptions{Container: container.Name})
	s, err := res.Stream(ctx.Context)
	if err == nil {
		suffix := 0
		logName := fmt.Sprintf("%s%s%s(%s).log", directory, pathSep, pod.Name, container.Name)
		_, err = os.Stat(logName)
		for err == nil {
			suffix++
			logName = fmt.Sprintf("%s%s%s(%s)-%d.log", directory, pathSep, pod.Name, container.Name, suffix)
			_, err = os.Stat(logName)
		}
		out, err := os.Create(logName)
		if err == nil {
			if _, err = io.Copy(out, s); err != nil {
				ctx.Logger.Info("cannot capture logs for Pod " + pod.Name + " container " + container.Name + " due to " + err.Error())
			}
		} else {
			ctx.Logger.Info("cannot capture logs for Pod " + pod.Name + " container " + container.Name + " due to " + err.Error())
		}
	} else {
		ctx.Logger.Info("cannot capture logs for Pod " + pod.Name + " container " + container.Name + " due to " + err.Error())
	}
}

// GetTestSslSecret gets the test k8s secret that can be used for SSL testing.
func GetTestSslSecret() (*OperatorSSL, *coh.SSLSpec, error) {
	return CreateSslSecret(nil, GetTestNamespace(), GetTestSSLSecretName())
}

// CreateSslSecret creates a k8s secret that can be used for SSL testing.
func CreateSslSecret(kubeClient kubernetes.Interface, namespace, name string) (*OperatorSSL, *coh.SSLSpec, error) {
	certs, err := FindTestCertsDir()
	if err != nil {
		return nil, nil, err
	}

	// ensure the certs dir exists
	_, err = os.Stat(certs)
	if err != nil {
		return nil, nil, err
	}

	keystore := "keystore.jks"
	storepass := "storepass.txt"
	keypass := "keypass.txt"
	truststore := "truststore.jks"
	trustpass := "trustpass.txt"
	keyFile := "operator.key"
	certFile := "operator.crt"
	caCert := "operator-ca.crt"

	secret := &corev1.Secret{}
	secret.SetNamespace(namespace)
	secret.SetName(name)

	secret.Data = make(map[string][]byte)

	opSSL := OperatorSSL{}

	opSSL.Secrets = &name
	opSSL.KeyFile = &keyFile
	secret.Data[keyFile], err = readCertFile(certs + "/groot.key")
	if err != nil {
		return nil, nil, err
	}

	opSSL.CertFile = &certFile
	secret.Data[certFile], err = readCertFile(certs + "/groot.crt")
	if err != nil {
		return nil, nil, err
	}

	opSSL.CaFile = &caCert
	secret.Data[caCert], err = readCertFile(certs + "/guardians-ca.crt")
	if err != nil {
		return nil, nil, err
	}

	cohSSL := coh.SSLSpec{}

	cohSSL.Secrets = &name
	cohSSL.KeyStore = &keystore
	secret.Data[keystore], err = readCertFile(certs + "/groot.jks")
	if err != nil {
		return nil, nil, err
	}

	cohSSL.KeyStorePasswordFile = &storepass
	secret.Data[storepass], err = readCertFile(certs + "/storepassword.txt")
	if err != nil {
		return nil, nil, err
	}

	cohSSL.KeyPasswordFile = &keypass
	secret.Data[keypass], err = readCertFile(certs + "/keypassword.txt")
	if err != nil {
		return nil, nil, err
	}

	cohSSL.TrustStore = &truststore
	secret.Data[keypass], err = readCertFile(certs + "/truststore-guardians.jks")
	if err != nil {
		return nil, nil, err
	}

	cohSSL.TrustStorePasswordFile = &trustpass
	secret.Data[trustpass], err = readCertFile(certs + "/trustpassword.txt")
	if err != nil {
		return nil, nil, err
	}

	// We do not want to overwrite the existing test secret
	if kubeClient != nil && name != GetTestSSLSecretName() {
		_, err = kubeClient.CoreV1().Secrets(namespace).Create(context.TODO(), secret, metav1.CreateOptions{})
	}

	return &opSSL, &cohSSL, err
}

func readCertFile(name string) ([]byte, error) {
	f, err := os.Create(name)
	if err != nil {
		return nil, err
	}

	_, err = f.Stat()
	if err != nil {
		return nil, err
	}

	return io.ReadAll(f)
}

// DumpOperatorLogs dumps the operator logs
func DumpOperatorLogs(t *testing.T, ctx TestContext) {
	namespace := GetTestNamespace()
	DumpOperatorLog(ctx, namespace, t.Name())
	DumpState(ctx, namespace, t.Name())
	clusterNamespace := GetTestClusterNamespace()
	if clusterNamespace != namespace {
		DumpState(ctx, clusterNamespace, t.Name()+string(os.PathSeparator)+clusterNamespace)
	}
}

func DumpState(ctx TestContext, namespace, dir string) {
	dumpEvents(namespace, dir, ctx)
	dumpCoherences(namespace, dir, ctx)
	dumpStatefulSets(namespace, dir, ctx)
	dumpJobs(namespace, dir, ctx)
	dumpServices(namespace, dir, ctx)
	dumpPods(namespace, dir, ctx)
	dumpRbacRoles(namespace, dir, ctx)
	dumpRbacRoleBindings(namespace, dir, ctx)
	dumpServiceAccounts(namespace, dir, ctx)
	dumpClientNamespace(dir, ctx)
}

func dumpClientNamespace(dir string, ctx TestContext) {
	namespace := GetTestClientNamespace()
	_, err := ctx.KubeClient.CoreV1().Namespaces().Get(ctx.Context, namespace, metav1.GetOptions{})
	if err == nil {
		dumpPods(namespace, dir, ctx)
		dumpJobs(namespace, dir, ctx)
	}
}

func dumpEvents(namespace, dir string, ctx TestContext) {
	const message = "Could not dump events for namespace %s due to %s"
	list, err := ctx.KubeClient.CoreV1().Events(namespace).List(ctx.Context, metav1.ListOptions{})
	if err != nil {
		ctx.Logf(message, namespace, err.Error())
		return
	}

	logsDir, err := EnsureLogsDir(dir)
	if err != nil {
		ctx.Logf(message, namespace, err.Error())
		return
	}

	listFile, err := os.Create(logsDir + string(os.PathSeparator) + "events.json")
	if err != nil {
		ctx.Logf(message, namespace, err.Error())
		return
	}

	_, err = fmt.Fprint(listFile, "[\n")
	if err != nil {
		ctx.Logf(message, namespace, err.Error())
		return
	}

	fn := func() { _ = listFile.Close() }
	defer fn()

	if len(list.Items) > 0 {
		for i, item := range list.Items {
			if i > 0 {
				_, err = fmt.Fprint(listFile, ",\n")
				if err != nil {
					ctx.Logf(message, namespace, err.Error())
					return
				}
			}

			d, err := json.MarshalIndent(item, "", "    ")
			if err != nil {
				ctx.Logf(message, namespace, err.Error())
				return
			}
			_, err = fmt.Fprint(listFile, string(d))
			if err != nil {
				ctx.Logf(message, namespace, err.Error())
				return
			}
		}
	}

	_, err = fmt.Fprint(listFile, "]")
	if err != nil {
		ctx.Logf(message, namespace, err.Error())
		return
	}
}

func dumpCoherences(namespace, dir string, ctx TestContext) {
	const message = "Could not dump Coherence resource for namespace %s due to %s"

	list := coh.CoherenceList{}
	err := ctx.Client.List(ctx.Context, &list, client.InNamespace(namespace))
	if err != nil {
		ctx.Logf(message, namespace, err.Error())
		return
	}

	logsDir, err := EnsureLogsDir(dir)
	if err != nil {
		ctx.Logf(message, namespace, err.Error())
		return
	}

	listFile, err := os.Create(logsDir + string(os.PathSeparator) + "deployments-list.txt")
	if err != nil {
		ctx.Logf(message, namespace, err.Error())
		return
	}

	fn := func() { _ = listFile.Close() }
	defer fn()

	if len(list.Items) > 0 {
		for _, item := range list.Items {
			_, err = fmt.Fprint(listFile, item.GetName()+"\n")
			if err != nil {
				ctx.Logf(message, namespace, err.Error())
				return
			}

			d, err := json.MarshalIndent(item, "", "    ")
			if err != nil {
				ctx.Logf(message, namespace, err.Error())
				return
			}

			file, err := os.Create(logsDir + string(os.PathSeparator) + "Coherence-" + item.GetName() + ".json")
			if err != nil {
				ctx.Logf(message, namespace, err.Error())
				return
			}

			_, err = fmt.Fprint(file, string(d))
			_ = file.Close()
			if err != nil {
				ctx.Logf(message, namespace, err.Error())
				return
			}
		}
	} else {
		_, _ = fmt.Fprint(listFile, "No Coherence resources found in namespace "+namespace)
	}
}

func dumpStatefulSets(namespace, dir string, ctx TestContext) {
	const message = "Could not dump StatefulSets for namespace %s due to %s"

	list, err := ctx.KubeClient.AppsV1().StatefulSets(namespace).List(ctx.Context, metav1.ListOptions{})
	if err != nil {
		ctx.Logf(message, namespace, err.Error())
		return
	}

	logsDir, err := EnsureLogsDir(dir)
	if err != nil {
		ctx.Logf(message, namespace, err.Error())
		return
	}

	listFile, err := os.Create(logsDir + string(os.PathSeparator) + "sts-list.txt")
	if err != nil {
		ctx.Logf(message, namespace, err.Error())
		return
	}

	fn := func() { _ = listFile.Close() }
	defer fn()

	if len(list.Items) > 0 {
		for _, item := range list.Items {
			_, err = fmt.Fprint(listFile, item.GetName()+"\n")
			if err != nil {
				ctx.Logf(message, namespace, err.Error())
				return
			}

			d, err := json.MarshalIndent(item, "", "    ")
			if err != nil {
				ctx.Logf(message, namespace, err.Error())
				return
			}

			file, err := os.Create(logsDir + string(os.PathSeparator) + "StatefulSet-" + item.GetName() + ".json")
			if err != nil {
				ctx.Logf(message, namespace, err.Error())
				return
			}

			_, err = fmt.Fprint(file, string(d))
			_ = file.Close()
			if err != nil {
				ctx.Logf(message, namespace, err.Error())
				return
			}
		}
	} else {
		_, _ = fmt.Fprint(listFile, "No StatefulSet resources found in namespace "+namespace)
	}
}

func dumpServices(namespace, dir string, ctx TestContext) {
	const message = "Could not dump Services for namespace %s due to %s"

	list, err := ctx.KubeClient.CoreV1().Services(namespace).List(ctx.Context, metav1.ListOptions{})
	if err != nil {
		ctx.Logf(message, namespace, err.Error())
		return
	}

	logsDir, err := EnsureLogsDir(dir)
	if err != nil {
		ctx.Logf(message, namespace, err.Error())
		return
	}

	listFile, err := os.Create(logsDir + string(os.PathSeparator) + "svc-list.txt")
	if err != nil {
		ctx.Logf(message, namespace, err.Error())
		return
	}

	fn := func() { _ = listFile.Close() }
	defer fn()

	if len(list.Items) > 0 {
		for _, item := range list.Items {
			_, err = fmt.Fprint(listFile, item.GetName()+"\n")
			if err != nil {
				ctx.Logf(message, namespace, err.Error())
				return
			}

			d, err := json.MarshalIndent(item, "", "    ")
			if err != nil {
				ctx.Logf(message, namespace, err.Error())
				return
			}

			file, err := os.Create(logsDir + string(os.PathSeparator) + "Service-" + item.GetName() + ".json")
			if err != nil {
				ctx.Logf(message, namespace, err.Error())
				return
			}

			_, err = fmt.Fprint(file, string(d))
			_ = file.Close()
			if err != nil {
				ctx.Logf(message, namespace, err.Error())
				return
			}
		}
	} else {
		_, _ = fmt.Fprintf(listFile, "No Service resources found in namespace %s", namespace)
	}
}

func dumpJobs(namespace, dir string, ctx TestContext) {
	const message = "Could not dump Jobs for namespace %s due to %s"

	list, err := ctx.KubeClient.BatchV1().Jobs(namespace).List(ctx.Context, metav1.ListOptions{})
	if err != nil {
		ctx.Logf(message, namespace, err.Error())
		return
	}

	logsDir, err := EnsureLogsDir(dir)
	if err != nil {
		ctx.Logf(message, namespace, err.Error())
		return
	}

	listFile, err := os.Create(logsDir + string(os.PathSeparator) + "jobs-list.txt")
	if err != nil {
		ctx.Logf(message, namespace, err.Error())
		return
	}

	fn := func() { _ = listFile.Close() }
	defer fn()

	if len(list.Items) > 0 {
		for _, item := range list.Items {
			_, err = fmt.Fprint(listFile, item.GetName()+"\n")
			if err != nil {
				ctx.Logf(message, namespace, err.Error())
				return
			}

			d, err := json.MarshalIndent(item, "", "    ")
			if err != nil {
				ctx.Logf(message, namespace, err.Error())
				return
			}

			file, err := os.Create(logsDir + string(os.PathSeparator) + "Job-" + item.GetName() + ".json")
			if err != nil {
				ctx.Logf(message, namespace, err.Error())
				return
			}

			_, err = fmt.Fprint(file, string(d))
			_ = file.Close()
			if err != nil {
				ctx.Logf(message, namespace, err.Error())
				return
			}
		}
	} else {
		_, _ = fmt.Fprintf(listFile, "No Job resources found in namespace %s", namespace)
	}
}

func dumpRbacRoles(namespace, dir string, ctx TestContext) {
	const message = "Could not dump RBAC Roles for namespace %s due to %s"

	list, err := ctx.KubeClient.RbacV1().Roles(namespace).List(ctx.Context, metav1.ListOptions{})
	if err != nil {
		ctx.Logf(message, namespace, err.Error())
		return
	}

	logsDir, err := EnsureLogsDir(dir)
	if err != nil {
		ctx.Logf(message, namespace, err.Error())
		return
	}

	listFile, err := os.Create(logsDir + string(os.PathSeparator) + "rbac-role-list.txt")
	if err != nil {
		ctx.Logf(message, namespace, err.Error())
		return
	}

	fn := func() { _ = listFile.Close() }
	defer fn()

	if len(list.Items) > 0 {
		for _, item := range list.Items {
			_, err = fmt.Fprint(listFile, item.GetName()+"\n")
			if err != nil {
				ctx.Logf(message, namespace, err.Error())
				return
			}

			d, err := json.MarshalIndent(item, "", "    ")
			if err != nil {
				ctx.Logf(message, namespace, err.Error())
				return
			}

			file, err := os.Create(logsDir + string(os.PathSeparator) + "RBAC-Role-" + item.GetName() + ".json")
			if err != nil {
				ctx.Logf(message, namespace, err.Error())
				return
			}

			_, err = fmt.Fprint(file, string(d))
			_ = file.Close()
			if err != nil {
				ctx.Logf(message, namespace, err.Error())
				return
			}
		}
	} else {
		_, _ = fmt.Fprint(listFile, "No RBAC Role resources found in namespace "+namespace)
	}
}

func dumpRbacRoleBindings(namespace, dir string, ctx TestContext) {
	const message = "Could not dump RBAC RoleBindings for namespace %s due to %s"

	list, err := ctx.KubeClient.RbacV1().RoleBindings(namespace).List(ctx.Context, metav1.ListOptions{})
	if err != nil {
		ctx.Logf(message, namespace, err.Error())
		return
	}

	logsDir, err := EnsureLogsDir(dir)
	if err != nil {
		ctx.Logf(message, namespace, err.Error())
		return
	}

	listFile, err := os.Create(logsDir + string(os.PathSeparator) + "rbac-role-binding-list.txt")
	if err != nil {
		ctx.Logf(message, namespace, err.Error())
		return
	}

	fn := func() { _ = listFile.Close() }
	defer fn()

	if len(list.Items) > 0 {
		for _, item := range list.Items {
			_, err = fmt.Fprint(listFile, item.GetName()+"\n")
			if err != nil {
				ctx.Logf(message, namespace, err.Error())
				return
			}

			d, err := json.MarshalIndent(item, "", "    ")
			if err != nil {
				ctx.Logf(message, namespace, err.Error())
				return
			}

			file, err := os.Create(logsDir + string(os.PathSeparator) + "RBAC-RoleBinding-" + item.GetName() + ".json")
			if err != nil {
				ctx.Logf(message, namespace, err.Error())
				return
			}

			_, err = fmt.Fprint(file, string(d))
			_ = file.Close()
			if err != nil {
				ctx.Logf(message, namespace, err.Error())
				return
			}
		}
	} else {
		_, _ = fmt.Fprint(listFile, "No RBAC RoleBinding resources found in namespace "+namespace)
	}
}

func dumpServiceAccounts(namespace, dir string, ctx TestContext) {
	const message = "Could not dump ServiceAccounts for namespace %s due to %s"

	list, err := ctx.KubeClient.CoreV1().ServiceAccounts(namespace).List(ctx.Context, metav1.ListOptions{})
	if err != nil {
		ctx.Logf(message, namespace, err.Error())
		return
	}

	logsDir, err := EnsureLogsDir(dir)
	if err != nil {
		ctx.Logf(message, namespace, err.Error())
		return
	}

	listFile, err := os.Create(logsDir + string(os.PathSeparator) + "service-accounts-list.txt")
	if err != nil {
		ctx.Logf(message, namespace, err.Error())
		return
	}

	fn := func() { _ = listFile.Close() }
	defer fn()

	if len(list.Items) > 0 {
		for _, item := range list.Items {
			_, err = fmt.Fprint(listFile, item.GetName()+"\n")
			if err != nil {
				ctx.Logf(message, namespace, err.Error())
				return
			}

			d, err := json.MarshalIndent(item, "", "    ")
			if err != nil {
				ctx.Logf(message, namespace, err.Error())
				return
			}

			file, err := os.Create(logsDir + string(os.PathSeparator) + "ServiceAccount-" + item.GetName() + ".json")
			if err != nil {
				ctx.Logf(message, namespace, err.Error())
				return
			}

			_, err = fmt.Fprint(file, string(d))
			_ = file.Close()
			if err != nil {
				ctx.Logf(message, namespace, err.Error())
				return
			}
		}
	} else {
		_, _ = fmt.Fprint(listFile, "No ServiceAccount resources found in namespace "+namespace)
	}
}

func DumpPodsForTest(ctx TestContext, t *testing.T) {
	namespace := GetTestNamespace()
	dumpPods(namespace, t.Name(), ctx)
}

func dumpPods(namespace, dir string, ctx TestContext) {
	const message = "Could not dump Pods for namespace %s due to %s"

	list, err := ctx.KubeClient.CoreV1().Pods(namespace).List(ctx.Context, metav1.ListOptions{})
	if err != nil {
		ctx.Logf(message, namespace, err.Error())
		return
	}

	logsDir, err := EnsureLogsDir(dir)
	if err != nil {
		ctx.Logf(message, namespace, err.Error())
		return
	}

	listFile, err := os.Create(logsDir + string(os.PathSeparator) + "pod-list.txt")
	if err != nil {
		ctx.Logf(message, namespace, err.Error())
		return
	}

	fn := func() { _ = listFile.Close() }
	defer fn()

	if len(list.Items) > 0 {
		for i := range list.Items {
			item := list.Items[i]
			_, err = fmt.Fprint(listFile, item.GetName()+"\n")
			if err != nil {
				ctx.Logf(message, namespace, err.Error())
				return
			}

			d, err := json.MarshalIndent(item, "", "    ")
			if err != nil {
				ctx.Logf(message, namespace, err.Error())
				return
			}

			file, err := os.Create(logsDir + string(os.PathSeparator) + "Pod-" + item.GetName() + ".json")
			if err != nil {
				ctx.Logf(message, namespace, err.Error())
				return
			}

			_, err = fmt.Fprint(file, string(d))
			_ = file.Close()
			if err != nil {
				ctx.Logf(message, namespace, err.Error())
				return
			}

			DumpPodLog(ctx, &item, dir)
		}
	} else {
		_, _ = fmt.Fprint(listFile, "No Pod resources found in namespace "+namespace)
	}
}

func EnsureLogsDir(subDir string) (string, error) {
	logs, err := FindTestLogsDir()
	if err != nil {
		return "", err
	}

	dir := logs + string(os.PathSeparator) + subDir
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return "", err
	}

	return dir, err
}

// GetLastPodReadyTime returns the latest ready time from all the specified Pods for a given deployment
func GetLastPodReadyTime(pods []corev1.Pod, deployment string) metav1.Time {
	t := metav1.NewTime(time.Time{})
	for _, p := range pods {
		if p.GetLabels()[coh.LabelCoherenceDeployment] == deployment {
			for _, c := range p.Status.Conditions {
				if c.Type == corev1.PodReady && t.Before(&c.LastTransitionTime) {
					t = c.LastTransitionTime
				}
			}
		}
	}
	return t
}

// GetFirstPodReadyTime returns the first ready time from all the specified Pods for a given deployment
func GetFirstPodReadyTime(pods []corev1.Pod, deployment string) metav1.Time {
	t := metav1.NewTime(time.Now())
	for _, p := range pods {
		if p.GetLabels()[coh.LabelCoherenceDeployment] == deployment {
			for _, c := range p.Status.Conditions {
				if c.Type == corev1.PodReady && t.After(c.LastTransitionTime.Time) {
					t = c.LastTransitionTime
				}
			}
		}
	}
	return t
}

// GetFirstPodScheduledTime returns the earliest scheduled time from all the specified Pods for a given deployment
func GetFirstPodScheduledTime(pods []corev1.Pod, deployment string) metav1.Time {
	t := metav1.NewTime(time.Now())
	for _, p := range pods {
		if p.GetLabels()[coh.LabelCoherenceDeployment] == deployment {
			for _, c := range p.Status.Conditions {
				if c.Type == corev1.PodScheduled && t.After(c.LastTransitionTime.Time) {
					t = c.LastTransitionTime
				}
			}
		}
	}
	return t
}

// AssertSingleDeployment tests that a cluster can be created using the specified yaml.
func AssertSingleDeployment(ctx TestContext, t *testing.T, yamlFile string) (coh.Coherence, error) {
	c, _ := AssertDeployments(ctx, t, yamlFile)
	for _, v := range c {
		return v, nil
	}
	// should not actually get here
	return coh.Coherence{}, fmt.Errorf("there were no Coherence resources found")
}

// AssertDeployments tests that one or more clusters can be created using the specified yaml.
func AssertDeployments(ctx TestContext, t *testing.T, yamlFile string) (map[string]coh.Coherence, []corev1.Pod) {
	return AssertDeploymentsInNamespace(ctx, t, yamlFile, GetTestNamespace())
}

// AssertDeploymentsInNamespace tests that a cluster can be created using the specified yaml.
func AssertDeploymentsInNamespace(ctx TestContext, t *testing.T, yamlFile, namespace string) (map[string]coh.Coherence, []corev1.Pod) {
	// initialise Gomega so we can use matchers
	g := NewGomegaWithT(t)

	deployments, err := NewCoherenceFromYaml(namespace, yamlFile)
	g.Expect(err).NotTo(HaveOccurred())

	// we must have at least one deployment
	g.Expect(len(deployments)).NotTo(BeZero())

	// assert all deployments have the same cluster name
	clusterName := deployments[0].GetCoherenceClusterName()
	for _, d := range deployments {
		g.Expect(d.GetCoherenceClusterName()).To(Equal(clusterName))
	}

	// work out the expected cluster size
	expectedClusterSize := 0
	expectedWkaSize := 0
	for _, d := range deployments {
		ctx.Logf("Deployment %s has replica count %d", d.Name, d.GetReplicas())
		replicas := int(d.GetReplicas())
		expectedClusterSize += replicas
		if d.Spec.Coherence.IsWKAMember() {
			expectedWkaSize += replicas
		}
	}
	ctx.Logf("Expected cluster size is %d", expectedClusterSize)

	for i := range deployments {
		d := deployments[i]
		ctx.Logf("Deploying %s", d.Name)
		// deploy the Coherence resource
		err = ctx.Client.Create(ctx.Context, &d)
		g.Expect(err).NotTo(HaveOccurred())
	}

	// Assert that a StatefulSet or Job of the correct number or replicas is created for each roleSpec in the cluster
	for _, d := range deployments {
		ctx.Logf("Waiting for StatefulSet for deployment %s", d.Name)
		// Wait for the StatefulSet for the roleSpec to be ready - wait five minutes max
		sts, err := WaitForStatefulSet(ctx, namespace, d.Name, d.GetReplicas(), time.Second*10, time.Minute*5)
		g.Expect(err).NotTo(HaveOccurred())
		g.Expect(sts.Status.ReadyReplicas).To(Equal(d.GetReplicas()))
		ctx.Logf("Have StatefulSet for deployment %s", d.Name)
	}

	// Assert that the finalizer has been added to all the deployments that do not have AllowUnsafeDelete=false
	for _, d := range deployments {
		ctx.Logf("Deploying %s", d.Name)
		// deploy the Coherence resource
		actual := &coh.Coherence{}
		err = ctx.Client.Get(ctx.Context, d.GetNamespacedName(), actual)
		g.Expect(err).NotTo(HaveOccurred())
		if d.Spec.AllowUnsafeDelete != nil && *d.Spec.AllowUnsafeDelete {
			g.Expect(controllerutil.ContainsFinalizer(actual, coh.CoherenceFinalizer)).To(BeFalse())
		} else {
			g.Expect(controllerutil.ContainsFinalizer(actual, coh.CoherenceFinalizer)).To(BeTrue())
		}
	}

	// Get all the Pods in the cluster
	ctx.Logf("Getting all Pods for cluster '%s'", clusterName)
	pods, err := ListCoherencePodsForCluster(ctx, namespace, clusterName)
	g.Expect(err).NotTo(HaveOccurred())
	ctx.Logf("Found %d Pods for cluster '%s'", len(pods), clusterName)

	// assert that the correct number of Pods is returned
	g.Expect(len(pods)).To(Equal(expectedClusterSize))

	// Verify that the WKA service has the same number of endpoints as the cluster size.
	serviceName := deployments[0].GetWkaServiceName()

	ep, err := ctx.KubeClient.CoreV1().Endpoints(namespace).Get(ctx.Context, serviceName, metav1.GetOptions{})
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(len(ep.Subsets)).NotTo(BeZero())

	subset := ep.Subsets[0]
	g.Expect(len(subset.Addresses)).To(Equal(expectedWkaSize))

	m := make(map[string]coh.Coherence)
	for _, d := range deployments {
		opts := client.ObjectKey{Namespace: namespace, Name: d.Name}
		dpl := coh.Coherence{}
		err = ctx.Client.Get(ctx.Context, opts, &dpl)
		g.Expect(err).NotTo(HaveOccurred())
		m[dpl.Name] = dpl
	}

	// Obtain the expected WKA list of Pod IP addresses
	var wkaPods []string
	for _, d := range deployments {
		if d.Spec.Coherence.IsWKAMember() {
			pods, err := ListCoherencePodsForDeployment(ctx, d.Namespace, d.Name)
			g.Expect(err).NotTo(HaveOccurred())
			for _, pod := range pods {
				wkaPods = append(wkaPods, pod.Status.PodIP)
			}
		}
	}

	// Verify that the WKA service endpoints list for each deployment has all the required the Pod IP addresses.
	for _, d := range deployments {
		serviceName := d.GetWkaServiceName()
		ep, err = ctx.KubeClient.CoreV1().Endpoints(namespace).Get(ctx.Context, serviceName, metav1.GetOptions{})
		g.Expect(err).NotTo(HaveOccurred())
		g.Expect(len(ep.Subsets)).NotTo(BeZero())

		subset := ep.Subsets[0]
		g.Expect(len(subset.Addresses)).To(Equal(len(wkaPods)))
		var actualWKA []string
		for _, address := range subset.Addresses {
			actualWKA = append(actualWKA, address.IP)
		}
		g.Expect(actualWKA).To(ConsistOf(wkaPods))
	}

	return m, pods
}

// AssertCoherenceJobs tests that one or more CoherenceJobs can be created using the specified yaml.
func AssertCoherenceJobs(ctx TestContext, t *testing.T, yamlFile string) (map[string]coh.CoherenceJob, []corev1.Pod) {
	return AssertCoherenceJobsInNamespace(ctx, t, yamlFile, GetTestNamespace())
}

// AssertCoherenceJobsInNamespace tests that a CoherenceJob can be created using the specified yaml.
func AssertCoherenceJobsInNamespace(ctx TestContext, t *testing.T, yamlFile, namespace string) (map[string]coh.CoherenceJob, []corev1.Pod) {
	// initialise Gomega so we can use matchers
	g := NewGomegaWithT(t)

	jobs, err := NewCoherenceJobFromYaml(namespace, yamlFile)
	g.Expect(err).NotTo(HaveOccurred())

	return AssertCoherenceJobsSpecInNamespace(ctx, t, jobs, namespace)
}

// AssertCoherenceJobsSpec tests that a CoherenceJob can be created.
func AssertCoherenceJobsSpec(ctx TestContext, t *testing.T, jobs []coh.CoherenceJob) (map[string]coh.CoherenceJob, []corev1.Pod) {
	return AssertCoherenceJobsSpecInNamespace(ctx, t, jobs, GetTestNamespace())
}

// AssertCoherenceJobsSpecInNamespace tests that a CoherenceJob can be created.
func AssertCoherenceJobsSpecInNamespace(ctx TestContext, t *testing.T, jobs []coh.CoherenceJob, namespace string) (map[string]coh.CoherenceJob, []corev1.Pod) {
	// initialise Gomega so we can use matchers
	g := NewGomegaWithT(t)

	var err error

	// we must have at least one deployment
	g.Expect(len(jobs)).NotTo(BeZero())

	// assert all jobs have the same cluster name
	clusterName := jobs[0].GetCoherenceClusterName()
	for _, d := range jobs {
		g.Expect(d.GetCoherenceClusterName()).To(Equal(clusterName))
	}

	// work out the expected cluster size
	expectedClusterSize := 0
	expectedWkaSize := 0
	for _, d := range jobs {
		ctx.Logf("CoherenceJob %s has replica count %d", d.Name, d.GetReplicas())
		replicas := int(d.GetReplicas())
		expectedClusterSize += replicas
		if d.Spec.Coherence.IsWKAMember() {
			expectedWkaSize += replicas
		}
	}
	ctx.Logf("Expected cluster size is %d", expectedClusterSize)

	for i := range jobs {
		d := jobs[i]
		ctx.Logf("Deploying CoherenceJob %s", d.Name)
		// deploy the CoherenceJob resource
		err = ctx.Client.Create(ctx.Context, &d)
		g.Expect(err).NotTo(HaveOccurred())
	}

	// Assert that a Job with the correct number of replicas is created for each Spec in the cluster
	for _, d := range jobs {
		ctx.Logf("Waiting for Job for deployment %s", d.Name)
		// Wait for the Job for the roleSpec to be ready - wait five minutes max
		sts, err := WaitForJob(ctx, namespace, d.Name, d.GetReplicas(), time.Second*10, time.Minute*5)
		g.Expect(err).NotTo(HaveOccurred())
		g.Expect(sts.Status).NotTo(BeNil())
		g.Expect(sts.Status.Ready).NotTo(BeNil())
		g.Expect(*sts.Status.Ready).To(Equal(d.GetReplicas()))
		ctx.Logf("Have Job for deployment %s", d.Name)
		//}
	}

	// Get all the Pods in the cluster
	ctx.Logf("Getting all Pods for Job '%s'", clusterName)
	pods, err := ListCoherencePodsForCluster(ctx, namespace, clusterName)
	g.Expect(err).NotTo(HaveOccurred())
	ctx.Logf("Found %d Pods for Job '%s'", len(pods), clusterName)

	// assert that the correct number of Pods is returned
	g.Expect(len(pods)).To(Equal(expectedClusterSize))

	// Verify that the WKA service has the same number of endpoints as the cluster size.
	serviceName := jobs[0].GetWkaServiceName()

	ep, err := ctx.KubeClient.CoreV1().Endpoints(namespace).Get(ctx.Context, serviceName, metav1.GetOptions{})
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(len(ep.Subsets)).NotTo(BeZero())

	subset := ep.Subsets[0]
	g.Expect(len(subset.Addresses)).To(Equal(expectedWkaSize))

	m := make(map[string]coh.CoherenceJob)
	for _, d := range jobs {
		opts := client.ObjectKey{Namespace: namespace, Name: d.Name}
		dpl := coh.CoherenceJob{}
		err = ctx.Client.Get(ctx.Context, opts, &dpl)
		g.Expect(err).NotTo(HaveOccurred())
		m[dpl.Name] = dpl
	}

	// Obtain the expected WKA list of Pod IP addresses
	var wkaPods []string
	for _, d := range jobs {
		if d.Spec.Coherence.IsWKAMember() {
			pods, err := ListCoherencePodsForDeployment(ctx, d.Namespace, d.Name)
			g.Expect(err).NotTo(HaveOccurred())
			for _, pod := range pods {
				wkaPods = append(wkaPods, pod.Status.PodIP)
			}
		}
	}

	// Verify that the WKA service endpoints list for each deployment has all the required the Pod IP addresses.
	for _, d := range jobs {
		serviceName := d.GetWkaServiceName()
		ep, err = ctx.KubeClient.CoreV1().Endpoints(namespace).Get(ctx.Context, serviceName, metav1.GetOptions{})
		g.Expect(err).NotTo(HaveOccurred())
		g.Expect(len(ep.Subsets)).NotTo(BeZero())

		subset := ep.Subsets[0]
		g.Expect(len(subset.Addresses)).To(Equal(len(wkaPods)))
		var actualWKA []string
		for _, address := range subset.Addresses {
			actualWKA = append(actualWKA, address.IP)
		}
		g.Expect(actualWKA).To(ConsistOf(wkaPods))
	}

	return m, pods
}

// WaitForDelete waits for the provided runtime object to be deleted from cluster
func WaitForDelete(ctx TestContext, obj client.Object) error {
	key := ObjectKey(obj)
	ctx.Logf("Waiting for obj %s/%s to be finally deleted", key.Namespace, key.Name)

	// Wait for resources to be deleted.
	return wait.PollUntilContextTimeout(ctx.Context, 1*time.Second, 30*time.Second, true, func(context.Context) (done bool, err error) {
		err = ctx.Client.Get(ctx.Context, key, obj.DeepCopyObject().(client.Object))
		ctx.Logf("Fetched %s/%s to wait for delete: %v", key.Namespace, key.Name, err)

		if err != nil && apierrors.IsNotFound(err) {
			return true, nil
		}
		return false, err
	})
}

// ObjectKey returns an instantiated ObjectKey for the provided object.
func ObjectKey(obj runtime.Object) client.ObjectKey {
	m, _ := meta.Accessor(obj)
	return client.ObjectKey{
		Name:      m.GetName(),
		Namespace: m.GetNamespace(),
	}
}

func DeletePersistentVolumes(ctx TestContext, namespace string) {
	opts := metav1.ListOptions{}

	claims, err := ctx.KubeClient.CoreV1().PersistentVolumeClaims(namespace).List(ctx.Context, opts)
	if err != nil {
		ctx.Logf("Failed to list PVCs in namespace %s %v", namespace, err)
		return
	}

	var pvs []string
	delOpts := metav1.DeleteOptions{}

	for _, claim := range claims.Items {
		ctx.Logf("Deleting PVC %s/%s", namespace, claim.Name)
		err := ctx.KubeClient.CoreV1().PersistentVolumeClaims(namespace).Delete(ctx.Context, claim.Name, delOpts)
		if err != nil {
			ctx.Logf("Failed to delete PVC %s/%s %v", namespace, claim.Name, err)
		}
		pvs = append(pvs, claim.Spec.VolumeName)
	}

	for _, pv := range pvs {
		ctx.Logf("Deleting PV %s/%s", namespace, pv)
		err := ctx.KubeClient.CoreV1().PersistentVolumes().Delete(ctx.Context, pv, delOpts)
		if err != nil {
			ctx.Logf("Failed to delete PV %s %v", pv, err)
		}
	}
}

func AddLoopbackTestHostnameLabel(c *coh.Coherence) {
	if c.Spec.Labels == nil {
		c.Spec.Labels = map[string]string{}
	}
	c.Spec.Labels[operator.LabelTestHostName] = "127.0.0.1"
}
