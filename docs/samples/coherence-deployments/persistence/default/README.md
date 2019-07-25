# Use the Default Persistent Volume Claim

By default, when you enable Coherence Persistence, the required infrastructure in terms
of persistent volumes (PV) and persistent volume claims (PVC) is set up automatically. Also, the persistence-mode
is set to `active`. This allows the Coherence cluster to be restarted and the data to be retained.

This sample shows how to enable Persistence with all the defaults
under `store.persistence` in the `coherence` chart, `values.yaml`. Refer to [this sample](../pvc/README.md)
for more details about setting other values such as `storageClasses`.

[Return to Persistence samples](../) / [Return to Coherence Deployments samples](../../) / [Return to samples](../../../README.md#list-of-samples)

## Prerequisites

Ensure that you have installed Coherence Operator by following the instructions [here](../../../README.md#install-the-coherence-operator).

## Installation Steps

1. Install the Coherence cluster.

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
      --set store.snapshot.enabled=true \
      coherence/coherence
   ```

   You may also change the size of the default directories:

   * `/persistence` default 2Gi - `--set store.persistence.size=10Gi`

   * `/snapshot` default 2Gi - `--set store.snapshot.size-10Gi`                  

2. Ensure that the pods are running, by using the following command:

   ```bash
   $ kubectl get pods  -n sample-coherence-ns
   ```
    When all the three pods are in the `Running` state, proceed to check the PVC.
   ```console
   NAME                                  READY   STATUS    RESTARTS   AGE
   coherence-operator-66f9bb7b75-hqk4l   1/1     Running   0          13m
   storage-coherence-0                   1/1     Running   0          3m
   storage-coherence-1                   1/1     Running   0          2m
   storage-coherence-2                   0/1     Running   0          44s
   ```


3. Run the following command to check the PVC created:

   ```bash
   $ kubectl get pvc -n sample-coherence-ns
   ```
   ```console
   NAME                                     STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS   AGE
   persistence-volume-storage-coherence-0   Bound    pvc-a3fa6ce3-6588-11e9-bad6-025000000001   2Gi        RWO            hostpath       44m
   persistence-volume-storage-coherence-1   Bound    pvc-d2f732e7-6588-11e9-bad6-025000000001   2Gi        RWO            hostpath       43m
   persistence-volume-storage-coherence-2   Bound    pvc-fc175ea1-6588-11e9-bad6-025000000001   2Gi        RWO            hostpath       41m
   snapshot-volume-storage-coherence-0      Bound    pvc-a3fb2172-6588-11e9-bad6-025000000001   2Gi        RWO            hostpath       44m
   snapshot-volume-storage-coherence-1      Bound    pvc-d2f89ce9-6588-11e9-bad6-025000000001   2Gi        RWO            hostpath       43m
   snapshot-volume-storage-coherence-2      Bound    pvc-fc13ae4b-6588-11e9-bad6-025000000001   2Gi        RWO            hostpath       41m
   ```

4. Add data to the cluster.

   a. Connect to the Coherence console using the following command to create a cache:

   ```bash
   $ kubectl exec -it --namespace sample-coherence-ns storage-coherence-0 bash /scripts/startCoherence.sh console
   ```   

   b. At the `Map (?):` prompt, type `cache test`.  This creates a cache in the service, `PartitionedCache`.

   c. Use the following to add 50,000 objects of size 1024 bytes, starting at index 0 and using batches of 100:

   ```bash
   bulkput 50000 1024 0 100

   Mon Apr 15 07:37:03 GMT 2019: adding 500000 items (starting with #0) each 1024 bytes ...
   Mon Apr 15 07:37:15 GMT 2019: done putting (11578ms, 8878KB/sec, 8637 items/sec)
   ```

  At the prompt, type `size` and it should show 50000.

  d. Create a snapshot of the `PartitionedCache` service which contains the cache `test`. This is for later use.

   ```bash
   snapshot create test-snapshot
   Issuing createSnapshot for service PartitionedCache and snapshot test-snapshot
   Success
   ```

   ```bash
   snapshot list
    Snapshots for service PartitionedCache
       test-snapshot
   ```

   To exit the console, type `bye` or CTRL-C.

5. Delete the Coherence cluster by using `helm delete`:

    ```bash
    $ helm delete storage --purge
    ```

   Before continuing, ensure that the pods are deleted using the following command:

   ```bash
   $ kubectl get pods -n sample-coherence-ns
   ```   
  Ensure that there are no pods running with the name `coherence-storage-*`.

   Run the following command to ensure that the PVC still exists:

   ```bash
   $ kubectl get pvc -n sample-coherence-ns
   ```
   ```console
   NAME                                     STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS   AGE
   persistence-volume-storage-coherence-0   Bound    pvc-a3fa6ce3-6588-11e9-bad6-025000000001   2Gi        RWO            hostpath       44m
   persistence-volume-storage-coherence-1   Bound    pvc-d2f732e7-6588-11e9-bad6-025000000001   2Gi        RWO            hostpath       43m
   persistence-volume-storage-coherence-2   Bound    pvc-fc175ea1-6588-11e9-bad6-025000000001   2Gi        RWO            hostpath       41m
   snapshot-volume-storage-coherence-0      Bound    pvc-a3fb2172-6588-11e9-bad6-025000000001   2Gi        RWO            hostpath       44m
   snapshot-volume-storage-coherence-1      Bound    pvc-d2f89ce9-6588-11e9-bad6-025000000001   2Gi        RWO            hostpath       43m
   snapshot-volume-storage-coherence-2      Bound    pvc-fc13ae4b-6588-11e9-bad6-025000000001   2Gi        RWO            hostpath       41m
   ```

6. Reinstall the Coherence cluster with Persistence enabled, using the following command:


   ```bash
   $ helm install \
      --namespace sample-coherence-ns \
      --name storage \
      --set clusterSize=3 \
      --set cluster=storage-tier-cluster \
      --set imagePullSecrets=sample-coherence-secret \
      --set prometheusoperator.enabled=false \
      --set logCaptureEnabled=false \
      --set store.persistence.enabled=true \
      --set store.snapshot.enabled=true \
      coherence/coherence
   ```   

   Wait until all the three pods are running before you continue to the next step.

7. Ensure that the data previously added still exists.

   a. Connect to the Coherence console using the following:

   ```bash
   $ kubectl exec -it --namespace sample-coherence-ns storage-coherence-0 bash /scripts/startCoherence.sh console
   ```   

   b. At the `Map (?):` prompt, type `cache test`.  This creates/uses a cache in the service, `PartitionedCache`.

   c. At the prompt, type `size` and it should show 50000.

  This shows that the previous data entered has automatically been recovered as the PVC was honoured.

   > **Note:** There is currently a bug with default /persistence mount not being created on Docker for Mac environments,
   > and therefore the size may show zero.  

   d. Clear the cache using the `clear` command and confirm that the cache size is zero.

    Recover the `test-snapshot` using:

   ```bash
   snapshot recover test-snapshot
   ```

   The size of the cache should now be 50000.

   To exit the console, type `bye`.

## Uninstall the Charts

Run the following commands to delete the chart installed in this sample.

```bash
$ helm delete storage --purge
```

After deleting the pods, run the following command to delete the PVC:

```bash
$ kubectl get pvc -n sample-coherence-ns | sed 1d | awk '{print $1}' | xargs kubectl delete pvc -n sample-coherence-ns
```
```console
persistentvolumeclaim "persistence-volume-storage-coherence-0" deleted
persistentvolumeclaim "persistence-volume-storage-coherence-1" deleted
persistentvolumeclaim "persistence-volume-storage-coherence-2" deleted
persistentvolumeclaim "snapshot-volume-storage-coherence-0" deleted
persistentvolumeclaim "snapshot-volume-storage-coherence-1" deleted
persistentvolumeclaim "snapshot-volume-storage-coherence-2" deleted

$ kubectl get pvc -n sample-coherence-ns
No resources found.
```

Before starting another sample, ensure that all  pods are removed from the previous sample. To remove `coherence-operator`, use the `helm delete` command.
