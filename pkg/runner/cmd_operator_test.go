/*
 * Copyright (c) 2020, 2026, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package runner

import (
	"crypto/tls"
	"net/http"
	"strings"
	"testing"

	. "github.com/onsi/gomega"
	coh "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/pkg/operator"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilnet "k8s.io/apimachinery/pkg/util/net"
	restclient "k8s.io/client-go/rest"
	clientgotransport "k8s.io/client-go/transport"
	ctrl "sigs.k8s.io/controller-runtime"
)

var (
	testLog = ctrl.Log.WithName("test")
)

type testWrappedRoundTripper struct {
	wrapped http.RoundTripper
}

func (t testWrappedRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return t.wrapped.RoundTrip(req)
}

func (t testWrappedRoundTripper) WrappedRoundTripper() http.RoundTripper {
	return t.wrapped
}

type testUnsupportedRoundTripper struct{}

func (t testUnsupportedRoundTripper) RoundTrip(_ *http.Request) (*http.Response, error) {
	return nil, nil
}

type testSelfWrappingRoundTripper struct{}

func (t testSelfWrappingRoundTripper) RoundTrip(_ *http.Request) (*http.Response, error) {
	return nil, nil
}

func (t testSelfWrappingRoundTripper) WrappedRoundTripper() http.RoundTripper {
	return t
}

type testRecordingRoundTripper struct {
	wrapped         http.RoundTripper
	roundTripCalled *bool
}

func (t *testRecordingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// The chaining test needs to prove delegation without opening sockets, so return
	// a synthetic response while still exposing the wrapped transport for TLS mutation.
	*t.roundTripCalled = true
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       http.NoBody,
		Request:    req,
	}, nil
}

func (t *testRecordingRoundTripper) WrappedRoundTripper() http.RoundTripper {
	return t.wrapped
}

func wrapRoundTripper(rt http.RoundTripper, depth int) http.RoundTripper {
	for range depth {
		rt = testWrappedRoundTripper{wrapped: rt}
	}
	return rt
}

func TestConfigureTransportTLSWithHTTPTransport(t *testing.T) {
	g := NewGomegaWithT(t)
	transport := &http.Transport{}

	configureTransportTLS(transport, func(c *tls.Config) {
		// The helper must create missing TLS config so operator TLS options still
		// affect plain client-go transports instead of being silently skipped.
		c.MinVersion = tls.VersionTLS13
	})

	g.Expect(transport.TLSClientConfig).NotTo(BeNil())
	g.Expect(transport.TLSClientConfig.MinVersion).To(Equal(uint16(tls.VersionTLS13)))
}

func TestConfigureTransportTLSWithWrappedHTTPTransport(t *testing.T) {
	g := NewGomegaWithT(t)
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			ServerName: "api-server",
		},
	}
	wrapped := testWrappedRoundTripper{wrapped: transport}

	configureTransportTLS(wrapped, func(c *tls.Config) {
		// Wrapped transports are used by newer client-go cache layers, so unwrapping
		// preserves the expected cipher and protocol settings without changing the chain.
		c.NextProtos = []string{"http/1.1"}
	})

	g.Expect(transport.TLSClientConfig.ServerName).To(Equal("api-server"))
	g.Expect(transport.TLSClientConfig.NextProtos).To(Equal([]string{"http/1.1"}))
}

func TestConfigureTransportTLSWithMaxDepthWrappedHTTPTransport(t *testing.T) {
	g := NewGomegaWithT(t)
	transport := &http.Transport{}
	wrapped := wrapRoundTripper(transport, maxRoundTripperUnwrapDepth)

	configureTransportTLS(wrapped, func(c *tls.Config) {
		// The depth constant is intended to allow this many real wrapper layers;
		// only deeper or cyclic chains should fall back to best-effort skip behavior.
		c.ServerName = "api-server"
	})

	g.Expect(transport.TLSClientConfig).NotTo(BeNil())
	g.Expect(transport.TLSClientConfig.ServerName).To(Equal("api-server"))
}

func TestConfigureTransportTLSSkipsOverMaxDepthWrappedHTTPTransport(t *testing.T) {
	g := NewGomegaWithT(t)
	transport := &http.Transport{}
	wrapped := wrapRoundTripper(transport, maxRoundTripperUnwrapDepth+1)
	called := false

	configureTransportTLS(wrapped, func(*tls.Config) {
		called = true
	})

	// Over-limit chains are treated like hidden transports: leave the client alone
	// rather than risk hanging startup while chasing a pathological wrapper graph.
	g.Expect(called).To(BeFalse())
	g.Expect(transport.TLSClientConfig).To(BeNil())
}

func TestConfigureTransportTLSWithUnsupportedRoundTripper(t *testing.T) {
	g := NewGomegaWithT(t)
	called := false

	g.Expect(func() {
		configureTransportTLS(testUnsupportedRoundTripper{}, func(*tls.Config) {
			called = true
		})
	}).NotTo(Panic())
	// Unknown transport implementations should be left alone so custom clients keep working
	// even when their TLS internals are not visible to this package.
	g.Expect(called).To(BeFalse())
}

func TestConfigureTransportTLSWithSelfWrappingRoundTripper(t *testing.T) {
	g := NewGomegaWithT(t)
	called := false

	g.Expect(func() {
		configureTransportTLS(testSelfWrappingRoundTripper{}, func(*tls.Config) {
			called = true
		})
	}).NotTo(Panic())
	// The unwrap depth guard keeps malformed wrapper cycles from hanging startup;
	// TLS configuration remains best-effort when no concrete transport is reachable.
	g.Expect(called).To(BeFalse())
}

func TestConfigureTransportTLSWithClientGoCachedTransport(t *testing.T) {
	g := NewGomegaWithT(t)
	cfg := &clientgotransport.Config{
		TLS: clientgotransport.TLSConfig{
			Insecure: true,
		},
	}
	cfg.Wrap(func(rt http.RoundTripper) http.RoundTripper {
		// client-go v0.36 creates a cache-managed wrapper before invoking WrapTransport;
		// this verifies the production helper works against that real wrapper path.
		configureTransportTLS(rt, func(c *tls.Config) {
			c.ServerName = "api-server"
		})
		return rt
	})

	rt, err := clientgotransport.New(cfg)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(rt).NotTo(BeNil())

	tlsConfig, err := utilnet.TLSClientConfig(rt)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(tlsConfig).NotTo(BeNil())
	g.Expect(tlsConfig.ServerName).To(Equal("api-server"))
}

func TestConfigureTransportTLSWithConfigWrapPreservesChain(t *testing.T) {
	g := NewGomegaWithT(t)
	cfg := &restclient.Config{}
	base := &http.Transport{}
	compositionOrder := make([]string, 0, 2)
	roundTripCalled := false
	var preExistingLayer *testRecordingRoundTripper
	var operatorReceived http.RoundTripper
	var operatorReturned http.RoundTripper

	cfg.Wrap(func(rt http.RoundTripper) http.RoundTripper {
		compositionOrder = append(compositionOrder, "pre-existing")
		preExistingLayer = &testRecordingRoundTripper{
			wrapped:         rt,
			roundTripCalled: &roundTripCalled,
		}
		return preExistingLayer
	})
	cfg.Wrap(func(rt http.RoundTripper) http.RoundTripper {
		compositionOrder = append(compositionOrder, "operator")
		operatorReceived = rt
		// Match production's callback contract: mutate the reachable TLS transport
		// synchronously at WrapTransport composition time, then return the same rt.
		configureTransportTLS(rt, func(c *tls.Config) {
			c.MinVersion = tls.VersionTLS13
		})
		operatorReturned = rt
		return rt
	})

	finalRoundTripper := cfg.WrapTransport(base)

	g.Expect(compositionOrder).To(Equal([]string{"pre-existing", "operator"}))
	g.Expect(operatorReceived).To(BeIdenticalTo(preExistingLayer))
	g.Expect(operatorReturned).To(BeIdenticalTo(operatorReceived))
	g.Expect(finalRoundTripper).To(BeIdenticalTo(preExistingLayer))
	g.Expect(base.TLSClientConfig).NotTo(BeNil())
	g.Expect(base.TLSClientConfig.MinVersion).To(Equal(uint16(tls.VersionTLS13)))

	req, err := http.NewRequest(http.MethodGet, "https://example.com", nil)
	g.Expect(err).NotTo(HaveOccurred())
	resp, err := finalRoundTripper.RoundTrip(req)
	g.Expect(err).NotTo(HaveOccurred())
	defer resp.Body.Close()
	g.Expect(resp.StatusCode).To(Equal(http.StatusOK))
	g.Expect(roundTripCalled).To(BeTrue())
}

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
