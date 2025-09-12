/*
 * Copyright (c) 2020, 2025, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package fakes

import (
	"context"
	"reflect"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

// ClientWithErrors is a client.Client that can be configured
// to return errors from calls.
type ClientWithErrors interface {
	client.Client
	EnableErrors(errors ClientErrors)
	DisableErrors()
	GetOperations() []ClientOperation
	GetCreates() []runtime.Object
	GetUpdates() []runtime.Object
	GetDeletes() []runtime.Object
	GetPatches() []ClientOperation
	GetStatefulSet(namespace, name string) (*appsv1.StatefulSet, error)
	GetJob(namespace, name string) (*batchv1.Job, error)
	GetService(namespace, name string) (*corev1.Service, error)
}

// clientWithErrors an internal implementation of ClientWithErrors
type clientWithErrors struct {
	wrapped    client.Client
	errorsOn   bool
	errors     ClientErrors
	operations []ClientOperation
}

func (c *clientWithErrors) Apply(ctx context.Context, obj runtime.ApplyConfiguration, opts ...client.ApplyOption) error {
	return c.wrapped.Apply(ctx, obj, opts...)
}

func (c *clientWithErrors) GroupVersionKindFor(obj runtime.Object) (schema.GroupVersionKind, error) {
	return c.wrapped.GroupVersionKindFor(obj)
}

func (c *clientWithErrors) IsObjectNamespaced(obj runtime.Object) (bool, error) {
	return c.wrapped.IsObjectNamespaced(obj)
}

// ClientErrors is the configuration used by ClientWithErrors to
// decide whether to return an error from a method call.
type ClientErrors struct {
	getErrors    map[ErrorOpts]error
	createErrors map[ErrorOpts]error
	updateErrors map[ErrorOpts]error
	deleteErrors map[ErrorOpts]error
}

// ErrorOpts is used to determine whether a particular client
// method call should return an error.
// If all fields are nil then all calls will match.
type ErrorOpts interface {
	Matches(key client.ObjectKey, obj runtime.Object) bool
}

// ErrorIf is a simple implementation of ErrorOpts
type ErrorIf struct {
	// The optional key to match when deciding whether to return an error.
	KeyIs *client.ObjectKey
	// The optional type to match when deciding whether to return an error.
	TypeIs runtime.Object
}

// Matches returns true if this ErrorOpts is a match for the key and object.
func (o ErrorIf) Matches(key client.ObjectKey, obj runtime.Object) bool {
	if o.KeyIs == nil && o.TypeIs == nil {
		return true
	}

	var keyMatch bool
	var typeMatch bool

	if o.KeyIs != nil {
		keyMatch = o.KeyIs.Namespace == key.Namespace && o.KeyIs.Name == key.Name
	} else {
		keyMatch = true
	}

	if o.TypeIs != nil {
		t1 := reflect.TypeOf(o.TypeIs)
		t2 := reflect.TypeOf(obj)
		typeMatch = t1.String() == t2.String()
	} else {
		typeMatch = true
	}

	return keyMatch && typeMatch
}

// FakeError is a simple error implementation.
type FakeError struct {
	Msg string
}

func (f FakeError) Error() string {
	return f.Msg
}

// Two internal vars to make sure that we actually do implement the client.Client
// interface. We will get compile erros if we do not.
var _ client.Client = &clientWithErrors{}
var _ ClientWithErrors = &clientWithErrors{}

// NewFakeClient creates a new ClientWithErrors and initialises it
// with the specified objects.
func NewFakeClient(initObjs ...runtime.Object) ClientWithErrors {
	var subs []client.Object
	for _, o := range initObjs {
		subs = append(subs, o.(client.Object))
	}
	cl := fake.NewClientBuilder().WithScheme(scheme.Scheme).
		WithRuntimeObjects(initObjs...).
		WithStatusSubresource(subs...).
		Build()

	c := clientWithErrors{wrapped: cl}
	return &c
}

func (c *clientWithErrors) Scheme() *runtime.Scheme {
	return c.wrapped.Scheme()
}

func (c *clientWithErrors) RESTMapper() meta.RESTMapper {
	return c.wrapped.RESTMapper()
}

func (c *clientWithErrors) Get(ctx context.Context, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
	if c.errorsOn {
		if err := c.errors.IsGetError(key, obj); err != nil {
			return err
		}
	}
	return c.wrapped.Get(ctx, key, obj)
}

func (c *clientWithErrors) GetStatefulSet(namespace, name string) (*appsv1.StatefulSet, error) {
	sts := &appsv1.StatefulSet{}
	err := c.Get(context.TODO(), types.NamespacedName{Namespace: namespace, Name: name}, sts)
	return sts, err
}

func (c *clientWithErrors) GetJob(namespace, name string) (*batchv1.Job, error) {
	job := &batchv1.Job{}
	err := c.Get(context.TODO(), types.NamespacedName{Namespace: namespace, Name: name}, job)
	return job, err
}

func (c *clientWithErrors) GetService(namespace, name string) (*corev1.Service, error) {
	svc := &corev1.Service{}
	err := c.Get(context.TODO(), types.NamespacedName{Namespace: namespace, Name: name}, svc)
	return svc, err
}

func (c *clientWithErrors) List(ctx context.Context, list client.ObjectList, opts ...client.ListOption) error {
	return c.wrapped.List(ctx, list, opts...)
}

func (c *clientWithErrors) Create(ctx context.Context, obj client.Object, opts ...client.CreateOption) error {
	if c.errorsOn {
		accessor, _ := meta.Accessor(obj)
		key := types.NamespacedName{Namespace: accessor.GetNamespace(), Name: accessor.GetName()}
		if err := c.errors.IsGetError(key, obj); err != nil {
			return err
		}
	}
	err := c.wrapped.Create(ctx, obj, opts...)
	if err == nil {
		if _, ok := obj.(metav1.Object); ok {
			mo := obj.(metav1.Object)
			mo.SetCreationTimestamp(metav1.Time{Time: time.Now()})
		}
		c.operations = append(c.operations, ClientOperation{Action: ClientActionCreate, Object: obj})
	}
	return err
}

func (c *clientWithErrors) Delete(ctx context.Context, obj client.Object, opts ...client.DeleteOption) error {
	if c.errorsOn {
		accessor, _ := meta.Accessor(obj)
		key := types.NamespacedName{Namespace: accessor.GetNamespace(), Name: accessor.GetName()}
		if err := c.errors.IsDeleteError(key, obj); err != nil {
			return err
		}
	}
	err := c.wrapped.Delete(ctx, obj, opts...)
	if err == nil {
		if _, ok := obj.(metav1.Object); ok {
			t := metav1.Time{Time: time.Now()}
			mo := obj.(metav1.Object)
			mo.SetDeletionTimestamp(&t)
			c.operations = append(c.operations, ClientOperation{Action: ClientActionDelete, Object: obj})
		}
	}
	return err
}

func (c *clientWithErrors) Update(ctx context.Context, obj client.Object, opts ...client.UpdateOption) error {
	if c.errorsOn {
		accessor, _ := meta.Accessor(obj)
		key := types.NamespacedName{Namespace: accessor.GetNamespace(), Name: accessor.GetName()}
		if err := c.errors.IsUpdateError(key, obj); err != nil {
			return err
		}
	}
	err := c.wrapped.Update(ctx, obj, opts...)
	if err == nil {
		if _, ok := obj.(metav1.Object); ok {
			mo := obj.(metav1.Object)
			mo.SetGeneration(obj.(metav1.Object).GetGeneration() + 1)
			c.operations = append(c.operations, ClientOperation{Action: ClientActionUpdate, Object: obj})
		}
	}
	return err
}

func (c *clientWithErrors) Patch(ctx context.Context, obj client.Object, patch client.Patch, opts ...client.PatchOption) error {
	err := c.wrapped.Patch(ctx, obj, patch, opts...)
	if err == nil {
		if _, ok := obj.(metav1.Object); ok {
			mo := obj.(metav1.Object)
			mo.SetGeneration(obj.(metav1.Object).GetGeneration() + 1)
			cpy := obj.DeepCopyObject()
			_ = c.wrapped.Get(context.TODO(), types.NamespacedName{Namespace: mo.GetNamespace(), Name: mo.GetName()}, cpy.(client.Object))
			c.operations = append(c.operations, ClientOperation{Action: ClientActionPatch, Object: cpy, Patch: patch})
		}
	}
	return err
}

func (c *clientWithErrors) Status() client.StatusWriter {
	return c.wrapped.Status()
}

func (c *clientWithErrors) DeleteAllOf(ctx context.Context, obj client.Object, opts ...client.DeleteAllOfOption) error {

	return c.wrapped.DeleteAllOf(ctx, obj, opts...)
}

func (c *clientWithErrors) EnableErrors(errors ClientErrors) {
	c.errors = errors
	c.errorsOn = true
}

func (c *clientWithErrors) DisableErrors() {
	c.errorsOn = false
}

func (c *clientWithErrors) GetOperations() []ClientOperation {
	var ops []ClientOperation
	if c != nil {
		ops = append(ops, c.operations...)
	}
	return ops
}

func (c *clientWithErrors) GetCreates() []runtime.Object {
	var objects []runtime.Object
	if c != nil {
		for _, op := range c.operations {
			if op.Action == ClientActionCreate {
				objects = append(objects, op.Object)
			}
		}
	}
	return objects
}

func (c *clientWithErrors) GetUpdates() []runtime.Object {
	var objects []runtime.Object
	if c != nil {
		for _, op := range c.operations {
			if op.Action == ClientActionUpdate {
				objects = append(objects, op.Object)
			}
		}
	}
	return objects
}

func (c *clientWithErrors) GetDeletes() []runtime.Object {
	var objects []runtime.Object
	if c != nil {
		for _, op := range c.operations {
			if op.Action == ClientActionDelete {
				objects = append(objects, op.Object)
			}
		}
	}
	return objects
}

func (c *clientWithErrors) GetPatches() []ClientOperation {
	var objects []ClientOperation
	if c != nil {
		for _, op := range c.operations {
			if op.Action == ClientActionPatch {
				objects = append(objects, op)
			}
		}
	}
	return objects
}

func (c *ClientErrors) isError(key client.ObjectKey, obj runtime.Object, errors map[ErrorOpts]error) error {
	for opts, err := range errors {
		if opts.Matches(key, obj) {
			return err
		}
	}
	return nil
}

func (c *ClientErrors) IsCreateError(key client.ObjectKey, obj runtime.Object) error {
	return c.isError(key, obj, c.createErrors)
}

func (c *ClientErrors) AddCreateError(opts ErrorOpts, err error) {
	if c.createErrors == nil {
		c.createErrors = make(map[ErrorOpts]error)
	}
	c.createErrors[opts] = err
}

func (c *ClientErrors) IsGetError(key client.ObjectKey, obj runtime.Object) error {
	return c.isError(key, obj, c.getErrors)
}

func (c *ClientErrors) AddGetError(opts ErrorOpts, err error) {
	if c.getErrors == nil {
		c.getErrors = make(map[ErrorOpts]error)
	}
	c.getErrors[opts] = err
}

func (c *ClientErrors) IsUpdateError(key client.ObjectKey, obj runtime.Object) error {
	return c.isError(key, obj, c.updateErrors)
}

func (c *ClientErrors) AddUpdateError(opts ErrorOpts, err error) {
	if c.updateErrors == nil {
		c.updateErrors = make(map[ErrorOpts]error)
	}
	c.updateErrors[opts] = err
}

func (c *ClientErrors) IsDeleteError(key client.ObjectKey, obj runtime.Object) error {
	return c.isError(key, obj, c.deleteErrors)
}

func (c *ClientErrors) AddDeleteError(opts ErrorOpts, err error) {
	if c.deleteErrors == nil {
		c.deleteErrors = make(map[ErrorOpts]error)
	}
	c.deleteErrors[opts] = err
}

func (c *clientWithErrors) SubResource(_ string) client.SubResourceClient {
	panic("implement me")
}

type ClientAction string

const (
	ClientActionCreate ClientAction = "Create"
	ClientActionUpdate ClientAction = "Update"
	ClientActionDelete ClientAction = "Delete"
	ClientActionPatch  ClientAction = "ThreeWayPatch"
)

type ClientOperation struct {
	Action ClientAction
	Object runtime.Object
	Patch  client.Patch
}
