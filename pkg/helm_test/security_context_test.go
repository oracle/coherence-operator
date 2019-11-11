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

func TestPodSecurityContextWhenNotSetInMinimalCluster(t *testing.T) {
	g := NewGomegaWithT(t)

	// Use the specified yaml files to create a CoherenceCluster and trigger a fake end-to-end
	// reconcile to obtain the resources that would have been created by the Helm operator.
	result, cluster, err := CreateCluster("minimal-cluster.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	// Obtain the StatefulSet that Helm would have created
	sts, err := findStatefulSet(result, cluster, cohv1.DefaultRoleName)
	g.Expect(err).NotTo(HaveOccurred())

	// Assert that the security context is not set
	ctx := sts.Spec.Template.Spec.SecurityContext
	g.Expect(ctx).To(BeNil())
}

func TestPodSecurityContextWhenSetForImplicitRole(t *testing.T) {
	g := NewGomegaWithT(t)

	result, cluster, err := CreateCluster("security-context-implicit-role.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	sts, err := findStatefulSet(result, cluster, cohv1.DefaultRoleName)
	g.Expect(err).NotTo(HaveOccurred())

	// Assert that the security context is set
	ctx := sts.Spec.Template.Spec.SecurityContext
	g.Expect(ctx).NotTo(BeNil())
	g.Expect(ctx.RunAsUser).To(Equal(pointer.Int64Ptr(1001)))
	g.Expect(ctx.RunAsNonRoot).To(Equal(pointer.BoolPtr(true)))
	g.Expect(ctx.FSGroup).To(BeNil())
	g.Expect(ctx.RunAsGroup).To(BeNil())
	g.Expect(ctx.SELinuxOptions).To(BeNil())
	g.Expect(ctx.SupplementalGroups).To(BeNil())
	g.Expect(ctx.Sysctls).To(BeNil())
}

func TestPodSecurityContextWhenNotSetInClusterWithExplicitRoles(t *testing.T) {
	g := NewGomegaWithT(t)

	result, cluster, err := CreateCluster("minimal-explicit-roles.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	for name := range cluster.GetRoles() {
		sts, err := findStatefulSet(result, cluster, name)
		g.Expect(err).NotTo(HaveOccurred())

		// Assert that the security context is not set
		ctx := sts.Spec.Template.Spec.SecurityContext
		g.Expect(ctx).To(BeNil())
	}
}

func TestPodSecurityContextWhenDefaultSetInClusterWithExplicitRoles(t *testing.T) {
	g := NewGomegaWithT(t)

	result, cluster, err := CreateCluster("security-context-explicit-roles-with-default.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	// Assert that the security context is set for the data role
	stsData, err := findStatefulSet(result, cluster, "data")
	g.Expect(err).NotTo(HaveOccurred())
	ctxData := stsData.Spec.Template.Spec.SecurityContext
	g.Expect(ctxData).NotTo(BeNil())
	g.Expect(ctxData.RunAsUser).To(Equal(pointer.Int64Ptr(1001)))
	g.Expect(ctxData.RunAsNonRoot).To(Equal(pointer.BoolPtr(true)))
	g.Expect(ctxData.FSGroup).To(BeNil())
	g.Expect(ctxData.RunAsGroup).To(BeNil())
	g.Expect(ctxData.SELinuxOptions).To(BeNil())
	g.Expect(ctxData.SupplementalGroups).To(BeNil())
	g.Expect(ctxData.Sysctls).To(BeNil())

	// Assert that the security context is set for the data role
	stsProxy, err := findStatefulSet(result, cluster, "proxy")
	g.Expect(err).NotTo(HaveOccurred())
	ctxProxy := stsProxy.Spec.Template.Spec.SecurityContext
	g.Expect(ctxProxy).NotTo(BeNil())
	g.Expect(ctxProxy.RunAsUser).To(Equal(pointer.Int64Ptr(1001)))
	g.Expect(ctxProxy.RunAsNonRoot).To(Equal(pointer.BoolPtr(true)))
	g.Expect(ctxProxy.FSGroup).To(BeNil())
	g.Expect(ctxProxy.RunAsGroup).To(BeNil())
	g.Expect(ctxProxy.SELinuxOptions).To(BeNil())
	g.Expect(ctxProxy.SupplementalGroups).To(BeNil())
	g.Expect(ctxProxy.Sysctls).To(BeNil())
}

func TestPodSecurityContextWhenDefaultSetAndOverriddenInClusterWithExplicitRoles(t *testing.T) {
	g := NewGomegaWithT(t)

	result, cluster, err := CreateCluster("security-context-explicit-roles-override-default.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	// Assert that the security context is set for the data role
	stsData, err := findStatefulSet(result, cluster, "data")
	g.Expect(err).NotTo(HaveOccurred())
	ctxData := stsData.Spec.Template.Spec.SecurityContext
	g.Expect(ctxData).NotTo(BeNil())
	g.Expect(ctxData.RunAsUser).To(Equal(pointer.Int64Ptr(1001)))
	g.Expect(ctxData.RunAsNonRoot).To(Equal(pointer.BoolPtr(true)))
	g.Expect(ctxData.FSGroup).To(BeNil())
	g.Expect(ctxData.RunAsGroup).To(BeNil())
	g.Expect(ctxData.SELinuxOptions).To(BeNil())
	g.Expect(ctxData.SupplementalGroups).To(BeNil())
	g.Expect(ctxData.Sysctls).To(BeNil())

	// Assert that the security context is set for the data role
	stsProxy, err := findStatefulSet(result, cluster, "proxy")
	g.Expect(err).NotTo(HaveOccurred())
	ctxProxy := stsProxy.Spec.Template.Spec.SecurityContext
	g.Expect(ctxProxy).NotTo(BeNil())
	g.Expect(ctxProxy.RunAsUser).To(Equal(pointer.Int64Ptr(2000)))
	g.Expect(ctxProxy.RunAsNonRoot).To(BeNil())
	g.Expect(ctxProxy.FSGroup).To(BeNil())
	g.Expect(ctxProxy.RunAsGroup).To(BeNil())
	g.Expect(ctxProxy.SELinuxOptions).To(BeNil())
	g.Expect(ctxProxy.SupplementalGroups).To(BeNil())
	g.Expect(ctxProxy.Sysctls).To(BeNil())
}

func TestPodSecurityContextWhenSetInClusterWithExplicitRoles(t *testing.T) {
	g := NewGomegaWithT(t)

	result, cluster, err := CreateCluster("security-context-explicit-roles.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	// Assert that the security context is set for the data role
	stsData, err := findStatefulSet(result, cluster, "data")
	g.Expect(err).NotTo(HaveOccurred())
	ctxData := stsData.Spec.Template.Spec.SecurityContext
	g.Expect(ctxData).NotTo(BeNil())
	g.Expect(ctxData.RunAsUser).To(Equal(pointer.Int64Ptr(1001)))
	g.Expect(ctxData.RunAsNonRoot).To(Equal(pointer.BoolPtr(true)))
	g.Expect(ctxData.FSGroup).To(BeNil())
	g.Expect(ctxData.RunAsGroup).To(BeNil())
	g.Expect(ctxData.SELinuxOptions).To(BeNil())
	g.Expect(ctxData.SupplementalGroups).To(BeNil())
	g.Expect(ctxData.Sysctls).To(BeNil())

	// Assert that the security context is set for the data role
	stsProxy, err := findStatefulSet(result, cluster, "proxy")
	g.Expect(err).NotTo(HaveOccurred())
	ctxProxy := stsProxy.Spec.Template.Spec.SecurityContext
	g.Expect(ctxProxy).NotTo(BeNil())
	g.Expect(ctxProxy.RunAsUser).To(Equal(pointer.Int64Ptr(2002)))
	g.Expect(ctxProxy.RunAsNonRoot).To(BeNil())
	g.Expect(ctxProxy.FSGroup).To(BeNil())
	g.Expect(ctxProxy.RunAsGroup).To(BeNil())
	g.Expect(ctxProxy.SELinuxOptions).To(BeNil())
	g.Expect(ctxProxy.SupplementalGroups).To(BeNil())
	g.Expect(ctxProxy.Sysctls).To(BeNil())
}
