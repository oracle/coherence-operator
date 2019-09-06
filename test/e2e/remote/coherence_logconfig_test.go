/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package remote

import (
	goctx "context"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	coherence "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	"github.com/oracle/coherence-operator/pkg/management"
	"github.com/oracle/coherence-operator/test/e2e/helper"

	"io/ioutil"
	"net/http"

	corev1 "k8s.io/api/core/v1"
	"os"
	"testing"

	. "github.com/onsi/gomega"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	loggingPropsFileName = "logging.properties"
	loggingConfigMapName = "logging-config-map"
)

/*
Deploy a CoherenceCluster with the log level changed in the "spec" section.
*/
func TestCoherenceLogLevelSpecSection(t *testing.T) {
	var (
		clusterName        = "mycluster"
		roleName           = coherence.DefaultRoleName
		roleFullName       = clusterName + "-" + roleName
		replicas     int32 = 3
	)
	g := NewGomegaWithT(t)

	f := framework.Global

	ctx := helper.CreateTestContext(t)
	defer ctx.Cleanup()

	namespace, err := ctx.GetNamespace()
	g.Expect(err).NotTo(HaveOccurred())

	cluster := coherence.CoherenceCluster{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      clusterName,
		},
		Spec: getCoherenceClusterSpec(true, roleName, replicas, 9, "", ""),
	}

	err = f.Client.Create(goctx.TODO(), &cluster, helper.DefaultCleanup(ctx))
	g.Expect(err).NotTo(HaveOccurred())

	role, err := helper.WaitForCoherenceRole(f, namespace, roleFullName,
		helper.RetryInterval, helper.Timeout, t)

	sts, err := helper.WaitForStatefulSetForRole(f.KubeClient, namespace, &cluster,
		role.Spec, helper.RetryInterval, helper.Timeout, t)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(sts.Status.ReadyReplicas).To(Equal(replicas))

	verifyLogConfigOnMembers(g, clusterName, roleName, int(replicas), int(9))
}

/*
Deploy a CoherenceCluster with the log level changed in the "roles" section.
*/
func TestCoherenceLogLevelRolesSection(t *testing.T) {
	var (
		clusterName        = "mycluster"
		roleName           = "storage"
		roleFullName       = clusterName + "-" + roleName
		replicas     int32 = 1
	)
	g := NewGomegaWithT(t)

	f := framework.Global

	ctx := helper.CreateTestContext(t)
	defer ctx.Cleanup()

	namespace, err := ctx.GetNamespace()
	g.Expect(err).NotTo(HaveOccurred())

	cluster := coherence.CoherenceCluster{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      clusterName,
		},
		Spec: getCoherenceClusterSpec(false, roleName, replicas, 9, "", ""),
	}

	err = f.Client.Create(goctx.TODO(), &cluster, helper.DefaultCleanup(ctx))
	g.Expect(err).NotTo(HaveOccurred())

	role, err := helper.WaitForCoherenceRole(f, namespace, roleFullName,
		helper.RetryInterval, helper.Timeout, t)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(role.Spec.GetRoleName()).To(Equal(roleName))
	g.Expect(role.Spec.GetReplicas()).To(Equal(replicas))

	sts, err := helper.WaitForStatefulSetForRole(f.KubeClient, namespace, &cluster,
		role.Spec, helper.RetryInterval, helper.Timeout, t)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(sts.Status.ReadyReplicas).To(Equal(replicas))

	verifyLogConfigOnMembers(g, clusterName, roleName, int(replicas), int(9))
}

/*
Deploy a CoherenceCluster with the logging config set to a custom file in the
"spec" section.
*/
func TestCoherenceLogConfigFileSpecSection(t *testing.T) {
	var (
		clusterName        = "mycluster"
		roleName           = coherence.DefaultRoleName
		roleFullName       = clusterName + "-" + roleName
		replicas     int32 = 3
	)
	g := NewGomegaWithT(t)

	f := framework.Global

	ctx := helper.CreateTestContext(t)
	defer ctx.Cleanup()

	namespace, err := ctx.GetNamespace()
	g.Expect(err).NotTo(HaveOccurred())

	cluster := coherence.CoherenceCluster{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      clusterName,
		},
		Spec: getCoherenceClusterSpec(true, roleName, replicas, 0, loggingPropsFileName, ""),
	}

	err = f.Client.Create(goctx.TODO(), &cluster, helper.DefaultCleanup(ctx))
	g.Expect(err).NotTo(HaveOccurred())

	role, err := helper.WaitForCoherenceRole(f, namespace, roleFullName,
		helper.RetryInterval, helper.Timeout, t)

	sts, err := helper.WaitForStatefulSetForRole(f.KubeClient, namespace, &cluster,
		role.Spec, helper.RetryInterval, helper.Timeout, t)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(sts.Status.ReadyReplicas).To(Equal(replicas))

	verifyLogConfigOnMembers(g, clusterName, roleName, int(replicas), int(5))
}

