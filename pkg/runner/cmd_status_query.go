/*
 * Copyright (c) 2021, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package runner

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	v1 "github.com/oracle/coherence-operator/api/v1"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

const (
	// CommandStatus is the argument to execute a deployment status query.
	CommandStatus = "status"

	// ArgOperatorUrl is the Operator URL status command argument.
	ArgOperatorUrl = "operator-url"
	// ArgNamespace is the Coherence resource namespace status command argument.
	ArgNamespace = "namespace"
	// ArgName is the Coherence resource name status command argument.
	ArgName = "name"
	// ArgCondition is the required condition status command argument.
	ArgCondition = "condition"
	// ArgTimeout is the timeout status command argument.
	ArgTimeout = "timeout"
	// ArgInterval is the retry interval status command argument.
	ArgInterval = "interval"
	// ArgSkipInsecure is the skip insecure https checks status command argument.
	ArgSkipInsecure = "insecure-skip-tls-verify"
	// ArgCertAuthority is the location of the CA file status command argument.
	ArgCertAuthority = "certificate-authority"
	// ArgCert is the location of the cert file status command argument.
	ArgCert = "client-certificate"
	// ArgKey is the location of the key file status command argument.
	ArgKey = "client-key"
)

// statusCommand creates the Corba "status" sub-command
func statusCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   CommandStatus,
		Short: "Run a Coherence resource status query",
		Long:  "Run a Coherence resource status query to verify the status or the resource",
		RunE: func(cmd *cobra.Command, args []string) error {
			return statusQuery(cmd)
		},
	}

	flagSet := cmd.Flags()
	flagSet.String(ArgOperatorUrl, "http://coherence-operator-rest.coherence.svc.local:8000", "The Coherence Operator URL, typically the operator's REST service")
	flagSet.String(ArgNamespace, "", "The namespace the Coherence resource is deployed into")
	flagSet.String(ArgName, "", "The name of the Coherence resource")
	flagSet.String(ArgCondition, string(v1.ConditionTypeReady), "The required condition that the Coherence resource should be in")
	flagSet.Duration(ArgTimeout, time.Minute*5, "The maximum amount of time to wait for the Coherence resource to reach the required condition")
	flagSet.Duration(ArgInterval, time.Second*10, "The status check re-try interval")
	flagSet.Bool(ArgSkipInsecure, false, "If true, the Operator's certificate will not be checked for validity")
	flagSet.String(ArgCertAuthority, "", "Path to a cert file for the certificate authority")
	flagSet.String(ArgCert, "", "Path to a client certificate file for TLS")
	flagSet.String(ArgKey, "", "Path to a client key file for TLS")

	return cmd
}

// statusQuery performs a status query for a given Coherence deployment
func statusQuery(cmd *cobra.Command) error {
	var err error

	flagSet := cmd.Flags()
	url, err := flagSet.GetString(ArgOperatorUrl)
	if err != nil {
		return err
	}
	ns, err := flagSet.GetString(ArgNamespace)
	if err != nil {
		return err
	}
	n, err := flagSet.GetString(ArgName)
	if err != nil {
		return err
	}
	condition, err := flagSet.GetString(ArgCondition)
	if err != nil {
		return err
	}
	timeout, err := flagSet.GetDuration(ArgTimeout)
	if err != nil {
		return err
	}
	interval, err := flagSet.GetDuration(ArgInterval)
	if err != nil {
		return err
	}
	insecureSkipVerify, err := flagSet.GetBool(ArgSkipInsecure)
	if err != nil {
		return err
	}
	clientCertFile, err := flagSet.GetString(ArgCert)
	if err != nil {
		return err
	}
	clientKeyFile, err := flagSet.GetString(ArgKey)
	if err != nil {
		return err
	}
	caCertFile, err := flagSet.GetString(ArgCertAuthority)
	if err != nil {
		return err
	}

	var certs []tls.Certificate
	var caCertPool *x509.CertPool

	if clientCertFile != "" && clientKeyFile != "" {
		cert, err := tls.LoadX509KeyPair(clientCertFile, clientKeyFile)
		if err != nil {
			return errors.Wrapf(err, "creating x509 keypair from client cert file '%s' and client key file '%s'", clientCertFile, clientKeyFile)
		}
		certs = []tls.Certificate{cert}
	}

	if caCertFile != "" {
		caCert, err := ioutil.ReadFile(caCertFile)
		if err != nil {
			return errors.Wrapf(err, "opening cert file %s", caCertFile)
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			Certificates:       certs,
			RootCAs:            caCertPool,
			InsecureSkipVerify: insecureSkipVerify,
		},
	}

	client := http.Client{Transport: tr}
	request := fmt.Sprintf("%s/status/%s/%s?phase=%s", url, ns, n, condition)

	maxTime := time.Now().Add(timeout)

	// Retry until we get a 200 response, or we time out
	for time.Now().Before(maxTime) {
		msg, status, err := httpGet(request, client)
		if err == nil && status == http.StatusOK {
			return nil
		}
		log.Info("Status request backoff", "url", request, "retry", interval.String(), "response", msg, "status", strconv.Itoa(status))
		time.Sleep(interval)
	}

	// If we get here we timed out
	return fmt.Errorf("failed to receive a 200 response within a timeout of %s", timeout.String())
}
