/*
 * Copyright (c) 2019, 2020 Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package legacy

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NOTE: This file is used to generate the CRDs use by the Operator. The CRD files should not be manually edited
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// CoherenceRole is the Schema for the coherenceroles API
//
//
//
//
//
//
//
//
type CoherenceRole struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CoherenceRoleSpec   `json:"spec,omitempty"`
	Status CoherenceRoleStatus `json:"status,omitempty"`
}

// CoherenceRoleList contains a list of CoherenceRole
type CoherenceRoleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CoherenceRole `json:"items"`
}

// CoherenceRoleStatus defines the observed state of CoherenceRole
type CoherenceRoleStatus struct {
	// The name of the cluster.
	ClusterName string `json:"clusterName,omitempty"`
	// The current status.
	Status RoleStatus `json:"status,omitempty"`
	// Replicas is the desired size of the Coherence cluster.
	Replicas int32 `json:"replicas"`
	// CurrentReplicas is the current size of the Coherence cluster.
	CurrentReplicas int32 `json:"currentReplicas"`
	// ReadyReplicas is the number of Pods created by the StatefulSet.
	ReadyReplicas int32 `json:"readyReplicas"`
	// label query over pods that should match the replicas count. This is same
	// as the label selector but in the string format to avoid introspection
	// by clients. The string will be in the same format as the query-param syntax.
	// More info about label selectors: http://kubernetes.io/docs/user-guide/labels#label-selectors
	Selector string `json:"selector,omitempty"`
	// The status of the start quorums for this role.
	StartQuorum []StartQuorumStatus `json:"startQuorum,omitempty"`
}

// GetCoherenceClusterName obtains the Coherence cluster name from the label for a CoherenceRole.
func (in *CoherenceRole) GetCoherenceClusterName() string {
	if in == nil {
		return ""
	}

	if in.Labels != nil {
		if name, ok := in.Labels[CoherenceClusterLabel]; ok {
			return name
		}
	}

	l := len(in.Name) - len(in.Spec.GetRoleName())
	name := in.Name[0 : l-1]

	return name
}

// RoleStatus is the status value for a CoherenceRoleStatus.
type RoleStatus string

const (
	RoleStatusCreated        RoleStatus = "Created"
	RoleStatusReady          RoleStatus = "Ready"
	RoleStatusScaling        RoleStatus = "Scaling"
	RoleStatusRollingUpgrade RoleStatus = "RollingUpgrade"
	RoleStatusFailed         RoleStatus = "Failed"
	RoleStatusWaiting        RoleStatus = "Waiting"
)
