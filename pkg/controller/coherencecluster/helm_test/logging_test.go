/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package helm_test

import (
	. "github.com/onsi/gomega"
	cohv1 "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	"github.com/oracle/coherence-operator/pkg/fakes"
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

	assertLoggingDefaults(t, cohv1.DefaultRoleName, result, cluster)
}

func assertLoggingDefaults(t *testing.T, roleName string, result *fakes.HelmInstallResult, cluster *cohv1.CoherenceCluster) {
	g := NewGomegaWithT(t)

	// Obtain the StatefulSet that Helm would have created
	sts, err := findStatefulSet(result, cluster, roleName)
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
	assertLoggingWithLogLevel(t, "9", cohv1.DefaultRoleName, "logging-level.yaml")
}

func TestLoggingWithLogLevelSetForExplicitRole(t *testing.T) {
	assertLoggingWithLogLevel(t, "9", "data", "logging-level-explicit-role.yaml")
}

func TestLoggingWithLogLevelSetForExplicitRoleWithDefault(t *testing.T) {
	assertLoggingWithLogLevel(t, "9", "data", "logging-level-explicit-role-with-default.yaml")
	assertLoggingWithLogLevel(t, "9", "proxy", "logging-level-explicit-role-with-default.yaml")
}

func TestLoggingWithLogLevelSetForExplicitRoleWithOverride(t *testing.T) {
	assertLoggingWithLogLevel(t, "9", "data", "logging-level-explicit-role-with-override.yaml")
	assertLoggingWithLogLevel(t, "6", "proxy", "logging-level-explicit-role-with-override.yaml")
}

