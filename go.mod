module github.com/oracle/coherence-operator

require (
	github.com/ghodss/yaml v1.0.0
	github.com/go-logr/logr v0.1.0
	github.com/go-openapi/loads v0.17.0
	github.com/go-openapi/spec v0.19.0
	github.com/go-openapi/strfmt v0.18.0
	github.com/go-openapi/validate v0.18.0
	github.com/onsi/ginkgo v1.7.0
	github.com/onsi/gomega v1.4.3
	github.com/operator-framework/operator-sdk v0.9.1-0.20190724001845-d6e1aba9fa51
	github.com/pborman/uuid v1.2.0
	github.com/pkg/errors v0.8.1
	github.com/spf13/pflag v1.0.3
	github.com/tebeka/go2xunit v1.4.10
	golang.org/x/net v0.0.0-20190404232315-eb5bcb51f2a3
	golang.org/x/tools v0.0.0-20190408170212-12dd9f86f350
	k8s.io/api v0.0.0-20190612125737-db0771252981
	k8s.io/apiextensions-apiserver v0.0.0-20190228180357-d002e88f6236
	k8s.io/apimachinery v0.0.0-20190612125636-6a5db36e93ad
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/helm v2.13.1+incompatible
	k8s.io/klog v0.3.1
	k8s.io/kube-openapi v0.0.0-20190603182131-db7b694dc208 // indirect
	k8s.io/kubernetes v1.11.8-beta.0.0.20190124204751-3a10094374f2
	k8s.io/utils v0.0.0-20190308190857-21c4ce38f2a7
	sigs.k8s.io/controller-runtime v0.1.12
	sigs.k8s.io/controller-tools v0.1.10
	sigs.k8s.io/yaml v1.1.0
)

// Pinned to kubernetes-1.13.4
replace (
	k8s.io/api => k8s.io/api v0.0.0-20190222213804-5cb15d344471
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.0.0-20190228180357-d002e88f6236
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20190221213512-86fb29eff628
	k8s.io/apiserver => k8s.io/apiserver v0.0.0-20190228174905-79427f02047f
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.0.0-20190228180923-a9e421a79326
	k8s.io/client-go => k8s.io/client-go v0.0.0-20190228174230-b40b2a5939e4
	k8s.io/code-generator => k8s.io/code-generator v0.0.0-20181117043124-c2090bec4d9b
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.0.0-20190228175259-3e0149950b0e
	k8s.io/kube-openapi => k8s.io/kube-openapi v0.0.0-20180711000925-0cf8f7e6ed1d
	k8s.io/kubernetes => k8s.io/kubernetes v1.13.4
)

replace (
	github.com/coreos/prometheus-operator => github.com/coreos/prometheus-operator v0.29.0
	k8s.io/kube-state-metrics => k8s.io/kube-state-metrics v1.6.0
	sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.1.12
	sigs.k8s.io/controller-tools => sigs.k8s.io/controller-tools v0.1.11-0.20190411181648-9d55346c2bde
)

replace github.com/operator-framework/operator-sdk => github.com/operator-framework/operator-sdk v0.9.0
