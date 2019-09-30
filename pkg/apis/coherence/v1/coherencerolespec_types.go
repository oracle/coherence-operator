/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
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
	// The optional application definition
	// +optional
	Application *ApplicationSpec `json:"application,omitempty"`
	// The optional application definition
	// +optional
	Coherence *CoherenceSpec `json:"coherence,omitempty"`
	// The configuration for the Coherence utils image
	// +optional
	CoherenceUtils *ImageSpec `json:"coherenceUtils,omitempty"`
	// Logging allows configuration of Coherence and java util logging.
	// +optional
	Logging *LoggingSpec `json:"logging,omitempty"`
	// The JVM specific options
	// +optional
	JVM *JVMSpec `json:"jvm,omitempty"`
	// Ports specifies additional port mappings for the Pod and additional Services for those ports
	// +optional
	Ports []NamedPortSpec `json:"ports,omitempty"`
	// Env is additional environment variable mappings that will be passed to
	// the Coherence container in the Pod
	// To specify extra variables add them as name value pairs the same as they
	// would be added to a Pod containers spec, for example these values:
	//
	// env:
	//   - name "FOO"
	//     value: "foo-value"
	//   - name: "BAR"
	//     value "bar-value"
	//
	// will add the environment variable mappings FOO="foo-value" and BAR="bar-value"
	// +optional
	Env []corev1.EnvVar `json:"env,omitempty"`
	// The port that the health check endpoint will bind to.
	// +optional
	HealthPort *int32 `json:"healthPort,omitempty"`
	// The readiness probe config to be used for the Pods in this role.
	// ref: https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-probes/
	// +optional
	ReadinessProbe *ReadinessProbeSpec `json:"readinessProbe,omitempty"`
	// The liveness probe config to be used for the Pods in this role.
	// ref: https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-probes/
	// +optional
	LivenessProbe *ReadinessProbeSpec `json:"livenessProbe,omitempty"`
	// The configuration to control safe scaling.
	// +optional
	Scaling *ScalingSpec `json:"scaling,omitempty"`
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
	// +optional
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`
	// Annotations are free-form yaml that will be added to the store release as annotations
	// Any annotations should be placed BELOW this annotations: key. For example if we wanted to
	// include annotations for Prometheus it would look like this:
	//
	// annotations:
	//   prometheus.io/scrape: "true"
	//   prometheus.io/port: "2408"
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`
	// The extra labels to add to the all of the Pods in this roles.
	// Labels here will add to or override those defined for the cluster.
	// More info: http://kubernetes.io/docs/user-guide/labels
	// +optional
	Labels map[string]string `json:"labels,omitempty"`
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
	// Affinity controls Pod scheduling preferences.
	//   ref: https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#affinity-and-anti-affinity
	// +optional
	Affinity *corev1.Affinity `json:"affinity,omitempty"`
	// NodeSelector is the Node labels for pod assignment
	//   ref: https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#nodeselector
	// +optional
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
	// +optional
	Tolerations []corev1.Toleration `json:"tolerations,omitempty"`
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
	if in == nil {
		return ""
	}

	return cluster.GetFullRoleName(in.GetRoleName())
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

func (in *CoherenceRoleSpec) GetEffectiveScalingPolicy() ScalingPolicy {
	if in == nil {
		return SafeScaling
	}

	var policy ScalingPolicy

	if in.Scaling == nil || in.Scaling.Policy == nil {
		// the scaling policy is not set the look at the storage enabled flag
		if in.Coherence.StorageEnabled == nil || *in.Coherence.StorageEnabled {
			// storage enabled is either not set or is true so do safe scaling
			policy = ParallelUpSafeDownScaling
		} else {
			// storage enabled is false so do parallel scaling
			policy = ParallelScaling
		}
	} else {
		// scaling policy is set so use it
		policy = *in.Scaling.Policy
	}

	return policy
}

// Returns the port that the health check endpoint will bind to.
func (in *CoherenceRoleSpec) GetHealthPort() int32 {
	if in == nil || in.HealthPort == nil || *in.HealthPort <= 0 {
		return DefaultHealthPort
	}
	return *in.HealthPort
}

// Returns the ScalingProbe to use for checking Status HA for the role.
// This method will not return nil.
func (in *CoherenceRoleSpec) GetScalingProbe() *ScalingProbe {
	if in == nil || in.Scaling == nil || in.Scaling.Probe == nil {
		return in.GetDefaultScalingProbe()
	}
	return in.Scaling.Probe
}

// Obtain a default ScalingProbe
func (in *CoherenceRoleSpec) GetDefaultScalingProbe() *ScalingProbe {
	timeout := 10

	defaultStatusHA := ScalingProbe{
		TimeoutSeconds: &timeout,
		Handler: corev1.Handler{
			HTTPGet: &corev1.HTTPGetAction{
				Path: "/ha",
				Port: intstr.FromString("health"),
			},
		},
	}

	return defaultStatusHA.DeepCopy()
}

