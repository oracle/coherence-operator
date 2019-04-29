# Add application jars/Config to a Coherence deployment

A common scenario for Coherence deployments is to include specific configuration files, such as cache configuration 
operational override files, as well as user classes.

Using the `coherence-operator`, this is achieved by using the `sidecar` approach 
with the docker image having the following directories which are copied to the `coherence` container on startup
and available in the classpath.

[Return to Storage-Disabled clients samples](../) / [Return to Coherence Deployments samples](../../) / [Return to samples](../../../README.md#list-of-samples)

## Sample files

* [src/main/docker/Dockerfile](src/main/docker/Dockerfile) - Dockerfile for creating side-car image from which configuration
  and server side jar will be read from at pod startup

* [src/main/resources/conf/storage-cache-config.xml](src/main/resources/conf/storage-cache-config.xml) - cache config for storage-enabled tier

* [src/main/resources/conf/storage-pof-config.xml](src/main/resources/conf/storage-pof-config.xml) - POF cache config for storage-enabled tier

* [src/main/java/demo/Person.java](src/main/java/demo/Person.java) - domain class for storing Person

* [src/main/java/demo/SampleClient.java](src/main/java/demo/SampleClient.java) - Java client to connect via extend and run entry processor. 


Note if you wish to enable Prometheus or log capture, change the following in the helm installs to `true`. Their default values are false, but they are set to `false` in the samples below for completeness.

* Prometheus: `--set prometheusoperator.enabled=true`

* Log capture: `--set logCaptureEnabled=true`

## Prerequisites

Ensure you have already installed the Coherence Operator by using the instructions [here](../../README.md#install-the-coherence-operator).

## Installation Steps

1. Change to the `samples/coherence-deployments/sidecar` directory and ensure you have your you have your maven build     
   environment set for JDK11 and build the project.

   ```bash
   mvn clean install -P docker
   ```

   **Note:** If you are running against a remote Kubernetes cluster you will need to
   push the above image to your repository accessible to that cluster.

1. The result of the above is the docker image will be built with the cache configuration files
   and compiled Java classes with the name in the format sidecar-sample:${version}.

   You may change the `docker.repo` property in the main [Samples pom.xml](../../../pom.xml).

   For Example:

   ```bash
   sidecar-sample:1.0.0-SNAPSHOT
   ```

   **Note:** If you are running against a remote Kubernetes cluster you will need to
   push the above image to your repository accessible to that cluster.


1. Install the Coherence cluster

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

1. Port-forward Coherence*Extend 

   Issue the following to port-forward default Coherence*Extend port.
   
   ```bash
   $ kubectl port-forward storage-coherence-0 -n sample-coherence-ns 20000:20000
   ```
   
1. Run the SampleClient and connect via Coherence*Extend

   Issue the following to run the `SampleClient.java` class to insert a Person and
   run a server-side Lambda entry processor to change name and address to uppercase. 
   The execution of this entry processor shows that the Coherence cluster is aware of the Person object.
   
   ```bash
   $ mvn exec:java

   ...

   2019-04-16 13:32:35.835/5.091 Oracle Coherence GE 12.2.1.4.0 <Info> (thread=Proxy:TcpInitiator, member=n/a): Loaded POF configuration from "file:/Users/timmiddleton/Documents/CoherenceEngineering/github/samples-project/samples/coherence-deployments/sidecar/target/classes/conf/storage-pof-config.xml"
   2019-04-16 13:32:35.859/5.115 Oracle Coherence GE 12.2.1.4.0 <Info> (thread=Proxy:TcpInitiator, member=n/a): Loaded included POF configuration from "jar:file:/Users/timmiddleton/.m2/repository/com/oracle/coherence/coherence/12.2.1-4-0-73500/coherence-12.2.1-4-0-73500.jar!/coherence-pof-config.xml"

   New Person is: Person{Id=1, Name='Tom Jones', Address='123 Hollywood Ave, California, USA'}
   Person after entry processor is: Person{Id=1, Name='TOM JONES', Address='123 HOLLYWOOD AVE, CALIFORNIA, USA'}
   ```
   
## Verifying Grafana Data (If you enabled Prometheus)

Access Grafana using the instructions [here](../../../README.md#access-grafana).

## Verifying Kibana Logs (if you enabled log capture)

Access Kibana using the instructions [here](../../../README.md#access-kibana).

## Uninstalling the Charts

Carry out the following commands to delete the chart installed in this sample.

```bash
helm delete storage --purge
```

Before starting another sample, ensure that all the pods are gone from previous sample.

If you wish to remove the `coherence-operator`, then include it in the `helm delete` command above.
