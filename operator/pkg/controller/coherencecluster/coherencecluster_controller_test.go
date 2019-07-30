package coherencecluster

import (
	"context"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	coherence "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	stubs "github.com/oracle/coherence-operator/pkg/controller/fakes"
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
	)

	JustBeforeEach(func() {
		mgr = stubs.NewFakeManager(existing...)

		if cluster != nil {
			_ = mgr.Client.Create(context.TODO(), cluster)
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
				name := cluster.Spec.DefaultRole.GetFullRoleName(cluster)
				role := mgr.AssertCoherenceRoleExists(testNamespace, name)
				Expect(role.Spec).To(Equal(coherence.CoherenceRoleSpec{}))
			})

			It("should fire a successful CoherenceRole create event", func() {
				name := cluster.Spec.DefaultRole.GetFullRoleName(cluster)
				msg := fmt.Sprintf(createEventMessage, name, testClusterName)
				event := mgr.AssertEvent()

				Expect(event.Owner).To(Equal(cluster))
				Expect(event.Type).To(Equal(v1.EventTypeNormal))
				Expect(event.Reason).To(Equal(EventReasonCreated))
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
				RoleName: "storage",
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
				Expect(event.Type).To(Equal(v1.EventTypeNormal))
				Expect(event.Reason).To(Equal(EventReasonCreated))
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
				Expect(role.Spec).To(Equal(roleSpec))
			})
		})
	})

	When("a CoherenceCluster that has two roles is added and reconcile called", func() {
		var roleSpecOne coherence.CoherenceRoleSpec
		var roleSpecTwo coherence.CoherenceRoleSpec

		BeforeEach(func() {
			roleSpecOne = coherence.CoherenceRoleSpec{
				RoleName: "storage",
			}

			roleSpecTwo = coherence.CoherenceRoleSpec{
				RoleName: "proxy",
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
				Expect(eventOne.Type).To(Equal(v1.EventTypeNormal))
				Expect(eventOne.Reason).To(Equal(EventReasonCreated))
				Expect(eventOne.Message).To(Equal(msgOne))

				Expect(eventTwo.Owner).To(Equal(cluster))
				Expect(eventTwo.Type).To(Equal(v1.EventTypeNormal))
				Expect(eventTwo.Reason).To(Equal(EventReasonCreated))
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
				Expect(role.Spec).To(Equal(roleSpecOne))
			})

			It("should create a CoherenceRole for the second role", func() {
				name := roleSpecTwo.GetFullRoleName(cluster)
				role := mgr.AssertCoherenceRoleExists(testNamespace, name)
				Expect(role.Spec).To(Equal(roleSpecTwo))
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
			var updatedgRoleSpec coherence.CoherenceRoleSpec

			BeforeEach(func() {
				existingRoleSpec = coherence.CoherenceRoleSpec{
					RoleName: "storage",
				}

				imageName := "coherence:1.2.3"

				updatedgRoleSpec = coherence.CoherenceRoleSpec{
					RoleName: "storage",
					Images: &coherence.Images{
						Coherence: &coherence.ImageSpec{Image: &imageName},
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
						Roles: []coherence.CoherenceRoleSpec{updatedgRoleSpec},
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
					Expect(role.Spec).To(Equal(updatedgRoleSpec))
				})

				It("should fire a successful CoherenceRole update event", func() {
					roleName := testClusterName + "-storage"
					msg := fmt.Sprintf(updateEventMessage, roleName, testClusterName)
					event := mgr.AssertEvent()

					Expect(event.Owner).To(Equal(cluster))
					Expect(event.Type).To(Equal(v1.EventTypeNormal))
					Expect(event.Reason).To(Equal(EventReasonUpdated))
					Expect(event.Message).To(Equal(msg))

					mgr.AssertNoRemainingEvents()
				})
			})
		})

		When("an existing role is updated to zero replicas", func() {
			var zero int32 = 0
			var three int32 = 3

			BeforeEach(func() {
				existing = []runtime.Object{
					&coherence.CoherenceRole{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: testNamespace,
							Name:      testClusterName + "-storage",
							Labels:    map[string]string{coherence.CoherenceClusterLabel: testClusterName},
						},
						Spec: coherence.CoherenceRoleSpec{
							RoleName: "storage",
							Replicas: &three,
						},
					},
				}

				cluster = &coherence.CoherenceCluster{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: testNamespace,
						Name:      testClusterName,
					},
					Spec: coherence.CoherenceClusterSpec{
						Roles: []coherence.CoherenceRoleSpec{
							{RoleName: "storage", Replicas: &zero},
						},
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

				It("should delete the CoherenceRole", func() {
					mgr.AssertCoherenceRoles(testNamespace, 0)
				})

				It("should fire a successful CoherenceRole delete event", func() {
					roleName := testClusterName + "-storage"
					msg := fmt.Sprintf(deleteEventMessage, roleName, testClusterName)
					event := mgr.AssertEvent()

					Expect(event.Owner).To(Equal(cluster))
					Expect(event.Type).To(Equal(v1.EventTypeNormal))
					Expect(event.Reason).To(Equal(EventReasonDeleted))
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
						RoleName: "storage",
					},
				}

				proxy = coherence.CoherenceRole{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: testNamespace,
						Name:      testClusterName + "-proxy",
						Labels:    map[string]string{coherence.CoherenceClusterLabel: testClusterName},
					},
					Spec: coherence.CoherenceRoleSpec{
						RoleName: "proxy",
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
					Expect(role).To(Equal(&storage))
				})

				It("should fire a successful CoherenceRole delete event", func() {
					msg := fmt.Sprintf(deleteEventMessage, proxy.Name, testClusterName)
					event := mgr.AssertEvent()

					Expect(event.Owner).To(Equal(cluster))
					Expect(event.Type).To(Equal(v1.EventTypeNormal))
					Expect(event.Reason).To(Equal(EventReasonDeleted))
					Expect(event.Message).To(Equal(msg))

					mgr.AssertNoRemainingEvents()
				})
			})
		})
	})
})
