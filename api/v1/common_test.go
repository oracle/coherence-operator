/*
 * Copyright (c) 2019, 2020, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1_test

import (
	"encoding/json"
	"fmt"
	"github.com/go-test/deep"
	. "github.com/onsi/gomega"
	coh "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/pkg/operator"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	"github.com/spf13/viper"
	"io/ioutil"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
	"os"
	"sort"
	"strings"
	"testing"
)

const (
	testCoherenceImage = "oracle/coherence-ce:1.2.3"
	testUtilsImage     = "oracle/operator:1.2.3-utils"
)

// Returns a pointer to an int32
func int32Ptr(x int32) *int32 {
	return &x
}

// Returns a pointer to an int32
func boolPtr(x bool) *bool {
	return &x
}

// Returns a pointer to a string
func stringPtr(x string) *string {
	return &x
}

func assertEnvironmentVariables(t *testing.T, stsActual, stsExpected *appsv1.StatefulSet) {
	g := NewGomegaWithT(t)

	for _, contExpected := range stsExpected.Spec.Template.Spec.InitContainers {
		contActual := coh.FindInitContainer(contExpected.Name, stsActual)
		g.Expect(contActual).NotTo(BeNil(), "Error asserting environment variables, could not find init-container with name "+contExpected.Name)
		assertEnvironmentVariablesForContainer(t, contActual, &contExpected)
	}

	for _, contExpected := range stsExpected.Spec.Template.Spec.Containers {
		contActual := coh.FindContainer(contExpected.Name, stsActual)
		g.Expect(contActual).NotTo(BeNil(), "Error asserting environment variables, could not find container with name "+contExpected.Name)
		assertEnvironmentVariablesForContainer(t, contActual, &contExpected)
	}
}

func assertEnvironmentVariablesForContainer(t *testing.T, c, cExpected *corev1.Container) {
	g := NewGomegaWithT(t)

	env := envVarsToMap(c)
	envExpected := envVarsToMap(cExpected)

	equal := deep.Equal(env, envExpected)
	g.Expect(equal).To(BeNil(), fmt.Sprintf("Environment variable mis-match for container '%s'", cExpected.Name))
}

func envVarsToMap(c *corev1.Container) map[string]corev1.EnvVar {
	var m = make(map[string]corev1.EnvVar)
	for _, e := range c.Env {
		m[e.Name] = e
	}
	return m
}

func assertStatefulSet(t *testing.T, res coh.Resource, stsExpected *appsv1.StatefulSet) {
	g := NewGomegaWithT(t)

	dir, err := helper.EnsureLogsDir(t.Name())
	g.Expect(err).NotTo(HaveOccurred())

	g.Expect(res.Kind).To(Equal(coh.ResourceTypeStatefulSet))
	g.Expect(res.Name).To(Equal(stsExpected.GetName()))

	stsActual := res.Spec.(*appsv1.StatefulSet)

	// sort env vars before diff
	sortEnvVars(stsActual)
	sortEnvVars(stsExpected)

	// sort volume mounts before diff
	sortVolumeMounts(stsActual)
	sortVolumeMounts(stsExpected)

	// sort volumes before diff
	sortVolumes(stsActual)
	sortVolumes(stsExpected)

	// sort ports before diff
	sortPorts(stsActual)
	sortPorts(stsExpected)

	// Dump the json for the actual StatefulSet for debugging failures
	jsonActual, err := json.MarshalIndent(stsActual, "", "    ")
	g.Expect(err).NotTo(HaveOccurred())
	err = ioutil.WriteFile(fmt.Sprintf("%s%c%s-Actual.json", dir, os.PathSeparator, stsActual.Name), jsonActual, os.ModePerm)
	g.Expect(err).NotTo(HaveOccurred())

	// Dump the json for the expected StatefulSet for debugging failures
	jsonExpected, err := json.MarshalIndent(stsExpected, "", "    ")
	g.Expect(err).NotTo(HaveOccurred())
	err = ioutil.WriteFile(fmt.Sprintf("%s%c%s-Expected.json", dir, os.PathSeparator, stsActual.Name), jsonExpected, os.ModePerm)
	g.Expect(err).NotTo(HaveOccurred())

	assertEnvironmentVariables(t, stsActual, stsExpected)
	assertEnvironmentVariables(t, stsActual, stsExpected)

	diffs := deep.Equal(*stsActual, *stsExpected)
	msg := "StatefulSets not equal:"
	if len(diffs) > 0 {
		// Dump the diffs
		err = ioutil.WriteFile(fmt.Sprintf("%s%c%s-Diff.txt", dir, os.PathSeparator, stsActual.Name), []byte(strings.Join(diffs, "\n")), os.ModePerm)
		g.Expect(err).NotTo(HaveOccurred())
		for _, diff := range diffs {
			msg = msg + "\n" + diff
		}
		t.Errorf(msg)
	}
}

// Create the expected default StatefulSet for a spec with nothing but the minimal fields set.
func createMinimalExpectedStatefulSet(deployment *coh.Coherence) *appsv1.StatefulSet {
	spec := deployment.Spec
	labels := deployment.CreateCommonLabels()
	labels[coh.LabelComponent] = coh.LabelComponentCoherenceStatefulSet
	selector := deployment.CreateCommonLabels()
	selector[coh.LabelComponent] = coh.LabelComponentCoherencePod
	podLabels := deployment.CreateCommonLabels()
	podLabels[coh.LabelComponent] = coh.LabelComponentCoherencePod
	podLabels[coh.LabelCoherenceWKAMember] = "true"

	emptyVolume := corev1.VolumeSource{
		EmptyDir: &corev1.EmptyDirVolumeSource{},
	}

	// The Coherence Container
	cohContainer := corev1.Container{
		Name:    coh.ContainerNameCoherence,
		Image:   testCoherenceImage,
		Command: []string{coh.RunnerCommand, "server"},
		Ports: []corev1.ContainerPort{
			{
				Name:          "coherence",
				ContainerPort: 7,
				Protocol:      "TCP",
			},
			{
				Name:          "health",
				ContainerPort: spec.GetHealthPort(),
				Protocol:      "TCP",
			},
		},
		Resources:      spec.CreateDefaultResources(),
		ReadinessProbe: spec.UpdateDefaultReadinessProbeAction(spec.CreateDefaultReadinessProbe()),
		LivenessProbe:  spec.UpdateDefaultLivenessProbeAction(spec.CreateDefaultLivenessProbe()),
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      coh.VolumeNameJVM,
				MountPath: coh.VolumeMountPathJVM,
				ReadOnly:  false,
			},
			{
				Name:      coh.VolumeNameUtils,
				MountPath: coh.VolumeMountPathUtils,
				ReadOnly:  false,
			},
		},
		Env: []corev1.EnvVar{
			{
				Name:  "COH_CLUSTER_NAME",
				Value: deployment.Name,
			},
			{
				Name:  "COH_HEALTH_PORT",
				Value: fmt.Sprintf("%d", spec.GetHealthPort()),
			},
			{
				Name: "COH_MACHINE_NAME",
				ValueFrom: &corev1.EnvVarSource{
					FieldRef: &corev1.ObjectFieldSelector{
						FieldPath: "spec.nodeName",
					},
				},
			},
			{
				Name: "COH_MEMBER_NAME",
				ValueFrom: &corev1.EnvVarSource{
					FieldRef: &corev1.ObjectFieldSelector{
						FieldPath: "metadata.name",
					},
				},
			},
			{
				Name:  "COH_METRICS_ENABLED",
				Value: "false",
			},
			{
				Name:  "COH_MGMT_ENABLED",
				Value: "false",
			},
			{
				Name: "COH_POD_UID",
				ValueFrom: &corev1.EnvVarSource{
					FieldRef: &corev1.ObjectFieldSelector{
						FieldPath: "metadata.uid",
					},
				},
			},
			{
				Name:  "COH_RACK_INFO_LOCATION",
				Value: "http://$(OPERATOR_HOST)/rack/$(COH_MACHINE_NAME)",
			},
			{
				Name:  "COH_ROLE",
				Value: deployment.GetRoleName(),
			},
			{
				Name:  "COH_SITE_INFO_LOCATION",
				Value: "http://$(OPERATOR_HOST)/site/$(COH_MACHINE_NAME)",
			},
			{
				Name:  "COH_UTIL_DIR",
				Value: coh.VolumeMountPathUtils,
			},
			{
				Name:  "COH_WKA",
				Value: deployment.GetWkaServiceName() + ".svc.cluster.local",
			},
			{
				Name:  "JVM_GC_LOGGING",
				Value: "false",
			},
			{
				Name:  "JVM_USE_CONTAINER_LIMITS",
				Value: "true",
			},
			{
				Name: "OPERATOR_HOST",
				ValueFrom: &corev1.EnvVarSource{
					SecretKeyRef: &corev1.SecretKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{Name: coh.OperatorConfigName},
						Key:                  coh.OperatorConfigKeyHost,
						Optional:             pointer.BoolPtr(true),
					},
				},
			},
			{
				Name:  "OPERATOR_REQUEST_TIMEOUT",
				Value: "120",
			},
			{
				Name:  coh.EnvVarCohIdentity,
				Value: deployment.Name + "@" + deployment.Namespace,
			},
		},
	}

	if cohImage := spec.GetCoherenceImage(); cohImage != nil {
		cohContainer.Image = *cohImage
	}

	// The Utils Init-Container
	utilsContainer := corev1.Container{
		Name:    coh.ContainerNameUtils,
		Image:   testUtilsImage,
		Command: []string{coh.UtilsInitCommand, coh.RunnerInit},
		Env: []corev1.EnvVar{
			{
				Name:  "COH_CLUSTER_NAME",
				Value: deployment.Name,
			},
			{
				Name:  "COH_UTIL_DIR",
				Value: coh.VolumeMountPathUtils,
			},
		},
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      coh.VolumeNameJVM,
				MountPath: coh.VolumeMountPathJVM,
				ReadOnly:  false,
			},
			{
				Name:      coh.VolumeNameUtils,
				MountPath: coh.VolumeMountPathUtils,
				ReadOnly:  false,
			},
		},
	}

	if utilsImage := spec.GetCoherenceUtilsImage(); utilsImage != nil {
		utilsContainer.Image = *utilsImage
	}

	// The StatefulSet
	sts := appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:   deployment.Name,
			Labels: labels,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: pointer.Int32Ptr(spec.GetReplicas()),
			Selector: &metav1.LabelSelector{
				MatchLabels: selector,
			},
			ServiceName:          deployment.GetHeadlessServiceName(),
			RevisionHistoryLimit: pointer.Int32Ptr(5),
			UpdateStrategy:       appsv1.StatefulSetUpdateStrategy{Type: appsv1.RollingUpdateStatefulSetStrategyType},
			PodManagementPolicy:  appsv1.ParallelPodManagement,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: podLabels,
				},
				Spec: corev1.PodSpec{
					InitContainers: []corev1.Container{utilsContainer},
					Containers:     []corev1.Container{cohContainer},
					Volumes: []corev1.Volume{
						{
							Name:         coh.VolumeNameJVM,
							VolumeSource: emptyVolume,
						},
						{
							Name:         coh.VolumeNameUtils,
							VolumeSource: emptyVolume,
						},
					},
					Affinity: spec.CreateDefaultPodAffinity(deployment),
				},
			},
		},
		Status: appsv1.StatefulSetStatus{
			Replicas: 0,
		},
	}
	return &sts
}

func sortEnvVars(sts *appsv1.StatefulSet) {
	for _, c := range sts.Spec.Template.Spec.InitContainers {
		sort.SliceStable(c.Env, func(i, j int) bool {
			return c.Env[i].Name < c.Env[j].Name
		})
	}
	for _, c := range sts.Spec.Template.Spec.Containers {
		sort.SliceStable(c.Env, func(i, j int) bool {
			return c.Env[i].Name < c.Env[j].Name
		})
	}
}

func sortVolumeMounts(sts *appsv1.StatefulSet) {
	for _, c := range sts.Spec.Template.Spec.InitContainers {
		sort.SliceStable(c.VolumeMounts, func(i, j int) bool {
			return c.VolumeMounts[i].Name < c.VolumeMounts[j].Name
		})
	}
	for _, c := range sts.Spec.Template.Spec.Containers {
		sort.SliceStable(c.VolumeMounts, func(i, j int) bool {
			return c.VolumeMounts[i].Name < c.VolumeMounts[j].Name
		})
	}
}

func sortVolumes(sts *appsv1.StatefulSet) {
	sort.SliceStable(sts.Spec.Template.Spec.Volumes, func(i, j int) bool {
		return sts.Spec.Template.Spec.Volumes[i].Name < sts.Spec.Template.Spec.Volumes[j].Name
	})
}

func sortPorts(sts *appsv1.StatefulSet) {
	for _, c := range sts.Spec.Template.Spec.InitContainers {
		sort.SliceStable(c.Ports, func(i, j int) bool {
			return c.Ports[i].Name < c.Ports[j].Name
		})
	}
	for _, c := range sts.Spec.Template.Spec.Containers {
		sort.SliceStable(c.Ports, func(i, j int) bool {
			return c.Ports[i].Name < c.Ports[j].Name
		})
	}
}

func addEnvVars(sts *appsv1.StatefulSet, containerName string, envVars ...corev1.EnvVar) {
	for i, c := range sts.Spec.Template.Spec.InitContainers {
		if c.Name == containerName {
			addEnvVarsToContainer(&c, envVars...)
			sts.Spec.Template.Spec.InitContainers[i] = c
		}
	}
	for i, c := range sts.Spec.Template.Spec.Containers {
		if c.Name == containerName {
			addEnvVarsToContainer(&c, envVars...)
			sts.Spec.Template.Spec.Containers[i] = c
		}
	}
}

func addEnvVarsToContainer(c *corev1.Container, envVars ...corev1.EnvVar) {
	for _, evAdd := range envVars {
		found := false
		for e, ev := range c.Env {
			if ev.Name == evAdd.Name {
				ev.Value = evAdd.Value
				ev.ValueFrom = evAdd.ValueFrom
				c.Env[e] = ev
				found = true
				break
			}
		}
		if !found {
			c.Env = append(c.Env, evAdd)
		}
	}
}

func addPorts(sts *appsv1.StatefulSet, containerName string, ports ...corev1.ContainerPort) {
	for i, c := range sts.Spec.Template.Spec.InitContainers {
		if c.Name == containerName {
			addPortsToContainer(&c, ports...)
			sts.Spec.Template.Spec.InitContainers[i] = c
		}
	}
	for i, c := range sts.Spec.Template.Spec.Containers {
		if c.Name == containerName {
			addPortsToContainer(&c, ports...)
			sts.Spec.Template.Spec.Containers[i] = c
		}
	}
}

func addPortsToContainer(c *corev1.Container, ports ...corev1.ContainerPort) {
	for _, portAdd := range ports {
		found := false
		for p, port := range c.Ports {
			if port.Name == portAdd.Name {
				port.ContainerPort = portAdd.ContainerPort
				port.HostIP = portAdd.HostIP
				port.HostPort = portAdd.HostPort
				port.Protocol = portAdd.Protocol
				c.Ports[p] = port
				found = true
				break
			}
		}
		if !found {
			c.Ports = append(c.Ports, portAdd)
		}
	}
}

func createTestDeployment(spec coh.CoherenceResourceSpec) *coh.Coherence {
	return &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-deployment",
		},
		Spec: spec,
	}
}

func assertStatefulSetCreation(t *testing.T, deployment *coh.Coherence, stsExpected *appsv1.StatefulSet) {
	viper.Set(operator.FlagCoherenceImage, testCoherenceImage)
	viper.Set(operator.FlagUtilsImage, testUtilsImage)

	res := deployment.Spec.CreateStatefulSet(deployment)
	assertStatefulSet(t, res, stsExpected)
}
