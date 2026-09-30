/*
 * Copyright (c) 2026, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */
package runner

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	coh "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/pkg/runner/run_details"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/ptr"
)

func TestManagedHealthCoherenceTLSMappingWithExistingJAAS(t *testing.T) {
	directory := t.TempDir()
	material := writeManagedHealthMaterial(t, directory)
	connection := managedHealthConnection(material, "existing")
	connection.ServerTLS.KeyStoreProvider = ptr.To("SunJSSE")
	connection.ServerTLS.TrustStore = ptr.To("legacy-truststore-not-mounted.jks")
	connection.ServerTLS.TrustStorePasswordFile = ptr.To("legacy-truststore-password-not-mounted")
	connection.ServerTLS.TrustStoreType = ptr.To("JKS")
	connection.ServerTLS.TrustStoreAlgorithm = ptr.To("SunX509")
	connection.ServerTLS.TrustStoreProvider = ptr.To("LegacyProvider")
	details := newManagedHealthDetails(t, directory, connection, map[string]string{
		coh.EnvVarJvmArgs: "-Djava.security.auth.login.config=/application/jaas.conf",
	})

	if err := configureHealthMutators(details); err != nil {
		t.Fatal(err)
	}
	arguments := details.GetArguments()
	args := strings.Join(arguments, " ")
	for _, required := range []string{
		"-Dcoherence.health.http.mutators.enabled=true",
		"-Dcoherence.health.http.auth=basic",
		"-Dcoherence.health.http.provider=HealthSSLProvider",
		"-Dcoherence.health.security.keystore=file:" + material,
		"-Dcoherence.health.security.keystore.password=" + material,
		"-Dcoherence.health.security.key.password=" + material,
		"-Dcoherence.health.security.keystore.type=JKS",
		"-Dcoherence.health.security.keystore.algorithm=SunX509",
		"-Dcoherence.health.security.keystore.provider=SunJSSE",
	} {
		if count := argumentCount(arguments, required); count != 1 {
			t.Errorf("expected exactly one %s argument, found %d in %v", required, count, arguments)
		}
	}
	if fileQualified := "-Dcoherence.health.security.keystore.provider=file:SunJSSE"; argumentCount(arguments, fileQualified) != 0 {
		t.Fatalf("keystore provider was emitted as a file URL: %v", arguments)
	}
	if strings.Contains(args, "coherence.operator.health") {
		t.Fatal("managed Coherence security configured the legacy Operator server")
	}
	if strings.Contains(args, "coherence.health.security.truststore") {
		t.Fatal("one-way health TLS configured a server truststore from legacy fields")
	}
}

func TestManagedHealthOmitsKeyStoreProviderWhenUnset(t *testing.T) {
	directory := t.TempDir()
	connection := managedHealthConnection(writeManagedHealthMaterial(t, directory), "existing")
	details := newManagedHealthDetails(t, directory, connection, map[string]string{
		coh.EnvVarJvmArgs: "-Djava.security.auth.login.config=/application/jaas.conf",
	})

	if err := configureHealthMutators(details); err != nil {
		t.Fatal(err)
	}
	for _, arg := range details.GetArguments() {
		if strings.HasPrefix(arg, "-Dcoherence.health.security.keystore.provider=") {
			t.Fatalf("unset keystore provider emitted argument %q", arg)
		}
	}
}

