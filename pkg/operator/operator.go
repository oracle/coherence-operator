/*
 * Copyright (c) 2019, 2020, 2021 Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

// The operator package contains types and functions used directly by the Operator main
package operator

import (
	"fmt"
	"github.com/oracle/coherence-operator/pkg/data"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	ctrl "sigs.k8s.io/controller-runtime"
	"strings"
	"time"
)

const (
	DefaultRestHost        = "0.0.0.0"
	DefaultRestPort  int32 = 8000

	// DefaultCertValidity makes new certificates default to a 1 year expiration
	DefaultCertValidity = 365 * 24 * time.Hour
	// DefaultRotateBefore defines how long before expiration a certificate
	// should be re-issued
	DefaultRotateBefore = 24 * time.Hour

	// CertFileName is used for Certificates inside a secret
	CertFileName = "tls.crt"
	// KeyFileName is used for Private Keys inside a secret
	KeyFileName = "tls.key"

	CertTypeSelfSigned  = "self-signed"
	CertTypeCertManager = "cert-manager"
	CertTypeManual      = "manual"

	DefaultMutatingWebhookName   = "coherence-operator-mutating-webhook-configuration"
	DefaultValidatingWebhookName = "coherence-operator-validating-webhook-configuration"

	FlagCACertRotateBefore    = "ca-cert-rotate-before"
	FlagCACertValidity        = "ca-cert-validity"
	FlagCertType              = "cert-type"
	FlagCoherenceImage        = "coherence-image"
	FlagDevMode               = "coherence-dev-mode"
	FlagEnableWebhook         = "enable-webhook"
	FlagMutatingWebhookName   = "mutating-webhook-name"
	FlagOperatorNamespace     = "operator-namespace"
	FlagRackLabel             = "rack-label"
	FlagRestHost              = "rest-host"
	FlagRestPort              = "rest-port"
	FlagServiceName           = "service-name"
	FlagServicePort           = "service-port"
	FlagSiteLabel             = "site-label"
	FlagSkipServiceSuspend    = "skip-service-suspend"
	FlagUtilsImage            = "utils-image"
	FlagValidatingWebhookName = "validating-webhook-name"
	FlagWebhookCertDir        = "webhook-cert-dir"
	FlagWebhookSecret         = "webhook-secret"
	FlagWebhookService        = "webhook-service"
)

var setupLog = ctrl.Log.WithName("setup")

var (
	DefaultSiteLabel = []string{"topology.kubernetes.io/zone", "failure-domain.beta.kubernetes.io/zone"}
	DefaultRackLabel = []string{"topology.kubernetes.io/region", "failure-domain.beta.kubernetes.io/region", "topology.kubernetes.io/zone", "failure-domain.beta.kubernetes.io/zone"}
)

func SetupFlags(cmd *cobra.Command) {
	f, err := data.Assets.Open("config.json")
	if err != nil {
		setupLog.Error(err, "finding data.json asset")
		os.Exit(1)
	}
	defer f.Close()

	viper.SetConfigType("json")
	if err := viper.ReadConfig(f); err != nil {
		setupLog.Error(err, "reading configuration file")
		os.Exit(1)
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
		DefaultRackLabel,
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
		DefaultSiteLabel,
		"The node label to use when obtaining a value for a Pod's Coherence site.",
	)
	cmd.Flags().Bool(
		FlagSkipServiceSuspend,
		false,
		"Suspend Coherence services on a cluster prior to shutdown or scaling to zero. "+
			"This option is rarely set to false outside of testing.",
	)
	cmd.Flags().String(
		FlagUtilsImage,
		"",
		"The default Coherence Operator utils image to use if none is specified.",
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

func GetDefaultUtilsImage() string {
	return viper.GetString(FlagUtilsImage)
}

func GetRestHost() string {
	return viper.GetString(FlagRestHost)
}

func GetRestPort() int32 {
	return viper.GetInt32(FlagRestPort)
}

func GetRestServiceName() string {
	return viper.GetString(FlagServiceName)
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
