/*
 * Copyright (c) 2019, 2020 Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package fakes

import (
	"context"
	. "github.com/onsi/gomega"
	coherence "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	apitypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// NewFakeManager creates a fake manager.Manager for use when testing controllers.
func NewFakeManager(initObjs ...runtime.Object) (*FakeManager, error) {
	gv := schema.GroupVersion{Group: "coherence.oracle.com", Version: "v1"}

	s := scheme.Scheme
	s.AddKnownTypes(gv, &coherence.CoherenceCluster{})
	s.AddKnownTypes(gv, &coherence.CoherenceClusterList{})
	s.AddKnownTypes(gv, &coherence.CoherenceRole{})
	s.AddKnownTypes(gv, &coherence.CoherenceRoleList{})

	gvk := coherence.GetCoherenceInternalGroupVersionKind(s)
	s.AddKnownTypeWithName(gvk, &unstructured.Unstructured{})

	//s.AddKnownTypes(gv, &coherence.CoherenceInternalList{})
	s.AddKnownTypeWithName(schema.GroupVersionKind{Group: gv.Group, Version: gv.Version, Kind: "CoherenceInternalList"}, &unstructured.UnstructuredList{})

	cfg, _, err := helper.GetKubeconfigAndNamespace("")
	if err != nil {
		return nil, err
	}

	options := manager.Options{
		Namespace:      helper.GetTestNamespace(),
		MapperProvider: restMapper,
		LeaderElection: false,
	}

	// Create the mapper provider
	mapper, err := options.MapperProvider(cfg)
	if err != nil {
		return nil, err
	}

	cl := NewFakeClient(initObjs...)

	mgr := FakeManager{Scheme: s, Client: cl, Events: NewFakeEventRecorder(20), Mapper: mapper, Config: cfg}

	return &mgr, nil
}

var restMapper = apiutil.NewDiscoveryRESTMapper

var _ manager.Manager = &FakeManager{}

type FakeManager struct {
	Scheme *runtime.Scheme
	Client ClientWithErrors
	Events *FakeEventRecorder
	Mapper meta.RESTMapper
	Config *rest.Config
}

func (f *FakeManager) GetEventRecorderFor(name string) record.EventRecorder {
	return f.Events
}

func (f *FakeManager) GetAPIReader() client.Reader {
	panic("implement me")
}

func (f *FakeManager) GetWebhookServer() *webhook.Server {
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
	panic("fake method not implemented")
}

func (f *FakeManager) AddReadyzCheck(name string, check healthz.Checker) error {
	panic("fake method not implemented")
}

func (f *FakeManager) Add(manager.Runnable) error {
	panic("fake method not implemented")
}

func (f *FakeManager) SetFields(interface{}) error {
	panic("fake method not implemented")
}

func (f *FakeManager) Start(<-chan struct{}) error {
	panic("fake method not implemented")
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
	panic("fake method not implemented")
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

// AssertEvent asserts that there is an event in the event channel and returns it.
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

// AssertCoherenceClusters asserts that the specified number of CoherenceCluster resources exist in a namespace
func (f *FakeManager) AssertCoherenceClusters(namespace string, count int) {
	list := &coherence.CoherenceClusterList{}

	err := f.Client.List(context.TODO(), list, client.InNamespace(namespace))
	Expect(err).To(BeNil())
	Expect(len(list.Items)).To(Equal(count))
}

// AssertCoherenceRoleExists asserts that the specified CoherenceRole exists in the namespace and returns it
func (f *FakeManager) AssertCoherenceRoleExists(namespace, name string) *coherence.CoherenceRole {
	role := &coherence.CoherenceRole{}
	err := f.Client.Get(context.TODO(), apitypes.NamespacedName{Namespace: namespace, Name: name}, role)
	Expect(err).NotTo(HaveOccurred())
	return role
}

// AssertCoherenceRoles asserts that the specified number of CoherenceRole resources exist in a namespace
func (f *FakeManager) GetCoherenceRoles(namespace string) (coherence.CoherenceRoleList, error) {
	list := coherence.CoherenceRoleList{}

	err := f.Client.List(context.TODO(), &list, client.InNamespace(namespace))
	return list, err
}

// AssertCoherenceRoles asserts that the specified number of CoherenceRole resources exist in a namespace
func (f *FakeManager) AssertCoherenceRoles(namespace string, count int) {
	list, err := f.GetCoherenceRoles(namespace)
	Expect(err).NotTo(HaveOccurred())
	Expect(len(list.Items)).To(Equal(count))
}

// AssertCoherenceRoles asserts that the specified number of CoherenceRole resources exist in a namespace for a cluster name
func (f *FakeManager) AssertCoherenceRolesForCluster(namespace, clusterName string, count int) {
	list := &coherence.CoherenceRoleList{}

	labels := client.MatchingLabels{}
	labels[coherence.CoherenceClusterLabel] = clusterName

	err := f.Client.List(context.TODO(), list, client.InNamespace(namespace), labels)
	Expect(err).To(BeNil())
	Expect(len(list.Items)).To(Equal(count))
}

// AssertCoherenceRoleExists asserts that the specified CoherenceRole exists in the namespace and returns it
func (f *FakeManager) AssertCoherenceInternalExists(namespace, name string) (*unstructured.Unstructured, error) {
	gvk := coherence.GetCoherenceInternalGroupVersionKind(f.Scheme)
	role := &unstructured.Unstructured{}
	role.SetGroupVersionKind(gvk)

	err := f.Client.Get(context.TODO(), apitypes.NamespacedName{Namespace: namespace, Name: name}, role)
	return role, err
}

// AssertCoherenceRoles asserts that the specified number of CoherenceRole resources exist in a namespace
func (f *FakeManager) AssertCoherenceInternals(namespace string, count int) {
	list := f.GetCoherenceInternals(namespace)
	Expect(len(list.Items)).To(Equal(count))
}

// GetCoherenceInternals obtains the CoherenceInternals for the specified namespace
func (f *FakeManager) GetCoherenceInternals(namespace string) unstructured.UnstructuredList {
	gvk := coherence.GetCoherenceInternalGroupVersionKind(f.Scheme)
	list := unstructured.UnstructuredList{}

	list.SetGroupVersionKind(gvk)
	list.SetKind("CoherenceInternalList")

	_ = f.Client.List(context.TODO(), &list, client.InNamespace(namespace))
	return list
}

// AssertWkaService asserts that a headless service to use for WKA exists for a given cluster in a namespace.
func (f *FakeManager) AssertWkaService(namespace string, cluster *coherence.CoherenceCluster) {
	service, err := f.GetService(namespace, cluster.GetWkaServiceName())
	Expect(err).NotTo(HaveOccurred())
	Expect(service).NotTo(BeNil())
	Expect(service.Spec.Selector[coherence.CoherenceClusterLabel]).To(Equal(cluster.Name))
}
