package coherencerole

import (
	"context"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	coherence "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	stubs "github.com/oracle/coherence-operator/pkg/controller/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// These tests use fakes and stubs for the k8s and operator-sdk so that the
// tests will run without requiring a k8s cluster.
var _ = Describe("coherencerole_controller", func() {
	const (
		testNamespace   = "coherence-test"
		testClusterName = "test-cluster"
		roleName        = "storage"
		fullRoleName    = testClusterName + "-" + roleName
	)

	var (
		mgr         *stubs.FakeManager
		cluster     *coherence.CoherenceCluster
		roleNew     *coherence.CoherenceRole
		roleCurrent *coherence.CoherenceRole
		existing    []runtime.Object
		result      stubs.ReconcileResult

		controller *ReconcileCoherenceRole

		defaultCluster = &coherence.CoherenceCluster{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: testNamespace,
				Name:      testClusterName,
			},
		}
	)

	JustBeforeEach(func() {
		mgr = stubs.NewFakeManager(existing...)
		controller = newReconciler(mgr)

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

			cohIntern := controller.CreateCoherenceInternal(cluster, roleNew, specMap)

			err = mgr.Client.Create(context.TODO(), cohIntern)
			Expect(err).NotTo(HaveOccurred())
		}

		request := reconcile.Request{
			NamespacedName: types.NamespacedName{
				Namespace: testNamespace,
				Name:      fullRoleName,
			},
		}

		r, err := controller.Reconcile(request)
		result = stubs.ReconcileResult{Result: r, Error: err}
	})

	When("a CoherenceRole does not exist (has been deleted)", func() {
		BeforeEach(func() {
			cluster = defaultCluster
			roleNew = nil
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

			It("should not create any CoherenceInternals", func() {
				mgr.AssertCoherenceInternals(testNamespace, 0)
			})
		})
	})

	When("a CoherenceRole is created", func() {
		BeforeEach(func() {
			cluster = defaultCluster
			roleNew = &coherence.CoherenceRole{
				ObjectMeta: metav1.ObjectMeta{Namespace: testNamespace, Name: fullRoleName},
				Spec: coherence.CoherenceRoleSpec{
					Role: roleName,
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

			It("should fire a create event", func() {
				msg := fmt.Sprintf(createEventMessage, roleNew.Name, roleNew.Name)
				event := mgr.AssertEvent()

				Expect(event.Type).To(Equal(corev1.EventTypeNormal))
				Expect(event.Reason).To(Equal(eventReasonCreated))
				Expect(event.Message).To(Equal(msg))

				mgr.AssertNoRemainingEvents()
			})

			It("should create a CoherenceInternal", func() {
				u := mgr.AssertCoherenceInternalExists(testNamespace, fullRoleName)
				roleSpec := UnstructuredToCoherenceInternalSpec(u)
				expected := coherence.NewCoherenceInternalSpec(cluster, roleNew)
				Expect(roleSpec).To(Equal(expected))
			})
		})
	})

	When("a CoherenceRole is updated", func() {
		BeforeEach(func() {
			cluster = defaultCluster

			imageOrig := "coherence:1.0"
			imageNew := "coherence:2.0"
			var replicas int32 = 3

			roleCurrent = &coherence.CoherenceRole{
				ObjectMeta: metav1.ObjectMeta{Namespace: testNamespace, Name: fullRoleName},
				Spec: coherence.CoherenceRoleSpec{
					Role:     roleName,
					Replicas: &replicas,
					Images: &coherence.Images{
						Coherence: &coherence.ImageSpec{
							Image: &imageOrig,
						},
					},
				},
			}

			roleNew = &coherence.CoherenceRole{
				ObjectMeta: metav1.ObjectMeta{Namespace: testNamespace, Name: fullRoleName},
				Spec: coherence.CoherenceRoleSpec{
					Role:     roleName,
					Replicas: &replicas,
					Images: &coherence.Images{
						Coherence: &coherence.ImageSpec{
							Image: &imageNew,
						},
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

			It("should fire an update event", func() {
				msg := fmt.Sprintf(updateEventMessage, roleNew.Name, roleNew.Name)
				event := mgr.AssertEvent()

				Expect(event.Type).To(Equal(corev1.EventTypeNormal))
				Expect(event.Reason).To(Equal(eventReasonUpdated))
				Expect(event.Message).To(Equal(msg))

				mgr.AssertNoRemainingEvents()
			})

			It("should update the CoherenceInternal", func() {
				u := mgr.AssertCoherenceInternalExists(testNamespace, fullRoleName)
				roleSpec := UnstructuredToCoherenceInternalSpec(u)
				expected := coherence.NewCoherenceInternalSpec(cluster, roleNew)
				Expect(roleSpec).To(Equal(expected))
			})
		})
	})

	When("a CoherenceRole is unchanged and the StatefulSet is unchanged", func() {
		var sts appsv1.StatefulSet

		BeforeEach(func() {
			cluster = defaultCluster

			var replicas int32 = 3

			roleCurrent = &coherence.CoherenceRole{
				ObjectMeta: metav1.ObjectMeta{Namespace: testNamespace, Name: fullRoleName},
				Spec: coherence.CoherenceRoleSpec{
					Role:     roleName,
					Replicas: &replicas,
				},
			}

			roleNew = &coherence.CoherenceRole{
				ObjectMeta: metav1.ObjectMeta{Namespace: testNamespace, Name: fullRoleName},
				Spec: coherence.CoherenceRoleSpec{
					Role:     roleName,
					Replicas: &replicas,
				},
				Status: coherence.CoherenceRoleStatus{
					Status:          coherence.RoleStatusReady,
					Replicas:        replicas,
					CurrentReplicas: replicas,
					ReadyReplicas:   replicas,
				},
			}

			sts = appsv1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{Namespace: testNamespace, Name: fullRoleName},
				Status: appsv1.StatefulSetStatus{
					Replicas:        replicas,
					ReadyReplicas:   replicas,
					CurrentReplicas: replicas,
				},
			}

			existing = []runtime.Object{&sts}
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

			It("should not create any CoherenceInternals", func() {
				mgr.AssertCoherenceInternals(testNamespace, 0)
			})
		})
	})

	When("a CoherenceRole is unchanged and the StatefulSet replicas has changed to the desired size", func() {
		var sts appsv1.StatefulSet
		var replicas int32 = 3

		BeforeEach(func() {
			cluster = defaultCluster

			roleCurrent = &coherence.CoherenceRole{
				ObjectMeta: metav1.ObjectMeta{Namespace: testNamespace, Name: fullRoleName},
				Spec: coherence.CoherenceRoleSpec{
					Role:     roleName,
					Replicas: &replicas,
				},
			}

			roleNew = &coherence.CoherenceRole{
				ObjectMeta: metav1.ObjectMeta{Namespace: testNamespace, Name: fullRoleName},
				Spec: coherence.CoherenceRoleSpec{
					Role:     roleName,
					Replicas: &replicas,
				},
				Status: coherence.CoherenceRoleStatus{
					Status:          coherence.RoleStatusCreated,
					Replicas:        3,
					CurrentReplicas: 2,
					ReadyReplicas:   2,
				},
			}

			sts = appsv1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{Namespace: testNamespace, Name: fullRoleName},
				Status: appsv1.StatefulSetStatus{
					Replicas:        replicas,
					ReadyReplicas:   replicas,
					CurrentReplicas: replicas,
				},
			}

			existing = []runtime.Object{&sts}
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

			It("should not create any CoherenceInternals", func() {
				mgr.AssertCoherenceInternals(testNamespace, 0)
			})

			It("should update the CoherenceRole's status replicas", func() {
				role := mgr.AssertCoherenceRoleExists(testNamespace, fullRoleName)
				Expect(role.Status.CurrentReplicas).To(Equal(sts.Status.CurrentReplicas))
				Expect(role.Status.ReadyReplicas).To(Equal(sts.Status.ReadyReplicas))
			})

			It("should update the CoherenceRole's status value", func() {
				role := mgr.AssertCoherenceRoleExists(testNamespace, fullRoleName)
				Expect(role.Status.Status).To(Equal(coherence.RoleStatusReady))
			})
		})
	})

	When("a CoherenceRole is unchanged and the StatefulSet replicas has changed but not to the desired size", func() {
		var sts appsv1.StatefulSet
		var replicas int32 = 3

		BeforeEach(func() {
			cluster = defaultCluster

			roleCurrent = &coherence.CoherenceRole{
				ObjectMeta: metav1.ObjectMeta{Namespace: testNamespace, Name: fullRoleName},
				Spec: coherence.CoherenceRoleSpec{
					Role:     roleName,
					Replicas: &replicas,
				},
			}

			roleNew = &coherence.CoherenceRole{
				ObjectMeta: metav1.ObjectMeta{Namespace: testNamespace, Name: fullRoleName},
				Spec: coherence.CoherenceRoleSpec{
					Role:     roleName,
					Replicas: &replicas,
				},
				Status: coherence.CoherenceRoleStatus{
					Status:          coherence.RoleStatusCreated,
					Replicas:        3,
					CurrentReplicas: 1,
					ReadyReplicas:   1,
				},
			}

			sts = appsv1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{Namespace: testNamespace, Name: fullRoleName},
				Status: appsv1.StatefulSetStatus{
					Replicas:        replicas,
					ReadyReplicas:   2,
					CurrentReplicas: 2,
				},
			}

			existing = []runtime.Object{&sts}
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

			It("should not create any CoherenceInternals", func() {
				mgr.AssertCoherenceInternals(testNamespace, 0)
			})

			It("should update the CoherenceRole's status replica counts", func() {
				role := mgr.AssertCoherenceRoleExists(testNamespace, fullRoleName)
				Expect(role.Status.CurrentReplicas).To(Equal(sts.Status.CurrentReplicas))
				Expect(role.Status.ReadyReplicas).To(Equal(sts.Status.ReadyReplicas))
			})

			It("should not update the CoherenceRole's status value", func() {
				role := mgr.AssertCoherenceRoleExists(testNamespace, fullRoleName)
				Expect(role.Status.Status).To(Equal(roleNew.Status.Status))
			})
		})
	})
})
