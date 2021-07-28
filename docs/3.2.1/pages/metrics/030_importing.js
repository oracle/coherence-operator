<doc-view>

<v-layout row wrap>
<v-flex xs12 sm10 lg10>
<v-card class="section-def" v-bind:color="$store.state.currentColor">
<v-card-text class="pa-3">
<v-card class="section-def__card">
<v-card-text>
<dl>
<dt slot=title>Importing Grafana Dashboards</dt>
<dd slot="desc"><p>The Operator has a set of Grafana dashboards that can be imported into a Grafana instance.</p>

<div class="admonition note">
<p class="admonition-inline">Note: Use of metrics is available only when using the operator with clusters running
Coherence 12.2.1.4 or later version.</p>
</div></dd>
</dl>
</v-card-text>
</v-card>
</v-card-text>
</v-card>
</v-flex>
</v-layout>

<h2 id="_obtain_the_coherence_dashboards">Obtain the Coherence Dashboards</h2>
<div class="section">
<p>The Coherence Operator provides a set of dashboards for Coherence that may be imported into Grafana.
There are two ways to obtain the dashboards:</p>

<ul class="ulist">
<li>
<p>Clone the Coherence Operator GitHub repo, checkout the branch or tag for the version you want to use and
then obtain the dashboards from the <code>dashboards/</code> directory.</p>

</li>
<li>
<p>Download the <code>.tar.gz</code> dashboards package for the release you want to use.</p>

</li>
</ul>
<markup
lang="bash"

>curl https://oracle.github.io/coherence-operator/dashboards/latest/coherence-dashboards.tar.gz \
    -o coherence-dashboards.tar.gz
tar -zxvf coherence-dashboards.tar.gz</markup>

<p>The above commands will download the <code>coherence-dashboards.tar.gz</code> file and unpack it resulting in a
directory named <code>dashboards/</code> in the current working directory. This <code>dashboards/</code> directory will contain
the various Coherence dashboard files.</p>

</div>

<h2 id="_importing_grafana_dashboards">Importing Grafana Dashboards.</h2>
<div class="section">
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
<p ><p>There are three sets of dashboards available</p>

<ul class="ulist">
<li>
<p>Default - these are dashboards under the <code>dashboards/grafana/</code> directory that are compatible with
the metric names produced by the Coherence metrics publisher</p>

</li>
<li>
<p>Microprofile - these are dashboards under the <code>dashboards/grafana-microprofile/</code> directory that are compatible with
the metric names produced by the Coherence MP Metrics module.</p>

</li>
<li>
<p>Micrometer - these are dashboards under the <code>dashboards/grafana-micrometer/</code> directory that are compatible with
the metric names produced by the Coherence Micrometer Metrics module. These are metrics commonly used with Microprofile applications such as Helidon.</p>

</li>
<li>
<p>Micrometer - these are dashboards under the <code>dashboards/grafana-micrometer/</code> directory that are compatible with the metric names produced by the Coherence Micrometer Metrics module. Micrometer is a common metrics framework used with applications such as Spring, Micronaut etc.</p>

</li>
</ul>
<p>If you do not see metrics on the dashboards as expected you might be using the wrong dashboards version for how
Coherence has been configured.
The simplest way to find out which version corresponds to your Coherence cluster
is to query the metrics endpoint with something like <code>curl</code>.
If the metric names are in the format <code>vendor:coherence_cluster_size</code>, i.e. prefixed with <code>vendor:</code> then this is
the default Coherence format.
If metric names are in the format <code>vendor_Coherence_Cluster_Size</code>, i.e. prefixed with <code>vendor_</code> then this is
Microprofile format.</p>
</p>
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
</div>

<h2 id="_automatically_importing_dashboards">Automatically Importing Dashboards</h2>
<div class="section">
<p>There are ways to automatically import dashboards into Grafana, the exact method would depend on how Grafana is to
be installed and run.
The Coherence Operator test pipeline, for example, uses the
<a id="" title="" target="_blank" href="https://github.com/coreos/prometheus-operator">Prometheus Operator</a>
to install Prometheus and Grafana and automatically imports the Coherence dashboards from a <code>ConfigMap</code>.<br>
The examples below show how to create the dashboards as a <code>ConfigMap</code> and then install them into a Grafana
instances started with the Prometheus Operator.</p>

<p>There are two ways to create the <code>ConfigMap</code> containing the dashboard files:</p>

