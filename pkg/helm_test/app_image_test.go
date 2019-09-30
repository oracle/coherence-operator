/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package helm_test

import (
	. "github.com/onsi/gomega"
	cohv1 "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"testing"
)

func TestAppImageWhenNotSetInMinimalCluster(t *testing.T) {
	g := NewGomegaWithT(t)

	// Use the specified yaml files to create a CoherenceCluster and trigger a fake end-to-end
	// reconcile to obtain the resources that would have been created by the Helm operator.
	result, cluster, err := CreateCluster("minimal-cluster.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	// Obtain the StatefulSet that Helm would have created
	sts, err := findStatefulSet(result, cluster, cohv1.DefaultRoleName)
	g.Expect(err).NotTo(HaveOccurred())

	// Obtain the Coherence container from the StatefulSet
	_, err = findInitContainer(sts, applicationContainer)
	g.Expect(errors.IsNotFound(err)).To(BeTrue())
}

func TestAppImageWhenSetForImplicitRole(t *testing.T) {
	g := NewGomegaWithT(t)

	result, cluster, err := CreateCluster("app-image-implicit-role.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	sts, err := findStatefulSet(result, cluster, cohv1.DefaultRoleName)
	g.Expect(err).NotTo(HaveOccurred())

	container, err := findInitContainer(sts, applicationContainer)
	g.Expect(err).NotTo(HaveOccurred())

	g.Expect(container.Image).To(Equal("app:1.0"))
}

func TestAppImageWhenNotSetInClusterWithExplicitRoles(t *testing.T) {
	g := NewGomegaWithT(t)

	result, cluster, err := CreateCluster("minimal-explicit-roles.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	for name := range cluster.GetRoles() {
		sts, err := findStatefulSet(result, cluster, name)
		g.Expect(err).NotTo(HaveOccurred())

		_, err = findInitContainer(sts, applicationContainer)
		g.Expect(errors.IsNotFound(err)).To(BeTrue())
	}
}

func TestAppImageWhenDefaultSetInClusterWithExplicitRoles(t *testing.T) {
	g := NewGomegaWithT(t)

	result, cluster, err := CreateCluster("app-image-explicit-roles-with-default.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	dataContainer, err := findInitContainerForRole(result, cluster, "data", applicationContainer)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(dataContainer.Image).To(Equal("app:1.0"))

	proxyContainer, err := findInitContainerForRole(result, cluster, "proxy", applicationContainer)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(proxyContainer.Image).To(Equal("app:1.0"))
}

func TestAppImageWhenDefaultSetAndOverriddenInClusterWithExplicitRoles(t *testing.T) {
	g := NewGomegaWithT(t)

	result, cluster, err := CreateCluster("app-image-explicit-roles-override-default.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	dataContainer, err := findInitContainerForRole(result, cluster, "data", applicationContainer)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(dataContainer.Image).To(Equal("app:1.0"))

	proxyContainer, err := findInitContainerForRole(result, cluster, "proxy", applicationContainer)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(proxyContainer.Image).To(Equal("app:2.0"))
}

func TestAppImageWhenSetInClusterWithExplicitRoles(t *testing.T) {
	g := NewGomegaWithT(t)

	result, cluster, err := CreateCluster("app-image-explicit-roles.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	dataContainer, err := findInitContainerForRole(result, cluster, "data", applicationContainer)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(dataContainer.Image).To(Equal("app:1.0"))

	proxyContainer, err := findInitContainerForRole(result, cluster, "proxy", applicationContainer)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(proxyContainer.Image).To(Equal("app:2.0"))
}
