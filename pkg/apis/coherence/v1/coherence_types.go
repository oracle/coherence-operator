package v1

import (
	v1 "k8s.io/api/core/v1"
)

// Common Coherence types

// NOTE: json tags are required. Any new fields you add must have json tags for the fields to be serialized.

// ----- constants ----------------------------------------------------------

const (
	// The default number of replicas that will be created for a role if no value is specified in the spec
	DefaultReplicas = 3

	// The defaultrole name that will be used for a role if no value is specified in the spec
	DefaultRoleName = "storage"

	// The suffix appended to a cluster name to give the WKA service name
	WKAServiceNameSuffix = "-wka"

	// The key of the label used to hold the Coherence cluster name
	CoherenceClusterLabel string = "coherenceCluster"

	// The key of the label used to hold the Coherence role name
	CoherenceRoleLabel string = "coherenceRole"

	// The key of the label used to hold the component name
	CoherenceComponentLabel string = "component"
)

// ----- Images struct ------------------------------------------------------

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

// DeepCopyWithDefaults returns a copy of this Images struct with any nil or not set values set
// by the corresponding value in the defaults Images struct.
func (in *Images) DeepCopyWithDefaults(defaults *Images) *Images {
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

	clone := Images{}
	clone.Coherence = in.Coherence.DeepCopyWithDefaults(defaults.Coherence)
	clone.CoherenceUtils = in.CoherenceUtils.DeepCopyWithDefaults(defaults.CoherenceUtils)

	return &clone
}

// ----- ImageSpec struct ---------------------------------------------------

// CoherenceInternalImageSpec defines the settings for a Docker image
// +k8s:openapi-gen=true
type ImageSpec struct {
	// Docker image name.
	// More info: https://kubernetes.io/docs/concepts/containers/images
	// +optional
	Image *string `json:"image,omitempty"`
	// Image pull policy.
	// One of Always, Never, IfNotPresent.
	// More info: https://kubernetes.io/docs/concepts/containers/images#updating-images
	// +optional
	ImagePullPolicy *v1.PullPolicy `json:"imagePullPolicy,omitempty"`
}

// DeepCopyWithDefaults returns a copy of this ImageSpec struct with any nil or not set values set
// by the corresponding value in the defaults ImageSpec struct.
func (in *ImageSpec) DeepCopyWithDefaults(defaults *ImageSpec) *ImageSpec {
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

	clone := ImageSpec{}

	if in.Image != nil {
		clone.Image = in.Image
	} else {
		clone.Image = defaults.Image
	}

	if in.ImagePullPolicy != nil {
		clone.ImagePullPolicy = in.ImagePullPolicy
	} else {
		clone.ImagePullPolicy = defaults.ImagePullPolicy
	}

	return &clone
}

// ----- ReadinessProbeSpec struct ------------------------------------------

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

// DeepCopyWithDefaults returns a copy of this ReadinessProbeSpec struct with any nil or not set values set
// by the corresponding value in the defaults ReadinessProbeSpec struct.
func (in *ReadinessProbeSpec) DeepCopyWithDefaults(defaults *ReadinessProbeSpec) *ReadinessProbeSpec {
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

	clone := ReadinessProbeSpec{}

	if in.InitialDelaySeconds != nil {
		clone.InitialDelaySeconds = in.InitialDelaySeconds
	} else {
		clone.InitialDelaySeconds = defaults.InitialDelaySeconds
	}

	if in.TimeoutSeconds != nil {
		clone.TimeoutSeconds = in.TimeoutSeconds
	} else {
		clone.TimeoutSeconds = defaults.TimeoutSeconds
	}

	if in.PeriodSeconds != nil {
		clone.PeriodSeconds = in.PeriodSeconds
	} else {
		clone.PeriodSeconds = defaults.PeriodSeconds
	}

	if in.SuccessThreshold != nil {
		clone.SuccessThreshold = in.SuccessThreshold
	} else {
		clone.SuccessThreshold = defaults.SuccessThreshold
	}

	if in.FailureThreshold != nil {
		clone.FailureThreshold = in.FailureThreshold
	} else {
		clone.FailureThreshold = defaults.FailureThreshold
	}

	return &clone
}

