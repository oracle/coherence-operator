# Using multiple Coherence*Extend proxies       

This sample shows how to run multiple Proxy services on a single cluster node. 
To support this within the Coherence Operator, a custom cache configuration needs to be added as well as exposing
an additional port on the Pod.

[Return to Coherence*Extend samples](../) / [Return to Coherence Deployments samples](../../) / [Return to samples](../../../README.md#list-of-samples)

## Sample files

* [src/main/docker/Dockerfile](src/main/docker/Dockerfile) - Dockerfile for creating side-car image from which configuration
  will be read from at pod startup

* [src/main/resources/client-cache-config.xml](src/main/resources/client-cache-config.xml) - cache config for extend client

* [src/main/resources/conf/storage-cache-config.xml](src/main/resources/conf/storage-cache-config.xml) - cache config for storage-tier

## Prerequisites

Ensure you have already installed the Coherence Operator by using the instructions [here](../../../README.md#install-the-coherence-operator).

## Installation Steps

1. Change to the `samples/coherence-deployments/extend/multiple` directory and ensure you have your maven build     
   environment set for JDK11 and build the project.

   ```bash
   $ mvn clean install -P docker
   ```

1. The result of the above is the docker image will be built with the cache configuration files
   with the name in the format multiple-proxy-sample:${version}.

   For Example:

   ```bash
   multiple-proxy-sample:1.0.0-SNAPSHOT
   ```

   > Note: If you are running against a remote Kubernetes cluster you will need to
   > push the above image to your repository accessible to that cluster. You will also need to 
   > prefix the image name in your `helm` command below.

1. Install the Coherence cluster

   The following additional options are set:
   
   * `--set store.ports.custom-port=20001` - port for second proxy server
   
   * `--set store.javaOpts="-Dcoherence.extend.port2=20001"` - set property for cache config

   ```bash
   $ helm install \
      --namespace sample-coherence-ns \
      --name storage \
      --set clusterSize=3 \
      --set cluster=multiple-proxy-cluster \
      --set imagePullSecrets=sample-coherence-secret \
      --set store.cacheConfig=storage-cache-config.xml \
      --set prometheusoperator.enabled=false \
      --set logCaptureEnabled=false \
      --set userArtifacts.image=multiple-proxy-sample:1.0.0-SNAPSHOT \
      --set store.javaOpts="-Dcoherence.extend.port2=20001" \
      --set store.ports.custom-port=20001 \
      --version 1.0.0-SNAPSHOT coherence-community/coherence
   ```
   
   Use `kubectl get pods -n sample-coherence-ns` to ensure that all pods are running.
   All 3 storage-coherence-0/1/2 pods should be running and ready, as below:

   ```bash
   NAME                   READY   STATUS    RESTARTS   AGE
   storage-coherence-0    1/1     Running   0          4m
   storage-coherence-1    1/1     Running   0          2m   
   storage-coherence-2    1/1     Running   0          2m
   ```
   
1. Ensure both Proxy services are running

   ```bash
   $ kubectl logs storage-coherence-0   -n sample-coherence-ns | grep 'TcpAcceptor now listening' | grep ProxyService   
   2019-05-01 01:55:38.856/8.215 Oracle Coherence GE 12.2.1.4.0 <Info> (thread=Proxy:ProxyService1:TcpAcceptor, member=1): TcpAcceptor now listening for connections on storage-coherence-0.coherence.sample-coherence-ns.svc.cluster.local:20000
   2019-05-01 01:55:38.955/8.313 Oracle Coherence GE 12.2.1.4.0 <Info> (thread=Proxy:ProxyService2:TcpAcceptor, member=1): TcpAcceptor now listening for connections on storage-coherence-0.coherence.sample-coherence-ns.svc.cluster.local:20001
   ```   

1. Port forward the proxy ports on the proxy-tier

   ```bash
   $ kubectl port-forward -n sample-coherence-ns storage-coherence-0 20000:20000
   ```

   ```bash
   $ kubectl port-forward -n sample-coherence-ns storage-coherence-0 20001:20001
   ```

1. Connect via QueryPlus using port 20000

   Issue the following command to run QueryPlus and connect to proxy port 20000

   ```bash
   $ mvn exec:java -Dproxy.port=20000
   ```
   
   Run the following CohQL commands to insert data into the cluster via proxy port 20000.

   ```
   CohQL> insert into 'test' key('key-1') value('value-1');

   CohQL> select key(), value() from 'test';
   Results
   ["key-1", "value-1"]

   CohQL> select count() from 'test';
   Results
   1
   ```
  
   You should see a message showing the connection to 127.0.0.1:20000
   ```bash
   2019-05-01 10:01:08.007/5.684 Oracle Coherence GE 12.2.1.4.0 <D6> (thread=com.tangosol.coherence.dslquery.QueryPlus.main(), member=n/a): Connecting Socket to 127.0.0.1:20000
   ```
   
1. Connect via QueryPlus using port 20001

   Issue the following command to run QueryPlus and connect to proxy port 20001

   ```bash
   $ mvn exec:java -Dproxy.port=20001
   ```

   Run the following CohQL commands to insert data into the cluster via proxy port 20001.

   ```
   CohQL> select key(), value() from 'test';
   Results
   ["key-1", "value-1"]

   CohQL> select count() from 'test';
   Results
   1
   ```
   
   You should see a message showing the connection to 127.0.0.1:20001
   ```bash
   2019-05-01 10:05:18.764/20.659 Oracle Coherence GE 12.2.1.4.0 <D6> (thread=com.tangosol.coherence.dslquery.QueryPlus.main(), member=n/a): Connecting Socket to 127.0.0.1:20001   
   ```
   
## Verifying Grafana Data (If you enabled Prometheus)

Access Grafana using the instructions [here](../../../README.md#access-grafana).

## Verifying Kibana Logs (if you enabled log capture)

Access Kibana using the instructions [here](../../../README.md#access-kibana).

## Uninstalling the Charts

Carry out the following commands to delete the two charts installed in this sample.

```bash
$ helm delete storage --purge
```

Before starting another sample, ensure that all the pods are gone from previous sample.

If you wish to remove the `coherence-operator`, then include it in the `helm delete` command above.
