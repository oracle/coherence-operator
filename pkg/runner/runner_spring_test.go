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
	"testing"
)

func TestSpringBootApplication(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Application: &coh.ApplicationSpec{
					Type: ptr.To(coh.AppTypeSpring2),
				},
			},
		},
	}

	wd := ensureTestUtilsDir(t)
	expectedCP := GetOperatorClasspathWithUtilsDir(wd)
	expectedFileArgs := GetExpectedArgsFileContent()
	verifyConfigFilesWithArgsAndClasspath(t, d, expectedFileArgs, expectedCP)

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(t, d)

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(GetJavaCommand()))
	g.Expect(e.OsCmd.Args).To(ConsistOf(GetMinimalExpectedSpringBootArgs(t, coh.SpringBootMain2)))
}

func TestSpringBoot3Application(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Application: &coh.ApplicationSpec{
					Type: ptr.To(coh.AppTypeSpring3),
				},
			},
		},
	}

	wd := ensureTestUtilsDir(t)
	expectedCP := GetOperatorClasspathWithUtilsDir(wd)
	expectedFileArgs := GetExpectedArgsFileContent()
	verifyConfigFilesWithArgsAndClasspath(t, d, expectedFileArgs, expectedCP)

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(t, d)

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(GetJavaCommand()))
	g.Expect(e.OsCmd.Args).To(ConsistOf(GetMinimalExpectedSpringBootArgs(t, coh.SpringBootMain3)))
}

func TestSpringBootFatJarApplication(t *testing.T) {
	g := NewGomegaWithT(t)

	jar := "/apps/lib/foo.jar"
	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Application: &coh.ApplicationSpec{
					Type:             ptr.To(coh.AppTypeSpring2),
					SpringBootFatJar: &jar,
				},
			},
		},
	}

	wd := ensureTestUtilsDir(t)
	expectedCP := GetOperatorClasspathWithUtilsDir(wd)
	expectedFileArgs := GetExpectedArgsFileContent()
	verifyConfigFilesWithArgsAndClasspath(t, d, expectedFileArgs, expectedCP)

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(t, d)

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(GetJavaCommand()))
	expected := append(GetMinimalExpectedSpringBootArgs(t, coh.SpringBootMain2), "--class-path", "/apps/lib/foo.jar")
	g.Expect(e.OsCmd.Args).To(ConsistOf(expected))
}

func TestSpringBoot3FatJarApplication(t *testing.T) {
	g := NewGomegaWithT(t)

	jar := "/apps/lib/foo.jar"
	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Application: &coh.ApplicationSpec{
					Type:             ptr.To(coh.AppTypeSpring3),
					SpringBootFatJar: &jar,
				},
			},
		},
	}

	wd := ensureTestUtilsDir(t)
	expectedCP := GetOperatorClasspathWithUtilsDir(wd)
	expectedFileArgs := GetExpectedArgsFileContent()
	verifyConfigFilesWithArgsAndClasspath(t, d, expectedFileArgs, expectedCP)

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(t, d)

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(GetJavaCommand()))
	expected := append(GetMinimalExpectedSpringBootArgs(t, coh.SpringBootMain3), "--class-path", "/apps/lib/foo.jar")
	g.Expect(e.OsCmd.Args).To(ConsistOf(expected))
}

