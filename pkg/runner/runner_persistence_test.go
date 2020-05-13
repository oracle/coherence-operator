/*
 * Copyright (c) 2020 Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package runner

import (
	. "github.com/onsi/gomega"
	coh "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
	"testing"
)

func TestServerWithPersistenceMode(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.CoherenceDeployment{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceDeploymentSpec{
			Coherence: &coh.CoherenceSpec{
				Persistence: &coh.PersistenceSpec{
					Mode: pointer.StringPtr("active"),
				},
			},
		},
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
		"-Dcoherence.distributed.persistence-mode=active",
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

func TestServerWithPersistenceDirectory(t *testing.T) {
	g := NewGomegaWithT(t)

	args := []string{"server"}
	env := map[string]string{
		coh.EnvVarCohPersistenceDir: coh.VolumeMountPathPersistence,
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
		"-Dcoherence.distributed.persistence.base.dir=/persistence",
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

func TestServerWithSnapshotDirectory(t *testing.T) {
	g := NewGomegaWithT(t)

	args := []string{"server"}
	env := map[string]string{
		coh.EnvVarCohSnapshotDir: coh.VolumeMountPathSnapshots,
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
		"-Dcoherence.distributed.persistence.snapshot.dir=/snapshot",
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
