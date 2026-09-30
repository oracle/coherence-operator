/*
 * Copyright (c) 2020, 2026, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package finalizer

import (
	"context"
	"strconv"

	"github.com/go-logr/logr"
	coh "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/pkg/events"
	"github.com/oracle/coherence-operator/pkg/operator"
	"github.com/oracle/coherence-operator/pkg/patching"
	"github.com/oracle/coherence-operator/pkg/probe"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	events2 "k8s.io/client-go/tools/events"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// FinalizerManager manages finalizers for Coherence resources
type FinalizerManager struct {
	Client        client.Client
	Log           logr.Logger
	EventRecorder events2.EventRecorder
	Patcher       patching.ResourcePatcher
}

// EnsureFinalizerApplied ensures the finalizer is applied to the Coherence resource
func (fm *FinalizerManager) EnsureFinalizerApplied(ctx context.Context, c *coh.Coherence) (bool, error) {
	if !controllerutil.ContainsFinalizer(c, coh.CoherenceFinalizer) {
		// Re-fetch the Coherence resource to ensure we have the most recent copy
		latest := &coh.Coherence{}
		c.DeepCopyInto(latest)
		controllerutil.AddFinalizer(latest, coh.CoherenceFinalizer)

		callback := func() {
			fm.Log.Info("Added finalizer to Coherence resource", "Namespace", c.Namespace, "Name", c.Name, "Finalizer", coh.CoherenceFinalizer)
		}

		// Perform a three-way patch to apply the finalizer
		applied, err := fm.Patcher.ThreeWayPatchWithCallback(ctx, c.Name, c, c, latest, callback)
		if err != nil {
			return false, errors.Wrapf(err, "failed to update Coherence resource %s/%s with finalizer", c.Namespace, c.Name)
		}
		return applied, nil
	}
	return false, nil
}

// EnsureFinalizerRemoved ensures the finalizer is removed from the Coherence resource
func (fm *FinalizerManager) EnsureFinalizerRemoved(ctx context.Context, c *coh.Coherence) error {
	if controllerutil.RemoveFinalizer(c, coh.CoherenceFinalizer) {
		err := fm.Client.Update(ctx, c)
		if err != nil {
			fm.Log.Info("Failed to remove the finalizer from the Coherence resource, it looks like it had already been deleted")
			return err
		}
	}
	return nil
}

// FinalizeDeployment performs any required finalizer tasks for the Coherence resource
func (fm *FinalizerManager) FinalizeDeployment(ctx context.Context, c *coh.Coherence, findStatefulSet func(ctx context.Context, namespace, name string) (*appsv1.StatefulSet, bool, error)) error {
	// Preserve the existing deletion escape hatch, including presence-only semantics.
	annotations := c.GetAnnotations()
	if _, bypass := annotations["coherence.oracle.com/finalizer-bypass"]; bypass {
		fm.Log.Info("Bypassing service suspension due to finalizer-bypass annotation",
			"Namespace", c.Namespace, "Name", c.Name)
		fm.EventRecorder.Eventf(c, nil, corev1.EventTypeNormal, "FinalizerBypassed", "Finalize",
			"Service suspension bypassed due to finalizer-bypass annotation")
		return nil
	}

	// determine whether we can skip service suspension
	if viper.GetBool(operator.FlagSkipServiceSuspend) {
		fm.Log.Info("Skipping suspension of Coherence services in deployment " + c.Name +
			operator.FlagSkipServiceSuspend + " is set to true")
		return nil
	}
	if !c.Spec.IsSuspendServicesOnShutdown() {
		fm.Log.Info("Skipping suspension of Coherence services in deployment " + c.Name +
			" Spec.SuspendServicesOnShutdown is set to false")
		return nil
	}
	if c.GetReplicas() == 0 {
		fm.Log.Info("Skipping suspension of Coherence services in deployment " + c.Name +
			" Spec.Replicas is zero")
		return nil
	}

	fm.Log.Info("Finalizing Coherence resource", "Namespace", c.Namespace, "Name", c.Name)
	// Get the StatefulSet
	sts, stsExists, err := findStatefulSet(ctx, c.Namespace, c.Name)
	if err != nil {
		return fm.suspensionFailure(c, &probe.SuspensionCauseError{Kind: probe.SuspensionCauseWorkloadLookupFailed, Cause: err})
	}
	if stsExists {
		if sts.Status.ReadyReplicas == 0 {
			fm.Log.Info("Skipping suspension of Coherence services in deployment " + c.Name + " - No Pods are ready")
		} else {
			p := probe.CoherenceProbe{Client: fm.Client, EventRecorder: events.NewOwnedEventRecorder(c, fm.EventRecorder)}
			_, err = p.SuspendServicesWithError(ctx, c, sts)
			if err != nil {
				// Retain the established recovery policy. This is a recorded general
				// error count, not a count of suspension attempts made here.
				errorCount, parseErr := strconv.Atoi(annotations["coherence.oracle.com/error-count"])
				if parseErr == nil && errorCount > 3 {
					fm.Log.Info("Service suspension failed multiple times, allowing deletion to proceed",
						"Namespace", c.Namespace, "Name", c.Name, "ErrorCount", errorCount)
					fm.EventRecorder.Eventf(c, nil, corev1.EventTypeWarning, "ServiceSuspensionFailed", "SuspendServices",
						"Service suspension failed multiple times, allowing deletion to proceed anyway")
					return nil
				}
			}
			return err
		}
	}
	return nil
}

func (fm *FinalizerManager) suspensionFailure(c *coh.Coherence, cause error) error {
	err := &probe.SuspensionError{Cause: cause}
	if fm.EventRecorder != nil {
		fm.EventRecorder.Eventf(c, nil, corev1.EventTypeWarning, "ServiceSuspendFailed", "SuspendServices", "%s", err.Error())
	}
	return err
}
