/*
 * Copyright (c) 2020, 2024, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package statefulset_test

import (
	. "github.com/onsi/gomega"
	coh "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/controllers/statefulset"
	"github.com/oracle/coherence-operator/pkg/probe"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	"testing"
)

func TestUseUpgradeStrategyByPodIfNotSet(t *testing.T) {
	g := NewGomegaWithT(t)

	c := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-deployment",
		},
		Spec: coh.CoherenceStatefulSetResourceSpec{},
	}

	p := probe.CoherenceProbe{}
	s := statefulset.GetUpgradeStrategy(c, p)

	g.Expect(s).To(BeAssignableToTypeOf(statefulset.ByPodUpgradeStrategy{}))
	g.Expect(s.IsOperatorManaged()).To(BeFalse())
}

func TestUseUpgradeStrategyByPod(t *testing.T) {
	g := NewGomegaWithT(t)

	c := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-deployment",
		},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			RollingUpdateStrategy: ptr.To(coh.UpgradeByPod),
		},
	}

	p := probe.CoherenceProbe{}
	s := statefulset.GetUpgradeStrategy(c, p)

	g.Expect(s).To(BeAssignableToTypeOf(statefulset.ByPodUpgradeStrategy{}))
	g.Expect(s.IsOperatorManaged()).To(BeFalse())
}

func TestUseUpgradeStrategyByNode(t *testing.T) {
	g := NewGomegaWithT(t)

	c := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-deployment",
		},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			RollingUpdateStrategy: ptr.To(coh.UpgradeByNode),
		},
	}

	p := probe.CoherenceProbe{}
	s := statefulset.GetUpgradeStrategy(c, p)

	g.Expect(s).To(BeAssignableToTypeOf(statefulset.ByNodeUpgradeStrategy{}))
	g.Expect(s.IsOperatorManaged()).To(BeTrue())
}

func TestUseUpgradeStrategyManual(t *testing.T) {
	g := NewGomegaWithT(t)

	c := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-deployment",
		},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			RollingUpdateStrategy: ptr.To(coh.UpgradeManual),
		},
	}

	p := probe.CoherenceProbe{}
	s := statefulset.GetUpgradeStrategy(c, p)

	g.Expect(s).To(BeAssignableToTypeOf(statefulset.ManualUpgradeStrategy{}))
	g.Expect(s.IsOperatorManaged()).To(BeFalse())
}
