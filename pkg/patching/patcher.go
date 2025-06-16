/*
 * Copyright (c) 2020, 2025, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package patching

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

const (
	// PatchIgnore - If the patching json is this we can skip the patching.
	PatchIgnore = "{\"metadata\":{\"creationTimestamp\":null},\"status\":{\"replicas\":0}}"
)

type ResourcePatcher interface {
	// TwoWayPatch performs a two-way merge patching on the resource.
	TwoWayPatch(context.Context, string, client.Object, client.Object) (bool, error)
	// CreateTwoWayPatch creates a two-way patching between the original state, the current state and the desired state of a k8s resource.
	CreateTwoWayPatch(string, runtime.Object, runtime.Object, ...string) (client.Patch, error)
	// CreateTwoWayPatchOfType creates a two-way patching between the current state and the desired state of a k8s resource.
	CreateTwoWayPatchOfType(types.PatchType, string, runtime.Object, runtime.Object, ...string) (client.Patch, error)
	// ThreeWayPatch performs a three-way merge patching on the resource returning true if a patching was required otherwise false.
	ThreeWayPatch(context.Context, string, client.Object, client.Object, client.Object) (bool, error)
	// ThreeWayPatchWithCallback performs a three-way merge patching on the resource returning true if a patching was required otherwise false.
	ThreeWayPatchWithCallback(context.Context, string, client.Object, client.Object, client.Object, func()) (bool, error)
	// ApplyThreeWayPatchWithCallback performs a three-way merge patching on the resource returning true if a patching was required otherwise false.
	ApplyThreeWayPatchWithCallback(context.Context, string, client.Object, client.Patch, []byte, func()) (bool, error)
	// CreateThreeWayPatch creates a three-way patching between the original state, the current state and the desired state of a k8s resource.
	CreateThreeWayPatch(string, runtime.Object, runtime.Object, runtime.Object, ...string) (client.Patch, []byte, error)
	// CreateThreeWayPatchData creates a three-way patching between the original state, the current state and the desired state of a k8s resource.
	CreateThreeWayPatchData(original, desired, current runtime.Object) ([]byte, error)
	// GetPatchType returns the patching type this patcher uses
	GetPatchType() types.PatchType
	// SetPatchType sets the patching type this patcher uses
	SetPatchType(pt types.PatchType)
}

// NewResourcePatcher creates a new ResourcePatcher
func NewResourcePatcher(mgr manager.Manager, logger logr.Logger, patchType types.PatchType) ResourcePatcher {
	return &patcher{
		mgr:       mgr,
		logger:    logger,
		patchType: patchType,
	}
}

// compile time check to verify `patcher` implements `ResourcePatcher`
var _ ResourcePatcher = &patcher{}

type patcher struct {
	mgr       manager.Manager
	logger    logr.Logger
	patchType types.PatchType
}

func (in *patcher) GetPatchType() types.PatchType { return in.patchType }

func (in *patcher) SetPatchType(pt types.PatchType) { in.patchType = pt }

// TwoWayPatch performs a two-way merge patching on the resource.
func (in *patcher) TwoWayPatch(ctx context.Context, name string, current, desired client.Object) (bool, error) {
	patch, err := in.CreateTwoWayPatch(name, desired, current, PatchIgnore)
	if err != nil {
		kind := current.GetObjectKind().GroupVersionKind().Kind
		return false, errors.Wrapf(err, "failed to create patching for %s/%s", kind, name)
	}

	if patch == nil {
		// nothing to patching so just return
		return false, nil
	}

	err = in.mgr.GetClient().Patch(ctx, current, patch)
	if err != nil {
		kind := current.GetObjectKind().GroupVersionKind().Kind
		return false, errors.Wrapf(err, "cannot patching  %s/%s", kind, name)
	}

	return true, nil
}

// CreateTwoWayPatch creates a two-way patching between the original state, the current state and the desired state of a k8s resource.
func (in *patcher) CreateTwoWayPatch(name string, desired, current runtime.Object, ignore ...string) (client.Patch, error) {
	return in.CreateTwoWayPatchOfType(in.patchType, name, desired, current, ignore...)
}

// CreateTwoWayPatchOfType creates a two-way patching between the current state and the desired state of a k8s resource.
func (in *patcher) CreateTwoWayPatchOfType(patchType types.PatchType, name string, desired, current runtime.Object, ignore ...string) (client.Patch, error) {
	currentData, err := json.Marshal(current)
	if err != nil {
		return nil, errors.Wrap(err, "serializing current configuration")
	}
	desiredData, err := json.Marshal(desired)
	if err != nil {
		return nil, errors.Wrap(err, "serializing desired configuration")
	}

	// Get a versioned object
	versionedObject := in.asVersioned(desired)

	patchMeta, err := strategicpatch.NewPatchMetaFromStruct(versionedObject)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create patching metadata from object")
	}

	data, err := strategicpatch.CreateTwoWayMergePatchUsingLookupPatchMeta(currentData, desiredData, patchMeta)
	if err != nil {
		return nil, errors.Wrap(err, "creating three-way patching")
	}

	// check whether the patching counts as no-patching
	ignore = append(ignore, "{}")
	for _, s := range ignore {
		if string(data) == s {
			// empty patching
			return nil, err
		}
	}

	// log the patching
	kind := current.GetObjectKind().GroupVersionKind().Kind

	in.logger.V(2).Info(fmt.Sprintf("Created patching for %s/%s\n%s", kind, name, string(data)))

	return client.RawPatch(patchType, data), nil
}

// ThreeWayPatch performs a three-way merge patching on the resource returning true if a patching was required otherwise false.
func (in *patcher) ThreeWayPatch(ctx context.Context, name string, current, original, desired client.Object) (bool, error) {
	return in.ThreeWayPatchWithCallback(ctx, name, current, original, desired, nil)
}

// ThreeWayPatchWithCallback performs a three-way merge patching on the resource returning true if a patching was required otherwise false.
func (in *patcher) ThreeWayPatchWithCallback(ctx context.Context, name string, current, original, desired client.Object, callback func()) (bool, error) {
	kind := current.GetObjectKind().GroupVersionKind().Kind
	// fix the CreationTimestamp so that it is not in the patching
	desired.(metav1.Object).SetCreationTimestamp(current.(metav1.Object).GetCreationTimestamp())
	// create the patching
	patch, data, err := in.CreateThreeWayPatch(name, original, desired, current, PatchIgnore)
	if err != nil {
		return false, errors.Wrapf(err, "failed to create patching for %s/%s", kind, name)
	}

	if patch == nil {
		// nothing to patching so just return
		return false, nil
	}

	return in.ApplyThreeWayPatchWithCallback(ctx, name, current, patch, data, callback)
}

// ApplyThreeWayPatchWithCallback performs a three-way merge patching on the resource returning true if a patching was required otherwise false.
func (in *patcher) ApplyThreeWayPatchWithCallback(ctx context.Context, name string, current client.Object, patch client.Patch, data []byte, callback func()) (bool, error) {
	kind := current.GetObjectKind().GroupVersionKind().Kind

	// execute any callback
	if callback != nil {
		callback()
	}

	in.logger.WithValues().Info(fmt.Sprintf("Patching %s/%s", kind, name), "Patch", string(data))
	err := in.mgr.GetClient().Patch(ctx, current, patch)
	if err != nil {
		in.logger.WithValues().Info(fmt.Sprintf("Failed to patch %s/%s", kind, name), "Patch", string(data), "Error", err.Error())
		return false, errors.Wrapf(err, "failed to patch  %s/%s with %s", kind, name, string(data))
	}

	return true, nil
}

// CreateThreeWayPatch creates a three-way patching between the original state, the current state and the desired state of a k8s resource.
func (in *patcher) CreateThreeWayPatch(name string, original, desired, current runtime.Object, ignore ...string) (client.Patch, []byte, error) {
	data, err := in.CreateThreeWayPatchData(original, desired, current)
	if err != nil {
		return nil, data, errors.Wrap(err, "creating three-way patching")
	}

	// check whether the patching counts as no-patching
	ignore = append(ignore, "{}")
	for _, s := range ignore {
		if string(data) == s {
			// empty patching
			return nil, data, err
		}
	}

	// log the patching
	kind := current.GetObjectKind().GroupVersionKind().Kind
	in.logger.Info(fmt.Sprintf("Created patch for %s/%s", kind, name), "Patch", string(data))

	return client.RawPatch(in.patchType, data), data, nil
}

// CreateThreeWayPatchData creates a three-way patching between the original state, the current state and the desired state of a k8s resource.
func (in *patcher) CreateThreeWayPatchData(original, desired, current runtime.Object) ([]byte, error) {
	originalData, err := json.Marshal(original)
	if err != nil {
		return nil, errors.Wrap(err, "serializing original configuration")
	}
	currentData, err := json.Marshal(current)
	if err != nil {
		return nil, errors.Wrap(err, "serializing current configuration")
	}
	desiredData, err := json.Marshal(desired)
	if err != nil {
		return nil, errors.Wrap(err, "serializing desired configuration")
	}

	// Get a versioned object
	versionedObject := in.asVersioned(desired)

	patchMeta, err := strategicpatch.NewPatchMetaFromStruct(versionedObject)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create patching metadata from object")
	}

	return strategicpatch.CreateThreeWayMergePatch(originalData, desiredData, currentData, patchMeta, true)
}

// asVersioned converts the given object into a runtime.Template with the correct group and version set.
func (in *patcher) asVersioned(obj runtime.Object) runtime.Object {
	var gv = runtime.GroupVersioner(schema.GroupVersions(scheme.Scheme.PrioritizedVersionsAllGroups()))
	if obj, err := runtime.ObjectConvertor(scheme.Scheme).ConvertToVersion(obj, gv); err == nil {
		return obj
	}
	return obj
}
