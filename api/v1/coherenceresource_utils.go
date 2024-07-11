/*
 * Copyright (c) 2020, 2024, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1

import (
	"context"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/go-logr/logr"
	"github.com/oracle/coherence-operator/pkg/data"
	"github.com/oracle/coherence-operator/pkg/operator"
	"github.com/pkg/errors"
	"io"
	crdv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/version"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// Utility and helper functions for the Coherence API

// EnsureCRDs ensures that the Operator configuration secret exists in the namespace.
// CRDs will be created depending on the server version of k8s.
func EnsureCRDs(ctx context.Context, v *version.Version, scheme *runtime.Scheme, cl client.Client) error {
	logger := logf.Log.WithName("operator")
	logger.Info(fmt.Sprintf("Ensuring operator CRDs are present (K8s version %v)", v))
	return EnsureV1CRDs(ctx, logger, scheme, cl)
}

// EnsureV1CRDs ensures that the Operator configuration secret exists in the namespace.
func EnsureV1CRDs(ctx context.Context, logger logr.Logger, scheme *runtime.Scheme, cl client.Client) error {
	err := ensureV1CRDs(ctx, logger, scheme, cl, "apiextensions.k8s.io_v1_customresourcedefinition_coherence.coherence.oracle.com.yaml")
	if err != nil {
		return err
	}
	if operator.ShouldInstallJobCRD() {
		return ensureV1CRDs(ctx, logger, scheme, cl, "apiextensions.k8s.io_v1_customresourcedefinition_coherencejob.coherence.oracle.com.yaml")
	}
	return nil
}

// ensureV1CRDs ensures that the specified V1 CRDs are loaded using the specified embedded CRD files
func ensureV1CRDs(ctx context.Context, logger logr.Logger, scheme *runtime.Scheme, cl client.Client, fileNames ...string) error {
	if err := crdv1.AddToScheme(scheme); err != nil {
		return err
	}
	for _, fileName := range fileNames {
		if err := ensureV1CRD(ctx, logger, cl, fileName); err != nil {
			return err
		}
	}
	return nil
}

// ensureV1CRD ensures that the specified V1 CRD is loaded using the specified embedded CRD file
func ensureV1CRD(ctx context.Context, logger logr.Logger, cl client.Client, fileName string) error {
	f, err := data.Assets.Open("assets/" + fileName)
	if err != nil {
		return errors.Wrap(err, "opening embedded CRD asset "+fileName)
	}
	//goland:noinspection GoUnhandledErrorResult
	defer f.Close()

	yml, err := io.ReadAll(f)
	if err != nil {
		return errors.Wrap(err, "reading embedded CRD asset "+fileName)
	}

	u := unstructured.Unstructured{}
	err = yaml.Unmarshal(yml, &u)
	if err != nil {
		return err
	}

	oldCRD := crdv1.CustomResourceDefinition{}
	newCRD := crdv1.CustomResourceDefinition{}
	err = yaml.Unmarshal(yml, &newCRD)
	if err != nil {
		return err
	}

	logger.Info("Loading operator CRD yaml from '" + fileName + "'")

	// Get the existing CRD
	err = cl.Get(ctx, client.ObjectKey{Name: newCRD.Name}, &oldCRD)
	switch {
	case err == nil:
		// CRD exists so update it
		logger.Info("Updating operator CRD '" + newCRD.Name + "'")
		newCRD.ResourceVersion = oldCRD.ResourceVersion
		err = cl.Update(ctx, &newCRD)
		if err != nil {
			return errors.Wrapf(err, "updating Coherence CRD %s", newCRD.Name)
		}
	case apierrors.IsNotFound(err):
		// CRD does not exist so create it
		logger.Info("Creating operator CRD '" + newCRD.Name + "'")
		err = cl.Create(ctx, &newCRD)
		if err != nil {
			return errors.Wrapf(err, "creating Coherence CRD %s", newCRD.Name)
		}
	default:
		// An error occurred
		logger.Error(err, "checking for existing Coherence CRD "+newCRD.Name)
		return errors.Wrapf(err, "checking for existing Coherence CRD %s", newCRD.Name)
	}

	return nil
}
