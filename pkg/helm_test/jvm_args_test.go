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

func TestJvmArgWhenNotSetInMinimalCluster(t *testing.T) {
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

	// Assert that the JVM_ARGS env-var is not set
	g.Expect(container.Env).NotTo(HaveEnvVarNamed("JVM_ARGS"))
}

func TestJvmArgWhenNotSetInClusterWithExplicitRoles(t *testing.T) {
	g := NewGomegaWithT(t)

	result, cluster, err := CreateCluster("minimal-explicit-roles.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	for name := range cluster.GetRoles() {
		sts, err := findStatefulSet(result, cluster, name)
		g.Expect(err).NotTo(HaveOccurred())

		container, err := findContainer(sts, coherenceContainer)
		g.Expect(err).NotTo(HaveOccurred())

		g.Expect(container.Env).NotTo(HaveEnvVarNamed("JVM_ARGS"))
	}
}

func TestJvmSingleArgWhenSetForImplicitRole(t *testing.T) {
	g := NewGomegaWithT(t)

	result, cluster, err := CreateCluster("jvm-single-arg-implicit-role.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	sts, err := findStatefulSet(result, cluster, cohv1.DefaultRoleName)
	g.Expect(err).NotTo(HaveOccurred())

	container, err := findContainer(sts, coherenceContainer)
	g.Expect(err).NotTo(HaveOccurred())

	g.Expect(container.Env).To(HaveEnvVar(corev1.EnvVar{Name: "JVM_ARGS", Value: "-DargOne=1"}))
}

func TestJvmSingleArgWhenDefaultSetInClusterWithExplicitRoles(t *testing.T) {
	g := NewGomegaWithT(t)

	result, cluster, err := CreateCluster("jvm-single-arg-explicit-roles-with-default.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	dataContainer, err := findContainerForRole(result, cluster, "data", coherenceContainer)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(dataContainer.Env).To(HaveEnvVar(corev1.EnvVar{Name: "JVM_ARGS", Value: "-DargOne=1"}))

	proxyContainer, err := findContainerForRole(result, cluster, "proxy", coherenceContainer)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(proxyContainer.Env).To(HaveEnvVar(corev1.EnvVar{Name: "JVM_ARGS", Value: "-DargOne=1"}))
}

func TestJvmSingleArgWhenDefaultSetAndOverriddenInClusterWithExplicitRoles(t *testing.T) {
	g := NewGomegaWithT(t)

	result, cluster, err := CreateCluster("jvm-single-arg-explicit-roles-override-default.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	dataContainer, err := findContainerForRole(result, cluster, "data", coherenceContainer)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(dataContainer.Env).To(HaveEnvVar(corev1.EnvVar{Name: "JVM_ARGS", Value: "-DargOne=1"}))

	proxyContainer, err := findContainerForRole(result, cluster, "proxy", coherenceContainer)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(proxyContainer.Env).To(HaveEnvVar(corev1.EnvVar{Name: "JVM_ARGS", Value: "-DargTwo=2 -DargOne=1"}))
}

func TestJvmSingleArgWhenSetInClusterWithExplicitRoles(t *testing.T) {
	g := NewGomegaWithT(t)

	result, cluster, err := CreateCluster("jvm-single-arg-explicit-roles.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	dataContainer, err := findContainerForRole(result, cluster, "data", coherenceContainer)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(dataContainer.Env).To(HaveEnvVar(corev1.EnvVar{Name: "JVM_ARGS", Value: "-DargOne=1"}))

	proxyContainer, err := findContainerForRole(result, cluster, "proxy", coherenceContainer)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(proxyContainer.Env).To(HaveEnvVar(corev1.EnvVar{Name: "JVM_ARGS", Value: "-DargTwo=2"}))
}

func TestJvmMultipleArgsWhenSetForImplicitRole(t *testing.T) {
	g := NewGomegaWithT(t)

	result, cluster, err := CreateCluster("jvm-multiple-args-implicit-role.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	sts, err := findStatefulSet(result, cluster, cohv1.DefaultRoleName)
	g.Expect(err).NotTo(HaveOccurred())

	container, err := findContainer(sts, coherenceContainer)
	g.Expect(err).NotTo(HaveOccurred())

	g.Expect(container.Env).To(HaveEnvVar(corev1.EnvVar{Name: "JVM_ARGS", Value: "-DargOne=1 -DargTwo=2"}))
}

func TestJvmMultipleArgsWhenDefaultSetInClusterWithExplicitRoles(t *testing.T) {
	g := NewGomegaWithT(t)

	result, cluster, err := CreateCluster("jvm-multiple-args-explicit-roles-with-default.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	dataContainer, err := findContainerForRole(result, cluster, "data", coherenceContainer)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(dataContainer.Env).To(HaveEnvVar(corev1.EnvVar{Name: "JVM_ARGS", Value: "-DargOne=1 -DargTwo=2"}))

	proxyContainer, err := findContainerForRole(result, cluster, "proxy", coherenceContainer)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(proxyContainer.Env).To(HaveEnvVar(corev1.EnvVar{Name: "JVM_ARGS", Value: "-DargOne=1 -DargTwo=2"}))
}

func TestJvmMultipleArgsWhenDefaultSetAndOverriddenInClusterWithExplicitRoles(t *testing.T) {
	g := NewGomegaWithT(t)

	result, cluster, err := CreateCluster("jvm-multiple-args-explicit-roles-override-default.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	dataContainer, err := findContainerForRole(result, cluster, "data", coherenceContainer)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(dataContainer.Env).To(HaveEnvVar(corev1.EnvVar{Name: "JVM_ARGS", Value: "-DargOne=1 -DargTwo=2"}))

	proxyContainer, err := findContainerForRole(result, cluster, "proxy", coherenceContainer)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(proxyContainer.Env).To(HaveEnvVar(corev1.EnvVar{Name: "JVM_ARGS", Value: "-DargThree=3 -DargFour=4 -DargOne=1 -DargTwo=2"}))
}

func TestJvmMultipleArgsWhenSetInClusterWithExplicitRoles(t *testing.T) {
	g := NewGomegaWithT(t)

	result, cluster, err := CreateCluster("jvm-multiple-args-explicit-roles.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	dataContainer, err := findContainerForRole(result, cluster, "data", coherenceContainer)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(dataContainer.Env).To(HaveEnvVar(corev1.EnvVar{Name: "JVM_ARGS", Value: "-DargOne=1 -DargTwo=2"}))

	proxyContainer, err := findContainerForRole(result, cluster, "proxy", coherenceContainer)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(proxyContainer.Env).To(HaveEnvVar(corev1.EnvVar{Name: "JVM_ARGS", Value: "-DargThree=3 -DargFour=4"}))
}