// DeepCopyWithDefaults returns a copy of this CoherenceRoleSpec with any nil or not set values set
// by the corresponding value in the defaults spec.
func (in *CoherenceRoleSpec) DeepCopyWithDefaults(defaults *CoherenceRoleSpec) *CoherenceRoleSpec {
	if in == nil {
		if defaults != nil {
			return defaults.DeepCopy()
		}
		return nil
	}

	if defaults == nil {
		return in.DeepCopy()
	}

	clone := CoherenceRoleSpec{}

	// Copy EVERY field from "in" to the clone.
	// If a field is not set use the value from the default
	// If the field is a struct it should implement DeepCopyWithDefaults so call that method

	// Affinity is NOT merged
	if in.Affinity != nil {
		clone.Affinity = in.Affinity
	} else {
		clone.Affinity = defaults.Affinity
	}

	// Annotations are a map and are merged
	clone.Annotations = in.mergeMap(in.Annotations, defaults.Annotations)
	// Application is merged
	clone.Application = in.Application.DeepCopyWithDefaults(defaults.Application)
	clone.Coherence = in.Coherence.DeepCopyWithDefaults(defaults.Coherence)
	clone.CoherenceUtils = in.CoherenceUtils.DeepCopyWithDefaults(defaults.CoherenceUtils)
	// Environment variables are merged
	clone.Env = in.mergeEnvVar(in.Env, defaults.Env)
	clone.JVM = in.JVM.DeepCopyWithDefaults(defaults.JVM)
	// Labels are a map and are merged
	clone.Labels = in.mergeMap(in.Labels, defaults.Labels)
	clone.Logging = in.Logging.DeepCopyWithDefaults(defaults.Logging)

	// NodeSelector is a map and is NOT merged
	clone.NodeSelector = in.mergeMap(in.NodeSelector, defaults.NodeSelector)
	if in.NodeSelector != nil {
		clone.NodeSelector = in.NodeSelector
	} else {
		clone.NodeSelector = defaults.NodeSelector
	}

	// Ports are named ports in an array and are merged
	if in.Ports != nil {
		clone.Ports = MergeNamedPortSpecs(in.Ports, defaults.Ports)
	} else {
		clone.Ports = defaults.Ports
	}

	// ReadinessProbe is merged
	clone.ReadinessProbe = in.ReadinessProbe.DeepCopyWithDefaults(defaults.ReadinessProbe)

	// Application is NOT merged
	if in.Replicas != nil {
		clone.Replicas = in.Replicas
	} else {
		clone.Replicas = defaults.Replicas
	}

	// Resources is NOT merged
	if in.Resources != nil {
		clone.Resources = in.Resources
	} else {
		clone.Resources = defaults.Resources
	}

	// Role is NOT merged
	if in.Role != "" {
		clone.Role = in.Role
	} else {
		clone.Role = defaults.Role
	}

	// Tolerations is an array but is NOT merged
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

	// VolumeClaimTemplates is an array of named PersistentVolumeClaims and is merged
	clone.VolumeClaimTemplates = in.mergePersistentVolumeClaims(in.VolumeClaimTemplates, defaults.VolumeClaimTemplates)
	// VolumeMounts is an array of named VolumeMounts and is merged
	clone.VolumeMounts = in.mergeVolumeMounts(in.VolumeMounts, defaults.VolumeMounts)
	// Volumes is an array of named VolumeMounts and is merged
	clone.Volumes = in.mergeVolumes(in.Volumes, defaults.Volumes)

	return &clone
}

func (in *CoherenceRoleSpec) mergeEnvVar(primary, secondary []corev1.EnvVar) []corev1.EnvVar {
	if primary == nil {
		return secondary
	}

	if secondary == nil {
		return primary
	}

	if len(primary) == 0 && len(secondary) == 0 {
		return []corev1.EnvVar{}
	}

	var merged []corev1.EnvVar
	merged = append(merged, primary...)

	for _, p := range secondary {
		found := false
		for _, pp := range primary {
			if pp.Name == p.Name {
				found = true
				break
			}
		}

		if !found {
			merged = append(merged, p)
		}
	}

	return merged
}

func (in *CoherenceRoleSpec) mergePersistentVolumeClaims(primary, secondary []corev1.PersistentVolumeClaim) []corev1.PersistentVolumeClaim {
	if primary == nil {
		return secondary
	}

	if secondary == nil {
		return primary
	}

	if len(primary) == 0 && len(secondary) == 0 {
		return []corev1.PersistentVolumeClaim{}
	}

	var merged []corev1.PersistentVolumeClaim
	merged = append(merged, primary...)

	for _, p := range secondary {
		found := false
		for _, pp := range primary {
			if pp.Name == p.Name {
				found = true
				break
			}
		}

		if !found {
			merged = append(merged, p)
		}
	}

	return merged
}

func (in *CoherenceRoleSpec) mergeVolumeMounts(primary, secondary []corev1.VolumeMount) []corev1.VolumeMount {
	if primary == nil {
		return secondary
	}

	if secondary == nil {
		return primary
	}

	if len(primary) == 0 && len(secondary) == 0 {
		return []corev1.VolumeMount{}
	}

	var merged []corev1.VolumeMount
	merged = append(merged, primary...)

	for _, p := range secondary {
		found := false
		for _, pp := range primary {
			if pp.Name == p.Name {
				found = true
				break
			}
		}

		if !found {
			merged = append(merged, p)
		}
	}

	return merged
}

func (in *CoherenceRoleSpec) mergeVolumes(primary, secondary []corev1.Volume) []corev1.Volume {
	if primary == nil {
		return secondary
	}

	if secondary == nil {
		return primary
	}

	if len(primary) == 0 && len(secondary) == 0 {
		return []corev1.Volume{}
	}

	var merged []corev1.Volume
	merged = append(merged, primary...)

	for _, p := range secondary {
		found := false
		for _, pp := range primary {
			if pp.Name == p.Name {
				found = true
				break
			}
		}

		if !found {
			merged = append(merged, p)
		}
	}

	return merged
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

	for k, v := range m2 {
		if v != "" {
			merged[k] = v
		}
	}

	for k, v := range m1 {
		if v != "" {
			merged[k] = v
		} else {
			delete(merged, k)
		}
	}

	return merged
}
