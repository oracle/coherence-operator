/*
 * Copyright (c) 2020, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1

import (
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
	"time"
)

func TestSetRoleStatusNewRole(t *testing.T) {
	g := NewGomegaWithT(t)

	cluster := &CoherenceCluster{}
	start := time.Now()
	cluster.SetRoleStatus("foo", true, 3, RoleStatusReady)
	end := time.Now()

	status := cluster.GetRoleStatus("foo")
	g.Expect(len(cluster.Status.RoleStatus)).To(Equal(1))
	g.Expect(status.Role).To(Equal("foo"))
	g.Expect(status.Ready).To(Equal(true))
	g.Expect(status.Count).To(Equal(int32(3)))
	g.Expect(status.Status).To(Equal(RoleStatusReady))

	c := status.GetCondition(RoleStatusReady)
	g.Expect(start.After(c.LastTransitionTime.Time)).To(Equal(false))
	g.Expect(end.Before(c.LastTransitionTime.Time)).To(Equal(false))
}

func TestSetRoleStatusForExistingRole(t *testing.T) {
	g := NewGomegaWithT(t)

	cluster := &CoherenceCluster{
		Status: CoherenceClusterStatus{
			RoleStatus: []ClusterRoleStatus{
				{
					Role:   "foo",
					Ready:  false,
					Count:  0,
					Status: RoleStatusCreated,
					Conditions: []ClusterRoleStatusCondition{
						{
							Status:             RoleStatusReady,
							LastTransitionTime: metav1.NewTime(time.Time{}),
						},
					},
				},
			}},
	}

	start := time.Now()
	cluster.SetRoleStatus("foo", true, 3, RoleStatusReady)
	end := time.Now()

	status := cluster.GetRoleStatus("foo")
	g.Expect(len(cluster.Status.RoleStatus)).To(Equal(1))
	g.Expect(status.Role).To(Equal("foo"))
	g.Expect(status.Ready).To(Equal(true))
	g.Expect(status.Count).To(Equal(int32(3)))
	g.Expect(status.Status).To(Equal(RoleStatusReady))

	c := status.GetCondition(RoleStatusReady)
	g.Expect(start.After(c.LastTransitionTime.Time)).To(Equal(false))
	g.Expect(end.Before(c.LastTransitionTime.Time)).To(Equal(false))
}
