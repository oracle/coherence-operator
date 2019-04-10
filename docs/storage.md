# Recommended Storage for Coherence

### Active Persistence

For active persistence, local storage is recommended as its IO has a lower latency.

The following are steps for setting up Coherence active persistence with local storage.
#### 1. Create a local storage class.

Use the following yaml to create a local storage class, `localsc`, if it is necessary.

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

#### 2. Create local Persistent Volumes.

Since local volumes can [only by used as a statically created](https://kubernetes.io/docs/concepts/storage/#local)
PersistentVolume, persistent volumes need to be created before they are used.

A sample yaml for creating Persistence Volume `mylocalsc-pv0` with label `type: local` is as follows. 

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

with the following command
```
kubectl create -f mypv0.yaml
```

Note that the number of Persistent Volumes created need to be the same as Coherence cluster size.
For instance, in case of a cluster of size 3, one also need to create `mylocalsc-pv1` and `mylocalsc-pv2`.

#### 3. Specify parameters for Coherence Helm install.

##### a. Set nodeSelector for Coherence clusters.
         
In order to use the local storage, the Coherence cluster needs to run a specified set
of nodes. Labels will be used to identify the nodes that we want used.

Note that the set of nodes identified needs to be the same as the Coherence cluster size.

Suppose the label `name=pool1` is sufficient to identify the nodes that is required,
then the following Helm parameter needs to be set:

```
--set nodeSelector.name=pool1
```

##### b. Enable active persistence in Coherence.

```
--set store.persistence.enabled=true 
```

##### c. Specify Persistent volumes used for Coherence.

###### i. Set storage class.

```
--set store.persistence.storageClass=mylocalsc
```

###### ii. Set labels for persistent volumes.
```
--set store.persistence.selector.matchLabels.type=local 
```

#### 4. Uninstall the Coherence with persistence.
##### i. Helm delete the installation.
##### ii. Delete the pvc created one by one.
Note that the pvc is created in the same namespace as your Helm installation.
And the number of pvc equals the Coherence cluster size.
The name of pvc can be found as follows:
```
kubectl get pvc -n your_namespace
```
Then the pvc can be deleted one by one as follows:
```
kubectl delete pvc -n your_namespace your_pvc_name
```

##### iii. Delete the pv created one by one.
```
kubectl delete pv mylocalsc-pv0
```
etc.


### Snapshot
By default, Coherence snapshot use the same location for persistence snapshot data as in active persistence.

If it is desired to use a different location, then add the following parameter in Helm installation:

```
--set store.snapshot.enabled=true
```

In this case, one should use use block storage and configure `store.snapshot.storageClass` 
and `store.snapshot.selector.matchLabels` as in 3.c. above if it is necessary.

In OCI environement, by default, OCI block storage is used. There is no need to set the above two properties.
