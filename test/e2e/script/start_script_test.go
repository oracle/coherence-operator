/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package script

import (
	. "github.com/onsi/gomega"
	v1 "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/pointer"
	"testing"
)

// ----- Coherence configuration tests --------------------------------------

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

	g.Expect(appData.HasSystemProperty("coherence.distributed.local.storage")).To(BeFalse())

	// Assert default environment variables
	g.Expect(appData.GetEnv("COH_MGMT_ENABLED")).To(Equal("false"))
	g.Expect(appData.GetEnv("COH_MGMT_HTTP_PORT")).To(Equal("30000"))
	g.Expect(appData.GetEnv("COH_METRICS_ENABLED")).To(Equal("false"))
	g.Expect(appData.GetEnv("COH_METRICS_PORT")).To(Equal("9612"))

	// Last but one arg should be the Operator's custom main
	g.Expect(appData.Args[len(appData.Args)-2]).To(Equal("com.oracle.coherence.k8s.Main"))
	// Last arg should be DCS
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

func TestCoherenceStorageEnabled(t *testing.T) {
	g := NewGomegaWithT(t)

	role := v1.CoherenceRoleSpec{
		Coherence: &v1.CoherenceSpec{
			StorageEnabled: pointer.BoolPtr(true),
		},
	}

	appData, _, err := RunScript(t, role)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(appData.GetSystemProperty("coherence.distributed.localstorage")).To(Equal("true"))
}

func TestCoherenceStorageDisabled(t *testing.T) {
	g := NewGomegaWithT(t)

	role := v1.CoherenceRoleSpec{
		Coherence: &v1.CoherenceSpec{
			StorageEnabled: pointer.BoolPtr(false),
		},
	}

	appData, _, err := RunScript(t, role)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(appData.GetSystemProperty("coherence.distributed.localstorage")).To(Equal("false"))
}

func TestCoherenceLogLevel(t *testing.T) {
	g := NewGomegaWithT(t)

	role := v1.CoherenceRoleSpec{
		Coherence: &v1.CoherenceSpec{
			LogLevel: pointer.Int32Ptr(9),
		},
	}

	appData, _, err := RunScript(t, role)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(appData.GetSystemProperty("coherence.log.level")).To(Equal("9"))
}

// ----- Application configuration tests ------------------------------------

func TestApplicationMain(t *testing.T) {
	g := NewGomegaWithT(t)

	main := "TestMain"
	role := v1.CoherenceRoleSpec{
		Application: &v1.ApplicationSpec{
			Main: &main,
		},
	}

	appData, _, err := RunScript(t, role)
	g.Expect(err).NotTo(HaveOccurred())
	// last arg should be main
	g.Expect(appData.Args[len(appData.Args)-1]).To(Equal(main))
}

func TestApplicationArgs(t *testing.T) {
	g := NewGomegaWithT(t)

	args := []string{"One", "Two"}
	role := v1.CoherenceRoleSpec{
		Application: &v1.ApplicationSpec{
			Args: args,
		},
	}

	appData, _, err := RunScript(t, role)
	g.Expect(err).NotTo(HaveOccurred())

	g.Expect(appData.Args[len(appData.Args)-4]).To(Equal("com.oracle.coherence.k8s.Main"))
	g.Expect(appData.Args[len(appData.Args)-3]).To(Equal("com.tangosol.net.DefaultCacheServer"))

	// last two args should match the program args
	g.Expect(appData.Args[len(appData.Args)-2:]).To(Equal(args))
}

// ----- JVM configuration tests --------------------------------------------

func TestJvmHeapSize(t *testing.T) {
	g := NewGomegaWithT(t)

	role := v1.CoherenceRoleSpec{
		JVM: &v1.JVMSpec{
			Memory: &v1.JvmMemorySpec{
				HeapSize: pointer.StringPtr("10g"),
			},
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

func TestEnvironmentVariables(t *testing.T) {
	g := NewGomegaWithT(t)

	ev1 := corev1.EnvVar{Name: "One", Value: "1"}
	ev2 := corev1.EnvVar{Name: "Two", Value: "2"}

	role := v1.CoherenceRoleSpec{
		Env: []corev1.EnvVar{ev1, ev2},
	}

	appData, _, err := RunScript(t, role)
	g.Expect(err).NotTo(HaveOccurred())

	g.Expect(appData.GetEnv("One")).To(Equal("1"))
	g.Expect(appData.GetEnv("Two")).To(Equal("2"))
}
