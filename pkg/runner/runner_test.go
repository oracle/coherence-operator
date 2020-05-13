/*
 * Copyright (c) 2020 Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package runner

import (
	. "github.com/onsi/gomega"
	coh "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	"github.com/oracle/coherence-operator/pkg/flags"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
	"os/exec"
	"testing"
)

func TestMinimalServer(t *testing.T) {
	g := NewGomegaWithT(t)

	args := []string{"server"}
	env := map[string]string{}

	expectedCommand := GetJavaCommand()
	expectedArgs := []string{
		"java",
		"-cp",
		"/lib/coherence-utils.jar:/app/libs/*:/app/classes:/app/resources:/scripts",
		"-Dcoherence.cacheconfig=coherence-cache-config.xml",
		"-Dcoherence.health.port=6676",
		"-Dcoherence.management.http.port=30000",
		"-Dcoherence.metrics.http.port=9612",
		"-Dcoherence.distributed.persistence-mode=on-demand",
		"-Dcoherence.override=k8s-coherence-nossl-override.xml",
		"-XX:HeapDumpPath=/jvm/unknown/unknown/heap-dumps/unknown-unknown.hprof",
		"-XX:+UseG1GC",
		"-Dcoherence.ttl=0",
		"-XX:+UnlockDiagnosticVMOptions",
		"-XX:+UnlockExperimentalVMOptions",
		"-XX:ErrorFile=/jvm/unknown/unknown/hs-err-unknown-unknown.log",
		"com.tangosol.net.DefaultCacheServer",
	}

	_, cmd, err := DryRun(args, env)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(cmd).NotTo(BeNil())

	g.Expect(cmd.Dir).To(Equal(""))
	g.Expect(cmd.Path).To(Equal(expectedCommand))
	g.Expect(cmd.Args).To(Equal(expectedArgs))
}

func TestMinimalDeployment(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.CoherenceDeployment{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
	}

	args := []string{"server"}
	env := EnvVarsFromDeployment(d)

	expectedCommand := GetJavaCommand()
	expectedArgs := []string{
		"java",
		"-cp",
		"/utils/lib/coherence-utils.jar:/app/libs/*:/app/classes:/app/resources:/utils/scripts",
		"-Dcoherence.wka=test-wka",
		"-Dcoherence.cluster=test",
		"-Dcoherence.cacheconfig=coherence-cache-config.xml",
		"-Dcoherence.health.port=6676",
		"-Dcoherence.management.http.port=30000",
		"-Dcoherence.metrics.http.port=9612",
		"-Dcoherence.distributed.persistence-mode=on-demand",
		"-Dcoherence.override=k8s-coherence-nossl-override.xml",
		"-XX:HeapDumpPath=/jvm/unknown/unknown/heap-dumps/unknown-unknown.hprof",
		"-XX:+UseG1GC",
		"-Dcoherence.ttl=0",
		"-XX:+UnlockDiagnosticVMOptions",
		"-XX:+UnlockExperimentalVMOptions",
		"-XX:ErrorFile=/jvm/unknown/unknown/hs-err-unknown-unknown.log",
		"-XX:+UseContainerSupport",
		"com.tangosol.net.DefaultCacheServer",
	}

	_, cmd, err := DryRun(args, env)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(cmd).NotTo(BeNil())

	g.Expect(cmd.Dir).To(Equal(""))
	g.Expect(cmd.Path).To(Equal(expectedCommand))
	g.Expect(cmd.Args).To(Equal(expectedArgs))
}

func TestMinimalServerSkipCoherenceVersionCheck(t *testing.T) {
	g := NewGomegaWithT(t)

	args := []string{"server"}
	env := map[string]string{
		coh.EnvVarCohSkipVersionCheck: "true",
	}

	expectedCommand := GetJavaCommand()
	expectedArgs := []string{
		"java",
		"-cp",
		"/lib/coherence-utils.jar:/app/libs/*:/app/classes:/app/resources:/scripts",
		"-Dcoherence.cacheconfig=coherence-cache-config.xml",
		"-Dcoherence.health.port=6676",
		"-Dcoherence.management.http.port=30000",
		"-Dcoherence.metrics.http.port=9612",
		"-Dcoherence.distributed.persistence-mode=on-demand",
		"-Dcoherence.override=k8s-coherence-override.xml",
		"-XX:HeapDumpPath=/jvm/unknown/unknown/heap-dumps/unknown-unknown.hprof",
		"-XX:+UseG1GC",
		"-Dcoherence.ttl=0",
		"-XX:+UnlockDiagnosticVMOptions",
		"-XX:+UnlockExperimentalVMOptions",
		"-XX:ErrorFile=/jvm/unknown/unknown/hs-err-unknown-unknown.log",
		"com.tangosol.net.DefaultCacheServer",
	}

	_, cmd, err := DryRun(args, env)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(cmd).NotTo(BeNil())

	g.Expect(cmd.Dir).To(Equal(""))
	g.Expect(cmd.Path).To(Equal(expectedCommand))
	g.Expect(cmd.Args).To(Equal(expectedArgs))
}

func GetJavaCommand() string {
	cmd := exec.Command("java")
	return cmd.Path
}

func EnvVarsFromDeployment(d *coh.CoherenceDeployment) map[string]string {
	envVars := make(map[string]string)

	if d.Spec.JVM == nil {
		d.Spec.JVM = &coh.JVMSpec{}
	}

	if d.Spec.JVM.FlightRecorder == nil {
		d.Spec.JVM.FlightRecorder = pointer.BoolPtr(false)
	}

	opFlags := &flags.CoherenceOperatorFlags{}
	res := d.Spec.CreateStatefulSet(d, opFlags)
	sts := res.Spec.(*appsv1.StatefulSet)
	c := coh.FindContainer(coh.ContainerNameCoherence, sts)
	for _, ev := range c.Env {
		if ev.ValueFrom == nil {
			envVars[ev.Name] = ev.Value
		}
	}

	return envVars
}
