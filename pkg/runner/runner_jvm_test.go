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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	"testing"
)

func TestJvmArgsEmpty(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				JVM: &coh.JVMSpec{
					Args: []string{},
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

func TestJvmArgs(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				JVM: &coh.JVMSpec{
					Args: []string{"-Dfoo=foo-value", "-Dbar=bar-value"},
				},
			},
		},
	}

	verifyConfigFilesWithArgs(t, d, GetExpectedArgsFileContentWith("-Dfoo=foo-value", "-Dbar=bar-value"))

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(t, d)

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(GetJavaCommand()))
	g.Expect(e.OsCmd.Args).To(ConsistOf(GetMinimalExpectedArgsWith(t, "-Dfoo=foo-value", "-Dbar=bar-value")))
}

func TestJvmArgsWithEnvExpansion(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				JVM: &coh.JVMSpec{
					Args: []string{"-Dfoo=$FOO", "-Dbar=${BAR}"},
				},
				Env: []corev1.EnvVar{
					{Name: "FOO", Value: "foo-value"},
					{Name: "BAR", Value: "bar-value"},
				},
			},
		},
	}

	verifyConfigFilesWithArgs(t, d, GetExpectedArgsFileContentWith("-Dfoo=foo-value", "-Dbar=bar-value"))

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(t, d)

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(GetJavaCommand()))
	g.Expect(e.OsCmd.Args).To(ConsistOf(GetMinimalExpectedArgsWith(t, "-Dfoo=foo-value", "-Dbar=bar-value")))
}

func TestJvmUseContainerLimitsFalse(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				JVM: &coh.JVMSpec{
					UseContainerLimits: ptr.To(false),
				},
			},
		},
	}

	expectedArgs := GetExpectedArgsFileContentWithoutPrefix("-XX:+UseContainerSupport")
	verifyConfigFilesWithArgs(t, d, expectedArgs)

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(t, d)

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(GetJavaCommand()))
	expected := RemoveArgWithPrefix(GetMinimalExpectedArgsWith(t), "-XX:+UseContainerSupport")
	g.Expect(e.OsCmd.Args).To(ConsistOf(expected))
}

func TestJvmUseContainerLimitsTrue(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				JVM: &coh.JVMSpec{
					UseContainerLimits: ptr.To(true),
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

func TestJvmGarbageCollectorG1(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				JVM: &coh.JVMSpec{
					Gc: &coh.JvmGarbageCollectorSpec{
						Collector: ptr.To("g1"),
					},
				},
			},
		},
	}

	expectedArgs := append(GetExpectedArgsFileContentWithoutPrefix("-XX:+UseG1GC"), "-XX:+UseG1GC")
	verifyConfigFilesWithArgs(t, d, expectedArgs)

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(t, d)

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(GetJavaCommand()))
	expected := append(RemoveArgWithPrefix(GetMinimalExpectedArgs(t), "-XX:+UseG1GC"), "-XX:+UseG1GC")
	g.Expect(e.OsCmd.Args).To(ConsistOf(expected))
}

func TestJvmGarbageCollectorCMS(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				JVM: &coh.JVMSpec{
					Gc: &coh.JvmGarbageCollectorSpec{
						Collector: ptr.To("cms"),
					},
				},
			},
		},
	}

	expectedArgs := append(GetExpectedArgsFileContentWithoutPrefix("-XX:+UseG1GC"), "-XX:+UseConcMarkSweepGC")
	verifyConfigFilesWithArgs(t, d, expectedArgs)

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(t, d)

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(GetJavaCommand()))
	expected := append(RemoveArgWithPrefix(GetMinimalExpectedArgs(t), "-XX:+UseG1GC"), "-XX:+UseConcMarkSweepGC")
	g.Expect(e.OsCmd.Args).To(ConsistOf(expected))
}

