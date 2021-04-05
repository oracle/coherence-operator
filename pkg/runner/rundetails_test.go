/*
 * Copyright (c) 2020, 2021, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package runner

import (
	. "github.com/onsi/gomega"
	"testing"
)

func TestRunDetailsGetenvWhenMissing(t *testing.T) {
	g := NewGomegaWithT(t)

	r := RunDetails{}
	g.Expect(r.Getenv("foo")).To(Equal(""))
}

func TestRunDetailsGetenvWhenPresent(t *testing.T) {
	g := NewGomegaWithT(t)

	r := RunDetails{
		Env: map[string]string{"foo": "bar"},
	}
	g.Expect(r.Getenv("foo")).To(Equal("bar"))
}

func TestRunDetailsGetenvWithPrefixWhenMissing(t *testing.T) {
	g := NewGomegaWithT(t)

	r := RunDetails{}
	g.Expect(r.GetenvWithPrefix("foo", "_test")).To(Equal(""))
}

func TestRunDetailsGetenvWithPrefixWhenPresent(t *testing.T) {
	g := NewGomegaWithT(t)

	r := RunDetails{
		Env: map[string]string{"foo_test": "bar"},
	}
	g.Expect(r.GetenvWithPrefix("foo", "_test")).To(Equal("bar"))
}

func TestRunDetailsAddClasspath(t *testing.T) {
	g := NewGomegaWithT(t)

	r := RunDetails{}
	r.AddClasspath("foo")
	g.Expect(r.Classpath).To(Equal("foo"))
}

func TestRunDetailsAddClasspathWithExpansion(t *testing.T) {
	g := NewGomegaWithT(t)

	env := make(map[string]string)
	env["FOO"] = "foo-value"

	r := RunDetails{
		Env: env,
	}

	r.AddClasspath("${FOO}")
	g.Expect(r.Classpath).To(Equal("foo-value"))
}

func TestRunDetailsAddClasspathMultipleTimes(t *testing.T) {
	g := NewGomegaWithT(t)

	r := RunDetails{}
	r.AddClasspath("foo")
	r.AddClasspath("bar")
	g.Expect(r.Classpath).To(Equal("foo:bar"))
}

func TestRunDetailsAddClasspathEmptyString(t *testing.T) {
	g := NewGomegaWithT(t)

	r := RunDetails{}
	r.AddClasspath("foo")
	r.AddClasspath("")
	g.Expect(r.Classpath).To(Equal("foo"))
}

func TestRunDetailsAddToFrontClasspath(t *testing.T) {
	g := NewGomegaWithT(t)

	r := RunDetails{
		Classpath: "foo",
	}
	r.AddToFrontOfClasspath("bar")
	g.Expect(r.Classpath).To(Equal("bar:foo"))
}

func TestRunDetailsAddToFrontClasspathMultipleTimes(t *testing.T) {
	g := NewGomegaWithT(t)

	r := RunDetails{
		Classpath: "foo",
	}
	r.AddToFrontOfClasspath("bar1")
	r.AddToFrontOfClasspath("bar2")
	g.Expect(r.Classpath).To(Equal("bar2:bar1:foo"))
}

func TestRunDetailsAddToFrontOfClasspathEmptyString(t *testing.T) {
	g := NewGomegaWithT(t)

	r := RunDetails{
		Classpath: "foo",
	}
	r.AddToFrontOfClasspath("")
	g.Expect(r.Classpath).To(Equal("foo"))
}

func TestRunDetailsGetClasspath(t *testing.T) {
	g := NewGomegaWithT(t)

	r := RunDetails{}
	r.AddClasspath("foo")
	g.Expect(r.GetClasspath()).To(Equal("foo"))
}

func TestRunDetailsGetClasspathWithCoherenceHome(t *testing.T) {
	g := NewGomegaWithT(t)

	r := RunDetails{
		CoherenceHome: "/u01/oracle/coherence",
	}
	r.AddClasspath("foo")
	g.Expect(r.GetClasspath()).To(Equal("foo:/u01/oracle/coherence/conf:/u01/oracle/coherence/lib/coherence.jar"))
}

func TestRunDetailsGetJavaEmpty(t *testing.T) {
	g := NewGomegaWithT(t)

	r := RunDetails{}
	g.Expect(r.GetJavaExecutable()).To(Equal("java"))
}

func TestRunDetailsGetJavaWhenJavaHomeSet(t *testing.T) {
	g := NewGomegaWithT(t)

	r := RunDetails{
		JavaHome: "/local/bin/jdk11",
	}
	g.Expect(r.GetJavaExecutable()).To(Equal("/local/bin/jdk11/bin/java"))
}

func TestRunDetailsGetCommandWhenEmpty(t *testing.T) {
	g := NewGomegaWithT(t)

	r := RunDetails{}
	var expected []string
	g.Expect(r.GetCommand()).To(Equal(expected))
}

func TestExpandEnv(t *testing.T) {
	g := NewGomegaWithT(t)

	env := make(map[string]string)
	env["A"] = "value-a"
	env["B"] = "value-b"
	env["C"] = "value-c"

	r := RunDetails{}
	result := r.Expand("$(A) ${B} $C", func(s string) string {
		return env[s]
	})

	g.Expect(result).To(Equal("value-a value-b value-c"))
}
