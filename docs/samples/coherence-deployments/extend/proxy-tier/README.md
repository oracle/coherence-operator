# Access Coherence via separate proxy tier

This sample shows how we can deploy 2 tiers, a storage-enabled data tier and
storage-disabled proxy tier. This is a common scenario when using Coherence*Extend
to connect to a cluster and you with to separate the proxy tier from the data tier.

This is achieved by using 2 helm installs, which both include a sidecar container
for the data tier and proxy tier cache config.

[Return to Coherence*Extend samples](../) / [Return to Coherence Deployments samples](../../) / [Return to samples](../../../README.md#list-of-samples)

## Sample files

* [src/main/docker/Dockerfile](src/main/docker/Dockerfile) - Dockerfile for creating side-car image from which configuration
  will be read from at pod startup

* [src/main/resources/client-cache-config.xml](src/main/resources/client-cache-config.xml) - cache config for extend client

* [src/main/resources/conf/proxy-cache-config.xml](src/main/resources/conf/proxy-cache-config.xml) - cache config for proxy-tier

* [src/main/resources/conf/storage-cache-config.xml](src/main/resources/conf/storage-cache-config.xml) - cache config for storage-tier

## Prerequisites

Ensure you have already installed the Coherence Operator by using the instructions [here](../../../README.md#install-the-coherence-operator).

## Installation Steps

1. Change to the `samples/coherence-deployments/extend/proxy-tier` directory and ensure you have your maven build     
   environment set for JDK8 and build the project.

   ```bash
   $ mvn clean install -P docker
   ```

1. The result of the above is the docker image will be built with the cache configuration files
   with the name in the format proxy-tier-sample:${version}.

   For Example:

   ```bash
   proxy-tier-sample:1.0.0-SNAPSHOT
   ```

   > Note: If you are running against a remote Kubernetes cluster you will need to
   > push the above image to your repository accessible to that cluster. You will also need to 
   > prefix the image name in your `helm` command below.

1. Install the Coherence cluster

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
   
   Use `kubectl get pods -n sample-coherence-ns` to ensure that all pods are running.
   All 3 storage-coherence-0/1/2 pods should be running and ready, as below:

   ```bash
   NAME                   READY   STATUS    RESTARTS   AGE
   storage-coherence-0    1/1     Running   0          4m
   storage-coherence-1    1/1     Running   0          2m   
   storage-coherence-2    1/1     Running   0          2m
   ```

1. Install the storage-disabled proxy-tier

   The following are set to ensure this release will connect to the Coherence clusterSize
   created by the `storage` release:

   * `--set cluster=proxy-tier-cluster` - same cluster name

   * `--set store.wka=storage-coherence-headless` - ensures it can contact the cluster
   
   * `--set cluster=proxy-tier-cluster` - ensures the cluster name is the same

   * `--set prometheusoperator.enabled=false` - set storage to false

   * `--set store.cacheConfig=proxy-cache-config.xml` - uses proxy cache config from sidecar
   
   > Note: We are using a clusterSize for the proxy-tier of just 1, to save resources. You could
   > also scale out the proxy-tier for high availability purposes.

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
   
   To confirm the proxy-tier has joined the cluster you can look at the logs using:

   ```bash
   $ kubectl logs proxy-tier-coherence-0 -n sample-coherence-ns | grep ActualMemberSet
   ```

   This should return the following, which indicates there are now 4 members:

   ```bash
   ActualMemberSet=MemberSet(Size=4
   ```

   You should now see 3 charts installed:

   ```bash
   $ helm ls
   NAME              	REVISION	UPDATED                 	STATUS  	CHART                            	APP VERSION   	NAMESPACE          
   coherence-operator	1       	Wed Mar 20 14:12:31 2019	DEPLOYED	coherence-operator-1.0.0-SNAPSHOT	1.0.0-SNAPSHOT	sample-coherence-ns
   proxy-tier        	1       	Wed Mar 20 14:54:57 2019	DEPLOYED	coherence-1.0.0-SNAPSHOT         	1.0.0-SNAPSHOT	sample-coherence-ns
   storage           	1       	Wed Mar 20 14:53:58 2019	DEPLOYED	coherence-1.0.0-SNAPSHOT         	1.0.0-SNAPSHOT	sample-coherence-ns
   ```

1. Port forward the proxy port on the proxy-tier

   ```bash
   $ kubectl port-forward -n sample-coherence-ns proxy-tier-coherence-0 20000:20000
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

## Uninstalling the Charts

Carry out the following commands to delete the two charts installed in this sample.

```bash
$ helm delete storage proxy-tier --purge
```

Before starting another sample, ensure that all the pods are gone from previous sample.

If you wish to remove the `coherence-operator`, then include it in the `helm delete` command above.
