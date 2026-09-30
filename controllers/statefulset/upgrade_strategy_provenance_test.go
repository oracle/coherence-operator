/*
 * Copyright (c) 2026, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package statefulset

import (
	"testing"

	coh "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/pkg/probe"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/ptr"
)

func TestUpgradeStrategiesPreserveScalingProbeSource(t *testing.T) {
	custom := &coh.Probe{ProbeHandler: corev1.ProbeHandler{HTTPGet: &corev1.HTTPGetAction{
		Path:   "/ha",
		Port:   intstr.FromString(coh.PortNameHealth),
		Scheme: corev1.URISchemeHTTP,
	}}}

	for _, tc := range []struct {
		name      string
		strategy  coh.RollingUpdateStrategyType
		generated bool
		custom    *coh.Probe
	}{
		{name: "by-node generated", strategy: coh.UpgradeByNode, generated: true},
		{name: "by-node custom", strategy: coh.UpgradeByNode, custom: custom},
		{name: "by-node-label generated", strategy: coh.UpgradeByNodeLabel, generated: true},
		{name: "by-node-label custom", strategy: coh.UpgradeByNodeLabel, custom: custom},
	} {
		t.Run(tc.name, func(t *testing.T) {
			spec := coh.CoherenceStatefulSetResourceSpec{
				RollingUpdateStrategy: ptr.To(tc.strategy),
				RollingUpdateLabel:    ptr.To("topology.kubernetes.io/zone"),
			}
			if tc.custom != nil {
				spec.Scaling = &coh.ScalingSpec{Probe: tc.custom}
			}
			deployment := &coh.Coherence{Spec: spec}
			strategy := GetUpgradeStrategy(deployment, probe.CoherenceProbe{})

			var resolved probe.ResolvedScalingProbe
			switch actual := strategy.(type) {
			case ByNodeUpgradeStrategy:
				resolved = actual.scalingProbe
			case ByNodeLabelUpgradeStrategy:
				resolved = actual.scalingProbe
			default:
				t.Fatalf("unexpected strategy type %T", strategy)
			}

			if resolved.Generated != tc.generated {
				t.Fatalf("generated = %t, want %t", resolved.Generated, tc.generated)
			}
			if tc.custom != nil && resolved.Probe != tc.custom {
				t.Fatal("custom scaling probe identity was not preserved")
			}
			if resolved.Probe == nil {
				t.Fatal("resolved scaling probe is nil")
			}
		})
	}
}
