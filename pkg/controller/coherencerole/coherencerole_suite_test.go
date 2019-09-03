/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package coherencerole

import (
	"encoding/json"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	coherence "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"

	"testing"
)

func TestCoherenceRoleControler(t *testing.T) {
	RegisterFailHandler(Fail)
	junitReporter := reporters.NewJUnitReporter("test-report.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "CoherenceRole Controller Suite", []Reporter{junitReporter})
}

// UnstructuredToCoherenceInternalSpec converts an Unstructured to a CoherenceInternalSpec.
// The "spec" field of the Unstructured is marshalled to json and then
// un-marshalled back to a CoherenceInternalSpec.
func UnstructuredToCoherenceInternalSpec(u *unstructured.Unstructured) *coherence.CoherenceInternalSpec {
	spec := &coherence.CoherenceInternalSpec{}
	data, err := json.Marshal(u.Object["spec"])
	Expect(err).ToNot(HaveOccurred())
	err = json.Unmarshal(data, spec)
	Expect(err).ToNot(HaveOccurred())
	return spec
}
