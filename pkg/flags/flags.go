/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package flags

import (
	"github.com/operator-framework/operator-sdk/pkg/log/zap"
	"github.com/spf13/pflag"
	"strings"
)

const (
	DefaultSiteLabel       = "failure-domain.beta.kubernetes.io/zone"
	DefaultRackLabel       = "failure-domain.beta.kubernetes.io/region"
	DefaultRestHost        = "0.0.0.0"
	DefaultRestPort  int32 = 8000

	FlagRestHost       = "rest-host"
	FlagRestPort       = "rest-port"
	FlagServiceName    = "service-name"
	FlagServicePort    = "service-port"
	FlagLogIntegration = "log-integration"
	FlagEsHost         = "es-host"
	FlagEsPort         = "es-port"
	FlagEsUser         = "es-user"
	FlagEsPassword     = "es-password"
	FlagSiteLabel      = "site-label"
	FlagRackLabel      = "rack-label"
	FlagAlwaysPullTags = "force-always-pull-tags"
)

// CoherenceOperatorFlags - Options to be used by a Coherence operator.
type CoherenceOperatorFlags struct {
	// The host name that the ReST server binds to.
	RestHost string
	// The port that the ReST server binds to.
	RestPort int32
	// The service name that the operator ReST clients should use.
	ServiceName string
	// The service port that the operator ReST clients should use. If not set defaults to the same as the ReST port.
	ServicePort int32
	// A flag indicating whether logging (e.g. ELK) integration is enabled.
	LogIntegrationEnabled bool
	// The host name for the ElasticSearch server.
	ElasticSearchHost string
	// The port for the ElasticSearch server.
	ElasticSearchPort int32
	// The username to use to connect to ElasticSearch.
	ElasticSearchUser string
	// The credentials to use to connect to ElasticSearch.
	ElasticSearchCredentials string
	// The label to use to obtain the site value for a Node.
	SiteLabel string
	// The label to use to obtain the rack value for a Node.
	RackLabel string
	// If any image names in the CoherenceCluster spec end with any suffix in the specified comma-delimited list the imagePullPolicy will be forced to ALWAYS.
	AlwaysPullSuffixes string
}

// cohf is the struct containing the command line flags.
var cohf = &CoherenceOperatorFlags{}

// AddTo - Add the reconcile period and watches file flags to the the flag-set
// helpTextPrefix will allow you add a prefix to default help text. Joined by a space.
func (f *CoherenceOperatorFlags) AddTo(flagSet *pflag.FlagSet, helpTextPrefix ...string) {
	flagSet.StringVar(&f.RestHost,
		FlagRestHost,
		DefaultRestHost,
		strings.Join(append(helpTextPrefix, "The address that the ReST server will bind to"), " "),
	)
	flagSet.Int32Var(&f.RestPort,
		FlagRestPort,
		DefaultRestPort,
		strings.Join(append(helpTextPrefix, "The port that the ReST server will bind to"), " "),
	)
	flagSet.StringVar(&f.ServiceName,
		FlagServiceName,
		"",
		strings.Join(append(helpTextPrefix, "The service name that operator clients use as the host name to make ReST calls back to the operator."), " "),
	)
	flagSet.Int32Var(&f.ServicePort,
		FlagServicePort,
		-1,
		strings.Join(append(helpTextPrefix, "The service port that operator clients use in the host name to make ReST calls back to the operator. If not set defaults to the same as the ReST port"), " "),
	)
	flagSet.BoolVar(&f.LogIntegrationEnabled,
		FlagLogIntegration,
		false,
		strings.Join(append(helpTextPrefix, "A boolean indicating whether logging integration (e.g. EFK) is enabled"), " "),
	)
	flagSet.StringVar(&f.ElasticSearchHost,
		FlagEsHost,
		"",
		strings.Join(append(helpTextPrefix, "The host name of the ElasticSearch server"), " "),
	)
	flagSet.Int32Var(&f.ElasticSearchPort,
		FlagEsPort,
		-1,
		strings.Join(append(helpTextPrefix, "The port to use to connect to the ElasticSearch server"), " "),
	)
	flagSet.StringVar(&f.ElasticSearchUser,
		FlagEsUser,
		"",
		strings.Join(append(helpTextPrefix, "The user name to use to connect to the ElasticSearch server"), " "),
	)
	flagSet.StringVar(&f.ElasticSearchCredentials,
		FlagEsPassword,
		"",
		strings.Join(append(helpTextPrefix, "The credentials to use to connect to the ElasticSearch server"), " "),
	)
	flagSet.StringVar(&f.SiteLabel,
		FlagSiteLabel,
		DefaultSiteLabel,
		strings.Join(append(helpTextPrefix, "The node label to use when obtaining a value for a Pod's Coherence site."), " "),
	)
	flagSet.StringVar(&f.RackLabel,
		FlagRackLabel,
		DefaultRackLabel,
		strings.Join(append(helpTextPrefix, "The node label to use when obtaining a value for a Pod's Coherence rack."), " "),
	)
	flagSet.StringVar(&f.AlwaysPullSuffixes,
		FlagAlwaysPullTags,
		"",
		strings.Join(append(helpTextPrefix, "If any image names in the CoherenceCluster spec end with any suffix in the specified comma-delimited list the imagePullPolicy will be forced to ALWAYS."), " "),
	)
}

// GetOperatorFlags returns the Operator command line flags.
func GetOperatorFlags() *CoherenceOperatorFlags {
	return cohf
}

// AddTo - Add the Coherence operator flags to the the flagset
// helpTextPrefix will allow you add a prefix to default help text. Joined by a space.
func AddTo(flagSet *pflag.FlagSet, helpTextPrefix ...string) *CoherenceOperatorFlags {
	cohf.AddTo(flagSet, helpTextPrefix...)
	flagSet.AddFlagSet(zap.FlagSet())
	return cohf
}
