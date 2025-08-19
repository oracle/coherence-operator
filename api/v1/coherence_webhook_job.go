/*
 * Copyright (c) 2020, 2025, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1

import (
	"context"
	"fmt"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	apiequality "k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

func (in *CoherenceJob) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(in).
		Complete()
}

// The path in this annotation MUST match the const below
// +kubebuilder:webhook:path=/mutate-coherence-oracle-com-v1-coherencejob,mutating=true,failurePolicy=fail,groups=coherence.oracle.com,resources=coherencejob,verbs=create;update,versions=v1,name=mcoherencejob.kb.io

// JobMutatingWebHookPath This const MUST match the path in the kubebuilder annotation above
const JobMutatingWebHookPath = "/mutate-coherence-oracle-com-v1-coherencejob"

// An anonymous var to ensure that the CoherenceJob struct implements webhook.Defaulter
// there will be a compile time error here if this fails.
var _ webhook.CustomDefaulter = &CoherenceJob{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (in *CoherenceJob) Default(_ context.Context, obj runtime.Object) error {
	coh, ok := obj.(*CoherenceJob)
	if !ok {
		return fmt.Errorf("expected a CoherenceJob instance but got a %T", obj)
	}

	spec, _ := coh.GetJobResourceSpec()
	coherenceSpec := spec.Coherence
	if spec.Coherence == nil {
		coherenceSpec = &CoherenceSpec{}
		spec.Coherence = coherenceSpec
	}

	// default to storage disabled to false
	if coherenceSpec.StorageEnabled == nil {
		coherenceSpec.StorageEnabled = ptr.To(false)
	}

	// default the restart policy to never
	if spec.RestartPolicy == nil {
		spec.RestartPolicy = spec.RestartPolicyPointer(corev1.RestartPolicyNever)
	}

	co := spec.Coherence
	if co != nil {
		if co.StorageEnabled == nil {
			co.StorageEnabled = ptr.To(false)
		}
	}

	// set the default replicas if not present
	if spec.Replicas == nil {
		spec.SetReplicas(spec.GetReplicas())
	}

	SetCommonDefaults(coh)
	return nil
}

// The path in this annotation MUST match the const below
// +kubebuilder:webhook:verbs=create;update,path=/validate-coherence-oracle-com-v1-coherencejob,mutating=false,failurePolicy=fail,groups=coherence.oracle.com,resources=coherencejob,versions=v1,name=vcoherencejob.kb.io

// JobValidatingWebHookPath This const MUST match the path in the kubebuilder annotation above
const JobValidatingWebHookPath = "/validate-coherence-oracle-com-v1-coherencejob"

// An anonymous var to ensure that the Coherence struct implements webhook.Validator
// there will be a compile time error here if this fails.
var _ webhook.CustomValidator = &CoherenceJob{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (in *CoherenceJob) ValidateCreate(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	coh, ok := obj.(*CoherenceJob)
	if !ok {
		return nil, fmt.Errorf("expected a CoherenceJob instance but got a %T", obj)
	}

	var err error
	var warnings admission.Warnings

	webhookLogger.Info("validate create", "name", coh.Name)
	if err = commonWebHook.validateReplicas(coh); err != nil {
		return warnings, err
	}
	if err = commonWebHook.validateImages(coh); err != nil {
		return warnings, err
	}
	if err = commonWebHook.validateNodePorts(coh); err != nil {
		return warnings, err
	}
	return warnings, nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (in *CoherenceJob) ValidateUpdate(_ context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	cohNew, ok := newObj.(*CoherenceJob)
	if !ok {
		return nil, fmt.Errorf("expected a CoherenceJob instance for new value but got a %T", newObj)
	}
	cohPrev, ok := oldObj.(*CoherenceJob)
	if !ok {
		return nil, fmt.Errorf("expected a CoherenceJob instance for old value but got a %T", newObj)
	}

	webhookLogger.Info("validate update", "name", cohNew.Name)
	var warnings admission.Warnings

	if err := commonWebHook.validateReplicas(cohNew); err != nil {
		return warnings, err
	}
	if err := commonWebHook.validateImages(cohNew); err != nil {
		return warnings, err
	}

	if err := commonWebHook.validatePersistence(cohNew, cohPrev); err != nil {
		return warnings, err
	}
	if err := commonWebHook.validateNodePorts(cohNew); err != nil {
		return warnings, err
	}

	var errorList field.ErrorList
	job := cohNew.Spec.CreateJob(cohNew)
	jobOld := cohPrev.Spec.CreateJob(cohPrev)
	errorList = ValidateJobUpdate(&job, &jobOld)

	if len(errorList) > 0 {
		return warnings, fmt.Errorf("rejecting update as it would have resulted in an invalid Job: %v", errorList)
	}

	return warnings, nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (in *CoherenceJob) ValidateDelete(context.Context, runtime.Object) (admission.Warnings, error) {
	// we do not need to validate deletions
	return nil, nil
}

// ValidateJobUpdate tests if required fields in the Job are set.
func ValidateJobUpdate(job, oldJob *batchv1.Job) field.ErrorList {
	var allErrs field.ErrorList

	newJobClone := job.DeepCopy()
	newJobClone.Spec.ActiveDeadlineSeconds = oldJob.Spec.ActiveDeadlineSeconds     // +k8s:verify-mutation:reason=clone
	newJobClone.Spec.BackoffLimit = oldJob.Spec.BackoffLimit                       // +k8s:verify-mutation:reason=clone
	newJobClone.Spec.CompletionMode = oldJob.Spec.CompletionMode                   // +k8s:verify-mutation:reason=clone
	newJobClone.Spec.Parallelism = oldJob.Spec.Parallelism                         // +k8s:verify-mutation:reason=clone
	newJobClone.Spec.Suspend = oldJob.Spec.Suspend                                 // +k8s:verify-mutation:reason=clone
	newJobClone.Spec.Template = oldJob.Spec.Template                               // +k8s:verify-mutation:reason=clone
	newJobClone.Spec.TTLSecondsAfterFinished = oldJob.Spec.TTLSecondsAfterFinished // +k8s:verify-mutation:reason=clone
	newJobClone.Spec.PodFailurePolicy = oldJob.Spec.PodFailurePolicy               // +k8s:verify-mutation:reason=clone
	newJobClone.Spec.Completions = oldJob.Spec.Completions                         // +k8s:verify-mutation:reason=clone

	if !apiequality.Semantic.DeepEqual(newJobClone.Spec, oldJob.Spec) {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec"), "updates to Job spec for fields 'selector', 'manualSelector', are forbidden"))
	}
	return allErrs
}
