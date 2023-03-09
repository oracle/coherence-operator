/*
 * Copyright (c) 2020, 2023, Oracle and/or its affiliates.
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
	"k8s.io/utils/pointer"
	"os"
	"testing"
)

func TestJibClasspath(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceResourceSpec{
			JVM: &coh.JVMSpec{
				UseJibClasspath: pointer.Bool(true),
			},
		},
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

func TestJibClasspathFile(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceResourceSpec{
			JVM: &coh.JVMSpec{
				UseJibClasspath: pointer.Bool(true),
			},
		},
	}

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(d)

	expectedCommand := GetJavaCommand()
	f := createJibClasspathFile()
	defer os.Remove(f.Name())
	expectedArgs := GetMinimalExpectedArgsWithAppClasspathFile()

	e, err := ExecuteWithArgs(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(expectedCommand))
	g.Expect(e.OsCmd.Args).To(ConsistOf(expectedArgs))
}

func TestJibMainClassFile(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceResourceSpec{
			JVM: &coh.JVMSpec{
				UseJibClasspath: pointer.Bool(true),
			},
		},
	}

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(d)

	expectedCommand := GetJavaCommand()
	f := createJibMainClassFile()
	defer os.Remove(f.Name())
	expectedArgs := GetMinimalExpectedArgsWithAppMainClassFile()

	e, err := ExecuteWithArgs(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(expectedCommand))
	g.Expect(e.OsCmd.Args).To(ConsistOf(expectedArgs))
}

func TestJibClasspathFileAndMainClassFile(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceResourceSpec{
			JVM: &coh.JVMSpec{
				UseJibClasspath: pointer.Bool(true),
			},
		},
	}

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(d)

	expectedCommand := GetJavaCommand()
	f1 := createJibClasspathFile()
	defer os.Remove(f1.Name())
	f2 := createJibMainClassFile()
	defer os.Remove(f2.Name())
	expectedArgs := GetMinimalExpectedArgsWithAppClasspathFileAndMainClassFile()

	e, err := ExecuteWithArgs(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(expectedCommand))
	g.Expect(e.OsCmd.Args).To(ConsistOf(expectedArgs))
}

func createJibClasspathFile() *os.File {
	f, err := os.Create(TestAppDir + string(os.PathSeparator) + "jib-classpath-file")
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	_, err = f.WriteString(fmt.Sprintf("%s/classpath/*:%s/libs/*", TestAppDir, TestAppDir))
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

	_, err = f.WriteString("com.tangosol.net.DefaultCacheServer")
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
	args := []string{
		GetJavaCommand(),
		"-cp",
		cp + ":/coherence-operator/utils/lib/coherence-operator.jar:/coherence-operator/utils/config",
	}

	return append(AppendCommonExpectedArgs(args),
		"com.oracle.coherence.k8s.Main",
		"$DEFAULT$")
}

func GetMinimalExpectedArgsWithAppMainClassFile() []string {
	cp := fmt.Sprintf("%s/resources:%s/classes:%s/classpath/bar2.JAR:%s/classpath/foo2.jar:%s/libs/bar1.JAR:%s/libs/foo1.jar",
		TestAppDir, TestAppDir, TestAppDir, TestAppDir, TestAppDir, TestAppDir)

	args := []string{
		GetJavaCommand(),
		"-cp",
		cp + ":/coherence-operator/utils/lib/coherence-operator.jar:/coherence-operator/utils/config",
	}

	fileName := fmt.Sprintf("%s/jib-main-class-file", TestAppDir)
	mainCls := readFirstLine(fileName)
	return append(AppendCommonExpectedArgs(args),
		"com.oracle.coherence.k8s.Main",
		mainCls)
}

func GetMinimalExpectedArgsWithAppClasspathFileAndMainClassFile() []string {
	fileName := fmt.Sprintf("%s/jib-classpath-file", TestAppDir)
	cp := readFirstLine(fileName)

	args := []string{
		GetJavaCommand(),
		"-cp",
		cp + ":/coherence-operator/utils/lib/coherence-operator.jar:/coherence-operator/utils/config",
	}

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
