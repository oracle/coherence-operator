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
	"github.com/go-logr/logr"
	"github.com/operator-framework/operator-sdk/pkg/status"
	coh "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/controllers"
	"golang.org/x/net/context"
	"io"
	"io/ioutil"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/deprecated/scheme"
	"k8s.io/client-go/rest"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
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
)

// A context for end-to-end tests
type TestContext struct {
	Config     *rest.Config
	Client     client.Client
	KubeClient kubernetes.Interface
	TestEnv    *envtest.Environment
	Manager    ctrl.Manager
	Logger     logr.Logger
}

func (in TestContext) Logf(format string, a ...interface{}) {
	in.Logger.Info(fmt.Sprintf(format, a...))
}

func (in TestContext) Close() {
	in.Logger.Info("tearing down the test environment")

	ns := GetTestNamespace()
	err := in.Client.DeleteAllOf(context.Background(), &coh.Coherence{}, client.InNamespace(ns))
	if err != nil {
		in.Logger.Info("error tearing down the test environment: " + err.Error())
	}

	err = in.TestEnv.Stop()
	if err != nil {
		in.Logger.Info("error stopping test environment: " + err.Error())
	}
}

func NewContext() (TestContext, error) {
	testLogger := zap.New(zap.UseDevMode(true), zap.WriteTo(os.Stdout))

	logf.SetLogger(testLogger)

	// We need a real cluster for these tests
	useCluster := true

	testLogger.WithName("test").Info("bootstrapping test environment")
	testEnv := &envtest.Environment{
		UseExistingCluster:       &useCluster,
		AttachControlPlaneOutput: true,
	}

	var err error
	k8sCfg, err := testEnv.Start()
	if err != nil {
		return TestContext{}, err
	}

	err = corev1.AddToScheme(scheme.Scheme)
	if err != nil {
		return TestContext{}, err
	}
	err = coh.AddToScheme(scheme.Scheme)
	if err != nil {
		return TestContext{}, err
	}

	testEnv.CRDs = []runtime.Object{}

	k8sManager, err := ctrl.NewManager(k8sCfg, ctrl.Options{Scheme: scheme.Scheme})
	if err != nil {
		return TestContext{}, err
	}

	// Create the Coherence controller
	err = (&controllers.CoherenceReconciler{
		Client:    k8sManager.GetClient(),
		Log:       ctrl.Log.WithName("controllers").WithName("Coherence"),
	}).SetupWithManager(k8sManager)
	if err != nil {
		return TestContext{}, err
	}

	// Ensure CRDs exist
	err = coh.EnsureCRDs(k8sCfg)
	if err != nil {
		return TestContext{}, err
	}

	// Start the manager, which will start the controller
	go func() {
		err = k8sManager.Start(ctrl.SetupSignalHandler())
	}()

	k8sClient := k8sManager.GetClient()
	kubeClient, err := kubernetes.NewForConfig(k8sCfg)
	if err != nil {
		return TestContext{}, err
	}

	return TestContext{
		Config:     k8sCfg,
		Client:     k8sClient,
		KubeClient: kubeClient,
		TestEnv:    testEnv,
		Manager:    k8sManager,
		Logger:     testLogger.WithName("test"),
	}, nil
}


type Logger interface {
	Log(args ...interface{})
	Logf(format string, args ...interface{})
}

// WaitForStatefulSetForDeployment waits for a StatefulSet to be created for the specified deployment.
func WaitForStatefulSetForDeployment(ctx TestContext, namespace string, deployment *coh.Coherence, retryInterval, timeout time.Duration) (*appsv1.StatefulSet, error) {
	return WaitForStatefulSet(ctx, namespace, deployment.Name, deployment.Spec.GetReplicas(), retryInterval, timeout)
}

