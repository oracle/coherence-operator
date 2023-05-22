/*
 * Copyright (c) 2019, 2023, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1

import (
	"golang.org/x/mod/semver"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/pointer"
)

// ----- Coherence type ------------------------------------------------------------------

// CoherenceJob is the top level schema for the Coherence API and custom resource definition (CRD)
// for configuring Coherence Job workloads.
//
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:subresource:scale:specpath=.spec.replicas,statuspath=.status.replicas,selectorpath=.status.selector
// +kubebuilder:resource:path=coherencejob,scope=Namespaced,shortName=cohjob,categories=coherence
// +kubebuilder:printcolumn:name="Cluster",type="string",JSONPath=".status.coherenceCluster",description="The name of the Coherence cluster that this deployment belongs to"
// +kubebuilder:printcolumn:name="Role",type="string",JSONPath=".status.role",description="The role of this deployment in a Coherence cluster"
// +kubebuilder:printcolumn:name="Replicas",type="integer",JSONPath=".status.replicas",description="The number of Coherence deployments for this deployment"
// +kubebuilder:printcolumn:name="Ready",type="integer",JSONPath=".status.readyReplicas",description="The number of ready Coherence deployments for this deployment"
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase",description="The status of this deployment"
// +kubebuilder:printcolumn:name="Active",priority=1,type="integer",JSONPath=".status.active",description="When the Coherence resource is running a Job, the number of pending and running pods"
// +kubebuilder:printcolumn:name="Succeeded",priority=1,type="integer",JSONPath=".status.succeeded",description="When the Coherence resource is running a Job, the number of pods which reached phase Succeeded"
// +kubebuilder:printcolumn:name="Failed",priority=1,type="integer",JSONPath=".status.failed",description="When the Coherence resource is running a Job, the number of pods which reached phase Failed"
type CoherenceJob struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CoherenceJobResourceSpec `json:"spec,omitempty"`
	Status CoherenceResourceStatus  `json:"status,omitempty"`
}

var _ CoherenceResource = &CoherenceJob{}

// GetCoherenceClusterName obtains the Coherence cluster name for the Coherence resource.
func (in *CoherenceJob) GetCoherenceClusterName() string {
	if in == nil {
		return ""
	}

	if in.Spec.Cluster == "" {
		return in.Name
	}
	return in.Spec.Cluster
}

func (in *CoherenceJob) GetAPIVersion() string {
	return in.APIVersion
}

// GetSpec returns this resource's CoherenceResourceSpec
func (in *CoherenceJob) GetSpec() *CoherenceResourceSpec {
	return &in.Spec.CoherenceResourceSpec
}

// GetJobResourceSpec returns this resource's CoherenceJobResourceSpec
func (in *CoherenceJob) GetJobResourceSpec() (*CoherenceJobResourceSpec, bool) {
	return &in.Spec, true
}

// GetStatefulSetSpec always returns nil and false
func (in *CoherenceJob) GetStatefulSetSpec() (*CoherenceStatefulSetResourceSpec, bool) {
	return nil, false
}

func (in *CoherenceJob) AddAnnotation(key, value string) {
	if in != nil {
		if in.Annotations == nil {
			in.Annotations = make(map[string]string)
		}
		in.Annotations[key] = value
	}
}

// GetStatus returns this resource's CoherenceResourceSpec
func (in *CoherenceJob) GetStatus() *CoherenceResourceStatus {
	return &in.Status
}

// CreateKubernetesResources returns this resource's CoherenceResourceSpec
func (in *CoherenceJob) CreateKubernetesResources() (Resources, error) {
	res := in.Spec.CreateKubernetesResources(in)
	res = append(res, in.Spec.CreateJobResource(in))
	return Resources{Items: res}, nil
}

func (in *CoherenceJob) DeepCopyResource() CoherenceResource {
	return in.DeepCopy()
}

// GetWkaServiceName returns the name of the headless Service used for Coherence WKA.
func (in *CoherenceJob) GetWkaServiceName() string {
	if in == nil {
		return ""
	}
	return in.Name + WKAServiceNameSuffix
}

// GetHeadlessServiceName returns the name of the headless Service used for the StatefulSet.
func (in *CoherenceJob) GetHeadlessServiceName() string {
	if in == nil {
		return ""
	}
	return in.Name + HeadlessServiceNameSuffix
}

// GetReplicas returns the number of replicas required for a deployment.
// The Replicas field is a pointer and may be nil so this method will
// return either the actual Replicas value or the default (DefaultReplicas const)
// if the Replicas field is nil.
func (in *CoherenceJob) GetReplicas() int32 {
	if in == nil || in.Spec.Replicas == nil {
		return DefaultJobReplicas
	}
	return in.Spec.GetReplicas()
}

// SetReplicas sets the number of replicas required for a deployment.
func (in *CoherenceJob) SetReplicas(replicas int32) {
	if in != nil {
		in.Spec.Replicas = &replicas
	}
}

// FindFullyQualifiedPortServiceNames returns a map of the exposed ports of this resource mapped to their Service's
// fully qualified domain name.
func (in *CoherenceJob) FindFullyQualifiedPortServiceNames() map[string]string {
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
func (in *CoherenceJob) FindFullyQualifiedPortServiceName(name string) (string, bool) {
	n, found := in.FindPortServiceName(name)
	if found {
		n = n + "." + in.GetNamespace() + ".svc"
	}
	return n, found
}

// FindPortServiceNames returns a map of the port names for this resource mapped to their Service names.
func (in *CoherenceJob) FindPortServiceNames() map[string]string {
	if in == nil {
		return make(map[string]string)
	}
	return in.Spec.FindPortServiceNames(in)
}

// FindPortServiceName returns the name of the Service used to expose a named port and a bool indicating
// whether the named port has a Service.
func (in *CoherenceJob) FindPortServiceName(name string) (string, bool) {
	if in == nil {
		return "", false
	}
	return in.Spec.FindPortServiceName(name, in)
}

// CreateCommonLabels creates the deployment's common label set.
func (in *CoherenceJob) CreateCommonLabels() map[string]string {
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
func (in *CoherenceJob) CreateAnnotations() map[string]string {
	var annotations map[string]string
	if in.Spec.JobAnnotations != nil {
		annotations = make(map[string]string)
		for k, v := range in.Spec.JobAnnotations {
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

// GetNamespacedName returns the namespace/name key to look up this resource.
func (in *CoherenceJob) GetNamespacedName() types.NamespacedName {
	return types.NamespacedName{
		Namespace: in.Namespace,
		Name:      in.Name,
	}
}

// GetRoleName returns the role name for a deployment.
// If the Spec.Role field is set that is used for the role name
// otherwise the deployment name is used as the role name.
func (in *CoherenceJob) GetRoleName() string {
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
func (in *CoherenceJob) GetType() CoherenceType {
	return CoherenceTypeJob
}

// GetWKA returns the host name Coherence should for WKA.
func (in *CoherenceJob) GetWKA() string {
	if in == nil {
		return ""
	}
	return in.Spec.Coherence.GetWKA(in)
}

// GetVersionAnnotation if the returns the value of the Operator version annotation and true,
// if the version annotation is present. If the version annotation is not present this method
// returns empty string and false.
func (in *CoherenceJob) GetVersionAnnotation() (string, bool) {
	if in == nil || in.Annotations == nil {
		return "", false
	}
	version, found := in.Annotations[AnnotationOperatorVersion]
	return version, found
}

// IsBeforeVersion returns true if this Coherence resource Operator version annotation value is
// before the specified version, or is not set.
// The version parameter must be a valid SemVer value.
func (in *CoherenceJob) IsBeforeVersion(version string) bool {
	if actual, found := in.GetVersionAnnotation(); found {
		return semver.Compare(actual, version) < 0
	}
	return true
}

// ----- CoherenceJobList type ----------------------------------------------

// CoherenceJobResourceSpec defines the specification of a CoherenceJob resource.
// +k8s:openapi-gen=true
type CoherenceJobResourceSpec struct {
	CoherenceResourceSpec `json:",inline"`

	// The name of the Coherence cluster that this CoherenceJob resource belongs to.
	// A CoherenceJob will typically be part of an existing cluster, so this field is required.
	Cluster string `json:"cluster,omitempty"`

	// Specifies the desired number of successfully finished pods the
	// job should be run with.  Setting to nil means that the success of any
	// pod signals the success of all pods, and allows parallelism to have any positive
	// value.  Setting to 1 means that parallelism is limited to 1 and the success of that
	// pod signals the success of the job.
	// More info: https://kubernetes.io/docs/concepts/workloads/controllers/jobs-run-to-completion/
	// +optional
	Completions *int32 `json:"completions,omitempty"`

	// SyncCompletions is a flag to indicate that the Operator should always set the
	// Completions value to be the same as the Replicas value.
	// When a Job is then scaled, the Completions value will also be changed.
	// +optional
	SyncCompletionsToReplicas *bool `json:"syncCompletionsToReplicas,omitempty"`

	// Specifies the policy of handling failed pods. In particular, it allows to
	// specify the set of actions and conditions which need to be
	// satisfied to take the associated action.
	// If empty, the default behaviour applies - the counter of failed pods,
	// represented by the job's .status.failed field, is incremented, and it is
	// checked against the backoffLimit. This field cannot be used in combination
	// with restartPolicy=OnFailure.
	//
	// This field is alpha-level. To use this field, you must enable the
	// `JobPodFailurePolicy` feature gate (disabled by default).
	// +optional
	PodFailurePolicy *batchv1.PodFailurePolicy `json:"podFailurePolicy,omitempty"`

	// Specifies the number of retries before marking this job failed.
	// Defaults to 6
	// +optional
	BackoffLimit *int32 `json:"backoffLimit,omitempty"`

	// ttlSecondsAfterFinished limits the lifetime of a Job that has finished
	// execution (either Complete or Failed). If this field is set,
	// ttlSecondsAfterFinished after the Job finishes, it is eligible to be
	// automatically deleted. When the Job is being deleted, its lifecycle
	// guarantees (e.g. finalizers) will be honored. If this field is unset,
	// the Job won't be automatically deleted. If this field is set to zero,
	// the Job becomes eligible to be deleted immediately after it finishes.
	// +optional
	TTLSecondsAfterFinished *int32 `json:"ttlSecondsAfterFinished,omitempty"`

	// CompletionMode specifies how Pod completions are tracked. It can be
	// `NonIndexed` (default) or `Indexed`.
	//
	// `NonIndexed` means that the Job is considered complete when there have
	// been .spec.completions successfully completed Pods. Each Pod completion is
	// homologous to each other.
	//
	// `Indexed` means that the Pods of a
	// Job get an associated completion index from 0 to (.spec.completions - 1),
	// available in the annotation batch.kubernetes.io/job-completion-index.
	// The Job is considered complete when there is one successfully completed Pod
	// for each index.
	// When value is `Indexed`, .spec.completions must be specified and
	// `.spec.parallelism` must be less than or equal to 10^5.
	// In addition, The Pod name takes the form
	// `$(job-name)-$(index)-$(random-string)`,
	// the Pod hostname takes the form `$(job-name)-$(index)`.
	//
	// More completion modes can be added in the future.
	// If the Job controller observes a mode that it doesn't recognize, which
	// is possible during upgrades due to version skew, the controller
	// skips updates for the Job.
	// +optional
	CompletionMode *batchv1.CompletionMode `json:"completionMode,omitempty"`

	// Suspend specifies whether the Job controller should create Pods or not. If
	// a Job is created with suspend set to true, no Pods are created by the Job
	// controller. If a Job is suspended after creation (i.e. the flag goes from
	// false to true), the Job controller will delete all active Pods associated
	// with this Job. Users must design their workload to gracefully handle this.
	// Suspending a Job will reset the StartTime field of the Job, effectively
	// resetting the ActiveDeadlineSeconds timer too. Defaults to false.
	//
	// +optional
	Suspend *bool `json:"suspend,omitempty"`

	// JobAnnotations are free-form yaml that will be added to the Coherence workload's
	// `Job` as annotations.
	// Any annotations should be placed BELOW this "annotations:" key, for example:
	//
	// The default behaviour is to copy all annotations from the `Coherence` resource to the
	// `Job`, specifying any annotations in the `JobAnnotations` will override
	// this behaviour and only include the `JobAnnotations`.
	//
	// annotations:
	//   foo.io/one: "value1"
	//   foo.io/two: "value2"
	//
	// see: https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/[Kubernetes Annotations]
	// +optional
	JobAnnotations map[string]string `json:"jobAnnotations,omitempty"`

	// ReadyAction is a probe that will be executed when one or more Pods
	// reach the ready state. The probe will be executed on every Pod that
	// is ready. One the required number of ready Pods is reached the probe
	// will also be executed on every Pod that becomes ready after that time.
	// +optional
	ReadyAction *CoherenceJobProbe `json:"ReadyAction,omitempty"`
}

// GetRestartPolicy returns the name of the application image to use
func (in *CoherenceJobResourceSpec) GetRestartPolicy() *corev1.RestartPolicy {
	if in == nil || in.RestartPolicy == nil {
		return in.RestartPolicyPointer(corev1.RestartPolicyNever)
	}
	return in.RestartPolicy
}

// GetReplicas returns the number of replicas required for a deployment.
// The Replicas field is a pointer and may be nil so this method will
// return either the actual Replica value or the default (DefaultReplicas const)
// if the Replicas field is nil.
func (in *CoherenceJobResourceSpec) GetReplicas() int32 {
	if in == nil || in.CoherenceResourceSpec.Replicas == nil {
		return DefaultJobReplicas
	}
	return *in.CoherenceResourceSpec.Replicas
}

// UpdateJob updates a JobSpec from the fields in this spec
func (in *CoherenceJobResourceSpec) UpdateJob(spec *batchv1.JobSpec) {
	if in == nil {
		return
	}

	if in.IsSyncCompletions() {
		spec.Completions = pointer.Int32(in.GetReplicas())
	} else {
		spec.Completions = in.Completions
	}

	spec.PodFailurePolicy = in.PodFailurePolicy
	spec.BackoffLimit = in.BackoffLimit
	spec.TTLSecondsAfterFinished = in.TTLSecondsAfterFinished
	spec.CompletionMode = in.CompletionMode
	spec.Suspend = in.Suspend
}

// IsSyncCompletions returns true if Completions should always match Parallelism
func (in *CoherenceJobResourceSpec) IsSyncCompletions() bool {
	if in == nil {
		return false
	}
	return in.SyncCompletionsToReplicas != nil && *in.SyncCompletionsToReplicas
}

// CreateJobResource creates the deployment's Job resource.
func (in *CoherenceJobResourceSpec) CreateJobResource(deployment CoherenceResource) Resource {
	job := in.CreateJob(deployment)

	return Resource{
		Kind: ResourceTypeJob,
		Name: job.GetName(),
		Spec: &job,
	}
}

// CreateJob creates the deployment's Job.
func (in *CoherenceJobResourceSpec) CreateJob(deployment CoherenceResource) batchv1.Job {
	job := batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:   deployment.GetNamespace(),
			Name:        deployment.GetName(),
			Labels:      deployment.CreateCommonLabels(),
			Annotations: deployment.CreateAnnotations(),
		},
	}

	replicas := in.GetReplicas()
	podTemplate := in.CreatePodTemplateSpec(deployment)

	restartPolicy := in.GetRestartPolicy()
	if restartPolicy != nil {
		podTemplate.Spec.RestartPolicy = *restartPolicy
	}

	// Add the component label
	job.Labels[LabelComponent] = LabelComponentCoherenceStatefulSet

	job.Spec = batchv1.JobSpec{
		Parallelism: &replicas,
		Template:    podTemplate,
	}

	job.Spec.ActiveDeadlineSeconds = in.ActiveDeadlineSeconds

	in.UpdateJob(&job.Spec)

	return job
}

// ----- CoherenceResourceStatus type ---------------------------------------

type CoherenceJobStatus struct {
	CoherenceResourceStatus `json:",inline"`
	ProbeStatus             []CoherenceJobProbeStatus `json:"probeStatus,omitempty"`
}

// ----- CoherenceJobList type ----------------------------------------------

// +kubebuilder:object:root=true

// CoherenceJobList is a list of CoherenceJob resources.
type CoherenceJobList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CoherenceJob `json:"items"`
}

// ----- CoherenceJobProbe type ---------------------------------------------

type CoherenceJobProbe struct {
	Probe `json:",inline"`
	// The number of job Pods that should be ready before executing the Probe.
	// If not set the default will be the same as the job's Completions value.
	// The probe will be executed on all Pods
	// +optional
	ReadyCount *int32 `json:"readyCount,omitempty"`
}

// ----- CoherenceJobProbeStatus type ----------------------------------------

type CoherenceJobProbeStatus struct {
	Pod           string       `json:"pod,omitempty"`
	LastReadyTime *metav1.Time `json:"lastReadyTime,omitempty"`
	LastProbeTime *metav1.Time `json:"lastProbeTime,omitempty"`
	Success       *bool        `json:"success,omitempty"`
}
