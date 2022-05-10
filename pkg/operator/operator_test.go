/*
 * Copyright (c) 2019, 2022 Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package operator_test

import (
	"context"
	"github.com/go-logr/logr"
	. "github.com/onsi/gomega"
	v1 "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/pkg/fakes"
	crdv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"testing"
)

func TestShouldCreateV1CRDs(t *testing.T) {
	var err error

	g := NewGomegaWithT(t)
	mgr, err := fakes.NewFakeManager()
	g.Expect(err).NotTo(HaveOccurred())

	err = crdv1.AddToScheme(mgr.Scheme)
	g.Expect(err).NotTo(HaveOccurred())

	ctx := context.TODO()
	log := logr.New(fakes.TestLogSink{T: t})

	err = v1.EnsureV1CRDs(ctx, log, mgr.Scheme, mgr.Client)
	g.Expect(err).NotTo(HaveOccurred())

	crdList := crdv1.CustomResourceDefinitionList{}
	err = mgr.Client.List(ctx, &crdList)
	g.Expect(err).NotTo(HaveOccurred())

	g.Expect(len(crdList.Items)).To(Equal(1))

	expected := map[string]bool{
		"coherence.coherence.oracle.com": false,
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
		"coherence.coherence.oracle.com": nil,
	}

	for name := range oldCRDs {
		crd := crdv1.CustomResourceDefinition{}
		crd.SetName(name)
		crd.SetResourceVersion("1")
		oldCRDs[name] = &crd
		_ = mgr.GetClient().Create(context.TODO(), &crd)
	}

	ctx := context.TODO()
	log := logr.New(fakes.TestLogSink{T: t})

	err = v1.EnsureV1CRDs(ctx, log, mgr.Scheme, mgr.Client)
	g.Expect(err).NotTo(HaveOccurred())

	crdList := crdv1.CustomResourceDefinitionList{}
	err = mgr.Client.List(ctx, &crdList)
	g.Expect(err).NotTo(HaveOccurred())

	g.Expect(len(crdList.Items)).To(Equal(1))

	for _, crd := range crdList.Items {
		oldCRD := oldCRDs[crd.Name]
		g.Expect(crd).NotTo(Equal(oldCRD))
		g.Expect(crd.GetResourceVersion()).To(Equal(oldCRD.GetResourceVersion()))
	}
}
