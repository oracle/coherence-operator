/*
 * Copyright (c) 2020, 2024, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1

import (
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	apiequality "k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// ValidateStatefulSetUpdate tests if required fields in the StatefulSet are set.
func ValidateStatefulSetUpdate(statefulSet, oldStatefulSet *appsv1.StatefulSet) field.ErrorList {
	var allErrs field.ErrorList

	// StatefulSet updates aren't super common and general updates are likely to be touching spec, so we'll do this
	// deep copy right away.  This avoids mutating our inputs
	newStatefulSetClone := statefulSet.DeepCopy()
	newStatefulSetClone.Spec.Replicas = oldStatefulSet.Spec.Replicas               // +k8s:verify-mutation:reason=clone
	newStatefulSetClone.Spec.Template = oldStatefulSet.Spec.Template               // +k8s:verify-mutation:reason=clone
	newStatefulSetClone.Spec.UpdateStrategy = oldStatefulSet.Spec.UpdateStrategy   // +k8s:verify-mutation:reason=clone
	newStatefulSetClone.Spec.MinReadySeconds = oldStatefulSet.Spec.MinReadySeconds // +k8s:verify-mutation:reason=clone

	newStatefulSetClone.Spec.PersistentVolumeClaimRetentionPolicy = oldStatefulSet.Spec.PersistentVolumeClaimRetentionPolicy // +k8s:verify-mutation:reason=clone
	if !apiequality.Semantic.DeepEqual(newStatefulSetClone.Spec, oldStatefulSet.Spec) {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec"), "updates to StatefulSet spec for fields other than 'replicas', 'template', 'updateStrategy', 'persistentVolumeClaimRetentionPolicy' and 'minReadySeconds' are forbidden"))
	}

	return allErrs
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
