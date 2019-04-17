> Note, use of Prometheus and Grafana is only available when using the
> operator with Coherence 12.2.1.4.


# Monitoring Coherence services via Grafana dashboards

The Oracle Coherence Operator includes the Prometheus Operator as an optional subchart named `prometheusoperator`.
To configure the Prometheus Operator and monitor Coherence services via grafana dashboards, 
please follow the instructions below.

## 1. Installing the Charts

When you install the `coherence-operator` chart, you must specify the following
additional set value for `helm` to install subchart `prometheusoperator`.

```bash
--set prometheusoperator.enabled=true
```

All `coherence` charts installed in `coherence-operator` `targetNamespaces` are monitored by 
Prometheus. The servicemonitor `<releasename>-coherence-service-monitor` 
configures Prometheus to scrape all components of `coherence-service`.


## 2. Port Forward Grafana

Once you have installed the charts, use the following script to port forward the Grafana pod.

```bash
#!/bin/bash
  
while :
do
  kubectl port-forward $(kubectl get pods --selector=app=grafana -n namespace --output=jsonpath="{.items..metadata.name}") -n namespace 9200:3000
done

```

## 3. Login to Grafana

In browser, go to url `http://localhost:9200`.

At the Grafana login screen, the login is `admin` and the password is `prom-operator`.

Click `Home` in the upper left corner of screen to get a list of preconfigured dashboards.
Click ` Coherence Dashboard Main`.


## 4. Default Dashboards

There are a number of dashboard created via the import process.

* Coherence Dashboard main for inspecting coherence cluster(s)

* Coherence Cluster Members Summary and Details

* Coherence Cluster Members Machines Summary

* Coherence Cache Summary and Details

* Coherence Services Summary and Details

* Coherence Proxy Servers Summary and Details

* Coherence Elastic Data Summary

* Coherence Cache Persistence Summary

* Coherence Http Servers Summary

## 5. Troubleshooting

## Helm install of `coherence-operator` fails creating a custom resource definition (CRD).

Follow recommendation from [Prometheus Operator: helm fails to create CRDs](https://github.com/helm/charts/tree/master/stable/prometheus-operator#user-content-helm-fails-to-create-crds)
to manually install the Prometheus Operator CRDs, then install the `coherence-operator` chart with these additional set values. 

```bash
--set prometheusoperator.enabled=true --set prometheusoperator.prometheusOperator.createCustomResource=false
```

### No datasource found

Manually create a datasource by clicking on Grafana Home `Create your first data source` button 
and fill in these fields.
  
```bash
   Name:      Prometheus 
   HTTP URL:  http://release-name-prometheus:9090/
```

CLick `Save & Test` button on bottom of page.
