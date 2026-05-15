/*
 * Copyright (c) 2026, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package coherenceoperator

import (
	"os"
	"testing"

	semver "github.com/Masterminds/semver/v3"
	"github.com/ghodss/yaml"
)

type chartMetadata struct {
	KubeVersion string `json:"kubeVersion"`
}

func TestKubeVersionConstraint(t *testing.T) {
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

	constraint, err := semver.NewConstraint(chart.KubeVersion)
	if err != nil {
		t.Fatalf("failed to parse kubeVersion constraint %q: %v", chart.KubeVersion, err)
	}

	tests := []struct {
		name    string
		version string
		want    bool
	}{
		{
			name:    "accepts supported vendor suffixed Kubernetes versions",
			version: "1.32.3-vke.8",
			want:    true,
		},
		{
			name:    "accepts the documented support floor",
			version: "1.29.0",
			want:    true,
		},
		{
			name:    "rejects versions below the documented support floor",
			version: "1.28.0",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			version, err := semver.NewVersion(tt.version)
			if err != nil {
				t.Fatalf("failed to parse Kubernetes version %q: %v", tt.version, err)
			}

			// Exercise the chart's actual constraint so future edits preserve both the 1.29 floor and vendor suffix support.
			if got := constraint.Check(version); got != tt.want {
				t.Fatalf("kubeVersion %q match for %q = %v, want %v", chart.KubeVersion, tt.version, got, tt.want)
			}
		})
	}
}
