/*
 * Copyright (c) 2019, 2020 Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package operator_test

import (
	"context"
	. "github.com/onsi/gomega"
	"github.com/oracle/coherence-operator/pkg/fakes"
	"github.com/oracle/coherence-operator/pkg/flags"
	"github.com/oracle/coherence-operator/pkg/operator"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	"github.com/spf13/pflag"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"testing"
)

func TestShouldCreateCRDs(t *testing.T) {
	var err error

	g := NewGomegaWithT(t)
	mgr, err := fakes.NewFakeManager()
	g.Expect(err).NotTo(HaveOccurred())

	err = v1beta1.AddToScheme(mgr.Scheme)
	g.Expect(err).NotTo(HaveOccurred())

	log := fakes.TestLogger{T: t}
	cohFlags := flags.CoherenceOperatorFlags{}

	crdDir, err := helper.FindCrdDir()
	g.Expect(err).NotTo(HaveOccurred())

	args := []string{"--crd-files", crdDir}
	flagSet := pflag.FlagSet{}
	cohFlags = flags.CoherenceOperatorFlags{}
	cohFlags.AddTo(&flagSet)
	err = flagSet.Parse(args)
	g.Expect(err).NotTo(HaveOccurred())

	crdClient := FakeCustomResourceDefinitionInterface{Mgr: mgr}

	err = operator.EnsureCRDsUsingClient(mgr, &cohFlags, log, crdClient)
	g.Expect(err).NotTo(HaveOccurred())

	crdList := v1beta1.CustomResourceDefinitionList{}
	err = mgr.Client.List(context.TODO(), &crdList)
	g.Expect(err).NotTo(HaveOccurred())

	g.Expect(len(crdList.Items)).To(Equal(3))

	expected := map[string]bool{
		"coherenceclusters.coherence.oracle.com":  false,
		"coherenceinternals.coherence.oracle.com": false,
		"coherenceroles.coherence.oracle.com":     false,
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

func TestShouldUpdateCRDs(t *testing.T) {
	var err error

	g := NewGomegaWithT(t)
	mgr, err := fakes.NewFakeManager()
	g.Expect(err).NotTo(HaveOccurred())
	err = v1beta1.AddToScheme(mgr.Scheme)
	g.Expect(err).NotTo(HaveOccurred())

	oldCRDs := map[string]*v1beta1.CustomResourceDefinition{
		"coherenceclusters.coherence.oracle.com":  nil,
		"coherenceinternals.coherence.oracle.com": nil,
		"coherenceroles.coherence.oracle.com":     nil,
	}

	for name := range oldCRDs {
		crd := v1beta1.CustomResourceDefinition{}
		crd.SetName(name)
		crd.SetResourceVersion(name + "-1234")
		oldCRDs[name] = &crd
		_ = mgr.GetClient().Create(context.TODO(), &crd)
	}

	log := fakes.TestLogger{T: t}
	cohFlags := flags.CoherenceOperatorFlags{}

	crdDir, err := helper.FindCrdDir()
	g.Expect(err).NotTo(HaveOccurred())

	args := []string{"--crd-files", crdDir}
	flagSet := pflag.FlagSet{}
	cohFlags = flags.CoherenceOperatorFlags{}
	cohFlags.AddTo(&flagSet)
	err = flagSet.Parse(args)
	g.Expect(err).NotTo(HaveOccurred())

	crdClient := FakeCustomResourceDefinitionInterface{Mgr: mgr}

	err = operator.EnsureCRDsUsingClient(mgr, &cohFlags, log, crdClient)
	g.Expect(err).NotTo(HaveOccurred())

	crdList := v1beta1.CustomResourceDefinitionList{}
	err = mgr.Client.List(context.TODO(), &crdList)
	g.Expect(err).NotTo(HaveOccurred())

	g.Expect(len(crdList.Items)).To(Equal(3))

	for _, crd := range crdList.Items {
		oldCRD := oldCRDs[crd.Name]
		g.Expect(crd).NotTo(Equal(oldCRD))
		g.Expect(crd.GetResourceVersion()).To(Equal(oldCRD.GetResourceVersion()))
	}
}

type FakeCustomResourceDefinitionInterface struct {
	Mgr manager.Manager
}

func (f FakeCustomResourceDefinitionInterface) Get(name string, options metav1.GetOptions) (*v1beta1.CustomResourceDefinition, error) {
	crd := &v1beta1.CustomResourceDefinition{}
	err := f.Mgr.GetClient().Get(context.TODO(), types.NamespacedName{Name: name}, crd)
	return crd, err
}

func (f FakeCustomResourceDefinitionInterface) Create(crd *v1beta1.CustomResourceDefinition) (*v1beta1.CustomResourceDefinition, error) {
	err := f.Mgr.GetClient().Create(context.TODO(), crd)
	if err == nil {
		return f.Get(crd.Name, metav1.GetOptions{})
	}
	return nil, err
}

func (f FakeCustomResourceDefinitionInterface) Update(crd *v1beta1.CustomResourceDefinition) (*v1beta1.CustomResourceDefinition, error) {
	err := f.Mgr.GetClient().Update(context.TODO(), crd)
	if err == nil {
		return f.Get(crd.Name, metav1.GetOptions{})
	}
	return nil, err
}

func (f FakeCustomResourceDefinitionInterface) UpdateStatus(*v1beta1.CustomResourceDefinition) (*v1beta1.CustomResourceDefinition, error) {
	panic("implement me")
}

func (f FakeCustomResourceDefinitionInterface) Delete(name string, options *metav1.DeleteOptions) error {
	panic("implement me")
}

func (f FakeCustomResourceDefinitionInterface) DeleteCollection(options *metav1.DeleteOptions, listOptions metav1.ListOptions) error {
	panic("implement me")
}

func (f FakeCustomResourceDefinitionInterface) List(opts metav1.ListOptions) (*v1beta1.CustomResourceDefinitionList, error) {
	panic("implement me")
}

func (f FakeCustomResourceDefinitionInterface) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	panic("implement me")
}

func (f FakeCustomResourceDefinitionInterface) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1beta1.CustomResourceDefinition, err error) {
	panic("implement me")
}
