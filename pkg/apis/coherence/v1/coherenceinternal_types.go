/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1

import (
	"encoding/json"
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
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
	// +optional
	ServiceAccountName string `json:"serviceAccountName,omitempty"`
	// The secrets to be used when pulling images. Secrets must be manually created in the target namespace.
	// ref: https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry/
	// +optional
	ImagePullSecrets []string `json:"imagePullSecrets,omitempty"`
	// The Coherence Docker image settings
	// +optional
	Coherence *ImageSpec `json:"coherence,omitempty"`
	// The Coherence Utilities Docker image settings
	// +optional
	CoherenceUtils *ImageSpec `json:"coherenceUtils,omitempty"`
	// The store settings
	// +optional
	Store *CoherenceInternalStoreSpec `json:"store,omitempty"`
	// Affinity controls Pod scheduling preferences.
	//   ref: https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#affinity-and-anti-affinity
	Affinity *corev1.Affinity `json:"affinity,omitempty"`
	// NodeSelector is the Node labels for pod assignment
	//   ref: https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#nodeselector
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// Tolerations is for nodes that have taints on them.
	//   Useful if you want to dedicate nodes to just run the coherence container
	// For example:
	//   tolerations:
	//   - key: "key"
	//     operator: "Equal"
	//     value: "value"
	//     effect: "NoSchedule"
	//
	//   ref: https://kubernetes.io/docs/concepts/configuration/taint-and-toleration/
	Tolerations []corev1.Toleration `json:"tolerations,omitempty"`
	// Resources is the optional resource requests and limits for the containers
	//  ref: http://kubernetes.io/docs/user-guide/compute-resources/
	//
	// By default the cpu requests is set to zero and the cpu limit set to 32. This
	// is because it appears that K8s defaults cpu to one and since Java 10 the JVM
	// now correctly picks up cgroup cpu limits then the JVM will only see one cpu.
	// By setting resources.requests.cpu=0 and resources.limits.cpu=32 it ensures that
	// the JVM will see the either the number of cpus on the host if this is <= 32 or
	// the JVM will see 32 cpus if the host has > 32 cpus. The limit is set to zero
	// so that there is no hard-limit applied.
	//
	// No default memory limits are applied.
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`
	// Controls whether or not log capture via EFK stack is enabled.
	// +optional
	LogCaptureEnabled bool `json:"logCaptureEnabled,omitempty"`
	// Specify the fluentd image
	// These parameters are ignored if 'LogCaptureEnabled' is false.
	// +optional
	Fluentd *FluentdImageSpec `json:"fluentd,omitempty"`
	// The user artifacts image settings
	// +optional
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
	// CacheConfig is the name of the cache configuration file to use
	// +optional
	CacheConfig *string `json:"cacheConfig,omitempty"`
	// PofConfig is the name of the POF configuration file to use when using POF serializer
	// +optional
	PofConfig *string `json:"pofConfig,omitempty"`
	// OverrideConfig is name of the Coherence operational configuration override file,
	// the default is tangosol-coherence-override.xml
	// +optional
	OverrideConfig *string `json:"overrideConfig,omitempty"`
	// Logging allows configuration of Coherence and java util logging.
	// +optional
	Logging *LoggingSpec `json:"logging,omitempty"`
	// Main allows specification of Coherence container main class.
	// +optional
	Main *MainSpec `json:"main,omitempty"`
	// MaxHeap is the min/max heap value to pass to the JVM.
	// The format should be the same as that used for Java's -Xms and -Xmx JVM options.
	// If not set the JVM defaults are used.
	// +optional
	MaxHeap *string `json:"maxHeap,omitempty"`
	// JvmArgs specifies the options to pass to the Coherence JVM. The default is
	// to use the G1 collector.
	// +optional
	JvmArgs *string `json:"jvmArgs,omitempty"`
	// JavaOpts is miscellaneous JVM options to pass to the Coherence store container
	// This options will override the system options computed in the start up script.
	// +optional
	JavaOpts *string `json:"javaOpts,omitempty"`
	// Ports specifies additional port mappings for the Pod and additional Services for those ports
	// +optional
	Ports []NamedPortSpec `json:"ports,omitempty"`
	// Env is additional environment variable mappings that will be passed to
	// the Coherence container in the Pod
	// To specify extra variables add them as name value pairs the same as they
	// would be added to a Pod containers spec, for example these values:
	//
	// env:
	//   FOO: "foo-value"
	//   BAR: "bar-value"
	//
	// will add the environment variable mappings FOO="foo-value" and BAR="bar-value"
	// +optional
	Env map[string]string `json:"env,omitempty"`
	// Annotations are free-form yaml that will be added to the store release as annotations
	// Any annotations should be placed BELOW this annotations: key. For example if we wanted to
	// include annotations for Prometheus it would look like this:
	//
	// annotations:
	//   prometheus.io/scrape: "true"
	//   prometheus.io/port: "2408"
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`
	// PodManagementPolicy sets the podManagementPolicy value for the
	// Coherence cluster StatefulSet.  The default value is Parallel, to
	// cause Pods to be started and stopped in parallel, which can be
	// useful for faster cluster start-up in certain scenarios such as
	// testing but could cause data loss if multiple Pods are stopped in
	// parallel.  This can be changed to OrderedReady which causes Pods
	// to start and stop in sequence.
	// +optional
	PodManagementPolicy *appv1.PodManagementPolicyType `json:"podManagementPolicy,omitempty"`
	// RevisionHistoryLimit is the number of deployment revision K8s keeps after rolling upgrades.
	// The default value if not set is 3.
	// +optional
	RevisionHistoryLimit *int32 `json:"revisionHistoryLimit,omitempty"`
	// Persistence values configure the on-disc data persistence settings.
	// The bool Enabled enables or disabled on disc persistence of data.
	// +optional
	Persistence *PersistentStorageSpec `json:"persistence,omitempty"`
	// Snapshot values configure the on-disc persistence data snapshot (backup) settings.
	// The bool Enabled enables or disabled a different location for
	// persistence snapshot data. If set to false then snapshot files will be written
	// to the same volume configured for persistence data in the Persistence section.
	// +optional
	Snapshot *PersistentStorageSpec `json:"snapshot,omitempty"`
	// Management configures Coherence management over REST
	//   Note: Coherence management over REST will be available in 12.2.1.4.
	// +optional
	Management *PortSpecWithSSL `json:"management,omitempty"`
	// Metrics configures Coherence metrics publishing
	//   Note: Coherence metrics publishing will be available in 12.2.1.4.
	// +optional
	Metrics *PortSpecWithSSL `json:"metrics,omitempty"`
	// JMX defines the values used to enable and configure a separate set of cluster members
	//   that will act as MBean server members and expose a JMX port via a dedicated service.
	//   The JMX port exposed will be using the JMXMP transport as RMI does not work properly in containers.
	// +optional
	JMX *JMXSpec `json:"jmx,omitempty"`
	// Volumes defines extra volume mappings that will be added to the Coherence Pod.
	//   The content of this yaml should match the normal k8s volumes section of a Pod definition
	//   as described in https://kubernetes.io/docs/concepts/storage/volumes/
	// +optional
	Volumes []corev1.Volume `json:"volumes,omitempty"`
	// VolumeClaimTemplates defines extra PVC mappings that will be added to the Coherence Pod.
	//   The content of this yaml should match the normal k8s volumeClaimTemplates section of a Pod definition
	//   as described in https://kubernetes.io/docs/concepts/storage/persistent-volumes/
	// +optional
	VolumeClaimTemplates []corev1.PersistentVolumeClaim `json:"volumeClaimTemplates,omitempty"`
	// VolumeMounts defines extra volume mounts to map to the additional volumes or PVCs declared above
	//   in store.volumes and store.volumeClaimTemplates
	// +optional
	VolumeMounts []corev1.VolumeMount `json:"volumeMounts,omitempty"`
	// The timeout in seconds used by curl when requesting site and rack info.
	// +optional
	CurlTimeout *int `json:"curlTimeout,omitempty"`
}

