/*
 * Copyright (c) 2021, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package clients

import (
	"github.com/oracle/coherence-operator/test/e2e/helper"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

// CoherenceCluster represents a running test Coherence cluster
type CoherenceCluster struct {
	// Manifest is the name of the yaml file used to deploy the cluster
	Manifest       string
	Services       map[string]string
	ServiceFQDNs   map[string]string
	ServiceIngress map[string][]corev1.LoadBalancerIngress
}

func DeployTestCluster(testContext helper.TestContext, t *testing.T, yaml string) (*CoherenceCluster, error) {
	// Start the Coherence cluster
	cluster, err := helper.AssertSingleDeployment(testContext, t, yaml)
	if err != nil {
		return nil, err
	}

	ns := helper.GetTestNamespace()

	svcNames := cluster.FindPortServiceNames()
	svcFQDN := cluster.FindFullyQualifiedPortServiceNames()

	ingress := make(map[string][]corev1.LoadBalancerIngress)
	for _, n := range svcNames {
		extendSvc, err := testContext.KubeClient.CoreV1().Services(ns).Get(testContext.Context, n, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		ingress[n] = extendSvc.Status.LoadBalancer.Ingress
	}

	return &CoherenceCluster{
		Manifest:       yaml,
		Services:       svcNames,
		ServiceFQDNs:   svcFQDN,
		ServiceIngress: ingress,
	}, nil
}

// CreateClientJob creates a k8s Job based on the test case.
func CreateClientJob(ns, name string, tc ClientTestCase) *batchv1.Job {
	image := helper.GetClientImage()

	cfg := tc.CacheConfig
	if cfg == "" {
		cfg = "test-cache-config.xml"
	}

	envVars := []corev1.EnvVar{
		{
			Name:  "CLIENT_TYPE",
			Value: string(tc.ClientType),
		},
		{
			Name:  "COHERENCE_CACHECONFIG",
			Value: cfg,
		},
		{
			Name:  "COHERENCE_DISTRIBUTED_LOCALSTORAGE",
			Value: "false",
		},
	}

	switch tc.ClientType {
	case ClientTypeExtend:
		envVars = append(envVars, corev1.EnvVar{
			Name:  "COHERENCE_EXTEND_ADDRESS",
			Value: tc.Cluster.ServiceFQDNs["extend"],
		})
	case ClientTypeGrpc:
		envVars = append(envVars, corev1.EnvVar{
			Name:  "COHERENCE_GRPC_CHANNELS_DEFAULT_HOST",
			Value: tc.Cluster.ServiceFQDNs["grpc"],
		})
	}

	job := batchv1.Job{}
	job.SetNamespace(ns)
	job.SetName(name)
	job.Spec = batchv1.JobSpec{
		Template: corev1.PodTemplateSpec{
			Spec: corev1.PodSpec{
				RestartPolicy: corev1.RestartPolicyNever,
				Containers: []corev1.Container{
					{
						Name:  "test",
						Image: image,
						Env:   envVars,
					},
				},
			},
		},
	}

	return &job
}
