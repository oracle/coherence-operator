/*
 * Copyright (c) 2020, 2025, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package runner

import (
	"context"
	"fmt"
	coh "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/controllers"
	"github.com/oracle/coherence-operator/pkg/clients"
	"github.com/oracle/coherence-operator/pkg/operator"
	"github.com/oracle/coherence-operator/pkg/rest"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	apiruntime "k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/version"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	rest2 "k8s.io/client-go/rest"
	"k8s.io/utils/ptr"
	"net/http"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/config"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
	// +kubebuilder:scaffold:imports
)

const (
	// CommandOperator is the argument to launch the Operator manager.
	CommandOperator = "operator"

	// lockName is the name of the lock used for leadership election.
	// This value should not be changed, otherwise a rolling upgrade of the Operator
	// would have two leaders.
	lockName = "ca804aa8.oracle.com"
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

// operatorCommand runs the Coherence Operator manager
func operatorCommand(v *viper.Viper) *cobra.Command {
	cmd := &cobra.Command{
		Use:   CommandOperator,
		Short: "Run the Coherence Operator",
		Long:  "Run the Coherence Operator",
		RunE: func(cmd *cobra.Command, args []string) error {
			return execute()
		},
	}

	operator.SetupOperatorManagerFlags(cmd, v)

	return cmd
}

func execute() error {
	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))

	setupLog.Info(fmt.Sprintf("Operator Coherence Image: %s", operator.GetDefaultCoherenceImage()))
	setupLog.Info(fmt.Sprintf("Operator Image: %s", operator.GetDefaultOperatorImage()))

	cfg := ctrl.GetConfigOrDie()
	cs, err := clients.NewForConfig(cfg)
	if err != nil {
		return errors.Wrap(err, "unable to create client set")
	}

	// create the client here as we use it to install CRDs then inject it into the Manager
	setupLog.Info("Creating Kubernetes client", "Host", cfg.Host)
	cl, err := client.New(cfg, client.Options{Scheme: scheme})
	if err != nil {
		return errors.Wrap(err, "unable to create client")
	}

	vpr := operator.GetViper()
	f := vpr.GetBool(operator.FlagDryRun)

	options := ctrl.Options{
		Scheme:                 scheme,
		HealthProbeBindAddress: viper.GetString(operator.FlagHealthAddress),
		Metrics: metricsserver.Options{
			BindAddress: viper.GetString(operator.FlagMetricsAddress),
		},
		LeaderElection:   viper.GetBool(operator.FlagLeaderElection),
		LeaderElectionID: lockName,
		Controller: config.Controller{
			SkipNameValidation: ptr.To(f),
		},
	}

	// Determine the Operator scope...
	watchNamespaces := operator.GetWatchNamespace()
	switch len(watchNamespaces) {
	case 0:
		// Watching all namespaces
		setupLog.Info("Operator will watch all namespaces")
	case 1:
		// Watch a single namespace
		setupLog.Info("Operator will watch single namespace: " + watchNamespaces[0])
		options.NewCache = func(config *rest2.Config, opts cache.Options) (cache.Cache, error) {
			opts.DefaultNamespaces = map[string]cache.Config{
				watchNamespaces[0]: {},
			}
			return cache.New(config, opts)
		}
	default:
		// Watch a multiple namespaces
		setupLog.Info(fmt.Sprintf("Operator will watch multiple namespaces: %v", watchNamespaces))
		options.NewCache = func(config *rest2.Config, opts cache.Options) (cache.Cache, error) {
			nsMap := make(map[string]cache.Config)
			for _, ns := range watchNamespaces {
				nsMap[ns] = cache.Config{}
			}
			opts.DefaultNamespaces = nsMap
			return cache.New(config, opts)
		}
	}

	mgr, err := manager.New(cfg, options)
	if err != nil {
		return errors.Wrap(err, "unable to create controller manager")
	}

	v, err := operator.DetectKubernetesVersion(cs)
	if err != nil {
		return errors.Wrap(err, "unable to detect Kubernetes version")
	}

	ctx := context.TODO()
	initialiseOperator(ctx, v, cl)

	// Set up the Coherence reconciler
	if err = (&controllers.CoherenceReconciler{
		Client:    mgr.GetClient(),
		ClientSet: cs,
		Log:       ctrl.Log.WithName("controllers").WithName("Coherence"),
		Scheme:    mgr.GetScheme(),
	}).SetupWithManager(mgr, cs); err != nil {
		return errors.Wrap(err, "unable to create Coherence controller")
	}

	// Set up the CoherenceJob reconciler
	if operator.ShouldInstallJobCRD() {
		if err = (&controllers.CoherenceJobReconciler{
			Client:    mgr.GetClient(),
			ClientSet: cs,
			Log:       ctrl.Log.WithName("controllers").WithName("CoherenceJob"),
			Scheme:    mgr.GetScheme(),
		}).SetupWithManager(mgr, cs); err != nil {
			return errors.Wrap(err, "unable to create CoherenceJob controller")
		}
	}

	dryRun := operator.IsDryRun()
	if !dryRun {
		// We intercept the signal handler here so that we can do clean-up before the Manager stops
		handler := ctrl.SetupSignalHandler()

		// Create the REST server
		restServer := rest.NewServer(cs.KubeClient)
		if err := restServer.SetupWithManager(mgr); err != nil {
			return errors.Wrap(err, " unable to start REST server")
		}

		var health healthz.Checker = func(_ *http.Request) error {
			<-restServer.Running()
			return nil
		}

		if err := mgr.AddHealthzCheck("health", health); err != nil {
			return errors.Wrap(err, "unable to set up health check")
		}
		if err := mgr.AddReadyzCheck("ready", health); err != nil {
			return errors.Wrap(err, "unable to set up ready check")
		}

		// +kubebuilder:scaffold:builder

		setupLog.Info("starting manager")
		if err := mgr.Start(handler); err != nil {
			setupLog.Error(err, "problem running manager")
			os.Exit(1)
		}
	}

	return nil
}

func initialiseOperator(ctx context.Context, v *version.Version, cl client.Client) {
	opLog := ctrl.Log.WithName("operator")

	// Ensure that the CRDs exist
	if operator.ShouldInstallCRDs() {
		err := coh.EnsureCRDs(ctx, v, scheme, cl)
		if err != nil {
			opLog.Error(err, "")
			os.Exit(1)
		}
	}
}
