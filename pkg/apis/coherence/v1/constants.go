/*
 * Copyright (c) 2020, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1

const (
	// The default number of replicas that will be created for a deployment if no value is specified in the spec
	DefaultReplicas int32 = 3
	// The suffix appended to a cluster name to give the WKA service name
	WKAServiceNameSuffix = "-wka"

	// Label keys used to label k8s resources
	LabelCoherenceDeployment = "coherenceDeployment"
	LabelCoherenceCluster    = "coherenceCluster"
	LabelCoherenceRole       = "coherenceRole"
	LabelComponent           = "coherenceComponent"
	LabelPort                = "coherencePort"
	LabelCoherenceWKAMember  = "coherenceWKAMember"

	// Values used for the component label in k8s resources
	LabelComponentCoherenceStatefulSet = "coherence"
	LabelComponentCoherencePod         = "coherencePod"
	LabelComponentCoherenceHeadless    = "coherence-headless"
	LabelComponentEfkConfig            = "coherence-efk-config"
	LabelComponentPVC                  = "coherence-volume"
	LabelComponentPortService          = "coherence-service"
	LabelComponentPortServiceMonitor   = "coherence-service-monitor"
	LabelComponentWKA                  = "coherenceWkaService"

	EfkConfigMapNameTemplate = "%s-efk-config"

	StatusSelectorTemplate = LabelCoherenceCluster + "=%s," + LabelCoherenceDeployment + "=%s"

	// The default k8s service account name.
	DefaultServiceAccount = "default"

	// The affinity topology key for fault domains.
	AffinityTopologyKey = "failure-domain.beta.kubernetes.io/zone"

	// Container Names
	ContainerNameCoherence = "coherence"
	ContainerNameUtils     = "coherence-k8s-utils"
	ContainerNameFluentd   = "fluentd"

	// Volume names
	VolumeNamePersistence     = "persistence-volume"
	VolumeNameSnapshots       = "snapshot-volume"
	VolumeNameLogs            = "log-dir"
	VolumeNameUtils           = "utils-dir"
	VolumeNameJVM             = "jvm"
	VolumeNameFluentdConfig   = "fluentd-coherence-conf"
	VolumeNameFluentdEsConfig = "fluentd-es-config"
	VolumeNameManagementSSL   = "management-ssl-config"
	VolumeNameMetricsSSL      = "metrics-ssl-config"
	VolumeNameLoggingConfig   = "logging-config"

	// Volume mount paths
	VolumeMountPathPersistence       = "/persistence"
	VolumeMountPathSnapshots         = "/snapshot"
	VolumeMountPathUtils             = UtilsDir
	VolumeMountPathJVM               = "/jvm"
	VolumeMountPathLogs              = "/logs"
	VolumeMountPathManagementCerts   = "/coherence/certs/management"
	VolumeMountPathMetricsCerts      = "/coherence/certs/metrics"
	VolumeMountPathLoggingConfig     = "/loggingconfig"
	VolumeMountPathFluentdConfigBase = "/fluentd/etc/"
	VolumeMountSubPathFluentdConfig  = "fluentd-coherence.conf"
	VolumeMountPathFluentdConfig     = VolumeMountPathFluentdConfigBase + VolumeMountSubPathFluentdConfig

	UtilFilesDir     = "/files"
	UtilsDir         = "/utils"
	ScriptsDir       = UtilsDir + "/scripts"
	UtilsInitCommand = UtilFilesDir + "/utils-init"

	ServiceMonitorKind         = "ServiceMonitor"
	ServiceMonitorGroup        = "monitoring.coreos.com"
	ServiceMonitorVersion      = "v1"
	ServiceMonitorGroupVersion = ServiceMonitorGroup + "/" + ServiceMonitorVersion

	// Port names
	PortNameCoherence  = "coherence"
	PortNameDebug      = "debug-port"
	PortNameHealth     = "health"
	PortNameManagement = "management"
	PortNameMetrics    = "metrics"

	DefaultLoggingConfig = ScriptsDir + "/logging.properties"

	DefaultDebugPort      int32 = 5005
	DefaultManagementPort int32 = 30000
	DefaultMetricsPort    int32 = 9612
	DefaultJmxmpPort      int32 = 9099
	DefaultHealthPort     int32 = 6676

	OperatorConfigName = "coherence-operator-config"

	OperatorConfigKeyHost = "operatorhost"

	DefaultReadinessPath = "/ready"
	DefaultLivenessPath  = "/healthz"

	DefaultFluentdImage = "fluent/fluentd-kubernetes-daemonset:v1.3.3-debian-elasticsearch-1.3"

	EnvVarAppType                     = "APP_TYPE"
	EnvVarAppMainClass                = "COH_MAIN_CLASS"
	EnvVarAppMainArgs                 = "COH_MAIN_ARGS"
	EnvVarOperatorHost                = "OPERATOR_HOST"
	EnvVarOperatorTimeout             = "OPERATOR_REQUEST_TIMEOUT"
	EnvVarCoherenceHome               = "COHERENCE_HOME"
	EnvVarCohDependencyModules        = "DEPENDENCY_MODULES"
	EnvVarCohSkipVersionCheck         = "COH_SKIP_VERSION_CHECK"
	EnvVarCohClusterName              = "COH_CLUSTER_NAME"
	EnvVarCohWka                      = "COH_WKA"
	EnvVarCohAppDir                   = "COH_APP_DIR"
	EnvVarCohExtraClassPath           = "COH_EXTRA_CLASSPATH"
	EnvVarCohMachineName              = "COH_MACHINE_NAME"
	EnvVarCohMemberName               = "COH_MEMBER_NAME"
	EnvVarCohPodUID                   = "COH_POD_UID"
	EnvVarCohSite                     = "COH_SITE_INFO_LOCATION"
	EnvVarCohRack                     = "COH_RACK_INFO_LOCATION"
	EnvVarCohRole                     = "COH_ROLE"
	EnvVarCohUtilDir                  = "COH_UTIL_DIR"
	EnvVarCohHealthPort               = "COH_HEALTH_PORT"
	EnvVarCohCacheConfig              = "COH_CACHE_CONFIG"
	EnvVarCohOverride                 = "COH_OVERRIDE_CONFIG"
	EnvVarCohLogLevel                 = "COH_LOG_LEVEL"
	EnvVarCohStorage                  = "COH_STORAGE_ENABLED"
	EnvVarCohPersistenceMode          = "COH_PERSISTENCE_MODE"
	EnvVarCohPersistenceDir           = "COH_PERSISTENCE_DIR"
	EnvVarCohSnapshotDir              = "COH_SNAPSHOT_DIR"
	EnvVarCohLoggingConfig            = "COH_LOGGING_CONFIG"
	EnvVarCohMgmtPrefix               = "COH_MGMT"
	EnvVarCohMetricsPrefix            = "COH_METRICS"
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
	EnvVarJavaHome                    = "JAVA_HOME"
	EnvVarJavaClasspath               = "CLASSPATH"
	EnvVarJvmClasspathJib             = "JVM_USE_JIB_CLASSPATH"
	EnvVarJvmExtraClasspath           = "JVM_EXTRA_CLASSPATH"
	EnvVarJvmExtraModulepath          = "JVM_EXTRA_MODULEPATH"
	EnvVarJvmArgs                     = "JVM_ARGS"
	EnvVarJvmUseContainerLimits       = "JVM_USE_CONTAINER_LIMITS"
	EnvVarJvmFlightRecorder           = "JVM_FLIGHT_RECORDER"
	EnvVarJvmDebugEnabled             = "JVM_DEBUG_ENABLED"
	EnvVarJvmDebugPort                = "JVM_DEBUG_PORT"
	EnvVarJvmDebugSuspended           = "JVM_DEBUG_SUSPEND"
	EnvVarJvmDebugAttach              = "JVM_DEBUG_ATTACH"
	EnvVarJvmGcArgs                   = "JVM_GC_ARGS"
	EnvVarJvmGcCollector              = "JVM_GC_COLLECTOR"
	EnvVarJvmGcLogging                = "JVM_GC_LOGGING"
	EnvVarJvmMemoryHeap               = "JVM_HEAP_SIZE"
	EnvVarJvmMemoryDirect             = "JVM_DIRECT_MEMORY_SIZE"
	EnvVarJvmMemoryStack              = "JVM_STACK_SIZE"
	EnvVarJvmMemoryMeta               = "JVM_METASPACE_SIZE"
	EnvVarJvmMemoryNativeTracking     = "JVM_NATIVE_MEMORY_TRACKING"
	EnvVarJvmOomExit                  = "JVM_OOM_EXIT"
	EnvVarJvmOomHeapDump              = "JVM_OOM_HEAP_DUMP"
	EnvVarJvmJmxmpEnabled             = "JVM_JMXMP_ENABLED"
	EnvVarJvmJmxmpPort                = "JVM_JMXMP_PORT"
	EnvVarFluentdPodID                = "COHERENCE_POD_ID"
	EnvVarFluentdConf                 = "FLUENTD_CONF"
	EnvVarFluentdSedDisable           = "FLUENT_ELASTICSEARCH_SED_DISABLE"
	EnvVarFluentdEsHosts              = "ELASTICSEARCH_HOSTS"
	EnvVarFluentdEsUser               = "ELASTICSEARCH_USER"
	EnvVarFluentdEsCreds              = "ELASTICSEARCH_PASSWORD"
)

const EfkConfig = `# Coherence fluentd configuration
{{- if .Logging.Fluentd }}
{{-   if .Logging.Fluentd.ConfigFileInclude }}
@include {{ .Logging.Fluentd.ConfigFileInclude }}
{{-   end }}
{{- end }}

# Ignore fluentd messages
<match fluent.**>
  @type null
</match>

# Coherence Logs
<source>
  @type tail
  path /logs/coherence-*.log
  pos_file /tmp/cohrence.log.pos
  read_from_head true
  tag coherence-cluster
  multiline_flush_interval 20s
  <parse>
    @type multiline
    format_firstline /^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}.\d{3}/
    format1 /^(?<time>\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}.\d{3})\/(?<uptime>[0-9\.]+) (?<product>.+) <(?<level>[^\s]+)> \(thread=(?<thread>.+), member=(?<member>.+)\):[\S\s](?<log>.*)/
  </parse>
</source>

<filter coherence-cluster>
  @type record_transformer
  <record>
    cluster "{{ .Cluster }}"
    deployment "{{ .DeploymentName }}"
    role "{{ .RoleName }}"
    host "#{ENV['HOSTNAME']}"
    pod-uid "#{ENV['COHERENCE_POD_ID']}"
  </record>
</filter>

<match coherence-cluster>
  @type elasticsearch
  hosts "#{ENV['ELASTICSEARCH_HOSTS']}"
  user "#{ENV['ELASTICSEARCH_USER']}"
  password "#{ENV['ELASTICSEARCH_PASSWORD']}"
  logstash_format true
  logstash_prefix coherence-cluster
{{- if .Logging.Fluentd }}
{{-   if .Logging.Fluentd.SSLVerify }}
  ssl_verify {{ .Logging.Fluentd.SSLVerify }}
{{-   end }}
{{-   if .Logging.Fluentd.SSLVersion }}
  ssl_version {{ .Logging.Fluentd.SSLVersion }}
{{-   end }}
{{-   if .Logging.Fluentd.SSLMinVersion }}
  ssl_min_version {{ .Logging.Fluentd.SSLMinVersion }}
{{-   end }}
{{-   if .Logging.Fluentd.SSLMaxVersion }}
  ssl_max_version {{ .Logging.Fluentd.SSLMaxVersion }}
{{-   end }}
{{- end }}
</match>

{{- if .Logging.Fluentd }}
{{-   if .Logging.Fluentd.Tag }}
<match {{ .Logging.Fluentd.Tag }} >
  @type elasticsearch
  hosts "#{ENV['ELASTICSEARCH_HOSTS']}"
  user "#{ENV['ELASTICSEARCH_USER']}"
  password "#{ENV['ELASTICSEARCH_PASSWORD']}"
  logstash_format true
  logstash_prefix {{ .Logging.Fluentd.Tag }}
{{- if .Logging.Fluentd }}
{{-   if .Logging.Fluentd.SSLVerify }}
  ssl_verify {{ .Logging.Fluentd.SSLVerify }}
{{-   end }}
{{-   if .Logging.Fluentd.SSLVersion }}
  ssl_version {{ .Logging.Fluentd.SSLVersion }}
{{-   end }}
{{-   if .Logging.Fluentd.SSLMinVersion }}
  ssl_min_version {{ .Logging.Fluentd.SSLMinVersion }}
{{-   end }}
{{-   if .Logging.Fluentd.SSLMaxVersion }}
  ssl_max_version {{ .Logging.Fluentd.SSLMaxVersion }}
{{-   end }}
{{- end }}
</match>
{{-   end }}
{{- end }}`
