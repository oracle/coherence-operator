/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1_test

import (
	"encoding/json"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
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
		clusterIP    = corev1.ServiceTypeClusterIP
		loadBalancer = corev1.ServiceTypeLoadBalancer
		cluster      coherence.CoherenceCluster
		role         coherence.CoherenceRole
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
		block := corev1.PersistentVolumeBlock

		portOne := coherence.NamedPortSpec{
			Name: "one",
			PortSpec: coherence.PortSpec{
				Port: 100,
				Service: &coherence.ServiceSpec{
					Enabled:        boolPtr(true),
					Port:           int32Ptr(1100),
					Type:           &clusterIP,
					LoadBalancerIP: stringPtr("10.10.100.1"),
				}}}

		portTwo := coherence.NamedPortSpec{
			Name: "two",
			PortSpec: coherence.PortSpec{
				Port: 200,
				Service: &coherence.ServiceSpec{
					Enabled:        boolPtr(true),
					Port:           int32Ptr(1200),
					Type:           &loadBalancer,
					LoadBalancerIP: stringPtr("10.10.100.2"),
				}}}

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
				Labels:         map[string]string{"one": "1", "two": "2"},
				CacheConfig:    stringPtr("cache-config.xml"),
				PofConfig:      stringPtr("pof-config.xml"),
				OverrideConfig: stringPtr("coherence-override.xml"),
				Logging: &coherence.LoggingSpec{
					Level:         int32Ptr(9),
					ConfigFile:    stringPtr("logging.properties"),
					ConfigMapName: stringPtr("loggingMap"),
					Fluentd: &coherence.FluentdSpec{
						ImageSpec: coherence.ImageSpec{
							Image:           stringPtr("fluentd:1.0"),
							ImagePullPolicy: &always,
						},
						Enabled:    boolPtr(true),
						ConfigFile: stringPtr("one.yaml"),
						Tag:        stringPtr("tag-one"),
					},
				},
				Main: &coherence.MainSpec{
					Class:     stringPtr("com.tangosol.net.DefaultCacheServer"),
					Arguments: stringPtr("-Dcoherence.localhost=192.168.0.301"),
				},
				MaxHeap:     stringPtr("-Xmx1G"),
				JvmArgs:     stringPtr("-XX:+UseG1GC"),
				JavaOpts:    stringPtr("-Dcoherence.log.level=9"),
				Ports:       []coherence.NamedPortSpec{portOne, portTwo},
				Env:         []corev1.EnvVar{{Name: "FOO", Value: "foo-value"}, {Name: "BAR", Value: "bar-value"}},
				Annotations: map[string]string{"prometheus.io/scrape": "true", "prometheus.io/port": "2408"},
				Persistence: &coherence.PersistentStorageSpec{
					Enabled: boolPtr(true),
					PersistentVolumeClaim: &corev1.PersistentVolumeClaimSpec{
						Resources: corev1.ResourceRequirements{
							Requests: map[corev1.ResourceName]resource.Quantity{"storage": resource.MustParse("2Gi")},
						},
						StorageClassName: stringPtr("sc1"),
						DataSource:       &corev1.TypedLocalObjectReference{Name: "pvc1", Kind: "PersistentVolumeClaim"},
						VolumeMode:       &block,
						VolumeName:       "name1",
						Selector:         &metav1.LabelSelector{MatchLabels: map[string]string{"component": "coh1"}},
					},
					Volume: &corev1.Volume{
						VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}},
					},
				},
				Snapshot: &coherence.PersistentStorageSpec{
					Enabled: boolPtr(true),
					PersistentVolumeClaim: &corev1.PersistentVolumeClaimSpec{
						Resources: corev1.ResourceRequirements{
							Requests: map[corev1.ResourceName]resource.Quantity{"storage": resource.MustParse("2Gi")},
						},
						StorageClassName: stringPtr("sc1"),
						DataSource:       &corev1.TypedLocalObjectReference{Name: "pvc2", Kind: "PersistentVolumeClaim"},
						VolumeMode:       &block,
						VolumeName:       "name",
						Selector:         &metav1.LabelSelector{MatchLabels: map[string]string{"component": "coh1"}},
					},
					Volume: &corev1.Volume{
						VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}},
					},
				},
				Management: &coherence.PortSpecWithSSL{
					Enabled: boolPtr(true),
					SSL: &coherence.SSLSpec{
						Enabled:                boolPtr(true),
						Secrets:                stringPtr("ssl-secret"),
						KeyStore:               stringPtr("keystore.jks"),
						KeyStorePasswordFile:   stringPtr("storepassword.txt"),
						KeyPasswordFile:        stringPtr("keypassword.txt"),
						KeyStoreAlgorithm:      stringPtr("SunX509"),
						KeyStoreProvider:       stringPtr("fooJCA"),
						KeyStoreType:           stringPtr("JKS"),
						TrustStore:             stringPtr("truststore-guardians.jks"),
						TrustStorePasswordFile: stringPtr("trustpassword.txt"),
						TrustStoreAlgorithm:    stringPtr("SunX509"),
						TrustStoreProvider:     stringPtr("fooJCA"),
						TrustStoreType:         stringPtr("JKS"),
						RequireClientCert:      boolPtr(true),
					},
				},
				Metrics: &coherence.PortSpecWithSSL{
					Enabled: boolPtr(false),
					SSL: &coherence.SSLSpec{
						Enabled:                boolPtr(true),
						Secrets:                stringPtr("ssl-secret"),
						KeyStore:               stringPtr("keystore.jks"),
						KeyStorePasswordFile:   stringPtr("storepassword.txt"),
						KeyPasswordFile:        stringPtr("keypassword.txt"),
						KeyStoreAlgorithm:      stringPtr("SunX509"),
						KeyStoreProvider:       stringPtr("fooJCA"),
						KeyStoreType:           stringPtr("JKS"),
						TrustStore:             stringPtr("truststore-guardians.jks"),
						TrustStorePasswordFile: stringPtr("trustpassword.txt"),
						TrustStoreAlgorithm:    stringPtr("SunX509"),
						TrustStoreProvider:     stringPtr("fooJCA"),
						TrustStoreType:         stringPtr("JKS"),
						RequireClientCert:      boolPtr(true),
					},
				},
				JMX: &coherence.JMXSpec{
					Enabled:  boolPtr(true),
					Replicas: int32Ptr(3),
					MaxHeap:  stringPtr("2Gi"),
					Service: &coherence.ServiceSpec{
						Type:           &loadBalancer,
						LoadBalancerIP: stringPtr("10.10.10.20"),
						Annotations:    map[string]string{"foo": "1"},
						Port:           int32Ptr(9099),
					},
				},
				Volumes: []corev1.Volume{
					{
						Name:         "vol1",
						VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}},
					},
				},
				VolumeClaimTemplates: []corev1.PersistentVolumeClaim{
					{
						ObjectMeta: metav1.ObjectMeta{Name: "test-mount-1"},
						Spec: corev1.PersistentVolumeClaimSpec{
							AccessModes: []corev1.PersistentVolumeAccessMode{"ReadWriteOnce"},
							Resources: corev1.ResourceRequirements{
								Requests: map[corev1.ResourceName]resource.Quantity{"storage": resource.MustParse("2Gi")},
							},
						},
					},
				},
				VolumeMounts: []corev1.VolumeMount{
					{Name: "vol-mount-1", ReadOnly: false, MountPath: "/mountpath1"},
				},
				Affinity: &corev1.Affinity{
					NodeAffinity: &corev1.NodeAffinity{
						RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
							NodeSelectorTerms: []corev1.NodeSelectorTerm{
								{
									MatchExpressions: []corev1.NodeSelectorRequirement{
										{
											Key:      "kubernetes.io/e2e-az-name",
											Operator: corev1.NodeSelectorOpIn,
											Values:   []string{"e2e-az1", "e2e-az2"},
										},
									},
								},
							},
						},
					},
				},
				NodeSelector: map[string]string{"one": "1", "two": "2"},
				Tolerations: []corev1.Toleration{
					{Key: "key", Operator: "Equal", Value: "value", Effect: "NoSchedule"},
				},
				Resources: &corev1.ResourceRequirements{
					Requests: map[corev1.ResourceName]resource.Quantity{"storage": resource.MustParse("4Gi")},
				},
			},
		}
	})

	Context("Creating a CoherenceInternal from a CoherenceCluster and CoherenceRole", func() {
		var (
			result *coherence.CoherenceInternalSpec
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

				for k, v := range role.Spec.Labels {
					expected[k] = v
				}
				expected[coherence.CoherenceClusterLabel] = "test-cluster"
				expected[coherence.CoherenceRoleLabel] = "storage"

				Expect(result.Store.Labels).To(Equal(expected))
			})

			It("should set the Store CacheConfig", func() {
				Expect(result.Store.CacheConfig).To(Equal(role.Spec.CacheConfig))
			})

			It("should set the Store PofConfig", func() {
				Expect(result.Store.PofConfig).To(Equal(role.Spec.PofConfig))
			})

			It("should set the Store OverrideConfig", func() {
				Expect(result.Store.OverrideConfig).To(Equal(role.Spec.OverrideConfig))
			})

			It("should set the Store Logging", func() {
				Expect(result.Store.Logging).To(Equal(role.Spec.Logging))
			})

			It("should set the Store Main", func() {
				Expect(result.Store.Main).To(Equal(role.Spec.Main))
			})

			It("should set the Store MaxHeap", func() {
				Expect(result.Store.MaxHeap).To(Equal(role.Spec.MaxHeap))
			})

			It("should set the Store JvmArgs", func() {
				Expect(result.Store.JvmArgs).To(Equal(role.Spec.JvmArgs))
			})

			It("should set the Store JavaOpts", func() {
				Expect(result.Store.JavaOpts).To(Equal(role.Spec.JavaOpts))
			})

			It("should set the Store Ports", func() {
				Expect(result.Store.Ports).To(Equal(role.Spec.Ports))
			})

			It("should set the Store Env", func() {
				Expect(result.Store.Env).To(Equal(role.Spec.Env))
			})

			It("should set the Store Annotations", func() {
				expected := make(map[string]string)

				for k, v := range role.Spec.Annotations {
					expected[k] = v
				}

				Expect(result.Store.Annotations).To(Equal(expected))
			})

			It("should set the Store PodManagementPolicy", func() {
				Expect(result.Store.PodManagementPolicy).To(Equal(role.Spec.PodManagementPolicy))
			})

			It("should set the Store RevisionHistoryLimit", func() {
				Expect(result.Store.RevisionHistoryLimit).To(Equal(role.Spec.RevisionHistoryLimit))
			})

			It("should set the Store Persistence", func() {
				expectedPersistence := role.Spec.Persistence.DeepCopy()
				if expectedPersistence.Volume != nil {
					expectedPersistence.Volume.Name = "persistence-volume"
				}
				Expect(result.Store.Persistence).To(Equal(expectedPersistence))
			})

			It("should set the Store Snapshot", func() {
				expectedSnapshot := role.Spec.Snapshot.DeepCopy()
				if expectedSnapshot.Volume != nil {
					expectedSnapshot.Volume.Name = "snapshot-volume"
				}
				Expect(result.Store.Snapshot).To(Equal(expectedSnapshot))
			})

			It("should set the Store Management", func() {
				Expect(result.Store.Management).To(Equal(role.Spec.Management))
			})

			It("should set the Store Metrics", func() {
				Expect(result.Store.Metrics).To(Equal(role.Spec.Metrics))
			})

			It("should set the Store JMX", func() {
				Expect(result.Store.JMX).To(Equal(role.Spec.JMX))
			})

			It("should set the Store Volumes", func() {
				Expect(result.Store.Volumes).To(Equal(role.Spec.Volumes))
			})

			It("should set the Store VolumeClaimTemplates", func() {
				Expect(result.Store.VolumeClaimTemplates).To(Equal(role.Spec.VolumeClaimTemplates))
			})

			It("should set the Store VolumeMounts", func() {
				Expect(result.Store.VolumeMounts).To(Equal(role.Spec.VolumeMounts))
			})

			It("should set the Affinity", func() {
				Expect(result.Affinity).To(Equal(role.Spec.Affinity))
			})

			It("should set the NodeSelector", func() {
				Expect(result.NodeSelector).To(Equal(role.Spec.NodeSelector))
			})

			It("should set the Tolerations", func() {
				Expect(result.Tolerations).To(Equal(role.Spec.Tolerations))
			})

			It("should set the Resources", func() {
				Expect(result.Resources).To(Equal(role.Spec.Resources))
			})
		})
	})

	Context("Creating a CoherenceInternal as a Map from a CoherenceCluster and CoherenceRole", func() {
		var (
			cohMap map[string]interface{}
			cohInt *coherence.CoherenceInternalSpec
			err    error
		)

		JustBeforeEach(func() {
			cohInt = coherence.NewCoherenceInternalSpec(&cluster, &role)
			cohMap, err = coherence.NewCoherenceInternalSpecAsMap(&cluster, &role)
		})

		It("should not produce an error", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("should produce a map that serializes back to the expected CoherenceInternal", func() {
			data, e := json.Marshal(cohMap)
			Expect(e).ToNot(HaveOccurred())

			result := &coherence.CoherenceInternalSpec{}

			e = json.Unmarshal(data, result)
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
