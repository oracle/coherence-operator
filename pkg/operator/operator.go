/*
 * Copyright (c) 2019, 2020 Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

// The operator package contains types and functions used directly by the Operator main
package operator

import (
	"context"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/go-logr/logr"
	coh "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/pkg/data"
	"github.com/oracle/coherence-operator/pkg/rest"
	"github.com/pkg/errors"
	"io/ioutil"
	corev1 "k8s.io/api/core/v1"
	crdv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	crdbeta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	v1client "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1"
	v1beta1client "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1beta1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/version"
	"k8s.io/client-go/discovery"
	rest2 "k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// EnsureCRDs ensures that the Operator configuration secret exists in the namespace.
// CRDs will be created depending on the server version of k8s. For k8s v1.16.0 and above
// the v1 CRDs will be created and for lower than v1.16.0 the v1beta1 CRDs will be created.
func EnsureCRDs(cfg *rest2.Config) error {
	// Create the CRD client
	c, err := apiextensions.NewForConfig(cfg)
	if err != nil {
		return err
	}

	cl, err := discovery.NewDiscoveryClientForConfig(cfg)
	if err != nil {
		return err
	}
	sv, err := cl.ServerVersion()
	if err != nil {
		return err
	}
	v, err := version.ParseSemantic(sv.GitVersion)
	if err != nil {
		return err
	}

	logger := logf.Log.WithName("operator")

	if v.Major() > 1 || (v.Major() == 1 && v.Minor() >= 16) {
		// k8s v1.16.0 or above - install v1 CRD
		crdClient := c.ApiextensionsV1().CustomResourceDefinitions()
		return EnsureV1CRDs(logger, crdClient)
	}
	// k8s lower than v1.16.0 - install v1beta1 CRD
	crdClient := c.ApiextensionsV1beta1().CustomResourceDefinitions()
	return EnsureV1Beta1CRDs(logger, crdClient)
}

// EnsureCRDs ensures that the Operator configuration secret exists in the namespace.
func EnsureV1CRDs(logger logr.Logger, crdClient v1client.CustomResourceDefinitionInterface) error {
	logger.Info("Ensuring operator v1 CRDs are present")
	return ensureV1CRDs(logger, crdClient, "crd_v1.yaml")
}

// EnsureCRD ensures that the specified V1 CRDs are loaded using the specified embedded CRD files
func ensureV1CRDs(logger logr.Logger, crdClient v1client.CustomResourceDefinitionInterface, fileNames ...string) error {
	for _, fileName := range fileNames {
		if err := ensureV1CRD(logger, crdClient, fileName); err != nil {
			return err
		}
	}
	return nil
}

// EnsureCRD ensures that the specified V1 CRD is loaded using the specified embedded CRD file
func ensureV1CRD(logger logr.Logger, crdClient v1client.CustomResourceDefinitionInterface, fileName string) error {
	logger.Info("Ensuring operator v1 CRDs are present")

	f, err := data.Assets.Open(fileName)
	if err != nil {
		return errors.Wrap(err, "opening embedded CRD asset " + fileName)
	}
	defer f.Close()

	yml, err := ioutil.ReadAll(f);
	if err != nil {
		return errors.Wrap(err, "reading embedded CRD asset " + fileName)
	}

	u := unstructured.Unstructured{}
	err = yaml.Unmarshal(yml, &u)
	if err != nil {
		return err
	}

	newCRD := crdv1.CustomResourceDefinition{}
	err = yaml.Unmarshal(yml, &newCRD)
	if err != nil {
		return err
	}

	logger.Info("Loading operator CRD yaml from '" + fileName + "'")

	// Get the existing CRD
	oldCRD, err := crdClient.Get(context.TODO(), newCRD.Name, metav1.GetOptions{})
	switch {
	case err == nil:
		// CRD exists so update it
		logger.Info("Updating operator CRD '" + newCRD.Name + "'")
		newCRD.ResourceVersion = oldCRD.ResourceVersion
		_, err = crdClient.Update(context.TODO(), &newCRD, metav1.UpdateOptions{})
		if err != nil {
			return errors.Wrapf(err, "updating Coherence CRD %s", newCRD.Name)
		}
	case apierrors.IsNotFound(err):
		// CRD does not exist so create it
		logger.Info("Creating operator CRD '" + newCRD.Name + "'")
		_, err = crdClient.Create(context.TODO(), &newCRD, metav1.CreateOptions{})
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

// EnsureCRDs ensures that the Operator configuration secret exists in the namespace.
func EnsureV1Beta1CRDs(logger logr.Logger, crdClient v1beta1client.CustomResourceDefinitionInterface) error {
	logger.Info("Ensuring operator v1beta1 CRDs are present")
	return ensureV1Beta1CRDs(logger, crdClient, "crd_v1beta1.yaml")
}

// EnsureCRD ensures that the specified V1 CRDs are loaded using the specified embedded CRD files
func ensureV1Beta1CRDs(logger logr.Logger, crdClient v1beta1client.CustomResourceDefinitionInterface, fileNames ...string) error {
	for _, fileName := range fileNames {
		if err := ensureV1Beta1CRD(logger, crdClient, fileName); err != nil {
			return err
		}
	}
	return nil
}

// EnsureCRD ensures that the specified V1 CRD is loaded using the specified embedded CRD file
func ensureV1Beta1CRD(logger logr.Logger, crdClient v1beta1client.CustomResourceDefinitionInterface, fileName string) error {
	logger.Info("Ensuring operator v1 CRDs are present")

	f, err := data.Assets.Open(fileName)
	if err != nil {
		return errors.Wrap(err, "opening embedded CRD asset " + fileName)
	}
	defer f.Close()

	yml, err := ioutil.ReadAll(f);
	if err != nil {
		return errors.Wrap(err, "reading embedded CRD asset " + fileName)
	}

	u := unstructured.Unstructured{}
	err = yaml.Unmarshal(yml, &u)
	if err != nil {
		return err
	}

	newCRD := crdbeta1.CustomResourceDefinition{}
	err = yaml.Unmarshal(yml, &newCRD)
	if err != nil {
		return err
	}

	logger.Info("Loading operator CRD yaml from '" + fileName + "'")

	// Get the existing CRD
	oldCRD, err := crdClient.Get(context.TODO(), newCRD.Name, metav1.GetOptions{})
	switch {
	case err == nil:
		// CRD exists so update it
		logger.Info("Updating operator CRD '" + newCRD.Name + "'")
		newCRD.ResourceVersion = oldCRD.ResourceVersion
		_, err = crdClient.Update(context.TODO(), &newCRD, metav1.UpdateOptions{})
		if err != nil {
			return errors.Wrapf(err, "updating Coherence CRD %s", newCRD.Name)
		}
	case apierrors.IsNotFound(err):
		// CRD does not exist so create it
		logger.Info("Creating operator CRD '" + newCRD.Name + "'")
		_, err = crdClient.Create(context.TODO(), &newCRD, metav1.CreateOptions{})
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

// EnsureOperatorSecret ensures that the Operator configuration secret exists in the namespace.
func EnsureOperatorSecret(namespace string, c client.Client, log logr.Logger) error {
	log.Info("Ensuring configuration secret")

	secret := &corev1.Secret{}

	err := c.Get(context.TODO(), types.NamespacedName{Name: coh.OperatorConfigName, Namespace: namespace}, secret)
	if err != nil && !apierrors.IsNotFound(err) {
		return err
	}

	restHostAndPort := rest.GetServerHostAndPort()

	log.Info(fmt.Sprintf("Operator Configuration: '%s' value set to %s", coh.OperatorConfigKeyHost, restHostAndPort))

	secret.SetNamespace(namespace)
	secret.SetName(coh.OperatorConfigName)

	if secret.StringData == nil {
		secret.StringData = make(map[string]string)
	}

	secret.StringData[coh.OperatorConfigKeyHost] = restHostAndPort

	if apierrors.IsNotFound(err) {
		// for some reason we're getting here even if the secret exists so delete it!!
		_ = c.Delete(context.TODO(), secret)
		log.Info("Creating secret " + coh.OperatorConfigName + " in namespace " + namespace)
		err = c.Create(context.TODO(), secret)
	} else {
		log.Info("Updating secret " + coh.OperatorConfigName + " in namespace " + namespace)
		err = c.Update(context.TODO(), secret)
	}

	return err
}
