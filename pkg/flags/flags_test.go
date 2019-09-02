package flags_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/oracle/coherence-operator/pkg/flags"
	"github.com/spf13/pflag"
	"reflect"
)

var _ = Describe("Coherence Operator Flags tests", func() {
	Context("When flags are valid", func() {
		var args []string
		var cohFlags flags.CoherenceOperatorFlags

		JustBeforeEach(func() {
			flagSet := pflag.FlagSet{}
			cohFlags = flags.CoherenceOperatorFlags{}
			cohFlags.AddTo(&flagSet)
			err := flagSet.Parse(args)
			Expect(err).NotTo(HaveOccurred())
		})

		When("no flags set", func() {
			BeforeEach(func() {
				args = []string{}
			})

			It("should have empty always pull suffixes", func() {
				Expect(cohFlags.AlwaysPullSuffixes).To(Equal(""))
			})

			It("should have empty elastic search credentials", func() {
				Expect(cohFlags.ElasticSearchCredentials).To(Equal(""))
			})

			It("should have empty elastic search host", func() {
				Expect(cohFlags.ElasticSearchHost).To(Equal(""))
			})

			It("should have negative elastic search port", func() {
				Expect(cohFlags.ElasticSearchPort).To(Equal(int32(-1)))
			})

			It("should have empty elastic search user", func() {
				Expect(cohFlags.ElasticSearchUser).To(Equal(""))
			})

			It("should set log integration to false ", func() {
				Expect(cohFlags.LogIntegrationEnabled).To(Equal(false))
			})

			It("should have default rack label flag", func() {
				Expect(cohFlags.RackLabel).To(Equal(flags.DefaultRackLabel))
			})

			It("should have default ReST host ", func() {
				Expect(cohFlags.RestHost).To(Equal(flags.DefaultRestHost))
			})

			It("should have default ReST port", func() {
				Expect(cohFlags.RestPort).To(Equal(flags.DefaultRestPort))
			})

			It("should have empty service name", func() {
				Expect(cohFlags.ServiceName).To(Equal(""))
			})

			It("should have negative service port", func() {
				Expect(cohFlags.ServicePort).To(Equal(int32(-1)))
			})

			It("should have empty SSL CA File", func() {
				Expect(cohFlags.SSLCAFile).To(Equal(""))
			})

			It("should have empty SSL Cert File", func() {
				Expect(cohFlags.SSLCertFile).To(Equal(""))
			})

			It("should have empty SSL Key File", func() {
				Expect(cohFlags.SSLKeyFile).To(Equal(""))
			})
		})

		When("rest-host set", func() {
			var expected flags.CoherenceOperatorFlags

			BeforeEach(func() {
				restHost := "10.10.123.0"
				args = []string{"--rest-host", restHost}
				expected = flags.CoherenceOperatorFlags{
					RestHost:                 restHost,
					RestPort:                 flags.DefaultRestPort,
					ServiceName:              "",
					ServicePort:              -1,
					LogIntegrationEnabled:    false,
					ElasticSearchHost:        "",
					ElasticSearchPort:        -1,
					ElasticSearchUser:        "",
					ElasticSearchCredentials: "",
					SSLKeyFile:               "",
					SSLCertFile:              "",
					SSLCAFile:                "",
					SiteLabel:                flags.DefaultSiteLabel,
					RackLabel:                flags.DefaultRackLabel,
					AlwaysPullSuffixes:       "",
				}
			})

			It("should have the correct flags", func() {
				Expect(reflect.DeepEqual(cohFlags, expected)).To(BeTrue())
			})
		})

		When("rest-port set", func() {
			var expected flags.CoherenceOperatorFlags

			BeforeEach(func() {
				args = []string{"--rest-port", "9000"}
				expected = flags.CoherenceOperatorFlags{
					RestHost:                 flags.DefaultRestHost,
					RestPort:                 9000,
					ServiceName:              "",
					ServicePort:              -1,
					LogIntegrationEnabled:    false,
					ElasticSearchHost:        "",
					ElasticSearchPort:        -1,
					ElasticSearchUser:        "",
					ElasticSearchCredentials: "",
					SSLKeyFile:               "",
					SSLCertFile:              "",
					SSLCAFile:                "",
					SiteLabel:                flags.DefaultSiteLabel,
					RackLabel:                flags.DefaultRackLabel,
					AlwaysPullSuffixes:       "",
				}
			})

			It("should have the correct flags", func() {
				Expect(reflect.DeepEqual(cohFlags, expected)).To(BeTrue())
			})
		})

		When("service-name set", func() {
			var expected flags.CoherenceOperatorFlags

			BeforeEach(func() {
				args = []string{"--service-name", "foo.com"}
				expected = flags.CoherenceOperatorFlags{
					RestHost:                 flags.DefaultRestHost,
					RestPort:                 flags.DefaultRestPort,
					ServiceName:              "foo.com",
					ServicePort:              -1,
					LogIntegrationEnabled:    false,
					ElasticSearchHost:        "",
					ElasticSearchPort:        -1,
					ElasticSearchUser:        "",
					ElasticSearchCredentials: "",
					SSLKeyFile:               "",
					SSLCertFile:              "",
					SSLCAFile:                "",
					SiteLabel:                flags.DefaultSiteLabel,
					RackLabel:                flags.DefaultRackLabel,
					AlwaysPullSuffixes:       "",
				}
			})

			It("should have the correct flags", func() {
				Expect(reflect.DeepEqual(cohFlags, expected)).To(BeTrue())
			})
		})

		When("service-port set", func() {
			var expected flags.CoherenceOperatorFlags

			BeforeEach(func() {
				args = []string{"--service-port", "80"}
				expected = flags.CoherenceOperatorFlags{
					RestHost:                 flags.DefaultRestHost,
					RestPort:                 flags.DefaultRestPort,
					ServiceName:              "",
					ServicePort:              80,
					LogIntegrationEnabled:    false,
					ElasticSearchHost:        "",
					ElasticSearchPort:        -1,
					ElasticSearchUser:        "",
					ElasticSearchCredentials: "",
					SSLKeyFile:               "",
					SSLCertFile:              "",
					SSLCAFile:                "",
					SiteLabel:                flags.DefaultSiteLabel,
					RackLabel:                flags.DefaultRackLabel,
					AlwaysPullSuffixes:       "",
				}
			})

			It("should have the correct flags", func() {
				Expect(reflect.DeepEqual(cohFlags, expected)).To(BeTrue())
			})
		})

		When("log-integration set", func() {
			var expected flags.CoherenceOperatorFlags

			BeforeEach(func() {
				args = []string{"--log-integration"}
				expected = flags.CoherenceOperatorFlags{
					RestHost:                 flags.DefaultRestHost,
					RestPort:                 flags.DefaultRestPort,
					ServiceName:              "",
					ServicePort:              -1,
					LogIntegrationEnabled:    true,
					ElasticSearchHost:        "",
					ElasticSearchPort:        -1,
					ElasticSearchUser:        "",
					ElasticSearchCredentials: "",
					SSLKeyFile:               "",
					SSLCertFile:              "",
					SSLCAFile:                "",
					SiteLabel:                flags.DefaultSiteLabel,
					RackLabel:                flags.DefaultRackLabel,
					AlwaysPullSuffixes:       "",
				}
			})

			It("should have the correct flags", func() {
				Expect(reflect.DeepEqual(cohFlags, expected)).To(BeTrue())
			})
		})

		When("log-integration set to true", func() {
			var expected flags.CoherenceOperatorFlags

			BeforeEach(func() {
				args = []string{"--log-integration", "true"}
				expected = flags.CoherenceOperatorFlags{
					RestHost:                 flags.DefaultRestHost,
					RestPort:                 flags.DefaultRestPort,
					ServiceName:              "",
					ServicePort:              -1,
					LogIntegrationEnabled:    true,
					ElasticSearchHost:        "",
					ElasticSearchPort:        -1,
					ElasticSearchUser:        "",
					ElasticSearchCredentials: "",
					SSLKeyFile:               "",
					SSLCertFile:              "",
					SSLCAFile:                "",
					SiteLabel:                flags.DefaultSiteLabel,
					RackLabel:                flags.DefaultRackLabel,
					AlwaysPullSuffixes:       "",
				}
			})

			It("should have the correct flags", func() {
				Expect(reflect.DeepEqual(cohFlags, expected)).To(BeTrue())
			})
		})

		When("log-integration set to false", func() {
			var expected flags.CoherenceOperatorFlags

			BeforeEach(func() {
				args = []string{"--log-integration=false"}
				expected = flags.CoherenceOperatorFlags{
					RestHost:                 flags.DefaultRestHost,
					RestPort:                 flags.DefaultRestPort,
					ServiceName:              "",
					ServicePort:              -1,
					LogIntegrationEnabled:    false,
					ElasticSearchHost:        "",
					ElasticSearchPort:        -1,
					ElasticSearchUser:        "",
					ElasticSearchCredentials: "",
					SSLKeyFile:               "",
					SSLCertFile:              "",
					SSLCAFile:                "",
					SiteLabel:                flags.DefaultSiteLabel,
					RackLabel:                flags.DefaultRackLabel,
					AlwaysPullSuffixes:       "",
				}
			})

			It("should have the correct flags", func() {
				Expect(reflect.DeepEqual(cohFlags, expected)).To(BeTrue())
			})
		})

		When("es-host set", func() {
			var expected flags.CoherenceOperatorFlags

			BeforeEach(func() {
				args = []string{"--es-host", "es.com"}
				expected = flags.CoherenceOperatorFlags{
					RestHost:                 flags.DefaultRestHost,
					RestPort:                 flags.DefaultRestPort,
					ServiceName:              "",
					ServicePort:              -1,
					LogIntegrationEnabled:    false,
					ElasticSearchHost:        "es.com",
					ElasticSearchPort:        -1,
					ElasticSearchUser:        "",
					ElasticSearchCredentials: "",
					SSLKeyFile:               "",
					SSLCertFile:              "",
					SSLCAFile:                "",
					SiteLabel:                flags.DefaultSiteLabel,
					RackLabel:                flags.DefaultRackLabel,
					AlwaysPullSuffixes:       "",
				}
			})

			It("should have the correct flags", func() {
				Expect(reflect.DeepEqual(cohFlags, expected)).To(BeTrue())
			})
		})

		When("es-port set", func() {
			var expected flags.CoherenceOperatorFlags

			BeforeEach(func() {
				args = []string{"--es-port", "2000"}
				expected = flags.CoherenceOperatorFlags{
					RestHost:                 flags.DefaultRestHost,
					RestPort:                 flags.DefaultRestPort,
					ServiceName:              "",
					ServicePort:              -1,
					LogIntegrationEnabled:    false,
					ElasticSearchHost:        "",
					ElasticSearchPort:        2000,
					ElasticSearchUser:        "",
					ElasticSearchCredentials: "",
					SSLKeyFile:               "",
					SSLCertFile:              "",
					SSLCAFile:                "",
					SiteLabel:                flags.DefaultSiteLabel,
					RackLabel:                flags.DefaultRackLabel,
					AlwaysPullSuffixes:       "",
				}
			})

			It("should have the correct flags", func() {
				Expect(reflect.DeepEqual(cohFlags, expected)).To(BeTrue())
			})
		})

		When("es-user set", func() {
			var expected flags.CoherenceOperatorFlags

			BeforeEach(func() {
				args = []string{"--es-user", "admin"}
				expected = flags.CoherenceOperatorFlags{
					RestHost:                 flags.DefaultRestHost,
					RestPort:                 flags.DefaultRestPort,
					ServiceName:              "",
					ServicePort:              -1,
					LogIntegrationEnabled:    false,
					ElasticSearchHost:        "",
					ElasticSearchPort:        -1,
					ElasticSearchUser:        "admin",
					ElasticSearchCredentials: "",
					SSLKeyFile:               "",
					SSLCertFile:              "",
					SSLCAFile:                "",
					SiteLabel:                flags.DefaultSiteLabel,
					RackLabel:                flags.DefaultRackLabel,
					AlwaysPullSuffixes:       "",
				}
			})

			It("should have the correct flags", func() {
				Expect(reflect.DeepEqual(cohFlags, expected)).To(BeTrue())
			})
		})

		When("es-password set", func() {
			var expected flags.CoherenceOperatorFlags

			BeforeEach(func() {
				args = []string{"--es-password", "secret"}
				expected = flags.CoherenceOperatorFlags{
					RestHost:                 flags.DefaultRestHost,
					RestPort:                 flags.DefaultRestPort,
					ServiceName:              "",
					ServicePort:              -1,
					LogIntegrationEnabled:    false,
					ElasticSearchHost:        "",
					ElasticSearchPort:        -1,
					ElasticSearchUser:        "",
					ElasticSearchCredentials: "secret",
					SSLKeyFile:               "",
					SSLCertFile:              "",
					SSLCAFile:                "",
					SiteLabel:                flags.DefaultSiteLabel,
					RackLabel:                flags.DefaultRackLabel,
					AlwaysPullSuffixes:       "",
				}
			})

			It("should have the correct flags", func() {
				Expect(reflect.DeepEqual(cohFlags, expected)).To(BeTrue())
			})
		})

		When("site-label set", func() {
			var expected flags.CoherenceOperatorFlags

			BeforeEach(func() {
				args = []string{"--site-label", "foo"}
				expected = flags.CoherenceOperatorFlags{
					RestHost:                 flags.DefaultRestHost,
					RestPort:                 flags.DefaultRestPort,
					ServiceName:              "",
					ServicePort:              -1,
					LogIntegrationEnabled:    false,
					ElasticSearchHost:        "",
					ElasticSearchPort:        -1,
					ElasticSearchUser:        "",
					ElasticSearchCredentials: "",
					SSLKeyFile:               "",
					SSLCertFile:              "",
					SSLCAFile:                "",
					SiteLabel:                "foo",
					RackLabel:                flags.DefaultRackLabel,
					AlwaysPullSuffixes:       "",
				}
			})

			It("should have the correct flags", func() {
				Expect(reflect.DeepEqual(cohFlags, expected)).To(BeTrue())
			})
		})

		When("rack-label set", func() {
			var expected flags.CoherenceOperatorFlags

			BeforeEach(func() {
				args = []string{"--rack-label", "foo"}
				expected = flags.CoherenceOperatorFlags{
					RestHost:                 flags.DefaultRestHost,
					RestPort:                 flags.DefaultRestPort,
					ServiceName:              "",
					ServicePort:              -1,
					LogIntegrationEnabled:    false,
					ElasticSearchHost:        "",
					ElasticSearchPort:        -1,
					ElasticSearchUser:        "",
					ElasticSearchCredentials: "",
					SSLKeyFile:               "",
					SSLCertFile:              "",
					SSLCAFile:                "",
					SiteLabel:                flags.DefaultSiteLabel,
					RackLabel:                "foo",
					AlwaysPullSuffixes:       "",
				}
			})

			It("should have the correct flags", func() {
				Expect(reflect.DeepEqual(cohFlags, expected)).To(BeTrue())
			})
		})

		When("force-always-pull-tags set", func() {
			var expected flags.CoherenceOperatorFlags

			BeforeEach(func() {
				args = []string{"--force-always-pull-tags", "-ci,latest"}
				expected = flags.CoherenceOperatorFlags{
					RestHost:                 flags.DefaultRestHost,
					RestPort:                 flags.DefaultRestPort,
					ServiceName:              "",
					ServicePort:              -1,
					LogIntegrationEnabled:    false,
					ElasticSearchHost:        "",
					ElasticSearchPort:        -1,
					ElasticSearchUser:        "",
					ElasticSearchCredentials: "",
					SSLKeyFile:               "",
					SSLCertFile:              "",
					SSLCAFile:                "",
					SiteLabel:                flags.DefaultSiteLabel,
					RackLabel:                flags.DefaultRackLabel,
					AlwaysPullSuffixes:       "-ci,latest",
				}
			})

			It("should have the correct flags", func() {
				Expect(reflect.DeepEqual(cohFlags, expected)).To(BeTrue())
			})
		})

		When("ssl-key-file set", func() {
			var expected flags.CoherenceOperatorFlags

			BeforeEach(func() {
				args = []string{"--ssl-key-file", "my.key"}
				expected = flags.CoherenceOperatorFlags{
					RestHost:                 flags.DefaultRestHost,
					RestPort:                 flags.DefaultRestPort,
					ServiceName:              "",
					ServicePort:              -1,
					LogIntegrationEnabled:    false,
					ElasticSearchHost:        "",
					ElasticSearchPort:        -1,
					ElasticSearchUser:        "",
					ElasticSearchCredentials: "",
					SSLKeyFile:               "my.key",
					SSLCertFile:              "",
					SSLCAFile:                "",
					SiteLabel:                flags.DefaultSiteLabel,
					RackLabel:                flags.DefaultRackLabel,
					AlwaysPullSuffixes:       "",
				}
			})

			It("should have the correct flags", func() {
				Expect(reflect.DeepEqual(cohFlags, expected)).To(BeTrue())
			})
		})

		When("ssl-cert-file set", func() {
			var expected flags.CoherenceOperatorFlags

			BeforeEach(func() {
				args = []string{"--ssl-cert-file", "my.cert"}
				expected = flags.CoherenceOperatorFlags{
					RestHost:                 flags.DefaultRestHost,
					RestPort:                 flags.DefaultRestPort,
					ServiceName:              "",
					ServicePort:              -1,
					LogIntegrationEnabled:    false,
					ElasticSearchHost:        "",
					ElasticSearchPort:        -1,
					ElasticSearchUser:        "",
					ElasticSearchCredentials: "",
					SSLKeyFile:               "",
					SSLCertFile:              "my.cert",
					SSLCAFile:                "",
					SiteLabel:                flags.DefaultSiteLabel,
					RackLabel:                flags.DefaultRackLabel,
					AlwaysPullSuffixes:       "",
				}
			})

			It("should have the correct flags", func() {
				Expect(reflect.DeepEqual(cohFlags, expected)).To(BeTrue())
			})
		})

		When("ssl-ca-file set", func() {
			var expected flags.CoherenceOperatorFlags

			BeforeEach(func() {
				args = []string{"--ssl-ca-file", "my-ca.cert"}
				expected = flags.CoherenceOperatorFlags{
					RestHost:                 flags.DefaultRestHost,
					RestPort:                 flags.DefaultRestPort,
					ServiceName:              "",
					ServicePort:              -1,
					LogIntegrationEnabled:    false,
					ElasticSearchHost:        "",
					ElasticSearchPort:        -1,
					ElasticSearchUser:        "",
					ElasticSearchCredentials: "",
					SSLKeyFile:               "",
					SSLCertFile:              "",
					SSLCAFile:                "my-ca.cert",
					SiteLabel:                flags.DefaultSiteLabel,
					RackLabel:                flags.DefaultRackLabel,
					AlwaysPullSuffixes:       "",
				}
			})

			It("should have the correct flags", func() {
				Expect(reflect.DeepEqual(cohFlags, expected)).To(BeTrue())
			})
		})
	})
})
