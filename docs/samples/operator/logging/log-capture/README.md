# Enable log capture to view logs in Kiabana   

The Coherence Operator manages data logging through the Elasticsearch, Fluentd and Kibana (EFK) stack. The log capture feature is disabled by default.

This sample shows how to enable log capture and access the Kibana user interface (UI) to view the captured logs.

[Return to Logging samples](../) / [Return to Coherence Operator samples](../../) / [Return to samples](../../../README.md#list-of-samples)

## Installation Steps

When you install `coherence-operator` and `coherence` charts, you must specify the following
option for `helm` , for both charts, to ensure that the EFK stack (Elasitcsearch, Fluentd and Kibana) 
is installed and correctly configured.

```bash
--set logCaptureEnabled=true 
```

1. Install Coherence Operator

   Use the following command to install `coherence-operator` with log capture enabled:
   
   > **Note:** If you have already installed the `coherence-operator` without log capture enabled, you
   > must first delete it using `helm delete coherence-operator --purge` command and then continue.
   
   ```bash
   $ helm install \
      --namespace sample-coherence-ns \
      --set imagePullSecrets=sample-coherence-secret \
      --name coherence-operator \
      --set logCaptureEnabled=true \
      --set "targetNamespaces={sample-coherence-ns}" \
      coherence/coherence-operator
   ```

   After the installation completes, list the pods:

   ```bash
   $ kubectl get pods -n sample-coherence-ns
   ```
 
   ```console
   NAME                                  READY   STATUS    RESTARTS   AGE
   coherence-operator-7f596c6796-nns9n   2/2     Running   0          41s
   elasticsearch-5b5474865c-86888        1/1     Running   0          41s
   kibana-f6955c4b9-4ndsh                1/1     Running   0          41s
   ```
   
   In the output, you can see the pods for Elasticsearch and Kibana along with the operator.
   
1. Install Coherence cluster with log capture enabled:

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
   
   After the installation completes, list the pods:

   ```bash
   $ kubectl get pods -n sample-coherence-ns
   ```
   ```console
   NAME                                  READY   STATUS    RESTARTS   AGE
   coherence-operator-7f596c6796-nns9n   2/2     Running   0          22m
   elasticsearch-5b5474865c-86888        1/1     Running   0          22m
   kibana-f6955c4b9-4ndsh                1/1     Running   0          22m
   storage-coherence-0                   2/2     Running   0          17m
   storage-coherence-1                   2/2     Running   0          16m
   storage-coherence-2                   2/2     Running   0          16m
   ```
  
   The `coherence-operator` and all the `coherence` pods have two containers.
   
   To view the logs, you must specify the container `coherence` or `fluentd`.
   
   ```bash
   $ kubectl logs storage-coherence-0 -n sample-coherence-ns
   ```
   ```console
   Error from server (BadRequest): a container name must be specified for pod storage-coherence-0, choose one of:
     [coherence fluentd] or one of the init containers: [coherence-k8s-utils]
   ```
   
   ```bash
   $ kubectl logs storage-coherence-0 -n sample-coherence-ns coherence | tail -5
   ```
   ```console
   2019-04-16 01:45:18.316/92.963 Oracle Coherence GE 12.2.1.4.0 <Info> (thread=Proxy, member=1): Member 3 joined Service Proxy with senior member 1
   2019-04-16 01:45:18.501/93.148 Oracle Coherence GE 12.2.1.4.0 <Info> (thread=Proxy:MetricsHttpProxy, member=1): Member 3 joined Service MetricsHttpProxy with senior member 1
   2019-04-16 01:45:19.281/93.928 Oracle Coherence GE 12.2.1.4.0 <Info> (thread=DistributedCache:PartitionedCache, member=1): Transferring 44B of backup[1] for PartitionSet{172..215} to member 3
   2019-04-16 01:45:19.437/94.084 Oracle Coherence GE 12.2.1.4.0 <Info> (thread=DistributedCache:PartitionedCache, member=1): Transferring primary PartitionSet{128..171} to member 3 requesting 44
   2019-04-16 01:45:19.650/94.297 Oracle Coherence GE 12.2.1.4.0 <Info> (thread=DistributedCache:PartitionedCache, member=1): Partition ownership has stabilized with 3 nodes
   ```

## Access Kibana

Run the `port-forward-kibana.sh` script in the
[common](../../../common) directory to view the log messages.

1. Start the port-forward:

   ```bash
   $ ./port-forward-kibana.sh sample-coherence-ns
   ```
   ```console
   Forwarding from 127.0.0.1:5601 -> 5601
   Forwarding from [::1]:5601 -> 5601
   ```
1. Access Kibana using the following URL:

   [http://127.0.0.1:5601/](http://127.0.0.1:5601/)
   
   > **Note:** It takes approximately 5 minutes for the data to reach the Elasticsearch instance.
   
## Default Kibana Dashboards

There are a number of Kibana dashboards created via the import process.

| Dashboard | Options| Description|
|-----------|------------|-----------------|
| Coherence Operator | All Messages| Shows all Coherence Operator messages                                                |
| Coherence Cluster  | All Messages           | Shows all messages                                                                   |
| Coherence Cluster  | Errors and Warnings    | Shows only errors and warnings                                                       |
| Coherence Cluster  | Persistence            | Shows partition related messages                                                    |
| Coherence Cluster  | Message Sources        | Allows visualization of messages via the message source (Thread)                     |
| Coherence Cluster  | Configuration Messages | Shows configuration related messages                                                 |
| Coherence Cluster  | Network                | Shows network related messages, such as communication delays and TCP ring disconnects |

## Default Queries

There are many queries related to common Coherence messages, warnings, and errors that are loaded and can be accessed via the `Discover` side-bar.

## Uninstalling the Charts

Use the following commands to delete the chart installed in this sample:

```bash
$ helm delete storage --purge
```

> **Note**: If you are using Kubernetes 1.13.0 or older version, you cannot delete the pods. This is a known issue and you need to add the options `--force --grace-period=0` to force delete the pods.
>
> Refer to [https://github.com/kubernetes/kubernetes/issues/45688](https://github.com/kubernetes/kubernetes/issues/45688).

Before starting another sample, ensure that all the pods are deleted from the previous sample.

If you want to remove the `coherence-operator`, then use the `helm delete` command.
