/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package helm_test

import (
	"fmt"
	. "github.com/onsi/gomega"
	cohv1 "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	corev1 "k8s.io/api/core/v1"
	"testing"
)

func TestHostAliasesWhenNotSetInMinimalCluster(t *testing.T) {
	g := NewGomegaWithT(t)

	// Use the specified yaml files to create a CoherenceCluster and trigger a fake end-to-end
	// reconcile to obtain the resources that would have been created by the Helm operator.
	result, cluster, err := CreateCluster("minimal-cluster.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	// Obtain the StatefulSet that Helm would have created
	sts, err := findStatefulSet(result, cluster, cohv1.DefaultRoleName)
	g.Expect(err).NotTo(HaveOccurred())

	g.Expect(sts.Spec.Template.Spec.HostAliases).To(BeNil())
}

func TestHostAliasesWhenSetForImplicitRole(t *testing.T) {
	g := NewGomegaWithT(t)

	result, cluster, err := CreateCluster("host-aliases-implicit-role.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	sts, err := findStatefulSet(result, cluster, cohv1.DefaultRoleName)
	g.Expect(err).NotTo(HaveOccurred())

	assertHostAliases(t, sts.Spec.Template.Spec, []corev1.HostAlias{
		{IP: "10.10.10.100", Hostnames: []string{"foo1.com", "foo2.com"}},
		{IP: "10.10.10.200", Hostnames: []string{"bar1.com", "bar2.com"}},
	})
}

func TestHostAliasesWhenNotSetInClusterWithExplicitRoles(t *testing.T) {
	g := NewGomegaWithT(t)

	result, cluster, err := CreateCluster("minimal-explicit-roles.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	for name := range cluster.GetRoles() {
		sts, err := findStatefulSet(result, cluster, name)
		g.Expect(err).NotTo(HaveOccurred())
		g.Expect(sts.Spec.Template.Spec.HostAliases).To(BeNil())
	}
}

func TestHostAliasesWhenDefaultSetInClusterWithExplicitRoles(t *testing.T) {
	g := NewGomegaWithT(t)

	result, cluster, err := CreateCluster("host-aliases-explicit-roles-with-default.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	stsData, err := findStatefulSet(result, cluster, "data")
	g.Expect(err).NotTo(HaveOccurred())

	assertHostAliases(t, stsData.Spec.Template.Spec, []corev1.HostAlias{
		{IP: "10.10.10.100", Hostnames: []string{"foo1.com", "foo2.com"}},
		{IP: "10.10.10.200", Hostnames: []string{"bar1.com", "bar2.com"}},
	})

	stsProxy, err := findStatefulSet(result, cluster, "proxy")
	g.Expect(err).NotTo(HaveOccurred())

	assertHostAliases(t, stsProxy.Spec.Template.Spec, []corev1.HostAlias{
		{IP: "10.10.10.100", Hostnames: []string{"foo1.com", "foo2.com"}},
		{IP: "10.10.10.200", Hostnames: []string{"bar1.com", "bar2.com"}},
	})
}

func TestHostAliasesWhenDefaultSetAndOverriddenInClusterWithExplicitRoles(t *testing.T) {
	g := NewGomegaWithT(t)

	result, cluster, err := CreateCluster("host-aliases-explicit-roles-override-default.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	stsData, err := findStatefulSet(result, cluster, "data")
	g.Expect(err).NotTo(HaveOccurred())

	assertHostAliases(t, stsData.Spec.Template.Spec, []corev1.HostAlias{
		{IP: "10.10.10.100", Hostnames: []string{"foo1.com", "foo2.com"}},
		{IP: "10.10.10.200", Hostnames: []string{"bar1.com", "bar2.com"}},
	})

	stsProxy, err := findStatefulSet(result, cluster, "proxy")
	g.Expect(err).NotTo(HaveOccurred())

	assertHostAliases(t, stsProxy.Spec.Template.Spec, []corev1.HostAlias{
		{IP: "10.10.10.100", Hostnames: []string{"foo1.com", "foo2.com"}},
		{IP: "10.10.10.200", Hostnames: []string{"bar1.com", "bar2.com"}},
		{IP: "10.10.10.300", Hostnames: []string{"foo3.com", "foo4.com"}},
		{IP: "10.10.10.400", Hostnames: []string{"bar3.com", "bar4.com"}},
	})

	stsWeb, err := findStatefulSet(result, cluster, "web")
	g.Expect(err).NotTo(HaveOccurred())

	assertHostAliases(t, stsWeb.Spec.Template.Spec, []corev1.HostAlias{
		{IP: "10.10.10.100", Hostnames: []string{"foo11.com"}},
		{IP: "10.10.10.200", Hostnames: []string{"bar1.com", "bar2.com"}},
		{IP: "10.10.10.400", Hostnames: []string{"bar3.com", "bar4.com"}},
	})
}

func TestHostAliasesWhenSetInClusterWithExplicitRoles(t *testing.T) {
	g := NewGomegaWithT(t)

	result, cluster, err := CreateCluster("host-aliases-explicit-roles.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	stsData, err := findStatefulSet(result, cluster, "data")
	g.Expect(err).NotTo(HaveOccurred())

	assertHostAliases(t, stsData.Spec.Template.Spec, []corev1.HostAlias{
		{IP: "10.10.10.100", Hostnames: []string{"foo1.com", "foo2.com"}},
		{IP: "10.10.10.200", Hostnames: []string{"bar1.com", "bar2.com"}},
	})

	stsProxy, err := findStatefulSet(result, cluster, "proxy")
	g.Expect(err).NotTo(HaveOccurred())

	assertHostAliases(t, stsProxy.Spec.Template.Spec, []corev1.HostAlias{
		{IP: "10.10.10.300", Hostnames: []string{"foo3.com", "foo4.com"}},
		{IP: "10.10.10.400", Hostnames: []string{"bar3.com", "bar4.com"}},
	})
}

func assertHostAliases(t *testing.T, pod corev1.PodSpec, expected []corev1.HostAlias) {
	g := NewGomegaWithT(t)

	g.Expect(len(pod.HostAliases)).To(Equal(len(expected)))

	m := make(map[string]corev1.HostAlias)
	for _, ha := range pod.HostAliases {
		m[ha.IP] = ha
	}

	for _, ex := range expected {
		opt, found := m[ex.IP]
		g.Expect(found).To(BeTrue(), fmt.Sprintf("Did not find IP %s", ex.IP))
		g.Expect(opt.IP).To(Equal(ex.IP))
		g.Expect(opt.Hostnames).To(Equal(ex.Hostnames))
	}
}
