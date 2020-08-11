/*
 * Copyright (c) 2019, 2020 Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package main

import (
	"flag"
	"fmt"
	"github.com/oracle/coherence-operator/pkg/flags"
	"github.com/oracle/coherence-operator/pkg/operator"
	"github.com/oracle/coherence-operator/pkg/rest"
	"github.com/spf13/pflag"
	"os"
	"runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
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

var (
	scheme   = apiruntime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(coh.AddToScheme(scheme))
	// +kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool

	pflag.CommandLine.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	pflag.CommandLine.BoolVar(&enableLeaderElection, "enable-leader-election", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")

	// Add flags registered by imported packages (e.g. glog and
	// controller-runtime)
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()

	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))
	printVersion()
	initialiseOperator()

	options := ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: metricsAddr,
		Port:               9443,
		LeaderElection:     enableLeaderElection,
		LeaderElectionID:   "ca804aa8.oracle.com",
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

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), options)
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	if err = (&controllers.CoherenceReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("Coherence"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Coherence")
		os.Exit(1)
	}
	// +kubebuilder:scaffold:builder

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}

func initialiseOperator() {
	log := ctrl.Log.WithName("operator")

	// Get a config to talk to the apiserver
	cfg, err := config.GetConfig()
	if err != nil {
		log.Error(err, "")
		os.Exit(1)
	}

	// Ensure that the CRDs exist
	err = coh.EnsureCRDs(cfg)
	if err != nil {
		log.Error(err, "")
		os.Exit(1)
	}

	opFlags := flags.GetOperatorFlags()

	// Create the REST server
	s, err := rest.EnsureServer(cfg, opFlags)
	if err != nil {
		log.Error(err, "failed to create REST server")
		os.Exit(1)
	}
	// Add the REST server to the Manager so that is is started after the Manager is initialized
	err = s.Start()
	if err != nil {
		log.Error(err, "failed to start the REST server")
		os.Exit(1)
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
	cfg, err := operator.GetOperatorConfig()
	if err != nil {
		panic(err)
	}

	log := ctrl.Log.WithName("operator")
	log.Info(fmt.Sprintf("Operator Version: %s", Version))
	log.Info(fmt.Sprintf("Operator Build Date: %s", Date))
	log.Info(fmt.Sprintf("Operator Git Commit: %s", Commit))
	log.Info(fmt.Sprintf("Operator Coherence Image: %s", cfg.GetDefaultCoherenceImage()))
	log.Info(fmt.Sprintf("Operator Utils Image: %s", cfg.GetDefaultUtilsImage()))
	log.Info(fmt.Sprintf("Go Version: %s", runtime.Version()))
	log.Info(fmt.Sprintf("Go OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH))
}

// ---- Coherence Operator additions ---------------------------------------------------------------

var (
	// BuildInfo is a pipe delimited string of build information injected by the Go linker at build time.
	BuildInfo string
	Version   string
	Commit    string
	Date      string
)

func init() {
	// Use the Go init function to add Operator specific functionality to main
	// Add the Operator flags
	pflag.CommandLine.AddFlagSet(flags.FlagSet())

	if BuildInfo != "" {
		parts := strings.Split(BuildInfo, "|")

		if len(parts) > 0 {
			Version = parts[0]
		}

		if len(parts) > 1 {
			Commit = parts[1]
		}

		if len(parts) > 2 {
			Date = strings.Replace(parts[2], ".", " ", -1)
		}
	}
}
