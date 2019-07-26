package coherencerole

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
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
var _ = Describe("CoherenceRole Controller", func() {
	const (
		namespace    = "coherence"
		clusterName  = "test-cluster"
		roleName     = "storage"
		fullRoleName = clusterName + "-" + roleName
	)
	var (
		cluster *coherence.CoherenceCluster
		role    *coherence.CoherenceRole
		req     reconcile.Request
	)

	BeforeEach(func() {
		// create the test CoherenceCluster
		cluster = &coherence.CoherenceCluster{
			ObjectMeta: metav1.ObjectMeta{
				Name:      clusterName,
				Namespace: namespace,
				Labels:    map[string]string{coherence.CoherenceClusterLabel: clusterName},
			},
		}

		// create the test CoherenceRole
		role = &coherence.CoherenceRole{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fullRoleName,
				Namespace: namespace,
				Labels:    map[string]string{coherence.CoherenceClusterLabel: clusterName},
			},
			Spec: coherence.CoherenceRoleSpec{
				RoleName: roleName,
			},
		}

		// create the test reconcile request
		req = reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name:      fullRoleName,
				Namespace: namespace,
			},
		}
	})

	Context("when calling Reconcile", func() {
		It("should return and not requeue if no parent CoherenceCluster", func() {
			mgr := stubs.NewFakeManager(role)
			r   := newReconciler(mgr)

			res, err := r.Reconcile(req)

			Expect(res).To(Equal(reconcile.Result{Requeue: false}))
			Expect(err).To(BeNil())

			// Assert that no CoherenceInternal was created for the role
			_, err = r.getCoherenceInternal(role)
			Expect(err).To(Not(BeNil()))
			Expect(errors.IsNotFound(err)).To(BeTrue())

			// Assert that a failure event was raised for the role
			event, found := mgr.NextEvent()
			Expect(found).To(BeTrue())
			Expect(event.Object.GetObjectKind()).To(Equal(role.GetObjectKind()))
			Expect(event.Type).To(Equal(corev1.EventTypeNormal))
			Expect(event.Reason).To(Equal("Failed"))
			Expect(event.Message).To(HavePrefix("Invalid CoherenceRole"))

			// Assert that the role has an error status
			roleUpdated := coherence.CoherenceRole{}
			err = mgr.Client.Get(context.TODO(), types.NamespacedName{Namespace: role.Namespace, Name: role.Name}, &roleUpdated)
			if err != nil {
				Fail(err.Error())
			}
			Expect(roleUpdated.Status.Status).To(Equal(coherence.RoleStatusFailed))
		})


		It("should create new CoherenceInternal if one does not exist", func() {
			mgr := stubs.NewFakeManager(cluster, role)
			r   := newReconciler(mgr)

			res, err := r.Reconcile(req)

			Expect(err).To(BeNil())
			Expect(res).To(Equal(reconcile.Result{}))

			// Assert that there is a CoherenceInternal for the role
			cohInt, err := r.getCoherenceInternal(role)
			Expect(err).To(BeNil())
			Expect(cohInt).To(Not(BeNil()))

			// Assert fields of CoherenceInternal e.g. assert labels etc...
			labels := cohInt.GetLabels();
			Expect(labels).To(Not(BeNil()))
			Expect(labels[coherence.CoherenceRoleLabel]).To(Equal(roleName))

			// Assert that a success event was raised for the role
			event, found := mgr.NextEvent()
			Expect(found).To(BeTrue())
			Expect(event.Object.GetObjectKind()).To(Equal(role.GetObjectKind()))
			Expect(event.Type).To(Equal(corev1.EventTypeNormal))
			Expect(event.Reason).To(Equal("SuccessfulCreate"))
			Expect(event.Message).To(HavePrefix("create Helm install"))

			// Assert that the role has an created status
			roleUpdated := coherence.CoherenceRole{}
			err = mgr.Client.Get(context.TODO(), types.NamespacedName{Namespace: role.Namespace, Name: role.Name}, &roleUpdated)
			if err != nil {
				Fail(err.Error())
			}
			Expect(roleUpdated.Status.Status).To(Equal(coherence.RoleStatusCreated))
			Expect(roleUpdated.Status.Replicas).To(Equal(role.Spec.GetReplicas()))
			Expect(roleUpdated.Status.CurrentReplicas).To(BeZero())
			Expect(roleUpdated.Status.ReadyReplicas).To(BeZero())
			Expect(roleUpdated.Status.Selector).To(Equal(fmt.Sprintf(selectorTemplate, cluster.Name, role.Spec.RoleName)))
		})
	})
})

