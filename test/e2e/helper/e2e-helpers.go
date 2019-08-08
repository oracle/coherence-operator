// helper package contains various helpers for use in end-to-end testing.
package helper

import (
	goctx "context"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	"github.com/oracle/coherence-operator/pkg/apis"
	coherence "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
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

var Readiness = &coherence.ReadinessProbeSpec{
	InitialDelaySeconds: &tenSeconds,
	PeriodSeconds:       &tenSeconds,
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
		return waitForCleanup(t, f, namespace)
	})

	cleanup := framework.CleanupOptions{TestContext: testCtx, Timeout: CleanupTimeout, RetryInterval: CleanupRetryInterval}

	err = testCtx.InitializeClusterResources(&cleanup)
	if err != nil && !errors.IsAlreadyExists(err) {
		t.Fatal(err)
		return nil
	}

	clusterList := &coherence.CoherenceClusterList{}

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

// WaitForStatefulSet waits for a StatefulSet to be created with the specified number of replicas.
func WaitForStatefulSet(t *testing.T, kubeclient kubernetes.Interface, namespace, name string, replicas int32, retryInterval, timeout time.Duration) (*appsv1.StatefulSet, error) {
	var sts *appsv1.StatefulSet

	err := wait.Poll(retryInterval, timeout, func() (done bool, err error) {
		sts, err = kubeclient.AppsV1().StatefulSets(namespace).Get(name, metav1.GetOptions{IncludeUninitialized: true})
		if err != nil {
			if apierrors.IsNotFound(err) {
				t.Logf("Waiting for availability of %s StatefulSet - NotFound\n", name)
				return false, nil
			}
			t.Logf("Waiting for availability of %s StatefulSet - %s\n", name, err.Error())
			return false, err
		}

		if sts.Status.ReadyReplicas == replicas {
			return true, nil
		}
		t.Logf("Waiting for full availability of %s StatefulSet (%d/%d)\n", name, sts.Status.ReadyReplicas, replicas)
		return false, nil
	})

	return sts, err
}

// WaitForCoherenceRole waits for a CoherenceRole to be created.
func WaitForCoherenceRole(t *testing.T, f *framework.Framework, namespace, name string, retryInterval, timeout time.Duration) (*coherence.CoherenceRole, error) {
	var role *coherence.CoherenceRole

	err := wait.Poll(retryInterval, timeout, func() (done bool, err error) {
		role, err = GetCoherenceRole(f, namespace, name)
		if err != nil {
			if apierrors.IsNotFound(err) {
				t.Logf("Waiting for availability of CoherenceRole %s - NotFound\n", name)
				return false, nil
			}
			t.Logf("Waiting for availability of CoherenceRole %s - %s\n", name, err.Error())
			return false, err
		}

		return true, nil
	})

	return role, err
}

// GetCoherenceRole gets the specified CoherenceRole
func GetCoherenceRole(f *framework.Framework, namespace, name string) (*coherence.CoherenceRole, error) {
	opts := client.ObjectKey{Namespace: namespace, Name: name}
	role := &coherence.CoherenceRole{
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

// List the Operator Pods that exist - this is Pods with the label "name=coherence-operator"
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
func waitForCleanup(t *testing.T, f *framework.Framework, namespace string) error {
	t.Logf("Waiting for clean-up of CoherenceInternal resources in namespace %s\n", namespace)

	// wait for all CoherenceInternal resources to be deleted
	opts := &client.ListOptions{}
	opts.InNamespace(namespace)

	list := &coherence.CoherenceClusterList{}
	err := f.Client.List(goctx.TODO(), opts, list)
	if err != nil {
		return err
	}

	for _, c := range list.Items {
		t.Logf("Deleting CoherenceCluster %s in namespace %s", c.Name, c.Namespace)
		err = f.Client.Delete(goctx.TODO(), &c)
		if err != nil {
			t.Logf("Error deleting CoherenceCluster %s - %s\n", c.Name, err.Error())
		}
	}

	u := unstructured.UnstructuredList{}

	u.SetGroupVersionKind(schema.GroupVersionKind{Group: "coherence.oracle.com", Version: "v1", Kind: "CoherenceInternal"})

	err = f.Client.Client.List(goctx.TODO(), opts, &u)
	if err != nil {
		t.Logf("Error listing CoherenceInternal resources - %s\n", err.Error())
	} else {
		if len(u.Items) > 0 {
			t.Logf("Remaining CoherenceInternal resources in namespace %s:\n", namespace)
			for _, ci := range u.Items {
				t.Log(ci.GetName())
			}
		} else {
			t.Logf("Zero CoherenceInternal resources remain in namespace %s\n", namespace)
		}
	}

	var empty []string

	err = wait.Poll(RetryInterval, Timeout, func() (done bool, err error) {
		err = f.Client.List(goctx.TODO(), opts, &u)
		if err == nil {
			if len(u.Items) > 0 {
				t.Logf("Waiting for deletion of %d CoherenceInternal resources\n", len(u.Items))

				for _, ci := range u.Items {
					ci.SetFinalizers(empty)
					t.Logf("Removing finalizers on CoherenceInternal resources %s\n", ci.GetName())
					_ = f.Client.Update(goctx.TODO(), &ci)
				}

				return false, nil
			}
			return true, nil
		} else {
			t.Logf("Error waiting for deletion of CoherenceInternal resources: %s\n", err.Error())
			return false, nil
		}
	})

	return err
}
