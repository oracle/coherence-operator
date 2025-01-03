/*
 * Copyright (c) 2021, 2025, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package clients

import (
	v1 "github.com/oracle/coherence-operator/api/v1"
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
	Coherence      v1.Coherence
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
		Coherence:      cluster,
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
		{
			Name:  "COHERENCE_CLUSTER",
			Value: tc.Cluster.Coherence.GetCoherenceClusterName(),
		},
	}

	switch tc.ClientType {
	case ClientTypeExtend:
		var extendAddress string
		var extendPort string

		if tc.Name == "ExtendInternalNS" {
			// Use name service
			extendAddress = tc.Cluster.ServiceFQDNs["wka"]
			extendPort = "7574"
		} else {
			// use direct socket address
			extendAddress = tc.Cluster.ServiceFQDNs["extend"]
			extendPort = "20000"
		}

		envVars = append(envVars, corev1.EnvVar{
			Name:  "COHERENCE_EXTEND_ADDRESS",
			Value: extendAddress,
		})
		envVars = append(envVars, corev1.EnvVar{
			Name:  "COHERENCE_EXTEND_PORT",
			Value: extendPort,
		})
	case ClientTypeGrpc:
		grpcAddress := tc.Cluster.ServiceFQDNs["grpc"]
		envVars = append(envVars, corev1.EnvVar{
			Name:  "COHERENCE_GRPC_ADDRESS",
			Value: grpcAddress,
		})
		envVars = append(envVars, corev1.EnvVar{
			Name:  "COHERENCE_GRPC_PORT",
			Value: "1408",
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
