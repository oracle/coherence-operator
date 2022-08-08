/*
 * Copyright (c) 2019, 2022, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/oracle/coherence-operator/controllers/webhook"
	"github.com/oracle/coherence-operator/pkg/clients"
	"github.com/oracle/coherence-operator/pkg/operator"
	"github.com/oracle/coherence-operator/pkg/rest"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"k8s.io/apimachinery/pkg/util/version"
	"net/http"
	"os"
	"runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"strings"

	apiruntime "k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	coh "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/controllers"
	// +kubebuilder:scaffold:imports
)

const (
	flagMetricsAddress = "metrics-addr"
	flagHealthAddress  = "health-addr"
	flagLeaderElection = "enable-leader-election"

	// lockName is the name of the lock used for leadership election.
	// This value should not be changed, otherwise a rolling upgrade of the Operator
	// would have two leaders.
	lockName = "ca804aa8.oracle.com"
)

var (
	scheme   = apiruntime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")

	// Cmd is the cobra command to start the manager.
	Cmd = &cobra.Command{
		Use:   "manager",
		Short: "Start the operator manager",
		Long:  "manager starts the manager for this operator, which will in turn create the necessary controller.",
		Run: func(cmd *cobra.Command, args []string) {
			execute()
		},
	}
)

func init() {
	Cmd.Flags().String(flagMetricsAddress, ":8080", "The address the metric endpoint binds to.")
	Cmd.Flags().String(flagHealthAddress, ":8088", "The address the health endpoint binds to.")
	Cmd.Flags().Bool(flagLeaderElection, false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	operator.SetupFlags(Cmd)

	// Add flags registered by imported packages (e.g. glog and controller-runtime)
	flagSet := pflag.NewFlagSet("operator", pflag.ContinueOnError)
	flagSet.AddGoFlagSet(flag.CommandLine)
	if err := viper.BindPFlags(flagSet); err != nil {
		setupLog.Error(err, "binding flags")
		os.Exit(1)
	}

	// Validate the command line flags and environment variables
	if err := operator.ValidateFlags(); err != nil {
		fmt.Println(err.Error())
		_ = Cmd.Help()
		os.Exit(1)
	}

	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(coh.AddToScheme(scheme))
	// +kubebuilder:scaffold:scheme
}

func main() {
	if err := Cmd.Execute(); err != nil {
		logf.Log.WithName("main").Error(err, "Unexpected error while executing command")
	}
}

func execute() {
	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))
	printVersion()

	cfg := ctrl.GetConfigOrDie()
	cs, err := clients.NewForConfig(cfg)
	if err != nil {
		setupLog.Error(err, "unable to create clientset")
		os.Exit(1)
	}

	// create the client here as we use it to install CRDs then inject it into the Manager
	cl, err := client.New(cfg, client.Options{Scheme: scheme})
	if err != nil {
		setupLog.Error(err, "unable to create client")
		os.Exit(1)
	}

	options := ctrl.Options{
		Scheme:                 scheme,
		HealthProbeBindAddress: viper.GetString(flagHealthAddress),
		MetricsBindAddress:     viper.GetString(flagMetricsAddress),
		Port:                   9443,
		LeaderElection:         viper.GetBool(flagLeaderElection),
		LeaderElectionID:       lockName,
	}

	// Determine the Operator scope...
	watchNamespaces := getWatchNamespace()
	switch len(watchNamespaces) {
	case 0:
		// Watching all namespaces
		setupLog.Info("Operator will watch all namespaces")
	case 1:
		// Watch a single namespace
		setupLog.Info("Operator will watch single namespace: " + watchNamespaces[0])
		options.Namespace = watchNamespaces[0]
	default:
		// Watch a multiple namespaces
		setupLog.Info(fmt.Sprintf("Operator will watch multiple namespaces: %v", watchNamespaces))
		options.NewCache = cache.MultiNamespacedCacheBuilder(watchNamespaces)
	}

	mgr, err := ctrl.NewManager(cfg, options)
	if err != nil {
		setupLog.Error(err, "unable to create controller manager")
		os.Exit(1)
	}

	v, err := operator.DetectKubernetesVersion(cs)
	if err != nil {
		setupLog.Error(err, "unable to detect Kubernetes version")
		os.Exit(1)
	}

	ctx := context.TODO()
	initialiseOperator(ctx, v, cl)

	// Set up the Coherence reconciler
	if err = (&controllers.CoherenceReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("Coherence"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Coherence")
		os.Exit(1)
	}

	// We intercept the signal handler here so that we can do clean-up before the Manager stops
	handler := ctrl.SetupSignalHandler()

	// Set-up webhooks if required
	var cr *webhook.CertReconciler
	if operator.ShouldEnableWebhooks() {
		// Set up the webhook certificate reconciler
		cr = &webhook.CertReconciler{
			Clientset: cs,
		}
		if err := cr.SetupWithManager(handler, mgr); err != nil {
			setupLog.Error(err, " unable to create webhook certificate controller", "controller", "Certs")
			os.Exit(1)
		}

		// Set up the webhooks
		if err = (&coh.Coherence{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, " unable to create webhook", "webhook", "Coherence")
			os.Exit(1)
		}
	}

	// Create the REST server
	restServer := rest.NewServer(cs)
	if err := restServer.SetupWithManager(mgr); err != nil {
		setupLog.Error(err, " unable to start REST server")
		os.Exit(1)
	}

	var health healthz.Checker = func(_ *http.Request) error {
		<-restServer.Running()
		return nil
	}

	if err := mgr.AddHealthzCheck("health", health); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("ready", health); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	// +kubebuilder:scaffold:builder

	go func() {
		<-handler.Done()
		if cr != nil {
			cr.Cleanup()
		}
	}()

	setupLog.Info("starting manager")
	if err := mgr.Start(handler); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}

func initialiseOperator(ctx context.Context, v *version.Version, cl client.Client) {
	opLog := ctrl.Log.WithName("operator")

	operator.SetVersion(Version)

	// Ensure that the CRDs exist
	if operator.ShouldInstallCRDs() {
		err := coh.EnsureCRDs(ctx, v, scheme, cl)
		if err != nil {
			opLog.Error(err, "")
			os.Exit(1)
		}
	}
}

// getWatchNamespace returns the Namespace(s) the operator should be watching for changes
func getWatchNamespace() []string {
	// WatchNamespaceEnvVar is the constant for env variable WATCH_NAMESPACE
	// which specifies the Namespace to watch.
	// An empty value means the operator is running with cluster scope.
	var watchNamespaceEnvVar = "WATCH_NAMESPACE"
	var watches []string

	ns, found := os.LookupEnv(watchNamespaceEnvVar)
	if !found || ns == "" || strings.TrimSpace(ns) == "" {
		return watches
	}

	for _, s := range strings.Split(ns, ",") {
		watches = append(watches, strings.TrimSpace(s))
	}
	return watches
}

func printVersion() {
	opLog := ctrl.Log.WithName("operator")
	opLog.Info(fmt.Sprintf("Operator Version: %s", Version))
	opLog.Info(fmt.Sprintf("Operator Build Date: %s", Date))
	opLog.Info(fmt.Sprintf("Operator Built By: %s", Author))
	opLog.Info(fmt.Sprintf("Operator Git Commit: %s", Commit))
	opLog.Info(fmt.Sprintf("Operator Coherence Image: %s", viper.GetString(operator.FlagCoherenceImage)))
	opLog.Info(fmt.Sprintf("Operator Image: %s", viper.GetString(operator.FlagOperatorImage)))
	opLog.Info(fmt.Sprintf("Go Version: %s", runtime.Version()))
	opLog.Info(fmt.Sprintf("Go OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH))
}

var (
	// Version is the Operator version injected by the Go linker at build time.
	Version string
	// Commit is the git commit hash injected by the Go linker at build time.
	Commit string
	// Date is the build timestamp injected by the Go linker at build time.
	Date string
	// Author is the user name of the account at build time
	Author string
)