/*
Deploy a CoherenceCluster with the logging config set to a custom file in the
"roles" section.
*/
func TestCoherenceLogConfigFileRolesSection(t *testing.T) {
	var (
		clusterName        = "mycluster"
		roleName           = coherence.DefaultRoleName
		roleFullName       = clusterName + "-" + roleName
		replicas     int32 = 3
	)
	g := NewGomegaWithT(t)

	f := framework.Global

	ctx := helper.CreateTestContext(t)
	defer ctx.Cleanup()

	namespace, err := ctx.GetNamespace()
	g.Expect(err).NotTo(HaveOccurred())

	cluster := coherence.CoherenceCluster{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      clusterName,
		},
		Spec: getCoherenceClusterSpec(false, roleName, replicas, 0, loggingPropsFileName, ""),
	}

	err = f.Client.Create(goctx.TODO(), &cluster, helper.DefaultCleanup(ctx))
	g.Expect(err).NotTo(HaveOccurred())

	role, err := helper.WaitForCoherenceRole(f, namespace, roleFullName,
		helper.RetryInterval, helper.Timeout, t)

	sts, err := helper.WaitForStatefulSetForRole(f.KubeClient, namespace, &cluster,
		role.Spec, helper.RetryInterval, helper.Timeout, t)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(sts.Status.ReadyReplicas).To(Equal(replicas))

	verifyLogConfigOnMembers(g, clusterName, roleName, int(replicas), int(5))
}

/*
Create a ConfigMap containing the custom logging configuration.
Deploy a CoherenceCluster with the logging config set to the ConfigMap file in the
"spec" section.
*/
func TestCoherenceLogConfigMapSpecSection(t *testing.T) {
	var (
		clusterName        = "mycluster"
		roleName           = coherence.DefaultRoleName
		roleFullName       = clusterName + "-" + roleName
		replicas     int32 = 3
	)
	g := NewGomegaWithT(t)
	f := framework.Global

	ctx := helper.CreateTestContext(t)
	defer ctx.Cleanup()

	namespace, err := ctx.GetNamespace()
	g.Expect(err).NotTo(HaveOccurred())

	createLoggingConfigMap(g)

	cluster := coherence.CoherenceCluster{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      clusterName,
		},
		Spec: getCoherenceClusterSpec(true, roleName, replicas, 0, loggingPropsFileName,
			loggingConfigMapName),
	}

	err = f.Client.Create(goctx.TODO(), &cluster, helper.DefaultCleanup(ctx))
	g.Expect(err).NotTo(HaveOccurred())

	role, err := helper.WaitForCoherenceRole(f, namespace, roleFullName,
		helper.RetryInterval, helper.Timeout, t)

	sts, err := helper.WaitForStatefulSetForRole(f.KubeClient, namespace, &cluster,
		role.Spec, helper.RetryInterval, helper.Timeout, t)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(sts.Status.ReadyReplicas).To(Equal(replicas))

	verifyLogConfigOnMembers(g, clusterName, roleName, int(replicas), int(5))
	deleteLoggingConfigMap(namespace)
}

/*
Create a ConfigMap containing the custom logging configuration.
Deploy a CoherenceCluster with the logging config set to the ConfigMap file in the
"roles" section.
*/
func TestCoherenceLogConfigMapRolesSection(t *testing.T) {
	var (
		clusterName        = "mycluster"
		roleName           = coherence.DefaultRoleName
		roleFullName       = clusterName + "-" + roleName
		replicas     int32 = 3
	)
	g := NewGomegaWithT(t)

	f := framework.Global

	ctx := helper.CreateTestContext(t)
	defer ctx.Cleanup()

	namespace, err := ctx.GetNamespace()
	g.Expect(err).NotTo(HaveOccurred())

	createLoggingConfigMap(g)

	cluster := coherence.CoherenceCluster{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      clusterName,
		},
		// ClusterRole with logging config file
		Spec: getCoherenceClusterSpec(false, roleName, replicas, 0, loggingPropsFileName,
			loggingConfigMapName),
	}

	err = f.Client.Create(goctx.TODO(), &cluster, helper.DefaultCleanup(ctx))
	g.Expect(err).NotTo(HaveOccurred())

	role, err := helper.WaitForCoherenceRole(f, namespace, roleFullName,
		helper.RetryInterval, helper.Timeout, t)

	sts, err := helper.WaitForStatefulSetForRole(f.KubeClient, namespace, &cluster,
		role.Spec, helper.RetryInterval, helper.Timeout, t)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(sts.Status.ReadyReplicas).To(Equal(replicas))

	verifyLogConfigOnMembers(g, clusterName, roleName, int(replicas), int(5))
	deleteLoggingConfigMap(namespace)
}

func int32Ptr(x int32) *int32 {
	return &x
}

func stringPtr(x string) *string {
	return &x
}

