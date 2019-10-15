/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package helm_test

import (
	. "github.com/onsi/gomega"
	cohv1 "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	. "github.com/oracle/coherence-operator/test/e2e/helper/matchers"
	corev1 "k8s.io/api/core/v1"
	"testing"
)

func TestSnapshotWhenNotSetInMinimalCluster(t *testing.T) {
	g := NewGomegaWithT(t)

	// Use the specified yaml files to create a CoherenceCluster and trigger a fake end-to-end
	// reconcile to obtain the resources that would have been created by the Helm operator.
	result, cluster, err := CreateCluster("minimal-cluster.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	// Obtain the StatefulSet that Helm would have created
	sts, err := findStatefulSet(result, cluster, cohv1.DefaultRoleName)
	g.Expect(err).NotTo(HaveOccurred())

	// Obtain the Coherence container from the StatefulSet
	container, err := findContainer(sts, coherenceContainer)
	g.Expect(err).NotTo(HaveOccurred())

	// Assert that the COH_SNAPSHOT_ENABLED env-var is not set
	g.Expect(container.Env).NotTo(HaveEnvVarNamed("COH_SNAPSHOT_ENABLED"))
}

func TestSnapshotWhenNotSetInClusterWithExplicitRoles(t *testing.T) {
	g := NewGomegaWithT(t)

	result, cluster, err := CreateCluster("minimal-explicit-roles.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	for name := range cluster.GetRoles() {
		sts, err := findStatefulSet(result, cluster, name)
		g.Expect(err).NotTo(HaveOccurred())

		container, err := findContainer(sts, coherenceContainer)
		g.Expect(err).NotTo(HaveOccurred())

		g.Expect(container.Env).NotTo(HaveEnvVarNamed("COH_SNAPSHOT_ENABLED"))
	}
}

func TestSnapshotWithPVCWhenSetForImplicitRole(t *testing.T) {
	g := NewGomegaWithT(t)

	result, cluster, err := CreateCluster("snapshot-test-pvc-implicit-role.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	sts, err := findStatefulSet(result, cluster, cohv1.DefaultRoleName)
	g.Expect(err).NotTo(HaveOccurred())

	container, err := findContainer(sts, coherenceContainer)
	g.Expect(err).NotTo(HaveOccurred())

	g.Expect(container.Env).To(HaveEnvVar(corev1.EnvVar{Name: "COH_SNAPSHOT_ENABLED", Value: "true"}))

	mount, err := findVolumeMount(container, "snapshot-volume")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(mount.MountPath).To(Equal("/snapshot"))

	_, err = findVolume(sts, "snapshot-volume")
	g.Expect(err).To(HaveOccurred())

	pvc, err := findPersistentVolumeClaim(sts, "snapshot-volume")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(*pvc.Spec.StorageClassName).To(Equal("foo"))
}

func TestSnapshotWithPVCWhenDefaultSetInClusterWithExplicitRoles(t *testing.T) {
	g := NewGomegaWithT(t)

	result, cluster, err := CreateCluster("snapshot-test-pvc-explicit-roles-with-default.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	stsData, err := findStatefulSet(result, cluster, "data")
	g.Expect(err).NotTo(HaveOccurred())

	stsProxy, err := findStatefulSet(result, cluster, "proxy")
	g.Expect(err).NotTo(HaveOccurred())

	dataContainer, err := findContainerForRole(result, cluster, "data", coherenceContainer)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(dataContainer.Env).To(HaveEnvVar(corev1.EnvVar{Name: "COH_SNAPSHOT_ENABLED", Value: "true"}))

	dataMount, err := findVolumeMount(dataContainer, "snapshot-volume")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(dataMount.MountPath).To(Equal("/snapshot"))

	_, err = findVolume(stsData, "snapshot-volume")
	g.Expect(err).To(HaveOccurred())

	dataPVC, err := findPersistentVolumeClaim(stsData, "snapshot-volume")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(*dataPVC.Spec.StorageClassName).To(Equal("foo"))

	proxyContainer, err := findContainerForRole(result, cluster, "proxy", coherenceContainer)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(proxyContainer.Env).To(HaveEnvVar(corev1.EnvVar{Name: "COH_SNAPSHOT_ENABLED", Value: "true"}))

	proxyMount, err := findVolumeMount(dataContainer, "snapshot-volume")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(proxyMount.MountPath).To(Equal("/snapshot"))

	_, err = findVolume(stsProxy, "snapshot-volume")
	g.Expect(err).To(HaveOccurred())

	proxyPVC, err := findPersistentVolumeClaim(stsProxy, "snapshot-volume")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(*proxyPVC.Spec.StorageClassName).To(Equal("foo"))
}

func TestSnapshotWithPVCWhenDefaultSetAndOverriddenInClusterWithExplicitRoles(t *testing.T) {
	g := NewGomegaWithT(t)

	result, cluster, err := CreateCluster("snapshot-test-pvc-explicit-roles-override-default.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	stsData, err := findStatefulSet(result, cluster, "data")
	g.Expect(err).NotTo(HaveOccurred())

	dataContainer, err := findContainerForRole(result, cluster, "data", coherenceContainer)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(dataContainer.Env).To(HaveEnvVar(corev1.EnvVar{Name: "COH_SNAPSHOT_ENABLED", Value: "true"}))

	dataMount, err := findVolumeMount(dataContainer, "snapshot-volume")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(dataMount.MountPath).To(Equal("/snapshot"))

	_, err = findVolume(stsData, "snapshot-volume")
	g.Expect(err).To(HaveOccurred())

	dataPVC, err := findPersistentVolumeClaim(stsData, "snapshot-volume")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(*dataPVC.Spec.StorageClassName).To(Equal("foo"))

	proxyContainer, err := findContainerForRole(result, cluster, "proxy", coherenceContainer)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(proxyContainer.Env).NotTo(HaveEnvVarNamed("COH_SNAPSHOT_ENABLED"))
}

func TestSnapshotWithPVCWhenSetInClusterWithExplicitRoles(t *testing.T) {
	g := NewGomegaWithT(t)

	result, cluster, err := CreateCluster("snapshot-test-pvc-explicit-roles.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	stsData, err := findStatefulSet(result, cluster, "data")
	g.Expect(err).NotTo(HaveOccurred())

	dataContainer, err := findContainerForRole(result, cluster, "data", coherenceContainer)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(dataContainer.Env).To(HaveEnvVar(corev1.EnvVar{Name: "COH_SNAPSHOT_ENABLED", Value: "true"}))

	dataMount, err := findVolumeMount(dataContainer, "snapshot-volume")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(dataMount.MountPath).To(Equal("/snapshot"))

	_, err = findVolume(stsData, "snapshot-volume")
	g.Expect(err).To(HaveOccurred())

	dataPVC, err := findPersistentVolumeClaim(stsData, "snapshot-volume")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(*dataPVC.Spec.StorageClassName).To(Equal("foo"))

	proxyContainer, err := findContainerForRole(result, cluster, "proxy", coherenceContainer)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(proxyContainer.Env).NotTo(HaveEnvVarNamed("COH_SNAPSHOT_ENABLED"))
}

func TestSnapshotWithVolumeWhenSetForImplicitRole(t *testing.T) {
	g := NewGomegaWithT(t)

	result, cluster, err := CreateCluster("snapshot-test-volume-implicit-role.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	sts, err := findStatefulSet(result, cluster, cohv1.DefaultRoleName)
	g.Expect(err).NotTo(HaveOccurred())

	container, err := findContainer(sts, coherenceContainer)
	g.Expect(err).NotTo(HaveOccurred())

	g.Expect(container.Env).To(HaveEnvVar(corev1.EnvVar{Name: "COH_SNAPSHOT_ENABLED", Value: "true"}))

	mount, err := findVolumeMount(container, "snapshot-volume")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(mount.MountPath).To(Equal("/snapshot"))

	vol, err := findVolume(sts, "snapshot-volume")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(vol.HostPath).NotTo(BeNil())
	g.Expect(vol.HostPath.Path).To(Equal("/data"))

	_, err = findPersistentVolumeClaim(sts, "snapshot-volume")
	g.Expect(err).To(HaveOccurred())
}

func TestSnapshotWithVolumeWhenDefaultSetInClusterWithExplicitRoles(t *testing.T) {
	g := NewGomegaWithT(t)

	result, cluster, err := CreateCluster("snapshot-test-volume-explicit-roles-with-default.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	stsData, err := findStatefulSet(result, cluster, "data")
	g.Expect(err).NotTo(HaveOccurred())

	stsProxy, err := findStatefulSet(result, cluster, "proxy")
	g.Expect(err).NotTo(HaveOccurred())

	dataContainer, err := findContainerForRole(result, cluster, "data", coherenceContainer)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(dataContainer.Env).To(HaveEnvVar(corev1.EnvVar{Name: "COH_SNAPSHOT_ENABLED", Value: "true"}))

	dataMount, err := findVolumeMount(dataContainer, "snapshot-volume")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(dataMount.MountPath).To(Equal("/snapshot"))

	dataVol, err := findVolume(stsData, "snapshot-volume")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(dataVol.HostPath).NotTo(BeNil())
	g.Expect(dataVol.HostPath.Path).To(Equal("/data"))

	_, err = findPersistentVolumeClaim(stsData, "snapshot-volume")
	g.Expect(err).To(HaveOccurred())

	proxyContainer, err := findContainerForRole(result, cluster, "proxy", coherenceContainer)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(proxyContainer.Env).To(HaveEnvVar(corev1.EnvVar{Name: "COH_SNAPSHOT_ENABLED", Value: "true"}))

	proxyMount, err := findVolumeMount(dataContainer, "snapshot-volume")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(proxyMount.MountPath).To(Equal("/snapshot"))

	proxyVol, err := findVolume(stsProxy, "snapshot-volume")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(proxyVol.HostPath).NotTo(BeNil())
	g.Expect(proxyVol.HostPath.Path).To(Equal("/data"))

	_, err = findPersistentVolumeClaim(stsProxy, "snapshot-volume")
	g.Expect(err).To(HaveOccurred())
}

func TestSnapshotWithVolumeWhenDefaultSetAndOverriddenInClusterWithExplicitRoles(t *testing.T) {
	g := NewGomegaWithT(t)

	result, cluster, err := CreateCluster("snapshot-test-volume-explicit-roles-override-default.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	stsData, err := findStatefulSet(result, cluster, "data")
	g.Expect(err).NotTo(HaveOccurred())

	dataContainer, err := findContainerForRole(result, cluster, "data", coherenceContainer)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(dataContainer.Env).To(HaveEnvVar(corev1.EnvVar{Name: "COH_SNAPSHOT_ENABLED", Value: "true"}))

	dataMount, err := findVolumeMount(dataContainer, "snapshot-volume")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(dataMount.MountPath).To(Equal("/snapshot"))

	dataVol, err := findVolume(stsData, "snapshot-volume")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(dataVol.HostPath).NotTo(BeNil())
	g.Expect(dataVol.HostPath.Path).To(Equal("/data"))

	_, err = findPersistentVolumeClaim(stsData, "snapshot-volume")
	g.Expect(err).To(HaveOccurred())

	proxyContainer, err := findContainerForRole(result, cluster, "proxy", coherenceContainer)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(proxyContainer.Env).NotTo(HaveEnvVarNamed("COH_SNAPSHOT_ENABLED"))
}

func TestSnapshotWithVolumeWhenSetInClusterWithExplicitRoles(t *testing.T) {
	g := NewGomegaWithT(t)

	result, cluster, err := CreateCluster("snapshot-test-volume-explicit-roles.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	stsData, err := findStatefulSet(result, cluster, "data")
	g.Expect(err).NotTo(HaveOccurred())

	dataContainer, err := findContainerForRole(result, cluster, "data", coherenceContainer)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(dataContainer.Env).To(HaveEnvVar(corev1.EnvVar{Name: "COH_SNAPSHOT_ENABLED", Value: "true"}))

	dataMount, err := findVolumeMount(dataContainer, "snapshot-volume")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(dataMount.MountPath).To(Equal("/snapshot"))

	dataVol, err := findVolume(stsData, "snapshot-volume")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(dataVol.HostPath).NotTo(BeNil())
	g.Expect(dataVol.HostPath.Path).To(Equal("/data"))

	_, err = findPersistentVolumeClaim(stsData, "snapshot-volume")
	g.Expect(err).To(HaveOccurred())

	proxyContainer, err := findContainerForRole(result, cluster, "proxy", coherenceContainer)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(proxyContainer.Env).NotTo(HaveEnvVarNamed("COH_SNAPSHOT_ENABLED"))
}
