# Deploy Coherence Operator with Prometheus Enabled and View Metrics in Grafana

The Coherence Operator includes the Prometheus Operator as an optional subchart named `prometheusoperator`.

This sample shows you how to configure the Prometheus Operator and monitor Coherence services through Grafana.

> **Note**: Use of Prometheus and Grafana is available only when using the operator with Oracle Coherence 12.2.1.4.0 version.

[Return to Metrics samples](../) / [Return to Coherence Operator samples](../../) / [Return to samples](../../../README.md#list-of-samples)

## Installation Steps

1. Install Coherence Operator

   When you install the `coherence-operator` chart, you must specify the following additional set value for `helm` to install subchart `prometheusoperator`:
  
   ```bash
   --set prometheusoperator.enabled=true
   ```
  
   All `coherence` charts installed in `coherence-operator` `targetNamespaces` are monitored by Prometheus. The ServiceMonitor `<releasename>-coherence-service-monitor` 
   configures Prometheus to scrape all components of `coherence-service`.

   Use the following command to install `coherence-operator` with `prometheusoperator` enabled:
   
   ```bash
   $ helm install \
      --namespace sample-coherence-ns \
      --set imagePullSecrets=sample-coherence-secret \
      --name coherence-operator \
      --set prometheusoperator.enabled=true \
      --set prometheusoperator.prometheusOperator.createCustomResource=false \
      --set "targetNamespaces={sample-coherence-ns}" \
      coherence/coherence-operator
   ```
   
   After the installation completes, list the pods:
   
   ```bash
   $ kubectl get pods -n sample-coherence-ns
   ```
   ```console
   NAME                                                     READY   STATUS    RESTARTS   AGE
   coherence-operator-66f9bb7b75-q2w8t                      1/1     Running   0          34s
   coherence-operator-grafana-769bb4d5cb-xwm9w              3/3     Running   0          35s
   coherence-operator-kube-state-metrics-5d5f6855bd-hh7cv   1/1     Running   0          35s
   coherence-operator-prometh-operator-58bd58ddfd-rldqk     1/1     Running   0          34s
   coherence-operator-prometheus-node-exporter-n9ls7        1/1     Running   0          35s
   prometheus-coherence-operator-prometh-prometheus-0       3/3     Running   1          21s
   ```

   You can see `grafana` and other `promethues` related pods along with the `coherence-operator` pod.

1. Install the Coherence cluster with `prometheusoperator` enabled:

   ```bash
   $ helm install \
      --namespace sample-coherence-ns \
      --name storage \
      --set clusterSize=3 \
      --set cluster=storage-tier-cluster \
      --set imagePullSecrets=sample-coherence-secret \
      --set prometheusoperator.enabled=true \
      --set "targetNamespaces={sample-coherence-ns}" \
      coherence/coherence
   ```

   > **Note:** If the Coherence Operator chart version does not have the default Coherence image as 12.2.1.4, then you need to set this using `--set coherence.image=your-12.2.1.4-image`.

   After the installation completes, list the pods:

   ```bash
   $ kubectl get pods -n sample-coherence-ns
   ```
   ```console
   NAME                                                     READY   STATUS    RESTARTS   AGE
   coherence-operator-66f9bb7b75-q2w8t                      1/1     Running   0          9m
   coherence-operator-grafana-769bb4d5cb-xwm9w              3/3     Running   0          9m
   coherence-operator-kube-state-metrics-5d5f6855bd-hh7cv   1/1     Running   0          9m
   coherence-operator-prometh-operator-58bd58ddfd-rldqk     1/1     Running   0          9m
   coherence-operator-prometheus-node-exporter-n9ls7        1/1     Running   0          9m
   prometheus-coherence-operator-prometh-prometheus-0       3/3     Running   1          9m
   storage-coherence-0                                      1/1     Running   0          3m
   storage-coherence-1                                      1/1     Running   0          2m
   storage-coherence-2                                      1/1     Running   0          1m
   ```
 
## Access Grafana

Run the `port-forward-grafana.sh` script in the [common](../../../common) directory to view metrics.

1. Start the port-forward:

   ```bash
   $ ./port-forward-grafana.sh sample-coherence-ns

   Forwarding from 127.0.0.1:3000 -> 3000
   Forwarding from [::1]:3000 -> 3000
   ```
   
1. Access Grafana using the following URL:

   [http://127.0.0.1:3000/d/coh-main/coherence-dashboard-main](http://127.0.0.1:3000/d/coh-main/coherence-dashboard-main)

   * Username: admin  

   * Password: prom-operator

## Default Dashboards

There are a number of dashboard created via the import process:

* Coherence Dashboard main for inspecting coherence clusters

* Coherence Cluster Members Summary and Details

* Coherence Cluster Members Machines Summary

* Coherence Cache Summary and Details

* Coherence Services Summary and Details

* Coherence Proxy Servers Summary and Details

* Coherence Elastic Data Summary

* Coherence Cache Persistence Summary

* Coherence Http Servers Summary

## Uninstall the Charts

Use the following command to delete the charts installed in this sample:

```bash
$ helm delete storage --purge
```

Before starting another sample, ensure that all the pods are removed from previous sample.

If you want to remove the `coherence-operator`, then use the `helm delete` command.

## Troubleshooting

### Helm install of `coherence-operator` fails creating a custom resource definition (CRD)

See [Prometheus Operator: helm fails to create CRDs](https://github.com/helm/charts/tree/master/stable/prometheus-operator#user-content-helm-fails-to-create-crds)
to manually install the Prometheus Operator CRDs and then install the `coherence-operator` chart with these additional set values.

```bash
--set prometheusoperator.enabled=true --set prometheusoperator.prometheusOperator.createCustomResource=false
```

### No datasource found in Grafana

On the Grafana home page, click **Create your first data source** to create a datasource manually and fill in these fields:
  
```bash
   Name:      Prometheus 
   HTTP URL:  http://{release-name}-prometheus:9090/
```

CLick **Save & Test**.
