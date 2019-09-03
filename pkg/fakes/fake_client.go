/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package fakes

import (
	"context"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

// ClientWithErrors is a client.Client that can be configured
// to return errors from calls.
type ClientWithErrors interface {
	client.Client
	EnableErrors(errors ClientErrors)
	DisableErrors()
}

// clientWithErrors an internal implementation of ClientWithErrors
type clientWithErrors struct {
	wrapped  client.Client
	errorsOn bool
	errors   ClientErrors
}

// ClientErrors is the configuration used by ClientWithErrors to
// decide whether or not to return an error from a method call.
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

// matches returns true if this ErrorOpts is a match for the key and object.
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
	c := clientWithErrors{wrapped: fake.NewFakeClient(initObjs...)}
	return &c
}

func (c *clientWithErrors) Get(ctx context.Context, key client.ObjectKey, obj runtime.Object) error {
	if c.errorsOn {
		if err := c.errors.IsGetError(key, obj); err != nil {
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
		if err := c.errors.IsGetError(key, obj); err != nil {
			return err
		}
	}
	return c.wrapped.Create(ctx, obj)
}

func (c *clientWithErrors) Delete(ctx context.Context, obj runtime.Object, opts ...client.DeleteOptionFunc) error {
	if c.errorsOn {
		accessor, _ := meta.Accessor(obj)
		key := types.NamespacedName{Namespace: accessor.GetNamespace(), Name: accessor.GetName()}
		if err := c.errors.IsDeleteError(key, obj); err != nil {
			return err
		}
	}
	return c.wrapped.Delete(ctx, obj, opts...)
}

func (c *clientWithErrors) Update(ctx context.Context, obj runtime.Object) error {
	if c.errorsOn {
		accessor, _ := meta.Accessor(obj)
		key := types.NamespacedName{Namespace: accessor.GetNamespace(), Name: accessor.GetName()}
		if err := c.errors.IsUpdateError(key, obj); err != nil {
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
