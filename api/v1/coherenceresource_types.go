/*
 * Copyright (c) 2020, 2021, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1

import (
	"fmt"
	"github.com/operator-framework/operator-lib/status"
	"golang.org/x/mod/semver"
	appsv1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// Coherence resource Condition Types
// The different eight types of state that a deployment may be in.
//
// Transitions are:
// Initialized    -> Waiting
//                -> Created
// Waiting        -> Created
// Created        -> Ready
//                -> Stopped
// Ready          -> Scaling
//                -> RollingUpgrade
//                -> Stopped
// Scaling        -> Ready
//                -> Stopped
// RollingUpgrade -> Ready
// Stopped        -> Created
const (
	ConditionTypeInitialized    status.ConditionType = "Initialized"
	ConditionTypeWaiting        status.ConditionType = "Waiting"
	ConditionTypeCreated        status.ConditionType = "Created"
	ConditionTypeReady          status.ConditionType = "Ready"
	ConditionTypeScaling        status.ConditionType = "Scaling"
	ConditionTypeRollingUpgrade status.ConditionType = "RollingUpgrade"
	ConditionTypeFailed         status.ConditionType = "Failed"
	ConditionTypeStopped        status.ConditionType = "Stopped"
)

// The package init function that will automatically register the Coherence resource types with
// the default k8s Scheme.
func init() {
	SchemeBuilder.Register(&Coherence{}, &CoherenceList{})
}

// ----- Coherence type ------------------------------------------------------------------

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
type Coherence struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CoherenceResourceSpec   `json:"spec,omitempty"`
	Status CoherenceResourceStatus `json:"status,omitempty"`
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
	if in == nil {
		return 0
	}
	if in.Spec.Replicas == nil {
		return DefaultReplicas
	}
	return *in.Spec.Replicas
}

// SetReplicas sets the number of replicas required for a deployment.
func (in *Coherence) SetReplicas(replicas int32) {
	if in != nil {
		in.Spec.Replicas = &replicas
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
		m[k] = v + "." + in.GetNamespace() + ".svc.cluster.local"
	}
	return m
}

// FindFullyQualifiedPortServiceName returns the fully qualified name of the Service used to expose a named port and a bool indicating
// whether the named port has a Service.
func (in *Coherence) FindFullyQualifiedPortServiceName(name string) (string, bool) {
	n, found := in.FindPortServiceName(name)
	if found {
		n = n + "." + in.GetNamespace() + ".svc.cluster.local"
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

// GetNamespacedName returns the namespace/name key to lookup this resource.
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
	// Created:        The deployments secondary resources, (e.g. the StatefulSet, Services etc) have been created.
	// Ready:          The StatefulSet for the deployment has the correct number of replicas and ready replicas.
	// Waiting:        The deployment's start quorum conditions have not yet been met.
	// Scaling:        The number of replicas in the deployment is being scaled up or down.
	// RollingUpgrade: The StatefulSet is performing a rolling upgrade.
	// Stopped:        The replica count has been set to zero.
	// Failed:         An error occurred reconciling the deployment and its secondary resources.
	//
	// +optional
	Phase status.ConditionType `json:"phase,omitempty"`
	// The name of the Coherence cluster that this deployment is part of.
	// +optional
	CoherenceCluster string `json:"coherenceCluster,omitempty"`
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
	// The effective role name for this deployment.
	// This will come from the Spec.Role field if set otherwise the deployment name
	// will be used for the role name
	// +optional
	Role string `json:"role,omitempty"`
	// label query over deployments that should match the replicas count. This is same
	// as the label selector but in the string format to avoid introspection
	// by clients. The string will be in the same format as the query-param syntax.
	// More info about label selectors: http://kubernetes.io/docs/user-guide/labels#label-selectors
	// +optional
	Selector string `json:"selector,omitempty"`
	// The status conditions.
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	Conditions status.Conditions `json:"conditions,omitempty"`
	// Hash is the hash of the latest applied Coherence spec
	// +optional
	Hash string `json:"hash,omitempty"`
	// ActionsExecuted tracks whether actions were executed
	// +optional
	ActionsExecuted bool `json:"actionsExecuted,omitempty"`
}

// UpdatePhase updates the current Phase
// TODO not used?
func (in *CoherenceResourceStatus) UpdatePhase(deployment *Coherence, phase status.ConditionType) bool {
	return in.SetCondition(deployment, status.Condition{Type: phase, Status: coreV1.ConditionTrue})
}

// SetCondition sets the current Status Condition
func (in *CoherenceResourceStatus) SetCondition(deployment *Coherence, c status.Condition) bool {
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

// set a status phase.
func (in *CoherenceResourceStatus) setPhase(phase status.ConditionType) bool {
	if in.Phase == phase {
		return false
	}

	switch {
	case in.Phase == ConditionTypeReady && phase != ConditionTypeReady:
		// we're transitioning out of Ready state
		in.Conditions.SetCondition(status.Condition{Type: ConditionTypeReady, Status: coreV1.ConditionFalse})
	case in.Phase == ConditionTypeScaling && phase != ConditionTypeScaling:
		// we're transitioning out of Scaling state
		in.Conditions.SetCondition(status.Condition{Type: ConditionTypeScaling, Status: coreV1.ConditionFalse})
	case in.Phase == ConditionTypeRollingUpgrade && phase != ConditionTypeRollingUpgrade:
		// we're transitioning out of Upgrading state
		in.Conditions.SetCondition(status.Condition{Type: ConditionTypeRollingUpgrade, Status: coreV1.ConditionFalse})
	case in.Phase == ConditionTypeWaiting && phase != ConditionTypeWaiting:
		// we're transitioning out of Waiting state
		in.Conditions.SetCondition(status.Condition{Type: ConditionTypeWaiting, Status: coreV1.ConditionFalse})
	case in.Phase == ConditionTypeStopped && phase != ConditionTypeStopped:
		// we're transitioning out of Stopped state
		in.Conditions.SetCondition(status.Condition{Type: ConditionTypeStopped, Status: coreV1.ConditionFalse})
	}
	in.Phase = phase
	in.Conditions.SetCondition(status.Condition{Type: phase, Status: coreV1.ConditionTrue})
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

	return updated
}
