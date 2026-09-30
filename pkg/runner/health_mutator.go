/*
 * Copyright (c) 2026, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */
package runner

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	v1 "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/pkg/runner/run_details"
)

const healthSecuritySupportClass = "com.oracle.coherence.k8s.HealthSecuritySupport"

var (
	healthAuthDirectory = "/coherence/health-auth"
	healthTLSDirectory  = "/coherence/health-tls"
)

func configureHealthMutators(details *run_details.RunDetails) error {
	raw := details.Getenv(v1.HealthMutatorEnv)
	if raw == "" {
		return nil
	}
	connection := &v1.HealthMutatorConnection{}
	if err := json.Unmarshal([]byte(raw), connection); err != nil {
		return fmt.Errorf("invalid managed health configuration")
	}
	if err := connection.Validate(); err != nil {
		return err
	}
	if !connection.Secure {
		return nil
	}
	if details.ResolvedOperatorHealth {
		return fmt.Errorf("managed secure health mutators require the Coherence health server")
	}
	contextName := "CoherenceREST"
	mode := "packaged"
	if connection.JAAS != nil {
		if connection.JAAS.Mode != "" {
			mode = connection.JAAS.Mode
		}
	}
	// Probe capability without invoking a mutator. Unknown classpaths fail closed.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cmd, err := createHealthSecurityCheckCommand(ctx, details)
	if err != nil {
		return err
	}
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("coherence health server security capability could not be established")
	}
	userArgs := details.Getenv(v1.EnvVarJvmArgs) + " " + details.Getenv("JAVA_TOOL_OPTIONS") + " " + details.Getenv("JDK_JAVA_OPTIONS") + " " + details.Getenv("_JAVA_OPTIONS")
	for _, property := range []string{"coherence.health.", "coherence.health.security."} {
		if strings.Contains(userArgs, "-D"+property) {
			return fmt.Errorf("direct health properties conflict with managed security")
		}
	}
	if mode == "packaged" {
		if strings.Contains(userArgs, "java.security.auth.login.config") {
			return fmt.Errorf("existing JAAS configuration requires jaas.mode=existing")
		}
		for _, name := range []string{"username", "password"} {
			content, err := os.ReadFile(filepath.Join(healthAuthDirectory, name, "value"))
			if err != nil || len(content) == 0 {
				return fmt.Errorf("health credential file unavailable")
			}
		}
		config := contextName + " { com.oracle.coherence.k8s.SecretLoginModule required usernameFile=\"" + filepath.Join(healthAuthDirectory, "username", "value") + "\" passwordFile=\"" + filepath.Join(healthAuthDirectory, "password", "value") + "\"; };\n"
		name := filepath.Join(details.UtilsDir, "health-jaas.conf")
		if err := os.WriteFile(name, []byte(config), 0600); err != nil {
			return fmt.Errorf("cannot write health JAAS configuration")
		}
		details.AddSystemPropertyArg("java.security.auth.login.config", name)
	} else if !strings.Contains(userArgs, "java.security.auth.login.config") {
		return fmt.Errorf("existing mode requires an explicit JAAS login configuration")
	}
	ssl := connection.ServerTLS
	base := ""
	if ssl.Secrets != nil && *ssl.Secrets != "" {
		base = healthTLSDirectory
	}
	file := func(value *string) (string, error) {
		path := *value
		if base != "" {
			if filepath.Base(path) != path || path == "." || path == ".." {
				return "", fmt.Errorf("invalid health store Secret key")
			}
			path = filepath.Join(base, path)
		}
		info, err := os.Stat(path)
		if err != nil || !info.Mode().IsRegular() {
			return "", fmt.Errorf("health TLS file unavailable")
		}
		return path, nil
	}
	mappings := []struct {
		value *string
		name  string
		url   bool
	}{
		{ssl.KeyStore, "keystore", true},
		{ssl.KeyStorePasswordFile, "keystore.password", false},
		{ssl.KeyPasswordFile, "key.password", false},
	}
	for _, mapping := range mappings {
		value, err := file(mapping.value)
		if err != nil {
			return err
		}
		if mapping.url {
			value = "file:" + value
		}
		details.AddSystemPropertyArg("coherence.health.security."+mapping.name, value)
	}
	for _, mapping := range []struct {
		value *string
		name  string
	}{
		{ssl.KeyStoreType, "keystore.type"},
		{ssl.KeyStoreAlgorithm, "keystore.algorithm"},
		{ssl.KeyStoreProvider, "keystore.provider"},
	} {
		if mapping.value != nil {
			details.AddSystemPropertyArg("coherence.health.security."+mapping.name, *mapping.value)
		}
	}
	details.AddSystemPropertyArg("coherence.health.http.mutators.enabled", "true")
	details.AddSystemPropertyArg("coherence.health.http.auth", "basic")
	details.AddSystemPropertyArg("coherence.health.http.provider", "HealthSSLProvider")
	return nil
}

func createHealthSecurityCheckCommand(ctx context.Context, details *run_details.RunDetails) (*exec.Cmd, error) {
	if details.IsBuildPacks() {
		return nil, fmt.Errorf("managed secure health mutators do not support the Cloud Native Buildpack launcher")
	}

	var args []string
	if details.IsSpringBoot() {
		jar, _ := details.LookupEnv(v1.EnvVarSpringBootFatJar)
		if jar == "" {
			return nil, fmt.Errorf("managed secure health mutators require a configured Spring Boot fat JAR")
		}
		classpath := strings.ReplaceAll(details.GetClasspath(), ":", ",")
		args = append(args,
			"-Dloader.path="+classpath,
			"-Dcoherence.operator.springboot.listener=false",
			"-Dloader.main="+healthSecuritySupportClass,
			v1.JvmOptClassPath, jar)
		if details.AppType == v1.AppTypeSpring2 {
			args = append(args, v1.SpringBootMain2)
		} else {
			args = append(args, v1.SpringBootMain3)
		}
	} else {
		args = append(args, v1.JvmOptClassPath, details.GetClasspath(), healthSecuritySupportClass)
	}
	cmd := exec.CommandContext(ctx, details.GetJavaExecutable(), args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd, nil
}
