/*
 * Copyright (c) 2020, 2025, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// CoherenceResource is a common interface implemented by different coherence resources.
// +kubebuilder:object:generate=false
type CoherenceResource interface {
	client.Object
	// GetCoherenceClusterName obtains the Coherence cluster name for the Coherence resource.
	GetCoherenceClusterName() string
	// GetWkaServiceName returns the name of the headless Service used for Coherence WKA.
	GetWkaServiceName() string
	// GetWkaIPFamily returns the IP Family of the headless Service used for Coherence WKA.
	GetWkaIPFamily() corev1.IPFamily
	// GetHeadlessServiceName returns the name of the headless Service used for the StatefulSet.
	GetHeadlessServiceName() string
	// GetHeadlessServiceIPFamily always returns an empty array as this is not applicable to Jobs.
	GetHeadlessServiceIPFamily() []corev1.IPFamily
	// GetReplicas returns the number of replicas required for a deployment.
	// The Replicas field is a pointer and may be nil so this method will
	// return either the actual Replicas value or the default (DefaultReplicas const)
	// if the Replicas field is nil.
	GetReplicas() int32
	// SetReplicas sets the number of replicas required for a deployment.
	SetReplicas(replicas int32)
	// FindFullyQualifiedPortServiceNames returns a map of the exposed ports of this resource mapped to their Service's
	// fully qualified domain name.
	FindFullyQualifiedPortServiceNames() map[string]string
	// FindFullyQualifiedPortServiceName returns the fully qualified name of the Service used to expose
	// a named port and a bool indicating whether the named port has a Service.
	FindFullyQualifiedPortServiceName(name string) (string, bool)
	// FindPortServiceNames returns a map of the port names for this resource mapped to their Service names.
	FindPortServiceNames() map[string]string
	// FindPortServiceName returns the name of the Service used to expose a named port and a bool indicating
	// whether the named port has a Service.
	FindPortServiceName(name string) (string, bool)
	// CreateCommonLabels creates the deployment's common label set.
	CreateCommonLabels() map[string]string
	// CreateGlobalLabels creates the common label set for all resources.
	CreateGlobalLabels() map[string]string
	// CreateGlobalAnnotations creates the common annotation set for all resources.
	CreateGlobalAnnotations() map[string]string
	// CreateAnnotations returns the annotations to apply to this cluster's
	// deployment (StatefulSet).
	CreateAnnotations() map[string]string
	// GetNamespacedName returns the namespace/name key to look up this resource.
	GetNamespacedName() types.NamespacedName
	// GetRoleName returns the role name for a deployment.
	// If the Spec.Role field is set that is used for the role name
	// otherwise the deployment name is used as the role name.
	GetRoleName() string
	// GetType returns the type for a deployment.
	GetType() CoherenceType
	// GetWKA returns the host name Coherence should for WKA.
	GetWKA() string
	// GetVersionAnnotation if the returns the value of the Operator version annotation and true,
	// if the version annotation is present. If the version annotation is not present this method
	// returns empty string and false.
	GetVersionAnnotation() (string, bool)
	// IsBeforeVersion returns true if this Coherence resource Operator version annotation value is
	// before the specified version, or is not set.
	// The version parameter must be a valid SemVer value.
	IsBeforeVersion(version string) bool
	// IsBeforeOrSameVersion returns true if this Coherence resource Operator version annotation value is
	// the same or before the specified version, or is not set.
	// The version parameter must be a valid SemVer value.
	IsBeforeOrSameVersion(version string) bool
	// GetSpec returns this resource's CoherenceResourceSpec
	GetSpec() *CoherenceResourceSpec
	// GetJobResourceSpec returns this resource's CoherenceJobResourceSpec.
	// If the spec is not a CoherenceJobResourceSpec the bool return value will be false.
	GetJobResourceSpec() (*CoherenceJobResourceSpec, bool)
	// GetStatefulSetSpec returns this resource's CoherenceStatefulSetResourceSpec
	// If the spec is not a CoherenceStatefulSetResourceSpec the bool return value will be false.
	GetStatefulSetSpec() (*CoherenceStatefulSetResourceSpec, bool)
	// GetStatus returns this resource's CoherenceResourceStatus
	GetStatus() *CoherenceResourceStatus
	// AddAnnotation adds an annotation to this resource
	AddAnnotation(key, value string)
	// AddAnnotationIfMissing adds an annotation to this resource if it is not already present
	AddAnnotationIfMissing(key, value string)
	// GetAnnotations returns the annotations on this resource
	GetAnnotations() map[string]string
	// CreateKubernetesResources creates the kubernetes resources defined by this resource
	CreateKubernetesResources() (Resources, error)
	// DeepCopyObject copies this CoherenceResource.
	DeepCopyObject() runtime.Object
	// DeepCopyResource copies this CoherenceResource.
	DeepCopyResource() CoherenceResource
	// GetAPIVersion returns the TypeMeta API version
	GetAPIVersion() string
	// IsForceExit is a flag to determine whether the Operator calls System.exit when the main class finishes.
	IsForceExit() bool
	// GetEnvVarFrom returns the array of EnvVarSource configurations
	GetEnvVarFrom() []corev1.EnvFromSource
	// GetGlobalSpec returns the attributes to be applied to all resources
	GetGlobalSpec() *GlobalSpec
	// GetInitResources returns the optional resource requirements for the init container
	GetInitResources() *corev1.ResourceRequirements
	// GetGenerationString returns the resource metadata generation as a string
	GetGenerationString() string
	// HashLabelMatches determines whether the hash label on the specified metav1.Object
	// matched the generation string of this resource
	HashLabelMatches(m metav1.Object) bool
	// UpdateStatusVersion update the version field of the status
	UpdateStatusVersion(v string)
}
