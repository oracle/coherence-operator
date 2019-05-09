# Enable SSL in Coherence 12.2.1.3.X

This sample shows how to secure Coherence*Extend traffic via 2-way SSL when using the Coherence 
Operator with Coherence 12.2.1.3.X.

Please see the [Coherence Documentation](https://docs.oracle.com/middleware/12213/coherence/secure/securing-extend-client-connections.htm)
for more information on using SSL with Coherence.

> Note: If you are using Coherence 12.2.1.4, the instructions are slightly different and
> you should refer to them [here](../12214/).

[Return to Coherence*Extend SSL samples](../) / [Return to Coherence*Extend samples](../../) / [Return to Coherence Deployments samples](../../../) / [Return to samples](../../../../README.md#list-of-samples)

## Sample files

* [src/main/docker/Dockerfile](src/main/docker/Dockerfile) - Dockerfile for creating side-car image from which configuration
  will be read from at pod startup
  
* [src/main/resources/certs/keys.sh](src/main/resources/conf/certs/keys.sh) - File for generating keys for this example  

* [src/main/resources/client-cache-config.xml](src/main/resources/client-cache-config.xml) - cache config for extend client

* [src/main/resources/conf/storage-cache-config.xml](src/main/resources/conf/storage-cache-config.xml) - cache config for storage-tier

## Prerequisites

Ensure you have already installed the Coherence Operator by using the instructions [here](../../../README.md#install-the-coherence-operator).

## Installation Steps

1. Change to the `samples/coherence-deployments/extend/ssl/12213` directory and ensure you have your 
   maven build environment set for JDK11 and build the project.
   
   > Note: This sample uses self-signed certificates and simple passwords. They are for sample
   > purposes only and should **NOT** use these in a production environment.
   > You should use and generate proper certificates with appropriate passwords.

   ```bash
   $ mvn clean install -P docker
   ```

1. The result of the above is the docker image will be built with the cache configuration files
   with the name in the format proxy-ssl-sample-12213:${version}.

   For Example:

   ```bash
   proxy-ssl-sample-12213:1.0.0-SNAPSHOT
   ```

   > Note: If you are running against a remote Kubernetes cluster you will need to
   > push the above image to your repository accessible to that cluster. You will also need to 
   > prefix the image name in your `helm` command below.
   
1. Install the Coherence cluster

   The following are set:
   
   * `--set store.cacheConfig=storage-cache-config.xml` - ensure correct cache config is used
   
   ```bash
   $ helm install \
      --namespace sample-coherence-ns \
      --name storage \
      --set clusterSize=3 \
      --set cluster=proxy-ssl-cluster \
      --set imagePullSecrets=sample-coherence-secret \
      --set store.cacheConfig=storage-cache-config.xml \
      --set prometheusoperator.enabled=false \
      --set logCaptureEnabled=false \
      --set userArtifacts.image=proxy-ssl-sample-12213:1.0.0-SNAPSHOT \
      coherence-community/coherence
   ```
   
   Use `kubectl get pods -n sample-coherence-ns` to ensure that all pods are running.
   All 3 storage-coherence-0/1/2 pods should be running and ready, as below:

   ```bash
   NAME                   READY   STATUS    RESTARTS   AGE
   storage-coherence-0    1/1     Running   0          4m
   storage-coherence-1    1/1     Running   0          2m   
   storage-coherence-2    1/1     Running   0          2m
   ```   

1. Port forward the proxy port on the proxy-tier

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
   ```
   
   ```bash
   2019-05-06 10:58:49.752/5.105 Oracle Coherence GE 12.2.1.3.2 <D5> (thread=com.tangosol.coherence.dslquery.QueryPlus.main(), member=n/a): instantiated SSLSocketProviderDependencies: SSLSocketProvider(auth=two-way, \
       identity=SunX509/file:conf/certs/groot.jks, trust=SunX509/file:conf/certs/truststore-all.jks)
   ```
   
   You should notice above that there should be a message indicating the Coherence*Extend client is using a `SSLSocketProvider` with 2-way auth. 

   ```sql
   select key(), value() from 'test';
   Results
   ["key-1", "value-1"]

   select count() from 'test';
   Results
   1
   ```
   
   Type `bye` or CTRL-C to exit QueryPlus.  

## Uninstalling the Charts

Carry out the following commands to delete the two charts installed in this sample.

```bash
$ helm delete storage --purge
```

Before starting another sample, ensure that all the pods are gone from previous sample.

If you wish to remove the `coherence-operator`, then include it in the `helm delete` command above.
