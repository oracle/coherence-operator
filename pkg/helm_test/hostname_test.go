/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package helm_test

import (
	. "github.com/onsi/gomega"
	cohv1 "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	"testing"
)

func TestHostnameWhenNotSetInMinimalCluster(t *testing.T) {
	g := NewGomegaWithT(t)

	// Use the specified yaml files to create a CoherenceCluster and trigger a fake end-to-end
	// reconcile to obtain the resources that would have been created by the Helm operator.
	result, cluster, err := CreateCluster("minimal-cluster.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	// Obtain the StatefulSet that Helm would have created
	sts, err := findStatefulSet(result, cluster, cohv1.DefaultRoleName)
	g.Expect(err).NotTo(HaveOccurred())

	g.Expect(sts.Spec.Template.Spec.Hostname).To(Equal(""))
}

func TestHostnameWhenSetForImplicitRole(t *testing.T) {
	g := NewGomegaWithT(t)

	result, cluster, err := CreateCluster("hostname-implicit-role.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	sts, err := findStatefulSet(result, cluster, cohv1.DefaultRoleName)
	g.Expect(err).NotTo(HaveOccurred())

	g.Expect(sts.Spec.Template.Spec.Hostname).To(Equal("foo.com"))
}

func TestHostnameWhenNotSetInClusterWithExplicitRoles(t *testing.T) {
	g := NewGomegaWithT(t)

	result, cluster, err := CreateCluster("minimal-explicit-roles.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	for name := range cluster.GetRoles() {
		sts, err := findStatefulSet(result, cluster, name)
		g.Expect(err).NotTo(HaveOccurred())
		g.Expect(sts.Spec.Template.Spec.Hostname).To(Equal(""))
	}
}

func TestHostnameWhenDefaultSetInClusterWithExplicitRoles(t *testing.T) {
	g := NewGomegaWithT(t)

	result, cluster, err := CreateCluster("hostname-explicit-roles-with-default.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	stsData, err := findStatefulSet(result, cluster, "data")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(stsData.Spec.Template.Spec.Hostname).To(Equal("foo.com"))

	stsProxy, err := findStatefulSet(result, cluster, "proxy")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(stsProxy.Spec.Template.Spec.Hostname).To(Equal("foo.com"))
}

func TestHostnameWhenDefaultSetAndOverriddenInClusterWithExplicitRoles(t *testing.T) {
	g := NewGomegaWithT(t)

	result, cluster, err := CreateCluster("hostname-explicit-roles-override-default.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	stsData, err := findStatefulSet(result, cluster, "data")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(stsData.Spec.Template.Spec.Hostname).To(Equal("foo.com"))

	stsProxy, err := findStatefulSet(result, cluster, "proxy")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(stsProxy.Spec.Template.Spec.Hostname).To(Equal("bar.com"))
}

func TestHostnameWhenSetInClusterWithExplicitRoles(t *testing.T) {
	g := NewGomegaWithT(t)

	result, cluster, err := CreateCluster("hostname-explicit-roles.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	stsData, err := findStatefulSet(result, cluster, "data")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(stsData.Spec.Template.Spec.Hostname).To(Equal("foo.com"))

	stsProxy, err := findStatefulSet(result, cluster, "proxy")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(stsProxy.Spec.Template.Spec.Hostname).To(Equal("bar.com"))
}