// WaitForStatefulSet waits for a StatefulSet to be created with the specified number of replicas.
func WaitForStatefulSet(ctx TestContext, namespace, stsName string, replicas int32, retryInterval, timeout time.Duration) (*appsv1.StatefulSet, error) {
	var sts *appsv1.StatefulSet

	err := wait.Poll(retryInterval, timeout, func() (done bool, err error) {
		sts, err = ctx.KubeClient.AppsV1().StatefulSets(namespace).Get(context.TODO(), stsName, metav1.GetOptions{})
		if err != nil {
			if apierrors.IsNotFound(err) {
				ctx.Logf("Waiting for availability of StatefulSet %s - NotFound", stsName)
				return false, nil
			}
			ctx.Logf("Waiting for availability of %s StatefulSet - %s", stsName, err.Error())
			return false, err
		}

		if sts.Status.ReadyReplicas == replicas {
			return true, nil
		}
		ctx.Logf("Waiting for full availability of %s StatefulSet (%d/%d)", stsName, sts.Status.ReadyReplicas, replicas)
		return false, nil
	})

	if err != nil && sts != nil {
		d, _ := json.MarshalIndent(sts, "", "    ")
		ctx.Logf("Error waiting for StatefulSet%s", string(d))
	}
	return sts, err
}

