package v1

import (
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
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
	Labels map[string]string `json:"labels,omitempty"`
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
	// Pod scheduling values: Affinity, NodeSelector, Tolerations
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

// Set the number of replicas required for a role.
func (in *CoherenceRoleSpec) SetReplicas(replicas int32) {
	if in != nil {
		in.Replicas = &replicas
	}
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

// DeepCopyWithDefaults returns a copy of this CoherenceRoleSpec with any nil or not set values set
// by the corresponding value in the defaults spec.
func (in *CoherenceRoleSpec) GetEffectiveScalingPolicy() ScalingPolicy {
	if in == nil {
		return SafeScaling
	}

	policy := SafeScaling

	if in.ScalingPolicy == nil {
		// the scaling policy is not set the look at the storage enabled flag
		if in.StorageEnabled == nil || *in.StorageEnabled {
			// storage enabled is either not set or is true so do safe scaling
			policy = ParallelUpSafeDownScaling
		} else {
			// storage enabled is false so do parallel scaling
			policy = ParallelScaling
		}
	} else {
		// scaling policy is set so use it
		policy = *in.ScalingPolicy
	}

	return policy
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

	if in.CacheConfig != nil {
		clone.CacheConfig = in.CacheConfig
	} else {
		clone.CacheConfig = defaults.CacheConfig
	}

	if in.PofConfig != nil {
		clone.PofConfig = in.PofConfig
	} else {
		clone.PofConfig = defaults.PofConfig
	}

	if in.OverrideConfig != nil {
		clone.OverrideConfig = in.OverrideConfig
	} else {
		clone.OverrideConfig = defaults.OverrideConfig
	}

	if in.MaxHeap != nil {
		clone.MaxHeap = in.MaxHeap
	} else {
		clone.MaxHeap = defaults.MaxHeap
	}

	if in.JvmArgs != nil {
		clone.JvmArgs = in.JvmArgs
	} else {
		clone.JvmArgs = defaults.JvmArgs
	}

	if in.JavaOpts != nil {
		clone.JavaOpts = in.JavaOpts
	} else {
		clone.JavaOpts = defaults.JavaOpts
	}

	if in.PodManagementPolicy != nil {
		clone.PodManagementPolicy = in.PodManagementPolicy
	} else {
		clone.PodManagementPolicy = defaults.PodManagementPolicy
	}

	if in.RevisionHistoryLimit != nil {
		clone.RevisionHistoryLimit = in.RevisionHistoryLimit
	} else {
		clone.RevisionHistoryLimit = defaults.RevisionHistoryLimit
	}

	if in.Affinity != nil {
		clone.Affinity = in.Affinity
	} else {
		clone.Affinity = defaults.Affinity
	}

	if in.Resources != nil {
		clone.Resources = in.Resources
	} else {
		clone.Resources = defaults.Resources
	}

	if in.Ports != nil {
		clone.Ports = in.Ports
	} else {
		clone.Ports = defaults.Ports
	}

	if in.Volumes != nil {
		clone.Volumes = make([]corev1.Volume, len(in.Volumes))
		for i := 0; i < len(in.Volumes); i++ {
			clone.Volumes[i] = *in.Volumes[i].DeepCopy()
		}
	} else if defaults.Volumes != nil {
		clone.Volumes = make([]corev1.Volume, len(defaults.Volumes))
		for i := 0; i < len(defaults.Volumes); i++ {
			clone.Volumes[i] = *defaults.Volumes[i].DeepCopy()
		}
	}

	if in.VolumeClaimTemplates != nil {
		clone.VolumeClaimTemplates = make([]corev1.PersistentVolumeClaim, len(in.VolumeClaimTemplates))
		for i := 0; i < len(in.VolumeClaimTemplates); i++ {
			clone.VolumeClaimTemplates[i] = *in.VolumeClaimTemplates[i].DeepCopy()
		}
	} else if defaults.VolumeClaimTemplates != nil {
		clone.VolumeClaimTemplates = make([]corev1.PersistentVolumeClaim, len(defaults.VolumeClaimTemplates))
		for i := 0; i < len(defaults.Volumes); i++ {
			clone.VolumeClaimTemplates[i] = *defaults.VolumeClaimTemplates[i].DeepCopy()
		}
	}

	if in.VolumeMounts != nil {
		clone.VolumeMounts = make([]corev1.VolumeMount, len(in.VolumeMounts))
		for i := 0; i < len(in.VolumeMounts); i++ {
			clone.VolumeMounts[i] = *in.VolumeMounts[i].DeepCopy()
		}
	} else if defaults.VolumeMounts != nil {
		clone.VolumeMounts = make([]corev1.VolumeMount, len(defaults.VolumeMounts))
		for i := 0; i < len(defaults.VolumeMounts); i++ {
			clone.VolumeMounts[i] = *defaults.VolumeMounts[i].DeepCopy()
		}
	}

	if in.Tolerations != nil {
		clone.Tolerations = make([]corev1.Toleration, len(in.Tolerations))
		for i := 0; i < len(in.Tolerations); i++ {
			clone.Tolerations[i] = *in.Tolerations[i].DeepCopy()
		}
	} else if defaults.Tolerations != nil {
		clone.Tolerations = make([]corev1.Toleration, len(defaults.Tolerations))
		for i := 0; i < len(defaults.Tolerations); i++ {
			clone.Tolerations[i] = *defaults.Tolerations[i].DeepCopy()
		}
	}

	clone.Labels = in.mergeMap(in.Labels, defaults.Labels)
	clone.Env = in.mergeMap(in.Env, defaults.Env)
	clone.Annotations = in.mergeMap(in.Annotations, defaults.Annotations)
	clone.NodeSelector = in.mergeMap(in.NodeSelector, defaults.NodeSelector)

	clone.Images = in.Images.DeepCopyWithDefaults(defaults.Images)
	clone.Logging = in.Logging.DeepCopyWithDefaults(defaults.Logging)
	clone.Main = in.Main.DeepCopyWithDefaults(defaults.Main)
	clone.Persistence = in.Persistence.DeepCopyWithDefaults(defaults.Persistence)
	clone.Snapshot = in.Snapshot.DeepCopyWithDefaults(defaults.Snapshot)
	clone.Management = in.Management.DeepCopyWithDefaults(defaults.Management)
	clone.Metrics = in.Metrics.DeepCopyWithDefaults(defaults.Metrics)
	clone.JMX = in.JMX.DeepCopyWithDefaults(defaults.JMX)
	clone.ReadinessProbe = in.ReadinessProbe.DeepCopyWithDefaults(defaults.ReadinessProbe)

	return &clone
}

// Return a map that is two maps merged.
// If both maps are nil then nil is returned.
// Where there are duplicate keys those in m1 take precedence.
// Keys that map to "" will not be added to the merged result
func (in *CoherenceRoleSpec) mergeMap(m1, m2 map[string]string) map[string]string {
	if m1 == nil && m2 == nil {
		return nil
	}

	merged := make(map[string]string)

	if m2 != nil {
		for k, v := range m2 {
			if v != "" {
				merged[k] = v
			}
		}
	}

	if m1 != nil {
		for k, v := range m1 {
			if v != "" {
				merged[k] = v
			} else {
				delete(merged, k)
			}
		}
	}

	return merged
}

// Return a map that is two maps merged.
// If both maps are nil then nil is returned.
// Where there are duplicate keys those in m1 take precedence.
func (in *CoherenceRoleSpec) mergeMapInt32(m1, m2 map[string]int32) map[string]int32 {
	if m1 == nil && m2 == nil {
		return nil
	}

	merged := make(map[string]int32)

	if m2 != nil {
		for k, v := range m2 {
			merged[k] = v
		}
	}

	if m1 != nil {
		for k, v := range m1 {
			merged[k] = v
		}
	}

	return merged
}
