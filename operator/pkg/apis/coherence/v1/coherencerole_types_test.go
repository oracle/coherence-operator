package v1_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	coherence "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
)

var _ = Describe("Testing CoherenceRole struct", func() {

	Context("Copying an CoherenceRoleSpec using DeepCopyWithDefaults", func() {
		var (
			original *coherence.CoherenceRoleSpec
			defaults *coherence.CoherenceRoleSpec
			clone    *coherence.CoherenceRoleSpec

			roleNameOne = "one"
			roleNameTwo = "two"

			storageOne = false
			storageTwo = true

			replicasOne = int32Pointer(19)
			replicasTwo = int32Pointer(20)

			probeOne = &coherence.ReadinessProbeSpec{InitialDelaySeconds: int32Pointer(10)}
			probeTwo = &coherence.ReadinessProbeSpec{InitialDelaySeconds: int32Pointer(99)}

			imagesOne = &coherence.Images{Coherence: &coherence.ImageSpec{Image: stringPointer("foo:1.0")}}
			imagesTwo = &coherence.Images{Coherence: &coherence.ImageSpec{Image: stringPointer("foo:2.0")}}

			scalingOne = coherence.SafeScaling
			scalingTwo = coherence.ParallelScaling

			labelsOne = map[string]string{"one": "two"}
			labelsTwo = map[string]string{"three": "four"}
		)

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
				original = &coherence.CoherenceRoleSpec{
					RoleName:       roleNameOne,
					Replicas:       replicasOne,
					StorageEnabled: &storageOne,
					ReadinessProbe: probeOne,
					Images:         imagesOne,
					Labels:         &labelsOne,
					ScalingPolicy:  &scalingOne,
				}

				defaults = nil
			})

			It("clone should be the original", func() {
				Expect(clone).To(Equal(original))
			})

			It("should copy the original RoleName field", func() {
				Expect(clone.RoleName).To(Equal(original.RoleName))
			})

			It("should copy the original Replicas field", func() {
				Expect(*clone.Replicas).To(Equal(*original.Replicas))
			})

			It("should copy the original StorageEnabled field", func() {
				Expect(*clone.StorageEnabled).To(Equal(*original.StorageEnabled))
			})

			It("should copy the original ReadinessProbe field", func() {
				Expect(*clone.ReadinessProbe).To(Equal(*original.ReadinessProbe))
			})

			It("should copy the original Images field", func() {
				Expect(*clone.Images).To(Equal(*original.Images))
			})

			It("should copy the original Labels field", func() {
				Expect(*clone.Labels).To(Equal(*original.Labels))
			})

			It("should copy the original ScalingPolicy field", func() {
				Expect(*clone.ScalingPolicy).To(Equal(*original.ScalingPolicy))
			})
		})

		When("original is nil", func() {
			BeforeEach(func() {
				defaults = &coherence.CoherenceRoleSpec{
					RoleName:       roleNameTwo,
					Replicas:       replicasTwo,
					StorageEnabled: &storageTwo,
					ReadinessProbe: probeTwo,
					Images:         imagesTwo,
					Labels:         &labelsTwo,
					ScalingPolicy:  &scalingTwo,
				}

				original = nil
			})

			It("should copy the defaults RoleName field", func() {
				Expect(clone.RoleName).To(Equal(defaults.RoleName))
			})

			It("should copy the defaults Replicas field", func() {
				Expect(*clone.Replicas).To(Equal(*defaults.Replicas))
			})

			It("should copy the defaults StorageEnabled field", func() {
				Expect(*clone.StorageEnabled).To(Equal(*defaults.StorageEnabled))
			})

			It("should copy the defaults ReadinessProbe field", func() {
				Expect(*clone.ReadinessProbe).To(Equal(*defaults.ReadinessProbe))
			})

			It("should copy the defaults Images field", func() {
				Expect(*clone.Images).To(Equal(*defaults.Images))
			})

			It("should copy the defaults Labels field", func() {
				Expect(*clone.Labels).To(Equal(*defaults.Labels))
			})

			It("should copy the defaults ScalingPolicy field", func() {
				Expect(*clone.ScalingPolicy).To(Equal(*defaults.ScalingPolicy))
			})
		})

		When("all fields in the original are set", func() {
			BeforeEach(func() {
				original = &coherence.CoherenceRoleSpec{
					RoleName:       roleNameOne,
					Replicas:       replicasOne,
					StorageEnabled: &storageOne,
					ReadinessProbe: probeOne,
					Images:         imagesOne,
					Labels:         &labelsOne,
					ScalingPolicy:  &scalingOne,
				}

				defaults = &coherence.CoherenceRoleSpec{
					RoleName:       roleNameTwo,
					Replicas:       replicasTwo,
					StorageEnabled: &storageTwo,
					ReadinessProbe: probeTwo,
					Images:         imagesTwo,
					Labels:         &labelsTwo,
					ScalingPolicy:  &scalingTwo,
				}
			})

			It("should copy the original RoleName field", func() {
				Expect(clone.RoleName).To(Equal(original.RoleName))
			})

			It("should copy the original Replicas field", func() {
				Expect(*clone.Replicas).To(Equal(*original.Replicas))
			})

			It("should copy the original StorageEnabled field", func() {
				Expect(*clone.StorageEnabled).To(Equal(*original.StorageEnabled))
			})

			It("should copy the original ReadinessProbe field", func() {
				Expect(*clone.ReadinessProbe).To(Equal(*original.ReadinessProbe))
			})

			It("should copy the original Images field", func() {
				Expect(*clone.Images).To(Equal(*original.Images))
			})

			It("should copy the original Labels field", func() {
				Expect(*clone.Labels).To(Equal(*original.Labels))
			})

			It("should copy the original ScalingPolicy field", func() {
				Expect(*clone.ScalingPolicy).To(Equal(*original.ScalingPolicy))
			})
		})

		When("the original RoleName is not set", func() {
			BeforeEach(func() {
				original = &coherence.CoherenceRoleSpec{
					Replicas:       replicasOne,
					StorageEnabled: &storageOne,
					ReadinessProbe: probeOne,
					Images:         imagesOne,
					Labels:         &labelsOne,
					ScalingPolicy:  &scalingOne,
				}

				defaults = &coherence.CoherenceRoleSpec{
					RoleName:       roleNameTwo,
					Replicas:       replicasTwo,
					StorageEnabled: &storageTwo,
					ReadinessProbe: probeTwo,
					Images:         imagesTwo,
					Labels:         &labelsTwo,
					ScalingPolicy:  &scalingTwo,
				}
			})

			It("should copy the defaults RoleName field", func() {
				Expect(clone.RoleName).To(Equal(defaults.RoleName))
			})

			It("should copy the original Replicas field", func() {
				Expect(*clone.Replicas).To(Equal(*original.Replicas))
			})

			It("should copy the original StorageEnabled field", func() {
				Expect(*clone.StorageEnabled).To(Equal(*original.StorageEnabled))
			})

			It("should copy the original ReadinessProbe field", func() {
				Expect(*clone.ReadinessProbe).To(Equal(*original.ReadinessProbe))
			})

			It("should copy the original Images field", func() {
				Expect(*clone.Images).To(Equal(*original.Images))
			})

			It("should copy the original Labels field", func() {
				Expect(*clone.Labels).To(Equal(*original.Labels))
			})

			It("should copy the original ScalingPolicy field", func() {
				Expect(*clone.ScalingPolicy).To(Equal(*original.ScalingPolicy))
			})
		})

		When("the original Replicas is not set", func() {
			BeforeEach(func() {
				original = &coherence.CoherenceRoleSpec{
					RoleName:       roleNameOne,
					StorageEnabled: &storageOne,
					ReadinessProbe: probeOne,
					Images:         imagesOne,
					Labels:         &labelsOne,
					ScalingPolicy:  &scalingOne,
				}

				defaults = &coherence.CoherenceRoleSpec{
					RoleName:       roleNameTwo,
					Replicas:       replicasTwo,
					StorageEnabled: &storageTwo,
					ReadinessProbe: probeTwo,
					Images:         imagesTwo,
					Labels:         &labelsTwo,
					ScalingPolicy:  &scalingTwo,
				}
			})

			It("should copy the original RoleName field", func() {
				Expect(clone.RoleName).To(Equal(original.RoleName))
			})

			It("should copy the defaults Replicas field", func() {
				Expect(*clone.Replicas).To(Equal(*defaults.Replicas))
			})

			It("should copy the original StorageEnabled field", func() {
				Expect(*clone.StorageEnabled).To(Equal(*original.StorageEnabled))
			})

			It("should copy the original ReadinessProbe field", func() {
				Expect(*clone.ReadinessProbe).To(Equal(*original.ReadinessProbe))
			})

			It("should copy the original Images field", func() {
				Expect(*clone.Images).To(Equal(*original.Images))
			})

			It("should copy the original Labels field", func() {
				Expect(*clone.Labels).To(Equal(*original.Labels))
			})

			It("should copy the original ScalingPolicy field", func() {
				Expect(*clone.ScalingPolicy).To(Equal(*original.ScalingPolicy))
			})
		})

		When("the original ReadinessProbe is not set", func() {
			BeforeEach(func() {
				original = &coherence.CoherenceRoleSpec{
					RoleName:       roleNameOne,
					Replicas:       replicasOne,
					StorageEnabled: &storageOne,
					Images:         imagesOne,
					Labels:         &labelsOne,
					ScalingPolicy:  &scalingOne,
				}

				defaults = &coherence.CoherenceRoleSpec{
					RoleName:       roleNameTwo,
					Replicas:       replicasTwo,
					StorageEnabled: &storageTwo,
					ReadinessProbe: probeTwo,
					Images:         imagesTwo,
					Labels:         &labelsTwo,
					ScalingPolicy:  &scalingTwo,
				}
			})

			It("should copy the original RoleName field", func() {
				Expect(clone.RoleName).To(Equal(original.RoleName))
			})

			It("should copy the original Replicas field", func() {
				Expect(*clone.Replicas).To(Equal(*original.Replicas))
			})

			It("should copy the original StorageEnabled field", func() {
				Expect(*clone.StorageEnabled).To(Equal(*original.StorageEnabled))
			})

			It("should copy the defaults ReadinessProbe field", func() {
				Expect(*clone.ReadinessProbe).To(Equal(*defaults.ReadinessProbe))
			})

			It("should copy the original Images field", func() {
				Expect(*clone.Images).To(Equal(*original.Images))
			})

			It("should copy the original Labels field", func() {
				Expect(*clone.Labels).To(Equal(*original.Labels))
			})

			It("should copy the original ScalingPolicy field", func() {
				Expect(*clone.ScalingPolicy).To(Equal(*original.ScalingPolicy))
			})
		})

		When("the original Images is not set", func() {
			BeforeEach(func() {
				original = &coherence.CoherenceRoleSpec{
					RoleName:       roleNameOne,
					Replicas:       replicasOne,
					StorageEnabled: &storageOne,
					ReadinessProbe: probeOne,
					Labels:         &labelsOne,
					ScalingPolicy:  &scalingOne,
				}

				defaults = &coherence.CoherenceRoleSpec{
					RoleName:       roleNameTwo,
					Replicas:       replicasTwo,
					StorageEnabled: &storageTwo,
					ReadinessProbe: probeTwo,
					Images:         imagesTwo,
					Labels:         &labelsTwo,
					ScalingPolicy:  &scalingTwo,
				}
			})

			It("should copy the original RoleName field", func() {
				Expect(clone.RoleName).To(Equal(original.RoleName))
			})

			It("should copy the original Replicas field", func() {
				Expect(*clone.Replicas).To(Equal(*original.Replicas))
			})

			It("should copy the original StorageEnabled field", func() {
				Expect(*clone.StorageEnabled).To(Equal(*original.StorageEnabled))
			})

			It("should copy the original ReadinessProbe field", func() {
				Expect(*clone.ReadinessProbe).To(Equal(*original.ReadinessProbe))
			})

			It("should copy the defaults Images field", func() {
				Expect(*clone.Images).To(Equal(*defaults.Images))
			})

			It("should copy the original Labels field", func() {
				Expect(*clone.Labels).To(Equal(*original.Labels))
			})

			It("should copy the original ScalingPolicy field", func() {
				Expect(*clone.ScalingPolicy).To(Equal(*original.ScalingPolicy))
			})
		})

		When("the original Labels is not set", func() {
			BeforeEach(func() {
				original = &coherence.CoherenceRoleSpec{
					RoleName:       roleNameOne,
					Replicas:       replicasOne,
					StorageEnabled: &storageOne,
					ReadinessProbe: probeOne,
					Images:         imagesOne,
					ScalingPolicy:  &scalingOne,
				}

				defaults = &coherence.CoherenceRoleSpec{
					RoleName:       roleNameTwo,
					Replicas:       replicasTwo,
					StorageEnabled: &storageTwo,
					ReadinessProbe: probeTwo,
					Images:         imagesTwo,
					Labels:         &labelsTwo,
					ScalingPolicy:  &scalingTwo,
				}
			})

			It("should copy the original RoleName field", func() {
				Expect(clone.RoleName).To(Equal(original.RoleName))
			})

			It("should copy the original Replicas field", func() {
				Expect(*clone.Replicas).To(Equal(*original.Replicas))
			})

			It("should copy the original StorageEnabled field", func() {
				Expect(*clone.StorageEnabled).To(Equal(*original.StorageEnabled))
			})

			It("should copy the original ReadinessProbe field", func() {
				Expect(*clone.ReadinessProbe).To(Equal(*original.ReadinessProbe))
			})

			It("should copy the original Images field", func() {
				Expect(*clone.Images).To(Equal(*original.Images))
			})

			It("should copy the defaults Labels field", func() {
				Expect(*clone.Labels).To(Equal(*defaults.Labels))
			})

			It("should copy the original ScalingPolicy field", func() {
				Expect(*clone.ScalingPolicy).To(Equal(*original.ScalingPolicy))
			})
		})

		When("the original ScalingPolicy is not set", func() {
			BeforeEach(func() {
				original = &coherence.CoherenceRoleSpec{
					RoleName:       roleNameOne,
					Replicas:       replicasOne,
					StorageEnabled: &storageOne,
					ReadinessProbe: probeOne,
					Images:         imagesOne,
					Labels:         &labelsOne,
				}

				defaults = &coherence.CoherenceRoleSpec{
					RoleName:       roleNameTwo,
					Replicas:       replicasTwo,
					StorageEnabled: &storageTwo,
					ReadinessProbe: probeTwo,
					Images:         imagesTwo,
					Labels:         &labelsTwo,
					ScalingPolicy:  &scalingTwo,
				}
			})

			It("should copy the original RoleName field", func() {
				Expect(clone.RoleName).To(Equal(original.RoleName))
			})

			It("should copy the original Replicas field", func() {
				Expect(*clone.Replicas).To(Equal(*original.Replicas))
			})

			It("should copy the original StorageEnabled field", func() {
				Expect(*clone.StorageEnabled).To(Equal(*original.StorageEnabled))
			})

			It("should copy the original ReadinessProbe field", func() {
				Expect(*clone.ReadinessProbe).To(Equal(*original.ReadinessProbe))
			})

			It("should copy the original Images field", func() {
				Expect(*clone.Images).To(Equal(*original.Images))
			})

			It("should copy the original Labels field", func() {
				Expect(*clone.Labels).To(Equal(*original.Labels))
			})

			It("should copy the defaults ScalingPolicy field", func() {
				Expect(*clone.ScalingPolicy).To(Equal(*defaults.ScalingPolicy))
			})
		})
	})
})
