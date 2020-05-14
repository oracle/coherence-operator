/*
 * Copyright (c) 2019, 2020, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package helper

import (
	"errors"
	"github.com/ghodss/yaml"
	coh "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	"io/ioutil"
	coreV1 "k8s.io/api/core/v1"
	"os"
)

type OperatorValues struct {
	// The name to use for the service account to use when RBAC is enabled
	// The role bindings must already have been created as this chart does not create them it just
	// sets the serviceAccountName value in the Pod spec.
	// +optional
	ServiceAccountName string `json:"serviceAccountName,omitempty"`
	// The secrets to be used when pulling images. Secrets must be manually created in the target namespace.
	// ref: https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry/
	// +optional
	ImagePullSecrets []coh.LocalObjectReference `json:"imagePullSecrets,omitempty"`
	// Affinity controls Pod scheduling preferences.
	// ref: https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#affinity-and-anti-affinity
	// +optional
	Affinity *coreV1.Affinity `json:"affinity,omitempty"`
	// NodeSelector is the Node labels for pod assignment
	// ref: https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#nodeselector
	// +optional
	NodeSelector *coreV1.NodeSelector `json:"nodeSelector,omitempty"`
	// If specified, the pod's tolerations.
	// +optional
	Tolerations *[]coreV1.Toleration `json:"tolerations,omitempty"`
	// CoherenceOperator groups the values used to configure the operator
	// +optional
	CoherenceOperator *OperatorSpec `json:"coherenceOperator,omitempty"`
	// Specifies whether to generate the ClusterRole yaml.
	// +optional
	EnableClusterRole *bool `json:"enableClusterRole,omitempty"`
	// The Helm full name override
	// +optional
	FullnameOverride *string `json:"fullnameOverride,omitempty"`
}

// OperatorSpec defines the settings for the Operator.
type OperatorSpec struct {
	coh.ImageSpec `json:",inline"`
	SSL           *OperatorSSL `json:"ssl,omitempty"`
}

// OperatorSSL defines the SSL settings for the Operator.
type OperatorSSL struct {
	Secrets  *string `json:"secrets,omitempty"`
	KeyFile  *string `json:"keyFile,omitempty"`
	CertFile *string `json:"certFile,omitempty"`
	CaFile   *string `json:"caFile,omitempty"`
}

// Set whether to generate the ClusterRole yaml.
func (v *OperatorValues) SetEnableClusterRole(enabled bool) {
	if v != nil {
		v.EnableClusterRole = &enabled
	}
}

// LoadFromYaml loads the data from the specified YAML file into this OperatorValues
func (v *OperatorValues) LoadFromYaml(file string) error {
	if v == nil {
		return errors.New("attempted to load yaml into a nil OperatorValues reference")
	}
	_, err := os.Stat(file)
	if err != nil {
		return err
	}

	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, v)
}

// ToYaml marshals this OperatorValues to yaml
func (v *OperatorValues) ToYaml() ([]byte, error) {
	if v == nil {
		return nil, errors.New("attempted to marshall nil OperatorValues to yaml")
	}

	return yaml.Marshal(v)
}

// ToYaml marshals this OperatorValues to yaml
func (v *OperatorValues) ToMap(m *map[string]interface{}) error {
	if v == nil {
		return errors.New("attempted to convert nil OperatorValues to a map")
	}

	d, err := v.ToYaml()
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(d, m)
	return err
}
