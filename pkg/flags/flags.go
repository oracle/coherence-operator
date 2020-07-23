/*
 * Copyright (c) 2019, 2020, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package flags

import (
	"github.com/spf13/pflag"
	"os"
)

const (
	zoneLabel = "failure-domain.beta.kubernetes.io/zone"
	//	regionLabel            = "failure-domain.beta.kubernetes.io/region"
	DefaultSiteLabel       = zoneLabel
	DefaultRackLabel       = zoneLabel
	DefaultRestHost        = "0.0.0.0"
	DefaultRestPort  int32 = 8000

	// The environment variable holding the default Coherence image name
	coherenceImageEnv = "COHERENCE_IMAGE"
	// The environment variable holding the default Coherence Utils image name
	utilsImageEnv = "UTILS_IMAGE"

	FlagRestHost       = "rest-host"
	FlagRestPort       = "rest-port"
	FlagServiceName    = "service-name"
	FlagServicePort    = "service-port"
	FlagSiteLabel      = "site-label"
	FlagRackLabel      = "rack-label"
	FlagCoherenceImage = "coherence-image"
	FlagUtilsImage     = "utils-image"
)

// The default CRD location
var (
	flagSet  *pflag.FlagSet
	cohFlags CoherenceOperatorFlags
)

func init() {
	flagSet = pflag.NewFlagSet("coh", pflag.ExitOnError)
	cohFlags = CoherenceOperatorFlags{}
	cohFlags.AddTo(flagSet)
}

// FlagSet - The Coherence flag set.
func FlagSet() *pflag.FlagSet {
	return flagSet
}

func GetFlags() *CoherenceOperatorFlags {
	return &cohFlags
}

// CoherenceOperatorFlags - Options to be used by a Coherence operator.
type CoherenceOperatorFlags struct {
	// The directory where the Operator's CRD file are located.
	CrdFiles string
	// The host name that the ReST server binds to.
	RestHost string
	// The port that the ReST server binds to.
	RestPort int32
	// The service name that the operator ReST clients should use.
	ServiceName string
	// The service port that the operator ReST clients should use. If not set defaults to the same as the ReST port.
	ServicePort int32
	// The label to use to obtain the site value for a Node.
	SiteLabel string
	// The label to use to obtain the rack value for a Node.
	RackLabel string
	// The default Coherence image to use if one is not specified for a deployment.
	CoherenceImage string
	// The default Coherence Utils image to use if one is not specified for a deployment.
	CoherenceUtilsImage string
}

// AddTo - Add the reconcile period and watches file flags to the the flag-set
// helpTextPrefix will allow you add a prefix to default help text. Joined by a space.
func (f *CoherenceOperatorFlags) AddTo(flagSet *pflag.FlagSet) {
	flagSet.StringVar(&cohFlags.RestHost,
		FlagRestHost,
		DefaultRestHost,
		"The address that the ReST server will bind to",
	)
	flagSet.Int32Var(&cohFlags.RestPort,
		FlagRestPort,
		DefaultRestPort,
		"The port that the ReST server will bind to",
	)
	flagSet.StringVar(&cohFlags.ServiceName,
		FlagServiceName,
		"",
		"The service name that operator clients use as the host name to make ReST calls back to the operator.",
	)
	flagSet.Int32Var(&cohFlags.ServicePort,
		FlagServicePort,
		-1,
		"The service port that operator clients use in the host name to make ReST calls back to the operator. If not set defaults to the same as the ReST port",
	)
	flagSet.StringVar(&cohFlags.SiteLabel,
		FlagSiteLabel,
		DefaultSiteLabel,
		"The node label to use when obtaining a value for a Pod's Coherence site.",
	)
	flagSet.StringVar(&cohFlags.RackLabel,
		FlagRackLabel,
		DefaultRackLabel,
		"The node label to use when obtaining a value for a Pod's Coherence rack.",
	)

	cohImg := os.Getenv(coherenceImageEnv)
	flagSet.StringVar(&cohFlags.CoherenceImage,
		FlagCoherenceImage,
		cohImg,
		"The Coherence image to use if one is not specified for a deployment.",
	)

	utilsImg := os.Getenv(utilsImageEnv)
	flagSet.StringVar(&cohFlags.CoherenceUtilsImage,
		FlagUtilsImage,
		utilsImg,
		"The Coherence Utils image to use if one is not specified for a deployment.",
	)
}

func GetDefaultCoherenceImage() string {
	img, ok := os.LookupEnv(coherenceImageEnv)
	if ok {
		return img
	}
	return ""
}

func (f *CoherenceOperatorFlags) GetCoherenceImage() string {
	if f.CoherenceImage != "" {
		return f.CoherenceImage
	}
	return GetDefaultCoherenceImage()
}

func GetDefaultCoherenceUtilsImage() string {
	img, ok := os.LookupEnv(utilsImageEnv)
	if ok {
		return img
	}
	return ""
}

func (f *CoherenceOperatorFlags) GetCoherenceUtilsImage() string {
	if f.CoherenceUtilsImage != "" {
		return f.CoherenceUtilsImage
	}
	return GetDefaultCoherenceUtilsImage()
}

// GetOperatorFlags returns the Operator command line flags.
func GetOperatorFlags() *CoherenceOperatorFlags {
	return &cohFlags
}
