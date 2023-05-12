/*
 * Copyright (c) 2019, 2023, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1

import (
	"fmt"
	"golang.org/x/mod/semver"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"
)

// Coherence resource Condition Types
// The different eight types of state that a deployment may be in.
//
// Transitions are:
// Initialized    -> Waiting
//
//	-> Created
//
// Waiting        -> Created
// Created        -> Ready
//
//	-> Stopped
//
// Ready          -> Scaling
//
//	-> RollingUpgrade
//	-> Stopped
//
// Scaling        -> Ready
//
//	-> Stopped
//
// RollingUpgrade -> Ready
// Stopped        -> Created
const (
	ConditionTypeInitialized    ConditionType = "Initialized"
	ConditionTypeWaiting        ConditionType = "Waiting"
	ConditionTypeCreated        ConditionType = "Created"
	ConditionTypeReady          ConditionType = "Ready"
	ConditionTypeScaling        ConditionType = "Scaling"
	ConditionTypeRollingUpgrade ConditionType = "RollingUpgrade"
	ConditionTypeFailed         ConditionType = "Failed"
	ConditionTypeStopped        ConditionType = "Stopped"
	ConditionTypeCompleted      ConditionType = "Completed"

	CoherenceTypeUnknown     CoherenceType = "Unknown"
	CoherenceTypeStatefulSet CoherenceType = "StatefulSet"
	CoherenceTypeJob         CoherenceType = "Job"
)

type CoherenceType string

// The package init function that will automatically register the Coherence resource types with
// the default k8s Scheme.
func init() {
	SchemeBuilder.Register(&Coherence{}, &CoherenceList{}, &CoherenceJob{}, &CoherenceJobList{})
}

// ----- Coherence type ------------------------------------------------------------------

var _ CoherenceResource = &Coherence{}

// Coherence is the top level schema for the Coherence API and custom resource definition (CRD).
//
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:subresource:scale:specpath=.spec.replicas,statuspath=.status.replicas,selectorpath=.status.selector
// +kubebuilder:resource:path=coherence,scope=Namespaced,shortName=coh,categories=coherence
// +kubebuilder:printcolumn:name="Cluster",type="string",JSONPath=".status.coherenceCluster",description="The name of the Coherence cluster that this deployment belongs to"
// +kubebuilder:printcolumn:name="Role",type="string",JSONPath=".status.role",description="The role of this deployment in a Coherence cluster"
// +kubebuilder:printcolumn:name="Replicas",type="integer",JSONPath=".status.replicas",description="The number of Coherence deployments for this deployment"
// +kubebuilder:printcolumn:name="Ready",type="integer",JSONPath=".status.readyReplicas",description="The number of ready Coherence deployments for this deployment"
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase",description="The status of this deployment"
// +kubebuilder:printcolumn:name="Image",priority=1,type="string",JSONPath=".spec.image",description="The image name"
type Coherence struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CoherenceStatefulSetResourceSpec `json:"spec,omitempty"`
	Status CoherenceResourceStatus          `json:"status,omitempty"`
}

// GetCoherenceClusterName obtains the Coherence cluster name for the Coherence resource.
func (in *Coherence) GetCoherenceClusterName() string {
	if in == nil {
		return ""
	}

	if in.Spec.Cluster == nil {
		return in.Name
	}
	return *in.Spec.Cluster
}

// GetSpec returns this resource's CoherenceResourceSpec
func (in *Coherence) GetSpec() *CoherenceResourceSpec {
	return &in.Spec.CoherenceResourceSpec
}

// GetStatefulSetSpec returns this resource's CoherenceStatefulSetResourceSpec
func (in *Coherence) GetStatefulSetSpec() (*CoherenceStatefulSetResourceSpec, bool) {
	return &in.Spec, true
}

// GetJobResourceSpec always returns nil and false
func (in *Coherence) GetJobResourceSpec() (*CoherenceJobResourceSpec, bool) {
	return nil, false
}

// GetStatus returns this resource's CoherenceResourceSpec
func (in *Coherence) GetStatus() *CoherenceResourceStatus {
	return &in.Status
}

// CreateKubernetesResources returns this resource's CoherenceResourceSpec
func (in *Coherence) CreateKubernetesResources() (Resources, error) {
	res := in.Spec.CreateKubernetesResources(in)

	// Create the headless Service
	res = append(res, in.Spec.CreateHeadlessService(in))
	// Create the StatefulSet
	res = append(res, in.Spec.CreateStatefulSetResource(in))
	return Resources{Items: res}, nil
}

// GetAPIVersion returns the TypeMeta API version
func (in *Coherence) GetAPIVersion() string {
	return in.APIVersion
}

func (in *Coherence) DeepCopyResource() CoherenceResource {
	return in.DeepCopy()
}

// GetWkaServiceName returns the name of the headless Service used for Coherence WKA.
func (in *Coherence) GetWkaServiceName() string {
	if in == nil {
		return ""
	}
	return in.Name + WKAServiceNameSuffix
}

// GetHeadlessServiceName returns the name of the headless Service used for the StatefulSet.
func (in *Coherence) GetHeadlessServiceName() string {
	if in == nil {
		return ""
	}
	return in.Name + HeadlessServiceNameSuffix
}

// GetReplicas returns the number of replicas required for a deployment.
// The Replicas field is a pointer and may be nil so this method will
// return either the actual Replicas value or the default (DefaultReplicas const)
// if the Replicas field is nil.
func (in *Coherence) GetReplicas() int32 {
	if in == nil || in.Spec.Replicas == nil {
		return DefaultReplicas
	}
	return in.Spec.GetReplicas()
}

// SetReplicas sets the number of replicas required for a deployment.
func (in *Coherence) SetReplicas(replicas int32) {
	if in != nil {
		in.Spec.CoherenceResourceSpec.Replicas = &replicas
	}
}

// FindFullyQualifiedPortServiceNames returns a map of the exposed ports of this resource mapped to their Service's
// fully qualified domain name.
func (in *Coherence) FindFullyQualifiedPortServiceNames() map[string]string {
	if in == nil {
		return make(map[string]string)
	}
	m := in.Spec.FindPortServiceNames(in)
	for k, v := range m {
		m[k] = v + "." + in.GetNamespace() + ".svc"
	}
	return m
}

// FindFullyQualifiedPortServiceName returns the fully qualified name of the Service used to expose a named port and a bool indicating
// whether the named port has a Service.
func (in *Coherence) FindFullyQualifiedPortServiceName(name string) (string, bool) {
	n, found := in.FindPortServiceName(name)
	if found {
		n = n + "." + in.GetNamespace() + ".svc"
	}
	return n, found
}

// FindPortServiceNames returns a map of the port names for this resource mapped to their Service names.
func (in *Coherence) FindPortServiceNames() map[string]string {
	if in == nil {
		return make(map[string]string)
	}
	return in.Spec.FindPortServiceNames(in)
}

// FindPortServiceName returns the name of the Service used to expose a named port and a bool indicating
// whether the named port has a Service.
func (in *Coherence) FindPortServiceName(name string) (string, bool) {
	if in == nil {
		return "", false
	}
	return in.Spec.FindPortServiceName(name, in)
}

// CreateCommonLabels creates the deployment's common label set.
func (in *Coherence) CreateCommonLabels() map[string]string {
	labels := make(map[string]string)
	labels[LabelCoherenceDeployment] = in.Name
	labels[LabelCoherenceCluster] = in.GetCoherenceClusterName()
	labels[LabelCoherenceRole] = in.GetRoleName()

	if in.Spec.AppLabel != nil {
		labels[LabelApp] = *in.Spec.AppLabel
	}

	if in.Spec.VersionLabel != nil {
		labels[LabelVersion] = *in.Spec.VersionLabel
	}

	return labels
}

// CreateAnnotations returns the annotations to apply to this cluster's
// deployment (StatefulSet).
func (in *Coherence) CreateAnnotations() map[string]string {
	var annotations map[string]string
	if in.Spec.StatefulSetAnnotations != nil {
		annotations = make(map[string]string)
		for k, v := range in.Spec.StatefulSetAnnotations {
			annotations[k] = v
		}
	} else if in.Annotations != nil {
		annotations = make(map[string]string)
		for k, v := range in.Annotations {
			annotations[k] = v
		}
	}
	return annotations
}

func (in *Coherence) AddAnnotation(key, value string) {
	if in != nil {
		if in.Annotations == nil {
			in.Annotations = make(map[string]string)
		}
		in.Annotations[key] = value
	}
}

// GetNamespacedName returns the namespace/name key to look up this resource.
func (in *Coherence) GetNamespacedName() types.NamespacedName {
	return types.NamespacedName{
		Namespace: in.Namespace,
		Name:      in.Name,
	}
}

// GetRoleName returns the role name for a deployment.
// If the Spec.Role field is set that is used for the role name
// otherwise the deployment name is used as the role name.
func (in *Coherence) GetRoleName() string {
	switch {
	case in == nil:
		return ""
	case in.Spec.Role != "":
		return in.Spec.Role
	default:
		return in.Name
	}
}

// GetType returns the type for a deployment.
func (in *Coherence) GetType() CoherenceType {
	return CoherenceTypeStatefulSet
}

// GetWKA returns the host name Coherence should for WKA.
func (in *Coherence) GetWKA() string {
	if in == nil {
		return ""
	}
	return in.Spec.Coherence.GetWKA(in)
}

// GetVersionAnnotation if the returns the value of the Operator version annotation and true,
// if the version annotation is present. If the version annotation is not present this method
// returns empty string and false.
func (in *Coherence) GetVersionAnnotation() (string, bool) {
	if in == nil || in.Annotations == nil {
		return "", false
	}
	version, found := in.Annotations[AnnotationOperatorVersion]
	return version, found
}

// IsBeforeVersion returns true if this Coherence resource Operator version annotation value is
// before the specified version, or is not set.
// The version parameter must be a valid SemVer value.
func (in *Coherence) IsBeforeVersion(version string) bool {
	if actual, found := in.GetVersionAnnotation(); found {
		return semver.Compare(actual, version) < 0
	}
	return true
}

// ----- CoherenceStatefulSetResourceSpec type -----------------------------------------------------

// CoherenceStatefulSetResourceSpec defines the specification of a Coherence resource. A Coherence resource is
// typically one or more Pods that perform the same functionality, for example storage members.
// +k8s:openapi-gen=true
type CoherenceStatefulSetResourceSpec struct {
	CoherenceResourceSpec `json:",inline"`
	// StatefulSetAnnotations are free-form yaml that will be added to the Coherence cluster
	// `StatefulSet` as annotations.
	// Any annotations should be placed BELOW this "annotations:" key, for example:
	//
	// The default behaviour is to copy all annotations from the `Coherence` resource to the
	// `StatefulSet`, specifying any annotations in the `StatefulSetAnnotations` will override
	// this behaviour and only include the `StatefulSetAnnotations`.
	//
	// annotations:
	//   foo.io/one: "value1"
	//   foo.io/two: "value2"
	//
	// see: https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/[Kubernetes Annotations]
	// +optional
	StatefulSetAnnotations map[string]string `json:"statefulSetAnnotations,omitempty"`
	// VolumeClaimTemplates defines extra PVC mappings that will be added to the Coherence Pod.
	// The content of this yaml should match the normal k8s volumeClaimTemplates section of a StatefulSet spec
	// as described in https://kubernetes.io/docs/concepts/storage/persistent-volumes/
	// Every claim in this list must have at least one matching (by name) volumeMount in one
	// container in the template. A claim in this list takes precedence over any volumes in the
	// template, with the same name.
	// +listType=atomic
	// +optional
	VolumeClaimTemplates []PersistentVolumeClaim `json:"volumeClaimTemplates,omitempty"`
	// The configuration to control safe scaling.
	// +optional
	Scaling *ScalingSpec `json:"scaling,omitempty"`
	// The configuration of the probe used to signal that services must be suspended
	// before a deployment is stopped.
	// +optional
	SuspendProbe *Probe `json:"suspendProbe,omitempty"`
	// A flag controlling whether storage enabled cache services in this deployment
	// will be suspended before the deployment is shutdown or scaled to zero.
	// The action of suspending storage enabled services when the whole deployment is being
	// stopped ensures that cache services with persistence enabled will shut down cleanly
	// without the possibility of Coherence trying to recover and re-balance partitions
	// as Pods are stopped.
	// The default value if not specified is true.
	// +optional
	SuspendServicesOnShutdown *bool `json:"suspendServicesOnShutdown,omitempty"`
	// ResumeServicesOnStartup allows the Operator to resume suspended Coherence services when
	// the Coherence container is started. This only applies to storage enabled distributed cache
	// services. This ensures that services that are suspended due to the shutdown of a storage
	// tier, but those services are still running (albeit suspended) in other storage disabled
	// deployments, will be resumed when storage comes back.
	// Note that starting Pods with suspended partitioned cache services may stop the Pod reaching the ready state.
	// The default value if not specified is true.
	// +optional
	ResumeServicesOnStartup *bool `json:"resumeServicesOnStartup,omitempty"`
	// AutoResumeServices is a map of Coherence service names to allow more fine-grained control over
	// which services may be auto-resumed by the operator when a Coherence Pod starts.
	// The key to the map is the name of the Coherence service. This should be the fully qualified name
	// if scoped services are being used in Coherence. The value is a bool, set to `true` to allow the
	// service to be auto-resumed or `false` to not allow the service to be auto-resumed.
	// Adding service names to this list will override any value set in `ResumeServicesOnStartup`, so if the
	// `ResumeServicesOnStartup` field is `false` but there are service names in the `AutoResumeServices`, mapped
	// to `true`, those services will still be resumed.
	// Note that starting Pods with suspended partitioned cache services may stop the Pod reaching the ready state.
	// +optional
	AutoResumeServices map[string]bool `json:"autoResumeServices,omitempty"`
	// SuspendServiceTimeout sets the number of seconds to wait for the service suspend
	// call to return (the default is 60 seconds)
	// +optional
	SuspendServiceTimeout *int `json:"suspendServiceTimeout,omitempty"`
	// Whether to perform a StatusHA test on the cluster before performing an update or deletion.
	// This field can be set to "false" to force through an update even when a Coherence deployment is in
	// an unstable state.
	// The default is true, to always check for StatusHA before updating a Coherence deployment.
	// +optional
	HABeforeUpdate *bool `json:"haBeforeUpdate,omitempty"`
	// AllowUnsafeDelete controls whether the Operator will add a finalizer to the Coherence resource
	// so that it can intercept deletion of the resource and initiate a controlled shutdown of the
	// Coherence cluster. The default value is `false`.
	// The primary use for setting this flag to `true` is in CI/CD environments so that cleanup jobs
	// can delete a whole namespace without requiring the Operator to have removed finalizers from
	// any Coherence resources deployed into that namespace.
	// It is not recommended to set this flag to `true` in a production environment, especially when
	// using Coherence persistence features.
	// +optional
	AllowUnsafeDelete *bool `json:"allowUnsafeDelete,omitempty"`
	// Actions to execute once all the Pods are ready after an initial deployment
	// +optional
	Actions []Action `json:"actions,omitempty"`
}

// CreateStatefulSetResource creates the deployment's StatefulSet resource.
func (in *CoherenceStatefulSetResourceSpec) CreateStatefulSetResource(deployment *Coherence) Resource {
	sts := in.CreateStatefulSet(deployment)

	return Resource{
		Kind: ResourceTypeStatefulSet,
		Name: sts.GetName(),
		Spec: &sts,
	}
}

// CreateStatefulSet creates the deployment's StatefulSet.
func (in *CoherenceStatefulSetResourceSpec) CreateStatefulSet(deployment *Coherence) appsv1.StatefulSet {
	sts := appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:   deployment.GetNamespace(),
			Name:        deployment.GetName(),
			Labels:      deployment.CreateCommonLabels(),
			Annotations: deployment.CreateAnnotations(),
		},
	}

	replicas := in.GetReplicas()
	podTemplate := in.CreatePodTemplateSpec(deployment)

	// Add the component label
	sts.Labels[LabelComponent] = LabelComponentCoherenceStatefulSet
	sts.Spec = appsv1.StatefulSetSpec{
		Replicas:            &replicas,
		PodManagementPolicy: appsv1.ParallelPodManagement,
		UpdateStrategy: appsv1.StatefulSetUpdateStrategy{
			Type: appsv1.RollingUpdateStatefulSetStrategyType,
		},
		RevisionHistoryLimit: pointer.Int32(5),
		ServiceName:          deployment.GetHeadlessServiceName(),
		Selector: &metav1.LabelSelector{
			MatchLabels: in.CreatePodSelectorLabels(deployment),
		},
		Template: podTemplate,
	}

	// Add any Coherence settings
	in.Coherence.UpdateStatefulSet(deployment, &sts)

	// append any additional PVCs
	for _, v := range in.VolumeClaimTemplates {
		sts.Spec.VolumeClaimTemplates = append(sts.Spec.VolumeClaimTemplates, v.ToPVC())
	}

	return sts
}

// CheckHABeforeUpdate returns true if a StatusHA check should be made before updating a deployment.
func (in *CoherenceStatefulSetResourceSpec) CheckHABeforeUpdate() bool {
	return in.HABeforeUpdate == nil || *in.HABeforeUpdate
}

// IsSuspendServicesOnShutdown returns true if services should be suspended before a cluster is shutdown.
func (in *CoherenceStatefulSetResourceSpec) IsSuspendServicesOnShutdown() bool {
	return in.SuspendServicesOnShutdown == nil || *in.SuspendServicesOnShutdown
}

// GetEffectiveScalingPolicy returns the scaling policy to be used.
func (in *CoherenceStatefulSetResourceSpec) GetEffectiveScalingPolicy() ScalingPolicy {
	if in == nil {
		return SafeScaling
	}

	var policy ScalingPolicy

	if in.Scaling == nil || in.Scaling.Policy == nil {
		// the scaling policy is not set the look at the storage enabled flag
		if in.Coherence == nil || in.Coherence.StorageEnabled == nil || *in.Coherence.StorageEnabled {
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

// GetScalingProbe returns the Probe to use for checking Phase HA for the deployment.
// This method will not return nil.
func (in *CoherenceStatefulSetResourceSpec) GetScalingProbe() *Probe {
	if in == nil || in.Scaling == nil || in.Scaling.Probe == nil {
		return in.GetDefaultScalingProbe()
	}
	return in.Scaling.Probe
}

// GetSuspendProbe returns the Probe to use for signaling to a deployment that services should be suspended
// prior to the deployment being stopped.
// This method will not return nil.
func (in *CoherenceStatefulSetResourceSpec) GetSuspendProbe() *Probe {
	if in == nil || in.SuspendProbe == nil {
		return in.GetDefaultSuspendProbe()
	}
	return in.SuspendProbe
}

// GetDefaultSuspendProbe returns the default Suspend probe
func (in *CoherenceStatefulSetResourceSpec) GetDefaultSuspendProbe() *Probe {
	timeout := in.SuspendServiceTimeout
	if timeout == nil {
		oneMinute := 60
		timeout = &oneMinute
	}

	probe := Probe{
		TimeoutSeconds: timeout,
		ProbeHandler: corev1.ProbeHandler{
			HTTPGet: &corev1.HTTPGetAction{
				Path: "/suspend",
				Port: intstr.FromString(PortNameHealth),
			},
		},
	}

	return probe.DeepCopy()
}

// ----- CoherenceList type ------------------------------------------------------------------------

// +kubebuilder:object:root=true

// CoherenceList is a list of Coherence resources.
type CoherenceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Coherence `json:"items"`
}

// ----- CoherenceResourceStatus type --------------------------------------------------------------

// CoherenceResourceStatus defines the observed state of Coherence resource.
type CoherenceResourceStatus struct {
	// The phase of a Coherence resource is a simple, high-level summary of where the
	// Coherence resource is in its lifecycle.
	// The conditions array, the reason and message fields, and the individual container status
	// arrays contain more detail about the pod's status.
	// There are eight possible phase values:
	//
	// Initialized:    The deployment has been accepted by the Kubernetes system.
	// Created:        The deployments secondary resources, (e.g. the StatefulSet, Services etc.) have been created.
	// Ready:          The StatefulSet for the deployment has the correct number of replicas and ready replicas.
	// Waiting:        The deployment's start quorum conditions have not yet been met.
	// Scaling:        The number of replicas in the deployment is being scaled up or down.
	// RollingUpgrade: The StatefulSet is performing a rolling upgrade.
	// Stopped:        The replica count has been set to zero.
	// Completed:      The Coherence resource is running a Job and the Job has completed.
	// Failed:         An error occurred reconciling the deployment and its secondary resources.
	//
	// +optional
	Phase ConditionType `json:"phase,omitempty"`
	// The name of the Coherence cluster that this deployment is part of.
	// +optional
	CoherenceCluster string `json:"coherenceCluster,omitempty"`
	// The type of the Coherence resource.
	// +optional
	Type CoherenceType `json:"type,omitempty"`
	// Replicas is the desired number of members in the Coherence deployment
	// represented by the Coherence resource.
	// +optional
	Replicas int32 `json:"replicas"`
	// CurrentReplicas is the current number of members in the Coherence deployment
	// represented by the Coherence resource.
	// +optional
	CurrentReplicas int32 `json:"currentReplicas"`
	// ReadyReplicas is the number of members in the Coherence deployment
	// represented by the Coherence resource that are in the ready state.
	// +optional
	ReadyReplicas int32 `json:"readyReplicas"`
	// When the Coherence resource is running a Job, the number of pending and running pods.
	// +optional
	Active int32 `json:"active,omitempty"`
	// When the Coherence resource is running a Job, the number of pods which reached phase Succeeded.
	// +optional
	Succeeded int32 `json:"succeeded,omitempty"`
	// When the Coherence resource is running a Job, the number of pods which reached phase Failed.
	// +optional
	Failed int32 `json:"failed,omitempty"`
	// The effective role name for this deployment.
	// This will come from the Spec.Role field if set otherwise the deployment name
	// will be used for the role name
	// +optional
	Role string `json:"role,omitempty"`
	// label query over deployments that should match the replicas count. This is same
	// as the label selector but in the string format to avoid introspection
	// by clients. The string will be in the same format as the query-param syntax.
	// More info about label selectors: https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/
	// +optional
	Selector string `json:"selector,omitempty"`
	// The status conditions.
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	Conditions Conditions `json:"conditions,omitempty"`
	// Hash is the hash of the latest applied Coherence spec
	// +optional
	Hash string `json:"hash,omitempty"`
	// ActionsExecuted tracks whether actions were executed
	// +optional
	ActionsExecuted bool `json:"actionsExecuted,omitempty"`
}

// UpdatePhase updates the current Phase
// TODO not used?
func (in *CoherenceResourceStatus) UpdatePhase(deployment *Coherence, phase ConditionType) bool {
	return in.SetCondition(deployment, Condition{Type: phase, Status: corev1.ConditionTrue})
}

// SetCondition sets the current Status Condition
func (in *CoherenceResourceStatus) SetCondition(deployment *Coherence, c Condition) bool {
	deployment.Status.DeepCopyInto(in)
	updated := in.ensureInitialized(deployment)
	if in.Phase != "" && in.Phase == c.Type {
		// already at the desired phase
		return updated
	}
	// set the requested condition's type as the current phase
	updated = in.setPhase(c.Type) || updated
	return updated
}

// Update the status based on the condition of the StatefulSet status.
func (in *CoherenceResourceStatus) Update(deployment *Coherence, sts *appsv1.StatefulSetStatus) bool {
	// ensure that there is an Initialized condition
	updated := in.ensureInitialized(deployment)

	if sts != nil {
		// update CurrentReplicas from StatefulSet if required
		if in.CurrentReplicas != sts.CurrentReplicas {
			in.CurrentReplicas = sts.CurrentReplicas
			updated = true
		}

		// update ReadyReplicas from StatefulSet if required
		if in.ReadyReplicas != sts.ReadyReplicas {
			in.ReadyReplicas = sts.ReadyReplicas
			updated = true
		}

		if sts.CurrentRevision == sts.UpdateRevision {
			// both revisions are the same so the StatefulSet is not updating
			// If the current phase is not Ready check to see whether it should be ready.
			if in.Phase != ConditionTypeReady && in.Replicas == in.ReadyReplicas && in.Replicas == in.CurrentReplicas {
				updated = in.setPhase(ConditionTypeReady)
			}
		} else {
			// the revisions are different so the StatefulSet is updating, ensure the phase is set correctly
			if in.Phase != ConditionTypeRollingUpgrade {
				updated = in.setPhase(ConditionTypeRollingUpgrade)
			}
		}
	} else {
		// update CurrentReplicas to zero
		if in.CurrentReplicas != 0 {
			in.CurrentReplicas = 0
			updated = true
		}
		// update ReadyReplicas to zero
		if in.ReadyReplicas != 0 {
			in.ReadyReplicas = 0
			updated = true
		}
	}

	if deployment.Spec.GetReplicas() == 0 {
		// scaled to zero
		if in.Phase != ConditionTypeStopped {
			updated = in.setPhase(ConditionTypeStopped)
		}
	}

	return updated
}

// UpdateFromJob the status based on the condition of the Job status.
func (in *CoherenceResourceStatus) UpdateFromJob(deployment *Coherence, jobStatus *batchv1.JobStatus) bool {
	// ensure that there is an Initialized condition
	updated := in.ensureInitialized(deployment)

	fmt.Printf("***** UpdateFromJob %v\n", jobStatus)

	if jobStatus != nil {
		count := jobStatus.Active + jobStatus.Succeeded
		// update CurrentReplicas from Job if required
		if in.CurrentReplicas != count {
			in.CurrentReplicas = count
			updated = true
		}

		// update ReadyReplicas from Job if required
		if jobStatus.Ready != nil && in.ReadyReplicas != *jobStatus.Ready {
			in.ReadyReplicas = *jobStatus.Ready
			updated = true
		}

		if in.Phase != ConditionTypeReady && in.Replicas == in.ReadyReplicas && in.Replicas == in.CurrentReplicas {
			updated = in.setPhase(ConditionTypeReady)
		}

		if jobStatus.CompletionTime != nil {
			updated = in.setPhase(ConditionTypeCompleted)
		}

		if in.Active != jobStatus.Active {
			in.Active = jobStatus.Active
			updated = true
		}

		if in.Succeeded != jobStatus.Succeeded {
			in.Succeeded = jobStatus.Succeeded
			updated = true
		}

		if in.Failed != jobStatus.Failed {
			in.Failed = jobStatus.Failed
			updated = true
		}
	} else {
		// update CurrentReplicas to zero
		if in.CurrentReplicas != 0 {
			in.CurrentReplicas = 0
			updated = true
		}
		// update ReadyReplicas to zero
		if in.ReadyReplicas != 0 {
			in.ReadyReplicas = 0
			updated = true
		}

		if in.Active != 0 {
			in.Active = 0
			updated = true
		}

		if in.Succeeded != 0 {
			in.Succeeded = 0
			updated = true
		}

		if in.Failed != 0 {
			in.Failed = 0
			updated = true
		}
	}

	if deployment.Spec.GetReplicas() == 0 {
		// scaled to zero
		if in.Phase != ConditionTypeStopped {
			updated = in.setPhase(ConditionTypeStopped)
		}
	}

	return updated
}

// set a status phase.
func (in *CoherenceResourceStatus) setPhase(phase ConditionType) bool {
	if in.Phase == phase {
		return false
	}

	switch {
	case in.Phase == ConditionTypeReady && phase != ConditionTypeReady:
		// we're transitioning out of Ready state
		in.Conditions.SetCondition(Condition{Type: ConditionTypeReady, Status: corev1.ConditionFalse})
	case in.Phase == ConditionTypeScaling && phase != ConditionTypeScaling:
		// we're transitioning out of Scaling state
		in.Conditions.SetCondition(Condition{Type: ConditionTypeScaling, Status: corev1.ConditionFalse})
	case in.Phase == ConditionTypeRollingUpgrade && phase != ConditionTypeRollingUpgrade:
		// we're transitioning out of Upgrading state
		in.Conditions.SetCondition(Condition{Type: ConditionTypeRollingUpgrade, Status: corev1.ConditionFalse})
	case in.Phase == ConditionTypeWaiting && phase != ConditionTypeWaiting:
		// we're transitioning out of Waiting state
		in.Conditions.SetCondition(Condition{Type: ConditionTypeWaiting, Status: corev1.ConditionFalse})
	case in.Phase == ConditionTypeStopped && phase != ConditionTypeStopped:
		// we're transitioning out of Stopped state
		in.Conditions.SetCondition(Condition{Type: ConditionTypeStopped, Status: corev1.ConditionFalse})
	}

	// if we're complete we don't change the phase again
	if in.Phase != ConditionTypeCompleted {
		in.Phase = phase
	}

	in.Conditions.SetCondition(Condition{Type: phase, Status: corev1.ConditionTrue})
	return true
}

// ensure that the initial state conditions are present
func (in *CoherenceResourceStatus) ensureInitialized(deployment *Coherence) bool {
	updated := false

	// update Hash if required
	if in.Hash != deployment.Status.Hash {
		in.Hash = deployment.Status.Hash
		updated = true
	}

	// update Replicas if required
	if in.Replicas != deployment.Spec.GetReplicas() {
		in.Replicas = deployment.Spec.GetReplicas()
		updated = true
	}

	// update cluster name if required
	if in.CoherenceCluster != deployment.GetCoherenceClusterName() {
		in.CoherenceCluster = deployment.GetCoherenceClusterName()
		updated = true
	}

	// ensure that there is an Initialized condition
	if in.Conditions.GetCondition(ConditionTypeInitialized) == nil {
		// there is not an Initialized condition - this is probably the first status update
		updated = in.setPhase(ConditionTypeInitialized)
	}

	// update Selector if required
	if in.Selector == "" {
		in.Selector = fmt.Sprintf(StatusSelectorTemplate, deployment.GetCoherenceClusterName(), deployment.Name)
		updated = true
	}

	// update Role if required
	if in.Role != deployment.GetRoleName() {
		in.Role = deployment.GetRoleName()
		updated = true
	}

	// update the type if required
	t := deployment.GetType()
	if in.Type != t {
		in.Type = t
		updated = true
	}

	return updated
}
