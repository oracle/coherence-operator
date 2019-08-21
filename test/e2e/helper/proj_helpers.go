package helper

import (
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"strings"
)

const (
	TestNamespaceEnv      = "TEST_NAMESPACE"
	TestManifestEnv       = "TEST_MANIFEST"
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
