/*
 * Copyright (c) 2020, 2026, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

// Package v1 contains API Schema definitions for the coherence.oracle.com v1 API group
// +kubebuilder:object:generate=true
// +groupName=coherence.oracle.com
package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	// GroupVersion is group version used to register these objects
	GroupVersion = schema.GroupVersion{Group: "coherence.oracle.com", Version: "v1"}

	// SchemeBuilder stays on apimachinery's registration path so importing the API package
	// avoids the deprecated controller-runtime helper while preserving AddToScheme behavior.
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)

	// AddToScheme adds the types in this group-version to the given scheme.
	AddToScheme = SchemeBuilder.AddToScheme
)

func addKnownTypes(scheme *runtime.Scheme) error {
	// Registering the root API objects here keeps scheme setup explicit for this group/version,
	// which replaces the deprecated controller-runtime object-registration helper.
	scheme.AddKnownTypes(GroupVersion, &Coherence{}, &CoherenceList{}, &CoherenceJob{}, &CoherenceJobList{})
	// AddToGroupVersion records the API metadata so serialized objects keep the expected
	// coherence.oracle.com/v1 identity after the registration path changes.
	metav1.AddToGroupVersion(scheme, GroupVersion)
	return nil
}
