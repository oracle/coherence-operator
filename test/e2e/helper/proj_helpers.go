/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package helper

import (
	"fmt"
	"github.com/ghodss/yaml"
	coh "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

const (
	TestNamespaceEnv      = "TEST_NAMESPACE"
	TestManifestEnv       = "TEST_MANIFEST"
	TestLocalManifestEnv  = "TEST_LOCAL_MANIFEST"
	TestGlobalManifestEnv = "TEST_GLOBAL_MANIFEST"
	TestSslSecretEnv      = "TEST_SSL_SECRET"
	TestManifestValuesEnv = "TEST_MANIFEST_VALUES"
	ImagePullSecretsEnv   = "IMAGE_PULL_SECRETS"
	CoherenceVersionEnv   = "COHERENCE_VERSION"

	defaultNamespace = "operator-test"

	buildDir       = "build"
	outDir         = buildDir + string(os.PathSeparator) + "_output"
	chartDir       = outDir + string(os.PathSeparator) + "helm-charts"
	coherenceChart = chartDir + string(os.PathSeparator) + "coherence"
	operatorChart  = chartDir + string(os.PathSeparator) + "coherence-operator"
	testLogs       = outDir + string(os.PathSeparator) + "test-logs"
	certs          = outDir + string(os.PathSeparator) + "certs"
	deploy         = "deploy"
	crds           = deploy + string(os.PathSeparator) + "crds"
	manifest       = outDir + string(os.PathSeparator) + "manifest"
)

func GetTestNamespace() string {
	ns := os.Getenv(TestNamespaceEnv)
	if ns == "" {
		ns = defaultNamespace
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

func FindBuildDir() (string, error) {
	pd, err := FindProjectRootDir()
	if err != nil {
		return "", err
	}
	return pd + string(os.PathSeparator) + buildDir, nil
}

func FindCrdDir() (string, error) {
	pd, err := FindProjectRootDir()
	if err != nil {
		return "", err
	}
	return pd + string(os.PathSeparator) + crds, nil
}

func FindBuildOutputDir() (string, error) {
	pd, err := FindProjectRootDir()
	if err != nil {
		return "", err
	}

	return pd + string(os.PathSeparator) + outDir, nil
}

func FindHelmChartsDir() (string, error) {
	pd, err := FindProjectRootDir()
	if err != nil {
		return "", err
	}

	return pd + string(os.PathSeparator) + chartDir, nil
}

func FindCoherenceHelmChartDir() (string, error) {
	pd, err := FindProjectRootDir()
	if err != nil {
		return "", err
	}

	return pd + string(os.PathSeparator) + coherenceChart, nil
}

func FindOperatorHelmChartDir() (string, error) {
	pd, err := FindProjectRootDir()
	if err != nil {
		return "", err
	}

	return pd + string(os.PathSeparator) + operatorChart, nil
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

// NewCoherenceCluster creates a new CoherenceCluster from the default yaml file.
func NewCoherenceCluster(namespace string) (coh.CoherenceCluster, error) {
	return createCoherenceClusterFromYaml(namespace)
}

// NewCoherenceClusterFromYaml creates a new CoherenceCluster from a yaml file.
func NewCoherenceClusterFromYaml(namespace string, files ...string) (coh.CoherenceCluster, error) {
	if len(files) == 0 {
		return coh.CoherenceCluster{}, fmt.Errorf("no yaml files specified (did you specify a file instead of a namespace as the first argument?)")
	}
	return createCoherenceClusterFromYaml(namespace, files...)
}

// createCoherenceClusterFromYaml creates a new CoherenceCluster from a yaml file.
func createCoherenceClusterFromYaml(namespace string, files ...string) (coh.CoherenceCluster, error) {
	c := coh.CoherenceCluster{}

	l := coherenceClusterLoader{}
	err := l.loadYaml(&c, files...)

	if namespace != "" {
		c.SetNamespace(namespace)
	}

	return c, err
}

type coherenceClusterLoader struct {
}

// Load this CoherenceCluster from the specified yaml file
func (in *coherenceClusterLoader) FromYaml(cluster *coh.CoherenceCluster, files ...string) error {
	return in.loadYaml(cluster, files...)
}

func (in *coherenceClusterLoader) loadYaml(cluster *coh.CoherenceCluster, files ...string) error {
	if in == nil || files == nil {
		return nil
	}

	// try loading common-coherence-cluster.yaml first as this contains various values common
	// to all test structures as well as values replaced by test environment variables.
	_, c, _, _ := runtime.Caller(0)
	dir := filepath.Dir(c)
	common := dir + string(os.PathSeparator) + "common-coherence-cluster.yaml"
	err := in.loadYamlFromFile(cluster, common)
	if err != nil {
		return err
	}

	// Append any
	secrets := GetImagePullSecrets()
	cluster.Spec.ImagePullSecrets = append(cluster.Spec.ImagePullSecrets, secrets...)

	for _, file := range files {
		err := in.loadYamlFromFile(cluster, file)
		if err != nil {
			return err
		}
	}

	return nil
}

func (in *coherenceClusterLoader) loadYamlFromFile(cluster *coh.CoherenceCluster, file string) error {
	if in == nil || file == "" {
		return nil
	}

	actualFile, err := in.findActualFile(file)
	if err != nil {
		return err
	}

	data, err := ioutil.ReadFile(actualFile)
	if err != nil {
		return errors.New("Failed to read file " + actualFile + " caused by " + err.Error())
	}

	// expand any ${env-var} references in the yaml file
	s := os.ExpandEnv(string(data))

	err = yaml.Unmarshal([]byte(s), cluster)
	if err != nil {
		return errors.New("Failed to parse yaml file " + actualFile + " caused by " + err.Error())
	}

	return nil
}

func (in *coherenceClusterLoader) findActualFile(file string) (string, error) {
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
