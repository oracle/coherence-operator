/*
 * Copyright (c) 2020, 2021 Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package runner

import (
	. "github.com/onsi/gomega"
	coh "github.com/oracle/coherence-operator/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
	"testing"
)

func TestServerWithPersistenceMode(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceResourceSpec{
			Coherence: &coh.CoherenceSpec{
				Persistence: &coh.PersistenceSpec{
					Mode: pointer.StringPtr("active"),
				},
			},
		},
	}

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(d)

	expectedCommand := GetJavaCommand()
	expectedArgs := append(GetMinimalExpectedArgsWithoutPrefix("-Dcoherence.distributed.persistence-mode="),
		"-Dcoherence.distributed.persistence-mode=active")

	e, err := executeWithArgs(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(expectedCommand))
	g.Expect(e.OsCmd.Args).To(ConsistOf(expectedArgs))
}

func TestServerWithPersistenceDirectory(t *testing.T) {
	g := NewGomegaWithT(t)

	args := []string{"server", "--dry-run"}
	env := map[string]string{
		coh.EnvVarCohPersistenceDir: coh.VolumeMountPathPersistence,
	}

	expectedCommand := GetJavaCommand()

	e, err := executeWithArgs(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(expectedCommand))
	g.Expect(e.OsCmd.Args).To(ContainElement("-Dcoherence.distributed.persistence.base.dir=" + coh.VolumeMountPathPersistence))
}

func TestServerWithSnapshotDirectory(t *testing.T) {
	g := NewGomegaWithT(t)

	args := []string{"server", "--dry-run"}
	env := map[string]string{
		coh.EnvVarCohSnapshotDir: coh.VolumeMountPathSnapshots,
	}

	expectedCommand := GetJavaCommand()

	e, err := executeWithArgs(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(expectedCommand))
	g.Expect(e.OsCmd.Args).To(ContainElement("-Dcoherence.distributed.persistence.snapshot.dir=" + coh.VolumeMountPathSnapshots))
}
