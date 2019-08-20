package helper

import (
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"strings"
)

const (
	TestNamespaceEnv    = "TEST_NAMESPACE"
	ImagePullSecretsEnv = "IMAGE_PULL_SECRETS"

	defaultNamespace = "operator-test"

	buildDir       = "build"
	outDir         = buildDir + string(os.PathSeparator) + "_output"
	chartDir       = outDir + string(os.PathSeparator) + "helm-charts"
	coherenceChart = chartDir + string(os.PathSeparator) + "coherence"
	operatorChart  = chartDir + string(os.PathSeparator) + "coherence-operator"
	testLogs       = outDir + string(os.PathSeparator) + "test-logs"
	certs          = outDir + string(os.PathSeparator) + "certs"
)

func GetTestNamespace() string {
	ns := os.Getenv(TestNamespaceEnv)
	if ns == "" {
		ns = defaultNamespace
	}
	return ns
}

func GetImagePullSecrets() []string {
	return strings.Split(os.Getenv(ImagePullSecretsEnv), ",")
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
