package v1_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	appv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	coherence "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
)

var _ = Describe("Testing CoherenceRoleSpec struct", func() {

	Context("Copying an CoherenceRoleSpec using DeepCopyWithDefaults", func() {
		var (
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

			JavaOptsOne = "-Dcoherence.log.level=9"
			JavaOptsTwo = "-Dcoherence.log.level=5"

			portsOne = map[string]int32{"port1": 8081, "port2": 8082}
			portsTwo = map[string]int32{"port3": 8083, "port4": 8084}

			envOne = map[string]string{"foo1": "1", "foo2": "2"}
			envTwo = map[string]string{"foo3": "3", "foo4": "4"}

			annotationsOne = map[string]string{"anno1": "1", "anno2": "2"}
			annotationsTwo = map[string]string{"anno3": "3", "anno4": "4"}

			podManagementPolicyOne appv1.PodManagementPolicyType = "Parallel"
			podManagementPolicyTwo appv1.PodManagementPolicyType = "OrderedReady"

			revisionHistoryLimitOne = int32Ptr(3)
			revisionHistoryLimitTwo = int32Ptr(5)

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
				JavaOpts:             &JavaOptsOne,
				Ports:                nil,
				Env:                  nil,
				Annotations:          nil,
				PodManagementPolicy:  &podManagementPolicyOne,
				RevisionHistoryLimit: revisionHistoryLimitOne,
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
				JavaOpts:             &JavaOptsTwo,
				Ports:                nil,
				Env:                  nil,
				Annotations:          nil,
				PodManagementPolicy:  &podManagementPolicyTwo,
				RevisionHistoryLimit: revisionHistoryLimitTwo,
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
				defaults.Ports = portsTwo
			})

			It("clone should be equal to the original with the Ports field from the defaults", func() {
				// expected is a deep copy of original so that we can change the
				// expected without changing original
				expected := original.DeepCopy()
				expected.Ports = defaults.Ports

				Expect(clone).To(Equal(expected))
			})
		})

		When("the original Ports is set and default Ports is not set", func() {
			BeforeEach(func() {
				// original and defaults are deep copies so that we can change them
				defaults = roleSpecTwo.DeepCopy()
				original = roleSpecOne.DeepCopy()
				original.Ports = portsOne
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
				original.Ports = portsOne
				defaults.Ports = portsTwo
			})

			It("clone should have the combined Ports", func() {
				// expected is a deep copy of original so that we can change the
				// expected without changing original
				expected := original.DeepCopy()
				expected.Ports = map[string]int32{"port1": 8081, "port2": 8082, "port3": 8083, "port4": 8084}

				Expect(clone).To(Equal(expected))
			})
		})

		When("the original Ports has duplicate keys to the default Ports", func() {
			BeforeEach(func() {
				// original and defaults are deep copies so that we can change them
				defaults = roleSpecTwo.DeepCopy()
				original = roleSpecOne.DeepCopy()

				original.Ports = map[string]int32{"port1": 8081, "port2": 8082, "port3": 18083}
				defaults.Ports = map[string]int32{"port3": 8083, "port4": 8084}
			})

			It("clone should have the combined Ports with the duplicate key mapped to the originals value", func() {
				// expected is a deep copy of original so that we can change the
				// expected without changing original
				expected := original.DeepCopy()
				expected.Ports = map[string]int32{"port1": 8081, "port2": 8082, "port3": 18083, "port4": 8084}

				Expect(clone).To(Equal(expected))
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
				expected.Env = map[string]string{"foo1": "1", "foo2": "2", "foo3": "3", "foo4": "4"}

				Expect(clone).To(Equal(expected))
			})
		})

		When("the original Env has duplicate keys to the default Env", func() {
			BeforeEach(func() {
				// original and defaults are deep copies so that we can change them
				defaults = roleSpecTwo.DeepCopy()
				original = roleSpecOne.DeepCopy()

				original.Env = map[string]string{"foo1": "1", "foo2": "2", "foo3": "changed"}
				defaults.Env = map[string]string{"foo3": "3", "foo4": "4"}
			})

			It("clone should have the combined Env with the duplicate key mapped to the originals value", func() {
				// expected is a deep copy of original so that we can change the
				// expected without changing original
				expected := original.DeepCopy()
				expected.Env = map[string]string{"foo1": "1", "foo2": "2", "foo3": "changed", "foo4": "4"}

				Expect(clone).To(Equal(expected))
			})
		})

		When("the original Env has duplicate keys to the default Env where the value is empty string", func() {
			BeforeEach(func() {
				// original and defaults are deep copies so that we can change them
				defaults = roleSpecTwo.DeepCopy()
				original = roleSpecOne.DeepCopy()

				original.Env = map[string]string{"foo1": "1", "foo2": "2", "foo3": ""}
				defaults.Env = map[string]string{"foo3": "3", "foo4": "4"}
			})

			It("clone should have the combined Env with the duplicate key removed", func() {
				// expected is a deep copy of original so that we can change the
				// expected without changing original
				expected := original.DeepCopy()
				expected.Env = map[string]string{"foo1": "1", "foo2": "2", "foo4": "4"}

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
