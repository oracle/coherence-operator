/*
 * Copyright (c) 2020, 2025, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package runner

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"testing"

	. "github.com/onsi/gomega"
	coh "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	buildDir  = "build"
	outDir    = buildDir + string(os.PathSeparator) + "_output"
	testFiles = outDir + string(os.PathSeparator) + "init-test-files"
)

func TestInitialise(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
	}

	env := EnvVarsFromDeployment(t, d)

	args, err := createArgs()
	g.Expect(err).NotTo(HaveOccurred())

	ex, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(ex).NotTo(BeNil())
	g.Expect(ex.OsCmd).To(BeNil())
}

func TestInitialiseWithCommand(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
	}

	env := EnvVarsFromDeployment(t, d)

	args, err := createArgs()
	g.Expect(err).NotTo(HaveOccurred())

	args = append(args, "--cmd", "server,--dry-run")

	ex, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(ex).NotTo(BeNil())
	g.Expect(ex.OsCmd).NotTo(BeNil())
	g.Expect(strings.HasSuffix(ex.OsCmd.Path, "java")).To(BeTrue())
}

func createArgs() ([]string, error) {
	testRoot, err := FindTestFilesRootDir()
	if err != nil {
		return nil, err
	}

	sep := string(os.PathSeparator)
	root := testRoot + sep + "root"
	utils := testRoot + sep + "utils"
	persistence := testRoot + sep + "persistence"
	snapshots := testRoot + sep + "snapshots"

	err = os.RemoveAll(utils)
	if err != nil {
		return nil, err
	}

	err = CreateInitTestFiles(root)
	if err != nil {
		return nil, err
	}

	return []string{"init", "--root", root, "--utils", utils, "--persistence", persistence, "--snapshots", snapshots}, nil
}

func CreateInitTestFiles(root string) error {
	sep := string(os.PathSeparator)

	files := root + sep + "files"
	_ = os.RemoveAll(files)

	// remove any left-over lib files from previous test
	lib := files + sep + "lib"
	if err := os.RemoveAll(lib); err != nil {
		return err
	}
	if err := os.MkdirAll(lib, os.ModePerm); err != nil {
		return err
	}

	// create some random fake files in lib
	if err := createFakeFileInDir(lib); err != nil {
		return err
	}

	logging := files + sep + "logging"
	if err := os.MkdirAll(logging, os.ModePerm); err != nil {
		return err
	}

	// create some random fake files in logging
	if err := createFakeFileInDir(logging); err != nil {
		return err
	}

	// create fake runner
	runner := files + sep + "runner"
	err := createFakeFile(runner)

	return err
}

// FindTestFilesRootDir returns the test file root directory.
func FindTestFilesRootDir() (string, error) {
	pd, err := helper.FindProjectRootDir()
	if err != nil {
		return "", err
	}

	return pd + string(os.PathSeparator) + testFiles, nil
}

func createFakeFileInDir(dir string) error {
	id := rand.Intn(1000)
	return createFakeFile(fmt.Sprintf("%s%sfile-%d.txt", dir, string(os.PathSeparator), id))
}

func createFakeFile(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer closeIgnoreError(f)

	data := rand.Intn(100)
	if _, err = fmt.Fprintf(f, "%d", data); err != nil {
		return err
	}

	return nil
}

func closeIgnoreError(f *os.File) {
	_ = f.Close()

}
