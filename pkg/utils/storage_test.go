/*
 * Copyright (c) 2020, 2025, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package utils

import (
	"testing"

	. "github.com/onsi/gomega"
	coh "github.com/oracle/coherence-operator/api/v1"
)

func TestGetDeletionsWhenLatestAndPrevAreNil(t *testing.T) {
	g := NewGomegaWithT(t)
	s := &secretStore{}
	g.Expect(s.GetDeletions()).To(BeNil())
}

func TestGetDeletionsWhenLatestAndPrevItemsAreNil(t *testing.T) {
	g := NewGomegaWithT(t)
	s := &secretStore{
		latest:   coh.Resources{},
		previous: coh.Resources{},
	}
	g.Expect(s.GetDeletions()).To(BeNil())
}

func TestGetDeletionsWhenLatestAndPrevItemsAreEmpty(t *testing.T) {
	g := NewGomegaWithT(t)
	s := &secretStore{
		latest: coh.Resources{
			Items: make([]coh.Resource, 0),
		},
		previous: coh.Resources{
			Items: make([]coh.Resource, 0),
		},
	}
	g.Expect(s.GetDeletions()).To(BeNil())
}

func TestGetDeletionsWhenPrevItemsAreEmpty(t *testing.T) {
	g := NewGomegaWithT(t)
	s := &secretStore{
		latest: coh.Resources{
			Items: []coh.Resource{
				{Name: "foo", Kind: coh.ResourceTypeService},
				{Name: "svc1", Kind: coh.ResourceTypeService},
				{Name: "foo", Kind: coh.ResourceTypeSecret},
				{Name: "bar", Kind: coh.ResourceTypeSecret},
				{Name: "sec2", Kind: coh.ResourceTypeSecret},
			},
		},
		previous: coh.Resources{
			Items: make([]coh.Resource, 0),
		},
	}
	g.Expect(s.GetDeletions()).To(BeNil())
}

func TestGetDeletions(t *testing.T) {
	g := NewGomegaWithT(t)

	s := &secretStore{
		latest: coh.Resources{
			Items: []coh.Resource{
				{Name: "svc1", Kind: coh.ResourceTypeService},
				{Name: "foo", Kind: coh.ResourceTypeSecret},
				{Name: "bar", Kind: coh.ResourceTypeSecret},
				{Name: "sec2", Kind: coh.ResourceTypeSecret},
			},
		},
		previous: coh.Resources{
			Items: []coh.Resource{
				{Name: "foo", Kind: coh.ResourceTypeService},
				{Name: "bar", Kind: coh.ResourceTypeService},
				{Name: "svc1", Kind: coh.ResourceTypeService},
				{Name: "svc2", Kind: coh.ResourceTypeService},
				{Name: "foo", Kind: coh.ResourceTypeSecret},
				{Name: "sec1", Kind: coh.ResourceTypeSecret},
				{Name: "sec2", Kind: coh.ResourceTypeSecret},
			},
		},
	}
	g.Expect(s.GetDeletions()).NotTo(BeNil())

	expected := []coh.Resource{
		{Name: "foo", Kind: coh.ResourceTypeService},
		{Name: "bar", Kind: coh.ResourceTypeService},
		{Name: "svc2", Kind: coh.ResourceTypeService},
		{Name: "sec1", Kind: coh.ResourceTypeSecret},
	}

	g.Expect(s.GetDeletions()).To(Equal(expected))
}
