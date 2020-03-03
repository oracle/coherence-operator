/*
 * Copyright (c) 2019, 2020 Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package coherencerole

import (
	"context"
	"fmt"
	"github.com/oracle/coherence-operator/pkg/flags"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	coherence "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	stubs "github.com/oracle/coherence-operator/pkg/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// These tests use fakes and stubs for the k8s and operator-sdk so that the
// tests will run without requiring a k8s cluster.
var _ = Describe("coherencerole_controller", func() {
	const (
		namespace    = "coherence-test"
		clusterName  = "test-cluster"
		roleName     = "storage"
		fullRoleName = clusterName + "-" + roleName
	)

	var (
		mgr         *stubs.FakeManager
		cluster     *coherence.CoherenceCluster
		roleNew     *coherence.CoherenceRole
		roleCurrent *coherence.CoherenceRole
		existing    []runtime.Object
		result      stubs.ReconcileResult
		errors      stubs.ClientErrors
		err         error

		controller *ReconcileCoherenceRole

		defaultCluster = &coherence.CoherenceCluster{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: namespace,
				Name:      clusterName,
			},
		}
	)

	JustBeforeEach(func() {
		mgr, err = stubs.NewFakeManager(existing...)
		Expect(err).NotTo(HaveOccurred())
		controller = newReconciler(mgr, &flags.CoherenceOperatorFlags{})
		// skip initialization for unit tests
		controller.SetInitialized(true)

		mgr.Client.DisableErrors()

		if cluster != nil {
			_ = mgr.Client.Create(context.TODO(), cluster)
		}

		if roleNew != nil {
			_ = mgr.Client.Create(context.TODO(), roleNew)
		}

		if roleCurrent != nil {
			spec := coherence.NewCoherenceInternalSpec(cluster, roleCurrent)
			specMap, err := coherence.CoherenceInternalSpecAsMapFromSpec(spec)
			Expect(err).NotTo(HaveOccurred())

			cohIntern := controller.CreateHelmValues(cluster, roleNew, specMap)

			err = mgr.Client.Create(context.TODO(), cohIntern)
			Expect(err).NotTo(HaveOccurred())
		}

		request := reconcile.Request{
			NamespacedName: types.NamespacedName{
				Namespace: namespace,
				Name:      fullRoleName,
			},
		}

		mgr.Client.EnableErrors(errors)
		controller.client = mgr.Client
		r, err := controller.Reconcile(request)
		result = stubs.ReconcileResult{Result: r, Error: err}
	})

	JustAfterEach(func() {
		// reset the errors
		errors = stubs.ClientErrors{}
	})

	When("The k8s client returns an error getting the CoherenceRole", func() {
		var err error = stubs.FakeError{Msg: "error getting role"}

		BeforeEach(func() {
			cluster = defaultCluster
			roleNew = nil
			errors.AddGetError(stubs.ErrorIf{KeyIs: &types.NamespacedName{Namespace: namespace, Name: fullRoleName}}, err)
		})

		When("reconcile is called", func() {
			It("should not return error", func() {
				Expect(result.Error).NotTo(HaveOccurred())
			})

			It("should re-queue the request", func() {
				Expect(result.Result).To(Equal(reconcile.Result{Requeue: true}))
			})

			It("should not fire an event", func() {
				mgr.AssertNoRemainingEvents()
			})

			It("should not create any CoherenceInternals", func() {
				mgr.AssertCoherenceInternals(namespace, 0)
			})
		})
	})

	When("the the parent CoherenceCluster does not exist", func() {
		BeforeEach(func() {
			cluster = nil

			roleNew = &coherence.CoherenceRole{
				ObjectMeta: metav1.ObjectMeta{Namespace: namespace, Name: fullRoleName},
				Spec: coherence.CoherenceRoleSpec{
					Role: roleName,
				},
			}
		})

		When("reconcile is called", func() {
			It("should not return error", func() {
				Expect(result.Error).NotTo(HaveOccurred())
			})

			It("should not re-queue the request", func() {
				Expect(result.Result).To(Equal(reconcile.Result{}))
			})

			It("should fire a failed event", func() {
				msg := fmt.Sprintf(invalidRoleEventMessage, fullRoleName, clusterName)
				event := mgr.AssertEvent()

				Expect(event.Type).To(Equal(corev1.EventTypeNormal))
				Expect(event.Reason).To(Equal(eventReasonFailed))
				Expect(event.Message).To(Equal(msg))

				mgr.AssertNoRemainingEvents()
			})

			It("should update the CoherenceRole's status to failed", func() {
				role := mgr.AssertCoherenceRoleExists(namespace, fullRoleName)
				Expect(role.Status.Status).To(Equal(coherence.RoleStatusFailed))
			})

			It("should not create any CoherenceInternals", func() {
				mgr.AssertCoherenceInternals(namespace, 0)
			})
		})
	})

	When("the k8s client returns an error getting the parent CoherenceCluster", func() {
		var err error = stubs.FakeError{Msg: "error getting cluster"}

		BeforeEach(func() {
			cluster = defaultCluster

			roleNew = &coherence.CoherenceRole{
				ObjectMeta: metav1.ObjectMeta{Namespace: namespace, Name: fullRoleName},
				Spec: coherence.CoherenceRoleSpec{
					Role: roleName,
				},
			}

			errors.AddGetError(stubs.ErrorIf{KeyIs: &types.NamespacedName{Namespace: namespace, Name: clusterName}}, err)
		})

		When("reconcile is called", func() {
			It("should not return error", func() {
				Expect(result.Error).NotTo(HaveOccurred())
			})

			It("should re-queue the request", func() {
				Expect(result.Result).To(Equal(reconcile.Result{Requeue: true}))
			})

			It("should fire a failed event", func() {
				msg := fmt.Sprintf(failedToGetParentCluster, clusterName, fullRoleName, err.Error())
				event := mgr.AssertEvent()

				Expect(event.Type).To(Equal(corev1.EventTypeNormal))
				Expect(event.Reason).To(Equal(eventReasonFailed))
				Expect(event.Message).To(Equal(msg))

				mgr.AssertNoRemainingEvents()
			})

			It("should not create any CoherenceInternals", func() {
				mgr.AssertCoherenceInternals(namespace, 0)
			})
		})
	})

	When("the k8s client returns an error getting the Helm values for a CoherenceRole", func() {
		var err error = stubs.FakeError{Msg: "error getting values"}

		BeforeEach(func() {
			cluster = defaultCluster

			roleNew = &coherence.CoherenceRole{
				ObjectMeta: metav1.ObjectMeta{Namespace: namespace, Name: fullRoleName},
				Spec: coherence.CoherenceRoleSpec{
					Role: roleName,
				},
			}

			errors.AddGetError(stubs.ErrorIf{
				KeyIs:  &types.NamespacedName{Namespace: namespace, Name: fullRoleName},
				TypeIs: &unstructured.Unstructured{},
			}, err)
		})

		When("reconcile is called", func() {
			It("should not return error", func() {
				Expect(result.Error).NotTo(HaveOccurred())
			})

			It("should re-queue the request", func() {
				Expect(result.Result).To(Equal(reconcile.Result{Requeue: true}))
			})

			It("should fire a failed event", func() {
				msg := fmt.Sprintf(failedToGetHelmValuesMessage, fullRoleName, err.Error())
				event := mgr.AssertEvent()

				Expect(event.Type).To(Equal(corev1.EventTypeNormal))
				Expect(event.Reason).To(Equal(eventReasonFailed))
				Expect(event.Message).To(Equal(msg))

				mgr.AssertNoRemainingEvents()
			})

			It("should not create any CoherenceInternals", func() {
				mgr.AssertCoherenceInternals(namespace, 0)
			})
		})
	})
})
