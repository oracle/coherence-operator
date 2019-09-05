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
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/utils/pointer"

	stubs "github.com/oracle/coherence-operator/pkg/fakes"
)

var _ = Describe("Persistence CoherenceCluster to Helm install verification suite", func() {
	const (
		testNamespace   = "test-namespace"
		testClusterName = "test-cluster"
		roleOneName     = "one"
	)

	var (
		// A fake manager to use to obtain the k8s client
		mgr *stubs.FakeManager
		// The CoherenceCluster to create the Helm install from
		cluster *cohv1.CoherenceCluster
		// The result of the Helm install
		result *stubs.HelmInstallResult

        // ----- helpers ------------------------------------------------------------
		NewPersistenceStorageSpec = func(enabled, pvc bool) *cohv1.PersistentStorageSpec {
			if (!enabled) {
				return &cohv1.PersistentStorageSpec{ Enabled: pointer.BoolPtr(false), }
			}
			if pvc {
				return &cohv1.PersistentStorageSpec{ Enabled: pointer.BoolPtr(true),
					PersistentVolumeClaim: &corev1.PersistentVolumeClaimSpec {
						AccessModes: []corev1.PersistentVolumeAccessMode{ "ReadWriteOnce",},
						Resources: corev1.ResourceRequirements{
							Requests: map[corev1.ResourceName]resource.Quantity{"storage": resource.MustParse("2Gi"),},
						},
					},
				}
			}

			twoGiResource := resource.MustParse("2Gi")
			return &cohv1.PersistentStorageSpec{ Enabled: pointer.BoolPtr(true),
				Volume: &corev1.Volume{
					VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{ SizeLimit: &twoGiResource,}},
				},
			}
		}

		createCluster = func(persistence *cohv1.PersistentStorageSpec, snapshot *cohv1.PersistentStorageSpec) {
			roleOne := cohv1.CoherenceRoleSpec{Role: roleOneName, Replicas: pointer.Int32Ptr(1),
				Persistence: persistence, Snapshot: snapshot, }

			cluster = &cohv1.CoherenceCluster{}
			cluster.SetNamespace(testNamespace)
			cluster.SetName(testClusterName)
			cluster.Spec.Roles = []cohv1.CoherenceRoleSpec{roleOne}
		}

		assertInstall = func(expectedPvcs []string, expectedVolMap map[string]bool) func() {
			return func() {
				By("Checking only have one StatefulSet")
				list := appsv1.StatefulSetList{}
				err := result.List(&list)
				Expect(err).NotTo(HaveOccurred())
				Expect(len(list.Items)).To(Equal(1))

				By("Checking StatefulSet name")
				sts, err := findStatefulSet(result, cluster, roleOneName)
				Expect(err).NotTo(HaveOccurred())
				Expect(sts.GetName()).To(Equal(cluster.GetFullRoleName(roleOneName)))

				By("Checking StatefulSet PVC")
				Expect(len(sts.Spec.VolumeClaimTemplates)).To(Equal(len(expectedPvcs)))
				expectedPvcSet := make(map[string]bool)
				for _, p := range(expectedPvcs) {
					expectedPvcSet[p] = true
				}
				pvcSet := make(map[string]bool)
				for _, k := range sts.Spec.VolumeClaimTemplates {
					pvcSet[k.Name] = true
				}
				Expect(expectedPvcSet).To(Equal(pvcSet))

				By("Checking StatefulSet volume")
				volSet := make(map[string]bool)
				for _, k := range sts.Spec.Template.Spec.Volumes {
					volSet[k.Name] = true
				}
				for vol, exist := range expectedVolMap {
					_, ok := volSet[vol]
					Expect(exist).To(Equal(ok))
				}
			}
		}
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

	When("installing a CoherenceCluster with persistence/pvc and snapshot/pvc", func() {
		// Create a valid CoherenceCluster with a role to use for the Helm install
		BeforeEach(func() {
			createCluster(NewPersistenceStorageSpec(true, true), NewPersistenceStorageSpec(true, true))
		})

		It("should have two pvc and no corresponding volumes", assertInstall(
			[]string{"persistence-volume", "snapshot-volume"},
			map[string]bool{ "persistence-volume": false, "snapshot-volume": false}))
	})

	When("installing a CoherenceCluster with persistence/volume, and without snapshot", func() {
		// Create a valid CoherenceCluster with a role to use for the Helm install
		BeforeEach(func() {
			createCluster(NewPersistenceStorageSpec(true, false), nil)
		})

		It("should have no pvc and only corresponding persistence-volume", assertInstall(
			[]string{},
			map[string]bool{ "persistence-volume": true, "snapshot-volume": false}))
	})

	When("installing a CoherenceCluster with persistence disabled, and snapshot/pvc", func() {
		// Create a valid CoherenceCluster with a role to use for the Helm install
		BeforeEach(func() {
			createCluster(NewPersistenceStorageSpec(false, false), NewPersistenceStorageSpec(true, true))
		})

		It("should have no pvc and only corresponding persistence-volume", assertInstall(
			[]string{ "snapshot-volume" },
			map[string]bool{ "persistence-volume": false, "snapshot-volume": false}))
	})

	When("installing a CoherenceCluster without persistence and without snapshotc", func() {
		// Create a valid CoherenceCluster with a role to use for the Helm install
		BeforeEach(func() {
			createCluster(nil, nil)
		})

		It("should have no pvc and only corresponding persistence-volume", assertInstall(
			[]string{},
			map[string]bool{ "persistence-volume": false, "snapshot-volume": false}))
	})

	When("installing a CoherenceCluster without persistence and snapshot/volume", func() {
		// Create a valid CoherenceCluster with a role to use for the Helm install
		BeforeEach(func() {
			createCluster(nil, NewPersistenceStorageSpec(true, false))
		})

		It("should have no pvc and only corresponding persistence-volume", assertInstall(
			[]string{},
			map[string]bool{ "persistence-volume": false, "snapshot-volume": true}))
	})
})

