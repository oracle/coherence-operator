# Across across separate Kubernetes clusters

The Federated Caching feature federates cache data asynchronously across multiple geographically 
dispersed clusters. Cached data is federated across clusters to provide redundancy, 
off-site backup, and multiple points of access for application users in different 
geographical locations.

This sample shows how to setup two Federated Coherence clusters across separate Kubernetes clusters.

> **Note**: To setup two Federated Coherence Clusters within a single Kubernetes cluster, please 
> see additional information [here](../across-clusters/README.md).

[Return to Federation samples](../) / [Return to Coherence Deployments samples](../../) / [Return to samples](../../../README.md#list-of-samples)

## Installation Steps

You need to follow the steps described [here](../across-clusters/README.md) (Federation within a single Kubernetes
cluster) but you will need to change the following:

1. Expose Port 40000 on each of the clusters to a visible IP or load balancer

   You will need to expose the cluster port `40000` (or an alternative) via an IP directly or load balancer
   on each of the clusters.  This will allow the clusters to communicate.
   
   See the [Kubernetes Documentation](https://kubernetes.io/docs/tutorials/stateless-application/expose-external-ip-address/)
   for more information.
   
1. Ignore any port-forward commands
   
1. Install the Coherence Clusters

   Once you have an IP/PORT for each cluster, you must change the the following in the
   `--set store.javaOpts` in the [3rd and 5th installation steps](../within-cluster/README.md#installation-steps):
   
   * `priamry.cluster.host` - primary cluster external IP
   
   * `primary.cluster.port` - primary cluster external port (default to 40000) 
   
   * `secondary.cluster.host` - secondary cluster external IP
   
   * `secondary.cluster.port` - secondary cluster external port (default to 40000) 

   Eg.
   
   ```bash
   --set store.javaOpts="-Dprimary.cluster=PrimaryCluster -Dprimary.cluster.port=40000 -Dprimary.cluster.host=PRIMARY-CLUSTER-IP -Dsecondary.cluster=SecondaryCluster -Dsecondary.cluster.port=40000 -Dsecondary.cluster.host=SECONDARY-CLUSTER-IP"  \
   ```

