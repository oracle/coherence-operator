/*
 * Copyright (c) 2019, 2020 Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

// helper package contains various helpers for use in end-to-end testing.
package helper

import (
	goctx "context"
	"encoding/json"
	"fmt"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	"github.com/oracle/coherence-operator/pkg/apis"
	coh "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	"golang.org/x/net/context"
	"io"
	"io/ioutil"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"os"
	"strings"
	"testing"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"time"
)

const operatorPodSelector = "name=coherence-operator"

var (
	RetryInterval        = time.Second * 5
	Timeout              = time.Minute * 3
	CleanupRetryInterval = time.Second * 1
	CleanupTimeout       = time.Second * 5

	tenSeconds int32 = 10

	Readiness = &coh.ReadinessProbeSpec{
		InitialDelaySeconds: &tenSeconds,
		PeriodSeconds:       &tenSeconds,
	}
)

type Logger interface {
	Log(args ...interface{})
	Logf(format string, args ...interface{})
}

func CreateTestContext(t *testing.T) *framework.TestCtx {
	f := framework.Global
	testCtx := framework.NewTestCtx(t)

	namespace, err := testCtx.GetNamespace()
	if err != nil {
		t.Fatal(err)
		return nil
	}

	testCtx.AddCleanupFn(func() error {
		return WaitForCoherenceInternalCleanup(f, namespace)
	})

	cleanup := framework.CleanupOptions{TestContext: testCtx, Timeout: CleanupTimeout, RetryInterval: CleanupRetryInterval}

	err = testCtx.InitializeClusterResources(&cleanup)
	if err != nil && !errors.IsAlreadyExists(err) {
		t.Fatal(err)
		return nil
	}

	clusterList := &coh.CoherenceClusterList{}

	err = framework.AddToFrameworkScheme(apis.AddToScheme, clusterList)
	if err != nil {
		t.Fatal(err)
		return nil
	}

	return testCtx
}

func NewUnstructuredCoherenceInternalList() unstructured.UnstructuredList {
	u := unstructured.UnstructuredList{}
	u.SetGroupVersionKind(schema.GroupVersionKind{Group: "coherence.oracle.com", Version: "v1", Kind: "CoherenceInternal"})

	return u
}

func NewUnstructuredCoherenceInternal() unstructured.Unstructured {
	u := unstructured.Unstructured{}
	u.SetGroupVersionKind(schema.GroupVersionKind{Group: "coherence.oracle.com", Version: "v1", Kind: "CoherenceInternal"})

	return u
}

func DefaultCleanup(ctx *framework.TestCtx) *framework.CleanupOptions {
	return &framework.CleanupOptions{TestContext: ctx, Timeout: Timeout, RetryInterval: RetryInterval}
}

// WaitForStatefulSetForRole waits for a StatefulSet to be created for the specified role.
func WaitForStatefulSetForRole(kubeclient kubernetes.Interface, namespace string, cluster *coh.CoherenceCluster, role coh.CoherenceRoleSpec, retryInterval, timeout time.Duration, logger Logger) (*appsv1.StatefulSet, error) {
	return WaitForStatefulSet(kubeclient, namespace, role.GetFullRoleName(cluster), role.GetReplicas(), retryInterval, timeout, logger)
}

// WaitForStatefulSet waits for a StatefulSet to be created with the specified number of replicas.
func WaitForStatefulSet(kubeclient kubernetes.Interface, namespace, stsName string, replicas int32, retryInterval, timeout time.Duration, logger Logger) (*appsv1.StatefulSet, error) {
	var sts *appsv1.StatefulSet

	err := wait.Poll(retryInterval, timeout, func() (done bool, err error) {
		sts, err = kubeclient.AppsV1().StatefulSets(namespace).Get(stsName, metav1.GetOptions{})
		if err != nil {
			if apierrors.IsNotFound(err) {
				logger.Logf("Waiting for availability of %s StatefulSet - NotFound\n", stsName)
				return false, nil
			}
			logger.Logf("Waiting for availability of %s StatefulSet - %s\n", stsName, err.Error())
			return false, err
		}

		if sts.Status.ReadyReplicas == replicas {
			return true, nil
		}
		logger.Logf("Waiting for full availability of %s StatefulSet (%d/%d)\n", stsName, sts.Status.ReadyReplicas, replicas)
		return false, nil
	})

	if err != nil && sts != nil {
		d, _ := json.Marshal(sts)
		logger.Logf("Error waiting for StatefulSet\n%s", string(d))
	}
	return sts, err
}

// A function that takes a role and determines whether it meets a condition
type RoleStateCondition interface {
	Test(*coh.CoherenceRole) bool
	String() string
}

// An always true RoleStateCondition
type alwayRoleCondition struct{}

func (a alwayRoleCondition) Test(*coh.CoherenceRole) bool {
	return true
}

func (a alwayRoleCondition) String() string {
	return "true"
}

func AlwayRoleCondition() RoleStateCondition {
	return alwayRoleCondition{}
}

// An always true RoleStateCondition
type replicasRoleCondition struct {
	replicas int32
}

func (a replicasRoleCondition) Test(*coh.CoherenceRole) bool {
	return true
}

func (a replicasRoleCondition) String() string {
	return fmt.Sprintf("replicas == %d", a.replicas)
}

func ReplicasRoleCondition(replicas int32) RoleStateCondition {
	return replicasRoleCondition{replicas: replicas}
}

// WaitForCoherenceRole waits for a CoherenceRole to be created.
func WaitForCoherenceRole(f *framework.Framework, namespace, name string, retryInterval, timeout time.Duration, logger Logger) (*coh.CoherenceRole, error) {
	return WaitForCoherenceRoleCondition(f, namespace, name, alwayRoleCondition{}, retryInterval, timeout, logger)
}

// WaitForCoherenceRole waits for a CoherenceRole to be created.
func WaitForCoherenceRoleCondition(f *framework.Framework, namespace, name string, conditon RoleStateCondition, retryInterval, timeout time.Duration, logger Logger) (*coh.CoherenceRole, error) {
	var role *coh.CoherenceRole

	err := wait.Poll(retryInterval, timeout, func() (done bool, err error) {
		role, err = GetCoherenceRole(f, namespace, name)
		if err != nil {
			if apierrors.IsNotFound(err) {
				logger.Logf("Waiting for availability of CoherenceRole %s - NotFound\n", name)
				return false, nil
			}
			logger.Logf("Waiting for availability of CoherenceRole %s - %s\n", name, err.Error())
			return false, err
		}
		valid := true
		if conditon != nil {
			valid = conditon.Test(role)
			if !valid {
				logger.Logf("Waiting for CoherenceRole %s to meet condition '%s'\n", name, conditon.String())
			}
		}
		return valid, nil
	})

	return role, err
}

// GetCoherenceRole gets the specified CoherenceRole
func GetCoherenceRole(f *framework.Framework, namespace, name string) (*coh.CoherenceRole, error) {
	opts := client.ObjectKey{Namespace: namespace, Name: name}
	role := &coh.CoherenceRole{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
		},
	}

	err := f.Client.Get(goctx.TODO(), opts, role)

	return role, err
}

// WaitForOperatorPods waits for a Coherence Operator Pods to be created.
func WaitForOperatorPods(k8s kubernetes.Interface, namespace string, retryInterval, timeout time.Duration) ([]corev1.Pod, error) {
	return WaitForPodsWithSelector(k8s, namespace, operatorPodSelector, retryInterval, timeout)
}

// WaitForOperatorPods waits for a Coherence Operator Pods to be created.
func WaitForPodsWithSelector(k8s kubernetes.Interface, namespace, selector string, retryInterval, timeout time.Duration) ([]corev1.Pod, error) {
	var pods []corev1.Pod

	err := wait.Poll(retryInterval, timeout, func() (done bool, err error) {
		pods, err = ListPodsWithLabelSelector(k8s, namespace, selector)
		if err != nil {
			return false, err
		}
		return len(pods) > 0, nil
	})
	return pods, err
}

// WaitForOperatorDeletion waits for deletion of the Operator Pods.
func WaitForOperatorDeletion(k8s kubernetes.Interface, namespace string, retryInterval, timeout time.Duration, logger Logger) error {
	return WaitForDeleteOfPodsWithSelector(k8s, namespace, operatorPodSelector, retryInterval, timeout, logger)
}

// WaitForDeleteOfPodsWithSelector waits for Pods to be deleted.
func WaitForDeleteOfPodsWithSelector(k8s kubernetes.Interface, namespace, selector string, retryInterval, timeout time.Duration, logger Logger) error {
	logger.Logf("Waiting for Pods in namespace %s with selector '%s' to be deleted", namespace, selector)
	var pods []corev1.Pod

	err := wait.Poll(retryInterval, timeout, func() (done bool, err error) {
		logger.Logf("List Pods in namespace %s with selector '%s'", namespace, selector)
		pods, err = ListPodsWithLabelSelector(k8s, namespace, selector)
		if err != nil {
			logger.Logf("Error listing Pods in namespace %s with selector '%s' - %s", namespace, selector, err.Error())
			return false, err
		}
		logger.Logf("Found %d Pods in namespace %s with selector '%s'", len(pods), namespace, selector)
		return len(pods) == 0, nil
	})

	logger.Logf("Finished waiting for Pods in namespace %s with selector '%s' to be deleted. Error=%s", namespace, selector, err)
	return err
}

// WaitForDeletion waits for deletion of the specified resource.
func WaitForDeletion(f *framework.Framework, namespace, name string, resource runtime.Object, retryInterval, timeout time.Duration) error {
	key := types.NamespacedName{Namespace: namespace, Name: name}

	err := wait.Poll(retryInterval, timeout, func() (done bool, err error) {
		err = f.Client.Get(context.TODO(), key, resource)
		switch {
		case err != nil && errors.IsNotFound(err):
			return true, nil
		case err != nil && !errors.IsNotFound(err):
			return false, err
		default:
			fmt.Printf("Waiting for deletion of %s in namespace %s\n", name, namespace)
			return false, nil
		}
	})

	return err
}

// List the Operator Pods that exist - this is Pods with the label "name=coh-operator"
func ListOperatorPods(client kubernetes.Interface, namespace string) ([]corev1.Pod, error) {
	return ListPodsWithLabelSelector(client, namespace, operatorPodSelector)
}

// List the Coherence Cluster Pods that exist for a cluster - this is Pods with the label "coherenceCluster=<cluster>"
func ListCoherencePodsForCluster(client kubernetes.Interface, namespace, cluster string) ([]corev1.Pod, error) {
	return ListPodsWithLabelSelector(client, namespace, fmt.Sprintf("coherenceCluster=%s", cluster))
}

// List the Coherence Cluster Pods that exist for a role - this is Pods with the label "coherenceCluster=<cluster>,coherenceRole=<role>"
func ListCoherencePodsForRole(client kubernetes.Interface, namespace, cluster, role string) ([]corev1.Pod, error) {
	return ListPodsWithLabelSelector(client, namespace, fmt.Sprintf("coherenceCluster=%s,coherenceRole=%s", cluster, role))
}

// WaitForPodsWithLabel waits for at least the required number of Pods matching the specified labels selector to be created.
func WaitForPodsWithLabel(k8s kubernetes.Interface, namespace, selector string, count int, retryInterval, timeout time.Duration) ([]corev1.Pod, error) {
	var pods []corev1.Pod

	err := wait.Poll(retryInterval, timeout, func() (done bool, err error) {
		pods, err = ListPodsWithLabelSelector(k8s, namespace, selector)
		if err != nil {
			fmt.Printf("Waiting for at least %d Pods with label selector '%s' - failed due to %s\n", count, selector, err.Error())
			return false, err
		}
		found := len(pods) >= count
		if !found {
			fmt.Printf("Waiting for at least %d Pods with label selector '%s' - found %d\n", count, selector, len(pods))
		}
		return found, nil
	})

	return pods, err
}

// List the Coherence Cluster Pods that exist for a given label selector.
func ListPodsWithLabelSelector(k8s kubernetes.Interface, namespace, selector string) ([]corev1.Pod, error) {
	opts := metav1.ListOptions{LabelSelector: selector}

	list, err := k8s.CoreV1().Pods(namespace).List(opts)
	if err != nil {
		return []corev1.Pod{}, err
	}

	return list.Items, nil
}

// WaitForPodReady waits for a Pods to be ready.
func WaitForPodReady(k8s kubernetes.Interface, namespace, name string, retryInterval, timeout time.Duration) error {
	err := wait.Poll(retryInterval, timeout, func() (done bool, err error) {
		p, err := k8s.CoreV1().Pods(namespace).Get(name, metav1.GetOptions{})
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

// waitForCleanup waits until there are no CoherenceInternal resources left in the test namespace.
// The default clean-up hooks only wait for deletion of resources directly created via the test client
// but CoherenceInternal resources and corresponding Helm installs are created internally.
func WaitForCoherenceInternalCleanup(f *framework.Framework, namespace string) error {
	fmt.Printf("Waiting for clean-up of CoherenceInternal resources in namespace %s\n", namespace)

	// wait for all CoherenceInternal resources to be deleted
	clusters := &coh.CoherenceClusterList{}
	err := f.Client.List(goctx.TODO(), clusters, client.InNamespace(namespace))
	if err != nil {
		return err
	}

	// Delete all of the CoherenceClusters
	for _, c := range clusters.Items {
		fmt.Printf("Deleting CoherenceCluster %s in namespace %s\n", c.Name, c.Namespace)
		err = f.Client.Delete(goctx.TODO(), &c)
		if err != nil {
			fmt.Printf("Error deleting CoherenceCluster %s - %s\n", c.Name, err.Error())
		}
	}

	// Wait for removal of the CoherenceClusters
	err = wait.Poll(RetryInterval, Timeout, func() (done bool, err error) {
		err = f.Client.List(goctx.TODO(), clusters, client.InNamespace(namespace))
		if err == nil {
			if len(clusters.Items) > 0 {
				fmt.Printf("Waiting for deletion of %d CoherenceCluster resources\n", len(clusters.Items))
				return false, nil
			}
			return true, nil
		} else {
			fmt.Printf("Error waiting for deletion of CoherenceCluster resources: %s\n", err.Error())
			return false, nil
		}
	})

	roles := &coh.CoherenceRoleList{}
	err = f.Client.List(goctx.TODO(), roles, client.InNamespace(namespace))
	if err != nil {
		return err
	}

	// Delete all of the CoherenceRoles
	for _, r := range roles.Items {
		fmt.Printf("Deleting CoherenceRoles %s in namespace %s\n", r.Name, r.Namespace)
		err = f.Client.Delete(goctx.TODO(), &r)
		if err != nil {
			fmt.Printf("Error deleting CoherenceRoles %s - %s\n", r.Name, err.Error())
		}
	}

	// Wait for removal of the CoherenceRoles
	err = wait.Poll(RetryInterval, Timeout, func() (done bool, err error) {
		err = f.Client.List(goctx.TODO(), roles, client.InNamespace(namespace))
		if err == nil || isNoResources(err) || errors.IsNotFound(err) {
			if len(roles.Items) > 0 {
				fmt.Printf("Waiting for deletion of %d CoherenceRole resources\n", len(roles.Items))
				return false, nil
			}
			return true, nil
		} else {
			fmt.Printf("Error waiting for deletion of CoherenceRole resources: %s\n%v\n", err.Error(), err)
			return false, nil
		}
	})

	uList := NewUnstructuredCoherenceInternalList()
	err = f.Client.List(goctx.TODO(), &uList, client.InNamespace(namespace))
	if err != nil {
		return err
	}

	// Delete all of the CoherenceInternals
	for _, ci := range uList.Items {
		fmt.Printf("Deleting CoherenceInternals %s in namespace %s\n", ci.GetName(), ci.GetNamespace())
		err = f.Client.Delete(goctx.TODO(), &ci)
		if err != nil {
			fmt.Printf("Error deleting CoherenceInternals %s - %s\n", ci.GetName(), err.Error())
		}
	}

	// Wait for removal of the CoherenceInternals
	err = wait.Poll(time.Second*5, time.Minute*1, func() (done bool, err error) {
		err = f.Client.List(goctx.TODO(), &uList, client.InNamespace(namespace))
		if err == nil || isNoResources(err) || errors.IsNotFound(err) {
			if len(uList.Items) > 0 {
				fmt.Printf("Waiting for deletion of %d CoherenceInternal resources\n", len(uList.Items))
				return false, nil
			}
			return true, nil
		} else {
			fmt.Printf("Error waiting for deletion of CoherenceInternal resources: %s\n", err.Error())
			return false, nil
		}
	})

	// List and print the remaining CoherenceInternals
	err = f.Client.Client.List(goctx.TODO(), &uList, client.InNamespace(namespace))
	if err != nil && !isNoResources(err) && !errors.IsNotFound(err) {
		fmt.Printf("Error listing CoherenceInternal resources - %s\n", err.Error())
	} else {
		if len(uList.Items) > 0 {
			fmt.Printf("Remaining CoherenceInternal resources in namespace %s (%d):\n", namespace, len(uList.Items))
			for i, ci := range uList.Items {
				fmt.Printf("%d: %s\n", i, ci.GetName())
			}
		} else {
			fmt.Printf("Zero CoherenceInternal resources remain in namespace %s\n", namespace)
		}
	}

	var empty []string

	// Force delete the remaining CoherenceInternals
	err = wait.Poll(RetryInterval, Timeout, func() (done bool, err error) {
		err = f.Client.List(goctx.TODO(), &uList, client.InNamespace(namespace))
		if err == nil {
			if len(uList.Items) > 0 {
				fmt.Printf("Waiting for deletion of %d CoherenceInternal resources\n", len(uList.Items))

				for _, ci := range uList.Items {
					ci.SetFinalizers(empty)
					fmt.Printf("Removing finalizers and deleting CoherenceInternal resources %s\n", ci.GetName())
					_ = f.Client.Update(goctx.TODO(), &ci)
					_ = f.Client.Delete(goctx.TODO(), &ci)
				}

				return false, nil
			}
			return true, nil
		} else {
			fmt.Printf("Error waiting for deletion of CoherenceInternal resources: %s\n", err.Error())
			return false, nil
		}
	})

	return err
}

func isNoResources(err error) bool {
	return err != nil && strings.HasPrefix(err.Error(), "no matches for kind")
}

// WaitForOperatorCleanup waits until there are no Operator Pods in the test namespace.
func WaitForOperatorCleanup(kubeClient kubernetes.Interface, namespace string, logger Logger) error {
	logger.Logf("Waiting for deletion of Coherence Operator Pods\n")
	// wait for all Operator Pods to be deleted
	err := wait.Poll(RetryInterval, Timeout, func() (done bool, err error) {
		list, err := ListOperatorPods(kubeClient, namespace)
		if err == nil {
			if len(list) > 0 {
				logger.Logf("Waiting for deletion of %d Coherence Operator Pods\n", len(list))
				return false, nil
			}
			return true, nil
		} else {
			logger.Logf("Error waiting for deletion of Coherence Operator Pods: %s\n", err.Error())
			return false, nil
		}
	})

	logger.Logf("Coherence Operator Pods deleted\n")
	return err
}

// Dump the Operator Pod log to a file.
func DumpOperatorLog(kubeClient kubernetes.Interface, namespace, directory string, logger Logger) {
	list, err := kubeClient.CoreV1().Pods(namespace).List(metav1.ListOptions{LabelSelector: "name=coherence-operator"})
	if err == nil {
		if len(list.Items) > 0 {
			pod := list.Items[0]
			DumpPodLog(kubeClient, &pod, directory, logger)
		} else {
			logger.Log("Could not capture Operator Pod log. No Pods found.")
		}
	}

	if err != nil {
		logger.Logf("Could not capture Operator Pod log due to error: %s\n", err.Error())
	}
}

// Dump the Pod log to a file.
func DumpPodLog(kubeClient kubernetes.Interface, pod *corev1.Pod, directory string, logger Logger) {
	logs, err := FindTestLogsDir()
	if err != nil {
		logger.Log("cannot capture logs due to " + err.Error())
		return
	}

	pathSep := string(os.PathSeparator)

	for _, container := range pod.Spec.Containers {
		res := kubeClient.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, &corev1.PodLogOptions{Container: container.Name})
		s, err := res.Stream()
		if err == nil {
			name := logs + pathSep + directory
			err = os.MkdirAll(name, os.ModePerm)
			if err == nil {
				suffix := 0
				logName := fmt.Sprintf("%s%s%s(%s).log", name, pathSep, pod.Name, container.Name)
				_, err = os.Stat(logName)
				for err == nil {
					suffix++
					logName = fmt.Sprintf("%s%s%s(%s)-%d.log", name, pathSep, pod.Name, container.Name, suffix)
					_, err = os.Stat(logName)
				}
				out, err := os.Create(logName)
				if err == nil {
					_, err = io.Copy(out, s)
				}
			}
		}
	}
}

// Ensure that the k8s secret has been deleted
func EnsureSecretDeleted(kubeClient kubernetes.Interface, namespace, name string) error {
	err := kubeClient.CoreV1().Secrets(namespace).Delete(name, &metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return nil
}

// Get the test k8s secret that can be used for SSL testing.
func GetTestSslSecret() (*OperatorSSL, *coh.SSLSpec, error) {
	return CreateSslSecret(nil, GetTestNamespace(), GetTestSSLSecretName())
}

// Create a k8s secret that can be used for SSL testing.
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
		_, err = kubeClient.CoreV1().Secrets(namespace).Create(secret)
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

	return ioutil.ReadAll(f)
}

// Dump the operator logs and clean-up the test context
func DumpOperatorLogsAndCleanup(t *testing.T, ctx *framework.TestCtx) {
	DumpOperatorLogs(t, ctx)
	ctx.Cleanup()
}

// Dump the operator logs and clean-up the test context
func DumpOperatorLogs(t *testing.T, ctx *framework.TestCtx) {
	namespace, err := ctx.GetNamespace()
	if err == nil {
		DumpOperatorLog(framework.Global.KubeClient, namespace, t.Name(), t)
		DumpState(namespace, t.Name(), t)
	} else {
		t.Logf("Could not dump logs and state\n")
		t.Log(err)
	}
	ctx.Cleanup()
}

func DumpState(namespace, dir string, logger Logger) {
	dumpCoherenceClusters(namespace, dir, logger)
	dumpCoherenceRoles(namespace, dir, logger)
	dumpCoherenceInternals(namespace, dir, logger)
	dumpStatefulSets(namespace, dir, logger)
	dumpServices(namespace, dir, logger)
	dumpPods(namespace, dir, logger)
	dumpRoles(namespace, dir, logger)
	dumpRoleBindings(namespace, dir, logger)
	dumpServiceAccounts(namespace, dir, logger)
}

func dumpCoherenceClusters(namespace, dir string, logger Logger) {
	const message = "Could not dump CoherenceClusters for namespace %s due to %s\n"

	f := framework.Global
	list := coh.CoherenceClusterList{}
	err := f.Client.List(context.TODO(), &list, client.InNamespace(namespace))
	if err != nil {
		fmt.Printf(message, namespace, err.Error())
		return
	}

	logsDir, err := ensureLogsDir(dir)
	if err != nil {
		fmt.Printf(message, namespace, err.Error())
		return
	}

	listFile, err := os.Create(logsDir + string(os.PathSeparator) + "cluster-list.txt")
	if err != nil {
		fmt.Printf(message, namespace, err.Error())
		return
	}

	fn := func() { _ = listFile.Close() }
	defer fn()

	if len(list.Items) > 0 {
		for _, item := range list.Items {
			_, err = fmt.Fprint(listFile, item.GetName()+"\n")
			if err != nil {
				logger.Logf(message, namespace, err.Error())
				return
			}

			d, err := json.Marshal(item)
			if err != nil {
				logger.Logf(message, namespace, err.Error())
				return
			}

			file, err := os.Create(logsDir + string(os.PathSeparator) + "CoherenceCluster-" + item.GetName() + ".json")
			if err != nil {
				logger.Logf(message, namespace, err.Error())
				return
			}

			_, err = fmt.Fprint(file, string(d))
			_ = file.Close()
			if err != nil {
				logger.Logf(message, namespace, err.Error())
				return
			}
		}
	} else {
		_, _ = fmt.Fprint(listFile, "No CoherenceClusters resources found in namespace "+namespace)
	}
}

func dumpCoherenceRoles(namespace, dir string, logger Logger) {
	const message = "Could not dump CoherenceRoles for namespace %s due to %s\n"

	f := framework.Global
	list := coh.CoherenceRoleList{}
	err := f.Client.List(context.TODO(), &list, client.InNamespace(namespace))
	if err != nil {
		fmt.Printf(message, namespace, err.Error())
		return
	}

	logsDir, err := ensureLogsDir(dir)
	if err != nil {
		fmt.Printf(message, namespace, err.Error())
		return
	}

	listFile, err := os.Create(logsDir + string(os.PathSeparator) + "role-list.txt")
	if err != nil {
		fmt.Printf(message, namespace, err.Error())
		return
	}

	fn := func() { _ = listFile.Close() }
	defer fn()

	if len(list.Items) > 0 {
		for _, item := range list.Items {
			_, err = fmt.Fprint(listFile, item.GetName()+"\n")
			if err != nil {
				logger.Logf(message, namespace, err.Error())
				return
			}

			d, err := json.Marshal(item)
			if err != nil {
				logger.Logf(message, namespace, err.Error())
				return
			}

			file, err := os.Create(logsDir + string(os.PathSeparator) + "CoherenceRole-" + item.GetName() + ".json")
			if err != nil {
				logger.Logf(message, namespace, err.Error())
				return
			}

			_, err = fmt.Fprint(file, string(d))
			_ = file.Close()
			if err != nil {
				logger.Logf(message, namespace, err.Error())
				return
			}
		}
	} else {
		_, _ = fmt.Fprint(listFile, "No CoherenceRoles resources found in namespace "+namespace)
	}
}

func dumpCoherenceInternals(namespace, dir string, logger Logger) {
	const message = "Could not dump CoherenceInternals for namespace %s due to %s\n"

	f := framework.Global
	list := NewUnstructuredCoherenceInternalList()
	err := f.Client.List(context.TODO(), &list, client.InNamespace(namespace))
	if err != nil {
		logger.Logf(message, namespace, err.Error())
		return
	}

	logsDir, err := ensureLogsDir(dir)
	if err != nil {
		logger.Logf(message, namespace, err.Error())
		return
	}

	listFile, err := os.Create(logsDir + string(os.PathSeparator) + "internal-list.txt")
	if err != nil {
		logger.Logf(message, namespace, err.Error())
		return
	}

	fn := func() { _ = listFile.Close() }
	defer fn()

	if len(list.Items) > 0 {
		for _, item := range list.Items {
			_, err = fmt.Fprint(listFile, item.GetName()+"\n")
			if err != nil {
				logger.Logf(message, namespace, err.Error())
				return
			}

			d, err := json.Marshal(item)
			if err != nil {
				logger.Logf(message, namespace, err.Error())
				return
			}

			file, err := os.Create(logsDir + string(os.PathSeparator) + "CoherenceInternal-" + item.GetName() + ".json")
			if err != nil {
				logger.Logf(message, namespace, err.Error())
				return
			}

			_, err = fmt.Fprint(file, string(d))
			_ = file.Close()
			if err != nil {
				logger.Logf(message, namespace, err.Error())
				return
			}
		}
	} else {
		_, _ = fmt.Fprint(listFile, "No CoherenceInternals resources found in namespace "+namespace)
	}
}

func dumpStatefulSets(namespace, dir string, logger Logger) {
	const message = "Could not dump StatefulSets for namespace %s due to %s\n"

	f := framework.Global
	list, err := f.KubeClient.AppsV1().StatefulSets(namespace).List(metav1.ListOptions{})
	if err != nil {
		logger.Logf(message, namespace, err.Error())
		return
	}

	logsDir, err := ensureLogsDir(dir)
	if err != nil {
		logger.Logf(message, namespace, err.Error())
		return
	}

	listFile, err := os.Create(logsDir + string(os.PathSeparator) + "sts-list.txt")
	if err != nil {
		logger.Logf(message, namespace, err.Error())
		return
	}

	fn := func() { _ = listFile.Close() }
	defer fn()

	if len(list.Items) > 0 {
		for _, item := range list.Items {
			_, err = fmt.Fprint(listFile, item.GetName()+"\n")
			if err != nil {
				logger.Logf(message, namespace, err.Error())
				return
			}

			d, err := json.Marshal(item)
			if err != nil {
				logger.Logf(message, namespace, err.Error())
				return
			}

			file, err := os.Create(logsDir + string(os.PathSeparator) + "StatefulSet-" + item.GetName() + ".json")
			if err != nil {
				logger.Logf(message, namespace, err.Error())
				return
			}

			_, err = fmt.Fprint(file, string(d))
			_ = file.Close()
			if err != nil {
				logger.Logf(message, namespace, err.Error())
				return
			}
		}
	} else {
		_, _ = fmt.Fprint(listFile, "No StatefulSet resources found in namespace "+namespace)
	}
}

func dumpServices(namespace, dir string, logger Logger) {
	const message = "Could not dump Services for namespace %s due to %s\n"

	f := framework.Global
	list, err := f.KubeClient.CoreV1().Services(namespace).List(metav1.ListOptions{})
	if err != nil {
		logger.Logf(message, namespace, err.Error())
		return
	}

	logsDir, err := ensureLogsDir(dir)
	if err != nil {
		logger.Logf(message, namespace, err.Error())
		return
	}

	listFile, err := os.Create(logsDir + string(os.PathSeparator) + "svc-list.txt")
	if err != nil {
		logger.Logf(message, namespace, err.Error())
		return
	}

	fn := func() { _ = listFile.Close() }
	defer fn()

	if len(list.Items) > 0 {
		for _, item := range list.Items {
			_, err = fmt.Fprint(listFile, item.GetName()+"\n")
			if err != nil {
				logger.Logf(message, namespace, err.Error())
				return
			}

			d, err := json.Marshal(item)
			if err != nil {
				logger.Logf(message, namespace, err.Error())
				return
			}

			file, err := os.Create(logsDir + string(os.PathSeparator) + "Service-" + item.GetName() + ".json")
			if err != nil {
				logger.Logf(message, namespace, err.Error())
				return
			}

			_, err = fmt.Fprint(file, string(d))
			_ = file.Close()
			if err != nil {
				logger.Logf(message, namespace, err.Error())
				return
			}
		}
	} else {
		_, _ = fmt.Fprintf(listFile, "No Service resources found in namespace %s", namespace)
	}
}

func dumpRoles(namespace, dir string, logger Logger) {
	const message = "Could not dump Roles for namespace %s due to %s\n"

	f := framework.Global
	list, err := f.KubeClient.RbacV1().Roles(namespace).List(metav1.ListOptions{})
	if err != nil {
		logger.Logf(message, namespace, err.Error())
		return
	}

	logsDir, err := ensureLogsDir(dir)
	if err != nil {
		logger.Logf(message, namespace, err.Error())
		return
	}

	listFile, err := os.Create(logsDir + string(os.PathSeparator) + "role-list.txt")
	if err != nil {
		logger.Logf(message, namespace, err.Error())
		return
	}

	fn := func() { _ = listFile.Close() }
	defer fn()

	if len(list.Items) > 0 {
		for _, item := range list.Items {
			_, err = fmt.Fprint(listFile, item.GetName()+"\n")
			if err != nil {
				logger.Logf(message, namespace, err.Error())
				return
			}

			d, err := json.Marshal(item)
			if err != nil {
				logger.Logf(message, namespace, err.Error())
				return
			}

			file, err := os.Create(logsDir + string(os.PathSeparator) + "Role-" + item.GetName() + ".json")
			if err != nil {
				logger.Logf(message, namespace, err.Error())
				return
			}

			_, err = fmt.Fprint(file, string(d))
			_ = file.Close()
			if err != nil {
				logger.Logf(message, namespace, err.Error())
				return
			}
		}
	} else {
		_, _ = fmt.Fprint(listFile, "No Role resources found in namespace "+namespace)
	}
}

func dumpRoleBindings(namespace, dir string, logger Logger) {
	const message = "Could not dump RoleBindings for namespace %s due to %s\n"

	f := framework.Global
	list, err := f.KubeClient.RbacV1().RoleBindings(namespace).List(metav1.ListOptions{})
	if err != nil {
		logger.Logf(message, namespace, err.Error())
		return
	}

	logsDir, err := ensureLogsDir(dir)
	if err != nil {
		logger.Logf(message, namespace, err.Error())
		return
	}

	listFile, err := os.Create(logsDir + string(os.PathSeparator) + "role-binding-list.txt")
	if err != nil {
		logger.Logf(message, namespace, err.Error())
		return
	}

	fn := func() { _ = listFile.Close() }
	defer fn()

	if len(list.Items) > 0 {
		for _, item := range list.Items {
			_, err = fmt.Fprint(listFile, item.GetName()+"\n")
			if err != nil {
				logger.Logf(message, namespace, err.Error())
				return
			}

			d, err := json.Marshal(item)
			if err != nil {
				logger.Logf(message, namespace, err.Error())
				return
			}

			file, err := os.Create(logsDir + string(os.PathSeparator) + "RoleBinding-" + item.GetName() + ".json")
			if err != nil {
				logger.Logf(message, namespace, err.Error())
				return
			}

			_, err = fmt.Fprint(file, string(d))
			_ = file.Close()
			if err != nil {
				logger.Logf(message, namespace, err.Error())
				return
			}
		}
	} else {
		_, _ = fmt.Fprint(listFile, "No RoleBinding resources found in namespace "+namespace)
	}
}

func dumpServiceAccounts(namespace, dir string, logger Logger) {
	const message = "Could not dump ServiceAccounts for namespace %s due to %s\n"

	f := framework.Global
	list, err := f.KubeClient.CoreV1().ServiceAccounts(namespace).List(metav1.ListOptions{})
	if err != nil {
		logger.Logf(message, namespace, err.Error())
		return
	}

	logsDir, err := ensureLogsDir(dir)
	if err != nil {
		logger.Logf(message, namespace, err.Error())
		return
	}

	listFile, err := os.Create(logsDir + string(os.PathSeparator) + "service-accounts-list.txt")
	if err != nil {
		logger.Logf(message, namespace, err.Error())
		return
	}

	fn := func() { _ = listFile.Close() }
	defer fn()

	if len(list.Items) > 0 {
		for _, item := range list.Items {
			_, err = fmt.Fprint(listFile, item.GetName()+"\n")
			if err != nil {
				logger.Logf(message, namespace, err.Error())
				return
			}

			d, err := json.Marshal(item)
			if err != nil {
				logger.Logf(message, namespace, err.Error())
				return
			}

			file, err := os.Create(logsDir + string(os.PathSeparator) + "ServiceAccount-" + item.GetName() + ".json")
			if err != nil {
				logger.Logf(message, namespace, err.Error())
				return
			}

			_, err = fmt.Fprint(file, string(d))
			_ = file.Close()
			if err != nil {
				logger.Logf(message, namespace, err.Error())
				return
			}
		}
	} else {
		_, _ = fmt.Fprint(listFile, "No ServiceAccount resources found in namespace "+namespace)
	}
}

func DumpPodsForTest(t *testing.T, ctx *framework.TestCtx) {
	namespace, err := ctx.GetNamespace()
	if err == nil {
		dumpPods(namespace, t.Name(), t)
	} else {
		t.Logf("Could not dump Pod logs and state\n")
		t.Log(err)
	}
}

func dumpPods(namespace, dir string, logger Logger) {
	const message = "Could not dump Pods for namespace %s due to %s\n"

	f := framework.Global
	list, err := f.KubeClient.CoreV1().Pods(namespace).List(metav1.ListOptions{})
	if err != nil {
		logger.Logf(message, namespace, err.Error())
		return
	}

	logsDir, err := ensureLogsDir(dir)
	if err != nil {
		logger.Logf(message, namespace, err.Error())
		return
	}

	listFile, err := os.Create(logsDir + string(os.PathSeparator) + "pod-list.txt")
	if err != nil {
		logger.Logf(message, namespace, err.Error())
		return
	}

	fn := func() { _ = listFile.Close() }
	defer fn()

	if len(list.Items) > 0 {
		for _, item := range list.Items {
			_, err = fmt.Fprint(listFile, item.GetName()+"\n")
			if err != nil {
				logger.Logf(message, namespace, err.Error())
				return
			}

			d, err := json.Marshal(item)
			if err != nil {
				logger.Logf(message, namespace, err.Error())
				return
			}

			file, err := os.Create(logsDir + string(os.PathSeparator) + "Pod-" + item.GetName() + ".json")
			if err != nil {
				logger.Logf(message, namespace, err.Error())
				return
			}

			_, err = fmt.Fprint(file, string(d))
			_ = file.Close()
			if err != nil {
				logger.Logf(message, namespace, err.Error())
				return
			}

			DumpPodLog(f.KubeClient, &item, dir, logger)
		}
	} else {
		_, _ = fmt.Fprint(listFile, "No Pod resources found in namespace "+namespace)
	}
}

func ensureLogsDir(subDir string) (string, error) {
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

// Obtain the latest ready time from all of the specified Pods for a given role
func GetLastPodReadyTime(pods []corev1.Pod, role string) metav1.Time {
	t := metav1.NewTime(time.Time{})
	for _, p := range pods {
		if p.GetLabels()["coherenceRole"] == role {
			for _, c := range p.Status.Conditions {
				if c.Type == corev1.PodReady && t.Before(&c.LastTransitionTime) {
					t = c.LastTransitionTime
				}
			}
		}
	}
	return t
}

// Obtain the first ready time from all of the specified Pods for a given role
func GetFirstPodReadyTime(pods []corev1.Pod, role string) metav1.Time {
	t := metav1.NewTime(time.Now())
	for _, p := range pods {
		if p.GetLabels()["coherenceRole"] == role {
			for _, c := range p.Status.Conditions {
				if c.Type == corev1.PodReady && t.After(c.LastTransitionTime.Time) {
					t = c.LastTransitionTime
				}
			}
		}
	}
	return t
}

// Obtain the earliest scheduled time from all of the specified Pods for a given role
func GetFirstPodScheduledTime(pods []corev1.Pod, role string) metav1.Time {
	t := metav1.NewTime(time.Now())
	for _, p := range pods {
		if p.GetLabels()["coherenceRole"] == role {
			for _, c := range p.Status.Conditions {
				if c.Type == corev1.PodScheduled && t.After(c.LastTransitionTime.Time) {
					t = c.LastTransitionTime
				}
			}
		}
	}
	return t
}

func UninstallCrds(t *testing.T) error {
	if err := UninstallCrd(t, "coherenceclusters.coherence.oracle.com"); err != nil {
		return err
	}

	if err := UninstallCrd(t, "coherenceroles.coherence.oracle.com"); err != nil {
		return err
	}

	if err := UninstallCrd(t, "coherenceinternals.coherence.oracle.com"); err != nil {
		return err
	}

	return nil
}

func UninstallCrd(t *testing.T, name string) error {
	var err error

	t.Logf("Will delete CRD %s", name)

	f := framework.Global
	crd := &v1beta1.CustomResourceDefinition{}
	crd.SetName(name)

	mergePatch, err := json.Marshal(map[string]interface{}{
		"metadata": map[string]interface{}{
			"finalizers": []string{},
		},
	})
	if err != nil {
		return err
	}

	patch := client.ConstantPatch(types.MergePatchType, mergePatch)

	t.Logf("Removing finalizer from CRD %s", name)
	if err = f.Client.Patch(context.TODO(), crd, patch); err != nil {
		return err
	}

	t.Logf("Actually deleting finalizer from CRD %s", name)
	return f.Client.Delete(context.TODO(), crd)
}
