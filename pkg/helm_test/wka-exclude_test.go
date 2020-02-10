/*
 * Copyright (c) 2020, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package helm_test

import (
	. "github.com/onsi/gomega"
	cohv1 "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	"testing"
)

func TestWKAExcludeWhenNotSetInMinimalCluster(t *testing.T) {
	g := NewGomegaWithT(t)

	// Use the specified yaml files to create a CoherenceCluster and trigger a fake end-to-end
	// reconcile to obtain the resources that would have been created by the Helm operator.
	result, cluster, err := CreateCluster("minimal-cluster.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	// Obtain the StatefulSet that Helm would have created
	sts, err := findStatefulSet(result, cluster, cohv1.DefaultRoleName)
	g.Expect(err).NotTo(HaveOccurred())

	// Obtain the coherenceWKAMember label from the StatefulSet Pod template
	isMember := sts.Spec.Template.GetLabels()["coherenceWKAMember"]
	g.Expect(isMember).To(Equal("true"))
}

func TestWKAExcludeWhenSetForImplicitRole(t *testing.T) {
	g := NewGomegaWithT(t)

	result, cluster, err := CreateCluster("wka-exclude-test-implicit-role.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	sts, err := findStatefulSet(result, cluster, cohv1.DefaultRoleName)
	g.Expect(err).NotTo(HaveOccurred())

	// Obtain the coherenceWKAMember label from the StatefulSet Pod template
	isMember := sts.Spec.Template.GetLabels()["coherenceWKAMember"]
	g.Expect(isMember).To(Equal("false"))
}

func TestWKAExcludeWhenNotSetInClusterWithExplicitRoles(t *testing.T) {
	g := NewGomegaWithT(t)

	result, cluster, err := CreateCluster("minimal-explicit-roles.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	for name := range cluster.GetRoles() {
		sts, err := findStatefulSet(result, cluster, name)
		g.Expect(err).NotTo(HaveOccurred())

		// Obtain the coherenceWKAMember label from the StatefulSet Pod template
		isMember := sts.Spec.Template.GetLabels()["coherenceWKAMember"]
		g.Expect(isMember).To(Equal("true"))
	}
}

func TestWKAExcludeWhenDefaultSetInClusterWithExplicitRoles(t *testing.T) {
	g := NewGomegaWithT(t)

	result, cluster, err := CreateCluster("wka-exclude-test-explicit-roles-with-default.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	// Obtain the data role's StatefulSet
	stsData, err := findStatefulSet(result, cluster, "data")
	g.Expect(err).NotTo(HaveOccurred())

	// Obtain the coherenceWKAMember label from the data role's StatefulSet Pod template
	isMemberData := stsData.Spec.Template.GetLabels()["coherenceWKAMember"]
	g.Expect(isMemberData).To(Equal("false"))

	// Obtain the proxy role's StatefulSet
	stsProxy, err := findStatefulSet(result, cluster, "proxy")
	g.Expect(err).NotTo(HaveOccurred())

	// Obtain the coherenceWKAMember label from the proxy role's StatefulSet Pod template
	isMemberProxy := stsProxy.Spec.Template.GetLabels()["coherenceWKAMember"]
	g.Expect(isMemberProxy).To(Equal("false"))
}

func TestWKAExcludeWhenDefaultSetAndOverriddenInClusterWithExplicitRoles(t *testing.T) {
	g := NewGomegaWithT(t)

	result, cluster, err := CreateCluster("wka-exclude-test-explicit-roles-override-default.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	// Obtain the data role's StatefulSet
	stsData, err := findStatefulSet(result, cluster, "data")
	g.Expect(err).NotTo(HaveOccurred())

	// Obtain the coherenceWKAMember label from the data role's StatefulSet Pod template
	isMemberData := stsData.Spec.Template.GetLabels()["coherenceWKAMember"]
	g.Expect(isMemberData).To(Equal("true"))

	// Obtain the proxy role's StatefulSet
	stsProxy, err := findStatefulSet(result, cluster, "proxy")
	g.Expect(err).NotTo(HaveOccurred())

	// Obtain the coherenceWKAMember label from the proxy role's StatefulSet Pod template
	isMemberProxy := stsProxy.Spec.Template.GetLabels()["coherenceWKAMember"]
	g.Expect(isMemberProxy).To(Equal("false"))
}

func TestWKAExcludeWhenSetInClusterWithExplicitRoles(t *testing.T) {
	g := NewGomegaWithT(t)

	result, cluster, err := CreateCluster("wka-exclude-test-explicit-roles.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	// Obtain the data role's StatefulSet
	stsData, err := findStatefulSet(result, cluster, "data")
	g.Expect(err).NotTo(HaveOccurred())

	// Obtain the coherenceWKAMember label from the data role's StatefulSet Pod template
	isMemberData := stsData.Spec.Template.GetLabels()["coherenceWKAMember"]
	g.Expect(isMemberData).To(Equal("true"))

	// Obtain the proxy role's StatefulSet
	stsProxy, err := findStatefulSet(result, cluster, "proxy")
	g.Expect(err).NotTo(HaveOccurred())

	// Obtain the coherenceWKAMember label from the proxy role's StatefulSet Pod template
	isMemberProxy := stsProxy.Spec.Template.GetLabels()["coherenceWKAMember"]
	g.Expect(isMemberProxy).To(Equal("false"))
}
