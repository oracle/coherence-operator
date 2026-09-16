/*
 * Copyright (c) 2026, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1_test

import (
	"encoding/json"
	"fmt"
	monitoring "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
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

func TestNamedPortSpecCreateServiceMonitorUseDefaultLabelDrop(t *testing.T) {
	broad := monitoring.RelabelConfig{Action: "labeldrop", Regex: "(endpoint|instance|job|service)"}
	custom := []monitoring.RelabelConfig{{Action: "replace", SourceLabels: []monitoring.LabelName{"pod"}, TargetLabel: "origin"}, {Action: "labeldrop", Regex: "(endpoint|service)"}, broad}
	for _, mode := range []struct {
		name  string
		value *bool
	}{{"omitted", nil}, {"true", boolPtr(true)}, {"false", boolPtr(false)}} {
		for _, rules := range [][]monitoring.RelabelConfig{nil, {}, custom} {
			t.Run(fmt.Sprintf("%s/%d/%t", mode.name, len(rules), rules == nil), func(t *testing.T) {
				g := NewGomegaWithT(t)
				spec := &coh.ServiceMonitorSpec{Enabled: boolPtr(true), UseDefaultLabelDrop: mode.value, Relabelings: rules, MetricRelabelings: custom, HonorLabels: true, Interval: "15s", Path: "/metrics", ScrapeTimeout: "10s"}
				port := coh.NamedPortSpec{Name: "metrics", ServiceMonitor: spec}
				expected := spec.CreateEndpoint()
				expected.Port = "metrics"
				expected.RelabelConfigs = append([]monitoring.RelabelConfig(nil), rules...)
				if rules != nil && len(rules) == 0 {
					expected.RelabelConfigs = []monitoring.RelabelConfig{}
				}
				if mode.value == nil || *mode.value {
					expected.RelabelConfigs = append(expected.RelabelConfigs, broad)
				}
				for _, parent := range []coh.CoherenceResource{&coh.Coherence{}, &coh.CoherenceJob{}} {
					sm := port.CreateServiceMonitor(parent)
					g.Expect(sm).NotTo(BeNil())
					g.Expect(sm.Spec.Endpoints).To(Equal([]monitoring.Endpoint{expected}))
				}
			})
		}
	}
}

func TestNamedPortSpecCreateServiceMonitorRelabelConfigsAreIndependent(t *testing.T) {
	for _, flag := range []*bool{nil, boolPtr(true), boolPtr(false)} {
		g := NewGomegaWithT(t)
		rules := make([]monitoring.RelabelConfig, 1, 8)
		rules[0] = monitoring.RelabelConfig{Action: "labeldrop", Regex: "custom"}
		backing := rules[:cap(rules)]
		port := coh.NamedPortSpec{Name: "metrics", ServiceMonitor: &coh.ServiceMonitorSpec{Enabled: boolPtr(true), UseDefaultLabelDrop: flag, Relabelings: rules}}
		first := port.CreateServiceMonitor(&coh.Coherence{})
		second := port.CreateServiceMonitor(&coh.Coherence{})
		g.Expect(second).To(Equal(first))
		first.Spec.Endpoints[0].RelabelConfigs[0].Regex = "changed"
		g.Expect(rules[0].Regex).To(Equal("custom"))
		g.Expect(second.Spec.Endpoints[0].RelabelConfigs[0].Regex).To(Equal("custom"))
		g.Expect(backing[1:]).To(Equal(make([]monitoring.RelabelConfig, 7)))
	}
}

func TestNamedPortSpecCreateServiceMonitorDisabled(t *testing.T) {
	for _, port := range []*coh.NamedPortSpec{nil, {}, {ServiceMonitor: &coh.ServiceMonitorSpec{UseDefaultLabelDrop: boolPtr(false)}}, {ServiceMonitor: &coh.ServiceMonitorSpec{Enabled: boolPtr(false)}}, {Service: &coh.ServiceSpec{Enabled: boolPtr(false)}, ServiceMonitor: &coh.ServiceMonitorSpec{Enabled: boolPtr(true), UseDefaultLabelDrop: boolPtr(false)}}} {
		NewGomegaWithT(t).Expect(port.CreateServiceMonitor(&coh.Coherence{})).To(BeNil())
	}
}

func TestServiceMonitorSpecUseDefaultLabelDropJSONAndDeepCopy(t *testing.T) {
	for _, flag := range []*bool{nil, boolPtr(true), boolPtr(false)} {
		g := NewGomegaWithT(t)
		original := &coh.ServiceMonitorSpec{UseDefaultLabelDrop: flag}
		data, err := json.Marshal(original)
		g.Expect(err).NotTo(HaveOccurred())
		var decoded coh.ServiceMonitorSpec
		g.Expect(json.Unmarshal(data, &decoded)).To(Succeed())
		g.Expect(&decoded).To(Equal(original))
		if flag == nil {
			g.Expect(string(data)).To(Equal("{}"))
		} else {
			g.Expect(string(data)).To(ContainSubstring(fmt.Sprintf(`"useDefaultLabelDrop":%t`, *flag)))
		}
		copied := original.DeepCopy()
		g.Expect(copied).To(Equal(original))
		if flag != nil {
			*copied.UseDefaultLabelDrop = !*flag
			g.Expect(*original.UseDefaultLabelDrop).NotTo(Equal(*copied.UseDefaultLabelDrop))
		}
	}
}

func TestServiceMonitorPodTemplatesAndServicesUnchanged(t *testing.T) {
	g := NewGomegaWithT(t)
	ports := []coh.NamedPortSpec{{Name: "metrics", ServiceMonitor: &coh.ServiceMonitorSpec{Enabled: boolPtr(true)}}, {Name: "other", ServiceMonitor: &coh.ServiceMonitorSpec{Enabled: boolPtr(true)}}}
	deployment := createTestCoherenceDeployment(coh.CoherenceStatefulSetResourceSpec{CoherenceResourceSpec: coh.CoherenceResourceSpec{Ports: ports}})
	job := createTestCoherenceJob(coh.CoherenceResourceSpec{Ports: ports}).DeepCopy()
	stsBefore := deployment.Spec.CreateStatefulSet(deployment)
	jobBefore := job.Spec.CreateJob(job)
	serviceBefore := deployment.Spec.Ports[0].CreateService(deployment)
	deployment.Spec.Ports[0].ServiceMonitor.UseDefaultLabelDrop = boolPtr(false)
	job.Spec.Ports[0].ServiceMonitor.UseDefaultLabelDrop = boolPtr(false)
	stsAfter := deployment.Spec.CreateStatefulSet(deployment)
	jobAfter := job.Spec.CreateJob(job)
	g.Expect(stsAfter.Spec.Template).To(Equal(stsBefore.Spec.Template))
	g.Expect(jobAfter.Spec.Template).To(Equal(jobBefore.Spec.Template))
	g.Expect(deployment.Spec.Ports[0].CreateService(deployment)).To(Equal(serviceBefore))
	for _, parent := range []coh.CoherenceResource{deployment, job} {
		g.Expect(parent.GetSpec().Ports[0].CreateServiceMonitor(parent).Spec.Endpoints[0].RelabelConfigs).To(BeEmpty())
		g.Expect(parent.GetSpec().Ports[1].CreateServiceMonitor(parent).Spec.Endpoints[0].RelabelConfigs).To(HaveLen(1))
	}
}
