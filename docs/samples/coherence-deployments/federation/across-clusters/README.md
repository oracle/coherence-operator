# Across Separate Kubernetes Clusters

The Federated Caching feature federates cached data asynchronously across multiple geographically
dispersed clusters. Cached data is federated across clusters to provide redundancy,
off-site backup, and multiple points of access for application users in different
geographical locations.

This sample shows how to set up two Federated Coherence clusters across separate Kubernetes clusters.

> **Note**: To set up two Federated Coherence clusters within a single Kubernetes cluster, refer to the additional information [here](../across-clusters/README.md).

[Return to Federation samples](../) / [Return to Coherence Deployments samples](../../) / [Return to samples](../../../README.md#list-of-samples)

## Installation Steps

Follow the steps described in [Federation Within a Single Kubernetes
Cluster](../within-cluster/README.md), with the following changes:

1. Expose port 40000 on each cluster to a visible IP or load balancer.

   You must expose the cluster port 40000 (or an alternative) via an IP directly or via load balancer
   on each cluster. This allows the clusters to communicate.

   Refer to the [Kubernetes documentation](https://kubernetes.io/docs/tutorials/stateless-application/expose-external-ip-address/)
   for more information.

1. Ignore any port-forward commands.

1. Install the Coherence clusters.

   Once you have an IP/PORT for each cluster, you must change the following in the
   `--set store.javaOpts` option in the [installation steps 2 and 4](../within-cluster/README.md#installation-steps):

   * `priamry.cluster.host` - The external IP of the primary cluster

   * `primary.cluster.port` - The external port (default to 40000) of the primary cluster

   * `secondary.cluster.host` - The external IP of the secondary cluster

   * `secondary.cluster.port` - The external port (default to 40000) of the secondary cluster

   Eg.

   ```bash
   --set store.javaOpts="-Dprimary.cluster=PrimaryCluster -Dprimary.cluster.port=40000 -Dprimary.cluster.host=PRIMARY-CLUSTER-IP -Dsecondary.cluster=SecondaryCluster -Dsecondary.cluster.port=40000 -Dsecondary.cluster.host=SECONDARY-CLUSTER-IP"  \
   ```
