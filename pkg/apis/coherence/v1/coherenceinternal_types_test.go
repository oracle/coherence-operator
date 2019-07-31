package v1_test

import (
	"encoding/json"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"

	"github.com/oracle/coherence-operator/pkg/apis"
	coherence "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"

	"k8s.io/apimachinery/pkg/runtime/schema"

	"reflect"
)

var _ = Describe("Testing CoherenceInternal struct", func() {
	var (
		cluster coherence.CoherenceCluster
		role    coherence.CoherenceRole
	)

	BeforeEach(func() {
		cluster = coherence.CoherenceCluster{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "test-namespace",
				Name:      "test-cluster",
			},
			Spec: coherence.CoherenceClusterSpec{
				ImagePullSecrets:   []string{"test-secret"},
				ServiceAccountName: "foo-account",
				CoherenceRoleSpec:  coherence.CoherenceRoleSpec{},
				Roles:              nil,
			},
		}

		safeScaling := coherence.SafeScaling
		always := corev1.PullAlways
		ifNotPresent := corev1.PullIfNotPresent

		// Fully populated CoherenceRole
		role = coherence.CoherenceRole{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "test-namespace",
				Name:      "test-cluster-storage",
			},
			Spec: coherence.CoherenceRoleSpec{
				Role:     "storage",
				Replicas: int32Ptr(5),
				Images: &coherence.Images{
					Coherence: &coherence.ImageSpec{
						Image:           stringPtr("coherence:1.0"),
						ImagePullPolicy: &ifNotPresent,
					},
					CoherenceUtils: &coherence.ImageSpec{
						Image:           stringPtr("coherence-utils:1.0"),
						ImagePullPolicy: &always,
					},
					UserArtifacts: &coherence.UserArtifactsImageSpec{
						ImageSpec: coherence.ImageSpec{
							Image:           stringPtr("custom:1.0"),
							ImagePullPolicy: &always,
						},
						LibDir:    stringPtr("/lib"),
						ConfigDir: stringPtr("/conf"),
					},
					Fluentd: &coherence.FluentdImageSpec{
						ImageSpec: coherence.ImageSpec{
							Image:           stringPtr("fluentd:1.0"),
							ImagePullPolicy: &ifNotPresent,
						},
						Application: &coherence.FluentdApplicationSpec{
							ConfigFile: stringPtr("fluent.yaml"),
							Tag:        stringPtr("fluentd-tag"),
						},
					},
				},
				StorageEnabled: boolPtr(false),
				ScalingPolicy:  &safeScaling,
				ReadinessProbe: &coherence.ReadinessProbeSpec{
					InitialDelaySeconds: int32Ptr(10),
					TimeoutSeconds:      int32Ptr(20),
					PeriodSeconds:       int32Ptr(30),
					SuccessThreshold:    int32Ptr(40),
					FailureThreshold:    int32Ptr(50),
				},
				Labels: mapPtr(map[string]string{"one": "1", "two": "2"}),
			},
		}
	})

	Context("Creating a CoherenceInternal from a CoherenceCluster and CoherenceRole", func() {
		var (
			result coherence.CoherenceInternalSpec
		)

		JustBeforeEach(func() {
			result = coherence.NewCoherenceInternalSpec(&cluster, &role)
		})

		When("all fields are set", func() {
			It("should set the FullNameOverride field", func() {
				Expect(result.FullnameOverride).To(Equal("test-cluster-storage"))
			})

			It("should set the ClusterSize", func() {
				var expected int32 = 5
				Expect(result.ClusterSize).To(Equal(expected))
			})

			It("should set the Cluster", func() {
				Expect(result.Cluster).To(Equal("test-cluster"))
			})

			It("should set the ServiceAccountName", func() {
				Expect(result.ServiceAccountName).To(Equal("foo-account"))
			})

			It("should set the ImagePullSecrets", func() {
				Expect(result.ImagePullSecrets).To(Equal(cluster.Spec.ImagePullSecrets))
			})

			It("should set the Role to the role's role name", func() {
				Expect(result.Role).To(Equal("storage"))
			})

			It("should set the Coherence Image", func() {
				Expect(result.Coherence).To(Equal(role.Spec.Images.Coherence))
			})

			It("should set the Coherence Utils Image", func() {
				Expect(result.CoherenceUtils).To(Equal(role.Spec.Images.CoherenceUtils))
			})

			It("should set the User Artifacts Image", func() {
				Expect(result.UserArtifacts).To(Equal(role.Spec.Images.UserArtifacts))
			})

			It("should set the Fluentd Image", func() {
				Expect(result.Fluentd).To(Equal(role.Spec.Images.Fluentd))
			})

			It("should set the Store WKA", func() {
				Expect(result.Store.WKA).To(Equal("test-cluster-wka"))
			})

			It("should set the Store StorageEnabled", func() {
				Expect(result.Store.StorageEnabled).To(Equal(role.Spec.StorageEnabled))
			})

			It("should set the Store ReadinessProbe", func() {
				Expect(result.Store.ReadinessProbe).To(Equal(role.Spec.ReadinessProbe))
			})

			It("should set the Store Labels", func() {
				expected := make(map[string]string)

				for k, v := range *role.Spec.Labels {
					expected[k] = v
				}
				expected[coherence.CoherenceClusterLabel] = "test-cluster"
				expected[coherence.CoherenceRoleLabel] = "storage"

				Expect(result.Store.Labels).To(Equal(expected))
			})

		})
	})

	Context("Creating a CoherenceInternal as a Map from a CoherenceCluster and CoherenceRole", func() {
		var (
			cohMap map[string]interface{}
			cohInt coherence.CoherenceInternalSpec
			err    error
		)

		JustBeforeEach(func() {
			cohInt = coherence.NewCoherenceInternalSpec(&cluster, &role)
			cohMap, err = coherence.NewCoherenceInternalSpecAsMap(&cluster, &role)
		})

		It("should not produce an error", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("should not produce a map that serializes back to the expected CoherenceInternal", func() {
			data, e := json.Marshal(cohMap)
			Expect(e).ToNot(HaveOccurred())

			result := coherence.CoherenceInternalSpec{}

			e = json.Unmarshal(data, &result)
			Expect(e).ToNot(HaveOccurred())

			Expect(result).To(Equal(cohInt))
		})
	})

	When("Getting the GroupVersionKind", func() {
		var s *runtime.Scheme
		var gvk schema.GroupVersionKind

		BeforeEach(func() {
			s = scheme.Scheme

			_ = apis.AddToScheme(s)

			gvk = coherence.GetCoherenceInternalGroupVersionKind(s)
		})

		It("should have the correct Group", func() {
			Expect(gvk.Group).To(Equal("coherence.oracle.com"))
		})

		It("should have the correct Version", func() {
			Expect(gvk.Version).To(Equal("v1"))
		})

		It("should have the correct Kind", func() {
			Expect(gvk.Kind).To(Equal(reflect.TypeOf(coherence.CoherenceInternal{}).Name()))
		})
	})
})
