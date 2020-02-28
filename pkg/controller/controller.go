/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package controller

import (
	"github.com/oracle/coherence-operator/pkg/flags"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// AddToManagerFuncs is a list of functions to add all Controllers to the Manager
var AddToManagerFuncs []func(manager.Manager, *flags.CoherenceOperatorFlags) error

// AddToManager adds all Controllers to the Manager
func AddToManager(m manager.Manager, opFlags *flags.CoherenceOperatorFlags) error {
	for _, f := range AddToManagerFuncs {
		if err := f(m, opFlags); err != nil {
			return err
		}
	}
	return nil
}
