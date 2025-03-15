/*
 * Copyright (c) 2020, 2025, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package runner

import (
	. "github.com/onsi/gomega"
	coh "github.com/oracle/coherence-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	"testing"
)

func TestServerWithPersistenceMode(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Coherence: &coh.CoherenceSpec{
					Persistence: &coh.PersistenceSpec{
						Mode: ptr.To("active"),
					},
				},
			},
		},
	}

	expectedFileArgs := append(GetExpectedArgsFileContentWithoutPrefix("-Dcoherence.distributed.persistence-mode="),
		"-Dcoherence.distributed.persistence-mode=active")
	verifyConfigFilesWithArgs(t, d, expectedFileArgs)

	args := []string{"server", "--dry-run"}
	env := EnvVarsFromDeployment(t, d)

	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(GetJavaCommand()))
	expected := RemoveArgWithPrefix(GetMinimalExpectedArgs(t), "-Dcoherence.distributed.persistence-mode=")
	expected = append(expected, "-Dcoherence.distributed.persistence-mode=active")
	g.Expect(e.OsCmd.Args).To(ConsistOf(expected))
}

func TestServerWithPersistenceDirectory(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Env: []corev1.EnvVar{
					{
						Name:  coh.EnvVarCohPersistenceDir,
						Value: coh.VolumeMountPathPersistence,
					},
				},
			},
		},
	}

	verifyConfigFilesWithArgs(t, d, GetExpectedArgsFileContentWith("-Dcoherence.distributed.persistence.base.dir="+coh.VolumeMountPathPersistence))

	env := EnvVarsFromDeployment(t, d)

	args := []string{"server", "--dry-run"}
	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(GetJavaCommand()))
	g.Expect(e.OsCmd.Args).To(ConsistOf(GetMinimalExpectedArgsWith(t, "-Dcoherence.distributed.persistence.base.dir="+coh.VolumeMountPathPersistence)))
}

func TestServerWithSnapshotDirectory(t *testing.T) {
	g := NewGomegaWithT(t)

	d := &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: coh.CoherenceStatefulSetResourceSpec{
			CoherenceResourceSpec: coh.CoherenceResourceSpec{
				Env: []corev1.EnvVar{
					{
						Name:  coh.EnvVarCohSnapshotDir,
						Value: coh.VolumeMountPathSnapshots,
					},
				},
			},
		},
	}

	verifyConfigFilesWithArgs(t, d, GetExpectedArgsFileContentWith("-Dcoherence.distributed.persistence.snapshot.dir="+coh.VolumeMountPathSnapshots))

	env := EnvVarsFromDeployment(t, d)

	args := []string{"server", "--dry-run"}
	e, err := ExecuteWithArgsAndNewViper(env, args)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(e).NotTo(BeNil())
	g.Expect(e.OsCmd).NotTo(BeNil())

	g.Expect(e.OsCmd.Dir).To(Equal(TestAppDir))
	g.Expect(e.OsCmd.Path).To(Equal(GetJavaCommand()))
	g.Expect(e.OsCmd.Args).To(ConsistOf(GetMinimalExpectedArgsWith(t, "-Dcoherence.distributed.persistence.snapshot.dir="+coh.VolumeMountPathSnapshots)))
}
