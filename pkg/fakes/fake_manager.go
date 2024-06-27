/*
 * Copyright (c) 2019, 2024, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package fakes

import (
	"context"
	"github.com/go-logr/logr"
	. "github.com/onsi/gomega"
	coh "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/api/meta/testrestmapper"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	apitypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"
	"net/http"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/config"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// NewFakeManager creates a fake manager.Manager for use when testing controllers.
func NewFakeManager(initObjs ...runtime.Object) (*FakeManager, error) {
	gv := schema.GroupVersion{Group: "coherence.oracle.com", Version: "v1"}

	s := scheme.Scheme
	s.AddKnownTypes(gv, &coh.Coherence{})
	s.AddKnownTypes(gv, &coh.CoherenceList{})
	s.AddKnownTypes(gv, &coh.CoherenceJob{})
	s.AddKnownTypes(gv, &coh.CoherenceJobList{})

	cfg, _, err := helper.GetKubeconfigAndNamespace("")
	if err != nil {
		return nil, err
	}

	// Use a TestOnlyStaticRESTMapper so that the tests can run without needing a real k8s server
	restMapper := func(c *rest.Config, httpClient *http.Client) (meta.RESTMapper, error) {
		return testrestmapper.TestOnlyStaticRESTMapper(s), nil
	}

	options := manager.Options{
		MapperProvider: restMapper,
		LeaderElection: false,
		NewCache: func(config *rest.Config, opts cache.Options) (cache.Cache, error) {
			opts.DefaultNamespaces = map[string]cache.Config{
				helper.GetTestNamespace(): {},
			}
			return cache.New(config, opts)
		},
	}

	// Create the mapper provider
	mapper, err := options.MapperProvider(cfg, &http.Client{})
	if err != nil {
		return nil, err
	}

	cl := NewFakeClient(initObjs...)

	mgr := FakeManager{Scheme: s, Client: cl, Events: NewFakeEventRecorder(20), Mapper: mapper, Config: cfg}

	return &mgr, nil
}

var _ manager.Manager = &FakeManager{}

type FakeManager struct {
	Scheme *runtime.Scheme
	Client ClientWithErrors
	Events *FakeEventRecorder
	Mapper meta.RESTMapper
	Config *rest.Config
}

func (f *FakeManager) GetHTTPClient() *http.Client {
	return &http.Client{}
}

func (f *FakeManager) GetControllerOptions() config.Controller {
	return config.Controller{}
}

func (f *FakeManager) Start(ctx context.Context) error {
	return nil
}

func (f *FakeManager) GetLogger() logr.Logger {
	return ctrl.Log.WithName("manager")
}

func (f *FakeManager) Elected() <-chan struct{} {
	panic("implement me")
}

func (f *FakeManager) AddMetricsServerExtraHandler(path string, handler http.Handler) error {
	panic("implement me")
}

func (f *FakeManager) AddMetricsExtraHandler(path string, handler http.Handler) error {
	panic("implement me")
}

func (f *FakeManager) GetEventRecorderFor(name string) record.EventRecorder {
	return f.Events
}

func (f *FakeManager) GetAPIReader() client.Reader {
	panic("implement me")
}

func (f *FakeManager) GetWebhookServer() webhook.Server {
	panic("implement me")
}

type ReconcileResult struct {
	Result reconcile.Result
	Error  error
}

// Reset creates a new client and event recorder for this manager.
func (f *FakeManager) Reset(initObjs ...runtime.Object) {
	f.Client = NewFakeClient(initObjs...)
	f.Events = NewFakeEventRecorder(20)
}

func (f *FakeManager) AddHealthzCheck(name string, check healthz.Checker) error {
	return nil
}

func (f *FakeManager) AddReadyzCheck(name string, check healthz.Checker) error {
	return nil
}

func (f *FakeManager) Add(manager.Runnable) error {
	return nil
}

func (f *FakeManager) SetFields(interface{}) error {
	return nil
}

func (f *FakeManager) GetConfig() *rest.Config {
	return f.Config
}

func (f *FakeManager) GetScheme() *runtime.Scheme {
	return f.Scheme
}

func (f *FakeManager) GetClient() client.Client {
	return f.Client
}

func (f *FakeManager) GetFieldIndexer() client.FieldIndexer {
	panic("fake method not implemented")
}

func (f *FakeManager) GetCache() cache.Cache {
	return fakeCache{}
}

func (f *FakeManager) GetRESTMapper() meta.RESTMapper {
	return f.Mapper
}

// AssertEvent asserts that there is an event in the event channel and returns it.
func (f *FakeManager) AssertEvent() FakeEvent {
	event, found := f.NextEvent()
	Expect(found).To(BeTrue())
	return event
}

// AssertNoRemainingEvents asserts that there is an event in the event channel and returns it.
func (f *FakeManager) AssertNoRemainingEvents() {
	_, found := f.NextEvent()
	Expect(found).To(BeFalse())
}

// NextEvent returns the next event available in the event channel.
// If there is an event present then it is returned along with a bool of true.
// If the channel is empty then nil and false are returned.
func (f *FakeManager) NextEvent() (FakeEvent, bool) {
	var ok bool
	var item FakeEvent

	select {
	case item = <-f.Events.Events:
		ok = true
	default:
		ok = false
	}
	// at this point, "ok" is:
	//   true  => item was popped off the queue (or queue was closed, see below)
	//   false => not popped, would have blocked because of queue empty
	return item, ok
}

// GetService obtains the specified service
func (f *FakeManager) GetService(namespace, name string) (*v1.Service, error) {
	service := &v1.Service{}
	err := f.Client.Get(context.TODO(), apitypes.NamespacedName{Namespace: namespace, Name: name}, service)
	return service, err
}

// ServiceExists determines whether a service exists
func (f *FakeManager) ServiceExists(namespace, name string) bool {
	service := &v1.Service{}
	err := f.Client.Get(context.TODO(), apitypes.NamespacedName{Namespace: namespace, Name: name}, service)
	return err == nil
}

// AssertCoherences asserts that the specified number of Coherence resources exist in a namespace
func (f *FakeManager) AssertCoherences(namespace string, count int) {
	list := &coh.CoherenceList{}

	err := f.Client.List(context.TODO(), list, client.InNamespace(namespace))
	Expect(err).To(BeNil())
	Expect(len(list.Items)).To(Equal(count))
}

// AssertCoherenceExists asserts that the specified Coherence resource exists in the namespace and returns it
func (f *FakeManager) AssertCoherenceExists(namespace, name string) *coh.Coherence {
	role := &coh.Coherence{}
	err := f.Client.Get(context.TODO(), apitypes.NamespacedName{Namespace: namespace, Name: name}, role)
	Expect(err).NotTo(HaveOccurred())
	return role
}

func (f *FakeManager) GetCoherences(namespace string) (coh.CoherenceList, error) {
	list := coh.CoherenceList{}

	err := f.Client.List(context.TODO(), &list, client.InNamespace(namespace))
	return list, err
}

// AssertCoherenceRoles asserts that the specified number of CoherenceRole resources exist in a namespace
func (f *FakeManager) AssertCoherenceRoles(namespace string, count int) {
	list, err := f.GetCoherences(namespace)
	Expect(err).NotTo(HaveOccurred())
	Expect(len(list.Items)).To(Equal(count))
}

// AssertCoherencesForCluster asserts that the specified number of Coherence resources exist in a namespace for a cluster name
func (f *FakeManager) AssertCoherencesForCluster(namespace, clusterName string, count int) {
	list := &coh.CoherenceList{}

	labels := client.MatchingLabels{}
	labels[coh.LabelCoherenceCluster] = clusterName

	err := f.Client.List(context.TODO(), list, client.InNamespace(namespace), labels)
	Expect(err).To(BeNil())
	Expect(len(list.Items)).To(Equal(count))
}

// AssertWkaService asserts that a headless service to use for WKA exists for a given cluster in a namespace.
func (f *FakeManager) AssertWkaService(namespace string, deployment *coh.Coherence) {
	service, err := f.GetService(namespace, deployment.GetWkaServiceName())
	Expect(err).NotTo(HaveOccurred())
	Expect(service).NotTo(BeNil())
	Expect(service.Spec.Selector[coh.LabelCoherenceCluster]).To(Equal(deployment.Name))
}

var _ cache.Cache = fakeCache{}

type fakeCache struct {
}

func (f fakeCache) RemoveInformer(ctx context.Context, obj client.Object) error {
	panic("implement me")
}

func (f fakeCache) GetInformer(ctx context.Context, obj client.Object, opts ...cache.InformerGetOption) (cache.Informer, error) {
	panic("implement me")
}

func (f fakeCache) GetInformerForKind(ctx context.Context, gvk schema.GroupVersionKind, opts ...cache.InformerGetOption) (cache.Informer, error) {
	panic("implement me")
}

func (f fakeCache) Get(ctx context.Context, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
	panic("implement me")
}

func (f fakeCache) List(ctx context.Context, list client.ObjectList, opts ...client.ListOption) error {
	panic("implement me")
}

func (f fakeCache) Start(ctx context.Context) error {
	panic("implement me")
}

func (f fakeCache) WaitForCacheSync(ctx context.Context) bool {
	panic("implement me")
}

func (f fakeCache) IndexField(ctx context.Context, obj client.Object, field string, extractValue client.IndexerFunc) error {
	panic("implement me")
}
