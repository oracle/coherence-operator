/*
 * Copyright (c) 2026, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1

import (
	"testing"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
)

func TestUpdateFromJobClearsProbeErrorAfterLaterReadyTransition(t *testing.T) {
	firstReady := metav1.NewTime(time.Now())
	secondReady := metav1.NewTime(firstReady.Add(time.Second))
	secondProbe := metav1.NewTime(secondReady.Add(time.Second))
	status := CoherenceResourceStatus{JobProbes: []CoherenceJobProbeStatus{{
		Pod:           "job-0",
		LastReadyTime: &firstReady,
		LastProbeTime: &firstReady,
		Success:       ptr.To(false),
		Error:         ptr.To("HTTP probe failed: transport timeout"),
	}}}

	updated := status.UpdateFromJob(&CoherenceJob{}, nil, []CoherenceJobProbeStatus{{
		Pod:           "job-0",
		LastReadyTime: &secondReady,
		LastProbeTime: &secondProbe,
		Success:       ptr.To(true),
	}})
	if !updated {
		t.Fatal("successful attempt after a later Ready transition did not update status")
	}
	if len(status.JobProbes) != 1 {
		t.Fatalf("expected the existing probe status to be updated, got %d entries", len(status.JobProbes))
	}
	probe := status.MaybeFindJobProbeStatus("job-0")
	if probe == nil || probe.Success == nil || !*probe.Success || probe.Error != nil {
		t.Fatalf("unexpected probe status: %+v", probe)
	}
}
