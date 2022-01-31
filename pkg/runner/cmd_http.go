/*
 * Copyright (c) 2022, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package runner

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
)

const (
	// ArgURL is the Operator URL status command argument.
	ArgURL = "url"
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

func createHttpClient(cmd *cobra.Command) (http.Client, error) {
	client := http.Client{}
	flagSet := cmd.Flags()

	i, err := flagSet.GetBool(ArgSkipInsecure)
	if err != nil {
		return client, err
	}
	clientCertFile, err := flagSet.GetString(ArgCert)
	if err != nil {
		return client, err
	}
	clientKeyFile, err := flagSet.GetString(ArgKey)
	if err != nil {
		return client, err
	}
	caCertFile, err := flagSet.GetString(ArgCertAuthority)
	if err != nil {
		return client, err
	}

	var certs []tls.Certificate
	var caCertPool *x509.CertPool

	if clientCertFile != "" && clientKeyFile != "" {
		cert, err := tls.LoadX509KeyPair(clientCertFile, clientKeyFile)
		if err != nil {
			return client, errors.Wrapf(err, "creating x509 keypair from client cert file '%s' and client key file '%s'", clientCertFile, clientKeyFile)
		}
		certs = []tls.Certificate{cert}
	}

	if caCertFile != "" {
		caCert, err := ioutil.ReadFile(caCertFile)
		if err != nil {
			return client, errors.Wrapf(err, "opening cert file %s", caCertFile)
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			Certificates: certs,
			RootCAs:      caCertPool,
		},
	}

	tr.TLSClientConfig.InsecureSkipVerify = i
	client.Transport = tr

	return client, nil
}

// executeQuery performs a http on a URL
func executeQuery(cmd *cobra.Command) error {
	var err error

	flagSet := cmd.Flags()
	url, err := flagSet.GetString(ArgURL)
	if err != nil {
		return err
	}

	client, err := createHttpClient(cmd)
	if err != nil {
		return err
	}

	_, status, err := httpGet(url, client)
	if err == nil && status == http.StatusOK {
		return nil
	}

	return fmt.Errorf("failed to receive a 200 response from %s", url)
}
