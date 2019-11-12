/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package helm_test

import (
	. "github.com/onsi/gomega"
	cohv1 "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	"k8s.io/utils/pointer"
	"testing"
)

/*
 * These tests verify the various scenarios for setting Pod Security Policy
 * in a CoherenceCluster.
 */

func TestShareProcessNamespaceWhenNotSetInMinimalCluster(t *testing.T) {
	g := NewGomegaWithT(t)

	// Use the specified yaml files to create a CoherenceCluster and trigger a fake end-to-end
	// reconcile to obtain the resources that would have been created by the Helm operator.
	result, cluster, err := CreateCluster("minimal-cluster.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	// Obtain the StatefulSet that Helm would have created
	sts, err := findStatefulSet(result, cluster, cohv1.DefaultRoleName)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(sts.Spec.Template.Spec.ShareProcessNamespace).To(BeNil())
}

func TestShareProcessNamespaceWhenSetForImplicitRole(t *testing.T) {
	g := NewGomegaWithT(t)

	result, cluster, err := CreateCluster("share-process-namespace-implicit-role.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	sts, err := findStatefulSet(result, cluster, cohv1.DefaultRoleName)
	g.Expect(err).NotTo(HaveOccurred())

	g.Expect(sts.Spec.Template.Spec.ShareProcessNamespace).To(Equal(pointer.BoolPtr(true)))
}

func TestShareProcessNamespaceWhenNotSetInClusterWithExplicitRoles(t *testing.T) {
	g := NewGomegaWithT(t)

	result, cluster, err := CreateCluster("minimal-explicit-roles.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	for name := range cluster.GetRoles() {
		sts, err := findStatefulSet(result, cluster, name)
		g.Expect(err).NotTo(HaveOccurred())
		g.Expect(sts.Spec.Template.Spec.ShareProcessNamespace).To(BeNil())
	}
}

func TestShareProcessNamespaceWhenDefaultSetInClusterWithExplicitRoles(t *testing.T) {
	g := NewGomegaWithT(t)

	result, cluster, err := CreateCluster("share-process-namespace-explicit-roles-with-default.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	stsData, err := findStatefulSet(result, cluster, "data")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(stsData.Spec.Template.Spec.ShareProcessNamespace).To(Equal(pointer.BoolPtr(true)))

	stsProxy, err := findStatefulSet(result, cluster, "proxy")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(stsProxy.Spec.Template.Spec.ShareProcessNamespace).To(Equal(pointer.BoolPtr(true)))
}

func TestShareProcessNamespaceWhenDefaultSetAndOverriddenInClusterWithExplicitRoles(t *testing.T) {
	g := NewGomegaWithT(t)

	result, cluster, err := CreateCluster("share-process-namespace-explicit-roles-override-default.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	stsData, err := findStatefulSet(result, cluster, "data")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(stsData.Spec.Template.Spec.ShareProcessNamespace).To(Equal(pointer.BoolPtr(true)))

	stsProxy, err := findStatefulSet(result, cluster, "proxy")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(stsProxy.Spec.Template.Spec.ShareProcessNamespace).To(Equal(pointer.BoolPtr(false)))
}

func TestShareProcessNamespaceWhenSetInClusterWithExplicitRoles(t *testing.T) {
	g := NewGomegaWithT(t)

	result, cluster, err := CreateCluster("share-process-namespace-explicit-roles.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	stsData, err := findStatefulSet(result, cluster, "data")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(stsData.Spec.Template.Spec.ShareProcessNamespace).To(Equal(pointer.BoolPtr(true)))

	stsProxy, err := findStatefulSet(result, cluster, "proxy")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(stsProxy.Spec.Template.Spec.ShareProcessNamespace).To(Equal(pointer.BoolPtr(false)))
}
