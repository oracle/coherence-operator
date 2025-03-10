/*
 * Copyright (c) 2020, 2025, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1_test

import (
	coh "github.com/oracle/coherence-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	"strconv"
	"testing"
)

func TestCreateStatefulSetWithCoherenceMetricsEmpty(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		Coherence: &coh.CoherenceSpec{
			Metrics: &coh.PortSpecWithSSL{},
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithCoherenceMetricsEnabledFalse(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		Coherence: &coh.CoherenceSpec{
			Metrics: &coh.PortSpecWithSSL{
				Enabled: boolPtr(false),
				Port:    int32Ptr(1234),
				SSL: &coh.SSLSpec{
					Enabled:                boolPtr(true),
					Secrets:                stringPtr("ssl-secret"),
					KeyStore:               stringPtr("ssl-keystore.jks"),
					KeyStorePasswordFile:   stringPtr("ssl-key-pass.txt"),
					KeyPasswordFile:        stringPtr("ssl-pass.txt"),
					KeyStoreAlgorithm:      stringPtr("ssl-key-algo"),
					KeyStoreProvider:       stringPtr("ssl-key-provider"),
					KeyStoreType:           stringPtr("ssl-key-type"),
					TrustStore:             stringPtr("ssl-trust.jks"),
					TrustStorePasswordFile: stringPtr("ssl-trust-pass.txt"),
					TrustStoreAlgorithm:    stringPtr("ssl-key-algo"),
					TrustStoreProvider:     stringPtr("ssl-trust-provider"),
					TrustStoreType:         stringPtr("ssl-trust-type"),
					RequireClientCert:      boolPtr(true),
				},
			},
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithCoherenceMetricsEnabledTrue(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		Coherence: &coh.CoherenceSpec{
			Metrics: &coh.PortSpecWithSSL{
				Enabled: boolPtr(true),
			},
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	addEnvVars(stsExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: "COHERENCE_METRICS_ENABLED", Value: "true"})
	addEnvVars(stsExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: "COHERENCE_METRICS_PORT", Value: strconv.FormatInt(int64(coh.DefaultMetricsPort), 10)})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithCoherenceMetricsEnabledWithPort(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		Coherence: &coh.CoherenceSpec{
			Metrics: &coh.PortSpecWithSSL{
				Enabled: boolPtr(true),
				Port:    int32Ptr(1234),
			},
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	addEnvVars(stsExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: "COHERENCE_METRICS_ENABLED", Value: "true"})
	addEnvVars(stsExpected, coh.ContainerNameCoherence, corev1.EnvVar{Name: "COHERENCE_METRICS_PORT", Value: "1234"})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithCoherenceMetricsWithSSLEnabledWithoutSecret(t *testing.T) {

	spec := coh.CoherenceResourceSpec{
		Coherence: &coh.CoherenceSpec{
			Metrics: &coh.PortSpecWithSSL{
				Enabled: boolPtr(true),
				SSL: &coh.SSLSpec{
					Enabled:                boolPtr(true),
					KeyStore:               stringPtr("ssl-keystore.jks"),
					KeyStorePasswordFile:   stringPtr("ssl-key-pass.txt"),
					KeyPasswordFile:        stringPtr("ssl-pass.txt"),
					KeyStoreAlgorithm:      stringPtr("ssl-key-algo"),
					KeyStoreProvider:       stringPtr("ssl-key-provider"),
					KeyStoreType:           stringPtr("ssl-key-type"),
					TrustStore:             stringPtr("ssl-trust.jks"),
					TrustStorePasswordFile: stringPtr("ssl-trust-pass.txt"),
					TrustStoreAlgorithm:    stringPtr("ssl-trust-algo"),
					TrustStoreProvider:     stringPtr("ssl-trust-provider"),
					TrustStoreType:         stringPtr("ssl-trust-type"),
					RequireClientCert:      boolPtr(true),
				},
			},
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	addEnvVars(stsExpected, coh.ContainerNameCoherence,
		corev1.EnvVar{Name: "COHERENCE_METRICS_ENABLED", Value: "true"},
		corev1.EnvVar{Name: "COHERENCE_METRICS_PORT", Value: strconv.FormatInt(int64(coh.DefaultMetricsPort), 10)},
		corev1.EnvVar{Name: "COHERENCE_METRICS_SSL_ENABLED", Value: "true"},
		corev1.EnvVar{Name: "COHERENCE_METRICS_SSL_KEYSTORE", Value: "ssl-keystore.jks"},
		corev1.EnvVar{Name: "COHERENCE_METRICS_SSL_KEYSTORE_PASSWORD_FILE", Value: "ssl-key-pass.txt"},
		corev1.EnvVar{Name: "COHERENCE_METRICS_SSL_KEY_PASSWORD_FILE", Value: "ssl-pass.txt"},
		corev1.EnvVar{Name: "COHERENCE_METRICS_SSL_KEYSTORE_TYPE", Value: "ssl-key-type"},
		corev1.EnvVar{Name: "COHERENCE_METRICS_SSL_KEYSTORE_ALGORITHM", Value: "ssl-key-algo"},
		corev1.EnvVar{Name: "COHERENCE_METRICS_SSL_KEYSTORE_PROVIDER", Value: "ssl-key-provider"},
		corev1.EnvVar{Name: "COHERENCE_METRICS_SSL_TRUSTSTORE", Value: "ssl-trust.jks"},
		corev1.EnvVar{Name: "COHERENCE_METRICS_SSL_TRUSTSTORE_PASSWORD_FILE", Value: "ssl-trust-pass.txt"},
		corev1.EnvVar{Name: "COHERENCE_METRICS_SSL_TRUSTSTORE_TYPE", Value: "ssl-trust-type"},
		corev1.EnvVar{Name: "COHERENCE_METRICS_SSL_TRUSTSTORE_ALGORITHM", Value: "ssl-trust-algo"},
		corev1.EnvVar{Name: "COHERENCE_METRICS_SSL_TRUSTSTORE_PROVIDER", Value: "ssl-trust-provider"},
		corev1.EnvVar{Name: "COHERENCE_METRICS_SSL_REQUIRE_CLIENT_CERT", Value: "true"})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}

func TestCreateStatefulSetWithCoherenceMetricsWithSSLEnabledWithSecret(t *testing.T) {

	secretName := "test-ssl-secret"

	spec := coh.CoherenceResourceSpec{
		Coherence: &coh.CoherenceSpec{
			Metrics: &coh.PortSpecWithSSL{
				Enabled: boolPtr(true),
				SSL: &coh.SSLSpec{
					Enabled:                boolPtr(true),
					Secrets:                &secretName,
					KeyStore:               stringPtr("ssl-keystore.jks"),
					KeyStorePasswordFile:   stringPtr("ssl-key-pass.txt"),
					KeyPasswordFile:        stringPtr("ssl-pass.txt"),
					KeyStoreAlgorithm:      stringPtr("ssl-key-algo"),
					KeyStoreProvider:       stringPtr("ssl-key-provider"),
					KeyStoreType:           stringPtr("ssl-key-type"),
					TrustStore:             stringPtr("ssl-trust.jks"),
					TrustStorePasswordFile: stringPtr("ssl-trust-pass.txt"),
					TrustStoreAlgorithm:    stringPtr("ssl-trust-algo"),
					TrustStoreProvider:     stringPtr("ssl-trust-provider"),
					TrustStoreType:         stringPtr("ssl-trust-type"),
					RequireClientCert:      boolPtr(true),
				},
			},
		},
	}

	// Create the test deployment
	deployment := createTestDeployment(spec)
	// Create expected StatefulSet
	stsExpected := createMinimalExpectedStatefulSet(deployment)
	addEnvVars(stsExpected, coh.ContainerNameCoherence,
		corev1.EnvVar{Name: "COHERENCE_METRICS_SSL_CERTS", Value: coh.VolumeMountPathMetricsCerts},
		corev1.EnvVar{Name: "COHERENCE_METRICS_ENABLED", Value: "true"},
		corev1.EnvVar{Name: "COHERENCE_METRICS_PORT", Value: strconv.FormatInt(int64(coh.DefaultMetricsPort), 10)},
		corev1.EnvVar{Name: "COHERENCE_METRICS_SSL_ENABLED", Value: "true"},
		corev1.EnvVar{Name: "COHERENCE_METRICS_SSL_KEYSTORE", Value: "ssl-keystore.jks"},
		corev1.EnvVar{Name: "COHERENCE_METRICS_SSL_KEYSTORE_PASSWORD_FILE", Value: "ssl-key-pass.txt"},
		corev1.EnvVar{Name: "COHERENCE_METRICS_SSL_KEY_PASSWORD_FILE", Value: "ssl-pass.txt"},
		corev1.EnvVar{Name: "COHERENCE_METRICS_SSL_KEYSTORE_TYPE", Value: "ssl-key-type"},
		corev1.EnvVar{Name: "COHERENCE_METRICS_SSL_KEYSTORE_ALGORITHM", Value: "ssl-key-algo"},
		corev1.EnvVar{Name: "COHERENCE_METRICS_SSL_KEYSTORE_PROVIDER", Value: "ssl-key-provider"},
		corev1.EnvVar{Name: "COHERENCE_METRICS_SSL_TRUSTSTORE", Value: "ssl-trust.jks"},
		corev1.EnvVar{Name: "COHERENCE_METRICS_SSL_TRUSTSTORE_PASSWORD_FILE", Value: "ssl-trust-pass.txt"},
		corev1.EnvVar{Name: "COHERENCE_METRICS_SSL_TRUSTSTORE_TYPE", Value: "ssl-trust-type"},
		corev1.EnvVar{Name: "COHERENCE_METRICS_SSL_TRUSTSTORE_ALGORITHM", Value: "ssl-trust-algo"},
		corev1.EnvVar{Name: "COHERENCE_METRICS_SSL_TRUSTSTORE_PROVIDER", Value: "ssl-trust-provider"},
		corev1.EnvVar{Name: "COHERENCE_METRICS_SSL_REQUIRE_CLIENT_CERT", Value: "true"})

	// add the metrics ConfigMap volume mount
	stsExpected.Spec.Template.Spec.Containers[0].VolumeMounts = append(stsExpected.Spec.Template.Spec.Containers[0].VolumeMounts, corev1.VolumeMount{
		Name:      coh.VolumeNameMetricsSSL,
		MountPath: coh.VolumeMountPathMetricsCerts,
		ReadOnly:  true,
	})
	// add the metrics ConfigMap volume
	stsExpected.Spec.Template.Spec.Volumes = append(stsExpected.Spec.Template.Spec.Volumes, corev1.Volume{
		Name: coh.VolumeNameMetricsSSL,
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName:  secretName,
				DefaultMode: int32Ptr(0777),
			},
		},
	})

	// assert that the StatefulSet is as expected
	assertStatefulSetCreation(t, deployment, stsExpected)
}
