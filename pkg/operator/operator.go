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
	v1 "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	"github.com/oracle/coherence-operator/pkg/flags"
	"io/ioutil"
	v12 "k8s.io/api/apps/v1"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"os"
	"path/filepath"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"

	corev1 "k8s.io/api/core/v1"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	crdclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

const (
	// configName is the name of the internal Coherence Operator configuration secret.
	configName = "coherence-operator-config"
	// The environment variable holding the default Coherence image name
	coherenceImageEnv = "HELM_COHERENCE_IMAGE"
	// The environment variable holding the default Coherence Utils image name
	utilsImageEnv = "UTILS_IMAGE"

	// The name of the Coherence container in the Coherence Pods
	CoherenceContainerName = "coherence"
	// The name of the Coherence Utils container in the Coherence Pods
	CoherenceUtilsContainerName = "coherence-k8s-utils"
)

var restHostAndPort string

func SetHostAndPort(hostAndPort string) {
	restHostAndPort = hostAndPort
}

// EnsureCRDs ensures that the Operator configuration secret exists in the namespace.
func EnsureCRDs(mgr manager.Manager, cohFlags *flags.CoherenceOperatorFlags, log logr.Logger) error {
	// Create the CRD client
	c, err := apiextensions.NewForConfig(mgr.GetConfig())
	if err != nil {
		return err
	}

	crdClient := c.ApiextensionsV1beta1().CustomResourceDefinitions()

	return EnsureCRDsUsingClient(mgr, cohFlags, log, crdClient)
}

// EnsureCRDs ensures that the Operator configuration secret exists in the namespace.
func EnsureCRDsUsingClient(mgr manager.Manager, cohFlags *flags.CoherenceOperatorFlags, log logr.Logger, crdClient crdclient.CustomResourceDefinitionInterface) error {
	log.Info("Ensuring operator CRDs are present")

	if cohFlags.CrdFiles == "" {
		log.Info("The CRD files location is blank - cannot ensure that CRDs are present in Kubernetes")
		return nil
	}

	_, err := os.Stat(cohFlags.CrdFiles)
	if err != nil {
		return fmt.Errorf("the CRD files location '%s' does not exist", cohFlags.CrdFiles)
	}

	log.Info("Loading operator CRDs from '" + cohFlags.CrdFiles + "'")
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

	if !mgr.GetScheme().IsGroupRegistered(v1beta1.GroupName) {
		err = v1beta1.AddToScheme(mgr.GetScheme())
		if err != nil {
			return err
		}
	}

	for _, file := range files {
		log.Info("Loading operator CRD yaml from '" + file + "'")
		// Load the CRD from the yaml file
		yml, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}
		newCRD := v1beta1.CustomResourceDefinition{}
		err = yaml.Unmarshal(yml, &newCRD)
		if err != nil {
			return err
		}

		// Get the existing CRD
		oldCRD, err := crdClient.Get(newCRD.Name, metav1.GetOptions{})
		switch {
		case err == nil:
			// CRD exists so update it
			log.Info("Updating operator CRD '" + newCRD.Name + "'")
			newCRD.ResourceVersion = oldCRD.ResourceVersion
			_, err = crdClient.Update(&newCRD)
			if err != nil {
				return err
			}
		case errors.IsNotFound(err):
			// CRD does not exist so create it
			log.Info("Creating operator CRD '" + newCRD.Name + "'")
			_, err = crdClient.Create(&newCRD)
			if err != nil {
				return err
			}
		default:
			// An error occurred
			return err
		}
	}

	return nil
}

// EnsureOperatorSecret ensures that the Operator configuration secret exists in the namespace.
func EnsureOperatorSecret(namespace string, c client.Client, log logr.Logger) error {
	log.Info("Ensuring configuration secret")

	err := c.Get(context.TODO(), types.NamespacedName{Name: configName, Namespace: namespace}, &corev1.Secret{})
	if err != nil && !errors.IsNotFound(err) {
		return err
	}

	log.Info("Operator Configuration: 'operatorhost' value set to " + restHostAndPort)

	secret := &corev1.Secret{}
	secret.SetNamespace(namespace)
	secret.SetName(configName)

	if secret.StringData == nil {
		secret.StringData = make(map[string]string)
	}

	secret.StringData["operatorhost"] = restHostAndPort

	if errors.IsNotFound(err) {
		// for some reason we're getting here even if the secret exists so delete it!!
		_ = c.Delete(context.TODO(), secret)
		log.Info("Creating secret " + configName + " in namespace " + namespace)
		err = c.Create(context.TODO(), secret)
	} else {
		log.Info("Updating secret " + configName + " in namespace " + namespace)
		err = c.Update(context.TODO(), secret)
	}

	return err
}

func GetDefaultCoherenceImage() *string {
	img, ok := os.LookupEnv(coherenceImageEnv)
	if ok {
		return &img
	}
	return nil
}

func GetDefaultCoherenceUtilsImages() *string {
	img, ok := os.LookupEnv(utilsImageEnv)
	if ok {
		return &img
	}
	return nil
}

// If the CoherenceInternalSpec does not have a Coherence image specified we set the default here.
// This ensures that the image is fixed to either that specified in the cluster spec or to the current default
// and means that the Helm controller does not upgrade the images if the Operator is upgraded.
func EnsureImages(ci *v1.CoherenceInternalSpec, logger logr.Logger) {
	coherenceImage := GetDefaultCoherenceImage()
	if ci.EnsureCoherenceImage(coherenceImage) {
		logger.Info(fmt.Sprintf("Injected Coherence image name into role: '%s'", *coherenceImage))
	}

	utilsImage := GetDefaultCoherenceUtilsImages()
	if ci.EnsureCoherenceUtilsImage(utilsImage) {
		logger.Info(fmt.Sprintf("Injected Coherence Utils image name into role: '%s'", *utilsImage))
	}
}

// Create the desired CoherenceInternalSpec for a given role.
func CreateDesiredRole(cluster *v1.CoherenceCluster, role *v1.CoherenceRole, existing *v1.CoherenceInternalSpec, sts *v12.StatefulSet) *v1.CoherenceInternalSpec {
	desiredRole := v1.NewCoherenceInternalSpec(cluster, role)

	coherenceImage := existing.GetCoherenceImage()

	if sts != nil && desiredRole.GetCoherenceImage() == nil {
		// if the desired Coherence image is still nil then this could be an update to a cluster
		// started with a much older Operator so we'll obtain the current image from the StatefulSet
		for _, c := range sts.Spec.Template.Spec.Containers {
			if c.Name == CoherenceContainerName {
				coherenceImage = &c.Image
			}
		}
	}

	if coherenceImage == nil {
		// If the Coherence image is still nil then use the default
		coherenceImage = GetDefaultCoherenceUtilsImages()
	}

	utilsImage := existing.GetCoherenceImage()

	if sts != nil && utilsImage == nil {
		// if the desired Coherence Utils image is still nil then this could be an update to a cluster
		// started with a much older Operator so we'll obtain the current image from the StatefulSet
		for _, c := range sts.Spec.Template.Spec.InitContainers {
			if c.Name == CoherenceUtilsContainerName {
				utilsImage = &c.Image
			}
		}
	}

	if utilsImage == nil {
		// If the utils image is still nil then use the default
		utilsImage = GetDefaultCoherenceUtilsImages()
	}

	desiredRole.EnsureCoherenceImage(coherenceImage)
	desiredRole.EnsureCoherenceUtilsImage(utilsImage)

	return desiredRole
}
