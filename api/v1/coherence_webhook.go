/*
 * Copyright (c) 2020, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1

import (
	"fmt"
	"github.com/go-test/deep"
	"github.com/oracle/coherence-operator/pkg/operator"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var (
	webhookLogger = logf.Log.WithName("coherence-webhook")
)

func (in *Coherence) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(in).
		Complete()
}


// The path in this annotation MUST match the const below
// +kubebuilder:webhook:path=/mutate-coherence-oracle-com-v1-coherence,mutating=true,failurePolicy=fail,groups=coherence.oracle.com,resources=coherences,verbs=create;update,versions=v1,name=mcoherence.kb.io

// This const MUST match the path in the kubebuilder annotation above
const MutatingWebHookPath = "/mutate-coherence-oracle-com-v1-coherence"

// An anonymous var to ensure that the Coherence struct implements webhook.Defaulter
// there will be a compile time error here if this fails.
var _ webhook.Defaulter = &Coherence{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (in *Coherence) Default() {
	webhookLogger.Info("default", "name", in.Name)

	if in.Spec.Replicas == nil {
		in.Spec.SetReplicas(3)
	}

	coherenceImage := operator.GetDefaultCoherenceImage()
	in.Spec.EnsureCoherenceImage(&coherenceImage)
	utilsImage := operator.GetDefaultUtilsImage()
	in.Spec.EnsureCoherenceUtilsImage(&utilsImage)
}

// The path in this annotation MUST match the const below
// +kubebuilder:webhook:verbs=create;update,path=/validate-coherence-oracle-com-v1-coherence,mutating=false,failurePolicy=fail,groups=coherence.oracle.com,resources=coherences,versions=v1,name=vcoherence.kb.io

// This const MUST match the path in the kubebuilder annotation above
const ValidatingWebHookPath = "/validate-coherence-oracle-com-v1-coherence"

// An anonymous var to ensure that the Coherence struct implements webhook.Validator
// there will be a compile time error here if this fails.
var _ webhook.Validator = &Coherence{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (in *Coherence) ValidateCreate() error {
	webhookLogger.Info("validate create", "name", in.Name)
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (in *Coherence) ValidateUpdate(previous runtime.Object) error {
	webhookLogger.Info("validate update", "name", in.Name)
	prev := previous.(*Coherence)
	if err := in.validatePersistence(prev); err != nil {
		return err
	}
	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (in *Coherence) ValidateDelete() error {
	// we do not need to validate deletions
	return nil
}

func (in *Coherence) validatePersistence(previous *Coherence) error {
	if in.GetReplicas() == 0 || previous.GetReplicas() == 0 {
		// changes are allowed if current or previous replicas == 0
		return nil
	}

	diff := deep.Equal(previous.Spec.GetCoherencePersistence(), in.Spec.GetCoherencePersistence())
	if len(diff) != 0 {
		return fmt.Errorf("changes cannot be made to spec.coherence.persistence unless replicas == 0")
	}
	return nil
}