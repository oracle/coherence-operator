package v1_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	coherence "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
)

var _ = Describe("Testing CoherenceCluster", func() {

	It("cluster produces wka service name", func() {
		cluster := coherence.CoherenceCluster{}
		cluster.Name = "foo"

		Expect(cluster.GetWkaServiceName()).To(Equal("foo" + coherence.WKAServiceNameSuffix))
	})

	It("cluster has role with name", func() {
		roleOne := coherence.CoherenceRoleSpec{}
		roleOne.Role = "storage"
		roleTwo := coherence.CoherenceRoleSpec{}
		roleTwo.Role = "proxy"

		cluster := coherence.CoherenceCluster{}
		cluster.Spec.Roles = make([]coherence.CoherenceRoleSpec, 2)
		cluster.Spec.Roles[0] = roleOne
		cluster.Spec.Roles[1] = roleTwo

		Expect(cluster.GetRole("storage")).To(Equal(roleOne))
	})

	It("cluster does not have role with name", func() {
		roleOne := coherence.CoherenceRoleSpec{}
		roleOne.Role = "storage"
		roleTwo := coherence.CoherenceRoleSpec{}
		roleTwo.Role = "proxy"

		cluster := coherence.CoherenceCluster{}
		cluster.Spec.Roles = make([]coherence.CoherenceRoleSpec, 2)
		cluster.Spec.Roles[0] = roleOne
		cluster.Spec.Roles[1] = roleTwo

		Expect(cluster.GetRole("foo")).To(Equal(coherence.CoherenceRoleSpec{}))
	})

	It("set cluster role", func() {
		roleOne := coherence.CoherenceRoleSpec{}
		roleOne.Role = "storage"
		roleTwo := coherence.CoherenceRoleSpec{}
		roleTwo.Role = "proxy"

		roleUpdate := coherence.CoherenceRoleSpec{}
		roleUpdate.Role = "storage"
		roleUpdate.Replicas = int32Ptr(19)

		cluster := coherence.CoherenceCluster{}
		cluster.Spec.Roles = make([]coherence.CoherenceRoleSpec, 2)
		cluster.Spec.Roles[0] = roleOne
		cluster.Spec.Roles[1] = roleTwo

		cluster.SetRole(roleUpdate)

		Expect(cluster.GetRole("storage")).To(Equal(roleUpdate))
	})

	It("set cluster role when role not in cluster", func() {
		roleOne := coherence.CoherenceRoleSpec{}
		roleOne.Role = "storage"
		roleTwo := coherence.CoherenceRoleSpec{}
		roleTwo.Role = "proxy"

		roleUpdate := coherence.CoherenceRoleSpec{}
		roleUpdate.Role = "foo"
		roleUpdate.Replicas = int32Ptr(19)

		cluster := coherence.CoherenceCluster{}
		cluster.Spec.Roles = make([]coherence.CoherenceRoleSpec, 2)
		cluster.Spec.Roles[0] = roleOne
		cluster.Spec.Roles[1] = roleTwo

		cluster.SetRole(roleUpdate)

		Expect(len(cluster.Spec.Roles)).To(Equal(2))
		Expect(cluster.GetRole("storage")).To(Equal(roleOne))
		Expect(cluster.GetRole("proxy")).To(Equal(roleTwo))
	})

	Context("loading CoherenceCluster from yaml file", func() {
		var (
			cluster coherence.CoherenceCluster
			file    []string
			err     error
		)

		JustBeforeEach(func() {
			cluster, err = coherence.NewCoherenceClusterFromYaml("test-ns", file...)
		})

		When("file is valid", func() {
			BeforeEach(func() {
				file = []string{"test-coherence-cluster-one.yaml"}
			})

			It("should load fields", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(cluster).NotTo(BeNil())

				// values come from the test-coherence-cluster-one.yaml file
				Expect(cluster.Name).To(Equal("test-cluster"))
				Expect(cluster.Spec.ReadinessProbe).NotTo(BeNil())
				Expect(cluster.Spec.ReadinessProbe.InitialDelaySeconds).To(Equal(int32Ptr(10)))
				Expect(cluster.Spec.ReadinessProbe.PeriodSeconds).To(Equal(int32Ptr(30)))

				Expect(cluster.Spec.Roles).ToNot(BeNil())
				Expect(len(cluster.Spec.Roles)).To(Equal(1))

				role := cluster.Spec.Roles[0]
				Expect(role.GetRoleName()).To(Equal("one"))
				Expect(role.GetReplicas()).To(Equal(int32(1)))
				Expect(role.CacheConfig).To(Equal(stringPtr("test-cache-config.xml")))

				Expect(role.Images).NotTo(BeNil())
				Expect(role.Images.Coherence).NotTo(BeNil())
				Expect(role.Images.Coherence.Image).To(Equal(stringPtr("test/coherence:1.0")))
			})
		})

		When("file does not exist", func() {
			BeforeEach(func() {
				file = []string{"foo.yaml"}
			})

			It("should return error", func() {
				Expect(err).To(HaveOccurred())
			})
		})

		When("multiple yaml files", func() {
			BeforeEach(func() {
				file = []string{"test-coherence-cluster-one.yaml", "test-coherence-cluster-two.yaml"}
			})

			It("should load fields", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(cluster).NotTo(BeNil())

				// values come from the test-coherence-cluster-one.yaml file and are then
				// overridden or added to by the test-coherence-cluster-two.yaml file
				Expect(cluster.Name).To(Equal("test-cluster-two"))
				Expect(cluster.Spec.ReadinessProbe).NotTo(BeNil())
				Expect(cluster.Spec.ReadinessProbe.InitialDelaySeconds).To(Equal(int32Ptr(60)))
				Expect(cluster.Spec.ReadinessProbe.PeriodSeconds).To(Equal(int32Ptr(30)))

				Expect(cluster.Spec.Roles).ToNot(BeNil())
				Expect(len(cluster.Spec.Roles)).To(Equal(2))

				roleOne := cluster.Spec.Roles[0]
				Expect(roleOne.GetRoleName()).To(Equal("one"))
				Expect(roleOne.GetReplicas()).To(Equal(int32(3)))
				Expect(roleOne.CacheConfig).To(Equal(stringPtr("test-cache-config.xml")))

				Expect(roleOne.Images).NotTo(BeNil())
				Expect(roleOne.Images.Coherence).NotTo(BeNil())
				Expect(roleOne.Images.Coherence.Image).To(Equal(stringPtr("test/coherence:1.0")))

				roleTwo := cluster.Spec.Roles[1]
				Expect(roleTwo.GetRoleName()).To(Equal("two"))
				Expect(roleTwo.GetReplicas()).To(Equal(int32(3)))
			})
		})
	})
})
