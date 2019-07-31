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
	// The size of the cluster
	ClusterSize int32 `json:"clusterSize,omitempty"`
	// The name of the cluster
	Cluster string `json:"cluster"`
	// The role name of a Coherence cluster member
	Role string `json:"role,omitempty"`
	// The name to use for the service account to use when RBAC is enabled
	// The role bindings must already have been created as this chart does not create them it just
	// sets the serviceAccountName value in the Pod spec.
	ServiceAccountName string `json:"serviceAccountName,omitempty"`
	// The secrets to be used when pulling images. Secrets must be manually created in the target namespace.
	// ref: https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry/
	ImagePullSecrets []string `json:"imagePullSecrets,omitempty"`
	// The Coherence Docker image settings
	Coherence *ImageSpec `json:"coherence,omitempty"`
	// The Coherence Utilities Docker image settings
	CoherenceUtils *ImageSpec `json:"coherenceUtils,omitempty"`
	// The store settings
	Store *CoherenceInternalStoreSpec `json:"store,omitempty"`
	// Controls whether or not log capture via EFK stack is enabled.
	LogCaptureEnabled bool `json:"logCaptureEnabled,omitempty"`
	// Specify the fluentd image
	// These parameters are ignored if 'LogCaptureEnabled' is false.
	Fluentd *FluentdImageSpec `json:"fluentd,omitempty"`
	// The user artifacts image settings
	UserArtifacts *UserArtifactsImageSpec `json:"userArtifacts,omitempty"`
}

// CoherenceInternalStoreSpec defines the desired state of CoherenceInternal stores
// +k8s:openapi-gen=true
type CoherenceInternalStoreSpec struct {
	// A boolean flag indicating whether members of this role are storage enabled.
	// If not specified the default value is true.
	// +optional
	StorageEnabled *bool `json:"storageEnabled,omitempty"`
	// The name of the headless service used for Coherence WKA
	WKA string `json:"wka,omitempty"`
	// The extra labels to add to the Coherence Pod.
	// More info: http://kubernetes.io/docs/user-guide/labels
	// +optional
	Labels map[string]string `json:"labels,omitempty"`
	// The readiness probe config.
	// ref: https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-probes/
	// +optional
	ReadinessProbe *ReadinessProbeSpec `json:"readinessProbe,omitempty"`
}

// CoherenceInternalStatus defines the observed state of CoherenceInternal
// +k8s:openapi-gen=true
type CoherenceInternalStatus struct {
}

// NewCoherenceInternalSpec creates a new CoherenceInternalSpec from the specified cluster and role
func NewCoherenceInternalSpec(cluster *CoherenceCluster, role *CoherenceRole) CoherenceInternalSpec {
	out := CoherenceInternalSpec{}

	out.FullnameOverride = role.Name
	out.ClusterSize = role.Spec.GetReplicas()
	out.Cluster = cluster.Name
	out.ServiceAccountName = cluster.Spec.ServiceAccountName
	out.ImagePullSecrets = cluster.Spec.ImagePullSecrets
	out.Role = role.Spec.GetRoleName()

	// Set the images from the cluster and role
	if role.Spec.Images != nil {
		out.Coherence = role.Spec.Images.Coherence
		out.CoherenceUtils = role.Spec.Images.CoherenceUtils
		out.UserArtifacts = role.Spec.Images.UserArtifacts
		out.Fluentd = role.Spec.Images.Fluentd
	}

	// Set the Store fields
	out.Store = &CoherenceInternalStoreSpec{}
	out.Store.WKA = cluster.GetWkaServiceName()
	out.Store.StorageEnabled = role.Spec.StorageEnabled
	out.Store.ReadinessProbe = role.Spec.ReadinessProbe

	// Set the labels
	labels := make(map[string]string)
	if role.Spec.Labels != nil {
		for k, v := range *role.Spec.Labels {
			labels[k] = v
		}
	}
	// Add the cluster and role labels
	labels[CoherenceClusterLabel] = cluster.Name
	labels[CoherenceRoleLabel] = role.Spec.GetRoleName()

	out.Store.Labels = labels

	return out
}

// NewCoherenceInternalSpecAsMap creates a new CoherenceInternalSpec as a map from the specified cluster and role
func NewCoherenceInternalSpecAsMap(cluster *CoherenceCluster, role *CoherenceRole) (map[string]interface{}, error) {
	spec := NewCoherenceInternalSpec(cluster, role)

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
