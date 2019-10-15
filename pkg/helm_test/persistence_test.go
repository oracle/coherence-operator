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

func TestPersistenceWhenNotSetInMinimalCluster(t *testing.T) {
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

	// Assert that the COH_PERSISTENCE_ENABLED env-var is not set
	g.Expect(container.Env).NotTo(HaveEnvVarNamed("COH_PERSISTENCE_ENABLED"))
}

func TestPersistenceWhenNotSetInClusterWithExplicitRoles(t *testing.T) {
	g := NewGomegaWithT(t)

	result, cluster, err := CreateCluster("minimal-explicit-roles.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	for name := range cluster.GetRoles() {
		sts, err := findStatefulSet(result, cluster, name)
		g.Expect(err).NotTo(HaveOccurred())

		container, err := findContainer(sts, coherenceContainer)
		g.Expect(err).NotTo(HaveOccurred())

		g.Expect(container.Env).NotTo(HaveEnvVarNamed("COH_PERSISTENCE_ENABLED"))
	}
}

func TestPersistenceWithPVCWhenSetForImplicitRole(t *testing.T) {
	g := NewGomegaWithT(t)

	result, cluster, err := CreateCluster("persistence-test-pvc-implicit-role.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	sts, err := findStatefulSet(result, cluster, cohv1.DefaultRoleName)
	g.Expect(err).NotTo(HaveOccurred())

	container, err := findContainer(sts, coherenceContainer)
	g.Expect(err).NotTo(HaveOccurred())

	g.Expect(container.Env).To(HaveEnvVar(corev1.EnvVar{Name: "COH_PERSISTENCE_ENABLED", Value: "true"}))

	mount, err := findVolumeMount(container, "persistence-volume")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(mount.MountPath).To(Equal("/persistence"))

	_, err = findVolume(sts, "persistence-volume")
	g.Expect(err).To(HaveOccurred())

	pvc, err := findPersistentVolumeClaim(sts, "persistence-volume")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(*pvc.Spec.StorageClassName).To(Equal("foo"))
}

func TestPersistenceWithPVCWhenDefaultSetInClusterWithExplicitRoles(t *testing.T) {
	g := NewGomegaWithT(t)

	result, cluster, err := CreateCluster("persistence-test-pvc-explicit-roles-with-default.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	stsData, err := findStatefulSet(result, cluster, "data")
	g.Expect(err).NotTo(HaveOccurred())

	stsProxy, err := findStatefulSet(result, cluster, "proxy")
	g.Expect(err).NotTo(HaveOccurred())

	dataContainer, err := findContainerForRole(result, cluster, "data", coherenceContainer)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(dataContainer.Env).To(HaveEnvVar(corev1.EnvVar{Name: "COH_PERSISTENCE_ENABLED", Value: "true"}))

	dataMount, err := findVolumeMount(dataContainer, "persistence-volume")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(dataMount.MountPath).To(Equal("/persistence"))

	_, err = findVolume(stsData, "persistence-volume")
	g.Expect(err).To(HaveOccurred())

	dataPVC, err := findPersistentVolumeClaim(stsData, "persistence-volume")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(*dataPVC.Spec.StorageClassName).To(Equal("foo"))

	proxyContainer, err := findContainerForRole(result, cluster, "proxy", coherenceContainer)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(proxyContainer.Env).To(HaveEnvVar(corev1.EnvVar{Name: "COH_PERSISTENCE_ENABLED", Value: "true"}))

	proxyMount, err := findVolumeMount(dataContainer, "persistence-volume")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(proxyMount.MountPath).To(Equal("/persistence"))

	_, err = findVolume(stsProxy, "persistence-volume")
	g.Expect(err).To(HaveOccurred())

	proxyPVC, err := findPersistentVolumeClaim(stsProxy, "persistence-volume")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(*proxyPVC.Spec.StorageClassName).To(Equal("foo"))
}

func TestPersistenceWithPVCWhenDefaultSetAndOverriddenInClusterWithExplicitRoles(t *testing.T) {
	g := NewGomegaWithT(t)

	result, cluster, err := CreateCluster("persistence-test-pvc-explicit-roles-override-default.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	stsData, err := findStatefulSet(result, cluster, "data")
	g.Expect(err).NotTo(HaveOccurred())

	dataContainer, err := findContainerForRole(result, cluster, "data", coherenceContainer)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(dataContainer.Env).To(HaveEnvVar(corev1.EnvVar{Name: "COH_PERSISTENCE_ENABLED", Value: "true"}))

	dataMount, err := findVolumeMount(dataContainer, "persistence-volume")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(dataMount.MountPath).To(Equal("/persistence"))

	_, err = findVolume(stsData, "persistence-volume")
	g.Expect(err).To(HaveOccurred())

	dataPVC, err := findPersistentVolumeClaim(stsData, "persistence-volume")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(*dataPVC.Spec.StorageClassName).To(Equal("foo"))

	proxyContainer, err := findContainerForRole(result, cluster, "proxy", coherenceContainer)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(proxyContainer.Env).NotTo(HaveEnvVarNamed("COH_PERSISTENCE_ENABLED"))
}

func TestPersistenceWithPVCWhenSetInClusterWithExplicitRoles(t *testing.T) {
	g := NewGomegaWithT(t)

	result, cluster, err := CreateCluster("persistence-test-pvc-explicit-roles.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	stsData, err := findStatefulSet(result, cluster, "data")
	g.Expect(err).NotTo(HaveOccurred())

	dataContainer, err := findContainerForRole(result, cluster, "data", coherenceContainer)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(dataContainer.Env).To(HaveEnvVar(corev1.EnvVar{Name: "COH_PERSISTENCE_ENABLED", Value: "true"}))

	dataMount, err := findVolumeMount(dataContainer, "persistence-volume")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(dataMount.MountPath).To(Equal("/persistence"))

	_, err = findVolume(stsData, "persistence-volume")
	g.Expect(err).To(HaveOccurred())

	dataPVC, err := findPersistentVolumeClaim(stsData, "persistence-volume")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(*dataPVC.Spec.StorageClassName).To(Equal("foo"))

	proxyContainer, err := findContainerForRole(result, cluster, "proxy", coherenceContainer)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(proxyContainer.Env).NotTo(HaveEnvVarNamed("COH_PERSISTENCE_ENABLED"))
}

func TestPersistenceWithVolumeWhenSetForImplicitRole(t *testing.T) {
	g := NewGomegaWithT(t)

	result, cluster, err := CreateCluster("persistence-test-volume-implicit-role.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	sts, err := findStatefulSet(result, cluster, cohv1.DefaultRoleName)
	g.Expect(err).NotTo(HaveOccurred())

	container, err := findContainer(sts, coherenceContainer)
	g.Expect(err).NotTo(HaveOccurred())

	g.Expect(container.Env).To(HaveEnvVar(corev1.EnvVar{Name: "COH_PERSISTENCE_ENABLED", Value: "true"}))

	mount, err := findVolumeMount(container, "persistence-volume")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(mount.MountPath).To(Equal("/persistence"))

	vol, err := findVolume(sts, "persistence-volume")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(vol.HostPath).NotTo(BeNil())
	g.Expect(vol.HostPath.Path).To(Equal("/data"))

	_, err = findPersistentVolumeClaim(sts, "persistence-volume")
	g.Expect(err).To(HaveOccurred())
}

func TestPersistenceWithVolumeWhenDefaultSetInClusterWithExplicitRoles(t *testing.T) {
	g := NewGomegaWithT(t)

	result, cluster, err := CreateCluster("persistence-test-volume-explicit-roles-with-default.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	stsData, err := findStatefulSet(result, cluster, "data")
	g.Expect(err).NotTo(HaveOccurred())

	stsProxy, err := findStatefulSet(result, cluster, "proxy")
	g.Expect(err).NotTo(HaveOccurred())

	dataContainer, err := findContainerForRole(result, cluster, "data", coherenceContainer)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(dataContainer.Env).To(HaveEnvVar(corev1.EnvVar{Name: "COH_PERSISTENCE_ENABLED", Value: "true"}))

	dataMount, err := findVolumeMount(dataContainer, "persistence-volume")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(dataMount.MountPath).To(Equal("/persistence"))

	dataVol, err := findVolume(stsData, "persistence-volume")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(dataVol.HostPath).NotTo(BeNil())
	g.Expect(dataVol.HostPath.Path).To(Equal("/data"))

	_, err = findPersistentVolumeClaim(stsData, "persistence-volume")
	g.Expect(err).To(HaveOccurred())

	proxyContainer, err := findContainerForRole(result, cluster, "proxy", coherenceContainer)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(proxyContainer.Env).To(HaveEnvVar(corev1.EnvVar{Name: "COH_PERSISTENCE_ENABLED", Value: "true"}))

	proxyMount, err := findVolumeMount(dataContainer, "persistence-volume")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(proxyMount.MountPath).To(Equal("/persistence"))

	proxyVol, err := findVolume(stsProxy, "persistence-volume")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(proxyVol.HostPath).NotTo(BeNil())
	g.Expect(proxyVol.HostPath.Path).To(Equal("/data"))

	_, err = findPersistentVolumeClaim(stsProxy, "persistence-volume")
	g.Expect(err).To(HaveOccurred())
}

func TestPersistenceWithVolumeWhenDefaultSetAndOverriddenInClusterWithExplicitRoles(t *testing.T) {
	g := NewGomegaWithT(t)

	result, cluster, err := CreateCluster("persistence-test-volume-explicit-roles-override-default.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	stsData, err := findStatefulSet(result, cluster, "data")
	g.Expect(err).NotTo(HaveOccurred())

	dataContainer, err := findContainerForRole(result, cluster, "data", coherenceContainer)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(dataContainer.Env).To(HaveEnvVar(corev1.EnvVar{Name: "COH_PERSISTENCE_ENABLED", Value: "true"}))

	dataMount, err := findVolumeMount(dataContainer, "persistence-volume")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(dataMount.MountPath).To(Equal("/persistence"))

	dataVol, err := findVolume(stsData, "persistence-volume")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(dataVol.HostPath).NotTo(BeNil())
	g.Expect(dataVol.HostPath.Path).To(Equal("/data"))

	_, err = findPersistentVolumeClaim(stsData, "persistence-volume")
	g.Expect(err).To(HaveOccurred())

	proxyContainer, err := findContainerForRole(result, cluster, "proxy", coherenceContainer)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(proxyContainer.Env).NotTo(HaveEnvVarNamed("COH_PERSISTENCE_ENABLED"))
}

func TestPersistenceWithVolumeWhenSetInClusterWithExplicitRoles(t *testing.T) {
	g := NewGomegaWithT(t)

	result, cluster, err := CreateCluster("persistence-test-volume-explicit-roles.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	stsData, err := findStatefulSet(result, cluster, "data")
	g.Expect(err).NotTo(HaveOccurred())

	dataContainer, err := findContainerForRole(result, cluster, "data", coherenceContainer)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(dataContainer.Env).To(HaveEnvVar(corev1.EnvVar{Name: "COH_PERSISTENCE_ENABLED", Value: "true"}))

	dataMount, err := findVolumeMount(dataContainer, "persistence-volume")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(dataMount.MountPath).To(Equal("/persistence"))

	dataVol, err := findVolume(stsData, "persistence-volume")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(dataVol.HostPath).NotTo(BeNil())
	g.Expect(dataVol.HostPath.Path).To(Equal("/data"))

	_, err = findPersistentVolumeClaim(stsData, "persistence-volume")
	g.Expect(err).To(HaveOccurred())

	proxyContainer, err := findContainerForRole(result, cluster, "proxy", coherenceContainer)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(proxyContainer.Env).NotTo(HaveEnvVarNamed("COH_PERSISTENCE_ENABLED"))
}
