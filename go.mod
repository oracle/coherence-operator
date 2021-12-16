module github.com/oracle/coherence-operator

go 1.16

require (
	github.com/coreos/prometheus-operator v0.38.1-0.20200424145508-7e176fda06cc
	github.com/elastic/go-elasticsearch/v7 v7.6.0
	github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32
	github.com/go-logr/logr v0.4.0
	github.com/go-test/deep v1.0.7
	github.com/onsi/gomega v1.15.0
	github.com/operator-framework/operator-lib v0.1.0
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.1.3
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.1
	golang.org/x/net v0.0.0-20210520170846-37e1c6afe023
	k8s.io/api v0.22.2
	k8s.io/apiextensions-apiserver v0.22.2
	k8s.io/apimachinery v0.22.2
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/kubernetes v1.13.0
	k8s.io/utils v0.0.0-20210819203725-bdf08cb9a70a
	sigs.k8s.io/controller-runtime v0.10.0
	sigs.k8s.io/testing_frameworks v0.1.2
)

replace k8s.io/client-go => k8s.io/client-go v0.22.2
