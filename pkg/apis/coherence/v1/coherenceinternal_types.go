/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1

import (
	"encoding/json"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"reflect"
)

// NOTE: This file is used to generate the CRDs use by the Operator. The CRD files should not be manually edited
// NOTE: json tags are required. Any new fields you add must have json tags for the fields to be serialized.

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CoherenceInternal is the Schema for the coherenceinternal API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type CoherenceInternal struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec contains the specification for a Coherence cluster. The format is the same
	// as the values file for the Coherence Helm chart.
	Spec   CoherenceInternalSpec   `json:"spec,omitempty"`
	Status CoherenceInternalStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CoherenceInternalList contains a list of CoherenceInternal
type CoherenceInternalList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CoherenceInternal `json:"items"`
}

// CoherenceInternalSpec defines the desired state of CoherenceInternal
// +k8s:openapi-gen=true
type CoherenceInternalSpec struct {
	FullnameOverride string `json:"fullnameOverride,omitempty"`
	NameOverride     string `json:"nameOverride,omitempty"`
	// The cluster name
	Cluster string `json:"cluster"`
	// The name of the headless service used for Coherence WKA
	WKA string `json:"wka,omitempty"`
	// The name to use for the service account to use when RBAC is enabled
	// The role bindings must already have been created as this chart does not create them it just
	// sets the serviceAccountName value in the Pod spec.
	// +optional
	ServiceAccountName string `json:"serviceAccountName,omitempty"`
	// The secrets to be used when pulling images. Secrets must be manually created in the target namespace.
	// ref: https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry/
	// +optional
	ImagePullSecrets  []string `json:"imagePullSecrets,omitempty"`
	CoherenceRoleSpec `json:",inline"`
}

// CoherenceInternalStatus defines the observed state of CoherenceInternal
// +k8s:openapi-gen=true
type CoherenceInternalStatus struct {
}

// NewCoherenceInternalSpec creates a new CoherenceInternalSpec from the specified cluster and role
func NewCoherenceInternalSpec(cluster *CoherenceCluster, role *CoherenceRole) *CoherenceInternalSpec {
	out := CoherenceInternalSpec{}

	out.FullnameOverride = role.Name
	out.Cluster = cluster.Name
	out.ServiceAccountName = cluster.Spec.ServiceAccountName
	out.WKA = cluster.GetWkaServiceName()

	out.CoherenceRoleSpec = CoherenceRoleSpec{}
	role.Spec.DeepCopyInto(&out.CoherenceRoleSpec)

	if cluster.Spec.ImagePullSecrets != nil {
		imagePullSecrets := make([]string, len(cluster.Spec.ImagePullSecrets))
		copy(imagePullSecrets, cluster.Spec.ImagePullSecrets)
		out.ImagePullSecrets = imagePullSecrets
	}

	return &out
}

// NewCoherenceInternalSpecAsMap creates a new CoherenceInternalSpec as a map from the specified cluster and role
func NewCoherenceInternalSpecAsMap(cluster *CoherenceCluster, role *CoherenceRole) (map[string]interface{}, error) {
	spec := NewCoherenceInternalSpec(cluster, role)
	return CoherenceInternalSpecAsMapFromSpec(spec)
}

// CoherenceInternalSpecAsMapFromSpec creates a new CoherenceInternalSpec as a map from the specified and role
func CoherenceInternalSpecAsMapFromSpec(spec *CoherenceInternalSpec) (map[string]interface{}, error) {
	b, _ := json.Marshal(spec)
	jsonMap := make(map[string]interface{})
	err := json.Unmarshal(b, &jsonMap)
	return jsonMap, err
}

// GetCoherenceInternalGroupVersionKind obtains the GroupVersionKind for the CoherenceInternal struct
func GetCoherenceInternalGroupVersionKind(scheme *runtime.Scheme) schema.GroupVersionKind {
	kinds, _, _ := scheme.ObjectKinds(&CoherenceCluster{})

	return schema.GroupVersionKind{
		Group:   kinds[0].Group,
		Version: kinds[0].Version,
		Kind:    reflect.TypeOf(CoherenceInternal{}).Name(),
	}
}
