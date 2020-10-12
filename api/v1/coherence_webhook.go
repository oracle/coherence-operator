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
// +kubebuilder:webhook:path=/mutate-coherence-oracle-com-v1-coherence,mutating=true,failurePolicy=fail,groups=coherence.oracle.com,resources=coherence,verbs=create;update,versions=v1,name=mcoherence.kb.io

// This const MUST match the path in the kubebuilder annotation above
const MutatingWebHookPath = "/mutate-coherence-oracle-com-v1-coherence"

// An anonymous var to ensure that the Coherence struct implements webhook.Defaulter
// there will be a compile time error here if this fails.
var _ webhook.Defaulter = &Coherence{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (in *Coherence) Default() {
	logger := webhookLogger.WithValues("namespace", in.Namespace, "name", in.Name)
	if in.Status.Phase == "" {
		logger.Info("setting defaults for resource")

		if in.Spec.Replicas == nil {
			in.Spec.SetReplicas(3)
		}

		// only set defaults for image names new Coherence instances
		coherenceImage := operator.GetDefaultCoherenceImage()
		in.Spec.EnsureCoherenceImage(&coherenceImage)
		utilsImage := operator.GetDefaultUtilsImage()
		in.Spec.EnsureCoherenceUtilsImage(&utilsImage)

		// Set the features supported by this version
		in.Annotations[ANNOTATION_FEATURE_SUSPEND] = "true"
	} else {
		logger.Info("skipping defaulting for existing resource")
	}
}

// The path in this annotation MUST match the const below
// +kubebuilder:webhook:verbs=create;update,path=/validate-coherence-oracle-com-v1-coherence,mutating=false,failurePolicy=fail,groups=coherence.oracle.com,resources=coherence,versions=v1,name=vcoherence.kb.io

// This const MUST match the path in the kubebuilder annotation above
const ValidatingWebHookPath = "/validate-coherence-oracle-com-v1-coherence"

// An anonymous var to ensure that the Coherence struct implements webhook.Validator
// there will be a compile time error here if this fails.
var _ webhook.Validator = &Coherence{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (in *Coherence) ValidateCreate() error {
	webhookLogger.Info("validate create", "name", in.Name)
	if err := in.validateReplicas(); err != nil {
		return err
	}
	if err := in.validateNodePorts(); err != nil {
		return err
	}
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (in *Coherence) ValidateUpdate(previous runtime.Object) error {
	webhookLogger.Info("validate update", "name", in.Name)
	if err := in.validateReplicas(); err != nil {
		return err
	}
	prev := previous.(*Coherence)
	if err := in.validatePersistence(prev); err != nil {
		return err
	}
	if err := in.validateVolumeClaimTemplates(prev); err != nil {
		return err
	}
	if err := in.validateNodePorts(); err != nil {
		return err
	}
	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (in *Coherence) ValidateDelete() error {
	// we do not need to validate deletions
	return nil
}

// validateReplicas validates that spec.replicas >= 0
func (in *Coherence) validateReplicas() error {
	replicas := in.GetReplicas()
	if replicas < 0 {
		return fmt.Errorf("the Coherence resource \"%s\" is invalid: spec.replicas: Invalid value: %d: "+
			"must be greater than or equal to 0", in.Name, replicas)
	}
	return nil
}

func (in *Coherence) validatePersistence(previous *Coherence) error {
	if in.GetReplicas() == 0 || previous.GetReplicas() == 0 {
		// changes are allowed if current or previous replicas == 0
		return nil
	}

	diff := deep.Equal(previous.Spec.GetCoherencePersistence(), in.Spec.GetCoherencePersistence())
	if len(diff) != 0 {
		return fmt.Errorf("the Coherence resource \"%s\" is invalid: "+
			"changes cannot be made to spec.coherence.persistence unless spec.replicas == 0 or the previous"+
			" instance of the resource has spec.replicas == 0", in.Name)
	}
	return nil
}

func (in *Coherence) validateVolumeClaimTemplates(previous *Coherence) error {
	if in.GetReplicas() == 0 || previous.GetReplicas() == 0 {
		// changes are allowed if current or previous replicas == 0
		return nil
	}

	if len(in.Spec.VolumeClaimTemplates) == 0 && len(previous.Spec.VolumeClaimTemplates) == 0 {
		// no PVCs in either deployment
		return nil
	}

	diff := deep.Equal(previous.Spec.VolumeClaimTemplates, in.Spec.VolumeClaimTemplates)
	if len(diff) != 0 {
		return fmt.Errorf("the Coherence resource \"%s\" is invalid: "+
			"changes cannot be made to spec.volumeclaimtemplates unless spec.replicas == 0 or the previous"+
			" instance of the resource has spec.replicas == 0", in.Name)
	}
	return nil
}

func (in *Coherence) validateNodePorts() error {
	var badPorts []string

	for _, port := range in.Spec.Ports {
		if port.NodePort != nil {
			p := *port.NodePort
			if p < 30000 || p > 32767 {
				badPorts = append(badPorts, port.Name)
			}
		}
	}

	if len(badPorts) > 0 {
		return fmt.Errorf("the following NodePort values are invalid, valid port range is 30000-32767 - %v", badPorts)
	}

	return nil
}