func assertLoggingWithLogLevel(t *testing.T, level, roleName, yamlFile string) {
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
	g.Expect(coherence.Env).To(HaveEnvVar(corev1.EnvVar{Name: "COH_LOG_LEVEL", Value: level}))
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

// ----- Logging Config File ------------------------------------------------

func TestLoggingWithLoggingConfigSetForImplicitRole(t *testing.T) {
	assertLoggingWithLoggingConfig(t, "test/logging.properties", cohv1.DefaultRoleName, "logging-config.yaml")
}

func TestLoggingWithLoggingConfigSetForExplicitRole(t *testing.T) {
	assertLoggingWithLoggingConfig(t, "test/logging.properties", "data", "logging-config-explicit-role.yaml")
}

func TestLoggingWithLogConfigSetForExplicitRoleWithDefault(t *testing.T) {
	assertLoggingWithLoggingConfig(t, "test/logging.properties", "data", "logging-config-explicit-role-with-default.yaml")
	assertLoggingWithLoggingConfig(t, "test/logging.properties", "proxy", "logging-config-explicit-role-with-default.yaml")
}

func TestLoggingWithLoggingConfigSetForExplicitRoleWithOverride(t *testing.T) {
	assertLoggingWithLoggingConfig(t, "test/logging.properties", "data", "logging-config-explicit-role-with-override.yaml")
	assertLoggingWithLoggingConfig(t, "test/default-logging.properties", "proxy", "logging-config-explicit-role-with-override.yaml")
}

func assertLoggingWithLoggingConfig(t *testing.T, config, roleName, yamlFile string) {
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

	// Assert that the COH_LOG_LEVEL env-var is not set
	g.Expect(coherence.Env).NotTo(HaveEnvVarNamed("COH_LOG_LEVEL"))
	// Assert that the COH_LOGGING_CONFIG is the default
	g.Expect(coherence.Env).To(HaveEnvVar(corev1.EnvVar{Name: "COH_LOGGING_CONFIG", Value: config}))

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

// ----- Logging ConfigMap --------------------------------------------------

func TestLoggingWithLoggingConfigMapSetForImplicitRole(t *testing.T) {
	assertLoggingWithLoggingConfigMap(t, "cm-logging", cohv1.DefaultRoleName, "logging-configmap.yaml")
}

func TestLoggingWithLoggingConfigMapSetForExplicitRole(t *testing.T) {
	assertLoggingWithLoggingConfigMap(t, "cm-logging", "data", "logging-configmap-explicit-role.yaml")
}

func TestLoggingWithLogConfigMapSetForExplicitRoleWithDefault(t *testing.T) {
	assertLoggingWithLoggingConfigMap(t, "cm-logging", "data", "logging-configmap-explicit-role-with-default.yaml")
	assertLoggingWithLoggingConfigMap(t, "cm-logging", "proxy", "logging-configmap-explicit-role-with-default.yaml")
}

func TestLoggingWithLoggingConfigMapSetForExplicitRoleWithOverride(t *testing.T) {
	assertLoggingWithLoggingConfigMap(t, "cm-logging", "data", "logging-configmap-explicit-role-with-override.yaml")
	assertLoggingWithLoggingConfigMap(t, "cm-default-logging", "proxy", "logging-configmap-explicit-role-with-override.yaml")
}

func assertLoggingWithLoggingConfigMap(t *testing.T, configMap, roleName, yamlFile string) {
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

	// Assert that the COH_LOG_LEVEL env-var is not set
	g.Expect(coherence.Env).NotTo(HaveEnvVarNamed("COH_LOG_LEVEL"))
	// Assert that the COH_LOGGING_CONFIG is the default
	g.Expect(coherence.Env).To(HaveEnvVar(corev1.EnvVar{Name: "COH_LOGGING_CONFIG", Value: "/scripts/logging.properties"}))

	// Assert that the logging config volume mount is present
	vm, err := findVolumeMount(coherence, "logging-config")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(vm.MountPath).To(Equal("/loggingconfig"))

	// Assert that the logging config volume is not present
	v, err := findVolume(sts, "logging-config")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(v.ConfigMap).NotTo(BeNil())
	g.Expect(v.ConfigMap.Name).To(Equal(configMap))

	// Assert that the fluentd config volume is not present
	_, err = findVolume(sts, "fluentd-coherence-conf")
	g.Expect(errors.IsNotFound(err)).To(BeTrue())

	// The Fluentd container should not be present
	_, err = findContainer(sts, fluentdContainer)
	g.Expect(errors.IsNotFound(err)).To(BeTrue())
}

// ----- Fluentd Enabled ----------------------------------------------------

func TestLoggingWithFluentdEnabledForImplicitRole(t *testing.T) {
	assertLoggingWithFluentdEnabled(t, true, cohv1.DefaultRoleName, "logging-fluentd-enabled.yaml")
}

func TestLoggingWithFluentdEnabledForExplicitRole(t *testing.T) {
	assertLoggingWithFluentdEnabled(t, true, "data", "logging-fluentd-enabled-explicit-role.yaml")
}

func TestLoggingWithFluentdEnabledForExplicitRoleWithDefault(t *testing.T) {
	assertLoggingWithFluentdEnabled(t, true, "data", "logging-fluentd-enabled-explicit-role-with-default.yaml")
	assertLoggingWithFluentdEnabled(t, true, "proxy", "logging-fluentd-enabled-explicit-role-with-default.yaml")
}

func TestLoggingWithFluentdEnabledForExplicitRoleWithOverride(t *testing.T) {
	assertLoggingWithFluentdEnabled(t, true, "data", "logging-fluentd-enabled-explicit-role-with-override.yaml")
	assertLoggingWithFluentdEnabled(t, false, "proxy", "logging-fluentd-enabled-explicit-role-with-override.yaml")
}

func assertLoggingWithFluentdEnabled(t *testing.T, enabled bool, roleName, yamlFile string) {
	g := NewGomegaWithT(t)

	// Use the specified yaml files to create a CoherenceCluster and trigger a fake end-to-end
	// reconcile to obtain the resources that would have been created by the Helm operator.
	result, cluster, err := CreateCluster(yamlFile)
	g.Expect(err).NotTo(HaveOccurred())

	if enabled {
		assertLoggingWithFluentd(t, roleName, result, cluster, fluentdImage, corev1.PullIfNotPresent)
	} else {
		assertLoggingDefaults(t, roleName, result, cluster)
	}
}

func assertLoggingWithFluentd(t *testing.T, roleName string, result *fakes.HelmInstallResult, cluster *cohv1.CoherenceCluster, imageName string, pullPolicy corev1.PullPolicy) {
	g := NewGomegaWithT(t)

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

	configMapName := sts.Name + "-efk-config"

	// Assert that the fluentd config volume is present
	v, err := findVolume(sts, "fluentd-coherence-conf")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(v.ConfigMap).NotTo(BeNil())
	g.Expect(v.ConfigMap.Name).To(Equal(configMapName))

	_, err = findConfigMap(result, configMapName)
	g.Expect(err).NotTo(HaveOccurred())

	// Obtain the Fluentd container from the StatefulSet
	fluentd, err := findContainer(sts, fluentdContainer)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(fluentd).NotTo(BeNil())
	g.Expect(fluentd.Image).To(Equal(imageName))
	g.Expect(fluentd.ImagePullPolicy).To(Equal(pullPolicy))
}

// ----- Fluentd Image ------------------------------------------------------

func TestLoggingWithFluentdImageForImplicitRole(t *testing.T) {
	assertLoggingWithFluentdImage(t, "fluentd:1.0", cohv1.DefaultRoleName, "logging-fluentd-image.yaml")
}

func TestLoggingWithFluentImageForExplicitRole(t *testing.T) {
	assertLoggingWithFluentdImage(t, "fluentd:1.0", "data", "logging-fluentd-image-explicit-role.yaml")
}

func TestLoggingWithFluentdImageForExplicitRoleWithDefault(t *testing.T) {
	assertLoggingWithFluentdImage(t, "fluentd:1.0", "data", "logging-fluentd-image-explicit-role-with-default.yaml")
	assertLoggingWithFluentdImage(t, "fluentd:1.0", "proxy", "logging-fluentd-image-explicit-role-with-default.yaml")
}

func TestLoggingWithFluentdImageForExplicitRoleWithOverride(t *testing.T) {
	assertLoggingWithFluentdImage(t, "fluentd:1.0", "data", "logging-fluentd-image-explicit-role-with-override.yaml")
	assertLoggingWithFluentdImage(t, "fluentd:2.0", "proxy", "logging-fluentd-image-explicit-role-with-override.yaml")
}

func assertLoggingWithFluentdImage(t *testing.T, imageName, roleName, yamlFile string) {
	g := NewGomegaWithT(t)

	result, cluster, err := CreateCluster(yamlFile)
	g.Expect(err).NotTo(HaveOccurred())

	assertLoggingWithFluentd(t, roleName, result, cluster, imageName, corev1.PullIfNotPresent)
}

// ----- Fluentd ImagePullPolicy --------------------------------------------

func TestLoggingWithFluentdImagePullPolicyForImplicitRole(t *testing.T) {
	assertLoggingWithFluentdImagePullPolicy(t, corev1.PullNever, cohv1.DefaultRoleName, "logging-fluentd-image-pull-policy.yaml")
}

func TestLoggingWithFluentImagePullPolicyForExplicitRole(t *testing.T) {
	assertLoggingWithFluentdImagePullPolicy(t, corev1.PullNever, "data", "logging-fluentd-image-pull-policy-explicit-role.yaml")
}

func TestLoggingWithFluentdImagePullPolicyForExplicitRoleWithDefault(t *testing.T) {
	assertLoggingWithFluentdImagePullPolicy(t, corev1.PullNever, "data", "logging-fluentd-image-pull-policy-explicit-role-with-default.yaml")
	assertLoggingWithFluentdImagePullPolicy(t, corev1.PullNever, "proxy", "logging-fluentd-image-pull-policy-explicit-role-with-default.yaml")
}

func TestLoggingWithFluentdImagePullPolicyForExplicitRoleWithOverride(t *testing.T) {
	assertLoggingWithFluentdImagePullPolicy(t, corev1.PullNever, "data", "logging-fluentd-image-pull-policy-explicit-role-with-override.yaml")
	assertLoggingWithFluentdImagePullPolicy(t, corev1.PullAlways, "proxy", "logging-fluentd-image-pull-policy-explicit-role-with-override.yaml")
}

func assertLoggingWithFluentdImagePullPolicy(t *testing.T, policy corev1.PullPolicy, roleName, yamlFile string) {
	g := NewGomegaWithT(t)

	result, cluster, err := CreateCluster(yamlFile)
	g.Expect(err).NotTo(HaveOccurred())

	assertLoggingWithFluentd(t, roleName, result, cluster, fluentdImage, policy)
}
