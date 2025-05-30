/*
 * Copyright (c) 2020, 2025, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package runner

import (
	"crypto/tls"
	. "github.com/onsi/gomega"
	coh "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/pkg/operator"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"strings"
	"testing"
)

var (
	testLog = ctrl.Log.WithName("test")
)

func TestBasicOperator(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
	}

	args := []string{"operator", "--dry-run"}
	env := EnvVarsFromDeployment(t, d)

	e, err := ExecuteWithArgsAndNewViper(env, args)
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
	env := EnvVarsFromDeployment(t, d)

	e, err := ExecuteWithArgsAndNewViper(env, args)
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
	env := EnvVarsFromDeployment(t, d)

	e, err := ExecuteWithArgsAndNewViper(env, args)
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
	env := EnvVarsFromDeployment(t, d)

	e, err := ExecuteWithArgsAndNewViper(env, args)
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
	env := EnvVarsFromDeployment(t, d)

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())

	l := operator.GetGlobalAnnotationsNoError()
	g.Expect(l).NotTo(BeNil())
	g.Expect(len(l)).To(Equal(3))
	g.Expect(l["one"]).To(Equal("value-one"))
	g.Expect(l["two"]).To(Equal("value-two"))
	g.Expect(l["three"]).To(Equal("value-three"))
}

func TestOperatorWithCipherAllowList(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
	}

	toAdd := tls.InsecureCipherSuites()[0]

	args := []string{"operator", "--dry-run",
		"--cipher-allow-list", strings.ToLower(toAdd.Name),
	}
	env := EnvVarsFromDeployment(t, d)

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())

	l := operator.GetTlsCipherAllowList(e.V)
	g.Expect(l).NotTo(BeNil())
	g.Expect(len(l)).To(Equal(1))
	g.Expect(l[0]).To(Equal(toAdd.Name))

	fn, err := operator.NewCipherSuiteConfig(e.V, testLog)
	g.Expect(err).To(BeNil())
	g.Expect(fn).NotTo(BeNil())
	cfg := &tls.Config{}
	fn(cfg)

	expected := append(defaultCiphers(), toAdd.ID)
	g.Expect(cfg.CipherSuites).To(Equal(expected))
}

func TestOperatorWithInvalidCipherNameInAllowList(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
	}

	args := []string{"operator", "--dry-run",
		"--cipher-allow-list", "Foo",
	}
	env := EnvVarsFromDeployment(t, d)

	_, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).To(HaveOccurred())
}

func TestOperatorWithMultipleCipherAllowList(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
	}

	toAdd1 := tls.InsecureCipherSuites()[1]
	toAdd2 := tls.InsecureCipherSuites()[0]

	args := []string{"operator", "--dry-run",
		"--cipher-allow-list", strings.ToLower(toAdd1.Name),
		"--cipher-allow-list", toAdd2.Name,
	}
	env := EnvVarsFromDeployment(t, d)

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())

	l := operator.GetTlsCipherAllowList(e.V)
	g.Expect(l).NotTo(BeNil())
	g.Expect(len(l)).To(Equal(2))
	g.Expect(l[0]).To(Equal(toAdd1.Name))
	g.Expect(l[1]).To(Equal(toAdd2.Name))

	fn, err := operator.NewCipherSuiteConfig(e.V, testLog)
	g.Expect(err).To(BeNil())
	g.Expect(fn).NotTo(BeNil())
	cfg := &tls.Config{}
	fn(cfg)

	expected := append(defaultCiphers(), toAdd1.ID, toAdd2.ID)
	g.Expect(cfg.CipherSuites).To(Equal(expected))
}

func TestOperatorWithMultipleCommaDelimitedCipherAllowList(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
	}

	toAdd1 := tls.InsecureCipherSuites()[1]
	toAdd2 := tls.InsecureCipherSuites()[0]

	args := []string{"operator", "--dry-run",
		"--cipher-allow-list", strings.ToLower(toAdd1.Name) + "," + toAdd2.Name,
	}
	env := EnvVarsFromDeployment(t, d)

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())

	l := operator.GetTlsCipherAllowList(e.V)
	g.Expect(l).NotTo(BeNil())
	g.Expect(len(l)).To(Equal(2))
	g.Expect(l[0]).To(Equal(toAdd1.Name))
	g.Expect(l[1]).To(Equal(toAdd2.Name))

	fn, err := operator.NewCipherSuiteConfig(e.V, testLog)
	g.Expect(err).To(BeNil())
	g.Expect(fn).NotTo(BeNil())
	cfg := &tls.Config{}
	fn(cfg)

	expected := append(defaultCiphers(), toAdd1.ID, toAdd2.ID)
	g.Expect(cfg.CipherSuites).To(Equal(expected))
}

func TestOperatorWithCipherDenyList(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
	}

	toRemove := tls.CipherSuites()[0]

	args := []string{"operator", "--dry-run",
		"--cipher-deny-list", strings.ToLower(toRemove.Name),
	}
	env := EnvVarsFromDeployment(t, d)

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())

	l := operator.GetTlsCipherDenyList(e.V)
	g.Expect(l).NotTo(BeNil())
	g.Expect(len(l)).To(Equal(1))
	g.Expect(l[0]).To(Equal(toRemove.Name))

	fn, err := operator.NewCipherSuiteConfig(e.V, testLog)
	g.Expect(err).To(BeNil())
	g.Expect(fn).NotTo(BeNil())
	cfg := &tls.Config{}
	fn(cfg)

	expected := operator.RemoveFromUInt16Array(defaultCiphers(), toRemove.ID)
	g.Expect(cfg.CipherSuites).To(Equal(expected))
}

func TestOperatorWithInvalidCipherNameInDenyList(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
	}

	args := []string{"operator", "--dry-run",
		"--cipher-deny-list", "Foo",
	}
	env := EnvVarsFromDeployment(t, d)

	_, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).To(HaveOccurred())
}

func TestOperatorWithMultipleCipherDenyList(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
	}

	toRemove1 := tls.CipherSuites()[1]
	toRemove2 := tls.CipherSuites()[0]

	args := []string{"operator", "--dry-run",
		"--cipher-deny-list", strings.ToLower(toRemove1.Name),
		"--cipher-deny-list", toRemove2.Name,
	}
	env := EnvVarsFromDeployment(t, d)

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())

	l := operator.GetTlsCipherDenyList(e.V)
	g.Expect(l).NotTo(BeNil())
	g.Expect(len(l)).To(Equal(2))
	g.Expect(l[0]).To(Equal(toRemove1.Name))
	g.Expect(l[1]).To(Equal(toRemove2.Name))

	fn, err := operator.NewCipherSuiteConfig(e.V, testLog)
	g.Expect(err).To(BeNil())
	g.Expect(fn).NotTo(BeNil())
	cfg := &tls.Config{}
	fn(cfg)

	expected := operator.RemoveAllFromUInt16Array(defaultCiphers(), toRemove1.ID, toRemove2.ID)
	g.Expect(cfg.CipherSuites).To(Equal(expected))
}

func TestOperatorWithMultipleCommaDelimitedCipherDenyList(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
	}

	toRemove1 := tls.CipherSuites()[1]
	toRemove2 := tls.CipherSuites()[0]

	args := []string{"operator", "--dry-run",
		"--cipher-deny-list", strings.ToLower(toRemove1.Name) + "," + toRemove2.Name,
	}
	env := EnvVarsFromDeployment(t, d)

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())

	l := operator.GetTlsCipherDenyList(e.V)
	g.Expect(l).NotTo(BeNil())
	g.Expect(len(l)).To(Equal(2))
	g.Expect(l[0]).To(Equal(toRemove1.Name))
	g.Expect(l[1]).To(Equal(toRemove2.Name))

	fn, err := operator.NewCipherSuiteConfig(e.V, testLog)
	g.Expect(err).To(BeNil())
	g.Expect(fn).NotTo(BeNil())
	cfg := &tls.Config{}
	fn(cfg)

	expected := operator.RemoveAllFromUInt16Array(defaultCiphers(), toRemove1.ID, toRemove2.ID)
	g.Expect(cfg.CipherSuites).To(Equal(expected))
}

func TestOperatorWithAllowListAndDenyList(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
	}

	toAdd := tls.InsecureCipherSuites()[1]
	toRemove := tls.CipherSuites()[0]

	args := []string{"operator", "--dry-run",
		"--cipher-allow-list", toAdd.Name,
		"--cipher-deny-list", toRemove.Name,
	}
	env := EnvVarsFromDeployment(t, d)

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())

	bl := operator.GetTlsCipherDenyList(e.V)
	g.Expect(bl).NotTo(BeNil())
	g.Expect(len(bl)).To(Equal(1))
	g.Expect(bl[0]).To(Equal(toRemove.Name))

	wl := operator.GetTlsCipherAllowList(e.V)
	g.Expect(wl).NotTo(BeNil())
	g.Expect(len(wl)).To(Equal(1))
	g.Expect(wl[0]).To(Equal(toAdd.Name))

	fn, err := operator.NewCipherSuiteConfig(e.V, testLog)
	g.Expect(err).To(BeNil())
	g.Expect(fn).NotTo(BeNil())
	cfg := &tls.Config{}
	fn(cfg)

	expected := operator.RemoveFromUInt16Array(defaultCiphers(), toRemove.ID)
	expected = append(expected, toAdd.ID)
	g.Expect(cfg.CipherSuites).To(Equal(expected))
}

func TestOperatorWithAllowListAndDenyListSameCipher(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
	}

	cipher := tls.InsecureCipherSuites()[1]

	args := []string{"operator", "--dry-run",
		"--cipher-allow-list", cipher.Name,
		"--cipher-deny-list", cipher.Name,
	}
	env := EnvVarsFromDeployment(t, d)

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())

	bl := operator.GetTlsCipherDenyList(e.V)
	g.Expect(bl).NotTo(BeNil())
	g.Expect(len(bl)).To(Equal(1))
	g.Expect(bl[0]).To(Equal(cipher.Name))

	wl := operator.GetTlsCipherAllowList(e.V)
	g.Expect(wl).NotTo(BeNil())
	g.Expect(len(wl)).To(Equal(1))
	g.Expect(wl[0]).To(Equal(cipher.Name))

	fn, err := operator.NewCipherSuiteConfig(e.V, testLog)
	g.Expect(err).To(BeNil())
	g.Expect(fn).NotTo(BeNil())
	cfg := &tls.Config{}
	fn(cfg)

	expected := defaultCiphers()
	g.Expect(cfg.CipherSuites).To(Equal(expected))
}

func TestOperatorWithDenyAllCiphers(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
	}

	cipher := tls.InsecureCipherSuites()[1]

	args := []string{"operator", "--dry-run",
		"--cipher-allow-list", cipher.Name,
		"--cipher-deny-list", "ALL",
	}
	env := EnvVarsFromDeployment(t, d)

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())

	bl := operator.GetTlsCipherDenyList(e.V)
	g.Expect(bl).NotTo(BeNil())
	g.Expect(len(bl)).To(Equal(1))
	g.Expect(bl[0]).To(Equal("ALL"))

	wl := operator.GetTlsCipherAllowList(e.V)
	g.Expect(wl).NotTo(BeNil())
	g.Expect(len(wl)).To(Equal(1))
	g.Expect(wl[0]).To(Equal(cipher.Name))

	fn, err := operator.NewCipherSuiteConfig(e.V, testLog)
	g.Expect(err).To(BeNil())
	g.Expect(fn).NotTo(BeNil())
	cfg := &tls.Config{}
	fn(cfg)

	expected := []uint16{cipher.ID}
	g.Expect(cfg.CipherSuites).To(Equal(expected))
}

func TestOperatorWithDenyAllCiphersButNoAllowList(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
	}

	args := []string{"operator", "--dry-run", "--cipher-deny-list", "ALL"}
	env := EnvVarsFromDeployment(t, d)

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(BeNil())

	bl := operator.GetTlsCipherDenyList(e.V)
	g.Expect(bl).NotTo(BeNil())
	g.Expect(len(bl)).To(Equal(1))
	g.Expect(bl[0]).To(Equal("ALL"))

	wl := operator.GetTlsCipherAllowList(e.V)
	g.Expect(wl).To(BeNil())

	_, err = operator.NewCipherSuiteConfig(e.V, testLog)
	g.Expect(err).NotTo(BeNil())
}

func defaultCiphers() []uint16 {
	var ciphers []uint16
	for _, i := range tls.CipherSuites() {
		ciphers = append(ciphers, i.ID)
	}
	ciphers = operator.RemoveAllFromUInt16Array(ciphers, operator.DefaultCipherDenyList()...)
	return ciphers
}
