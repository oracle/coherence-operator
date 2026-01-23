/*
 * Copyright (c) 2020, 2026, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package errorhandling

import (
	"context"
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/go-logr/logr"
	coh "github.com/oracle/coherence-operator/api/v1"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/events"
	"k8s.io/client-go/util/retry"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// ErrorCategory represents the category of an error
type ErrorCategory string

const (
	// ErrorCategoryTransient represents a transient error that may resolve itself with time
	ErrorCategoryTransient ErrorCategory = "Transient"
	// ErrorCategoryPermanent represents a permanent error that requires manual intervention
	ErrorCategoryPermanent ErrorCategory = "Permanent"
	// ErrorCategoryRecoverable represents an error that can be automatically recovered from
	ErrorCategoryRecoverable ErrorCategory = "Recoverable"
	// ErrorCategoryUnknown represents an error of unknown category
	ErrorCategoryUnknown ErrorCategory = "Unknown"

	// AnnotationLastError is the annotation key for the last error message
	AnnotationLastError = "coherence.oracle.com/last-error"
	// AnnotationErrorCount is the annotation key for the error count
	AnnotationErrorCount = "coherence.oracle.com/error-count"
	// AnnotationLastRecoveryAttempt is the annotation key for the last recovery attempt timestamp
	AnnotationLastRecoveryAttempt = "coherence.oracle.com/last-recovery-attempt"
)

// ErrorHandler handles errors in the reconciliation loop
type ErrorHandler struct {
	Client        client.Client
	Log           logr.Logger
	EventRecorder events.EventRecorder
}

// HandleError handles an error in the reconciliation loop
func (eh *ErrorHandler) HandleError(ctx context.Context, err error, resource coh.CoherenceResource, msg string) (reconcile.Result, error) {
	if err == nil {
		return reconcile.Result{}, nil
	}

	// Add stack trace if not already present
	err = WithStack(err)

	// Categorize the error
	category := eh.categorizeError(err)

	// Add caller information to the log
	callerInfo := GetCallerInfo(1)

	// Update error tracking information
	if trackErr := eh.updateErrorTracking(ctx, resource, err.Error(), category); trackErr != nil {
		eh.Log.Error(trackErr, "Failed to update error tracking information")
	}

	// Log the error with category and caller info
	eh.Log.Error(err, msg, "category", category, "caller", callerInfo)

	// Record an event
	eventType := corev1.EventTypeWarning
	eventReason := "ReconcileError"
	eventMsg := fmt.Sprintf("%s: %s (Category: %s)", msg, err.Error(), category)
	eh.EventRecorder.Eventf(resource, nil, eventType, eventReason, "HandleError", eventMsg)

	// Update the status to reflect the error
	if statusErr := eh.updateStatus(ctx, resource, category); statusErr != nil {
		eh.Log.Error(statusErr, "Failed to update status")
	}

	// Handle the error based on its category
	switch category {
	case ErrorCategoryTransient:
		// For transient errors, requeue with backoff
		return eh.handleTransientError(resource)
	case ErrorCategoryRecoverable:
		// For recoverable errors, attempt recovery
		return eh.attemptRecovery(ctx, err, resource)
	case ErrorCategoryPermanent:
		// For permanent errors, don't requeue
		return reconcile.Result{}, nil
	default:
		// For unknown errors, requeue with a short delay
		return reconcile.Result{RequeueAfter: 5 * time.Second}, nil
	}
}

// categorizeError categorizes an error based on its type and content
func (eh *ErrorHandler) categorizeError(err error) ErrorCategory {
	// Check for Kubernetes API errors
	if apierrors.IsNotFound(err) {
		return ErrorCategoryTransient // Resource might appear later
	}
	if apierrors.IsConflict(err) {
		return ErrorCategoryTransient // Conflict can resolve with retry
	}
	if apierrors.IsServerTimeout(err) || apierrors.IsTimeout(err) {
		return ErrorCategoryTransient // Timeouts are typically transient
	}
	if apierrors.IsTooManyRequests(err) {
		return ErrorCategoryTransient // Rate limiting is transient
	}
	if apierrors.IsServiceUnavailable(err) {
		return ErrorCategoryTransient // Service might become available later
	}
	if apierrors.IsInternalError(err) {
		return ErrorCategoryTransient // Internal errors might resolve
	}
	if apierrors.IsResourceExpired(err) {
		return ErrorCategoryTransient // Resource version expired, can retry
	}

	// Check for specific error strings that indicate transient errors
	errStr := err.Error()
	if contains(errStr, "connection refused", "network is unreachable", "connection reset", "EOF", "broken pipe") {
		return ErrorCategoryTransient
	}
	if contains(errStr, "etcdserver: leader changed", "etcdserver: request timed out") {
		return ErrorCategoryTransient
	}
	if contains(errStr, "net/http: TLS handshake timeout", "i/o timeout", "deadline exceeded") {
		return ErrorCategoryTransient
	}
	if contains(errStr, "the object has been modified", "optimistic concurrency error") {
		return ErrorCategoryTransient
	}

	// Check for specific error strings that indicate permanent errors
	if contains(errStr, "invalid", "not supported", "forbidden", "unauthorized", "denied") {
		return ErrorCategoryPermanent
	}
	if contains(errStr, "admission webhook", "validation failed") && contains(errStr, "denied") {
		return ErrorCategoryPermanent
	}
	if contains(errStr, "field is immutable") {
		return ErrorCategoryPermanent
	}

	// Check for specific error strings that indicate recoverable errors
	if contains(errStr, "failed to suspend services") {
		return ErrorCategoryRecoverable
	}
	if contains(errStr, "pod disruption budget") && contains(errStr, "not available") {
		return ErrorCategoryRecoverable
	}
	if contains(errStr, "cannot patch") && contains(errStr, "StatefulSet") {
		return ErrorCategoryRecoverable
	}
	if contains(errStr, "cannot create pods") && contains(errStr, "quota exceeded") {
		return ErrorCategoryRecoverable
	}
	if contains(errStr, "cannot schedule") || contains(errStr, "unschedulable") {
		return ErrorCategoryRecoverable
	}

	// Default to unknown
	return ErrorCategoryUnknown
}

// contains checks if any of the substrings are contained in the string
func contains(s string, substrings ...string) bool {
	for _, substr := range substrings {
		if strings.Contains(s, substr) {
			return true
		}
	}
	return false
}

// Common error creation helpers for specific operations

// NewCreateResourceError creates an error for resource creation failures
func NewCreateResourceError(resource string, namespace string, err error) error {
	return NewResourceError("create", resource, namespace, err)
}

// NewUpdateResourceError creates an error for resource update failures
func NewUpdateResourceError(resource string, namespace string, err error) error {
	return NewResourceError("update", resource, namespace, err)
}

// NewDeleteResourceError creates an error for resource deletion failures
func NewDeleteResourceError(resource string, namespace string, err error) error {
	return NewResourceError("delete", resource, namespace, err)
}

// NewGetResourceError creates an error for resource retrieval failures
func NewGetResourceError(resource string, namespace string, err error) error {
	return NewResourceError("get", resource, namespace, err)
}

// NewListResourceError creates an error for resource listing failures
func NewListResourceError(resource string, namespace string, err error) error {
	return NewResourceError("list", resource, namespace, err)
}

// NewPatchResourceError creates an error for resource patching failures
func NewPatchResourceError(resource string, namespace string, err error) error {
	return NewResourceError("patch", resource, namespace, err)
}

// NewReconcileError creates an error for reconciliation failures
func NewReconcileError(resource string, namespace string, err error) error {
	return NewResourceError("reconcile", resource, namespace, err)
}

// NewValidationError creates an error for validation failures
func NewValidationError(resource string, namespace string, err error) error {
	return NewResourceError("validate", resource, namespace, err).WithContext("reason", "validation")
}

// NewTimeoutError creates an error for timeout failures
func NewTimeoutError(operation string, resource string, namespace string, err error) error {
	return NewResourceError(operation, resource, namespace, err).WithContext("reason", "timeout")
}

// NewConnectionError creates an error for connection failures
func NewConnectionError(endpoint string, err error) error {
	return NewOperationError("connect", err).WithContext("endpoint", endpoint)
}

// NewAuthenticationError creates an error for authentication failures
func NewAuthenticationError(resource string, err error) error {
	return NewOperationError("authenticate", err).WithContext("resource", resource)
}

// NewAuthorizationError creates an error for authorization failures
func NewAuthorizationError(resource string, err error) error {
	return NewOperationError("authorize", err).WithContext("resource", resource)
}

// updateErrorTracking updates the error tracking information on the resource
func (eh *ErrorHandler) updateErrorTracking(ctx context.Context, resource coh.CoherenceResource, errMsg string, category ErrorCategory) error {
	// Get the latest version of the resource
	latest := resource.DeepCopyResource()
	if err := eh.Client.Get(ctx, types.NamespacedName{
		Namespace: resource.GetNamespace(),
		Name:      resource.GetName(),
	}, latest); err != nil {
		return err
	}

	// Get the annotations
	annotations := latest.GetAnnotations()
	if annotations == nil {
		annotations = make(map[string]string)
	}

	// Update the last error message
	annotations[AnnotationLastError] = errMsg

	// Update the error count
	count := 1
	if countStr, ok := annotations[AnnotationErrorCount]; ok {
		if parsedCount, err := strconv.Atoi(countStr); err == nil {
			count = parsedCount + 1
		}
	}
	annotations[AnnotationErrorCount] = strconv.Itoa(count)

	// Set the annotations back on the resource
	latest.SetAnnotations(annotations)

	// Update the resource
	return eh.Client.Update(ctx, latest)
}

// updateStatus updates the status of the resource based on the error category
func (eh *ErrorHandler) updateStatus(ctx context.Context, resource coh.CoherenceResource, category ErrorCategory) error {
	// Get the latest version of the resource
	latest := resource.DeepCopyResource()
	if err := eh.Client.Get(ctx, types.NamespacedName{
		Namespace: resource.GetNamespace(),
		Name:      resource.GetName(),
	}, latest); err != nil {
		return err
	}

	// Update the status based on the error category
	status := latest.GetStatus()

	// Set the condition based on the error category
	var condition coh.Condition
	switch category {
	case ErrorCategoryTransient:
		condition = coh.Condition{
			Type:    coh.ConditionTypeFailed,
			Status:  corev1.ConditionTrue,
			Reason:  "TransientError",
			Message: "Encountered a transient error, will retry",
		}
	case ErrorCategoryRecoverable:
		condition = coh.Condition{
			Type:    coh.ConditionTypeFailed,
			Status:  corev1.ConditionTrue,
			Reason:  "RecoverableError",
			Message: "Encountered a recoverable error, attempting recovery",
		}
	case ErrorCategoryPermanent:
		condition = coh.Condition{
			Type:    coh.ConditionTypeFailed,
			Status:  corev1.ConditionTrue,
			Reason:  "PermanentError",
			Message: "Encountered a permanent error, manual intervention required",
		}
	default:
		condition = coh.Condition{
			Type:    coh.ConditionTypeFailed,
			Status:  corev1.ConditionTrue,
			Reason:  "UnknownError",
			Message: "Encountered an error of unknown category",
		}
	}

	// Set the condition on the status
	status.SetCondition(latest, condition)

	// Update the resource status
	return eh.Client.Status().Update(ctx, latest)
}

// handleTransientError handles a transient error
func (eh *ErrorHandler) handleTransientError(resource coh.CoherenceResource) (reconcile.Result, error) {
	// Get the error count from the annotations
	annotations := resource.GetAnnotations()
	count := 1
	if countStr, ok := annotations[AnnotationErrorCount]; ok {
		if parsedCount, err := strconv.Atoi(countStr); err == nil {
			count = parsedCount
		}
	}

	// Calculate backoff duration based on error count (exponential backoff)
	// Start with 5 seconds, double each time, cap at 5 minutes
	backoff := time.Duration(min(5*pow(2, count-1), 300)) * time.Second

	eh.Log.Info("Requeuing with backoff due to transient error",
		"resource", resource.GetName(),
		"namespace", resource.GetNamespace(),
		"errorCount", count,
		"backoff", backoff)

	return reconcile.Result{RequeueAfter: backoff}, nil
}

// attemptRecovery attempts to recover from a recoverable error
func (eh *ErrorHandler) attemptRecovery(ctx context.Context, err error, resource coh.CoherenceResource) (reconcile.Result, error) {
	// Get the latest version of the resource
	latest := resource.DeepCopyResource()
	if err := eh.Client.Get(ctx, types.NamespacedName{
		Namespace: resource.GetNamespace(),
		Name:      resource.GetName(),
	}, latest); err != nil {
		return reconcile.Result{RequeueAfter: time.Minute}, err
	}

	// Get the annotations
	annotations := latest.GetAnnotations()
	if annotations == nil {
		annotations = make(map[string]string)
	}

	// Check if we've already attempted recovery recently
	if lastAttempt, ok := annotations[AnnotationLastRecoveryAttempt]; ok {
		// Parse the last attempt time
		lastAttemptTime, parseErr := time.Parse(time.RFC3339, lastAttempt)
		if parseErr == nil {
			// If the last attempt was less than 5 minutes ago, don't try again yet
			if time.Since(lastAttemptTime) < 5*time.Minute {
				eh.Log.Info("Skipping recovery attempt - too soon since last attempt",
					"resource", resource.GetName(),
					"namespace", resource.GetNamespace(),
					"lastAttempt", lastAttempt)
				return reconcile.Result{RequeueAfter: 1 * time.Minute}, nil
			}
		}
	}

	// Record the recovery attempt
	annotations[AnnotationLastRecoveryAttempt] = time.Now().Format(time.RFC3339)
	latest.SetAnnotations(annotations)

	// Update the resource
	if err := eh.Client.Update(ctx, latest); err != nil {
		return reconcile.Result{RequeueAfter: time.Minute}, err
	}

	// Log the recovery attempt with detailed diagnostics
	eh.Log.Info("Attempting recovery from error",
		"resource", resource.GetName(),
		"namespace", resource.GetNamespace(),
		"error", err.Error(),
		"resourceVersion", latest.GetResourceVersion(),
		"generation", latest.GetGeneration(),
		"deletionTimestamp", latest.GetDeletionTimestamp())

	// Record an event with detailed information
	eh.EventRecorder.Eventf(resource, nil, corev1.EventTypeNormal, "RecoveryAttempt", "AttemptRecovery",
		fmt.Sprintf("Attempting to recover from error: %s", err.Error()))

	// Implement recovery logic based on the error
	errStr := err.Error()

	// Service suspension failure
	if strings.Contains(errStr, "failed to suspend services") {
		return eh.recoverFromServiceSuspensionFailure(ctx, resource)
	}

	// Pod disruption budget issues
	if strings.Contains(errStr, "pod disruption budget") && strings.Contains(errStr, "not available") {
		return eh.recoverFromPDBIssue(ctx, resource)
	}

	// StatefulSet patching issues
	if strings.Contains(errStr, "cannot patch") && strings.Contains(errStr, "StatefulSet") {
		return eh.recoverFromStatefulSetPatchIssue(ctx, resource)
	}

	// Resource quota issues
	if strings.Contains(errStr, "cannot create pods") && strings.Contains(errStr, "quota exceeded") {
		return eh.recoverFromQuotaIssue(ctx, resource)
	}

	// Scheduling issues
	if strings.Contains(errStr, "cannot schedule") || strings.Contains(errStr, "unschedulable") {
		return eh.recoverFromSchedulingIssue(ctx, resource)
	}

	// If we don't have specific recovery logic for this error, log it and requeue with backoff
	eh.Log.Info("No specific recovery mechanism for this error, will retry later",
		"resource", resource.GetName(),
		"namespace", resource.GetNamespace(),
		"error", err.Error())

	// Add diagnostic information to the resource
	annotations["coherence.oracle.com/last-unhandled-error"] = err.Error()
	annotations["coherence.oracle.com/diagnostic-info"] = fmt.Sprintf(
		"No specific recovery mechanism available. Error: %s, Time: %s",
		err.Error(),
		time.Now().Format(time.RFC3339))

	if updateErr := eh.Client.Update(ctx, latest); updateErr != nil {
		eh.Log.Error(updateErr, "Failed to update resource with diagnostic information")
	}

	return reconcile.Result{RequeueAfter: 30 * time.Second}, nil
}

// recoverFromServiceSuspensionFailure attempts to recover from a service suspension failure
func (eh *ErrorHandler) recoverFromServiceSuspensionFailure(ctx context.Context, resource coh.CoherenceResource) (reconcile.Result, error) {
	// 1. Log the recovery attempt
	eh.Log.Info("Attempting to recover from service suspension failure",
		"resource", resource.GetName(),
		"namespace", resource.GetNamespace())

	// 2. Record an event
	eh.EventRecorder.Eventf(resource, nil, corev1.EventTypeNormal, "RecoveryAttempt", "RecoverServiceSuspension",
		"Attempting to recover from service suspension failure")

	// 3. Implement the recovery logic
	// For service suspension failures, we'll try to force remove the finalizer
	// This allows the resource to be deleted even if service suspension failed
	latest := resource.DeepCopyResource()
	if err := eh.Client.Get(ctx, types.NamespacedName{
		Namespace: resource.GetNamespace(),
		Name:      resource.GetName(),
	}, latest); err != nil {
		return reconcile.Result{RequeueAfter: time.Minute}, err
	}

	// Check if this is a deletion and has the Coherence finalizer
	if latest.GetDeletionTimestamp() != nil {
		// This is a deletion, so we'll try to remove the finalizer
		annotations := latest.GetAnnotations()
		if annotations == nil {
			annotations = make(map[string]string)
		}

		// Add an annotation to indicate we're bypassing the finalizer
		annotations["coherence.oracle.com/finalizer-bypass"] = "true"
		latest.SetAnnotations(annotations)

		// Update the resource with the annotation
		if err := eh.Client.Update(ctx, latest); err != nil {
			eh.Log.Error(err, "Failed to add finalizer bypass annotation")
			return reconcile.Result{RequeueAfter: time.Minute}, err
		}

		eh.Log.Info("Added finalizer bypass annotation to resource",
			"resource", resource.GetName(),
			"namespace", resource.GetNamespace())

		eh.EventRecorder.Eventf(resource, nil, corev1.EventTypeNormal, "RecoveryAction", "BypassFinalizer",
			"Added finalizer bypass annotation to allow deletion despite service suspension failure")
	}

	// 4. Return a result that requeues after a short delay to check if recovery was successful
	return reconcile.Result{RequeueAfter: 10 * time.Second}, nil
}

// recoverFromPDBIssue attempts to recover from Pod Disruption Budget issues
func (eh *ErrorHandler) recoverFromPDBIssue(ctx context.Context, resource coh.CoherenceResource) (reconcile.Result, error) {
	// Log the recovery attempt
	eh.Log.Info("Attempting to recover from Pod Disruption Budget issue",
		"resource", resource.GetName(),
		"namespace", resource.GetNamespace())

	// Record an event
	eh.EventRecorder.Eventf(resource, nil, corev1.EventTypeNormal, "RecoveryAttempt", "RecoverPDB",
		"Attempting to recover from Pod Disruption Budget issue")

	// For PDB issues, we'll add an annotation to indicate that we're aware of the issue
	// This can be used by operators or automation to take appropriate action
	latest := resource.DeepCopyResource()
	if err := eh.Client.Get(ctx, types.NamespacedName{
		Namespace: resource.GetNamespace(),
		Name:      resource.GetName(),
	}, latest); err != nil {
		return reconcile.Result{RequeueAfter: time.Minute}, err
	}

	// Add diagnostic information
	annotations := latest.GetAnnotations()
	if annotations == nil {
		annotations = make(map[string]string)
	}

	annotations["coherence.oracle.com/pdb-issue-detected"] = "true"
	annotations["coherence.oracle.com/pdb-issue-time"] = time.Now().Format(time.RFC3339)
	latest.SetAnnotations(annotations)

	// Update the resource with the annotation
	if err := eh.Client.Update(ctx, latest); err != nil {
		eh.Log.Error(err, "Failed to add PDB issue annotation")
		return reconcile.Result{RequeueAfter: time.Minute}, err
	}

	eh.Log.Info("Added PDB issue annotation to resource",
		"resource", resource.GetName(),
		"namespace", resource.GetNamespace())

	eh.EventRecorder.Eventf(resource, nil, corev1.EventTypeNormal, "RecoveryAction", "AnnotatePDBIssue",
		"Added PDB issue annotation to resource")

	// Return a result that requeues after a delay
	return reconcile.Result{RequeueAfter: 30 * time.Second}, nil
}

// recoverFromStatefulSetPatchIssue attempts to recover from StatefulSet patching issues
func (eh *ErrorHandler) recoverFromStatefulSetPatchIssue(ctx context.Context, resource coh.CoherenceResource) (reconcile.Result, error) {
	// Log the recovery attempt
	eh.Log.Info("Attempting to recover from StatefulSet patching issue",
		"resource", resource.GetName(),
		"namespace", resource.GetNamespace())

	// Record an event
	eh.EventRecorder.Eventf(resource, nil, corev1.EventTypeNormal, "RecoveryAttempt", "RecoverStatefulSetPatch",
		"Attempting to recover from StatefulSet patching issue")

	// For StatefulSet patching issues, we'll add an annotation to force a recreation
	latest := resource.DeepCopyResource()
	if err := eh.Client.Get(ctx, types.NamespacedName{
		Namespace: resource.GetNamespace(),
		Name:      resource.GetName(),
	}, latest); err != nil {
		return reconcile.Result{RequeueAfter: time.Minute}, err
	}

	// Add diagnostic information
	annotations := latest.GetAnnotations()
	if annotations == nil {
		annotations = make(map[string]string)
	}

	annotations["coherence.oracle.com/force-reconcile"] = time.Now().Format(time.RFC3339)
	latest.SetAnnotations(annotations)

	// Update the resource with the annotation
	if err := eh.Client.Update(ctx, latest); err != nil {
		eh.Log.Error(err, "Failed to add force-reconcile annotation")
		return reconcile.Result{RequeueAfter: time.Minute}, err
	}

	eh.Log.Info("Added force-reconcile annotation to resource",
		"resource", resource.GetName(),
		"namespace", resource.GetNamespace())

	eh.EventRecorder.Eventf(resource, nil, corev1.EventTypeNormal, "RecoveryAction", "ForceReconcile",
		"Added force-reconcile annotation to resource to address StatefulSet patching issue")

	// Return a result that requeues after a delay
	return reconcile.Result{RequeueAfter: 10 * time.Second}, nil
}

// recoverFromQuotaIssue attempts to recover from resource quota issues
func (eh *ErrorHandler) recoverFromQuotaIssue(ctx context.Context, resource coh.CoherenceResource) (reconcile.Result, error) {
	// Log the recovery attempt
	eh.Log.Info("Attempting to recover from resource quota issue",
		"resource", resource.GetName(),
		"namespace", resource.GetNamespace())

	// Record an event
	eh.EventRecorder.Eventf(resource, nil, corev1.EventTypeWarning, "RecoveryAttempt", "RecoverQuota",
		"Attempting to recover from resource quota issue - this may require manual intervention")

	// For quota issues, we'll add an annotation to indicate the issue
	latest := resource.DeepCopyResource()
	if err := eh.Client.Get(ctx, types.NamespacedName{
		Namespace: resource.GetNamespace(),
		Name:      resource.GetName(),
	}, latest); err != nil {
		return reconcile.Result{RequeueAfter: time.Minute}, err
	}

	// Add diagnostic information
	annotations := latest.GetAnnotations()
	if annotations == nil {
		annotations = make(map[string]string)
	}

	annotations["coherence.oracle.com/quota-issue-detected"] = "true"
	annotations["coherence.oracle.com/quota-issue-time"] = time.Now().Format(time.RFC3339)
	latest.SetAnnotations(annotations)

	// Update the resource with the annotation
	if err := eh.Client.Update(ctx, latest); err != nil {
		eh.Log.Error(err, "Failed to add quota issue annotation")
		return reconcile.Result{RequeueAfter: time.Minute}, err
	}

	eh.Log.Info("Added quota issue annotation to resource",
		"resource", resource.GetName(),
		"namespace", resource.GetNamespace())

	eh.EventRecorder.Eventf(resource, nil, corev1.EventTypeWarning, "ResourceQuotaIssue", "QuotaIssue",
		"Resource quota exceeded - manual intervention may be required to increase quota or reduce resource requests")

	// Return a result that requeues after a longer delay
	return reconcile.Result{RequeueAfter: 2 * time.Minute}, nil
}

// recoverFromSchedulingIssue attempts to recover from pod scheduling issues
func (eh *ErrorHandler) recoverFromSchedulingIssue(ctx context.Context, resource coh.CoherenceResource) (reconcile.Result, error) {
	// Log the recovery attempt
	eh.Log.Info("Attempting to recover from pod scheduling issue",
		"resource", resource.GetName(),
		"namespace", resource.GetNamespace())

	// Record an event
	eh.EventRecorder.Eventf(resource, nil, corev1.EventTypeWarning, "RecoveryAttempt", "RecoverScheduling",
		"Attempting to recover from pod scheduling issue - this may require manual intervention")

	// For scheduling issues, we'll add an annotation to indicate the issue
	latest := resource.DeepCopyResource()
	if err := eh.Client.Get(ctx, types.NamespacedName{
		Namespace: resource.GetNamespace(),
		Name:      resource.GetName(),
	}, latest); err != nil {
		return reconcile.Result{RequeueAfter: time.Minute}, err
	}

	// Add diagnostic information
	annotations := latest.GetAnnotations()
	if annotations == nil {
		annotations = make(map[string]string)
	}

	annotations["coherence.oracle.com/scheduling-issue-detected"] = "true"
	annotations["coherence.oracle.com/scheduling-issue-time"] = time.Now().Format(time.RFC3339)
	latest.SetAnnotations(annotations)

	// Update the resource with the annotation
	if err := eh.Client.Update(ctx, latest); err != nil {
		eh.Log.Error(err, "Failed to add scheduling issue annotation")
		return reconcile.Result{RequeueAfter: time.Minute}, err
	}

	eh.Log.Info("Added scheduling issue annotation to resource",
		"resource", resource.GetName(),
		"namespace", resource.GetNamespace())

	eh.EventRecorder.Eventf(resource, nil, corev1.EventTypeWarning, "SchedulingIssue", "SchedulingIssue",
		"Pod scheduling issue detected - manual intervention may be required to address node resources or affinity rules")

	// Return a result that requeues after a longer delay
	return reconcile.Result{RequeueAfter: 2 * time.Minute}, nil
}

// RetryOnError retries the given function with exponential backoff on error
func (eh *ErrorHandler) RetryOnError(operation string, fn func() error) error {
	return retry.OnError(retry.DefaultRetry, func(err error) bool {
		category := eh.categorizeError(err)
		retryable := category == ErrorCategoryTransient
		if retryable {
			eh.Log.Info("Retrying operation due to transient error",
				"operation", operation,
				"error", err.Error())
		}
		return retryable
	}, fn)
}

// RetryWithContext retries the given function with context and additional metadata
func (eh *ErrorHandler) RetryWithContext(ctx context.Context, operation string, resource string, namespace string, fn func() error) error {
	return retry.OnError(retry.DefaultRetry, func(err error) bool {
		category := eh.categorizeError(err)
		retryable := category == ErrorCategoryTransient
		if retryable {
			eh.Log.Info("Retrying operation due to transient error",
				"operation", operation,
				"resource", resource,
				"namespace", namespace,
				"error", err.Error())
		}
		return retryable
	}, func() error {
		err := fn()
		if err != nil {
			// Wrap the error with operation context
			return NewResourceError(operation, resource, namespace, err)
		}
		return nil
	})
}

// HandleResourceError handles an error related to a specific resource
func (eh *ErrorHandler) HandleResourceError(ctx context.Context, err error, resource coh.CoherenceResource, operation string, msg string) (reconcile.Result, error) {
	if err == nil {
		return reconcile.Result{}, nil
	}

	// Create a resource error with context
	resourceErr := NewResourceError(
		operation,
		resource.GetName(),
		resource.GetNamespace(),
		err,
	)

	// Add additional context if available
	if opErr, ok := err.(*OperationError); ok {
		for k, v := range opErr.Context {
			_ = resourceErr.WithContext(k, v)
		}
	}

	return eh.HandleError(ctx, resourceErr, resource, msg)
}

// MustNotError panics if the error is not nil
func MustNotError(err error) {
	if err != nil {
		panic(err)
	}
}

// LogAndReturnError logs the error and returns it
func (eh *ErrorHandler) LogAndReturnError(err error, msg string) error {
	if err == nil {
		return nil
	}

	eh.Log.Error(err, msg)
	return err
}

// LogAndWrapError logs the error, wraps it with the message, and returns it
func (eh *ErrorHandler) LogAndWrapError(err error, msg string) error {
	if err == nil {
		return nil
	}

	wrappedErr := WrapError(err, msg)
	eh.Log.Error(wrappedErr, msg)
	return wrappedErr
}

// LogAndWrapErrorf logs the error, wraps it with the formatted message, and returns it
func (eh *ErrorHandler) LogAndWrapErrorf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	msg := fmt.Sprintf(format, args...)
	wrappedErr := WrapError(err, msg)
	eh.Log.Error(wrappedErr, msg)
	return wrappedErr
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// pow returns a^b for integers
func pow(a, b int) int {
	result := 1
	for i := 0; i < b; i++ {
		result *= a
	}
	return result
}

// OperationError represents an error that occurred during an operation
type OperationError struct {
	Operation string
	Resource  string
	Namespace string
	Err       error
	Context   map[string]string
}

// Error returns the error message
func (e *OperationError) Error() string {
	msg := fmt.Sprintf("operation '%s' failed", e.Operation)
	if e.Resource != "" {
		if e.Namespace != "" {
			msg += fmt.Sprintf(" for resource '%s' in namespace '%s'", e.Resource, e.Namespace)
		} else {
			msg += fmt.Sprintf(" for resource '%s'", e.Resource)
		}
	}

	if len(e.Context) > 0 {
		contextStr := ""
		for k, v := range e.Context {
			if contextStr != "" {
				contextStr += ", "
			}
			contextStr += fmt.Sprintf("%s=%s", k, v)
		}
		msg += fmt.Sprintf(" (context: %s)", contextStr)
	}

	if e.Err != nil {
		msg += fmt.Sprintf(": %v", e.Err)
	}

	return msg
}

// Unwrap returns the underlying error
func (e *OperationError) Unwrap() error {
	return e.Err
}

// Cause returns the underlying error (for compatibility with github.com/pkg/errors)
func (e *OperationError) Cause() error {
	return e.Err
}

// WithContext adds context to the error
func (e *OperationError) WithContext(key, value string) *OperationError {
	if e.Context == nil {
		e.Context = make(map[string]string)
	}
	e.Context[key] = value
	return e
}

// NewOperationError creates a new OperationError
func NewOperationError(operation string, err error) *OperationError {
	return &OperationError{
		Operation: operation,
		Err:       err,
		Context:   make(map[string]string),
	}
}

// NewResourceError creates a new OperationError for a specific resource
func NewResourceError(operation string, resource string, namespace string, err error) *OperationError {
	return &OperationError{
		Operation: operation,
		Resource:  resource,
		Namespace: namespace,
		Err:       err,
		Context:   make(map[string]string),
	}
}

// WrapError wraps an error with context information
func WrapError(err error, message string) error {
	if err == nil {
		return nil
	}
	return errors.Wrap(err, message)
}

// WrapErrorf wraps an error with formatted context information
func WrapErrorf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return errors.Wrapf(err, format, args...)
}

// NewError creates a new error with the given message
func NewError(message string) error {
	return errors.New(message)
}

// NewErrorf creates a new error with the given formatted message
func NewErrorf(format string, args ...interface{}) error {
	return errors.Errorf(format, args...)
}

// IsNotFound returns true if the error is a not found error
func IsNotFound(err error) bool {
	return apierrors.IsNotFound(err)
}

// IsAlreadyExists returns true if the error is an already exists error
func IsAlreadyExists(err error) bool {
	return apierrors.IsAlreadyExists(err)
}

// IsConflict returns true if the error is a conflict error
func IsConflict(err error) bool {
	return apierrors.IsConflict(err)
}

// IsInvalid returns true if the error is an invalid error
func IsInvalid(err error) bool {
	return apierrors.IsInvalid(err)
}

// IsForbidden returns true if the error is a forbidden error
func IsForbidden(err error) bool {
	return apierrors.IsForbidden(err)
}

// IsTimeout returns true if the error is a timeout error
func IsTimeout(err error) bool {
	return apierrors.IsTimeout(err)
}

// IsServerTimeout returns true if the error is a server timeout error
func IsServerTimeout(err error) bool {
	return apierrors.IsServerTimeout(err)
}

// IsTooManyRequests returns true if the error is a too many requests error
func IsTooManyRequests(err error) bool {
	return apierrors.IsTooManyRequests(err)
}

// IsServiceUnavailable returns true if the error is a service unavailable error
func IsServiceUnavailable(err error) bool {
	return apierrors.IsServiceUnavailable(err)
}

// WithStack adds a stack trace to the error if it doesn't already have one
func WithStack(err error) error {
	if err == nil {
		return nil
	}
	return errors.WithStack(err)
}

// GetCallerInfo returns the file and line number of the caller
func GetCallerInfo(skip int) string {
	_, file, line, ok := runtime.Caller(skip + 1)
	if !ok {
		return "unknown:0"
	}

	// Extract just the filename from the full path
	parts := strings.Split(file, "/")
	file = parts[len(parts)-1]

	return fmt.Sprintf("%s:%d", file, line)
}

// NewErrorHandler creates a new ErrorHandler
func NewErrorHandler(client client.Client, log logr.Logger, recorder events.EventRecorder) *ErrorHandler {
	return &ErrorHandler{
		Client:        client,
		Log:           log,
		EventRecorder: recorder,
	}
}
