/*
 * Copyright (c) 2020, 2023, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1

import "github.com/oracle/coherence-operator/pkg/operator"

const (
	// DefaultReplicas is the default number of replicas that will be created for a deployment if no value is specified in the spec
	DefaultReplicas int32 = 3
	// DefaultJobReplicas is the default number of replicas that will be created for a Job deployment if no value is specified in the spec
	DefaultJobReplicas int32 = 1
	// WKAServiceNameSuffix is the suffix appended to a deployment name to give the WKA service name
	WKAServiceNameSuffix = "-wka"
	// HeadlessServiceNameSuffix is the suffix appended to a deployment name to give the StatefulSet headless-service name
	HeadlessServiceNameSuffix = "-sts"

	// CoherenceFinalizer is the name of the finalizer that the Operator adds to Coherence deployments
	CoherenceFinalizer = "coherence.oracle.com/operator"

	// LabelCoherenceDeployment is the label containing the name of the owning Coherence resource
	LabelCoherenceDeployment = "coherenceDeployment"
	// LabelCoherenceCluster is the label containing the Coherence cluster name
	LabelCoherenceCluster = "coherenceCluster"
	// LabelCoherenceRole is the label containing a Coherence role name
	LabelCoherenceRole = "coherenceRole"
	// LabelComponent is the label containing a component name
	LabelComponent = "coherenceComponent"
	// LabelPort is the label associated to an exposed port name
	LabelPort = "coherencePort"
	// LabelCoherenceWKAMember is the label applied to WKA members
	LabelCoherenceWKAMember = "coherenceWKAMember"
	// LabelApp is an optional application label that can be applied to resources
	LabelApp = "app"
	// LabelVersion is the label containing a resource version
	LabelVersion = "version"
	// LabelCoherenceHash is the label for the Coherence resource spec hash
	LabelCoherenceHash = "coherence-hash"

	// LabelComponentCoherenceStatefulSet is the component label value for a Coherence StatefulSet resource
	LabelComponentCoherenceStatefulSet = "coherence"
	// LabelComponentCoherenceHeadless is the component label value for  a Coherence StatefulSet headless Service resource
	LabelComponentCoherenceHeadless = "coherence-headless"
	// LabelComponentCoherencePod is the component label value for a Coherence Pod
	LabelComponentCoherencePod = "coherencePod"
	// LabelComponentPVC is the component label value for Coherence PersistentVolumeClaim
	LabelComponentPVC = "coherence-volume"
	// LabelComponentPortService is the component label value for a Coherence Service
	LabelComponentPortService = "coherence-service"
	// LabelComponentPortServiceMonitor is the component label value for a Coherence ServiceMonitor
	LabelComponentPortServiceMonitor = "coherence-service-monitor"
	// LabelComponentWKA is the component label value for a Coherence WKA Service
	LabelComponentWKA = "coherenceWkaService"
	// LabelCoherenceStore is the component label value for a Coherence state storage Secret
	LabelCoherenceStore = "coherence-storage"

	// StatusSelectorTemplate is the string template for a WKA service selector
	StatusSelectorTemplate = LabelCoherenceCluster + "=%s," + LabelCoherenceDeployment + "=%s"

	// AnnotationFeatureSuspend is the feature annotations
	AnnotationFeatureSuspend = "com.oracle.coherence.operator/feature.suspend"
	// AnnotationOperatorVersion is the Operator version annotations
	AnnotationOperatorVersion = "com.oracle.coherence.operator/version"

	// DefaultServiceAccount is the default k8s service account name.
	DefaultServiceAccount = "default"

	// ContainerNameCoherence is the Coherence container name
	ContainerNameCoherence = "coherence"
	// ContainerNameOperatorInit is the Operator init-container name
	ContainerNameOperatorInit = "coherence-k8s-utils"

	// VolumeNamePersistence is the name of the persistence volume
	VolumeNamePersistence = "persistence-volume"
	// VolumeNameSnapshots is the name of the snapshots volume
	VolumeNameSnapshots = "snapshot-volume"
	// VolumeNameUtils is the name of the utils volume
	VolumeNameUtils = "coh-utils"
	// VolumeNameJVM is the name of the JVM diagnostics volume
	VolumeNameJVM = "jvm"
	// VolumeNameManagementSSL is the name of the management TLS volume
	VolumeNameManagementSSL = "management-ssl-config"
	// VolumeNameMetricsSSL is the name of the metrics TLS volume
	VolumeNameMetricsSSL = "metrics-ssl-config"

	// VolumePathAttributes is the container attributes file volume
	VolumePathAttributes = "attributes"

	// VolumeMountRoot is the root path for volume mounts
	VolumeMountRoot = "/coherence-operator"
	// VolumeMountPathPersistence is the persistence volume mount
	VolumeMountPathPersistence = VolumeMountRoot + "/persistence"
	// VolumeMountPathSnapshots is the snapshot's volume mount
	VolumeMountPathSnapshots = VolumeMountRoot + "/snapshot"
	// VolumeMountPathUtils is the utils volume mount
	VolumeMountPathUtils = VolumeMountRoot + "/utils"
	// VolumeMountPathJVM is the JVM diagnostics volume mount
	VolumeMountPathJVM = VolumeMountRoot + "/jvm"
	// VolumeMountPathManagementCerts is the management certs volume mount
	VolumeMountPathManagementCerts = VolumeMountRoot + "/coherence/certs/management"
	// VolumeMountPathMetricsCerts is the metrics certs volume mount
	VolumeMountPathMetricsCerts = VolumeMountRoot + "/coherence/certs/metrics"

	// RunnerCommand is the start command for the runner
	RunnerCommand = VolumeMountPathUtils + "/runner"

	// RunnerInitCommand is the start command for the Operator init-container
	RunnerInitCommand = "/files/runner"
	// RunnerInit is the command line argument for the Operator init-container
	RunnerInit = "init"

	// ServiceMonitorKind is the Prometheus ServiceMonitor resource API Kind
	ServiceMonitorKind = "ServiceMonitor"
	// ServiceMonitorGroup is the Prometheus ServiceMonitor resource API Group
	ServiceMonitorGroup = "monitoring.coreos.com"
	// ServiceMonitorVersion is the Prometheus ServiceMonitor resource API version
	ServiceMonitorVersion = "v1"
	// ServiceMonitorGroupVersion is the Prometheus ServiceMonitor resource API group version
	ServiceMonitorGroupVersion = ServiceMonitorGroup + "/" + ServiceMonitorVersion

	// PortNameCoherence is the name of the Coherence port
	PortNameCoherence = "coherence"
	// PortNameDebug is the name of the debug port
	PortNameDebug = "debug-port"
	// PortNameHealth is the name of the health port
	PortNameHealth = "health"
	// PortNameMetrics is the name of the Coherence management port
	PortNameManagement = "management"
	// PortNameMetrics is the name of the Coherence metrics port
	PortNameMetrics = "metrics"

	// DefaultDebugPort is the default debug port
	DefaultDebugPort int32 = 5005
	// DefaultManagementPort is the Coherence manaement debug port
	DefaultManagementPort int32 = 30000
	// DefaultMetricsPort is the default Coherence metrics port
	DefaultMetricsPort int32 = 9612
	// DefaultHealthPort is the default health port
	DefaultHealthPort int32 = 6676
	// DefaultUnicastPort is the default Coherence unicast port
	DefaultUnicastPort int32 = 7575
	// DefaultUnicastPortAdjust is the default Coherence unicast port adjust value
	DefaultUnicastPortAdjust int32 = 7576

	// OperatorConfigName is the Operator configuration Secret name
	OperatorConfigName = "coherence-operator-config"
	// OperatorConfigKeyHost is the key used in the Operator configuration Secret
	OperatorConfigKeyHost = "operatorhost"
	// OperatorSiteURL is the default Operator site query URL
	OperatorSiteURL = "http://$(OPERATOR_HOST)/site/$(COH_MACHINE_NAME)"
	// OperatorRackURL is the default Operator rack query URL
	OperatorRackURL = "http://$(OPERATOR_HOST)/rack/$(COH_MACHINE_NAME)"

	// DefaultReadinessPath is the default readiness endpoint path
	DefaultReadinessPath = "/ready"
	// DefaultLivenessPath is the default liveness endpoint path
	DefaultLivenessPath = "/healthz"

	// DefaultCnbpLauncher is the Cloud Native Build Pack launcher executable
	DefaultCnbpLauncher = "/cnb/lifecycle/launcher"

	EnvVarAppType                     = "COH_APP_TYPE"
	EnvVarAppMainClass                = "COH_MAIN_CLASS"
	EnvVarAppMainArgs                 = "COH_MAIN_ARGS"
	EnvVarOperatorHost                = "OPERATOR_HOST"
	EnvVarOperatorTimeout             = "OPERATOR_REQUEST_TIMEOUT"
	EnvVarOperatorAllowResume         = "OPERATOR_ALLOW_RESUME"
	EnvVarOperatorResumeServices      = "OPERATOR_RESUME_SERVICES"
	EnvVarUseOperatorHealthCheck      = "OPERATOR_HEALTH_CHECK"
	EnvVarCoherenceHome               = "COHERENCE_HOME"
	EnvVarCohDependencyModules        = "DEPENDENCY_MODULES"
	EnvVarCohSkipVersionCheck         = "COH_SKIP_VERSION_CHECK"
	EnvVarCohClusterName              = "COH_CLUSTER_NAME"
	EnvVarCohIdentity                 = "COH_IDENTITY"
	EnvVarCohWka                      = "COH_WKA"
	EnvVarCohAppDir                   = "COH_APP_DIR"
	EnvVarCohMachineName              = "COH_MACHINE_NAME"
	EnvVarCohMemberName               = "COH_MEMBER_NAME"
	EnvVarCohPodUID                   = "COH_POD_UID"
	EnvVarCohSkipSite                 = "COH_SKIP_SITE"
	EnvVarCohSite                     = "COH_SITE_INFO_LOCATION"
	EnvVarCohRack                     = "COH_RACK_INFO_LOCATION"
	EnvVarCohRole                     = "COH_ROLE"
	EnvVarCohUtilDir                  = "COH_UTIL_DIR"
	EnvVarCohUtilLibDir               = "COH_UTIL_LIB_DIR"
	EnvVarCohHealthPort               = "COH_HEALTH_PORT"
	EnvVarCohCacheConfig              = "COH_CACHE_CONFIG"
	EnvVarCohOverride                 = "COH_OVERRIDE_CONFIG"
	EnvVarCohLogLevel                 = "COH_LOG_LEVEL"
	EnvVarCohStorage                  = "COH_STORAGE_ENABLED"
	EnvVarCohPersistenceMode          = "COH_PERSISTENCE_MODE"
	EnvVarCohPersistenceDir           = "COH_PERSISTENCE_DIR"
	EnvVarCohSnapshotDir              = "COH_SNAPSHOT_DIR"
	EnvVarCohTracingRatio             = "COH_TRACING_RATIO"
	EnvVarCohAllowEndangered          = "COH_ALLOW_ENDANGERED"
	EnvVarCohMgmtPrefix               = "COH_MGMT"
	EnvVarCohMetricsPrefix            = "COH_METRICS"
	EnvVarCohEnabledSuffix            = "_ENABLED"
	EnvVarCohPortSuffix               = "_PORT"
	EnvVarCohForceExit                = "COH_FORCE_EXIT"
	EnvVarCoherenceLocalPort          = "COHERENCE_LOCALPORT"
	EnvVarCoherenceLocalPortAdjust    = "COHERENCE_LOCALPORT_ADJUST"
	EnvVarEnableIPMonitor             = "COH_ENABLE_IPMONITOR"
	EnvVarSuffixSSLEnabled            = "_SSL_ENABLED"
	EnvVarSuffixSSLCerts              = "_SSL_CERTS"
	EnvVarSuffixSSLKeyStore           = "_SSL_KEYSTORE"
	EnvVarSuffixSSLKeyStoreCredFile   = "_SSL_KEYSTORE_PASSWORD_FILE"
	EnvVarSuffixSSLKeyCredFile        = "_SSL_KEY_PASSWORD_FILE"
	EnvVarSuffixSSLKeyStoreAlgo       = "_SSL_KEYSTORE_ALGORITHM"
	EnvVarSuffixSSLKeyStoreProvider   = "_SSL_KEYSTORE_PROVIDER"
	EnvVarSuffixSSLKeyStoreType       = "_SSL_KEYSTORE_TYPE"
	EnvVarSuffixSSLTrustStore         = "_SSL_TRUSTSTORE"
	EnvVarSuffixSSLTrustStoreCredFile = "_SSL_TRUSTSTORE_PASSWORD_FILE"
	EnvVarSuffixSSLTrustStoreAlgo     = "_SSL_TRUSTSTORE_ALGORITHM"
	EnvVarSuffixSSLTrustStoreProvider = "_SSL_TRUSTSTORE_PROVIDER"
	EnvVarSuffixSSLTrustStoreType     = "_SSL_TRUSTSTORE_TYPE"
	EnvVarSuffixSSLRequireClientCert  = "_SSL_REQUIRE_CLIENT_CERT"
	EnvVarJavaHome                    = "JAVA_HOME"
	EnvVarJavaClasspath               = "CLASSPATH"
	EnvVarJvmClasspathJib             = "JVM_USE_JIB_CLASSPATH"
	EnvVarJvmExtraClasspath           = "JVM_EXTRA_CLASSPATH"
	EnvVarJvmArgs                     = "JVM_ARGS"
	EnvVarJvmUseContainerLimits       = "JVM_USE_CONTAINER_LIMITS"
	EnvVarJvmShowSettings             = "JVM_SHOW_SETTINGS"
	EnvVarJvmDebugEnabled             = "JVM_DEBUG_ENABLED"
	EnvVarJvmDebugPort                = "JVM_DEBUG_PORT"
	EnvVarJvmDebugSuspended           = "JVM_DEBUG_SUSPEND"
	EnvVarJvmDebugAttach              = "JVM_DEBUG_ATTACH"
	EnvVarJvmGcArgs                   = "JVM_GC_ARGS"
	EnvVarJvmGcCollector              = "JVM_GC_COLLECTOR"
	EnvVarJvmGcLogging                = "JVM_GC_LOGGING"
	EnvVarJvmMemoryHeap               = "JVM_HEAP_SIZE"
	EnvVarJvmMemoryInitialHeap        = "JVM_INITIAL_HEAP_SIZE"
	EnvVarJvmMemoryMaxHeap            = "JVM_MAX_HEAP_SIZE"
	EnvVarJvmMaxRAM                   = "JVM_MAX_RAM"
	EnvVarJvmRAMPercentage            = "JVM_RAM_PERCENTAGE"
	EnvVarJvmInitialRAMPercentage     = "JVM_INITIAL_RAM_PERCENTAGE"
	EnvVarJvmMaxRAMPercentage         = "JVM_MAX_RAM_PERCENTAGE"
	EnvVarJvmMinRAMPercentage         = "JVM_MIN_RAM_PERCENTAGE"
	EnvVarJvmMemoryDirect             = "JVM_DIRECT_MEMORY_SIZE"
	EnvVarJvmMemoryStack              = "JVM_STACK_SIZE"
	EnvVarJvmMemoryMeta               = "JVM_METASPACE_SIZE"
	EnvVarJvmMemoryNativeTracking     = "JVM_NATIVE_MEMORY_TRACKING"
	EnvVarJvmOomExit                  = "JVM_OOM_EXIT"
	EnvVarJvmOomHeapDump              = "JVM_OOM_HEAP_DUMP"
	EnvVarSpringBootFatJar            = "COH_SPRING_BOOT_FAT_JAR"
	EnvVarCnbpEnabled                 = "COH_CNBP_ENABLED"
	EnvVarCnbpLauncher                = "COH_CNBP_LAUNCHER"
)

var (
	// AffinityTopologyKey is the affinity topology key for fault domains.
	AffinityTopologyKey = operator.DefaultSiteLabels[0]
)
