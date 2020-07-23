/*
 * Copyright (c) 2019, 2020 Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package legacy

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

// NOTE: This file is used to generate the CRDs use by the Operator. The CRD files should not be manually edited
// NOTE: json tags are required. Any new fields you add must have json tags for the fields to be serialized.

// CoherenceClusterSpec defines the desired state of CoherenceCluster
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

// CoherenceCluster is the Schema for the coherenceclusters API
//
//
//
//
type CoherenceCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CoherenceClusterSpec   `json:"spec,omitempty"`
	Status CoherenceClusterStatus `json:"status,omitempty"`
}

// CoherenceClusterList contains a list of CoherenceCluster
type CoherenceClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CoherenceCluster `json:"items"`
}

// CoherenceClusterStatus defines the observed state of CoherenceCluster
type CoherenceClusterStatus struct {
	// The number of roles in this cluster
	Roles int32 `json:"roles,omitempty"`
	// The number of roles in this cluster in the Ready state
	Ready int32 `json:"ready,omitempty"`
	// The status of the roles in the cluster
	// +listType=map
	// +listMapKey=role
	RoleStatus []ClusterRoleStatus `json:"roleStatus,omitempty"`
}

// Set the CoherenceRoleSpec
func (in *CoherenceClusterStatus) SetRoleStatus(roleName string, ready bool, pods int32, status RoleStatus) {
	found := false

	if len(in.RoleStatus) > 0 {
		for index, role := range in.RoleStatus {
			if role.Role == roleName {
				s := in.RoleStatus[index]
				s.Transition(ready, pods, status)
				in.RoleStatus[index] = s
				found = true
				break
			}
		}
	}

	if !found {
		s := ClusterRoleStatus{Role: roleName}
		s.Transition(ready, pods, status)
		in.RoleStatus = append(in.RoleStatus, s)
	}

	// update the ready role count
	readyCount := int32(0)
	for _, r := range in.RoleStatus {
		if r.Ready {
			readyCount++
		}
	}
	in.Ready = readyCount
}

// ClusterRoleStatus defines the observed state of role within the cluster
type ClusterRoleStatus struct {
	// The role name
	Role string `json:"role,omitempty"`
	// A flag indicating the role's ready state
	Ready bool `json:"ready,omitempty"`
	// The number of ready Pods.
	Count int32 `json:"count,omitempty"`
	// A status description
	Status RoleStatus `json:"status,omitempty"`
	// The status transitions for the role
	// +optional
	// +listType=map
	// +listMapKey=status
	// +patchMergeKey=type
	// +patchStrategy=merge
	Conditions []ClusterRoleStatusCondition `json:"conditions,omitempty"`
}

func (in *ClusterRoleStatus) Transition(ready bool, count int32, status RoleStatus) {
	in.Ready = ready
	lastStatus := in.Status
	in.Status = status
	in.Count = count

	if lastStatus != status {
		now := metav1.NewTime(time.Now())
		found := false
		for index, condition := range in.Conditions {
			if condition.Status == status {
				condition.LastTransitionTime = now
				in.Conditions[index] = condition
				found = true
				break
			}
		}

		if !found {
			in.Conditions = append(in.Conditions, ClusterRoleStatusCondition{Status: status, LastTransitionTime: now})
		}
	}
}

func (in *ClusterRoleStatus) GetCondition(status RoleStatus) ClusterRoleStatusCondition {
	for _, c := range in.Conditions {
		if c.Status == status {
			return c
		}
	}
	return ClusterRoleStatusCondition{
		Status:             status,
		LastTransitionTime: metav1.Time{},
	}
}

// ClusterRoleStatusCondition defines a specific role status condition
type ClusterRoleStatusCondition struct {
	// The status description
	Status RoleStatus `json:"status,omitempty"`
	// Last time the condition transitioned from one status to another.
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
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
	return CoherenceRoleSpec{Role: name}
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

// Obtain the ClusterRoleStatus for the specified role name
func (in *CoherenceCluster) GetRoleStatus(name string) ClusterRoleStatus {
	if len(in.Status.RoleStatus) > 0 {
		for _, role := range in.Status.RoleStatus {
			if role.Role == name {
				return role
			}
		}
	}
	role := ClusterRoleStatus{Role: name}
	return role
}

// Update the status for a role
func (in *CoherenceCluster) SetRoleStatus(roleName string, ready bool, pods int32, status RoleStatus) {
	if in != nil {
		in.Status.SetRoleStatus(roleName, ready, pods, status)
	}
}
