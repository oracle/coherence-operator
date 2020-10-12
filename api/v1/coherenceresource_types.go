/*
 * Copyright (c) 2020, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1

import (
	"context"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/go-logr/logr"
	"github.com/operator-framework/operator-lib/status"
	"github.com/oracle/coherence-operator/pkg/clients"
	"github.com/oracle/coherence-operator/pkg/data"
	"github.com/oracle/coherence-operator/pkg/rest"
	"github.com/pkg/errors"
	"io/ioutil"
	appsv1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	crdv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	crdbeta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	v1client "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1"
	v1beta1client "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1beta1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/version"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// The package init function that will automatically register the Coherence resource types with
// the default k8s Scheme.
func init() {
	SchemeBuilder.Register(&Coherence{}, &CoherenceList{})
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
	Replicas int32 `json:"replicas,omitempty"`
	// CurrentReplicas is the current number of members in the Coherence deployment
	// represented by the Coherence resource.
	CurrentReplicas int32 `json:"currentReplicas,omitempty"`
	// ReadyReplicas is the number of number of members in the Coherence deployment
	// represented by the Coherence resource that are in the ready state.
	// +optional
	ReadyReplicas int32 `json:"readyReplicas,omitempty"`
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
}

// Update the current Phase
func (in *CoherenceResourceStatus) UpdatePhase(deployment *Coherence, phase status.ConditionType) bool {
	return in.SetCondition(deployment, status.Condition{Type: phase, Status: coreV1.ConditionTrue})
}

// Set the current Status Condition
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

// ----- Coherence type ------------------------------------------------------------------

// Coherence is the Schema for the Coherence API.
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

// Obtain the name of the headless Service used for Coherence WKA.
func (in *Coherence) GetWkaServiceName() string {
	if in == nil {
		return ""
	}
	return in.Name + WKAServiceNameSuffix
}

// Obtain the name of the headless Service used for the StatefulSet.
func (in *Coherence) GetHeadlessServiceName() string {
	if in == nil {
		return ""
	}
	return in.Name + HeadlessServiceNameSuffix
}

// Obtain the number of replicas required for a deployment.
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

// Set the number of replicas required for a deployment.
func (in *Coherence) SetReplicas(replicas int32) {
	if in != nil {
		in.Spec.Replicas = &replicas
	}
}

// Create the deployment's common label set.
func (in *Coherence) CreateCommonLabels() map[string]string {
	labels := make(map[string]string)
	labels[LabelCoherenceDeployment] = in.Name
	labels[LabelCoherenceCluster] = in.GetCoherenceClusterName()
	labels[LabelCoherenceRole] = in.GetRoleName()
	return labels
}

func (in *Coherence) GetNamespacedName() types.NamespacedName {
	return types.NamespacedName{
		Namespace: in.Namespace,
		Name:      in.Name,
	}
}

// Obtain the role name for a deployment.
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

// ----- CoherenceList type ------------------------------------------------------------------------

// +kubebuilder:object:root=true

// CoherenceList contains a list of Coherence resources.
type CoherenceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Coherence `json:"items"`
}

// EnsureCRDs ensures that the Operator configuration secret exists in the namespace.
// CRDs will be created depending on the server version of k8s. For k8s v1.16.0 and above
// the v1 CRDs will be created and for lower than v1.16.0 the v1beta1 CRDs will be created.
func EnsureCRDs(c clients.ClientSet) error {
	sv, err := c.DiscoveryClient.ServerVersion()
	if err != nil {
		return err
	}
	v, err := version.ParseSemantic(sv.GitVersion)
	if err != nil {
		return err
	}

	logger := logf.Log.WithName("operator")

	if v.Major() > 1 || (v.Major() == 1 && v.Minor() >= 16) {
		// k8s v1.16.0 or above - install v1 CRD
		crdClient := c.ExtClient.ApiextensionsV1().CustomResourceDefinitions()
		return EnsureV1CRDs(logger, crdClient)
	}
	// k8s lower than v1.16.0 - install v1beta1 CRD
	crdClient := c.ExtClient.ApiextensionsV1beta1().CustomResourceDefinitions()
	return EnsureV1Beta1CRDs(logger, crdClient)
}

// EnsureCRDs ensures that the Operator configuration secret exists in the namespace.
func EnsureV1CRDs(logger logr.Logger, crdClient v1client.CustomResourceDefinitionInterface) error {
	return ensureV1CRDs(logger, crdClient, "crd_v1.yaml")
}

// EnsureCRD ensures that the specified V1 CRDs are loaded using the specified embedded CRD files
func ensureV1CRDs(logger logr.Logger, crdClient v1client.CustomResourceDefinitionInterface, fileNames ...string) error {
	logger.Info("Ensuring operator v1 CRDs are present")
	for _, fileName := range fileNames {
		if err := ensureV1CRD(logger, crdClient, fileName); err != nil {
			return err
		}
	}
	return nil
}

// EnsureCRD ensures that the specified V1 CRD is loaded using the specified embedded CRD file
func ensureV1CRD(logger logr.Logger, crdClient v1client.CustomResourceDefinitionInterface, fileName string) error {
	f, err := data.Assets.Open(fileName)
	if err != nil {
		return errors.Wrap(err, "opening embedded CRD asset "+fileName)
	}
	defer f.Close()

	yml, err := ioutil.ReadAll(f)
	if err != nil {
		return errors.Wrap(err, "reading embedded CRD asset "+fileName)
	}

	u := unstructured.Unstructured{}
	err = yaml.Unmarshal(yml, &u)
	if err != nil {
		return err
	}

	newCRD := crdv1.CustomResourceDefinition{}
	err = yaml.Unmarshal(yml, &newCRD)
	if err != nil {
		return err
	}

	logger.Info("Loading operator CRD yaml from '" + fileName + "'")

	// Get the existing CRD
	oldCRD, err := crdClient.Get(context.TODO(), newCRD.Name, metav1.GetOptions{})
	switch {
	case err == nil:
		// CRD exists so update it
		logger.Info("Updating operator CRD '" + newCRD.Name + "'")
		newCRD.ResourceVersion = oldCRD.ResourceVersion
		_, err = crdClient.Update(context.TODO(), &newCRD, metav1.UpdateOptions{})
		if err != nil {
			return errors.Wrapf(err, "updating Coherence CRD %s", newCRD.Name)
		}
	case apierrors.IsNotFound(err):
		// CRD does not exist so create it
		logger.Info("Creating operator CRD '" + newCRD.Name + "'")
		_, err = crdClient.Create(context.TODO(), &newCRD, metav1.CreateOptions{})
		if err != nil {
			return errors.Wrapf(err, "creating Coherence CRD %s", newCRD.Name)
		}
	default:
		// An error occurred
		logger.Error(err, "checking for existing Coherence CRD "+newCRD.Name)
		return errors.Wrapf(err, "checking for existing Coherence CRD %s", newCRD.Name)
	}

	return nil
}

// EnsureCRDs ensures that the Operator configuration secret exists in the namespace.
func EnsureV1Beta1CRDs(logger logr.Logger, crdClient v1beta1client.CustomResourceDefinitionInterface) error {
	return ensureV1Beta1CRDs(logger, crdClient, "crd_v1beta1.yaml")
}

// EnsureCRD ensures that the specified V1 CRDs are loaded using the specified embedded CRD files
func ensureV1Beta1CRDs(logger logr.Logger, crdClient v1beta1client.CustomResourceDefinitionInterface, fileNames ...string) error {
	logger.Info("Ensuring operator v1beta1 CRDs are present")
	for _, fileName := range fileNames {
		if err := ensureV1Beta1CRD(logger, crdClient, fileName); err != nil {
			return err
		}
	}
	return nil
}

// EnsureCRD ensures that the specified V1 CRD is loaded using the specified embedded CRD file
func ensureV1Beta1CRD(logger logr.Logger, crdClient v1beta1client.CustomResourceDefinitionInterface, fileName string) error {
	f, err := data.Assets.Open(fileName)
	if err != nil {
		return errors.Wrap(err, "opening embedded CRD asset "+fileName)
	}
	defer f.Close()

	yml, err := ioutil.ReadAll(f)
	if err != nil {
		return errors.Wrap(err, "reading embedded CRD asset "+fileName)
	}

	u := unstructured.Unstructured{}
	err = yaml.Unmarshal(yml, &u)
	if err != nil {
		return err
	}

	newCRD := crdbeta1.CustomResourceDefinition{}
	err = yaml.Unmarshal(yml, &newCRD)
	if err != nil {
		return err
	}

	logger.Info("Loading operator CRD yaml from '" + fileName + "'")

	// Get the existing CRD
	oldCRD, err := crdClient.Get(context.TODO(), newCRD.Name, metav1.GetOptions{})
	switch {
	case err == nil:
		// CRD exists so update it
		logger.Info("Updating operator CRD '" + newCRD.Name + "'")
		newCRD.ResourceVersion = oldCRD.ResourceVersion
		_, err = crdClient.Update(context.TODO(), &newCRD, metav1.UpdateOptions{})
		if err != nil {
			return errors.Wrapf(err, "updating Coherence CRD %s", newCRD.Name)
		}
	case apierrors.IsNotFound(err):
		// CRD does not exist so create it
		logger.Info("Creating operator CRD '" + newCRD.Name + "'")
		_, err = crdClient.Create(context.TODO(), &newCRD, metav1.CreateOptions{})
		if err != nil {
			return errors.Wrapf(err, "creating Coherence CRD %s", newCRD.Name)
		}
	default:
		// An error occurred
		logger.Error(err, "checking for existing Coherence CRD "+newCRD.Name)
		return errors.Wrapf(err, "checking for existing Coherence CRD %s", newCRD.Name)
	}

	return nil
}

// EnsureOperatorSecret ensures that the Operator configuration secret exists in the namespace.
func EnsureOperatorSecret(namespace string, c client.Client, log logr.Logger) error {
	log.Info("Ensuring configuration secret")

	secret := &coreV1.Secret{}

	err := c.Get(context.TODO(), types.NamespacedName{Name: OperatorConfigName, Namespace: namespace}, secret)
	if err != nil && !apierrors.IsNotFound(err) {
		return err
	}

	restHostAndPort := rest.GetServerHostAndPort()

	log.Info(fmt.Sprintf("Operator Configuration: '%s' value set to %s", OperatorConfigKeyHost, restHostAndPort))

	secret.SetNamespace(namespace)
	secret.SetName(OperatorConfigName)

	if secret.StringData == nil {
		secret.StringData = make(map[string]string)
	}

	secret.StringData[OperatorConfigKeyHost] = restHostAndPort

	if apierrors.IsNotFound(err) {
		// for some reason we're getting here even if the secret exists so delete it!!
		_ = c.Delete(context.TODO(), secret)
		log.Info("Creating secret " + OperatorConfigName + " in namespace " + namespace)
		err = c.Create(context.TODO(), secret)
	} else {
		log.Info("Updating secret " + OperatorConfigName + " in namespace " + namespace)
		err = c.Update(context.TODO(), secret)
	}

	return err
}
