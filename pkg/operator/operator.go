/*
 * Copyright (c) 2019, 2023, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

// Package operator package contains types and functions used directly by the Operator main
package operator

import (
	"flag"
	"fmt"
	"github.com/oracle/coherence-operator/pkg/clients"
	"github.com/oracle/coherence-operator/pkg/data"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/version"
	"os"
	"path/filepath"
	ctrl "sigs.k8s.io/controller-runtime"
	"strings"
	"time"
)

const (
	DefaultRestHost       = "0.0.0.0"
	DefaultRestPort int32 = 8000

	// DefaultCertValidity makes new certificates default to a one-year expiration
	DefaultCertValidity = 365 * 24 * time.Hour
	// DefaultRotateBefore defines how long before expiration a certificate
	// should be re-issued
	DefaultRotateBefore = 24 * time.Hour

	// CertFileName is used for Certificates inside a secret
	CertFileName = "tls.crt"
	// KeyFileName is used for Private Keys inside a secret
	KeyFileName = "tls.key"

	CertTypeSelfSigned    = "self-signed"
	CertTypeCertManager   = "cert-manager"
	CertTypeManual        = "manual"
	CertManagerIssuerName = "coherence-webhook-server-issuer"

	DefaultMutatingWebhookName   = "coherence-operator-mutating-webhook-configuration"
	DefaultValidatingWebhookName = "coherence-operator-validating-webhook-configuration"

	FlagCACertRotateBefore    = "ca-cert-rotate-before"
	FlagCACertValidity        = "ca-cert-validity"
	FlagCertType              = "cert-type"
	FlagCertIssuer            = "cert-issuer"
	FlagCoherenceImage        = "coherence-image"
	FlagDevMode               = "coherence-dev-mode"
	FlagCRD                   = "install-crd"
	FlagEnableWebhook         = "enable-webhook"
	FlagHealthAddress         = "health-addr"
	FlagLeaderElection        = "enable-leader-election"
	FlagMetricsAddress        = "metrics-addr"
	FlagMutatingWebhookName   = "mutating-webhook-name"
	FlagOperatorNamespace     = "operator-namespace"
	FlagRackLabel             = "rack-label"
	FlagRestHost              = "rest-host"
	FlagRestPort              = "rest-port"
	FlagServiceName           = "service-name"
	FlagServicePort           = "service-port"
	FlagSiteLabel             = "site-label"
	FlagSkipServiceSuspend    = "skip-service-suspend"
	FlagOperatorImage         = "operator-image"
	FlagValidatingWebhookName = "validating-webhook-name"
	FlagWebhookCertDir        = "webhook-cert-dir"
	FlagWebhookSecret         = "webhook-secret"
	FlagWebhookService        = "webhook-service"

	// EnvVarWatchNamespace is the environment variable to use to set the watch namespace(s)
	EnvVarWatchNamespace = "WATCH_NAMESPACE"
	// EnvVarCoherenceImage is the environment variable to use to set the default Coherence image
	EnvVarCoherenceImage = "COHERENCE_IMAGE"

	// OCI Node Labels

	// LabelOciNodeFaultDomain is the OCI Node label for the fault domain.
	LabelOciNodeFaultDomain = "oci.oraclecloud.com/fault-domain"

	// LabelHostName is the Node label for the Node's hostname.
	LabelHostName = "kubernetes.io/hostname"
)

var setupLog = ctrl.Log.WithName("setup")

var (
	operatorVersion   = "999.0.0"
	DefaultSiteLabels = []string{corev1.LabelTopologyZone, corev1.LabelFailureDomainBetaZone}
	DefaultRackLabels = []string{LabelOciNodeFaultDomain, corev1.LabelTopologyZone, corev1.LabelFailureDomainBetaZone}
)

func SetupOperatorManagerFlags(cmd *cobra.Command) {
	flags := cmd.Flags()
	flags.String(FlagMetricsAddress, ":8080", "The address the metric endpoint binds to.")
	flags.String(FlagHealthAddress, ":8088", "The address the health endpoint binds to.")
	flags.Bool(FlagLeaderElection, false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")

	SetupFlags(cmd)

	// Add flags registered by imported packages (e.g. glog and controller-runtime)
	flagSet := pflag.NewFlagSet("operator", pflag.ContinueOnError)
	flagSet.AddGoFlagSet(flag.CommandLine)
	if err := viper.BindPFlags(flagSet); err != nil {
		setupLog.Error(err, "binding flags")
		os.Exit(1)
	}

	// Validate the command line flags and environment variables
	if err := ValidateFlags(); err != nil {
		fmt.Println(err.Error())
		_ = cmd.Help()
		os.Exit(1)
	}

}

func SetupFlags(cmd *cobra.Command) {
	f, err := data.Assets.Open("assets/config.json")
	if err != nil {
		setupLog.Error(err, "finding config.json asset")
		os.Exit(1)
	}

	viper.SetConfigType("json")
	if err := viper.ReadConfig(f); err != nil {
		setupLog.Error(err, "reading configuration file")
		os.Exit(1)
	} else {
		_ = f.Close()
	}

	cmd.Flags().Duration(
		FlagCACertRotateBefore,
		DefaultRotateBefore,
		"Duration representing how long before expiration CA certificates should be reissued",
	)
	cmd.Flags().Duration(
		FlagCACertValidity,
		DefaultCertValidity,
		"Duration representing how long before a newly created CA cert expires",
	)
	cmd.Flags().String(
		FlagCertType,
		CertTypeSelfSigned,
		fmt.Sprintf("The type of certificate management used for webhook certificates. "+
			"Valid options are %v", []string{CertTypeSelfSigned, CertTypeCertManager, CertTypeManual}),
	)
	cmd.Flags().String(
		FlagCertIssuer,
		CertManagerIssuerName,
		fmt.Sprintf("The name of an existing cert-manager issuer (in the same namespace) used for webhook certificates. "+
			"Valid options are %v", []string{CertTypeSelfSigned, CertTypeCertManager, CertTypeManual}),
	)
	cmd.Flags().String(
		FlagCoherenceImage,
		"",
		"The default Coherence image to use if none is specified.",
	)
	cmd.Flags().Bool(
		FlagDevMode,
		false,
		"Run in dev mode. This should only be used during testing outside of a k8s cluster",
	)
	cmd.Flags().Bool(
		FlagCRD,
		true,
		"Enables automatic installation/update of Coherence CRDs",
	)
	cmd.Flags().Bool(
		FlagEnableWebhook,
		true,
		"Enables the defaulting and validating web-hooks",
	)
	cmd.Flags().String(
		FlagMutatingWebhookName,
		DefaultMutatingWebhookName,
		"Name of the Kubernetes ValidatingWebhookConfiguration resource. Only used when enable-webhook is true.",
	)
	cmd.Flags().String(
		FlagOperatorNamespace,
		"operator-test",
		"The K8s namespace the operator is running in",
	)
	cmd.Flags().StringSlice(
		FlagRackLabel,
		DefaultRackLabels,
		"The node label to use when obtaining a value for a Pod's Coherence rack.",
	)
	cmd.Flags().String(
		FlagRestHost,
		DefaultRestHost,
		"The address that the REST server will bind to",
	)
	cmd.Flags().Int32(
		FlagRestPort,
		DefaultRestPort,
		"The port that the REST server will bind to",
	)
	cmd.Flags().String(
		FlagServiceName,
		"",
		"The service name that operator clients use as the host name to make REST calls back to the operator.",
	)
	cmd.Flags().Int32(
		FlagServicePort,
		-1,
		"The service port that operator clients use in the host name to make REST calls back to the operator. "+
			"If not set defaults to the same as the REST port",
	)
	cmd.Flags().StringSlice(
		FlagSiteLabel,
		DefaultSiteLabels,
		"The node label to use when obtaining a value for a Pod's Coherence site.",
	)
	cmd.Flags().Bool(
		FlagSkipServiceSuspend,
		false,
		"Suspend Coherence services on a cluster prior to shutdown or scaling to zero. "+
			"This option is rarely set to false outside of testing.",
	)
	cmd.Flags().String(
		FlagOperatorImage,
		"",
		"The default Coherence Operator image to use if none is specified.",
	)
	cmd.Flags().String(
		FlagValidatingWebhookName,
		DefaultValidatingWebhookName,
		"Name of the Kubernetes ValidatingWebhookConfiguration resource. Only used when enable-webhook is true.",
	)
	cmd.Flags().String(
		FlagWebhookCertDir,
		filepath.Join(os.TempDir(), "k8s-webhook-server", "serving-certs"),
		"The name of the directory containing the webhook server key and certificate",
	)
	cmd.Flags().String(
		FlagWebhookSecret,
		"coherence-webhook-server-cert",
		"K8s secret to be used for webhook certificates",
	)
	cmd.Flags().String(
		FlagWebhookService,
		"webhook-service",
		"The K8s service used for the webhook",
	)

	// enable using dashed notation in flags and underscores in env
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		setupLog.Error(err, "binding flags")
		os.Exit(1)
	}

	viper.AutomaticEnv()
}

func ValidateFlags() error {
	certValidity := viper.GetDuration(FlagCACertValidity)
	certRotateBefore := viper.GetDuration(FlagCACertRotateBefore)
	if certRotateBefore > certValidity {
		return fmt.Errorf("%s must be larger than %s", FlagCACertValidity, FlagCACertRotateBefore)
	}

	certType := viper.GetString(FlagCertType)
	if certType != CertTypeSelfSigned && certType != CertTypeCertManager && certType != CertTypeManual {
		return fmt.Errorf("%s parameter is invalid", FlagCertType)
	}

	return nil
}

func IsDevMode() bool {
	return viper.GetBool(FlagDevMode)
}

func GetDefaultCoherenceImage() string {
	return viper.GetString(FlagCoherenceImage)
}

func GetDefaultOperatorImage() string {
	return viper.GetString(FlagOperatorImage)
}

func GetRestHost() string {
	return viper.GetString(FlagRestHost)
}

func GetRestPort() int32 {
	return viper.GetInt32(FlagRestPort)
}

func GetRestServiceName() string {
	s := viper.GetString(FlagServiceName)
	if s != "" {
		ns := GetNamespace()
		return s + "." + ns + ".svc"
	}
	return ""
}

func GetRestServicePort() int32 {
	return viper.GetInt32(FlagServicePort)
}
func GetSiteLabel() []string {
	return viper.GetStringSlice(FlagSiteLabel)
}

func GetRackLabel() []string {
	return viper.GetStringSlice(FlagRackLabel)
}

func ShouldInstallCRDs() bool {
	return viper.GetBool(FlagCRD)
}

func ShouldEnableWebhooks() bool {
	return viper.GetBool(FlagEnableWebhook)
}

func ShouldUseSelfSignedCerts() bool {
	return viper.GetString(FlagCertType) == CertTypeSelfSigned
}

func ShouldUseCertManager() bool {
	return viper.GetString(FlagCertType) == CertTypeCertManager
}

func GetNamespace() string {
	return viper.GetString(FlagOperatorNamespace)
}

func GetWebhookCertDir() string {
	return viper.GetString(FlagWebhookCertDir)
}

func GetCACertRotateBefore() time.Duration {
	return viper.GetDuration(FlagCACertRotateBefore)
}

func GetWebhookServiceDNSNames() []string {
	var dns []string
	s := viper.GetString(FlagWebhookService)
	if IsDevMode() {
		dns = []string{s}
	} else {
		ns := GetNamespace()
		return []string{
			fmt.Sprintf("%s.%s", s, ns),
			fmt.Sprintf("%s.%s.svc", s, ns),
			fmt.Sprintf("%s.%s.svc.cluster.local", s, ns),
		}
	}
	return dns
}

func DetectKubernetesVersion(cs clients.ClientSet) (*version.Version, error) {
	sv, err := cs.DiscoveryClient.ServerVersion()
	if err != nil {
		return nil, err
	}
	return version.ParseSemantic(sv.GitVersion)
}

// GetVersion returns the Operator version.
// The Operator version is injected at compile time.
// In development environments, for example running in an IDE where the version has not been injected
// the version 999.0.0 will be returned
func GetVersion() string {
	if operatorVersion == "" {
		return "999.0.0"
	}
	return operatorVersion
}

func SetVersion(v string) {
	operatorVersion = v
}

// GetWatchNamespace returns the Namespace(s) the operator should be watching for changes
func GetWatchNamespace() []string {
	// WatchNamespaceEnvVar is the constant for env variable WATCH_NAMESPACE
	// which specifies the Namespace to watch.
	// An empty value means the operator is running with cluster scope.
	var watches []string

	ns, found := os.LookupEnv(EnvVarWatchNamespace)
	if !found || ns == "" || strings.TrimSpace(ns) == "" {
		return watches
	}

	for _, s := range strings.Split(ns, ",") {
		watches = append(watches, strings.TrimSpace(s))
	}
	return watches
}
