/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	coherence "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
)

var _ = Describe("Testing PortSpecWithSSL struct", func() {

	Context("Copying a PortSpecWithSSL using DeepCopyWithDefaults", func() {
		var original *coherence.PortSpecWithSSL
		var defaults *coherence.PortSpecWithSSL
		var clone *coherence.PortSpecWithSSL
		var expected *coherence.PortSpecWithSSL

		NewPortSpecOne := func() *coherence.PortSpecWithSSL {
			return &coherence.PortSpecWithSSL{
				PortSpec: coherence.PortSpec{
					Port: 8080,
					Service: &coherence.ServiceSpec{
						Name: stringPtr("foo"),
						Port: int32Ptr(80),
					},
				},
				SSL: &coherence.SSLSpec{
					Enabled:                boolPtr(true),
					Secrets:                stringPtr("ssl-secret"),
					KeyStore:               stringPtr("keystore.jks"),
					KeyStorePasswordFile:   stringPtr("storepassword.txt"),
					KeyPasswordFile:        stringPtr("keypassword.txt"),
					KeyStoreAlgorithm:      stringPtr("SunX509"),
					KeyStoreProvider:       stringPtr("fooJCA"),
					KeyStoreType:           stringPtr("JKS"),
					TrustStore:             stringPtr("truststore-guardians.jks"),
					TrustStorePasswordFile: stringPtr("trustpassword.txt"),
					TrustStoreAlgorithm:    stringPtr("SunX509"),
					TrustStoreProvider:     stringPtr("fooJCA"),
					TrustStoreType:         stringPtr("JKS"),
					RequireClientCert:      boolPtr(true),
				},
			}
		}

		NewPortSpecTwo := func() *coherence.PortSpecWithSSL {
			return &coherence.PortSpecWithSSL{
				PortSpec: coherence.PortSpec{
					Port: 9090,
					Service: &coherence.ServiceSpec{
						Name: stringPtr("bar"),
						Port: int32Ptr(90),
					},
				},
				SSL: &coherence.SSLSpec{
					Enabled:                boolPtr(true),
					Secrets:                stringPtr("ssl-secret2"),
					KeyStore:               stringPtr("keystore.jks"),
					KeyStorePasswordFile:   stringPtr("storepassword2.txt"),
					KeyPasswordFile:        stringPtr("keypassword2.txt"),
					KeyStoreAlgorithm:      stringPtr("SunX509"),
					KeyStoreProvider:       stringPtr("barJCA"),
					KeyStoreType:           stringPtr("JKS"),
					TrustStore:             stringPtr("truststore-guardians2.jks"),
					TrustStorePasswordFile: stringPtr("trustpassword2.txt"),
					TrustStoreAlgorithm:    stringPtr("SunX509"),
					TrustStoreProvider:     stringPtr("barJCA"),
					TrustStoreType:         stringPtr("JKS"),
					RequireClientCert:      boolPtr(false),
				},
			}
		}

		ValidateResult := func() {
			It("should have correct Port", func() {
				Expect(clone.Port).To(Equal(expected.Port))
			})

			It("should have correct SSL", func() {
				Expect(*clone.SSL).To(Equal(*expected.SSL))
			})
		}

		JustBeforeEach(func() {
			clone = original.DeepCopyWithDefaults(defaults)
		})

		When("original and defaults are nil", func() {
			BeforeEach(func() {
				original = nil
				defaults = nil
			})

			It("the copy should be nil", func() {
				Expect(clone).Should(BeNil())
			})
		})

		When("defaults is nil", func() {
			BeforeEach(func() {
				original = NewPortSpecOne()
				defaults = nil
				expected = original
			})

			ValidateResult()
		})

		When("original is nil", func() {
			BeforeEach(func() {
				defaults = NewPortSpecOne()
				original = nil
				expected = defaults
			})

			ValidateResult()
		})

		When("all original fields are set", func() {
			BeforeEach(func() {
				original = NewPortSpecOne()
				defaults = NewPortSpecTwo()
				expected = original
			})

			ValidateResult()
		})

		When("original Port is nil", func() {
			BeforeEach(func() {
				original = NewPortSpecOne()
				original.Port = 0
				defaults = NewPortSpecTwo()

				expected = NewPortSpecOne()
				expected.Port = defaults.Port
			})

			ValidateResult()
		})

		When("original SSL is nil", func() {
			BeforeEach(func() {
				original = NewPortSpecOne()
				original.SSL = nil
				defaults = NewPortSpecTwo()

				expected = NewPortSpecOne()
				expected.SSL = defaults.SSL
			})

			ValidateResult()
		})
	})
})
