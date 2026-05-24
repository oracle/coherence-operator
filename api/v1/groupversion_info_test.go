/*
 * Copyright (c) 2026, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1_test

import (
	"testing"

	. "github.com/onsi/gomega"
	coh "github.com/oracle/coherence-operator/api/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestAddToSchemeRegistersCoherenceAPITypes(t *testing.T) {
	g := NewGomegaWithT(t)

	// Use an empty scheme so the test proves AddToScheme owns the API type registration
	// contract instead of inheriting it from broader test setup.
	scheme := runtime.NewScheme()
	g.Expect(coh.AddToScheme(scheme)).To(Succeed())

	tests := []struct {
		name     string
		kind     string
		expected runtime.Object
	}{
		{
			name:     "coherence",
			kind:     "Coherence",
			expected: &coh.Coherence{},
		},
		{
			name:     "coherence list",
			kind:     "CoherenceList",
			expected: &coh.CoherenceList{},
		},
		{
			name:     "coherence job",
			kind:     "CoherenceJob",
			expected: &coh.CoherenceJob{},
		},
		{
			name:     "coherence job list",
			kind:     "CoherenceJobList",
			expected: &coh.CoherenceJobList{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGomegaWithT(t)

			obj, err := scheme.New(coh.GroupVersion.WithKind(tt.kind))

			g.Expect(err).NotTo(HaveOccurred())
			g.Expect(obj).To(BeAssignableToTypeOf(tt.expected))
		})
	}
}
