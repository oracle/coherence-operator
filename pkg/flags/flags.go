package flags

import (
	"github.com/operator-framework/operator-sdk/pkg/log/zap"
	"github.com/spf13/pflag"
	"strings"
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
	// The name of the image to use for Operator Utilities, e.g. the backup Pod.
	ImageName string
	// keyFile is the name of the client key file in the k8s secret used by the Operator when connecting to a Coherence
	//  Pod's management endpoint if that Pod has SSL enabled.
	SSLKeyFile string
	// certFile is the name of the client certificate file in the k8s secret used by the Operator when connecting to a
	//  Coherence Pod's management endpoint if that Pod has SSL enabled.
	SSLCertFile string
	// caFile is the name of the cert file in the k8s secret for the certificate authority used by the Operator when
	//  connecting to a Coherence Pod's management endpoint if that Pod has SSL enabled.
	SSLCAFile string
}

// cohf is the struct containing the command line flags.
var cohf = &CoherenceOperatorFlags{}

// AddTo - Add the reconcile period and watches file flags to the the flag-set
// helpTextPrefix will allow you add a prefix to default help text. Joined by a space.
func (f *CoherenceOperatorFlags) AddTo(flagSet *pflag.FlagSet, helpTextPrefix ...string) {
	flagSet.StringVar(&f.RestHost,
		"rest-host",
		"0.0.0.0",
		strings.Join(append(helpTextPrefix, "The address that the ReST server will bind to"), " "),
	)
	flagSet.Int32Var(&f.RestPort,
		"rest-port",
		8000,
		strings.Join(append(helpTextPrefix, "The port that the ReST server will bind to"), " "),
	)
	flagSet.StringVar(&f.ServiceName,
		"service-name",
		"",
		strings.Join(append(helpTextPrefix, "The service name that operator clients use as the host name to make ReST calls back to the operator."), " "),
	)
	flagSet.Int32Var(&f.ServicePort,
		"service-port",
		-1,
		strings.Join(append(helpTextPrefix, "The service port that operator clients use in the host name to make ReST calls back to the operator. If not set defaults to the same as the ReST port"), " "),
	)
	flagSet.BoolVar(&f.LogIntegrationEnabled,
		"log-integration",
		false,
		strings.Join(append(helpTextPrefix, "A boolean indicating whether logging integration (e.g. EFK) is enabled"), " "),
	)
	flagSet.StringVar(&f.ElasticSearchHost,
		"es-host",
		"",
		strings.Join(append(helpTextPrefix, "The host name of the ElasticSearch server"), " "),
	)
	flagSet.Int32Var(&f.ElasticSearchPort,
		"es-port",
		-1,
		strings.Join(append(helpTextPrefix, "The port to use to connect to the ElasticSearch server"), " "),
	)
	flagSet.StringVar(&f.ElasticSearchUser,
		"es-user",
		"",
		strings.Join(append(helpTextPrefix, "The user name to use to connect to the ElasticSearch server"), " "),
	)
	flagSet.StringVar(&f.ElasticSearchUser,
		"es-password",
		"",
		strings.Join(append(helpTextPrefix, "The credentials to use to connect to the ElasticSearch server"), " "),
	)
	flagSet.StringVar(&f.ImageName,
		"utils-image",
		"",
		strings.Join(append(helpTextPrefix, "The name of the Operator Utils Docker image"), " "),
	)

	flagSet.StringVar(&f.SSLKeyFile,
		"ssl-key-file",
		"",
		strings.Join(append(helpTextPrefix, "The name of the client key file in the k8s secret used by the Operator when connecting to a Coherence Pod's management endpoint if that Pod has SSL enabled"), " "),
	)

	flagSet.StringVar(&f.SSLCertFile,
		"ssl-cert-file",
		"",
		strings.Join(append(helpTextPrefix, "The name of the client certificate file in the k8s secret used by the Operator when connecting to a Coherence Pod's management endpoint if that Pod has SSL enabled"), " "),
	)

	flagSet.StringVar(&f.SSLCAFile,
		"ssl-ca-file",
		"",
		strings.Join(append(helpTextPrefix, "The name of the cert file in the k8s secret for the certificate authority used by the Operator when connecting to a Coherence Pod's management endpoint if that Pod has SSL enabled"), " "),
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