func TestJvmGarbageCollectorParallel(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				JVM: &coh.JVMSpec{
					Gc: &coh.JvmGarbageCollectorSpec{
						Collector: ptr.To("parallel"),
					},
				},
			},
		},
	}

	expectedArgs := append(GetExpectedArgsFileContentWithoutPrefix("-XX:+UseG1GC"), "-XX:+UseParallelGC")
	verifyConfigFilesWithArgs(t, d, expectedArgs)

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(t, d)

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(GetJavaCommand()))
	expected := append(RemoveArgWithPrefix(GetMinimalExpectedArgs(t), "-XX:+UseG1GC"), "-XX:+UseParallelGC")
	g.Expect(e.OsCmd.Args).To(ConsistOf(expected))
}

func TestJvmGarbageCollectorSerial(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				JVM: &coh.JVMSpec{
					Gc: &coh.JvmGarbageCollectorSpec{
						Collector: ptr.To("serial"),
					},
				},
			},
		},
	}

	expectedArgs := append(GetExpectedArgsFileContentWithoutPrefix("-XX:+UseG1GC"), "-XX:+UseSerialGC")
	verifyConfigFilesWithArgs(t, d, expectedArgs)

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(t, d)

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(GetJavaCommand()))
	expected := append(RemoveArgWithPrefix(GetMinimalExpectedArgs(t), "-XX:+UseG1GC"), "-XX:+UseSerialGC")
	g.Expect(e.OsCmd.Args).To(ConsistOf(expected))
}

func TestJvmGarbageCollectorZGC(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				JVM: &coh.JVMSpec{
					Gc: &coh.JvmGarbageCollectorSpec{
						Collector: ptr.To("zgc"),
					},
				},
			},
		},
	}

	expectedArgs := append(GetExpectedArgsFileContentWithoutPrefix("-XX:+UseG1GC"), "-XX:+UseZGC")
	verifyConfigFilesWithArgs(t, d, expectedArgs)

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(t, d)

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(GetJavaCommand()))
	expected := append(RemoveArgWithPrefix(GetMinimalExpectedArgs(t), "-XX:+UseG1GC"), "-XX:+UseZGC")
	g.Expect(e.OsCmd.Args).To(ConsistOf(expected))
}

func TestJvmGarbageCollectorLoggingTrue(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				JVM: &coh.JVMSpec{
					Gc: &coh.JvmGarbageCollectorSpec{
						Logging: ptr.To(true),
					},
				},
			},
		},
	}

	expectedArgs := append(GetExpectedArgsFileContent(),
		"-verbose:gc",
		"-XX:+PrintGCDetails",
		"-XX:+PrintGCTimeStamps",
		"-XX:+PrintHeapAtGC",
		"-XX:+PrintTenuringDistribution",
		"-XX:+PrintGCApplicationStoppedTime",
		"-XX:+PrintGCApplicationConcurrentTime")

	verifyConfigFilesWithArgs(t, d, expectedArgs)

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(t, d)

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(GetJavaCommand()))
	g.Expect(e.OsCmd.Args).To(ConsistOf(GetMinimalExpectedArgsWith(t,
		"-verbose:gc",
		"-XX:+PrintGCDetails",
		"-XX:+PrintGCTimeStamps",
		"-XX:+PrintHeapAtGC",
		"-XX:+PrintTenuringDistribution",
		"-XX:+PrintGCApplicationStoppedTime",
		"-XX:+PrintGCApplicationConcurrentTime")))
}

func TestJvmGarbageCollectorArgsEmpty(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				JVM: &coh.JVMSpec{
					Gc: &coh.JvmGarbageCollectorSpec{
						Args: []string{},
					},
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

func TestJvmGarbageCollectorArgs(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				JVM: &coh.JVMSpec{
					Gc: &coh.JvmGarbageCollectorSpec{
						Args: []string{"-XX:Arg1", "-XX:Arg2"},
					},
				},
			},
		},
	}

	verifyConfigFilesWithArgs(t, d, GetExpectedArgsFileContentWith("-XX:Arg1", "-XX:Arg2"))

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(t, d)

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(GetJavaCommand()))
	g.Expect(e.OsCmd.Args).To(ConsistOf(GetMinimalExpectedArgsWith(t, "-XX:Arg1", "-XX:Arg2")))
}
