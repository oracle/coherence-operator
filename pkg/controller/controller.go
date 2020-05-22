/*
 * Copyright (c) 2019, 2020, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package controller

import (
	"github.com/oracle/coherence-operator/pkg/flags"
	"github.com/oracle/coherence-operator/pkg/rest"
	"github.com/pkg/errors"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// AddToManagerFuncs is a list of functions to add all Controllers to the Manager
var AddToManagerFuncs []func(manager.Manager) error

// AddToManager adds all Controllers to the Manager
func AddToManager(m manager.Manager) error {
	opFlags := flags.GetOperatorFlags()

	// Create the REST server
	s, err := rest.EnsureServer(m, opFlags)
	if err != nil {
		return errors.Wrap(err, "failed to create REST server")
	}
	// Add the REST server to the Manager so that is is started after the Manager is initialized
	err = s.Start()
	if err != nil {
		return errors.Wrap(err, "failed to start the REST server")
	}

	for _, f := range AddToManagerFuncs {
		if err := f(m); err != nil {
			return err
		}
	}
	return nil
}