<ul class="ulist">
<li>
<p>Use the <code>ConfigMap</code> yaml available directly from GitHub</p>

</li>
<li>
<p>Obtain the dashboards as described above and create a <code>ConfigMap</code> from those files.</p>

</li>
</ul>

<h3 id="_create_a_configmap_from_github_yaml">Create a ConfigMap from GitHub Yaml</h3>
<div class="section">
<p>To create a <code>ConfigMap</code> with the Grafana dashboards directly from the <code>ConfigMap</code> yaml for a specific Operator release
the following commands can be used:</p>

<markup
lang="bash"

>kubectl -n monitoring create \
    -f https://oracle.github.io/coherence-operator/dashboards/latest/coherence-grafana-dashboards.yaml</markup>

<p>In this example the dashboards will be installed into the <code>monitoring</code> namespace.</p>

<p>The example above installs the dashboards configured to use the default Coherence metrics format.
Coherence provides integrations with Microprofile metrics and <a id="" title="" target="_blank" href="https://micrometer.io">Micrometer</a> metrics, which
produce metrics with slightly different name formats.
The operator provides dashboards compatible with both of these formats.</p>

<ul class="ulist">
<li>
<p>Microprofile change the URL to <code>coherence-grafana-microprofile-dashboards.yaml</code>, for example:</p>

</li>
</ul>
<markup
lang="bash"

>kubectl -n monitoring create \
    -f https://oracle.github.io/coherence-operator/dashboards/latest/coherence-grafana-microprofile-dashboards.yaml</markup>

<ul class="ulist">
<li>
<p>Micrometer change the URL to <code>coherence-grafana-micrometer-dashboards.yaml</code>, for example:</p>

</li>
</ul>
<markup
lang="bash"

>kubectl -n monitoring create \
    -f https://oracle.github.io/coherence-operator/dashboards/latest/coherence-grafana-micrometer-dashboards.yaml</markup>

</div>

<h3 id="_create_a_configmap_from_the_dashboard_package_file">Create a ConfigMap from the Dashboard Package File</h3>
<div class="section">
<p>To create a <code>ConfigMap</code> with the Grafana dashboards in directly from <code>.tar.gz</code> dashboard package for a specific
Operator release the following commands can be used:</p>

<markup
lang="bash"

>curl https://oracle.github.io/coherence-operator/dashboards/latest/coherence-dashboards.tar.gz \
    -o coherence-dashboards.tar.gz
tar -zxvf coherence-dashboards.tar.gz
kubectl -n monitoring create configmap coherence-grafana-dashboards --from-file=dashboards/grafana</markup>

<p>The <code>VERSION</code> variable has been set to the version of the dashboards to be used (this corresponds to an
Operator release version but dashboards can be used independently of the Operator).<br>
In this example the dashboards <code>ConfigMap</code> named <code>coherence-grafana-dashboards</code> will be installed into
the <code>monitoring</code> namespace.</p>

</div>

<h3 id="_label_the_configmap">Label the ConfigMap</h3>
<div class="section">
<p>In this example Grafana will be configured to import dashboards from <code>ConfigMaps</code> with the
label <code>grafana_dashboard</code>, so the <code>ConfigMap</code> created above needs to be labelled:</p>

<markup
lang="bash"

>kubectl -n monitoring label configmap coherence-grafana-dashboards grafana_dashboard=1</markup>

</div>

<h3 id="_install_the_prometheus_operator">Install the Prometheus Operator</h3>
<div class="section">
<p>The Prometheus Operator will be installed using its Helm chart.
Create a Helm values file like the following:</p>

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
<li data-value="1">Grafana will be enabled.</li>
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

<markup
lang="bash"

>GRAFANA_POD=$(kubectl -n monitoring get pod -l app.kubernetes.io/name=grafana -o name)
kubectl -n monitoring port-forward ${GRAFANA_POD} 3000:3000</markup>

<div class="admonition note">
<p class="admonition-inline">The default username for Grafana installed by the Prometheus Operator is <code>admin</code>
the default password is <code>prom-operator</code></p>
</div>
<p>If a Coherence cluster has been started with the Operator as described in the <router-link to="/metrics/020_metrics">Publish Metrics</router-link>
page, its metrics will eventually appear in Prometheus and Grafana. It can sometimes take a minute or so for
Prometheus to start scraping metrics and for them to appear in Grafana.</p>

</div>
</div>
</doc-view>