func TestSpringBootFatJarConsole(t *testing.T) {
	g := NewGomegaWithT(t)

	jar := "/apps/lib/foo.jar"
	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Application: &coh.ApplicationSpec{
					Type:             ptr.To(coh.AppTypeSpring2),
					SpringBootFatJar: &jar,
				},
			},
		},
	}

	wd := ensureTestUtilsDir(t)
	expectedCP := GetOperatorClasspathWithUtilsDir(wd)
	expectedFileArgs := GetExpectedArgsFileContent()
	verifyConfigFilesWithArgsAndClasspath(t, d, expectedFileArgs, expectedCP)

	args := []string{"console", "--dry-run"}
	env := EnvVarsFromDeployment(t, d)

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(GetJavaCommand()))

	expectedArgs := append(GetMinimalExpectedSpringBootArgs(t, coh.SpringBootMain2), "--class-path", jar)
	expectedArgs = append(expectedArgs, "-Dcoherence.role=console",
		"-Dloader.main=com.tangosol.net.CacheFactory",
		"-Dcoherence.distributed.localstorage=false",
		"-Dcoherence.localport.adjust=true",
		"-Dcoherence.management.http=none",
		"-Dcoherence.management.http.port=0",
		"-Dcoherence.metrics.http.enabled=false",
		"-Dcoherence.metrics.http.port=0",
		"-Dcoherence.operator.health.enabled=false",
		"-Dcoherence.health.http.port=0",
		"-Dcoherence.grpc.enabled=false",
		"-XX:NativeMemoryTracking=off")

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
					Type:             ptr.To(coh.AppTypeSpring2),
					SpringBootFatJar: &jar,
				},
			},
		},
	}

	wd := ensureTestUtilsDir(t)
	expectedCP := GetOperatorClasspathWithUtilsDir(wd)
	expectedFileArgs := GetExpectedArgsFileContent()
	verifyConfigFilesWithArgsAndClasspath(t, d, expectedFileArgs, expectedCP)

	args := []string{"console", "--dry-run", "--", "foo", "bar"}
	env := EnvVarsFromDeployment(t, d)

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(GetJavaCommand()))

	expectedArgs := append(GetMinimalExpectedSpringBootArgs(t, coh.SpringBootMain2), "--class-path", jar)
	expectedArgs = append(expectedArgs, "-Dcoherence.role=console",
		"-Dloader.main=com.tangosol.net.CacheFactory",
		"-Dcoherence.distributed.localstorage=false",
		"-Dcoherence.localport.adjust=true",
		"-Dcoherence.management.http=none",
		"-Dcoherence.management.http.port=0",
		"-Dcoherence.metrics.http.enabled=false",
		"-Dcoherence.metrics.http.port=0",
		"-Dcoherence.operator.health.enabled=false",
		"-Dcoherence.health.http.port=0",
		"-Dcoherence.grpc.enabled=false",
		"-XX:NativeMemoryTracking=off",
		"foo", "bar")

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
					Type:             ptr.To(coh.AppTypeSpring2),
					SpringBootFatJar: &jar,
					Main:             ptr.To("foo.Bar"),
				},
			},
		},
	}

	wd := ensureTestUtilsDir(t)
	expectedCP := GetOperatorClasspathWithUtilsDir(wd)
	expectedFileArgs := GetExpectedArgsFileContent()
	verifyConfigFilesWithArgsAndClasspath(t, d, expectedFileArgs, expectedCP)

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(t, d)

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(GetJavaCommand()))

	expectedArgs := append(GetMinimalExpectedSpringBootArgs(t, coh.SpringBootMain2),
		"--class-path", jar, "-Dloader.main=foo.Bar")

	g.Expect(e.OsCmd.Args).To(ConsistOf(expectedArgs))
}

func TestSpringBootBuildpacks(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Application: &coh.ApplicationSpec{
					Type: ptr.To(coh.AppTypeSpring2),
					CloudNativeBuildPack: &coh.CloudNativeBuildPackSpec{
						Enabled: ptr.To(true),
					},
				},
			},
		},
	}

	wd := ensureTestUtilsDir(t)
	expectedCP := GetOperatorClasspathWithUtilsDir(wd)
	expectedFileArgs := GetExpectedArgsFileContent()
	verifyConfigFilesWithArgsAndClasspath(t, d, expectedFileArgs, expectedCP)

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(t, d)

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(""))
	g.Expect(e.OsCmd.Path).To(Equal(getBuildpackLauncher()))

	g.Expect(len(e.OsCmd.Args)).To(Equal(4))
	g.Expect(e.OsCmd.Args[0]).To(Equal(coh.DefaultCnbpLauncher))
	g.Expect(e.OsCmd.Args[1]).To(Equal("java"))
	g.Expect(e.OsCmd.Args[3]).To(Equal(coh.SpringBootMain2))
}

func GetMinimalExpectedSpringBootArgs(t *testing.T, main string) []string {
	utils := ensureTestUtilsDir(t)

	cp := utils + "/lib/coherence-operator.jar"
	cfg := utils + "config"
	if _, err := os.Stat(cfg); err == nil {
		cp = cp + "," + cfg
	}
	args := []string{
		GetJavaArg(),
		"-Dloader.path=" + cp,
	}
	return append(AppendCommonExpectedArgs(args), main)
}
