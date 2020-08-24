/*
 * Copyright (c) 2020, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 *
 */

package predicates

import (
	"sigs.k8s.io/controller-runtime/pkg/event"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

var _ predicate.Predicate = SecondaryPredicate{}

// SecondaryPredicate is a predicate that filters events for resources
// created as dependents of a primary resource. It follows the following
// rules:
//
//   - Create events are ignored because it is assumed that the controller
//     reconciling the parent is the client creating the dependent
//     resources.
//   - Update events are always handled.
//   - Deletion events are always handled because a controller will
//     typically want to recreate deleted dependent resources if the
//     primary resource is not deleted.
//   - Generic events are ignored.
//
// SecondaryPredicate is most often used in conjunction with
// controller-runtime's handler.EnqueueRequestForOwner
type SecondaryPredicate struct {
	predicate.Funcs
}

// Create filters out all events. It assumes that the controller
// reconciling the parent is the only client creating the dependent
// resources.
func (SecondaryPredicate) Create(e event.CreateEvent) bool {
	return false
}

// Update passes all events through.
func (SecondaryPredicate) Update(e event.UpdateEvent) bool {
	return true
}

// Delete passes all events through. This allows the controller to
// recreate deleted dependent resources if the primary resource is
// not deleted.
func (SecondaryPredicate) Delete(e event.DeleteEvent) bool {
	return true
}

// Generic filters out all events.
func (SecondaryPredicate) Generic(e event.GenericEvent) bool {
	return false
}