// ----- UserArtifactsImageSpec struct --------------------------------------

// UserArtifactsImageSpec defines the settings for the user artifacts image
// +k8s:openapi-gen=true
type UserArtifactsImageSpec struct {
	ImageSpec `json:",inline"`
	// The folder in the custom artifacts Docker image containing jar
	// files to be added to the classpath of the Coherence container.
	// If not set the libDir is "/files/lib".
	LibDir *string `json:"libDir,omitempty"`
	// The folder in the custom artifacts Docker image containing
	// configuration files to be added to the classpath of the Coherence container.
	// If not set the configDir is "/files/conf".
	ConfigDir *string `json:"configDir,omitempty"`
}

// DeepCopyWithDefaults returns a copy of this UserArtifactsImageSpec struct with any nil or not set values set
// by the corresponding value in the defaults UserArtifactsImageSpec struct.
func (in *UserArtifactsImageSpec) DeepCopyWithDefaults(defaults *UserArtifactsImageSpec) *UserArtifactsImageSpec {
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

	clone := UserArtifactsImageSpec{}

	if in.LibDir != nil {
		clone.LibDir = in.LibDir
	} else {
		clone.LibDir = defaults.LibDir
	}

	if in.ConfigDir != nil {
		clone.ConfigDir = in.ConfigDir
	} else {
		clone.ConfigDir = defaults.ConfigDir
	}

	return &clone
}

// ----- FluentdImageSpec struct --------------------------------------------

// FluentdImageSpec defines the settings for the fluentd image
// +k8s:openapi-gen=true
type FluentdImageSpec struct {
	ImageSpec `json:",inline"`
	// The fluentd application configuration
	Application *FluentdApplicationSpec `json:"application,omitempty"`
}

// DeepCopyWithDefaults returns a copy of this FluentdImageSpec struct with any nil or not set values set
// by the corresponding value in the defaults FluentdImageSpec struct.
func (in *FluentdImageSpec) DeepCopyWithDefaults(defaults *FluentdImageSpec) *FluentdImageSpec {
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

	clone := FluentdImageSpec{}
	clone.ImageSpec = *in.ImageSpec.DeepCopyWithDefaults(&defaults.ImageSpec)
	clone.Application = in.Application.DeepCopyWithDefaults(defaults.Application)

	return &clone
}

// ----- FluentdApplicationSpec struct --------------------------------------

// FluentdImageSpec defines the settings for the fluentd application
// +k8s:openapi-gen=true
type FluentdApplicationSpec struct {
	// The fluentd configuration file configuring source for application log.
	ConfigFile *string `json:"configFile,omitempty"`
	// This value should be source.tag from fluentd.application.configFile.
	Tag *string `json:"tag,omitempty"`
}

// DeepCopyWithDefaults returns a copy of this FluentdApplicationSpec struct with any nil or not set values set
// by the corresponding value in the defaults FluentdApplicationSpec struct.
func (in *FluentdApplicationSpec) DeepCopyWithDefaults(defaults *FluentdApplicationSpec) *FluentdApplicationSpec {
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

	clone := FluentdApplicationSpec{}

	if in.ConfigFile != nil {
		clone.ConfigFile = in.ConfigFile
	} else {
		clone.ConfigFile = defaults.ConfigFile
	}

	if in.Tag != nil {
		clone.Tag = in.Tag
	} else {
		clone.Tag = defaults.Tag
	}

	return &clone
}

// ----- ScalingPolicy type -------------------------------------------------

// ScalingPolicy describes a policy for scaling a cluster role
type ScalingPolicy string

// Scaling policy constants
const (
	// Safe means that a role will be scaled up or down in a safe manner to ensure no data loss.
	SafeScaling ScalingPolicy = "Safe"
	// Parallel means that a role will be scaled up or down by adding or removing members in parallel.
	// If the members of the role are storage enabled then this could cause data loss
	ParallelScaling ScalingPolicy = "Parallel"
	// ParallelUpSafeDown means that a role will be scaled up by adding or removing members in parallel
	// but will be scaled down in a safe manner to ensure no data loss.
	ParallelUpSafeDown ScalingPolicy = "ParallelUpSafeDown"
)
