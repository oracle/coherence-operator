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
	"strings"
	"testing"
)

func TestMinimalDeployment(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
	}

	args := []string{"runner", "server"}
	env := EnvVarsFromDeployment(d)

	expectedCommand := GetJavaCommand()
	expectedArgs := GetMinimalExpectedArgs()

	_, cmd, err := DryRun(args, env)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(cmd).NotTo(BeNil())

	g.Expect(cmd.Dir).To(Equal(""))
	g.Expect(cmd.Path).To(Equal(expectedCommand))
	g.Expect(cmd.Args).To(Equal(expectedArgs))
}

func TestMinimalServerSkipCoherenceVersionCheck(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceResourceSpec{
			Coherence: &coh.CoherenceSpec{
				SkipVersionCheck: pointer.BoolPtr(true),
			},
		},
	}

	args := []string{"runner", "server"}
	env := EnvVarsFromDeployment(d)

	expectedCommand := GetJavaCommand()
	expectedArgs := append(GetMinimalExpectedArgsWithoutPrefix("-Dcoherence.override="),
		"-Dcoherence.override=k8s-coherence-override.xml")

	_, cmd, err := DryRun(args, env)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(cmd).NotTo(BeNil())

	g.Expect(cmd.Dir).To(Equal(""))
	g.Expect(cmd.Path).To(Equal(expectedCommand))
	g.Expect(cmd.Args).To(ConsistOf(expectedArgs))
}

func GetMinimalExpectedArgsWithoutPrefix(prefix string) []string {
	return RemoveArgWithPrefix(GetMinimalExpectedArgs(), prefix)
}

func ReplaceArg(args []string, toReplace, replaceWith string) []string {
	for i, a := range args {
		if a == toReplace {
			args[i] = replaceWith
		}
	}
	return args
}

func GetMinimalExpectedArgs() []string {
	return []string{
		"java",
		"-cp",
		"/utils/lib/coherence-utils.jar:/app/libs/*:/app/classes:/app/resources:/utils/scripts",
		"-Dcoherence.role=test",
		"-XshowSettings:all",
		"-XX:+PrintCommandLineFlags",
		"-XX:+PrintFlagsFinal",
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
		"com.oracle.coherence.k8s.Main",
		"com.tangosol.net.DefaultCacheServer",
	}
}

func RemoveArgWithPrefix(args []string, prefix string) []string {
	result := args
	found := true

	for found {
		found = false
		for i, a := range result {
			if strings.HasPrefix(a, prefix) {
				result[i] = result[len(result)-1] // Copy last element to index i.
				result[len(result)-1] = ""        // Erase last element (write zero value).
				result = result[:len(result)-1]   // Truncate slice.
				found = true
				break
			}
		}
	}
	return result
}

func GetJavaCommand() string {
	cmd := exec.Command("java")
	return cmd.Path
}

func EnvVarsFromDeployment(d *coh.Coherence) map[string]string {
	envVars := make(map[string]string)

	if d.Spec.JVM == nil {
		d.Spec.JVM = &coh.JVMSpec{}
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
