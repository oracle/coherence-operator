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
		roleCurrent *coherence.CoherenceRole
		existing    []runtime.Object
		result      stubs.ReconcileResult
		errors      stubs.ClientErrors

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
				Namespace: testNamespace,
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

	When("The client will return an error getting the CoherenceRole", func() {
		var err error = stubs.FakeError{Msg: "error getting role"}

		BeforeEach(func() {
			cluster = defaultCluster
			roleNew = nil
			errors.AddGetError(stubs.ErrorIf{KeyIs: &types.NamespacedName{Namespace: testNamespace, Name: fullRoleName}}, err)
		})

		When("reconcile is called", func() {
			It("should return error", func() {
				Expect(result.Error).To(HaveOccurred())
				Expect(result.Error).To(Equal(err))
			})

			It("should re-queue the request", func() {
				Expect(result.Result).To(Equal(reconcile.Result{Requeue: true}))
			})

			It("should not fire an event", func() {
				mgr.AssertNoRemainingEvents()
			})

			It("should not create any CoherenceInternals", func() {
				mgr.AssertCoherenceInternals(testNamespace, 0)
			})
		})
	})

	When("the the parent CoherenceCluster does not exist", func() {
		BeforeEach(func() {
			cluster = nil

			roleNew = &coherence.CoherenceRole{
				ObjectMeta: metav1.ObjectMeta{Namespace: testNamespace, Name: fullRoleName},
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
				msg := fmt.Sprintf(invalidRoleEventMessage, fullRoleName, testClusterName)
				event := mgr.AssertEvent()

				Expect(event.Type).To(Equal(corev1.EventTypeNormal))
				Expect(event.Reason).To(Equal(eventReasonFailed))
				Expect(event.Message).To(Equal(msg))

				mgr.AssertNoRemainingEvents()
			})

			It("should update the CoherenceRole's status to failed", func() {
				role := mgr.AssertCoherenceRoleExists(testNamespace, fullRoleName)
				Expect(role.Status.Status).To(Equal(coherence.RoleStatusFailed))
			})

			It("should not create any CoherenceInternals", func() {
				mgr.AssertCoherenceInternals(testNamespace, 0)
			})
		})
	})

	When("the k8s client returns an error getting the parent CoherenceCluster", func() {
		var err error = stubs.FakeError{Msg: "error getting cluster"}

		BeforeEach(func() {
			cluster = defaultCluster

			roleNew = &coherence.CoherenceRole{
				ObjectMeta: metav1.ObjectMeta{Namespace: testNamespace, Name: fullRoleName},
				Spec: coherence.CoherenceRoleSpec{
					Role: roleName,
				},
			}

			errors.AddGetError(stubs.ErrorIf{KeyIs: &types.NamespacedName{Namespace: testNamespace, Name: fullRoleName}}, err)
		})

		When("reconcile is called", func() {
			It("should return error", func() {
				Expect(result.Error).To(HaveOccurred())
				Expect(result.Error).To(Equal(err))
			})

			It("should re-queue the request", func() {
				Expect(result.Result).To(Equal(reconcile.Result{Requeue: true}))
			})

			It("should not fire an event", func() {
				mgr.AssertNoRemainingEvents()
			})

			It("should not create any CoherenceInternals", func() {
				mgr.AssertCoherenceInternals(testNamespace, 0)
			})
		})
	})
})
