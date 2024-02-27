/*
 * Copyright (c) 2020, 2024, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package runner

import (
	. "github.com/onsi/gomega"
	coh "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/pkg/operator"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestBasicOperator(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
	}

	args := []string{"operator", "--dry-run"}
	env := EnvVarsFromDeployment(d)

	e, err := ExecuteWithArgs(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())

	l := operator.GetGlobalLabelsNoError()
	g.Expect(l).NotTo(BeNil())
	g.Expect(len(l)).To(Equal(0))

	a := operator.GetGlobalAnnotationsNoError()
	g.Expect(a).NotTo(BeNil())
	g.Expect(len(a)).To(Equal(0))
}

func TestOperatorWithSingleGlobalLabel(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
	}

	args := []string{"operator", "--dry-run", "--global-label", "one=value-one"}
	env := EnvVarsFromDeployment(d)

	e, err := ExecuteWithArgs(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())

	l := operator.GetGlobalLabelsNoError()
	g.Expect(l).NotTo(BeNil())
	g.Expect(len(l)).To(Equal(1))
	g.Expect(l["one"]).To(Equal("value-one"))
}

func TestOperatorWithMultipleGlobalLabels(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
	}

	args := []string{"operator", "--dry-run",
		"--global-label", "one=value-one",
		"--global-label", "two=value-two",
		"--global-label", "three=value-three",
	}
	env := EnvVarsFromDeployment(d)

	e, err := ExecuteWithArgs(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())

	l := operator.GetGlobalLabelsNoError()
	g.Expect(l).NotTo(BeNil())
	g.Expect(len(l)).To(Equal(3))
	g.Expect(l["one"]).To(Equal("value-one"))
	g.Expect(l["two"]).To(Equal("value-two"))
	g.Expect(l["three"]).To(Equal("value-three"))
}

func TestOperatorWithSingleGlobalAnnotation(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
	}

	args := []string{"operator", "--dry-run", "--global-annotation", "one=value-one"}
	env := EnvVarsFromDeployment(d)

	e, err := ExecuteWithArgs(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())

	l := operator.GetGlobalAnnotationsNoError()
	g.Expect(l).NotTo(BeNil())
	g.Expect(len(l)).To(Equal(1))
	g.Expect(l["one"]).To(Equal("value-one"))
}

func TestOperatorWithMultipleGlobalAnnotations(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
	}

	args := []string{"operator", "--dry-run",
		"--global-annotation", "one=value-one",
		"--global-annotation", "two=value-two",
		"--global-annotation", "three=value-three",
	}
	env := EnvVarsFromDeployment(d)

	e, err := ExecuteWithArgs(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())

	l := operator.GetGlobalAnnotationsNoError()
	g.Expect(l).NotTo(BeNil())
	g.Expect(len(l)).To(Equal(3))
	g.Expect(l["one"]).To(Equal("value-one"))
	g.Expect(l["two"]).To(Equal("value-two"))
	g.Expect(l["three"]).To(Equal("value-three"))
}
