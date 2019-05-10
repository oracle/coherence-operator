# Scrape metrics from your own Prometheus instance

You may wish to scrape the metrics from your own Prometheus instance rather than use the 
included prometheus and Grafana. 

This sample shows you how to scrape metrics from your own Prometheus instance.

> Note, use of Prometheus and Grafana is only available when using the
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
   
   Once the install has completed issue the following command to list the pods:
   ```bash
   $ kubectl get pods -n sample-coherence-ns
   NAME                                  READY   STATUS    RESTARTS   AGE
   coherence-operator-665489854f-jr7qj   1/1     Running   0          9s
   ```
   
   There should only be a single `coherence-operator` pod.
   
1. Install the Coherence cluster with `prometheusoperator` enabled

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
   
   Once the install has completed issue the following command to list the pods:

   ```bash
   $ kubectl get pods -n sample-coherence-ns
   NAME                                  READY   STATUS    RESTARTS   AGE
   coherence-operator-665489854f-jr7qj   1/1     Running   0          3m
   storage-coherence-0                   1/1     Running   0          1m
   storage-coherence-1                   1/1     Running   0          1m
   storage-coherence-2                   0/1     Running   0          22s
   ```
 
## Configure your Prometheus to scrape Coherence pods

Normally if `prometheusoperator.enabled=true` for the `coherence-operator`, Prometheus will be installed and setup 
to scrape all the coherence pods. By just setting `prometheusoperator.enabled=true` for coherence install, 
each pod will expose metrics on :9095/metrics. 

You then need to point your prometheus to these targets. 

You can simulate this by using port-forward:

```bash
$ kubectl port-forward storage-coherence-0 -n sample-coherence-ns 9095:9095
```

Use `curl` or a browser to retrieve the metrics from the url `http://127.0.0.1:9095/metrics`.

```bash
$ curl http://127.0.0.1:9095/metrics | tail -2
coherence_jmx_scrape_duration_ms{member="storage-coherence-0", nodeId="1", cluster="storage-tier-cluster", site="coherence.sample-coherence-ns.svc.cluster.local", machine="docker-for-desktop", role="CoherenceServer"} 6
coherence_jmx_scrape_total{member="storage-coherence-0", nodeId="1", cluster="storage-tier-cluster", site="coherence.sample-coherence-ns.svc.cluster.local", machine="docker-for-desktop", role="CoherenceServer"} 1
```

## Uninstalling the Charts

Carry out the following commands to delete the chart installed in this sample.

```bash
$ helm delete storage --purge
```

Before starting another sample, ensure that all the pods are gone from previous sample.

If you wish to remove the `coherence-operator`, then include it in the `helm delete` command above.



