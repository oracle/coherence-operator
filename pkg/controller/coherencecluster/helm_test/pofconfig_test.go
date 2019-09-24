/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package helm_test

import (
	. "github.com/onsi/gomega"
	cohv1 "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	. "github.com/oracle/coherence-operator/test/e2e/helper/matchers"
	corev1 "k8s.io/api/core/v1"
	"testing"
)

/*
 * These tests verify the various scenarios for setting Coherence cache configuration
 * in a CoherenceCluster.
 */

func TestPofConfigWhenNotSetInMinimalCluster(t *testing.T) {
	g := NewGomegaWithT(t)

	// Use the specified yaml files to create a CoherenceCluster and trigger a fake end-to-end
	// reconcile to obtain the resources that would have been created by the Helm operator.
	result, cluster, err := CreateCluster("minimal-cluster.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	// Obtain the StatefulSet that Helm would have created
	sts, err := findStatefulSet(result, cluster, cohv1.DefaultRoleName)
	g.Expect(err).NotTo(HaveOccurred())

	// Obtain the Coherence container from the StatefulSet
	container, err := findContainer(sts, coherenceContainer)
	g.Expect(err).NotTo(HaveOccurred())

	// Assert that the COH_POF_CONFIG env-var is not set
	g.Expect(container.Env).NotTo(HaveEnvVarNamed("COH_POF_CONFIG"))
}

func TestPofConfigWhenSetForImplicitRole(t *testing.T) {
	g := NewGomegaWithT(t)

	result, cluster, err := CreateCluster("pofconfig-test-implicit-role.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	sts, err := findStatefulSet(result, cluster, cohv1.DefaultRoleName)
	g.Expect(err).NotTo(HaveOccurred())

	container, err := findContainer(sts, coherenceContainer)
	g.Expect(err).NotTo(HaveOccurred())

	g.Expect(container.Env).To(HaveEnvVar(corev1.EnvVar{Name: "COH_POF_CONFIG", Value: "test-pof-config.xml"}))
}

func TestPofConfigWhenNotSetInClusterWithExplicitRoles(t *testing.T) {
	g := NewGomegaWithT(t)

	result, cluster, err := CreateCluster("minimal-explicit-roles.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	for name := range cluster.GetRoles() {
		sts, err := findStatefulSet(result, cluster, name)
		g.Expect(err).NotTo(HaveOccurred())

		container, err := findContainer(sts, coherenceContainer)
		g.Expect(err).NotTo(HaveOccurred())

		g.Expect(container.Env).NotTo(HaveEnvVarNamed("COH_POF_CONFIG"))
	}
}

func TestPofConfigWhenDefaultSetInClusterWithExplicitRoles(t *testing.T) {
	g := NewGomegaWithT(t)

	result, cluster, err := CreateCluster("pofconfig-test-explicit-roles-with-default.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	dataContainer, err := findContainerForRole(result, cluster, "data", coherenceContainer)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(dataContainer.Env).To(HaveEnvVar(corev1.EnvVar{Name: "COH_POF_CONFIG", Value: "test-pof-config.xml"}))

	proxyContainer, err := findContainerForRole(result, cluster, "proxy", coherenceContainer)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(proxyContainer.Env).To(HaveEnvVar(corev1.EnvVar{Name: "COH_POF_CONFIG", Value: "test-pof-config.xml"}))
}

func TestPofConfigWhenDefaultSetAndOverriddentInClusterWithExplicitRoles(t *testing.T) {
	g := NewGomegaWithT(t)

	result, cluster, err := CreateCluster("pofconfig-test-explicit-roles-override-default.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	dataContainer, err := findContainerForRole(result, cluster, "data", coherenceContainer)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(dataContainer.Env).To(HaveEnvVar(corev1.EnvVar{Name: "COH_POF_CONFIG", Value: "test-pof-config.xml"}))

	proxyContainer, err := findContainerForRole(result, cluster, "proxy", coherenceContainer)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(proxyContainer.Env).To(HaveEnvVar(corev1.EnvVar{Name: "COH_POF_CONFIG", Value: "proxy-pof-config.xml"}))
}

func TestPofConfigWhenSetInClusterWithExplicitRoles(t *testing.T) {
	g := NewGomegaWithT(t)

	result, cluster, err := CreateCluster("pofconfig-test-explicit-roles.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	dataContainer, err := findContainerForRole(result, cluster, "data", coherenceContainer)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(dataContainer.Env).To(HaveEnvVar(corev1.EnvVar{Name: "COH_POF_CONFIG", Value: "data-pof-config.xml"}))

	proxyContainer, err := findContainerForRole(result, cluster, "proxy", coherenceContainer)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(proxyContainer.Env).To(HaveEnvVar(corev1.EnvVar{Name: "COH_POF_CONFIG", Value: "proxy-pof-config.xml"}))
}
