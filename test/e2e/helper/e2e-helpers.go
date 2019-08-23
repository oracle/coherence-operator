// helper package contains various helpers for use in end-to-end testing.
package helper

import (
	goctx "context"
	"fmt"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	"github.com/oracle/coherence-operator/pkg/apis"
	coh "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	"io"
	"io/ioutil"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"os"
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

var (
	RetryInterval        = time.Second * 5
	Timeout              = time.Second * 60
	CleanupRetryInterval = time.Second * 1
	CleanupTimeout       = time.Second * 5
)

var tenSeconds int32 = 10

var Readiness = &coh.ReadinessProbeSpec{
	InitialDelaySeconds: &tenSeconds,
	PeriodSeconds:       &tenSeconds,
}

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

func DefaultCleanup(ctx *framework.TestCtx) *framework.CleanupOptions {
	return &framework.CleanupOptions{TestContext: ctx, Timeout: Timeout, RetryInterval: RetryInterval}
}

// WaitForStatefulSetForRole waits for a StatefulSet to be created for the specified role.
func WaitForStatefulSetForRole(kubeclient kubernetes.Interface, namespace string, cluster *coh.CoherenceCluster, role coh.CoherenceRoleSpec, retryInterval, timeout time.Duration, logger Logger) (*appsv1.StatefulSet, error) {
	return WaitForStatefulSet(kubeclient, namespace, role.GetFullRoleName(cluster), role.GetReplicas(), retryInterval, timeout, logger)
}

// WaitForStatefulSet waits for a StatefulSet to be created with the specified number of replicas.
func WaitForStatefulSet(kubeclient kubernetes.Interface, namespace, name string, replicas int32, retryInterval, timeout time.Duration, logger Logger) (*appsv1.StatefulSet, error) {
	var sts *appsv1.StatefulSet

	err := wait.Poll(retryInterval, timeout, func() (done bool, err error) {
		sts, err = kubeclient.AppsV1().StatefulSets(namespace).Get(name, metav1.GetOptions{IncludeUninitialized: true})
		if err != nil {
			if apierrors.IsNotFound(err) {
				logger.Logf("Waiting for availability of %s StatefulSet - NotFound\n", name)
				return false, nil
			}
			logger.Logf("Waiting for availability of %s StatefulSet - %s\n", name, err.Error())
			return false, err
		}

		if sts.Status.ReadyReplicas == replicas {
			return true, nil
		}
		logger.Logf("Waiting for full availability of %s StatefulSet (%d/%d)\n", name, sts.Status.ReadyReplicas, replicas)
		return false, nil
	})

	return sts, err
}

// WaitForCoherenceRole waits for a CoherenceRole to be created.
func WaitForCoherenceRole(f *framework.Framework, namespace, name string, retryInterval, timeout time.Duration, logger Logger) (*coh.CoherenceRole, error) {
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

		return true, nil
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
func WaitForOperatorPods(kubeclient kubernetes.Interface, namespace string, retryInterval, timeout time.Duration) ([]corev1.Pod, error) {
	var pods []corev1.Pod

	err := wait.Poll(retryInterval, timeout, func() (done bool, err error) {
		pods, err = ListOperatorPods(kubeclient, namespace)
		if err != nil {
			return false, err
		}
		return len(pods) > 0, nil
	})
	return pods, err
}

// List the Operator Pods that exist - this is Pods with the label "name=coh-operator"
func ListOperatorPods(kubeclient kubernetes.Interface, namespace string) ([]corev1.Pod, error) {
	opts := metav1.ListOptions{
		LabelSelector: "name=coherence-operator",
	}

	list, err := kubeclient.CoreV1().Pods(namespace).List(opts)
	if err != nil {
		return []corev1.Pod{}, err
	}

	return list.Items, nil
}

// List the Coherence Cluster Pods that exist - this is Pods with the label "coherenceCluster=<cluster>,coherenceRole=<role>"
func ListCoherencePods(kubeclient kubernetes.Interface, namespace, cluster, role string) ([]corev1.Pod, error) {
	opts := metav1.ListOptions{LabelSelector: "coherenceCluster=" + cluster + ",coherenceRole=" + role}

	list, err := kubeclient.CoreV1().Pods(namespace).List(opts)
	if err != nil {
		return []corev1.Pod{}, err
	}

	return list.Items, nil
}

// WaitForPodReady waits for a Pods to be ready.
func WaitForPodReady(kubeclient kubernetes.Interface, namespace, name string, retryInterval, timeout time.Duration) error {
	err := wait.Poll(retryInterval, timeout, func() (done bool, err error) {
		p, err := kubeclient.CoreV1().Pods(namespace).Get(name, metav1.GetOptions{})
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
	opts := &client.ListOptions{}
	opts.InNamespace(namespace)

	list := &coh.CoherenceClusterList{}
	err := f.Client.List(goctx.TODO(), opts, list)
	if err != nil {
		return err
	}

	for _, c := range list.Items {
		fmt.Printf("Deleting CoherenceCluster %s in namespace %s", c.Name, c.Namespace)
		err = f.Client.Delete(goctx.TODO(), &c)
		if err != nil {
			fmt.Printf("Error deleting CoherenceCluster %s - %s\n", c.Name, err.Error())
		}
	}

	u := unstructured.UnstructuredList{}

	u.SetGroupVersionKind(schema.GroupVersionKind{Group: "coherence.oracle.com", Version: "v1", Kind: "CoherenceInternal"})

	err = f.Client.Client.List(goctx.TODO(), opts, &u)
	if err != nil {
		fmt.Printf("Error listing CoherenceInternal resources - %s\n", err.Error())
	} else {
		if len(u.Items) > 0 {
			fmt.Printf("Remaining CoherenceInternal resources in namespace %s (%d):\n", namespace, len(u.Items))
			for i, ci := range u.Items {
				fmt.Printf("%d: %s\n", i, ci.GetName())
			}
		} else {
			fmt.Printf("Zero CoherenceInternal resources remain in namespace %s\n", namespace)
		}
	}

	var empty []string

	err = wait.Poll(RetryInterval, Timeout, func() (done bool, err error) {
		err = f.Client.List(goctx.TODO(), opts, &u)
		if err == nil {
			if len(u.Items) > 0 {
				fmt.Printf("Waiting for deletion of %d CoherenceInternal resources\n", len(u.Items))

				for _, ci := range u.Items {
					ci.SetFinalizers(empty)
					fmt.Printf("Removing finalizers on CoherenceInternal resources %s\n", ci.GetName())
					_ = f.Client.Update(goctx.TODO(), &ci)
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

// WaitForOperatorCleanup waits until there are no Operator Pods in the test namespace.
func WaitForOperatorCleanup(kubeClient kubernetes.Interface, namespace string, logger Logger) error {
	// wait for all CoherenceInternal resources to be deleted
	opts := &client.ListOptions{}
	opts.InNamespace(namespace)

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
	logs := os.Getenv("TEST_LOGS")
	if logs == "" {
		logger.Log("cannot capture logs as log folder env var TEST_LOGS is not set")
	}

	res := kubeClient.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, &corev1.PodLogOptions{})
	s, err := res.Stream()
	if err == nil {
		name := logs + "/" + directory
		err = os.MkdirAll(name, os.ModePerm)
		if err == nil {
			out, err := os.Create(name + "/" + pod.Name + ".log")
			if err == nil {
				_, err = io.Copy(out, s)
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
	namespace, err := ctx.GetNamespace()
	if err == nil {
		DumpOperatorLog(framework.Global.KubeClient, namespace, t.Name(), t)
	}
	ctx.Cleanup()
}
