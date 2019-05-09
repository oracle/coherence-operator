# Access Coherence console and CohQL on a cluster node

The Coherence Console and CohQL (Coherence Query Language) are developer tools for interacting with a 
Coherence cluster.

This samples shows using `kubectl exec' to connect to any of the pods and start 
either of these tools as a storage-disabled client.

[Return to Management samples](../) / [Return to samples](../../README.md#list-of-samples)

## Prerequisites

Ensure you have already installed the Coherence Operator by using the instructions [here](../../../README.md#install-the-coherence-operator).

## Installation Steps

1. Install the Coherence cluster

   Issue the following to install the cluster with Persistence and Snapshot enabled:

   ```bash
   $ helm install \
      --namespace sample-coherence-ns \
      --name storage \
      --set clusterSize=3 \
      --set cluster=coherence-cluster \
      --set imagePullSecrets=sample-coherence-secret \
      --set prometheusoperator.enabled=false \
      --set logCaptureEnabled=false \
      coherence-community/coherence
   ```
   
1. Ensure the pods are running:

   ```bash
   $ kubectl get pods -n sample-coherence-ns
   NAME                                  READY   STATUS    RESTARTS   AGE
   coherence-operator-66f9bb7b75-hqk4l   1/1     Running   0          13m
   storage-coherence-0                   1/1     Running   0          3m
   storage-coherence-1                   1/1     Running   0          2m
   storage-coherence-2                   1/1     Running   0          44s
   ```
   
1. Add data to the cluster via the Coherence Console

   Connect to the Coherence `console` using the following to create a cache.

   ```bash
   $ kubectl exec -it --namespace sample-coherence-ns storage-coherence-0 bash /scripts/startCoherence.sh console
   ```   
   
   At the `Map (?):` prompt, type `cache test`.  This will create a cache in the service `PartitionedCache`.
   
   Use the following to add 50,000 objects of size 1024 bytes, starting at index 0 and using batches of 100.
   
   ```bash
   bulkput 50000 1024 0 100
   Wed Apr 24 01:17:44 GMT 2019: adding 50000 items (starting with #0) each 1024 bytes ...
   Wed Apr 24 01:18:11 GMT 2019: done putting (26802ms, 1917KB/sec, 1865 items/sec)
   ```
   
   At the prompt, type `size` and it should show 50000.
   
   The `help` command will show all commands available.
   
   Then type `bye` to exit the `console`.
      
1. Add data to the cluster via CohQL      

   Connect to the `CohQL` client using the following:.

   ```bash
   $ kubectl exec -it --namespace sample-coherence-ns storage-coherence-0 bash /scripts/startCoherence.sh queryplus
   ```   
   
   Use the following at the `CohQL>` prompt to view and manipulate data:
   
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
   
   At the `CohQL>` prompt type `commands` to show view of all commands available, while `help` will show
   detailed help.
   
   Type`quit` at the prompt to exit CohQL.
   
   Please see [Coherence Documentation](https://docs.oracle.com/middleware/1221/coherence/develop-applications/api_cq.htm#COHDG5264) 
   for more information on CohQL.

## Uninstalling the Charts

Carry out the following commands to delete the chart installed in this sample.

```bash
$ helm delete storage --purge
```

Before starting another sample, ensure that all the pods are gone from previous sample.

If you wish to remove the `coherence-operator`, then include it in the `helm delete` command above.
   
