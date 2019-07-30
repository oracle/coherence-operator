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
		roleOne.RoleName = "storage"
		roleTwo := coherence.CoherenceRoleSpec{}
		roleTwo.RoleName = "proxy"

		cluster := coherence.CoherenceCluster{}
		cluster.Spec.Roles = make([]coherence.CoherenceRoleSpec, 2)
		cluster.Spec.Roles[0] = roleOne
		cluster.Spec.Roles[1] = roleTwo

		Expect(cluster.GetRole("storage")).To(Equal(roleOne))
	})

	It("cluster does not have role with name", func() {
		roleOne := coherence.CoherenceRoleSpec{}
		roleOne.RoleName = "storage"
		roleTwo := coherence.CoherenceRoleSpec{}
		roleTwo.RoleName = "proxy"

		cluster := coherence.CoherenceCluster{}
		cluster.Spec.Roles = make([]coherence.CoherenceRoleSpec, 2)
		cluster.Spec.Roles[0] = roleOne
		cluster.Spec.Roles[1] = roleTwo

		Expect(cluster.GetRole("foo")).To(Equal(coherence.CoherenceRoleSpec{}))
	})

	It("set cluster role", func() {
		roleOne := coherence.CoherenceRoleSpec{}
		roleOne.RoleName = "storage"
		roleTwo := coherence.CoherenceRoleSpec{}
		roleTwo.RoleName = "proxy"

		roleUpdate := coherence.CoherenceRoleSpec{}
		roleUpdate.RoleName = "storage"
		roleUpdate.Replicas = int32Pointer(19)

		cluster := coherence.CoherenceCluster{}
		cluster.Spec.Roles = make([]coherence.CoherenceRoleSpec, 2)
		cluster.Spec.Roles[0] = roleOne
		cluster.Spec.Roles[1] = roleTwo

		cluster.SetRole(roleUpdate)

		Expect(cluster.GetRole("storage")).To(Equal(roleUpdate))
	})

	It("set cluster role when role not in cluster", func() {
		roleOne := coherence.CoherenceRoleSpec{}
		roleOne.RoleName = "storage"
		roleTwo := coherence.CoherenceRoleSpec{}
		roleTwo.RoleName = "proxy"

		roleUpdate := coherence.CoherenceRoleSpec{}
		roleUpdate.RoleName = "foo"
		roleUpdate.Replicas = int32Pointer(19)

		cluster := coherence.CoherenceCluster{}
		cluster.Spec.Roles = make([]coherence.CoherenceRoleSpec, 2)
		cluster.Spec.Roles[0] = roleOne
		cluster.Spec.Roles[1] = roleTwo

		cluster.SetRole(roleUpdate)

		Expect(len(cluster.Spec.Roles)).To(Equal(2))
		Expect(cluster.GetRole("storage")).To(Equal(roleOne))
		Expect(cluster.GetRole("proxy")).To(Equal(roleTwo))
	})
})
