/*
 * Copyright (c) 2020, 2025, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1

import (
	"context"
	"fmt"
	"github.com/distribution/reference"
	"github.com/go-test/deep"
	"github.com/oracle/coherence-operator/pkg/operator"
	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	apiequality "k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var (
	webhookLogger = logf.Log.WithName("coherence-webhook")
)

func (in *Coherence) SetupWebhookWithManager(mgr ctrl.Manager) error {
	hook := &Coherence{}
	return ctrl.NewWebhookManagedBy(mgr).
		For(in).
		WithDefaulter(hook).
		WithValidator(hook).
		Complete()
}

// The path in this annotation MUST match the const below
// +kubebuilder:webhook:path=/mutate-coherence-oracle-com-v1-coherence,mutating=true,failurePolicy=fail,groups=coherence.oracle.com,resources=coherence,verbs=create;update,versions=v1,name=mcoherence.kb.io

// MutatingWebHookPath This const MUST match the path in the kubebuilder annotation above
const MutatingWebHookPath = "/mutate-coherence-oracle-com-v1-coherence"

// An anonymous var to ensure that the Coherence struct implements webhook.Defaulter
// there will be a compile time error here if this fails.
var _ webhook.CustomDefaulter = &Coherence{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (in *Coherence) Default(_ context.Context, obj runtime.Object) error {
	coh, ok := obj.(*Coherence)
	if !ok {
		return fmt.Errorf("expected a Coherence instance but got a %T", obj)
	}

	spec, _ := coh.GetStatefulSetSpec()
	// set the default replicas if not present
	if spec.Replicas == nil {
		spec.SetReplicas(spec.GetReplicas())
	}
	SetCommonDefaults(coh)
	return nil
}

// SetCommonDefaults sets defaults common to both a Job and StatefulSet
func SetCommonDefaults(in CoherenceResource) {
	logger := webhookLogger.WithValues("namespace", in.GetNamespace(), "name", in.GetName())
	status := in.GetStatus()
	spec := in.GetSpec()
	dt := in.GetDeletionTimestamp()

	if dt != nil {
		// the deletion timestamp is set so do nothing
		logger.Info("Skipping updating defaults for deleted resource", "deletionTimestamp", *dt)
		return
	}

	if status.Phase == "" {
		logger.Info("Setting defaults for new resource")

		stsSpec, found := in.GetStatefulSetSpec()
		if found {
			// ensure the operator finalizer is present
			if stsSpec.AllowUnsafeDelete != nil && *stsSpec.AllowUnsafeDelete {
				if controllerutil.ContainsFinalizer(in, CoherenceFinalizer) {
					controllerutil.RemoveFinalizer(in, CoherenceFinalizer)
					logger.Info("Removed Finalizer from Coherence resource as AllowUnsafeDelete has been set to true")
				} else {
					logger.Info("Finalizer not added to Coherence resource as AllowUnsafeDelete has been set to true")
				}
			} else {
				controllerutil.AddFinalizer(in, CoherenceFinalizer)
			}
		}

		// set the default Coherence local port and local port adjust if not present
		if spec.Coherence == nil {
			var lpa = intstr.FromInt32(DefaultUnicastPortAdjust)
			spec.Coherence = &CoherenceSpec{
				LocalPort:       ptr.To(DefaultUnicastPort),
				LocalPortAdjust: &lpa,
			}
		} else {
			if spec.Coherence.LocalPort == nil {
				spec.Coherence.LocalPort = ptr.To(DefaultUnicastPort)
			}
			if spec.Coherence.LocalPortAdjust == nil {
				lpa := intstr.FromInt32(DefaultUnicastPortAdjust)
				spec.Coherence.LocalPortAdjust = &lpa
			}
		}

		// only set defaults for image names in new Coherence instances
		coherenceImage := operator.GetDefaultCoherenceImage()
		spec.EnsureCoherenceImage(&coherenceImage)

		// Set the features supported by this version
		in.AddAnnotation(AnnotationFeatureSuspend, "true")
	} else {
		// this is an update
		logger.Info("Updating defaults for existing resource")
	}
}

// The path in this annotation MUST match the const below
// +kubebuilder:webhook:verbs=create;update,path=/validate-coherence-oracle-com-v1-coherence,mutating=false,failurePolicy=fail,groups=coherence.oracle.com,resources=coherence,versions=v1,name=vcoherence.kb.io

// ValidatingWebHookPath This const MUST match the path in the kubebuilder annotation above
const ValidatingWebHookPath = "/validate-coherence-oracle-com-v1-coherence"

// An anonymous var to ensure that the Coherence struct implements webhook.Validator
// there will be a compile time error here if this fails.
var _ webhook.CustomValidator = &Coherence{}
var commonWebHook = CommonWebHook{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
// The optional warnings will be added to the response as warning messages.
// Return an error if the object is invalid.
func (in *Coherence) ValidateCreate(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	coh, ok := obj.(*Coherence)
	if !ok {
		return nil, fmt.Errorf("expected a Coherence instance but got a %T", obj)
	}

	logger := webhookLogger.WithValues("namespace", coh.GetNamespace(), "name", coh.GetName())
	var warnings admission.Warnings

	dt := coh.GetDeletionTimestamp()
	if dt != nil {
		// the deletion timestamp is set so do nothing
		logger.Info("Skipping validation for deleted resource", "deletionTimestamp", *dt)
		return warnings, nil
	}

	webhookLogger.Info("validate create", "name", coh.Name)
	if err := commonWebHook.validateReplicas(coh); err != nil {
		return warnings, err
	}
	if err := commonWebHook.validateImages(coh); err != nil {
		return warnings, err
	}
	if err := commonWebHook.validateNodePorts(coh); err != nil {
		return warnings, err
	}
	return warnings, nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
// The optional warnings will be added to the response as warning messages.
// Return an error if the object is invalid.
func (in *Coherence) ValidateUpdate(_ context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	cohNew, ok := newObj.(*Coherence)
	if !ok {
		return nil, fmt.Errorf("expected a Coherence instance for new value but got a %T", newObj)
	}
	cohPrev, ok := oldObj.(*Coherence)
	if !ok {
		return nil, fmt.Errorf("expected a Coherence instance for old value but got a %T", newObj)
	}

	webhookLogger.Info("validate update", "name", cohNew.Name)
	logger := webhookLogger.WithValues("namespace", cohNew.GetNamespace(), "name", cohNew.GetName())
	var warnings admission.Warnings

	dt := cohNew.GetDeletionTimestamp()
	if dt != nil {
		// the deletion timestamp is set so do nothing
		logger.Info("Skipping validation for deleted resource", "deletionTimestamp", *dt)
		return warnings, nil
	}

	if err := commonWebHook.validateReplicas(cohNew); err != nil {
		return warnings, err
	}
	if err := commonWebHook.validateImages(cohNew); err != nil {
		return warnings, err
	}

	if err := commonWebHook.validatePersistence(cohNew, cohPrev); err != nil {
		return warnings, err
	}
	if err := cohNew.validateVolumeClaimTemplates(cohNew, cohPrev); err != nil {
		return warnings, err
	}
	if err := commonWebHook.validateNodePorts(cohNew); err != nil {
		return warnings, err
	}

	var errorList field.ErrorList
	sts := cohNew.Spec.CreateStatefulSet(cohNew)
	stsOld := cohPrev.Spec.CreateStatefulSet(cohPrev)
	errorList = ValidateStatefulSetUpdate(&sts, &stsOld)

	if len(errorList) > 0 {
		return warnings, fmt.Errorf("rejecting update as it would have resulted in an invalid StatefulSet: %v", errorList)
	}

	return warnings, nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
// The optional warnings will be added to the response as warning messages.
// Return an error if the object is invalid.
func (in *Coherence) ValidateDelete(context.Context, runtime.Object) (admission.Warnings, error) {
	// we do not need to validate deletions
	return nil, nil
}

func (in *Coherence) validateVolumeClaimTemplates(cohNew, cohPrev *Coherence) error {
	if cohNew.GetReplicas() == 0 || cohPrev.GetReplicas() == 0 {
		// changes are allowed if current or previous replicas == 0
		return nil
	}

	if len(cohNew.Spec.VolumeClaimTemplates) == 0 && len(cohPrev.Spec.VolumeClaimTemplates) == 0 {
		// no PVCs in either deployment
		return nil
	}

	diff := deep.Equal(cohPrev.Spec.VolumeClaimTemplates, cohNew.Spec.VolumeClaimTemplates)
	if len(diff) != 0 {
		return fmt.Errorf("the Coherence resource \"%s\" is invalid: "+
			"changes cannot be made to spec.volumeclaimtemplates unless spec.replicas == 0 or the previous"+
			" instance of the resource has spec.replicas == 0", cohNew.Name)
	}
	return nil
}

// ----- Common Validator ---------------------------------------------------

type CommonWebHook struct {
}

// validateImages validates image names
func (in *CommonWebHook) validateImages(c CoherenceResource) error {
	var err error
	spec := c.GetSpec()
	if spec != nil {
		img := spec.GetCoherenceImage()
		if img != nil {
			_, err = reference.Parse(*img)
			if err != nil {
				return errors.Errorf("invalid spec.image field, %s", err.Error())
			}
		}
		for _, c := range spec.InitContainers {
			_, err = reference.Parse(c.Image)
			if err != nil {
				return errors.Errorf("invalid image name in init-container %s, %s", c.Name, err.Error())
			}
		}
		for _, c := range spec.SideCars {
			_, err = reference.Parse(c.Image)
			if err != nil {
				return errors.Errorf("invalid image name in side-car container %s, %s", c.Name, err.Error())
			}
		}
	}
	return err
}

// validateReplicas validates that spec.replicas >= 0
func (in *CommonWebHook) validateReplicas(c CoherenceResource) error {
	replicas := c.GetReplicas()
	if replicas < 0 {
		return fmt.Errorf("the Coherence resource \"%s\" is invalid: spec.replicas: Invalid value: %d: "+
			"must be greater than or equal to 0", c.GetName(), replicas)
	}
	return nil
}

func (in *CommonWebHook) validatePersistence(current, previous CoherenceResource) error {
	if current.GetReplicas() == 0 || previous.GetReplicas() == 0 {
		// changes are allowed if current or previous replicas == 0
		return nil
	}

	currentSpec := current.GetSpec()
	previousSpec := previous.GetSpec()
	diff := deep.Equal(previousSpec.GetCoherencePersistence(), currentSpec.GetCoherencePersistence())
	if len(diff) != 0 {
		return fmt.Errorf("the Coherence resource \"%s\" is invalid: "+
			"changes cannot be made to spec.coherence.persistence unless spec.replicas == 0 or the previous"+
			" instance of the resource has spec.replicas == 0", current.GetName())
	}
	return nil
}

func (in *CommonWebHook) validateNodePorts(current CoherenceResource) error {
	var badPorts []string
	spec := current.GetSpec()
	for _, port := range spec.Ports {
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

// ValidateStatefulSetUpdate tests if required fields in the StatefulSet are set.
func ValidateStatefulSetUpdate(statefulSet, oldStatefulSet *appsv1.StatefulSet) field.ErrorList {
	var allErrs field.ErrorList

	// statefulset updates aren't super common and general updates are likely to be touching spec, so we'll do this
	// deep copy right away.  This avoids mutating our inputs
	newStatefulSetClone := statefulSet.DeepCopy()
	newStatefulSetClone.Spec.Replicas = oldStatefulSet.Spec.Replicas               // +k8s:verify-mutation:reason=clone
	newStatefulSetClone.Spec.Template = oldStatefulSet.Spec.Template               // +k8s:verify-mutation:reason=clone
	newStatefulSetClone.Spec.UpdateStrategy = oldStatefulSet.Spec.UpdateStrategy   // +k8s:verify-mutation:reason=clone
	newStatefulSetClone.Spec.MinReadySeconds = oldStatefulSet.Spec.MinReadySeconds // +k8s:verify-mutation:reason=clone

	newStatefulSetClone.Spec.PersistentVolumeClaimRetentionPolicy = oldStatefulSet.Spec.PersistentVolumeClaimRetentionPolicy // +k8s:verify-mutation:reason=clone
	if !apiequality.Semantic.DeepEqual(newStatefulSetClone.Spec, oldStatefulSet.Spec) {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec"), "updates to statefulset spec for fields other than 'replicas', 'template', 'updateStrategy', 'persistentVolumeClaimRetentionPolicy' and 'minReadySeconds' are forbidden"))
	}

	return allErrs
}
