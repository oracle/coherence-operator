package v1_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

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

			scalingOne = coherence.SafeScaling
			scalingTwo = coherence.ParallelScaling

			labelsOne = map[string]string{"one": "1", "two": "2"}
			labelsTwo = map[string]string{"three": "3", "four": "4"}

			roleSpecOne = &coherence.CoherenceRoleSpec{
				Role:           roleNameOne,
				Replicas:       replicasOne,
				Images:         imagesOne,
				StorageEnabled: &storageOne,
				ScalingPolicy:  &scalingOne,
				ReadinessProbe: probeOne,
				Labels:         nil,
			}

			roleSpecTwo = &coherence.CoherenceRoleSpec{
				Role:           roleNameTwo,
				Replicas:       replicasTwo,
				Images:         imagesTwo,
				StorageEnabled: &storageTwo,
				ScalingPolicy:  &scalingTwo,
				ReadinessProbe: probeTwo,
				Labels:         nil,
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
				defaults.Labels = &labelsTwo
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
				original.Labels = &labelsOne
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
				original.Labels = &labelsOne
				defaults.Labels = &labelsTwo
			})

			It("clone should have the combined Labels", func() {
				// expected is a deep copy of original so that we can change the
				// expected without changing original
				expected := original.DeepCopy()
				expected.Labels = mapPtr(map[string]string{"one": "1", "two": "2", "three": "3", "four": "4"})

				Expect(clone).To(Equal(expected))
			})
		})

		When("the original Labels has duplicate keys to the default Labels", func() {
			BeforeEach(func() {
				// original and defaults are deep copies so that we can change them
				defaults = roleSpecTwo.DeepCopy()
				original = roleSpecOne.DeepCopy()

				original.Labels = mapPtr(map[string]string{"one": "1", "two": "2", "three": "changed"})
				defaults.Labels = mapPtr(map[string]string{"three": "3", "four": "4"})
			})

			It("clone should have the combined Labels with the duplicate key mapped to the originals value", func() {
				// expected is a deep copy of original so that we can change the
				// expected without changing original
				expected := original.DeepCopy()
				expected.Labels = mapPtr(map[string]string{"one": "1", "two": "2", "three": "changed", "four": "4"})

				Expect(clone).To(Equal(expected))
			})
		})

		When("the original Labels has duplicate keys to the default Labels where the value is empty string", func() {
			BeforeEach(func() {
				// original and defaults are deep copies so that we can change them
				defaults = roleSpecTwo.DeepCopy()
				original = roleSpecOne.DeepCopy()

				original.Labels = mapPtr(map[string]string{"one": "1", "two": "2", "three": ""})
				defaults.Labels = mapPtr(map[string]string{"three": "3", "four": "4"})
			})

			It("clone should have the combined Labels with the duplicate key removed", func() {
				// expected is a deep copy of original so that we can change the
				// expected without changing original
				expected := original.DeepCopy()
				expected.Labels = mapPtr(map[string]string{"one": "1", "two": "2", "four": "4"})

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
