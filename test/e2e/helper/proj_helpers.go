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
	"strings"
)

const (
	TestNamespaceEnv      = "TEST_NAMESPACE"
	TestManifestEnv       = "TEST_MANIFEST"
	TestGlobalManifestEnv = "TEST_GLOBAL_MANIFEST"
	TestSslSecretEnv      = "TEST_SSL_SECRET"
	TestManifestValuesEnv = "TEST_MANIFEST_VALUES"
	ImagePullSecretsEnv   = "IMAGE_PULL_SECRETS"

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

func GetImagePullSecrets() []string {
	s := os.Getenv(ImagePullSecretsEnv)
	if s == "" {
		return nil
	}
	return strings.Split(s, ",")
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

// NewCoherenceClusterFromYaml creates a new CoherenceCluster from a yaml file.
func NewCoherenceClusterFromYaml(namespace string, files ...string) (coh.CoherenceCluster, error) {
	c := coh.CoherenceCluster{}

	if len(files) == 0 {
		return c, fmt.Errorf("no yaml files specified (did you specify a file instead of a namespace as the first argument?)")
	}

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