// WaitForEndpoints waits for Enpoints for a Service to be created.
func WaitForEndpoints(ctx TestContext, namespace, service string, retryInterval, timeout time.Duration) (*corev1.Endpoints, error) {
	var ep *corev1.Endpoints

	err := wait.Poll(retryInterval, timeout, func() (done bool, err error) {
		ep, err = ctx.KubeClient.CoreV1().Endpoints(namespace).Get(context.TODO(), service, metav1.GetOptions{})
		if err != nil {
			if apierrors.IsNotFound(err) {
				ctx.Logf("Waiting for availability of Endpoints %s - NotFound", service)
				return false, nil
			}
			ctx.Logf("Waiting for availability of %s Endpoints - %s", service, err.Error())
			return false, err
		}
		return true, nil
	})

	if err != nil && ep != nil {
		d, _ := json.MarshalIndent(ep, "", "    ")
		ctx.Logf("Error waiting for Endpoints%s", string(d))
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
func WaitForCoherence(ctx TestContext, namespace, name string, retryInterval, timeout time.Duration, logger Logger) (*coh.Coherence, error) {
	return WaitForCoherenceCondition(ctx, namespace, name, alwaysCondition{}, retryInterval, timeout, logger)
}

// WaitForCoherence waits for a Coherence resource to be created.
func WaitForCoherenceCondition(ctx TestContext, namespace, name string, condition DeploymentStateCondition, retryInterval, timeout time.Duration, logger Logger) (*coh.Coherence, error) {
	var deployment *coh.Coherence

	err := wait.Poll(retryInterval, timeout, func() (done bool, err error) {
		deployment, err = GetCoherence(ctx, namespace, name)
		if err != nil {
			if apierrors.IsNotFound(err) {
				ctx.Logf("Waiting for availability of Coherence resource %s - NotFound", name)
				return false, nil
			}
			ctx.Logf("Waiting for availability of Coherence resource %s - %s", name, err.Error())
			return false, nil
		}
		valid := true
		if condition != nil {
			valid = condition.Test(deployment)
			if !valid {
				ctx.Logf("Waiting for Coherence resource %s to meet condition '%s'", name, condition.String())
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

	err := ctx.Client.Get(goctx.TODO(), opts, d)

	return d, err
}

// WaitForOperatorPods waits for a Coherence Operator Pods to be created.
func WaitForOperatorPods(ctx TestContext, namespace string, retryInterval, timeout time.Duration) ([]corev1.Pod, error) {
	return WaitForPodsWithSelector(ctx, namespace, operatorPodSelector, retryInterval, timeout)
}

// WaitForOperatorPods waits for a Coherence Operator Pods to be created.
func WaitForPodsWithSelector(ctx TestContext, namespace, selector string, retryInterval, timeout time.Duration) ([]corev1.Pod, error) {
	var pods []corev1.Pod

	err := wait.Poll(retryInterval, timeout, func() (done bool, err error) {
		pods, err = ListPodsWithLabelSelector(ctx, namespace, selector)
		if err != nil {
			return false, err
		}
		return len(pods) > 0, nil
	})
	return pods, err
}

// WaitForOperatorDeletion waits for deletion of the Operator Pods.
func WaitForOperatorDeletion(ctx TestContext, namespace string, retryInterval, timeout time.Duration) error {
	return WaitForDeleteOfPodsWithSelector(ctx, namespace, operatorPodSelector, retryInterval, timeout)
}

// WaitForDeleteOfPodsWithSelector waits for Pods to be deleted.
func WaitForDeleteOfPodsWithSelector(ctx TestContext, namespace, selector string, retryInterval, timeout time.Duration) error {
	ctx.Logf("Waiting for Pods in namespace %s with selector '%s' to be deleted", namespace, selector)
	var pods []corev1.Pod

	err := wait.Poll(retryInterval, timeout, func() (done bool, err error) {
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
func WaitForDeletion(ctx TestContext, namespace, name string, resource runtime.Object, retryInterval, timeout time.Duration, logger Logger) error {
	gvk, _ := apiutil.GVKForObject(resource, ctx.Manager.GetScheme())
	ctx.Logf("Waiting for deletion of %v %s/%s", gvk, namespace, name)

	key := types.NamespacedName{Namespace: namespace, Name: name}

	err := wait.Poll(retryInterval, timeout, func() (done bool, err error) {
		err = ctx.Client.Get(context.TODO(), key, resource)
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

// List the Operator Pods that exist - this is Pods with the label "name=coh-operator"
func ListOperatorPods(ctx TestContext, namespace string) ([]corev1.Pod, error) {
	return ListPodsWithLabelSelector(ctx, namespace, operatorPodSelector)
}

// List the Coherence Cluster Pods that exist for a cluster - this is Pods with the label "coherenceCluster=<cluster>"
func ListCoherencePodsForCluster(ctx TestContext, namespace, cluster string) ([]corev1.Pod, error) {
	return ListPodsWithLabelSelector(ctx, namespace, fmt.Sprintf("%s=%s", coh.LabelCoherenceCluster, cluster))
}

// WaitForPodsWithLabel waits for at least the required number of Pods matching the specified labels selector to be created.
func WaitForPodsWithLabel(ctx TestContext, namespace, selector string, count int, retryInterval, timeout time.Duration) ([]corev1.Pod, error) {
	var pods []corev1.Pod

	err := wait.Poll(retryInterval, timeout, func() (done bool, err error) {
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

// List the Pods that exist for a deployment - this is Pods with the label "coherenceDeployment=<deployment>"
func ListCoherencePodsForDeployment(ctx TestContext, namespace, deployment string) ([]corev1.Pod, error) {
	selector := fmt.Sprintf("%s=%s", coh.LabelCoherenceDeployment, deployment)
	return ListPodsWithLabelSelector(ctx, namespace, selector)
}

// List the Coherence Cluster Pods that exist for a given label selector.
func ListPodsWithLabelSelector(ctx TestContext, namespace, selector string) ([]corev1.Pod, error) {
	opts := metav1.ListOptions{LabelSelector: selector}

	list, err := ctx.KubeClient.CoreV1().Pods(namespace).List(context.TODO(), opts)
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
func WaitForCoherenceCleanup(ctx TestContext, namespace string) error {
	ctx.Logf("Waiting for clean-up of Coherence resources in namespace %s", namespace)

	list := &coh.CoherenceList{}
	err := ctx.Client.List(goctx.TODO(), list, client.InNamespace(namespace))
	if err != nil {
		return err
	}

	// Delete all of the Coherence resources
	for _, r := range list.Items {
		ctx.Logf("Deleting Coherence resource %s in namespace %s", r.Name, r.Namespace)
		err = ctx.Client.Delete(goctx.TODO(), &r)
		if err != nil {
			ctx.Logf("Error deleting Coherence resource %s - %s", r.Name, err.Error())
		}
	}

	// Wait for removal of the Coherence resources
	err = wait.Poll(RetryInterval, Timeout, func() (done bool, err error) {
		err = ctx.Client.List(goctx.TODO(), list, client.InNamespace(namespace))
		if err == nil || isNoResources(err) || errors.IsNotFound(err) {
			if len(list.Items) > 0 {
				ctx.Logf("Waiting for deletion of %d Coherence resources", len(list.Items))
				return false, nil
			}
			return true, nil
		} else {
			ctx.Logf("Error waiting for deletion of Coherence resources: %s\n%v", err.Error(), err)
			return false, nil
		}
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
	err := wait.Poll(RetryInterval, Timeout, func() (done bool, err error) {
		list, err := ListOperatorPods(ctx, namespace)
		if err == nil {
			if len(list) > 0 {
				ctx.Logf("Waiting for deletion of %d Coherence Operator Pods", len(list))
				return false, nil
			}
			return true, nil
		} else {
			ctx.Logf("Error waiting for deletion of Coherence Operator Pods: %s", err.Error())
			return false, nil
		}
	})

	ctx.Logger.Info("Coherence Operator Pods deleted")
	return err
}

// Dump the Operator Pod log to a file.
func DumpOperatorLog(ctx TestContext, namespace, directory string) {
	list, err := ctx.KubeClient.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: "name=coherence-operator"})
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

// Dump the Pod log to a file.
func DumpPodLog(ctx TestContext, pod *corev1.Pod, directory string) {
	logs, err := FindTestLogsDir()
	if err != nil {
		ctx.Logger.Info("cannot capture logs due to " + err.Error())
		return
	}

	ctx.Logger.Info("Capturing Pod logs for " + pod.Name)
	pathSep := string(os.PathSeparator)

	for _, container := range pod.Spec.Containers {
		var err error
		res := ctx.KubeClient.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, &corev1.PodLogOptions{Container: container.Name})
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
					ctx.Logger.Info("cannot capture logs for Pod " + pod.Name + " container " + container.Name + " due to " + err.Error())
				}
			} else {
				ctx.Logger.Info("cannot capture logs for Pod " + pod.Name + " container " + container.Name + " due to " + err.Error())
			}
		} else {
			ctx.Logger.Info("cannot capture logs for Pod " + pod.Name + " container " + container.Name + " due to " + err.Error())
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

// Dump the operator logs
func DumpOperatorLogs(t *testing.T, ctx TestContext) {
	opNamespace := GetOperatorTestNamespace()
	DumpOperatorLog(ctx, opNamespace, t.Name())
	watchNamespace := GetTestNamespace()
	DumpState(ctx, watchNamespace, t.Name())
}

func DumpState(ctx TestContext, namespace, dir string) {
	dumpCoherences(namespace, dir, ctx)
	dumpStatefulSets(namespace, dir, ctx)
	dumpServices(namespace, dir, ctx)
	dumpPods(namespace, dir, ctx)
	dumpRbacRoles(namespace, dir, ctx)
	dumpRbacRoleBindings(namespace, dir, ctx)
	dumpServiceAccounts(namespace, dir, ctx)
}

func dumpCoherences(namespace, dir string, ctx TestContext) {
	const message = "Could not dump Coherence resource for namespace %s due to %s"

	list := coh.CoherenceList{}
	err := ctx.Client.List(context.TODO(), &list, client.InNamespace(namespace))
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

	list, err := ctx.KubeClient.AppsV1().StatefulSets(namespace).List(context.TODO(), metav1.ListOptions{})
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

	list, err := ctx.KubeClient.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{})
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

func dumpRbacRoles(namespace, dir string, ctx TestContext) {
	const message = "Could not dump RBAC Roles for namespace %s due to %s"

	list, err := ctx.KubeClient.RbacV1().Roles(namespace).List(context.TODO(), metav1.ListOptions{})
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

	list, err := ctx.KubeClient.RbacV1().RoleBindings(namespace).List(context.TODO(), metav1.ListOptions{})
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

	list, err := ctx.KubeClient.CoreV1().ServiceAccounts(namespace).List(context.TODO(), metav1.ListOptions{})
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
	namespace := GetOperatorTestNamespace()
	dumpPods(namespace, t.Name(), ctx)
}

func dumpPods(namespace, dir string, ctx TestContext) {
	const message = "Could not dump Pods for namespace %s due to %s"

	list, err := ctx.KubeClient.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
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
