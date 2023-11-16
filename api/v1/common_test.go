/*
 * Copyright (c) 2019, 2023, Oracle and/or its affiliates.
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
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
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
	testOperatorImage  = "oracle/operator:1.2.3"
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
	assertEnvironmentVariablesForPodTemplate(t, &stsActual.Spec.Template, &stsExpected.Spec.Template)
}

func assertEnvironmentVariablesForJob(t *testing.T, actual, expected *batchv1.Job) {
	assertEnvironmentVariablesForPodTemplate(t, &actual.Spec.Template, &expected.Spec.Template)
}

func assertEnvironmentVariablesForPodTemplate(t *testing.T, actual, expected *corev1.PodTemplateSpec) {
	g := NewGomegaWithT(t)

	for _, contExpected := range expected.Spec.InitContainers {
		contActual := coh.FindInitContainerInPodTemplate(contExpected.Name, actual)
		g.Expect(contActual).NotTo(BeNil(), "Error asserting environment variables, could not find init-container with name "+contExpected.Name)
		assertEnvironmentVariablesForContainer(t, contActual, &contExpected)
	}

	for _, contExpected := range expected.Spec.Containers {
		contActual := coh.FindContainerInPodTemplate(contExpected.Name, actual)
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
	err = os.WriteFile(fmt.Sprintf("%s%c%s-Actual.json", dir, os.PathSeparator, stsActual.Name), jsonActual, os.ModePerm)
	g.Expect(err).NotTo(HaveOccurred())

	// Dump the json for the expected StatefulSet for debugging failures
	jsonExpected, err := json.MarshalIndent(stsExpected, "", "    ")
	g.Expect(err).NotTo(HaveOccurred())
	err = os.WriteFile(fmt.Sprintf("%s%c%s-Expected.json", dir, os.PathSeparator, stsActual.Name), jsonExpected, os.ModePerm)
	g.Expect(err).NotTo(HaveOccurred())

	assertEnvironmentVariables(t, stsActual, stsExpected)
	assertEnvironmentVariables(t, stsActual, stsExpected)

	diffs := deep.Equal(*stsActual, *stsExpected)
	msg := "StatefulSets not equal:"
	if len(diffs) > 0 {
		// Dump the diffs
		err = os.WriteFile(fmt.Sprintf("%s%c%s-Diff.txt", dir, os.PathSeparator, stsActual.Name), []byte(strings.Join(diffs, "\n")), os.ModePerm)
		g.Expect(err).NotTo(HaveOccurred())
		for _, diff := range diffs {
			msg = msg + "\n" + diff
		}
		t.Errorf(msg)
	}
}

func assertJob(t *testing.T, res coh.Resource, expected *batchv1.Job) {
	g := NewGomegaWithT(t)

	dir, err := helper.EnsureLogsDir(t.Name())
	g.Expect(err).NotTo(HaveOccurred())

	g.Expect(res.Kind).To(Equal(coh.ResourceTypeJob))
	g.Expect(res.Name).To(Equal(expected.GetName()))

	jobActual := res.Spec.(*batchv1.Job)

	// sort env vars before diff
	sortEnvVarsForJob(jobActual)
	sortEnvVarsForJob(expected)

	// sort volume mounts before diff
	sortVolumeMountsForJob(jobActual)
	sortVolumeMountsForJob(expected)

	// sort volumes before diff
	sortVolumesForJob(jobActual)
	sortVolumesForJob(expected)

	// sort ports before diff
	sortPortsForJob(jobActual)
	sortPortsForJob(expected)

	// Dump the json for the actual StatefulSet for debugging failures
	jsonActual, err := json.MarshalIndent(jobActual, "", "    ")
	g.Expect(err).NotTo(HaveOccurred())
	err = os.WriteFile(fmt.Sprintf("%s%c%s-Actual.json", dir, os.PathSeparator, jobActual.Name), jsonActual, os.ModePerm)
	g.Expect(err).NotTo(HaveOccurred())

	// Dump the json for the expected StatefulSet for debugging failures
	jsonExpected, err := json.MarshalIndent(expected, "", "    ")
	g.Expect(err).NotTo(HaveOccurred())
	err = os.WriteFile(fmt.Sprintf("%s%c%s-Expected.json", dir, os.PathSeparator, jobActual.Name), jsonExpected, os.ModePerm)
	g.Expect(err).NotTo(HaveOccurred())

	assertEnvironmentVariablesForJob(t, jobActual, expected)
	assertEnvironmentVariablesForJob(t, jobActual, expected)

	diffs := deep.Equal(*jobActual, *expected)
	msg := "Jobs not equal:"
	if len(diffs) > 0 {
		// Dump the diffs
		err = os.WriteFile(fmt.Sprintf("%s%c%s-Diff.txt", dir, os.PathSeparator, jobActual.Name), []byte(strings.Join(diffs, "\n")), os.ModePerm)
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
	podTemplate := createMinimalExpectedPodSpec(deployment)

	// The StatefulSet
	sts := appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:   deployment.Name,
			Labels: labels,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: pointer.Int32(spec.GetReplicas()),
			Selector: &metav1.LabelSelector{
				MatchLabels: selector,
			},
			ServiceName:          deployment.GetHeadlessServiceName(),
			RevisionHistoryLimit: pointer.Int32(5),
			UpdateStrategy:       appsv1.StatefulSetUpdateStrategy{Type: appsv1.RollingUpdateStatefulSetStrategyType},
			PodManagementPolicy:  appsv1.ParallelPodManagement,
			Template:             podTemplate,
		},
		Status: appsv1.StatefulSetStatus{
			Replicas: 0,
		},
	}

	return &sts
}

// Create the expected default Job for a spec with nothing but the minimal fields set.
func createMinimalExpectedJob(deployment *coh.CoherenceJob) *batchv1.Job {
	spec := deployment.Spec
	labels := deployment.CreateCommonLabels()
	labels[coh.LabelComponent] = coh.LabelComponentCoherenceStatefulSet
	podTemplate := createMinimalExpectedPodSpec(deployment)

	podTemplate.Spec.RestartPolicy = corev1.RestartPolicyNever

	job := batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:   deployment.Name,
			Labels: labels,
		},
		Spec: batchv1.JobSpec{
			Parallelism: pointer.Int32(spec.GetReplicas()),
			Template:    podTemplate,
		},
		Status: batchv1.JobStatus{},
	}

	return &job
}

// Create the expected default PodTemplateSpec for a spec with nothing but the minimal fields set.
func createMinimalExpectedPodSpec(deployment coh.CoherenceResource) corev1.PodTemplateSpec {
	spec := deployment.GetSpec()
	podLabels := deployment.CreateCommonLabels()
	podLabels[coh.LabelComponent] = coh.LabelComponentCoherencePod
	podLabels[coh.LabelCoherenceWKAMember] = "true"

	emptyVolume := corev1.VolumeSource{
		EmptyDir: &corev1.EmptyDirVolumeSource{},
	}

	envVars := []corev1.EnvVar{
		{
			Name:  "COH_CLUSTER_NAME",
			Value: deployment.GetName(),
		},
		{
			Name:  "COHERENCE_LOCALPORT",
			Value: fmt.Sprintf("%d", coh.DefaultUnicastPort),
		},
		{
			Name:  "COHERENCE_LOCALPORT_ADJUST",
			Value: fmt.Sprintf("%d", coh.DefaultUnicastPortAdjust),
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
			Value: deployment.GetWKA(),
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
					Optional:             pointer.Bool(true),
				},
			},
		},
		{
			Name:  "OPERATOR_REQUEST_TIMEOUT",
			Value: "120",
		},
	}

	if deployment.GetType() == coh.CoherenceTypeJob {
		envVars = append(envVars, corev1.EnvVar{
			Name:  "COH_STORAGE_ENABLED",
			Value: "false",
		})
	}

	lp, _ := deployment.GetSpec().Coherence.GetLocalPorts()

	// The Coherence Container
	cohContainer := corev1.Container{
		Name:    coh.ContainerNameCoherence,
		Image:   testCoherenceImage,
		Command: []string{coh.RunnerCommand, "server"},
		Ports: []corev1.ContainerPort{
			{
				Name:          coh.PortNameCoherence,
				ContainerPort: 7,
				Protocol:      "TCP",
			},
			{
				Name:          coh.PortNameHealth,
				ContainerPort: spec.GetHealthPort(),
				Protocol:      "TCP",
			},
			{
				Name:          coh.PortNameCoherenceLocal,
				ContainerPort: lp,
				Protocol:      corev1.ProtocolTCP,
			},
			{
				Name:          coh.PortNameCoherenceCluster,
				ContainerPort: coh.DefaultClusterPort,
				Protocol:      corev1.ProtocolTCP,
			},
		},
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
		Env: envVars,
	}

	if cohImage := spec.GetCoherenceImage(); cohImage != nil {
		cohContainer.Image = *cohImage
	}

	// The Operator Init-Container
	initContainer := corev1.Container{
		Name:    coh.ContainerNameOperatorInit,
		Image:   testOperatorImage,
		Command: []string{coh.RunnerInitCommand, coh.RunnerInit},
		Env: []corev1.EnvVar{
			{
				Name:  "COH_CLUSTER_NAME",
				Value: deployment.GetName(),
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

	if operatorImage := spec.GetCoherenceOperatorImage(); operatorImage != nil {
		initContainer.Image = *operatorImage
	}

	podTemplate := corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels: podLabels,
		},
		Spec: corev1.PodSpec{
			InitContainers: []corev1.Container{initContainer},
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
	}

	return podTemplate
}

func sortEnvVars(sts *appsv1.StatefulSet) {
	if sts != nil {
		sortEnvVarsForPodSpec(&sts.Spec.Template)
	}
}

func sortEnvVarsForJob(job *batchv1.Job) {
	if job != nil {
		sortEnvVarsForPodSpec(&job.Spec.Template)
	}
}

func sortEnvVarsForPodSpec(template *corev1.PodTemplateSpec) {
	for _, c := range template.Spec.InitContainers {
		sort.SliceStable(c.Env, func(i, j int) bool {
			return c.Env[i].Name < c.Env[j].Name
		})
	}
	for _, c := range template.Spec.Containers {
		sort.SliceStable(c.Env, func(i, j int) bool {
			return c.Env[i].Name < c.Env[j].Name
		})
	}
}

func sortVolumeMounts(sts *appsv1.StatefulSet) {
	if sts != nil {
		sortVolumeMountsForPodSpec(&sts.Spec.Template)
	}
}

func sortVolumeMountsForJob(job *batchv1.Job) {
	if job != nil {
		sortVolumeMountsForPodSpec(&job.Spec.Template)
	}
}

func sortVolumeMountsForPodSpec(template *corev1.PodTemplateSpec) {
	for _, c := range template.Spec.InitContainers {
		sort.SliceStable(c.VolumeMounts, func(i, j int) bool {
			return c.VolumeMounts[i].Name < c.VolumeMounts[j].Name
		})
	}
	for _, c := range template.Spec.Containers {
		sort.SliceStable(c.VolumeMounts, func(i, j int) bool {
			return c.VolumeMounts[i].Name < c.VolumeMounts[j].Name
		})
	}
}

func sortVolumes(sts *appsv1.StatefulSet) {
	if sts != nil {
		sortVolumesForPodTemplate(&sts.Spec.Template)
	}
}

func sortVolumesForJob(job *batchv1.Job) {
	if job != nil {
		sortVolumesForPodTemplate(&job.Spec.Template)
	}
}

func sortVolumesForPodTemplate(template *corev1.PodTemplateSpec) {
	sort.SliceStable(template.Spec.Volumes, func(i, j int) bool {
		return template.Spec.Volumes[i].Name < template.Spec.Volumes[j].Name
	})
}

func sortPorts(sts *appsv1.StatefulSet) {
	if sts != nil {
		sortPortsForPodTemplate(&sts.Spec.Template)
	}
}

func sortPortsForJob(job *batchv1.Job) {
	if job != nil {
		sortPortsForPodTemplate(&job.Spec.Template)
	}
}

func sortPortsForPodTemplate(template *corev1.PodTemplateSpec) {
	for _, c := range template.Spec.InitContainers {
		sort.SliceStable(c.Ports, func(i, j int) bool {
			return c.Ports[i].Name < c.Ports[j].Name
		})
	}
	for _, c := range template.Spec.Containers {
		sort.SliceStable(c.Ports, func(i, j int) bool {
			return c.Ports[i].Name < c.Ports[j].Name
		})
	}
}

func addEnvVars(sts *appsv1.StatefulSet, containerName string, envVars ...corev1.EnvVar) {
	if sts != nil {
		addEnvVarsToPodSpec(&sts.Spec.Template, containerName, envVars...)
	}
}

func addEnvVarsToJob(job *batchv1.Job, containerName string, envVars ...corev1.EnvVar) {
	if job != nil {
		addEnvVarsToPodSpec(&job.Spec.Template, containerName, envVars...)
	}
}

func addEnvVarsToPodSpec(template *corev1.PodTemplateSpec, containerName string, envVars ...corev1.EnvVar) {
	for i, c := range template.Spec.InitContainers {
		if c.Name == containerName {
			addEnvVarsToContainer(&c, envVars...)
			template.Spec.InitContainers[i] = c
		}
	}
	for i, c := range template.Spec.Containers {
		if c.Name == containerName {
			addEnvVarsToContainer(&c, envVars...)
			template.Spec.Containers[i] = c
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

func addEnvVarsFrom(sts *appsv1.StatefulSet, containerName string, envVars ...corev1.EnvFromSource) {
	if sts != nil {
		addEnvVarsFromToPodSpec(&sts.Spec.Template, containerName, envVars...)
	}
}

func addEnvVarsFromToJob(job *batchv1.Job, containerName string, envVars ...corev1.EnvFromSource) {
	if job != nil {
		addEnvVarsFromToPodSpec(&job.Spec.Template, containerName, envVars...)
	}
}

func addEnvVarsFromToPodSpec(template *corev1.PodTemplateSpec, containerName string, envVars ...corev1.EnvFromSource) {
	for i, c := range template.Spec.InitContainers {
		if c.Name == containerName {
			addEnvVarsFromToContainer(&c, envVars...)
			template.Spec.InitContainers[i] = c
		}
	}
	for i, c := range template.Spec.Containers {
		if c.Name == containerName {
			addEnvVarsFromToContainer(&c, envVars...)
			template.Spec.Containers[i] = c
		}
	}
}

func addEnvVarsFromToContainer(c *corev1.Container, envVars ...corev1.EnvFromSource) {
	c.EnvFrom = append(c.EnvFrom, envVars...)
}

func addPorts(sts *appsv1.StatefulSet, containerName string, ports ...corev1.ContainerPort) {
	if sts != nil {
		addPortsToPodSpec(&sts.Spec.Template, containerName, ports...)
	}
}

func addPortsForJob(job *batchv1.Job, containerName string, ports ...corev1.ContainerPort) {
	if job != nil {
		addPortsToPodSpec(&job.Spec.Template, containerName, ports...)
	}
}

func addPortsToPodSpec(template *corev1.PodTemplateSpec, containerName string, ports ...corev1.ContainerPort) {
	for i, c := range template.Spec.InitContainers {
		if c.Name == containerName {
			addPortsToContainer(&c, ports...)
			template.Spec.InitContainers[i] = c
		}
	}
	for i, c := range template.Spec.Containers {
		if c.Name == containerName {
			addPortsToContainer(&c, ports...)
			template.Spec.Containers[i] = c
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
	s := coh.CoherenceStatefulSetResourceSpec{
		CoherenceResourceSpec: spec,
	}
	return createTestCoherenceDeployment(s)
}

func createTestCoherenceDeployment(spec coh.CoherenceStatefulSetResourceSpec) *coh.Coherence {
	return &coh.Coherence{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-deployment",
		},
		Spec: spec,
	}
}

func createTestCoherenceJob(spec coh.CoherenceResourceSpec) *coh.CoherenceJob {
	s := coh.CoherenceJobResourceSpec{
		CoherenceResourceSpec: spec,
	}
	return createTestCoherenceJobDeployment(s)
}

func createTestCoherenceJobDeployment(spec coh.CoherenceJobResourceSpec) *coh.CoherenceJob {
	return &coh.CoherenceJob{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-deployment",
		},
		Spec: spec,
	}
}

func assertStatefulSetCreation(t *testing.T, deployment *coh.Coherence, stsExpected *appsv1.StatefulSet) {
	viper.Set(operator.FlagCoherenceImage, testCoherenceImage)
	viper.Set(operator.FlagOperatorImage, testOperatorImage)

	res := deployment.Spec.CreateStatefulSetResource(deployment)
	assertStatefulSet(t, res, stsExpected)
}

func assertJobCreation(t *testing.T, deployment *coh.CoherenceJob, jobExpected *batchv1.Job) {
	viper.Set(operator.FlagCoherenceImage, testCoherenceImage)
	viper.Set(operator.FlagOperatorImage, testOperatorImage)

	res := deployment.Spec.CreateJobResource(deployment)
	assertJob(t, res, jobExpected)
}

func assertResourceCreation(t *testing.T, deployment *coh.Coherence) coh.Resources {
	g := NewGomegaWithT(t)
	viper.Set(operator.FlagCoherenceImage, testCoherenceImage)
	viper.Set(operator.FlagOperatorImage, testOperatorImage)

	res, err := deployment.CreateKubernetesResources()
	g.Expect(err).NotTo(HaveOccurred())
	return res
}
