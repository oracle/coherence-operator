# Scaling a Coherence deployment via kubectl

The Coherence Operator leverages Kubernetes Statefulsets to ensure that
scale up and scale down operations allow the underlying Coherence
cluster nodes sufficient time to re-balance the cluster data to ensure no data is lost.

This sample shows you how to scale a statefulset using `kubectl`. 

[Return to Coherence Operator samples](../) / [Return to samples](../../README.md#list-of-samples)

## Prerequisites

Ensure you have already installed the Coherence Operator by using the instructions [here](../../../README.md#install-the-coherence-operator).

## Installation Steps

1. Install the Coherence cluster

   Issue the following to install the cluster with only 2 replicas.

   ```bash
   $ helm install \
      --namespace sample-coherence-ns \
      --name storage \
      --set clusterSize=2 \
      --set cluster=coherence-cluster \
      --set imagePullSecrets=sample-coherence-secret \
      --set prometheusoperator.enabled=false \
      --set logCaptureEnabled=false \
      --version 1.0.0-SNAPSHOT coherence-community/coherence
   ```
   
1. Ensure both of the pods are running:

   ```bash
   $ kubectl get pods -n sample-coherence-ns
   NAME                                  READY   STATUS    RESTARTS   AGE
   coherence-operator-66f9bb7b75-hqk4l   1/1     Running   0          13m
   storage-coherence-0                   1/1     Running   0          3m
   storage-coherence-1                   1/1     Running   0          2m
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

1. Scale to 4 nodes using via `kubect`

   Scale the statefulset to 4 nodes.
  
   ```bash
   $ kubectl scale statefulsets storage-coherence --namespace sample-coherence-ns --replicas=4
   ```  
   
   ```bash
   $ kubectl get pods -n sample-coherence-ns
   NAME                                  READY   STATUS    RESTARTS   AGE 
   storage-coherence-0                   1/1     Running   0          10m
   storage-coherence-1                   1/1     Running   0          9m
   storage-coherence-2                   1/1     Running   0          3m
   storage-coherence-3                   1/1     Running   0          1m
   ```
   
   Wait for the number of `coherence-storage` pods to be 4 and all of them to be ready.
   
1. Check the size of the cache

   Use instructions above to access the `console` and to confirm the size is still 50000. 
   Ensure you type `cache test` at the `Map:` prompt and then`size`.
   
1. Scale to 2 nodes via `kubectl`    
  
   ```bash
   $ kubectl scale statefulsets storage-coherence --namespace sample-coherence-ns --replicas=2
   ``` 
   
   Before each node is stopped to reach the desired number of replicas, the cluster is checked to ensure
   the cluster is balanced to ensure that no data is lost.
   
    
   ```bash
   $ kubectl get pods -n sample-coherence-ns
   NAME                                  READY   STATUS    RESTARTS   AGE 
   storage-coherence-0                   1/1     Running   0          22m
   storage-coherence-1                   1/1     Running   0          21m
   ```
   
1. Check the size of the cache

   Use instructions above to access the `console` and to confirm the size is still 50000. 
   Ensure you type `cache test` at the `Map:` prompt and then`size`.
      
## Verifying Grafana Data (If you enabled Prometheus)

Access Grafana using the instructions [here](../../../README.md#access-grafana).

## Verifying Kibana Logs (if you enabled log capture)

Access Kibana using the instructions [here](../../../README.md#access-kibana).

## Uninstalling the Charts

Carry out the following commands to delete the chart and PV's created in this sample.

```bash
$ helm delete storage --purge
```
    
   
   
   