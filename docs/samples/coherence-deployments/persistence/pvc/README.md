# Use specific persistent volumes

In this sample we explain how to use specific Persistent Volumes (PV) for Coherence when 
using `active` Persistence mode. Local storage is the recommended storage type for achieving best  
performance for `active` Persistence, but this sample can be modified to use any storage class.

*Note*: We only show how to set `store.persistence.*` chart values which apply for `active` persistence (/persistence mount point) only.
It is equally applicable to `store.snapshot.*` chart values with apply to the `/snapshot` volume.
  
[Return to Persistence samples](../) / [Return to Coherence Deployments samples](../../) / [Return to samples](../../../README.md#list-of-samples)

## Sample files

* [local-sc.yaml](local-sc.yaml) - Yaml for creating local storage class


## Prerequisites

Ensure you have already installed the Coherence Operator by using the instructions [here](../../../README.md#install-the-coherence-operator).

## Installation Steps

1. Create a local storage class

   Using the `local-sc.yaml` file, create the local storage class called `localsc`.
   
   ```bash
   $ kubectl create -f local-sc.yaml

    storageclass.storage.k8s.io/localsc created

   $ kubectl get storageclass 
   NAME                 PROVISIONER                    AGE
   hostpath (default)   docker.io/hostpath             26d
   localsc              kubernetes.io/no-provisioner   31s
   ```

1. Create the Persistent Volumes (PV)

   Note: the PV have a label of `coherenceCluster=persistence-cluster` which will
   be used by a nodeSelector to match PV with Coherence clusters.

   ```bash
   $ kubectl create -f mylocal-pv0.yaml -n sample-coherence-ns
   persistentvolume/mylocal-pv0 created   

   $ kubectl create -f mylocal-pv1.yaml -n sample-coherence-ns
   persistentvolume/mylocal-pv2 created   

   $ kubectl create -f mylocal-pv2.yaml -n sample-coherence-ns
   persistentvolume/mylocal-pv2 created

   $ kubectl get pv -n  sample-coherence-ns
   NAME          CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS      CLAIM   STORAGECLASS   REASON   AGE
   mylocal-pv0   2Gi        RWO            Retain           Available           mylocalsc               1m
   mylocal-pv1   2Gi        RWO            Retain           Available           mylocalsc               14s
   mylocal-pv2   2Gi        RWO            Retain           Available           mylocalsc               9s
   ```

   *Note:*  The number of Persistent Volumes created need to be the same as Coherence cluster size.
    For this example whe have assumed a cluster size of three.
    
    
1. Install the Coherence cluster    

   Issue the following to install the cluster with Persistence enabled and select correct storage-class:
   
   * `--set store.persistence.storageClass=mylocalsc` - Specifies the storage-class
   
   * `--set store.persistence.selector.matchLabels.coherenceCluster=persistence-cluster` - Will ensure that PV are only chosen
   where labels match.

   ```bash
   $ helm install \
      --namespace sample-coherence-ns \
      --name storage \
      --set clusterSize=3 \
      --set cluster=persistence-cluster \
      --set imagePullSecrets=sample-coherence-secret \
      --set prometheusoperator.enabled=false \
      --set logCaptureEnabled=false \
      --set store.persistence.enabled=true \
      --set store.persistence.storageClass=mylocalsc \
      --set store.persistence.selector.matchLabels.coherenceCluster=persistence-cluster \
      --version 1.0.0-SNAPSHOT coherence-community/coherence
   ```
   
   Check the pods are running using:
   
   ```bash
   $ kubectl get pods -n sample-coherence-ns
   NAME                                  READY   STATUS    RESTARTS   AGE
   coherence-operator-66f9bb7b75-55msb   1/1     Running   0          23m
   storage-coherence-0                   1/1     Running   0          5m
   storage-coherence-1                   1/1     Running   0          4m
   storage-coherence-2                   1/1     Running   0          3m
   ```
   
   Ensure the PV are matched to the pods:
   
   ```bash
   $ kubectl get pv -n sample-coherence-ns
   NAME          CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS   CLAIM                                                        STORAGECLASS   REASON   AGE
   mylocal-pv0   2Gi        RWO            Retain           Bound    sample-coherence-ns/persistence-volume-storage-coherence-0   mylocalsc               10m
   mylocal-pv1   2Gi        RWO            Retain           Bound    sample-coherence-ns/persistence-volume-storage-coherence-1   mylocalsc               8m
   mylocal-pv2   2Gi        RWO            Retain           Bound    sample-coherence-ns/persistence-volume-storage-coherence-2   mylocalsc               8m
   ```
   
1. Add data to the cluster

   Connect to the Coherence `console` using the following to create a cache.

   ```bash
   $ kubectl exec -it --namespace sample-coherence-ns storage-coherence-0 bash /scripts/startCoherence.sh console
   ```   
   
   At the `Map (?):` prompt, type `cache test`.  This will create a cache in the service `PartitionedCache`.
   
   Use the following to add 100,000 objects of size 1024 bytes, starting at index 0 and using batches of 100.
   
   ```bash
   bulkput 100000 1024 0 100

   Mon Apr 15 07:37:03 GMT 2019: adding 100000 items (starting with #0) each 1024 bytes ...
   Mon Apr 15 07:37:15 GMT 2019: done putting (11578ms, 8878KB/sec, 8637 items/sec)
   ```
   
   At the prompt, type `size` and it should show 100000.
   
   Then type `bye` to exit the `console`.
   
1. Delete the Coherence cluster

   Issue the following to delete the Coherence cluster:
   
   ```bash
   $ helm delete storage --purge
   ```
   
   Ensure the pods are deleted before continuing. This can be done using the following and
   ensuring there are no more pods running with the name `coherence-storage-*`.
   
   ```bash
   $ kubectl get pods -n sample-coherence-ns
   ```   
  
   Issue the following to show that the PVC still exists:   
      
   ```bash
   $ kubectl get pv -n sample-coherence-ns
   NAME          CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS   CLAIM                                                        STORAGECLASS   REASON   AGE
   mylocal-pv0   2Gi        RWO            Retain           Bound    sample-coherence-ns/persistence-volume-storage-coherence-0   mylocalsc               10m
   mylocal-pv1   2Gi        RWO            Retain           Bound    sample-coherence-ns/persistence-volume-storage-coherence-1   mylocalsc               8m
   mylocal-pv2   2Gi        RWO            Retain           Bound    sample-coherence-ns/persistence-volume-storage-coherence-2   mylocalsc               8m
   ```

1. Re-Install the Coherence cluster    

   Issue the following to re-install the cluster:

   ```bash
   $ helm install \
      --namespace sample-coherence-ns \
      --name storage \
      --set clusterSize=3 \
      --set cluster=persistence-cluster \
      --set imagePullSecrets=sample-coherence-secret \
      --set prometheusoperator.enabled=false \
      --set logCaptureEnabled=false \
      --set store.persistence.enabled=true \
      --set store.persistence.storageClass=mylocalsc \
      --set store.persistence.selector.matchLabels.coherenceCluster=persistence-cluster \
      --version 1.0.0-SNAPSHOT coherence-community/coherence
   ```
   
  Wait until all three pods are running before you continue to the next step.
   
1. Ensure the data added previously is still present

   Connect to the Coherence `console` using the following:

   ```bash
   $ kubectl exec -it --namespace sample-coherence-ns storage-coherence-0 bash /scripts/startCoherence.sh console
   ```   
   
   At the `Map (?):` prompt, type `cache test`.  This will create/use a cache in the service `PartitionedCache`.
   
   At the prompt, type `size` and it should show 100000. 
   
   *This shows that the previous data entered has automatically been recovered due the PVC being honoured*.
    
   Then type `bye` to exit the `console`.   
   
## Verifying Grafana Data (If you enabled Prometheus)

Access Grafana using the instructions [here](../../../README.md#access-grafana).

## Verifying Kibana Logs (if you enabled log capture)

Access Kibana using the instructions [here](../../../README.md#access-kibana).

## Uninstalling the Charts

Carry out the following commands to delete the chart and PV's created in this sample.

```bash
$ helm delete storage --purge
```

Once the pods are deleted, issue the following to delete the PVC.:

```bash
$ kubectl delete pvc persistence-volume-storage-coherence-0 persistence-volume-storage-coherence-1 persistence-volume-storage-coherence-2 -n sample-coherence-ns
persistentvolumeclaim "snapshot-volume-storage-coherence-0" deleted
persistentvolumeclaim "snapshot-volume-storage-coherence-1" deleted
persistentvolumeclaim "snapshot-volume-storage-coherence-2" deleted

$ kubectl delete pv mylocal-pv0 mylocal-pv1 mylocal-pv2 
persistentvolume "mylocal-pv0" deleted
persistentvolume "mylocal-pv1" deleted
persistentvolume "mylocal-pv2" deleted

$ kubectl get pvc -n sample-coherence-ns
No resources found.
```

Before starting another sample, ensure that all the pods are gone from previous sample.

If you wish to remove the `coherence-operator`, then include it in the `helm delete` command above.

  