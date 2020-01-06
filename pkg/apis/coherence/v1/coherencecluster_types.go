/*
 * Copyright (c) 2019, 2020 Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NOTE: This file is used to generate the CRDs use by the Operator. The CRD files should not be manually edited
// NOTE: json tags are required. Any new fields you add must have json tags for the fields to be serialized.

// CoherenceClusterSpec defines the desired state of CoherenceCluster
// +k8s:openapi-gen=true
type CoherenceClusterSpec struct {
	// ImagePullSecrets is an optional list of references to secrets in the same namespace to use for pulling any
	// of the images used by this PodSpec.
	// If specified, these secrets will be passed to individual puller implementations for them to use. For example,
	// in the case of docker, only DockerConfig type secrets are honored.
	// More info: https://kubernetes.io/docs/concepts/containers/images#specifying-imagepullsecrets-on-a-pod
	// +listType=map
	// +listMapKey=name
	// +optional
	ImagePullSecrets []LocalObjectReference `json:"imagePullSecrets,omitempty"`
	// The name to use for the service account to use when RBAC is enabled
	// The role bindings must already have been created as this chart does not create them it just
	// sets the serviceAccountName value in the Pod spec.
	// +optional
	ServiceAccountName string `json:"serviceAccountName,omitempty"`
	// Whether or not to auto-mount the Kubernetes API credentials for a service account
	// +optional
	AutomountServiceAccountToken *bool `json:"automountServiceAccountToken,omitempty"`
	// The timeout to apply to rest requests made back to the operator from Coherence Pods.
	// +optional
	OperatorRequestTimeout *int32 `json:"operatorRequestTimeout,omitempty"`
	// This spec is either the spec of a single role cluster or is used as the
	// default values applied to roles in Roles array.
	CoherenceRoleSpec `json:",inline"`
	// Roles is the list of different roles in the cluster
	// There must be at least one role in a cluster.
	// +listType=map
	// +listMapKey=role
	// +optional
	Roles []CoherenceRoleSpec `json:"roles,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CoherenceCluster is the Schema for the coherenceclusters API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:categories=coherence
// +kubebuilder:resource:shortName=coh
// +kubebuilder:printcolumn:name="Roles",type="integer",JSONPath=".status.roles",description="The number of roles in this Coherence cluster"
type CoherenceCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CoherenceClusterSpec   `json:"spec,omitempty"`
	Status CoherenceClusterStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CoherenceClusterList contains a list of CoherenceCluster
type CoherenceClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CoherenceCluster `json:"items"`
}

// CoherenceClusterStatus defines the observed state of CoherenceCluster
// +k8s:openapi-gen=true
type CoherenceClusterStatus struct {
	// The number of roles in this cluster
	Roles int `json:"roles,omitempty"`
}

func init() {
	SchemeBuilder.Register(&CoherenceCluster{}, &CoherenceClusterList{})
}

func (in *CoherenceCluster) GetWkaServiceName() string {
	if in == nil {
		return ""
	}
	return in.Name + WKAServiceNameSuffix
}

// Obtain a map of the CoherenceRoleSpec structs in the cluster.
// These CoherenceRoleSpec instances are copies of those in this
// cluster and not references.
func (in *CoherenceCluster) GetRoles() map[string]CoherenceRoleSpec {
	m := make(map[string]CoherenceRoleSpec)
	if in == nil {
		return m
	}

	if len(in.Spec.Roles) == 0 {
		spec := in.Spec.CoherenceRoleSpec
		m[spec.GetRoleName()] = *spec.DeepCopy()
	} else {
		defaults := in.Spec.CoherenceRoleSpec
		for _, role := range in.Spec.Roles {
			spec := role.DeepCopyWithDefaults(&defaults)
			m[spec.GetRoleName()] = *spec
		}
	}

	return m
}

// Obtain the full name for  a role.
func (in *CoherenceCluster) GetFullRoleName(role string) string {
	if in == nil {
		return role
	}
	return in.Name + "-" + role
}

// Obtain the CoherenceRoleSpec for the first role spec.
// This method is useful to obtain the role from a cluster
// that only has a single role spec.
func (in *CoherenceCluster) GetFirstRole() CoherenceRoleSpec {
	if in == nil {
		return CoherenceRoleSpec{}
	}

	if len(in.Spec.Roles) == 0 {
		return in.Spec.CoherenceRoleSpec
	}
	return in.Spec.Roles[0]
}

// Obtain the CoherenceRoleSpec for the specified role name
func (in *CoherenceCluster) GetRole(name string) CoherenceRoleSpec {
	if len(in.Spec.Roles) > 0 {
		for _, role := range in.Spec.Roles {
			if role.GetRoleName() == name {
				return role
			}
		}
	} else if name == in.Spec.CoherenceRoleSpec.GetRoleName() {
		return in.Spec.CoherenceRoleSpec
	}
	return CoherenceRoleSpec{}
}

// Set the CoherenceRoleSpec
func (in *CoherenceCluster) SetRole(spec CoherenceRoleSpec) {
	name := spec.GetRoleName()
	if len(in.Spec.Roles) > 0 {
		for index, role := range in.Spec.Roles {
			if role.GetRoleName() == name {
				in.Spec.Roles[index] = spec
				break
			}
		}
	} else if name == in.Spec.CoherenceRoleSpec.GetRoleName() {
		in.Spec.CoherenceRoleSpec = spec
	}
}

// Obtain the total number of replicas across all roles in the cluster
func (in *CoherenceCluster) GetClusterSize() int {
	var size = 0
	for _, role := range in.GetRoles() {
		size += int(role.GetReplicas())
	}
	return size
}
