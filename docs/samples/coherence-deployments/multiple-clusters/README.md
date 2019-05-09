# Installing Multiple Coherence clusters with one Operator

This sample shows how the Coherence Operator can manage two or more Coherence clusters and
how you can see the logs from both clusters using Kibana.

[Return to Coherence Deployments samples](../) / [Return to samples](../../README.md#list-of-samples)

## Installation Steps

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
      coherence-community/coherence-operator  
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

1. Install first Coherence cluster `cluster-a`

   for each of the clusters, we only create 2 pods to save on resources.

   ```bash
   $ helm install \
      --namespace sample-coherence-ns \
      --name cluster-a \
      --set clusterSize=2 \
      --set cluster=cluster-a \
      --set imagePullSecrets=sample-coherence-secret \
      --set prometheusoperator.enabled=false \
      --set logCaptureEnabled=true \
      coherence-community/coherence
   ```
   
1. Install second Coherence cluster `cluster-b`

   for each of the clusters, we only create 2 pods to save on resources.

   ```bash
   $ helm install \
      --namespace sample-coherence-ns \
      --name cluster-b \
      --set clusterSize=2 \
      --set cluster=cluster-b \
      --set imagePullSecrets=sample-coherence-secret \
      --set prometheusoperator.enabled=false \
      --set logCaptureEnabled=true \
      coherence-community/coherence
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
   
   Once logged in, click on `Dashboard` on the left and choose `Coherence Cluster - All Messages`.
   You should see messages from both `cluster-a` and `cluster-b`.
   
   ![Coherence Cluster - All Messages](img/kibana-dashboard.png)
   
   
## Uninstalling the Charts

Carry out the following commands to delete the charts installed in this sample.

```bash
$ helm delete cluster-a cluster-b --purge
```

Before starting another sample, ensure that all the pods are gone from previous sample.

If you wish to remove the `coherence-operator`, then include it in the `helm delete` command above.
