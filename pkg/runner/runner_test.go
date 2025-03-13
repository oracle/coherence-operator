/*
 * Copyright (c) 2020, 2025, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package runner

import (
	"fmt"
	. "github.com/onsi/gomega"
	coh "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestMinimalDeployment(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
	}

	verifyConfigFilesWithArgs(t, d, GetExpectedArgsFileContent())

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(t, d)

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(GetJavaCommand()))
	g.Expect(e.OsCmd.Args).To(ConsistOf(GetMinimalExpectedArgs(t)))
}

func TestMinimalServerSkipCoherenceVersionCheck(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Coherence: &coh.CoherenceSpec{
					SkipVersionCheck: ptr.To(true),
				},
			},
		},
	}

	verifyConfigFilesWithArgs(t, d, GetExpectedArgsFileContent())

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(t, d)

	expectedCommand := GetJavaCommand()
	//expectedArgs := append(GetExpectedArgsFileContentWithoutPrefix("-Dcoherence.override="), "-Dcoherence.override=k8s-coherence-override.xml")

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(expectedCommand))
	g.Expect(e.OsCmd.Args).To(ConsistOf(GetMinimalExpectedArgs(t)))
}

func ReplaceArg(args []string, toReplace, replaceWith string) []string {
	for i, a := range args {
		if a == toReplace {
			args[i] = replaceWith
		}
	}
	return args
}

func RemoveArg(args []string, toRemove string) []string {
	result := make([]string, 0, len(args))
	for _, a := range args {
		if a != toRemove {
			result = append(result, a)
		}
	}
	return result
}

func GetJavaArg() string {
	var javaCmd = "java"
	javaHome, found := os.LookupEnv("JAVA_HOME")
	if found {
		javaCmd = javaHome + "/bin/java"
	}
	return javaCmd
}

func GetMinimalExpectedArgs(t *testing.T) []string {
	return GetMinimalExpectedArgsWithWorkingDir(t, TestAppDir)
}

func GetMinimalExpectedArgsWithWorkingDir(t *testing.T, wd string) []string {
	cp := ""
	if wd == TestAppDir {
		cp = fmt.Sprintf("%s/resources:%s/classes:%s/classpath/bar2.JAR:%s/classpath/foo2.jar:%s/libs/bar1.JAR:%s/libs/foo1.jar:",
			TestAppDir, TestAppDir, TestAppDir, TestAppDir, TestAppDir, TestAppDir)
	}

	utils := ensureTestUtilsDir(t)
	jar := fmt.Sprintf("%s/lib/coherence-operator.jar", utils)
	cfg := fmt.Sprintf(":%s/config", utils)
	cp = cp + jar
	if _, err := os.Stat(cfg); err == nil {
		cp = cp + ":" + cfg
	}

	args := []string{GetJavaArg(), "--class-path", cp}
	return append(AppendCommonExpectedArgs(args),
		"com.oracle.coherence.k8s.Main",
		"$DEFAULT$")
}

func GetExpectedArgsWithoutPrefix(t *testing.T, prefix ...string) []string {
	args := GetMinimalExpectedArgs(t)
	for _, p := range prefix {
		args = RemoveArgWithPrefix(args, p)
	}
	return args
}

func GetMinimalExpectedArgsWith(t *testing.T, args ...string) []string {
	return append(GetMinimalExpectedArgs(t), args...)
}

func GetMinimalExpectedArgsWithMainClass(t *testing.T, clz string) []string {
	return append(RemoveArg(GetMinimalExpectedArgs(t), coh.DefaultMain), clz)
}

func GetMinimalExpectedArgsWithoutCP() []string {
	args := []string{GetJavaArg()}
	return append(AppendCommonExpectedArgs(args),
		"com.oracle.coherence.k8s.Main",
		"$DEFAULT$")
}

func GetMinimalExpectedArgsWithoutAppClasspath() []string {
	cp := "/coherence-operator/utils/lib/coherence-operator.jar"
	if _, err := os.Stat("/coherence-operator/utils/config"); err == nil {
		cp = cp + ":/coherence-operator/utils/config"
	}
	args := []string{GetJavaArg(), "--class-path", cp}

	return append(AppendCommonExpectedArgs(args),
		"com.oracle.coherence.k8s.Main",
		"$DEFAULT$")
}

func GetExpectedClasspathWithUtilsDir(utils string) string {
	cp := fmt.Sprintf("%s/resources:%s/classes:%s/classpath/bar2.JAR:%s/classpath/foo2.jar:%s/libs/bar1.JAR:%s/libs/foo1.jar",
		TestAppDir, TestAppDir, TestAppDir, TestAppDir, TestAppDir, TestAppDir)
	return cp + ":" + GetOperatorClasspathWithUtilsDir(utils)
}

func GetOperatorClasspathWithUtilsDir(utils string) string {
	cp := utils + "/lib/coherence-operator.jar"
	cfg := utils + "/config"
	if _, err := os.Stat(cfg); err == nil {
		cp = cp + ":" + cfg
	}
	return cp
}

func GetExpectedArgsFileContent() []string {
	return GetExpectedArgsFileContentWith()
}

func GetExpectedArgsFileContentWith(args ...string) []string {
	var expected []string
	expected = AppendCommonExpectedArgs(expected)
	expected = append(expected, args...)
	return expected
}

func GetExpectedArgsFileContentWithoutPrefix(prefix string) []string {
	return RemoveArgWithPrefix(GetExpectedArgsFileContent(), prefix)
}

func AppendCommonExpectedArgs(args []string) []string {
	return append(AppendCommonExpectedNonServerArgs(args, "test"),
		"-XshowSettings:all",
		"-XX:+PrintCommandLineFlags",
		"-XX:+PrintFlagsFinal",
	)
}

func AppendCommonExpectedNonServerArgs(args []string, role string) []string {
	if role != "" {
		args = append(args, "-Dcoherence.role="+role)
	}
	return append(args,
		"-Dcoherence.wka=test-wka..svc",
		"-Dcoherence.cluster=test",
		"-Dcoherence.operator.health.port=6676",
		"-Dcoherence.health.http.port=6676",
		"-Dcoherence.operator.health.enabled=false",
		"-Dcoherence.management.http.port=30000",
		"-Dcoherence.metrics.http.port=9612",
		"-Dcoherence.distributed.persistence-mode=on-demand",
		"-Dcoherence.override=k8s-coherence-override.xml",
		"-Dcoherence.ipmonitor.pingtimeout=0",
		"-Dcoherence.operator.diagnostics.dir=/coherence-operator/jvm/unknown/unknown",
		"-XX:HeapDumpPath=/coherence-operator/jvm/unknown/unknown/heap-dumps/unknown-unknown.hprof",
		"-Dcoherence.operator.can.resume.services=true",
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
	javaHome, found := os.LookupEnv("JAVA_HOME")
	if found {
		return javaHome + "/bin/java"
	}
	path, _ := exec.LookPath("java")
	if path != "" {
		return path
	}
	return "java"
}

func EnvVarsFromDeployment(t *testing.T, d *coh.Coherence) map[string]string {
	return EnvVarsForContainerWithSkipSite(t, d, coh.ContainerNameCoherence, true)
}

func EnvVarsFromDeploymentWithSkipSite(t *testing.T, d *coh.Coherence, skipSite bool) map[string]string {
	return EnvVarsForContainerWithSkipSite(t, d, coh.ContainerNameCoherence, skipSite)
}

func EnvVarsForConfigContainerWithSkipSite(t *testing.T, d *coh.Coherence, skipSite bool) map[string]string {
	return EnvVarsForContainerWithSkipSite(t, d, coh.ContainerNameOperatorConfig, skipSite)
}

func EnvVarsForContainerWithSkipSite(t *testing.T, d *coh.Coherence, containerName string, skipSite bool) map[string]string {
	envVars := make(map[string]string)

	if d.Spec.JVM == nil {
		d.Spec.JVM = &coh.JVMSpec{}
	}

	res := d.Spec.CreateStatefulSetResource(d)
	sts := res.Spec.(*appsv1.StatefulSet)
	c := coh.FindContainer(containerName, sts)
	if c == nil {
		c = coh.FindInitContainer(containerName, sts)
	}
	if c == nil {
		return nil
	}

	for _, ev := range c.Env {
		if ev.ValueFrom == nil {
			envVars[ev.Name] = ev.Value
		}
	}

	if d.Spec.Application != nil && d.Spec.Application.WorkingDir != nil && *d.Spec.Application.WorkingDir != "" {
		envVars[coh.EnvVarCohAppDir] = *d.Spec.Application.WorkingDir
	} else {
		envVars[coh.EnvVarCohAppDir] = TestAppDir
	}

	dir := ensureTestUtilsDir(t)
	envVars[coh.EnvVarCohUtilDir] = dir
	envVars[coh.EnvVarCohCtlHome] = dir
	envVars[coh.EnvVarCohSkipSite] = fmt.Sprintf("%t", skipSite)

	return envVars
}

func ensureTestUtilsDir(t *testing.T) string {
	g := NewGomegaWithT(t)
	dir, err := helper.EnsureLogsDir(t.Name())
	g.Expect(err).NotTo(HaveOccurred())
	return dir
}

func verifyConfigFilesWithArgs(t *testing.T, d *coh.Coherence, expectedArgs []string) {
	verifyConfigFilesWithArgsWithSkipSite(t, d, expectedArgs, true)
}

func verifyConfigFilesWithArgsWithSkipSite(t *testing.T, d *coh.Coherence, expectedArgs []string, skipSite bool) {
	dir := ensureTestUtilsDir(t)
	expectedCP := GetExpectedClasspathWithUtilsDir(dir)
	verifyConfigFilesWithArgsAndClasspathWithSkipSite(t, d, expectedArgs, expectedCP, skipSite)
}

func verifyConfigFilesWithArgsAndClasspath(t *testing.T, d *coh.Coherence, expectedArgs []string, expectedCP string) {
	verifyConfigFilesWithArgsAndClasspathWithSkipSite(t, d, expectedArgs, expectedCP, true)
}

func verifyConfigFilesWithArgsAndClasspathWithSkipSite(t *testing.T, d *coh.Coherence, expectedArgs []string, expectedCP string, skipSite bool) {
	cfgEnv := EnvVarsForConfigContainerWithSkipSite(t, d, skipSite)
	verifyConfigFilesWithArgsAndClasspathUsingEnv(t, cfgEnv, expectedArgs, expectedCP)
}

func verifyConfigFilesWithArgsAndClasspathUsingEnv(t *testing.T, cfgEnv map[string]string, expectedArgs []string, expectedCP string) {
	var err error

	g := NewGomegaWithT(t)
	dir := ensureTestUtilsDir(t)
	cfgEnv[coh.EnvVarCohUtilDir] = dir
	cfgEnv[coh.EnvVarCohCtlHome] = dir

	_, err = ExecuteWithArgsAndNewViper(cfgEnv, []string{coh.RunnerConfig})
	g.Expect(err).NotTo(HaveOccurred())

	cpName := fmt.Sprintf("%s/%s", dir, coh.OperatorClasspathFile)
	_, err = os.Stat(cpName)
	g.Expect(err).NotTo(HaveOccurred())
	dataCP, err := os.ReadFile(cpName)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(dataCP).NotTo(BeNil())
	cp := string(dataCP)
	g.Expect(cp).To(Equal(expectedCP))

	argsName := fmt.Sprintf("%s/%s", dir, coh.OperatorJvmArgsFile)
	_, err = os.Stat(argsName)
	g.Expect(err).NotTo(HaveOccurred())
	dataArgs, err := os.ReadFile(argsName)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(dataArgs).NotTo(BeNil())

	args := filterNonEmptyStringArray(strings.Split(string(dataArgs), "\n"))

	g.Expect(args).To(ConsistOf(expectedArgs))
}

func filterNonEmptyStringArray(ss []string) (ret []string) {
	test := func(s string) bool { return s != "" }
	return filterStringArray(ss, test)
}

func filterStringArray(ss []string, test func(string) bool) (ret []string) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}
