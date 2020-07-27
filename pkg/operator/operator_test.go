/*
 * Copyright (c) 2019, 2020 Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package operator_test

import (
	"context"
	. "github.com/onsi/gomega"
	v1 "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/pkg/fakes"
	crdv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	v1client "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1"
	v1beta1client "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"testing"
)

func TestShouldCreateV1CRDs(t *testing.T) {
	var err error

	g := NewGomegaWithT(t)
	mgr, err := fakes.NewFakeManager()
	g.Expect(err).NotTo(HaveOccurred())

	err = crdv1.AddToScheme(mgr.Scheme)
	g.Expect(err).NotTo(HaveOccurred())

	log := fakes.TestLogger{T: t}
	crdClient := FakeV1Client{Mgr: mgr}

	err = v1.EnsureV1CRDs(log, crdClient)
	g.Expect(err).NotTo(HaveOccurred())

	crdList := crdv1.CustomResourceDefinitionList{}
	err = mgr.Client.List(context.TODO(), &crdList)
	g.Expect(err).NotTo(HaveOccurred())

	g.Expect(len(crdList.Items)).To(Equal(1))

	expected := map[string]bool{
		"coherences.coherence.oracle.com": false,
	}

	for _, crd := range crdList.Items {
		expected[crd.Name] = true
	}

	for crd, found := range expected {
		if !found {
			t.Errorf("Failed to create CRD " + crd)
		}
	}
}

func TestShouldUpdateV1CRDs(t *testing.T) {
	var err error

	g := NewGomegaWithT(t)
	mgr, err := fakes.NewFakeManager()
	g.Expect(err).NotTo(HaveOccurred())
	err = crdv1.AddToScheme(mgr.Scheme)
	g.Expect(err).NotTo(HaveOccurred())

	oldCRDs := map[string]*crdv1.CustomResourceDefinition{
		"coherences.coherence.oracle.com": nil,
	}

	for name := range oldCRDs {
		crd := crdv1.CustomResourceDefinition{}
		crd.SetName(name)
		crd.SetResourceVersion("1")
		oldCRDs[name] = &crd
		_ = mgr.GetClient().Create(context.TODO(), &crd)
	}

	log := fakes.TestLogger{T: t}
	crdClient := FakeV1Client{Mgr: mgr}

	err = v1.EnsureV1CRDs(log, crdClient)
	g.Expect(err).NotTo(HaveOccurred())

	crdList := crdv1.CustomResourceDefinitionList{}
	err = mgr.Client.List(context.TODO(), &crdList)
	g.Expect(err).NotTo(HaveOccurred())

	g.Expect(len(crdList.Items)).To(Equal(1))

	for _, crd := range crdList.Items {
		oldCRD := oldCRDs[crd.Name]
		g.Expect(crd).NotTo(Equal(oldCRD))
		g.Expect(crd.GetResourceVersion()).To(Equal(oldCRD.GetResourceVersion()))
	}
}

func TestShouldCreateV1beta1CRDs(t *testing.T) {
	var err error

	g := NewGomegaWithT(t)
	mgr, err := fakes.NewFakeManager()
	g.Expect(err).NotTo(HaveOccurred())

	err = v1beta1.AddToScheme(mgr.Scheme)
	g.Expect(err).NotTo(HaveOccurred())

	log := fakes.TestLogger{T: t}
	crdClient := FakeV1beta1Client{Mgr: mgr}

	err = v1.EnsureV1Beta1CRDs(log, crdClient)
	g.Expect(err).NotTo(HaveOccurred())

	crdList := v1beta1.CustomResourceDefinitionList{}
	err = mgr.Client.List(context.TODO(), &crdList)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(len(crdList.Items)).To(Equal(1))

	expected := map[string]bool{
		"coherences.coherence.oracle.com": false,
	}

	for _, crd := range crdList.Items {
		expected[crd.Name] = true
	}

	for crd, found := range expected {
		if !found {
			t.Errorf("Failed to create CRD " + crd)
		}
	}
}

func TestShouldUpdateV1beta1CRDs(t *testing.T) {
	var err error

	g := NewGomegaWithT(t)
	mgr, err := fakes.NewFakeManager()
	g.Expect(err).NotTo(HaveOccurred())
	err = v1beta1.AddToScheme(mgr.Scheme)
	g.Expect(err).NotTo(HaveOccurred())

	oldCRDs := map[string]*v1beta1.CustomResourceDefinition{
		"coherences.coherence.oracle.com": nil,
	}

	for name := range oldCRDs {
		crd := v1beta1.CustomResourceDefinition{}
		crd.SetName(name)
		crd.SetResourceVersion("1")
		oldCRDs[name] = &crd
		_ = mgr.GetClient().Create(context.TODO(), &crd)
	}

	log := fakes.TestLogger{T: t}
	crdClient := FakeV1beta1Client{Mgr: mgr}

	err = v1.EnsureV1Beta1CRDs(log, crdClient)
	g.Expect(err).NotTo(HaveOccurred())

	crdList := v1beta1.CustomResourceDefinitionList{}
	err = mgr.Client.List(context.TODO(), &crdList)
	g.Expect(err).NotTo(HaveOccurred())

	g.Expect(len(crdList.Items)).To(Equal(1))

	for _, crd := range crdList.Items {
		oldCRD := oldCRDs[crd.Name]
		g.Expect(crd).NotTo(Equal(oldCRD))
		g.Expect(crd.GetResourceVersion()).To(Equal(oldCRD.GetResourceVersion()))
	}
}

// ----- FakeV1Client --------------------------------------------------------------------------------------------------

var _ v1client.CustomResourceDefinitionInterface = FakeV1Client{}

type FakeV1Client struct {
	Mgr manager.Manager
}

func (f FakeV1Client) Get(ctx context.Context, name string, options metav1.GetOptions) (*crdv1.CustomResourceDefinition, error) {
	crd := &crdv1.CustomResourceDefinition{}
	err := f.Mgr.GetClient().Get(context.TODO(), types.NamespacedName{Name: name}, crd)
	return crd, err
}

func (f FakeV1Client) Create(ctx context.Context, crd *crdv1.CustomResourceDefinition, opts metav1.CreateOptions) (*crdv1.CustomResourceDefinition, error) {
	err := f.Mgr.GetClient().Create(context.TODO(), crd)
	if err == nil {
		return f.Get(ctx, crd.Name, metav1.GetOptions{})
	}
	return nil, err
}

func (f FakeV1Client) Update(ctx context.Context, crd *crdv1.CustomResourceDefinition, opts metav1.UpdateOptions) (*crdv1.CustomResourceDefinition, error) {
	err := f.Mgr.GetClient().Update(context.TODO(), crd)
	if err == nil {
		return f.Get(ctx, crd.Name, metav1.GetOptions{})
	}
	return nil, err
}

func (f FakeV1Client) UpdateStatus(ctx context.Context, customResourceDefinition *crdv1.CustomResourceDefinition, opts metav1.UpdateOptions) (*crdv1.CustomResourceDefinition, error) {
	panic("implement me")
}

func (f FakeV1Client) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	panic("implement me")
}

func (f FakeV1Client) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	panic("implement me")
}

func (f FakeV1Client) List(ctx context.Context, opts metav1.ListOptions) (*crdv1.CustomResourceDefinitionList, error) {
	panic("implement me")
}

func (f FakeV1Client) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	panic("implement me")
}

func (f FakeV1Client) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *crdv1.CustomResourceDefinition, err error) {
	panic("implement me")
}

// ----- FakeV1beta1Client ---------------------------------------------------------------------------------------------

var _ v1beta1client.CustomResourceDefinitionInterface = FakeV1beta1Client{}

type FakeV1beta1Client struct {
	Mgr manager.Manager
}

func (f FakeV1beta1Client) Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1beta1.CustomResourceDefinition, error) {
	crd := &v1beta1.CustomResourceDefinition{}
	err := f.Mgr.GetClient().Get(context.TODO(), types.NamespacedName{Name: name}, crd)
	return crd, err
}

func (f FakeV1beta1Client) Create(ctx context.Context, crd *v1beta1.CustomResourceDefinition, opts metav1.CreateOptions) (*v1beta1.CustomResourceDefinition, error) {
	err := f.Mgr.GetClient().Create(context.TODO(), crd)
	if err == nil {
		return f.Get(ctx, crd.Name, metav1.GetOptions{})
	}
	return nil, err
}

func (f FakeV1beta1Client) Update(ctx context.Context, crd *v1beta1.CustomResourceDefinition, opts metav1.UpdateOptions) (*v1beta1.CustomResourceDefinition, error) {
	err := f.Mgr.GetClient().Update(context.TODO(), crd)
	if err == nil {
		return f.Get(ctx, crd.Name, metav1.GetOptions{})
	}
	return nil, err
}

func (f FakeV1beta1Client) UpdateStatus(ctx context.Context, customResourceDefinition *v1beta1.CustomResourceDefinition, opts metav1.UpdateOptions) (*v1beta1.CustomResourceDefinition, error) {
	panic("implement me")
}

func (f FakeV1beta1Client) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	panic("implement me")
}

func (f FakeV1beta1Client) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	panic("implement me")
}

func (f FakeV1beta1Client) List(ctx context.Context, opts metav1.ListOptions) (*v1beta1.CustomResourceDefinitionList, error) {
	panic("implement me")
}

func (f FakeV1beta1Client) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	panic("implement me")
}

func (f FakeV1beta1Client) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1beta1.CustomResourceDefinition, err error) {
	panic("implement me")
}
