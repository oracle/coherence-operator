> **Note:** Use of Prometheus and Grafana is only available when using the
> operator with Coherence 12.2.1.4.


# Monitoring Coherence services via Grafana dashboards

The Oracle Coherence Operator includes the Prometheus Operator as an optional subchart named `prometheusoperator`.
To configure the Prometheus Operator and monitor Coherence services via grafana dashboards, 
please follow the instructions below.

This use-case is covered [in the samples](docs/samples/operator/metrics/enable-metrics/).

## Installing the charts

When you install the `coherence-operator` chart, you must specify the following
additional set value for `helm` to install subchart `prometheusoperator`.

```bash
--set prometheusoperator.enabled=true
```

All `coherence` charts installed in `coherence-operator` `targetNamespaces` are monitored by 
Prometheus. The servicemonitor `<releasename>-coherence-service-monitor` 
configures Prometheus to scrape all components of `coherence-service`.


## Port forward Grafana

Once you have installed the charts, use the following script to port forward the Grafana pod.

```bash
#!/bin/bash

trap "exit" INT
  
while :
do
  kubectl port-forward $(kubectl get pods --selector=app=grafana -n namespace --output=jsonpath="{.items..metadata.name}") -n namespace 3000:3000
done
```

> **Note:** We add place the port-forward in a while to ensure it restarts any time it exists as 
> port-forwarding is sometimes unreliable and should only be used as a development tool. 

## Login to Grafana

In browser, go to the url `http://127.0.0.1:3000/d/coh-main/coherence-dashboard-main` to access the main Coherence dashboard.

At the Grafana login screen, the login is `admin` and the password is `prom-operator`.

Click `Home` in the upper left corner of screen to get a list of preconfigured dashboards.
Click ` Coherence Dashboard Main`.


## Default dashboards

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

## Navigating the dashboards

The Grafana dashboards created to monitor Coherence Clusters have some common UI elements and navigation patterns:

1. Variables and Annotations

   At the top left under the dashboard name any [variables](https://grafana.com/docs/reference/templating/), which are changeable and affect the
   queries in the dashboards, are displayed. Also [annotations](https://grafana.com/docs/reference/annotations/), which
   indicate events on the dashboard are also able to be enabled or disabled.
   
   ![Variables and Annotations](img/variables-and-annotations.png)
   
   `ClusterName` is a common variable which can be changed to choose the cluster do display information for.
   
   `Show Cluster Size Changed` is an annotation which shows anytime the cluster size has changed. All
   annotations appear as a red vertical line as shown below:
   
   ![Show Cluster Size Changed Annotation](img/annotation.png)

1. Access other dashboards

   On the right of the page you can click to show all the dashboards available for viewing.
   
   ![All Dashboards](img/all-dashboards.png)
   
## Configure your Prometheus Operator to scrape Coherence pods

This section assumes that you do not want the coherence-operator's helm subchart PrometheusOperator installed.
It provides information on how to configure what is automated by using coherence-operator helm chart parameter
`prometheusoperator.enabled`=`true`. 

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

Note that there are a number of Coherence Grafana dashboards bundled in the coherence-operator under dashboards.
While Grafana will have to be configured to the location of your prometheus datasource, one can still take advantage
of these Coherence dashboards by extracting them from the coherence-operator helm chart. 
    
## Troubleshooting

## Helm install of `coherence-operator` fails creating a custom resource definition (CRD).

Follow recommendation from [Prometheus Operator: helm fails to create CRDs](https://github.com/helm/charts/tree/master/stable/prometheus-operator#user-content-helm-fails-to-create-crds)
to manually install the Prometheus Operator CRDs, then install the `coherence-operator` chart with these additional set values. 

```bash
--set prometheusoperator.enabled=true --set prometheusoperator.prometheusOperator.createCustomResource=false
```

### No datasource found

Manually create a datasource by clicking on Grafana Home `Create your first data source` button 
and fill in these fields. Ensure the datasource is set as the default.
  
```bash
   Name:      Prometheus 
   HTTP URL:  http://prometheus-operated.<namespace>.svc.cluster.local:9090
```

Click `Save & Test` button on bottom of page.
