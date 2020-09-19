/*
 * Copyright (c) 2020 Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package runner

import (
	. "github.com/onsi/gomega"
	coh "github.com/oracle/coherence-operator/api/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
	"testing"
)

func TestJvmHeapSize(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceResourceSpec{
			JVM: &coh.JVMSpec{
				Memory: &coh.JvmMemorySpec{
					HeapSize: pointer.StringPtr("10g"),
				},
			},
		},
	}

	args := []string{"runner", "server"}
	env := EnvVarsFromDeployment(d)

	expectedCommand := GetJavaCommand()
	expectedArgs := append(GetMinimalExpectedArgs(), "-XX:InitialHeapSize=10g", "-XX:MaxHeapSize=10g")

	_, cmd, err := DryRun(args, env)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(cmd).NotTo(BeNil())

	g.Expect(cmd.Dir).To(Equal(TestAppDir))
	g.Expect(cmd.Path).To(Equal(expectedCommand))
	g.Expect(cmd.Args).To(ConsistOf(expectedArgs))
}

func TestJvmInitialHeapSize(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceResourceSpec{
			JVM: &coh.JVMSpec{
				Memory: &coh.JvmMemorySpec{
					InitialHeapSize: pointer.StringPtr("10g"),
				},
			},
		},
	}

	args := []string{"runner", "server"}
	env := EnvVarsFromDeployment(d)

	expectedCommand := GetJavaCommand()
	expectedArgs := append(GetMinimalExpectedArgs(), "-XX:InitialHeapSize=10g")

	_, cmd, err := DryRun(args, env)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(cmd).NotTo(BeNil())

	g.Expect(cmd.Dir).To(Equal(TestAppDir))
	g.Expect(cmd.Path).To(Equal(expectedCommand))
	g.Expect(cmd.Args).To(ConsistOf(expectedArgs))
}

func TestJvmMaxHeapSize(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceResourceSpec{
			JVM: &coh.JVMSpec{
				Memory: &coh.JvmMemorySpec{
					MaxHeapSize: pointer.StringPtr("10g"),
				},
			},
		},
	}

	args := []string{"runner", "server"}
	env := EnvVarsFromDeployment(d)

	expectedCommand := GetJavaCommand()
	expectedArgs := append(GetMinimalExpectedArgs(), "-XX:MaxHeapSize=10g")

	_, cmd, err := DryRun(args, env)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(cmd).NotTo(BeNil())

	g.Expect(cmd.Dir).To(Equal(TestAppDir))
	g.Expect(cmd.Path).To(Equal(expectedCommand))
	g.Expect(cmd.Args).To(ConsistOf(expectedArgs))
}

func TestJvmHeapSizeOverridesInitialAndMaxHeapSize(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceResourceSpec{
			JVM: &coh.JVMSpec{
				Memory: &coh.JvmMemorySpec{
					HeapSize:        pointer.StringPtr("5g"),
					InitialHeapSize: pointer.StringPtr("1g"),
					MaxHeapSize:     pointer.StringPtr("10g"),
				},
			},
		},
	}

	args := []string{"runner", "server"}
	env := EnvVarsFromDeployment(d)

	expectedCommand := GetJavaCommand()
	expectedArgs := append(GetMinimalExpectedArgs(), "-XX:InitialHeapSize=5g", "-XX:MaxHeapSize=5g")

	_, cmd, err := DryRun(args, env)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(cmd).NotTo(BeNil())

	g.Expect(cmd.Dir).To(Equal(TestAppDir))
	g.Expect(cmd.Path).To(Equal(expectedCommand))
	g.Expect(cmd.Args).To(ConsistOf(expectedArgs))
}

func TestJvmMaxRam(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceResourceSpec{
			JVM: &coh.JVMSpec{
				Memory: &coh.JvmMemorySpec{
					MaxRAM: pointer.StringPtr("10g"),
				},
			},
		},
	}

	args := []string{"runner", "server"}
	env := EnvVarsFromDeployment(d)

	expectedCommand := GetJavaCommand()
	expectedArgs := append(GetMinimalExpectedArgs(), "-XX:MaxRAM=10g")

	_, cmd, err := DryRun(args, env)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(cmd).NotTo(BeNil())

	g.Expect(cmd.Dir).To(Equal(TestAppDir))
	g.Expect(cmd.Path).To(Equal(expectedCommand))
	g.Expect(cmd.Args).To(ConsistOf(expectedArgs))
}

func TestJvmRamPercent(t *testing.T) {
	g := NewGomegaWithT(t)

	pct := resource.MustParse("5.5")
	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceResourceSpec{
			JVM: &coh.JVMSpec{
				Memory: &coh.JvmMemorySpec{
					Percentage: &pct,
				},
			},
		},
	}

	args := []string{"runner", "server"}
	env := EnvVarsFromDeployment(d)

	expectedCommand := GetJavaCommand()
	expectedArgs := append(GetMinimalExpectedArgs(), "-XX:InitialRAMPercentage=5.500",
		"-XX:MaxRAMPercentage=5.500", "-XX:MinRAMPercentage=5.500")

	_, cmd, err := DryRun(args, env)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(cmd).NotTo(BeNil())

	g.Expect(cmd.Dir).To(Equal(TestAppDir))
	g.Expect(cmd.Path).To(Equal(expectedCommand))
	g.Expect(cmd.Args).To(ConsistOf(expectedArgs))
}

func TestJvmInitialRamPercent(t *testing.T) {
	g := NewGomegaWithT(t)

	pct := resource.MustParse("5.5")
	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceResourceSpec{
			JVM: &coh.JVMSpec{
				Memory: &coh.JvmMemorySpec{
					InitialRAMPercentage: &pct,
				},
			},
		},
	}

	args := []string{"runner", "server"}
	env := EnvVarsFromDeployment(d)

	expectedCommand := GetJavaCommand()
	expectedArgs := append(GetMinimalExpectedArgs(), "-XX:InitialRAMPercentage=5.500")

	_, cmd, err := DryRun(args, env)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(cmd).NotTo(BeNil())

	g.Expect(cmd.Dir).To(Equal(TestAppDir))
	g.Expect(cmd.Path).To(Equal(expectedCommand))
	g.Expect(cmd.Args).To(ConsistOf(expectedArgs))
}

func TestJvmMaxRamPercent(t *testing.T) {
	g := NewGomegaWithT(t)

	pct := resource.MustParse("5.5")
	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceResourceSpec{
			JVM: &coh.JVMSpec{
				Memory: &coh.JvmMemorySpec{
					MaxRAMPercentage: &pct,
				},
			},
		},
	}

	args := []string{"runner", "server"}
	env := EnvVarsFromDeployment(d)

	expectedCommand := GetJavaCommand()
	expectedArgs := append(GetMinimalExpectedArgs(), "-XX:MaxRAMPercentage=5.500")

	_, cmd, err := DryRun(args, env)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(cmd).NotTo(BeNil())

	g.Expect(cmd.Dir).To(Equal(TestAppDir))
	g.Expect(cmd.Path).To(Equal(expectedCommand))
	g.Expect(cmd.Args).To(ConsistOf(expectedArgs))
}

func TestJvmMinRamPercent(t *testing.T) {
	g := NewGomegaWithT(t)

	pct := resource.MustParse("5.5")
	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceResourceSpec{
			JVM: &coh.JVMSpec{
				Memory: &coh.JvmMemorySpec{
					MinRAMPercentage: &pct,
				},
			},
		},
	}

	args := []string{"runner", "server"}
	env := EnvVarsFromDeployment(d)

	expectedCommand := GetJavaCommand()
	expectedArgs := append(GetMinimalExpectedArgs(), "-XX:MinRAMPercentage=5.500")

	_, cmd, err := DryRun(args, env)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(cmd).NotTo(BeNil())

	g.Expect(cmd.Dir).To(Equal(TestAppDir))
	g.Expect(cmd.Path).To(Equal(expectedCommand))
	g.Expect(cmd.Args).To(ConsistOf(expectedArgs))
}

func TestJvmRamPercentOverridesInitialMaxAndMin(t *testing.T) {
	g := NewGomegaWithT(t)

	pct := resource.MustParse("5.5")
	pctInit := resource.MustParse("1")
	pctMin := resource.MustParse("2")
	pctMax := resource.MustParse("10")
	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceResourceSpec{
			JVM: &coh.JVMSpec{
				Memory: &coh.JvmMemorySpec{
					Percentage:           &pct,
					InitialRAMPercentage: &pctInit,
					MinRAMPercentage:     &pctMin,
					MaxRAMPercentage:     &pctMax,
				},
			},
		},
	}

	args := []string{"runner", "server"}
	env := EnvVarsFromDeployment(d)

	expectedCommand := GetJavaCommand()
	expectedArgs := append(GetMinimalExpectedArgs(), "-XX:InitialRAMPercentage=5.500",
		"-XX:MaxRAMPercentage=5.500", "-XX:MinRAMPercentage=5.500")

	_, cmd, err := DryRun(args, env)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(cmd).NotTo(BeNil())

	g.Expect(cmd.Dir).To(Equal(TestAppDir))
	g.Expect(cmd.Path).To(Equal(expectedCommand))
	g.Expect(cmd.Args).To(ConsistOf(expectedArgs))
}

func TestJvmStackSize(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceResourceSpec{
			JVM: &coh.JVMSpec{
				Memory: &coh.JvmMemorySpec{
					StackSize: pointer.StringPtr("500k"),
				},
			},
		},
	}

	args := []string{"runner", "server"}
	env := EnvVarsFromDeployment(d)

	expectedCommand := GetJavaCommand()
	expectedArgs := append(GetMinimalExpectedArgs(), "-Xss500k")

	_, cmd, err := DryRun(args, env)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(cmd).NotTo(BeNil())

	g.Expect(cmd.Dir).To(Equal(TestAppDir))
	g.Expect(cmd.Path).To(Equal(expectedCommand))
	g.Expect(cmd.Args).To(ConsistOf(expectedArgs))
}

func TestJvmMetaspaceSize(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceResourceSpec{
			JVM: &coh.JVMSpec{
				Memory: &coh.JvmMemorySpec{
					MetaspaceSize: pointer.StringPtr("5g"),
				},
			},
		},
	}

	args := []string{"runner", "server"}
	env := EnvVarsFromDeployment(d)

	expectedCommand := GetJavaCommand()
	expectedArgs := append(GetMinimalExpectedArgs(), "-XX:MetaspaceSize=5g", "-XX:MaxMetaspaceSize=5g")

	_, cmd, err := DryRun(args, env)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(cmd).NotTo(BeNil())

	g.Expect(cmd.Dir).To(Equal(TestAppDir))
	g.Expect(cmd.Path).To(Equal(expectedCommand))
	g.Expect(cmd.Args).To(ConsistOf(expectedArgs))
}

func TestJvmDirectMemorySize(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceResourceSpec{
			JVM: &coh.JVMSpec{
				Memory: &coh.JvmMemorySpec{
					DirectMemorySize: pointer.StringPtr("5g"),
				},
			},
		},
	}

	args := []string{"runner", "server"}
	env := EnvVarsFromDeployment(d)

	expectedCommand := GetJavaCommand()
	expectedArgs := append(GetMinimalExpectedArgs(), "-XX:MaxDirectMemorySize=5g")

	_, cmd, err := DryRun(args, env)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(cmd).NotTo(BeNil())

	g.Expect(cmd.Dir).To(Equal(TestAppDir))
	g.Expect(cmd.Path).To(Equal(expectedCommand))
	g.Expect(cmd.Args).To(ConsistOf(expectedArgs))
}

func TestJvmNativeMemoryTracking(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceResourceSpec{
			JVM: &coh.JVMSpec{
				Memory: &coh.JvmMemorySpec{
					NativeMemoryTracking: pointer.StringPtr("detail"),
				},
			},
		},
	}

	args := []string{"runner", "server"}
	env := EnvVarsFromDeployment(d)

	expectedCommand := GetJavaCommand()
	expectedArgs := append(GetMinimalExpectedArgsWithoutPrefix("-XX:NativeMemoryTracking="), "-XX:NativeMemoryTracking=detail")

	_, cmd, err := DryRun(args, env)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(cmd).NotTo(BeNil())

	g.Expect(cmd.Dir).To(Equal(TestAppDir))
	g.Expect(cmd.Path).To(Equal(expectedCommand))
	g.Expect(cmd.Args).To(ConsistOf(expectedArgs))
}

func TestJvmOOMHeapDumpOff(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceResourceSpec{
			JVM: &coh.JVMSpec{
				Memory: &coh.JvmMemorySpec{
					OnOutOfMemory: &coh.JvmOutOfMemorySpec{
						HeapDump: pointer.BoolPtr(false),
					},
				},
			},
		},
	}

	args := []string{"runner", "server"}
	env := EnvVarsFromDeployment(d)

	expectedCommand := GetJavaCommand()
	expectedArgs := GetMinimalExpectedArgsWithoutPrefix("-XX:+HeapDumpOnOutOfMemoryError")

	_, cmd, err := DryRun(args, env)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(cmd).NotTo(BeNil())

	g.Expect(cmd.Dir).To(Equal(TestAppDir))
	g.Expect(cmd.Path).To(Equal(expectedCommand))
	g.Expect(cmd.Args).To(ConsistOf(expectedArgs))
}

func TestJvmOOMExitOff(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceResourceSpec{
			JVM: &coh.JVMSpec{
				Memory: &coh.JvmMemorySpec{
					OnOutOfMemory: &coh.JvmOutOfMemorySpec{
						Exit: pointer.BoolPtr(false),
					},
				},
			},
		},
	}

	args := []string{"runner", "server"}
	env := EnvVarsFromDeployment(d)

	expectedCommand := GetJavaCommand()
	expectedArgs := GetMinimalExpectedArgsWithoutPrefix("-XX:+ExitOnOutOfMemoryError")

	_, cmd, err := DryRun(args, env)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(cmd).NotTo(BeNil())

	g.Expect(cmd.Dir).To(Equal(TestAppDir))
	g.Expect(cmd.Path).To(Equal(expectedCommand))
	g.Expect(cmd.Args).To(ConsistOf(expectedArgs))
}
