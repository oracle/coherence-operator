/*
 * Copyright (c) 2022, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package runner

import (
	"fmt"
	v1 "github.com/oracle/coherence-operator/api/v1"
	"github.com/spf13/cobra"
	"os"
	"strconv"
)

const (
	// CommandReady is the argument to execute ready check.
	CommandReady = "ready"
)

// statusCommand creates the Corba "status" sub-command
func readyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   CommandReady,
		Short: "Run a Coherence ready check",
		Long:  "Run a Coherence ready check",
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeQuery(cmd)
		},
	}

	port := strconv.Itoa(int(v1.DefaultHealthPort))
	if p, found := os.LookupEnv(v1.EnvVarCohHealthPort); found {
		port = p
	}

	url := fmt.Sprintf("http://localhost:%s/ready", port)

	flagSet := cmd.Flags()
	flagSet.String(ArgURL, url, "The URL of the Coherence ready endpoint")
	flagSet.Bool(ArgSkipInsecure, false, "If true, the server's certificate will not be checked for validity")
	flagSet.String(ArgCertAuthority, "", "Path to a cert file for the certificate authority")
	flagSet.String(ArgCert, "", "Path to a client certificate file for TLS")
	flagSet.String(ArgKey, "", "Path to a client key file for TLS")

	return cmd
}
