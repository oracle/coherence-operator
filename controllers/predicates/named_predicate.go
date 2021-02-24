/*
 * Copyright (c) 2020, 2021, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package predicates

import (
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// NamedPredicate is an event handler that watches for a resource identified by
// a specific namespace and name.
type NamedPredicate struct {
	Name      string
	Namespace string
}

func (w NamedPredicate) Create(e event.CreateEvent) bool {
	return e.Object.GetNamespace() == w.Namespace && e.Object.GetName() == w.Name
}

func (w NamedPredicate) Delete(e event.DeleteEvent) bool {
	return e.Object.GetNamespace() == w.Namespace && e.Object.GetName() == w.Name
}

func (w NamedPredicate) Update(e event.UpdateEvent) bool {
	return e.ObjectNew.GetNamespace() == w.Namespace && e.ObjectNew.GetName() == w.Name
}

func (w NamedPredicate) Generic(e event.GenericEvent) bool {
	return e.Object.GetNamespace() == w.Namespace && e.Object.GetName() == w.Name
}

var _ predicate.Predicate = &NamedPredicate{}
