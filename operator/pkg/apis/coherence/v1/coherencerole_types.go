package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NOTE: This file is used to generate the CRDs use by the Operator. The CRD files should not be manually edited
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// CoherenceRoleSpec defines a role in a Coherence cluster. A role is one or
// more Pods that perform the same functionality, for example storage members.
// +k8s:openapi-gen=true
type CoherenceRoleSpec struct {
	// The name of this role.
	// This value will be used to set the Coherence role property for all members of this role
	RoleName string `json:"roleName,omitempty"`
	// The desired number of cluster members of this role.
	// This is a pointer to distinguish between explicit zero and not specified.
	// Default value is 3.
	// +optional
	Replicas *int32 `json:"replicas,omitempty"`
	// Details of the Docker images used in the role
	// +optional
	Images Images `json:"images,omitempty"`
	// ScalingPolicy describes how the replicas of the cluster role will be scaled.
	// The default is ParallelUpSafeDown
	// +optional
	ScalingPolicy *ScalingPolicy `json:"scalingPolicy,omitempty"`
	// The readiness probe config to be used for the Pods in this role.
	// ref: https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-probes/
	// +optional
	ReadinessProbe *ReadinessProbeSpec `json:"readinessProbe,omitempty"`
	// The extra labels to add to the all of the Pods in this roles.
	// Labels here will add to or override those defined for the cluster.
	// More info: http://kubernetes.io/docs/user-guide/labels
	// +optional
	Labels *map[string]string `json:"labels,omitempty"`
}

// CoherenceRoleStatus defines the observed state of CoherenceRole
// +k8s:openapi-gen=true
type CoherenceRoleStatus struct {
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
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CoherenceRole is the Schema for the coherenceroles API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:subresource:scale:specpath=.spec.replicas,statuspath=.status.replicas,selectorpath=.status.selector
type CoherenceRole struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CoherenceRoleSpec   `json:"spec,omitempty"`
	Status CoherenceRoleStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CoherenceRoleList contains a list of CoherenceRole
type CoherenceRoleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CoherenceRole `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CoherenceRole{}, &CoherenceRoleList{})
}

// Obtain the number of replicas required for a role.
// The Replicas field is a pointer and may be nil so this method will
// return either the actual Replica value or the default (DefaultReplicas const)
// if the Replicas field is nil.
func (in *CoherenceRoleSpec) GetReplicas() int32 {
	if in == nil {
		return 0
	}
	if in.Replicas == nil {
		return DefaultReplicas
	}
	return *in.Replicas
}

// GetCoherenceClusterName obtains the Coherence cluster name from the label for a CoherenceRole.
func (in *CoherenceRole) GetCoherenceClusterName() string {
	if in == nil {
		return ""
	}
	if in.Labels == nil {
		return ""
	}
	return in.Labels[CoherenceClusterLabel]
}

type RoleStatus string

const (
	RoleStatusCreated        RoleStatus = "Created"
	RoleStatusReady          RoleStatus = "Ready"
	RoleStatusScaling        RoleStatus = "Scaling"
	RoleStatusRollingUpgrade RoleStatus = "RollingUpgrade"
	RoleStatusFailed         RoleStatus = "Failed"
)