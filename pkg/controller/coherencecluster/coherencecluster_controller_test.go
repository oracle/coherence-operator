/*
 * Copyright (c) 2019, 2020 Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package coherencecluster

import (
	"context"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	coherence "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	stubs "github.com/oracle/coherence-operator/pkg/fakes"
)

var _ = Describe("coherencecluster_controller", func() {
	const (
		testNamespace   = "test-namespace"
		testClusterName = "test-cluster"
	)

	var (
		mgr      *stubs.FakeManager
		cluster  *coherence.CoherenceCluster
		existing []runtime.Object
		result   stubs.ReconcileResult
		err      error
	)

	JustBeforeEach(func() {
		mgr, err = stubs.NewFakeManager(existing...)
		Expect(err).NotTo(HaveOccurred())

		if cluster != nil {
			_ = mgr.Client.Create(context.TODO(), cluster)
			_ = mgr.Client.Get(context.TODO(), types.NamespacedName{Namespace: testNamespace, Name: testClusterName}, cluster)
		}

		request := reconcile.Request{
			NamespacedName: types.NamespacedName{
				Namespace: testNamespace,
				Name:      testClusterName,
			},
		}

		controller := newReconciler(mgr)
		r, err := controller.Reconcile(request)
		result = stubs.ReconcileResult{Result: r, Error: err}
	})

	When("a CoherenceCluster that has no Spec is added", func() {
		BeforeEach(func() {
			cluster = &coherence.CoherenceCluster{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: testNamespace,
					Name:      testClusterName,
				},
			}
		})

		When("reconcile is called", func() {
			It("should not return error", func() {
				Expect(result.Error).To(BeNil())
			})

			It("should not re-queue the request", func() {
				Expect(result.Result).To(Equal(reconcile.Result{}))
			})

			It("should create one default CoherenceRole", func() {
				mgr.AssertCoherenceRoles(testNamespace, 1)
				name := cluster.Spec.CoherenceRoleSpec.GetFullRoleName(cluster)
				role := mgr.AssertCoherenceRoleExists(testNamespace, name)
				Expect(role.Spec).To(Equal(coherence.CoherenceRoleSpec{Replicas: pointer.Int32Ptr(coherence.DefaultReplicas)}))
			})

			It("should fire a successful CoherenceRole create event", func() {
				name := cluster.Spec.CoherenceRoleSpec.GetFullRoleName(cluster)
				msg := fmt.Sprintf(createEventMessage, name, testClusterName)
				event := mgr.AssertEvent()

				Expect(event.Owner).To(Equal(cluster))
				Expect(event.Type).To(Equal(corev1.EventTypeNormal))
				Expect(event.Reason).To(Equal(eventReasonCreated))
				Expect(event.Message).To(Equal(msg))

				mgr.AssertNoRemainingEvents()
			})

			It("should create the WKA service", func() {
				mgr.AssertWkaService(testNamespace, cluster)
			})
		})
	})

	When("a CoherenceCluster that has one role is added", func() {
		var roleSpec coherence.CoherenceRoleSpec

		BeforeEach(func() {
			roleSpec = coherence.CoherenceRoleSpec{
				Role: "storage",
			}

			cluster = &coherence.CoherenceCluster{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: testNamespace,
					Name:      testClusterName,
				},
				Spec: coherence.CoherenceClusterSpec{
					Roles: []coherence.CoherenceRoleSpec{roleSpec},
				},
			}
		})

		When("reconcile is called", func() {
			It("should not return error", func() {
				Expect(result.Error).To(BeNil())
			})

			It("should not re-queue the request", func() {
				Expect(result.Result).To(Equal(reconcile.Result{}))
			})

			It("should fire a successful CoherenceRole create event", func() {
				name := roleSpec.GetFullRoleName(cluster)
				msg := fmt.Sprintf(createEventMessage, name, testClusterName)
				event := mgr.AssertEvent()

				Expect(event.Owner).To(Equal(cluster))
				Expect(event.Type).To(Equal(corev1.EventTypeNormal))
				Expect(event.Reason).To(Equal(eventReasonCreated))
				Expect(event.Message).To(Equal(msg))

				mgr.AssertNoRemainingEvents()
			})

			It("should create the WKA service", func() {
				mgr.AssertWkaService(testNamespace, cluster)
			})

			It("should create one CoherenceRole", func() {
				mgr.AssertCoherenceRoles(testNamespace, 1)
				name := roleSpec.GetFullRoleName(cluster)
				role := mgr.AssertCoherenceRoleExists(testNamespace, name)
				expected := roleSpec.DeepCopy()
				expected.SetReplicas(coherence.DefaultReplicas)
				Expect(role.Spec).To(Equal(*expected))
			})
		})
	})

	When("a CoherenceCluster that has two roles is added and reconcile called", func() {
		var roleSpecOne coherence.CoherenceRoleSpec
		var roleSpecTwo coherence.CoherenceRoleSpec

		BeforeEach(func() {
			roleSpecOne = coherence.CoherenceRoleSpec{
				Role: "storage",
			}

			roleSpecTwo = coherence.CoherenceRoleSpec{
				Role: "proxy",
			}

			cluster = &coherence.CoherenceCluster{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: testNamespace,
					Name:      testClusterName,
				},
				Spec: coherence.CoherenceClusterSpec{
					Roles: []coherence.CoherenceRoleSpec{roleSpecOne, roleSpecTwo},
				},
			}
		})

		When("reconcile is called", func() {
			It("should not return error", func() {
				Expect(result.Error).To(BeNil())
			})

			It("should not re-queue the request", func() {
				Expect(result.Result).To(Equal(reconcile.Result{}))
			})

			It("should fire two successful CoherenceRole create event", func() {
				roleNameOne := roleSpecOne.GetFullRoleName(cluster)
				msgOne := fmt.Sprintf(createEventMessage, roleNameOne, testClusterName)
				roleNameTwo := roleSpecTwo.GetFullRoleName(cluster)
				msgTwo := fmt.Sprintf(createEventMessage, roleNameTwo, testClusterName)

				eventOne := mgr.AssertEvent()
				eventTwo := mgr.AssertEvent()

				Expect(eventOne.Owner).To(Equal(cluster))
				Expect(eventOne.Type).To(Equal(corev1.EventTypeNormal))
				Expect(eventOne.Reason).To(Equal(eventReasonCreated))
				Expect(eventOne.Message).To(Equal(msgOne))

				Expect(eventTwo.Owner).To(Equal(cluster))
				Expect(eventTwo.Type).To(Equal(corev1.EventTypeNormal))
				Expect(eventTwo.Reason).To(Equal(eventReasonCreated))
				Expect(eventTwo.Message).To(Equal(msgTwo))

				mgr.AssertNoRemainingEvents()
			})

			It("should create the WKA service", func() {
				mgr.AssertWkaService(testNamespace, cluster)
			})

			It("should create two CoherenceRoles", func() {
				mgr.AssertCoherenceRoles(testNamespace, 2)
			})

			It("should create a CoherenceRole for the first role", func() {
				name := roleSpecOne.GetFullRoleName(cluster)
				role := mgr.AssertCoherenceRoleExists(testNamespace, name)
				expected := roleSpecOne.DeepCopy()
				expected.SetReplicas(coherence.DefaultReplicas)
				Expect(role.Spec).To(Equal(*expected))
			})

			It("should create a CoherenceRole for the second role", func() {
				name := roleSpecTwo.GetFullRoleName(cluster)
				role := mgr.AssertCoherenceRoleExists(testNamespace, name)
				expected := roleSpecTwo.DeepCopy()
				expected.SetReplicas(coherence.DefaultReplicas)
				Expect(role.Spec).To(Equal(*expected))
			})
		})
	})

	When("a CoherenceCluster does not exist", func() {
		BeforeEach(func() {
			cluster = nil
		})

		When("reconcile is called", func() {
			It("should not return error", func() {
				Expect(result.Error).To(BeNil())
			})

			It("should not re-queue the request", func() {
				Expect(result.Result).To(Equal(reconcile.Result{}))
			})

			It("should not fire an event", func() {
				_, found := mgr.NextEvent()
				Expect(found).To(BeFalse())
			})

			It("should not create the WKA service", func() {
				exists := mgr.ServiceExists(testNamespace, cluster.GetWkaServiceName())
				Expect(exists).To(BeFalse())
			})

			It("should not create any CoherenceRoles", func() {
				mgr.AssertCoherenceRoles(testNamespace, 0)
			})
		})
	})

	When("a CoherenceCluster is updated", func() {
		When("an existing role is updated", func() {
			var existingRoleSpec coherence.CoherenceRoleSpec
			var updatedRoleSpec coherence.CoherenceRoleSpec

			BeforeEach(func() {
				existingRoleSpec = coherence.CoherenceRoleSpec{
					Role: "storage",
				}

				imageName := "coherence:1.2.3"

				updatedRoleSpec = coherence.CoherenceRoleSpec{
					Role: "storage",
					Coherence: &coherence.CoherenceSpec{
						ImageSpec: coherence.ImageSpec{Image: &imageName},
					},
				}

				existing = []runtime.Object{
					&coherence.CoherenceRole{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: testNamespace,
							Name:      testClusterName + "-storage",
							Labels:    map[string]string{coherence.CoherenceClusterLabel: testClusterName},
						},
						Spec: existingRoleSpec,
					},
				}

				cluster = &coherence.CoherenceCluster{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: testNamespace,
						Name:      testClusterName,
					},
					Spec: coherence.CoherenceClusterSpec{
						Roles: []coherence.CoherenceRoleSpec{updatedRoleSpec},
					},
				}
			})

			When("reconcile is called", func() {
				It("should not return error", func() {
					Expect(result.Error).To(BeNil())
				})

				It("should not re-queue the request", func() {
					Expect(result.Result).To(Equal(reconcile.Result{}))
				})

				It("should not create the WKA service", func() {
					exists := mgr.ServiceExists(testNamespace, cluster.GetWkaServiceName())
					Expect(exists).To(BeFalse())
				})

				It("should update the CoherenceRole", func() {
					mgr.AssertCoherenceRoles(testNamespace, 1)
					name := existingRoleSpec.GetFullRoleName(cluster)
					role := mgr.AssertCoherenceRoleExists(testNamespace, name)
					expected := updatedRoleSpec.DeepCopy()
					expected.SetReplicas(coherence.DefaultReplicas)
					Expect(role.Spec).To(Equal(*expected))
				})

				It("should fire a successful CoherenceRole update event", func() {
					roleName := testClusterName + "-storage"
					msg := fmt.Sprintf(updateEventMessage, roleName, testClusterName)
					event := mgr.AssertEvent()

					Expect(event.Owner).To(Equal(cluster))
					Expect(event.Type).To(Equal(corev1.EventTypeNormal))
					Expect(event.Reason).To(Equal(eventReasonUpdated))
					Expect(event.Message).To(Equal(msg))

					mgr.AssertNoRemainingEvents()
				})
			})
		})

		When("an existing role is updated to zero replicas", func() {
			var zero int32 = 0
			var three int32 = 3

			var existingRoleSpec coherence.CoherenceRoleSpec
			var updatedRoleSpec coherence.CoherenceRoleSpec

			BeforeEach(func() {
				existingRoleSpec = coherence.CoherenceRoleSpec{
					Role:     "data",
					Replicas: &three,
				}

				updatedRoleSpec = coherence.CoherenceRoleSpec{
					Role:     "data",
					Replicas: &zero,
				}

				existing = []runtime.Object{
					&coherence.CoherenceRole{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: testNamespace,
							Name:      testClusterName + "-data",
							Labels:    map[string]string{coherence.CoherenceClusterLabel: testClusterName},
						},
						Spec: existingRoleSpec,
					},
				}

				cluster = &coherence.CoherenceCluster{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: testNamespace,
						Name:      testClusterName,
					},
					Spec: coherence.CoherenceClusterSpec{
						Roles: []coherence.CoherenceRoleSpec{updatedRoleSpec},
					},
				}
			})

			When("reconcile is called", func() {
				It("should not return error", func() {
					Expect(result.Error).To(BeNil())
				})

				It("should not re-queue the request", func() {
					Expect(result.Result).To(Equal(reconcile.Result{}))
				})

				It("should not create the WKA service", func() {
					exists := mgr.ServiceExists(testNamespace, cluster.GetWkaServiceName())
					Expect(exists).To(BeFalse())
				})

				It("should not delete the CoherenceRole", func() {
					mgr.AssertCoherenceRoles(testNamespace, 1)
					name := existingRoleSpec.GetFullRoleName(cluster)
					role := mgr.AssertCoherenceRoleExists(testNamespace, name)
					Expect(role.Spec).To(Equal(updatedRoleSpec))
				})

				It("should fire a successful CoherenceRole update event", func() {
					roleName := testClusterName + "-data"
					msg := fmt.Sprintf(updateEventMessage, roleName, testClusterName)
					event := mgr.AssertEvent()

					Expect(event.Owner).To(Equal(cluster))
					Expect(event.Type).To(Equal(corev1.EventTypeNormal))
					Expect(event.Reason).To(Equal(eventReasonUpdated))
					Expect(event.Message).To(Equal(msg))

					mgr.AssertNoRemainingEvents()
				})
			})
		})

		When("an existing role is removed from the CoherenceCluster", func() {
			var storage coherence.CoherenceRole
			var proxy coherence.CoherenceRole

			BeforeEach(func() {
				storage = coherence.CoherenceRole{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: testNamespace,
						Name:      testClusterName + "-storage",
						Labels:    map[string]string{coherence.CoherenceClusterLabel: testClusterName},
					},
					Spec: coherence.CoherenceRoleSpec{
						Role:     "storage",
						Replicas: pointer.Int32Ptr(2),
					},
				}

				proxy = coherence.CoherenceRole{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: testNamespace,
						Name:      testClusterName + "-proxy",
						Labels:    map[string]string{coherence.CoherenceClusterLabel: testClusterName},
					},
					Spec: coherence.CoherenceRoleSpec{
						Role:     "proxy",
						Replicas: pointer.Int32Ptr(2),
					},
				}

				existing = []runtime.Object{&storage, &proxy}

				cluster = &coherence.CoherenceCluster{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: testNamespace,
						Name:      testClusterName,
					},
					Spec: coherence.CoherenceClusterSpec{
						Roles: []coherence.CoherenceRoleSpec{storage.Spec},
					},
				}
			})

			When("reconcile is called", func() {
				It("should not return error", func() {
					Expect(result.Error).To(BeNil())
				})

				It("should not re-queue the request", func() {
					Expect(result.Result).To(Equal(reconcile.Result{}))
				})

				It("should not create the WKA service", func() {
					exists := mgr.ServiceExists(testNamespace, cluster.GetWkaServiceName())
					Expect(exists).To(BeFalse())
				})

				It("should delete the proxy CoherenceRole leaving storage unchanged", func() {
					mgr.AssertCoherenceRoles(testNamespace, 1)
					role := mgr.AssertCoherenceRoleExists(testNamespace, storage.Name)
					Expect(role.Spec).To(Equal(storage.Spec))
				})

				It("should fire a successful CoherenceRole delete event", func() {
					msg := fmt.Sprintf(deleteEventMessage, proxy.Name, testClusterName)
					event := mgr.AssertEvent()

					Expect(event.Owner).To(Equal(cluster))
					Expect(event.Type).To(Equal(corev1.EventTypeNormal))
					Expect(event.Reason).To(Equal(eventReasonDeleted))
					Expect(event.Message).To(Equal(msg))

					mgr.AssertNoRemainingEvents()
				})
			})
		})
	})
})
