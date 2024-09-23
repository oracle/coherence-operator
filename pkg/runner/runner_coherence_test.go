/*
 * Copyright (c) 2020, 2024, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package runner

import (
	. "github.com/onsi/gomega"
	coh "github.com/oracle/coherence-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	"testing"
)

func TestCoherenceClusterName(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			Cluster: ptr.To("test-cluster"),
		},
	}

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(d)

	expectedCommand := GetJavaCommand()
	expectedArgs := append(GetMinimalExpectedArgsWithoutPrefix("-Dcoherence.cluster="),
		"-Dcoherence.cluster=test-cluster")

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(expectedCommand))
	g.Expect(e.OsCmd.Args).To(ConsistOf(expectedArgs))
}

func TestCoherenceCacheConfig(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Coherence: &coh.CoherenceSpec{
					CacheConfig: ptr.To("test-config.xml"),
				},
			},
		},
	}

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(d)

	expectedCommand := GetJavaCommand()
	expectedArgs := append(GetMinimalExpectedArgsWithoutPrefix("-Dcoherence.cacheconfig="),
		"-Dcoherence.cacheconfig=test-config.xml")

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(expectedCommand))
	g.Expect(e.OsCmd.Args).To(ConsistOf(expectedArgs))
}

func TestCoherenceOperationalConfig(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Coherence: &coh.CoherenceSpec{
					OverrideConfig: ptr.To("test-override.xml"),
				},
			},
		},
	}

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(d)

	expectedCommand := GetJavaCommand()
	expectedArgs := append(GetMinimalExpectedArgsWithoutPrefix("-Dcoherence.k8s.override="),
		"-Dcoherence.k8s.override=test-override.xml")

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(expectedCommand))
	g.Expect(e.OsCmd.Args).To(ConsistOf(expectedArgs))
}

func TestCoherenceStorageEnabledTrue(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Coherence: &coh.CoherenceSpec{
					StorageEnabled: ptr.To(true),
				},
			},
		},
	}

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(d)

	expectedCommand := GetJavaCommand()
	expectedArgs := append(GetMinimalExpectedArgsWithoutPrefix("-Dcoherence.distributed.localstorage="),
		"-Dcoherence.distributed.localstorage=true")

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(expectedCommand))
	g.Expect(e.OsCmd.Args).To(ConsistOf(expectedArgs))
}

func TestCoherenceStorageEnabledFalse(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Coherence: &coh.CoherenceSpec{
					StorageEnabled: ptr.To(false),
				},
			},
		},
	}

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(d)

	expectedCommand := GetJavaCommand()
	expectedArgs := append(GetMinimalExpectedArgsWithoutPrefix("-Dcoherence.distributed.localstorage="),
		"-Dcoherence.distributed.localstorage=false")

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(expectedCommand))
	g.Expect(e.OsCmd.Args).To(ConsistOf(expectedArgs))
}

func TestCoherenceExcludeFromWKATrue(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Coherence: &coh.CoherenceSpec{
					ExcludeFromWKA: ptr.To(true),
				},
			},
		},
	}

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(d)

	expectedCommand := GetJavaCommand()
	expectedArgs := GetMinimalExpectedArgs()

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(expectedCommand))
	g.Expect(e.OsCmd.Args).To(ConsistOf(expectedArgs))
}

func TestCoherenceLogLevel(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Coherence: &coh.CoherenceSpec{
					LogLevel: ptr.To(int32(9)),
				},
			},
		},
	}

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(d)

	expectedCommand := GetJavaCommand()
	expectedArgs := append(GetMinimalExpectedArgsWithoutPrefix("-Dcoherence.log.level="),
		"-Dcoherence.log.level=9")

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(expectedCommand))
	g.Expect(e.OsCmd.Args).To(ConsistOf(expectedArgs))
}

func TestCoherenceTracingRatio(t *testing.T) {
	g := NewGomegaWithT(t)

	ratio := resource.MustParse("0.01234")
	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Coherence: &coh.CoherenceSpec{
					Tracing: &coh.CoherenceTracingSpec{Ratio: &ratio},
				},
			},
		},
	}

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(d)

	expectedCommand := GetJavaCommand()
	expectedArgs := append(GetMinimalExpectedArgs(), "-Dcoherence.tracing.ratio=0.012340")

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(expectedCommand))
	g.Expect(e.OsCmd.Args).To(ConsistOf(expectedArgs))
}

func TestCoherenceAllowEndangeredEmptyList(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Coherence: &coh.CoherenceSpec{
					AllowEndangeredForStatusHA: []string{},
				},
			},
		},
	}

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(d)

	expectedCommand := GetJavaCommand()
	expectedArgs := GetMinimalExpectedArgs()

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(expectedCommand))
	g.Expect(e.OsCmd.Args).To(ConsistOf(expectedArgs))
}

func TestCoherenceAllowEndangered(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Coherence: &coh.CoherenceSpec{
					AllowEndangeredForStatusHA: []string{"foo", "bar"},
				},
			},
		},
	}

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(d)

	expectedCommand := GetJavaCommand()
	expectedArgs := append(GetMinimalExpectedArgs(), "-Dcoherence.k8s.operator.statusha.allowendangered=foo,bar")

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(expectedCommand))
	g.Expect(e.OsCmd.Args).To(ConsistOf(expectedArgs))
}

func TestCoherenceExistingWKADeploymentSameNamespace(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Coherence: &coh.CoherenceSpec{
					WKA: &coh.CoherenceWKASpec{
						Deployment: "data",
						Namespace:  "foo",
					},
				},
			},
		},
	}

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(d)

	expectedCommand := GetJavaCommand()
	expectedArgs := GetMinimalExpectedArgsWithoutPrefix("-Dcoherence.wka")
	expectedArgs = append(expectedArgs, "-Dcoherence.wka=data"+coh.WKAServiceNameSuffix+".foo.svc")

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(expectedCommand))
	g.Expect(e.OsCmd.Args).To(ConsistOf(expectedArgs))
}

func TestCoherenceExistingWKADeploymentDifferentNamespace(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Coherence: &coh.CoherenceSpec{
					WKA: &coh.CoherenceWKASpec{
						Deployment: "data",
						Namespace:  "back-end",
					},
				},
			},
		},
	}

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(d)

	expectedCommand := GetJavaCommand()
	expectedArgs := GetMinimalExpectedArgsWithoutPrefix("-Dcoherence.wka")
	expectedArgs = append(expectedArgs, "-Dcoherence.wka=data"+coh.WKAServiceNameSuffix+".back-end.svc")

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(expectedCommand))
	g.Expect(e.OsCmd.Args).To(ConsistOf(expectedArgs))
}

func TestCoherenceEnableIpMonitor(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Coherence: &coh.CoherenceSpec{
					EnableIPMonitor: ptr.To(true),
				},
			},
		},
	}

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(d)

	expectedCommand := GetJavaCommand()
	expectedArgs := GetMinimalExpectedArgsWithoutPrefix("-Dcoherence.ipmonitor.pingtimeout")

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(expectedCommand))
	g.Expect(e.OsCmd.Args).To(ConsistOf(expectedArgs))
}

func TestCoherenceDisableIpMonitor(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Coherence: &coh.CoherenceSpec{
					EnableIPMonitor: ptr.To(false),
				},
			},
		},
	}

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(d)

	expectedCommand := GetJavaCommand()
	expectedArgs := GetMinimalExpectedArgs()

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(expectedCommand))
	g.Expect(e.OsCmd.Args).To(ConsistOf(expectedArgs))
}

func TestCoherenceDefaultIpMonitor(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Coherence: &coh.CoherenceSpec{
					EnableIPMonitor: nil,
				},
			},
		},
	}

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(d)

	expectedCommand := GetJavaCommand()
	expectedArgs := GetMinimalExpectedArgs()

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(expectedCommand))
	g.Expect(e.OsCmd.Args).To(ConsistOf(expectedArgs))
}

func TestCoherenceSiteEAndRackEnvVarSet(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Env: []corev1.EnvVar{
					{
						Name:  coh.EnvVarCoherenceSite,
						Value: "test-site",
					},
					{
						Name:  coh.EnvVarCoherenceRack,
						Value: "test-rack",
					},
				},
			},
		},
	}

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeploymentWithSkipSite(d, false)

	expectedCommand := GetJavaCommand()
	// site and rack system properties should not be set
	expectedArgs := RemoveArg(GetMinimalExpectedArgs(), "-Dcoherence.site")
	expectedArgs = RemoveArg(expectedArgs, "-Dcoherence.rack")

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(expectedCommand))
	g.Expect(e.OsCmd.Args).To(ConsistOf(expectedArgs))
}

func TestCoherenceSiteAndRackEnvVarSetFromOtherEnvVar(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Env: []corev1.EnvVar{
					{
						Name:  coh.EnvVarCoherenceSite,
						Value: "${TEST_SITE_VAR}",
					},
					{
						Name:  "TEST_SITE_VAR",
						Value: "test-site",
					},
					{
						Name:  coh.EnvVarCoherenceRack,
						Value: "${TEST_RACK_VAR}",
					},
					{
						Name:  "TEST_RACK_VAR",
						Value: "test-rack",
					},
				},
			},
		},
	}

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeploymentWithSkipSite(d, false)

	expectedCommand := GetJavaCommand()
	expectedArgs := append(GetMinimalExpectedArgs(), "-Dcoherence.site=test-site")
	expectedArgs = append(expectedArgs, "-Dcoherence.rack=test-rack")

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(expectedCommand))
	g.Expect(e.OsCmd.Args).To(ConsistOf(expectedArgs))
}

func TestCoherenceSiteAndRackEnvVarSetFromOtherEnvVarWhenRackIsMissing(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Env: []corev1.EnvVar{
					{
						Name:  coh.EnvVarCoherenceSite,
						Value: "${TEST_SITE_VAR}",
					},
					{
						Name:  "TEST_SITE_VAR",
						Value: "test-site",
					},
					{
						Name:  coh.EnvVarCoherenceRack,
						Value: "${TEST_RACK_VAR}",
					},
				},
			},
		},
	}

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeploymentWithSkipSite(d, false)

	expectedCommand := GetJavaCommand()
	expectedArgs := append(GetMinimalExpectedArgs(), "-Dcoherence.site=test-site")
	// rack should be set to site
	expectedArgs = append(expectedArgs, "-Dcoherence.rack=test-site")

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(expectedCommand))
	g.Expect(e.OsCmd.Args).To(ConsistOf(expectedArgs))
}

func TestCoherenceSiteAndRackFromFile(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Env: []corev1.EnvVar{
					{
						Name:  coh.EnvVarCohSite,
						Value: "test-site.txt",
					},
					{
						Name:  coh.EnvVarCohRack,
						Value: "test-rack.txt",
					},
				},
			},
		},
	}

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeploymentWithSkipSite(d, false)

	expectedCommand := GetJavaCommand()
	expectedArgs := append(GetMinimalExpectedArgs(), "-Dcoherence.site=site-from-file")
	expectedArgs = append(expectedArgs, "-Dcoherence.rack=rack-from-file")

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(expectedCommand))
	g.Expect(e.OsCmd.Args).To(ConsistOf(expectedArgs))
}
