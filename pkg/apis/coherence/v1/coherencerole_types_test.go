/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	coherence "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Testing CoherenceRole struct", func() {

	When("Getting a CoherenceRole Coherence Cluster Name", func() {
		var role coherence.CoherenceRole

		When("the cluster label is present", func() {
			BeforeEach(func() {
				role = coherence.CoherenceRole{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "test-namespace",
						Name:      "test-cluster-foo",
						Labels:    map[string]string{coherence.CoherenceClusterLabel: "foo-cluster"},
					},
				}
			})

			It("should get the cluster name from the label", func() {
				Expect(role.GetCoherenceClusterName()).To(Equal("foo-cluster"))
			})
		})

		When("the cluster label is not present", func() {
			BeforeEach(func() {
				role = coherence.CoherenceRole{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "test-namespace",
						Name:      "test-cluster-foo",
					},
					Spec: coherence.CoherenceRoleSpec{Role: "foo"},
				}
			})

			It("should get the cluster name from the role resource's name", func() {
				Expect(role.GetCoherenceClusterName()).To(Equal("test-cluster"))
			})
		})

	})
})
