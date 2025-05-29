/*
 * Copyright (c) 2020, 2025, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package resources

import (
	"context"
	"github.com/go-logr/logr"
	coh "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/pkg/rest"
	coreV1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// OperatorSecretManager manages the operator configuration secret
type OperatorSecretManager struct {
	Client client.Client
	Log    logr.Logger
}

// EnsureOperatorSecret ensures that the Operator configuration secret exists in the namespace.
func (osm *OperatorSecretManager) EnsureOperatorSecret(ctx context.Context, deployment *coh.Coherence) error {
	namespace := deployment.Namespace
	s := &coreV1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:        coh.OperatorConfigName,
			Namespace:   namespace,
			Labels:      deployment.CreateGlobalLabels(),
			Annotations: deployment.CreateGlobalAnnotations(),
		},
	}

	err := osm.Client.Get(ctx, types.NamespacedName{Name: coh.OperatorConfigName, Namespace: namespace}, s)
	if err != nil && !apierrors.IsNotFound(err) {
		return err
	}

	restHostAndPort := rest.GetServerHostAndPort()

	oldValue := s.Data[coh.OperatorConfigKeyHost]
	if oldValue == nil || string(oldValue) != restHostAndPort {
		// data is different so create/update

		if s.StringData == nil {
			s.StringData = make(map[string]string)
		}
		s.StringData[coh.OperatorConfigKeyHost] = restHostAndPort

		osm.Log.Info("Operator configuration updated", "Key", coh.OperatorConfigKeyHost, "OldValue", string(oldValue), "NewValue", restHostAndPort)
		if apierrors.IsNotFound(err) {
			// for some reason we're getting here even if the secret exists so delete it!!
			_ = osm.Client.Delete(ctx, s)
			osm.Log.Info("Creating configuration secret " + coh.OperatorConfigName)
			err = osm.Client.Create(ctx, s)
		} else {
			osm.Log.Info("Updating configuration secret " + coh.OperatorConfigName)
			err = osm.Client.Update(ctx, s)
		}
	}

	return err
}
