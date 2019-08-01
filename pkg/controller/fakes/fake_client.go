package fakes

import (
	"context"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

type FakeClient interface {
	client.Client
	EnableErrors(errors ClientErrors)
	DisableErrors()
}

type clientWithErrors struct {
	wrapped  client.Client
	errorsOn bool
	errors   ClientErrors
}

type ClientErrors struct {
	getErrors    map[client.ObjectKey]error
	createErrors map[client.ObjectKey]error
	updateErrors map[client.ObjectKey]error
	deleteErrors map[client.ObjectKey]error
}

type ErrorOpts struct {
	Key  client.ObjectKey
	Type runtime.Object
}

type FakeError struct {
	Msg string
}

func (f FakeError) Error() string {
	return f.Msg
}

var _ client.Client = &clientWithErrors{}
var _ FakeClient = &clientWithErrors{}

func NewFakeClient(initObjs ...runtime.Object) FakeClient {
	c := clientWithErrors{wrapped: fake.NewFakeClient(initObjs...)}
	return &c
}

func (c *clientWithErrors) Get(ctx context.Context, key client.ObjectKey, obj runtime.Object) error {
	if c.errorsOn {
		if err := c.errors.IsGetError(key); err != nil {
			return err
		}
	}
	return c.wrapped.Get(ctx, key, obj)
}

func (c *clientWithErrors) List(ctx context.Context, opts *client.ListOptions, list runtime.Object) error {
	return c.wrapped.List(ctx, opts, list)
}

func (c *clientWithErrors) Create(ctx context.Context, obj runtime.Object) error {
	if c.errorsOn {
		accessor, _ := meta.Accessor(obj)
		key := types.NamespacedName{Namespace: accessor.GetNamespace(), Name: accessor.GetName()}
		if err := c.errors.IsGetError(key); err != nil {
			return err
		}
	}
	return c.wrapped.Create(ctx, obj)
}

func (c *clientWithErrors) Delete(ctx context.Context, obj runtime.Object, opts ...client.DeleteOptionFunc) error {
	if c.errorsOn {
		accessor, _ := meta.Accessor(obj)
		key := types.NamespacedName{Namespace: accessor.GetNamespace(), Name: accessor.GetName()}
		if err := c.errors.IsDeleteError(key); err != nil {
			return err
		}
	}
	return c.wrapped.Delete(ctx, obj, opts...)
}

func (c *clientWithErrors) Update(ctx context.Context, obj runtime.Object) error {
	if c.errorsOn {
		accessor, _ := meta.Accessor(obj)
		key := types.NamespacedName{Namespace: accessor.GetNamespace(), Name: accessor.GetName()}
		if err := c.errors.IsUpdateError(key); err != nil {
			return err
		}
	}
	return c.wrapped.Update(ctx, obj)
}

func (c *clientWithErrors) Status() client.StatusWriter {
	return c.wrapped.Status()
}

func (c *clientWithErrors) EnableErrors(errors ClientErrors) {
	c.errors = errors
	c.errorsOn = true
}

func (c *clientWithErrors) DisableErrors() {
	c.errorsOn = false
}

func (c *ClientErrors) IsCreateError(key client.ObjectKey) error {
	err, found := c.createErrors[key]
	if found {
		return err
	}
	return nil
}

func (c *ClientErrors) AddCreateError(key client.ObjectKey, err error) {
	if c.createErrors == nil {
		c.createErrors = make(map[client.ObjectKey]error)
	}
	c.createErrors[key] = err
}

func (c *ClientErrors) IsGetError(key client.ObjectKey) error {
	err, found := c.getErrors[key]
	if found {
		return err
	}
	return nil
}

func (c *ClientErrors) AddGetError(key client.ObjectKey, err error) {
	if c.getErrors == nil {
		c.getErrors = make(map[client.ObjectKey]error)
	}
	c.getErrors[key] = err
}

func (c *ClientErrors) IsUpdateError(key client.ObjectKey) error {
	err, found := c.updateErrors[key]
	if found {
		return err
	}
	return nil
}

func (c *ClientErrors) AddUpdateError(key client.ObjectKey, err error) {
	if c.updateErrors == nil {
		c.updateErrors = make(map[client.ObjectKey]error)
	}
	c.updateErrors[key] = err
}

func (c *ClientErrors) IsDeleteError(key client.ObjectKey) error {
	err, found := c.deleteErrors[key]
	if found {
		return err
	}
	return nil
}

func (c *ClientErrors) AddDeleteError(key client.ObjectKey, err error) {
	if c.deleteErrors == nil {
		c.deleteErrors = make(map[client.ObjectKey]error)
	}
	c.deleteErrors[key] = err
}
