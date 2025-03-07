/*
 * Copyright (c) 2020, 2025, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package helper

import (
	"fmt"
	coh "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/pkg/operator"
	"github.com/pkg/errors"
	"io"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

const (
	// TestNamespaceEnv is environment variable holding the name of the test k8s namespace.
	TestNamespaceEnv = "OPERATOR_NAMESPACE"
	// TestClusterNamespaceEnv is environment variable holding the name of the coherence cluster test k8s namespace.
	TestClusterNamespaceEnv = "CLUSTER_NAMESPACE"
	// TestClientNamespaceEnv is environment variable holding the name of the client test k8s namespace.
	TestClientNamespaceEnv = "OPERATOR_NAMESPACE_CLIENT"
	// PrometheusNamespaceEnv is environment variable holding the name of the Prometheus k8s namespace.
	PrometheusNamespaceEnv = "PROMETHEUS_NAMESPACE"
	// OperatorImageRegistryEnv is environment variable holding the registry of the Operator image.
	OperatorImageRegistryEnv = "OPERATOR_IMAGE_REGISTRY"
	// OperatorImageNameEnv is environment variable holding the name part of the Operator image.
	OperatorImageNameEnv = "OPERATOR_IMAGE_NAME"
	// CoherenceImageRegistryEnv is environment variable holding the registry of the default Coherence image.
	CoherenceImageRegistryEnv = "COHERENCE_IMAGE_REGISTRY"
	// CoherenceImageNameEnv is environment variable holding the name part of the default Coherence image.
	CoherenceImageNameEnv = "COHERENCE_IMAGE_NAME"
	// CoherenceImageTagEnv is environment variable holding the tag part of the default Coherence image.
	CoherenceImageTagEnv = "COHERENCE_IMAGE_TAG"
	// OperatorImageEnv is environment variable holding the full name of the Operator image.
	OperatorImageEnv = "OPERATOR_IMAGE"
	// ClientImageEnv is environment variable holding the name of the client test image.
	ClientImageEnv = "TEST_APPLICATION_IMAGE_CLIENT"
	// CohCompatibilityImageEnv is environment variable holding the name of the compatibility test image.
	CohCompatibilityImageEnv = "TEST_COMPATIBILITY_IMAGE"
	// TestSslSecretEnv is environment variable holding the name of the SSL certs secret.
	TestSslSecretEnv = "TEST_SSL_SECRET"
	// ImagePullSecEnv is environment variable holding the name of the image pull secrets.
	ImagePullSecEnv = "IMAGE_PULL_SECRETS"
	// CoherenceVersionEnv is environment variable holding the Coherence version.
	CoherenceVersionEnv = "COHERENCE_VERSION"
	// BuildOutputEnv is environment variable holding the build output directory.
	BuildOutputEnv = "BUILD_OUTPUT"
	// VersionEnv is environment variable holding the version the Operator.
	VersionEnv = "VERSION"

	defaultNamespace        = "operator-test"
	defaultClusterNamespace = "coherence-test"
	defaultClientNamespace  = "operator-test-client"

	defaultBuildDirectory = "build/_output"

	buildDir      = "build"
	outDir        = buildDir + string(os.PathSeparator) + "_output"
	chartDir      = outDir + string(os.PathSeparator) + "helm-charts"
	operatorChart = chartDir + string(os.PathSeparator) + "coherence-operator"
	testLogs      = outDir + string(os.PathSeparator) + "test-logs"
	testCharts    = outDir + string(os.PathSeparator) + "test-charts"
	certs         = outDir + string(os.PathSeparator) + "certs"
)

func EnsureTestEnvVars() {
	ensureEnvVar("TEST_IMAGE_PULL_POLICY", "IfNotPresent")
	ensureEnvVar("TEST_SKIP_SITE", "false")

	ensureEnvVar("K3D_OPERATOR_IMAGE", "k3d-myregistry.localhost:12345/oracle/coherence-operator:1.0.0")

	ensureEnvVar("TEST_COMPATIBILITY_IMAGE", "container-registry.oracle.com/middleware/operator-test-compatibility:1.0.0")
	ensureEnvVar("TEST_APPLICATION_IMAGE_CLIENT", "container-registry.oracle.com/middleware/operator-test-client:1.0.0")
	ensureEnvVar("TEST_APPLICATION_IMAGE", "container-registry.oracle.com/middleware/operator-test:1.0.0")
	ensureEnvVar("TEST_APPLICATION_IMAGE_HELIDON", "container-registry.oracle.com/middleware/operator-test-helidon:1.0.0")
	ensureEnvVar("TEST_APPLICATION_IMAGE_SPRING", "container-registry.oracle.com/middleware/operator-test-spring:1.0.0")
	ensureEnvVar("TEST_APPLICATION_IMAGE_SPRING_FAT", "container-registry.oracle.com/middleware/operator-test-spring-fat:1.0.0")
	ensureEnvVar("TEST_APPLICATION_IMAGE_SPRING_CNBP", "container-registry.oracle.com/middleware/operator-test-spring-cnbp:1.0.0")
}

func ensureEnvVar(key, value string) {
	if _, found := os.LookupEnv(key); !found {
		_ = os.Setenv(key, value)
	}
}

// GetOperatorVersionEnvVar returns the Operator version.
func GetOperatorVersionEnvVar() string {
	return os.Getenv(VersionEnv)
}

// GetOperatorImage returns the full name of the Operator image.
func GetOperatorImage() string {
	return os.Getenv(OperatorImageEnv)
}

// GetOperatorImageRegistry returns the registry name of the Operator image.
func GetOperatorImageRegistry() string {
	if s, found := os.LookupEnv(OperatorImageRegistryEnv); found {
		return s
	}
	return "container-registry.oracle.com/middleware"
}

// GetOperatorImageName returns the name part of the Operator image.
func GetOperatorImageName() string {
	if s, found := os.LookupEnv(OperatorImageNameEnv); found {
		return s
	}
	return "coherence-operator"
}

// GetDefaultCoherenceImage returns the full name of the default Coherence image.
func GetDefaultCoherenceImage() string {
	return os.Getenv(operator.EnvVarCoherenceImage)
}

// GetDefaultCoherenceImageRegistry returns the registry name of the default Coherence image.
func GetDefaultCoherenceImageRegistry() string {
	if s, found := os.LookupEnv(CoherenceImageRegistryEnv); found {
		return s
	}
	return "container-registry.oracle.com/middleware"
}

// GetDefaultCoherenceImageName returns the name part of the default Coherence image.
func GetDefaultCoherenceImageName() string {
	if s, found := os.LookupEnv(CoherenceImageNameEnv); found {
		return s
	}
	return "coherence-ce"
}

// GetDefaultCoherenceImageTag returns the tag part of the default Coherence image.
func GetDefaultCoherenceImageTag() string {
	if s, found := os.LookupEnv(CoherenceImageTagEnv); found {
		return s
	}
	return ""
}

// GetClientImage returns the name of the client test image
func GetClientImage() string {
	return os.Getenv(ClientImageEnv)
}

// GetCoherenceCompatibilityImage returns the name of the compatibility test image.
func GetCoherenceCompatibilityImage() string {
	return os.Getenv(CohCompatibilityImageEnv)
}

// GetTestNamespace returns the name of the test namespace.
func GetTestNamespace() string {
	ns := os.Getenv(TestNamespaceEnv)
	if ns == "" {
		ns = defaultNamespace
	}
	return ns
}

// GetTestClusterNamespace returns the name of the test cluster namespace.
func GetTestClusterNamespace() string {
	ns := os.Getenv(TestClusterNamespaceEnv)
	if ns == "" {
		ns = defaultClusterNamespace
	}
	return ns
}

// GetBuildOutputDirectory returns the build output directory
func GetBuildOutputDirectory() (os.FileInfo, error) {
	name := os.Getenv(BuildOutputEnv)
	if name == "" {
		name = defaultBuildDirectory
	}
	return os.Stat(name)
}

// GetTestClientNamespace returns the name of the client test namespace.
func GetTestClientNamespace() string {
	ns := os.Getenv(TestClientNamespaceEnv)
	if ns == "" {
		ns = defaultClientNamespace
	}
	return ns
}

// GetPrometheusNamespace returns the name of the Prometheus namespace.
func GetPrometheusNamespace() string {
	ns := os.Getenv(PrometheusNamespaceEnv)
	if ns == "" {
		ns = "monitoring"
	}
	return ns
}

// GetTestSSLSecretName returns the name of the SSL cert secret.
func GetTestSSLSecretName() string {
	return os.Getenv(TestSslSecretEnv)
}

// GetImagePullSecrets returns the names of the image pull secrets.
func GetImagePullSecrets() []coh.LocalObjectReference {
	s := os.Getenv(ImagePullSecEnv)
	if s == "" {
		return nil
	}
	var secrets []coh.LocalObjectReference
	for _, s := range strings.Split(s, ",") {
		secrets = append(secrets, coh.LocalObjectReference{Name: s})
	}
	return secrets
}

// FindProjectRootDir returns the project root directory.
func FindProjectRootDir() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", errors.Wrap(err, "error while checking if current directory is the project root")
	}

	for wd != "/" && wd != "." {
		_, err := os.Stat(wd + "/go.mod")
		if err == nil {
			return wd, nil
		} else if !os.IsNotExist(err) {
			return "", errors.Wrap(err, "error while checking if current directory is the project root")
		}
		wd = filepath.Dir(wd)
	}

	return "", os.ErrNotExist
}

// FindOperatorHelmChartDir returns the Operator Helm chart directory.
func FindOperatorHelmChartDir() (string, error) {
	pd, err := FindProjectRootDir()
	if err != nil {
		return "", err
	}

	return pd + string(os.PathSeparator) + operatorChart, nil
}

// FindOperatorTestHelmChartDir returns the Operator test Helm chart directory
// where previous version charts are downloaded to.
func FindOperatorTestHelmChartDir() (string, error) {
	pd, err := FindProjectRootDir()
	if err != nil {
		return "", err
	}

	return pd + string(os.PathSeparator) + testCharts, nil
}

// FindTestLogsDir returns the test log directory.
func FindTestLogsDir() (string, error) {
	pd, err := FindProjectRootDir()
	if err != nil {
		return "", err
	}

	return pd + string(os.PathSeparator) + testLogs, nil
}

// FindToolsDir returns the build tools directory.
func FindToolsDir() (string, error) {
	pd, err := FindProjectRootDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(pd, "build", "tools"), nil
}

// FindK8sApiToolsDir returns the k8s API build tools directory.
func FindK8sApiToolsDir() (string, error) {
	pd, err := FindProjectRootDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(pd, "build", "tools", "bin", "k8s", fmt.Sprintf("1.31.0-%s-%s", runtime.GOOS, runtime.GOARCH)), nil
}

// FindRuntimeCrdDir returns the CRD directory under the runtime assets.
func FindRuntimeCrdDir() (string, error) {
	pd, err := FindProjectRootDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(pd, "pkg", "data", "assets"), nil
}

// FindTestCertsDir returns the test cert directory.
func FindTestCertsDir() (string, error) {
	pd, err := FindProjectRootDir()
	if err != nil {
		return "", err
	}

	return pd + string(os.PathSeparator) + certs, nil
}

// NewSingleCoherenceFromYaml creates a single new Coherence resource from a yaml file.
func NewSingleCoherenceFromYaml(namespace, file string) (coh.Coherence, error) {
	return NewSingleCoherenceFromYamlWithSuffix(namespace, file, "")
}

// NewSingleCoherenceFromYamlWithSuffix creates a single new Coherence resource from a yaml file
// and adds a specified suffix to the Coherence resource names.
func NewSingleCoherenceFromYamlWithSuffix(namespace, file, suffix string) (coh.Coherence, error) {
	c, err := NewCoherenceFromYamlWithSuffix(namespace, file, suffix)
	switch {
	case err == nil && len(c) == 0:
		return coh.Coherence{}, fmt.Errorf("no deployments created from yaml %s", file)
	case err != nil:
		return coh.Coherence{}, err
	default:
		return c[0], err
	}
}

// NewCoherenceFromYaml creates a new Coherence resource from a yaml file.
func NewCoherenceFromYaml(namespace, file string) ([]coh.Coherence, error) {
	return NewCoherenceFromYamlWithSuffix(namespace, file, "")
}

// NewCoherenceFromYamlWithSuffix creates a new Coherence resource from a yaml file.
// and adds a specified suffix to the Coherence resource names
func NewCoherenceFromYamlWithSuffix(namespace, file, suffix string) ([]coh.Coherence, error) {
	res, err := createCoherenceFromYaml(namespace, file)
	if err == nil && suffix != "" {
		for _, c := range res {
			c.Name = c.Name + suffix
		}
	}
	return res, err
}

// createCoherenceFromYaml creates a new Coherence resource from a yaml file.
func createCoherenceFromYaml(namespace, file string) ([]coh.Coherence, error) {
	l := CoherenceLoader{}
	return l.loadYaml(namespace, file)
}

// NewSingleCoherenceJobFromYaml creates a single new CoherenceJob resource from a yaml file.
func NewSingleCoherenceJobFromYaml(namespace, file string) (coh.CoherenceJob, error) {
	deps, err := NewCoherenceJobFromYaml(namespace, file)
	switch {
	case err == nil && len(deps) == 0:
		return coh.CoherenceJob{}, fmt.Errorf("no deployments created from yaml %s", file)
	case err != nil:
		return coh.CoherenceJob{}, err
	default:
		return deps[0], err
	}
}

// NewCoherenceJobFromYaml creates a new CoherenceJob resource from a yaml file.
func NewCoherenceJobFromYaml(namespace string, file string) ([]coh.CoherenceJob, error) {
	return createCoherenceJobFromYaml(namespace, file)
}

// createCoherenceJobFromYaml creates a new CoherenceJob resource from a yaml file.
func createCoherenceJobFromYaml(namespace string, file string) ([]coh.CoherenceJob, error) {
	l := CoherenceLoader{}
	return l.loadJobYaml(namespace, file)
}

// CoherenceLoader can load Coherence resources from yaml files.
type CoherenceLoader struct {
}

func (in *CoherenceLoader) loadYaml(namespace, file string) ([]coh.Coherence, error) {
	var deployments []coh.Coherence

	if in == nil {
		return deployments, nil
	}

	// try loading common-coherence-deployment.yaml first as this contains various values common
	// to all test structures as well as values replaced by test environment variables.
	_, c, _, _ := runtime.Caller(0)
	dir := filepath.Dir(c)
	common := dir + string(os.PathSeparator) + "common-coherence-deployment.yaml"
	templates, err := in.loadYamlFromFile(coh.Coherence{}, common)
	if err != nil {
		return deployments, err
	}

	if len(templates) == 0 {
		return deployments, fmt.Errorf("could not load any deployment templates")
	}
	template := templates[0]

	if namespace != "" {
		template.SetNamespace(namespace)
	}

	// Append any pull secrets
	secrets := GetImagePullSecrets()
	template.Spec.ImagePullSecrets = append(template.Spec.ImagePullSecrets, secrets...)

	if file != "" {
		deployments, err = in.loadYamlFromFile(template, file)
	} else {
		deployments = append(deployments, template)
	}

	// add environment variables
	skipSite := os.Getenv(coh.EnvVarCohSkipSite)
	if skipSite == "true" {
		for i, d := range deployments {
			d.Spec.AddEnvVarIfAbsent(corev1.EnvVar{Name: coh.EnvVarCohSkipSite, Value: "true"})
			deployments[i] = d
		}
	}

	return deployments, err
}

func (in *CoherenceLoader) loadJobYaml(namespace, file string) ([]coh.CoherenceJob, error) {
	var deployments []coh.CoherenceJob

	if in == nil {
		return deployments, nil
	}

	// try loading common-coherence-deployment.yaml first as this contains various values common
	// to all test structures as well as values replaced by test environment variables.
	_, c, _, _ := runtime.Caller(0)
	dir := filepath.Dir(c)
	common := dir + string(os.PathSeparator) + "common-coherencejob-deployment.yaml"
	templates, err := in.loadJobYamlFromFile(coh.CoherenceJob{}, common)
	if err != nil {
		return deployments, err
	}

	if len(templates) == 0 {
		return deployments, fmt.Errorf("could not load any deployment templates")
	}
	template := templates[0]

	if namespace != "" {
		template.SetNamespace(namespace)
	}

	// Append any pull secrets
	secrets := GetImagePullSecrets()
	template.Spec.ImagePullSecrets = append(template.Spec.ImagePullSecrets, secrets...)

	if file != "" {
		deployments, err = in.loadJobYamlFromFile(template, file)
	} else {
		deployments = append(deployments, template)
	}

	// add environment variables
	skipSite := os.Getenv(coh.EnvVarCohSkipSite)
	if skipSite == "true" {
		for i, d := range deployments {
			d.Spec.AddEnvVarIfAbsent(corev1.EnvVar{Name: coh.EnvVarCohSkipSite, Value: "true"})
			deployments[i] = d
		}
	}

	return deployments, err
}

func (in *CoherenceLoader) loadYamlFromFile(template coh.Coherence, file string) ([]coh.Coherence, error) {
	var deployments []coh.Coherence
	if in == nil || file == "" {
		return deployments, nil
	}

	actualFile, err := FindActualFile(file)
	if err != nil {
		return deployments, err
	}

	// read the whole file
	data, err := os.ReadFile(actualFile)
	if err != nil {
		return deployments, errors.New("Failed to read file " + actualFile + " caused by " + err.Error())
	}

	// expand any ${env-var} references in the yaml file

	s := os.ExpandEnv(string(data))

	// Get the yaml decoder
	decoder := yaml.NewYAMLToJSONDecoder(strings.NewReader(s))

	for err == nil {
		deployment := coh.Coherence{}
		template.DeepCopyInto(&deployment)
		err = decoder.Decode(&deployment)
		if err == nil && deployment.Name != "" {
			deployments = append(deployments, deployment)
		}
	}

	if err != io.EOF {
		return deployments, errors.New("Failed to parse yaml file " + actualFile + " caused by " + err.Error())
	}

	return deployments, nil
}

func (in *CoherenceLoader) loadJobYamlFromFile(template coh.CoherenceJob, file string) ([]coh.CoherenceJob, error) {
	var deployments []coh.CoherenceJob
	if in == nil || file == "" {
		return deployments, nil
	}

	actualFile, err := FindActualFile(file)
	if err != nil {
		return deployments, err
	}

	// read the whole file
	data, err := os.ReadFile(actualFile)
	if err != nil {
		return deployments, errors.New("Failed to read file " + actualFile + " caused by " + err.Error())
	}

	// expand any ${env-var} references in the yaml file
	s := os.ExpandEnv(string(data))

	// Get the yaml decoder
	decoder := yaml.NewYAMLToJSONDecoder(strings.NewReader(s))

	for err == nil {
		deployment := coh.CoherenceJob{}
		template.DeepCopyInto(&deployment)
		err = decoder.Decode(&deployment)
		if err == nil && deployment.Name != "" {
			deployments = append(deployments, deployment)
		}
	}

	if err != io.EOF {
		return deployments, errors.New("Failed to parse yaml file " + actualFile + " caused by " + err.Error())
	}

	return deployments, nil
}

func (in *CoherenceLoader) LoadYamlIntoTemplate(template interface{}, file string) error {
	if in == nil || file == "" {
		return nil
	}

	actualFile, err := FindActualFile(file)
	if err != nil {
		return err
	}

	// read the whole file
	data, err := os.ReadFile(actualFile)
	if err != nil {
		return errors.New("Failed to read file " + actualFile + " caused by " + err.Error())
	}

	// expand any ${env-var} references in the yaml file
	s := os.ExpandEnv(string(data))

	// Get the yaml decoder
	decoder := yaml.NewYAMLToJSONDecoder(strings.NewReader(s))

	err = decoder.Decode(template)

	if err != io.EOF {
		return errors.New("Failed to parse yaml file " + actualFile + " caused by " + err.Error())
	}

	return nil
}

// LoadFromYamlFile loads the specified value from the yaml file.
func LoadFromYamlFile(file string, o interface{}) error {
	actualFile, err := FindActualFile(file)
	if err != nil {
		return err
	}

	// read the whole file
	data, err := os.ReadFile(actualFile)
	if err != nil {
		return errors.New("Failed to read file " + actualFile + " caused by " + err.Error())
	}

	// expand any ${env-var} references in the yaml file
	s := os.ExpandEnv(string(data))

	decoder := yaml.NewYAMLToJSONDecoder(strings.NewReader(s))
	return decoder.Decode(o)
}

func FindActualFile(file string) (string, error) {
	_, err := os.Stat(file)
	if err == nil {
		return file, nil
	}

	// file does not exist
	if !strings.HasPrefix(file, "/") {
		// the file does not exist and is not absolute so try relative to a location
		// in the call stack by walking up the stack and trying each location.
		i := 0
		for {
			_, caller, _, ok := runtime.Caller(i)
			if ok {
				dir := filepath.Dir(caller)
				f := dir + string(os.PathSeparator) + file
				_, e := os.Stat(f)
				if e == nil {
					return f, nil
				}
			} else {
				// no more call stack
				break
			}
			i++
		}
	}

	return "", err
}

// SkipIfCoherenceVersionLessThan skips the specified test if the current Coherence version set in the COHERENCE_VERSION
// environment variable is less than the specified version.
func SkipIfCoherenceVersionLessThan(t *testing.T, version ...int) {
	ok, err := IsCoherenceVersionAtLeast(version...)
	switch {
	case err == nil && !ok:
		versionStr := os.Getenv(CoherenceVersionEnv)
		t.Skipf("Skipping test as COHERENCE_VERSION %s is less than requested version %v", versionStr, version)
	case err != nil:
		t.Fatalf("Failed to check COHERENCE_VERSION due to %s", err.Error())
	}
}

// IsCoherenceVersionAtLeast determines whether current Coherence version set in the COHERENCE_VERSION
// environment variable is greater than the specified version or the COHERENCE_VERSION environment
// variable has not been set.
func IsCoherenceVersionAtLeast(version ...int) (bool, error) {
	if len(version) == 0 {
		return true, nil
	}

	versionStr := os.Getenv(CoherenceVersionEnv)
	if versionStr == "" {
		return true, nil
	}
	parts := strings.Split(versionStr, ".")

	for i, v := range version {
		if i >= len(parts) {
			break
		}
		vp, err := strconv.Atoi(parts[i])
		if err != nil {
			return false, err
		}
		if vp < v {
			return false, nil
		}
	}

	return true, nil
}

// GetKubeconfigAndNamespace returns the *rest.Config and default namespace defined in the
// kubeconfig at the specified path. If no path is provided, returns the default *rest.Config
// and namespace
func GetKubeconfigAndNamespace(configPath string) (*rest.Config, string, error) {
	var clientConfig clientcmd.ClientConfig
	var apiConfig *clientcmdapi.Config
	var err error
	if configPath != "" {
		apiConfig, err = clientcmd.LoadFromFile(configPath)
		if err != nil {
			return nil, "", fmt.Errorf("failed to load user provided kubeconfig: %v", err)
		}
	} else {
		apiConfig, err = clientcmd.NewDefaultClientConfigLoadingRules().Load()
		if err != nil {
			return nil, "", fmt.Errorf("failed to get kubeconfig: %v", err)
		}
	}
	clientConfig = clientcmd.NewDefaultClientConfig(*apiConfig, &clientcmd.ConfigOverrides{})
	kubeconfig, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, "", err
	}
	namespace, _, err := clientConfig.Namespace()
	if err != nil {
		return nil, "", err
	}

	u, err := url.Parse(kubeconfig.Host)
	if err != nil {
		return nil, "", err
	}

	ip, err := net.LookupIP(u.Hostname())
	if err != nil {
		return nil, "", err
	}

	// If this is Docker on Mac the host name resolves to loopback
	// It seems that if we use the host name we may later get a x509 error
	// but if we change the host to the loopback IP 127.0.0.1 it works fine
	if ip[0].IsLoopback() {
		kubeconfig.Host = strings.Replace(kubeconfig.Host, u.Hostname(), "127.0.0.1", 1)
	}

	return kubeconfig, namespace, nil
}
