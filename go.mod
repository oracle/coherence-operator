module github.com/oracle/coherence-operator

go 1.13

require (
	github.com/coreos/prometheus-operator v0.38.0
	github.com/elastic/go-elasticsearch/v7 v7.6.0
	github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32
	github.com/go-logr/logr v0.1.0
	github.com/go-openapi/spec v0.19.4
	github.com/go-openapi/validate v0.19.5
	github.com/go-test/deep v1.0.3
	github.com/onsi/gomega v1.8.1
	github.com/operator-framework/operator-sdk v0.17.0
	github.com/pborman/uuid v1.2.0
	github.com/pkg/errors v0.9.1
	github.com/spf13/pflag v1.0.5
	github.com/tebeka/go2xunit v1.4.10
	golang.org/x/net v0.0.0-20200226121028-0de0cce0169b
	helm.sh/helm/v3 v3.1.2
	k8s.io/api v0.17.4
	k8s.io/apiextensions-apiserver v0.17.4
	k8s.io/apimachinery v0.17.4
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/helm v2.16.3+incompatible
	k8s.io/klog v1.0.0
	k8s.io/kube-openapi v0.0.0-20191107075043-30be4d16710a
	k8s.io/kubectl v0.17.4
	k8s.io/kubernetes v1.13.0
	k8s.io/utils v0.0.0-20191114200735-6ca3b61696b6
	sigs.k8s.io/controller-runtime v0.5.2
	sigs.k8s.io/testing_frameworks v0.1.2
	sigs.k8s.io/yaml v1.1.0
)

replace (
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v13.3.2+incompatible // Required by OLM
	k8s.io/client-go => k8s.io/client-go v0.17.4 // Required by prometheus-operator
)
