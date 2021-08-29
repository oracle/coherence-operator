/*
 * Copyright (c) 2020, 2021, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package runner

import (
	. "github.com/onsi/gomega"
	coh "github.com/oracle/coherence-operator/api/v1"
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
	"strings"
	"testing"
)

func TestSpringBootApplication(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceResourceSpec{
			Application: &coh.ApplicationSpec{
				Type: pointer.StringPtr(AppTypeSpring),
			},
		},
	}

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(d)

	expectedCommand := GetJavaCommand()
	expectedArgs := GetMinimalExpectedSpringBootArgs()

	e, err := ExecuteWithArgs(env, args)
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
		Spec: coh.CoherenceResourceSpec{
			Application: &coh.ApplicationSpec{
				Type:             pointer.StringPtr(AppTypeSpring),
				SpringBootFatJar: &jar,
			},
		},
	}

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(d)

	expectedCommand := GetJavaCommand()
	expectedArgs := GetMinimalExpectedSpringBootFatJarArgs(jar)

	e, err := ExecuteWithArgs(env, args)
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
		Spec: coh.CoherenceResourceSpec{
			Application: &coh.ApplicationSpec{
				Type:             pointer.StringPtr(AppTypeSpring),
				SpringBootFatJar: &jar,
				Main:             pointer.StringPtr("foo.Bar"),
			},
		},
	}

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(d)

	expectedCommand := GetJavaCommand()
	expectedArgs := append(GetMinimalExpectedSpringBootFatJarArgs(jar), "-Dloader.main=foo.Bar")

	e, err := ExecuteWithArgs(env, args)
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
		Spec: coh.CoherenceResourceSpec{
			Application: &coh.ApplicationSpec{
				Type: pointer.StringPtr(AppTypeSpring),
				CloudNativeBuildPack: &coh.CloudNativeBuildPackSpec{
					Enabled: pointer.BoolPtr(true),
				},
			},
		},
	}

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(d)

	expectedCommand := getBuildpackLauncher()

	e, err := ExecuteWithArgs(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(""))
	g.Expect(e.OsCmd.Path).To(Equal(expectedCommand))

	g.Expect(len(e.OsCmd.Args)).To(Equal(4))
	g.Expect(e.OsCmd.Args[0]).To(Equal(coh.DefaultCnbpLauncher))
	g.Expect(e.OsCmd.Args[1]).To(Equal("java"))
	g.Expect(e.OsCmd.Args[3]).To(Equal(SpringBootMain))

	g.Expect(e.OsCmd.Args[2]).To(HavePrefix("@"))
	fileName := e.OsCmd.Args[2][1:]
	data, err := ioutil.ReadFile(fileName)
	g.Expect(err).NotTo(HaveOccurred())

	actualOpts := strings.Split(string(data), "\n")
	expectedOpts := AppendCommonExpectedArgs([]string{"-Dloader.path=/coherence-operator/utils/lib/coherence-operator.jar,/coherence-operator/utils/config"})
	g.Expect(actualOpts).To(ConsistOf(expectedOpts))
}

func GetMinimalExpectedSpringBootArgs() []string {
	args := []string{
		GetJavaCommand(),
		"-Dloader.path=/coherence-operator/utils/lib/coherence-operator.jar,/coherence-operator/utils/config",
	}
	args = append(AppendCommonExpectedArgs(args), SpringBootMain)
	return args
}

func GetMinimalExpectedSpringBootFatJarArgs(jar string) []string {
	args := []string{
		GetJavaCommand(),
		"-cp",
		jar,
		"-Dloader.path=/coherence-operator/utils/lib/coherence-operator.jar,/coherence-operator/utils/config",
	}

	return append(AppendCommonExpectedArgs(args), SpringBootMain)
}