// CoherenceInternalStatus defines the observed state of CoherenceInternal
// +k8s:openapi-gen=true
type CoherenceInternalStatus struct {
}

// NewCoherenceInternalSpec creates a new CoherenceInternalSpec from the specified cluster and role
func NewCoherenceInternalSpec(cluster *CoherenceCluster, role *CoherenceRole) *CoherenceInternalSpec {
	out := CoherenceInternalSpec{}

	out.FullnameOverride = role.Name
	out.ClusterSize = role.Spec.GetReplicas()
	out.Cluster = cluster.Name
	out.ServiceAccountName = cluster.Spec.ServiceAccountName
	out.Role = role.Spec.GetRoleName()
	out.Affinity = role.Spec.Affinity
	out.Resources = role.Spec.Resources

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
	out.Store.CacheConfig = role.Spec.CacheConfig
	out.Store.PofConfig = role.Spec.PofConfig
	out.Store.OverrideConfig = role.Spec.OverrideConfig
	out.Store.Logging = role.Spec.Logging
	out.Store.Main = role.Spec.Main
	out.Store.MaxHeap = role.Spec.MaxHeap
	out.Store.JvmArgs = role.Spec.JvmArgs
	out.Store.JavaOpts = role.Spec.JavaOpts
	out.Store.PodManagementPolicy = role.Spec.PodManagementPolicy
	out.Store.RevisionHistoryLimit = role.Spec.RevisionHistoryLimit
	out.Store.CurlTimeout = role.Spec.CurlTimeout
	if role.Spec.Persistence != nil {
		out.Store.Persistence = role.Spec.Persistence.DeepCopy()
		if out.Store.Persistence.Volume != nil {
			// override the persistence volume name
			out.Store.Persistence.Volume.Name = "persistence-volume"
		}
	}
	if role.Spec.Snapshot != nil {
		out.Store.Snapshot = role.Spec.Snapshot.DeepCopy()
		if out.Store.Snapshot.Volume != nil {
			// override the snapshot volume name
			out.Store.Snapshot.Volume.Name = "snapshot-volume"
		}
	}
	out.Store.Management = role.Spec.Management
	out.Store.Metrics = role.Spec.Metrics
	out.Store.JMX = role.Spec.JMX
	out.Store.Ports = role.Spec.Ports

	// Set the labels
	labels := make(map[string]string)
	if role.Spec.Labels != nil {
		for k, v := range role.Spec.Labels {
			labels[k] = v
		}
	}
	// Add the cluster and role labels
	labels[CoherenceClusterLabel] = cluster.Name
	labels[CoherenceRoleLabel] = role.Spec.GetRoleName()

	out.Store.Labels = labels

	// Set the Env
	if role.Spec.Env != nil {
		env := make(map[string]string)
		for k, v := range role.Spec.Env {
			env[k] = v
		}
		out.Store.Env = env
	}

	// Set the Annotations
	if role.Spec.Annotations != nil {
		annotations := make(map[string]string)
		for k, v := range role.Spec.Annotations {
			annotations[k] = v
		}
		out.Store.Annotations = annotations
	}

	// Set the NodeSelector
	if role.Spec.NodeSelector != nil {
		nodeSelector := make(map[string]string)
		for k, v := range role.Spec.NodeSelector {
			nodeSelector[k] = v
		}
		out.NodeSelector = nodeSelector
	}

	// Set the ImagePullSecrets
	if cluster.Spec.ImagePullSecrets != nil {
		imagePullSecrets := make([]string, len(cluster.Spec.ImagePullSecrets))
		copy(imagePullSecrets, cluster.Spec.ImagePullSecrets)
		out.ImagePullSecrets = imagePullSecrets
	}
	// Set the Tolerations
	if role.Spec.Tolerations != nil {
		tolerations := make([]corev1.Toleration, len(role.Spec.Tolerations))
		for i := 0; i < len(role.Spec.Tolerations); i++ {
			tolerations[i] = *role.Spec.Tolerations[i].DeepCopy()
		}
		out.Tolerations = tolerations
	}

	// Set the Volumes
	if role.Spec.Volumes != nil {
		volumes := make([]corev1.Volume, len(role.Spec.Volumes))
		for i := 0; i < len(role.Spec.Volumes); i++ {
			volumes[i] = *role.Spec.Volumes[i].DeepCopy()
		}
		out.Store.Volumes = volumes
	}

	// Set the VolumeClaimTemplates
	if role.Spec.VolumeClaimTemplates != nil {
		volumeClaimTemplates := make([]corev1.PersistentVolumeClaim, len(role.Spec.VolumeClaimTemplates))
		for i := 0; i < len(role.Spec.VolumeClaimTemplates); i++ {
			volumeClaimTemplates[i] = *role.Spec.VolumeClaimTemplates[i].DeepCopy()
		}
		out.Store.VolumeClaimTemplates = volumeClaimTemplates
	}

	// Set the VolumeMounts
	if role.Spec.VolumeMounts != nil {
		volumeMounts := make([]corev1.VolumeMount, len(role.Spec.VolumeMounts))
		for i := 0; i < len(role.Spec.VolumeMounts); i++ {
			volumeMounts[i] = *role.Spec.VolumeMounts[i].DeepCopy()
		}
		out.Store.VolumeMounts = volumeMounts
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
