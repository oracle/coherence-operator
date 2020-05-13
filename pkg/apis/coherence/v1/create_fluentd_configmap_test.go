/*
 * Copyright (c) 2020, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1_test

import (
	"fmt"
	. "github.com/onsi/gomega"
	coh "github.com/oracle/coherence-operator/pkg/apis/coherence/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
	"testing"
)

func TestCreateFluentdConfigMapForNilLoggingSpec(t *testing.T) {
	g := NewGomegaWithT(t)

	spec := coh.CoherenceDeploymentSpec{}
	deployment := createTestDeployment(spec)

	res, err := deployment.Spec.Logging.CreateFluentdConfigMap(deployment)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(res).To(BeNil())
}

func TestCreateFluentdConfigMapForNilFluentdSpec(t *testing.T) {
	g := NewGomegaWithT(t)

	spec := coh.CoherenceDeploymentSpec{
		Logging: &coh.LoggingSpec{},
	}
	deployment := createTestDeployment(spec)

	res, err := deployment.Spec.Logging.CreateFluentdConfigMap(deployment)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(res).To(BeNil())
}

func TestCreateFluentdConfigMapForEmptyFluentdSpec(t *testing.T) {
	g := NewGomegaWithT(t)

	spec := coh.CoherenceDeploymentSpec{
		Logging: &coh.LoggingSpec{
			Fluentd: &coh.FluentdSpec{},
		},
	}
	deployment := createTestDeployment(spec)

	res, err := deployment.Spec.Logging.CreateFluentdConfigMap(deployment)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(res).To(BeNil())
}

func TestCreateFluentdConfigMapWithFluentdEnabledFalse(t *testing.T) {
	g := NewGomegaWithT(t)

	policy := corev1.PullIfNotPresent
	spec := coh.CoherenceDeploymentSpec{
		Logging: &coh.LoggingSpec{
			Fluentd: &coh.FluentdSpec{
				Enabled: pointer.BoolPtr(false),
				ImageSpec: coh.ImageSpec{
					Image:           pointer.StringPtr("foo:1.0"),
					ImagePullPolicy: &policy,
				},
				ConfigFileInclude: pointer.StringPtr("foo"),
				Tag:               pointer.StringPtr("bar"),
			},
		},
	}
	deployment := createTestDeployment(spec)

	res, err := deployment.Spec.Logging.CreateFluentdConfigMap(deployment)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(res).To(BeNil())
}

func TestCreateFluentdConfigMapWithFluentdEnabledTrue(t *testing.T) {
	clusterName := "test-cluster"
	deploymentName := "test"
	roleName := "fluentd-test"
	expectedConfig := `# Coherence fluentd configuration

# Ignore fluentd messages
<match fluent.**>
  @type null
</match>

# Coherence Logs
<source>
  @type tail
  path /logs/coherence-*.log
  pos_file /tmp/cohrence.log.pos
  read_from_head true
  tag coherence-cluster
  multiline_flush_interval 20s
  <parse>
    @type multiline
    format_firstline /^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}.\d{3}/
    format1 /^(?<time>\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}.\d{3})\/(?<uptime>[0-9\.]+) (?<product>.+) <(?<level>[^\s]+)> \(thread=(?<thread>.+), member=(?<member>.+)\):[\S\s](?<log>.*)/
  </parse>
</source>

<filter coherence-cluster>
  @type record_transformer
  <record>
    cluster "test-cluster"
    deployment "test"
    role "fluentd-test"
    host "#{ENV['HOSTNAME']}"
    pod-uid "#{ENV['COHERENCE_POD_ID']}"
  </record>
</filter>

<match coherence-cluster>
  @type elasticsearch
  hosts "#{ENV['ELASTICSEARCH_HOSTS']}"
  user "#{ENV['ELASTICSEARCH_USER']}"
  password "#{ENV['ELASTICSEARCH_PASSWORD']}"
  logstash_format true
  logstash_prefix coherence-cluster
</match>`

	g := NewGomegaWithT(t)

	spec := coh.CoherenceDeploymentSpec{
		Role:        roleName,
		ClusterName: pointer.StringPtr(clusterName),
		Logging: &coh.LoggingSpec{
			Fluentd: &coh.FluentdSpec{
				Enabled: pointer.BoolPtr(true),
			},
		},
	}
	deployment := createTestDeployment(spec)
	deployment.Name = deploymentName

	res, err := deployment.Spec.Logging.CreateFluentdConfigMap(deployment)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(res).NotTo(BeNil())
	g.Expect(res.Name).To(Equal(fmt.Sprintf(coh.EfkConfigMapNameTemplate, deployment.GetName())))
	g.Expect(res.IsDelete()).To(BeFalse())

	cm := res.Spec.(*corev1.ConfigMap)
	g.Expect(len(cm.Data)).To(Equal(1))
	actualConfig, found := cm.Data[coh.VolumeMountSubPathFluentdConfig]
	g.Expect(found).To(BeTrue())
	g.Expect(actualConfig).To(Equal(expectedConfig))

	expected := CreateExpectedFluentdConfigMap(deployment)
	expected.Data[coh.VolumeMountSubPathFluentdConfig] = expectedConfig
	g.Expect(res.Spec).To(Equal(expected))
}

func TestCreateFluentdConfigMapWithFluentdConfigInclude(t *testing.T) {
	clusterName := "test-cluster"
	deploymentName := "test"
	roleName := "fluentd-test"
	include := "test-include.conf"
	expectedConfig := `# Coherence fluentd configuration
@include test-include.conf

# Ignore fluentd messages
<match fluent.**>
  @type null
</match>

# Coherence Logs
<source>
  @type tail
  path /logs/coherence-*.log
  pos_file /tmp/cohrence.log.pos
  read_from_head true
  tag coherence-cluster
  multiline_flush_interval 20s
  <parse>
    @type multiline
    format_firstline /^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}.\d{3}/
    format1 /^(?<time>\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}.\d{3})\/(?<uptime>[0-9\.]+) (?<product>.+) <(?<level>[^\s]+)> \(thread=(?<thread>.+), member=(?<member>.+)\):[\S\s](?<log>.*)/
  </parse>
</source>

<filter coherence-cluster>
  @type record_transformer
  <record>
    cluster "test-cluster"
    deployment "test"
    role "fluentd-test"
    host "#{ENV['HOSTNAME']}"
    pod-uid "#{ENV['COHERENCE_POD_ID']}"
  </record>
</filter>

<match coherence-cluster>
  @type elasticsearch
  hosts "#{ENV['ELASTICSEARCH_HOSTS']}"
  user "#{ENV['ELASTICSEARCH_USER']}"
  password "#{ENV['ELASTICSEARCH_PASSWORD']}"
  logstash_format true
  logstash_prefix coherence-cluster
</match>`

	g := NewGomegaWithT(t)

	spec := coh.CoherenceDeploymentSpec{
		Role:        roleName,
		ClusterName: pointer.StringPtr(clusterName),
		Logging: &coh.LoggingSpec{
			Fluentd: &coh.FluentdSpec{
				Enabled:           pointer.BoolPtr(true),
				ConfigFileInclude: &include,
			},
		},
	}
	deployment := createTestDeployment(spec)
	deployment.Name = deploymentName

	res, err := deployment.Spec.Logging.CreateFluentdConfigMap(deployment)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(res).NotTo(BeNil())
	g.Expect(res.Name).To(Equal(fmt.Sprintf(coh.EfkConfigMapNameTemplate, deployment.GetName())))
	g.Expect(res.IsDelete()).To(BeFalse())

	cm := res.Spec.(*corev1.ConfigMap)
	g.Expect(len(cm.Data)).To(Equal(1))
	actualConfig, found := cm.Data[coh.VolumeMountSubPathFluentdConfig]
	g.Expect(found).To(BeTrue())
	g.Expect(actualConfig).To(Equal(expectedConfig))

	expected := CreateExpectedFluentdConfigMap(deployment)
	expected.Data[coh.VolumeMountSubPathFluentdConfig] = expectedConfig
	g.Expect(res.Spec).To(Equal(expected))
}

func TestCreateFluentdConfigMapWithFluentdTag(t *testing.T) {
	clusterName := "test-cluster"
	deploymentName := "test"
	roleName := "fluentd-test"
	tag := "foo"
	expectedConfig := `# Coherence fluentd configuration

# Ignore fluentd messages
<match fluent.**>
  @type null
</match>

# Coherence Logs
<source>
  @type tail
  path /logs/coherence-*.log
  pos_file /tmp/cohrence.log.pos
  read_from_head true
  tag coherence-cluster
  multiline_flush_interval 20s
  <parse>
    @type multiline
    format_firstline /^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}.\d{3}/
    format1 /^(?<time>\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}.\d{3})\/(?<uptime>[0-9\.]+) (?<product>.+) <(?<level>[^\s]+)> \(thread=(?<thread>.+), member=(?<member>.+)\):[\S\s](?<log>.*)/
  </parse>
</source>

<filter coherence-cluster>
  @type record_transformer
  <record>
    cluster "test-cluster"
    deployment "test"
    role "fluentd-test"
    host "#{ENV['HOSTNAME']}"
    pod-uid "#{ENV['COHERENCE_POD_ID']}"
  </record>
</filter>

<match coherence-cluster>
  @type elasticsearch
  hosts "#{ENV['ELASTICSEARCH_HOSTS']}"
  user "#{ENV['ELASTICSEARCH_USER']}"
  password "#{ENV['ELASTICSEARCH_PASSWORD']}"
  logstash_format true
  logstash_prefix coherence-cluster
</match>
<match foo >
  @type elasticsearch
  hosts "#{ENV['ELASTICSEARCH_HOSTS']}"
  user "#{ENV['ELASTICSEARCH_USER']}"
  password "#{ENV['ELASTICSEARCH_PASSWORD']}"
  logstash_format true
  logstash_prefix foo
</match>`

	g := NewGomegaWithT(t)

	spec := coh.CoherenceDeploymentSpec{
		Role:        roleName,
		ClusterName: pointer.StringPtr(clusterName),
		Logging: &coh.LoggingSpec{
			Fluentd: &coh.FluentdSpec{
				Enabled: pointer.BoolPtr(true),
				Tag:     &tag,
			},
		},
	}
	deployment := createTestDeployment(spec)
	deployment.Name = deploymentName

	res, err := deployment.Spec.Logging.CreateFluentdConfigMap(deployment)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(res).NotTo(BeNil())
	g.Expect(res.Name).To(Equal(fmt.Sprintf(coh.EfkConfigMapNameTemplate, deployment.GetName())))
	g.Expect(res.IsDelete()).To(BeFalse())

	cm := res.Spec.(*corev1.ConfigMap)
	g.Expect(len(cm.Data)).To(Equal(1))
	actualConfig, found := cm.Data[coh.VolumeMountSubPathFluentdConfig]
	g.Expect(found).To(BeTrue())
	g.Expect(actualConfig).To(Equal(expectedConfig))

	expected := CreateExpectedFluentdConfigMap(deployment)
	expected.Data[coh.VolumeMountSubPathFluentdConfig] = expectedConfig
	g.Expect(res.Spec).To(Equal(expected))
}

func CreateExpectedFluentdConfigMap(deployment *coh.CoherenceDeployment) *corev1.ConfigMap {
	labels := deployment.CreateCommonLabels()
	labels[coh.LabelComponent] = coh.LabelComponentEfkConfig

	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: deployment.Namespace,
			Name:      fmt.Sprintf(coh.EfkConfigMapNameTemplate, deployment.GetName()),
			Labels:    labels,
		},
		Data: map[string]string{},
	}
}
