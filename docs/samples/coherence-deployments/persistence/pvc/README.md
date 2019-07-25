# Use Specific Persistent Volumes

This sample shows how to use specific persistent volumes (PV) for Coherence when
using `active` persistence mode. Local storage is the recommended storage type for achieving the best  performance for `active` persistence, but this sample can be modified to use any storage class.

> **Note:** We only show how to set `store.persistence.*` chart values which apply for `active` persistence (/persistence mount point) only.
> It is equally applicable to the `store.snapshot.*` chart values that apply to the `/snapshot` volume.

[Return to Persistence samples](../) / [Return to Coherence Deployments samples](../../) / [Return to samples](../../../README.md#list-of-samples)

## Sample files

* [local-sc.yaml](local-sc.yaml) - YAML for creating local storage class

* [mylocal-pv0.yaml](mylocal-pv0.yaml) - YAML for creating persistent volume mylocal-pv0

* [mylocal-pv1.yaml](mylocal-pv2.yaml) - YAML for creating persistent volume mylocal-pv0

* [mylocal-pv2.yaml](mylocal-pv2.yaml) - YAML for creating persistent volume mylocal-pv0

## Prerequisites

Ensure that you have installed Coherence Operator by following the instructions [here](../../../README.md#install-the-coherence-operator).

## Installation Steps

1. Create a local storage class.

   Using the `local-sc.yaml` file, create a local storage class called `localsc`.

   ```bash
   $ kubectl create -f local-sc.yaml
   ```
   ```console
   storageclass.storage.k8s.io/localsc created
   ```
   Confirm the creation of the storage class:
   ```bash
   $ kubectl get storageclass
   ```
   ```console
   NAME                 PROVISIONER                    AGE
   hostpath (default)   docker.io/hostpath             26d
   localsc              kubernetes.io/no-provisioner   31s
   ```

2. Create persistent volumes (PV).

   > **Note:** The PV has the label,  `coherenceCluster=persistence-cluster`, which is
   >  used by a nodeSelector to match PV with Coherence clusters.

   ```bash
   $ kubectl create -f mylocal-pv0.yaml -n sample-coherence-ns
   persistentvolume/mylocal-pv0 created   

   $ kubectl create -f mylocal-pv1.yaml -n sample-coherence-ns
   persistentvolume/mylocal-pv2 created   

   $ kubectl create -f mylocal-pv2.yaml -n sample-coherence-ns
   persistentvolume/mylocal-pv2 created
   ```
   Confirm the creation of persistent volumes by running the `kubectl` command:
   ```bash
   $ kubectl get pv -n  sample-coherence-ns
   ```
   ```console
   NAME          CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS      CLAIM   STORAGECLASS   REASON   AGE
   mylocal-pv0   2Gi        RWO            Retain           Available           mylocalsc               1m
   mylocal-pv1   2Gi        RWO            Retain           Available           mylocalsc               14s
   mylocal-pv2   2Gi        RWO            Retain           Available           mylocalsc               9s
   ```

   > **Note:** The number of persistent volumes created must be the same as the Coherence cluster size.
   > For this example we have assumed a cluster size of three.


3. Install the Coherence cluster.    

   Run the following command to install the cluster with persistence enabled, and select the correct storage-class:

   * `--set store.persistence.storageClass=mylocalsc` - Specifies the storage-class

   * `--set store.persistence.selector.matchLabels.coherenceCluster=persistence-cluster` - Ensures that the persistent volumes are chosen only where labels match.

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
      coherence/coherence
   ```

   Check whether the pods are running:

   ```bash
   $ kubectl get pods -n sample-coherence-ns
   ```
   ```console
   NAME                                  READY   STATUS    RESTARTS   AGE
   coherence-operator-66f9bb7b75-55msb   1/1     Running   0          23m
   storage-coherence-0                   1/1     Running   0          5m
   storage-coherence-1                   1/1     Running   0          4m
   storage-coherence-2                   1/1     Running   0          3m
   ```

   Ensure that the persistent volumes match the pods:

   ```bash
   $ kubectl get pv -n sample-coherence-ns
   ```
   ```console
   NAME          CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS   CLAIM                                                        STORAGECLASS   REASON   AGE
   mylocal-pv0   2Gi        RWO            Retain           Bound    sample-coherence-ns/persistence-volume-storage-coherence-0   mylocalsc               10m
   mylocal-pv1   2Gi        RWO            Retain           Bound    sample-coherence-ns/persistence-volume-storage-coherence-1   mylocalsc               8m
   mylocal-pv2   2Gi        RWO            Retain           Bound    sample-coherence-ns/persistence-volume-storage-coherence-2   mylocalsc               8m
   ```

4. Add data to the cluster.

   a. Connect to the Coherence console using the following command to create a cache:

   ```bash
   $ kubectl exec -it --namespace sample-coherence-ns storage-coherence-0 bash /scripts/startCoherence.sh console
   ```   

   b. At the `Map (?):` prompt, type `cache test`.  This creates a cache in the service, `PartitionedCache`.

   c. Use the following to add 100,000 objects of size 1024 bytes, starting at index 0, and using batches of 100.

   ```bash
   bulkput 100000 1024 0 100

   Mon Apr 15 07:37:03 GMT 2019: adding 100000 items (starting with #0) each 1024 bytes ...
   Mon Apr 15 07:37:15 GMT 2019: done putting (11578ms, 8878KB/sec, 8637 items/sec)
   ```

   At the prompt, type `size` and it should show 100000.

   d. Create a snapshot of the `PartitionedCache` service which contains the cache `test`. This is for later use.

   ```bash
   snapshot create test-snapshot
   Issuing createSnapshot for service PartitionedCache and snapshot empty-service
   Success
   ```

   ```bash
   snapshot list
    Snapshots for service PartitionedCache
       test-snapshot
   ```

   To exit the console, type `bye`.

5. Delete the Coherence cluster by running the following command:
    ```bash
    $ helm delete storage --purge
    ```

   Before continuing, ensure that the pods are deleted. This can be achieved using the following `kubectl` command:

   ```bash
   $ kubectl get pods -n sample-coherence-ns
   ```  
   Ensure that there are no pods running with the name `coherence-storage-*`

   Run the following to ensure that the PVC still exists:   

   ```bash
   $ kubectl get pv -n sample-coherence-ns
   NAME          CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS   CLAIM                                                        STORAGECLASS   REASON   AGE
   mylocal-pv0   2Gi        RWO            Retain           Bound    sample-coherence-ns/persistence-volume-storage-coherence-0   mylocalsc               10m
   mylocal-pv1   2Gi        RWO            Retain           Bound    sample-coherence-ns/persistence-volume-storage-coherence-1   mylocalsc               8m
   mylocal-pv2   2Gi        RWO            Retain           Bound    sample-coherence-ns/persistence-volume-storage-coherence-2   mylocalsc               8m
   ```

6. Reinstall the Coherence cluster:   

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
      coherence/coherence
   ```

  Wait until all three pods are running before you continue to the next step.

7. Ensure that the data added previously still exists.

   a. Connect to the Coherence console using the following command:

   ```bash
   $ kubectl exec -it --namespace sample-coherence-ns storage-coherence-0 bash /scripts/startCoherence.sh console
   ```   

   b. At the `Map (?):` prompt, type `cache test`.  This create/use a cache in the service `PartitionedCache`.

   c. At the prompt, type `size` and it should show 100000.

   > This shows that the previous data entered has automatically been recovered as the PVC was honoured.

   d. Clear the cache using `clear` command and confirm the cache size is zero.

   Recover the `test-snapshot` using:

   ```bash
   snapshot recover test-snapshot
   ```

   The size of the cache should now be 100000.

   To exit the console, type `bye`.   

## Uninstall the Charts

Run the following command to delete the chart and persistent volumes installed in this sample:

```bash
$ helm delete storage --purge
```

Once the pods are deleted, run the following command to delete the PVC.:

```bash
$ kubectl delete pvc persistence-volume-storage-coherence-0 persistence-volume-storage-coherence-1 \
                     persistence-volume-storage-coherence-2 -n sample-coherence-ns
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

Before starting another sample, ensure that all  pods are removed from the previous sample. To remove `coherence-operator`, use the `helm delete` command.
