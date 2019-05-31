# Scaling a Coherence deployment via kubectl

The Coherence Operator leverages Kubernetes Statefulsets to ensure that
scale up and scale down operations are carried out one pod at a time.
 
When scaling down, you should only scale down by 1 pod at a time and check for `HAStatus` before 
continuing. THis will ensure the cluster nodes have sufficient time to re-balance 
the cluster data to ensure no data is lost.

This sample shows you how to scale-up a statefulset using `kubectl` as well as scaling-down
1 Pod at a time and checking StatusHA via JVisualVM.

> **Note:** In a future version of the Coherence Operator, this manual `HAStatus` check will not be required.

[Return to Coherence Operator samples](../) / [Return to samples](../../README.md#list-of-samples)

## Prerequisites

1. Install Coherence Operator
   Ensure you have already installed the Coherence Operator by using the instructions [here](../../README.md#install-the-coherence-operator).

1. Download the JMXMP connector jar

   Please refer to [these instructions](../../management/jmx/README.md#Prerequisites) to download the 
   `JMXMP connector jar`.

## Installation Steps

1. Install the Coherence cluster

   Issue the following to install the cluster with only 2 replicas and 1 MBean Server Pod, which we will
   use to check statusHA.

   ```bash
   $ helm install \
      --namespace sample-coherence-ns \
      --name storage \
      --set clusterSize=2 \
      --set cluster=coherence-cluster \
      --set imagePullSecrets=sample-coherence-secret \
      --set prometheusoperator.enabled=false \
      --set logCaptureEnabled=false \
      --set store.jmx.enabled=true \
      --set store.jmx.replicas=1 \
      coherence/coherence
   ```
   
1. Ensure both of the pods are running:

   ```bash
   $ kubectl get pods -n sample-coherence-ns
   NAME                                    READY   STATUS    RESTARTS   AGE
   coherence-operator-66f9bb7b75-hqk4l     1/1     Running   0          13m
   storage-coherence-0                     1/1     Running   0          3m
   storage-coherence-1                     1/1     Running   0          2m
   storage-coherence-jmx-54f5d779d-svh29   1/1     Running   0          2m
   ```
   
   You should see a pod prefixed with `storage-coherence-jmx` in the above output.

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
   NAME                  READY   STATUS    RESTARTS   AGE 
   storage-coherence-0   1/1     Running   0          10m
   storage-coherence-1   1/1     Running   0          9m
   storage-coherence-2   1/1     Running   0          3m
   storage-coherence-3   1/1     Running   0          1m
   ```
   
   Wait for the number of `coherence-storage` pods to be 4 and all of them to be ready.
   
1. Check the size of the cache

   Use instructions above to access the `console` and to confirm the size is still 50000. 
   Ensure you type `cache test` at the `Map:` prompt and then`size`.
   
1. Port-forward the MBean Server Pod 

   Use the [instructions here](../../management/jmx/README.md#installation-steps) in step
   3 to port forward to the MBean Server Pod.
   
1. Connect via JVisualVM or JConsole

   Use the [instructions here](../../management/jmx/README.md#installation-steps) in step
   5 or 6  to connect to the cluster using JVisualVM or JConsole. 
   
1. Check `HAStatus` of `PartitionedCache` service

   > **Note:** We are using JVisualVM in this sample.
   
   Connect to the Mbean tab in JVisualVM and expand `Coherence`->`PartitionAssignment`->`PartitionedCache`->
   `DistributionCoordinator`.
   
   You should see something similar to the following with the `PartitionedCache` service `HAStatus` being a value
   other than `ENDANGERED` as well as `ServiceNodeCount` of 4.
   
   ![JVisualVM with 4 Nodes Running](img/jvisualvm-4-nodes.png)
    
   Ensure the `HAStatus` is correct before continuing. 
 
1. Scale down 1 node via `kubectl`    
  
   ```bash
   $ kubectl scale statefulsets storage-coherence --namespace sample-coherence-ns --replicas=3
   ``` 
   
   Using the step above, wait until the service `PartitionedCache` has `HAStatus` other than `ENDANGERED` 
   as well as `ServiceNodeCount` of 3.
   
   Once the above is completed, scale down the replicas to 2 nodes.
   
   ```bash
   $ kubectl get pods -n sample-coherence-ns
   NAME                                  READY   STATUS    RESTARTS   AGE 
   storage-coherence-0                   1/1     Running   0          22m
   storage-coherence-1                   1/1     Running   0          21m
   ```
   
   Your JVisualVM Mbeans tab should show the following:
   
   ![JVisualVM with 2 Nodes Running](img/jvisualvm-2-nodes.png)
   
1. Check the size of the cache

   Use instructions above to access the `console` and to confirm the size is still 50000. 
   Ensure you type `cache test` at the `Map:` prompt and then`size`.
      
## Verifying Grafana Data (If you enabled Prometheus)

Access Grafana using the instructions [here](../../README.md#access-grafana).

## Verifying Kibana Logs (if you enabled log capture)

Access Kibana using the instructions [here](../../README.md#access-kibana).

## Uninstalling the Charts

Carry out the following commands to delete the chart installed in this sample.

```bash
$ helm delete storage --purge
```
    
Before starting another sample, ensure that all the pods are gone from previous sample.

If you wish to remove the `coherence-operator`, then include it in the `helm delete` command above.
   
   
   
