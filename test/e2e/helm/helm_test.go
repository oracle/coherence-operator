/*
 * Copyright (c) 2020, 2025, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package helm

import (
	goctx "context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	. "github.com/onsi/gomega"
	coh "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/pkg/operator"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	"github.com/oracle/coherence-operator/test/e2e/helper/matchers"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	crdv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestPodSecurityContext(t *testing.T) {
	g := NewGomegaWithT(t)

	result, err := helmInstall("--set", "podSecurityContext.runAsNonRoot=true",
		"--set", "podSecurityContext.runAsUser=1000")

	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(result).NotTo(BeNil())

	d := &appsv1.Deployment{}
	err = result.Get("coherence-operator", d)
	g.Expect(err).NotTo(HaveOccurred())

	// Should have a Pod securityContext
	g.Expect(d.Spec.Template.Spec.SecurityContext).NotTo(BeNil())
	ctx := *d.Spec.Template.Spec.SecurityContext
	g.Expect(ctx.RunAsNonRoot).NotTo(BeNil())
	g.Expect(*ctx.RunAsNonRoot).To(BeTrue())
	g.Expect(ctx.RunAsUser).NotTo(BeNil())
	g.Expect(*ctx.RunAsUser).To(Equal(int64(1000)))

	// Should not have a container securityContext
	for _, c := range d.Spec.Template.Spec.Containers {
		g.Expect(c.SecurityContext).To(BeNil())
	}
}

func TestContainerSecurityContext(t *testing.T) {
	g := NewGomegaWithT(t)

	result, err := helmInstall("--set", "securityContext.runAsNonRoot=true",
		"--set", "securityContext.runAsUser=1000")

	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(result).NotTo(BeNil())

	d := &appsv1.Deployment{}
	err = result.Get("coherence-operator", d)
	g.Expect(err).NotTo(HaveOccurred())

	// Should have a container securityContext
	for _, c := range d.Spec.Template.Spec.Containers {
		g.Expect(c.SecurityContext).NotTo(BeNil())
		ctx := *c.SecurityContext
		g.Expect(ctx.RunAsNonRoot).NotTo(BeNil())
		g.Expect(*ctx.RunAsNonRoot).To(BeTrue())
		g.Expect(ctx.RunAsUser).NotTo(BeNil())
		g.Expect(*ctx.RunAsUser).To(Equal(int64(1000)))
	}

	// Should not have a Pod securityContext
	g.Expect(d.Spec.Template.Spec.SecurityContext).To(BeNil())
}

func TestCreateWebhookCertSecretByDefault(t *testing.T) {
	g := NewGomegaWithT(t)
	result, err := helmInstall()
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(result).NotTo(BeNil())

	secret := &corev1.Secret{}
	err = result.Get("coherence-webhook-server-cert", secret)
	g.Expect(err).NotTo(HaveOccurred())
}

func TestCreateWebhookCertSecretWithName(t *testing.T) {
	g := NewGomegaWithT(t)
	result, err := helmInstall("--set", "webhookCertSecret=foo")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(result).NotTo(BeNil())

	secret := &corev1.Secret{}
	err = result.Get("foo", secret)
	g.Expect(err).NotTo(HaveOccurred())

	dep := &appsv1.Deployment{}
	err = result.Get("coherence-operator", dep)
	g.Expect(err).NotTo(HaveOccurred())

	c := findContainer("manager", dep)
	g.Expect(c).NotTo(BeNil())

	senv := findEnvVar("WEBHOOK_SECRET", c)
	g.Expect(senv).NotTo(BeNil())
	g.Expect(senv.Value).To(Equal("foo"))
}

func TestNotCreateWebhookCertSecretIfManualCertManager(t *testing.T) {
	g := NewGomegaWithT(t)
	result, err := helmInstall("--set", "webhookCertType=manual")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(result).NotTo(BeNil())

	secret := &corev1.Secret{}
	err = result.Get("coherence-webhook-server-cert", secret)
	g.Expect(err).To(HaveOccurred())
}

func TestDisableWebhooks(t *testing.T) {
	g := NewGomegaWithT(t)
	result, err := helmInstall("--set", "webhooks=false")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(result).NotTo(BeNil())

	dep := &appsv1.Deployment{}
	err = result.Get("coherence-operator", dep)
	g.Expect(err).NotTo(HaveOccurred())

	c := findContainer("manager", dep)
	g.Expect(c).NotTo(BeNil())

	g.Expect(c.Args).NotTo(BeNil())
	g.Expect(c.Args).Should(ContainElements("operator", "--enable-leader-election", "--enable-webhook=false"))
}

func TestDisableJobCRD(t *testing.T) {
	g := NewGomegaWithT(t)

	err := RemoveCRDs()
	g.Expect(err).NotTo(HaveOccurred())

	result, err := helmInstall("--set", "allowCoherenceJobs=false")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(result).NotTo(BeNil())

	dep := &appsv1.Deployment{}
	err = result.Get("coherence-operator", dep)
	g.Expect(err).NotTo(HaveOccurred())

	c := findContainer("manager", dep)
	g.Expect(c).NotTo(BeNil())

	g.Expect(c.Args).NotTo(BeNil())
	g.Expect(c.Args).Should(ContainElements("operator", "--enable-leader-election", "--enable-jobs=false"))
}

func TestSetOnlySameNamespace(t *testing.T) {
	g := NewGomegaWithT(t)
	result, err := helmInstall("--set", "onlySameNamespace=true")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(result).NotTo(BeNil())

	dep := &appsv1.Deployment{}
	err = result.Get("coherence-operator", dep)
	g.Expect(err).NotTo(HaveOccurred())

	c := findContainer("manager", dep)
	g.Expect(c).NotTo(BeNil())

	ns := helper.GetTestNamespace()
	g.Expect(c.Env).NotTo(BeNil())
	g.Expect(c.Env).To(matchers.HaveEnvVar(corev1.EnvVar{Name: operator.EnvVarWatchNamespace, Value: ns}))
}

func TestSetOnlySameNamespaceIgnoresWatchNamespaces(t *testing.T) {
	g := NewGomegaWithT(t)
	result, err := helmInstall("--set", "watchNamespaces=foo", "--set", "onlySameNamespace=true")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(result).NotTo(BeNil())

	dep := &appsv1.Deployment{}
	err = result.Get("coherence-operator", dep)
	g.Expect(err).NotTo(HaveOccurred())

	c := findContainer("manager", dep)
	g.Expect(c).NotTo(BeNil())

	ns := helper.GetTestNamespace()
	g.Expect(c.Env).NotTo(BeNil())
	g.Expect(c.Env).To(matchers.HaveEnvVar(corev1.EnvVar{Name: operator.EnvVarWatchNamespace, Value: ns}))
}

func TestSetWatchNamespaces(t *testing.T) {
	g := NewGomegaWithT(t)
	result, err := helmInstall("--set", "watchNamespaces=foo")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(result).NotTo(BeNil())

	dep := &appsv1.Deployment{}
	err = result.Get("coherence-operator", dep)
	g.Expect(err).NotTo(HaveOccurred())

	c := findContainer("manager", dep)
	g.Expect(c).NotTo(BeNil())

	g.Expect(c.Env).NotTo(BeNil())
	g.Expect(c.Env).To(matchers.HaveEnvVar(corev1.EnvVar{Name: operator.EnvVarWatchNamespace, Value: "foo"}))
}

func TestSetOperatorImageInFull(t *testing.T) {
	g := NewGomegaWithT(t)
	result, err := helmInstall("--set", "image=foo.com/bar:1.0.0")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(result).NotTo(BeNil())

	dep := &appsv1.Deployment{}
	err = result.Get("coherence-operator", dep)
	g.Expect(err).NotTo(HaveOccurred())

	c := findContainer("manager", dep)
	g.Expect(c).NotTo(BeNil())

	g.Expect(c.Image).To(Equal("foo.com/bar:1.0.0"))
}

func TestSetOperatorImageRegistry(t *testing.T) {
	g := NewGomegaWithT(t)
	result, err := helmInstall("--set", "image.registry=foo.com")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(result).NotTo(BeNil())

	dep := &appsv1.Deployment{}
	err = result.Get("coherence-operator", dep)
	g.Expect(err).NotTo(HaveOccurred())

	c := findContainer("manager", dep)
	g.Expect(c).NotTo(BeNil())

	expected := fmt.Sprintf("foo.com/%s:%s", helper.GetOperatorImageName(), helper.GetOperatorVersionEnvVar())
	g.Expect(c.Image).To(Equal(expected))
}

func TestSetOperatorImageNamePart(t *testing.T) {
	g := NewGomegaWithT(t)
	result, err := helmInstall("--set", "image.name=foo")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(result).NotTo(BeNil())

	dep := &appsv1.Deployment{}
	err = result.Get("coherence-operator", dep)
	g.Expect(err).NotTo(HaveOccurred())

	c := findContainer("manager", dep)
	g.Expect(c).NotTo(BeNil())

	expected := fmt.Sprintf("%s/foo:%s", helper.GetOperatorImageRegistry(), helper.GetOperatorVersionEnvVar())
	g.Expect(c.Image).To(Equal(expected))
}

func TestSetOperatorImageTag(t *testing.T) {
	g := NewGomegaWithT(t)
	result, err := helmInstall("--set", "image.tag=1.0.0")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(result).NotTo(BeNil())

	dep := &appsv1.Deployment{}
	err = result.Get("coherence-operator", dep)
	g.Expect(err).NotTo(HaveOccurred())

	c := findContainer("manager", dep)
	g.Expect(c).NotTo(BeNil())

	expected := fmt.Sprintf("%s/%s:1.0.0", helper.GetOperatorImageRegistry(), helper.GetOperatorImageName())
	g.Expect(c.Image).To(Equal(expected))
}

func TestSetDefaultCoherenceImageInFull(t *testing.T) {
	g := NewGomegaWithT(t)
	result, err := helmInstall("--set", "defaultCoherenceImage=foo.com/bar:1.0.0")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(result).NotTo(BeNil())

	dep := &appsv1.Deployment{}
	err = result.Get("coherence-operator", dep)
	g.Expect(err).NotTo(HaveOccurred())

	c := findContainer("manager", dep)
	g.Expect(c).NotTo(BeNil())

	g.Expect(c.Env).NotTo(BeNil())
	g.Expect(c.Env).To(matchers.HaveEnvVar(corev1.EnvVar{Name: operator.EnvVarCoherenceImage, Value: "foo.com/bar:1.0.0"}))
}

func TestSetDefaultCoherenceImageRegistry(t *testing.T) {
	g := NewGomegaWithT(t)
	result, err := helmInstall("--set", "defaultCoherenceImage.registry=foo.com")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(result).NotTo(BeNil())

	dep := &appsv1.Deployment{}
	err = result.Get("coherence-operator", dep)
	g.Expect(err).NotTo(HaveOccurred())

	c := findContainer("manager", dep)
	g.Expect(c).NotTo(BeNil())

	expected := fmt.Sprintf("foo.com/%s:%s", helper.GetDefaultCoherenceImageName(), helper.GetDefaultCoherenceImageTag())
	g.Expect(c.Env).NotTo(BeNil())
	g.Expect(c.Env).To(matchers.HaveEnvVar(corev1.EnvVar{Name: operator.EnvVarCoherenceImage, Value: expected}))
}

func TestSetDefaultCoherenceImageNamePart(t *testing.T) {
	g := NewGomegaWithT(t)
	result, err := helmInstall("--set", "defaultCoherenceImage.name=foo")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(result).NotTo(BeNil())

	dep := &appsv1.Deployment{}
	err = result.Get("coherence-operator", dep)
	g.Expect(err).NotTo(HaveOccurred())

	c := findContainer("manager", dep)
	g.Expect(c).NotTo(BeNil())

	expected := fmt.Sprintf("%s/foo:%s", helper.GetDefaultCoherenceImageRegistry(), helper.GetDefaultCoherenceImageTag())
	g.Expect(c.Env).NotTo(BeNil())
	g.Expect(c.Env).To(matchers.HaveEnvVar(corev1.EnvVar{Name: operator.EnvVarCoherenceImage, Value: expected}))
}

func TestSetDefaultCoherenceImageTag(t *testing.T) {
	g := NewGomegaWithT(t)
	result, err := helmInstall("--set", "defaultCoherenceImage.tag=1.0.0")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(result).NotTo(BeNil())

	dep := &appsv1.Deployment{}
	err = result.Get("coherence-operator", dep)
	g.Expect(err).NotTo(HaveOccurred())

	c := findContainer("manager", dep)
	g.Expect(c).NotTo(BeNil())

	expected := fmt.Sprintf("%s/%s:1.0.0", helper.GetDefaultCoherenceImageRegistry(), helper.GetDefaultCoherenceImageName())
	g.Expect(c.Env).NotTo(BeNil())
	g.Expect(c.Env).To(matchers.HaveEnvVar(corev1.EnvVar{Name: operator.EnvVarCoherenceImage, Value: expected}))
}

func TestBasicHelmInstall(t *testing.T) {
	g := NewGomegaWithT(t)
	cmd, err := createHelmCommand()
	g.Expect(err).NotTo(HaveOccurred())
	AssertHelmInstallWithSubTest(t, "basic", cmd, g, AssertThreeReplicas)
}

func TestSetReplicas(t *testing.T) {
	g := NewGomegaWithT(t)
	cmd, err := createHelmCommand("--set", "replicas=1")
	g.Expect(err).NotTo(HaveOccurred())
	AssertHelmInstallWithSubTest(t, "basic", cmd, g, AssertSingleReplica)
}

func TestSetResources(t *testing.T) {
	g := NewGomegaWithT(t)
	cmd, err := createHelmCommand("--set", "replicas=1", "--set", "resources.requests.cpu=250m",
		"--set", "resources.requests.memory=64Mi", "--set", "resources.limits.cpu=512m",
		"--set", "resources.limits.memory=128Mi")
	g.Expect(err).NotTo(HaveOccurred())
	AssertHelmInstallWithSubTest(t, "basic", cmd, g, AssertResources)
}

func TestSetNonRootUser(t *testing.T) {
	g := NewGomegaWithT(t)
	cmd, err := createHelmCommand("--set", "securityContext.runAsNonRoot=true",
		"--set", "securityContext.runAsUser=1000")

	g.Expect(err).NotTo(HaveOccurred())
	AssertHelmInstallWithSubTest(t, "basic", cmd, g, AssertThreeReplicas)
}

func TestSetAdditionalPodLabel(t *testing.T) {
	g := NewGomegaWithT(t)
	result, err := helmInstall("--set", "labels.foo=bar")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(result).NotTo(BeNil())

	dep := &appsv1.Deployment{}
	err = result.Get("coherence-operator", dep)
	g.Expect(err).NotTo(HaveOccurred())

	labels := dep.Spec.Template.Labels
	actual, found := labels["control-plane"]
	g.Expect(found).To(BeTrue())
	g.Expect(actual).To(Equal("coherence"))
	actual, found = labels["foo"]
	g.Expect(found).To(BeTrue())
	g.Expect(actual).To(Equal("bar"))

	labels = dep.Labels
	_, found = labels["foo"]
	g.Expect(found).To(BeFalse())
}

func TestSetAdditionalPodLabels(t *testing.T) {
	g := NewGomegaWithT(t)
	result, err := helmInstall("--set", "labels.one=value-one", "--set", "labels.two=value-two")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(result).NotTo(BeNil())

	dep := &appsv1.Deployment{}
	err = result.Get("coherence-operator", dep)
	g.Expect(err).NotTo(HaveOccurred())

	labels := dep.Spec.Template.Labels
	actual, found := labels["control-plane"]
	g.Expect(found).To(BeTrue())
	g.Expect(actual).To(Equal("coherence"))
	actual, found = labels["one"]
	g.Expect(found).To(BeTrue())
	g.Expect(actual).To(Equal("value-one"))
	actual, found = labels["two"]
	g.Expect(found).To(BeTrue())
	g.Expect(actual).To(Equal("value-two"))

	labels = dep.Labels
	_, found = labels["one"]
	g.Expect(found).To(BeFalse())
	_, found = labels["two"]
	g.Expect(found).To(BeFalse())
}

func TestSetAdditionalPodAnnotation(t *testing.T) {
	g := NewGomegaWithT(t)
	result, err := helmInstall("--set", "annotations.foo=bar")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(result).NotTo(BeNil())

	dep := &appsv1.Deployment{}
	err = result.Get("coherence-operator", dep)
	g.Expect(err).NotTo(HaveOccurred())

	annotations := dep.Spec.Template.Annotations
	g.Expect(len(annotations)).To(Equal(1))
	actual, found := annotations["foo"]
	g.Expect(found).To(BeTrue())
	g.Expect(actual).To(Equal("bar"))

	annotations = dep.Annotations
	g.Expect(len(annotations)).To(BeZero())
}

func TestSetAdditionalPodAnnotations(t *testing.T) {
	g := NewGomegaWithT(t)
	result, err := helmInstall("--set", "annotations.one=value-one", "--set", "annotations.two=value-two")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(result).NotTo(BeNil())

	dep := &appsv1.Deployment{}
	err = result.Get("coherence-operator", dep)
	g.Expect(err).NotTo(HaveOccurred())

	annotations := dep.Spec.Template.Annotations
	g.Expect(len(annotations)).To(Equal(2))
	actual, found := annotations["one"]
	g.Expect(found).To(BeTrue())
	g.Expect(actual).To(Equal("value-one"))
	actual, found = annotations["two"]
	g.Expect(found).To(BeTrue())
	g.Expect(actual).To(Equal("value-two"))

	annotations = dep.Annotations
	g.Expect(len(annotations)).To(BeZero())
}

func TestSetAdditionalDeploymentLabel(t *testing.T) {
	g := NewGomegaWithT(t)
	result, err := helmInstall("--set", "deploymentLabels.foo=bar")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(result).NotTo(BeNil())

	dep := &appsv1.Deployment{}
	err = result.Get("coherence-operator", dep)
	g.Expect(err).NotTo(HaveOccurred())

	labels := dep.Labels
	actual, found := labels["control-plane"]
	g.Expect(found).To(BeTrue())
	g.Expect(actual).To(Equal("coherence"))
	actual, found = labels["foo"]
	g.Expect(found).To(BeTrue())
	g.Expect(actual).To(Equal("bar"))

	labels = dep.Spec.Template.Labels
	_, found = labels["foo"]
	g.Expect(found).To(BeFalse())
}

func TestSetAdditionalDeploymentLabels(t *testing.T) {
	g := NewGomegaWithT(t)
	result, err := helmInstall("--set", "deploymentLabels.one=value-one", "--set", "deploymentLabels.two=value-two")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(result).NotTo(BeNil())

	dep := &appsv1.Deployment{}
	err = result.Get("coherence-operator", dep)
	g.Expect(err).NotTo(HaveOccurred())

	labels := dep.Labels
	actual, found := labels["control-plane"]
	g.Expect(found).To(BeTrue())
	g.Expect(actual).To(Equal("coherence"))
	actual, found = labels["one"]
	g.Expect(found).To(BeTrue())
	g.Expect(actual).To(Equal("value-one"))
	actual, found = labels["two"]
	g.Expect(found).To(BeTrue())
	g.Expect(actual).To(Equal("value-two"))

	labels = dep.Spec.Template.Labels
	_, found = labels["one"]
	g.Expect(found).To(BeFalse())
	_, found = labels["two"]
	g.Expect(found).To(BeFalse())
}

func TestSetAdditionalDeploymentAnnotation(t *testing.T) {
	g := NewGomegaWithT(t)
	result, err := helmInstall("--set", "deploymentAnnotations.foo=bar")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(result).NotTo(BeNil())

	dep := &appsv1.Deployment{}
	err = result.Get("coherence-operator", dep)
	g.Expect(err).NotTo(HaveOccurred())

	annotations := dep.Annotations
	g.Expect(len(annotations)).To(Equal(1))
	actual, found := annotations["foo"]
	g.Expect(found).To(BeTrue())
	g.Expect(actual).To(Equal("bar"))

	annotations = dep.Spec.Template.Annotations
	g.Expect(len(annotations)).To(BeZero())
}

func TestSetAdditionalDeploymentAnnotations(t *testing.T) {
	g := NewGomegaWithT(t)
	result, err := helmInstall("--set", "deploymentAnnotations.one=value-one", "--set", "deploymentAnnotations.two=value-two")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(result).NotTo(BeNil())

	dep := &appsv1.Deployment{}
	err = result.Get("coherence-operator", dep)
	g.Expect(err).NotTo(HaveOccurred())

	annotations := dep.Annotations
	g.Expect(len(annotations)).To(Equal(2))
	actual, found := annotations["one"]
	g.Expect(found).To(BeTrue())
	g.Expect(actual).To(Equal("value-one"))
	actual, found = annotations["two"]
	g.Expect(found).To(BeTrue())
	g.Expect(actual).To(Equal("value-two"))

	annotations = dep.Spec.Template.Annotations
	g.Expect(len(annotations)).To(BeZero())
}

func TestGlobalLabelsAndAnnotations(t *testing.T) {
	g := NewGomegaWithT(t)

	cmd, err := createHelmCommand("--set", "globalLabels.one=label-one",
		"--set", "globalLabels.two=label-two",
		"--set", "globalAnnotations.three=annotation-three",
		"--set", "globalAnnotations.four=annotation-four")
	g.Expect(err).NotTo(HaveOccurred())
	AssertHelmInstallWithStatefulSetSubTest(t, "basic", cmd, g, AssertLabelsAndAnnotations)
}

func TestGlobalLabelsOnOperatorResources(t *testing.T) {
	g := NewGomegaWithT(t)

	result, err := helmInstall("--set", "globalLabels.one=label-one",
		"--set", "globalLabels.two=label-two")

	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(result).NotTo(BeNil())

	dep := &appsv1.Deployment{}
	err = result.Get("coherence-operator", dep)
	g.Expect(err).NotTo(HaveOccurred())

	labels := dep.Labels
	actual, found := labels["one"]
	g.Expect(found).To(BeTrue())
	g.Expect(actual).To(Equal("label-one"))
	actual, found = labels["two"]
	g.Expect(found).To(BeTrue())
	g.Expect(actual).To(Equal("label-two"))

	labels = dep.Spec.Template.Labels
	actual, found = labels["one"]
	g.Expect(found).To(BeTrue())
	g.Expect(actual).To(Equal("label-one"))
	actual, found = labels["two"]
	g.Expect(found).To(BeTrue())
	g.Expect(actual).To(Equal("label-two"))

	svc := &corev1.Service{}
	err = result.Get("coherence-operator-rest", svc)
	g.Expect(err).NotTo(HaveOccurred())

	labels = svc.Labels
	actual, found = labels["one"]
	g.Expect(found).To(BeTrue())
	g.Expect(actual).To(Equal("label-one"))
	actual, found = labels["two"]
	g.Expect(found).To(BeTrue())
	g.Expect(actual).To(Equal("label-two"))
}

func TestGlobalAnnotationsOnOperatorResources(t *testing.T) {
	g := NewGomegaWithT(t)

	result, err := helmInstall("--set", "globalAnnotations.one=annotation-one",
		"--set", "globalAnnotations.two=annotation-two")

	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(result).NotTo(BeNil())

	dep := &appsv1.Deployment{}
	err = result.Get("coherence-operator", dep)
	g.Expect(err).NotTo(HaveOccurred())

	annotations := dep.Annotations
	actual, found := annotations["one"]
	g.Expect(found).To(BeTrue())
	g.Expect(actual).To(Equal("annotation-one"))
	actual, found = annotations["two"]
	g.Expect(found).To(BeTrue())
	g.Expect(actual).To(Equal("annotation-two"))

	annotations = dep.Spec.Template.Annotations
	actual, found = annotations["one"]
	g.Expect(found).To(BeTrue())
	g.Expect(actual).To(Equal("annotation-one"))
	actual, found = annotations["two"]
	g.Expect(found).To(BeTrue())
	g.Expect(actual).To(Equal("annotation-two"))

	svc := &corev1.Service{}
	err = result.Get("coherence-operator-webhook", svc)
	g.Expect(err).NotTo(HaveOccurred())

	annotations = svc.Annotations
	actual, found = annotations["one"]
	g.Expect(found).To(BeTrue())
	g.Expect(actual).To(Equal("annotation-one"))
	actual, found = annotations["two"]
	g.Expect(found).To(BeTrue())
	g.Expect(actual).To(Equal("annotation-two"))

	svc = &corev1.Service{}
	err = result.Get("coherence-operator-rest", svc)
	g.Expect(err).NotTo(HaveOccurred())

	annotations = svc.Annotations
	actual, found = annotations["one"]
	g.Expect(found).To(BeTrue())
	g.Expect(actual).To(Equal("annotation-one"))
	actual, found = annotations["two"]
	g.Expect(found).To(BeTrue())
	g.Expect(actual).To(Equal("annotation-two"))

	sec := &corev1.Secret{}
	err = result.Get("coherence-webhook-server-cert", sec)
	g.Expect(err).NotTo(HaveOccurred())

	annotations = sec.Annotations
	actual, found = annotations["one"]
	g.Expect(found).To(BeTrue())
	g.Expect(actual).To(Equal("annotation-one"))
	actual, found = annotations["two"]
	g.Expect(found).To(BeTrue())
	g.Expect(actual).To(Equal("annotation-two"))
}

func AssertLabelsAndAnnotations(t *testing.T, g *GomegaWithT, _ *coh.Coherence, sts *appsv1.StatefulSet) {
	g.Expect(sts.Labels).NotTo(BeNil())
	g.Expect(sts.Labels["one"]).To(Equal("label-one"))
	g.Expect(sts.Labels["two"]).To(Equal("label-two"))

	g.Expect(sts.Annotations).NotTo(BeNil())
	g.Expect(sts.Annotations["three"]).To(Equal("annotation-three"))
	g.Expect(sts.Annotations["four"]).To(Equal("annotation-four"))
}

func AssertResources() error {
	ns := helper.GetTestNamespace()
	pods, err := helper.ListOperatorPods(testContext, ns)
	if err != nil {
		return err
	}
	if len(pods) == 0 {
		return fmt.Errorf("expected at least one Coherence Operator Pod but found zero")
	}
	resources := pods[0].Spec.Containers[0].Resources

	var qty *resource.Quantity

	qty = resources.Requests.Name("cpu", resource.BinarySI)
	if qty == nil {
		return fmt.Errorf("expected a cpu requests quantity")
	}
	if qty.String() != "250m" {
		return fmt.Errorf("expected a cpu requests of 250m")
	}

	qty = resources.Requests.Name("memory", resource.BinarySI)
	if qty == nil {
		return fmt.Errorf("expected a memory requests quantity")
	}
	if qty.String() != "64Mi" {
		return fmt.Errorf("expected a memory requests of 64Mi")
	}

	qty = resources.Limits.Name("cpu", resource.BinarySI)
	if qty == nil {
		return fmt.Errorf("expected a cpu limits quantity")
	}
	if qty.String() != "512m" {
		return fmt.Errorf("expected a cpu limit of 512m")
	}

	qty = resources.Limits.Name("memory", resource.BinarySI)
	if qty == nil {
		return fmt.Errorf("expected a memory limits quantity")
	}
	if qty.String() != "128Mi" {
		return fmt.Errorf("expected a memory limit of 128Mi")
	}

	return nil
}

func TestHelmInstallWithServiceAccountName(t *testing.T) {
	g := NewGomegaWithT(t)
	cmd, err := createHelmCommand("--set", "serviceAccountName=test-operator-account")
	g.Expect(err).NotTo(HaveOccurred())
	AssertHelmInstallWithSubTest(t, "account", cmd, g, emptySubTest)
}

func TestHelmInstallWithoutClusterRoles(t *testing.T) {
	g := NewGomegaWithT(t)
	cmd, err := createHelmCommand("--set", "clusterRoles=false")
	g.Expect(err).NotTo(HaveOccurred())
	AssertHelmInstallWithSubTest(t, "no-roles", cmd, g, AssertNoClusterRoles)
}

func TestHelmInstallWithoutClusterRolesWithNodeRole(t *testing.T) {
	g := NewGomegaWithT(t)
	cmd, err := createHelmCommand("--set", "clusterRoles=false", "--set", "nodeRoles=true")
	g.Expect(err).NotTo(HaveOccurred())
	AssertHelmInstallWithSubTest(t, "node-role", cmd, g, AssertOnlyNodeClusterRoles)
}

func createHelmCommand(args ...string) (*exec.Cmd, error) {
	ns := helper.GetTestNamespace()

	chart, err := helper.FindOperatorHelmChartDir()
	if err != nil {
		return nil, err
	}

	argList := []string{"install",
		"--set", "image=" + helper.GetOperatorImage(),
		"--set", "defaultCoherenceUtilsImage=" + helper.GetOperatorImage()}

	argList = append(argList, args...)

	argList = append(argList, args...)
	argList = append(argList, "--namespace", ns, "--wait", "operator", chart)

	cmd := exec.Command("helm", argList...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd, nil
}

type SubTest func() error

type StatefulSetSubTest func(*testing.T, *GomegaWithT, *coh.Coherence, *appsv1.StatefulSet)

type SubTestRunner struct {
	Test SubTest
}

func (in SubTestRunner) run(_ *testing.T, g *GomegaWithT, _ *coh.Coherence, _ *appsv1.StatefulSet) {
	err := in.Test()
	g.Expect(err).NotTo(HaveOccurred())
}

var emptySubTest = func() error {
	return nil
}

func AssertNoClusterRoles() error {
	return AssertRBAC(false)
}

func AssertOnlyNodeClusterRoles() error {
	return AssertRBAC(true)
}

func AssertSingleReplica() error {
	ns := helper.GetTestNamespace()
	pods, err := helper.ListOperatorPods(testContext, ns)
	if err != nil {
		return err
	}
	count := len(pods)
	if count != 1 {
		return fmt.Errorf("expected a single Coherence Operator Pod but found %d", count)
	}
	return nil
}

func AssertThreeReplicas() error {
	ns := helper.GetTestNamespace()
	pods, err := helper.ListOperatorPods(testContext, ns)
	if err != nil {
		return err
	}
	count := len(pods)
	if count != 3 {
		return fmt.Errorf("expected three Coherence Operator Pods but found %d", count)
	}
	return nil
}

func AssertRBAC(allowNode bool) error {
	rbacClient := testContext.KubeClient.RbacV1()

	crList, err := rbacClient.ClusterRoles().List(goctx.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, cr := range crList.Items {
		if strings.HasPrefix(strings.ToLower(cr.Name), "coherence-operator") {
			if !allowNode || cr.Name != "coherence-operator-node-viewer" {
				return fmt.Errorf("no Coherence ClusterRole should exist but found %s", cr.Name)
			}
		}
	}

	crbList, err := rbacClient.ClusterRoleBindings().List(goctx.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, crb := range crbList.Items {
		if strings.HasPrefix(strings.ToLower(crb.Name), "coherence-operator") {
			if !allowNode || crb.Name != "coherence-operator-node-viewer" {
				return fmt.Errorf("no Coherence ClusterRoleBinding should exist but found %s", crb.Name)
			}
		}
	}

	ns := helper.GetTestNamespace()
	_, err = rbacClient.Roles(ns).Get(goctx.TODO(), "coherence-operator", metav1.GetOptions{})
	if err != nil {
		return err
	}

	_, err = rbacClient.RoleBindings(ns).Get(goctx.TODO(), "coherence-operator", metav1.GetOptions{})
	if err != nil {
		return err
	}

	return nil
}

func AssertHelmInstallWithSubTest(t *testing.T, id string, cmd *exec.Cmd, g *GomegaWithT, test SubTest) {
	runner := SubTestRunner{
		Test: test,
	}
	AssertHelmInstallWithStatefulSetSubTest(t, id, cmd, g, runner.run)
}

func AssertHelmInstallWithStatefulSetSubTest(t *testing.T, id string, cmd *exec.Cmd, g *GomegaWithT, test StatefulSetSubTest) {
	ns := helper.GetTestNamespace()

	t.Cleanup(func() {
		if t.Failed() {
			helper.DumpOperatorLogs(t, testContext)
		}
		Cleanup(ns, "operator")
	})

	t.Logf("Asserting Helm install. Removing Webhooks")
	err := RemoveWebHook()
	g.Expect(err).NotTo(HaveOccurred())
	t.Logf("Asserting Helm install. Removing RBAC")
	err = RemoveRBAC()
	g.Expect(err).NotTo(HaveOccurred())

	t.Logf("Asserting Helm install. Performing Helm install")
	err = cmd.Run()
	g.Expect(err).NotTo(HaveOccurred())

	t.Logf("Asserting Helm install. Ensure Operator Pod is ready")
	pods, err := helper.ListOperatorPods(testContext, ns)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(len(pods)).NotTo(Equal(0))

	pod := pods[0]
	err = helper.WaitForPodReady(testContext, pod.Namespace, pod.Name, 10*time.Second, 5*time.Minute)
	g.Expect(err).NotTo(HaveOccurred())

	t.Logf("Asserting Helm install. Deploying Coherence resource")
	deployment, err := helper.NewSingleCoherenceFromYaml(ns, "coherence.yaml")
	g.Expect(err).NotTo(HaveOccurred())

	name := deployment.GetName()
	deployment.SetName(name + "-" + id)

	defer deleteCoherence(t, &deployment)

	err = testContext.Client.Create(goctx.TODO(), &deployment)
	g.Expect(err).NotTo(HaveOccurred())

	var sts *appsv1.StatefulSet
	sts, err = helper.WaitForStatefulSetForDeployment(testContext, ns, &deployment, helper.RetryInterval, helper.Timeout)
	g.Expect(err).NotTo(HaveOccurred())

	test(t, g, &deployment, sts)
}

func deleteCoherence(t *testing.T, d *coh.Coherence) {
	if err := testContext.Client.Delete(goctx.TODO(), d); err != nil {
		t.Logf("Error deleting Coherence deployment %s - %s", d.GetNamespacedName(), err.Error())
	}
}

func Cleanup(namespace, name string) {
	_ = helper.WaitForCoherenceCleanup(testContext, namespace)

	cmd := exec.Command("helm", "uninstall", "--namespace", namespace, name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()
	_ = helper.WaitForDeleteOfPodsWithSelector(testContext, namespace, "control-plane=coherence", helper.RetryInterval, helper.Timeout)
}

// Remove the web-hooks that the Operator install creates to
// ensure that nothing is left from a previous test.
func RemoveWebHook() error {
	// DefaultValidatingWebhookName
	client := testContext.KubeClient.AdmissionregistrationV1()

	err := client.MutatingWebhookConfigurations().Delete(goctx.TODO(), operator.DefaultMutatingWebhookName, metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return err
	}

	err = client.ValidatingWebhookConfigurations().Delete(goctx.TODO(), operator.DefaultValidatingWebhookName, metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return err
	}

	return nil
}

// Remove the CRDs that the Operator install creates.
func RemoveCRDs() error {
	cohCrd := crdv1.CustomResourceDefinition{}
	cohCrd.Name = "coherence.coherence.oracle.com"
	err := testContext.Client.Delete(goctx.TODO(), &cohCrd)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}

	cohJobCrd := crdv1.CustomResourceDefinition{}
	cohJobCrd.Name = "coherencejob.coherence.oracle.com"
	err = testContext.Client.Delete(goctx.TODO(), &cohJobCrd)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}

	return nil
}

// RemoveRBAC removes all the RBAC rules that the Operator install creates to
// ensure that nothing is left from a previous test.
func RemoveRBAC() error {
	var err error
	rbacClient := testContext.KubeClient.RbacV1()
	clusterRolesClient := rbacClient.ClusterRoles()
	clusterRoleBindingsClient := rbacClient.ClusterRoleBindings()

	crList, err := clusterRolesClient.List(goctx.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, cr := range crList.Items {
		if strings.HasPrefix(strings.ToLower(cr.Name), "coherence-operator") {
			if err := clusterRolesClient.Delete(goctx.TODO(), cr.Name, metav1.DeleteOptions{}); err != nil {
				return err
			}
		}
	}

	crbList, err := clusterRoleBindingsClient.List(goctx.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, crb := range crbList.Items {
		if strings.HasPrefix(strings.ToLower(crb.Name), "coherence-operator") {
			if err := clusterRoleBindingsClient.Delete(goctx.TODO(), crb.Name, metav1.DeleteOptions{}); err != nil {
				return err
			}
		}
	}

	ns := helper.GetTestNamespace()
	rolesClient := rbacClient.Roles(ns)
	roleBindingsClient := rbacClient.RoleBindings(ns)

	if role, err := rolesClient.Get(goctx.TODO(), "coherence-operator", metav1.GetOptions{}); err == nil {
		if err := rolesClient.Delete(goctx.TODO(), role.Name, metav1.DeleteOptions{}); err != nil {
			return err
		}
	}

	if roleBinding, err := roleBindingsClient.Get(goctx.TODO(), "coherence-operator", metav1.GetOptions{}); err == nil {
		if err := roleBindingsClient.Delete(goctx.TODO(), roleBinding.Name, metav1.DeleteOptions{}); err != nil {
			return err
		}
	}

	return nil
}
