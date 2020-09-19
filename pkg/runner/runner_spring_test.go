/*
 *  Copyright (c) 2020, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package runner

import (
	. "github.com/onsi/gomega"
	coh "github.com/oracle/coherence-operator/api/v1"
	"io/ioutil"
	corev1 "k8s.io/api/core/v1"
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

	args := []string{"runner", "server"}
	env := EnvVarsFromDeployment(d)

	expectedCommand := GetJavaCommand()
	expectedArgs := GetMinimalExpectedSpringBootArgs()

	_, cmd, err := DryRun(args, env)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(cmd).NotTo(BeNil())

	g.Expect(cmd.Dir).To(Equal(TestAppDir))
	g.Expect(cmd.Path).To(Equal(expectedCommand))
	g.Expect(cmd.Args).To(ConsistOf(expectedArgs))
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

	args := []string{"runner", "server"}
	env := EnvVarsFromDeployment(d)

	expectedCommand := GetJavaCommand()
	expectedArgs := GetMinimalExpectedSpringBootFatJarArgs(jar)

	_, cmd, err := DryRun(args, env)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(cmd).NotTo(BeNil())

	g.Expect(cmd.Dir).To(Equal(TestAppDir))
	g.Expect(cmd.Path).To(Equal(expectedCommand))
	g.Expect(cmd.Args).To(ConsistOf(expectedArgs))
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

	args := []string{"runner", "server"}
	env := EnvVarsFromDeployment(d)

	expectedCommand := GetJavaCommand()
	expectedArgs := append(GetMinimalExpectedSpringBootFatJarArgs(jar), "-Dloader.main=foo.Bar")

	_, cmd, err := DryRun(args, env)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(cmd).NotTo(BeNil())

	g.Expect(cmd.Dir).To(Equal(TestAppDir))
	g.Expect(cmd.Path).To(Equal(expectedCommand))
	g.Expect(cmd.Args).To(ConsistOf(expectedArgs))
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

	args := []string{"runner", "server"}
	env := EnvVarsFromDeployment(d)

	expectedCommand := getBuildpackLauncher()
	expectedArgs := []string{expectedCommand, "java", SpringBootMain}

	_, cmd, err := DryRun(args, env)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(cmd).NotTo(BeNil())

	g.Expect(cmd.Dir).To(Equal(""))
	g.Expect(cmd.Path).To(Equal(expectedCommand))
	g.Expect(cmd.Args).To(ConsistOf(expectedArgs))

	opts, found := findWithPrefix(cmd.Env, "JAVA_TOOL_OPTIONS=")
	g.Expect(found).To(BeTrue())
	g.Expect(opts).To(HavePrefix("JAVA_TOOL_OPTIONS=@"))

	fileName := opts[19:]
	data, err := ioutil.ReadFile(fileName)
	g.Expect(err).NotTo(HaveOccurred())

	actualOpts := strings.Split(string(data), "\n")
	expectedOpts := AppendCommonExpectedArgs([]string{"-Dloader.path=/coherence-operator/utils/lib/coherence-utils.jar,/coherence-operator/utils/config"})
	g.Expect(actualOpts).To(ConsistOf(expectedOpts))
}

func TestSpringBootBuildpacksWithJavaToolOptions(t *testing.T) {
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
			Env: []corev1.EnvVar{
				{
					Name:  "JAVA_TOOL_OPTIONS",
					Value: "-Dfoo=bar",
				},
			},
		},
	}

	args := []string{"runner", "server"}
	env := EnvVarsFromDeployment(d)

	expectedCommand := getBuildpackLauncher()
	expectedArgs := []string{expectedCommand, "java", SpringBootMain}

	_, cmd, err := DryRun(args, env)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(cmd).NotTo(BeNil())

	g.Expect(cmd.Dir).To(Equal(""))
	g.Expect(cmd.Path).To(Equal(expectedCommand))
	g.Expect(cmd.Args).To(ConsistOf(expectedArgs))

	opts, found := findWithPrefix(cmd.Env, "JAVA_TOOL_OPTIONS=")
	g.Expect(found).To(BeTrue())
	g.Expect(opts).To(HavePrefix("JAVA_TOOL_OPTIONS=@"))
	g.Expect(opts).To(HaveSuffix(" -Dfoo=bar"))
}

func findWithPrefix(args []string, prefix string) (string, bool) {
	for _, s := range args {
		if strings.HasPrefix(s, prefix) {
			return s, true
		}
	}
	return "", false
}

func GetMinimalExpectedSpringBootArgs() []string {
	args := []string{
		"java",
		"-Dloader.path=/coherence-operator/utils/lib/coherence-utils.jar,/coherence-operator/utils/config",
	}
	args = append(AppendCommonExpectedArgs(args), SpringBootMain)
	return args
}

func GetMinimalExpectedSpringBootFatJarArgs(jar string) []string {
	args := []string{
		"java",
		"-cp",
		jar,
		"-Dloader.path=/coherence-operator/utils/lib/coherence-utils.jar,/coherence-operator/utils/config",
	}

	return append(AppendCommonExpectedArgs(args), SpringBootMain)
}
