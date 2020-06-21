/*
 * Copyright (c) 2019, 2020 Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	monitoringv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/go-logr/logr"
	"github.com/go-test/deep"
	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"strconv"
	"strings"
	"time"
)

// Common Coherence API structs

// NOTE: This file is used to generate the CRDs use by the Operator. The CRD files should not be manually edited
// NOTE: json tags are required. Any new fields you add must have json tags for the fields to be serialized.

// ----- helper functions ---------------------------------------------------

func notNilBool(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

func notNilString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func notNilInt32(i *int32) int32 {
	return notNilInt32OrDefault(i, 0)
}

func notNilInt32OrDefault(i *int32, dflt int32) int32 {
	if i == nil {
		return dflt
	}
	return *i
}

// Ensure that the StatefulSet has a container with the specified name
func EnsureContainer(name string, sts *appsv1.StatefulSet) *corev1.Container {
	c := FindContainer(name, sts)
	if c == nil {
		c = &corev1.Container{Name: name}
		sts.Spec.Template.Spec.Containers = append(sts.Spec.Template.Spec.Containers, *c)
	}
	return c
}

// Ensure that the StatefulSet has a container with the specified name
func ReplaceContainer(sts *appsv1.StatefulSet, cNew *corev1.Container) {
	for i, c := range sts.Spec.Template.Spec.Containers {
		if c.Name == cNew.Name {
			sts.Spec.Template.Spec.Containers[i] = *cNew
			return
		}
	}
	sts.Spec.Template.Spec.Containers = append(sts.Spec.Template.Spec.Containers, *cNew)
}

// Find the StatefulSet container with the specified name.
func FindContainer(name string, sts *appsv1.StatefulSet) *corev1.Container {
	for _, c := range sts.Spec.Template.Spec.Containers {
		if c.Name == name {
			return &c
		}
	}
	return nil
}

// Find the StatefulSet init-container with the specified name.
func FindInitContainer(name string, sts *appsv1.StatefulSet) *corev1.Container {
	for _, c := range sts.Spec.Template.Spec.InitContainers {
		if c.Name == name {
			return &c
		}
	}
	return nil
}

// Ensure that the StatefulSet has a volume with the specified name
func ReplaceVolume(sts *appsv1.StatefulSet, volNew corev1.Volume) {
	for i, v := range sts.Spec.Template.Spec.Volumes {
		if v.Name == volNew.Name {
			sts.Spec.Template.Spec.Volumes[i] = volNew
			return
		}
	}
	sts.Spec.Template.Spec.Volumes = append(sts.Spec.Template.Spec.Volumes, volNew)
}

// ----- ApplicationSpec struct ---------------------------------------------

// The specification of the application deployed into the Coherence.
// +k8s:openapi-gen=true
type ApplicationSpec struct {
	// The application type to execute.
	// This field would be set if using the Coherence Graal image and running a none-Java
	// application. For example if the application was a Node application this field
	// would be set to "node". The default is to run a plain Java application.
	// +optional
	Type *string `json:"type,omitempty"`
	// Class is the Coherence container main class.  The default value is
	// com.tangosol.net.DefaultCacheServer.
	// If the application type is non-Java this would be the name of the corresponding language specific
	// runnable, for example if the application type is "node" the main may be a Javascript file.
	// +optional
	Main *string `json:"main,omitempty"`
	// Args is the optional arguments to pass to the main class.
	// +listType=atomic
	// +optional
	Args []string `json:"args,omitempty"`
	// The application folder in the custom artifacts Docker image containing
	// application artifacts.
	// This will effectively become the working directory of the Coherence container.
	// If not set the application directory default value is "/app".
	// +optional
	WorkingDir *string `json:"workingDir,omitempty"`
}

func (in *ApplicationSpec) UpdateCoherenceContainer(c *corev1.Container) {
	if in == nil {
		return
	}

	if in.Type != nil {
		c.Env = append(c.Env, corev1.EnvVar{Name: EnvVarAppType, Value: *in.Type})
	}
	if in.Main != nil {
		c.Env = append(c.Env, corev1.EnvVar{Name: EnvVarAppMainClass, Value: *in.Main})
	}
	if in.WorkingDir != nil {
		c.Env = append(c.Env, corev1.EnvVar{Name: EnvVarCohAppDir, Value: *in.WorkingDir})
	}
	if in.Args != nil && len(in.Args) > 0 {
		args := strings.Join(in.Args, " ")
		c.Env = append(c.Env, corev1.EnvVar{Name: EnvVarAppMainArgs, Value: args})
	}
}

// ----- CoherenceSpec struct -----------------------------------------------

// This section of the CRD configures settings specific to Coherence.
// +coh:doc=coherence_settings/010_overview.adoc,Coherence Configuration
// +k8s:openapi-gen=true
type CoherenceSpec struct {
	// CacheConfig is the name of the cache configuration file to use
	// +coh:doc=coherence_settings/030_cache_config.adoc,Configure Cache Config File
	// +optional
	CacheConfig *string `json:"cacheConfig,omitempty"`
	// OverrideConfig is name of the Coherence operational configuration override file,
	// the default is tangosol-coherence-override.xml
	// +coh:doc=coherence_settings/040_override_file.adoc,Configure Operational Config File
	// +optional
	OverrideConfig *string `json:"overrideConfig,omitempty"`
	// A boolean flag indicating whether members of this deployment are storage enabled.
	// This value will set the corresponding coherence.distributed.localstorage System property.
	// If not specified the default value is true.
	// This flag is also used to configure the ScalingPolicy value if a value is not specified. If the
	// StorageEnabled field is not specified or is true the scaling will be safe, if StorageEnabled is
	// set to false scaling will be parallel.
	// +coh:doc=coherence_settings/050_storage_enabled.adoc,Configure Storage Enabled
	// +optional
	StorageEnabled *bool `json:"storageEnabled,omitempty"`
	// Persistence values configure the on-disc data persistence settings.
	// The bool Enabled enables or disabled on disc persistence of data.
	// +coh:doc=coherence_settings/080_persistence.adoc,Configure Persistence
	// +optional
	Persistence *PersistenceSpec `json:"persistence,omitempty"`
	// The Coherence log level, default being 5 (info level).
	// +coh:doc=coherence_settings/060_log_level.adoc,Configure Coherence log level
	// +optional
	LogLevel *int32 `json:"logLevel,omitempty"`
	// Management configures Coherence management over REST
	// Note: Coherence management over REST will is available in Coherence version >= 12.2.1.4.
	// +coh:doc=management_and_diagnostics/010_overview.adoc,Management & Diagnostics
	// +optional
	Management *PortSpecWithSSL `json:"management,omitempty"`
	// Metrics configures Coherence metrics publishing
	// Note: Coherence metrics publishing will is available in Coherence version >= 12.2.1.4.
	// +coh:doc=metrics/010_overview.adoc,Metrics
	// +optional
	Metrics *PortSpecWithSSL `json:"metrics,omitempty"`
	// Tracing is used to configure Coherence distributed tracing functionality.
	// +optional
	Tracing *CoherenceTracingSpec `json:"tracing,omitempty"`
	// AllowEndangeredForStatusHA is a list of Coherence partitioned cache service names
	// that are allowed to be in an endangered state when testing for StatusHA.
	// Instances where a StatusHA check is performed include the readiness probe and when
	// scaling a deployment.
	// This field would not typically be used except in cases where a cache service is
	// configured with a backup count greater than zero but it does not matter if caches in
	// those services loose data due to member departure. Normally, such cache services would
	// have a backup count of zero, which would automatically excluded them from the StatusHA
	// check.
	// +optional
	AllowEndangeredForStatusHA []string `json:"allowEndangeredForStatusHA,omitempty"`
	// Exclude members of this deployment from being part of the cluster's WKA list.
	// +coh:doc=coherence_settings/070_wka.adoc,Well Known Addressing
	// +optional
	ExcludeFromWKA *bool `json:"excludeFromWKA,omitempty"`
	// Specify an existing Coherence deployment to be used for WKA.
	// If an existing deployment is to be used for WKA the ExcludeFromWKA is
	// implicitly set to true.
	// +coh:doc=coherence_settings/070_wka.adoc,Well Known Addressing
	// +optional
	WKA *CoherenceWKASpec `json:"wka,omitempty"`
	// Certain features rely on a version check prior to starting the server, e.g. metrics requires >= 12.2.1.4.
	// The version check relies on the ability of the start script to find coherence.jar but if due to how the image
	// has been built this check is failing then setting this flag to true will skip version checking and assume
	// that the latest coherence.jar is being used.
	// +optional
	SkipVersionCheck *bool `json:"skipVersionCheck,omitempty"`
}

// IsWKAMember returns true if this deployment is a WKA list member.
func (in *CoherenceSpec) IsWKAMember() bool {
	if in != nil && in.ExcludeFromWKA != nil && *in.ExcludeFromWKA {
		return false
	}
	if in != nil && in.WKA != nil && in.WKA.Deployment != "" {
		return false
	}
	return true
}

// RequiresWKAService returns true if this deployment requires a WKA Service.
func (in *CoherenceSpec) RequiresWKAService() bool {
	if in != nil && in.WKA != nil && in.WKA.Deployment != "" {
		return false
	}
	return true
}

// GetWKA returns the host name Coherence should for WKA.
func (in *CoherenceSpec) GetWKA(deployment string) string {
	if in == nil || in.WKA == nil || in.WKA.Deployment == "" {
		// there is no WKA override so return the deployment name
		return deployment + WKAServiceNameSuffix
	}

	if in.WKA.Namespace != "" {
		// A WKA override is specified with a namespace
		return fmt.Sprintf("%s%s.%s.svc.cluster.local", in.WKA.Deployment, WKAServiceNameSuffix, in.WKA.Namespace)
	}

	// A WKA override is specified without a namespace
	return in.WKA.Deployment + WKAServiceNameSuffix
}

// Add the persistence and snapshot volume mounts to the specified container
func (in *CoherenceSpec) AddPersistenceVolumeMounts(c *corev1.Container) {
	if in != nil {
		in.Persistence.AddVolumeMounts(c)
	}
}

// Add the persistence and snapshot persistent volume claims
func (in *CoherenceSpec) AddPersistencePVCs(deployment *Coherence, sts *appsv1.StatefulSet) {
	// Add the persistence PVC if required
	pvcs := in.Persistence.CreatePersistentVolumeClaims(deployment)
	sts.Spec.VolumeClaimTemplates = append(sts.Spec.VolumeClaimTemplates, pvcs...)
}

// Add the persistence and snapshot volumes
func (in *CoherenceSpec) AddPersistenceVolumes(sts *appsv1.StatefulSet) {
	// Add the persistence volume if required
	vols := in.Persistence.CreatePersistenceVolumes()
	sts.Spec.Template.Spec.Volumes = append(sts.Spec.Template.Spec.Volumes, vols...)
}

// Apply Coherence settings to the StatefulSet.
func (in *CoherenceSpec) UpdateStatefulSet(deployment *Coherence, sts *appsv1.StatefulSet) {
	// Get the Coherence container
	c := EnsureContainer(ContainerNameCoherence, sts)
	defer ReplaceContainer(sts, c)

	if in == nil {
		// we're nil so disable management and metrics/
		c.Env = append(c.Env, corev1.EnvVar{Name: EnvVarCohMgmtPrefix + EnvVarCohEnabledSuffix, Value: "false"},
			corev1.EnvVar{Name: EnvVarCohMetricsPrefix + EnvVarCohEnabledSuffix, Value: "false"})
		return
	}

	if in.CacheConfig != nil && *in.CacheConfig != "" {
		c.Env = append(c.Env, corev1.EnvVar{Name: EnvVarCohCacheConfig, Value: *in.CacheConfig})
	}

	if in.OverrideConfig != nil && *in.OverrideConfig != "" {
		c.Env = append(c.Env, corev1.EnvVar{Name: EnvVarCohOverride, Value: *in.OverrideConfig})
	}

	if in.LogLevel != nil {
		c.Env = append(c.Env, corev1.EnvVar{Name: EnvVarCohLogLevel, Value: Int32PtrToString(in.LogLevel)})
	}

	if in.StorageEnabled != nil {
		c.Env = append(c.Env, corev1.EnvVar{Name: EnvVarCohStorage, Value: BoolPtrToString(in.StorageEnabled)})
	}

	if in.SkipVersionCheck != nil {
		c.Env = append(c.Env, corev1.EnvVar{Name: EnvVarCohSkipVersionCheck, Value: BoolPtrToString(in.SkipVersionCheck)})
	}

	if in.Tracing != nil && in.Tracing.Ratio != nil {
		c.Env = append(c.Env, corev1.EnvVar{Name: EnvVarCohTracingRatio, Value: *in.Tracing.Ratio})
	}

	if len(in.AllowEndangeredForStatusHA) != 0 {
		c.Env = append(c.Env, corev1.EnvVar{Name: EnvVarCohAllowEndangered, Value: strings.Join(in.AllowEndangeredForStatusHA, ",")})
	}

	in.Management.AddSSLVolumes(sts, c, VolumeNameManagementSSL, VolumeMountPathManagementCerts)
	c.Env = append(c.Env, in.Management.CreateEnvVars(EnvVarCohMgmtPrefix, VolumeMountPathManagementCerts, DefaultManagementPort)...)

	in.Metrics.AddSSLVolumes(sts, c, VolumeNameMetricsSSL, VolumeMountPathMetricsCerts)
	c.Env = append(c.Env, in.Metrics.CreateEnvVars(EnvVarCohMetricsPrefix, VolumeMountPathMetricsCerts, DefaultMetricsPort)...)

	// set the persistence mode
	if mode := in.Persistence.GetMode(); mode != nil {
		c.Env = append(c.Env, corev1.EnvVar{Name: EnvVarCohPersistenceMode, Value: *mode})
	}

	in.AddPersistenceVolumeMounts(c)
	in.AddPersistenceVolumes(sts)
	in.AddPersistencePVCs(deployment, sts)
}

func (in *CoherenceSpec) GetMetricsPort() int32 {
	switch {
	case in == nil:
		return 0
	case in.Metrics == nil || in.Metrics.Port == nil:
		return DefaultMetricsPort
	default:
		return *in.Metrics.Port
	}
}

func (in *CoherenceSpec) GetManagementPort() int32 {
	switch {
	case in == nil:
		return 0
	case in.Management == nil || in.Management.Port == nil:
		return DefaultMetricsPort
	default:
		return *in.Management.Port
	}
}

// ----- CoherenceWKASpec struct --------------------------------------------
// CoherenceWKASpec configures Coherence well-known-addressing to use an
// existing Coherence deployment for WKA.
// +k8s:openapi-gen=true
type CoherenceWKASpec struct {
	// The name of the existing Coherence deployment to use for WKA.
	Deployment string `json:"deployment"`
	// The optional namespace of the existing Coherence deployment to use for WKA
	// if different from this deployment's namespace.
	// +optional
	Namespace string `json:"namespace,omitempty"`
}

// ----- CoherenceTracingSpec struct ----------------------------------------

// CoherenceTracingSpec configures Coherence tracing.
// +k8s:openapi-gen=true
type CoherenceTracingSpec struct {
	// Ratio is the tracing sampling-ratio, which controls the likelihood of a tracing span being collected.
	// For instance, a value of 1.0 will result in all tracing spans being collected, while a value of 0.1
	// will result in roughly 1 out of every 10 tracing spans being collected.
	//
	// A value of 0 indicates that tracing spans should only be collected if they are already in the context
	// of another tracing span.  With such a configuration, Coherence will not initiate tracing on its own,
	// and it is up to the application to start an outer tracing span, from which Coherence will then collect
	// inner tracing spans.
	//
	// A value of -1 disables tracing completely.
	//
	// The Coherence default is -1 if not overridden. For values other than -1, numbers between 0 and 1 are
	// acceptable.
	//
	// Due to decimal values not being allowed in a CRD field the ratio value is held as a string.
	// Consequently there is no validation that the value entered is valid and the JVM may fail
	// to start properly in an invalid non-numeric value is entered.
	//
	// +optional
	Ratio *string `json:"ratio,omitempty"`
}

// ----- JVMSpec struct -----------------------------------------------------

// The JVM configuration.
// +k8s:openapi-gen=true
type JVMSpec struct {
	// Classpath specifies additional items to add to the classpath of the JVM.
	// +listType=atomic
	// +optional
	Classpath []string `json:"classpath,omitempty"`
	// Args specifies the options (System properties, -XX: args etc) to pass to the JVM.
	// +listType=atomic
	// +optional
	Args []string `json:"args,omitempty"`
	// The settings for enabling debug mode in the JVM.
	// +optional
	Debug *JvmDebugSpec `json:"debug,omitempty"`
	// If set to true Adds the  -XX:+UseContainerSupport JVM option to ensure that the JVM
	// respects any container resource limits.
	// The default value is true
	// +optional
	UseContainerLimits *bool `json:"useContainerLimits,omitempty"`
	// Set JVM garbage collector options.
	// +optional
	Gc *JvmGarbageCollectorSpec `json:"gc,omitempty"`
	// +optional
	DiagnosticsVolume *corev1.VolumeSource `json:"diagnosticsVolume,omitempty"`
	// Configure the JVM memory options.
	// +optional
	Memory *JvmMemorySpec `json:"memory,omitempty"`
	// Configure JMX using JMXMP.
	// +optional
	Jmxmp *JvmJmxmpSpec `json:"jmxmp,omitempty"`
	// A flag indicating whether to automatically add the default classpath for images
	// created by the JIB tool https://github.com/GoogleContainerTools/jib
	// If true then the /app/lib/* /app/classes and /app/resources
	// entries are added to the JVM classpath.
	// The default value fif not specified is true.
	// +optional
	UseJibClasspath *bool `json:"useJibClasspath,omitempty"`
}

// Update the StatefulSet with any JVM specific settings
func (in *JVMSpec) UpdateStatefulSet(sts *appsv1.StatefulSet) {
	c := EnsureContainer(ContainerNameCoherence, sts)
	defer ReplaceContainer(sts, c)

	var gc *JvmGarbageCollectorSpec

	if in != nil {
		// Add debug settings
		in.Debug.UpdateCoherenceContainer(c)

		// Add additional classpath items to the Coherence container
		if in.Classpath != nil && len(in.Classpath) > 0 {
			// always use the linux/unix path separator as we only ever run on linux
			cp := strings.Join(in.Classpath, ":")
			c.Env = append(c.Env, corev1.EnvVar{Name: EnvVarJvmExtraClasspath, Value: cp})
		}

		// Add JVM args variables to the Coherence container
		if in.Args != nil && len(in.Args) > 0 {
			args := strings.Join(in.Args, " ")
			c.Env = append(c.Env, corev1.EnvVar{Name: EnvVarJvmArgs, Value: args})
		}

		if in.Memory != nil {
			c.Env = append(c.Env, in.Memory.CreateEnvVars()...)
		}

		if in.Jmxmp != nil {
			c.Env = append(c.Env, in.Jmxmp.CreateEnvVars()...)
		}

		if in.Gc != nil {
			gc = in.Gc
		}
	}

	c.Env = append(c.Env, gc.CreateEnvVars()...)

	// Configure the JVM to use container limits (true by default)
	useContainerLimits := in == nil || in.UseContainerLimits == nil || *in.UseContainerLimits
	c.Env = append(c.Env, corev1.EnvVar{Name: EnvVarJvmUseContainerLimits, Value: strconv.FormatBool(useContainerLimits)})

	// Configure the JVM to use the JIB classpath if UseJibClasspath is not nil
	if in != nil && in.UseJibClasspath != nil {
		c.Env = append(c.Env, corev1.EnvVar{Name: EnvVarJvmClasspathJib, Value: strconv.FormatBool(*in.UseJibClasspath)})
	}

	// Add diagnostic volume if specified otherwise use an empty-volume
	if in != nil && in.DiagnosticsVolume != nil {
		sts.Spec.Template.Spec.Volumes = append(sts.Spec.Template.Spec.Volumes, corev1.Volume{
			Name:         VolumeNameJVM,
			VolumeSource: *in.DiagnosticsVolume,
		})
	} else {
		sts.Spec.Template.Spec.Volumes = append(sts.Spec.Template.Spec.Volumes, corev1.Volume{
			Name:         VolumeNameJVM,
			VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}},
		})
	}
}

// ----- ImageSpec struct ---------------------------------------------------

// ImageSpec defines the settings for a Docker image
// +k8s:openapi-gen=true
type ImageSpec struct {
	// The image name.
	// More info: https://kubernetes.io/docs/concepts/containers/images
	// +optional
	Image *string `json:"image,omitempty"`
	// Image pull policy.
	// One of Always, Never, IfNotPresent.
	// More info: https://kubernetes.io/docs/concepts/containers/images#updating-images
	// +optional
	ImagePullPolicy *corev1.PullPolicy `json:"imagePullPolicy,omitempty"`
}

// Ensure that the image value is set.
func (in *ImageSpec) EnsureImage(image *string) bool {
	if in != nil && in.Image == nil {
		in.Image = image
		return true
	}
	return false
}

// ----- PersistenceSpec struct ---------------------------------------------

// The spec for Coherence persistence.
// +k8s:openapi-gen=true
type PersistenceSpec struct {
	// The persistence mode to use.
	// Valid choices are "on-demand", "active", "active-async".
	// This field will set the coherence.distributed.persistence-mode System property
	// to "default-" + Mode.
	// +optional
	Mode *string `json:"mode,omitempty"`
	// The persistence storage spec for.
	PersistentStorageSpec `json:",inline"`
	// Snapshot values configure the on-disc persistence data snapshot (backup) settings.
	// These settings enable a different location for persistence snapshot data.
	// If not set then snapshot files will be written to the same volume configured for
	// persistence data in the Persistence section.
	// +optional
	Snapshots *PersistentStorageSpec `json:"snapshots,omitempty"`
}

// Obtain the persistence mode to be used.
func (in *PersistenceSpec) GetMode() *string {
	if in == nil || in.Mode == nil {
		return nil
	}
	return in.Mode
}

func (in *PersistenceSpec) CreatePersistentVolumeClaims(deployment *Coherence) []corev1.PersistentVolumeClaim {
	var pvcs []corev1.PersistentVolumeClaim
	if in != nil {
		// Only create the PVC if there is not a volume definition configured
		if pvc := in.CreatePersistentVolumeClaim(deployment, VolumeNamePersistence); pvc != nil {
			pvcs = append(pvcs, *pvc)
		}

		// Only create the snapshots PVC if there is not a snapshots volume definition configured
		if pvc := in.Snapshots.CreatePersistentVolumeClaim(deployment, VolumeNameSnapshots); pvc != nil {
			pvcs = append(pvcs, *pvc)
		}
	}
	return pvcs
}

// Add the persistence and snapshot volumes
func (in *PersistenceSpec) CreatePersistenceVolumes() []corev1.Volume {
	var vols []corev1.Volume

	if in != nil {
		if in.Volume != nil {
			// A Persistence Volume s configured so use it
			vols = append(vols, in.CreatePersistenceVolume(VolumeNamePersistence))
		}

		if in.Snapshots != nil && in.Snapshots.Volume != nil {
			// A Snapshots Volume s configured so use it
			vols = append(vols, in.Snapshots.CreatePersistenceVolume(VolumeNameSnapshots))
		}
	}
	return vols
}

// Add the persistence and snapshot volume mounts to the specified container
func (in *PersistenceSpec) AddVolumeMounts(c *corev1.Container) {
	if in == nil {
		return
	}

	if in.Volume != nil || in.PersistentVolumeClaim != nil {
		// Set the persistence location environment variable
		c.Env = append(c.Env, corev1.EnvVar{Name: EnvVarCohPersistenceDir, Value: VolumeMountPathPersistence})
		// Add the persistence volume mount
		c.VolumeMounts = append(c.VolumeMounts, corev1.VolumeMount{
			Name:      VolumeNamePersistence,
			MountPath: VolumeMountPathPersistence,
		})
	}

	// Add the snapshot volume mount if required
	if in != nil && in.Snapshots != nil && (in.Snapshots.Volume != nil || in.Snapshots.PersistentVolumeClaim != nil) {
		// Set the snapshot location environment variable
		c.Env = append(c.Env, corev1.EnvVar{Name: EnvVarCohSnapshotDir, Value: VolumeMountPathSnapshots})
		// Add the snapshot volume mount
		c.VolumeMounts = append(c.VolumeMounts, corev1.VolumeMount{
			Name:      VolumeNameSnapshots,
			MountPath: VolumeMountPathSnapshots,
		})
	}
}

// ----- PersistentStorageSpec struct ---------------------------------------

// PersistenceStorageSpec defines the persistence settings for the Coherence
// +k8s:openapi-gen=true
type PersistentStorageSpec struct {
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
	Volume *corev1.VolumeSource `json:"volume,omitempty"` // from k8s.io/api/core/v1
}

// Create a PersistentVolumeClaim if required
func (in *PersistentStorageSpec) CreatePersistentVolumeClaim(deployment *Coherence, name string) *corev1.PersistentVolumeClaim {
	if in == nil || in.Volume != nil || in.PersistentVolumeClaim == nil {
		// no pv required
		return nil
	}

	spec := corev1.PersistentVolumeClaimSpec{}
	if in.PersistentVolumeClaim != nil {
		in.PersistentVolumeClaim.DeepCopyInto(&spec)
	}

	labels := deployment.CreateCommonLabels()
	labels[LabelComponent] = LabelComponentPVC

	return &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: labels,
		},
		Spec: spec,
	}
}

// Create any persistence volumes required.
func (in *PersistentStorageSpec) CreatePersistenceVolume(name string) corev1.Volume {
	source := corev1.VolumeSource{}
	if in.Volume != nil {
		in.Volume.DeepCopyInto(&source)
	}
	return corev1.Volume{Name: name, VolumeSource: source}
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
	KeyStore *string `json:"keyStore,omitempty"`
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

// Create the SSL environment variables
func (in *SSLSpec) CreateEnvVars(prefix, secretMount string) []corev1.EnvVar {
	var envVars []corev1.EnvVar

	if in == nil {
		return envVars
	}

	if in.Enabled != nil && *in.Enabled {
		envVars = append(envVars, corev1.EnvVar{Name: prefix + EnvVarSuffixSSLEnabled, Value: "true"})
	}

	if in.Secrets != nil && *in.Secrets != "" {
		envVars = append(envVars, corev1.EnvVar{Name: prefix + EnvVarSuffixSSLCerts, Value: secretMount})
	}

	if in.KeyStore != nil && *in.KeyStore != "" {
		envVars = append(envVars, corev1.EnvVar{Name: prefix + EnvVarSuffixSSLKeyStore, Value: *in.KeyStore})
	}

	if in.KeyStorePasswordFile != nil && *in.KeyStorePasswordFile != "" {
		envVars = append(envVars, corev1.EnvVar{Name: prefix + EnvVarSuffixSSLKeyStoreCredFile, Value: *in.KeyStorePasswordFile})
	}

	if in.KeyPasswordFile != nil && *in.KeyPasswordFile != "" {
		envVars = append(envVars, corev1.EnvVar{Name: prefix + EnvVarSuffixSSLKeyCredFile, Value: *in.KeyPasswordFile})
	}

	if in.KeyStoreAlgorithm != nil && *in.KeyStoreAlgorithm != "" {
		envVars = append(envVars, corev1.EnvVar{Name: prefix + EnvVarSuffixSSLKeyStoreAlgo, Value: *in.KeyStoreAlgorithm})
	}

	if in.KeyStoreProvider != nil && *in.KeyStoreProvider != "" {
		envVars = append(envVars, corev1.EnvVar{Name: prefix + EnvVarSuffixSSLKeyStoreProvider, Value: *in.KeyStoreProvider})
	}

	if in.KeyStoreType != nil && *in.KeyStoreType != "" {
		envVars = append(envVars, corev1.EnvVar{Name: prefix + EnvVarSuffixSSLKeyStoreType, Value: *in.KeyStoreType})
	}

	if in.TrustStore != nil && *in.TrustStore != "" {
		envVars = append(envVars, corev1.EnvVar{Name: prefix + EnvVarSuffixSSLTrustStore, Value: *in.TrustStore})
	}

	if in.TrustStorePasswordFile != nil && *in.TrustStorePasswordFile != "" {
		envVars = append(envVars, corev1.EnvVar{Name: prefix + EnvVarSuffixSSLTrustStoreCredFile, Value: *in.TrustStorePasswordFile})
	}

	if in.TrustStoreAlgorithm != nil && *in.TrustStoreAlgorithm != "" {
		envVars = append(envVars, corev1.EnvVar{Name: prefix + EnvVarSuffixSSLTrustStoreAlgo, Value: *in.TrustStoreAlgorithm})
	}

	if in.TrustStoreProvider != nil && *in.TrustStoreProvider != "" {
		envVars = append(envVars, corev1.EnvVar{Name: prefix + EnvVarSuffixSSLTrustStoreProvider, Value: *in.TrustStoreProvider})
	}

	if in.TrustStoreType != nil && *in.TrustStoreType != "" {
		envVars = append(envVars, corev1.EnvVar{Name: prefix + EnvVarSuffixSSLTrustStoreType, Value: *in.TrustStoreType})
	}

	if in.RequireClientCert != nil && *in.RequireClientCert {
		envVars = append(envVars, corev1.EnvVar{Name: prefix + EnvVarSuffixSSLRequireClientCert, Value: "true"})
	}

	return envVars
}

// ----- NamedPortSpec struct ----------------------------------------------------
// NamedPortSpec defines a named port for a Coherence component
// +k8s:openapi-gen=true
type NamedPortSpec struct {
	// Name specifies the name of the port.
	Name string `json:"name"`
	// Port specifies the port used.
	// +optional
	Port int32 `json:"port,omitempty"`
	// Protocol for container port. Must be UDP or TCP. Defaults to "TCP"
	// +optional
	Protocol *corev1.Protocol `json:"protocol,omitempty"`
	// The port on each node on which this service is exposed when type=NodePort or LoadBalancer.
	// Usually assigned by the system. If specified, it will be allocated to the service
	// if unused or else creation of the service will fail.
	// Default is to auto-allocate a port if the ServiceType of this Service requires one.
	// More info: https://kubernetes.io/docs/concepts/services-networking/service/#type-nodeport
	// +optional
	NodePort *int32 `json:"nodePort,omitempty"`
	// Number of port to expose on the host.
	// If specified, this must be a valid port number, 0 < x < 65536.
	// If HostNetwork is specified, this must match ContainerPort.
	// Most containers do not need this.
	// +optional
	HostPort *int32 `json:"hostPort,omitempty"`
	// What host IP to bind the external port to.
	// +optional
	HostIP *string `json:"hostIP,omitempty"`
	// Service configures the Kubernetes Service used to expose the port.
	// +optional
	Service *ServiceSpec `json:"service,omitempty"`
	// The specification of a Prometheus ServiceMonitor resource
	// that will be created for the Service being exposed for this
	// port.
	// +optional
	ServiceMonitor *ServiceMonitorSpec `json:"serviceMonitor,omitempty"`
}

// Create the Kubernetes service to expose this port.
func (in *NamedPortSpec) CreateService(deployment *Coherence) *corev1.Service {
	if in == nil || !in.IsEnabled() {
		return nil
	}

	var name string
	if in.Service != nil && in.Service.Name != nil {
		name = in.Service.GetName()
	} else {
		name = fmt.Sprintf("%s-%s", deployment.Name, in.Name)
	}

	// The labels for the service
	svcLabels := deployment.CreateCommonLabels()
	svcLabels[LabelComponent] = LabelComponentPortService
	svcLabels[LabelPort] = in.Name
	if in.Service != nil {
		for k, v := range in.Service.Labels {
			svcLabels[k] = v
		}
	}

	// The service annotations
	var ann map[string]string
	if in.Service != nil && in.Service.Annotations != nil {
		ann = in.Service.Annotations
	}

	// Create the Service spec
	spec := in.Service.createServiceSpec()

	// Add the port
	spec.Ports = []corev1.ServicePort{
		{
			Name:       in.Name,
			Protocol:   in.GetProtocol(),
			Port:       in.GetPort(deployment),
			TargetPort: intstr.FromInt(int(in.Port)),
			NodePort:   in.GetNodePort(),
		},
	}

	// Add the service selector
	spec.Selector = deployment.Spec.CreatePodSelectorLabels(deployment)

	svc := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:   deployment.GetNamespace(),
			Name:        name,
			Labels:      svcLabels,
			Annotations: ann,
		},
		Spec: spec,
	}

	return &svc
}

// Create the Prometheus ServiceMonitor to expose this port if enabled.
func (in *NamedPortSpec) CreateServiceMonitor(deployment *Coherence) *monitoringv1.ServiceMonitor {
	if in == nil || !in.IsEnabled() {
		return nil
	}
	if in.ServiceMonitor == nil || in.ServiceMonitor.Enabled == nil || !*in.ServiceMonitor.Enabled {
		return nil
	}

	var name string
	if in.Service != nil && in.Service.Name != nil {
		name = in.Service.GetName()
	} else {
		name = fmt.Sprintf("%s-%s", deployment.Name, in.Name)
	}

	// The labels for the ServiceMonitor
	labels := deployment.CreateCommonLabels()
	labels[LabelComponent] = LabelComponentPortServiceMonitor
	for k, v := range in.ServiceMonitor.Labels {
		labels[k] = v
	}

	// The selector labels for the ServiceMonitor
	selector := deployment.CreateCommonLabels()
	selector[LabelComponent] = LabelComponentPortService
	selector[LabelPort] = in.Name

	endpoint := in.ServiceMonitor.CreateEndpoint()
	endpoint.Port = in.Name
	endpoint.TargetPort = nil
	endpoint.RelabelConfigs = append(endpoint.RelabelConfigs, &monitoringv1.RelabelConfig{
		Action: "labeldrop",
		Regex:  "(endpoint|instance|job|service)",
	})

	spec := in.ServiceMonitor.CreateServiceMonitor()
	spec.Selector = metav1.LabelSelector{MatchLabels: selector}
	spec.Endpoints = []monitoringv1.Endpoint{endpoint}

	return &monitoringv1.ServiceMonitor{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: deployment.GetNamespace(),
			Labels:    labels,
		},
		Spec: spec,
	}
}

func (in *NamedPortSpec) IsEnabled() bool {
	return in != nil && in.Service.IsEnabled()
}

func (in *NamedPortSpec) GetProtocol() corev1.Protocol {
	if in == nil || in.Protocol == nil {
		return corev1.ProtocolTCP
	}
	return *in.Protocol
}

func (in *NamedPortSpec) GetPort(d *Coherence) int32 {
	switch {
	case in == nil:
		return 0
	case in != nil && in.Service != nil && in.Service.Port != nil:
		return *in.Service.Port
	case in.Port == 0 && strings.ToLower(in.Name) == PortNameMetrics:
		// special case for well known port - metrics
		return d.Spec.GetMetricsPort()
	case in.Port == 0 && strings.ToLower(in.Name) == PortNameManagement:
		// special case for well known port - management
		return d.Spec.GetManagementPort()
	default:
		return in.Port
	}
}

func (in *NamedPortSpec) GetNodePort() int32 {
	if in == nil || in.NodePort == nil {
		return 0
	}
	return *in.NodePort
}

func (in *NamedPortSpec) CreatePort(d *Coherence) corev1.ContainerPort {
	return corev1.ContainerPort{
		Name:          in.Name,
		ContainerPort: in.GetPort(d),
		Protocol:      in.GetProtocol(),
		HostPort:      notNilInt32(in.HostPort),
		HostIP:        notNilString(in.HostIP),
	}
}

// ----- ServiceMonitorSpec struct ---------------------------------------------

// The ServiceMonitor spec for a port service.
// +k8s:openapi-gen=true
type ServiceMonitorSpec struct {
	// Enabled is a flag to enable or disable creation of a Prometheus ServiceMonitor for a port.
	// If Prometheus ServiceMonitor CR is not installed no ServiceMonitor then even if this flag
	// is true no ServiceMonitor will be created.
	// +optional
	Enabled *bool `json:"enabled,omitempty"`
	// Additional labels to add to the ServiceMonitor.
	// More info: http://kubernetes.io/docs/user-guide/labels
	// +optional
	Labels map[string]string `json:"labels,omitempty"`
	// The label to use to retrieve the job name from.
	// See https://coreos.com/operators/prometheus/docs/latest/api.html#servicemonitorspec
	// +optional
	JobLabel string `json:"jobLabel,omitempty"`
	// TargetLabels transfers labels on the Kubernetes Service onto the target.
	// See https://coreos.com/operators/prometheus/docs/latest/api.html#servicemonitorspec
	// +listType=atomic
	// +optional
	TargetLabels []string `json:"targetLabels,omitempty"`
	// PodTargetLabels transfers labels on the Kubernetes Pod onto the target.
	// See https://coreos.com/operators/prometheus/docs/latest/api.html#servicemonitorspec
	// +listType=atomic
	// +optional
	PodTargetLabels []string `json:"podTargetLabels,omitempty"`
	// SampleLimit defines per-scrape limit on number of scraped samples that will be accepted.
	// See https://coreos.com/operators/prometheus/docs/latest/api.html#servicemonitorspec
	// +optional
	SampleLimit uint64 `json:"sampleLimit,omitempty"`
	// HTTP path to scrape for metrics.
	// See https://coreos.com/operators/prometheus/docs/latest/api.html#endpoint
	// +optional
	Path string `json:"path,omitempty"`
	// HTTP scheme to use for scraping.
	// See https://coreos.com/operators/prometheus/docs/latest/api.html#endpoint
	// +optional
	Scheme string `json:"scheme,omitempty"`
	// Optional HTTP URL parameters
	// See https://coreos.com/operators/prometheus/docs/latest/api.html#endpoint
	// +optional
	Params map[string][]string `json:"params,omitempty"`
	// Interval at which metrics should be scraped
	// See https://coreos.com/operators/prometheus/docs/latest/api.html#endpoint
	// +optional
	Interval string `json:"interval,omitempty"`
	// Timeout after which the scrape is ended
	// See https://coreos.com/operators/prometheus/docs/latest/api.html#endpoint
	// +optional
	ScrapeTimeout string `json:"scrapeTimeout,omitempty"`
	// TLS configuration to use when scraping the endpoint
	// See https://coreos.com/operators/prometheus/docs/latest/api.html#endpoint
	// +optional
	TLSConfig *monitoringv1.TLSConfig `json:"tlsConfig,omitempty"`
	// File to read bearer token for scraping targets.
	// See https://coreos.com/operators/prometheus/docs/latest/api.html#endpoint
	// +optional
	BearerTokenFile string `json:"bearerTokenFile,omitempty"`
	// Secret to mount to read bearer token for scraping targets. The secret
	// needs to be in the same namespace as the service monitor and accessible by
	// the Prometheus Operator.
	// See https://coreos.com/operators/prometheus/docs/latest/api.html#endpoint
	// +optional
	BearerTokenSecret corev1.SecretKeySelector `json:"bearerTokenSecret,omitempty"`
	// HonorLabels chooses the metric's labels on collisions with target labels.
	// See https://coreos.com/operators/prometheus/docs/latest/api.html#endpoint
	// +optional
	HonorLabels bool `json:"honorLabels,omitempty"`
	// HonorTimestamps controls whether Prometheus respects the timestamps present in scraped data.
	// See https://coreos.com/operators/prometheus/docs/latest/api.html#endpoint
	// +optional
	HonorTimestamps *bool `json:"honorTimestamps,omitempty"`
	// BasicAuth allow an endpoint to authenticate over basic authentication
	// More info: https://prometheus.io/docs/operating/configuration/#endpoints
	// See https://coreos.com/operators/prometheus/docs/latest/api.html#endpoint
	// +optional
	BasicAuth *monitoringv1.BasicAuth `json:"basicAuth,omitempty"`
	// MetricRelabelings to apply to samples before ingestion.
	// See https://coreos.com/operators/prometheus/docs/latest/api.html#endpoint
	// +listType=atomic
	// +optional
	MetricRelabelings []*monitoringv1.RelabelConfig `json:"metricRelabelings,omitempty"`
	// Relabelings to apply to samples before scraping.
	// More info: https://prometheus.io/docs/prometheus/latest/configuration/configuration/#relabel_config
	// See https://coreos.com/operators/prometheus/docs/latest/api.html#endpoint
	// +listType=atomic
	// +optional
	Relabelings []*monitoringv1.RelabelConfig `json:"relabelings,omitempty"`
	// ProxyURL eg http://proxyserver:2195 Directs scrapes to proxy through this endpoint.
	// See https://coreos.com/operators/prometheus/docs/latest/api.html#endpoint
	// +optional
	ProxyURL *string `json:"proxyURL,omitempty"`
}

func (in *ServiceMonitorSpec) CreateServiceMonitor() monitoringv1.ServiceMonitorSpec {
	if in == nil {
		return monitoringv1.ServiceMonitorSpec{}
	}

	return monitoringv1.ServiceMonitorSpec{
		JobLabel:        in.JobLabel,
		TargetLabels:    in.TargetLabels,
		PodTargetLabels: in.PodTargetLabels,
		SampleLimit:     in.SampleLimit,
	}
}

func (in *ServiceMonitorSpec) CreateEndpoint() monitoringv1.Endpoint {
	if in == nil {
		return monitoringv1.Endpoint{}
	}

	return monitoringv1.Endpoint{
		Path:                 in.Path,
		Scheme:               in.Scheme,
		Params:               in.Params,
		Interval:             in.Interval,
		ScrapeTimeout:        in.ScrapeTimeout,
		TLSConfig:            in.TLSConfig,
		BearerTokenFile:      in.BearerTokenFile,
		BearerTokenSecret:    in.BearerTokenSecret,
		HonorLabels:          in.HonorLabels,
		HonorTimestamps:      in.HonorTimestamps,
		BasicAuth:            in.BasicAuth,
		MetricRelabelConfigs: in.MetricRelabelings,
		RelabelConfigs:       in.Relabelings,
		ProxyURL:             in.ProxyURL,
	}
}

// ----- JvmDebugSpec struct ---------------------------------------------------

// The JVM Debug specific configuration.
// +k8s:openapi-gen=true
type JvmDebugSpec struct {
	// Enabled is a flag to enable or disable running the JVM in debug mode. Default is disabled.
	// +optional
	Enabled *bool `json:"enabled,omitempty"`
	// A boolean true if the target VM is to be suspended immediately before the main class is loaded;
	// false otherwise. The default value is false.
	// +optional
	Suspend *bool `json:"suspend,omitempty"`
	// Attach specifies the address of the debugger that the JVM should attempt to connect back to
	// instead of listening on a port.
	// +optional
	Attach *string `json:"attach,omitempty"`
	// The port that the debugger will listen on; the default is 5005.
	// +optional
	Port *int32 `json:"port,omitempty"`
}

// Update the Coherence Container with any JVM specific settings
func (in *JvmDebugSpec) UpdateCoherenceContainer(c *corev1.Container) {
	if in == nil || in.Enabled == nil || !*in.Enabled {
		// nothing to do, debug is either nil or disabled
		return
	}

	c.Ports = append(c.Ports, corev1.ContainerPort{
		Name:          PortNameDebug,
		ContainerPort: notNilInt32OrDefault(in.Port, DefaultDebugPort),
	})

	c.Env = append(c.Env, in.CreateEnvVars()...)
}

// Create the JVM debugger environment variables for the Coherence container.
func (in *JvmDebugSpec) CreateEnvVars() []corev1.EnvVar {
	var envVars []corev1.EnvVar

	if in == nil || in.Enabled == nil || !*in.Enabled {
		return envVars
	}

	envVars = append(envVars,
		corev1.EnvVar{Name: EnvVarJvmDebugEnabled, Value: "true"},
		corev1.EnvVar{Name: EnvVarJvmDebugPort, Value: Int32PtrToStringWithDefault(in.Port, DefaultDebugPort)},
	)

	if in != nil && in.Suspend != nil && *in.Suspend {
		envVars = append(envVars, corev1.EnvVar{Name: EnvVarJvmDebugSuspended, Value: "true"})
	}

	if in != nil && in.Attach != nil {
		envVars = append(envVars, corev1.EnvVar{Name: EnvVarJvmDebugAttach, Value: *in.Attach})
	}

	return envVars
}

// ----- JVM GC struct ------------------------------------------------------

// Options for managing the JVM garbage collector.
// +k8s:openapi-gen=true
type JvmGarbageCollectorSpec struct {
	// The name of the JVM garbage collector to use.
	// G1 - adds the -XX:+UseG1GC option
	// CMS - adds the -XX:+UseConcMarkSweepGC option
	// Parallel - adds the -XX:+UseParallelGC
	// Default - use the JVMs default collector
	// The field value is case insensitive
	// If not set G1 is used.
	// If set to a value other than those above then
	// the default collector for the JVM will be used.
	// +optional
	Collector *string `json:"collector,omitempty"`
	// Args specifies the GC options to pass to the JVM.
	// +listType=atomic
	// +optional
	Args []string `json:"args,omitempty"`
	// Enable the following GC logging args  -verbose:gc -XX:+PrintGCDetails -XX:+PrintGCTimeStamps
	// -XX:+PrintHeapAtGC -XX:+PrintTenuringDistribution -XX:+PrintGCApplicationStoppedTime
	// -XX:+PrintGCApplicationConcurrentTime
	// Default is true
	// +optional
	Logging *bool `json:"logging,omitempty"`
}

// Create the GC environment variables for the Coherence container.
func (in *JvmGarbageCollectorSpec) CreateEnvVars() []corev1.EnvVar {
	var envVars []corev1.EnvVar

	// Add any GC args
	if in != nil && in.Args != nil && len(in.Args) > 0 {
		args := strings.Join(in.Args, " ")
		envVars = append(envVars, corev1.EnvVar{Name: EnvVarJvmGcArgs, Value: args})
	}

	// Set the collector to use
	if in != nil && in.Collector != nil && *in.Collector != "" {
		envVars = append(envVars, corev1.EnvVar{Name: EnvVarJvmGcCollector, Value: *in.Collector})
	}

	// Enable or disable GC logging
	if in != nil && in.Logging != nil {
		envVars = append(envVars, corev1.EnvVar{Name: EnvVarJvmGcLogging, Value: BoolPtrToString(in.Logging)})
	} else {
		envVars = append(envVars, corev1.EnvVar{Name: EnvVarJvmGcLogging, Value: "false"})
	}

	return envVars
}

// ----- JVM MemoryGC struct ------------------------------------------------

// Options for managing the JVM memory.
// +k8s:openapi-gen=true
type JvmMemorySpec struct {
	// HeapSize is the min/max heap value to pass to the JVM.
	// The format should be the same as that used for Java's -Xms and -Xmx JVM options.
	// If not set the JVM defaults are used.
	// +optional
	HeapSize *string `json:"heapSize,omitempty"`
	// Sets the JVM option `-XX:MaxRAM=N` which sets the maximum amount of memory used by
	// the JVM to `n`, where `n` is expressed in terms of megabytes (for example, `100m`)
	// or gigabytes (for example `2g`).
	// +optional
	MaxRAM *string `json:"maxRAM,omitempty"`
	// Set initial heap size as a percentage of total memory.
	//
	// This option will be ignored if HeapSize is set.
	//
	// Valid values are decimal numbers between 0 and 100.
	//
	// This field is a string value as CRDs do not support decimal numbers.
	// Consequently, there is no validation on the value entered so the
	// JVM may fail to start if an invalid value is entered.
	//
	// This field maps the the -XX:InitialRAMPercentage JVM option and will
	// be incompatible with some JVMs that do not have this option (e.g. Java 8).
	// +optional
	InitialRAMPercentage *string `json:"initialRAMPercentage,omitempty"`
	// Set maximum heap size as a percentage of total memory.
	//
	// This option will be ignored if HeapSize is set.
	//
	// Valid values are decimal numbers between 0 and 100.
	//
	// This field is a string value as CRDs do not support decimal numbers.
	// Consequently, there is no validation on the value entered so the
	// JVM may fail to start if an invalid value is entered.
	//
	// This field maps the the -XX:MaxRAMPercentage JVM option and will
	// be incompatible with some JVMs that do not have this option (e.g. Java 8).
	// +optional
	MaxRAMPercentage *string `json:"maxRAMPercentage,omitempty"`
	// Set the minimal JVM Heap size as a percentage of the total memory.
	//
	// This option will be ignored if HeapSize is set.
	//
	// Valid values are decimal numbers between 0 and 100.
	//
	// This field is a string value as CRDs do not support decimal numbers.
	// Consequently, there is no validation on the value entered so the
	// JVM may fail to start if an invalid value is entered.
	//
	// This field maps the the -XX:MinRAMPercentage JVM option and will
	// be incompatible with some JVMs that do not have this option (e.g. Java 8).
	// +optional
	MinRAMPercentage *string `json:"minRAMPercentage,omitempty"`
	// StackSize is the stack size value to pass to the JVM.
	// The format should be the same as that used for Java's -Xss JVM option.
	// If not set the JVM defaults are used.
	// +optional
	StackSize *string `json:"stackSize,omitempty"`
	// MetaspaceSize is the min/max metaspace size to pass to the JVM.
	// This sets the -XX:MetaspaceSize and -XX:MaxMetaspaceSize=size JVM options.
	// If not set the JVM defaults are used.
	// +optional
	MetaspaceSize *string `json:"metaspaceSize,omitempty"`
	// DirectMemorySize sets the maximum total size (in bytes) of the New I/O (the java.nio package) direct-buffer
	// allocations. This value sets the -XX:MaxDirectMemorySize JVM option.
	// If not set the JVM defaults are used.
	// +optional
	DirectMemorySize *string `json:"directMemorySize,omitempty"`
	// Adds the -XX:NativeMemoryTracking=mode  JVM options
	// where mode is on of "off", "summary" or "detail", the default is "summary"
	// If not set to "off" also add -XX:+PrintNMTStatistics
	// +optional
	NativeMemoryTracking *string `json:"nativeMemoryTracking,omitempty"`
	// Configure the JVM behaviour when an OutOfMemoryError occurs.
	// +optional
	OnOutOfMemory *JvmOutOfMemorySpec `json:"onOutOfMemory,omitempty"`
}

// Create the environment variables to add to the Coherence container
func (in *JvmMemorySpec) CreateEnvVars() []corev1.EnvVar {
	var envVars []corev1.EnvVar

	if in == nil {
		return envVars
	}

	if in.HeapSize != nil && *in.HeapSize != "" {
		envVars = append(envVars, corev1.EnvVar{Name: EnvVarJvmMemoryHeap, Value: *in.HeapSize})
	}

	if in.MaxRAM != nil && *in.MaxRAM != "" {
		envVars = append(envVars, corev1.EnvVar{Name: EnvVarJvmMaxRAM, Value: *in.MaxRAM})
	}

	if in.InitialRAMPercentage != nil && *in.InitialRAMPercentage != "" {
		envVars = append(envVars, corev1.EnvVar{Name: EnvVarJvmInitialRAMPercentage, Value: *in.InitialRAMPercentage})
	}

	if in.MaxRAMPercentage != nil && *in.MaxRAMPercentage != "" {
		envVars = append(envVars, corev1.EnvVar{Name: EnvVarJvmMaxRAMPercentage, Value: *in.MaxRAMPercentage})
	}

	if in.MinRAMPercentage != nil && *in.MinRAMPercentage != "" {
		envVars = append(envVars, corev1.EnvVar{Name: EnvVarJvmMinRAMPercentage, Value: *in.MinRAMPercentage})
	}

	if in.DirectMemorySize != nil && *in.DirectMemorySize != "" {
		envVars = append(envVars, corev1.EnvVar{Name: EnvVarJvmMemoryDirect, Value: *in.DirectMemorySize})
	}

	if in.StackSize != nil && *in.StackSize != "" {
		envVars = append(envVars, corev1.EnvVar{Name: EnvVarJvmMemoryStack, Value: *in.StackSize})
	}

	if in.MetaspaceSize != nil && *in.MetaspaceSize != "" {
		envVars = append(envVars, corev1.EnvVar{Name: EnvVarJvmMemoryMeta, Value: *in.MetaspaceSize})
	}

	if in.NativeMemoryTracking != nil && *in.NativeMemoryTracking != "" {
		envVars = append(envVars, corev1.EnvVar{Name: EnvVarJvmMemoryNativeTracking, Value: *in.NativeMemoryTracking})
	}

	envVars = append(envVars, in.OnOutOfMemory.CreateEnvVars()...)

	return envVars
}

// ----- JVM Out Of Memory struct -------------------------------------------

// Options for managing the JVM behaviour when an OutOfMemoryError occurs.
// +k8s:openapi-gen=true
type JvmOutOfMemorySpec struct {
	// If set to true the JVM will exit when an OOM error occurs.
	// Default is true
	// +optional
	Exit *bool `json:"exit,omitempty"`
	// If set to true adds the -XX:+HeapDumpOnOutOfMemoryError JVM option to cause a heap dump
	// to be created when an OOM error occurs.
	// Default is true
	// +optional
	HeapDump *bool `json:"heapDump,omitempty"`
}

// Create the environment variables for the out of memory actions
func (in *JvmOutOfMemorySpec) CreateEnvVars() []corev1.EnvVar {
	var envVars []corev1.EnvVar

	if in != nil {
		if in.Exit != nil {
			envVars = append(envVars, corev1.EnvVar{Name: EnvVarJvmOomExit, Value: BoolPtrToString(in.Exit)})
		}
		if in.HeapDump != nil {
			envVars = append(envVars, corev1.EnvVar{Name: EnvVarJvmOomHeapDump, Value: BoolPtrToString(in.HeapDump)})
		}
	}

	return envVars
}

// ----- JvmJmxmpSpec struct -------------------------------------------------------

// Options for configuring JMX using JMXMP.
// +k8s:openapi-gen=true
type JvmJmxmpSpec struct {
	// If set to true the JMXMP support will be enabled.
	// Default is false
	// +optional
	Enabled *bool `json:"enabled,omitempty"`
	// The port tht the JMXMP MBeanServer should bind to.
	// If not set the default port is 9099
	// +optional
	Port *int32 `json:"port,omitempty"`
}

// Create any required environment variables for the Coherence container
func (in *JvmJmxmpSpec) CreateEnvVars() []corev1.EnvVar {
	enabled := in != nil && in.Enabled != nil && *in.Enabled

	envVars := []corev1.EnvVar{{Name: EnvVarJvmJmxmpEnabled, Value: strconv.FormatBool(enabled)}}
	envVars = append(envVars, corev1.EnvVar{Name: EnvVarJvmJmxmpPort, Value: Int32PtrToStringWithDefault(in.Port, DefaultJmxmpPort)})

	return envVars
}

// ----- PortSpecWithSSL struct ----------------------------------------------------

// PortSpecWithSSL defines a port with SSL settings for a Coherence component
// +k8s:openapi-gen=true
type PortSpecWithSSL struct {
	// Enable or disable flag.
	// +optional
	Enabled *bool `json:"enabled,omitempty"`
	// The port to bind to.
	// +optional
	Port *int32 `json:"port,omitempty"`
	// SSL configures SSL settings for a Coherence component
	// +optional
	SSL *SSLSpec `json:"ssl,omitempty"`
}

// IsSSLEnabled returns true if this port is SSL enabled
func (in *PortSpecWithSSL) IsSSLEnabled() bool {
	if in == nil || in.SSL == nil {
		return false
	}
	return in.SSL.Enabled != nil && *in.SSL.Enabled
}

// Create environment variables for the Coherence container
func (in *PortSpecWithSSL) CreateEnvVars(prefix, secretMount string, defaultPort int32) []corev1.EnvVar {
	if in == nil || !notNilBool(in.Enabled) {
		// disabled
		return []corev1.EnvVar{{Name: prefix + EnvVarCohEnabledSuffix, Value: "false"}}
	}

	envVars := []corev1.EnvVar{{Name: prefix + EnvVarCohEnabledSuffix, Value: "true"}}
	envVars = append(envVars, in.SSL.CreateEnvVars(prefix, secretMount)...)

	// add the port environment variable
	port := notNilInt32OrDefault(in.Port, defaultPort)
	envVars = append(envVars, corev1.EnvVar{Name: prefix + EnvVarCohPortSuffix, Value: Int32ToString(port)})

	return envVars
}

// Add the SSL secret volume and volume mount if required
func (in *PortSpecWithSSL) AddSSLVolumes(sts *appsv1.StatefulSet, c *corev1.Container, volName, path string) {
	if in == nil || !notNilBool(in.Enabled) || in.SSL == nil || !notNilBool(in.SSL.Enabled) {
		// the port spec is nil or disabled or SSL is nil or disabled
		return
	}

	if in.SSL.Secrets != nil && *in.SSL.Secrets != "" {
		c.VolumeMounts = append(c.VolumeMounts, corev1.VolumeMount{
			Name:      volName,
			ReadOnly:  true,
			MountPath: path,
		})

		sts.Spec.Template.Spec.Volumes = append(sts.Spec.Template.Spec.Volumes, corev1.Volume{
			Name: volName,
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName:  *in.SSL.Secrets,
					DefaultMode: pointer.Int32Ptr(int32(0777)),
				},
			},
		})
	}

}

// ----- ServiceSpec struct -------------------------------------------------
// ServiceSpec defines the settings for a Service
// +k8s:openapi-gen=true
type ServiceSpec struct {
	// Enabled controls whether to create the service yaml or not
	// +optional
	Enabled *bool `json:"enabled,omitempty"`
	// An optional name to use to override the generated service name.
	// +optional
	Name *string `json:"name,omitempty"`
	// The service port value
	// +optional
	Port *int32 `json:"port,omitempty"`
	// Kind is the K8s service type (typically ClusterIP or LoadBalancer)
	// The default is "ClusterIP".
	// +optional
	Type *corev1.ServiceType `json:"type,omitempty"`
	// externalIPs is a list of IP addresses for which nodes in the cluster
	// will also accept traffic for this service.  These IPs are not managed by
	// Kubernetes.  The user is responsible for ensuring that traffic arrives
	// at a node with this IP.  A common example is external load-balancers
	// that are not part of the Kubernetes system.
	// +optional
	// +listType=atomic
	ExternalIPs []string `json:"externalIPs,omitempty"`
	// clusterIP is the IP address of the service and is usually assigned
	// randomly by the master. If an address is specified manually and is not in
	// use by others, it will be allocated to the service; otherwise, creation
	// of the service will fail. This field can not be changed through updates.
	// Valid values are "None", empty string (""), or a valid IP address. "None"
	// can be specified for headless services when proxying is not required.
	// Only applies to types ClusterIP, NodePort, and LoadBalancer. Ignored if
	// type is ExternalName.
	// More info: https://kubernetes.io/docs/concepts/services-networking/service/#virtual-ips-and-service-proxies
	// +optional
	ClusterIP *string `json:"clusterIP,omitempty"`
	// LoadBalancerIP is the IP address of the load balancer
	// +optional
	LoadBalancerIP *string `json:"loadBalancerIP,omitempty"`
	// The extra labels to add to the service.
	// More info: http://kubernetes.io/docs/user-guide/labels
	// +optional
	Labels map[string]string `json:"labels,omitempty"`
	// Annotations is free form yaml that will be added to the service annotations
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`
	// Supports "ClientIP" and "None". Used to maintain session affinity.
	// Enable client IP based session affinity.
	// Must be ClientIP or None.
	// Defaults to None.
	// More info: https://kubernetes.io/docs/concepts/services-networking/service/#virtual-ips-and-service-proxies
	// +optional
	SessionAffinity *corev1.ServiceAffinity `json:"sessionAffinity,omitempty"`
	// If specified and supported by the platform, this will restrict traffic through the cloud-provider
	// load-balancer will be restricted to the specified client IPs. This field will be ignored if the
	// cloud-provider does not support the feature."
	// More info: https://kubernetes.io/docs/tasks/access-application-cluster/configure-cloud-provider-firewall/
	// +listType=atomic
	// +optional
	LoadBalancerSourceRanges []string `json:"loadBalancerSourceRanges,omitempty"`
	// externalName is the external reference that kubedns or equivalent will
	// return as a CNAME record for this service. No proxying will be involved.
	// Must be a valid RFC-1123 hostname (https://tools.ietf.org/html/rfc1123)
	// and requires Kind to be ExternalName.
	// +optional
	ExternalName *string `json:"externalName,omitempty"`
	// externalTrafficPolicy denotes if this Service desires to route external
	// traffic to node-local or cluster-wide endpoints. "Local" preserves the
	// client source IP and avoids a second hop for LoadBalancer and Nodeport
	// type services, but risks potentially imbalanced traffic spreading.
	// "Cluster" obscures the client source IP and may cause a second hop to
	// another node, but should have good overall load-spreading.
	// +optional
	ExternalTrafficPolicy *corev1.ServiceExternalTrafficPolicyType `json:"externalTrafficPolicy,omitempty"`
	// healthCheckNodePort specifies the healthcheck nodePort for the service.
	// If not specified, HealthCheckNodePort is created by the service api
	// backend with the allocated nodePort. Will use user-specified nodePort value
	// if specified by the client. Only effects when Kind is set to LoadBalancer
	// and ExternalTrafficPolicy is set to Local.
	// +optional
	HealthCheckNodePort *int32 `json:"healthCheckNodePort,omitempty"`
	// publishNotReadyAddresses, when set to true, indicates that DNS implementations
	// must publish the notReadyAddresses of subsets for the Endpoints associated with
	// the Service. The default value is false.
	// The primary use case for setting this field is to use a StatefulSet's Headless Service
	// to propagate SRV records for its Pods without respect to their readiness for purpose
	// of peer discovery.
	// +optional
	PublishNotReadyAddresses *bool `json:"publishNotReadyAddresses,omitempty"`
	// sessionAffinityConfig contains the configurations of session affinity.
	// +optional
	SessionAffinityConfig *corev1.SessionAffinityConfig `json:"sessionAffinityConfig,omitempty"`
	// ipFamily specifies whether this Service has a preference for a particular IP family (e.g. IPv4 vs.
	// IPv6).  If a specific IP family is requested, the clusterIP field will be allocated from that family, if it is
	// available in the cluster.  If no IP family is requested, the cluster's primary IP family will be used.
	// Other IP fields (loadBalancerIP, loadBalancerSourceRanges, externalIPs) and controllers which
	// allocate external load-balancers should use the same IP family.  Endpoints for this Service will be of
	// this family.  This field is immutable after creation. Assigning a ServiceIPFamily not available in the
	// cluster (e.g. IPv6 in IPv4 only cluster) is an error condition and will fail during clusterIP assignment.
	// +optional
	IPFamily *corev1.IPFamily `json:"ipFamily,omitempty"`
}

// Set the Kind of the service.
func (in *ServiceSpec) GetName() string {
	if in == nil || in.Name == nil {
		return ""
	}
	return *in.Name
}

// Set the Kind of the service.
func (in *ServiceSpec) IsEnabled() bool {
	if in == nil || in.Enabled == nil {
		return true
	}
	return *in.Enabled
}

// Set the Kind of the service.
func (in *ServiceSpec) SetServiceType(t corev1.ServiceType) {
	if in != nil {
		in.Type = &t
	}
}

// Create the service spec for the port.
func (in *ServiceSpec) createServiceSpec() corev1.ServiceSpec {
	spec := corev1.ServiceSpec{}
	if in != nil {
		if in.Type != nil {
			spec.Type = *in.Type
		}
		if in.LoadBalancerIP != nil {
			spec.LoadBalancerIP = *in.LoadBalancerIP
		}
		if in.SessionAffinity != nil {
			spec.SessionAffinity = *in.SessionAffinity
		}
		spec.LoadBalancerSourceRanges = in.LoadBalancerSourceRanges
		if in.ExternalName != nil {
			spec.ExternalName = *in.ExternalName
		}
		if in.ExternalTrafficPolicy != nil {
			spec.ExternalTrafficPolicy = *in.ExternalTrafficPolicy
		}
		if in.HealthCheckNodePort != nil {
			spec.HealthCheckNodePort = *in.HealthCheckNodePort
		}
		if in.PublishNotReadyAddresses != nil {
			spec.PublishNotReadyAddresses = *in.PublishNotReadyAddresses
		}
		if in.ClusterIP != nil {
			spec.ClusterIP = *in.ClusterIP
		}
		spec.SessionAffinityConfig = in.SessionAffinityConfig
		spec.IPFamily = in.IPFamily
		spec.ExternalIPs = in.ExternalIPs
	}
	return spec
}

// ----- ScalingSpec -----------------------------------------------------

// The configuration to control safe scaling.
// +k8s:openapi-gen=true
type ScalingSpec struct {
	// ScalingPolicy describes how the replicas of the deployment will be scaled.
	// The default if not specified is based upon the value of the StorageEnabled field.
	// If StorageEnabled field is not specified or is true the default scaling will be safe, if StorageEnabled is
	// set to false the default scaling will be parallel.
	// +optional
	Policy *ScalingPolicy `json:"policy,omitempty"`
	// The probe to use to determine whether a deployment is Phase HA.
	// If not set the default handler will be used.
	// In most use-cases the default handler would suffice but in
	// advanced use-cases where the application code has a different
	// concept of Phase HA to just checking Coherence services then
	// a different handler may be specified.
	// +optional
	Probe *ScalingProbe `json:"probe,omitempty"`
}

// ----- ScalingProbe ----------------------------------------------------

// ScalingProbe is the handler that will be used to determine how to check for StatusHA in a Coherence.
// StatusHA checking is primarily used during scaling of a deployment, a deployment must be in a safe Phase HA
// state before scaling takes place. If StatusHA handler is disabled for a deployment (by specifically setting
// Enabled to false then no check will take place and a deployment will be assumed to be safe).
// +k8s:openapi-gen=true
type ScalingProbe struct {
	corev1.Handler `json:",inline"`
	// Number of seconds after which the handler times out (only applies to http and tcp handlers).
	// Defaults to 1 second. Minimum value is 1.
	// +optional
	TimeoutSeconds *int `json:"timeoutSeconds,omitempty"`
}

// Returns the timeout value in seconds.
func (in *ScalingProbe) GetTimeout() time.Duration {
	if in == nil || in.TimeoutSeconds == nil || *in.TimeoutSeconds <= 0 {
		return time.Second
	}

	return time.Second * time.Duration(*in.TimeoutSeconds)
}

// ----- ReadinessProbeSpec struct ------------------------------------------

// ReadinessProbeSpec defines the settings for the Coherence Pod readiness probe
// +k8s:openapi-gen=true
type ReadinessProbeSpec struct {
	// The action taken to determine the health of a container
	ProbeHandler `json:",inline"`
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

// The definition of a probe handler.
// +k8s:openapi-gen=true
type ProbeHandler struct {
	// One and only one of the following should be specified.
	// Exec specifies the action to take.
	// +optional
	Exec *corev1.ExecAction `json:"exec,omitempty"`
	// HTTPGet specifies the http request to perform.
	// +optional
	HTTPGet *corev1.HTTPGetAction `json:"httpGet,omitempty"`
	// TCPSocket specifies an action involving a TCP port.
	// TCP hooks not yet supported
	// +optional
	TCPSocket *corev1.TCPSocketAction `json:"tcpSocket,omitempty"`
}

// Update the specified probe spec with the required configuration
func (in *ReadinessProbeSpec) UpdateProbeSpec(port int32, path string, probe *corev1.Probe) {
	switch {
	case in != nil && in.Exec != nil:
		probe.Exec = in.Exec
	case in != nil && in.HTTPGet != nil:
		probe.HTTPGet = in.HTTPGet
	case in != nil && in.TCPSocket != nil:
		probe.TCPSocket = in.TCPSocket
	default:
		probe.HTTPGet = &corev1.HTTPGetAction{
			Path:   path,
			Port:   intstr.FromInt(int(port)),
			Scheme: corev1.URISchemeHTTP,
		}
	}

	if in != nil {
		if in.InitialDelaySeconds != nil {
			probe.InitialDelaySeconds = *in.InitialDelaySeconds
		}
		if in.PeriodSeconds != nil {
			probe.PeriodSeconds = *in.PeriodSeconds
		}
		if in.FailureThreshold != nil {
			probe.FailureThreshold = *in.FailureThreshold
		}
		if in.SuccessThreshold != nil {
			probe.SuccessThreshold = *in.SuccessThreshold
		}
		if in.TimeoutSeconds != nil {
			probe.TimeoutSeconds = *in.TimeoutSeconds
		}
	}
}

// ----- ScalingPolicy type -------------------------------------------------

// ScalingPolicy describes a policy for scaling a cluster deployment
type ScalingPolicy string

// Scaling policy constants
const (
	// Safe means that a deployment will be scaled up or down in a safe manner to ensure no data loss.
	SafeScaling ScalingPolicy = "Safe"
	// Parallel means that a deployment will be scaled up or down by adding or removing members in parallel.
	// If the members of the deployment are storage enabled then this could cause data loss
	ParallelScaling ScalingPolicy = "Parallel"
	// ParallelUpSafeDown means that a deployment will be scaled up by adding or removing members in parallel
	// but will be scaled down in a safe manner to ensure no data loss.
	ParallelUpSafeDownScaling ScalingPolicy = "ParallelUpSafeDown"
)

// ----- LocalObjectReference -----------------------------------------------

// LocalObjectReference contains enough information to let you locate the
// referenced object inside the same namespace.
type LocalObjectReference struct {
	// Name of the referent.
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`
}

// ----- NetworkSpec --------------------------------------------------------

// NetworkSpec configures various networking and DNS settings for Pods in a deployment.
// +k8s:openapi-gen=true
type NetworkSpec struct {
	// Specifies the DNS parameters of a pod. Parameters specified here will be merged to the
	// generated DNS configuration based on DNSPolicy.
	// +optional
	DNSConfig *PodDNSConfig `json:"dnsConfig,omitempty"`
	// Set DNS policy for the pod. Defaults to "ClusterFirst". Valid values are 'ClusterFirstWithHostNet',
	// 'ClusterFirst', 'Default' or 'None'. DNS parameters given in DNSConfig will be merged with the policy
	// selected with DNSPolicy. To have DNS options set along with hostNetwork, you have to specify DNS
	// policy explicitly to 'ClusterFirstWithHostNet'.
	// +optional
	DNSPolicy *corev1.DNSPolicy `json:"dnsPolicy,omitempty"`
	// HostAliases is an optional list of hosts and IPs that will be injected into the pod's hosts file if specified.
	// This is only valid for non-hostNetwork pods.
	// +listType=map
	// +listMapKey=ip
	// +optional
	HostAliases []corev1.HostAlias `json:"hostAliases,omitempty"`
	// Host networking requested for this pod. Use the host's network namespace. If this option is set,
	// the ports that will be used must be specified. Default to false.
	// +optional
	HostNetwork *bool `json:"hostNetwork,omitempty"`
	// Specifies the hostname of the Pod If not specified, the pod's hostname will be set to a system-defined value.
	// +optional
	Hostname *string `json:"hostname,omitempty"`
}

// Update the specified StatefulSet's network settings.
func (in *NetworkSpec) UpdateStatefulSet(sts *appsv1.StatefulSet) {
	if in == nil {
		return
	}

	in.DNSConfig.UpdateStatefulSet(sts)

	if in.DNSPolicy != nil {
		sts.Spec.Template.Spec.DNSPolicy = *in.DNSPolicy
	}

	sts.Spec.Template.Spec.HostAliases = in.HostAliases
	sts.Spec.Template.Spec.HostNetwork = notNilBool(in.HostNetwork)
	sts.Spec.Template.Spec.Hostname = notNilString(in.Hostname)
}

// ----- PodDNSConfig -------------------------------------------------------

// PodDNSConfig defines the DNS parameters of a pod in addition to
// those generated from DNSPolicy.
// +k8s:openapi-gen=true
type PodDNSConfig struct {
	// A list of DNS name server IP addresses.
	// This will be appended to the base nameservers generated from DNSPolicy.
	// Duplicated nameservers will be removed.
	// +listType=atomic
	// +optional
	Nameservers []string `json:"nameservers,omitempty"`
	// A list of DNS search domains for host-name lookup.
	// This will be appended to the base search paths generated from DNSPolicy.
	// Duplicated search paths will be removed.
	// +listType=atomic
	// +optional
	Searches []string `json:"searches,omitempty"`
	// A list of DNS resolver options.
	// This will be merged with the base options generated from DNSPolicy.
	// Duplicated entries will be removed. Resolution options given in Options
	// will override those that appear in the base DNSPolicy.
	// +listType=map
	// +listMapKey=name
	// +optional
	Options []corev1.PodDNSConfigOption `json:"options,omitempty"`
}

func (in *PodDNSConfig) UpdateStatefulSet(sts *appsv1.StatefulSet) {
	if in == nil {
		return
	}

	cfg := corev1.PodDNSConfig{}

	if in.Nameservers != nil && len(in.Nameservers) > 0 {
		cfg.Nameservers = in.Nameservers
		sts.Spec.Template.Spec.DNSConfig = &cfg
	}

	if in.Searches != nil && len(in.Searches) > 0 {
		cfg.Searches = in.Searches
		sts.Spec.Template.Spec.DNSConfig = &cfg
	}

	if in.Options != nil && len(in.Options) > 0 {
		cfg.Options = in.Options
		sts.Spec.Template.Spec.DNSConfig = &cfg
	}
}

// ----- StartQuorum --------------------------------------------------------

// StartQuorum defines the order that deployments will be started in a Coherence cluster
// made up of multiple deployments.
// +k8s:openapi-gen=true
type StartQuorum struct {
	// The name of deployment that this deployment depends on.
	Deployment string `json:"deployment"`
	// The namespace that the deployment that this deployment depends on is installed into.
	// Default to the same namespace as this deployment
	// +optional
	Namespace string `json:"namespace,omitempty"`
	// The number of the Pods that should have been started before this
	// deployments will be started, defaults to all Pods for the deployment.
	// +optional
	PodCount int32 `json:"podCount,omitempty"`
}

// ----- StartQuorumStatus --------------------------------------------------

// StartQuorumStatus tracks the state of a deployment's start quorums.
type StartQuorumStatus struct {
	// The inlined start quorum.
	StartQuorum `json:",inline"`
	// Whether this quorum's condition has been met
	Ready bool `json:"ready"`
}

// ----- ConfigMapVolumeSpec ------------------------------------------------

// Represents a ConfigMap that will be added to the deployment's Pods as an
// additional Volume and as a VolumeMount in the containers.
// +coh:doc=misc_pod_settings/050_configmap_volumes.adoc,Add ConfigMap Volumes
// +k8s:openapi-gen=true
type ConfigMapVolumeSpec struct {
	// The name of the ConfigMap to mount.
	// This will also be used as the name of the Volume added to the Pod
	// if the VolumeName field is not set.
	Name string `json:"name"`
	// Path within the container at which the volume should be mounted.  Must
	// not contain ':'.
	MountPath string `json:"mountPath"`
	// The optional name to use for the Volume added to the Pod.
	// If not set, the ConfigMap name will be used as the VolumeName.
	// +optional
	VolumeName string `json:"volumeName,omitempty"`
	// Mounted read-only if true, read-write otherwise (false or unspecified).
	// Defaults to false.
	// +optional
	ReadOnly bool `json:"readOnly,omitempty"`
	// Path within the volume from which the container's volume should be mounted.
	// Defaults to "" (volume's root).
	// +optional
	SubPath string `json:"subPath,omitempty"`
	// mountPropagation determines how mounts are propagated from the host
	// to container and the other way around.
	// When not set, MountPropagationNone is used.
	// +optional
	MountPropagation *corev1.MountPropagationMode `json:"mountPropagation,omitempty"`
	// Expanded path within the volume from which the container's volume should be mounted.
	// Behaves similarly to SubPath but environment variable references $(VAR_NAME) are expanded using the container's environment.
	// Defaults to "" (volume's root).
	// SubPathExpr and SubPath are mutually exclusive.
	// +optional
	SubPathExpr string `json:"subPathExpr,omitempty"`
	// If unspecified, each key-value pair in the Data field of the referenced
	// ConfigMap will be projected into the volume as a file whose name is the
	// key and content is the value. If specified, the listed keys will be
	// projected into the specified paths, and unlisted keys will not be
	// present. If a key is specified which is not present in the ConfigMap,
	// the volume setup will error unless it is marked optional. Paths must be
	// relative and may not contain the '..' path or start with '..'.
	// +listType=map
	// +listMapKey=key
	// +optional
	Items []corev1.KeyToPath `json:"items,omitempty"`
	// Optional: mode bits to use on created files by default. Must be a
	// value between 0 and 0777. Defaults to 0644.
	// Directories within the path are not affected by this setting.
	// This might be in conflict with other options that affect the file
	// mode, like fsGroup, and the result can be other mode bits set.
	// +optional
	DefaultMode *int32 `json:"defaultMode,omitempty"`
	// Specify whether the ConfigMap or its keys must be defined
	// +optional
	Optional *bool `json:"optional,omitempty"`
}

// Add the Volume and VolumeMount for this ConfigMap spec.
func (in *ConfigMapVolumeSpec) AddVolumes(sts *appsv1.StatefulSet) {
	if in == nil {
		return
	}
	// Add the volume mount to the init-containers
	for i, c := range sts.Spec.Template.Spec.InitContainers {
		in.AddVolumeMounts(&c)
		// replace the updated container in the init-container array
		sts.Spec.Template.Spec.InitContainers[i] = c
	}
	// Add the volume mount to the containers
	for i, c := range sts.Spec.Template.Spec.Containers {
		in.AddVolumeMounts(&c)
		// replace the updated container in the container array
		sts.Spec.Template.Spec.Containers[i] = c
	}
	var volName string
	if in.VolumeName == "" {
		volName = in.Name
	} else {
		volName = in.VolumeName
	}
	vol := corev1.Volume{
		Name: volName,
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: in.Name,
				},
				Items:       in.Items,
				DefaultMode: in.DefaultMode,
				Optional:    in.Optional,
			},
		},
	}
	sts.Spec.Template.Spec.Volumes = append(sts.Spec.Template.Spec.Volumes, vol)
}

func (in *ConfigMapVolumeSpec) AddVolumeMounts(c *corev1.Container) {
	if in == nil {
		return
	}
	var volName string
	if in.VolumeName == "" {
		volName = in.Name
	} else {
		volName = in.VolumeName
	}
	vm := corev1.VolumeMount{
		Name:             volName,
		ReadOnly:         in.ReadOnly,
		MountPath:        in.MountPath,
		SubPath:          in.SubPath,
		MountPropagation: in.MountPropagation,
		SubPathExpr:      in.SubPathExpr,
	}
	c.VolumeMounts = append(c.VolumeMounts, vm)
}

// ----- SecretVolumeSpec ------------------------------------------------

// Represents a Secret that will be added to the deployment's Pods as an
// additional Volume and as a VolumeMount in the containers.
// +coh:doc=misc_pod_settings/020_secret_volumes.adoc,Add Secret Volumes
// +k8s:openapi-gen=true
type SecretVolumeSpec struct {
	// The name of the Secret to mount.
	// This will also be used as the name of the Volume added to the Pod
	// if the VolumeName field is not set.
	Name string `json:"name"`
	// Path within the container at which the volume should be mounted.  Must
	// not contain ':'.
	MountPath string `json:"mountPath"`
	// The optional name to use for the Volume added to the Pod.
	// If not set, the Secret name will be used as the VolumeName.
	// +optional
	VolumeName string `json:"volumeName,omitempty"`
	// Mounted read-only if true, read-write otherwise (false or unspecified).
	// Defaults to false.
	// +optional
	ReadOnly bool `json:"readOnly,omitempty"`
	// Path within the volume from which the container's volume should be mounted.
	// Defaults to "" (volume's root).
	// +optional
	SubPath string `json:"subPath,omitempty"`
	// mountPropagation determines how mounts are propagated from the host
	// to container and the other way around.
	// When not set, MountPropagationNone is used.
	// +optional
	MountPropagation *corev1.MountPropagationMode `json:"mountPropagation,omitempty"`
	// Expanded path within the volume from which the container's volume should be mounted.
	// Behaves similarly to SubPath but environment variable references $(VAR_NAME) are expanded using the container's environment.
	// Defaults to "" (volume's root).
	// SubPathExpr and SubPath are mutually exclusive.
	// +optional
	SubPathExpr string `json:"subPathExpr,omitempty"`
	// If unspecified, each key-value pair in the Data field of the referenced
	// Secret will be projected into the volume as a file whose name is the
	// key and content is the value. If specified, the listed keys will be
	// projected into the specified paths, and unlisted keys will not be
	// present. If a key is specified which is not present in the Secret,
	// the volume setup will error unless it is marked optional. Paths must be
	// relative and may not contain the '..' path or start with '..'.
	// +listType=map
	// +listMapKey=key
	// +optional
	Items []corev1.KeyToPath `json:"items,omitempty"`
	// Optional: mode bits to use on created files by default. Must be a
	// value between 0 and 0777. Defaults to 0644.
	// Directories within the path are not affected by this setting.
	// This might be in conflict with other options that affect the file
	// mode, like fsGroup, and the result can be other mode bits set.
	// +optional
	DefaultMode *int32 `json:"defaultMode,omitempty"`
	// Specify whether the Secret or its keys must be defined
	// +optional
	Optional *bool `json:"optional,omitempty"`
}

// Add the Volume and VolumeMount for this Secret spec.
func (in *SecretVolumeSpec) AddVolumes(sts *appsv1.StatefulSet) {
	if in == nil {
		return
	}
	// Add the volume mount to the init-containers
	for i, c := range sts.Spec.Template.Spec.InitContainers {
		in.AddVolumeMounts(&c)
		// replace the updated container in the init-container array
		sts.Spec.Template.Spec.InitContainers[i] = c
	}
	// Add the volume mount to the containers
	for i, c := range sts.Spec.Template.Spec.Containers {
		in.AddVolumeMounts(&c)
		// replace the updated container in the container array
		sts.Spec.Template.Spec.Containers[i] = c
	}
	var volName string
	if in.VolumeName == "" {
		volName = in.Name
	} else {
		volName = in.VolumeName
	}
	vol := corev1.Volume{
		Name: volName,
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName:  in.Name,
				Items:       in.Items,
				DefaultMode: in.DefaultMode,
				Optional:    in.Optional,
			},
		},
	}
	sts.Spec.Template.Spec.Volumes = append(sts.Spec.Template.Spec.Volumes, vol)
}

func (in *SecretVolumeSpec) AddVolumeMounts(c *corev1.Container) {
	if in == nil {
		return
	}
	var volName string
	if in.VolumeName == "" {
		volName = in.Name
	} else {
		volName = in.VolumeName
	}
	vm := corev1.VolumeMount{
		Name:             volName,
		ReadOnly:         in.ReadOnly,
		MountPath:        in.MountPath,
		SubPath:          in.SubPath,
		MountPropagation: in.MountPropagation,
		SubPathExpr:      in.SubPathExpr,
	}
	c.VolumeMounts = append(c.VolumeMounts, vm)
}

// ----- ResourceType -------------------------------------------------------

type ResourceType string

func (t ResourceType) Name() string {
	return string(t)
}

const (
	ResourceTypeDeployment     ResourceType = "Coherence"
	ResourceTypeConfigMap      ResourceType = "ConfigMap"
	ResourceTypeSecret         ResourceType = "Secret"
	ResourceTypeService        ResourceType = "Service"
	ResourceTypeServiceMonitor ResourceType = ServiceMonitorKind
	ResourceTypeStatefulSet    ResourceType = "StatefulSet"
)

func ToResourceType(kind string) (ResourceType, error) {
	var t ResourceType
	var err error

	switch kind {
	case ResourceTypeDeployment.Name():
		t = ResourceTypeDeployment
	case ResourceTypeConfigMap.Name():
		t = ResourceTypeConfigMap
	case ResourceTypeSecret.Name():
		t = ResourceTypeSecret
	case ResourceTypeService.Name():
		t = ResourceTypeService
	case ResourceTypeServiceMonitor.Name():
		t = ResourceTypeServiceMonitor
	case ResourceTypeStatefulSet.Name():
		t = ResourceTypeStatefulSet
	default:
		err = fmt.Errorf("attempt to obtain ResourceType unsupported kind %s", kind)
	}
	return t, err
}

func (t ResourceType) toObject() (runtime.Object, error) {
	var o runtime.Object
	var err error

	switch t {
	case ResourceTypeDeployment:
		o = &Coherence{}
	case ResourceTypeConfigMap:
		o = &corev1.ConfigMap{}
	case ResourceTypeSecret:
		o = &corev1.Secret{}
	case ResourceTypeService:
		o = &corev1.Service{}
	case ResourceTypeServiceMonitor:
		o = &monitoringv1.ServiceMonitor{}
	case ResourceTypeStatefulSet:
		o = &appsv1.StatefulSet{}
	default:
		err = fmt.Errorf("attempt to obtain runtime.Object for unsupported type %s", t)
	}
	return o, err
}

// ----- Resource -----------------------------------------------------------

type Resource struct {
	Kind ResourceType   `json:"kind"`
	Name string         `json:"name"`
	Spec runtime.Object `json:"spec"`
}

func (in *Resource) GetFullName() string {
	if in == nil {
		return ""
	}
	return fmt.Sprintf("%s_%s", in.Kind, in.Name)
}

func (in *Resource) IsDelete() bool {
	if in == nil {
		return false
	}
	// this resource is a delete if the Spec is nil
	return in.Spec == nil
}

// Set the the controller/owner of the resource
func (in *Resource) SetController(object runtime.Object, scheme *runtime.Scheme) error {
	if in == nil || in.Spec == nil {
		return nil
	}
	owner, ok := object.(metav1.Object)
	if !ok {
		return fmt.Errorf("%s is not a metav1.Template cannot call SetControllerReference", in.GetFullName())
	}
	m, ok := in.Spec.(metav1.Object)
	if !ok {
		return fmt.Errorf("%s is not a metav1.Template cannot call SetControllerReference", in.GetFullName())
	}
	if err := controllerutil.SetControllerReference(owner, m, scheme); err != nil {
		if _, owned := err.(*controllerutil.AlreadyOwnedError); !owned {
			// if the error is not an AlreadyOwnedError then return
			err = errors.Wrap(err, fmt.Sprintf("setting resource owner/controller in %s", in.GetFullName()))
			return err
		}
	}
	return nil
}

// ----- Resources ------------------------------------------------------

type Resources struct {
	Version int32      `json:"version"`
	Items   []Resource `json:"items"`
}

func (in Resources) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")
	buffer.WriteString(fmt.Sprintf(`"apiVersion":"%d"`, in.Version))
	buffer.WriteString(`, "kind": "Resources"`)
	buffer.WriteString(`, "items":[`)

	for i, item := range in.Items {
		if !item.IsDelete() {
			if i > 0 {
				buffer.WriteString(", ")
			}
			b, err := json.Marshal(item.Spec)
			if err != nil {
				return nil, err
			}
			buffer.Write(b)
		}
	}
	buffer.WriteString("]}")
	return buffer.Bytes(), nil
}

func (in *Resources) UnmarshalJSON(b []byte) error {
	l := unstructured.UnstructuredList{}
	if err := l.UnmarshalJSON(b); err != nil {
		return err
	}
	v, err := strconv.Atoi(l.GetAPIVersion())
	if err != nil {
		return err
	}
	in.Version = int32(v)
	for _, u := range l.Items {
		var o runtime.Object
		kind := u.GetObjectKind().GroupVersionKind().Kind
		t, err := ToResourceType(kind)
		if err != nil {
			return err
		}
		o, err = t.toObject()
		if err != nil {
			return err
		}
		if err = runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, o); err != nil {
			return errors.Wrapf(err, "converting unstructured to %s", kind)
		}

		in.Items = append(in.Items, Resource{
			Kind: t,
			Name: u.GetName(),
			Spec: o,
		})
	}
	return nil
}

func (in Resources) EnsureGVK(s *runtime.Scheme) {
	for _, r := range in.Items {
		switch {
		case !r.IsDelete() && r.Kind == ResourceTypeServiceMonitor:
			gvk := schema.GroupVersionKind{
				Group:   ServiceMonitorGroup,
				Version: ServiceMonitorVersion,
				Kind:    ServiceMonitorKind,
			}
			r.Spec.GetObjectKind().SetGroupVersionKind(gvk)
		case !r.IsDelete():
			gvks, _, _ := s.ObjectKinds(r.Spec)
			if len(gvks) > 0 {
				r.Spec.GetObjectKind().SetGroupVersionKind(gvks[0])
			}
		}
	}
}

func (in Resources) GetResource(kind ResourceType, name string) (Resource, bool) {
	for _, r := range in.Items {
		if r.Kind == kind && r.Name == name {
			return r, true
		}
	}
	return Resource{}, false
}

func (in Resources) GetResourcesOfKind(kind ResourceType) []Resource {
	var items []Resource
	for _, r := range in.Items {
		if r.Kind == kind {
			items = append(items, r)
		}
	}
	return items
}

// Obtain the diff between the previous deployment resources and this deployment resources.
func (in Resources) Diff(previous Resources) []Resource {
	var diff []Resource

	// work out any deletions
	for _, r := range previous.Items {
		_, found := in.GetResource(r.Kind, r.Name)
		if !found {
			// previous resource is deleted from this Resources
			diff = append(diff, Resource{Kind: r.Kind, Name: r.Name})
		}
	}

	// work out any additions or updates
	for _, r := range in.Items {
		prev, found := previous.GetResource(r.Kind, r.Name)
		if found {
			if len(deep.Equal(prev, r)) != 0 {
				// r and prev are different so there is something to update
				diff = append(diff, r)
			}
		} else {
			// r is a new resource so append it to the diff
			diff = append(diff, r)
		}
	}
	diff = append(diff, in.Items...)

	return diff
}

// Obtain the diff between the previous deployment resources of a specific kind and this deployment resources.
func (in Resources) DiffForKind(kind ResourceType, previous Resources) []Resource {
	var diff []Resource

	// work out any deletions
	for _, r := range previous.GetResourcesOfKind(kind) {
		_, found := in.GetResource(kind, r.Name)
		if !found {
			// previous resource is deleted from this Resources
			diff = append(diff, Resource{Kind: r.Kind, Name: r.Name})
		}
	}

	// work out any additions or updates
	for _, r := range in.GetResourcesOfKind(kind) {
		prev, found := previous.GetResource(r.Kind, r.Name)
		if found {
			if len(deep.Equal(prev, r)) != 0 {
				// r and prev are different so there is something to update
				diff = append(diff, r)
			}
		} else {
			// r is a new resource so append it to the diff
			diff = append(diff, r)
		}
	}

	return diff
}

// Set the deployment as the controller/owner of all of the resources
func (in Resources) SetController(object runtime.Object, scheme *runtime.Scheme) error {
	for _, r := range in.Items {
		if err := r.SetController(object, scheme); err != nil {
			return err
		}
	}
	return nil
}

// Create the specified resource
func (in Resources) Create(kind ResourceType, name string, mgr manager.Manager, logger logr.Logger) error {
	logger.Info(fmt.Sprintf("Creating %s for deployment", kind))
	// Get the resource state
	resource, found := in.GetResource(kind, name)
	if !found {
		return fmt.Errorf("cannot create %s for deployment %s as state not present in store", kind, name)
	}
	// create the resource
	if err := mgr.GetClient().Create(context.TODO(), resource.Spec); err != nil {
		return errors.Wrapf(err, "failed to create %s", kind)
	}
	return nil
}

// ----- helper methods -----------------------------------------------------

// Convert an int32 pointer to a string using the default if the pointer is nil.
func Int32PtrToStringWithDefault(i *int32, d int32) string {
	if i == nil {
		return Int32ToString(d)
	}
	return Int32ToString(*i)
}

// Convert an int32 pointer to a string.
func Int32PtrToString(i *int32) string {
	return Int32ToString(*i)
}

// Convert an int32 to a string.
func Int32ToString(i int32) string {
	return strconv.FormatInt(int64(i), 10)
}

// Convert a bool pointer to a string.
func BoolPtrToString(b *bool) string {
	if b != nil && *b {
		return "true"
	}
	return "false"
}
