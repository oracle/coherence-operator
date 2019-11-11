/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package helm_test

import (
	. "github.com/onsi/gomega"
	cohv1 "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/pointer"
	"testing"
)

/*
 * These tests verify the various scenarios for setting Pod Security Policy
 * in a CoherenceCluster.
 */

func TestDnsConfigWhenNotSetInMinimalCluster(t *testing.T) {
	g := NewGomegaWithT(t)

	// Use the specified yaml files to create a CoherenceCluster and trigger a fake end-to-end
	// reconcile to obtain the resources that would have been created by the Helm operator.
	result, cluster, err := CreateCluster("minimal-cluster.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	// Obtain the StatefulSet that Helm would have created
	sts, err := findStatefulSet(result, cluster, cohv1.DefaultRoleName)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(sts.Spec.Template.Spec.DNSConfig).To(BeNil())
}

func TestDnsConfigWhenSetForImplicitRole(t *testing.T) {
	g := NewGomegaWithT(t)

	result, cluster, err := CreateCluster("dns-config-implicit-role.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	sts, err := findStatefulSet(result, cluster, cohv1.DefaultRoleName)
	g.Expect(err).NotTo(HaveOccurred())

	g.Expect(sts.Spec.Template.Spec.DNSConfig).NotTo(BeNil())
	g.Expect(sts.Spec.Template.Spec.DNSConfig.Nameservers).To(Equal([]string{"ns-one", "ns-two"}))
	g.Expect(sts.Spec.Template.Spec.DNSConfig.Searches).To(Equal([]string{"s-one", "s-two"}))
	g.Expect(sts.Spec.Template.Spec.DNSConfig.Options).
		To(Equal([]corev1.PodDNSConfigOption{
			{Name: "o1", Value: pointer.StringPtr("v1")},
			{Name: "o2", Value: pointer.StringPtr("v2")}}))
}

func TestDnsConfigWhenNotSetInClusterWithExplicitRoles(t *testing.T) {
	g := NewGomegaWithT(t)

	result, cluster, err := CreateCluster("minimal-explicit-roles.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	for name := range cluster.GetRoles() {
		sts, err := findStatefulSet(result, cluster, name)
		g.Expect(err).NotTo(HaveOccurred())
		g.Expect(sts.Spec.Template.Spec.DNSConfig).To(BeNil())
	}
}

func TestDnsConfigWhenDefaultSetInClusterWithExplicitRoles(t *testing.T) {
	g := NewGomegaWithT(t)

	result, cluster, err := CreateCluster("dns-config-explicit-roles-with-default.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	stsData, err := findStatefulSet(result, cluster, "data")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(stsData.Spec.Template.Spec.DNSConfig).NotTo(BeNil())
	g.Expect(stsData.Spec.Template.Spec.DNSConfig.Nameservers).To(Equal([]string{"ns-one", "ns-two"}))
	g.Expect(stsData.Spec.Template.Spec.DNSConfig.Searches).To(Equal([]string{"s-one", "s-two"}))
	g.Expect(stsData.Spec.Template.Spec.DNSConfig.Options).
		To(Equal([]corev1.PodDNSConfigOption{
			{Name: "o1", Value: pointer.StringPtr("v1")},
			{Name: "o2", Value: pointer.StringPtr("v2")}}))

	stsProxy, err := findStatefulSet(result, cluster, "proxy")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(stsProxy.Spec.Template.Spec.DNSConfig).NotTo(BeNil())
	g.Expect(stsProxy.Spec.Template.Spec.DNSConfig.Nameservers).To(Equal([]string{"ns-one", "ns-two"}))
	g.Expect(stsProxy.Spec.Template.Spec.DNSConfig.Searches).To(Equal([]string{"s-one", "s-two"}))
	g.Expect(stsProxy.Spec.Template.Spec.DNSConfig.Options).
		To(Equal([]corev1.PodDNSConfigOption{
			{Name: "o1", Value: pointer.StringPtr("v1")},
			{Name: "o2", Value: pointer.StringPtr("v2")}}))
}

func TestDnsConfigWhenDefaultSetAndOverriddenInClusterWithExplicitRoles(t *testing.T) {
	g := NewGomegaWithT(t)

	result, cluster, err := CreateCluster("dns-config-explicit-roles-override-default.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	stsData, err := findStatefulSet(result, cluster, "data")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(stsData.Spec.Template.Spec.DNSConfig).NotTo(BeNil())
	g.Expect(stsData.Spec.Template.Spec.DNSConfig.Nameservers).To(Equal([]string{"ns-one", "ns-two"}))
	g.Expect(stsData.Spec.Template.Spec.DNSConfig.Searches).To(Equal([]string{"s-one", "s-two"}))
	assertPodDNSConfigOption(t, stsData.Spec.Template.Spec.DNSConfig, []corev1.PodDNSConfigOption{
		{Name: "o1", Value: pointer.StringPtr("v1")},
		{Name: "o2", Value: pointer.StringPtr("v2")}})

	stsProxy, err := findStatefulSet(result, cluster, "proxy")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(stsProxy.Spec.Template.Spec.DNSConfig).NotTo(BeNil())
	g.Expect(stsProxy.Spec.Template.Spec.DNSConfig.Nameservers).To(Equal([]string{"ns-one", "ns-two", "ns-three", "ns-four"}))
	g.Expect(stsProxy.Spec.Template.Spec.DNSConfig.Searches).To(Equal([]string{"s-one", "s-two", "s-three", "s-four"}))
	assertPodDNSConfigOption(t, stsProxy.Spec.Template.Spec.DNSConfig, []corev1.PodDNSConfigOption{
		{Name: "o1", Value: pointer.StringPtr("v1")},
		{Name: "o2", Value: pointer.StringPtr("v2")},
		{Name: "o3", Value: pointer.StringPtr("v3")},
		{Name: "o4", Value: pointer.StringPtr("v4")}})

	stsWeb, err := findStatefulSet(result, cluster, "web")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(stsWeb.Spec.Template.Spec.DNSConfig).NotTo(BeNil())
	g.Expect(stsWeb.Spec.Template.Spec.DNSConfig.Nameservers).To(Equal([]string{"ns-one", "ns-two", "ns-three", "ns-four"}))
	g.Expect(stsWeb.Spec.Template.Spec.DNSConfig.Searches).To(Equal([]string{"s-one", "s-two", "s-three", "s-four"}))
	assertPodDNSConfigOption(t, stsWeb.Spec.Template.Spec.DNSConfig, []corev1.PodDNSConfigOption{
		{Name: "o1", Value: pointer.StringPtr("v11")},
		{Name: "o4", Value: pointer.StringPtr("v4")}})
}

func assertPodDNSConfigOption(t *testing.T, cfg *corev1.PodDNSConfig, expected []corev1.PodDNSConfigOption) {
	g := NewGomegaWithT(t)
	m := make(map[string]corev1.PodDNSConfigOption)
	for _, opt := range cfg.Options {
		m[opt.Name] = opt
	}

	for _, ex := range expected {
		opt, found := m[ex.Name]
		g.Expect(found).To(BeTrue())
		g.Expect(opt.Name).To(Equal(ex.Name))
		g.Expect(*opt.Value).To(Equal(*ex.Value))
	}
}

func TestDnsConfigWhenSetInClusterWithExplicitRoles(t *testing.T) {
	g := NewGomegaWithT(t)

	result, cluster, err := CreateCluster("dns-config-explicit-roles.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	stsData, err := findStatefulSet(result, cluster, "data")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(stsData.Spec.Template.Spec.DNSConfig).NotTo(BeNil())
	g.Expect(stsData.Spec.Template.Spec.DNSConfig.Nameservers).To(Equal([]string{"ns-one", "ns-two"}))
	g.Expect(stsData.Spec.Template.Spec.DNSConfig.Searches).To(Equal([]string{"s-one", "s-two"}))
	g.Expect(stsData.Spec.Template.Spec.DNSConfig.Options).
		To(Equal([]corev1.PodDNSConfigOption{
			{Name: "o1", Value: pointer.StringPtr("v1")},
			{Name: "o2", Value: pointer.StringPtr("v2")}}))

	stsProxy, err := findStatefulSet(result, cluster, "proxy")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(stsProxy.Spec.Template.Spec.DNSConfig).NotTo(BeNil())
	g.Expect(stsProxy.Spec.Template.Spec.DNSConfig.Nameservers).To(Equal([]string{"ns-three", "ns-four"}))
	g.Expect(stsProxy.Spec.Template.Spec.DNSConfig.Searches).To(Equal([]string{"s-three", "s-four"}))
	g.Expect(stsProxy.Spec.Template.Spec.DNSConfig.Options).
		To(Equal([]corev1.PodDNSConfigOption{
			{Name: "o3", Value: pointer.StringPtr("v3")},
			{Name: "o4", Value: pointer.StringPtr("v4")}}))
}
