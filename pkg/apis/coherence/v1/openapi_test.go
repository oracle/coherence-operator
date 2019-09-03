/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1_test

import (
	"github.com/go-openapi/validate"
	v1 "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	"io/ioutil"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	"k8s.io/apiextensions-apiserver/pkg/apiserver/validation"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
	"strings"
	"testing"

	. "github.com/onsi/gomega"
)

/*
The tests in this file validate the Open API specifications in the generated CRDs have not been broken.
The main point of the test is to validate that the required fields of a CRD have not been changed.
In our CRDs all fields are currently optional so these tests make sure that the minimal structure for
each CRD passes the Open API validator.
*/

// The base location for CRD files - this is relative to this test file's location.
const crdBase = "../../../../deploy/crds/"

func TestCoherenceClusterOpenApiSpec(t *testing.T) {
	g := NewGomegaWithT(t)

	v := createValidator(t, crdBase+"coherence_v1_coherencecluster_crd.yaml")

	// This is the minimal valid spec for a CoherenceCluster.
	// This structure should be valid against the CRD spec
	spec := v1.CoherenceCluster{
		ObjectMeta: metav1.ObjectMeta{Namespace: "test-ns", Name: "test-cluster"},
	}

	result := v.Validate(spec)
	g.Expect(result.IsValid()).To(BeTrue(), resultToString(result))
}

func TestCoherenceRoleOpenApiSpec(t *testing.T) {
	g := NewGomegaWithT(t)

	v := createValidator(t, crdBase+"coherence_v1_coherencerole_crd.yaml")

	// This is the minimal valid spec for a CoherenceRole.
	// This structure should be valid against the CRD spec
	spec := v1.CoherenceRole{
		ObjectMeta: metav1.ObjectMeta{Namespace: "test-ns", Name: "test-role"},
	}

	result := v.Validate(spec)
	g.Expect(result.IsValid()).To(BeTrue(), resultToString(result))
}

func TestCoherenceInternalOpenApiSpec(t *testing.T) {
	g := NewGomegaWithT(t)

	v := createValidator(t, crdBase+"coherence_v1_coherenceinternal_crd.yaml")

	// This is the minimal valid spec for a CoherenceInternal.
	// This structure should be valid against the CRD spec
	spec := v1.CoherenceInternal{
		ObjectMeta: metav1.ObjectMeta{Namespace: "test-ns", Name: "test-role"},
	}

	result := v.Validate(spec)
	g.Expect(result.IsValid()).To(BeTrue(), resultToString(result))
}

// ----- helper methods -----------------------------------------------------

// Create an Open API spec validator for a give CRD.
func createValidator(t *testing.T, crdPath string) *validate.SchemaValidator {
	g := NewGomegaWithT(t)

	yamlFile, err := ioutil.ReadFile(crdPath)
	g.Expect(err).ToNot(HaveOccurred())

	crd := apiextensions.CustomResourceDefinition{}
	err = yaml.Unmarshal(yamlFile, &crd)
	g.Expect(err).ToNot(HaveOccurred())

	v, _, err := validation.NewSchemaValidator(crd.Spec.Validation)
	g.Expect(err).ToNot(HaveOccurred())

	return v
}

// Convert a result to a string to use for test failure descriptions.
func resultToString(r *validate.Result) string {
	if r.IsValid() {
		return "CRD is valid"
	}

	b := strings.Builder{}

	b.WriteString("Expected CRD to be valid.\nErrors:")

	for _, e := range r.Errors {
		b.WriteString("\n    ")
		b.WriteString(e.Error())
	}

	b.WriteString("Warnings:\n")

	for _, w := range r.Warnings {
		b.WriteString("\n    ")
		b.WriteString(w.Error())
	}

	return b.String()
}
