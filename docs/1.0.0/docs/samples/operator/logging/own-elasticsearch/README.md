# Push Logs to Your Elasticsearch Instance  

The Coherence Operator manages data logging through the
ElasticSearch, Fluentd, and Kibana (EFK) stack.

This sample explains how to make Fluentd to push data to your own Elasticsearch instance.

[Return to Logging samples](../) / [Return to Coherence Operator samples](../../) / [Return to samples](../../../README.md#list-of-samples)

This will ensure that Elasticsearch and Kibana will be installed and configured.

## Installation Steps

1. Install the Coherence Operator

   You must install the Coherence Operator using the instructions that enable log capture and point to the host and port of your existing Elasticsearch instance.

   *Note*: If you have a running Coherence Operator, you should uninstall using the command `helm delete coherence-operator --purge`.

   Ensure you set the following:
  
   * `elasticsearchEndpoint.host` to your Elasticsearch host.
   
   * `elasticsearchEndpoint.port` to your Elasticsearch port.
   
   *Note*: If your Elasticsearch host requires username and password, set the following:
   
   * `elasticsearchEndpoint.user` to your Elasticsearch username.
   
   * `elasticsearchEndpoint.password` to your elasticsearch password.
  
   ```bash
   $ helm install \
      --namespace sample-coherence-ns \
      --set imagePullSecrets=sample-coherence-secret \
      --name coherence-operator \
      --set logCaptureEnabled=true \
      --set elasticsearchEndpoint.host=my-elastic-host \
      --set elasticsearchEndpoint.port=my-elastic-port \
      --set "targetNamespaces={sample-coherence-ns}" \
      coherence/coherence-operator
   ```

1. Verify that the Elasticsearch endpoint is set.

   ```bash
   $ kubectl get pods -n sample-coherence-ns
   ```
   ```console
   NAME                                  READY   STATUS    RESTARTS   AGE
   coherence-operator-856d5f8544-kgxgd   2/2     Running   0          8m
   ```
   ```bash
   $ kubectl logs coherence-operator-856d5f8544-kgxgd  -n sample-coherence-ns -c fluentd | grep -A3 'match coherence-operator'
   ```
   ```console
   
   <match coherence-operator>
    @type elasticsearch
    host "my-elastic-host"
    port my-elastic-port
   ```
   
   The host and port value must match the values you supplied in the `helm install` command.
                   
1. Install the Coherence cluster

   The following additional options are set:
   
   * `--set logCaptureEnabled=true` - This uses the configuration of the operator for the Elasticsearch endpoint for Fluentd.

   > **Note**: The Coherence Operator provides the Elasticsearch host and port values to install Coherence.
   
   ```bash
   $ helm install \
      --namespace sample-coherence-ns \
      --name storage \
      --set clusterSize=3 \
      --set cluster=storage-tier-cluster \
      --set imagePullSecrets=sample-coherence-secret \
      --set logCaptureEnabled=true \
      coherence/coherence
   ```
   
   After the installation completes, list the pods:

   ```bash
   $ kubectl get pods -n sample-coherence-ns
   ```
   ```console
   NAME                                  READY   STATUS    RESTARTS   AGE
   coherence-operator-7f596c6796-nns9n   2/2     Running   0          22m
   storage-coherence-0                   2/2     Running   0          17m
   storage-coherence-1                   2/2     Running   0          16m
   storage-coherence-2                   2/2     Running   0          16m
   ```
   
   The `coherence-operator` and all the `coherence` pods have two containers. The second container is for Fluentd.
   
   ```bash
   $ kubectl logs storage-coherence-0  -n sample-coherence-ns -c fluentd | grep -A3 'match coherence-cluster'
    <match coherence-cluster>
    ```
    ```console
     @type elasticsearch
     host "my-elastic-host"
     port my-elastic-port
   ```
   
   The host and port values must match the values you provided to install the `coherence-operator`.
   
## Uninstall the Charts

Use the following commands to delete the chart installed in this sample:

```bash
$ helm delete storage --purge
```

Before starting another sample, ensure that all the pods are removed from previous sample.

If you want to remove the `coherence-operator`, then use the `helm delete` command.
