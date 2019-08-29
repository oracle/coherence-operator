# Access Coherence via the Default Proxy Port

This sample shows how to access the Coherence cluster via the default proxy service
exposed on port 20000.

[Return to Coherence*Extend samples](../) / [Return to Coherence Deployments samples](../../) / [Return to samples](../../../README.md#list-of-samples)

## Sample files

* [src/main/resources/client-cache-config.xml](src/main/resources/client-cache-config.xml) - Client configuration for the extend client

## Prerequisites

Ensure that you have installed the Coherence Operator by following the instructions [here](../../../README.md#install-the-coherence-operator).

## Installation Steps

1. Change to the `samples/coherence-deployments/extend/default` directory. Ensure that you have your maven build environment set for JDK8, and build the project.

   ```bash
   $ mvn clean install
   ```

2. Install the Coherence cluster.

   ```bash
   $ helm install \
      --namespace sample-coherence-ns \
      --name storage \
      --set clusterSize=3 \
      --set cluster=storage-tier-cluster \
      --set imagePullSecrets=sample-coherence-secret \
      --set logCaptureEnabled=false \
      coherence/coherence
   ```

   Once the installation is complete, get the list of pods by using the `kubectl` command:

   ```bash
   $ kubectl get pods -n sample-coherence-ns
   ```

   All the three storage-coherence pods should be running and ready, as shown in the output:

   ```console
   NAME                  READY   STATUS    RESTARTS   AGE
   storage-coherence-0   1/1     Running   0          4m
   storage-coherence-1   1/1     Running   0          2m
   storage-coherence-2   1/1     Running   0          1m
   ```


3. Port forward the proxy port on the storage-coherence-0 pod using the `kubectl` command:

   ```bash
   $ kubectl port-forward -n sample-coherence-ns storage-coherence-0 20000:20000
   ```

4. Connect via CohQL and run the following commands:

   ```bash
   $ mvn exec:java
   ```

   Run the following `CohQL` commands to insert data into the cluster.

   ```sql
   insert into 'test' key('key-1') value('value-1');

   select key(), value() from 'test';
   Results
   ["key-1", "value-1"]

   select count() from 'test';
   Results
   1
   ```

## Uninstall the Chart

Run the following command to delete the chart installed in this sample.

```bash
$ helm delete storage --purge
```

Before starting another sample, ensure that all  pods are removed from the previous sample. To remove `coherence-operator`, use the `helm delete` command.
