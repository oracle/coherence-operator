// The fakes package contains methods to create fakes and stubs for controller tests
package fakes

import (
	coherence "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission/types"
)

// NewFakeManager creates a fake manager.Manager for use when testing controllers.
func NewFakeManager(initObjs ...runtime.Object) *FakeManager {
	gvk := schema.GroupVersion{Group: "coherence.oracle.com", Version: "v1",}

	s := scheme.Scheme
	s.AddKnownTypes(gvk, &coherence.CoherenceCluster{})
	s.AddKnownTypes(gvk, &coherence.CoherenceClusterList{})
	s.AddKnownTypes(gvk, &coherence.CoherenceRole{})
	s.AddKnownTypes(gvk, &coherence.CoherenceRoleList{})
	s.AddKnownTypeWithName(schema.GroupVersionKind{Group: "coherence.oracle.com", Version: "v1", Kind: "CoherenceInternal"}, &unstructured.Unstructured{})

	cl := fake.NewFakeClient(initObjs...)

	return &FakeManager{Scheme: s, Client: cl, Events: NewFakeEventRecorder(20)}
}

type FakeManager struct {
	Scheme *runtime.Scheme
	Client client.Client
	Events *FakeEventRecorder
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
	panic("fake method not implemented")
}

func (f *FakeManager) GetScheme() *runtime.Scheme {
	return f.Scheme
}

func (f *FakeManager) GetAdmissionDecoder() types.Decoder {
	panic("fake method not implemented")
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

func (f *FakeManager) GetRecorder(name string) record.EventRecorder {
	return f.Events
}

func (f *FakeManager) GetRESTMapper() meta.RESTMapper {
	panic("fake method not implemented")
}

// NextEvent returns the next event available in the event channel.
// If there is an event present then it is returned along with a bool of true.
// If the channel is empty then nil and false are returned.
func (f *FakeManager) NextEvent() (FakeEvent, bool) {
	var ok bool
	var item FakeEvent

	select {
	    case item = <- f.Events.Events:
	        ok = true
	    default:
	        ok = false
	}
	// at this point, "ok" is:
	//   true  => item was popped off the queue (or queue was closed, see below)
	//   false => not popped, would have blocked because of queue empty
	return item, ok
}


