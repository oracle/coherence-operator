/*
 * Copyright (c) 2026, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package coherenceoperator

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"testing"

	semver "github.com/Masterminds/semver/v3"
	"github.com/ghodss/yaml"
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
