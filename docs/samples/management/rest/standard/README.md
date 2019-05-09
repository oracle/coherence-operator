# Access management over REST

By default when the Coherence chart is installed the Management over REST endpoint will be exposed
as port 30000 on each of the Pods. 

This sample shows how you can access the Management of REST endpoint via the following URL
`http://host:300000`/management/coherence/cluster`.

> Note, use of Management over REST is only available when using the
> operator with Coherence 12.2.1.4.

[Return to Management over REST samples](../)  [Return to Management samples](../../) / [Return to samples](../../../README.md#list-of-samples)

## Prerequisites

Ensure you have already installed the Coherence Operator by using the instructions [here](../../../README.md#install-the-coherence-operator).

## Installation Steps

1. Install the Coherence cluster

   Issue the following to install the cluster with 1 MBean Server Pod:

   ```bash
   $ helm install \
      --namespace sample-coherence-ns \
      --name storage \
      --set clusterSize=3 \
      --set cluster=coherence-cluster \
      --set imagePullSecrets=sample-coherence-secret \
      --set prometheusoperator.enabled=false \
      --set logCaptureEnabled=false \
      coherence-community/coherence
   ```
   
   Use `kubectl get pods -n sample-coherence-ns` and wait until all 3 pods are running.
   
1. Port-Forward the Management over REST port

   ```bash
   $ kubectl port-forward storage-coherence-0 -n sample-coherence-ns 30000:30000
   Forwarding from [::1]:30000 -> 30000
   Forwarding from 127.0.0.1:30000 -> 30000
   ```   

1. Access Management Over REST

   Using `curl`, issue the following command to access the endpoint
   
   ```bash
   $ curl http://127.0.0.1:30000/management/coherence/cluster
   ```
   
   This will return the top-level JSON.  You can access the Swagger endpoint via `http://127.0.0.1:30000/api-docs`.
   
   You can specify individual attributes via the following:
   
   ```bash
   $ curl http://127.0.0.1:30000/management/coherence/cluster?fields=clusterName,running,version,clusterSize
   ``` 
   
   The output, minus the links element, should be similar to below:
   ```json
   {
   "links": [ ... ]
   "clusterSize":3,
   "version":"12.2.1.4.0",
   "running":true,
   "clusterName":"coherence-cluster"}
   }
   ```
## Uninstalling the Charts

Carry out the following commands to delete the chart installed in this sample.

```bash
$ helm delete storage --purge
```

Before starting another sample, ensure that all the pods are gone from previous samples.

If you wish to remove the `coherence-operator`, then include it in the `helm delete` command above. 
  

   

