<doc-view>

<h2 id="_importing_grafana_dashboards">Importing Grafana Dashboards.</h2>
<div class="section">
<p>The Operator has a set of Grafana dashboards that can be imported into a Grafana instance.</p>

<div class="admonition note">
<p class="admonition-inline">Note: Use of metrics is available only when using the operator with clusters running
Coherence 12.2.1.4 or later version.</p>
</div>
<p>This example shows you how to import the Grafana dashboards into your own Grafana instance.</p>

<p>By default, the Coherence dashboards require a datasource named <code>Prometheus</code> which
should also be the default datasource.</p>

<p>If this datasource is already used, and you cannot add another datasource as the default,
then please go to <router-link to="#different" @click.native="this.scrollFix('#different')">Importing with a different datasource</router-link>.</p>


<h3 id="importing">Manually Importing Using the Defaults</h3>
<div class="section">
<p>In your Grafana environment, ensure you either:</p>

<ul class="ulist">
<li>
<p>have a Prometheus datasource named <code>Prometheus</code> which is also marked as the default datasource</p>

</li>
<li>
<p>have added a new Prometheus datasource which you have set as the default</p>

</li>
</ul>
<p>Then continue below.</p>

<ul class="ulist">
<li>
<p>Clone the git repository using</p>

</li>
</ul>
<div class="listing">
<pre>git clone https://github.com/oracle/coherence-operator.git</pre>
</div>

<div class="admonition note">
<p class="admonition-textlabel">Note</p>
<p ><p>There are two sets of dashboards available</p>

<ul class="ulist">
<li>
<p>Legacy - these are dashboards under the <code>dashboards/grafana-legacy/</code> directory that are compatible with
the metric names produced by the Coherence metrics publisher</p>

</li>
<li>
<p>Microprofile - these are dashboards under the <code>dashboards/grafana/</code> directory that are compatible with
the metric names produced by the Coherence MP Metrics module.</p>

</li>
</ul></p>
</div>
<ul class="ulist">
<li>
<p>Decide which dashboards you will import, depending on how metrics are being published (see the note above).</p>

</li>
<li>
<p>Login to Grafana and for each dashboard in the chosen dashboard directory carry out the
following to upload to Grafana:</p>
<ul class="ulist">
<li>
<p>Highlight the '+' (plus) icons and click on import</p>

</li>
<li>
<p>Click `Upload Json file' button to select a dashboard</p>

</li>
<li>
<p>Leave all the default values and click on <code>Import</code></p>

</li>
</ul>
</li>
</ul>
</div>

<h3 id="different">Manually Importing With a Different Datasource</h3>
<div class="section">
<p>If your Grafana environment has a default datasource that you cannot change or already has a
datasource named <code>Prometheus</code>, follow these steps to import the dashboards:</p>

<ul class="ulist">
<li>
<p>Login to Grafana</p>

</li>
<li>
<p>Create a new datasource named <code>Coherence-Prometheus</code> and point to your Prometheus endpoint</p>

</li>
<li>
<p>Create a temporary directory and copy all the dashboards from the cloned directory
<code>&lt;DIR&gt;/dashboards/grafana</code> to this temporary directory</p>

</li>
<li>
<p>Change to this temporary directory and run the following to update the datasource to <code>Coherence-Prometheus</code> or the
datasource of your own choice:</p>

</li>
</ul>
<div class="listing">
<pre>for file in *.json
do
    sed -i '' -e 's/"datasource": "Prometheus"/"datasource": "Coherence-Prometheus"/g' \
              -e 's/"datasource": null/"datasource": "Coherence-Prometheus"/g' \
              -e 's/"datasource": "-- Grafana --"/"datasource": "Coherence-Prometheus"/g' $file;
done</pre>
</div>

<ul class="ulist">
<li>
<p>Once you have completed the script, proceed to upload the dashboards as described above.</p>

</li>
</ul>
</div>

<h3 id="_automatically_importing_dashboards">Automatically Importing Dashboards</h3>
<div class="section">
<p>There are ways to automatically import dashboards into Grafana, the exact method would depend on how Grafana is to
be run.
The Coherence Operator test pipeline, for example, uses the
<a id="" title="" target="_blank" href="https://github.com/coreos/prometheus-operator">Prometheus Operator</a>
to install Prometheus and Grafana and automatically imports the Coherence dashboards from a <code>ConfigMap</code>.</p>

<p>The following commands create a <code>ConfigMap</code> named <code>coherence-grafana-dashboards</code> in the <code>monitoring</code> namespace
from the dashboard files in the <code>dashboards/grafana/</code> directory, then a label <code>grafana_dashboard=1</code> is added to
the <code>ConfigMap</code>:</p>

<markup
lang="bash"

>kubectl -n monitoring create configmap coherence-grafana-dashboards --from-file=dashboards/grafana/
kubectl -n monitoring label configmap coherence-grafana-dashboards grafana_dashboard=1</markup>

<p>The Prometheus Operator will be installed using its Helm chart.
Create a values file like the following:</p>

<markup
lang="yaml"
title="prometheus-values.yaml"
>prometheus:
  prometheusSpec:
    serviceMonitorSelectorNilUsesHelmValues: false
alertmanager:
  enabled: false
nodeExporter:
  enabled: true
grafana:
  enabled: true                   <span class="conum" data-value="1" />
  sidecar:
    dashboards:                   <span class="conum" data-value="2" />
      enabled: true
      label: grafana_dashboard</markup>

<ul class="colist">
<li data-value="1">Grafana is enabled.</li>
<li data-value="2">Grafana will automatically import dashboards from <code>ConfigMaps</code> that have the label <code>grafana_dashboard</code>
(which was given to the <code>ConfigMap</code> created above).</li>
</ul>
<p>Prometheus can be installed into the <code>monitoring</code> namespace using the Helm command:</p>

<markup
lang="bash"

>helm install --namespace monitoring \
    --values prometheus-values.yaml \
    prometheus stable/prometheus-operator</markup>

<p>To actually start Prometheus a <code>Prometheus</code> CRD resource needs to be added to Kubernetes.
Create a <code>Prometheus</code> resource yaml file suitable for testing:</p>

<markup
lang="yaml"
title="prometheus.yaml"
>apiVersion: monitoring.coreos.com/v1
kind: Prometheus
metadata:
  name: prometheus
spec:
  serviceAccountName: prometheus
  serviceMonitorSelector:
    matchLabels:
      coherenceComponent: coherence-service-monitor  <span class="conum" data-value="1" />
  resources:
    requests:
      memory: 400Mi
  enableAdminAPI: true</markup>

<ul class="colist">
<li data-value="1">The <code>serviceMonitorSelector</code> tells Prometheus to use any <code>ServiceMonitor</code> that is created with the
<code>coherence-service-monitor</code> label, which is a label that the Coherence Operator adds to any <code>ServiceMonitor</code>
that it creates.</li>
</ul>
<p>Install the <code>prometheus.yaml</code> file into Kubernetes:</p>

<markup
lang="bash"

>kubectl -n monitoring create -f etc/prometheus.yaml</markup>

<p>In the <code>monitoring</code> namespace there should now be a number of <code>Pods</code> and <code>Services</code>, among them a <code>Prometheus</code>
instance, and a Grafana instance. It should be possible to reach the Grafana UI on the ports exposed by the <code>Pod</code>
and see the imported Coherence dashboards.</p>

</div>
</div>
</doc-view>
