# Within a single Kubernetes cluster

The federated caching feature federates cache data asynchronously across multiple geographically 
dispersed clusters. Cached data is federated across clusters to provide redundancy, 
off-site backup, and multiple points of access for application users in different 
geographical locations.

This sample shows how to setup a two Federated Coherence clusters within a single Kubernetes cluster.

Although this is not a recommended topology, due to the co-location of the Coherence clusters, this is 
included as a sample of how to use Federation.

The setup for this example uses 2 Coherence clusters in the same Kubernetes cluster.

* Primary Cluster
  * Release name: cluster-1
  * Cluster name: PrimaryCluster
* Secondary Cluster
  * Release name: cluster-2
  * Cluster name: SecondaryCluster

[Return to Federation samples](../) / [Return to Coherence Deployments samples](../..) / [Return to samples](../../../README.md#list-of-samples)

## Sample files

* [src/main/docker/Dockerfile](src/main/docker/Dockerfile) - Dockerfile for creating side-car image from which configuration
  will be read from at pod startup

* [src/main/resources/client-cache-config.xml](src/main/resources/client-cache-config.xml) - cache config for extend client

* [src/main/resources/conf/cache-config-federation.xml](src/main/resources/conf/cache-config-federation.xml) - cache config for Federation

* [src/main/resources/conf/tangosol-coherence-override-federation.xml](src/main/resources/conf/tangosol-coherence-override-federation.xml) - override for Federation


## Prerequisites

Ensure you have already installed the Coherence Operator by using the instructions [here](../../../README.md#install-the-coherence-operator).

## Installation Steps

1. Change to the `samples/coherence-deployments/federation/within-cluster` directory and ensure you have your maven build environment set for JDK8 and build the project.

   ```bash
   $ mvn clean install -P docker
   ```

1. The result of the above is the docker image will be built with the cache configuration files
   with the name in the format proxy-tier-sample:${version}.

   For Example:

   ```bash
   federation-within-cluster-sample:1.0.0-SNAPSHOT
   ```

   > **Note:** If you are running against a remote Kubernetes cluster you will need to
   > push the above image to your repository accessible to that cluster. You will also need to 
   > prefix the image name in your `helm` command below.

1. Install the **PrimaryCluster** 

   ```bash
   $ helm install \
      --namespace sample-coherence-ns \
      --name cluster-1 \
      --set clusterSize=2 \
      --set cluster=PrimaryCluster \
      --set imagePullSecrets=sample-coherence-secret \
      --set store.cacheConfig=cache-config-federation.xml \
      --set store.overrideConfig=tangosol-coherence-override-federation.xml \
      --set store.javaOpts="-Dprimary.cluster=PrimaryCluster -Dprimary.cluster.port=40000 -Dprimary.cluster.host=cluster-1-coherence-headless -Dsecondary.cluster=SecondaryCluster -Dsecondary.cluster.port=40000 -Dsecondary.cluster.host=cluster-2-coherence-headless"  \
      --set store.ports.federation=40000 \
      --set userArtifacts.image=federation-within-cluster-sample:1.0.0-SNAPSHOT \
      --set coherence.image=your-12.2.1.4.0-Coherence-image \
      coherence/coherence
   ```  
   
   > **Note:** Ensure you replace `your-12.2.1.4.0-Coherence-image` with the proper Coherence 12.2.1.4.0 Docker image.
   
   Once the install has completed, issue the following command to list the pods:

   ```bash
   $ kubectl get pods -n sample-coherence-ns
   NAME                                  READY   STATUS    RESTARTS   AGE
   cluster-1-coherence-0                 1/1     Running   0          24s
   cluster-1-coherence-1                 1/1     Running   0          24s
   coherence-operator-695b9456d5-bzbhl   1/1     Running   0          30m
   ```
   
   All 2 cluster-1-coherence-0/2 pods should be running and ready, as above.

1. Port Forward the Primary Cluster Coherence*Extend - Port **20000**

   ```bash
   $ kubectl port-forward --namespace sample-coherence-ns cluster-1-coherence-0  20000:20000
   ```

1. Install the **SecondaryCluster**

   ```bash
   $ helm install \
      --namespace sample-coherence-ns \
      --name cluster-2 \
      --set clusterSize=2 \
      --set cluster=SecondaryCluster \
      --set imagePullSecrets=sample-coherence-secret \
      --set store.cacheConfig=cache-config-federation.xml \
      --set store.overrideConfig=tangosol-coherence-override-federation.xml \
      --set store.javaOpts="-Dprimary.cluster=PrimaryCluster -Dprimary.cluster.port=40000 -Dprimary.cluster.host=cluster-1-coherence-headless -Dsecondary.cluster=SecondaryCluster -Dsecondary.cluster.port=40000 -Dsecondary.cluster.host=cluster-2-coherence-headless"  \
      --set store.ports.federation=40000 \
      --set userArtifacts.image=federation-within-cluster-sample:1.0.0-SNAPSHOT \
      --set coherence.image=your-12.2.1.4.0-Coherence-image \
      coherence/coherence
   ```   
   
   Once the install has completed, issue the following command to list the pods:

   ```bash
   $ kubectl get pods -n sample-coherence-ns
   NAME                                  READY   STATUS    RESTARTS   AGE
   cluster-1-coherence-0                 1/1     Running   0          4m
   cluster-1-coherence-1                 1/1     Running   0          4m
   cluster-2-coherence-0                 1/1     Running   0          36s
   cluster-2-coherence-1                 1/1     Running   0          36s
   coherence-operator-695b9456d5-bzbhl   1/1     Running   0          34m
   ```
   
   All 4 cluster-1 and cluster-2 pods should be running and ready, as above.

1. Port Forward the Secondary Cluster Coherence*Extend - Port **20001**

   ```bash
   $ kubectl port-forward --namespace sample-coherence-ns cluster-2-coherence-0  20001:20000
   ```
   
1. Connect via QueryPlus to each of the clusters

   Issue the following command to run QueryPlus against the **PrimaryCluster**:

   ```bash
   $ mvn exec:java -Dproxy-port=20000
   ``` 
   
   Open a second terminal and issue the following command to run QueryPlus against the **SecondaryCluster**:

   ```bash
   $ mvn exec:java -Dproxy-port=20001
   ``` 
   
1. Insert data into the **PrimaryCluster**

   Run the following CohQL command to insert data into the **PrimaryCluster**:
   
   ```sql
   insert into 'test' key('key-1') value('value-1');

   select key(), value() from 'test';
   Results
   ["key-1", "value-1"]

   select count() from 'test';
   Results
   1
   ```
   
1. Confirm the data has been federated to **SecondaryCluster**

   Run the following CohQL command to insert data into the **SecondaryCluster**:
    
   ```sql
   select key(), value() from 'test';
   Results
   ["key-1", "value-1"]
   ```
   
1. Insert data into the **SecondaryCluster**

   Run the following CohQL command to insert data into the **SecondaryCluster**:
   
   ```sql
   insert into 'test' key('key-2') value('value-2');

   select key(), value() from 'test';
   Results
   ["key-1", "value-1"]
   ["key-2", "value-2"]
   ```
   
1. Confirm the data has been federated to **PrimaryCluster**

   Run the following CohQL command to insert data into the **PrimaryCluster**:
    
   ```sql
   select key(), value() from 'test';
   Results
   ["key-1", "value-1"]
   ["key-2", "value-2"]
   ``` 
   
## Uninstalling the Charts

Carry out the following commands to delete the two charts installed in this sample.

```bash
$ helm delete cluster-1 cluster-2 --purge
```

Before starting another sample, ensure that all the pods are gone from previous sample.

If you wish to remove the `coherence-operator`, then include it in the `helm delete` command above.   
