package local

import (
	goctx "context"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	coherence "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	"testing"
	"time"

	. "github.com/onsi/gomega"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestClusterWithSingleRoleWithSingleMember(t *testing.T) {
	var (
		clusterName        = "test-cluster"
		roleName           = "one"
		roleFullName       = clusterName + "-" + roleName
		replicas     int32 = 1
	)
	g := NewGomegaWithT(t)

	f := framework.Global

	ctx := helper.CreateTestContext(t)
	defer helper.DumpOperatorLogsAndCleanup(t, ctx)

	namespace, err := ctx.GetNamespace()
	g.Expect(err).NotTo(HaveOccurred())

	roleOne := coherence.CoherenceRoleSpec{
		Role:     roleName,
		Replicas: &replicas,
	}

	cluster := coherence.CoherenceCluster{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      clusterName,
		},
		Spec: coherence.CoherenceClusterSpec{
			ImagePullSecrets: helper.GetImagePullSecrets(),
			CoherenceRoleSpec: coherence.CoherenceRoleSpec{
				ReadinessProbe: helper.Readiness,
			},
			Roles: []coherence.CoherenceRoleSpec{roleOne},
		},
	}

	err = f.Client.Create(goctx.TODO(), &cluster, helper.DefaultCleanup(ctx))
	g.Expect(err).NotTo(HaveOccurred())

	role, err := helper.WaitForCoherenceRole(f, namespace, roleFullName, helper.RetryInterval, helper.Timeout, t)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(role.Spec.GetRoleName()).To(Equal(roleName))
	g.Expect(role.Spec.GetReplicas()).To(Equal(replicas))

	sts, err := helper.WaitForStatefulSetForRole(f.KubeClient, namespace, &cluster, role.Spec, time.Second*10, time.Minute*5, t)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(sts.Status.ReadyReplicas).To(Equal(replicas))
}
