# Storage-disabled client in cluster via interceptor

In this sample we will deploy 2 tiers, a storage-enabled data tier and
storage-disabled client tier. The client tier will use an Interceptor to
startup a Coherence storage-disabled client and perform some processing.

The advantage of using this method, is because you are running DefaultCacheServer
process, all the Prometheus metrics will be collected for the storage-disabled member.

This is achieved by using 2 helm installs, which both include a sidecar container
for the data tier and client tier cache config as well as the interceptor code.

In this example we only have one storage-disabled client, but you could change this by setting `clusterSize` to a value other that one on the client chart install.

[Return to Storage-Disabled clients samples](../) / [Return to Coherence Deployments samples](../../) / [Return to samples](../../../README.md#list-of-samples)

## Sample files

* [src/main/docker/Dockerfile](src/main/docker/Dockerfile) - Dockerfile for creating side-car image from which configuration
  and server side jar will be read from at pod startup

* [src/main/resources/conf/storage-cache-config.xml](src/main/resources/conf/storage-cache-config.xml) - cache config for storage-enabled tier

* [src/main/resources/conf/interceptor-cache-config.xml](src/main/resources/conf/interceptor-cache-config.xml) - cache config for storage-disabled tier

* [src/main/java/com/oracle/coherence/examples/DemoInterceptor.java](src/main/java/com/oracle/coherence/examples/DemoInterceptor.java) - interceptor that will start our `mock` client

Note if you wish to enable Prometheus or log capture, change the following in the helm installs to `true`. Their default values are false, but they are set to `false` in the samples below for completeness.

* Prometheus: `--set prometheusoperator.enabled=true`

* Log capture: `--set logCaptureEnabled=true`

## Prerequisites

Ensure you have already installed the Coherence Operator by using the instructions [here](here).

## Installation Steps

1. Change to the `samples/coherence-deployments/storage-disabled/interceptor` directory and ensure you have your maven build     
   environment set for JDK11 and build the project.

   ```bash
   mvn clean install -P docker
   ```

   The above will build the Docker image with the cache configuration files and compiled Java classes.

   > Note: If you are running against a remote Kubernetes cluster you will need to
   > push the above image to your repository accessible to that cluster. You will also need to 
   > prefix the image name in your `helm` command below.

1. The result of the above is the docker image will be built with the cache configuration files
   and compiled Java classes with the name in the format interceptor-sample:${version}.

   For Example:

   ```bash
   interceptor-sample:1.0.0-SNAPSHOT
   ```

   **Note:** If you are running against a remote Kubernetes cluster you will need to
   push the above image to your repository accessible to that cluster.


1. Install the Coherence cluster

   We are also setting the cluster-name to `interceptor-cluster` which will we use in later steps.

   ```bash
   helm install \
      --namespace sample-coherence-ns \
      --name storage \
      --set clusterSize=3 \
      --set cluster=interceptor-cluster \
      --set imagePullSecrets=sample-coherence-secret \
      --set store.cacheConfig=storage-cache-config.xml \
      --set prometheusoperator.enabled=false \
      --set logCaptureEnabled=false \
      --set userArtifacts.image=interceptor-sample:1.0.0-SNAPSHOT \
      --version 1.0.0-SNAPSHOT coherence-community/coherence
   ```

   Because we use stateful sets, the coherence cluster will start one pod at a time.
   
   You can change this by using `--set store.podManagementPolicy=Parallel` in the above command.
    
   Use `kubectl get pods -n sample-coherence-ns` to ensure that all pods are running.
   All 3 storage-coherence-0/1/2 pods should be running and ready, as below:

   ```bash
   NAME                                                     READY   STATUS    RESTARTS   AGE
   storage-coherence-0                                      1/1     Running   0          4m
   storage-coherence-1                                      1/1     Running   0          2m
   storage-coherence-2                                      1/1     Running   0          1m
   ```

1. Install the storage-disabled client tier

   The following are set to ensure this release will connect to the Coherence clusterSize
   created by the `storage` release:

   * `--set cluster=interceptor-cluster` - same cluster name

   * `--set store.wka=storage-coherence-headless` - ensures it can contact the cluster

   * `--set prometheusoperator.enabled=false` - set storage to false

   * `--set store.cacheConfig=interceptor-cache-config.xml` - uses interceptor cache config from sidecar

   ```bash
   helm install \
     --namespace sample-coherence-ns \
     --set cluster=interceptor-cluster \
     --set clusterSize=1 \
     --set store.wka=storage-coherence-headless \
     --set prometheusoperator.enabled=true \
     --name interceptor-client-tier \
     --set imagePullSecrets=sample-coherence-secret \
     --set store.cacheConfig=interceptor-cache-config.xml \
     --set prometheusoperator.enabled=false \
     --set logCaptureEnabled=false \
     --set userArtifacts.image=interceptor-sample:1.0.0-SNAPSHOT \
     --version 1.0.0-SNAPSHOT coherence-community/coherence
   ```

   To confirm the storage-disabled client has joined the cluster you can look at the logs using:

   ```bash
   kubectl logs interceptor-client-tier-coherence-0 -n sample-coherence-ns -f
   ```

   This should continuously follow the log and you should see messages similar to the following
   indicating that the storage-disabled client is inserting data.

   ```bash
   2019-03-20 08:33:43.138/48.454 Oracle Coherence GE 12.2.1.4.0 <Info> (thread=pool-1-thread-1, member=4): Inserted key=40, value=08:33:43
   2019-03-20 08:33:44.143/49.459 Oracle Coherence GE 12.2.1.4.0 <Info> (thread=pool-1-thread-1, member=4): Inserted key=41, value=08:33:44
   ```

## Uninstalling the Charts

Carry out the following commands to delete the two charts installed in this sample.

```bash
helm delete storage interceptor-client-tier --purge
```

Before starting another sample, ensure that all the pods are gone from previous sample.

If you wish to remove the `coherence-operator`, then include it in the `helm delete` command above.
