# Access Coherence via default proxy port

This sample shows how to access the Coherence cluster via the default proxy service
exposed on port 20000.

[Return to Coherence*Extend samples](../) / [Return to Coherence Deployments samples](../../) / [Return to samples](../../../README.md#list-of-samples)

## Sample files

* [src/main/resources/client-cache-config.xml](src/main/resources/client-cache-config.xml) - Client config for extend client

## Prerequisites

Ensure you have already installed the Coherence Operator by using the instructions [here](../../../README.md#install-the-coherence-operator).

## Installation Steps

1. Change to the `samples/coherence-deployments/extend/default` directory and ensure you have your maven build     
   environment set for JDK11 and build the project.
   
   ```bash
   $ mvn clean install
   ```

1. Install the Coherence cluster

   ```bash
   $ helm install \
      --namespace sample-coherence-ns \
      --name storage \
      --set clusterSize=3 \
      --set cluster=storage-tier-cluster \
      --set imagePullSecrets=sample-coherence-secret \
      --set prometheusoperator.enabled=false \
      --set logCaptureEnabled=false \
      coherence-community/coherence
   ```

   Because we use stateful sets, the coherence cluster will start one pod at a time.
    
   Use `kubectl get pods -n sample-coherence-ns` to ensure that all pods are running. All 3 storage-coherence-0/1/2 pods should be running and ready, as below:

   ```bash
   NAME                  READY   STATUS    RESTARTS   AGE
   storage-coherence-0   1/1     Running   0          4m
   storage-coherence-1   1/1     Running   0          2m
   storage-coherence-2   1/1     Running   0          1m
   ```

1. Port forward the proxy port on the storage-coherence-0 pod.

   ```bash
   $ kubectl port-forward -n sample-coherence-ns storage-coherence-0 20000:20000
   ```

1. Connect via QueryPlus and issue CohQL commands

   Issue the following command to run QueryPlus:

   ```bash
   $ mvn exec:java
   ```

   Run the following CohQL commands to insert data into the cluster.

   ```sql
   insert into 'test' key('key-1') value('value-1');

   select key(), value() from 'test';
   Results
   ["key-1", "value-1"]

   select count() from 'test';
   Results
   1
   ```

## Uninstalling the Chart

Carry out the following commands to delete the chart installed in this sample.

```bash
$ helm delete storage --purge
```

Before starting another sample, ensure that all the pods are gone from previous samples.

If you wish to remove the `coherence-operator`, then include it in the `helm delete` command above.
