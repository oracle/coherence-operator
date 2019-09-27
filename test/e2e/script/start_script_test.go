/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package script

import (
	. "github.com/onsi/gomega"
	v1 "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	"k8s.io/utils/pointer"
	"testing"
)

func TestCoherenceDefaults(t *testing.T) {
	g := NewGomegaWithT(t)

	role := v1.CoherenceRoleSpec{}
	appData, cluster, err := RunScript(t, role)
	g.Expect(err).NotTo(HaveOccurred())

	// Assert default System properties
	g.Expect(appData.GetSystemProperty("coherence.ttl")).To(Equal("0"))
	g.Expect(appData.GetSystemProperty("coherence.cluster")).To(Equal(cluster.Name))
	g.Expect(appData.GetSystemProperty("coherence.role")).To(Equal(v1.DefaultRoleName))
	g.Expect(appData.GetSystemProperty("coherence.wka")).To(Equal(cluster.GetWkaServiceName()))
	g.Expect(appData.GetSystemProperty("coherence.distributed.persistence-mode")).To(Equal("on-demand"))

	// last arg should be DCS
	g.Expect(appData.Args[len(appData.Args)-1]).To(Equal("com.tangosol.net.DefaultCacheServer"))
}

func TestCoherenceRoleName(t *testing.T) {
	g := NewGomegaWithT(t)

	role := v1.CoherenceRoleSpec{
		Role: "data",
	}

	appData, _, err := RunScript(t, role)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(appData.GetSystemProperty("coherence.role")).To(Equal("data"))
}

func TestCoherenceCacheConfig(t *testing.T) {
	g := NewGomegaWithT(t)

	cfg := "test-cache-config.xml"
	role := v1.CoherenceRoleSpec{
		Coherence: &v1.CoherenceSpec{
			CacheConfig: &cfg,
		},
	}

	appData, _, err := RunScript(t, role)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(appData.GetSystemProperty("coherence.cacheconfig")).To(Equal(cfg))
}

func TestCoherenceOverride(t *testing.T) {
	g := NewGomegaWithT(t)

	cfg := "test-operational-config.xml"
	role := v1.CoherenceRoleSpec{
		Coherence: &v1.CoherenceSpec{
			OverrideConfig: &cfg,
		},
	}

	appData, _, err := RunScript(t, role)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(appData.GetSystemProperty("coherence.override")).To(Equal("k8s-coherence-override.xml"))
	g.Expect(appData.GetSystemProperty("coherence.k8s.override")).To(Equal(cfg))
}

func TestJvmHeapSize(t *testing.T) {
	g := NewGomegaWithT(t)

	role := v1.CoherenceRoleSpec{
		JVM: &v1.JVMSpec{
			HeapSize: pointer.StringPtr("10g"),
		},
	}

	appData, _, err := RunScript(t, role)
	g.Expect(err).NotTo(HaveOccurred())

	max := appData.FindJvmOption("-Xmx")
	g.Expect(len(max)).To(Equal(1))
	g.Expect(max[0]).To(Equal("-Xmx10g"))

	min := appData.FindJvmOption("-Xms")
	g.Expect(len(min)).To(Equal(1))
	g.Expect(min[0]).To(Equal("-Xms10g"))
}
