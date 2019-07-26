package v1_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	coherence "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	v1 "k8s.io/api/core/v1"
)

// The tests for CoherenceInternal
var _ = Describe("coherenceinternal_types", func() {
	const (
		namespace    = "test-ns"
		clusterName  = "test-cluster"
		roleName     = "storage"
		fullRoleName = "test-cluster-" + roleName
	)

	var (
		cluster  coherence.CoherenceCluster
		role     coherence.CoherenceRole
		expected coherence.CoherenceInternalSpec
	)

	BeforeEach(func() {
		cluster = coherence.CoherenceCluster{}
		cluster.Namespace = namespace
		cluster.Name = clusterName

		role = coherence.CoherenceRole{}
		role.Namespace = namespace
		role.Name = fullRoleName
		role.Spec.RoleName = roleName

		expected = coherence.CoherenceInternalSpec{
			FullnameOverride: fullRoleName,
			Cluster:          clusterName,
			ClusterSize:      coherence.DefaultReplicas,
			Role:             roleName,
			Store: &coherence.CoherenceInternalStoreSpec{
				WKA:    cluster.GetWkaServiceName(),
				Labels: &map[string]string{coherence.CoherenceRoleLabel: roleName},
			},
		}
	})

	Describe("when calling NewCoherenceInternalSpec", func() {
		AssertFields := func() {
			It("should have the correct populated fields", func() {
				spec := coherence.NewCoherenceInternalSpec(&cluster, &role)
				Expect(spec).To(Equal(expected))
			})
		}

		Context("from basic cluster spec and role spec", func() {
			BeforeEach(func() {
				AssertFields()
			})
		})

		Context("with serviceAccount set in the cluster spec", func() {
			BeforeEach(func() {
				cluster.Spec.ServiceAccountName = "testAccount"

				expected.ServiceAccountName = cluster.Spec.ServiceAccountName

				AssertFields()
			})
		})

		Context("with image pull secrets set in the cluster spec", func() {
			BeforeEach(func() {
				cluster.Spec.ImagePullSecrets = []string{"foo", "bar"}

				expected.ImagePullSecrets = cluster.Spec.ImagePullSecrets

				AssertFields()
			})
		})

		Context("with Coherence image set in the cluster spec", func() {
			BeforeEach(func() {
				cluster.Spec.Images.Coherence = &coherence.ImageSpec{
					Image: "oracle/coherence:1.2.3",
				}

				expected.Coherence = cluster.Spec.Images.Coherence

				AssertFields()
			})
		})

		Context("with Coherence image pull policy set in the cluster spec", func() {
			BeforeEach(func() {
				cluster.Spec.Images.Coherence = &coherence.ImageSpec{
					ImagePullPolicy: v1.PullAlways,
				}

				expected.Coherence = cluster.Spec.Images.Coherence

				AssertFields()
			})
		})

		Context("with Coherence image set in the role spec", func() {
			BeforeEach(func() {
				role.Spec.Images.Coherence = &coherence.ImageSpec{
					Image: "oracle/coherence:1.2.3",
				}

				expected.Coherence = role.Spec.Images.Coherence

				AssertFields()
			})
		})

		Context("with Coherence image pull policy set in the role spec", func() {
			BeforeEach(func() {
				cluster.Spec.Images.Coherence = &coherence.ImageSpec{
					ImagePullPolicy: v1.PullAlways,
				}

				expected.Coherence = role.Spec.Images.Coherence

				AssertFields()
			})
		})

		Context("with Coherence image set in the cluster spec and role spec, role takes precedence", func() {
			BeforeEach(func() {
				cluster.Spec.Images.Coherence = &coherence.ImageSpec{
					Image: "oracle/coherence:1.2.3",
				}
				cluster.Spec.Images.Coherence = &coherence.ImageSpec{
					Image: "oracle/coherence:5.6.7",
				}

				expected.Coherence = role.Spec.Images.Coherence

				AssertFields()
			})
		})

		Context("with Coherence image pull policy set in the cluster spec and role spec, role takes precedence", func() {
			BeforeEach(func() {
				cluster.Spec.Images.Coherence = &coherence.ImageSpec{
					ImagePullPolicy: v1.PullAlways,
				}
				role.Spec.Images.Coherence = &coherence.ImageSpec{
					ImagePullPolicy: v1.PullNever,
				}

				expected.Coherence = role.Spec.Images.Coherence

				AssertFields()
			})
		})

		Context("with role labels set in the cluster spec", func() {
			BeforeEach(func() {
				clusterLabels := map[string]string{"key1": "value1", "key2": "value2"}
				expectedLabels := map[string]string{"key1": "value1",
					"key2": "value2", coherence.CoherenceRoleLabel: roleName}

				cluster.Spec.Labels = &clusterLabels
				expected.Store.Labels = &expectedLabels

				AssertFields()
			})
		})

		Context("with role labels set in the role spec", func() {
			BeforeEach(func() {
				roleLabels := map[string]string{"key1": "value1", "key2": "value2"}
				expectedLabels := map[string]string{"key1": "value1",
					"key2": "value2", coherence.CoherenceRoleLabel: roleName}

				role.Spec.Labels = &roleLabels
				expected.Store.Labels = &expectedLabels

				AssertFields()
			})
		})

		Context("with role labels set in the cluster spec and role spec, role takes precedence", func() {
			BeforeEach(func() {
				clusterLabels := map[string]string{"key1": "value1", "key2": "value2"}
				roleLabels := map[string]string{"key1": "foo", "key3": "value3"}
				expectedLabels := map[string]string{"key1": "foo", "key2": "value2",
					"key3": "value3", coherence.CoherenceRoleLabel: roleName}

				cluster.Spec.Labels = &clusterLabels
				role.Spec.Labels = &roleLabels
				expected.Store.Labels = &expectedLabels

				AssertFields()
			})
		})

		Context("with readiness timeout set in the cluster spec", func() {
			BeforeEach(func() {
				cluster.Spec.ReadinessProbe = &coherence.ReadinessProbeSpec{
					TimeoutSeconds: int32Pointer(100),
				}

				expected.Store.ReadinessProbe = cluster.Spec.ReadinessProbe

				AssertFields()
			})
		})

		Context("with readiness timeout set in the role spec", func() {
			BeforeEach(func() {
				role.Spec.ReadinessProbe = &coherence.ReadinessProbeSpec{
					TimeoutSeconds: int32Pointer(100),
				}

				expected.Store.ReadinessProbe = role.Spec.ReadinessProbe

				AssertFields()
			})
		})

		Context("with readiness timeout set in the cluster spec and role spec, role takes precedence", func() {
			BeforeEach(func() {
				cluster.Spec.ReadinessProbe = &coherence.ReadinessProbeSpec{
					TimeoutSeconds: int32Pointer(100),
				}
				role.Spec.ReadinessProbe = &coherence.ReadinessProbeSpec{
					TimeoutSeconds: int32Pointer(99),
				}

				expected.Store.ReadinessProbe = role.Spec.ReadinessProbe

				AssertFields()
			})
		})

		Context("with readiness period set in the cluster spec", func() {
			BeforeEach(func() {
				cluster.Spec.ReadinessProbe = &coherence.ReadinessProbeSpec{
					PeriodSeconds: int32Pointer(100),
				}

				expected.Store.ReadinessProbe = cluster.Spec.ReadinessProbe

				AssertFields()
			})
		})

		Context("with readiness period set in the role spec", func() {
			BeforeEach(func() {
				role.Spec.ReadinessProbe = &coherence.ReadinessProbeSpec{
					PeriodSeconds: int32Pointer(100),
				}

				expected.Store.ReadinessProbe = role.Spec.ReadinessProbe

				AssertFields()
			})
		})

		Context("with readiness period set in the cluster spec and role spec, role takes precedence", func() {
			BeforeEach(func() {
				cluster.Spec.ReadinessProbe = &coherence.ReadinessProbeSpec{
					PeriodSeconds: int32Pointer(100),
				}
				role.Spec.ReadinessProbe = &coherence.ReadinessProbeSpec{
					PeriodSeconds: int32Pointer(99),
				}

				expected.Store.ReadinessProbe = role.Spec.ReadinessProbe

				AssertFields()
			})
		})

		Context("with readiness success threshold set in the cluster spec", func() {
			BeforeEach(func() {
				cluster.Spec.ReadinessProbe = &coherence.ReadinessProbeSpec{
					SuccessThreshold: int32Pointer(100),
				}

				expected.Store.ReadinessProbe = cluster.Spec.ReadinessProbe

				AssertFields()
			})
		})

		Context("with readiness success threshold set in the role spec", func() {
			BeforeEach(func() {
				role.Spec.ReadinessProbe = &coherence.ReadinessProbeSpec{
					SuccessThreshold: int32Pointer(100),
				}

				expected.Store.ReadinessProbe = role.Spec.ReadinessProbe

				AssertFields()
			})
		})

		Context("with readiness success threshold set in the cluster spec and role spec, role takes precedence", func() {
			BeforeEach(func() {
				cluster.Spec.ReadinessProbe = &coherence.ReadinessProbeSpec{
					SuccessThreshold: int32Pointer(100),
				}
				role.Spec.ReadinessProbe = &coherence.ReadinessProbeSpec{
					SuccessThreshold: int32Pointer(99),
				}

				expected.Store.ReadinessProbe = role.Spec.ReadinessProbe

				AssertFields()
			})
		})

		Context("with readiness failure threshold set in the cluster spec", func() {
			BeforeEach(func() {
				cluster.Spec.ReadinessProbe = &coherence.ReadinessProbeSpec{
					FailureThreshold: int32Pointer(100),
				}

				expected.Store.ReadinessProbe = cluster.Spec.ReadinessProbe

				AssertFields()
			})
		})

		Context("with readiness failure threshold set in the role spec", func() {
			BeforeEach(func() {
				role.Spec.ReadinessProbe = &coherence.ReadinessProbeSpec{
					FailureThreshold: int32Pointer(100),
				}

				expected.Store.ReadinessProbe = role.Spec.ReadinessProbe

				AssertFields()
			})
		})

		Context("with readiness failure threshold set in the cluster spec and role spec, role takes precedence", func() {
			BeforeEach(func() {
				cluster.Spec.ReadinessProbe = &coherence.ReadinessProbeSpec{
					FailureThreshold: int32Pointer(100),
				}
				role.Spec.ReadinessProbe = &coherence.ReadinessProbeSpec{
					FailureThreshold: int32Pointer(99),
				}

				expected.Store.ReadinessProbe = role.Spec.ReadinessProbe

				AssertFields()
			})
		})

		Context("with ", func() {
			BeforeEach(func() {

				AssertFields()
			})
		})

	})
})
