/*
 * Copyright (c) 2019, 2020 Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package helper

import (
	"fmt"
	coh "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
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
	TestNamespaceEnv       = "TEST_NAMESPACE"
	PrometheusNamespaceEnv = "PROMETHEUS_NAMESPACE"
	TestManifestEnv        = "TEST_MANIFEST"
	TestLocalManifestEnv   = "TEST_LOCAL_MANIFEST"
	TestGlobalManifestEnv  = "TEST_GLOBAL_MANIFEST"
	TestSslSecretEnv       = "TEST_SSL_SECRET"
	TestManifestValuesEnv  = "TEST_MANIFEST_VALUES"
	ImagePullSecretsEnv    = "IMAGE_PULL_SECRETS"
	CoherenceVersionEnv    = "COHERENCE_VERSION"
	CompatibleVersionsEnv  = "COMPATIBLE_VERSIONS"
	VersionEnv             = "VERSION_FULL"
	OperatorImageEnv       = "OPERATOR_IMAGE"
	SkipCompatibilityEnv   = "SKIP_COMPATIBILITY"

	defaultNamespace = "operator-test"

	buildDir           = "build"
	outDir             = buildDir + string(os.PathSeparator) + "_output"
	chartDir           = outDir + string(os.PathSeparator) + "helm-charts"
	compatibleChartDir = outDir + string(os.PathSeparator) + "previous-charts"
	operatorChart      = chartDir + string(os.PathSeparator) + "coherence-operator"
	compatibleCharts   = compatibleChartDir + string(os.PathSeparator) + "coherence-operator"
	testLogs           = outDir + string(os.PathSeparator) + "test-logs"
	certs              = outDir + string(os.PathSeparator) + "certs"
	deploy             = "deploy"
	crds               = deploy + string(os.PathSeparator) + "crds"
	manifest           = outDir + string(os.PathSeparator) + "manifest"
)

func GetTestNamespace() string {
	ns := os.Getenv(TestNamespaceEnv)
	if ns == "" {
		ns = defaultNamespace
	}
	return ns
}

func GetPrometheusNamespace() string {
	ns := os.Getenv(PrometheusNamespaceEnv)
	if ns == "" {
		ns = GetTestNamespace()
	}
	return ns
}

func GetTestManifestFileName() (string, error) {
	man := os.Getenv(TestManifestEnv)
	if man == "" {
		dir, err := FindTestManifestDir()
		if err != nil {
			return "", err
		}
		man = dir + string(os.PathSeparator) + "test-manifest.yaml"
	}
	return man, nil
}

func GetTestLocalManifestFileName() (string, error) {
	man := os.Getenv(TestLocalManifestEnv)
	if man == "" {
		dir, err := FindTestManifestDir()
		if err != nil {
			return "", err
		}
		man = dir + string(os.PathSeparator) + "local-manifest.yaml"
	}
	return man, nil
}

func GetTestGlobalManifestFileName() (string, error) {
	man := os.Getenv(TestGlobalManifestEnv)
	if man == "" {
		dir, err := FindTestManifestDir()
		if err != nil {
			return "", err
		}
		man = dir + string(os.PathSeparator) + "global-manifest.yaml"
	}
	return man, nil
}

func GetTestManifestValuesFileName() string {
	return os.Getenv(TestManifestValuesEnv)
}

func GetOperatorImage() string {
	return os.Getenv(OperatorImageEnv)
}

func GetOperatorVersion() string {
	return os.Getenv(VersionEnv)
}

func GetCompatibleOperatorVersions() []string {
	var versions []string
	list, ok := os.LookupEnv(CompatibleVersionsEnv)
	if ok {
		versions = strings.Split(list, " ")
	} else {
		versions = []string{}
	}
	return versions
}

func GetTestSSLSecretName() string {
	return os.Getenv(TestSslSecretEnv)
}

func GetImagePullSecrets() []coh.LocalObjectReference {
	s := os.Getenv(ImagePullSecretsEnv)
	if s == "" {
		return nil
	}
	var secrets []coh.LocalObjectReference
	for _, s := range strings.Split(s, ",") {
		secrets = append(secrets, coh.LocalObjectReference{Name: s})
	}
	return secrets
}

func FindProjectRootDir() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", errors.Wrap(err, "error while checking if current directory is the project root")
	}

	for wd != "/" && wd != "." {
		_, err := os.Stat(wd + "/build/Dockerfile")
		if err == nil {
			return wd, nil
		} else if !os.IsNotExist(err) {
			return "", errors.Wrap(err, "error while checking if current directory is the project root")
		}
		wd = filepath.Dir(wd)
	}

	return "", os.ErrNotExist
}

func AssumeRunningCompatibilityTests(t *testing.T) {
	s := os.Getenv(SkipCompatibilityEnv)
	if strings.ToLower(s) == "true" {
		t.Skipf("Skipping compatibility tests, %s environment variable set to '%s'", SkipCompatibilityEnv, s)
	}
}

func FindCrdDir() (string, error) {
	pd, err := FindProjectRootDir()
	if err != nil {
		return "", err
	}
	return pd + string(os.PathSeparator) + crds, nil
}

func FindOperatorHelmChartDir() (string, error) {
	pd, err := FindProjectRootDir()
	if err != nil {
		return "", err
	}

	return pd + string(os.PathSeparator) + operatorChart, nil
}

func FindPreviousOperatorHelmChartDir(v string) (string, error) {
	pd, err := FindProjectRootDir()
	if err != nil {
		return "", err
	}

	return pd + string(os.PathSeparator) + compatibleCharts + "-" + v, nil
}

func FindTestLogsDir() (string, error) {
	pd, err := FindProjectRootDir()
	if err != nil {
		return "", err
	}

	return pd + string(os.PathSeparator) + testLogs, nil
}

func FindTestCertsDir() (string, error) {
	pd, err := FindProjectRootDir()
	if err != nil {
		return "", err
	}

	return pd + string(os.PathSeparator) + certs, nil
}

func FindTestManifestDir() (string, error) {
	pd, err := FindProjectRootDir()
	if err != nil {
		return "", err
	}

	return pd + string(os.PathSeparator) + manifest, nil
}

// NewCoherence creates a new Coherence resource from the default minimal yaml file.
func NewCoherence(namespace string) (coh.Coherence, error) {
	return NewSingleCoherenceFromYaml(namespace, "")
}

// NewSingleCoherenceFromYaml creates a single new Coherence resource from a yaml file.
func NewSingleCoherenceFromYaml(namespace string, file string) (coh.Coherence, error) {
	deps, err := NewCoherenceFromYaml(namespace, file)
	switch {
	case err == nil && len(deps) == 0:
		return coh.Coherence{}, fmt.Errorf("no deployments created from yaml %s", file)
	case err != nil:
		return coh.Coherence{}, err
	default:
		return deps[0], err
	}
}

// NewCoherenceFromYaml creates a new Coherence resource from a yaml file.
func NewCoherenceFromYaml(namespace string, file string) ([]coh.Coherence, error) {
	return createCoherenceFromYaml(namespace, file)
}

// createCoherenceFromYaml creates a new Coherence resource from a yaml file.
func createCoherenceFromYaml(namespace string, file string) ([]coh.Coherence, error) {
	l := CoherenceLoader{}
	return l.loadYaml(namespace, file)
}

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

	// Append any
	secrets := GetImagePullSecrets()
	template.Spec.ImagePullSecrets = append(template.Spec.ImagePullSecrets, secrets...)

	if file != "" {
		deployments, err = in.loadYamlFromFile(template, file)
	} else {
		deployments = append(deployments, template)
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
	data, err := ioutil.ReadFile(actualFile)
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

// Load the specified value from the yaml file.
func LoadFromYamlFile(file string, o interface{}) error {
	actualFile, err := FindActualFile(file)
	if err != nil {
		return err
	}

	// read the whole file
	data, err := ioutil.ReadFile(actualFile)
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

	// files does not exist
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
			i = i + 1
		}
	}

	return "", err
}

// Skip the specified test if the current Coherence version set in the COHERENCE_VERSION
// environment variable is less than the specified version.
func SkipIfCoherenceVersionLessThan(t *testing.T, version ...int) {
	ok, err := IsCoherenceVersionAtLeast(version...)
	switch {
	case err == nil && !ok:
		versionStr := os.Getenv(CoherenceVersionEnv)
		t.Skip(fmt.Sprintf("Skipping test as COHERENCE_VERSION %s is less than requested version %v", versionStr, version))
	case err != nil:
		t.Fatalf(fmt.Sprintf("Failed to check COHERENCE_VERSION due to %s", err.Error()))
	}
}

// Determine whether current Coherence version set in the COHERENCE_VERSION
// environment variable is greater than the specified version or the
// COHERENCE_VERSION environment variable has not been set.
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
	// It seems that if we use the host name we may later get an x509 error
	// but if we change the host to the loopback IP 127.0.0.1 it works fine
	if ip[0].IsLoopback() {
		kubeconfig.Host = strings.Replace(kubeconfig.Host, u.Hostname(), "127.0.0.1", 1)
	}

	return kubeconfig, namespace, nil
}
