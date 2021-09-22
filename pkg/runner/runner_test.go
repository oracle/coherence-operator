/*
 * Copyright (c) 2020, 2021, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package runner

import (
	"fmt"
	. "github.com/onsi/gomega"
	coh "github.com/oracle/coherence-operator/api/v1"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
	"os"
	"strings"
	"testing"
)

func TestMinimalDeployment(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
	}

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(d)

	expectedCommand := GetJavaCommand()
	expectedArgs := GetMinimalExpectedArgs()

	e, err := ExecuteWithArgs(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(expectedCommand))
	g.Expect(e.OsCmd.Args).To(ConsistOf(expectedArgs))
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

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(d)

	expectedCommand := GetJavaCommand()
	expectedArgs := append(GetMinimalExpectedArgsWithoutPrefix("-Dcoherence.override="),
		"-Dcoherence.override=k8s-coherence-override.xml")

	e, err := ExecuteWithArgs(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(expectedCommand))
	g.Expect(e.OsCmd.Args).To(ConsistOf(expectedArgs))
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
	cp := fmt.Sprintf("%s/resources:%s/classes:%s/classpath/bar2.JAR:%s/classpath/foo2.jar:%s/libs/bar1.JAR:%s/libs/foo1.jar",
		TestAppDir, TestAppDir, TestAppDir, TestAppDir, TestAppDir, TestAppDir)

	args := []string{
		GetJavaCommand(),
		"-cp",
		cp + ":/coherence-operator/utils/lib/coherence-operator.jar:/coherence-operator/utils/config",
	}

	return append(AppendCommonExpectedArgs(args),
		"com.oracle.coherence.k8s.Main",
		"com.tangosol.net.DefaultCacheServer")
}

func GetMinimalExpectedArgsWithoutAppClasspath() []string {
	args := []string{
		GetJavaCommand(),
		"-cp",
		"/coherence-operator/utils/lib/coherence-operator.jar:/coherence-operator/utils/config",
	}

	return append(AppendCommonExpectedArgs(args),
		"com.oracle.coherence.k8s.Main",
		"com.tangosol.net.DefaultCacheServer")
}

func AppendCommonExpectedArgs(args []string) []string {
	return append(args,
		"-Dcoherence.role=test",
		"-XshowSettings:all",
		"-XX:+PrintCommandLineFlags",
		"-XX:+PrintFlagsFinal",
		"-Dcoherence.wka=test-wka..svc.cluster.local",
		"-Dcoherence.cluster=test",
		//"-Dcoherence.k8s.operator.identity=test@",
		"-Dcoherence.k8s.operator.health.port=6676",
		"-Dcoherence.management.http.port=30000",
		"-Dcoherence.metrics.http.port=9612",
		"-Dcoherence.distributed.persistence-mode=on-demand",
		"-Dcoherence.override=k8s-coherence-nossl-override.xml",
		"-Dcoherence.ipmonitor.pingtimeout=0",
		"-Dcoherence.k8s.operator.diagnostics.dir=/coherence-operator/jvm/unknown/unknown",
		"-XX:HeapDumpPath=/coherence-operator/jvm/unknown/unknown/heap-dumps/unknown-unknown.hprof",
		"-Dcoherence.k8s.operator.can.resume.services=true",
		"-XX:+UseG1GC",
		"-Dcoherence.ttl=0",
		"-XX:+UnlockDiagnosticVMOptions",
		"-XX:ErrorFile=/coherence-operator/jvm/unknown/unknown/hs-err-unknown-unknown.log",
		"-XX:+UseContainerSupport",
		"-XX:+HeapDumpOnOutOfMemoryError",
		"-XX:+ExitOnOutOfMemoryError",
		"-XX:NativeMemoryTracking=summary",
		"-XX:+PrintNMTStatistics",
	)
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
	return os.Getenv("JAVA_HOME") + "/bin/java"
	//cmd := exec.Command("java")
	//return cmd.Path
}

func EnvVarsFromDeployment(d *coh.Coherence) map[string]string {
	envVars := make(map[string]string)
	envVars[coh.EnvVarCohAppDir] = TestAppDir
	envVars[coh.EnvVarCohSkipSite] = "true"

	if d.Spec.JVM == nil {
		d.Spec.JVM = &coh.JVMSpec{}
	}

	res := d.Spec.CreateStatefulSet(d)
	sts := res.Spec.(*appsv1.StatefulSet)
	c := coh.FindContainer(coh.ContainerNameCoherence, sts)
	for _, ev := range c.Env {
		if ev.ValueFrom == nil {
			envVars[ev.Name] = ev.Value
		}
	}

	return envVars
}
