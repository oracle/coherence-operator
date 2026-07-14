/*
 * Copyright (c) 2026, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package coherenceoperator

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"testing"

	semver "github.com/Masterminds/semver/v3"
	"github.com/ghodss/yaml"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	// Paths are relative to helm-charts/coherence-operator; update this if the
	// config/manifests/bases OLM metadata layout moves.
	csvMetadataPath = "../../config/manifests/bases/coherence-operator.clusterserviceversion.yaml"
	// Paths are relative to helm-charts/coherence-operator; update this if the
	// repository root Makefile moves.
	makefilePath = "../../Makefile"
)

type chartMetadata struct {
	KubeVersion string `json:"kubeVersion"`
}

type csvMetadata struct {
	Spec struct {
		MinKubeVersion string `json:"minKubeVersion"`
	} `json:"spec"`
}

type chartValuesSchema struct {
	Properties struct {
		RestService struct {
			Properties struct {
				IPFamilies struct {
					MaxItems    int  `json:"maxItems"`
					UniqueItems bool `json:"uniqueItems"`
				} `json:"ipFamilies"`
			} `json:"properties"`
		} `json:"restService"`
	} `json:"properties"`
}

type kubernetesVersionFloor struct {
	Major int
	Minor int
}

func TestKubeVersionConstraint(t *testing.T) {
	chart := readChartMetadata(t)
	minimum, constraint := parseChartKubernetesConstraint(t, chart.KubeVersion)

	tests := []struct {
		name    string
		version string
		want    bool
	}{
		{
			name:    "accepts supported vendor suffixed Kubernetes versions",
			version: minimum.vendorSuffixedPatchString(),
			want:    true,
		},
		{
			name:    "accepts vendor suffixed Kubernetes versions at the support floor",
			version: minimum.patchString() + "-eks.1",
			want:    true,
		},
		{
			name:    "accepts the documented support floor",
			version: minimum.patchString(),
			want:    true,
		},
		{
			name:    "rejects versions below the documented support floor",
			version: minimum.previousMinorPatchString(t),
			want:    false,
		},
		{
			name:    "rejects vendor suffixed Kubernetes versions below the support floor",
			version: minimum.previousMinorPatchString(t) + "-gke.123",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			version, err := semver.NewVersion(tt.version)
			if err != nil {
				t.Fatalf("failed to parse Kubernetes version %q: %v", tt.version, err)
			}

			// Exercise the chart's actual constraint so future edits preserve both the declared
			// floor and vendor suffix handling at the boundary.
			if got := constraint.Check(version); got != tt.want {
				t.Fatalf("kubeVersion %q match for %q = %v, want %v", chart.KubeVersion, tt.version, got, tt.want)
			}
		})
	}
}

func TestKubernetesMinimumVersionsAreAligned(t *testing.T) {
	chartFloor, _ := parseChartKubernetesConstraint(t, readChartMetadata(t).KubeVersion)
	csvFloor := readCSVKubernetesFloor(t)
	makefileFloor := readMakefileKubernetesFloor(t)

	// These files gate different install paths, so compare the normalized floor to catch
	// drift before Helm, OLM, and validation builds start enforcing different versions.
	if chartFloor != csvFloor || chartFloor != makefileFloor {
		t.Fatalf("Kubernetes minimum versions differ: Chart.yaml=%s, CSV=%s, Makefile=%s",
			chartFloor.minorString(), csvFloor.minorString(), makefileFloor.minorString())
	}
}

func TestRestServiceDefaultRendering(t *testing.T) {
	service := renderRestService(t)
	if service.Spec.IPFamilyPolicy != nil {
		t.Fatalf("default ipFamilyPolicy = %q, want omitted", *service.Spec.IPFamilyPolicy)
	}
	if service.Spec.IPFamilies != nil {
		t.Fatalf("default ipFamilies = %v, want omitted", service.Spec.IPFamilies)
	}

	if service.Namespace != "default" || service.Annotations != nil {
		t.Fatalf("REST Service metadata changed unexpectedly: namespace=%q annotations=%v", service.Namespace, service.Annotations)
	}
	wantLabels := map[string]string{
		"control-plane":                "coherence",
		"app.kubernetes.io/name":       "coherence-operator",
		"app.kubernetes.io/instance":   "coherence-operator-rest",
		"app.kubernetes.io/version":    "${VERSION}",
		"app.kubernetes.io/component":  "rest",
		"app.kubernetes.io/part-of":    "coherence-operator",
		"app.kubernetes.io/managed-by": "helm",
	}
	if !reflect.DeepEqual(service.Labels, wantLabels) {
		t.Fatalf("REST Service labels = %v, want %v", service.Labels, wantLabels)
	}
	wantPorts := []corev1.ServicePort{{
		Name:       "http-rest",
		Port:       8000,
		TargetPort: intstr.FromInt32(8000),
	}}
	if !reflect.DeepEqual(service.Spec.Ports, wantPorts) {
		t.Fatalf("REST Service ports = %v, want %v", service.Spec.Ports, wantPorts)
	}
	wantSelector := map[string]string{
		"app.kubernetes.io/name":      "coherence-operator",
		"app.kubernetes.io/instance":  "coherence-operator-manager",
		"app.kubernetes.io/version":   "${VERSION}",
		"app.kubernetes.io/component": "manager",
	}
	if !reflect.DeepEqual(service.Spec.Selector, wantSelector) {
		t.Fatalf("REST Service selector = %v, want %v", service.Spec.Selector, wantSelector)
	}
}

func TestRestServiceFieldsRenderIndependently(t *testing.T) {
	t.Run("policy only", func(t *testing.T) {
		service := renderRestService(t, "--set", "restService.ipFamilyPolicy=PreferDualStack")
		if service.Spec.IPFamilyPolicy == nil || *service.Spec.IPFamilyPolicy != corev1.IPFamilyPolicyPreferDualStack {
			t.Fatalf("ipFamilyPolicy = %v, want PreferDualStack", service.Spec.IPFamilyPolicy)
		}
		if service.Spec.IPFamilies != nil {
			t.Fatalf("ipFamilies = %v, want omitted", service.Spec.IPFamilies)
		}
	})

	t.Run("families only", func(t *testing.T) {
		service := renderRestService(t, "--set", "restService.ipFamilies={IPv6}")
		if service.Spec.IPFamilyPolicy != nil {
			t.Fatalf("ipFamilyPolicy = %v, want omitted", service.Spec.IPFamilyPolicy)
		}
		assertIPFamilies(t, service, corev1.IPv6Protocol)
	})
}

func TestRestServiceIPFamilies(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		families []corev1.IPFamily
	}{
		{name: "IPv4 only", value: "{IPv4}", families: []corev1.IPFamily{corev1.IPv4Protocol}},
		{name: "IPv6 only", value: "{IPv6}", families: []corev1.IPFamily{corev1.IPv6Protocol}},
		{name: "IPv4 then IPv6", value: "{IPv4,IPv6}", families: []corev1.IPFamily{corev1.IPv4Protocol, corev1.IPv6Protocol}},
		{name: "IPv6 then IPv4", value: "{IPv6,IPv4}", families: []corev1.IPFamily{corev1.IPv6Protocol, corev1.IPv4Protocol}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := renderRestService(t, "--set", "restService.ipFamilies="+tt.value)
			assertIPFamilies(t, service, tt.families...)
		})
	}
}

func TestRestServiceIPFamilyPolicies(t *testing.T) {
	policies := []corev1.IPFamilyPolicy{
		corev1.IPFamilyPolicySingleStack,
		corev1.IPFamilyPolicyPreferDualStack,
		corev1.IPFamilyPolicyRequireDualStack,
	}

	for _, policy := range policies {
		t.Run(string(policy), func(t *testing.T) {
			service := renderRestService(t, "--set", "restService.ipFamilyPolicy="+string(policy))
			if service.Spec.IPFamilyPolicy == nil || *service.Spec.IPFamilyPolicy != policy {
				t.Fatalf("ipFamilyPolicy = %v, want %s", service.Spec.IPFamilyPolicy, policy)
			}
		})
	}
}

func TestRestServiceSchemaValidation(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{name: "invalid policy", args: []string{"--set", "restService.ipFamilyPolicy=SometimesDualStack"}},
		{name: "invalid family", args: []string{"--set", "restService.ipFamilies={IPv5}"}},
		{name: "duplicate families", args: []string{"--set", "restService.ipFamilies={IPv4,IPv4}"}},
		{name: "scalar families", args: []string{"--set", "restService.ipFamilies=IPv4"}},
		{name: "more than two families", args: []string{"--set", "restService.ipFamilies={IPv4,IPv6,IPv4}"}},
		{name: "unknown setting", args: []string{"--set", "restService.ipFamilyPolcy=PreferDualStack"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := helmTemplate(t, tt.args...)
			if err == nil {
				t.Fatal("helm template succeeded, want schema validation error")
			}
		})
	}
}

func TestRestServiceSchemaListConstraints(t *testing.T) {
	data, err := os.ReadFile("values.schema.json")
	if err != nil {
		t.Fatalf("failed to read Helm values schema: %v", err)
	}

	var schema chartValuesSchema
	if err = yaml.Unmarshal(data, &schema); err != nil {
		t.Fatalf("failed to parse Helm values schema: %v", err)
	}

	ipFamilies := schema.Properties.RestService.Properties.IPFamilies
	if ipFamilies.MaxItems != 2 {
		t.Fatalf("ipFamilies maxItems = %d, want 2", ipFamilies.MaxItems)
	}
	if !ipFamilies.UniqueItems {
		t.Fatal("ipFamilies uniqueItems = false, want true")
	}
}

func renderRestService(t *testing.T, args ...string) *corev1.Service {
	t.Helper()

	output, err := helmTemplate(t, args...)
	if err != nil {
		t.Fatalf("helm template failed: %v\n%s", err, output)
	}

	for _, document := range strings.Split(output, "\n---") {
		var resource struct {
			Kind     string `json:"kind"`
			Metadata struct {
				Name string `json:"name"`
			} `json:"metadata"`
		}
		if err = yaml.Unmarshal([]byte(document), &resource); err != nil {
			t.Fatalf("failed to parse rendered manifest: %v\n%s", err, document)
		}
		if resource.Kind == "Service" && resource.Metadata.Name == "coherence-operator-rest" {
			service := &corev1.Service{}
			if err = yaml.Unmarshal([]byte(document), service); err != nil {
				t.Fatalf("failed to parse rendered REST Service: %v\n%s", err, document)
			}
			return service
		}
	}

	t.Fatalf("rendered chart does not contain coherence-operator-rest Service:\n%s", output)
	return nil
}

func assertIPFamilies(t *testing.T, service *corev1.Service, want ...corev1.IPFamily) {
	t.Helper()

	if len(service.Spec.IPFamilies) != len(want) {
		t.Fatalf("ipFamilies = %v, want %v", service.Spec.IPFamilies, want)
	}
	for i := range want {
		if service.Spec.IPFamilies[i] != want[i] {
			t.Fatalf("ipFamilies = %v, want %v", service.Spec.IPFamilies, want)
		}
	}
}

func helmTemplate(t *testing.T, args ...string) (string, error) {
	t.Helper()

	chart := copyTestChart(t)
	commandArgs := []string{"template", "operator", chart, "--show-only", "templates/deployment.yaml"}
	commandArgs = append(commandArgs, args...)
	cmd := exec.Command("helm", commandArgs...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func copyTestChart(t *testing.T) string {
	t.Helper()

	destination := filepath.Join(t.TempDir(), "coherence-operator")
	err := filepath.Walk(".", func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		target := filepath.Join(destination, path)
		if info.IsDir() {
			return os.MkdirAll(target, info.Mode())
		}

		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		if path == "Chart.yaml" {
			data = []byte(strings.ReplaceAll(string(data), "${VERSION}", "1.0.0"))
		}
		return os.WriteFile(target, data, info.Mode())
	})
	if err != nil {
		t.Fatalf("failed to prepare Helm chart for testing: %v", err)
	}
	return destination
}

func readChartMetadata(t *testing.T) chartMetadata {
	t.Helper()

	data, err := os.ReadFile("Chart.yaml")
	if err != nil {
		t.Fatalf("failed to read Helm chart metadata: %v", err)
	}

	var chart chartMetadata
	if err = yaml.Unmarshal(data, &chart); err != nil {
		t.Fatalf("failed to parse Helm chart metadata: %v", err)
	}
	if chart.KubeVersion == "" {
		t.Fatal("expected Helm chart metadata to declare kubeVersion")
	}

	return chart
}

func readCSVKubernetesFloor(t *testing.T) kubernetesVersionFloor {
	t.Helper()

	data, err := os.ReadFile(csvMetadataPath)
	if err != nil {
		t.Fatalf("failed to read OLM CSV metadata: %v", err)
	}

	var csv csvMetadata
	if err = yaml.Unmarshal(data, &csv); err != nil {
		t.Fatalf("failed to parse OLM CSV metadata: %v", err)
	}

	return parsePatchKubernetesFloor(t, csv.Spec.MinKubeVersion, "CSV minKubeVersion")
}

func readMakefileKubernetesFloor(t *testing.T) kubernetesVersionFloor {
	t.Helper()

	data, err := os.ReadFile(makefilePath)
	if err != nil {
		t.Fatalf("failed to read Makefile: %v", err)
	}

	re := regexp.MustCompile(`(?m)^[ \t]*(?:export[ \t]+)?KUBERNETES_MIN_VERSION[ \t]*(?::{1,2}=|\?=|=)[ \t]*(\d+)\.(\d+)(?:\.0)?[ \t]*(?:#.*)?\r?$`)
	matches := re.FindStringSubmatch(string(data))
	if matches == nil {
		t.Fatal("expected Makefile to declare KUBERNETES_MIN_VERSION with =, :=, ?=, or ::= and value <major>.<minor>[.0]")
	}

	return parseKubernetesFloor(t, matches[1], matches[2], "Makefile KUBERNETES_MIN_VERSION")
}

func parseChartKubernetesConstraint(t *testing.T, version string) (kubernetesVersionFloor, *semver.Constraints) {
	t.Helper()

	constraint, err := semver.NewConstraint(version)
	if err != nil {
		t.Fatalf("failed to parse Chart.yaml kubeVersion constraint %q: %v", version, err)
	}

	re := regexp.MustCompile(`^[ \t]*>=[ \t]*(\d+)\.(\d+)(?:\.0)?-0(?:[ \t]*(?:,|\|\||[ \t]+).*)?[ \t]*$`)
	matches := re.FindStringSubmatch(version)
	if matches == nil {
		t.Fatalf("expected Chart.yaml kubeVersion %q to include a >=<major>.<minor>[.0]-0 lower bound", version)
	}

	return parseKubernetesFloor(t, matches[1], matches[2], "Chart.yaml kubeVersion"), constraint
}

func parsePatchKubernetesFloor(t *testing.T, version, field string) kubernetesVersionFloor {
	t.Helper()

	re := regexp.MustCompile(`^(\d+)\.(\d+)\.0$`)
	matches := re.FindStringSubmatch(version)
	if matches == nil {
		t.Fatalf("expected %s %q to use <major>.<minor>.0", field, version)
	}

	return parseKubernetesFloor(t, matches[1], matches[2], field)
}

func parseKubernetesFloor(t *testing.T, major, minor, field string) kubernetesVersionFloor {
	t.Helper()

	majorValue, err := strconv.Atoi(major)
	if err != nil {
		t.Fatalf("failed to parse %s major version %q: %v", field, major, err)
	}

	minorValue, err := strconv.Atoi(minor)
	if err != nil {
		t.Fatalf("failed to parse %s minor version %q: %v", field, minor, err)
	}

	return kubernetesVersionFloor{
		Major: majorValue,
		Minor: minorValue,
	}
}

func (v kubernetesVersionFloor) minorString() string {
	return fmt.Sprintf("%d.%d", v.Major, v.Minor)
}

func (v kubernetesVersionFloor) patchString() string {
	return fmt.Sprintf("%s.0", v.minorString())
}

func (v kubernetesVersionFloor) vendorSuffixedPatchString() string {
	return fmt.Sprintf("%d.%d.3-vke.8", v.Major, v.Minor)
}

func (v kubernetesVersionFloor) previousMinorPatchString(t *testing.T) string {
	t.Helper()

	switch {
	case v.Minor > 0:
		return fmt.Sprintf("%d.%d.0", v.Major, v.Minor-1)
	case v.Major > 0:
		return fmt.Sprintf("%d.0.0", v.Major-1)
	default:
		t.Fatal("cannot compute previous minor for Kubernetes 0.0 floor")
		return ""
	}
}
