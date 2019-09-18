/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	coherence "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
)

var _ = Describe("Testing CoherenceRoleSpec struct", func() {

	Context("Copying an CoherenceRoleSpec using DeepCopyWithDefaults", func() {
		var (
			clusterIP    = corev1.ServiceTypeClusterIP
			loadBalancer = corev1.ServiceTypeLoadBalancer

			roleNameOne = "one"
			roleNameTwo = "two"

			storageOne = false
			storageTwo = true

			replicasOne = int32Ptr(19)
			replicasTwo = int32Ptr(20)

			probeOne = &coherence.ReadinessProbeSpec{InitialDelaySeconds: int32Ptr(10)}
			probeTwo = &coherence.ReadinessProbeSpec{InitialDelaySeconds: int32Ptr(99)}

			imagesOne = &coherence.Images{Coherence: &coherence.ImageSpec{Image: stringPtr("foo:1.0")}}
			imagesTwo = &coherence.Images{Coherence: &coherence.ImageSpec{Image: stringPtr("foo:2.0")}}

			logginOne = &coherence.LoggingSpec{Level: int32Ptr(9)}
			logginTwo = &coherence.LoggingSpec{Level: int32Ptr(7)}

			mainOne = &coherence.MainSpec{Class: stringPtr("com.tangosol.net.DefaultCacheServer")}
			mainTwo = &coherence.MainSpec{Class: stringPtr("com.tangosol.net.DefaultCacheServer2")}

			scalingOne = coherence.SafeScaling
			scalingTwo = coherence.ParallelScaling

			labelsOne = map[string]string{"one": "1", "two": "2"}
			labelsTwo = map[string]string{"three": "3", "four": "4"}

			cacheConfigOne = "config-one.xml"
			cacheConfigTwo = "config-two.xml"

			pofConfigOne = "pof-one.xml"
			pofConfigTwo = "pof-two.xml"

			overrideConfigOne = "tangsolol-coherence-override-one.xml"
			overrideConfigTwo = "tangsolol-coherence-override-two.xml"

			maxHeapOne = "-Xmx1G"
			maxHeapTwo = "-Xmx2G"

			jvmArgsOne = "-XX:+UseG1GC"
			jvmArgsTwo = ""

			javaOptsOne = "-Dcoherence.log.level=9"
			javaOptsTwo = "-Dcoherence.log.level=5"

			envOne = []corev1.EnvVar{{Name: "foo1", Value: "1"}, {Name: "foo2", Value: "2"}}
			envTwo = []corev1.EnvVar{{Name: "foo3", Value: "3"}, {Name: "foo4", Value: "4"}}

			annotationsOne = map[string]string{"anno1": "1", "anno2": "2"}
			annotationsTwo = map[string]string{"anno3": "3", "anno4": "4"}

			podManagementPolicyOne appv1.PodManagementPolicyType = "Parallel"
			podManagementPolicyTwo appv1.PodManagementPolicyType = "OrderedReady"

			revisionHistoryLimitOne = int32Ptr(3)
			revisionHistoryLimitTwo = int32Ptr(5)

			block          = corev1.PersistentVolumeBlock
			filesystem     = corev1.PersistentVolumeFilesystem
			persistenceOne = &coherence.PersistentStorageSpec{
				Enabled: boolPtr(true),
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimSpec{
					AccessModes: []corev1.PersistentVolumeAccessMode{"ReadWriteOnce"},
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
					Name:         "vol1",
					VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}},
				},
			}
			persistenceTwo = &coherence.PersistentStorageSpec{
				Enabled: boolPtr(true),
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimSpec{
					AccessModes: []corev1.PersistentVolumeAccessMode{"ReadWriteOnce"},
					Resources: corev1.ResourceRequirements{
						Requests: map[corev1.ResourceName]resource.Quantity{"storage": resource.MustParse("4Gi")},
					},
					StorageClassName: stringPtr("sc2"),
					DataSource:       &corev1.TypedLocalObjectReference{Name: "pvc2", Kind: "PersistentVolumeClaim"},
					VolumeMode:       &filesystem,
					VolumeName:       "name2",
					Selector:         &metav1.LabelSelector{MatchLabels: map[string]string{"component": "coh2"}},
				},
				Volume: &corev1.Volume{
					Name:         "vol2",
					VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}},
				},
			}

			snapshotOne = &coherence.PersistentStorageSpec{
				Enabled: boolPtr(true),
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimSpec{
					AccessModes: []corev1.PersistentVolumeAccessMode{"ReadWriteOnce"},
					Resources: corev1.ResourceRequirements{
						Requests: map[corev1.ResourceName]resource.Quantity{"storage": resource.MustParse("2Gi")},
					},
					StorageClassName: stringPtr("sc1"),
					DataSource:       &corev1.TypedLocalObjectReference{Name: "pvc3", Kind: "PersistentVolumeClaim"},
					VolumeMode:       &block,
					VolumeName:       "name3",
					Selector:         &metav1.LabelSelector{MatchLabels: map[string]string{"component": "coh1"}},
				},
				Volume: &corev1.Volume{
					Name:         "vol3",
					VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}},
				},
			}
			snapshotTwo = &coherence.PersistentStorageSpec{
				Enabled: boolPtr(true),
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimSpec{
					AccessModes: []corev1.PersistentVolumeAccessMode{"ReadWriteOnce"},
					Resources: corev1.ResourceRequirements{
						Requests: map[corev1.ResourceName]resource.Quantity{"storage": resource.MustParse("4Gi")},
					},
					StorageClassName: stringPtr("sc2"),
					DataSource:       &corev1.TypedLocalObjectReference{Name: "pvc4", Kind: "PersistentVolumeClaim"},
					VolumeMode:       &filesystem,
					VolumeName:       "name4",
					Selector:         &metav1.LabelSelector{MatchLabels: map[string]string{"component": "coh2"}},
				},
				Volume: &corev1.Volume{
					Name:         "vol4",
					VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}},
				},
			}

			managementOne = &coherence.PortSpecWithSSL{
				Enabled: boolPtr(true),
				SSL: &coherence.SSLSpec{
					Enabled:                boolPtr(true),
					KeyStore:               stringPtr("keystore.jks"),
					KeyStorePasswordFile:   stringPtr("keypassword.txt"),
					KeyStoreType:           stringPtr("JKS"),
					TrustStore:             stringPtr("trustore-guaradians.jks"),
					TrustStorePasswordFile: stringPtr("trustpassword.txt"),
					TrustStoreType:         stringPtr("JKS"),
					RequireClientCert:      boolPtr(true),
				},
			}

			managementTwo = &coherence.PortSpecWithSSL{
				Enabled: boolPtr(false),
				SSL: &coherence.SSLSpec{
					Enabled: boolPtr(false),
				},
			}

			metricsOne = &coherence.PortSpecWithSSL{
				Enabled: boolPtr(true),
				SSL: &coherence.SSLSpec{
					Enabled:                boolPtr(true),
					KeyStore:               stringPtr("keystore.jks"),
					KeyStorePasswordFile:   stringPtr("keypassword.txt"),
					KeyStoreType:           stringPtr("JKS"),
					TrustStore:             stringPtr("trustore-guaradians.jks"),
					TrustStorePasswordFile: stringPtr("trustpassword.txt"),
					TrustStoreType:         stringPtr("JKS"),
					RequireClientCert:      boolPtr(true),
				},
			}
			metricsTwo = &coherence.PortSpecWithSSL{
				Enabled: boolPtr(false),
			}

			jmxOne = &coherence.JMXSpec{
				Enabled:  boolPtr(true),
				Replicas: int32Ptr(3),
				MaxHeap:  stringPtr("2Gi"),
				Service: &coherence.ServiceSpec{
					Type:           &loadBalancer,
					LoadBalancerIP: stringPtr("10.10.10.20"),
					Annotations:    map[string]string{"foo": "1"},
					Port:           int32Ptr(9099),
				},
			}
			jmxTwo = &coherence.JMXSpec{
				Enabled:  boolPtr(true),
				Replicas: int32Ptr(6),
				MaxHeap:  stringPtr("3Gi"),
				Service: &coherence.ServiceSpec{
					Type:           &clusterIP,
					LoadBalancerIP: stringPtr("10.10.10.21"),
					Annotations:    map[string]string{"foo": "2"},
					Port:           int32Ptr(9098),
				},
			}

			clientAffinity = corev1.ServiceAffinityClientIP
			tpLocal        = corev1.ServiceExternalTrafficPolicyTypeLocal

			portOne = coherence.NamedPortSpec{
				Name: "one",
				PortSpec: coherence.PortSpec{
					Port: 1000,
					Service: &coherence.ServiceSpec{
						Enabled:                  boolPtr(true),
						Name:                     stringPtr("svc-one"),
						Port:                     int32Ptr(1100),
						Type:                     &clusterIP,
						LoadBalancerIP:           stringPtr("10.10.100.10"),
						Annotations:              map[string]string{"ann-key": "one"},
						SessionAffinity:          &clientAffinity,
						LoadBalancerSourceRanges: []string{"key", "one"},
						ExternalName:             stringPtr("ext-one"),
						ExternalTrafficPolicy:    &tpLocal,
						HealthCheckNodePort:      int32Ptr(2000),
						PublishNotReadyAddresses: boolPtr(true),
						SessionAffinityConfig: &corev1.SessionAffinityConfig{
							ClientIP: &corev1.ClientIPConfig{TimeoutSeconds: int32Ptr(60)},
						},
					},
				},
			}

			portTwo = coherence.NamedPortSpec{
				Name: "two",
				PortSpec: coherence.PortSpec{
					Port: 1200,
					Service: &coherence.ServiceSpec{
						Enabled:                  boolPtr(true),
						Name:                     stringPtr("svc-two"),
						Port:                     int32Ptr(1120),
						Type:                     &clusterIP,
						LoadBalancerIP:           stringPtr("10.10.100.20"),
						Annotations:              map[string]string{"ann-key": "two"},
						SessionAffinity:          &clientAffinity,
						LoadBalancerSourceRanges: []string{"key", "two"},
						ExternalName:             stringPtr("ext-two"),
						ExternalTrafficPolicy:    &tpLocal,
						HealthCheckNodePort:      int32Ptr(2200),
						PublishNotReadyAddresses: boolPtr(true),
						SessionAffinityConfig: &corev1.SessionAffinityConfig{
							ClientIP: &corev1.ClientIPConfig{TimeoutSeconds: int32Ptr(120)},
						},
					},
				},
			}

			portThree = coherence.NamedPortSpec{
				Name: "three",
				PortSpec: coherence.PortSpec{
					Port: 1300,
					Service: &coherence.ServiceSpec{
						Enabled:                  boolPtr(true),
						Name:                     stringPtr("svc-three"),
						Port:                     int32Ptr(1130),
						Type:                     &clusterIP,
						LoadBalancerIP:           stringPtr("10.10.100.30"),
						Annotations:              map[string]string{"ann-key": "three"},
						SessionAffinity:          &clientAffinity,
						LoadBalancerSourceRanges: []string{"key", "three"},
						ExternalName:             stringPtr("ext-three"),
						ExternalTrafficPolicy:    &tpLocal,
						HealthCheckNodePort:      int32Ptr(2300),
						PublishNotReadyAddresses: boolPtr(true),
						SessionAffinityConfig: &corev1.SessionAffinityConfig{
							ClientIP: &corev1.ClientIPConfig{TimeoutSeconds: int32Ptr(180)},
						},
					},
				},
			}

			volumesOne = []corev1.Volume{
				{
					Name:         "vol1",
					VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}},
				},
			}
			volumesTwo = []corev1.Volume{
				{
					Name:         "vol2",
					VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}},
				},
			}

			volumeClaimTemplatesOne = []corev1.PersistentVolumeClaim{
				{
					ObjectMeta: metav1.ObjectMeta{Name: "test-mount-1"},
					Spec: corev1.PersistentVolumeClaimSpec{
						AccessModes: []corev1.PersistentVolumeAccessMode{"ReadWriteOnce"},
						Resources: corev1.ResourceRequirements{
							Requests: map[corev1.ResourceName]resource.Quantity{"storage": resource.MustParse("2Gi")},
						},
					},
				},
			}
			volumeClaimTemplatesTwo = []corev1.PersistentVolumeClaim{
				{
					ObjectMeta: metav1.ObjectMeta{Name: "test-mount-2"},
					Spec: corev1.PersistentVolumeClaimSpec{
						AccessModes: []corev1.PersistentVolumeAccessMode{"ReadWriteMany"},
						Resources: corev1.ResourceRequirements{
							Requests: map[corev1.ResourceName]resource.Quantity{"storage": resource.MustParse("4Gi")},
						},
					},
				},
			}

			volumeMountsOne = []corev1.VolumeMount{
				{Name: "vol-mount-1", ReadOnly: false, MountPath: "/mountpath1"},
			}
			volumeMountsTwo = []corev1.VolumeMount{
				{Name: "vol-mount-2", MountPath: "/mountpath2"},
			}

			affinityOne = &corev1.Affinity{
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
			}
			affinityTwo = &corev1.Affinity{
				NodeAffinity: &corev1.NodeAffinity{
					RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
						NodeSelectorTerms: []corev1.NodeSelectorTerm{
							{
								MatchExpressions: []corev1.NodeSelectorRequirement{
									{
										Key:      "kubernetes.io/e2e-az-name",
										Operator: corev1.NodeSelectorOpIn,
										Values:   []string{"e2e-az3", "e2e-az4"},
									},
								},
							},
						},
					},
				},
			}

			nodeSelectorOne = map[string]string{"one": "1", "two": "2"}
			nodeSelectorTwo = map[string]string{"three": "3", "four": "4"}

			tolerationsOne = []corev1.Toleration{
				{Key: "key", Operator: "Equal", Value: "value", Effect: "NoSchedule"},
			}
			tolerationsTwo = []corev1.Toleration{
				{Operator: "Exists"},
			}

			resourcesOne = &corev1.ResourceRequirements{
				Requests: map[corev1.ResourceName]resource.Quantity{"storage": resource.MustParse("4Gi")},
			}
			resourcesTwo = &corev1.ResourceRequirements{
				Requests: map[corev1.ResourceName]resource.Quantity{"storage": resource.MustParse("8Gi")},
			}

			statusHAOne = &coherence.StatusHAHandler{Handler: corev1.Handler{Exec: &corev1.ExecAction{Command: []string{"one"}}}}
			statusHATwo = &coherence.StatusHAHandler{Handler: corev1.Handler{Exec: &corev1.ExecAction{Command: []string{"two"}}}}

			roleSpecOne = &coherence.CoherenceRoleSpec{
				Role:                 roleNameOne,
				Replicas:             replicasOne,
				Images:               imagesOne,
				StorageEnabled:       &storageOne,
				ScalingPolicy:        &scalingOne,
				ReadinessProbe:       probeOne,
				Labels:               nil,
				CacheConfig:          &cacheConfigOne,
				PofConfig:            &pofConfigOne,
				OverrideConfig:       &overrideConfigOne,
				Logging:              logginOne,
				Main:                 mainOne,
				MaxHeap:              &maxHeapOne,
				JvmArgs:              &jvmArgsOne,
				JavaOpts:             &javaOptsOne,
				Ports:                nil,
				Env:                  nil,
				Annotations:          nil,
				PodManagementPolicy:  &podManagementPolicyOne,
				RevisionHistoryLimit: revisionHistoryLimitOne,
				Persistence:          persistenceOne,
				Snapshot:             snapshotOne,
				Management:           managementOne,
				Metrics:              metricsOne,
				JMX:                  jmxOne,
				Volumes:              volumesOne,
				VolumeClaimTemplates: volumeClaimTemplatesOne,
				VolumeMounts:         volumeMountsOne,
				Affinity:             affinityOne,
				NodeSelector:         nil,
				Tolerations:          tolerationsOne,
				Resources:            resourcesOne,
				StatusHA:             statusHAOne,
			}

			roleSpecTwo = &coherence.CoherenceRoleSpec{
				Role:                 roleNameTwo,
				Replicas:             replicasTwo,
				Images:               imagesTwo,
				StorageEnabled:       &storageTwo,
				ScalingPolicy:        &scalingTwo,
				ReadinessProbe:       probeTwo,
				Labels:               nil,
				CacheConfig:          &cacheConfigTwo,
				PofConfig:            &pofConfigTwo,
				OverrideConfig:       &overrideConfigTwo,
				Logging:              logginTwo,
				Main:                 mainTwo,
				MaxHeap:              &maxHeapTwo,
				JvmArgs:              &jvmArgsTwo,
				JavaOpts:             &javaOptsTwo,
				Ports:                nil,
				Env:                  nil,
				Annotations:          nil,
				PodManagementPolicy:  &podManagementPolicyTwo,
				RevisionHistoryLimit: revisionHistoryLimitTwo,
				Persistence:          persistenceTwo,
				Snapshot:             snapshotTwo,
				Management:           managementTwo,
				Metrics:              metricsTwo,
				JMX:                  jmxTwo,
				Volumes:              volumesTwo,
				VolumeClaimTemplates: volumeClaimTemplatesTwo,
				VolumeMounts:         volumeMountsTwo,
				Affinity:             affinityTwo,
				NodeSelector:         nil,
				Tolerations:          tolerationsTwo,
				Resources:            resourcesTwo,
				StatusHA:             statusHATwo,
			}

			original *coherence.CoherenceRoleSpec
			defaults *coherence.CoherenceRoleSpec
			clone    *coherence.CoherenceRoleSpec
		)

		// just before every "It" this method is executed to actually do the cloning
		JustBeforeEach(func() {
			clone = original.DeepCopyWithDefaults(defaults)
		})

		When("original and defaults are nil", func() {
			BeforeEach(func() {
				original = nil
				defaults = nil
			})

			It("the copy should be nil", func() {
				Expect(clone).Should(BeNil())
			})
		})

		When("defaults is nil", func() {
			BeforeEach(func() {
				original = roleSpecOne
				defaults = nil
			})

			It("clone should be equal to the original", func() {
				Expect(clone).To(Equal(original))
			})
		})

		When("original is nil", func() {
			BeforeEach(func() {
				original = nil
				defaults = roleSpecTwo
			})

			It("clone should be equal to the defaults", func() {
				Expect(clone).To(Equal(defaults))
			})
		})

		When("all fields in the original are set", func() {
			BeforeEach(func() {
				original = roleSpecOne
				defaults = roleSpecTwo
			})

			It("clone should be equal to the original", func() {
				Expect(clone).To(Equal(original))
			})
		})

		When("the original Role is not set", func() {
			BeforeEach(func() {
				defaults = roleSpecTwo
				// original is a deep copy of roleSpecOne so that we can change the
				// original without changing roleSpecOne
				original = roleSpecOne.DeepCopy()
				original.Role = ""
			})

			It("clone should be equal to the original with the Role field from the defaults", func() {
				// expected is a deep copy of original so that we can change the
				// expected without changing original
				expected := original.DeepCopy()
				expected.Role = defaults.Role

				Expect(clone).To(Equal(expected))
			})
		})

		When("the original Replicas is not set", func() {
			BeforeEach(func() {
				defaults = roleSpecTwo
				// original is a deep copy of roleSpecOne so that we can change the
				// original without changing roleSpecOne
				original = roleSpecOne.DeepCopy()
				original.Replicas = nil
			})

			It("clone should be equal to the original with the Replicas field from the defaults", func() {
				// expected is a deep copy of original so that we can change the
				// expected without changing original
				expected := original.DeepCopy()
				expected.Replicas = defaults.Replicas

				Expect(clone).To(Equal(expected))
			})
		})

		When("the original CacheConfig is not set", func() {
			BeforeEach(func() {
				defaults = roleSpecTwo
				// original is a deep copy of roleSpecOne so that we can change the
				// original without changing roleSpecOne
				original = roleSpecOne.DeepCopy()
				original.CacheConfig = nil
			})

			It("clone should be equal to the original with the CacheConfig field from the defaults", func() {
				// expected is a deep copy of original so that we can change the
				// expected without changing original
				expected := original.DeepCopy()
				expected.CacheConfig = defaults.CacheConfig

				Expect(clone).To(Equal(expected))
			})
		})

		When("the original PofConfig is not set", func() {
			BeforeEach(func() {
				defaults = roleSpecTwo
				// original is a deep copy of roleSpecOne so that we can change the
				// original without changing roleSpecOne
				original = roleSpecOne.DeepCopy()
				original.PofConfig = nil
			})

			It("clone should be equal to the original with the PofConfig field from the defaults", func() {
				// expected is a deep copy of original so that we can change the
				// expected without changing original
				expected := original.DeepCopy()
				expected.PofConfig = defaults.PofConfig

				Expect(clone).To(Equal(expected))
			})
		})

		When("the original OverrideConfig is not set", func() {
			BeforeEach(func() {
				defaults = roleSpecTwo
				// original is a deep copy of roleSpecOne so that we can change the
				// original without changing roleSpecOne
				original = roleSpecOne.DeepCopy()
				original.OverrideConfig = nil
			})

			It("clone should be equal to the original with the OverrideConfig field from the defaults", func() {
				// expected is a deep copy of original so that we can change the
				// expected without changing original
				expected := original.DeepCopy()
				expected.OverrideConfig = defaults.OverrideConfig

				Expect(clone).To(Equal(expected))
			})
		})

		When("the original Logging is not set", func() {
			BeforeEach(func() {
				defaults = roleSpecTwo
				// original is a deep copy of roleSpecOne so that we can change the
				// original without changing roleSpecOne
				original = roleSpecOne.DeepCopy()
				original.Logging = nil
			})

			It("clone should be equal to the original with the Logging field from the defaults", func() {
				// expected is a deep copy of original so that we can change the
				// expected without changing original
				expected := original.DeepCopy()
				expected.Logging = defaults.Logging

				Expect(clone).To(Equal(expected))
			})
		})

		When("the original Main is not set", func() {
			BeforeEach(func() {
				defaults = roleSpecTwo
				// original is a deep copy of roleSpecOne so that we can change the
				// original without changing roleSpecOne
				original = roleSpecOne.DeepCopy()
				original.Main = nil
			})

			It("clone should be equal to the original with the Main field from the defaults", func() {
				// expected is a deep copy of original so that we can change the
				// expected without changing original
				expected := original.DeepCopy()
				expected.Main = defaults.Main

				Expect(clone).To(Equal(expected))
			})
		})

		When("the original MaxHeap is not set", func() {
			BeforeEach(func() {
				defaults = roleSpecTwo
				// original is a deep copy of roleSpecOne so that we can change the
				// original without changing roleSpecOne
				original = roleSpecOne.DeepCopy()
				original.MaxHeap = nil
			})

			It("clone should be equal to the original with the MaxHeap field from the defaults", func() {
				// expected is a deep copy of original so that we can change the
				// expected without changing original
				expected := original.DeepCopy()
				expected.MaxHeap = defaults.MaxHeap

				Expect(clone).To(Equal(expected))
			})
		})

		When("the original JvmArgs is not set", func() {
			BeforeEach(func() {
				defaults = roleSpecTwo
				// original is a deep copy of roleSpecOne so that we can change the
				// original without changing roleSpecOne
				original = roleSpecOne.DeepCopy()
				original.JvmArgs = nil
			})

			It("clone should be equal to the original with the JvmArgs field from the defaults", func() {
				// expected is a deep copy of original so that we can change the
				// expected without changing original
				expected := original.DeepCopy()
				expected.JvmArgs = defaults.JvmArgs

				Expect(clone).To(Equal(expected))
			})
		})

		When("the original JavaOpts is not set", func() {
			BeforeEach(func() {
				defaults = roleSpecTwo
				// original is a deep copy of roleSpecOne so that we can change the
				// original without changing roleSpecOne
				original = roleSpecOne.DeepCopy()
				original.JavaOpts = nil
			})

			It("clone should be equal to the original with the JavaOpts field from the defaults", func() {
				// expected is a deep copy of original so that we can change the
				// expected without changing original
				expected := original.DeepCopy()
				expected.JavaOpts = defaults.JavaOpts

				Expect(clone).To(Equal(expected))
			})
		})

		When("the original PodManagementPolicy is not set", func() {
			BeforeEach(func() {
				defaults = roleSpecTwo
				// original is a deep copy of roleSpecOne so that we can change the
				// original without changing roleSpecOne
				original = roleSpecOne.DeepCopy()
				original.PodManagementPolicy = nil
			})

			It("clone should be equal to the original with the PodManagementPolicy field from the defaults", func() {
				// expected is a deep copy of original so that we can change the
				// expected without changing original
				expected := original.DeepCopy()
				expected.PodManagementPolicy = defaults.PodManagementPolicy

				Expect(clone).To(Equal(expected))
			})
		})

		When("the original RevisionHistoryLimit is not set", func() {
			BeforeEach(func() {
				defaults = roleSpecTwo
				// original is a deep copy of roleSpecOne so that we can change the
				// original without changing roleSpecOne
				original = roleSpecOne.DeepCopy()
				original.RevisionHistoryLimit = nil
			})

			It("clone should be equal to the original with the RevisionHistoryLimit field from the defaults", func() {
				// expected is a deep copy of original so that we can change the
				// expected without changing original
				expected := original.DeepCopy()
				expected.RevisionHistoryLimit = defaults.RevisionHistoryLimit

				Expect(clone).To(Equal(expected))
			})
		})

		When("the original Persistence is not set", func() {
			BeforeEach(func() {
				defaults = roleSpecTwo
				// original is a deep copy of roleSpecOne so that we can change the
				// original without changing roleSpecOne
				original = roleSpecOne.DeepCopy()
				original.Persistence = nil
			})

			It("clone should be equal to the original with the Persistence field from the defaults", func() {
				// expected is a deep copy of original so that we can change the
				// expected without changing original
				expected := original.DeepCopy()
				expected.Persistence = defaults.Persistence

				Expect(clone).To(Equal(expected))
			})
		})

		When("the original Snapshot is not set", func() {
			BeforeEach(func() {
				defaults = roleSpecTwo
				// original is a deep copy of roleSpecOne so that we can change the
				// original without changing roleSpecOne
				original = roleSpecOne.DeepCopy()
				original.Snapshot = nil
			})

			It("clone should be equal to the original with the Snapshot field from the defaults", func() {
				// expected is a deep copy of original so that we can change the
				// expected without changing original
				expected := original.DeepCopy()
				expected.Snapshot = defaults.Snapshot

				Expect(clone).To(Equal(expected))
			})
		})

		When("the original Management is not set", func() {
			BeforeEach(func() {
				defaults = roleSpecTwo
				// original is a deep copy of roleSpecOne so that we can change the
				// original without changing roleSpecOne
				original = roleSpecOne.DeepCopy()
				original.Management = nil
			})

			It("clone should be equal to the original with the Management field from the defaults", func() {
				// expected is a deep copy of original so that we can change the
				// expected without changing original
				expected := original.DeepCopy()
				expected.Management = defaults.Management

				Expect(clone).To(Equal(expected))
			})
		})

		When("the original Metrics is not set", func() {
			BeforeEach(func() {
				defaults = roleSpecTwo
				// original is a deep copy of roleSpecOne so that we can change the
				// original without changing roleSpecOne
				original = roleSpecOne.DeepCopy()
				original.Metrics = nil
			})

			It("clone should be equal to the original with the Metrics field from the defaults", func() {
				// expected is a deep copy of original so that we can change the
				// expected without changing original
				expected := original.DeepCopy()
				expected.Metrics = defaults.Metrics

				Expect(clone).To(Equal(expected))
			})
		})

		When("the original JMX is not set", func() {
			BeforeEach(func() {
				defaults = roleSpecTwo
				// original is a deep copy of roleSpecOne so that we can change the
				// original without changing roleSpecOne
				original = roleSpecOne.DeepCopy()
				original.JMX = nil
			})

			It("clone should be equal to the original with the JMX field from the defaults", func() {
				// expected is a deep copy of original so that we can change the
				// expected without changing original
				expected := original.DeepCopy()
				expected.JMX = defaults.JMX

				Expect(clone).To(Equal(expected))
			})
		})

		When("the original Volumes is not set", func() {
			BeforeEach(func() {
				defaults = roleSpecTwo
				// original is a deep copy of roleSpecOne so that we can change the
				// original without changing roleSpecOne
				original = roleSpecOne.DeepCopy()
				original.Volumes = nil
			})

			It("clone should be equal to the original with the Volumes field from the defaults", func() {
				// expected is a deep copy of original so that we can change the
				// expected without changing original
				expected := original.DeepCopy()
				expected.Volumes = defaults.Volumes

				Expect(clone).To(Equal(expected))
			})
		})

		When("the original VolumeClaimTemplates is not set", func() {
			BeforeEach(func() {
				defaults = roleSpecTwo
				// original is a deep copy of roleSpecOne so that we can change the
				// original without changing roleSpecOne
				original = roleSpecOne.DeepCopy()
				original.VolumeClaimTemplates = nil
			})

			It("clone should be equal to the original with the VolumeClaimTemplates field from the defaults", func() {
				// expected is a deep copy of original so that we can change the
				// expected without changing original
				expected := original.DeepCopy()
				expected.VolumeClaimTemplates = defaults.VolumeClaimTemplates

				Expect(clone).To(Equal(expected))
			})
		})

		When("the original VolumeMounts is not set", func() {
			BeforeEach(func() {
				defaults = roleSpecTwo
				// original is a deep copy of roleSpecOne so that we can change the
				// original without changing roleSpecOne
				original = roleSpecOne.DeepCopy()
				original.VolumeMounts = nil
			})

			It("clone should be equal to the original with the VolumeMounts field from the defaults", func() {
				// expected is a deep copy of original so that we can change the
				// expected without changing original
				expected := original.DeepCopy()
				expected.VolumeMounts = defaults.VolumeMounts

				Expect(clone).To(Equal(expected))
			})
		})

		When("the original Affinity is not set", func() {
			BeforeEach(func() {
				defaults = roleSpecTwo
				// original is a deep copy of roleSpecOne so that we can change the
				// original without changing roleSpecOne
				original = roleSpecOne.DeepCopy()
				original.Affinity = nil
			})

			It("clone should be equal to the original with the Affinity field from the defaults", func() {
				// expected is a deep copy of original so that we can change the
				// expected without changing original
				expected := original.DeepCopy()
				expected.Affinity = defaults.Affinity

				Expect(clone).To(Equal(expected))
			})
		})

		When("the original Tolerations is not set", func() {
			BeforeEach(func() {
				defaults = roleSpecTwo
				// original is a deep copy of roleSpecOne so that we can change the
				// original without changing roleSpecOne
				original = roleSpecOne.DeepCopy()
				original.Tolerations = nil
			})

			It("clone should be equal to the original with the Tolerations field from the defaults", func() {
				// expected is a deep copy of original so that we can change the
				// expected without changing original
				expected := original.DeepCopy()
				expected.Tolerations = defaults.Tolerations

				Expect(clone).To(Equal(expected))
			})
		})

		When("the original Resources is not set", func() {
			BeforeEach(func() {
				defaults = roleSpecTwo
				// original is a deep copy of roleSpecOne so that we can change the
				// original without changing roleSpecOne
				original = roleSpecOne.DeepCopy()
				original.Resources = nil
			})

			It("clone should be equal to the original with the Resources field from the defaults", func() {
				// expected is a deep copy of original so that we can change the
				// expected without changing original
				expected := original.DeepCopy()
				expected.Resources = defaults.Resources

				Expect(clone).To(Equal(expected))
			})
		})

		When("the original ReadinessProbe is not set", func() {
			BeforeEach(func() {
				defaults = roleSpecTwo
				// original is a deep copy of roleSpecOne so that we can change the
				// original without changing roleSpecOne
				original = roleSpecOne.DeepCopy()
				original.ReadinessProbe = nil
			})

			It("clone should be equal to the original with the ReadinessProbe field from the defaults", func() {
				// expected is a deep copy of original so that we can change the
				// expected without changing original
				expected := original.DeepCopy()
				expected.ReadinessProbe = defaults.ReadinessProbe

				Expect(clone).To(Equal(expected))
			})
		})

		When("the original Images is not set", func() {
			BeforeEach(func() {
				defaults = roleSpecTwo
				// original is a deep copy of roleSpecOne so that we can change the
				// original without changing roleSpecOne
				original = roleSpecOne.DeepCopy()
				original.Images = nil
			})

			It("clone should be equal to the original with the Images field from the defaults", func() {
				// expected is a deep copy of original so that we can change the
				// expected without changing original
				expected := original.DeepCopy()
				expected.Images = defaults.Images

				Expect(clone).To(Equal(expected))
			})
		})

		When("the original ScalingPolicy is not set", func() {
			BeforeEach(func() {
				defaults = roleSpecTwo
				// original is a deep copy of roleSpecOne so that we can change the
				// original without changing roleSpecOne
				original = roleSpecOne.DeepCopy()
				original.ScalingPolicy = nil
			})

			It("clone should be equal to the original with the ScalingPolicy field from the defaults", func() {
				// expected is a deep copy of original so that we can change the
				// expected without changing original
				expected := original.DeepCopy()
				expected.ScalingPolicy = defaults.ScalingPolicy

				Expect(clone).To(Equal(expected))
			})
		})

		When("the original Labels is not set and default Labels is set", func() {
			BeforeEach(func() {
				// original and defaults are deep copies so that we can change them
				defaults = roleSpecTwo.DeepCopy()
				original = roleSpecOne.DeepCopy()
				original.Labels = nil
				defaults.Labels = labelsTwo
			})

			It("clone should be equal to the original with the Labels field from the defaults", func() {
				// expected is a deep copy of original so that we can change the
				// expected without changing original
				expected := original.DeepCopy()
				expected.Labels = defaults.Labels

				Expect(clone).To(Equal(expected))
			})
		})

		When("the original Labels is set and default Labels is not set", func() {
			BeforeEach(func() {
				// original and defaults are deep copies so that we can change them
				defaults = roleSpecTwo.DeepCopy()
				original = roleSpecOne.DeepCopy()
				original.Labels = labelsOne
				defaults.Labels = nil
			})

			It("clone should be equal to the original", func() {
				Expect(clone).To(Equal(original))
			})
		})

		When("the original Labels is set and default Labels is set", func() {
			BeforeEach(func() {
				// original and defaults are deep copies so that we can change them
				defaults = roleSpecTwo.DeepCopy()
				original = roleSpecOne.DeepCopy()
				original.Labels = labelsOne
				defaults.Labels = labelsTwo
			})

			It("clone should have the combined Labels", func() {
				// expected is a deep copy of original so that we can change the
				// expected without changing original
				expected := original.DeepCopy()
				expected.Labels = map[string]string{"one": "1", "two": "2", "three": "3", "four": "4"}

				Expect(clone).To(Equal(expected))
			})
		})

		When("the original Labels has duplicate keys to the default Labels", func() {
			BeforeEach(func() {
				// original and defaults are deep copies so that we can change them
				defaults = roleSpecTwo.DeepCopy()
				original = roleSpecOne.DeepCopy()

				original.Labels = map[string]string{"one": "1", "two": "2", "three": "changed"}
				defaults.Labels = map[string]string{"three": "3", "four": "4"}
			})

			It("clone should have the combined Labels with the duplicate key mapped to the originals value", func() {
				// expected is a deep copy of original so that we can change the
				// expected without changing original
				expected := original.DeepCopy()
				expected.Labels = map[string]string{"one": "1", "two": "2", "three": "changed", "four": "4"}

				Expect(clone).To(Equal(expected))
			})
		})

		When("the original Labels has duplicate keys to the default Labels where the value is empty string", func() {
			BeforeEach(func() {
				// original and defaults are deep copies so that we can change them
				defaults = roleSpecTwo.DeepCopy()
				original = roleSpecOne.DeepCopy()

				original.Labels = map[string]string{"one": "1", "two": "2", "three": ""}
				defaults.Labels = map[string]string{"three": "3", "four": "4"}
			})

			It("clone should have the combined Labels with the duplicate key removed", func() {
				// expected is a deep copy of original so that we can change the
				// expected without changing original
				expected := original.DeepCopy()
				expected.Labels = map[string]string{"one": "1", "two": "2", "four": "4"}

				Expect(clone).To(Equal(expected))
			})
		})

		When("the original Ports is not set and default Ports is set", func() {
			BeforeEach(func() {
				// original and defaults are deep copies so that we can change them
				defaults = roleSpecTwo.DeepCopy()
				original = roleSpecOne.DeepCopy()
				original.Ports = nil
				defaults.Ports = []coherence.NamedPortSpec{portOne, portTwo}
			})

			It("clone should be equal to the original with the Ports field from the defaults", func() {
				// expected is a deep copy of original so that we can change the
				// expected without changing original
				expected := original.DeepCopy()
				expected.Ports = []coherence.NamedPortSpec{portOne, portTwo}

				Expect(clone).To(Equal(expected))
			})
		})

		When("the original Ports is set and default Ports is not set", func() {
			BeforeEach(func() {
				// original and defaults are deep copies so that we can change them
				defaults = roleSpecTwo.DeepCopy()
				original = roleSpecOne.DeepCopy()
				original.Ports = []coherence.NamedPortSpec{portOne, portTwo}
				defaults.Ports = nil
			})

			It("clone should be equal to the original", func() {
				Expect(clone).To(Equal(original))
			})
		})

		When("the original Ports is set and default Ports is set", func() {
			BeforeEach(func() {
				// original and defaults are deep copies so that we can change them
				defaults = roleSpecTwo.DeepCopy()
				original = roleSpecOne.DeepCopy()
				original.Ports = []coherence.NamedPortSpec{portOne, portTwo}
				defaults.Ports = []coherence.NamedPortSpec{portThree}
			})

			It("clone should have the combined Ports", func() {
				// expected is a deep copy of original so that we can change the
				// expected without changing original
				expected := original.DeepCopy()
				expected.Ports = []coherence.NamedPortSpec{portOne, portTwo, portThree}

				Expect(clone).To(Equal(expected))
			})
		})

		When("the original StatusHA is not set and default StatusHA is set", func() {
			BeforeEach(func() {
				// original and defaults are deep copies so that we can change them
				defaults = roleSpecTwo.DeepCopy()
				original = roleSpecOne.DeepCopy()
				original.StatusHA = nil
			})

			It("clone should be equal to the original with the StatusHA field from the defaults", func() {
				// expected is a deep copy of original so that we can change the
				// expected without changing original
				expected := original.DeepCopy()
				expected.StatusHA = defaults.StatusHA
				Expect(clone).To(Equal(expected))
			})

			It("clone StatusHA should equal the defaults StatusHA", func() {
				Expect(clone.GetStatusHAHandler()).To(Equal(defaults.GetStatusHAHandler()))
			})
		})

		When("the original StatusHA handler is set and default StatusHA handler is not set", func() {
			BeforeEach(func() {
				// original and defaults are deep copies so that we can change them
				defaults = roleSpecTwo.DeepCopy()
				original = roleSpecOne.DeepCopy()
				defaults.StatusHA = nil
			})

			It("clone should be equal to the original", func() {
				Expect(clone).To(Equal(original))
			})

			It("clone StatusHA handler should equal the original StatusHA handler", func() {
				Expect(clone.GetStatusHAHandler()).To(Equal(original.GetStatusHAHandler()))
			})
		})

		When("the original StatusHA is set and default StatusHA is set", func() {
			BeforeEach(func() {
				// original and defaults are deep copies so that we can change them
				defaults = roleSpecTwo.DeepCopy()
				original = roleSpecOne.DeepCopy()
			})

			It("clone should be equal to the original", func() {
				Expect(clone).To(Equal(original))
			})

			It("clone StatusHA handler should equal the original StatusHA handler", func() {
				Expect(clone.GetStatusHAHandler()).To(Equal(original.GetStatusHAHandler()))
			})
		})

		When("the original StatusHA is not set and default StatusHA is not set", func() {
			BeforeEach(func() {
				// original and defaults are deep copies so that we can change them
				defaults = roleSpecTwo.DeepCopy()
				original = roleSpecOne.DeepCopy()
				defaults.StatusHA = nil
				original.StatusHA = nil
			})

			It("clone should be equal to the original", func() {
				Expect(clone).To(Equal(original))
			})

			It("clone StatusHA handler should equal the global default StatusHA handler", func() {
				expected := coherence.GetDefaultStatusHAHandler()
				Expect(clone.GetStatusHAHandler()).To(Equal(expected))
			})
		})

		When("the original Env is not set and default Env is set", func() {
			BeforeEach(func() {
				// original and defaults are deep copies so that we can change them
				defaults = roleSpecTwo.DeepCopy()
				original = roleSpecOne.DeepCopy()
				original.Env = nil
				defaults.Env = envTwo
			})

			It("clone should be equal to the original with the Env field from the defaults", func() {
				// expected is a deep copy of original so that we can change the
				// expected without changing original
				expected := original.DeepCopy()
				expected.Env = defaults.Env

				Expect(clone).To(Equal(expected))
			})
		})

		When("the original Env is set and default Env is not set", func() {
			BeforeEach(func() {
				// original and defaults are deep copies so that we can change them
				defaults = roleSpecTwo.DeepCopy()
				original = roleSpecOne.DeepCopy()
				original.Env = envOne
				defaults.Env = nil
			})

			It("clone should be equal to the original", func() {
				Expect(clone).To(Equal(original))
			})
		})

		When("the original Env is set and default Env is set", func() {
			BeforeEach(func() {
				// original and defaults are deep copies so that we can change them
				defaults = roleSpecTwo.DeepCopy()
				original = roleSpecOne.DeepCopy()
				original.Env = envOne
				defaults.Env = envTwo
			})

			It("clone should have the combined Env", func() {
				// expected is a deep copy of original so that we can change the
				// expected without changing original
				expected := original.DeepCopy()
				expected.Env = []corev1.EnvVar{{Name: "foo1", Value: "1"}, {Name: "foo2", Value: "2"}, {Name: "foo3", Value: "3"}, {Name: "foo4", Value: "4"}}

				Expect(clone).To(Equal(expected))
			})
		})

		When("the original Env has duplicate keys to the default Env", func() {
			BeforeEach(func() {
				// original and defaults are deep copies so that we can change them
				defaults = roleSpecTwo.DeepCopy()
				original = roleSpecOne.DeepCopy()

				original.Env = []corev1.EnvVar{{Name: "foo1", Value: "1"}, {Name: "foo2", Value: "2"}, {Name: "foo3", Value: "changed"}}
				defaults.Env = []corev1.EnvVar{{Name: "foo3", Value: "3"}, {Name: "foo4", Value: "4"}}
			})

			It("clone should have the combined Env with the duplicate key mapped to the originals value", func() {
				// expected is a deep copy of original so that we can change the
				// expected without changing original
				expected := original.DeepCopy()
				expected.Env = []corev1.EnvVar{{Name: "foo1", Value: "1"}, {Name: "foo2", Value: "2"}, {Name: "foo3", Value: "changed"}, {Name: "foo4", Value: "4"}}

				Expect(clone).To(Equal(expected))
			})
		})

		When("the original Env has duplicate keys to the default Env where the value is empty string", func() {
			BeforeEach(func() {
				// original and defaults are deep copies so that we can change them
				defaults = roleSpecTwo.DeepCopy()
				original = roleSpecOne.DeepCopy()

				original.Env = []corev1.EnvVar{{Name: "foo1", Value: "1"}, {Name: "foo2", Value: "2"}, {Name: "foo3", Value: ""}}
				defaults.Env = []corev1.EnvVar{{Name: "foo3", Value: "3"}, {Name: "foo4", Value: "4"}}
			})

			It("clone should have the combined Env with the duplicate key removed", func() {
				// expected is a deep copy of original so that we can change the
				// expected without changing original
				expected := original.DeepCopy()
				expected.Env = []corev1.EnvVar{{Name: "foo1", Value: "1"}, {Name: "foo2", Value: "2"}, {Name: "foo4", Value: "4"}}

				Expect(clone).To(Equal(expected))
			})
		})

		When("the original Annotations is not set and default Annotations is set", func() {
			BeforeEach(func() {
				// original and defaults are deep copies so that we can change them
				defaults = roleSpecTwo.DeepCopy()
				original = roleSpecOne.DeepCopy()
				original.Annotations = nil
				defaults.Annotations = annotationsTwo
			})

			It("clone should be equal to the original with the Annotations field from the defaults", func() {
				// expected is a deep copy of original so that we can change the
				// expected without changing original
				expected := original.DeepCopy()
				expected.Annotations = defaults.Annotations

				Expect(clone).To(Equal(expected))
			})
		})

		When("the original Annotations is set and default Annotations is not set", func() {
			BeforeEach(func() {
				// original and defaults are deep copies so that we can change them
				defaults = roleSpecTwo.DeepCopy()
				original = roleSpecOne.DeepCopy()
				original.Annotations = annotationsOne
				defaults.Annotations = nil
			})

			It("clone should be equal to the original", func() {
				Expect(clone).To(Equal(original))
			})
		})

		When("the original Annotations is set and default Annotations is set", func() {
			BeforeEach(func() {
				// original and defaults are deep copies so that we can change them
				defaults = roleSpecTwo.DeepCopy()
				original = roleSpecOne.DeepCopy()
				original.Annotations = annotationsOne
				defaults.Annotations = annotationsTwo
			})

			It("clone should have the combined Annotations", func() {
				// expected is a deep copy of original so that we can change the
				// expected without changing original
				expected := original.DeepCopy()
				expected.Annotations = map[string]string{"anno1": "1", "anno2": "2", "anno3": "3", "anno4": "4"}

				Expect(clone).To(Equal(expected))
			})
		})

		When("the original Annotations has duplicate keys to the default Annotations", func() {
			BeforeEach(func() {
				// original and defaults are deep copies so that we can change them
				defaults = roleSpecTwo.DeepCopy()
				original = roleSpecOne.DeepCopy()

				original.Annotations = map[string]string{"anno1": "1", "anno2": "2", "anno3": "changed"}
				defaults.Annotations = map[string]string{"anno3": "3", "anno4": "4"}
			})

			It("clone should have the combined Annotations with the duplicate key mapped to the originals value", func() {
				// expected is a deep copy of original so that we can change the
				// expected without changing original
				expected := original.DeepCopy()
				expected.Annotations = map[string]string{"anno1": "1", "anno2": "2", "anno3": "changed", "anno4": "4"}

				Expect(clone).To(Equal(expected))
			})
		})

		When("the original Annotations has duplicate keys to the default Annotations where the value is empty string", func() {
			BeforeEach(func() {
				// original and defaults are deep copies so that we can change them
				defaults = roleSpecTwo.DeepCopy()
				original = roleSpecOne.DeepCopy()

				original.Annotations = map[string]string{"anno1": "1", "anno2": "2", "anno3": ""}
				defaults.Annotations = map[string]string{"anno3": "3", "anno4": "4"}
			})

			It("clone should have the combined Annotations with the duplicate key removed", func() {
				// expected is a deep copy of original so that we can change the
				// expected without changing original
				expected := original.DeepCopy()
				expected.Annotations = map[string]string{"anno1": "1", "anno2": "2", "anno4": "4"}

				Expect(clone).To(Equal(expected))
			})
		})

		When("the original NodeSelector is not set and default NodeSelector is set", func() {
			BeforeEach(func() {
				// original and defaults are deep copies so that we can change them
				defaults = roleSpecTwo.DeepCopy()
				original = roleSpecOne.DeepCopy()
				original.NodeSelector = nil
				defaults.NodeSelector = nodeSelectorTwo
			})

			It("clone should be equal to the original with the NodeSelector field from the defaults", func() {
				// expected is a deep copy of original so that we can change the
				// expected without changing original
				expected := original.DeepCopy()
				expected.NodeSelector = defaults.NodeSelector

				Expect(clone).To(Equal(expected))
			})
		})

		When("the original NodeSelector is set and default NodeSelector is not set", func() {
			BeforeEach(func() {
				// original and defaults are deep copies so that we can change them
				defaults = roleSpecTwo.DeepCopy()
				original = roleSpecOne.DeepCopy()
				original.NodeSelector = nodeSelectorOne
				defaults.NodeSelector = nil
			})

			It("clone should be equal to the original", func() {
				Expect(clone).To(Equal(original))
			})
		})

		When("the original NodeSelector is set and default NodeSelector is set", func() {
			BeforeEach(func() {
				// original and defaults are deep copies so that we can change them
				defaults = roleSpecTwo.DeepCopy()
				original = roleSpecOne.DeepCopy()
				original.NodeSelector = nodeSelectorOne
				defaults.NodeSelector = nodeSelectorTwo
			})

			It("clone should have the combined NodeSelector", func() {
				// expected is a deep copy of original so that we can change the
				// expected without changing original
				expected := original.DeepCopy()
				expected.NodeSelector = map[string]string{"one": "1", "two": "2", "three": "3", "four": "4"}

				Expect(clone).To(Equal(expected))
			})
		})

		When("the original NodeSelector has duplicate keys to the default NodeSelector", func() {
			BeforeEach(func() {
				// original and defaults are deep copies so that we can change them
				defaults = roleSpecTwo.DeepCopy()
				original = roleSpecOne.DeepCopy()

				original.NodeSelector = map[string]string{"one": "1", "two": "2", "three": "changed"}
				defaults.NodeSelector = map[string]string{"three": "3", "four": "4"}
			})

			It("clone should have the combined NodeSelector with the duplicate key mapped to the originals value", func() {
				// expected is a deep copy of original so that we can change the
				// expected without changing original
				expected := original.DeepCopy()
				expected.NodeSelector = map[string]string{"one": "1", "two": "2", "three": "changed", "four": "4"}

				Expect(clone).To(Equal(expected))
			})
		})

		When("the original NodeSelector has duplicate keys to the default NodeSelector where the value is empty string", func() {
			BeforeEach(func() {
				// original and defaults are deep copies so that we can change them
				defaults = roleSpecTwo.DeepCopy()
				original = roleSpecOne.DeepCopy()

				original.NodeSelector = map[string]string{"one": "1", "two": "2", "three": ""}
				defaults.NodeSelector = map[string]string{"three": "3", "four": "4"}
			})

			It("clone should have the combined NodeSelector with the duplicate key removed", func() {
				// expected is a deep copy of original so that we can change the
				// expected without changing original
				expected := original.DeepCopy()
				expected.NodeSelector = map[string]string{"one": "1", "two": "2", "four": "4"}

				Expect(clone).To(Equal(expected))
			})
		})
	})

	Context("Getting Replica count", func() {
		var role coherence.CoherenceRoleSpec
		var replicas *int32

		JustBeforeEach(func() {
			role = coherence.CoherenceRoleSpec{Replicas: replicas}
		})

		When("Replicas is not set", func() {
			BeforeEach(func() {
				replicas = nil
			})

			It("should return the default replica count", func() {
				Expect(role.GetReplicas()).To(Equal(coherence.DefaultReplicas))
			})
		})

		When("Replicas is set", func() {
			BeforeEach(func() {
				replicas = int32Ptr(100)
			})

			It("should return the specified replica count", func() {
				Expect(role.GetReplicas()).To(Equal(*replicas))
			})
		})
	})

	When("Getting the full role name", func() {
		var cluster coherence.CoherenceCluster
		var role coherence.CoherenceRoleSpec

		BeforeEach(func() {
			cluster = coherence.CoherenceCluster{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-cluster",
				},
			}

			role = coherence.CoherenceRoleSpec{
				Role: "storage",
			}
		})

		It("should return the specified replica count", func() {
			Expect(role.GetFullRoleName(&cluster)).To(Equal("test-cluster-storage"))
		})
	})

	Context("Getting the role name", func() {
		var role *coherence.CoherenceRoleSpec
		var name string

		JustBeforeEach(func() {
			name = role.GetRoleName()
		})

		When("role is not set", func() {
			BeforeEach(func() {
				role = &coherence.CoherenceRoleSpec{Role: ""}
			})

			It("should use the default name", func() {
				Expect(name).To(Equal(coherence.DefaultRoleName))
			})
		})

		When("role is set", func() {
			BeforeEach(func() {
				role = &coherence.CoherenceRoleSpec{Role: "test-role"}
			})

			It("should use the default name", func() {
				Expect(name).To(Equal("test-role"))
			})
		})

		When("role is nil", func() {
			BeforeEach(func() {
				role = nil
			})

			It("should use the default name", func() {
				Expect(name).To(Equal(coherence.DefaultRoleName))
			})
		})

	})

})
