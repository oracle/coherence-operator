package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NOTE: This file is used to generate the CRDs use by the Operator. The CRD files should not be manually edited
// NOTE: json tags are required. Any new fields you add must have json tags for the fields to be serialized.

// CoherenceClusterSpec defines the desired state of CoherenceCluster
// +k8s:openapi-gen=true
type CoherenceClusterSpec struct {
	// The secrets to be used when pulling images. Secrets must be manually created in the target namespace.
	// ref: https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry/
	// +optional
	ImagePullSecrets []string `json:"imagePullSecrets,omitempty"`
	// The name to use for the service account to use when RBAC is enabled
	// The role bindings must already have been created as this chart does not create them it just
	// sets the serviceAccountName value in the Pod spec.
	// +optional
	ServiceAccountName string `json:"serviceAccountName,omitempty"`
	// This spec is either the spec of a single role cluster or is used as the
	// default values applied to roles in Roles array.
	CoherenceRoleSpec `json:",inline"`
	// Roles is the list of different roles in the cluster
	// There must be at least one role in a cluster.
	// +optional
	Roles []CoherenceRoleSpec `json:"roles,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CoherenceCluster is the Schema for the coherenceclusters API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
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
}

func init() {
	SchemeBuilder.Register(&CoherenceCluster{}, &CoherenceClusterList{})
}

func (c *CoherenceCluster) GetWkaServiceName() string {
	if c == nil {
		return ""
	}
	return c.Name + WKAServiceNameSuffix
}

// Obtain the CoherenceRoleSpec for the specified role name
func (c *CoherenceCluster) GetRole(name string) CoherenceRoleSpec {
	for _, role := range c.Spec.Roles {
		if role.GetRoleName() == name {
			return role
		}
	}
	return c.Spec.CoherenceRoleSpec
}

// Set the CoherenceRoleSpec
func (c *CoherenceCluster) SetRole(spec CoherenceRoleSpec) {
	for index, role := range c.Spec.Roles {
		if role.GetRoleName() == spec.GetRoleName() {
			c.Spec.Roles[index] = spec
			break
		}
	}
}
