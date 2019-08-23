package management

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/url"
	"strings"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/client-go/tools/remotecommand"
	utilexec "k8s.io/client-go/util/exec"
)

type ExecRequest struct {
	Pod       string
	Container string
	Namespace string
	Command   []string
	Arg       []string
	Timeout   time.Duration
}

// Execute a command in a Pod.
func PodExec(req *ExecRequest, config *rest.Config) (int, string, string, error) {
	kubeClient := kubernetes.NewForConfigOrDie(config)

	timeout := req.Timeout
	if timeout < time.Second*10 {
		timeout = time.Second * 10
	}

	execRequest := kubeClient.CoreV1().RESTClient().Post().
		Timeout(timeout).
		Resource("pods").
		Name(req.Pod).
		Namespace(req.Namespace).
		SubResource("exec").
		Param("stdin", "true").
		Param("stdout", "true").
		Param("stderr", "true").
		Param("tty", "false")

	if req.Container != "" {
		execRequest.Param("container", req.Container)
	}

	for _, cmd := range req.Command {
		execRequest.Param("command", cmd)
	}

	exec, err := remotecommand.NewSPDYExecutor(config, "POST", execRequest.URL())
	if err != nil {
		return 1, "", "", err
	}

	stdIn := newStringReader(req.Arg)
	stdOut := new(streamCapture)
	stdErr := new(streamCapture)

	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  stdIn,
		Stdout: stdOut,
		Stderr: stdErr,
		Tty:    false,
	})

	outStr := strings.Join(stdOut.Str, "")
	errStr := strings.Join(stdErr.Str, "")

	var exitCode int

	if err == nil {
		exitCode = 0
	} else {
		if exitErr, ok := err.(utilexec.ExitError); ok && exitErr.Exited() {
			exitCode = exitErr.ExitStatus()
		} else {
			return 1, outStr, errStr, errors.New("failed to find exit code")
		}
	}

	return exitCode, outStr, errStr, nil
}

type streamCapture struct {
	Str []string
}

func (w *streamCapture) Write(p []byte) (n int, err error) {
	str := string(p)
	if len(str) > 0 {
		w.Str = append(w.Str, str)
	}
	return len(str), nil
}

func newStringReader(ss []string) io.Reader {
	formattedString := strings.Join(ss, "\n")
	reader := strings.NewReader(formattedString)
	return reader
}

// getKubeconfigAndNamespace returns the *rest.Config and default namespace defined in the
// kubeconfig at the specified path. If no path is provided, returns the default *rest.Config
// and namespace
func getKubeconfigAndNamespace() (*rest.Config, string, error) {
	var clientConfig clientcmd.ClientConfig
	var apiConfig *clientcmdapi.Config
	var err error

	apiConfig, err = clientcmd.NewDefaultClientConfigLoadingRules().Load()
	if err != nil {
		return nil, "", fmt.Errorf("failed to get kubeconfig: %v", err)
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
