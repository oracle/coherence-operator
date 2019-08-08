package v1_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	coherence "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
)

var _ = Describe("Testing SSLSpec struct", func() {

	Context("Copying a SSLSpec using DeepCopyWithDefaults", func() {
		var original *coherence.SSLSpec
		var defaults *coherence.SSLSpec
		var clone *coherence.SSLSpec
		var expected *coherence.SSLSpec

		NewSSLSpecOne := func() *coherence.SSLSpec {
			return &coherence.SSLSpec {
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
			}
		}

		NewSSLSpecTwo := func() *coherence.SSLSpec {
			return &coherence.SSLSpec {
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
			}
		}

		ValidateResult := func() {
			It("should have correct Enabled", func() {
				Expect(*clone.Enabled).To(Equal(*expected.Enabled))
			})

			It("should have correct Secrets", func() {
				Expect(*clone.Secrets).To(Equal(*expected.Secrets))
			})

			It("should have correct KeyStore", func() {
				Expect(*clone.KeyStore).To(Equal(*expected.KeyStore))
			})

			It("should have correct KeyStorePasswordFile", func() {
				Expect(*clone.KeyStorePasswordFile).To(Equal(*expected.KeyStorePasswordFile))
			})

			It("should have correct KeyPasswordFile", func() {
				Expect(*clone.KeyPasswordFile).To(Equal(*expected.KeyPasswordFile))
			})

			It("should have correct KeyStoreAlgorithm", func() {
				Expect(*clone.KeyStoreAlgorithm).To(Equal(*expected.KeyStoreAlgorithm))
			})

			It("should have correct KeyStoreProvider", func() {
				Expect(*clone.KeyStoreProvider).To(Equal(*expected.KeyStoreProvider))
			})

			It("should have correct KeyStoreType", func() {
				Expect(*clone.KeyStoreType).To(Equal(*expected.KeyStoreType))
			})

			It("should have correct TrustStore", func() {
				Expect(*clone.TrustStore).To(Equal(*expected.TrustStore))
			})

			It("should have correct TrustStorePasswordFile", func() {
				Expect(*clone.TrustStorePasswordFile).To(Equal(*expected.TrustStorePasswordFile))
			})

			It("should have correct TrustStoreAlgorithm", func() {
				Expect(*clone.TrustStoreAlgorithm).To(Equal(*expected.TrustStoreAlgorithm))
			})

			It("should have correct TrustStoreProvider", func() {
				Expect(*clone.TrustStoreProvider).To(Equal(*expected.TrustStoreProvider))
			})

			It("should have correct TrustStoreType", func() {
				Expect(*clone.TrustStoreType).To(Equal(*expected.TrustStoreType))
			})

			It("should have correct RequireClientCert", func() {
				Expect(*clone.RequireClientCert).To(Equal(*expected.RequireClientCert))
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
				original = NewSSLSpecOne()
				defaults = nil
				expected = original
			})

			ValidateResult()
		})

		When("original is nil", func() {
			BeforeEach(func() {
				defaults = NewSSLSpecOne()
				original = nil
				expected = defaults
			})

			ValidateResult()
		})

		When("all original fields are set", func() {
			BeforeEach(func() {
				original = NewSSLSpecOne()
				defaults = NewSSLSpecTwo()
				expected = original
			})

			ValidateResult()
		})

		When("original Enabled is nil", func() {
			BeforeEach(func() {
				original = NewSSLSpecOne()
				original.Enabled = nil
				defaults = NewSSLSpecTwo()

				expected = NewSSLSpecOne()
				expected.Enabled = defaults.Enabled
			})

			ValidateResult()
		})

		When("original Secrets is nil", func() {
			BeforeEach(func() {
				original = NewSSLSpecOne()
				original.Secrets = nil
				defaults = NewSSLSpecTwo()

				expected = NewSSLSpecOne()
				expected.Secrets = defaults.Secrets
			})

			ValidateResult()
		})

		When("original KeyStore is nil", func() {
			BeforeEach(func() {
				original = NewSSLSpecOne()
				original.KeyStore = nil
				defaults = NewSSLSpecTwo()

				expected = NewSSLSpecOne()
				expected.KeyStore = defaults.KeyStore
			})

			ValidateResult()
		})

		When("original KeyStorePasswordFile is nil", func() {
			BeforeEach(func() {
				original = NewSSLSpecOne()
				original.KeyStorePasswordFile = nil
				defaults = NewSSLSpecTwo()

				expected = NewSSLSpecOne()
				expected.KeyStorePasswordFile = defaults.KeyStorePasswordFile
			})

			ValidateResult()
		})

		When("original KeyPasswordFile is nil", func() {
			BeforeEach(func() {
				original = NewSSLSpecOne()
				original.KeyPasswordFile = nil
				defaults = NewSSLSpecTwo()

				expected = NewSSLSpecOne()
				expected.KeyPasswordFile = defaults.KeyPasswordFile
			})

			ValidateResult()
		})

		When("original KeyStoreAlgorithm is nil", func() {
			BeforeEach(func() {
				original = NewSSLSpecOne()
				original.KeyStoreAlgorithm = nil
				defaults = NewSSLSpecTwo()

				expected = NewSSLSpecOne()
				expected.KeyStoreAlgorithm = defaults.KeyStoreAlgorithm
			})

			ValidateResult()
		})

		When("original KeyStoreProvider is nil", func() {
			BeforeEach(func() {
				original = NewSSLSpecOne()
				original.KeyStoreProvider = nil
				defaults = NewSSLSpecTwo()

				expected = NewSSLSpecOne()
				expected.KeyStoreProvider = defaults.KeyStoreProvider
			})

			ValidateResult()
		})

		When("original KeyStoreType  is nil", func() {
			BeforeEach(func() {
				original = NewSSLSpecOne()
				original.KeyStoreType = nil
				defaults = NewSSLSpecTwo()

				expected = NewSSLSpecOne()
				expected.KeyStoreType = defaults.KeyStoreType
			})

			ValidateResult()
		})

		When("original TrustStore is nil", func() {
			BeforeEach(func() {
				original = NewSSLSpecOne()
				original.TrustStore = nil
				defaults = NewSSLSpecTwo()

				expected = NewSSLSpecOne()
				expected.TrustStore = defaults.TrustStore
			})

			ValidateResult()
		})

		When("original TrustStorePasswordFile is nil", func() {
			BeforeEach(func() {
				original = NewSSLSpecOne()
				original.TrustStorePasswordFile = nil
				defaults = NewSSLSpecTwo()

				expected = NewSSLSpecOne()
				expected.TrustStorePasswordFile = defaults.TrustStorePasswordFile
			})

			ValidateResult()
		})

		When("original TrustStoreAlgorithm is nil", func() {
			BeforeEach(func() {
				original = NewSSLSpecOne()
				original.TrustStoreAlgorithm = nil
				defaults = NewSSLSpecTwo()

				expected = NewSSLSpecOne()
				expected.TrustStoreAlgorithm = defaults.TrustStoreAlgorithm
			})

			ValidateResult()
		})

		When("original TrustStoreProvider is nil", func() {
			BeforeEach(func() {
				original = NewSSLSpecOne()
				original.TrustStoreProvider = nil
				defaults = NewSSLSpecTwo()

				expected = NewSSLSpecOne()
				expected.TrustStoreProvider = defaults.TrustStoreProvider
			})

			ValidateResult()
		})

		When("original TrustStoreType is nil", func() {
			BeforeEach(func() {
				original = NewSSLSpecOne()
				original.TrustStoreType = nil
				defaults = NewSSLSpecTwo()

				expected = NewSSLSpecOne()
				expected.TrustStoreType = defaults.TrustStoreType
			})

			ValidateResult()
		})

		When("original RequireClientCert is nil", func() {
			BeforeEach(func() {
				original = NewSSLSpecOne()
				original.RequireClientCert = nil
				defaults = NewSSLSpecTwo()

				expected = NewSSLSpecOne()
				expected.RequireClientCert = defaults.RequireClientCert
			})

			ValidateResult()
		})
	})
})
