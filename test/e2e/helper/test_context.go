/*
 * Copyright (c) 2020, 2024, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package helper

import (
	"flag"
	"fmt"
	"github.com/go-logr/logr"
	coh "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/controllers"
	"github.com/oracle/coherence-operator/pkg/clients"
	"github.com/oracle/coherence-operator/pkg/operator"
	oprest "github.com/oracle/coherence-operator/pkg/rest"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"net/http"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"time"
)

// TestContext is a context for end-to-end tests
type TestContext struct {
	Config        *rest.Config
	Client        client.Client
	KubeClient    kubernetes.Interface
	Manager       ctrl.Manager
	Logger        logr.Logger
	Context       context.Context
	testEnv       *envtest.Environment
	stop          chan struct{}
	Cancel        context.CancelFunc
	namespaces    []string
	RestServer    oprest.Server
	RestEndpoints map[string]func(w http.ResponseWriter, r *http.Request)
}

func (in *TestContext) Logf(format string, a ...interface{}) {
	in.Logger.Info(fmt.Sprintf(format, a...))
}

func (in *TestContext) CleanupAfterTest(t *testing.T) {
	t.Cleanup(func() {
		if t.Failed() {
			// dump the logs if the test failed
			DumpOperatorLogs(t, *in)
		}
		in.Cleanup()
	})
}

func (in *TestContext) Cleanup() {
	in.Logger.Info("tearing down the test environment")
	ns := GetTestNamespace()
	in.CleanupNamespace(ns)
	clusterNS := GetTestClusterNamespace()
	if clusterNS != ns {
		in.CleanupNamespace(clusterNS)
	}
	clientNS := GetTestClientNamespace()
	in.CleanupNamespace(clientNS)
	for i := range in.namespaces {
		_ = in.cleanAndDeleteNamespace(in.namespaces[i])
	}
	in.namespaces = nil
}

func (in *TestContext) CleanupNamespace(ns string) {
	in.Logger.Info("tearing down the test environment - namespace: " + ns)
	if err := WaitForCoherenceCleanup(*in, ns); err != nil {
		in.Logf("error waiting for cleanup to complete: %+v", err)
	}
	DeletePersistentVolumes(*in, ns)
}

func (in *TestContext) CreateNamespace(ns string) error {
	n := corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ns,
			Namespace: ns,
		},
		Spec: corev1.NamespaceSpec{},
	}
	_, err := in.KubeClient.CoreV1().Namespaces().Create(in.Context, &n, metav1.CreateOptions{})
	if err != nil {
		in.namespaces = append(in.namespaces, ns)
	}
	return err
}

func (in *TestContext) DeleteNamespace(ns string) error {
	for i := range in.namespaces {
		if in.namespaces[i] == ns {
			err := in.cleanAndDeleteNamespace(ns)
			last := len(in.namespaces) - 1
			in.namespaces[i] = in.namespaces[last] // Copy last element to index i.
			in.namespaces[last-1] = ""             // Erase last element (write zero value).
			in.namespaces = in.namespaces[:last]   // Truncate slice.
			return err
		}
	}
	return nil
}

func (in *TestContext) cleanAndDeleteNamespace(ns string) error {
	in.CleanupNamespace(ns)
	return in.KubeClient.CoreV1().Namespaces().Delete(in.Context, ns, metav1.DeleteOptions{})
}

func (in *TestContext) Close() {
	in.Cleanup()
	if in.stop != nil {
		close(in.stop)
	}
	if in.testEnv != nil {
		if err := in.testEnv.Stop(); err != nil {
			in.Logf("error stopping test environment: %+v", err)
		}
	}
}

func (in *TestContext) Start() error {
	var err error

	// Create the REST server
	in.RestServer = oprest.NewServerWithEndpoints(in.KubeClient, in.RestEndpoints)
	if err := in.RestServer.SetupWithManager(in.Manager); err != nil {
		return err
	}

	// Start the manager, which will start the controller and REST server
	in.stop = make(chan struct{})
	go func() {
		err = in.Manager.Start(in.Context)
	}()

	in.Manager.GetCache().WaitForCacheSync(in.Context)
	<-in.RestServer.Running()

	time.Sleep(5 * time.Second)
	return err
}

// NewStartedContext creates a new TestContext starts it.
func NewStartedContext(startController bool, watchNamespaces ...string) (TestContext, error) {
	ctx, err := NewContext(startController, watchNamespaces...)
	if err == nil {
		err = ctx.Start()
	}
	return ctx, err
}

// NewContext creates a new TestContext.
func NewContext(startController bool, watchNamespaces ...string) (TestContext, error) {
	testLogger := zap.New(zap.UseDevMode(true), zap.WriteTo(os.Stdout))
	logf.SetLogger(testLogger)

	// create a dummy command
	Cmd := &cobra.Command{
		Use:   "manager",
		Short: "Start the operator manager",
	}

	// configure viper for the flags and env-vars
	operator.SetupFlags(Cmd, viper.GetViper())
	flagSet := pflag.NewFlagSet("operator", pflag.ContinueOnError)
	flagSet.AddGoFlagSet(flag.CommandLine)
	if err := viper.BindPFlags(flagSet); err != nil {
		return TestContext{}, err
	}

	// We need a real cluster for these tests
	useCluster := true

	testLogger.WithName("test").Info("bootstrapping test environment")
	testEnv := &envtest.Environment{
		UseExistingCluster:       &useCluster,
		AttachControlPlaneOutput: true,
		CRDs:                     []*v1.CustomResourceDefinition{},
	}

	var err error

	err = corev1.AddToScheme(scheme.Scheme)
	if err != nil {
		return TestContext{}, err
	}
	err = coh.AddToScheme(scheme.Scheme)
	if err != nil {
		return TestContext{}, err
	}

	k8sCfg, err := testEnv.Start()
	if err != nil {
		return TestContext{}, err
	}

	cl, err := client.New(k8sCfg, client.Options{Scheme: scheme.Scheme})
	if err != nil {
		return TestContext{}, err
	}

	options := ctrl.Options{
		Scheme: scheme.Scheme,
	}

	if len(watchNamespaces) == 1 {
		// Watch a single namespace
		options.NewCache = func(config *rest.Config, opts cache.Options) (cache.Cache, error) {
			opts.DefaultNamespaces = map[string]cache.Config{
				watchNamespaces[0]: {},
			}
			return cache.New(config, opts)
		}
	} else if len(watchNamespaces) > 1 {
		// Watch a multiple namespaces
		options.NewCache = func(config *rest.Config, opts cache.Options) (cache.Cache, error) {
			nsMap := make(map[string]cache.Config)
			for _, ns := range watchNamespaces {
				nsMap[ns] = cache.Config{}
			}
			opts.DefaultNamespaces = nsMap
			return cache.New(config, opts)
		}
	}

	k8sManager, err := ctrl.NewManager(k8sCfg, options)
	if err != nil {
		return TestContext{}, err
	}

	k8sClient := k8sManager.GetClient()

	cs, err := clients.NewForConfig(k8sCfg)
	if err != nil {
		return TestContext{}, err
	}

	v, err := operator.DetectKubernetesVersion(cs)
	if err != nil {
		return TestContext{}, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	var stop chan struct{}

	if startController {
		// Ensure CRDs exist
		err = coh.EnsureCRDs(ctx, v, scheme.Scheme, cl)
		if err != nil {
			return TestContext{}, err
		}

		// Create the Coherence controller
		err = (&controllers.CoherenceReconciler{
			Client: k8sManager.GetClient(),
			Log:    ctrl.Log.WithName("controllers").WithName("Coherence"),
		}).SetupWithManager(k8sManager, cs)
		if err != nil {
			return TestContext{}, err
		}

		// Create the CoherenceJob controller
		err = (&controllers.CoherenceJobReconciler{
			Client: k8sManager.GetClient(),
			Log:    ctrl.Log.WithName("controllers").WithName("CoherenceJob"),
		}).SetupWithManager(k8sManager, cs)
		if err != nil {
			return TestContext{}, err
		}
	}

	ep := make(map[string]func(w http.ResponseWriter, r *http.Request))

	return TestContext{
		Config:        k8sCfg,
		Client:        k8sClient,
		KubeClient:    cs.KubeClient,
		Manager:       k8sManager,
		Logger:        testLogger.WithName("test"),
		Context:       ctx,
		testEnv:       testEnv,
		stop:          stop,
		Cancel:        cancel,
		RestEndpoints: ep,
	}, nil
}
