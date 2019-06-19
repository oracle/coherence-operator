# Scrape metrics from your own Prometheus instance

You may wish to scrape the metrics from your own Prometheus Operator instance rather than use the
prometheusopeartor subchart included with coherence-operator. 

This sample shows you how to scrape metrics from your own Prometheus instance.

> **Note:** Use of Prometheus and Grafana is only available when using the
> operator with Coherence 12.2.1.4.

[Return to Metrics samples](../) / [Return to Coherence Operator samples](../../) / [Return to samples](../../../README.md#list-of-samples)

## Installation Steps

1. Install Coherence Operator

   When you install the `coherence-operator` chart, you must ensure to specify `--set prometheusoperator.enabled=false`
   or leave out the option completely, which will default to false. 
  
   ```bash
   --set prometheusoperator.enabled=false
   ```

   Issue the following command to install `coherence-operator` with `prometheusoperator` enabled:
   
   ```bash
   $ helm install \
      --namespace sample-coherence-ns \
      --set imagePullSecrets=sample-coherence-secret \
      --name coherence-operator \
      --set prometheusoperator.enabled=false \
      coherence/coherence-operator
   ```
   
   Once the install has completed, issue the following command to list the pods:
   
   ```bash
   $ kubectl get pods -n sample-coherence-ns
   NAME                                  READY   STATUS    RESTARTS   AGE
   coherence-operator-665489854f-jr7qj   1/1     Running   0          9s
   ```
   
   There should only be a single `coherence-operator` pod.
   
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
   
   Once the install has completed, issue the following command to list the pods:

   ```bash
   $ kubectl get pods -n sample-coherence-ns
   NAME                                  READY   STATUS    RESTARTS   AGE
   coherence-operator-665489854f-jr7qj   1/1     Running   0          3m
   storage-coherence-0                   1/1     Running   0          1m
   storage-coherence-1                   1/1     Running   0          1m
   storage-coherence-2                   0/1     Running   0          22s
   ```
 
## Configure your Prometheus Operator to scrape Coherence pods

Please consult the Prometheus Operator documentation on how to configure and deploy a service monitor for 
your own Prometheus Operator installation.

This section only describes service monitor configuration as it relates to the Coherence helm chart.

coherence-service-monitor.yaml fragment:
```
...
spec:
  selector:
    matchLabels:
      component: "coherence-service"
...      
endpoints:
  - port: 9095
```

If the Coherence helm chart parameter `service.metricsHttpPort` is set when installing the Coherence helm chart,
replace `9095` above with the new value.
  
If the Coherence helm chart parameter `store.metrics.ssl.enabled` is `true`, additionally add `endpoints.scheme` value of `https`
to `coherence-service-monitor.yaml` fragment.

Note that there are a number of Coherence Grafana dashboards bundled in the coherence-operator helm chart under dashboards folder.
While Grafana will have to be configured to the location of your prometheus datasource, one can still take advantage
of these Coherence dashboards by extracting them from the coherence-operator helm chart.

## Uninstalling the Charts

Carry out the following commands to delete the chart installed in this sample.

```bash
$ helm delete storage --purge
```

Before starting another sample, ensure that all the pods are gone from previous sample.

If you wish to remove the `coherence-operator`, then include it in the `helm delete` command above.



