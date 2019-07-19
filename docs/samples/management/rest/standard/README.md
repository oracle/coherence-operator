# Access Management over REST

When the Coherence chart is installed, the Management over REST endpoint is exposed on port 30000 on each of the pods by default.

This sample shows how you can access the Management over REST endpoint using the following URL:

 `http://host:30000/management/coherence/cluster`.

You can view the Swagger document at:  

 `http://host:30000/management/coherence/cluster/metadata-catalog`.

> **Note**: Use of Management over REST is available only when using the operator with Oracle Coherence 12.2.1.4.0.

[Return to Management over REST samples](../)  [Return to Management samples](../../) / [Return to samples](../../../README.md#list-of-samples)

## Prerequisites

Ensure you have already installed the Coherence Operator by using the instructions [here](../../../README.md#install-the-coherence-operator).

## Installation Steps

1. Install the Coherence cluster

   Execute the following command to install the cluster:

   ```bash
   $ helm install \
      --namespace sample-coherence-ns \
      --name storage \
      --set clusterSize=3 \
      --set cluster=coherence-cluster \
      --set imagePullSecrets=sample-coherence-secret \
      --set prometheusoperator.enabled=false \
      --set logCaptureEnabled=false \
      --set coherence.image=your-12.2.1.4.0-Coherence-image \
      coherence/coherence
   ```
   
   > *Note:* If your version of the Coherence Operator does not default to using Coherence 12.2.1.4.0, then you need to replace `your-12.2.1.4.0-Coherence-image` with an appropriate 12.2.1.4.0 image.
   
   After the installation completes, list the pods:
   
   ```bash
   $ kubectl get pods -n sample-coherence-ns
   NAME                   READY   STATUS    RESTARTS   AGE
   storage-coherence-0    1/1     Running   0          4m
   storage-coherence-1    1/1     Running   0          2m   
   storage-coherence-2    1/1     Running   0          2m
   ```
   
1. Port forward the Management over REST port:

   ```bash
   $ kubectl port-forward storage-coherence-0 -n sample-coherence-ns 30000:30000
   ```
   ```console
   Forwarding from [::1]:30000 -> 30000
   Forwarding from 127.0.0.1:30000 -> 30000
   ```   

1. Access Management over REST

   Use Curl to access the endpoint:
   
   ```bash
   $ curl --noproxy http://127.0.0.1:30000/management/coherence/cluster
   ```
   
   This returns the top-level JSON. You can access the Swagger endpoint via `http://127.0.0.1:30000/management/coherence/cluster/metadata-catalog`.
   
   You can specify individual attributes via the following:
   
   ```bash
   $ curl http://127.0.0.1:30000/management/coherence/cluster?fields=clusterName,running,version,clusterSize
   ``` 

   The output, minus the links element, will be similar to:

   ```json
   {
   "links": [ ... ]
   "clusterSize":3,
   "version":"12.2.1.4.0",
   "running":true,
   "clusterName":"coherence-cluster"}
   }
   ```
## Uninstall the Charts

Use the following command to delete the chart installed in this sample:

```bash
$ helm delete storage --purge
```

Before starting another sample, ensure that all the pods are removed from previous samples.

If you want to remove the `coherence-operator`, then include it in the `helm delete` command.
