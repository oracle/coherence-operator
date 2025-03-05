/*
 * Copyright (c) 2020, 2025, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package runner

import (
	. "github.com/onsi/gomega"
	coh "github.com/oracle/coherence-operator/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	"os"
	"strings"
	"testing"
)

func TestSpringBootApplication(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Application: &coh.ApplicationSpec{
					Type: ptr.To(AppTypeSpring2),
				},
			},
		},
	}

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(d)

	expectedCommand := GetJavaCommand()
	expectedArgs := GetMinimalExpectedSpringBootArgs(SpringBootMain2)

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(expectedCommand))
	g.Expect(e.OsCmd.Args).To(ConsistOf(expectedArgs))
}

func TestSpringBoot3Application(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Application: &coh.ApplicationSpec{
					Type: ptr.To(AppTypeSpring3),
				},
			},
		},
	}

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(d)

	expectedCommand := GetJavaCommand()
	expectedArgs := GetMinimalExpectedSpringBootArgs(SpringBootMain3)

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(expectedCommand))
	g.Expect(e.OsCmd.Args).To(ConsistOf(expectedArgs))
}

func TestSpringBootFatJarApplication(t *testing.T) {
	g := NewGomegaWithT(t)

	jar := "/apps/lib/foo.jar"
	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Application: &coh.ApplicationSpec{
					Type:             ptr.To(AppTypeSpring2),
					SpringBootFatJar: &jar,
				},
			},
		},
	}

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(d)

	expectedCommand := GetJavaCommand()
	expectedArgs := GetMinimalExpectedSpringBootFatJarArgs(jar, SpringBootMain2)

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(expectedCommand))
	g.Expect(e.OsCmd.Args).To(ConsistOf(expectedArgs))
}

func TestSpringBoot3FatJarApplication(t *testing.T) {
	g := NewGomegaWithT(t)

	jar := "/apps/lib/foo.jar"
	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Application: &coh.ApplicationSpec{
					Type:             ptr.To(AppTypeSpring3),
					SpringBootFatJar: &jar,
				},
			},
		},
	}

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(d)

	expectedCommand := GetJavaCommand()
	expectedArgs := GetMinimalExpectedSpringBootFatJarArgs(jar, SpringBootMain3)

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(expectedCommand))
	g.Expect(e.OsCmd.Args).To(ConsistOf(expectedArgs))
}

func TestSpringBootFatJarConsole(t *testing.T) {
	g := NewGomegaWithT(t)

	jar := "/apps/lib/foo.jar"
	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Application: &coh.ApplicationSpec{
					Type:             ptr.To(AppTypeSpring2),
					SpringBootFatJar: &jar,
				},
			},
		},
	}

	args := []string{"console", "--dry-run"}
	env := EnvVarsFromDeployment(d)

	expectedCommand := GetJavaCommand()
	expectedArgs := GetMinimalExpectedSpringBootFatJarArgsForRole(jar, ConsoleMain, "")
	expectedArgs = ReplaceArg(expectedArgs, "-XX:NativeMemoryTracking=summary", "-XX:NativeMemoryTracking=off")
	expectedArgs = ReplaceArg(expectedArgs, "-Dcoherence.health.http.port=6676", "-Dcoherence.health.http.port=0")
	expectedArgs = RemoveArg(expectedArgs, "-Dcoherence.k8s.operator.health.port=6676")
	expectedArgs = append(expectedArgs, "-Dcoherence.localport.adjust=true",
		"-Dcoherence.metrics.http.enabled=false",
		"-Dcoherence.management.http=none",
		"-Dcoherence.management.http.port=0",
		"-Dcoherence.metrics.http.port=0",
		"-Dcoherence.k8s.operator.health.enabled=false",
		"-Dcoherence.health.http.port=0",
		"-Dcoherence.grpc.enabled=false",
		"-Dcoherence.k8s.operator.health.port=0")

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(expectedCommand))
	g.Expect(e.OsCmd.Args).To(ConsistOf(expectedArgs))
}

func TestSpringBootFatJarConsoleWithArgs(t *testing.T) {
	g := NewGomegaWithT(t)

	jar := "/apps/lib/foo.jar"
	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Application: &coh.ApplicationSpec{
					Type:             ptr.To(AppTypeSpring2),
					SpringBootFatJar: &jar,
				},
			},
		},
	}

	args := []string{"console", "--dry-run", "--", "foo", "bar"}
	env := EnvVarsFromDeployment(d)

	expectedCommand := GetJavaCommand()
	expectedArgs := GetMinimalExpectedSpringBootFatJarArgsForRole(jar, ConsoleMain, "")
	expectedArgs = ReplaceArg(expectedArgs, "-XX:NativeMemoryTracking=summary", "-XX:NativeMemoryTracking=off")
	expectedArgs = ReplaceArg(expectedArgs, "-Dcoherence.health.http.port=6676", "-Dcoherence.health.http.port=0")
	expectedArgs = RemoveArg(expectedArgs, "-Dcoherence.k8s.operator.health.port=6676")
	expectedArgs = append(expectedArgs, "-Dcoherence.localport.adjust=true",
		"-Dcoherence.metrics.http.enabled=false",
		"-Dcoherence.management.http=none",
		"-Dcoherence.management.http.port=0",
		"-Dcoherence.metrics.http.port=0",
		"-Dcoherence.k8s.operator.health.enabled=false",
		"-Dcoherence.grpc.enabled=false",
		"-Dcoherence.k8s.operator.health.port=0",
		"-Dcoherence.health.http.port=0",
		"foo", "bar")

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(expectedCommand))
	g.Expect(e.OsCmd.Args).To(ConsistOf(expectedArgs))
}

func TestSpringBootFatJarApplicationWithCustomMain(t *testing.T) {
	g := NewGomegaWithT(t)

	jar := "/apps/lib/foo.jar"
	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Application: &coh.ApplicationSpec{
					Type:             ptr.To(AppTypeSpring2),
					SpringBootFatJar: &jar,
					Main:             ptr.To("foo.Bar"),
				},
			},
		},
	}

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(d)

	expectedCommand := GetJavaCommand()
	expectedArgs := append(GetMinimalExpectedSpringBootFatJarArgs(jar, SpringBootMain2), "-Dloader.main=foo.Bar")

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(expectedCommand))
	g.Expect(e.OsCmd.Args).To(ConsistOf(expectedArgs))
}

func TestSpringBootBuildpacks(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Application: &coh.ApplicationSpec{
					Type: ptr.To(AppTypeSpring2),
					CloudNativeBuildPack: &coh.CloudNativeBuildPackSpec{
						Enabled: ptr.To(true),
					},
				},
			},
		},
	}

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(d)

	expectedCommand := getBuildpackLauncher()

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(""))
	g.Expect(e.OsCmd.Path).To(Equal(expectedCommand))

	g.Expect(len(e.OsCmd.Args)).To(Equal(4))
	g.Expect(e.OsCmd.Args[0]).To(Equal(coh.DefaultCnbpLauncher))
	g.Expect(e.OsCmd.Args[1]).To(Equal("java"))
	g.Expect(e.OsCmd.Args[3]).To(Equal(SpringBootMain2))

	g.Expect(e.OsCmd.Args[2]).To(HavePrefix("@"))
	fileName := e.OsCmd.Args[2][1:]
	data, err := os.ReadFile(fileName)
	g.Expect(err).NotTo(HaveOccurred())

	cp := "/coherence-operator/utils/lib/coherence-operator.jar"
	if _, err := os.Stat("/coherence-operator/utils/config"); err == nil {
		cp = cp + ":/coherence-operator/utils/config"
	}

	actualOpts := strings.Split(string(data), "\n")
	expectedOpts := AppendCommonExpectedArgs([]string{"-Dloader.path=" + cp})
	g.Expect(actualOpts).To(ConsistOf(expectedOpts))
}

func GetMinimalExpectedSpringBootArgs(main string) []string {
	cp := "/coherence-operator/utils/lib/coherence-operator.jar"
	if _, err := os.Stat("/coherence-operator/utils/config"); err == nil {
		cp = cp + ",/coherence-operator/utils/config"
	}
	args := []string{
		GetJavaArg(),
		"-Dloader.path=" + cp,
	}
	args = append(AppendCommonExpectedArgs(args), main)
	return args
}

func GetMinimalExpectedSpringBootFatJarArgs(jar, main string) []string {
	return GetMinimalExpectedSpringBootFatJarArgsWithMain(jar, main, "")
}

func GetMinimalExpectedSpringBootFatJarArgsWithMain(jar, springMain, main string) []string {
	cp := "/coherence-operator/utils/lib/coherence-operator.jar"
	if _, err := os.Stat("/coherence-operator/utils/config"); err == nil {
		cp = cp + ",/coherence-operator/utils/config"
	}
	args := []string{
		GetJavaArg(),
		"--class-path",
		jar,
		"-Dloader.path=" + cp,
	}

	if main != "" {
		args = append(args, "-Dloader.main="+main)
	}

	return append(AppendCommonExpectedArgs(args), springMain)
}

func GetMinimalExpectedSpringBootFatJarArgsForRole(jar, main, role string) []string {
	cp := "/coherence-operator/utils/lib/coherence-operator.jar"
	if _, err := os.Stat("/coherence-operator/utils/config"); err == nil {
		cp = cp + ",/coherence-operator/utils/config"
	}
	args := []string{
		GetJavaArg(),
		"--class-path",
		jar,
		"-Dloader.path=" + cp,
		"-Dcoherence.distributed.localstorage=false",
	}

	if main != "" {
		args = append(args, "-Dloader.main="+main)
	}

	return append(AppendCommonExpectedNonServerArgs(args, role), SpringBootMain2)
}
