# Installing Multiple Coherence Clusters with One Operator

This sample shows how Coherence Operator can manage two or more Coherence clusters, and
how you can see the logs from both clusters using Kibana.

[Return to Coherence Deployments samples](../) / [Return to samples](../../README.md#list-of-samples)

## Installation Steps

1. Install Coherence Operator.

   Run the following command to install `coherence-operator` with log capture enabled:

   > **Note:** If you have already installede `coherence-operator` without log capture enabled, you
   > must delete it using `helm delete coherence-operator --purge` before continuing.

   ```bash
   $ helm install \
      --namespace sample-coherence-ns \
      --set imagePullSecrets=sample-coherence-secret \
      --name coherence-operator \
      --set logCaptureEnabled=true \
      --set "targetNamespaces={sample-coherence-ns}" \
      coherence/coherence-operator  
   ```

   After installation, run the following command to list the pods:

   ```bash
   $ kubectl get pods -n sample-coherence-ns
   ```
   Along with `coherence-operator`, the command also returns the `elasticsearch` and `kibana` pods.
   ```console
   NAME                                  READY   STATUS    RESTARTS   AGE
   coherence-operator-7f596c6796-nns9n   2/2     Running   0          41s
   elasticsearch-5b5474865c-86888        1/1     Running   0          41s
   kibana-f6955c4b9-4ndsh                1/1     Running   0          41s
   ```
2. Install the first Coherence cluster, `cluster-a`.

   For each cluster, create only two pods to save  resources.

   ```bash
   $ helm install \
      --namespace sample-coherence-ns \
      --name cluster-a \
      --set clusterSize=2 \
      --set cluster=cluster-a \
      --set imagePullSecrets=sample-coherence-secret \
      --set prometheusoperator.enabled=false \
      --set logCaptureEnabled=true \
      coherence/coherence
   ```

3. Install the second Coherence cluster, `cluster-b`

   For each cluster, create only two pods to save resources.

   ```bash
   $ helm install \
      --namespace sample-coherence-ns \
      --name cluster-b \
      --set clusterSize=2 \
      --set cluster=cluster-b \
      --set imagePullSecrets=sample-coherence-secret \
      --set prometheusoperator.enabled=false \
      --set logCaptureEnabled=true \
      coherence/coherence
   ```

## Access Kibana

Use the `port-forward-kibana.sh` script in the
[../../common](../../../common) directory to view log messages.

1. Start port forwarding.

   ```bash
   $ ./port-forward-kibana.sh sample-coherence-ns
   ```
   ```console
   Forwarding from 127.0.0.1:5601 -> 5601
   Forwarding from [::1]:5601 -> 5601
   ```
2. Access Kibana using the following URL:

   [http://127.0.0.1:5601/](http://127.0.0.1:5601/)

   > **Note:** It may take up to five minutes for the data to reach the elasticsearch instance.   

   Once logged in, click `Dashboard` on the left and choose `Coherence Cluster - All Messages`.
   You should see messages from both clusters, `cluster-a` and `cluster-b`.

   ![Coherence Cluster - All Messages](img/kibana-dashboard.png)


## Uninstall the Charts

Run the following command to delete the charts installed in this sample.

```bash
$ helm delete cluster-a cluster-b --purge
```

Before starting another sample, ensure that all  pods are removed from the previous sample. To remove `coherence-operator`, use the `helm delete` command.
