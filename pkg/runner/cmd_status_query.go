/*
 * Copyright (c) 2021, 2022, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package runner

import (
	"fmt"
	v1 "github.com/oracle/coherence-operator/api/v1"
	"github.com/spf13/cobra"
	"net/http"
	"strconv"
	"time"
)

const (
	// CommandStatus is the argument to execute a deployment status query.
	CommandStatus = "status"

	// ArgOperatorURL is the Operator URL status command argument.
	ArgOperatorURL = "operator-url"
	// ArgNamespace is the Coherence resource namespace status command argument.
	ArgNamespace = "namespace"
	// ArgName is the Coherence resource name status command argument.
	ArgName = "name"
	// ArgCondition is the required condition status command argument.
	ArgCondition = "condition"
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
	flagSet.String(ArgOperatorURL, "http://coherence-operator-rest.coherence.svc.local:8000", "The Coherence Operator URL, typically the operator's REST service")
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
	url, err := flagSet.GetString(ArgOperatorURL)
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

	client, err := createHttpClient(cmd)
	if err != nil {
		return err
	}

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
