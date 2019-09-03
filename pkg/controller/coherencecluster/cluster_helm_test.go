/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package coherencecluster

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	cohv1 "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	"github.com/oracle/coherence-operator/pkg/controller/coherencerole"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/utils/pointer"

	stubs "github.com/oracle/coherence-operator/pkg/fakes"
)

var _ = Describe("CoherenceCluster to Helm install verification suite", func() {
	const (
		testNamespace   = "test-namespace"
		testClusterName = "test-cluster"
	)

	var (
		mgr     *stubs.FakeManager
		cluster *cohv1.CoherenceCluster
		result  *stubs.HelmInstallResult
	)

	// Before each test run the fake Helm install using the cluster variable
	// and capture the result to be asserted by the tests
	JustBeforeEach(func() {
		mgr = stubs.NewFakeManager()
		cr := NewClusterReconciler(mgr)
		rr := coherencerole.NewRoleReconciler(mgr)
		helm := stubs.NewFakeHelm(mgr, cr, rr)

		r, err := helm.HelmInstallFromCoherenceCluster(cluster)
		Expect(err).NotTo(HaveOccurred())
		result = r
	})

	When("installing a minimal CoherenceCluster", func() {
		// Create a minimal valid CoherenceCluster to use for the Helm install
		BeforeEach(func() {
			cluster = &cohv1.CoherenceCluster{}
			cluster.SetNamespace(testNamespace)
			cluster.SetName(testClusterName)
		})

		It("should have created one StatefulSet", func() {
			list := appsv1.StatefulSetList{}
			err := result.List(&list)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(list.Items)).To(Equal(1))
		})

		It("should have created a StatefulSet with the same name as the role", func() {
			// find the corresponding StatefulSet in the Helm results
			sts, err := findStatefulSet(result, cluster, cohv1.DefaultRoleName)

			// Assert that the StatefulSet exists with the correct full role name
			Expect(err).NotTo(HaveOccurred())
			Expect(sts.GetName()).To(Equal(cluster.GetFullRoleName(cohv1.DefaultRoleName)))
		})

		It("should have created a StatefulSet with the default replica count", func() {
			// find the corresponding StatefulSet in the Helm results
			sts, err := findStatefulSet(result, cluster, cohv1.DefaultRoleName)

			// Assert that the StatefulSet exists with the correct full role name
			Expect(err).NotTo(HaveOccurred())
			Expect(sts.Spec.Replicas).To(Equal(pointer.Int32Ptr(cohv1.DefaultReplicas)))
		})
	})

	When("installing a CoherenceCluster with two roles", func() {
		var (
			roleOneName     = "data"
			roleOneReplicas = pointer.Int32Ptr(5)
			roleTwoName     = "proxy"
			roleTwoReplicas = pointer.Int32Ptr(2)
		)

		// Create a valid CoherenceCluster with two roles to use for the Helm install
		BeforeEach(func() {
			roleOne := cohv1.CoherenceRoleSpec{Role: roleOneName, Replicas: roleOneReplicas}
			roleTwo := cohv1.CoherenceRoleSpec{Role: roleTwoName, Replicas: roleTwoReplicas}

			cluster = &cohv1.CoherenceCluster{}
			cluster.SetNamespace(testNamespace)
			cluster.SetName(testClusterName)
			cluster.Spec.Roles = []cohv1.CoherenceRoleSpec{roleOne, roleTwo}
		})

		It("should have created two StatefulSet", func() {
			list := appsv1.StatefulSetList{}
			err := result.List(&list)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(list.Items)).To(Equal(2))
		})

		It("should have created a StatefulSet with the same name as roleOne", func() {
			sts, err := findStatefulSet(result, cluster, roleOneName)
			Expect(err).NotTo(HaveOccurred())
			Expect(sts.GetName()).To(Equal(cluster.GetFullRoleName(roleOneName)))
		})

		It("should have created a StatefulSet with the same name as roleTwo", func() {
			sts, err := findStatefulSet(result, cluster, roleTwoName)
			Expect(err).NotTo(HaveOccurred())
			Expect(sts.GetName()).To(Equal(cluster.GetFullRoleName(roleTwoName)))
		})

		It("should have created a StatefulSet for roleOne with the correct replica count", func() {
			sts, err := findStatefulSet(result, cluster, roleOneName)
			Expect(err).NotTo(HaveOccurred())
			Expect(sts.Spec.Replicas).To(Equal(roleOneReplicas))
		})

		It("should have created a StatefulSet for roleTwo with the correct replica count", func() {
			sts, err := findStatefulSet(result, cluster, roleTwoName)
			Expect(err).NotTo(HaveOccurred())
			Expect(sts.Spec.Replicas).To(Equal(roleTwoReplicas))
		})
	})
})

// ----- helpers ------------------------------------------------------------

// Shared function to find a StatefulSet in a Helm result
var findStatefulSet = func(result *stubs.HelmInstallResult, cluster *cohv1.CoherenceCluster, roleName string) (appsv1.StatefulSet, error) {
	name := cluster.GetFullRoleName(roleName)
	sts := appsv1.StatefulSet{}
	err := result.Get(name, &sts)
	return sts, err
}
