/*
 * Copyright (c) 2022, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package runner

import (
	"context"
	"fmt"
	coh "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/controllers"
	"github.com/oracle/coherence-operator/controllers/webhook"
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
	"net/http"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
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
func operatorCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   CommandOperator,
		Short: "Run the Coherence Operator",
		Long:  "Run the Coherence Operator",
		RunE: func(cmd *cobra.Command, args []string) error {
			return execute()
		},
	}

	operator.SetupOperatorManagerFlags(cmd)

	return cmd
}

func execute() error {
	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))

	setupLog.Info(fmt.Sprintf("Operator Coherence Image: %s", viper.GetString(operator.FlagCoherenceImage)))
	setupLog.Info(fmt.Sprintf("Operator Image: %s", viper.GetString(operator.FlagOperatorImage)))

	cfg := ctrl.GetConfigOrDie()
	cs, err := clients.NewForConfig(cfg)
	if err != nil {
		return errors.Wrap(err, "unable to create clientset")
	}

	// create the client here as we use it to install CRDs then inject it into the Manager
	cl, err := client.New(cfg, client.Options{Scheme: scheme})
	if err != nil {
		return errors.Wrap(err, "unable to create client")
	}

	options := ctrl.Options{
		Scheme:                 scheme,
		HealthProbeBindAddress: viper.GetString(operator.FlagHealthAddress),
		MetricsBindAddress:     viper.GetString(operator.FlagMetricsAddress),
		Port:                   9443,
		LeaderElection:         viper.GetBool(operator.FlagLeaderElection),
		LeaderElectionID:       lockName,
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
		options.Namespace = watchNamespaces[0]
	default:
		// Watch a multiple namespaces
		setupLog.Info(fmt.Sprintf("Operator will watch multiple namespaces: %v", watchNamespaces))
		options.NewCache = cache.MultiNamespacedCacheBuilder(watchNamespaces)
	}

	mgr, err := ctrl.NewManager(cfg, options)
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
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("Coherence"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		return errors.Wrap(err, "unable to create controller")
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
			return errors.Wrap(err, " unable to create webhook certificate controller")
		}

		// Set up the webhooks
		if err = (&coh.Coherence{}).SetupWebhookWithManager(mgr); err != nil {
			return errors.Wrap(err, " unable to create webhook")
		}
	}

	// Create the REST server
	restServer := rest.NewServer(cs)
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
