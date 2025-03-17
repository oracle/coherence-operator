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
	"os"
	"testing"
)

func TestApplicationArgsEmpty(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Application: &coh.ApplicationSpec{
					Args: []string{},
				},
			},
		},
	}

	verifyConfigFilesWithArgs(t, d, GetExpectedArgsFileContent())

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(t, d)

	expectedCommand := GetJavaCommand()

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(expectedCommand))
	g.Expect(e.OsCmd.Args).To(ConsistOf(GetMinimalExpectedArgs(t)))
}

func TestApplicationArgs(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Application: &coh.ApplicationSpec{
					Args: []string{"Foo", "Bar"},
				},
			},
		},
	}

	verifyConfigFilesWithArgs(t, d, GetExpectedArgsFileContent())

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(t, d)

	expectedCommand := GetJavaCommand()

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(expectedCommand))
	g.Expect(e.OsCmd.Args).To(ConsistOf(GetMinimalExpectedArgsWith(t, "Foo", "Bar")))
}

func TestApplicationArgsWithEnvVarExpansion(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Application: &coh.ApplicationSpec{
					Args: []string{"${FOO}", "$BAR"},
				},
				Env: []corev1.EnvVar{
					{Name: "FOO", Value: "foo-value"},
					{Name: "BAR", Value: "bar-value"},
				},
			},
		},
	}

	verifyConfigFilesWithArgs(t, d, GetExpectedArgsFileContent())

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(t, d)

	expectedCommand := GetJavaCommand()

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(expectedCommand))
	g.Expect(e.OsCmd.Args).To(ConsistOf(GetMinimalExpectedArgsWith(t, "foo-value", "bar-value")))
}

func TestApplicationMain(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Application: &coh.ApplicationSpec{
					Main: ptr.To("com.oracle.test.Main"),
				},
			},
		},
	}

	verifyConfigFilesWithArgs(t, d, GetExpectedArgsFileContent())

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(t, d)

	expectedCommand := GetJavaCommand()

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(expectedCommand))
	g.Expect(e.OsCmd.Args).To(ConsistOf(GetMinimalExpectedArgsWithMainClass(t, "com.oracle.test.Main")))
}

func TestApplicationWorkingDirectory(t *testing.T) {
	g := NewGomegaWithT(t)

	utils := ensureTestUtilsDir(t)
	wd := utils + "/foo"
	err := os.MkdirAll(wd, os.ModePerm)
	g.Expect(err).NotTo(HaveOccurred())

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Application: &coh.ApplicationSpec{
					WorkingDir: &wd,
				},
			},
		},
	}

	expectedCP := GetOperatorClasspathWithUtilsDir(utils)
	expectedArgs := GetExpectedArgsFileContent()
	verifyConfigFilesWithArgsAndClasspath(t, d, expectedArgs, expectedCP)

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(t, d)

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(wd))
	g.Expect(e.OsCmd.Path).To(Equal(GetJavaCommand()))
	g.Expect(e.OsCmd.Args).To(ConsistOf(GetMinimalExpectedArgsWithWorkingDir(t, wd)))
}
