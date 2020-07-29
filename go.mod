module github.com/oracle/coherence-operator

go 1.14

require (
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/coreos/prometheus-operator v0.38.1-0.20200424145508-7e176fda06cc
	github.com/elastic/go-elasticsearch/v7 v7.6.0
	github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32
	github.com/go-logr/logr v0.1.0
	github.com/go-test/deep v1.0.3
	github.com/onsi/gomega v1.10.1
	github.com/operator-framework/operator-sdk v0.19.0
	github.com/pborman/uuid v1.2.0
	github.com/pkg/errors v0.9.1
	github.com/shurcooL/vfsgen v0.0.0-20200627165143-92b8a710ab6c
	github.com/spf13/pflag v1.0.5
	github.com/tebeka/go2xunit v1.4.10
	golang.org/x/net v0.0.0-20200520004742-59133d7f0dd7
	gotest.tools/gotestsum v0.5.2 // indirect
	helm.sh/helm/v3 v3.2.4
	k8s.io/api v0.18.2
	k8s.io/apiextensions-apiserver v0.18.2
	k8s.io/apimachinery v0.18.2
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/helm v2.16.7+incompatible
	k8s.io/klog v1.0.0
	k8s.io/kubectl v0.18.2
	k8s.io/kubernetes v1.13.0
	k8s.io/utils v0.0.0-20200324210504-a9aa75ae1b89
	sigs.k8s.io/controller-runtime v0.6.0
	sigs.k8s.io/testing_frameworks v0.1.2
	sigs.k8s.io/yaml v1.2.0
)

replace (
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v13.3.2+incompatible // Required by OLM
	k8s.io/client-go => k8s.io/client-go v0.18.2
)
