# Use default persistent volume claim

By default when you enable Coherence Persistence the required infrastructure in terms
of Persistent Volumes (PV) and claims (PVC) is setup automatically as well as the persistence-mode
being set to `active`. This allows restarting of a Coherence cluster and retaining
the data.

This sample shows how to enabled Persistence with the all the defaults 
under `store.persistence` in `coherence` chart `values.yaml`. Please see (this sample)[../pvc/README.md]
for more details on setting other values such as `storageClasses`.

[Return to Persistence samples](../) / [Return to Coherence Deployments samples](../../) / [Return to samples](../../../README.md#list-of-samples)
 
## Prerequisites

Ensure you have already installed the Coherence Operator by using the instructions [here](../../../README.md#install-the-coherence-operator).

## Installation Steps

1. Install the Coherence cluster

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
      --version 1.0.0-SNAPSHOT coherence-community/coherence
   ```
   
   You may also change the size of the default directories 
   
   * `/persistence` default 2Gi - `--set store.persistence.size=10Gi`
   
   * `/snapshot` default 2Gi - `--set store.snapshot.size-10Gi`                  

1. Ensure the pods are running:

   ```bash
   $ kubectl get pods  -n sample-coherence-ns
   NAME                                  READY   STATUS    RESTARTS   AGE
   coherence-operator-66f9bb7b75-hqk4l   1/1     Running   0          13m
   storage-coherence-0                   1/1     Running   0          3m
   storage-coherence-1                   1/1     Running   0          2m
   storage-coherence-2                   0/1     Running   0          44s
   ```
   
   Once all three pods are `Running`, the proceed to check the PVC, below.

1. Check the Persistence Volumes created

   Issue the following command to check the PVC created.
   
   ```bash
   $ kubectl get pvc -n sample-coherence-ns
   NAME                                     STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS   AGE
   persistence-volume-storage-coherence-0   Bound    pvc-a3fa6ce3-6588-11e9-bad6-025000000001   2Gi        RWO            hostpath       44m
   persistence-volume-storage-coherence-1   Bound    pvc-d2f732e7-6588-11e9-bad6-025000000001   2Gi        RWO            hostpath       43m
   persistence-volume-storage-coherence-2   Bound    pvc-fc175ea1-6588-11e9-bad6-025000000001   2Gi        RWO            hostpath       41m
   snapshot-volume-storage-coherence-0      Bound    pvc-a3fb2172-6588-11e9-bad6-025000000001   2Gi        RWO            hostpath       44m
   snapshot-volume-storage-coherence-1      Bound    pvc-d2f89ce9-6588-11e9-bad6-025000000001   2Gi        RWO            hostpath       43m
   snapshot-volume-storage-coherence-2      Bound    pvc-fc13ae4b-6588-11e9-bad6-025000000001   2Gi        RWO            hostpath       41m
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
   
   **TODO** - Add in create snapshot command
   
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
   $ kubectl get pvc -n sample-coherence-ns
   NAME                                     STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS   AGE
   persistence-volume-storage-coherence-0   Bound    pvc-a3fa6ce3-6588-11e9-bad6-025000000001   2Gi        RWO            hostpath       44m
   persistence-volume-storage-coherence-1   Bound    pvc-d2f732e7-6588-11e9-bad6-025000000001   2Gi        RWO            hostpath       43m
   persistence-volume-storage-coherence-2   Bound    pvc-fc175ea1-6588-11e9-bad6-025000000001   2Gi        RWO            hostpath       41m
   snapshot-volume-storage-coherence-0      Bound    pvc-a3fb2172-6588-11e9-bad6-025000000001   2Gi        RWO            hostpath       44m
   snapshot-volume-storage-coherence-1      Bound    pvc-d2f89ce9-6588-11e9-bad6-025000000001   2Gi        RWO            hostpath       43m
   snapshot-volume-storage-coherence-2      Bound    pvc-fc13ae4b-6588-11e9-bad6-025000000001   2Gi        RWO            hostpath       41m
   ``` 

1. Re-Install the Coherence cluster

   Issue the following to re-install the cluster with Persistence enabled:

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
   
   > This shows that the previous data entered has automatically been recovered due the PVC being honoured.
    
   Then type `bye` to exit the `console`.
      
## Verifying Grafana Data (If you enabled Prometheus)

Access Grafana using the instructions [here](../../../README.md#access-grafana).

## Verifying Kibana Logs (if you enabled log capture)

Access Kibana using the instructions [here](../../../README.md#access-kibana).

## Uninstalling the Charts

Carry out the following commands to delete the chart installed in this sample.

```bash
$ helm delete storage --purge
```

Once the pods are deleted, issue the following to delete the PVC.:

```bash
$ kubectl get pvc -n sample-coherence-ns | sed 1d | awk '{print $1}' | xargs kubectl delete pvc -n sample-coherence-ns
persistentvolumeclaim "persistence-volume-storage-coherence-0" deleted
persistentvolumeclaim "persistence-volume-storage-coherence-1" deleted
persistentvolumeclaim "persistence-volume-storage-coherence-2" deleted
persistentvolumeclaim "snapshot-volume-storage-coherence-0" deleted
persistentvolumeclaim "snapshot-volume-storage-coherence-1" deleted
persistentvolumeclaim "snapshot-volume-storage-coherence-2" deleted

$ kubectl get pvc -n sample-coherence-ns
No resources found.
```

Before starting another sample, ensure that all the pods are gone from previous sample.

If you wish to remove the `coherence-operator`, then include it in the `helm delete` command above.

