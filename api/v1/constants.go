/*
 * Copyright (c) 2020, 2025, Oracle and/or its affiliates.
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
	// AnnotationIstioConfig is the Istio config annotation applied to Pods.
	AnnotationIstioConfig = "proxy.istio.io/config"
	// DefaultIstioConfigAnnotationValue is the default for the istio config annotation.
	// This makes the Istio Sidecar the first container in the Pod to allow it to ideally
	// be started before the Coherence container
	DefaultIstioConfigAnnotationValue = "{ \"holdApplicationUntilProxyStarts\": true }"

	// DefaultServiceAccount is the default k8s service account name.
	DefaultServiceAccount = "default"

	// ContainerNameCoherence is the Coherence container name
	ContainerNameCoherence = "coherence"
	// ContainerNameOperatorInit is the Operator init-container name
	ContainerNameOperatorInit = "coherence-k8s-utils"
	// ContainerNameOperatorConfig is the Operator config files init-container name
	ContainerNameOperatorConfig = "coherence-k8s-config"

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
	// RunnerInit is the command line argument for the Operator intialize init-container
	RunnerInit = "init"
	// RunnerConfig is the command line argument for the Operator config init-container
	RunnerConfig = "config"

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
	// PortNameCoherenceLocal is the name of the Coherence local port
	PortNameCoherenceLocal = "coh-local"
	// PortNameCoherenceCluster is the name of the Coherence cluster port
	PortNameCoherenceCluster = "coh-cluster"
	// PortNameDebug is the name of the debug port
	PortNameDebug = "debug-port"
	// PortNameHealth is the name of the health port
	PortNameHealth = "health"
	// PortNameMetrics is the name of the Coherence management port
	PortNameManagement = "management"
	// PortNameMetrics is the name of the Coherence metrics port
	PortNameMetrics = "metrics"

	// AppProtocolTcp is the appProtocol value for ports that use tcp
	AppProtocolTcp = "tcp"
	// AppProtocolHttp is the appProtocol value for ports that use http
	AppProtocolHttp = "http"

	// DefaultDebugPort is the default debug port
	DefaultDebugPort int32 = 5005
	// DefaultManagementPort is the Coherence manaement debug port
	DefaultManagementPort int32 = 30000
	// DefaultMetricsPort is the default Coherence metrics port
	DefaultMetricsPort int32 = 9612
	// DefaultHealthPort is the default health port
	DefaultHealthPort int32 = 6676
	// DefaultClusterPort is the default Coherence cluster port
	DefaultClusterPort int32 = 7574
	// DefaultUnicastPort is the default Coherence unicast port
	DefaultUnicastPort int32 = 7575
	// DefaultUnicastPortAdjust is the default Coherence unicast port adjust value
	DefaultUnicastPortAdjust int32 = 7576

	// OperatorConfigName is the Operator configuration Secret name
	OperatorConfigName = "coherence-operator-config"
	// OperatorConfigKeyHost is the key used in the Operator configuration Secret
	OperatorConfigKeyHost = "operatorhost"
	// OperatorSiteURL is the default Operator site query URL
	OperatorSiteURL = "http://$(COHERENCE_OPERATOR_HOST)/site/$(COHERENCE_MACHINE)"
	// OperatorRackURL is the default Operator rack query URL
	OperatorRackURL = "http://$(COHERENCE_OPERATOR_HOST)/rack/$(COHERENCE_MACHINE)"

	// OperatorCoherenceArgsFile is the name of the file in the utils directory containing the full set of
	// JVM arguments to run the Coherence container
	OperatorCoherenceArgsFile = "coherence-container-args.txt"
	// OperatorJvmArgsFile is the name of the file in the utils directory containing the JVM arguments
	OperatorJvmArgsFile = "coherence-jvm-args.txt"
	// OperatorClasspathFile is the name of the file in the utils directory containing the JVM class path
	OperatorClasspathFile = "coherence-class-path.txt"
	// OperatorMainClassFile is the name of the file in the utils directory containing the main class name
	OperatorMainClassFile = "coherence-main-class.txt"
	// OperatorSpringBootArgsFile is the name of the file in the utils directory containing the SpringBoot JVM args
	OperatorSpringBootArgsFile = "coherence-spring-args.txt"
	// OperatorJarFileSuffix is the suffix to append to the utils directory to locate the Operator jar file.
	OperatorJarFileSuffix = "/lib/coherence-operator.jar"
	// OperatorConfigDirSuffix is the suffix to append to the utils directory to locate the Operator config directory.
	OperatorConfigDirSuffix = "/config"

	// FileNamePattern is a formatting pattern for a directory separator and file name
	FileNamePattern = "%s%c%s"
	// ArgumentFileNamePattern is a formatting pattern for a JDK argument fle name: directory separator and file name
	ArgumentFileNamePattern = "@" + FileNamePattern

	// DefaultReadinessPath is the default readiness endpoint path
	DefaultReadinessPath = "/ready"
	// DefaultLivenessPath is the default liveness endpoint path
	DefaultLivenessPath = "/healthz"

	// DefaultCnbpLauncher is the Cloud Native Build Pack launcher executable
	DefaultCnbpLauncher = "/cnb/lifecycle/launcher"

	EnvVarAppType                = "COHERENCE_OPERATOR_APP_TYPE"
	EnvVarAppMainClass           = "COHERENCE_OPERATOR_MAIN_CLASS"
	EnvVarAppMainArgs            = "COHERENCE_OPERATOR_MAIN_ARGS"
	EnvVarOperatorHost           = "COHERENCE_OPERATOR_HOST"
	EnvVarOperatorTimeout        = "COHERENCE_OPERATOR_REQUEST_TIMEOUT"
	EnvVarOperatorAllowResume    = "COHERENCE_OPERATOR_ALLOW_RESUME"
	EnvVarOperatorResumeServices = "COHERENCE_OPERATOR_RESUME_SERVICES"
	EnvVarUseOperatorHealthCheck = "COHERENCE_OPERATOR_HEALTH_CHECK"
	EnvVarCohDependencyModules   = "COHERENCE_OPERATOR_DEPENDENCY_MODULES"
	EnvVarCohSkipVersionCheck    = "COHERENCE_OPERATOR_SKIP_VERSION_CHECK"
	EnvVarCohPodUID              = "COHERENCE_OPERATOR_POD_UID"
	EnvVarCohIdentity            = "COHERENCE_OPERATOR_IDENTITY"
	EnvVarCohAppDir              = "COHERENCE_OPERATOR_APP_DIR"
	EnvVarCohSkipSite            = "COHERENCE_OPERATOR_SKIP_SITE"
	EnvVarCohSite                = "COHERENCE_OPERATOR_SITE_INFO_LOCATION"
	EnvVarCohRack                = "COHERENCE_OPERATOR_RACK_INFO_LOCATION"
	EnvVarCohUtilDir             = "COHERENCE_OPERATOR_UTIL_DIR"
	EnvVarCohUtilLibDir          = "COHERENCE_OPERATOR_UTIL_LIB_DIR"
	EnvVarCohAllowEndangered     = "COHERENCE_OPERATOR_ALLOW_ENDANGERED"
	EnvVarSpringBootFatJar       = "COHERENCE_OPERATOR_SPRING_BOOT_FAT_JAR"
	EnvVarCnbpEnabled            = "COHERENCE_OPERATOR_CNBP_ENABLED"
	EnvVarCnbpLauncher           = "COHERENCE_OPERATOR_CNBP_LAUNCHER"
	EnvVarCohForceExit           = "COHERENCE_OPERATOR_FORCE_EXIT"
	EnvVarCohCliProtocol         = "COHERENCE_OPERATOR_CLI_PROTOCOL"

	EnvVarCoherenceHome            = "COHERENCE_HOME"
	EnvVarCohClusterName           = "COHERENCE_CLUSTER"
	EnvVarCohWka                   = "COHERENCE_WKA"
	EnvVarCohMachineName           = "COHERENCE_MACHINE"
	EnvVarCohMemberName            = "COHERENCE_MEMBER"
	EnvVarCoherenceSite            = "COHERENCE_SITE"
	EnvVarCoherenceRack            = "COHERENCE_RACK"
	EnvVarCohRole                  = "COHERENCE_ROLE"
	EnvVarCohHealthPort            = "COHERENCE_HEALTH_HTTP_PORT"
	EnvVarCohCacheConfig           = "COHERENCE_CACHECONFIG"
	EnvVarCohOverride              = "COHERENCE_OVERRIDE"
	EnvVarCohLogLevel              = "COHERENCE_LOG_LEVEL"
	EnvVarCohStorage               = "COHERENCE_DISTRIBUTED_LOCALSTORAGE"
	EnvVarCohPersistenceMode       = "COHERENCE_DISTRIBUTED_PERSISTENCE_MODE"
	EnvVarCohPersistenceDir        = "COHERENCE_DISTRIBUTED_PERSISTENCE_BASE_DIR"
	EnvVarCohSnapshotDir           = "COHERENCE_DISTRIBUTED_PERSISTENCE_SNAPSHOT_DIR"
	EnvVarCohTracingRatio          = "COHERENCE_TRACING_RATIO"
	EnvVarCohMgmtPrefix            = "COHERENCE_MANAGEMENT"
	EnvVarCohMetricsPrefix         = "COHERENCE_METRICS"
	EnvVarCoherenceLocalPort       = "COHERENCE_LOCALPORT"
	EnvVarCoherenceLocalPortAdjust = "COHERENCE_LOCALPORT_ADJUST"
	EnvVarCoherenceTTL             = "COHERENCE_TTL"
	EnvVarEnableIPMonitor          = "COHERENCE_ENABLE_IPMONITOR"
	EnvVarIPMonitorPingTimeout     = "COHERENCE_IPMONITOR_PINGTIMEOUT"

	EnvVarCohCtlHome = "COHCTL_HOME"

	EnvVarCohEnabledSuffix            = "_ENABLED"
	EnvVarCohPortSuffix               = "_PORT"
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

	EnvVarJavaHome                = "JAVA_HOME"
	EnvVarJdkOptions              = "JDK_JAVA_OPTIONS"
	EnvVarJavaClasspath           = "CLASSPATH"
	EnvVarJvmClasspathJib         = "JVM_USE_JIB_CLASSPATH"
	EnvVarJvmExtraClasspath       = "JVM_EXTRA_CLASSPATH"
	EnvVarJvmArgs                 = "JVM_ARGS"
	EnvVarJvmUseContainerLimits   = "JVM_USE_CONTAINER_LIMITS"
	EnvVarJvmShowSettings         = "JVM_SHOW_SETTINGS"
	EnvVarJvmDebugEnabled         = "JVM_DEBUG_ENABLED"
	EnvVarJvmDebugPort            = "JVM_DEBUG_PORT"
	EnvVarJvmDebugSuspended       = "JVM_DEBUG_SUSPEND"
	EnvVarJvmDebugAttach          = "JVM_DEBUG_ATTACH"
	EnvVarJvmGcArgs               = "JVM_GC_ARGS"
	EnvVarJvmGcCollector          = "JVM_GC_COLLECTOR"
	EnvVarJvmGcLogging            = "JVM_GC_LOGGING"
	EnvVarJvmMemoryHeap           = "JVM_HEAP_SIZE"
	EnvVarJvmMemoryInitialHeap    = "JVM_INITIAL_HEAP_SIZE"
	EnvVarJvmMemoryMaxHeap        = "JVM_MAX_HEAP_SIZE"
	EnvVarJvmMaxRAM               = "JVM_MAX_RAM"
	EnvVarJvmRAMPercentage        = "JVM_RAM_PERCENTAGE"
	EnvVarJvmInitialRAMPercentage = "JVM_INITIAL_RAM_PERCENTAGE"
	EnvVarJvmMaxRAMPercentage     = "JVM_MAX_RAM_PERCENTAGE"
	EnvVarJvmMinRAMPercentage     = "JVM_MIN_RAM_PERCENTAGE"
	EnvVarJvmMemoryDirect         = "JVM_DIRECT_MEMORY_SIZE"
	EnvVarJvmMemoryStack          = "JVM_STACK_SIZE"
	EnvVarJvmMemoryMeta           = "JVM_METASPACE_SIZE"
	EnvVarJvmMemoryNativeTracking = "JVM_NATIVE_MEMORY_TRACKING"
	EnvVarJvmOomExit              = "JVM_OOM_EXIT"
	EnvVarJvmOomHeapDump          = "JVM_OOM_HEAP_DUMP"

	SystemPropertyPattern = "-D%s=%s"

	SysPropCoherenceCacheConfig             = "coherence.cacheconfig"
	SysPropCoherenceCluster                 = "coherence.cluster"
	SysPropCoherenceDistributedLocalStorage = "coherence.distributed.localstorage"
	SysPropCoherenceGrpcEnabled             = "coherence.grpc.enabled"
	SysPropCoherenceHealthHttpPort          = "coherence.health.http.port"
	SysPropCoherenceIpMonitor               = "coherence.ipmonitor.pingtimeout"
	SysPropCoherenceLocalPortAdjust         = "coherence.localport.adjust"
	SysPropCoherenceLogLevel                = "coherence.log.level"
	SysPropCoherenceMachine                 = "coherence.machine"
	SysPropCoherenceManagementHttp          = "coherence.management.http"
	SysPropCoherenceManagementHttpPort      = "coherence.management.http.port"
	SysPropCoherenceMember                  = "coherence.member"
	SysPropCoherenceMetricsHttpEnabled      = "coherence.metrics.http.enabled"
	SysPropCoherenceMetricsHttpPort         = "coherence.metrics.http.port"
	SysPropCoherenceOverride                = "coherence.override"
	SysPropCoherencePersistenceBaseDir      = "coherence.distributed.persistence.base.dir"
	SysPropCoherencePersistenceMode         = "coherence.distributed.persistence-mode"
	SysPropCoherencePersistenceSnapshotDir  = "coherence.distributed.persistence.snapshot.dir"
	SysPropCoherenceRole                    = "coherence.role"
	SysPropCoherenceRack                    = "coherence.rack"
	SysPropCoherenceSite                    = "coherence.site"
	SysPropCoherenceTracingRatio            = "coherence.tracing.ratio"
	SysPropCoherenceTTL                     = "coherence.ttl"
	SysPropCoherenceWKA                     = "coherence.wka"

	SysPropOperatorForceExit     = "coherence.operator.force.exit"
	SysPropOperatorHealthEnabled = "coherence.operator.health.enabled"
	SysPropOperatorHealthPort    = "coherence.operator.health.port"
	SysPropOperatorIdentity      = "coherence.operator.identity"
	SysPropOperatorOverride      = "coherence.k8s.override"

	SysPropSpringLoaderMain = "loader.main"
	SysPropSpringLoaderPath = "loader.path"

	JvmOptClassPath                 = "-cp"
	JvmOptUnlockDiagnosticVMOptions = "-XX:+UnlockDiagnosticVMOptions"
	JvmOptNativeMemoryTracking      = "-XX:NativeMemoryTracking"

	// AppTypeNone is the argument to specify no application type.
	AppTypeNone = ""
	// AppTypeJava is the argument to specify a Java application.
	AppTypeJava = "java"
	// AppTypeCoherence is the argument to specify a Coherence application.
	AppTypeCoherence = "coherence"
	// AppTypeHelidon is the argument to specify a Helidon application.
	AppTypeHelidon = "helidon"
	// AppTypeSpring2 is the argument to specify an exploded Spring Boot 2.x application.
	AppTypeSpring2 = "spring"
	// AppTypeSpring3 is the argument to specify an exploded Spring Boot 3.x application.
	AppTypeSpring3 = "spring3"
	// AppTypeOperator is the argument to specify running an Operator command.
	AppTypeOperator = "operator"
	// AppTypeJShell is the argument to specify a JShell application.
	AppTypeJShell = "jshell"

	// DefaultMain is an indicator to run the default main class.
	DefaultMain = "$DEFAULT$"
	// HelidonMain is the default Helidon main class name.
	HelidonMain = "io.helidon.microprofile.cdi.Main"
	// ServerMain is the default server main class name.
	ServerMain = "com.oracle.coherence.k8s.Main"
	// SpringBootMain2 is the default Spring Boot 2.x main class name.
	SpringBootMain2 = "org.springframework.boot.loader.PropertiesLauncher"
	// SpringBootMain3 is the default Spring Boot 3.x main class name.
	SpringBootMain3 = "org.springframework.boot.loader.launch.PropertiesLauncher"
	// ConsoleMain is the Coherence console main class
	ConsoleMain = "com.tangosol.net.CacheFactory"
	// QueryPlusMain is the main class to run Coherence Query Plus
	QueryPlusMain = "com.tangosol.coherence.dslquery.QueryPlus"
	// SleepMain is the main class to run Operator sleep command
	SleepMain = "com.oracle.coherence.k8s.Sleep"
)

var (
	// AffinityTopologyKey is the affinity topology key for fault domains.
	AffinityTopologyKey = operator.DefaultSiteLabels[0]
)
