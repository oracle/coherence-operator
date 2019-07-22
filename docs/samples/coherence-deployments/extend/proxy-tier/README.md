# Access Coherence via a Separate Proxy Tier

This sample shows how to deploy two tiers, a storage-enabled data tier and a storage-disabled proxy tier. This is a common scenario when you are using Coherence*Extend
to connect to a cluster and when you want to separate the proxy tier from the data tier.

This is achieved by using two `helm install` commands, both of which include a sidecar container
for the data tier and proxy tier cache configuration.

[Return to Coherence*Extend samples](../) / [Return to Coherence Deployments samples](../../) / [Return to samples](../../../README.md#list-of-samples)

## Sample files

* [src/main/docker/Dockerfile](src/main/docker/Dockerfile) - Dockerfile for creating sidecar image from which the configuration will be read at pod startup

* [src/main/resources/client-cache-config.xml](src/main/resources/client-cache-config.xml) - Cache configuration for the extend client

* [src/main/resources/conf/proxy-cache-config.xml](src/main/resources/conf/proxy-cache-config.xml) - Cache configuration for proxy-tier

* [src/main/resources/conf/storage-cache-config.xml](src/main/resources/conf/storage-cache-config.xml) - Cache configuration for storage-tier

## Prerequisites

Ensure that you have installed the Oracle Coherence Operator by following the instructions [here](../../../README.md#install-the-coherence-operator).

## Installation Steps

1. Change to the `samples/coherence-deployments/extend/proxy-tier` directory. Ensure that you have your maven build environment set for JDK8, and build the project.

   ```bash
   $ mvn clean install -P docker
   ```

   This builds the Docker image with the cache configuration files and compiled Java classes, with the name in the format `proxy-tier-sample:${version}`. For example,

    ```bash
    proxy-tier-sample:1.0.0-SNAPSHOT
     ```

   > **Note:** If you are running against a remote Kubernetes cluster, you must push the above image to your repository accessible to that cluster. You must also prefix the image name in your `helm` command as shown below.

2. Install the Coherence cluster.

   ```bash
   $ helm install \
      --namespace sample-coherence-ns \
      --name storage \
      --set clusterSize=3 \
      --set cluster=proxy-tier-cluster \
      --set imagePullSecrets=sample-coherence-secret \
      --set store.cacheConfig=storage-cache-config.xml \
      --set prometheusoperator.enabled=false \
      --set logCaptureEnabled=false \
      --set userArtifacts.image=proxy-tier-sample:1.0.0-SNAPSHOT \
      coherence/coherence
   ```

  Once the installation is complete, run the following command to retrieve the list of pods:

   ```bash
   $ kubectl get pods -n sample-coherence-ns
   ```
  All the three storage-coherence pods should be running and ready, as shown in the output:

   ```console
   NAME                   READY   STATUS    RESTARTS   AGE
   storage-coherence-0    1/1     Running   0          4m
   storage-coherence-1    1/1     Running   0          2m   
   storage-coherence-2    1/1     Running   0          2m
   ```

3. Install the storage-disabled proxy-tier.

   Set the following properties to ensure that this release connects to the Coherence clusterSize
   created by the `storage` release:

   * `--set cluster=proxy-tier-cluster` - Uses the same cluster name

   * `--set store.wka=storage-coherence-headless` - Ensures that it can contact the cluster

   * `--set cluster=proxy-tier-cluster` - Ensures that the cluster name is the same

   * `--set prometheusoperator.enabled=false` - Sets storage to false

   * `--set store.cacheConfig=proxy-cache-config.xml` - Uses proxy cache configuration from sidecar

   > **Note:** For the proxy-tier, we are using a clusterSize of one, to save resources. You can
   > scale out the proxy-tier for high availability purposes.

   ```bash
   $ helm install \
     --namespace sample-coherence-ns \
     --set cluster=proxy-tier-cluster \
     --set clusterSize=1 \
     --set store.storageEnabled=false \
     --set store.wka=storage-coherence-headless \
     --set prometheusoperator.enabled=false \
     --name proxy-tier \
     --set imagePullSecrets=sample-coherence-secret \
     --set store.cacheConfig=proxy-cache-config.xml \
     --set prometheusoperator.enabled=false \
     --set userArtifacts.image=proxy-tier-sample:1.0.0-SNAPSHOT \
     coherence/coherence
   ```

   To confirm that the proxy-tier has joined the cluster, you can look at the logs using:

   ```bash
   $ kubectl logs proxy-tier-coherence-0 -n sample-coherence-ns | grep ActualMemberSet
   ```

   This should return the following, which indicates that there are now four members:

   ```console
   ActualMemberSet=MemberSet(Size=4
   ```

   You should now see three charts installed:

   ```bash
   $ helm ls
   NAME              	REVISION	UPDATED                 	STATUS  	CHART                            	APP VERSION   	NAMESPACE          
   coherence-operator	1       	Wed Mar 20 14:12:31 2019	DEPLOYED	coherence-operator-1.0.0-SNAPSHOT	1.0.0-SNAPSHOT	sample-coherence-ns
   proxy-tier        	1       	Wed Mar 20 14:54:57 2019	DEPLOYED	coherence-1.0.0-SNAPSHOT         	1.0.0-SNAPSHOT	sample-coherence-ns
   storage           	1       	Wed Mar 20 14:53:58 2019	DEPLOYED	coherence-1.0.0-SNAPSHOT         	1.0.0-SNAPSHOT	sample-coherence-ns
   ```

4. Port forward the proxy port on the proxy-tier.

   ```bash
   $ kubectl port-forward -n sample-coherence-ns proxy-tier-coherence-0 20000:20000
   ```

5. Connect via QueryPlus and run the `CohQL` commands.

   Run the following command to execute QueryPlus:

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

## Uninstall the Charts

Run the following command to delete both the charts installed in this sample:

```bash
$ helm delete storage proxy-tier --purge
```

Before starting another sample, ensure that all  pods are removed from the previous sample. To remove `coherence-operator`, use the `helm delete` command.
