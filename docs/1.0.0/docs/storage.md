# Recommended Storage for Coherence

<script>
  window.location.href = "https://docs.coherence.community/coherence-operator/docs/latest/docs/coherence/080_persistence";
</script>

### Active Persistence

For active persistence, local storage is recommended as its IO has a lower latency.

Set up Coherence active persistence with local storage using the following steps:

#### 1. Create a Local Storage Class

If necessary, create a local storage class, `localsc`, using YAML.

```
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"storage.k8s.io/v1beta1","kind":"StorageClass","metadata":{"annotations":{"storageclass.beta.kubernetes.io/is-default-class":"false"},"name":"localsc","namespace":""},"provisioner":"kubernetes.io/no-provisioner"}
    storageclass.beta.kubernetes.io/is-default-class: "false"
  name: localsc
  selfLink: /apis/storage.k8s.io/v1/storageclasses/localsc
provisioner: kubernetes.io/no-provisioner
volumeBindingMode: WaitForFirstConsumer
```

#### 2. Create Local Persistent Volumes

Since local volumes can only be used as a statically created PersistentVolume, persistent volumes need to be created before they are used. For information about persistent volumes, read the [Kubernetes document](https://kubernetes.io/docs/concepts/storage/#local).

A sample YAML for creating PersistentVolume, `mylocalsc-pv0` with label `type: local` is as follows:

```
kind: PersistentVolume
apiVersion: v1
metadata:
  name: mylocalsc-pv0
  labels:
    type: local
spec:
  storageClassName: mylocalsc
  capacity:
    storage: 2Gi
  accessModes:
    - ReadWriteOnce
  hostPath:
    path: "/coh/mydata"
```


The sample YAML file can be used to create PersistentVolume using the following ```kubectl``` command:

```
kubectl create -f mypv0.yaml
```

**Note**: The number of Persistent Volumes created should be the same as the Coherence cluster size.
For instance, in case of a cluster of size 3, you must create Persistent Volumes, `mylocalsc-pv1`, `mylocalsc-pv2`, and `mylocalsc-pv3`.

#### 3. Specify Parameters for Coherence Helm Install

##### a. Set nodeSelector for Coherence Clusters

In order to use the local storage, the Coherence cluster must run a specified set
of nodes. Labels are used to identify the nodes that you want to use.

Note that the set of nodes identified must be the same as the Coherence cluster size.

Suppose the label `name=pool1` is sufficient to identify the nodes that are required,
then the following Helm parameter needs to be set:

```
--set nodeSelector.name=pool1
```

##### b. Enable Active Persistence in Coherence
Set the following Helm parameter to enable Active Persistence in Coherence:
```
--set store.persistence.enabled=true
```

##### c. Specify Persistent Volumes Used for Coherence

Set storage class and set labels for persistent volumes.

```
--set store.persistence.storageClass=mylocalsc
```
```
--set store.persistence.selector.matchLabels.type=local
```

#### 4. Uninstall Coherence with Persistence
##### i. Remove the Installation Using Helm Delete
##### ii. Delete the Persistent Volume Claim (PVC)
***Note***: The PVC is created in the same namespace as your Helm installation. The number of PVC equals the Coherence cluster size.

Retrieve the name of the PVC as follows:
```
kubectl get pvc -n your_namespace
```
Then, delete the PVC one at a time as follows:
```
kubectl delete pvc -n your_namespace your_pvc_name
```

##### iii. Delete the Persistent Volume
```
kubectl delete pv mylocalsc-pv0
```


### Snapshot
By default, Coherence snapshot uses the same location as active persistence to store snapshot data.


If you want to use a different location, then add the following parameter during Helm installation:

```
--set store.snapshot.enabled=true
```

In this case, you should use block storage and configure `store.snapshot.storageClass`
and `store.snapshot.selector.matchLabels` as in step 3c above.

 OCI block storage is used by default in an OCI environment. Therefore, it is not required to set the above two properties.
