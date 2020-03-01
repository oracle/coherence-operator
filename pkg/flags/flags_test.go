/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package flags_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/oracle/coherence-operator/pkg/flags"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	"github.com/spf13/pflag"
	"reflect"
)

var _ = Describe("Coherence Operator Flags tests", func() {
	Context("When flags are valid", func() {
		var args []string
		var cohFlags flags.CoherenceOperatorFlags
		var defaultCRDs string
		var err error
		var dfltCohImg = *flags.GetDefaultCoherenceImage()
		var dfltUtilsImg = *flags.GetDefaultCoherenceUtilsImage()

		JustBeforeEach(func() {
			defaultCRDs, err = helper.FindCrdDir()
			Expect(err).NotTo(HaveOccurred())
			flags.SetDefaultCrdFiles(defaultCRDs)

			flagSet := pflag.FlagSet{}
			cohFlags = flags.CoherenceOperatorFlags{}
			cohFlags.AddTo(&flagSet)
			err = flagSet.Parse(args)
			Expect(err).NotTo(HaveOccurred())
		})

		When("no flags set", func() {
			BeforeEach(func() {
				args = []string{}
			})

			It("should have default CRD directory", func() {
				crds := cohFlags.DefaultCrdFiles()
				Expect(cohFlags.CrdFiles).To(Equal(crds))
			})

			It("should have empty always pull suffixes", func() {
				Expect(cohFlags.AlwaysPullSuffixes).To(Equal(""))
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

			It("should have the default Coherence Image", func() {
				Expect(cohFlags.CoherenceImage).To(Equal(dfltCohImg))
			})

			It("should have negative service port", func() {
				Expect(cohFlags.CoherenceUtilsImage).To(Equal(dfltUtilsImg))
			})
		})

		When("crd-files set", func() {
			var expected flags.CoherenceOperatorFlags

			BeforeEach(func() {
				crds := "/test-crds"
				args = []string{"--crd-files", crds}
				expected = flags.CoherenceOperatorFlags{
					CrdFiles:            crds,
					RestHost:            flags.DefaultRestHost,
					RestPort:            flags.DefaultRestPort,
					ServiceName:         "",
					ServicePort:         -1,
					SiteLabel:           flags.DefaultSiteLabel,
					RackLabel:           flags.DefaultRackLabel,
					AlwaysPullSuffixes:  "",
					CoherenceImage:      dfltCohImg,
					CoherenceUtilsImage: dfltUtilsImg,
				}
			})

			It("should have the correct flags", func() {
				Expect(reflect.DeepEqual(cohFlags, expected)).To(BeTrue())
			})
		})

		When("rest-host set", func() {
			var expected flags.CoherenceOperatorFlags

			BeforeEach(func() {
				restHost := "10.10.123.0"
				args = []string{"--rest-host", restHost}
				expected = flags.CoherenceOperatorFlags{
					CrdFiles:            cohFlags.DefaultCrdFiles(),
					RestHost:            restHost,
					RestPort:            flags.DefaultRestPort,
					ServiceName:         "",
					ServicePort:         -1,
					SiteLabel:           flags.DefaultSiteLabel,
					RackLabel:           flags.DefaultRackLabel,
					AlwaysPullSuffixes:  "",
					CoherenceImage:      dfltCohImg,
					CoherenceUtilsImage: dfltUtilsImg,
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
					CrdFiles:            cohFlags.DefaultCrdFiles(),
					RestHost:            flags.DefaultRestHost,
					RestPort:            9000,
					ServiceName:         "",
					ServicePort:         -1,
					SiteLabel:           flags.DefaultSiteLabel,
					RackLabel:           flags.DefaultRackLabel,
					AlwaysPullSuffixes:  "",
					CoherenceImage:      dfltCohImg,
					CoherenceUtilsImage: dfltUtilsImg,
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
					CrdFiles:            cohFlags.DefaultCrdFiles(),
					RestHost:            flags.DefaultRestHost,
					RestPort:            flags.DefaultRestPort,
					ServiceName:         "foo.com",
					ServicePort:         -1,
					SiteLabel:           flags.DefaultSiteLabel,
					RackLabel:           flags.DefaultRackLabel,
					AlwaysPullSuffixes:  "",
					CoherenceImage:      dfltCohImg,
					CoherenceUtilsImage: dfltUtilsImg,
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
					CrdFiles:            cohFlags.DefaultCrdFiles(),
					RestHost:            flags.DefaultRestHost,
					RestPort:            flags.DefaultRestPort,
					ServiceName:         "",
					ServicePort:         80,
					SiteLabel:           flags.DefaultSiteLabel,
					RackLabel:           flags.DefaultRackLabel,
					AlwaysPullSuffixes:  "",
					CoherenceImage:      dfltCohImg,
					CoherenceUtilsImage: dfltUtilsImg,
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
					CrdFiles:            cohFlags.DefaultCrdFiles(),
					RestHost:            flags.DefaultRestHost,
					RestPort:            flags.DefaultRestPort,
					ServiceName:         "",
					ServicePort:         -1,
					SiteLabel:           "foo",
					RackLabel:           flags.DefaultRackLabel,
					AlwaysPullSuffixes:  "",
					CoherenceImage:      dfltCohImg,
					CoherenceUtilsImage: dfltUtilsImg,
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
					CrdFiles:            cohFlags.DefaultCrdFiles(),
					RestHost:            flags.DefaultRestHost,
					RestPort:            flags.DefaultRestPort,
					ServiceName:         "",
					ServicePort:         -1,
					SiteLabel:           flags.DefaultSiteLabel,
					RackLabel:           "foo",
					AlwaysPullSuffixes:  "",
					CoherenceImage:      dfltCohImg,
					CoherenceUtilsImage: dfltUtilsImg,
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
					CrdFiles:            cohFlags.DefaultCrdFiles(),
					RestHost:            flags.DefaultRestHost,
					RestPort:            flags.DefaultRestPort,
					ServiceName:         "",
					ServicePort:         -1,
					SiteLabel:           flags.DefaultSiteLabel,
					RackLabel:           flags.DefaultRackLabel,
					AlwaysPullSuffixes:  "-ci,latest",
					CoherenceImage:      dfltCohImg,
					CoherenceUtilsImage: dfltUtilsImg,
				}
			})

			It("should have the correct flags", func() {
				Expect(reflect.DeepEqual(cohFlags, expected)).To(BeTrue())
			})
		})
	})
})
