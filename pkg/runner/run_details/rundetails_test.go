/*
 * Copyright (c) 2020, 2025, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package run_details

import (
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
	ctrl "sigs.k8s.io/controller-runtime"
	"testing"
)

var (
	testLog = ctrl.Log.WithName("test")
)

func TestRunDetailsGetenvWhenMissing(t *testing.T) {
	g := NewGomegaWithT(t)

	r := RunDetails{}
	g.Expect(r.Getenv("foo")).To(Equal(""))
}

func TestRunDetailsGetenvWhenPresent(t *testing.T) {
	g := NewGomegaWithT(t)

	v := viper.New()
	v.Set("foo", "bar")

	r := NewRunDetails(v, testLog)
	g.Expect(r.Getenv("foo")).To(Equal("bar"))
}

func TestRunDetailsGetenvWithPrefixWhenMissing(t *testing.T) {
	g := NewGomegaWithT(t)

	r := NewRunDetails(viper.New(), testLog)
	g.Expect(r.GetenvWithPrefix("foo", "_test")).To(Equal(""))
}

func TestRunDetailsGetenvWithPrefixWhenPresent(t *testing.T) {
	g := NewGomegaWithT(t)

	v := viper.New()
	v.Set("foo_test", "bar")

	r := NewRunDetails(v, testLog)
	g.Expect(r.GetenvWithPrefix("foo", "_test")).To(Equal("bar"))
}

func TestRunDetailsAddClasspath(t *testing.T) {
	g := NewGomegaWithT(t)

	r := NewRunDetails(viper.New(), testLog)
	r.AddClasspath("foo")
	g.Expect(r.Classpath).To(Equal("foo"))
}

func TestRunDetailsAddClasspathWithExpansion(t *testing.T) {
	g := NewGomegaWithT(t)

	v := viper.New()
	v.Set("FOO", "foo-value")

	r := NewRunDetails(v, testLog)

	r.AddClasspath("${FOO}")
	g.Expect(r.Classpath).To(Equal("foo-value"))
}

func TestRunDetailsAddClasspathMultipleTimes(t *testing.T) {
	g := NewGomegaWithT(t)

	r := NewRunDetails(viper.New(), testLog)
	r.AddClasspath("foo")
	r.AddClasspath("bar")
	g.Expect(r.Classpath).To(Equal("foo:bar"))
}

func TestRunDetailsAddClasspathEmptyString(t *testing.T) {
	g := NewGomegaWithT(t)

	r := NewRunDetails(viper.New(), testLog)
	r.AddClasspath("foo")
	r.AddClasspath("")
	g.Expect(r.Classpath).To(Equal("foo"))
}

func TestRunDetailsAddToFrontClasspath(t *testing.T) {
	g := NewGomegaWithT(t)

	r := NewRunDetails(viper.New(), testLog)
	r.Classpath = "foo"

	r.AddToFrontOfClasspath("bar")
	g.Expect(r.Classpath).To(Equal("bar:foo"))
}

func TestRunDetailsAddToFrontClasspathMultipleTimes(t *testing.T) {
	g := NewGomegaWithT(t)

	r := NewRunDetails(viper.New(), testLog)
	r.Classpath = "foo"
	r.AddToFrontOfClasspath("bar1")
	r.AddToFrontOfClasspath("bar2")
	g.Expect(r.Classpath).To(Equal("bar2:bar1:foo"))
}

func TestRunDetailsAddToFrontOfClasspathEmptyString(t *testing.T) {
	g := NewGomegaWithT(t)

	r := NewRunDetails(viper.New(), testLog)
	r.Classpath = "foo"
	r.AddToFrontOfClasspath("")
	g.Expect(r.Classpath).To(Equal("foo"))
}

func TestRunDetailsGetClasspath(t *testing.T) {
	g := NewGomegaWithT(t)

	r := NewRunDetails(viper.New(), testLog)
	r.AddClasspath("foo")
	g.Expect(r.GetClasspath()).To(Equal("foo"))
}

func TestRunDetailsGetJavaEmpty(t *testing.T) {
	g := NewGomegaWithT(t)

	r := NewRunDetails(viper.New(), testLog)
	g.Expect(r.GetJavaExecutable()).To(Equal("java"))
}

func TestRunDetailsGetJavaWhenJavaHomeSet(t *testing.T) {
	g := NewGomegaWithT(t)

	r := NewRunDetails(viper.New(), testLog)
	r.JavaHome = "/local/bin/jdk11"
	g.Expect(r.GetJavaExecutable()).To(Equal("/local/bin/jdk11/bin/java"))
}

func TestRunDetailsGetCommandWhenEmpty(t *testing.T) {
	g := NewGomegaWithT(t)

	r := NewRunDetails(viper.New(), testLog)
	var expected []string
	g.Expect(r.GetCommand()).To(Equal(expected))
}

func TestExpandEnv(t *testing.T) {
	g := NewGomegaWithT(t)

	env := make(map[string]string)
	env["A"] = "value-a"
	env["B"] = "value-b"
	env["C"] = "value-c"

	r := NewRunDetails(viper.New(), testLog)
	result := r.Expand("$(A) ${B} $C", func(s string) string {
		return env[s]
	})

	g.Expect(result).To(Equal("value-a value-b value-c"))
}
