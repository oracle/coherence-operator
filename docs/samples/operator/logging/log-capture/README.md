# Enable log capture to view logs in Kiabana   

The Oracle Coherence Operator manages logging data through the EFK
(ElasticSearch, Fluentd and Kibana) stack. Log capture is disabled be default.

This sample shows how to enable log capture and access the Kibana user interface
to view the captured logs.

[Return to Logging samples](../) / [Return to Coherence Operator samples](../../) / [Return to samples](../../../README.md#list-of-samples)

## Installation Steps
                        
When you install the `coherence-operator` and `coherence` charts, you must specify the following
option for `helm` , for both charts, to ensure the EFK stack (Elasitcsearch, Fluentd and Kibana) 
is installed and correctly configured.

```bash
--set logCaptureEnabled=true 
```

1. Install Coherence Operator

   Issue the following command to install `coherence-operator` with log capture enabled:
   
   > Note: If you already have the `coherence-operator` installed without log capture enabled, you
   > must delete it via `helm delete coherence-operator --purge`, before continuing.
   
   ```bash
   $ helm install \
      --namespace sample-coherence-ns \
      --set imagePullSecrets=sample-coherence-secret \
      --name coherence-operator \
      --set logCaptureEnabled=true \
      --set "targetNamespaces={sample-coherence-ns}" \
      coherence/coherence-operator
   ```
   
   Once the install has completed issue the following command to list the pods:
   ```bash
   $ kubectl get pods -n sample-coherence-ns
   NAME                                  READY   STATUS    RESTARTS   AGE
   coherence-operator-7f596c6796-nns9n   2/2     Running   0          41s
   elasticsearch-5b5474865c-86888        1/1     Running   0          41s
   kibana-f6955c4b9-4ndsh                1/1     Running   0          41s
   ```
   
   Along with the `coherence-operator`, you should also see `elasticsearch` and `kibana` pods.
   
1. Install Coherence cluster with logCapture enabled

   ```bash
   $ helm install \
      --namespace sample-coherence-ns \
      --name storage \
      --set clusterSize=3 \
      --set cluster=storage-tier-cluster \
      --set imagePullSecrets=sample-coherence-secret \
      --set logCaptureEnabled=true \
      coherence/coherence
   ```
   
   Once the install has completed issue the following command to list the pods:

   ```bash
   $ kubectl get pods -n sample-coherence-ns
   NAME                                  READY   STATUS    RESTARTS   AGE
   coherence-operator-7f596c6796-nns9n   2/2     Running   0          22m
   elasticsearch-5b5474865c-86888        1/1     Running   0          22m
   kibana-f6955c4b9-4ndsh                1/1     Running   0          22m
   storage-coherence-0                   2/2     Running   0          17m
   storage-coherence-1                   2/2     Running   0          16m
   storage-coherence-2                   2/2     Running   0          16m
   ```
   
   Notice that the `coherence-operator` and all the `coherence` pods have two containers.
   
   If you try to view logs, you must specify the container `coherence` or `fluentd`. 
   
   ```bash
   $ kubectl logs storage-coherence-0 -n sample-coherence-ns
   Error from server (BadRequest): a container name must be specified for pod storage-coherence-0, choose one of:
     [coherence fluentd] or one of the init containers: [coherence-k8s-utils]
   ```
   
   ```bash
   $ kubectl logs storage-coherence-0 -n sample-coherence-ns coherence | tail -5
   2019-04-16 01:45:18.316/92.963 Oracle Coherence GE 12.2.1.4.0 <Info> (thread=Proxy, member=1): Member 3 joined Service Proxy with senior member 1
   2019-04-16 01:45:18.501/93.148 Oracle Coherence GE 12.2.1.4.0 <Info> (thread=Proxy:MetricsHttpProxy, member=1): Member 3 joined Service MetricsHttpProxy with senior member 1
   2019-04-16 01:45:19.281/93.928 Oracle Coherence GE 12.2.1.4.0 <Info> (thread=DistributedCache:PartitionedCache, member=1): Transferring 44B of backup[1] for PartitionSet{172..215} to member 3
   2019-04-16 01:45:19.437/94.084 Oracle Coherence GE 12.2.1.4.0 <Info> (thread=DistributedCache:PartitionedCache, member=1): Transferring primary PartitionSet{128..171} to member 3 requesting 44
   2019-04-16 01:45:19.650/94.297 Oracle Coherence GE 12.2.1.4.0 <Info> (thread=DistributedCache:PartitionedCache, member=1): Partition ownership has stabilized with 3 nodes
   ```
   
## Access Kibana

Use the `port-forward-kibana.sh` script in the
[../../common](../../../common) directory to view log messages.

1. Start the port-forward

   ```bash
   $ ./port-forward-kibana.sh sample-coherence-ns

   Forwarding from 127.0.0.1:5601 -> 5601
   Forwarding from [::1]:5601 -> 5601
   ```
1. Access Kibana using the following URL:

   [http://127.0.0.1:5601/](http://127.0.0.1:5601/)
   
   >Note: It may take up to 5 minutes for the data to reach the elasticsearch instance.   
   
## Default Dashboards

There are a number of dashboard created via the import process.

**Coherence Operator**

* *Coherence Operator - All Messages* - Shows all Coherence Operator messages

**Coherence Cluster**

* *Coherence Cluster - All Messages* - Shows all messages

* *Coherence Cluster - Errors and Warnings* - Shows only errors and warnings

* *Coherence Cluster - Persistence* - Shows persistence related messages 

* *Coherence Cluster - Partitions* - Shows Ppartition related messages 

* *Coherence Cluster - Message Sources* - Allows visualization of messages via the message source (Thread)

* *Coherence Cluster - Configuration Messages* - Shows configuration related messages

* *Coherence Cluster - Network* - Shows network related messages such as communication delays and TCP ring disconnects 

## Default Queries

There are many queries related to common Coherence messages, warnings and errors that are 
loaded and can be accessed via the `Discover` side-bar.

## Uninstalling the Charts

Carry out the following commands to delete the chart installed in this sample.

```bash
$ helm delete storage --purge
```

Before starting another sample, ensure that all the pods are gone from previous sample.

If you wish to remove the `coherence-operator`, then include it in the `helm delete` command above.

