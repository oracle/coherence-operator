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

	args := []string{"runner", "server"}
	env := EnvVarsFromDeployment(d)

	expectedCommand := GetJavaCommand()
	expectedArgs := append(GetMinimalExpectedArgsWithoutPrefix("-Dcoherence.distributed.persistence-mode="),
		"-Dcoherence.distributed.persistence-mode=active")

	_, cmd, err := DryRun(args, env)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(cmd).NotTo(BeNil())

	g.Expect(cmd.Dir).To(Equal(""))
	g.Expect(cmd.Path).To(Equal(expectedCommand))
	g.Expect(cmd.Args).To(ConsistOf(expectedArgs))
}

func TestServerWithPersistenceDirectory(t *testing.T) {
	g := NewGomegaWithT(t)

	args := []string{"server"}
	env := map[string]string{
		coh.EnvVarCohPersistenceDir: coh.VolumeMountPathPersistence,
	}

	expectedCommand := GetJavaCommand()

	_, cmd, err := DryRun(args, env)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(cmd).NotTo(BeNil())

	g.Expect(cmd.Dir).To(Equal(""))
	g.Expect(cmd.Path).To(Equal(expectedCommand))
	g.Expect(cmd.Args).To(ContainElement("-Dcoherence.distributed.persistence.base.dir=/persistence"))
}

func TestServerWithSnapshotDirectory(t *testing.T) {
	g := NewGomegaWithT(t)

	args := []string{"server"}
	env := map[string]string{
		coh.EnvVarCohSnapshotDir: coh.VolumeMountPathSnapshots,
	}

	expectedCommand := GetJavaCommand()

	_, cmd, err := DryRun(args, env)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(cmd).NotTo(BeNil())

	g.Expect(cmd.Dir).To(Equal(""))
	g.Expect(cmd.Path).To(Equal(expectedCommand))
	g.Expect(cmd.Args).To(ContainElement("-Dcoherence.distributed.persistence.snapshot.dir=/snapshot"))
}
