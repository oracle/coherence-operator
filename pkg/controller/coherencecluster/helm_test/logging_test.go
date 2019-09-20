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
	"k8s.io/apimachinery/pkg/api/errors"
	"testing"
)

/*
 * These tests verify the various scenarios for setting logging configuration
 * in a CoherenceCluster.
 */

func TestLoggingInMinimalCluster(t *testing.T) {
	g := NewGomegaWithT(t)

	// Use the specified yaml files to create a CoherenceCluster and trigger a fake end-to-end
	// reconcile to obtain the resources that would have been created by the Helm operator.
	result, cluster, err := CreateCluster("minimal-cluster.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	// Obtain the StatefulSet that Helm would have created
	sts, err := findStatefulSet(result, cluster, cohv1.DefaultRoleName)
	g.Expect(err).NotTo(HaveOccurred())

	// Obtain the Coherence coherence from the StatefulSet
	coherence, err := findContainer(sts, coherenceContainer)
	g.Expect(err).NotTo(HaveOccurred())

	// Assert that the COH_LOG_LEVEL env-var is not set
	g.Expect(coherence.Env).NotTo(HaveEnvVarNamed("COH_LOG_LEVEL"))
	// Assert that the COH_LOGGING_CONFIG is the default
	g.Expect(coherence.Env).To(HaveEnvVar(corev1.EnvVar{Name: "COH_LOGGING_CONFIG", Value: "/scripts/logging.properties"}))

	// Assert that the logging config volume mount is not present
	_, err = findVolumeMount(coherence, "logging-config")
	g.Expect(errors.IsNotFound(err)).To(BeTrue())

	// Assert that the logging config volume is not present
	_, err = findVolume(sts, "logging-config")
	g.Expect(errors.IsNotFound(err)).To(BeTrue())

	// Assert that the fluentd config volume is not present
	_, err = findVolume(sts, "fluentd-coherence-conf")
	g.Expect(errors.IsNotFound(err)).To(BeTrue())

	// The Fluentd container should not be present
	_, err = findContainer(sts, fluentdContainer)
	g.Expect(errors.IsNotFound(err)).To(BeTrue())
}

// ----- Coherence Log Level ------------------------------------------------

func TestLoggingWithLogLevelSetForImplicitRole(t *testing.T) {
	assertLoggingWithLogLevel(t, cohv1.DefaultRoleName, "logging-level.yaml")
}

func TestLoggingWithLogLevelSetForExplicitRole(t *testing.T) {
	assertLoggingWithLogLevel(t, "data", "logging-explicit-role-level.yaml")
}

func TestLoggingWithLogLevelSetForExplicitRoleWithDefault(t *testing.T) {
	assertLoggingWithLogLevel(t, "data", "logging-explicit-role-with-default-level.yaml")
	assertLoggingWithLogLevel(t, "proxy", "logging-explicit-role-with-default-level.yaml")
}

func TestLoggingWithLogLevelSetForExplicitRoleWithOverride(t *testing.T) {
	assertLoggingWithLogLevel(t, "data", "logging-explicit-role-with-override-level.yaml")
	assertLoggingWithLogLevel(t, "proxy", "logging-explicit-role-with-override-level.yaml")
}

func assertLoggingWithLogLevel(t *testing.T, roleName, yamlFile string) {
	g := NewGomegaWithT(t)

	// Use the specified yaml files to create a CoherenceCluster and trigger a fake end-to-end
	// reconcile to obtain the resources that would have been created by the Helm operator.
	result, cluster, err := CreateCluster(yamlFile)
	g.Expect(err).NotTo(HaveOccurred())

	// Obtain the StatefulSet that Helm would have created
	sts, err := findStatefulSet(result, cluster, roleName)
	g.Expect(err).NotTo(HaveOccurred())

	// Obtain the Coherence coherence from the StatefulSet
	coherence, err := findContainer(sts, coherenceContainer)
	g.Expect(err).NotTo(HaveOccurred())

	// Assert that the COH_LOG_LEVEL env-var is set correctly
	g.Expect(coherence.Env).To(HaveEnvVar(corev1.EnvVar{Name: "COH_LOG_LEVEL", Value: "9"}))
	// Assert that the COH_LOGGING_CONFIG is the default
	g.Expect(coherence.Env).To(HaveEnvVar(corev1.EnvVar{Name: "COH_LOGGING_CONFIG", Value: "/scripts/logging.properties"}))

	// Assert that the logging config volume mount is not present
	_, err = findVolumeMount(coherence, "logging-config")
	g.Expect(errors.IsNotFound(err)).To(BeTrue())

	// Assert that the logging config volume is not present
	_, err = findVolume(sts, "logging-config")
	g.Expect(errors.IsNotFound(err)).To(BeTrue())

	// Assert that the fluentd config volume is not present
	_, err = findVolume(sts, "fluentd-coherence-conf")
	g.Expect(errors.IsNotFound(err)).To(BeTrue())

	// The Fluentd container should not be present
	_, err = findContainer(sts, fluentdContainer)
	g.Expect(errors.IsNotFound(err)).To(BeTrue())
}

// ----- Fluentd Enabled ----------------------------------------------------

func TestLoggingWithFluentdEnabledForImplicitRole(t *testing.T) {
	assertLoggingWithFluentdEnabled(t, cohv1.DefaultRoleName, "logging-fluentd-enabled.yaml")
}

func TestLoggingWithFluentdEnabledForExplicitRole(t *testing.T) {
	assertLoggingWithFluentdEnabled(t, "data", "logging-explicit-role-fluentd-enabled.yaml")
}

func TestLoggingWithFluentdEnabledForExplicitRoleWithDefault(t *testing.T) {
	assertLoggingWithFluentdEnabled(t, "data", "logging-explicit-role-with-default-fluentd-enabled.yaml")
	assertLoggingWithFluentdEnabled(t, "proxy", "logging-explicit-role-with-default-fluentd-enabled.yaml")
}

func TestLoggingWithFluentdEnabledForExplicitRoleWithOverride(t *testing.T) {
	assertLoggingWithFluentdEnabled(t, "data", "logging-explicit-role-with-override-fluentd-enabled.yaml")
	assertLoggingWithFluentdEnabled(t, "proxy", "logging-explicit-role-with-override-fluentd-enabled.yaml")
}

func assertLoggingWithFluentdEnabled(t *testing.T, roleName, yamlFile string) {
	g := NewGomegaWithT(t)

	// Use the specified yaml files to create a CoherenceCluster and trigger a fake end-to-end
	// reconcile to obtain the resources that would have been created by the Helm operator.
	result, cluster, err := CreateCluster(yamlFile)
	g.Expect(err).NotTo(HaveOccurred())

	// Obtain the StatefulSet that Helm would have created
	sts, err := findStatefulSet(result, cluster, roleName)
	g.Expect(err).NotTo(HaveOccurred())

	// Obtain the Coherence container from the StatefulSet
	coherence, err := findContainer(sts, coherenceContainer)
	g.Expect(err).NotTo(HaveOccurred())

	// Assert that the COH_LOG_LEVEL env-var is not set
	g.Expect(coherence.Env).NotTo(HaveEnvVarNamed("COH_LOG_LEVEL"))
	// Assert that the COH_LOGGING_CONFIG is the default
	g.Expect(coherence.Env).To(HaveEnvVar(corev1.EnvVar{Name: "COH_LOGGING_CONFIG", Value: "/scripts/logging.properties"}))

	// Assert that the logging config volume mount is not present
	_, err = findVolumeMount(coherence, "logging-config")
	g.Expect(errors.IsNotFound(err)).To(BeTrue())

	// Assert that the logging config volume is not present
	_, err = findVolume(sts, "logging-config")
	g.Expect(errors.IsNotFound(err)).To(BeTrue())

	// Assert that the fluentd config volume is present
	v, err := findVolume(sts, "fluentd-coherence-conf")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(v.ConfigMap).NotTo(BeNil())
	g.Expect(v.ConfigMap.Name).To(Equal(sts.Name + "-efk-config"))

	// Obtain the Fluentd container from the StatefulSet
	fluentd, err := findContainer(sts, fluentdContainer)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(fluentd).NotTo(BeNil())
}
