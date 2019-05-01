# Push logs to your own Elasticsearch Instance  

The Oracle Coherence Operator manages logging data through the EFK
(ElasticSearch, Fluentd and Kibana) stack.

the sample explains how to instruct fluentd to your own elastic search instance.

[Return to Logging samples](../) / [Return to Coherence Operator samples](../../) / [Return to samples](../../../README.md#list-of-samples)

This will ensure that Elasticsearch and Kibana will be installed and configured.

## Installation Steps

1. Install the Coherence Operator

   You must install the Coherence Operator by using the instructions below to 
   enable logCapture and point to the host and port of your `existing elasticsearch` instance.

   *Note*: If you have an existing running Coherence Operator, you should uninstall 
   by using `helm delete coherence-operator --purge`.

   Ensure you set the following:
  
   * `elasticsearchEndpoint.host` to your elasticsearch host
   
   * `elasticsearchEndpoint.port` to your elasticsearch port 
   
   *Note*: If your Elasticsearch engine requires user and password, you can also set the following:
   
   * `elasticsearchEndpoint.user` to your elasticsearch username
   
   * `elasticsearchEndpoint.password` to your elasticsearch password 
  
   ```bash
   $ helm install \
      --namespace sample-coherence-ns \
      --set imagePullSecrets=sample-coherence-secret \
      --name coherence-operator \
      --set logCaptureEnabled=true \
      --set elasticsearchEndpoint.host=my-elastic-host \
      --set elasticsearchEndpoint.port=my-elastic-port \
      --set "targetNamespaces={sample-coherence-ns}" \
      --version 1.0.0-SNAPSHOT coherence-community/coherence-operator
   ```

1. Confirm the elasticdata endpoint is set

   ```bash
   $ kubectl get pods -n sample-coherence-ns
   NAME                                  READY   STATUS    RESTARTS   AGE
   coherence-operator-856d5f8544-kgxgd   2/2     Running   0          8m

   $ kubectl logs coherence-operator-856d5f8544-kgxgd  -n sample-coherence-ns -c fluentd | grep -A3 'match coherence-operator'
   
   <match coherence-operator>
    @type elasticsearch
    host "my-elastic-host"
    port 9200
   ```
   
   The above host and port should match the values you supplied in the above `helm install`.
                   
1. Install the Coherence cluster

   The following additional options are set:
   
   * `--set logCaptureEnabled=true` - this will then use the configuration of the operator 
     for the elasticsearch endpoint for fluentd.

   *Note*: The Coherence Operator will provide the Elasticsearch host and port values to the Coherence install.
   
   ```bash
   $ helm install \
      --namespace sample-coherence-ns \
      --name storage \
      --set clusterSize=3 \
      --set cluster=storage-tier-cluster \
      --set imagePullSecrets=sample-coherence-secret \
      --set logCaptureEnabled=true \
      --version 1.0.0-SNAPSHOT coherence-community/coherence
   ```
   
   Once the install has completed issue the following command to list the pods:

   ```bash
   $ kubectl get pods -n sample-coherence-ns
   NAME                                  READY   STATUS    RESTARTS   AGE
   coherence-operator-7f596c6796-nns9n   2/2     Running   0          22m
   storage-coherence-0                   2/2     Running   0          17m
   storage-coherence-1                   2/2     Running   0          16m
   storage-coherence-2                   2/2     Running   0          16m
   ```
   
   Notice that the `coherence-operator` and all the `coherence` pods have two containers.
   
   ```bash
   $ kubectl logs storage-coherence-0  -n sample-coherence-ns -c fluentd | grep -A3 'match coherence-cluster'
    <match coherence-cluster>
     @type elasticsearch
     host "my-elastic-host"
     port 9200
   ```
   
   The above host and port should match the values you supplied in the `coherence-operator` install.
   
## Uninstalling the Charts

Carry out the following commands to delete the chart installed in this sample.

```bash
$ helm delete storage --purge
```

Before starting another sample, ensure that all the pods are gone from previous sample.

If you wish to remove the `coherence-operator`, then include it in the `helm delete` command above.

