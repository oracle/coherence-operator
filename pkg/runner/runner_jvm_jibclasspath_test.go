/*
 * Copyright (c) 2020, 2025, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package runner

import (
	"bufio"
	"fmt"
	. "github.com/onsi/gomega"
	coh "github.com/oracle/coherence-operator/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	"os"
	"testing"
)

const (
	JibClassPath = "jib/classpath/*:jib/libs/*"
	JibMainClass = "com.test.JibMainClass"
)

func TestJibClasspathWhenNoJibFilePresent(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				JVM: &coh.JVMSpec{
					UseJibClasspath: ptr.To(true),
				},
			},
		},
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

func TestJibClasspathFileWhenJibFilePresent(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				JVM: &coh.JVMSpec{
					UseJibClasspath: ptr.To(true),
				},
			},
		},
	}

	f := createJibClasspathFile()
	defer os.Remove(f.Name())

	expectedCp := JibClassPath + ":" + GetOperatorClasspathWithUtilsDir(ensureTestUtilsDir(t))
	verifyConfigFilesWithArgsAndClasspath(t, d, GetExpectedArgsFileContent(), expectedCp)

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(t, d)

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(GetJavaCommand()))

	expected := append(GetMinimalExpectedArgsWithoutCP(), coh.JvmOptClassPath, expectedCp)
	g.Expect(e.OsCmd.Args).To(ConsistOf(expected))
}

func TestJibMainClassFile(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				JVM: &coh.JVMSpec{
					UseJibClasspath: ptr.To(true),
				},
			},
		},
	}

	f := createJibMainClassFile()
	defer os.Remove(f.Name())

	verifyConfigFilesWithArgs(t, d, GetExpectedArgsFileContent())

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(t, d)

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(GetJavaCommand()))
	g.Expect(e.OsCmd.Args).To(ConsistOf(GetMinimalExpectedArgsWithMainClass(t, JibMainClass)))
}

func TestJibClasspathFileAndMainClassFile(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				JVM: &coh.JVMSpec{
					UseJibClasspath: ptr.To(true),
				},
			},
		},
	}

	f1 := createJibClasspathFile()
	defer os.Remove(f1.Name())
	f2 := createJibMainClassFile()
	defer os.Remove(f2.Name())

	expectedCp := JibClassPath + ":" + GetOperatorClasspathWithUtilsDir(ensureTestUtilsDir(t))
	verifyConfigFilesWithArgsAndClasspath(t, d, GetExpectedArgsFileContent(), expectedCp)

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(t, d)

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(GetJavaCommand()))
	expected := append(GetMinimalExpectedArgsWithoutCP(), coh.JvmOptClassPath, expectedCp)
	expected = ReplaceArg(expected, coh.DefaultMain, JibMainClass)
	g.Expect(e.OsCmd.Args).To(ConsistOf(expected))
	//	g.Expect(e.OsCmd.Args).To(ConsistOf(GetMinimalExpectedArgsWithMainClass(t, "com.tangosol.net.DefaultCacheServer")))
}

func createJibClasspathFile() *os.File {
	f, err := os.Create(TestAppDir + string(os.PathSeparator) + "jib-classpath-file")
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	_, err = f.WriteString(JibClassPath)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	err = f.Close()
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	return f
}

func createJibMainClassFile() *os.File {
	f, err := os.Create(TestAppDir + string(os.PathSeparator) + "jib-main-class-file")
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	_, err = f.WriteString(JibMainClass)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	err = f.Close()
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	return f
}

func GetMinimalExpectedArgsWithAppClasspathFile() []string {
	fileName := fmt.Sprintf("%s/jib-classpath-file", TestAppDir)
	cp := readFirstLine(fileName)
	cp += ":/coherence-operator/utils/lib/coherence-operator.jar"
	if _, err := os.Stat("/coherence-operator/utils/config"); err == nil {
		cp += ":/coherence-operator/utils/config"
	}

	args := []string{GetJavaArg(), "--class-path", cp}

	return append(AppendCommonExpectedArgs(args),
		"com.oracle.coherence.k8s.Main",
		"$DEFAULT$")
}

func GetMinimalExpectedArgsWithAppMainClassFile() []string {
	cp := fmt.Sprintf("%s/resources:%s/classes:%s/classpath/bar2.JAR:%s/classpath/foo2.jar:%s/libs/bar1.JAR:%s/libs/foo1.jar",
		TestAppDir, TestAppDir, TestAppDir, TestAppDir, TestAppDir, TestAppDir)

	cp += ":/coherence-operator/utils/lib/coherence-operator.jar"
	if _, err := os.Stat("/coherence-operator/utils/config"); err == nil {
		cp += ":/coherence-operator/utils/config"
	}

	args := []string{GetJavaArg(), "--class-path", cp}

	fileName := fmt.Sprintf("%s/jib-main-class-file", TestAppDir)
	mainCls := readFirstLine(fileName)
	return append(AppendCommonExpectedArgs(args),
		"com.oracle.coherence.k8s.Main",
		mainCls)
}

func GetMinimalExpectedArgsWithAppClasspathFileAndMainClassFile() []string {
	fileName := fmt.Sprintf("%s/jib-classpath-file", TestAppDir)
	cp := readFirstLine(fileName)
	cp += ":/coherence-operator/utils/lib/coherence-operator.jar"
	if _, err := os.Stat("/coherence-operator/utils/config"); err == nil {
		cp += ":/coherence-operator/utils/config"
	}

	args := []string{GetJavaArg(), "--class-path", cp}

	fileName = fmt.Sprintf("%s/jib-main-class-file", TestAppDir)
	mainCls := readFirstLine(fileName)
	return append(AppendCommonExpectedArgs(args),
		"com.oracle.coherence.k8s.Main",
		mainCls)
}

func readFirstLine(fqfn string) string {
	file, _ := os.Open(fqfn)
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var text []string
	for scanner.Scan() {
		text = append(text, scanner.Text())
	}
	file.Close()
	if len(text) == 0 {
		return ""
	}
	return text[0]
}
