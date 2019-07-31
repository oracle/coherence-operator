package coherencerole

import (
	"context"
	"fmt"
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
		roleCurrent *coherence.CoherenceInternalSpec
		existing    []runtime.Object
		result      stubs.ReconcileResult

		defaultCluster = &coherence.CoherenceCluster{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: testNamespace,
				Name:      testClusterName,
			},
		}
	)

	JustBeforeEach(func() {
		mgr = stubs.NewFakeManager(existing...)

		controller := newReconciler(mgr)

		if cluster != nil {
			_ = mgr.Client.Create(context.TODO(), cluster)
		}

		if roleNew != nil {
			_ = mgr.Client.Create(context.TODO(), roleNew)
		}

		if roleCurrent != nil {
			spec, err := coherence.CoherenceInternalSpecAsMapFromSpec(*roleCurrent)
			Expect(err).NotTo(HaveOccurred())

			cohIntern, err := controller.CreateCoherenceInternal(cluster, roleNew, spec)
			Expect(err).NotTo(HaveOccurred())

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

	When("a CoherenceRole does not exist", func() {
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

				//Expect(event.Owner).To(Equal(roleNew))
				Expect(event.Type).To(Equal(corev1.EventTypeNormal))
				Expect(event.Reason).To(Equal(eventReasonCreated))
				Expect(event.Message).To(Equal(msg))

				mgr.AssertNoRemainingEvents()
			})

			It("should create a CoherenceInternals", func() {
				u := mgr.AssertCoherenceInternalExists(testNamespace, fullRoleName)
				Expect(u).To(Not(BeNil()))

			})
		})
	})

	When("a CoherenceRole is updated", func() {
		BeforeEach(func() {
			cluster = defaultCluster

			imageOrig := "coherence:1.0"
			imageNew := "coherence:2.0"
			var replicas int32 = 3

			roleCurrent = &coherence.CoherenceInternalSpec{
				ClusterSize: replicas,
				Coherence: &coherence.ImageSpec{
					Image: &imageOrig,
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

				//Expect(event.Owner).To(Equal(roleNew))
				Expect(event.Type).To(Equal(corev1.EventTypeNormal))
				Expect(event.Reason).To(Equal(eventReasonUpdated))
				Expect(event.Message).To(Equal(msg))

				mgr.AssertNoRemainingEvents()
			})

			It("should update the CoherenceInternals", func() {
				u := mgr.AssertCoherenceInternalExists(testNamespace, fullRoleName)
				Expect(u).To(Not(BeNil()))

				//spec, err := coherence.CoherenceInternalSpecAsMapFromSpec(*roleCurrent)
				//Expect(err).NotTo(HaveOccurred())
				//
				//Expect(u.Object["spec"]).To(Equal(spec))
			})
		})
	})
})
