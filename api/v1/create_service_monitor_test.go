/*
 * Copyright (c) 2026, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1_test

import (
	"encoding/json"
	"math"
	"testing"

	. "github.com/onsi/gomega"
	coh "github.com/oracle/coherence-operator/api/v1"
)

func TestServiceMonitorSpecCreateServiceMonitorWithNilSampleLimit(t *testing.T) {
	g := NewGomegaWithT(t)

	spec := (&coh.ServiceMonitorSpec{}).CreateServiceMonitor()

	g.Expect(spec.SampleLimit).To(BeNil())
}

func TestServiceMonitorSpecCreateServiceMonitorConvertsSampleLimit(t *testing.T) {
	tests := []struct {
		name  string
		value uint64
	}{
		{name: "zero", value: 0},
		{name: "positive", value: 12345},
		{name: "maximum CRD value", value: math.MaxInt64},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGomegaWithT(t)
			value := tt.value

			spec := (&coh.ServiceMonitorSpec{SampleLimit: &value}).CreateServiceMonitor()

			g.Expect(spec.SampleLimit).NotTo(BeNil())
			g.Expect(*spec.SampleLimit).To(Equal(int64(tt.value)))
		})
	}
}

func TestServiceMonitorSpecSampleLimitHasLegacyWireFormat(t *testing.T) {
	g := NewGomegaWithT(t)
	value := uint64(12345)
	spec := (&coh.ServiceMonitorSpec{SampleLimit: &value}).CreateServiceMonitor()

	data, err := json.Marshal(spec)
	g.Expect(err).NotTo(HaveOccurred())

	legacy := struct {
		SampleLimit *uint64 `json:"sampleLimit,omitempty"`
	}{}
	err = json.Unmarshal(data, &legacy)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(legacy.SampleLimit).NotTo(BeNil())
	g.Expect(*legacy.SampleLimit).To(Equal(value))
}
