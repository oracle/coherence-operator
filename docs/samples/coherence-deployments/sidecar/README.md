# Add Application Jars/Config to a Coherence Deployment

A common scenario for Coherence deployments is to include specific configuration files, such as cache configuration
operational override files, as well as user classes.

This can be achieved with `coherence-operator` by using the `sidecar` approach. You must ensure that the docker image has the following directories that are copied to the `coherence` container on startup
and are available in the classpath.

This sample shows how to create a sidecar image that contains cache configuration, POF configuration and
associated portable `Person` object. A client, `SampleClient` runs and uses Lambda to modify data
on the server side.

[Return to Coherence Deployments samples](../) / [Return to samples](../../README.md#list-of-samples)

## Sample files

* [src/main/docker/Dockerfile](src/main/docker/Dockerfile) - Dockerfile for creating a sidecar image from which the configuration and server side JAR will be read at pod startup

* [src/main/resources/conf/storage-cache-config.xml](src/main/resources/conf/storage-cache-config.xml) - Cache configuration for storage-enabled tier

* [src/main/resources/conf/storage-pof-config.xml](src/main/resources/conf/storage-pof-config.xml) - POF cache configuration for storage-enabled tier

* [src/main/java/com/oracle/coherence/examples/Person.java](src/main/java/com/oracle/coherence/examples/Person.java) - Domain class for storing `Person`

* [src/main/java/com/oracle/coherence/examples/SampleClient.java](src/main/java/com/oracle/coherence/examples/SampleClient.java) - Java client to connect via extend and run entry processor.

## Prerequisites

Ensure that you have installed Oracle Coherence Operator by following the instructions [here](../../README.md#install-the-coherence-operator).

## Installation Steps

1. Change to the `samples/coherence-deployments/sidecar` directory. Ensure that you have your maven build environment set for JDK8, and build the project.

   ```bash
   $ mvn clean install -P docker
   ```

  As a result, the Docker image will be built with the cache configuration files and compiled Java classes with the name in the format `sidecar-sample:${version}`. For example,
  ```bash
  sidecar-sample:1.0.0-SNAPSHOT`.
  ```

   > **Note:** If you are running against a remote Kubernetes cluster, you must push the above image to your repository accessible to that cluster. You must also prefix the image name in the `helm` command as shown below.

2. Install the Coherence cluster.

   ```bash
   $ helm install \
      --namespace sample-coherence-ns \
      --name storage \
      --set clusterSize=3 \
      --set cluster=sidecar-cluster \
      --set imagePullSecrets=sample-coherence-secret \
      --set store.cacheConfig=storage-cache-config.xml \
      --set store.pof.config=storage-pof-config.xml \
      --set prometheusoperator.enabled=false \
      --set logCaptureEnabled=false \
      --set userArtifacts.image=sidecar-sample:1.0.0-SNAPSHOT \
      coherence/coherence
   ```

   Once the installation is complete, get the list of pods by using the `kubectl` command:

   ```bash
   $ kubectl get pods -n sample-coherence-ns
   ```
   All the three storage-coherence pods should be running and ready, as shown in the output:

   ```console
   NAME                 READY   STATUS    RESTARTS   AGE
   storage-coherence-0  1/1     Running   0          4m
   storage-coherence-1  1/1     Running   0          2m
   storage-coherence-2  1/1     Running   0          1m
   ```


3. Port-forward Coherence*Extend.

   Run the following `kubectl` command to port-forward the default Coherence*Extend port:

   ```bash
   $ kubectl port-forward storage-coherence-0 -n sample-coherence-ns 20000:20000
   ```

4. Run the `SampleClient.java` class and connect via Coherence*Extend.

   Run the `SampleClient.java` class to insert a `Person` and
   run a server-side Lambda entry processor to change the name and address to uppercase.
   The execution of this entry processor shows that the Coherence cluster
   is aware of the `Person` object as specified in  `userArtifacts.image`.

   ```bash
   $ mvn exec:java
   ```
```console
   2019-04-16 13:32:35.835/5.091 Oracle Coherence GE 12.2.1.4.0 <Info> (thread=Proxy:TcpInitiator, member=n/a): Loaded POF configuration from "file:/Users/timmiddleton/Documents/CoherenceEngineering/github/samples-project/samples/coherence-deployments/sidecar/target/classes/conf/storage-pof-config.xml"
   2019-04-16 13:32:35.859/5.115 Oracle Coherence GE 12.2.1.4.0 <Info> (thread=Proxy:TcpInitiator, member=n/a): Loaded included POF configuration from "jar:file:/Users/timmiddleton/.m2/repository/com/oracle/coherence/coherence/12.2.1-4-0-73500/coherence-12.2.1-4-0-73500.jar!/coherence-pof-config.xml"

   New Person is: Person{Id=1, Name='Tom Jones', Address='123 Hollywood Ave, California, USA'}
   Person after entry processor is: Person{Id=1, Name='TOM JONES', Address='123 HOLLYWOOD AVE, CALIFORNIA, USA'}
   ```

## Uninstall the Charts

Run the following command to delete the chart installed in this sample:

```bash
$ helm delete storage --purge
```

Before starting another sample, ensure that all  pods are removed from the previous sample. To remove `coherence-operator`, use the `helm delete` command.
