/*
 * Copyright (c) 2020, 2025, Oracle and/or its affiliates.
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

	expectedArgsFile := append(GetExpectedArgsFileContentWithoutPrefix("-Dcoherence.cluster="),
		"-Dcoherence.cluster=test-cluster")

	verifyConfigFilesWithArgs(t, d, expectedArgsFile)

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(t, d)

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(GetJavaCommand()))
	expectedArgs := append(GetExpectedArgsWithoutPrefix(t, "-Dcoherence.cluster="),
		"-Dcoherence.cluster=test-cluster")
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

	expectedFileArgs := append(GetExpectedArgsFileContentWithoutPrefix("-Dcoherence.cacheconfig="),
		"-Dcoherence.cacheconfig=test-config.xml")

	verifyConfigFilesWithArgs(t, d, expectedFileArgs)

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(t, d)

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(GetJavaCommand()))
	expectedArgs := append(GetExpectedArgsWithoutPrefix(t, "-Dcoherence.cacheconfig="),
		"-Dcoherence.cacheconfig=test-config.xml")
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

	expectedFileArgs := append(GetExpectedArgsFileContentWithoutPrefix("-Dcoherence.k8s.override="),
		"-Dcoherence.k8s.override=test-override.xml")

	verifyConfigFilesWithArgs(t, d, expectedFileArgs)

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(t, d)

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(GetJavaCommand()))
	expectedArgs := append(GetExpectedArgsWithoutPrefix(t, "-Dcoherence.k8s.override="),
		"-Dcoherence.k8s.override=test-override.xml")
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

	expectedFileArgs := append(GetExpectedArgsFileContentWithoutPrefix("-Dcoherence.distributed.localstorage="),
		"-Dcoherence.distributed.localstorage=true")

	verifyConfigFilesWithArgs(t, d, expectedFileArgs)

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(t, d)

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(GetJavaCommand()))
	expectedArgs := append(GetExpectedArgsWithoutPrefix(t, "-Dcoherence.distributed.localstorage="),
		"-Dcoherence.distributed.localstorage=true")
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

	expectedFileArgs := append(GetExpectedArgsFileContentWithoutPrefix("-Dcoherence.distributed.localstorage="),
		"-Dcoherence.distributed.localstorage=false")

	verifyConfigFilesWithArgs(t, d, expectedFileArgs)

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(t, d)

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(GetJavaCommand()))
	expectedArgs := append(GetExpectedArgsWithoutPrefix(t, "-Dcoherence.distributed.localstorage="),
		"-Dcoherence.distributed.localstorage=false")
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

	verifyConfigFilesWithArgs(t, d, GetExpectedArgsFileContent())

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(t, d)

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(GetJavaCommand()))
	g.Expect(e.OsCmd.Args).To(ConsistOf(GetMinimalExpectedArgs(t)))
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

	expectedFileArgs := append(GetExpectedArgsFileContentWithoutPrefix("-Dcoherence.log.level="),
		"-Dcoherence.log.level=9")

	verifyConfigFilesWithArgs(t, d, expectedFileArgs)

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(t, d)

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(GetJavaCommand()))
	expectedArgs := append(GetExpectedArgsWithoutPrefix(t, "-Dcoherence.log.level="),
		"-Dcoherence.log.level=9")
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

	verifyConfigFilesWithArgs(t, d, GetExpectedArgsFileContentWith("-Dcoherence.tracing.ratio=0.012340"))

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(t, d)

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(GetJavaCommand()))
	expectedArgs := append(GetMinimalExpectedArgs(t), "-Dcoherence.tracing.ratio=0.012340")
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

	verifyConfigFilesWithArgs(t, d, GetExpectedArgsFileContent())

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(t, d)

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(GetJavaCommand()))
	g.Expect(e.OsCmd.Args).To(ConsistOf(GetMinimalExpectedArgs(t)))
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

	verifyConfigFilesWithArgs(t, d, GetExpectedArgsFileContentWith("-Dcoherence.operator.statusha.allowendangered=foo,bar"))

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(t, d)

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(GetJavaCommand()))
	expectedArgs := append(GetMinimalExpectedArgs(t), "-Dcoherence.operator.statusha.allowendangered=foo,bar")
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

	expectedFileArgs := GetExpectedArgsFileContentWithoutPrefix("-Dcoherence.wka")
	expectedFileArgs = append(expectedFileArgs, "-Dcoherence.wka=data"+coh.WKAServiceNameSuffix+".foo.svc")

	verifyConfigFilesWithArgs(t, d, expectedFileArgs)

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(t, d)

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(GetJavaCommand()))
	expectedArgs := append(GetExpectedArgsWithoutPrefix(t, "-Dcoherence.wka"),
		"-Dcoherence.wka=data"+coh.WKAServiceNameSuffix+".foo.svc")
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

	expectedFileArgs := GetExpectedArgsFileContentWithoutPrefix("-Dcoherence.wka")
	expectedFileArgs = append(expectedFileArgs, "-Dcoherence.wka=data"+coh.WKAServiceNameSuffix+".back-end.svc")

	verifyConfigFilesWithArgs(t, d, expectedFileArgs)

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(t, d)

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(GetJavaCommand()))
	expectedArgs := append(GetExpectedArgsWithoutPrefix(t, "-Dcoherence.wka"),
		"-Dcoherence.wka=data"+coh.WKAServiceNameSuffix+".back-end.svc")
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

	expectedArgs := GetExpectedArgsFileContentWithoutPrefix("-Dcoherence.ipmonitor.pingtimeout")

	verifyConfigFilesWithArgs(t, d, expectedArgs)

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(t, d)

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(GetJavaCommand()))
	g.Expect(e.OsCmd.Args).To(ConsistOf(GetExpectedArgsWithoutPrefix(t, "-Dcoherence.ipmonitor.pingtimeout")))
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

	verifyConfigFilesWithArgs(t, d, GetExpectedArgsFileContent())

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(t, d)

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(GetJavaCommand()))
	g.Expect(e.OsCmd.Args).To(ConsistOf(GetMinimalExpectedArgs(t)))
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

	verifyConfigFilesWithArgs(t, d, GetExpectedArgsFileContent())

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(t, d)

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(GetJavaCommand()))
	g.Expect(e.OsCmd.Args).To(ConsistOf(GetMinimalExpectedArgs(t)))
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
						Value: "site-foo",
					},
					{
						Name:  coh.EnvVarCoherenceRack,
						Value: "rack-bar",
					},
				},
			},
		},
	}

	verifyConfigFilesWithArgsWithSkipSite(t, d, GetExpectedArgsFileContentWith("-Dcoherence.site=site-foo",
		"-Dcoherence.rack=rack-bar"), false)

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeploymentWithSkipSite(t, d, false)

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(GetJavaCommand()))
	expectedArgs := append(GetMinimalExpectedArgs(t), "-Dcoherence.site=site-foo", "-Dcoherence.rack=rack-bar")
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
						Value: "x-site",
					},
					{
						Name:  coh.EnvVarCoherenceRack,
						Value: "${TEST_RACK_VAR}",
					},
					{
						Name:  "TEST_RACK_VAR",
						Value: "x-rack",
					},
				},
			},
		},
	}

	verifyConfigFilesWithArgsWithSkipSite(t, d, GetExpectedArgsFileContentWith("-Dcoherence.site=x-site",
		"-Dcoherence.rack=x-rack"), false)

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeploymentWithSkipSite(t, d, false)

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(GetJavaCommand()))
	expectedArgs := append(GetMinimalExpectedArgs(t), "-Dcoherence.site=x-site", "-Dcoherence.rack=x-rack")
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

	// rack should be set to site
	expectedFileArgs := append(GetExpectedArgsFileContent(), "-Dcoherence.site=test-site", "-Dcoherence.rack=test-site")

	verifyConfigFilesWithArgsWithSkipSite(t, d, expectedFileArgs, false)

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeploymentWithSkipSite(t, d, false)

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(GetJavaCommand()))
	expectedArgs := append(GetMinimalExpectedArgs(t), "-Dcoherence.site=test-site",
		"-Dcoherence.rack=test-site")
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

	expectedFileArgs := append(GetExpectedArgsFileContent(), "-Dcoherence.site=site-from-file",
		"-Dcoherence.rack=rack-from-file")

	verifyConfigFilesWithArgsWithSkipSite(t, d, expectedFileArgs, false)

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeploymentWithSkipSite(t, d, false)

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(GetJavaCommand()))
	expectedArgs := append(GetMinimalExpectedArgs(t), "-Dcoherence.site=site-from-file",
		"-Dcoherence.rack=rack-from-file")
	g.Expect(e.OsCmd.Args).To(ConsistOf(expectedArgs))
}