func getCoherenceClusterSpec(addConfigToDefaultRoleSpec bool, roleName string,
	replicas, logLevel int32, loggingFile, loggingConfigMap string) coherence.CoherenceClusterSpec {
	if addConfigToDefaultRoleSpec {
		if len(loggingFile) == 0 || len(loggingConfigMap) > 0 {
			return coherence.CoherenceClusterSpec{
				CoherenceRoleSpec: coherence.CoherenceRoleSpec{
					Logging: getCoherenceClusterLoggingSpec(logLevel, loggingFile, loggingConfigMap),
				},
			}
		} else {
			return coherence.CoherenceClusterSpec{
				CoherenceRoleSpec: coherence.CoherenceRoleSpec{
					Logging: getCoherenceClusterLoggingSpec(logLevel, loggingFile, loggingConfigMap),
					Images:  getUserArtifactsSpec(loggingFile),
				},
			}
		}
	} else {
		var storage coherence.CoherenceRoleSpec
		if len(loggingFile) == 0 || len(loggingConfigMap) > 0 {
			storage = coherence.CoherenceRoleSpec{
				Role:     roleName,
				Replicas: &replicas,
				Logging:  getCoherenceClusterLoggingSpec(logLevel, loggingFile, loggingConfigMap),
			}
		} else {
			storage = coherence.CoherenceRoleSpec{
				Role:     roleName,
				Replicas: &replicas,
				Logging:  getCoherenceClusterLoggingSpec(logLevel, loggingFile, loggingConfigMap),
				Images:   getUserArtifactsSpec(loggingFile),
			}
		}
		return coherence.CoherenceClusterSpec{
			Roles: []coherence.CoherenceRoleSpec{storage},
		}
	}
}

func getUserArtifactsSpec(loggingFile string) *coherence.Images {
	return &coherence.Images{
		UserArtifacts: &coherence.UserArtifactsImageSpec{
			ImageSpec: coherence.ImageSpec{
				Image: stringPtr(os.Getenv("TEST_USER_IMAGE")),
				//ImagePullPolicy: v1.PullPolicy(os.Getenv("TEST_IMAGE_PULL_POLICY")),
			},
			LibDir:    stringPtr("/files/lib"),
			ConfigDir: stringPtr("/files/conf"),
		},
	}
}

func deleteLoggingConfigMap(namespace string) {
	framework.Global.KubeClient.CoreV1().ConfigMaps(namespace).Delete(loggingConfigMapName,
		&metav1.DeleteOptions{})
}

func createLoggingConfigMap(g *GomegaWithT) {
	f := framework.Global
	namespace := os.Getenv("TEST_NAMESPACE")
	logConfigMap := &corev1.ConfigMap{}
	logConfigMap.SetNamespace(namespace)
	logConfigMap.SetName(loggingConfigMapName)

	logConfigMap.Data = make(map[string]string)
	rootDir, err := helper.FindProjectRootDir()
	g.Expect(err).NotTo(HaveOccurred())

	var cmData []byte
	cmData, err = ioutil.ReadFile(rootDir + "/java/operator-test/src/docker/logging.properties")
	g.Expect(err).NotTo(HaveOccurred())
	logConfigMap.Data["logging.properties"] = string(cmData)

	existingCM, _ := f.KubeClient.CoreV1().ConfigMaps(namespace).Get(loggingConfigMapName,
		metav1.GetOptions{})
	if existingCM != nil {
		deleteLoggingConfigMap(namespace)
	}

	_, err = f.KubeClient.CoreV1().ConfigMaps(namespace).Create(logConfigMap)
	g.Expect(err).NotTo(HaveOccurred())
}

func getCoherenceClusterLoggingSpec(logLevel int32, loggingFile string,
	loggingConfigMap string) *coherence.LoggingSpec {
	if logLevel != 0 {
		return &coherence.LoggingSpec{
			Level: int32Ptr(logLevel),
		}
	} else if loggingConfigMap == string("") {
		return &coherence.LoggingSpec{
			ConfigFile: &loggingFile,
		}
	} else {
		return &coherence.LoggingSpec{
			ConfigFile:    &loggingFile,
			ConfigMapName: &loggingConfigMap,
		}
	}
}

func verifyLogConfigOnMembers(g *GomegaWithT, clusterName string, roleName string,
	replicas int, logLevel int) {
	ok, _ := helper.IsCoherenceVersionAtLeast(12, 2, 1, 4)
	if ok {
		f := framework.Global
		namespace := os.Getenv("TEST_NAMESPACE")

		// Get the list of Pods
		pods, err := helper.ListCoherencePodsForRole(f.KubeClient, namespace, clusterName, roleName)
		g.Expect(err).NotTo(HaveOccurred())

		// Port forward to the first Pod
		pf, ports, err := helper.StartPortForwarderForPod(&pods[0])
		g.Expect(err).NotTo(HaveOccurred())
		defer pf.Close()

		cl := &http.Client{}
		members, _, err := management.GetMembers(cl, "127.0.0.1", ports[management.PortName])
		g.Expect(err).NotTo(HaveOccurred())
		g.Expect(len(members.Items)).To(Equal(replicas))

		// assert that each member has required logging configuration.
		for _, member := range members.Items {
			g.Expect(member.MachineName).NotTo(BeEmpty())
			g.Expect(member.LoggingLevel).To(Equal(logLevel))
		}
	}
}
