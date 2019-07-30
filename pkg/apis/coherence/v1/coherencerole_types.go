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
	// +optional
	Role string `json:"role,omitempty"`
	// The desired number of cluster members of this role.
	// This is a pointer to distinguish between explicit zero and not specified.
	// Default value is 3.
	// +optional
	Replicas *int32 `json:"replicas,omitempty"`
	// Details of the Docker images used in the role
	// +optional
	Images *Images `json:"images,omitempty"`
	// A boolean flag indicating whether members of this role are storage enabled.
	// This value will set the corresponding coherence.distributed.localstorage System property.
	// If not specified the default value is true.
	// This flag is also used to configure the ScalingPolicy value if a value is not specified. If the
	// StorageEnabled field is not specified or is true the scaling will be safe, if StorageEnabled is
	// set to false scaling will be parallel.
	// +optional
	StorageEnabled *bool `json:"storageEnabled,omitempty"`
	// ScalingPolicy describes how the replicas of the cluster role will be scaled.
	// The default if not specified is based upon the value of the StorageEnabled field.
	// If StorageEnabled field is not specified or is true the default scaling will be safe, if StorageEnabled is
	// set to false the default scaling will be parallel.
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

// Obtain the full name for  a role.
func (in *CoherenceRoleSpec) GetFullRoleName(cluster *CoherenceCluster) string {
	return cluster.Name + "-" + in.GetRoleName()
}

// Obtain the name for a role.
// If the Role field is not set the default name is returned.
func (in *CoherenceRoleSpec) GetRoleName() string {
	if in == nil {
		return DefaultRoleName
	}
	if in.Role == "" {
		return DefaultRoleName
	}
	return in.Role
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

	l := len(in.Name) - len(in.Spec.Role)
	name := in.Name[0 : l-1]

	return name
}

// DeepCopyWithDefaults returns a copy of this CoherenceRoleSpec with any nil or not set values set
// by the corresponding value in the defaults spec.
func (in *CoherenceRoleSpec) DeepCopyWithDefaults(defaults *CoherenceRoleSpec) *CoherenceRoleSpec {
	if in == nil {
		if defaults != nil {
			return defaults.DeepCopy()
		} else {
			return nil
		}
	}

	if defaults == nil {
		return in.DeepCopy()
	}

	clone := CoherenceRoleSpec{}

	// Copy EVERY field from "in" to the clone.
	// If a field is not set use the value from the default
	// If the field is a struct it should implement DeepCopyWithDefaults so call that method

	if in.Role != "" {
		clone.Role = in.Role
	} else {
		clone.Role = defaults.Role
	}

	if in.Replicas != nil {
		clone.Replicas = in.Replicas
	} else {
		clone.Replicas = defaults.Replicas
	}

	if in.StorageEnabled != nil {
		clone.StorageEnabled = in.StorageEnabled
	} else {
		clone.StorageEnabled = defaults.StorageEnabled
	}

	if in.ScalingPolicy != nil {
		clone.ScalingPolicy = in.ScalingPolicy
	} else {
		clone.ScalingPolicy = defaults.ScalingPolicy
	}

	if in.Labels != nil {
		labels := make(map[string]string)
		for k, v := range *in.Labels {
			labels[k] = v
		}
		clone.Labels = &labels
	} else if defaults.Labels != nil {
		labels := make(map[string]string)
		for k, v := range *defaults.Labels {
			labels[k] = v
		}
		clone.Labels = &labels
	}

	clone.Images = in.Images.DeepCopyWithDefaults(defaults.Images)
	clone.ReadinessProbe = in.ReadinessProbe.DeepCopyWithDefaults(defaults.ReadinessProbe)

	return &clone
}

// RoleStatus is the status value for a CoherenceRoleStatus.
type RoleStatus string

const (
	RoleStatusCreated        RoleStatus = "Created"
	RoleStatusReady          RoleStatus = "Ready"
	RoleStatusScaling        RoleStatus = "Scaling"
	RoleStatusRollingUpgrade RoleStatus = "RollingUpgrade"
	RoleStatusFailed         RoleStatus = "Failed"
)
