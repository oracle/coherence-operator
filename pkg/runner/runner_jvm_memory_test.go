/*
 * Copyright (c) 2020, 2025, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package runner

import (
	. "github.com/onsi/gomega"
	coh "github.com/oracle/coherence-operator/api/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	"testing"
)

func TestJvmHeapSize(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				JVM: &coh.JVMSpec{
					Memory: &coh.JvmMemorySpec{
						HeapSize: ptr.To("10g"),
					},
				},
			},
		},
	}

	verifyConfigFilesWithArgs(t, d, GetExpectedArgsFileContentWith("-XX:InitialHeapSize=10g", "-XX:MaxHeapSize=10g"))

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(t, d)

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(GetJavaCommand()))
	g.Expect(e.OsCmd.Args).To(ConsistOf(GetMinimalExpectedArgsWith(t, "-XX:InitialHeapSize=10g", "-XX:MaxHeapSize=10g")))
}

func TestJvmInitialHeapSize(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				JVM: &coh.JVMSpec{
					Memory: &coh.JvmMemorySpec{
						InitialHeapSize: ptr.To("10g"),
					},
				},
			},
		},
	}

	verifyConfigFilesWithArgs(t, d, GetExpectedArgsFileContentWith("-XX:InitialHeapSize=10g"))

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(t, d)

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(GetJavaCommand()))
	g.Expect(e.OsCmd.Args).To(ConsistOf(GetMinimalExpectedArgsWith(t, "-XX:InitialHeapSize=10g")))
}

func TestJvmMaxHeapSize(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				JVM: &coh.JVMSpec{
					Memory: &coh.JvmMemorySpec{
						MaxHeapSize: ptr.To("10g"),
					},
				},
			},
		},
	}

	verifyConfigFilesWithArgs(t, d, GetExpectedArgsFileContentWith("-XX:MaxHeapSize=10g"))

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(t, d)

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(GetJavaCommand()))
	g.Expect(e.OsCmd.Args).To(ConsistOf(GetMinimalExpectedArgsWith(t, "-XX:MaxHeapSize=10g")))
}

func TestJvmHeapSizeOverridesInitialAndMaxHeapSize(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				JVM: &coh.JVMSpec{
					Memory: &coh.JvmMemorySpec{
						HeapSize:        ptr.To("5g"),
						InitialHeapSize: ptr.To("1g"),
						MaxHeapSize:     ptr.To("10g"),
					},
				},
			},
		},
	}

	verifyConfigFilesWithArgs(t, d, GetExpectedArgsFileContentWith("-XX:InitialHeapSize=5g", "-XX:MaxHeapSize=5g"))

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(t, d)

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(GetJavaCommand()))
	g.Expect(e.OsCmd.Args).To(ConsistOf(GetMinimalExpectedArgsWith(t, "-XX:InitialHeapSize=5g", "-XX:MaxHeapSize=5g")))
}

func TestJvmMaxRam(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				JVM: &coh.JVMSpec{
					Memory: &coh.JvmMemorySpec{
						MaxRAM: ptr.To("10g"),
					},
				},
			},
		},
	}

	verifyConfigFilesWithArgs(t, d, GetExpectedArgsFileContentWith("-XX:MaxRAM=10g"))

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(t, d)

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(GetJavaCommand()))
	g.Expect(e.OsCmd.Args).To(ConsistOf(GetMinimalExpectedArgsWith(t, "-XX:MaxRAM=10g")))
}

func TestJvmRamPercent(t *testing.T) {
	g := NewGomegaWithT(t)

	pct := resource.MustParse("5.5")
	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				JVM: &coh.JVMSpec{
					Memory: &coh.JvmMemorySpec{
						Percentage: &pct,
					},
				},
			},
		},
	}

	verifyConfigFilesWithArgs(t, d, GetExpectedArgsFileContentWith("-XX:InitialRAMPercentage=5.500",
		"-XX:MaxRAMPercentage=5.500", "-XX:MinRAMPercentage=5.500"))

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(t, d)

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(GetJavaCommand()))
	g.Expect(e.OsCmd.Args).To(ConsistOf(GetMinimalExpectedArgsWith(t, "-XX:InitialRAMPercentage=5.500",
		"-XX:MaxRAMPercentage=5.500", "-XX:MinRAMPercentage=5.500")))
}

func TestJvmInitialRamPercent(t *testing.T) {
	g := NewGomegaWithT(t)

	pct := resource.MustParse("5.5")
	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				JVM: &coh.JVMSpec{
					Memory: &coh.JvmMemorySpec{
						InitialRAMPercentage: &pct,
					},
				},
			},
		},
	}

	verifyConfigFilesWithArgs(t, d, GetExpectedArgsFileContentWith("-XX:InitialRAMPercentage=5.500"))

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(t, d)

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(GetJavaCommand()))
	g.Expect(e.OsCmd.Args).To(ConsistOf(GetMinimalExpectedArgsWith(t, "-XX:InitialRAMPercentage=5.500")))
}

func TestJvmMaxRamPercent(t *testing.T) {
	g := NewGomegaWithT(t)

	pct := resource.MustParse("5.5")
	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				JVM: &coh.JVMSpec{
					Memory: &coh.JvmMemorySpec{
						MaxRAMPercentage: &pct,
					},
				},
			},
		},
	}

	verifyConfigFilesWithArgs(t, d, GetExpectedArgsFileContentWith("-XX:MaxRAMPercentage=5.500"))

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(t, d)

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(GetJavaCommand()))
	g.Expect(e.OsCmd.Args).To(ConsistOf(GetMinimalExpectedArgsWith(t, "-XX:MaxRAMPercentage=5.500")))
}

func TestJvmMinRamPercent(t *testing.T) {
	g := NewGomegaWithT(t)

	pct := resource.MustParse("5.5")
	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				JVM: &coh.JVMSpec{
					Memory: &coh.JvmMemorySpec{
						MinRAMPercentage: &pct,
					},
				},
			},
		},
	}

	verifyConfigFilesWithArgs(t, d, GetExpectedArgsFileContentWith("-XX:MinRAMPercentage=5.500"))

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(t, d)

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(GetJavaCommand()))
	g.Expect(e.OsCmd.Args).To(ConsistOf(GetMinimalExpectedArgsWith(t, "-XX:MinRAMPercentage=5.500")))
}

func TestJvmRamPercentOverridesInitialMaxAndMin(t *testing.T) {
	g := NewGomegaWithT(t)

	pct := resource.MustParse("5.5")
	pctInit := resource.MustParse("1")
	pctMin := resource.MustParse("2")
	pctMax := resource.MustParse("10")
	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				JVM: &coh.JVMSpec{
					Memory: &coh.JvmMemorySpec{
						Percentage:           &pct,
						InitialRAMPercentage: &pctInit,
						MinRAMPercentage:     &pctMin,
						MaxRAMPercentage:     &pctMax,
					},
				},
			},
		},
	}

	verifyConfigFilesWithArgs(t, d, GetExpectedArgsFileContentWith("-XX:InitialRAMPercentage=5.500",
		"-XX:MaxRAMPercentage=5.500", "-XX:MinRAMPercentage=5.500"))

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(t, d)

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(GetJavaCommand()))
	g.Expect(e.OsCmd.Args).To(ConsistOf(GetMinimalExpectedArgsWith(t, "-XX:InitialRAMPercentage=5.500",
		"-XX:MaxRAMPercentage=5.500", "-XX:MinRAMPercentage=5.500")))
}

func TestJvmStackSize(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				JVM: &coh.JVMSpec{
					Memory: &coh.JvmMemorySpec{
						StackSize: ptr.To("500k"),
					},
				},
			},
		},
	}

	verifyConfigFilesWithArgs(t, d, GetExpectedArgsFileContentWith("-Xss500k"))

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(t, d)

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(GetJavaCommand()))
	g.Expect(e.OsCmd.Args).To(ConsistOf(GetMinimalExpectedArgsWith(t, "-Xss500k")))
}

func TestJvmMetaspaceSize(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				JVM: &coh.JVMSpec{
					Memory: &coh.JvmMemorySpec{
						MetaspaceSize: ptr.To("5g"),
					},
				},
			},
		},
	}

	verifyConfigFilesWithArgs(t, d, GetExpectedArgsFileContentWith("-XX:MetaspaceSize=5g", "-XX:MaxMetaspaceSize=5g"))

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(t, d)

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(GetJavaCommand()))
	g.Expect(e.OsCmd.Args).To(ConsistOf(GetMinimalExpectedArgsWith(t,
		"-XX:MetaspaceSize=5g", "-XX:MaxMetaspaceSize=5g")))
}

func TestJvmDirectMemorySize(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				JVM: &coh.JVMSpec{
					Memory: &coh.JvmMemorySpec{
						DirectMemorySize: ptr.To("5g"),
					},
				},
			},
		},
	}

	verifyConfigFilesWithArgs(t, d, GetExpectedArgsFileContentWith("-XX:MaxDirectMemorySize=5g"))

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(t, d)

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(GetJavaCommand()))
	g.Expect(e.OsCmd.Args).To(ConsistOf(GetMinimalExpectedArgsWith(t, "-XX:MaxDirectMemorySize=5g")))
}

func TestJvmNativeMemoryTracking(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				JVM: &coh.JVMSpec{
					Memory: &coh.JvmMemorySpec{
						NativeMemoryTracking: ptr.To("detail"),
					},
				},
			},
		},
	}

	expectedFileArgs := append(GetExpectedArgsFileContentWithoutPrefix("-XX:NativeMemoryTracking="), "-XX:NativeMemoryTracking=detail")
	verifyConfigFilesWithArgs(t, d, expectedFileArgs)

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(t, d)

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(GetJavaCommand()))
	expected := append(RemoveArgWithPrefix(GetMinimalExpectedArgs(t), "-XX:NativeMemoryTracking="), "-XX:NativeMemoryTracking=detail")
	g.Expect(e.OsCmd.Args).To(ConsistOf(expected))
}

func TestJvmOOMHeapDumpOff(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				JVM: &coh.JVMSpec{
					Memory: &coh.JvmMemorySpec{
						OnOutOfMemory: &coh.JvmOutOfMemorySpec{
							HeapDump: ptr.To(false),
						},
					},
				},
			},
		},
	}

	expectedArgs := GetExpectedArgsFileContentWithoutPrefix("-XX:+HeapDumpOnOutOfMemoryError")
	verifyConfigFilesWithArgs(t, d, expectedArgs)

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(t, d)

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(GetJavaCommand()))
	expected := RemoveArgWithPrefix(GetMinimalExpectedArgs(t), "-XX:+HeapDumpOnOutOfMemoryError")
	g.Expect(e.OsCmd.Args).To(ConsistOf(expected))
}

func TestJvmOOMExitOff(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				JVM: &coh.JVMSpec{
					Memory: &coh.JvmMemorySpec{
						OnOutOfMemory: &coh.JvmOutOfMemorySpec{
							Exit: ptr.To(false),
						},
					},
				},
			},
		},
	}

	expectedArgs := GetExpectedArgsFileContentWithoutPrefix("-XX:+ExitOnOutOfMemoryError")
	verifyConfigFilesWithArgs(t, d, expectedArgs)

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(t, d)

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(GetJavaCommand()))
	expected := RemoveArgWithPrefix(GetMinimalExpectedArgs(t), "-XX:+ExitOnOutOfMemoryError")
	g.Expect(e.OsCmd.Args).To(ConsistOf(expected))
}
