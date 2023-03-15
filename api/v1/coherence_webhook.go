/*
 * Copyright (c) 2020, 2023, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1

import (
	"fmt"
	"github.com/go-test/deep"
	"github.com/oracle/coherence-operator/pkg/operator"
	appsv1 "k8s.io/api/apps/v1"
	apiequality "k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
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

// MutatingWebHookPath This const MUST match the path in the kubebuilder annotation above
const MutatingWebHookPath = "/mutate-coherence-oracle-com-v1-coherence"

// An anonymous var to ensure that the Coherence struct implements webhook.Defaulter
// there will be a compile time error here if this fails.
var _ webhook.Defaulter = &Coherence{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (in *Coherence) Default() {
	logger := webhookLogger.WithValues("namespace", in.Namespace, "name", in.Name)
	if in.Status.Phase == "" {
		logger.Info("Setting defaults for new resource")
		// ensure the operator finalizer is present
		if in.Spec.AllowUnsafeDelete != nil && *in.Spec.AllowUnsafeDelete {
			if controllerutil.ContainsFinalizer(in, CoherenceFinalizer) {
				controllerutil.RemoveFinalizer(in, CoherenceFinalizer)
				logger.Info("Removed Finalizer from Coherence resource as AllowUnsafeDelete has been set to true")
			} else {
				logger.Info("Finalizer not added to Coherence resource as AllowUnsafeDelete has been set to true")
			}
		} else {
			controllerutil.AddFinalizer(in, CoherenceFinalizer)
		}

		// set the default replicas if not present
		if in.Spec.Replicas == nil {
			in.Spec.SetReplicas(3)
		}

		// set the default Coherence local port and local port adjust if not present
		if in.Spec.Coherence == nil {
			lpa := intstr.FromInt(int(DefaultUnicastPortAdjust))
			in.Spec.Coherence = &CoherenceSpec{
				LocalPort:       pointer.Int32(DefaultUnicastPort),
				LocalPortAdjust: &lpa,
			}
		} else {
			if in.Spec.Coherence.LocalPort == nil {
				in.Spec.Coherence.LocalPort = pointer.Int32(DefaultUnicastPort)
			}
			if in.Spec.Coherence.LocalPort == nil {
				lpa := intstr.FromInt(int(DefaultUnicastPortAdjust))
				in.Spec.Coherence.LocalPortAdjust = &lpa
			}
		}

		// only set defaults for image names in new Coherence instances
		coherenceImage := operator.GetDefaultCoherenceImage()
		in.Spec.EnsureCoherenceImage(&coherenceImage)
		operatorImage := operator.GetDefaultOperatorImage()
		in.Spec.EnsureCoherenceOperatorImage(&operatorImage)

		// Set the features supported by this version
		in.AddAnnotation(AnnotationFeatureSuspend, "true")
	} else {
		logger.Info("Updating defaults for existing resource")
		// this is an update
	}

	// apply the Operator version annotation
	in.AddAnnotation(AnnotationOperatorVersion, operator.GetVersion())

	// apply a label with the hash of the spec - ths must be the last action here to make sure that
	// any modifications to the spec field are included in the hash
	if hash, applied := EnsureHashLabel(in); applied {
		logger.Info(fmt.Sprintf("Applied %s label", LabelCoherenceHash), "hash", hash)
	}
}

func (in *Coherence) AddAnnotation(key, value string) {
	if in != nil {
		if in.Annotations == nil {
			in.Annotations = make(map[string]string)
		}
		in.Annotations[key] = value
	}
}

// The path in this annotation MUST match the const below
// +kubebuilder:webhook:verbs=create;update,path=/validate-coherence-oracle-com-v1-coherence,mutating=false,failurePolicy=fail,groups=coherence.oracle.com,resources=coherence,versions=v1,name=vcoherence.kb.io

// ValidatingWebHookPath This const MUST match the path in the kubebuilder annotation above
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

	sts := in.Spec.CreateStatefulSet(in)
	stsOld := prev.Spec.CreateStatefulSet(prev)
	errorList := ValidateStatefulSetUpdate(&sts, &stsOld)
	if len(errorList) > 0 {
		return fmt.Errorf("rejecting update as it would have resulted in an invalid statefuleset: %v", errorList)
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
