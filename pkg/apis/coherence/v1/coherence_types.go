package v1

import (
	corev1 "k8s.io/api/core/v1"
)

// Common Coherence API structs

// NOTE: This file is used to generate the CRDs use by the Operator. The CRD files should not be manually edited
// NOTE: json tags are required. Any new fields you add must have json tags for the fields to be serialized.

// ----- constants ----------------------------------------------------------

const (
	// The default number of replicas that will be created for a role if no value is specified in the spec
	DefaultReplicas int32 = 3

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
	// UserArtifacts configures the image containing jar files and configuration files
	// that are added to the Coherence JVM's classpath.
	// +optional
	UserArtifacts *UserArtifactsImageSpec `json:"userArtifacts,omitempty"`
	// Fluentd defines the settings for the fluentd image
	// +optional
	Fluentd *FluentdImageSpec `json:"fluentd,omitempty"`
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
	clone.UserArtifacts = in.UserArtifacts.DeepCopyWithDefaults(defaults.UserArtifacts)
	clone.Fluentd = in.Fluentd.DeepCopyWithDefaults(defaults.Fluentd)

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
	ImagePullPolicy *corev1.PullPolicy `json:"imagePullPolicy,omitempty"`
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

// ----- LoggingSpec struct -------------------------------------------------
// LoggingSpec defines the settings for the Coherence Pod logging
// +k8s:openapi-gen=true
type LoggingSpec struct {
	// The default being 5 (info level).
	// +optional
	Level *int32 `json:"level,omitempty"`
	// ConfigFile allows the location of the Java util logging configuration file to be overridden.
	//  If this value is not set the logging.properties file embedded in this chart will be used.
	//  If this value is set the configuration will be located by trying the following locations in order:
	//    1. If store.logging.configMapName is set then the config map will be mounted as a volume and the logging
	//         properties file will be located as a file location relative to the ConfigMap volume mount point.
	//    2. If userArtifacts.imageName is set then using this value as a file name relative to the location of the
	//         configuration files directory in the user artifacts image.
	//    3. Using this value as an absolute file name.
	// +optional
	ConfigFile *string `json:"configFile,omitempty"`
	// ConfigMapName allows a config map to be mounted as a volume containing the logging
	//  configuration file to use.
	// +optional
	ConfigMapName *string `json:"configMapName,omitempty"`
}

// DeepCopyWithDefaults returns a copy of this LoggingSpec struct with any nil or not set values set
// by the corresponding value in the defaults LoggingSpec struct.
func (in *LoggingSpec) DeepCopyWithDefaults(defaults *LoggingSpec) *LoggingSpec {
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

	clone := LoggingSpec{}

	if in.Level != nil {
		clone.Level = in.Level
	} else {
		clone.Level = defaults.Level
	}

	if in.ConfigFile != nil {
		clone.ConfigFile = in.ConfigFile
	} else {
		clone.ConfigFile = defaults.ConfigFile
	}

	if in.ConfigMapName != nil {
		clone.ConfigMapName = in.ConfigMapName
	} else {
		clone.ConfigMapName = defaults.ConfigMapName
	}

	return &clone
}

// ----- MainSpec struct ----------------------------------------------------
// MainSpec defines the specification of Coherence container main class.
// +k8s:openapi-gen=true
type MainSpec struct {
	// Class is the Coherence container main class.  The default value is
	//   com.tangosol.net.DefaultCacheServer.
	// +optional
	Class *string `json:"class,omitempty"`
	// Arguments is the optional arguments for Coherence container main class.
	// +optional
	Arguments *string `json:"arguments,omitempty"`
}

// DeepCopyWithDefaults returns a copy of this MainSpec struct with any nil or not set values set
// by the corresponding value in the defaults MainSpecstruct.
func (in *MainSpec) DeepCopyWithDefaults(defaults *MainSpec) *MainSpec {
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

	clone := MainSpec{}

	if in.Class != nil {
		clone.Class = in.Class
	} else {
		clone.Class = defaults.Class
	}

	if in.Arguments != nil {
		clone.Arguments = in.Arguments
	} else {
		clone.Arguments = defaults.Arguments
	}

	return &clone
}

// ----- PersistentStorageSpec struct ---------------------------------------
// PersistenceStorageSpec defines the persistence settings for the Coherence
// +k8s:openapi-gen=true
type PersistentStorageSpec struct {
	// +optional
	Enabled *bool `json:"enabled,omitempty"`
	// PersistentVolumeClaim allows the configuration of a normal k8s persistent volume claim
	// for persistence data.
	// +optional
	PersistentVolumeClaim *corev1.PersistentVolumeClaimSpec `json:"persistentVolumeClaim,omitempty"` // from k8s.io/api/core/v1
	// Volume allows the configuration of a normal k8s volume mapping
	// for persistence data instead of a persistent volume claim. If a value is defined
	// for store.persistence.volume then no PVC will be created and persistence data
	// will instead be written to this volume. It is up to the deployer to understand
	// the consequences of this and how the guarantees given when using PVCs differ
	// to the storage guarantees for the particular volume type configured here.
	// +optional
	Volume *corev1.Volume `json:"volume,omitempty"` // from k8s.io/api/core/v1
}

// DeepCopyWithDefaults returns a copy of this PersistentStorageSpec struct with any nil or not set values set
// by the corresponding value in the defaults PersistentStorageSpec struct.
func (in *PersistentStorageSpec) DeepCopyWithDefaults(defaults *PersistentStorageSpec) *PersistentStorageSpec {
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

	clone := PersistentStorageSpec{}

	if in.Enabled != nil {
		clone.Enabled = in.Enabled
	} else {
		clone.Enabled = defaults.Enabled
	}

	if in.PersistentVolumeClaim != nil {
		clone.PersistentVolumeClaim = in.PersistentVolumeClaim
	} else {
		clone.PersistentVolumeClaim = defaults.PersistentVolumeClaim
	}

	if in.Volume != nil {
		clone.Volume = in.Volume
	} else {
		clone.Volume = defaults.Volume
	}

	return &clone
}

// ----- SSLSpec struct -----------------------------------------------------
// SSLSpec defines the SSL settings for a Coherence component over REST endpoint.
// +k8s:openapi-gen=true
type SSLSpec struct {
	// Enabled is a boolean flag indicating whether enables or disables SSL on the Coherence management
	// over REST endpoint, the default is false (disabled).
	// +optional
	Enabled *bool `json:"enabled,omitempty"`
	// Secrets is the name of the k8s secrets containing the Java key stores and password files.
	//   This value MUST be provided if SSL is enabled on the Coherence management over ReST endpoint.
	// +optional
	Secrets *string `json:"secrets,omitempty"`
	// Keystore is the name of the Java key store file in the k8s secret to use as the SSL keystore
	//   when configuring component over REST to use SSL.
	// +optional
	KeyStore *string `json:"keyStore,omitemtpy"`
	// KeyStorePasswordFile is the name of the file in the k8s secret containing the keystore
	//   password when configuring component over REST to use SSL.
	// +optional
	KeyStorePasswordFile *string `json:"keyStorePasswordFile,omitempty"`
	// KeyStorePasswordFile is the name of the file in the k8s secret containing the key
	//   password when configuring component over REST to use SSL.
	// +optional
	KeyPasswordFile *string `json:"keyPasswordFile,omitempty"`
	// KeyStoreAlgorithm is the name of the keystore algorithm for the keystore in the k8s secret
	//   used when configuring component over REST to use SSL. If not set the default is SunX509
	// +optional
	KeyStoreAlgorithm *string `json:"keyStoreAlgorithm,omitempty"`
	// KeyStoreProvider is the name of the keystore provider for the keystore in the k8s secret
	//   used when configuring component over REST to use SSL.
	// +optional
	KeyStoreProvider *string `json:"keyStoreProvider,omitempty"`
	// KeyStoreType is the name of the Java keystore type for the keystore in the k8s secret used
	//   when configuring component over REST to use SSL. If not set the default is JKS.
	// +optional
	KeyStoreType *string `json:"keyStoreType,omitempty"`
	// TrustStore is the name of the Java trust store file in the k8s secret to use as the SSL
	//   trust store when configuring component over REST to use SSL.
	// +optional
	TrustStore *string `json:"trustStore,omitempty"`
	// TrustStorePasswordFile is the name of the file in the k8s secret containing the trust store
	//   password when configuring component over REST to use SSL.
	// +optional
	TrustStorePasswordFile *string `json:"trustStorePasswordFile,omitempty"`
	// TrustStoreAlgorithm is the name of the keystore algorithm for the trust store in the k8s
	//   secret used when configuring component over REST to use SSL.  If not set the default is SunX509.
	// +optional
	TrustStoreAlgorithm *string `json:"trustStoreAlgorithm,omitempty"`
	// TrustStoreProvider is the name of the keystore provider for the trust store in the k8s
	//   secret used when configuring component over REST to use SSL.
	// +optional
	TrustStoreProvider *string `json:"trustStoreProvider,omitempty"`
	// TrustStoreType is the name of the Java keystore type for the trust store in the k8s secret
	//   used when configuring component over REST to use SSL. If not set the default is JKS.
	// +optional
	TrustStoreType *string `json:"trustStoreType,omitempty"`
	// RequireClientCert is a boolean flag indicating whether the client certificate will be
	//   authenticated by the server (two-way SSL) when configuring component over REST to use SSL.
	//   If not set the default is false
	// +optional
	RequireClientCert *bool `json:"requireClientCert,omitempty"`
}

// DeepCopyWithDefaults returns a copy of this SSLSpec struct with any nil or not set values set
// by the corresponding value in the defaults SSLSpec struct.
func (in *SSLSpec) DeepCopyWithDefaults(defaults *SSLSpec) *SSLSpec {
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

	clone := SSLSpec{}

	if in.Enabled != nil {
		clone.Enabled = in.Enabled
	} else {
		clone.Enabled = defaults.Enabled
	}

	if in.Secrets != nil {
		clone.Secrets = in.Secrets
	} else {
		clone.Secrets = defaults.Secrets
	}

	if in.KeyStore != nil {
		clone.KeyStore = in.KeyStore
	} else {
		clone.KeyStore = defaults.KeyStore
	}

	if in.KeyStorePasswordFile != nil {
		clone.KeyStorePasswordFile = in.KeyStorePasswordFile
	} else {
		clone.KeyStorePasswordFile = defaults.KeyStorePasswordFile
	}

	if in.KeyPasswordFile != nil {
		clone.KeyPasswordFile = in.KeyPasswordFile
	} else {
		clone.KeyPasswordFile = defaults.KeyPasswordFile
	}

	if in.KeyStoreAlgorithm != nil {
		clone.KeyStoreAlgorithm = in.KeyStoreAlgorithm
	} else {
		clone.KeyStoreAlgorithm = defaults.KeyStoreAlgorithm
	}

	if in.KeyStoreProvider != nil {
		clone.KeyStoreProvider = in.KeyStoreProvider
	} else {
		clone.KeyStoreProvider = defaults.KeyStoreProvider
	}

	if in.KeyStoreType != nil {
		clone.KeyStoreType = in.KeyStoreType
	} else {
		clone.KeyStoreType = defaults.KeyStoreType
	}

	if in.TrustStore != nil {
		clone.TrustStore = in.TrustStore
	} else {
		clone.TrustStore = defaults.TrustStore
	}

	if in.TrustStorePasswordFile != nil {
		clone.TrustStorePasswordFile = in.TrustStorePasswordFile
	} else {
		clone.TrustStorePasswordFile = defaults.TrustStorePasswordFile
	}

	if in.TrustStoreAlgorithm != nil {
		clone.TrustStoreAlgorithm = in.TrustStoreAlgorithm
	} else {
		clone.TrustStoreAlgorithm = defaults.TrustStoreAlgorithm
	}

	if in.TrustStoreProvider != nil {
		clone.TrustStoreProvider = in.TrustStoreProvider
	} else {
		clone.TrustStoreProvider = defaults.TrustStoreProvider
	}

	if in.TrustStoreType != nil {
		clone.TrustStoreType= in.TrustStoreType
	} else {
		clone.TrustStoreType = defaults.TrustStoreType
	}

	if in.RequireClientCert != nil {
		clone.RequireClientCert = in.RequireClientCert
	} else {
		clone.RequireClientCert = defaults.RequireClientCert
	}

	return &clone
}

// ----- PortSpec struct ----------------------------------------------------
// PortSpec defines the port settings for a Coherence component
// +k8s:openapi-gen=true
type PortSpec struct {
	// Port specifies the port used.
	// +optional
	Port *int32 `json:"port,omitempty"`
	// SSL configures SSL settings for a Coherence component
	// +optional
	SSL *SSLSpec `json:"ssl,omitempty"`
}

// DeepCopyWithDefaults returns a copy of this PortSpec struct with any nil or not set values set
// by the corresponding value in the defaults PortSpec struct.
func (in *PortSpec) DeepCopyWithDefaults(defaults *PortSpec) *PortSpec {
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

	clone := PortSpec{}

	if in.Port != nil {
		clone.Port = in.Port
	} else {
		clone.Port = defaults.Port
	}

	if in.SSL != nil {
		clone.SSL = in.SSL
	} else {
		clone.SSL = defaults.SSL
	}

	return &clone
}

// ----- ServiceSpec struct -------------------------------------------------
// ServiceSpec defines the settings for a Service
// +k8s:openapi-gen=true
type ServiceSpec struct {
	// Enabled controls whether to create the service yaml or not
	// +optional
	Enabled *bool `json:"enabled,omitempty"`
	// Type is the K8s service type (typically ClusterIP or LoadBalancer)
	// The default is "ClusterIP".
	// +optional
	Type *string `json:"type,omitempty"`
	// Domain is the external domain name
	// The default is "cluster.local".
	// +optioanl
	Domain *string `json:"domain,omitempty"`
	// LoadBalancerIP is the IP address of the load balancer
	// +optional
	LoadBalancerIP *string `json:"loadBalancerIP,omitempty"`
	// Annotations is free form yaml that will be added to the service annotations
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`
	// The service port value
	// +optional
	ExternalPort *int32 `json:"externalPort,omitempty"`
}

// DeepCopyWithDefaults returns a copy of this ServiceSpec struct with any nil or not set values set
// by the corresponding value in the defaults PortSpec struct.
func (in *ServiceSpec) DeepCopyWithDefaults(defaults *ServiceSpec) *ServiceSpec {
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

	clone := ServiceSpec{}

	if in.Enabled != nil {
		clone.Enabled = in.Enabled
	} else {
		clone.Enabled = defaults.Enabled
	}

	if in.Type != nil {
		clone.Type = in.Type
	} else {
		clone.Type = defaults.Type
	}

	if in.Domain != nil {
		clone.Domain = in.Domain
	} else {
		clone.Domain = defaults.Domain
	}

	if in.LoadBalancerIP != nil {
		clone.LoadBalancerIP = in.LoadBalancerIP
	} else {
		clone.LoadBalancerIP = defaults.LoadBalancerIP
	}

	if in.Annotations != nil {
		clone.Annotations = in.Annotations
	} else {
		clone.Annotations = defaults.Annotations
	}

	if in.ExternalPort != nil {
		clone.ExternalPort = in.ExternalPort
	} else {
		clone.ExternalPort = defaults.ExternalPort
	}

	return &clone
}

// ----- JMXSpec struct -----------------------------------------------------
// JMXSpec defines the values used to enable and configure a separate set of cluster members
//   that will act as MBean server members and expose a JMX port via a dedicated service.
//   The JMX port exposed will be using the JMXMP transport as RMI does not work properly in containers.
// +k8s:openapi-gen=true
type JMXSpec struct {
	// Enabled enables or disables running the MBean server nodes.
	//   If not set the default is false.
	// +optional
	Enabled *bool `json:"enabled,omitempty"`
	// Replicas is the number of MBean server nodes to run.
	//   If not set the default is one.
	// +optional
	Replicas *int32 `json:"replicas,omitempty"`
	// MaxHeap is the min/max heap value to pass to the MBean server JVM.
	//   The format should be the same as that used for Java's -Xms and -Xmx JVM options.
	//   If not set the JVM defaults are used.
	// +optional
	MaxHeap *string `json:"maxHeap,omitempty"`
	// Service groups the values used to configure the management service
	// The default service external port is 9099.
	Service *ServiceSpec `json:"service,omitempty"`
}

// DeepCopyWithDefaults returns a copy of this JMXSpec struct with any nil or not set values set
// by the corresponding value in the defaults PortSpec struct.
func (in *JMXSpec) DeepCopyWithDefaults(defaults *JMXSpec) *JMXSpec {
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

	clone := JMXSpec{}

	if in.Enabled != nil {
		clone.Enabled = in.Enabled
	} else {
		clone.Enabled = defaults.Enabled
	}

	if in.Replicas != nil {
		clone.Replicas = in.Replicas
	} else {
		clone.Replicas = defaults.Replicas
	}

	if in.MaxHeap != nil {
		clone.MaxHeap = in.MaxHeap
	} else {
		clone.MaxHeap = defaults.MaxHeap
	}

	if in.Service != nil {
		clone.Service = in.Service
	} else {
		clone.Service = defaults.Service
	}

	return &clone
}

// ----- CoherenceServiceSpec struct ----------------------------------------
// CoherenceServiceSpec groups the values used to configure the K8s service
// +k8s:openapi-gen=true
type CoherenceServiceSpec struct {
	// The default service external port is 30000.
	ServiceSpec `json:",inline"`
	// The management Http port as integer
	// Default: 30000
	// +optional
	ManagementHttpPort *int32 `json:"managementHttpPort,omitempty"`
	// The metrics http port as integer
	// Default: 9612
	// +optional
	MetricsHttpPort *int32 `json:"metricsHttpPort,omitempty"`
}

// DeepCopyWithDefaults returns a copy of this CoherenceServiceSpec struct with any nil or not set values set
// by the corresponding value in the defaults PortSpec struct.
func (in *CoherenceServiceSpec) DeepCopyWithDefaults(defaults *CoherenceServiceSpec) *CoherenceServiceSpec {
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

	clone := CoherenceServiceSpec{}

	if in.Enabled != nil {
		clone.Enabled = in.Enabled
	} else {
		clone.Enabled = defaults.Enabled
	}

	if in.Type != nil {
		clone.Type = in.Type
	} else {
		clone.Type = defaults.Type
	}

	if in.Domain != nil {
		clone.Domain = in.Domain
	} else {
		clone.Domain = defaults.Domain
	}

	if in.LoadBalancerIP != nil {
		clone.LoadBalancerIP = in.LoadBalancerIP
	} else {
		clone.LoadBalancerIP = defaults.LoadBalancerIP
	}

	if in.Annotations != nil {
		clone.Annotations = in.Annotations
	} else {
		clone.Annotations = defaults.Annotations
	}

	if in.ExternalPort != nil {
		clone.ExternalPort = in.ExternalPort
	} else {
		clone.ExternalPort = defaults.ExternalPort
	}

	if in.ManagementHttpPort != nil {
		clone.ManagementHttpPort = in.ManagementHttpPort
	} else {
		clone.ManagementHttpPort = defaults.ManagementHttpPort
	}

	if in.MetricsHttpPort != nil {
		clone.MetricsHttpPort = in.MetricsHttpPort
	} else {
		clone.MetricsHttpPort = defaults.MetricsHttpPort
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
	// +optional
	LibDir *string `json:"libDir,omitempty"`
	// The folder in the custom artifacts Docker image containing
	// configuration files to be added to the classpath of the Coherence container.
	// If not set the configDir is "/files/conf".
	// +optional
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
	// +optional
	ConfigFile *string `json:"configFile,omitempty"`
	// This value should be source.tag from fluentd.application.configFile.
	// +optional
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
	// ParallelUpSafeDownScaling means that a role will be scaled up by adding or removing members in parallel
	// but will be scaled down in a safe manner to ensure no data loss.
	ParallelUpSafeDownScaling ScalingPolicy = "ParallelUpSafeDownScaling"
)
