# Using Multiple Coherence*Extend Proxies       

This sample shows how to run multiple Proxy services on a single cluster node. To support this within the Coherence Operator, a custom cache configuration must be added as well as an additional pod must be exposed on the pod.

[Return to Coherence*Extend samples](../) / [Return to Coherence Deployments samples](../../) / [Return to samples](../../../README.md#list-of-samples)

## Sample files

* [src/main/docker/Dockerfile](src/main/docker/Dockerfile) - Dockerfile for creating sidecar image from which the configuration will be read at pod startup

* [src/main/resources/client-cache-config.xml](src/main/resources/client-cache-config.xml) - Cache configuration for the extend client

* [src/main/resources/conf/storage-cache-config.xml](src/main/resources/conf/storage-cache-config.xml) - Cache configuration for storage-tier

## Prerequisites

Ensure that you have installed Coherence Operator by following the instructions [here](../../../README.md#install-the-coherence-operator).

## Installation Steps

1. Change to the `samples/coherence-deployments/extend/multiple` directory. Ensure that you have your maven build environment set for JDK8, and build the project.

   ```bash
   $ mvn clean install -P docker
   ```

  As a result, the docker image will be built with the cache configuration files and compiled Java classes with the name in the format, `multiple-proxy-sample:${version}`. For example,

   ```bash
   multiple-proxy-sample:1.0.0-SNAPSHOT
   ```

   > **Note** If you are running against a remote Kubernetes cluster, then you must
   > push the above image to your repository accessible to that cluster. You must also
   > prefix the image name in your `helm` command below.

2. Install the Coherence cluster.

   Set the following additional properties:

  * `--set store.ports.custom-port=20001` - This property sets the port for the second proxy server.

  *  `--set store.javaOpts="-Dcoherence.extend.port2=20001"` -  Set this property for cache configuration.

   ```bash
   $ helm install \
      --namespace sample-coherence-ns \
      --name storage \
      --set clusterSize=3 \
      --set cluster=multiple-proxy-cluster \
      --set imagePullSecrets=sample-coherence-secret \
      --set store.cacheConfig=storage-cache-config.xml \
      --set logCaptureEnabled=false \
      --set userArtifacts.image=multiple-proxy-sample:1.0.0-SNAPSHOT \
      --set store.javaOpts="-Dcoherence.extend.port2=20001" \
      --set store.ports.custom-port=20001 \
      coherence/coherence
   ```

    Once the installation is complete, run the following command to retrieve the list of pods:

   ```bash
   $ kubectl get pods -n sample-coherence-ns
   ```
   ```console
   NAME                   READY   STATUS    RESTARTS   AGE
   storage-coherence-0    1/1     Running   0          4m
   storage-coherence-1    1/1     Running   0          2m   
   storage-coherence-2    1/1     Running   0          2m
   ```

3. Ensure that both Proxy services are running, by using the following `kubectl` command:

   ```bash
   $ kubectl logs storage-coherence-0   -n sample-coherence-ns | grep 'TcpAcceptor now listening' | grep ProxyService   
   2019-05-01 01:55:38.856/8.215 Oracle Coherence GE 12.2.1.4.0 <Info> (thread=Proxy:ProxyService1:TcpAcceptor, member=1): TcpAcceptor now listening for connections on storage-coherence-0.coherence.sample-coherence-ns.svc.cluster.local:20000
   2019-05-01 01:55:38.955/8.313 Oracle Coherence GE 12.2.1.4.0 <Info> (thread=Proxy:ProxyService2:TcpAcceptor, member=1): TcpAcceptor now listening for connections on storage-coherence-0.coherence.sample-coherence-ns.svc.cluster.local:20001
   ```   

4. Port forward the proxy ports on the proxy-tier.

   ```bash
   $ kubectl port-forward -n sample-coherence-ns storage-coherence-0 20000:20000
   ```

   ```bash
   $ kubectl port-forward -n sample-coherence-ns storage-coherence-0 20001:20001
   ```

5. Connect via CohQL using port 20000 and run the following command:

   ```bash
   $ mvn exec:java -Dproxy.port=20000
   ```

   Run the following CohQL commands to insert data into the cluster via proxy port 20000.

   ```
   insert into 'test' key('key-1') value('value-1');

   select key(), value() from 'test';
   Results
   ["key-1", "value-1"]

   > select count() from 'test';
   Results
   1
   ```

   You should see a message indicating the connection to 127.0.0.1:20000, as shown in the output:

   ```console
   2019-05-01 10:01:08.007/5.684 Oracle Coherence GE 12.2.1.4.0 <D6> (thread=com.tangosol.coherence.dslquery.QueryPlus.main(), member=n/a): Connecting Socket to 127.0.0.1:20000
   ```

6. Connect via CohQL using port 20001 using the following command:

   ```bash
   $ mvn exec:java -Dproxy.port=20001
   ```

   Run the following `CohQL` commands to insert data into the cluster via proxy port 20001:

   ```
   select key(), value() from 'test';
   Results
   ["key-1", "value-1"]

   select count() from 'test';
   Results
   1
   ```

   You should see a message indicating the connection to 127.0.0.1:20001, as shown in the output:

   ```bash
   2019-05-01 10:05:18.764/20.659 Oracle Coherence GE 12.2.1.4.0 <D6> (thread=com.tangosol.coherence.dslquery.QueryPlus.main(), member=n/a): Connecting Socket to 127.0.0.1:20001   
   ```

## Uninstall the Charts

Run the following command to delete both the charts installed in this sample:

```bash
$ helm delete storage --purge
```

Before starting another sample, ensure that all  pods are removed from the previous sample. To remove `coherence-operator`, use the `helm delete` command.
