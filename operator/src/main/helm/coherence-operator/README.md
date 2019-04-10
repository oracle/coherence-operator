# coherence-operator
Install coherence-operator to work with Coherence clusters on Kubernetes.

## Introduction

This chart install a coherence-operator deployment on a 
[Kubernetes](https://kubernetes.io) cluster using the [Helm](https://helm.sh)
package manager.

## Prerequisites
* Kubernetes 1.10.3 or above
* Helm 2.11.0 or above

## Installing the Chart
To install the chart with the release name `sample-coherence-operator`:

```
@ helm install --name sample-coherence-operator coherence-operator
```

The command deploys coherence-operator on the Kubernetes cluster in the
default configuration. The [configuration](#configuration) section list
parameters that can be configured during installation.

## Uninstalling the Chart
To uninstall the `sample-coherence-operator` deployment:

```
$ helm delete sample-coherence-operator
```

The command removes all the Kubernetes components associated with the chart
and deletes the release.

We also need to remove the internal config map in targetNamespaces, which is
`{ default }` by default.

```
$ kubectl delete configmap coherence-internal-config
```

## Configuration

The following table list the configurable parameters of the coherence-operator
chart and their default values.

| Parameter | Description | Default |
| --------- | ----------- | ------- |
| `targetNamespaces` | A list of target namespaces the operator manages | `["default"]` |
| `serviceAccount` | Service account used to access Kubernetes API | `default` |
| `imagePullSecrets` | Secret for pull images from private registries |  |
| `imagePullSecretsSeparator` | Separator for secret for pull images from private registries | `$` |
| `service.name` | Name of the service | `coherence-operator-service` |
| `service.type` | Kubernetes Service Type | `"ClusterIP"`|
| `service.domain` | External domain name | `"cluster.local"` |
| `service.loadBalancerIP` | IP address of the load balancer | |
| `service.annotations` | Service annotations yaml | |
| `coherenceOperator.image` | Coherence Operator image to be pulled | `"oracle/coherence-operator:1.0.0-SNAPSHOT"` |
| `coherenceOperator.imagePullPolicy` | Image pull policy | `"IfNotPresent"` |
| `javaLoggingLevel` | Java logging level | `"INFO"` |
| `logCaptureEnabled` | Whether log capture via EFK stack is enabled | `false` |
| `elasticsearch.image` | Elasticsearch Docker image url with tag | `docker.elastic.co/elasticsearch/elasticsearch-oss:6.6.0` |
| `elasticsearch.imagePullPolicy` | Elasticsearch image pull policy | `"IfNotPresent"` |
| `elasticsearchEndpoint.host` | Elasticsearch host installed separately | `"elasticsearch.${namespace}.svc.cluster.local` |
| `elasticsearchEndpoint.port` | Elasticsearch port intalled separately | `9200` |
| `logstash.image` | Logstash Docker image url with tag | `docker.elastic.co/logstash/logstash-oss:6.6.0` |
| `logstash.imagePullPolicy` | Logstash image pull policy | `"IfNotPresent"` |
| `kibana.image` | Kibana Docker image url with tag | `docker.elastic.co/beats/filebeat:6.2.4` |
| `kibana.imagePullPolicy` | Kibana image pull policy | `"IfNotPresent"` |
| `filebeat.image` | FileBeat Docker image url with tag | `docker.elastic.co/beats/filebeat:6.2.4` |
| `filebeat.imagePullPolicy` | Filebeat image pull policy | `"IfNotPresent"` |
| `prometheusoperator.enabled` | Whether Prometheus is enabled | `false` |
| `prometheusoperator.grafana.enabled` | Whether Grafana is enabled | `false` |