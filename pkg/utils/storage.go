/*
 * Copyright (c) 2020, 2025, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package utils

import (
	"context"
	"encoding/json"
	"fmt"
	coh "github.com/oracle/coherence-operator/api/v1"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var log = logf.Log.WithName("Storage")

const (
	storeKeyLatest   = "latest"
	storeKeyPrevious = "previous"
	labelLegacy      = coh.LabelCoherenceHash + "-legacy"
)

type Storage interface {
	// GetName returns the name of the storage secret.
	GetName() string
	// GetLatest obtains the deployment resources for the specified version
	GetLatest() coh.Resources
	// GetPrevious obtains the deployment resources for the version prior to the specified version
	GetPrevious() coh.Resources
	// Store will store the deployment resources, this will create a new version in the store
	Store(coh.Resources, coh.CoherenceResource) error
	// Destroy will destroy the store
	Destroy()
	// GetHash will return the hash label of the owning resource
	GetHash() string
	// IsJob returns true if the Coherence deployment is a Job
	IsJob(reconcile.Request) bool
	// ResetHash sets the hash for the stprage
	ResetHash(owner coh.CoherenceResource) error
	// CheckHashMatches verifies the specified hash against the storage hash
	CheckHashMatches(hash string) bool
	// CheckObjectHashLabelMatches verifies the specified hash label against the storage hash
	CheckObjectHashLabelMatches(o metav1.Object) bool
}

// NewStorage creates a new storage for the given key.
func NewStorage(key client.ObjectKey, mgr manager.Manager) (Storage, error) {
	return newStorage(key, mgr, nil)
}

// NewStorageForOwner creates a new storage for the given key.
func NewStorageForOwner(key client.ObjectKey, mgr manager.Manager, owner coh.CoherenceResource) (Storage, error) {
	return newStorage(key, mgr, owner)
}

func newStorage(key client.ObjectKey, mgr manager.Manager, owner coh.CoherenceResource) (Storage, error) {
	store := &secretStore{manager: mgr, key: key}
	err := store.loadVersions(owner)
	return store, err
}

type secretStore struct {
	manager    manager.Manager
	key        client.ObjectKey
	latest     coh.Resources
	previous   coh.Resources
	hash       *string
	legacyHash *string
}

func (in *secretStore) IsJob(request reconcile.Request) bool {
	if in == nil {
		return false
	}
	latest := in.GetLatest()
	_, found := latest.GetResource(coh.ResourceTypeJob, request.Name)
	return found
}

func (in *secretStore) createSecretStruct() *corev1.Secret {
	labels := make(map[string]string)
	labels[coh.LabelCoherenceStore] = "true"
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: in.key.Namespace,
			Name:      in.key.Name,
			Labels:    labels,
		},
	}
}

func (in *secretStore) GetName() string {
	if in == nil || in.hash == nil {
		return ""
	}
	return in.key.Name
}

func (in *secretStore) GetHash() string {
	if in == nil || in.hash == nil {
		return ""
	}
	return *in.hash
}

func (in *secretStore) Destroy() {
	secret := in.createSecretStruct()
	if err := in.manager.GetClient().Delete(context.TODO(), secret); err != nil {
		log.Error(err, "Error deleting storage secret", "Namespace", in.key.Namespace, "Name", in.key.Name)
	}
}

func (in *secretStore) GetLatest() coh.Resources {
	if in == nil {
		return coh.Resources{}
	}
	return in.latest
}

func (in *secretStore) GetPrevious() coh.Resources {
	if in == nil {
		return coh.Resources{}
	}
	return in.previous
}

func (in *secretStore) CheckObjectHashLabelMatches(o metav1.Object) bool {
	actual, found := o.GetLabels()[coh.LabelCoherenceHash]
	return found && in.CheckHashMatches(actual)
}

func (in *secretStore) CheckHashMatches(hash string) bool {
	if in == nil {
		return false
	}
	if in.hash != nil && *in.hash == hash {
		return true
	}
	if in.legacyHash != nil && *in.legacyHash == hash {
		return true
	}
	return false
}

func (in *secretStore) ResetHash(owner coh.CoherenceResource) error {
	secret, exists, err := in.getSecret()
	if err != nil {
		// an error occurred other than NotFound
		return err
	}
	labels := secret.GetLabels()
	if labels == nil {
		labels = make(map[string]string)
	}
	legacy, found := labels[coh.LabelCoherenceHash]
	if found {
		labels[labelLegacy] = legacy
	}
	labels[coh.LabelCoherenceHash] = owner.GetGenerationString()
	return in.save(owner, secret, exists)
}

func (in *secretStore) Store(r coh.Resources, owner coh.CoherenceResource) error {
	secret, exists, err := in.getSecret()
	if err != nil {
		// an error occurred other than NotFound
		return err
	}

	r.Version = in.latest.Version + 1

	if secret.Data == nil {
		secret.Data = make(map[string][]byte)
	}

	r.EnsureGVK(in.manager.GetScheme())

	oldLatest := secret.Data[storeKeyLatest]
	newLatest, err := json.Marshal(r)
	if err != nil {
		return err
	}

	labels := secret.GetLabels()
	if labels == nil {
		labels = make(map[string]string)
	}
	labels[coh.LabelCoherenceHash] = owner.GetGenerationString()
	delete(labels, labelLegacy)

	globalLabels := owner.CreateGlobalLabels()
	for k, v := range globalLabels {
		labels[k] = v
	}
	secret.SetLabels(labels)

	ann := secret.GetAnnotations()
	globalAnn := owner.CreateGlobalAnnotations()
	if globalAnn != nil {
		if ann == nil {
			ann = make(map[string]string)
		}
		for k, v := range globalAnn {
			ann[k] = v
		}
	}
	secret.SetAnnotations(ann)

	secret.Data[storeKeyLatest] = newLatest
	secret.Data[storeKeyPrevious] = oldLatest

	err = in.save(owner, secret, exists)

	if err == nil {
		// everything was updated successfully so update the storage state
		in.previous = in.latest
		in.latest = r
	}
	return err
}

func (in *secretStore) save(owner coh.CoherenceResource, secret *corev1.Secret, exists bool) error {
	var err error

	if !exists {
		// the resource does not exist so set the deployment as the controller/owner and create it
		err = controllerutil.SetControllerReference(owner, secret, in.manager.GetScheme())
		if err != nil {
			err = errors.Wrap(err, fmt.Sprintf("setting resource owner/controller in state store %s/%s", secret.Namespace, secret.Name))
		} else {
			err = in.manager.GetClient().Create(context.TODO(), secret)
		}
	} else {
		// the store secret exists so update it
		err = in.manager.GetClient().Update(context.TODO(), secret)
	}
	return err
}

func (in *secretStore) loadVersions(owner coh.CoherenceResource) error {
	secret, exists, err := in.getSecret()
	if err != nil {
		// an error occurred other than NotFound
		return err
	}

	if exists {
		var data []byte
		var found bool

		data, found = secret.Data[storeKeyLatest]
		if found && len(data) > 0 {
			if err = json.Unmarshal(data, &in.latest); err != nil {
				return errors.Wrap(err, "unmarshalling latest store state")
			}
		}
		data, found = secret.Data[storeKeyPrevious]
		if found && len(data) > 0 {
			if err = json.Unmarshal(data, &in.previous); err != nil {
				return errors.Wrap(err, "unmarshalling previous store state")
			}
		}

		labels := secret.GetLabels()
		if hashValue, found := labels[coh.LabelCoherenceHash]; found {
			in.hash = &hashValue
		} else {
			in.hash = nil
		}

		if hashValue, found := labels[labelLegacy]; found {
			in.legacyHash = &hashValue
		} else {
			in.legacyHash = nil
		}
	}
	return nil
}

// getSecret obtains the store Secret from k8s returning the Secret and a bool indicating whether the Secret exists in k8s and any error
func (in *secretStore) getSecret() (*corev1.Secret, bool, error) {
	secret := in.createSecretStruct()
	err := in.manager.GetClient().Get(context.TODO(), in.key, secret)
	switch {
	case err != nil && !apierrors.IsNotFound(err):
		// an error occurred other than NotFound
		return nil, false, err
	case err != nil && apierrors.IsNotFound(err):
		// secret does not exist in k8s
		return secret, false, nil
	default:
		return secret, true, nil
	}
}
