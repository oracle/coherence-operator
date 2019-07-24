# Within a Single Kubernetes Cluster

The Federated Caching feature federates cached data asynchronously across multiple geographically
dispersed clusters. Cached data is federated across clusters to provide redundancy,
off-site backup, and multiple points of access for application users in different
geographical locations.

This sample shows how to set up two Federated Coherence clusters within a single Kubernetes cluster.

Although this is not a recommended topology, due to the co-location of the Coherence clusters, this is
included as a sample for how to use Federation.

> **Note**: To set up two Federated Coherence clusters across different Kubernetes clusters, refer to the additional information [here](../across-clusters/README.md).

The setup for this example uses two Coherence clusters in the same Kubernetes cluster, with the following details:

* Primary Cluster
  * Release name: `cluster-1`
  * Cluster name: `PrimaryCluster`
* Secondary Cluster
  * Release name: `cluster-2`
  * Cluster name: `SecondaryCluster`

[Return to Federation samples](../) / [Return to Coherence Deployments samples](../..) / [Return to samples](../../../README.md#list-of-samples)

## Sample files

* [src/main/docker/Dockerfile](src/main/docker/Dockerfile) - Dockerfile for creating sidecar image from which the configuration will be read at pod startup

* [src/main/resources/client-cache-config.xml](src/main/resources/client-cache-config.xml) - Cache configuration for the extend client

* [src/main/resources/conf/cache-config-federation.xml](src/main/resources/conf/cache-config-federation.xml) - Cache configuration for Federation

* [src/main/resources/conf/tangosol-coherence-override-federation.xml](src/main/resources/conf/tangosol-coherence-override-federation.xml) - Override for Federation


## Prerequisites

Ensure that you have installed Oracle Coherence Operator by following the instructions [here](../../../README.md#install-the-coherence-operator).

## Installation Steps

1. Change to the `samples/coherence-deployments/federation/within-cluster` directory. Ensure you have your maven build environment set for JDK8, and build the project.

   ```bash
   $ mvn clean install -P docker
   ```

  This builds the Docker image with the cache configuration files, with the name in the format, `proxy-tier-sample:${version}`. For example,

   ```bash
   federation-within-cluster-sample:1.0.0-SNAPSHOT
   ```

   > **Note:** If you are running against a remote Kubernetes cluster, you must
   > push the above image to your repository accessible to that cluster. You must also
   > prefix the image name in the `helm` command, as shown below.

2. Install the primary Coherence cluster, **PrimaryCluster**.

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

   > **Note:** Ensure that you replace `your-12.2.1.4.0-Coherence-image` with the appropriate Coherence 12.2.1.4.0 Docker image.

    After the installation is complete, get the list of pods by running the following command:

   ```bash
   $ kubectl get pods -n sample-coherence-ns
   ```
   Both the cluster-1-coherence pods should be running and ready, as shown in the output:
   ```console
   NAME                                  READY   STATUS    RESTARTS   AGE
   cluster-1-coherence-0                 1/1     Running   0          24s
   cluster-1-coherence-1                 1/1     Running   0          24s
   coherence-operator-695b9456d5-bzbhl   1/1     Running   0          30m
   ```

3. Port Forward the **PrimaryCluster** Coherence*Extend - Port 20000.

   ```bash
   $ kubectl port-forward --namespace sample-coherence-ns cluster-1-coherence-0  20000:20000
   ```

4. Install the secondary Coherence cluster, **SecondaryCluster**.

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

   After the installation has completed, get the list of pods by running the following command:

   ```bash
   $ kubectl get pods -n sample-coherence-ns
   ```
   All the four cluster-1 and cluster-2 pods should be running and ready, as shown in the output:
   ```console
   NAME                                  READY   STATUS    RESTARTS   AGE
   cluster-1-coherence-0                 1/1     Running   0          4m
   cluster-1-coherence-1                 1/1     Running   0          4m
   cluster-2-coherence-0                 1/1     Running   0          36s
   cluster-2-coherence-1                 1/1     Running   0          36s
   coherence-operator-695b9456d5-bzbhl   1/1     Running   0          34m
   ```


5. Port Forward the **SecondaryCluster** Coherence*Extend - Port 20001.

   ```bash
   $ kubectl port-forward --namespace sample-coherence-ns cluster-2-coherence-0  20001:20000
   ```

6. Use QueryPlus to connect to each cluster.

   Run QueryPlus against the **PrimaryCluster** by using the following command:

   ```bash
   $ mvn exec:java -Dproxy.port=20000
   ```

   Open another terminal and run QueryPlus against the **SecondaryCluster**:

   ```bash
   $ mvn exec:java -Dproxy.port=20001
   ```

7. Insert data into the **PrimaryCluster**.

   Run the following `CohQL` command to insert data into the **PrimaryCluster**:

   ```sql
   insert into 'test' key('key-1') value('value-1');

   select key(), value() from 'test';
   Results
   ["key-1", "value-1"]

   select count() from 'test';
   Results
   1
   ```

8. Confirm that the data has been federated to the **SecondaryCluster**.

   Run the following `CohQL` command to insert data into the **SecondaryCluster**:

   ```sql
   select key(), value() from 'test';
   Results
   ["key-1", "value-1"]
   ```

9. Insert data into the **SecondaryCluster**

   Run the following `CohQL` command to insert data into the **SecondaryCluster**:

   ```sql
   insert into 'test' key('key-2') value('value-2');

   select key(), value() from 'test';
   Results
   ["key-1", "value-1"]
   ["key-2", "value-2"]
   ```

10. Confirm that the data has been federated to the **PrimaryCluster**.

   Run the following `CohQL` command to insert data into the **PrimaryCluster**:

   ```sql
   select key(), value() from 'test';
   Results
   ["key-1", "value-1"]
   ["key-2", "value-2"]
   ```

## Uninstall the Charts

Run the following command to delete both the charts installed in this sample:

```bash
$ helm delete cluster-1 cluster-2 --purge
```

Before starting another sample, ensure that all  pods are removed from the previous sample. To remove `coherence-operator`, use the `helm delete` command.
