/*
 * Copyright (c) 2019, 2020 Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package coherencerole

import (
	"context"
	"fmt"
	"github.com/go-test/deep"
	"github.com/oracle/coherence-operator/pkg/operator"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
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
		statefulSet *appsv1.StatefulSet
		existing    []runtime.Object
		result      stubs.ReconcileResult
		err         error

		controller *ReconcileCoherenceRole

		defaultCluster = &coherence.CoherenceCluster{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: testNamespace,
				Name:      testClusterName,
			},
		}
	)

	JustBeforeEach(func() {
		mgr, err = stubs.NewFakeManager(existing...)
		Expect(err).NotTo(HaveOccurred())
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

			cohIntern := controller.CreateHelmValues(cluster, roleNew, specMap)

			err = mgr.Client.Create(context.TODO(), cohIntern)
			Expect(err).NotTo(HaveOccurred())
		}

		if statefulSet != nil {
			_ = mgr.Client.Create(context.TODO(), statefulSet)
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
				msg := fmt.Sprintf(createMessage, roleNew.Name, roleNew.Name)
				event := mgr.AssertEvent()

				Expect(event.Type).To(Equal(corev1.EventTypeNormal))
				Expect(event.Reason).To(Equal(eventReasonCreated))
				Expect(event.Message).To(Equal(msg))

				mgr.AssertNoRemainingEvents()
			})

			It("should create a CoherenceInternal", func() {
				u, err := mgr.AssertCoherenceInternalExists(testNamespace, fullRoleName)
				Expect(err).NotTo(HaveOccurred())
				roleSpec := UnstructuredToCoherenceInternalSpec(u)
				expected := coherence.NewCoherenceInternalSpec(cluster, roleNew)
				expected.EnsureCoherenceImage(operator.GetDefaultCoherenceImage())
				expected.EnsureCoherenceUtilsImage(operator.GetDefaultCoherenceUtilsImage())

				Expect(roleSpec).To(Equal(expected))
			})
		})
	})

	When("a CoherenceRole is updated", func() {
		BeforeEach(func() {
			imageOrig := "coherence:1.0"
			imageNew := "coherence:2.0"
			utilsImage := "foo/bar:1.0"

			var replicas int32 = 3

			roleCurrent = &coherence.CoherenceRole{
				ObjectMeta: metav1.ObjectMeta{Namespace: testNamespace, Name: fullRoleName},
				Spec: coherence.CoherenceRoleSpec{
					Role:     roleName,
					Replicas: &replicas,
					Coherence: &coherence.CoherenceSpec{
						ImageSpec: coherence.ImageSpec{
							Image: &imageOrig,
						},
					},
					CoherenceUtils: &coherence.ImageSpec{
						Image: &utilsImage,
					},
				},
			}

			roleNew = &coherence.CoherenceRole{
				ObjectMeta: metav1.ObjectMeta{Namespace: testNamespace, Name: fullRoleName},
				Spec: coherence.CoherenceRoleSpec{
					Role:     roleName,
					Replicas: &replicas,
					Coherence: &coherence.CoherenceSpec{
						ImageSpec: coherence.ImageSpec{
							Image: &imageNew,
						},
					},
					CoherenceUtils: &coherence.ImageSpec{
						Image: &utilsImage,
					},
				},
			}

			// The cluster would have the new role state
			cluster = defaultCluster.DeepCopy()
			roleNew.Spec.DeepCopyInto(&cluster.Spec.CoherenceRoleSpec)

			statefulSet = &appsv1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{Namespace: testNamespace, Name: fullRoleName},
				Spec: appsv1.StatefulSetSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							InitContainers: []corev1.Container{
								{
									Name:  operator.CoherenceUtilsContainerName,
									Image: utilsImage,
								},
							},
							Containers: []corev1.Container{
								{
									Name:  operator.CoherenceContainerName,
									Image: imageOrig,
								},
							},
						},
					},
				},
				Status: appsv1.StatefulSetStatus{
					Replicas:        replicas,
					ReadyReplicas:   replicas,
					CurrentReplicas: replicas,
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
				msg := fmt.Sprintf(updateMessage, roleNew.Name, roleNew.Name)
				event := mgr.AssertEvent()

				Expect(event.Type).To(Equal(corev1.EventTypeNormal))
				Expect(event.Reason).To(Equal(eventReasonUpdated))
				Expect(event.Message).To(Equal(msg))

				mgr.AssertNoRemainingEvents()
			})

			It("should update the CoherenceInternal", func() {
				u, err := mgr.AssertCoherenceInternalExists(testNamespace, fullRoleName)
				Expect(err).NotTo(HaveOccurred())
				roleSpec := UnstructuredToCoherenceInternalSpec(u)
				expected := coherence.NewCoherenceInternalSpec(cluster, roleNew)
				expected.EnsureCoherenceImage(operator.GetDefaultCoherenceImage())
				expected.EnsureCoherenceUtilsImage(operator.GetDefaultCoherenceUtilsImage())
				Expect(roleSpec).To(Equal(expected), fmt.Sprintf("Expected roles to match:\n%s", deep.Equal(roleSpec, expected)))
			})
		})
	})

	When("a CoherenceRole is unchanged and the StatefulSet is unchanged", func() {
		BeforeEach(func() {
			var replicas int32 = 3
			var image = "foo/bar:1.0"

			roleCurrent = &coherence.CoherenceRole{
				ObjectMeta: metav1.ObjectMeta{Namespace: testNamespace, Name: fullRoleName},
				Spec: coherence.CoherenceRoleSpec{
					Role:     roleName,
					Replicas: &replicas,
					CoherenceUtils: &coherence.ImageSpec{
						Image: &image,
					},
					Coherence: &coherence.CoherenceSpec{
						ImageSpec: coherence.ImageSpec{
							Image: &image,
						},
					},
				},
			}

			cluster = &coherence.CoherenceCluster{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: testNamespace,
					Name:      testClusterName,
				},
				Spec: coherence.CoherenceClusterSpec{Roles: []coherence.CoherenceRoleSpec{
					{
						Role:     roleName,
						Replicas: &replicas,
						CoherenceUtils: &coherence.ImageSpec{
							Image: &image,
						},
						Coherence: &coherence.CoherenceSpec{
							ImageSpec: coherence.ImageSpec{
								Image: &image,
							},
						},
					},
				}},
			}

			roleNew = &coherence.CoherenceRole{
				ObjectMeta: metav1.ObjectMeta{Namespace: testNamespace, Name: fullRoleName},
				Spec: coherence.CoherenceRoleSpec{
					Role:     roleName,
					Replicas: &replicas,
					CoherenceUtils: &coherence.ImageSpec{
						Image: &image,
					},
					Coherence: &coherence.CoherenceSpec{
						ImageSpec: coherence.ImageSpec{
							Image: &image,
						},
					},
				},
				Status: coherence.CoherenceRoleStatus{
					Status:          coherence.RoleStatusReady,
					Replicas:        replicas,
					CurrentReplicas: replicas,
					ReadyReplicas:   replicas,
				},
			}

			statefulSet = &appsv1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{Namespace: testNamespace, Name: fullRoleName},
				Spec: appsv1.StatefulSetSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							InitContainers: []corev1.Container{
								{
									Name:  operator.CoherenceUtilsContainerName,
									Image: image,
								},
							},
							Containers: []corev1.Container{
								{
									Name:  operator.CoherenceContainerName,
									Image: image,
								},
							},
						},
					},
				},
				Status: appsv1.StatefulSetStatus{
					Replicas:        replicas,
					ReadyReplicas:   replicas,
					CurrentReplicas: replicas,
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

			It("should not fire an event", func() {
				_, found := mgr.NextEvent()
				Expect(found).To(BeFalse())
			})

			It("should create one CoherenceInternals", func() {
				mgr.AssertCoherenceInternals(testNamespace, 1)
			})
		})
	})

	When("a CoherenceRole is unchanged and the StatefulSet replicas has changed to the desired size", func() {
		var replicas int32 = 3
		var image = "foo/bar:1.0"

		BeforeEach(func() {
			cluster = defaultCluster

			roleCurrent = &coherence.CoherenceRole{
				ObjectMeta: metav1.ObjectMeta{Namespace: testNamespace, Name: fullRoleName},
				Spec: coherence.CoherenceRoleSpec{
					Role:     roleName,
					Replicas: &replicas,
					CoherenceUtils: &coherence.ImageSpec{
						Image: &image,
					},
					Coherence: &coherence.CoherenceSpec{
						ImageSpec: coherence.ImageSpec{
							Image: &image,
						},
					},
				},
			}

			roleNew = &coherence.CoherenceRole{
				ObjectMeta: metav1.ObjectMeta{Namespace: testNamespace, Name: fullRoleName},
				Spec: coherence.CoherenceRoleSpec{
					Role:     roleName,
					Replicas: &replicas,
					CoherenceUtils: &coherence.ImageSpec{
						Image: &image,
					},
					Coherence: &coherence.CoherenceSpec{
						ImageSpec: coherence.ImageSpec{
							Image: &image,
						},
					},
				},
				Status: coherence.CoherenceRoleStatus{
					Status:          coherence.RoleStatusCreated,
					Replicas:        3,
					CurrentReplicas: 2,
					ReadyReplicas:   2,
				},
			}

			// The cluster would have the new role state
			cluster = defaultCluster.DeepCopy()
			roleNew.Spec.DeepCopyInto(&cluster.Spec.CoherenceRoleSpec)

			statefulSet = &appsv1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{Namespace: testNamespace, Name: fullRoleName},
				Spec: appsv1.StatefulSetSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							InitContainers: []corev1.Container{
								{
									Name:  operator.CoherenceUtilsContainerName,
									Image: image,
								},
							},
							Containers: []corev1.Container{
								{
									Name:  operator.CoherenceContainerName,
									Image: image,
								},
							},
						},
					},
				},
				Status: appsv1.StatefulSetStatus{
					Replicas:        replicas,
					ReadyReplicas:   replicas,
					CurrentReplicas: replicas,
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

			It("should not fire an event", func() {
				_, found := mgr.NextEvent()
				Expect(found).To(BeFalse())
			})

			It("should create one CoherenceInternals", func() {
				mgr.AssertCoherenceInternals(testNamespace, 1)
			})

			It("should update the CoherenceRole's status replicas", func() {
				role := mgr.AssertCoherenceRoleExists(testNamespace, fullRoleName)
				Expect(role.Status.CurrentReplicas).To(Equal(statefulSet.Status.CurrentReplicas))
				Expect(role.Status.ReadyReplicas).To(Equal(statefulSet.Status.ReadyReplicas))
			})

			It("should update the CoherenceRole's status value", func() {
				role := mgr.AssertCoherenceRoleExists(testNamespace, fullRoleName)
				Expect(role.Status.Status).To(Equal(coherence.RoleStatusReady))
			})
		})
	})

	When("a CoherenceRole is unchanged and the StatefulSet replicas has changed but not to the desired size", func() {
		var replicas int32 = 3
		var image = "foo/bar:1.0"

		BeforeEach(func() {
			roleCurrent = &coherence.CoherenceRole{
				ObjectMeta: metav1.ObjectMeta{Namespace: testNamespace, Name: fullRoleName},
				Spec: coherence.CoherenceRoleSpec{
					Role:     roleName,
					Replicas: &replicas,
					CoherenceUtils: &coherence.ImageSpec{
						Image: &image,
					},
					Coherence: &coherence.CoherenceSpec{
						ImageSpec: coherence.ImageSpec{
							Image: &image,
						},
					},
				},
			}

			cluster = defaultCluster.DeepCopy()
			cluster.Spec.Roles = []coherence.CoherenceRoleSpec{roleCurrent.Spec}

			roleNew = &coherence.CoherenceRole{
				ObjectMeta: metav1.ObjectMeta{Namespace: testNamespace, Name: fullRoleName},
				Spec: coherence.CoherenceRoleSpec{
					Role:     roleName,
					Replicas: &replicas,
					CoherenceUtils: &coherence.ImageSpec{
						Image: &image,
					},
					Coherence: &coherence.CoherenceSpec{
						ImageSpec: coherence.ImageSpec{
							Image: &image,
						},
					},
				},
				Status: coherence.CoherenceRoleStatus{
					Status:          coherence.RoleStatusCreated,
					Replicas:        3,
					CurrentReplicas: 1,
					ReadyReplicas:   1,
				},
			}

			statefulSet = &appsv1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{Namespace: testNamespace, Name: fullRoleName},
				Spec: appsv1.StatefulSetSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							InitContainers: []corev1.Container{
								{
									Name:  operator.CoherenceUtilsContainerName,
									Image: image,
								},
							},
							Containers: []corev1.Container{
								{
									Name:  operator.CoherenceContainerName,
									Image: image,
								},
							},
						},
					},
				},
				Status: appsv1.StatefulSetStatus{
					Replicas:        replicas,
					ReadyReplicas:   2,
					CurrentReplicas: 2,
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

			It("should not fire an event", func() {
				_, found := mgr.NextEvent()
				Expect(found).To(BeFalse())
			})

			It("should create one CoherenceInternals", func() {
				mgr.AssertCoherenceInternals(testNamespace, 1)
			})

			It("should update the CoherenceRole's status replica counts", func() {
				role := mgr.AssertCoherenceRoleExists(testNamespace, fullRoleName)
				Expect(role.Status.CurrentReplicas).To(Equal(statefulSet.Status.CurrentReplicas))
				Expect(role.Status.ReadyReplicas).To(Equal(statefulSet.Status.ReadyReplicas))
			})

			It("should not update the CoherenceRole's status value", func() {
				role := mgr.AssertCoherenceRoleExists(testNamespace, fullRoleName)
				Expect(role.Status.Status).To(Equal(coherence.RoleStatusCreated))
			})
		})
	})
})
