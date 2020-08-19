/*
 * Copyright (c) 2020, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1

import (
	"github.com/oracle/coherence-operator/pkg/operator"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var (
	webhookLogger = logf.Log.WithName("coherence-webhook")
	coherenceImage string
	utilsImage string
)

func (r *Coherence) SetupWebhookWithManager(mgr ctrl.Manager) error {
	coherenceImage = operator.GetDefaultCoherenceImage()
	utilsImage = operator.GetDefaultUtilsImage()

	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
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
func (r *Coherence) Default() {
	webhookLogger.Info("default", "name", r.Name)

	if r.Spec.Replicas == nil {
		r.Spec.SetReplicas(3)
	}

	r.Spec.EnsureCoherenceImage(&coherenceImage)
	r.Spec.EnsureCoherenceUtilsImage(&utilsImage)
}

// The path in this annotation MUST match the const below
// +kubebuilder:webhook:verbs=create;update,path=/validate-coherence-oracle-com-v1-coherence,mutating=false,failurePolicy=fail,groups=coherence.oracle.com,resources=coherences,versions=v1,name=vcoherence.kb.io

// This const MUST match the path in the kubebuilder annotation above
const ValidatingWebHookPath = "/validate-coherence-oracle-com-v1-coherence"

// An anonymous var to ensure that the Coherence struct implements webhook.Validator
// there will be a compile time error here if this fails.
var _ webhook.Validator = &Coherence{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Coherence) ValidateCreate() error {
	webhookLogger.Info("validate create", "name", r.Name)
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Coherence) ValidateUpdate(previous runtime.Object) error {
	webhookLogger.Info("validate update", "name", r.Name)
	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Coherence) ValidateDelete() error {
	// we do not need to validate deletions
	return nil
}
