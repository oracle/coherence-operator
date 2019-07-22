# Scrape Metrics from Your Prometheus Instance

You can scrape the metrics from your own Prometheus Operator instance rather than using the `prometheusopeartor` subchart included with `coherence-operator`.

This sample shows you how to scrape metrics from your own Prometheus instance.

> **Note:** Use of Prometheus and Grafana is available only when using the operator with Oracle Coherence 12.2.1.4.0 version.

[Return to Metrics samples](../) / [Return to Coherence Operator samples](../../) / [Return to samples](../../../README.md#list-of-samples)

## Installation Steps

1. Install Coherence Operator

   When you install the `coherence-operator`, you must ensure to specify `--set prometheusoperator.enabled=false`
   or leave out the option completely, which also defaults to false.
  
   ```bash
   --set prometheusoperator.enabled=false
   ```

   Use the following command to install `coherence-operator` with `prometheusoperator` enabled:
   
   ```bash
   $ helm install \
      --namespace sample-coherence-ns \
      --set imagePullSecrets=sample-coherence-secret \
      --name coherence-operator \
      --set prometheusoperator.enabled=false \
      coherence/coherence-operator
   ```
   
   After the installation completes, list the pods:
   
   ```bash
   $ kubectl get pods -n sample-coherence-ns
   ```
   ```console
   NAME                                  READY   STATUS    RESTARTS   AGE
   coherence-operator-665489854f-jr7qj   1/1     Running   0          9s
   ```
   
   There is only a single `coherence-operator` pod.
   
2. Install the Coherence cluster

   ```bash
   $ helm install \
      --namespace sample-coherence-ns \
      --name storage \
      --set clusterSize=3 \
      --set cluster=storage-tier-cluster \
      --set imagePullSecrets=sample-coherence-secret \
      --set "targetNamespaces={sample-coherence-ns}" \
      coherence/coherence
   ```
   
   After the installation completes, list the pods:

   ```bash
   $ kubectl get pods -n sample-coherence-ns
   ```
   ```console
   NAME                                  READY   STATUS    RESTARTS   AGE
   coherence-operator-665489854f-jr7qj   1/1     Running   0          3m
   storage-coherence-0                   1/1     Running   0          1m
   storage-coherence-1                   1/1     Running   0          1m
   storage-coherence-2                   0/1     Running   0          22s
   ```

## Configure Your Prometheus Operator to Scrape Coherence Pods

Refer [Prometheus Operator](https://github.com/coreos/prometheus-operator/blob/master/Documentation/user-guides/getting-started.md) documentation for information about how to configure and deploy a service monitor for your Prometheus Operator installation.

This section describes only the service monitor configuration as it relates to the Coherence Helm chart.

`coherence-service-monitor.yaml` fragment:
```
...
spec:
  selector:
    matchLabels:
      component: "coherence-service"
...
endpoints:
  - port: 9612
```

If the parameter `service.metricsHttpPort` is set when installing the Coherence Helm chart, replace `port: 9612` with the new value.
  
If the Coherence Helm chart parameter `store.metrics.ssl.enabled` is `true`, add `endpoints.scheme` value of `https` to `coherence-service-monitor.yaml` fragment.

There are a number of Coherence Grafana dashboards bundled in the Coherence Operator Helm chart under dashboards folder.
While Grafana have to be configured to the location of your Prometheus datasource, you can take advantage of these Coherence dashboards by extracting them from the Coherence Operator Helm chart.

## Uninstall the Charts

Use the following commands to delete the chart installed in this sample:

```bash
$ helm delete storage --purge
```

Before starting another sample, ensure that all the pods are removed from previous sample.

If you want to remove the `coherence-operator`, use the `helm delete` command.
