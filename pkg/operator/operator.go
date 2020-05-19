/*
 * Copyright (c) 2019, 2020 Oracle and/or its affiliates. All rights reserved.
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
	coh "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	"github.com/oracle/coherence-operator/pkg/flags"
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
	"os"
	"path/filepath"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"strings"
)

// EnsureCRDs ensures that the Operator configuration secret exists in the namespace.
// CRDs will be created depending on the server version of k8s. For k8s v1.16.0 and above
// the v1 CRDs will be created and for lower than v1.16.0 the v1beta1 CRDs will be created.
func EnsureCRDs(mgr manager.Manager) error {
	// Create the CRD client
	c, err := apiextensions.NewForConfig(mgr.GetConfig())
	if err != nil {
		return err
	}

	cl, err := discovery.NewDiscoveryClientForConfig(mgr.GetConfig())
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
	cohFlags := flags.GetFlags()

	if v.Major() > 1 || (v.Major() == 1 && v.Minor() >= 16) {
		// k8s v1.16.0 or above - install v1 CRD
		crdClient := c.ApiextensionsV1().CustomResourceDefinitions()
		return EnsureV1CRDs(mgr, cohFlags, logger, crdClient)
	}
	// k8s lower than v1.16.0 - install v1beta1 CRD
	crdClient := c.ApiextensionsV1beta1().CustomResourceDefinitions()
	return EnsureV1Beta1CRDs(mgr, cohFlags, logger, crdClient)
}

// EnsureCRDs ensures that the Operator configuration secret exists in the namespace.
func EnsureV1CRDs(mgr manager.Manager, cohFlags *flags.CoherenceOperatorFlags, logger logr.Logger, crdClient v1client.CustomResourceDefinitionInterface) error {
	logger.Info("Ensuring operator v1 CRDs are present")

	if cohFlags.CrdFiles == "" {
		logger.Info("The CRD files location is blank - cannot ensure that CRDs are present in Kubernetes")
		return nil
	}

	_, err := os.Stat(cohFlags.CrdFiles)
	if err != nil {
		return fmt.Errorf("the CRD files location '%s' does not exist", cohFlags.CrdFiles)
	}

	logger.Info("Loading operator CRDs from '" + cohFlags.CrdFiles + "'")
	var files []string
	err = filepath.Walk(cohFlags.CrdFiles, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, "crd.yaml") && !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		return err
	}

	if !mgr.GetScheme().IsGroupRegistered(crdv1.GroupName) {
		err = crdv1.AddToScheme(mgr.GetScheme())
		if err != nil {
			return err
		}
	}

	for _, file := range files {
		// Load the CRD from the yaml file
		yml, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}

		u := unstructured.Unstructured{}
		err = yaml.Unmarshal(yml, &u)
		if err != nil {
			return err
		}

		if u.GetAPIVersion() != crdbeta1.GroupName+"/v1" {
			continue
		}

		newCRD := crdv1.CustomResourceDefinition{}
		err = yaml.Unmarshal(yml, &newCRD)
		if err != nil {
			return err
		}

		// make sure we're only loading v1 files
		if newCRD.APIVersion != crdbeta1.GroupName+"/v1" {
			continue
		}
		logger.Info("Loading operator CRD yaml from '" + file + "'")

		// Get the existing CRD
		oldCRD, err := crdClient.Get(newCRD.Name, metav1.GetOptions{})
		switch {
		case err == nil:
			// CRD exists so update it
			logger.Info("Updating operator CRD '" + newCRD.Name + "'")
			newCRD.ResourceVersion = oldCRD.ResourceVersion
			err = mgr.GetClient().Update(context.TODO(), &newCRD, &client.UpdateOptions{})
			if err != nil {
				return err
			}
		case apierrors.IsNotFound(err):
			// CRD does not exist so create it
			logger.Info("Creating operator CRD '" + newCRD.Name + "'")
			err = mgr.GetClient().Create(context.TODO(), &newCRD, &client.CreateOptions{})
			if err != nil {
				return err
			}
		default:
			// An error occurred
			logger.Error(err, "checking for existing Coherence CRD "+newCRD.Name)
			return errors.Wrapf(err, "checking for existing Coherence CRD %s", newCRD.Name)
		}
	}

	return nil
}

// EnsureCRDs ensures that the Operator configuration secret exists in the namespace.
func EnsureV1Beta1CRDs(mgr manager.Manager, cohFlags *flags.CoherenceOperatorFlags, logger logr.Logger, crdClient v1beta1client.CustomResourceDefinitionInterface) error {
	logger.Info("Ensuring operator v1beta1 CRDs are present")

	if cohFlags.CrdFiles == "" {
		logger.Info("The CRD files location is blank - cannot ensure that CRDs are present in Kubernetes")
		return nil
	}

	_, err := os.Stat(cohFlags.CrdFiles)
	if err != nil {
		return fmt.Errorf("the CRD files location '%s' does not exist", cohFlags.CrdFiles)
	}

	logger.Info("Loading operator CRDs from '" + cohFlags.CrdFiles + "'")
	var files []string
	err = filepath.Walk(cohFlags.CrdFiles, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, "crd.yaml") && !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		return err
	}

	if !mgr.GetScheme().IsGroupRegistered(crdbeta1.GroupName) {
		err = crdbeta1.AddToScheme(mgr.GetScheme())
		if err != nil {
			return err
		}
	}

	for _, file := range files {
		// Load the CRD from the yaml file
		yml, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}

		u := unstructured.Unstructured{}
		err = yaml.Unmarshal(yml, &u)
		if err != nil {
			return err
		}

		if u.GetAPIVersion() != crdbeta1.GroupName+"/v1beta1" {
			continue
		}

		newCRD := crdbeta1.CustomResourceDefinition{}
		err = yaml.Unmarshal(yml, &newCRD)
		if err != nil {
			return err
		}

		// make sure we're only loading v1beta1 files
		if newCRD.APIVersion != crdbeta1.GroupName+"/v1beta1" {
			continue
		}
		logger.Info("Loading operator CRD yaml from '" + file + "'")

		// Get the existing CRD
		oldCRD, err := crdClient.Get(newCRD.Name, metav1.GetOptions{})
		switch {
		case err == nil:
			// CRD exists so update it
			logger.Info("Updating operator CRD '" + newCRD.Name + "'")
			newCRD.ResourceVersion = oldCRD.ResourceVersion
			err = mgr.GetClient().Update(context.TODO(), &newCRD, &client.UpdateOptions{})
			if err != nil {
				return errors.Wrapf(err, "updating Coherence CRD %s", newCRD.Name)
			}
		case apierrors.IsNotFound(err):
			// CRD does not exist so create it
			logger.Info("Creating operator CRD '" + newCRD.Name + "'")
			err = mgr.GetClient().Create(context.TODO(), &newCRD, &client.CreateOptions{})
			if err != nil {
				return errors.Wrapf(err, "creating Coherence CRD %s", newCRD.Name)
			}
		default:
			// An error occurred
			logger.Error(err, "checking for existing Coherence CRD "+newCRD.Name)
			return errors.Wrapf(err, "checking for existing Coherence CRD %s", newCRD.Name)
		}
	}

	return nil
}

// EnsureOperatorSecret ensures that the Operator configuration secret exists in the namespace.
func EnsureOperatorSecret(namespace string, c client.Client, log logr.Logger) error {
	log.Info("Ensuring configuration secret")

	err := c.Get(context.TODO(), types.NamespacedName{Name: coh.OperatorConfigName, Namespace: namespace}, &corev1.Secret{})
	if err != nil && !apierrors.IsNotFound(err) {
		return err
	}

	restHostAndPort := rest.GetServerHostAndPort()

	log.Info(fmt.Sprintf("Operator Configuration: '%s' value set to %s", coh.OperatorConfigKeyHost, restHostAndPort))

	secret := &corev1.Secret{}
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
