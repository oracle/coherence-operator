# Access Coherence Console and CohQL on a Cluster Node

The Coherence Console and CohQL (Coherence Query Language) are developer tools used for interacting with a Coherence cluster.

This samples shows how to use `kubectl exec` to connect to any of the pods and start Coherence console or CohQL as a storage-disabled client.

[Return to Management samples](../) / [Return to samples](../../README.md#list-of-samples)

## Prerequisites

Ensure you have already installed the Coherence Operator using the instructions [here](../../../README.md#install-the-coherence-operator).

## Installation Steps

1. Install the Coherence cluster

  Install the cluster with Persistence and Snapshot enabled:

   ```bash
   $ helm install \
      --namespace sample-coherence-ns \
      --name storage \
      --set clusterSize=3 \
      --set cluster=coherence-cluster \
      --set imagePullSecrets=sample-coherence-secret \
      --set prometheusoperator.enabled=false \
      --set logCaptureEnabled=false \
      coherence/coherence
   ```
   
1. Ensure the pods are running:

   ```bash
   $ kubectl get pods -n sample-coherence-ns
   ```
   ```console
   NAME                                  READY   STATUS    RESTARTS   AGE
   coherence-operator-66f9bb7b75-hqk4l   1/1     Running   0          13m
   storage-coherence-0                   1/1     Running   0          3m
   storage-coherence-1                   1/1     Running   0          2m
   storage-coherence-2                   1/1     Running   0          44s
   ```
   
1. Add data to the cluster using the Coherence console

   Connect to the console and create a cache:

   ```bash
   $ kubectl exec -it --namespace sample-coherence-ns storage-coherence-0 bash /scripts/startCoherence.sh console
   ```   
   
   At the `Map (?):` prompt, type `cache test`.  This creates a cache in the service `PartitionedCache`.
   
   Use the following command to add 50000 objects of size 1024 bytes, starting at index 0 and using batches of 100.
   
   ```bash
   bulkput 50000 1024 0 100
   Wed Apr 24 01:17:44 GMT 2019: adding 50000 items (starting with #0) each 1024 bytes ...
   Wed Apr 24 01:18:11 GMT 2019: done putting (26802ms, 1917KB/sec, 1865 items/sec)
   ```
   
   At the prompt, type `size` and it will show 50000.
   
   The `help` command shows all the available command options.
   
   Then type `bye` to exit the console.
      
1. Add data to the cluster through CohQL

   Connect to the CohQL client using the following command:

   ```bash
   $ kubectl exec -it --namespace sample-coherence-ns storage-coherence-0 bash /scripts/startCoherence.sh queryplus
   ```   
   
   Use the following command at the CohQL prompt to view and manipulate data:
   
   ```bash
   select count() from 'test';

   ... Cluster and service join messages here

   Results
   50000

   select value() from 'test' where key() = 0
   Results
   [B@63a5e46c

   insert into 'test' key(1) value('value 1');
   Results
   [B@7a34b7b8

   select value() from 'test' where key() = 1
   Results
   "value 1"
   ```
   
   At the CohQL prompt, type `commands` to view all the commands that are available, while `help` shows detailed information.
   
   Type`quit` at the prompt to exit CohQL.
   
   Refer to the [Coherence Documentation](https://docs.oracle.com/middleware/1221/coherence/develop-applications/api_cq.htm#COHDG5264) for more information about CohQL.

## Uninstall the Charts

Use the following command to delete the chart installed in this sample:

```bash
$ helm delete storage --purge
```

Before starting another sample, ensure that all the pods are removed from previous sample.

If you want to remove the `coherence-operator`, then include it in the `helm delete` command.
