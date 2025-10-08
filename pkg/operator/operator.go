/*
 * Copyright (c) 2020, 2025, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

// Package operator package contains types and functions used directly by the Operator main
package operator

import (
	"crypto/tls"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-logr/logr"
	"github.com/oracle/coherence-operator/pkg/clients"
	"github.com/oracle/coherence-operator/pkg/data"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/version"
	ctrl "sigs.k8s.io/controller-runtime"
)

const (
	DefaultRestHost       = "0.0.0.0"
	DefaultRestPort int32 = 8000

	DefaultMutatingWebhookName   = "coherence-operator-mutating-webhook-configuration"
	DefaultValidatingWebhookName = "coherence-operator-validating-webhook-configuration"

	FlagCoherenceImage         = "coherence-image"
	FlagCRD                    = "install-crd"
	FlagJobCRD                 = "install-job-crd"
	FlagEnableCoherenceJobs    = "enable-jobs"
	FlagDevMode                = "coherence-dev-mode"
	FlagCipherDenyList         = "cipher-deny-list"
	FlagCipherAllowList        = "cipher-allow-list"
	FlagConfig                 = "config"
	FlagConfigType             = "config-type"
	FlagDryRun                 = "dry-run"
	FlagEnableWebhook          = "enable-webhook"
	FlagEnableHttp2            = "enable-http2"
	FlagGlobalAnnotation       = "global-annotation"
	FlagGlobalLabel            = "global-label"
	FlagHealthAddress          = "health-addr"
	FlagLeaderElection         = "enable-leader-election"
	FlagLeaderElectionDuration = "leader-election-duration"
	FlagLeaderElectionRenew    = "leader-election-renew-timeout"
	FlagMetricsAddress         = "metrics-addr"
	FlagOperatorNamespace      = "operator-namespace"
	FlagNodeLookupEnabled      = "node-lookup-enabled"
	FlagRackLabel              = "rack-label"
	FlagRestHost               = "rest-host"
	FlagRestPort               = "rest-port"
	FlagSecureMetrics          = "metrics-secure"
	FlagServiceName            = "service-name"
	FlagServicePort            = "service-port"
	FlagSiteLabel              = "site-label"
	FlagSkipServiceSuspend     = "skip-service-suspend"
	FlagOperatorImage          = "operator-image"
	FlagEnvVar                 = "env"
	FlagJvmArg                 = "jvm"

	// EnvVarWatchNamespace is the environment variable to use to set the watch namespace(s)
	EnvVarWatchNamespace = "WATCH_NAMESPACE"
	// EnvVarCoherenceImage is the environment variable to use to set the default Coherence image
	EnvVarCoherenceImage = "COHERENCE_IMAGE"

	// OCI Node Labels

	// LabelOciNodeFaultDomain is the OCI Node label for the fault domain.
	LabelOciNodeFaultDomain = "oci.oraclecloud.com/fault-domain"
	// LabelTopologySubZone is the k8s topology label for sub-zone.
	LabelTopologySubZone = "topology.kubernetes.io/subzone"

	// LabelHostName is the Node label for the Node's hostname.
	LabelHostName = "kubernetes.io/hostname"

	// LabelTestHostName is a label applied to Pods to set a testing host name
	LabelTestHostName = "coherence.oracle.com/test_hostname"
	// LabelTestHealthPort is a label applied to Pods to set a testing health check port
	LabelTestHealthPort = "coherence.oracle.com/test_health_port"
)

var setupLog = ctrl.Log.WithName("setup")

var currentViper *viper.Viper

var (
	operatorVersion   = "999.0.0"
	DefaultSiteLabels = []string{corev1.LabelTopologyZone, corev1.LabelFailureDomainBetaZone}
	DefaultRackLabels = []string{LabelTopologySubZone, LabelOciNodeFaultDomain, corev1.LabelTopologyZone, corev1.LabelFailureDomainBetaZone}
)

func SetupOperatorManagerFlags(cmd *cobra.Command, v *viper.Viper) {
	flags := cmd.Flags()
	flags.String(FlagMetricsAddress, ":8080", "The address the metric endpoint binds to.")
	flags.String(FlagHealthAddress, ":8088", "The address the health endpoint binds to.")
	flags.Bool(FlagLeaderElection, false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")

	SetupFlags(cmd, v)

	// Add flags registered by imported packages (e.g. glog and controller-runtime)
	flagSet := pflag.NewFlagSet("operator", pflag.ContinueOnError)
	flagSet.AddGoFlagSet(flag.CommandLine)
}

func SetupFlags(cmd *cobra.Command, v *viper.Viper) {
	f, err := data.Assets.Open("assets/config.json")
	if err != nil {
		setupLog.Error(err, "finding config.json asset")
		os.Exit(1)
	}

	v.SetConfigType("json")
	if err := v.ReadConfig(f); err != nil {
		setupLog.Error(err, "reading configuration file")
		os.Exit(1)
	}
	_ = f.Close()

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
		"This flag is deprecated and no longer has any function",
	)
	cmd.Flags().Bool(
		FlagJobCRD,
		true,
		"Enables CoherenceJob support, this flag is deprecated use --"+FlagEnableCoherenceJobs,
	)
	cmd.Flags().Bool(
		FlagEnableCoherenceJobs,
		true,
		"Enables CoherenceJob support",
	)
	cmd.Flags().Bool(
		FlagEnableHttp2,
		false,
		"If set, HTTP/2 will be enabled for the metrics and REST servers",
	)
	cmd.Flags().Bool(
		FlagEnableWebhook,
		false,
		"This flag is here for backward compatibility but is ignored",
	)
	cmd.Flags().Bool(
		FlagNodeLookupEnabled,
		true,
		"The Operator is allowed to lookup information about kubernetes nodes",
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
	cmd.Flags().Bool(
		FlagSecureMetrics,
		false,
		"FlagSecureMetrics",
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
	cmd.Flags().StringArray(
		FlagGlobalAnnotation,
		nil,
		"An annotation to apply to all resources managed by the Operator (can be used multiple times)")
	cmd.Flags().StringArray(
		FlagGlobalLabel,
		nil,
		"A label to apply to all resources managed by the Operator (can be used multiple times)")
	cmd.Flags().StringArray(
		FlagCipherDenyList,
		nil,
		"A list of TLS cipher names to be disabled")
	cmd.Flags().StringArray(
		FlagCipherAllowList,
		nil,
		"A list of TLS cipher names to be enabled (if a cipher appears in this list and the deny list it will be disabled)")
	cmd.Flags().Duration(
		FlagLeaderElectionDuration,
		time.Second*30,
		"The value the Operator uses for the leadership lease duration. "+
			"Setting this value too low can cause Pod restarts as the leader may lose leadership. "+
			"If the value entered is less than 10s, then 10s will be used")
	cmd.Flags().Duration(
		FlagLeaderElectionRenew,
		time.Second*20,
		"The duration the Operator uses for the leadership lease renewal timeout. "+
			"If the value entered is less than 10s, then 10s will be used")

	// enable using dashed notation in flags and underscores in env
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	if err := v.BindPFlags(cmd.Flags()); err != nil {
		setupLog.Error(err, "binding flags")
		os.Exit(1)
	}

	v.AutomaticEnv()
}

func SetViper(v *viper.Viper) {
	currentViper = v
}

func GetViper() *viper.Viper {
	if currentViper == nil {
		return viper.GetViper()
	}
	return currentViper
}

func GetDefaultCoherenceImage() string {
	return GetViper().GetString(FlagCoherenceImage)
}

func GetDefaultOperatorImage() string {
	return GetViper().GetString(FlagOperatorImage)
}

func GetRestHost() string {
	return GetViper().GetString(FlagRestHost)
}

func GetRestPort() int32 {
	return GetViper().GetInt32(FlagRestPort)
}

func GetRestServiceName() string {
	s := GetViper().GetString(FlagServiceName)
	if s != "" {
		ns := GetNamespace()
		return s + "." + ns + ".svc"
	}
	return ""
}

func GetRestServicePort() int32 {
	return GetViper().GetInt32(FlagServicePort)
}
func GetSiteLabel() []string {
	return GetViper().GetStringSlice(FlagSiteLabel)
}

func GetRackLabel() []string {
	return GetViper().GetStringSlice(FlagRackLabel)
}

func ShouldSupportCoherenceJob() bool {
	v := GetViper()
	return v.GetBool(FlagEnableCoherenceJobs) || v.GetBool(FlagJobCRD)
}

func IsDryRun() bool {
	return GetViper().GetBool(FlagDryRun)
}

func GetNamespace() string {
	return GetViper().GetString(FlagOperatorNamespace)
}

func IsNodeLookupEnabled() bool {
	return GetViper().GetBool(FlagNodeLookupEnabled)
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

func GetGlobalAnnotationsNoError() map[string]string {
	m, _ := GetGlobalAnnotations(GetViper())
	return m
}

func GetGlobalAnnotations(v *viper.Viper) (map[string]string, error) {
	args := v.GetStringSlice(FlagGlobalAnnotation)
	return stringSliceToMap(args, FlagGlobalAnnotation)
}

func GetGlobalLabelsNoError() map[string]string {
	m, _ := GetGlobalLabels(GetViper())
	return m
}

func GetGlobalLabels(v *viper.Viper) (map[string]string, error) {
	args := v.GetStringSlice(FlagGlobalLabel)
	return stringSliceToMap(args, FlagGlobalLabel)
}

// GetTlsCipherDenyList returns the names of the TLS cipher suites to be disabled.
func GetTlsCipherDenyList(v *viper.Viper) []string {
	return createCipherList(v.GetStringSlice(FlagCipherDenyList))
}

// GetTlsCipherAllowList returns the names of the TLS cipher suites to be enabled.
func GetTlsCipherAllowList(v *viper.Viper) []string {
	return createCipherList(v.GetStringSlice(FlagCipherAllowList))
}

func createCipherList(l []string) []string {
	var list []string
	for i := 0; i < len(l); i++ {
		for _, s := range strings.Split(l[i], ",") {
			list = append(list, strings.ToUpper(s))
		}
	}
	return list
}

// DefaultCipherDenyList returns the default list of ciphers disabled by Oracle's policies.
func DefaultCipherDenyList() []uint16 {
	return []uint16{tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA, tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA}
}

// NewCipherSuiteConfig returns a function that will configure the allowed cipher suites for a TLS configuration.
func NewCipherSuiteConfig(v *viper.Viper, log logr.Logger) (func(c *tls.Config), error) {
	var allowed []uint16
	all := make(map[uint16]*tls.CipherSuite)
	allByName := make(map[string]*tls.CipherSuite)

	for _, cs := range tls.CipherSuites() {
		all[cs.ID] = cs
		allByName[cs.Name] = cs
		allowed = append(allowed, cs.ID)
	}

	for _, cs := range tls.InsecureCipherSuites() {
		all[cs.ID] = cs
		allByName[cs.Name] = cs
	}

	allowed = RemoveAllFromUInt16Array(allowed, DefaultCipherDenyList()...)

	denyList := GetTlsCipherDenyList(v)
	allDenied := false
	allowList := GetTlsCipherAllowList(v)

	for _, name := range denyList {
		if strings.ToUpper(name) == "ALL" {
			// remove all allowed suites
			allowed = make([]uint16, 0)
			allDenied = true
			break
		}
	}

	if allDenied && len(allowList) == 0 {
		return func(c *tls.Config) {},
			fmt.Errorf("the --%s command line flag has denied all cipher suites but no allowed suites have been specified using the %s flag ", FlagCipherDenyList, FlagCipherAllowList)
	}

	for _, name := range allowList {
		if cs, found := allByName[name]; found {
			allowed = append(allowed, cs.ID)
		} else {
			return func(c *tls.Config) {},
				fmt.Errorf("invalid --%s command line flag, %s is not a recognised cipher name (see https://pkg.go.dev/crypto/tls#pkg-constants)", FlagCipherAllowList, name)
		}
	}

	if !allDenied {
		for _, name := range denyList {
			if cs, found := allByName[name]; found {
				allowed = RemoveFromUInt16Array(allowed, cs.ID)
			} else {
				return func(c *tls.Config) {},
					fmt.Errorf("invalid --%s command line flag, %s is not a recognised cipher name (see https://pkg.go.dev/crypto/tls#pkg-constants)", FlagCipherDenyList, name)
			}
		}
	}

	for i := range allowed {
		if cs, found := all[allowed[i]]; found && cs.Insecure {
			setupLog.Info("WARNING: An insecure cipher suite has been enabled", "SuiteName", cs.Name)
		}
	}

	return func(c *tls.Config) {
		c.CipherSuites = allowed
	}, nil
}

func GetExtraEnvVars() []string {
	return GetViper().GetStringSlice(FlagEnvVar)
}

func GetExtraJvmArgs() []string {
	return GetViper().GetStringSlice(FlagJvmArg)
}

func stringSliceToMap(args []string, flag string) (map[string]string, error) {
	var m map[string]string
	if args != nil {
		m = make(map[string]string)
		for _, arg := range args {
			kv := strings.SplitN(arg, "=", 2)
			if len(kv) <= 1 {
				return nil, fmt.Errorf("invalid argument --%s=%s - must be in the format --%s=key=value",
					flag, arg, flag)
			}
			key := strings.TrimSpace(kv[0])
			value := strings.TrimSpace(kv[1])
			if key == "" || value == "" {
				return nil, fmt.Errorf("invalid argument --%s=%s - must be in the format --%s=\"key=value\" where the key and value cannot be blank",
					flag, arg, flag)
			}
			m[key] = value
		}
	}
	return m, nil
}

// RemoveAllFromUInt16Array removes all the specified values from a uint16 array, returning the updated array.
func RemoveAllFromUInt16Array(arr []uint16, toRemove ...uint16) []uint16 {
	for _, i := range toRemove {
		arr = RemoveFromUInt16Array(arr, i)
	}
	return arr
}

// RemoveFromUInt16Array removes a value from a uint16 array, returning the updated array.
func RemoveFromUInt16Array(arr []uint16, toRemove uint16) []uint16 {
	for i, ui := range arr {
		if ui == toRemove {
			return append(arr[:i], arr[i+1:]...)
		}
	}
	return arr
}