func TestManagedHealthPackagedJAAS(t *testing.T) {
	directory := t.TempDir()
	authDirectory := filepath.Join(directory, "auth")
	tlsDirectory := filepath.Join(directory, "tls")
	previousAuth, previousTLS := healthAuthDirectory, healthTLSDirectory
	healthAuthDirectory, healthTLSDirectory = authDirectory, tlsDirectory
	t.Cleanup(func() {
		healthAuthDirectory, healthTLSDirectory = previousAuth, previousTLS
	})
	for name, value := range map[string]string{"username": "health-operator", "password": "secret"} {
		path := filepath.Join(authDirectory, name)
		if err := os.MkdirAll(path, 0700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(path, "value"), []byte(value), 0600); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.MkdirAll(tlsDirectory, 0700); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"server.jks", "store-password"} {
		if err := os.WriteFile(filepath.Join(tlsDirectory, name), []byte("fixture"), 0600); err != nil {
			t.Fatal(err)
		}
	}
	connection := managedHealthConnection("server.jks", "packaged")
	connection.ServerTLS.Secrets = ptr.To("health-stores-v1")
	connection.ServerTLS.KeyStorePasswordFile = ptr.To("store-password")
	connection.ServerTLS.KeyPasswordFile = ptr.To("store-password")
	details := newManagedHealthDetails(t, directory, connection, nil)

	if err := configureHealthMutators(details); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(details.UtilsDir, "health-jaas.conf"))
	if err != nil {
		t.Fatal(err)
	}
	config := string(data)
	if !strings.Contains(config, "CoherenceREST {") || !strings.Contains(config, "SecretLoginModule") || !strings.Contains(config, filepath.Join(authDirectory, "username", "value")) {
		t.Fatalf("invalid packaged JAAS configuration: %s", config)
	}
	if !strings.Contains(strings.Join(details.GetArguments(), " "), "-Djava.security.auth.login.config="+filepath.Join(details.UtilsDir, "health-jaas.conf")) {
		t.Fatal("packaged JAAS file was not added to the member JVM")
	}
}

