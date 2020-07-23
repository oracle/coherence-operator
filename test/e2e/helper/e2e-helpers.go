/*
 * Copyright (c) 2019, 2020 Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

// helper package contains various helpers for use in end-to-end testing.
package helper

import (
	goctx "context"
	"encoding/json"
	"fmt"
	"github.com/operator-framework/operator-sdk/pkg/status"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	coh "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/pkg/apis"
	"golang.org/x/net/context"
	"io"
	"io/ioutil"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"strings"
	"testing"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
)

type Logger interface {
	Log(args ...interface{})
	Logf(format string, args ...interface{})
}

func CreateTestContext(t *testing.T) *framework.Context {
	f := framework.Global
	testCtx := framework.NewContext(t)

	namespace, err := testCtx.GetWatchNamespace()
	if err != nil {
		t.Fatal(err)
		return nil
	}

	testCtx.AddCleanupFn(func() error {
		return WaitForCoherenceCleanup(f, namespace)
	})

	cleanup := framework.CleanupOptions{TestContext: testCtx, Timeout: CleanupTimeout, RetryInterval: CleanupRetryInterval}

	err = testCtx.InitializeClusterResources(&cleanup)
	if err != nil && !errors.IsAlreadyExists(err) {
		t.Fatal(err)
		return nil
	}

	list := &coh.CoherenceList{}

	err = framework.AddToFrameworkScheme(apis.AddToScheme, list)
	if err != nil {
		t.Fatal(err)
		return nil
	}

	return testCtx
}

func DefaultCleanup(ctx *framework.Context) *framework.CleanupOptions {
	return &framework.CleanupOptions{TestContext: ctx, Timeout: Timeout, RetryInterval: RetryInterval}
}

// WaitForStatefulSetForDeployment waits for a StatefulSet to be created for the specified deployment.
func WaitForStatefulSetForDeployment(kubeclient kubernetes.Interface, namespace string, deployment *coh.Coherence, retryInterval, timeout time.Duration, logger Logger) (*appsv1.StatefulSet, error) {
	return WaitForStatefulSet(kubeclient, namespace, deployment.Name, deployment.Spec.GetReplicas(), retryInterval, timeout, logger)
}

// WaitForStatefulSet waits for a StatefulSet to be created with the specified number of replicas.
func WaitForStatefulSet(kubeclient kubernetes.Interface, namespace, stsName string, replicas int32, retryInterval, timeout time.Duration, logger Logger) (*appsv1.StatefulSet, error) {
	var sts *appsv1.StatefulSet

	err := wait.Poll(retryInterval, timeout, func() (done bool, err error) {
		sts, err = kubeclient.AppsV1().StatefulSets(namespace).Get(context.TODO(), stsName, metav1.GetOptions{})
		if err != nil {
			if apierrors.IsNotFound(err) {
				logger.Logf("Waiting for availability of StatefulSet %s - NotFound\n", stsName)
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
		d, _ := json.MarshalIndent(sts, "", "    ")
		logger.Logf("Error waiting for StatefulSet\n%s", string(d))
	}
	return sts, err
}

// WaitForEndpoints waits for Enpoints for a Service to be created.
func WaitForEndpoints(kubeclient kubernetes.Interface, namespace, service string, retryInterval, timeout time.Duration, logger Logger) (*corev1.Endpoints, error) {
	var ep *corev1.Endpoints

	err := wait.Poll(retryInterval, timeout, func() (done bool, err error) {
		ep, err = kubeclient.CoreV1().Endpoints(namespace).Get(context.TODO(), service, metav1.GetOptions{})
		if err != nil {
			if apierrors.IsNotFound(err) {
				logger.Logf("Waiting for availability of Endpoints %s - NotFound\n", service)
				return false, nil
			}
			logger.Logf("Waiting for availability of %s Endpoints - %s\n", service, err.Error())
			return false, err
		}
		return true, nil
	})

	if err != nil && ep != nil {
		d, _ := json.MarshalIndent(ep, "", "    ")
		logger.Logf("Error waiting for Endpoints\n%s", string(d))
	}
	return ep, err
}

// A function that takes a deployment and determines whether it meets a condition
type DeploymentStateCondition interface {
	Test(*coh.Coherence) bool
	String() string
}

// An always true DeploymentStateCondition
type alwaysCondition struct{}

func (a alwaysCondition) Test(*coh.Coherence) bool {
	return true
}

func (a alwaysCondition) String() string {
	return "true"
}

type replicaCountCondition struct {
	replicas int32
}

func (in replicaCountCondition) Test(d *coh.Coherence) bool {
	return d.Status.ReadyReplicas == in.replicas
}

func (in replicaCountCondition) String() string {
	return fmt.Sprintf("replicas == %d", in.replicas)
}

func ReplicaCountCondition(replicas int32) DeploymentStateCondition {
	return replicaCountCondition{replicas: replicas}
}

type phaseCondition struct {
	phase status.ConditionType
}

func (in phaseCondition) Test(d *coh.Coherence) bool {
	return d.Status.Phase == in.phase
}

func (in phaseCondition) String() string {
	return fmt.Sprintf("phase == %s", in.phase)
}

func StatusPhaseCondition(phase status.ConditionType) DeploymentStateCondition {
	return phaseCondition{phase: phase}
}

// WaitForCoherence waits for a Coherence resource to be created.
func WaitForCoherence(f *framework.Framework, namespace, name string, retryInterval, timeout time.Duration, logger Logger) (*coh.Coherence, error) {
	return WaitForCoherenceCondition(f, namespace, name, alwaysCondition{}, retryInterval, timeout, logger)
}

// WaitForCoherence waits for a Coherence resource to be created.
func WaitForCoherenceCondition(f *framework.Framework, namespace, name string, condition DeploymentStateCondition, retryInterval, timeout time.Duration, logger Logger) (*coh.Coherence, error) {
	var deployment *coh.Coherence

	err := wait.Poll(retryInterval, timeout, func() (done bool, err error) {
		deployment, err = GetCoherence(f, namespace, name)
		if err != nil {
			if apierrors.IsNotFound(err) {
				logger.Logf("Waiting for availability of Coherence resource %s - NotFound\n", name)
				return false, nil
			}
			logger.Logf("Waiting for availability of Coherence resource %s - %s\n", name, err.Error())
			return false, nil
		}
		valid := true
		if condition != nil {
			valid = condition.Test(deployment)
			if !valid {
				logger.Logf("Waiting for Coherence resource %s to meet condition '%s'\n", name, condition.String())
			}
		}
		return valid, nil
	})

	return deployment, err
}

// GetCoherence gets the specified Coherence resource
func GetCoherence(f *framework.Framework, namespace, name string) (*coh.Coherence, error) {
	opts := client.ObjectKey{Namespace: namespace, Name: name}
	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
		},
	}

	err := f.Client.Get(goctx.TODO(), opts, d)

	return d, err
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
func WaitForDeletion(f *framework.Framework, namespace, name string, resource runtime.Object, retryInterval, timeout time.Duration, logger Logger) error {
	gvk, _ := apiutil.GVKForObject(resource, f.Scheme)
	logger.Logf("Waiting for deletion of %v %s/%s", gvk, namespace, name)

	key := types.NamespacedName{Namespace: namespace, Name: name}

	err := wait.Poll(retryInterval, timeout, func() (done bool, err error) {
		err = f.Client.Get(context.TODO(), key, resource)
		switch {
		case err != nil && errors.IsNotFound(err):
			return true, nil
		case err != nil && !errors.IsNotFound(err):
			logger.Logf("Waiting for deletion of %v %s/%s - Error=%s", gvk, namespace, name, err)
			return false, err
		default:
			logger.Logf("Still waiting for deletion of %v %s/%s", gvk, namespace, name)
			return false, nil
		}
	})

	logger.Logf("Finished waiting for deletion of %s %s/%s - Error=%s", gvk, namespace, name, err)
	return err
}

// List the Operator Pods that exist - this is Pods with the label "name=coh-operator"
func ListOperatorPods(client kubernetes.Interface, namespace string) ([]corev1.Pod, error) {
	return ListPodsWithLabelSelector(client, namespace, operatorPodSelector)
}

// List the Coherence Cluster Pods that exist for a cluster - this is Pods with the label "coherenceCluster=<cluster>"
func ListCoherencePodsForCluster(client kubernetes.Interface, namespace, cluster string) ([]corev1.Pod, error) {
	return ListPodsWithLabelSelector(client, namespace, fmt.Sprintf("%s=%s", coh.LabelCoherenceCluster, cluster))
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

// List the Pods that exist for a deployment - this is Pods with the label "coherenceDeployment=<deployment>"
func ListCoherencePodsForDeployment(client kubernetes.Interface, namespace, deployment string) ([]corev1.Pod, error) {
	selector := fmt.Sprintf("%s=%s", coh.LabelCoherenceDeployment, deployment)
	return ListPodsWithLabelSelector(client, namespace, selector)
}

// List the Coherence Cluster Pods that exist for a given label selector.
func ListPodsWithLabelSelector(k8s kubernetes.Interface, namespace, selector string) ([]corev1.Pod, error) {
	opts := metav1.ListOptions{LabelSelector: selector}

	list, err := k8s.CoreV1().Pods(namespace).List(context.TODO(), opts)
	if err != nil {
		return []corev1.Pod{}, err
	}

	return list.Items, nil
}

// WaitForPodReady waits for a Pods to be ready.
func WaitForPodReady(k8s kubernetes.Interface, namespace, name string, retryInterval, timeout time.Duration) error {
	err := wait.Poll(retryInterval, timeout, func() (done bool, err error) {
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
func WaitForCoherenceCleanup(f *framework.Framework, namespace string) error {
	fmt.Printf("Waiting for clean-up of Coherence resources in namespace %s\n", namespace)

	list := &coh.CoherenceList{}
	err := f.Client.List(goctx.TODO(), list, client.InNamespace(namespace))
	if err != nil {
		return err
	}

	// Delete all of the Coherence resources
	for _, r := range list.Items {
		fmt.Printf("Deleting Coherence resource %s in namespace %s\n", r.Name, r.Namespace)
		err = f.Client.Delete(goctx.TODO(), &r)
		if err != nil {
			fmt.Printf("Error deleting Coherence resource %s - %s\n", r.Name, err.Error())
		}
	}

	// Wait for removal of the Coherence resources
	err = wait.Poll(RetryInterval, Timeout, func() (done bool, err error) {
		err = f.Client.List(goctx.TODO(), list, client.InNamespace(namespace))
		if err == nil || isNoResources(err) || errors.IsNotFound(err) {
			if len(list.Items) > 0 {
				fmt.Printf("Waiting for deletion of %d Coherence resources\n", len(list.Items))
				return false, nil
			}
			return true, nil
		} else {
			fmt.Printf("Error waiting for deletion of Coherence resources: %s\n%v\n", err.Error(), err)
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
	list, err := kubeClient.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: "name=coherence-operator"})
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

	logger.Log("Capturing Pod logs for " + pod.Name)
	pathSep := string(os.PathSeparator)

	for _, container := range pod.Spec.Containers {
		var err error
		res := kubeClient.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, &corev1.PodLogOptions{Container: container.Name})
		s, err := res.Stream(context.TODO())
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
				} else {
					logger.Log("cannot capture logs for Pod " + pod.Name + " container " + container.Name + " due to " + err.Error())
				}
			} else {
				logger.Log("cannot capture logs for Pod " + pod.Name + " container " + container.Name + " due to " + err.Error())
			}
		} else {
			logger.Log("cannot capture logs for Pod " + pod.Name + " container " + container.Name + " due to " + err.Error())
		}
	}
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

	return ioutil.ReadAll(f)
}

// Dump the operator logs and clean-up the test context
func DumpOperatorLogsAndCleanup(t *testing.T, ctx *framework.Context) {
	DumpOperatorLogs(t, ctx)
	ctx.Cleanup()
}

// Dump the operator logs and clean-up the test context
func DumpOperatorLogs(t *testing.T, ctx *framework.Context) {
	opNamespace, err := ctx.GetOperatorNamespace()
	watchNamespace, err := ctx.GetWatchNamespace()
	if err == nil {
		DumpOperatorLog(framework.Global.KubeClient, opNamespace, t.Name(), t)
		DumpState(watchNamespace, t.Name(), t)
	} else {
		t.Logf("Could not dump logs and state\n")
		t.Log(err)
	}
	ctx.Cleanup()
}

func DumpState(namespace, dir string, logger Logger) {
	dumpCoherences(namespace, dir, logger)
	dumpStatefulSets(namespace, dir, logger)
	dumpServices(namespace, dir, logger)
	dumpPods(namespace, dir, logger)
	dumpRbacRoles(namespace, dir, logger)
	dumpRbacRoleBindings(namespace, dir, logger)
	dumpServiceAccounts(namespace, dir, logger)
}

func dumpCoherences(namespace, dir string, logger Logger) {
	const message = "Could not dump Coherence resource for namespace %s due to %s\n"

	f := framework.Global
	list := coh.CoherenceList{}
	err := f.Client.List(context.TODO(), &list, client.InNamespace(namespace))
	if err != nil {
		fmt.Printf(message, namespace, err.Error())
		return
	}

	logsDir, err := EnsureLogsDir(dir)
	if err != nil {
		fmt.Printf(message, namespace, err.Error())
		return
	}

	listFile, err := os.Create(logsDir + string(os.PathSeparator) + "deployments-list.txt")
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

			d, err := json.MarshalIndent(item, "", "    ")
			if err != nil {
				logger.Logf(message, namespace, err.Error())
				return
			}

			file, err := os.Create(logsDir + string(os.PathSeparator) + "Coherence-" + item.GetName() + ".json")
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
		_, _ = fmt.Fprint(listFile, "No Coherence resources found in namespace "+namespace)
	}
}

func dumpStatefulSets(namespace, dir string, logger Logger) {
	const message = "Could not dump StatefulSets for namespace %s due to %s\n"

	f := framework.Global
	list, err := f.KubeClient.AppsV1().StatefulSets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Logf(message, namespace, err.Error())
		return
	}

	logsDir, err := EnsureLogsDir(dir)
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

			d, err := json.MarshalIndent(item, "", "    ")
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
	list, err := f.KubeClient.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Logf(message, namespace, err.Error())
		return
	}

	logsDir, err := EnsureLogsDir(dir)
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

			d, err := json.MarshalIndent(item, "", "    ")
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

func dumpRbacRoles(namespace, dir string, logger Logger) {
	const message = "Could not dump RBAC Roles for namespace %s due to %s\n"

	f := framework.Global
	list, err := f.KubeClient.RbacV1().Roles(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Logf(message, namespace, err.Error())
		return
	}

	logsDir, err := EnsureLogsDir(dir)
	if err != nil {
		logger.Logf(message, namespace, err.Error())
		return
	}

	listFile, err := os.Create(logsDir + string(os.PathSeparator) + "rbac-role-list.txt")
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

			d, err := json.MarshalIndent(item, "", "    ")
			if err != nil {
				logger.Logf(message, namespace, err.Error())
				return
			}

			file, err := os.Create(logsDir + string(os.PathSeparator) + "RBAC-Role-" + item.GetName() + ".json")
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
		_, _ = fmt.Fprint(listFile, "No RBAC Role resources found in namespace "+namespace)
	}
}

func dumpRbacRoleBindings(namespace, dir string, logger Logger) {
	const message = "Could not dump RBAC RoleBindings for namespace %s due to %s\n"

	f := framework.Global
	list, err := f.KubeClient.RbacV1().RoleBindings(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Logf(message, namespace, err.Error())
		return
	}

	logsDir, err := EnsureLogsDir(dir)
	if err != nil {
		logger.Logf(message, namespace, err.Error())
		return
	}

	listFile, err := os.Create(logsDir + string(os.PathSeparator) + "rbac-role-binding-list.txt")
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

			d, err := json.MarshalIndent(item, "", "    ")
			if err != nil {
				logger.Logf(message, namespace, err.Error())
				return
			}

			file, err := os.Create(logsDir + string(os.PathSeparator) + "RBAC-RoleBinding-" + item.GetName() + ".json")
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
		_, _ = fmt.Fprint(listFile, "No RBAC RoleBinding resources found in namespace "+namespace)
	}
}

func dumpServiceAccounts(namespace, dir string, logger Logger) {
	const message = "Could not dump ServiceAccounts for namespace %s due to %s\n"

	f := framework.Global
	list, err := f.KubeClient.CoreV1().ServiceAccounts(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Logf(message, namespace, err.Error())
		return
	}

	logsDir, err := EnsureLogsDir(dir)
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

			d, err := json.MarshalIndent(item, "", "    ")
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

func DumpPodsForTest(t *testing.T, ctx *framework.Context) {
	namespace, err := ctx.GetWatchNamespace()
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
	list, err := f.KubeClient.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Logf(message, namespace, err.Error())
		return
	}

	logsDir, err := EnsureLogsDir(dir)
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

			d, err := json.MarshalIndent(item, "", "    ")
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

// Obtain the latest ready time from all of the specified Pods for a given deployment
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

// Obtain the first ready time from all of the specified Pods for a given deployment
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

// Obtain the earliest scheduled time from all of the specified Pods for a given deployment
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
