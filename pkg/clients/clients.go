package clients

import (
	apiextensions "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
)

type ClientSet struct {
	KubeClient            kubernetes.Interface
	ExtClient             apiextensions.Interface
	DynamicClient         dynamic.Interface
	DiscoveryClient       *discovery.DiscoveryClient
}

func New() (ClientSet, error) {
	cfg, err := ctrl.GetConfig()
	if err != nil {
		return ClientSet{}, err
	}
	return NewForConfig(cfg)
}

func NewForConfig(cfg *rest.Config) (ClientSet, error) {
	var err error
	c := ClientSet{}
	c.KubeClient, err = kubernetes.NewForConfig(cfg)
	if err != nil {
		return c, err
	}
	c.ExtClient, err = apiextensions.NewForConfig(cfg)
	if err != nil {
		return c, err
	}
	c.DynamicClient, err = dynamic.NewForConfig(cfg)
	if err != nil {
		return c, err
	}
	c.DiscoveryClient, err = discovery.NewDiscoveryClientForConfig(cfg)
	if err != nil {
		return c, err
	}
	return c, nil
}