func TestManagedHealthRejectsOperatorServer(t *testing.T) {
	directory := t.TempDir()
	connection := managedHealthConnection(writeManagedHealthMaterial(t, directory), "existing")
	details := newManagedHealthDetails(t, directory, connection, map[string]string{
		coh.EnvVarJvmArgs: "-Djava.security.auth.login.config=/application/jaas.conf",
	})
	details.ResolvedOperatorHealth = true
	if err := configureHealthMutators(details); err == nil || !strings.Contains(err.Error(), "require the Coherence health server") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestManagedHealthRejectsIncompleteOrConflictingConfiguration(t *testing.T) {
	for _, raw := range []string{`{"secure":true}`, `{"secure":false,"jaas":{"mode":"packaged"}}`, `invalid`} {
		if err := configureHealthMutators(newSSLRunDetails(map[string]string{coh.HealthMutatorEnv: raw})); err == nil {
			t.Fatalf("accepted incomplete configuration %s", raw)
		}
	}

	directory := t.TempDir()
	connection := managedHealthConnection(writeManagedHealthMaterial(t, directory), "existing")
	for _, property := range []string{
		"-Dcoherence.health.http.auth=none",
		"-Dcoherence.health.security.keystore=/application/server.jks",
	} {
		details := newManagedHealthDetails(t, directory, connection, map[string]string{
			coh.EnvVarJvmArgs: property + " -Djava.security.auth.login.config=/application/jaas.conf",
		})
		if err := configureHealthMutators(details); err == nil || !strings.Contains(err.Error(), "conflict") {
			t.Fatalf("accepted conflicting property %s: %v", property, err)
		}
	}
}

func TestHealthSecurityCheckUsesApplicationLauncher(t *testing.T) {
	directory := t.TempDir()
	details := newManagedHealthDetails(t, directory, nil, nil)
	details.AddClasspath("/application/classes")
	command, err := createHealthSecurityCheckCommand(context.Background(), details)
	if err != nil {
		t.Fatal(err)
	}
	joined := strings.Join(command.Args, " ")
	if !strings.Contains(joined, "-cp") || !strings.Contains(joined, "/application/classes") || !strings.HasSuffix(joined, healthSecuritySupportClass) {
		t.Fatalf("invalid Java capability command: %s", joined)
	}

	spring := newManagedHealthDetails(t, directory, nil, map[string]string{
		coh.EnvVarAppType:          coh.AppTypeSpring3,
		coh.EnvVarSpringBootFatJar: "/application/app.jar",
	})
	spring.AppType = coh.AppTypeSpring3
	command, err = createHealthSecurityCheckCommand(context.Background(), spring)
	if err != nil {
		t.Fatal(err)
	}
	joined = strings.Join(command.Args, " ")
	for _, value := range []string{"-Dloader.main=" + healthSecuritySupportClass, "-cp /application/app.jar", coh.SpringBootMain3} {
		if !strings.Contains(joined, value) {
			t.Fatalf("Spring Boot capability command %q is missing %q", joined, value)
		}
	}

	unpackagedSpring := newManagedHealthDetails(t, directory, nil, map[string]string{coh.EnvVarAppType: coh.AppTypeSpring3})
	unpackagedSpring.AppType = coh.AppTypeSpring3
	if _, err := createHealthSecurityCheckCommand(context.Background(), unpackagedSpring); err == nil || !strings.Contains(err.Error(), "configured Spring Boot fat JAR") {
		t.Fatalf("Spring Boot without a configured fat JAR was not rejected: %v", err)
	}

	buildpack := newManagedHealthDetails(t, directory, nil, map[string]string{
		coh.EnvVarCnbpEnabled: "true",
		coh.EnvVarAppType:     coh.AppTypeSpring3,
	})
	buildpack.AppType = coh.AppTypeSpring3
	if _, err := createHealthSecurityCheckCommand(context.Background(), buildpack); err == nil || !strings.Contains(err.Error(), "do not support the Cloud Native Buildpack launcher") {
		t.Fatalf("buildpack packaging was not rejected: %v", err)
	}
}

func managedHealthConnection(material, mode string) *coh.HealthMutatorConnection {
	selector := corev1.SecretKeySelector{LocalObjectReference: corev1.LocalObjectReference{Name: "auth-v1"}, Key: "value"}
	return &coh.HealthMutatorConnection{
		Secure:    true,
		BasicAuth: &coh.HTTPBasicAuth{Username: selector, Password: selector},
		TLS:       &coh.HTTPClientTLS{CASecret: selector, ServerName: "health.example"},
		JAAS:      &coh.HealthJAAS{Mode: mode},
		ServerTLS: &coh.SSLSpec{
			Enabled:              ptr.To(true),
			KeyStore:             ptr.To(material),
			KeyStorePasswordFile: ptr.To(material),
			KeyPasswordFile:      ptr.To(material),
			KeyStoreType:         ptr.To("JKS"),
			KeyStoreAlgorithm:    ptr.To("SunX509"),
			RequireClientCert:    ptr.To(false),
		},
	}
}

func newManagedHealthDetails(t *testing.T, directory string, connection *coh.HealthMutatorConnection, env map[string]string) *run_details.RunDetails {
	t.Helper()
	bin := filepath.Join(directory, "bin")
	if err := os.MkdirAll(bin, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(bin, "java"), []byte("#!/bin/sh\nexit 0\n"), 0700); err != nil {
		t.Fatal(err)
	}
	utils := filepath.Join(directory, "utils")
	if err := os.MkdirAll(utils, 0700); err != nil {
		t.Fatal(err)
	}
	values := map[string]string{coh.EnvVarCohUtilDir: utils}
	for name, value := range env {
		values[name] = value
	}
	if connection != nil {
		data, err := json.Marshal(connection)
		if err != nil {
			t.Fatal(err)
		}
		values[coh.HealthMutatorEnv] = string(data)
	}
	details := newSSLRunDetails(values)
	details.JavaHome = directory
	details.Classpath = filepath.Join(directory, "operator.jar")
	return details
}

func writeManagedHealthMaterial(t *testing.T, directory string) string {
	t.Helper()
	material := filepath.Join(directory, "store")
	if err := os.WriteFile(material, []byte("fixture"), 0600); err != nil {
		t.Fatal(err)
	}
	return material
}

func argumentCount(args []string, expected string) int {
	count := 0
	for _, arg := range args {
		if arg == expected {
			count++
		}
	}
	return count
}
