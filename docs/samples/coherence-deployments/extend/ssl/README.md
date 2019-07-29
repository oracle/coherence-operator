# Enable SSL in Coherence 12.2.1.3.X

This sample shows how to secure Coherence*Extend traffic via 2-way SSL when using the Oracle Coherence
Operator with Coherence 12.2.1.3.x.

Refer to the [Coherence Documentation](https://docs.oracle.com/middleware/12213/coherence/secure/securing-extend-client-connections.htm) for more information about using SSL with Coherence.

> **Note:** If you are using Coherence 12.2.1.4.0, the instructions are slightly different, and
> you should refer to them [here](../12214/).

[Return to Coherence*Extend SSL samples](../) / [Return to Coherence*Extend samples](../../) / [Return to Coherence Deployments samples](../../../) / [Return to samples](../../../README.md#list-of-samples)

## Sample files

* [src/main/docker/Dockerfile](src/main/docker/Dockerfile) - Dockerfile for creating sidecar image from which the configuration will be read at pod startup

* [src/main/resources/certs/keys.sh](src/main/resources/conf/certs/keys.sh) - File for generating keys for this example  

* [src/main/resources/client-cache-config.xml](src/main/resources/client-cache-config.xml) - Cache configuration for the extend client

* [src/main/resources/conf/storage-cache-config.xml](src/main/resources/conf/storage-cache-config.xml) - Cache configuration for storage-tier

## Prerequisites

Ensure that you have installed Coherence Operator by following the instructions [here](../../../README.md#install-the-coherence-operator).

## Installation Steps

1. Change to the `samples/coherence-deployments/extend/ssl/12213` directory. Ensure that you have your maven build environment set for JDK8, and build the project.

   > **Note**: This sample uses self-signed certificates and simple passwords. They are for demonstration
   > purposes only and should **NOT** be used in a production environment.
   > You should use and generate standard certificates with appropriate passwords.

   ```bash
   $ mvn clean install -P docker
   ```

   As a result, the docker image will be built with the cache configuration files and compiled Java classes with the name in the format, `proxy-ssl-sample:${version}`. For example,

   ```bash
   proxy-ssl-sample:1.0.0-SNAPSHOT
   ```

   > **Note:** If you are running against a remote Kubernetes cluster, then you must
   > push the above image to your repository accessible to that cluster. You must also
   > prefix the image name in the `helm` command, as shown below.

2. Install the Coherence cluster.

   Set the following property and ensure that correct cache configuration is used:

   `--set store.cacheConfig=storage-cache-config.xml`


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
      --set userArtifacts.image=proxy-ssl-sample:1.0.0-SNAPSHOT \
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

3. Port forward the proxy port on the proxy-tier.

   ```bash
   $ kubectl port-forward -n sample-coherence-ns storage-coherence-0 20000:20000
   ```

4. Connect via CohQL and run the commands:

   ```bash
   $ mvn exec:java
   ```

   Run the following `CohQL` commands to insert data into the cluster.

   ```sql
   insert into 'test' key('key-1') value('value-1');
   ```

    ```sql
    select key(), value() from 'test';
    Results
    ["key-1", "value-1"]

    select count() from 'test';
    Results
    1
    ```
   You should see a message indicating that the Coherence*Extend client is using  `SSLSocketProvider` with 2-way auth, as shown in the output:

   ```console
   2019-05-06 10:58:49.752/5.105 Oracle Coherence GE 12.2.1.3.2 <D5> (thread=com.tangosol.coherence.dslquery.QueryPlus.main(), member=n/a): instantiated SSLSocketProviderDependencies: SSLSocketProvider(auth=two-way, \
       identity=SunX509/file:conf/certs/groot.jks, trust=SunX509/file:conf/certs/truststore-all.jks)
   ```

   Type `bye` or `CTRL-C` to exit CohQL.  

## Uninstall the Charts

Run the following command to delete both the charts installed in this sample:

```bash
$ helm delete storage --purge
```

Before starting another sample, ensure that all  pods are removed from the previous sample. To remove `coherence-operator`, use the `helm delete` command.
