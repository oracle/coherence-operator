package v1

import (
	v1 "k8s.io/api/core/v1"
)

// Common Coherence types

// NOTE: json tags are required. Any new fields you add must have json tags for the fields to be serialized.

// The default number of replicas that will be created for a role if no value is specified in the spec
const DefaultReplicas = 3

// The suffix appended to a cluster name to give the WKA service name
const WKAServiceNameSuffix = "-wka"

// Images defines the different Docker images used in the role
// +k8s:openapi-gen=true
type Images struct {
	// CoherenceImage is the details of the Coherence image to be used
	// +optional
	Coherence *ImageSpec `json:"coherence,omitempty"`
	// CoherenceUtils is the details of the Coherence utilities image to be used
	// +optional
	CoherenceUtils *ImageSpec `json:"coherenceUtils,omitempty"`
}

// CoherenceInternalImageSpec defines the settings for a Docker image
// +k8s:openapi-gen=true
type ImageSpec struct {
	// Docker image name.
	// More info: https://kubernetes.io/docs/concepts/containers/images
	// +optional
	Image string `json:"image,omitempty"`
	// Image pull policy.
	// One of Always, Never, IfNotPresent.
	// More info: https://kubernetes.io/docs/concepts/containers/images#updating-images
	// +optional
	ImagePullPolicy v1.PullPolicy `json:"imagePullPolicy,omitempty"`
}

// ReadinessProbeSpec defines the settings for the Coherence Pod readiness probe
// +k8s:openapi-gen=true
type ReadinessProbeSpec struct {
	// Number of seconds after the container has started before liveness probes are initiated.
	// More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes
	// +optional
	InitialDelaySeconds *int32 `json:"initialDelaySeconds,omitempty"`
	// Number of seconds after which the probe times out.
	// More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes
	// +optional
	TimeoutSeconds *int32 `json:"timeoutSeconds,omitempty"`
	// How often (in seconds) to perform the probe.
	// +optional
	PeriodSeconds *int32 `json:"periodSeconds,omitempty"`
	// Minimum consecutive successes for the probe to be considered successful after having failed.
	// +optional
	SuccessThreshold *int32 `json:"successThreshold,omitempty"`
	// Minimum consecutive failures for the probe to be considered failed after having succeeded.
	// +optional
	FailureThreshold *int32 `json:"failureThreshold,omitempty"`
}

// UserArtifactsImageSpec defines the settings for the user artifacts image
// +k8s:openapi-gen=true
type UserArtifactsImageSpec struct {
	ImageSpec `json:",inline"`
	// The folder in the custom artifacts Docker image containing jar
	// files to be added to the classpath of the Coherence container.
	// If not set the libDir is "/files/lib".
	LibDir string `json:"libDir,omitempty"`
	// The folder in the custom artifacts Docker image containing
	// configuration files to be added to the classpath of the Coherence container.
	// If not set the configDir is "/files/conf".
	ConfigDir string `json:"configDir,omitempty"`
}

// FluentdImageSpec defines the settings for the fluentd image
// +k8s:openapi-gen=true
type FluentdImageSpec struct {
	ImageSpec `json:",inline"`
	// The fluentd application configuration
	Application *FluentdApplicationSpec `json:"application,omitempty"`
}

// FluentdImageSpec defines the settings for the fluentd application
// +k8s:openapi-gen=true
type FluentdApplicationSpec struct {
	// The fluentd configuration file configuring source for application log.
	ConfigFile string `json:"configFile,omitempty"`
	// This value should be source.tag from fluentd.application.configFile.
	Tag string `json:"tag,omitempty"`
}

// ScalingPolicy describes a policy for scaling a cluster role
type ScalingPolicy string

// Scaling policy constants
const (
	// Safe means that a role will be scaled up or down in a safe manner to ensure no data loss.
	PullAlways ScalingPolicy = "Safe"
	// Parallel means that a role will be scaled up or down by adding or removing members in parallel.
	// If the members of the role are storage enabled then this could cause data loss
	Parallel ScalingPolicy = "Parallel"
	// ParallelUpSafeDown means that a role will be scaled up by adding or removing members in parallel
	// but will be scaled down in a safe manner to ensure no data loss.
	ParallelUpSafeDown ScalingPolicy = "ParallelUpSafeDown"
)

// The key of the label used to hold the Coherence cluster name
const CoherenceClusterLabel string = "coherenceCluster"

// The key of the label used to hold the Coherence role name
const CoherenceRoleLabel string = "coherenceRole"

// The key of the label used to hold the component name
const CoherenceComponentLabel string = "component"
