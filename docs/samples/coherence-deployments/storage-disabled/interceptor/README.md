# Storage-Disabled Client in Cluster via Interceptor

This sample demonstrates deploying 2 tiers, a storage-enabled data tier, and a storage-disabled client tier. The client tier uses an interceptor to
start up a Coherence storage-disabled client and perform some processing.

The advantage of using this method is that when you run the DefaultCacheServer
process, all Prometheus metrics is collected for the storage-disabled member.

This is achieved by using two `helm install` commands, both of which include a sidecar container
for the data tier, and a client tier cache configuration as well as the interceptor code.

This sample uses only one storage-disabled client. You can change this by setting `clusterSize` to a value other than one on the client chart install.

[Return to Storage-Disabled clients samples](../) / [Return to Coherence Deployments samples](../../) / [Return to samples](../../../README.md#list-of-samples)

## Sample files

* [src/main/docker/Dockerfile](src/main/docker/Dockerfile) - Dockerfile for creating sidecar image from which the configuration and the server side jar will be read at pod startup

* [src/main/resources/conf/storage-cache-config.xml](src/main/resources/conf/storage-cache-config.xml) - Cache configuration for storage-enabled tier

* [src/main/resources/conf/interceptor-cache-config.xml](src/main/resources/conf/interceptor-cache-config.xml) - Cache configuration for storage-disabled tier

* [src/main/java/com/oracle/coherence/examples/DemoInterceptor.java](src/main/java/com/oracle/coherence/examples/DemoInterceptor.java) - Interceptor that starts our `mock` client

Note that if you want to enable Prometheus or log capture, set the following properties in the `helm install` command to `true`. These properties are set to false by default. However, in this sample, these properties have been set to `false` for completeness.

* Prometheus: `--set prometheusoperator.enabled=true`

* Log capture: `--set logCaptureEnabled=true`

## Prerequisites

Ensure that you have installed the Coherence Operator by following the instructions [here](../../../README.md#install-the-coherence-operator).

## Installation Steps

1. Change to the `samples/coherence-deployments/storage-disabled/interceptor` directory. Ensure that you have your maven build environment set for JDK8, and build the project.

   ```bash
   $ mvn clean install -P docker
   ```

   This builds the Docker image with the cache configuration files and compiled Java classes, with the name in the format, `interceptor-sample:${version}`. For example,

   ```bash
   interceptor-sample:1.0.0-SNAPSHOT
   ```

   > **Note:** If you are running against a remote Kubernetes cluster, you must push the above image to your repository accessible to that cluster. You must also
   > prefix the image name in the `helm` command as shown below.



2. Install the Coherence cluster.

   Set the cluster-name to `interceptor-cluster` for use in later steps.

   ```bash
   $ helm install \
      --namespace sample-coherence-ns \
      --name storage \
      --set clusterSize=3 \
      --set cluster=interceptor-cluster \
      --set imagePullSecrets=sample-coherence-secret \
      --set store.cacheConfig=storage-cache-config.xml \
      --set prometheusoperator.enabled=false \
      --set logCaptureEnabled=false \
      --set userArtifacts.image=interceptor-sample:1.0.0-SNAPSHOT \
      coherence/coherence
   ```

   Once the installation is complete, run the following command to retrieve the list of pods:

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


3. Install the storage-disabled client tier.

   Set the following properties to ensure that this release connects to the Coherence cluster
   created by the `storage` release:

   * `--set cluster=interceptor-cluster` - Uses the same cluster name

   * `--set store.wka=storage-coherence-headless` - Ensures that it can contact the cluster

   * `--set prometheusoperator.enabled=false` - Sets storage to false

   * `--set store.cacheConfig=interceptor-cache-config.xml` - Uses interceptor cache configuration from the sidecar image

   ```bash
   $ helm install \
     --namespace sample-coherence-ns \
     --set cluster=interceptor-cluster \
     --set clusterSize=1 \
     --set store.wka=storage-coherence-headless \
     --set prometheusoperator.enabled=true \
     --name interceptor-client-tier \
     --set imagePullSecrets=sample-coherence-secret \
     --set store.cacheConfig=interceptor-cache-config.xml \
     --set store.storageEnabled=false \
     --set prometheusoperator.enabled=false \
     --set logCaptureEnabled=false \
     --set userArtifacts.image=interceptor-sample:1.0.0-SNAPSHOT \
     coherence/coherence
   ```

   To confirm that the storage-disabled client has joined the cluster, you can look at the logs using the following `kubectl` commands:

   ```bash
  $  kubectl logs interceptor-client-tier-coherence-0 -n sample-coherence-ns -f
   ```

   This continuously follows the log, and you see messages similar to the following
   indicating that the storage-disabled client is inserting data:

   ```bash
   2019-03-20 08:33:43.138/48.454 Oracle Coherence GE 12.2.1.3.0 <Info> (thread=pool-1-thread-1, member=4): Inserted key=40, value=08:33:43
   2019-03-20 08:33:44.143/49.459 Oracle Coherence GE 12.2.1.3.0 <Info> (thread=pool-1-thread-1, member=4): Inserted key=41, value=08:33:44
   ```

## Uninstall the Charts

Run the following command to delete both the charts installed in this sample.

```bash
$ helm delete storage interceptor-client-tier --purge
```

Before starting another sample, ensure that all  pods are removed from the previous sample. To remove `coherence-operator`, use the `helm delete` command.
