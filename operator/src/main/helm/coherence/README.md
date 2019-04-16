# coherence
Install coherence on Kubernetes.

## Introduction

This chart install a Coherence deployment/statefulset on a 
[Kubernetes](https://kubernetes.io) cluster using the [Helm](https://helm.sh)
package manager.

## Prerequisites
* Kubernetes 1.10.3 or above
* Helm 2.11.0 or above

## Installing the Chart
To install the chart with the release name `sample-coherence`:

```
@ helm install --name sample-coherence coherence
```

The command deploys coherence on the Kubernetes cluster in the
default configuration. The [configuration](#configuration) section list
parameters that can be configured during installation.

## Uninstalling the Chart
To uninstall the `sample-coherence` deployment:

```
$ helm delete sample-coherence
```

The command removes all the Kubernetes components associated with the chart
and deletes the release.

## Configuration

The following table list the configurable parameters of the coherence
chart and their default values.

Note: Coherence management over REST and metrics publishing will be available in 
Coherence 12.2.1.4. 

| Parameter | Description | Default |
| --------- | ----------- | ------- |
| `clusterSize` | Initial size of the Coherence cluster | `3` |
| `imagePullSecrets` | Secret for pull images from private registries |  |
| `imagePullSecretsSeparator` | Separator for secret for pull images from private registries | `$` |
| `role` | Role name of Coherence cluster member |  |
| `cluster` | Cluster name of Coherence cluster member | Helm release name |
| `serviceAccountName` | Name of the service | `"default"` |
| `store.cacheConfig` | Name of the cache configuration file | |
| `store.pof.config` | Name of POF serialization configuration file | `pof-config.xml` |
| `store.storageEnabled` | Whether Coherence storage is enabled | |
| `store.logging.level` | Coherence log level | `5` |
| `store.logging.configFile` | The location of the java util logging configuration file | `logging.properties` embedded in this chart |
| `store.logging.configMapName` | The config map to be mounted as a volume containing the logging configuration file | |
| `store.maxHeap` | min/max heap value pass to JVM | |
| `store.jvmArgs` | Options pass to JVM | `"-XX:+UseG1GC"` |
| `store.javaOpts` | Miscellaneous JVM options overriding the computed startup JVM options | |
| `store.wkaRelease` | Used to override WKA address| Default WKA address |
| `store.wka` | Used to oveerride the default WKA address| The address of the headless service |
| `store.ports` | Additional port mappings (in Yaml) | |
| `store.env` | Additional environment variable mappings (in Yaml) | |
| `store.annotations` | Annotations (in Yaml) | |
| `store.management.ssl.enabled` | Whether SSL is enabled for Coherence management over REST endpoint | `false` |
| `store.management.ssl.secrets` | Name of the k8s secrets containing the Java key stores and password files | |
| `store.management.ssl.keyStore` | Name of the file in the k8s secret containing the keystore | |
| `store.management.ssl.keyStorePasswordFile` | Name of the file in the k8s secret containing the keystore password | |
| `store.management.ssl.keyPasswordFile` | Name of the file in the k8s secret containing the key password | |
| `store.management.ssl.keyStoreType` | Name of the Java keystore type for the keystore | `JKS` |
| `store.management.ssl.trustStore` | Name of the file in the k8s secret containing the trust store | |
| `store.management.ssl.trustStorePasswordFile` | Name of the file in the k8s secret containing the trust store password | `false` |
| `store.management.ssl.trustStoreType` | Name of the Java keystore type for the trust store | `JKS` |
| `store.management.ssl.requireClientCert` | Whether the client certificate will be authenticated by the server | `false` |
| `store.metrics.ssl.enabled` | Whether SSL is enabled for Coherence metrics endpoint | `false` |
| `store.metrics.ssl.secrets` | Name of the k8s secrets containing the Java key stores and password files | |
| `store.metrics.ssl.keyStore` | Name of the file in the k8s secret containing the keystore | |
| `store.metrics.ssl.keyStorePasswordFile` | Name of the file in the k8s secret containing the keystore password | |
| `store.metrics.ssl.keyPasswordFile` | Name of the file in the k8s secret containing the key password | |
| `store.metrics.ssl.keyStoreType` | Name of the Java keystore type for the keystore | `JKS` |
| `store.metrics.ssl.trustStore` | Name of the file in the k8s secret containing the trust store | |
| `store.metrics.ssl.trustStorePasswordFile` | Name of the file in the k8s secret containing the trust store password | `false` |
| `store.metrics.ssl.trustStoreType` | Name of the Java keystore type for the trust store | `JKS` |
| `store.metrics.ssl.requireClientCert` | Whether the client certificate will be authenticated by the server | `false` |
| `store.persistence.enabled` | Whether persistence on-disc is enabled | `false` |
| `store.persistence.size` | The size of the PersistentVolume to allocate to each Pod for persistence | `2Gi` |
| `store.persistence.storageClass` | the Persistent Volume Storage Class for persistence | |
| `store.persistence.datasource` | the dataSource attribute to use in the persistent PVC | |
| `store.persistence.volumeMode` | the volumeMode attribute to use in the persistent PVC | |
| `store.persistence.volumeName` | the volumeName attribute to use in the persistent PVC | |
| `store.persistence.selector` | YAML as selector section of the persistent PVC | |
| `store.persistence.volume` | Allows the configuration of a normal k8s volume mapping for persistence data instead of a PVC | |
| `store.podManagementPolicy` | Management Policy for Coherence members' pod | `Parallel` |
| `store.snapshot.enabled` | Enables or disabled a different location for persistence snapshot data | `false` |
| `store.snapshot.size` | The size of the PersistentVolume to allocate to each Pod for snapshot | `2Gi` |
| `store.snapshot.storageClass` | The Persistent Volume Storage Class for snapshot | |
| `store.snapshot.datasource` | the dataSource attribute to use in the snapshot PVC | |
| `store.snapshot.volumeMode` | the volumeMode attribute to use in the snapshotPVC | |
| `store.snapshot.volumeName` | the volumeName attribute to use in the snapshot PVC | |
| `store.snapshot.selector` | YAML as selector section of the snapshot PVC | |
| `store.snapshot.volume` | Allows the configuration of a normal k8s volume mapping for persistence snapshot data instead of a PVC. | |
| `store.readinessProbe.initialDelaySeconds` | Number of seconds after the container has started before liveness or readiness probes are initiated | 30 |
| `store.readinessProbe.periodSeconds` | How often (in seconds) to perform the probe. Minimum value is 1. | 60 |
| `store.readinessProbe.timeoutSeconds` | Number of seconds after which the probe times out. Defaults to 1 second. Minimum value is 1. | 5 |
| `store.readinessProbe.successThreshold` | Minimum consecutive successes for the probe to be considered successful after having failed. | Defaults to 1. Must be 1 for liveness. Minimum value is 1.|
| `store.readinessProbe.failureThreshold` | Kubernetes will try failureThreshold times before giving up. | 50 |
| `affinity` | Affinity that controls Pod scheduling preferences | `{}`|
| `nodeSelector` | Node lables for pod assignment | `{}` |
| `tolerations` | For nodes that have taints on them. See (https://kubernetes.io/docs/concepts/configuration/taint-and-toleration) | `[]` |
| `service.enabled` | Whether to create the service yaml or not| `true` |
| `service.type` |  The Kubernetes service type (typically ClusterIP or LoadBalancer) | `"ClusterIP"`|
| `service.domain` | The external domain name | `cluster.local` |
| `service.loadBalancerIP` | The IP address of the load balancer| |
| `service.annotations` | Service annotations (in Yaml) | |
| `service.externalPort` | The Extend service port| `20000`|
| `service.managementPort` | The management http port | `30000` |
| `service.metricsHttpPort` | The metrics http port | `9095` |
| `resources.requests.cpu` | See (http://kubernetes.io/docs/user-guide/compute-resources) | `0` |
| `resources.limits.cpu` | See (http://kubernetes.io/docs/user-guide/compute-resources) | `32`|
| `coherence.image` | Coherence image to be pulled | `"oracle/coherence:12.2.1.3"` |
| `coherence.imagepullPolicy` | Coherence Image pull policy | `"IfNotPresent"` |
| `coherenceUtils.image` | Coherence Utils image to be pulled | `"oracle/coherence-utils:1.0.0"` |
| `coherenceUtils.imagePullPolicy` | Coherence Utils Image pull policy | `"IfNotPresent"` |
| `logCaptureEnabled` | Whether log capture via EFK stack is enabled | `false` |
| `fluentd.image` | Fluentd image to be pulled | `"fluent/fluentd-kubernetes-daemonset:v1.3.3-debian-elasticsearch-1.3"' |
| `fluentd.imagepullPolicy` | Fluentd Image pull policy | `"IfNotPresent"` |
| `fluentd.application.configFile` | The location of fluentd application configuration file containing application source.  | |
| `fluentd.application.tag` | The fluentd tag for fluentd source specified in fluentd.application.configFile file |  |
| `logstash.image` | Logstash Docker image url with tag | `docker.elastic.co/logstash/logstash-oss:6.6.0` |
| `logstash.imagePullPolicy` | Logstash image pull policy | `"IfNotPresent"` |
| `filebeat.image` | FileBeat Docker image url with tag | `docker.elastic.co/beats/filebeat:6.2.4` |
| `filebeat.imagePullPolicy` | Filebeat image pull policy | `"IfNotPresent"` |
| `userArtifacts.image` | The name of the Docker image containing the custom jar files and configuration files to add to the classpath|| 
| `userArtifacts.imagePullPolicy` | Image pull policy of the user artifacts| `"IfNotPresent"`|
| `userArtifacts.libDir` | The folder in the custom artifacts Docker image containing jar files to be added to the classpath of the Coherence container | `"/files/lib"`|
| `userArtifacts.configDir` | the folder in the custom artifacts Docker image containing configuration files to be added to the classpath of the Coherence container | `"/files/conf"`|

## Current Limitations

* In case if cluster needs to be scaled down for any reason, care needs to be taken while performing scale down or reducing the number of replica count of stateful sets. Make sure only one pod brought down or terminated at any given point of time and make sure cluster has reached 'statusHA' before further scale down operation on a same cluster. This must be followed *strickly* in order NOT to loose any partition or data when scale down the number of cluster members.
