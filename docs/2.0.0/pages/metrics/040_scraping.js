<doc-view>

<v-layout row wrap>
<v-flex xs12 sm10 lg10>
<v-card class="section-def" v-bind:color="$store.state.currentColor">
<v-card-text class="pa-3">
<v-card class="section-def__card">
<v-card-text>
<dl>
<dt slot=title>Using Your Own Prometheus</dt>
<dd slot="desc"><p>If required, you can scrape the metrics from your own Prometheus Operator instance rather
than using the <code>prometheusopeartor</code> subchart included with the Coherence Operator.</p>
</dd>
</dl>
</v-card-text>
</v-card>
</v-card-text>
</v-card>
</v-flex>
</v-layout>

<h2 id="_scraping_metrics_from_your_own_prometheus_instance">Scraping metrics from your own Prometheus instance</h2>
<div class="section">
<div class="admonition note">
<p class="admonition-inline">Note: Use of metrics is available only when using the operator with clusters running
Coherence 12.2.1.4 or later version.</p>
</div>
<p>This example shows you how to scrape metrics from your own Prometheus instance.</p>


<h3 id="install">1. Install the Coherence Operator with Prometheus disabled</h3>
<div class="section">
<p>A more complete helm install command to enable Prometheus is as follows:</p>

<markup
lang="bash"

>helm install \
    --namespace &lt;namespace&gt; \  <span class="conum" data-value="1" />
    --name coherence-operator \
    coherence/coherence-operator</markup>

<ul class="colist">
<li data-value="1">Set <code>&lt;namespace&gt;</code> to the Kubernetes namespace that the Coherence Operator should be installed into.</li>
</ul>
<p>After the installation completes, list the pods in the namespace that the Operator was installed into:</p>

<markup
lang="bash"

>kubectl -n &lt;namespace&gt; get pods</markup>

<p>The results returned should only show the Coherence Operator.</p>

<markup
lang="bash"

>NAME                                                   READY   STATUS    RESTARTS   AGE
operator-coherence-operator-5d779ffc7-7xz7j            1/1     Running   0          53s</markup>

</div>

<h3 id="_2_install_prometheus_operator_optional">2. Install Prometheus Operator (Optional)</h3>
<div class="section">
<p>Id you do not already have a Prometheus environment installed, you can use the <code>Prometheus Operator</code>
chart from <a id="" title="" target="_blank" href="https://github.com/helm/charts/tree/master/stable/prometheus-operator">https://github.com/helm/charts/tree/master/stable/prometheus-operator</a> using the following:</p>

<markup
lang="bash"

>helm install stable/prometheus-operator --namespace &lt;namespace&gt; --name prometheus \
      --set prometheusOperator.createCustomResource=false</markup>

</div>

<h3 id="_3_create_a_servicemonitor">3. Create a ServiceMonitor</h3>
<div class="section">
<p>Create a ServiceMonitor with the following configuration to instruct Prometheus to scrape the Coherence Pods.</p>

<markup
lang="yaml"
title="service-monitor.yaml"
>apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: monitoring-coherence
  namespace: coherence-example
  labels:
    release: prometheus         <span class="conum" data-value="1" />
spec:
  selector:
    matchLabels:
      component: coherencePod   <span class="conum" data-value="2" />
  endpoints:
  - port: metrics               <span class="conum" data-value="3" />
    interval: 30s
  namespaceSelector:
    matchNames:
      - coherence-example</markup>

<ul class="colist">
<li data-value="1">Match the Prometheus Operator release name</li>
<li data-value="2">Scrape all Pods that match component=<code>coherencePod</code></li>
<li data-value="3">The metrics Pod to Scrape</li>
</ul>
<p>The yaml above can be installed into Kubernetes using <code>kubectl</code>:</p>

<markup
lang="bash"

>kubectl create -n &lt;namespace&gt;-f service-monitor.yaml</markup>

<p>See the <a id="" title="" target="_blank" href="https://github.com/coreos/prometheus-operator">Prometheus Operator</a> documentation
for more information on ServiceMonitors usage.</p>

</div>

<h3 id="install-coh">4. Install a Coherence Cluster with Metrics Enabled</h3>
<div class="section">
<p>Now that Prometheus is running Coherence clusters can be created that expose metrics on a port on each <code>Pod</code>.</p>

<p>Deploy a simple metrics enabled <code>CoherenceCluster</code> resource with a single role like this:</p>

<markup
lang="yaml"
title="metrics-cluster.yaml"
>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: metrics-cluster
spec:
  role: storage
  replicas: 3
  coherence:
    metrics:
      enabled: true
      port: 9612
  ports:
    - name: metrics
      port: 9612</markup>

<p>The yaml above can be installed into Kubernetes using <code>kubectl</code>:</p>

<markup
lang="bash"

>kubectl -n &lt;namespace&gt; create -f metrics-cluster.yaml</markup>

<markup
lang="bash"

>kubectl -n &lt;namespace&gt; get pods

NAME                                                     READY   STATUS    RESTARTS   AGE
alertmanager-prometheus-prometheus-oper-alertmanager-0   2/2     Running   0          51m
coherence-operator-8465cf7d88-7hw4g                      1/1     Running   0          81m
metrics-cluster-storage-0                                1/1     Running   0          12m
metrics-cluster-storage-1                                1/1     Running   0          12m
metrics-cluster-storage-2                                1/1     Running   0          12m
prometheus-grafana-757f7c9f6d-brqvb                      2/2     Running   0          51m
prometheus-kube-state-metrics-5ffdf76ddd-86qg4           1/1     Running   0          51m
prometheus-prometheus-node-exporter-4d9qx                1/1     Running   0          51m
prometheus-prometheus-oper-operator-64cd6c6c45-5ql2k     2/2     Running   0          51m
prometheus-prometheus-prometheus-oper-prometheus-0       3/3     Running   1          51m</markup>

</div>

<h3 id="_5_validate_that_prometheus_can_see_the_pods">5. Validate that Prometheus can see the Pods</h3>
<div class="section">
<p>Port-forward the Prometheus port using the following:</p>

<markup
lang="bash"

>kubectl -n &lt;namespace&gt; port-forward prometheus-prometheus-prometheus-oper-prometheus-0 9090:9090

Forwarding from 127.0.0.1:9090 -&gt; 9090
Forwarding from [::1]:9090 -&gt; 9090</markup>

<p>Access the following endpoint to confirm that Prometheus is scraping the pods:</p>

<p><a id="" title="" target="_blank" href="http://127.0.0.1:9090/targets">http://127.0.0.1:9090/targets</a></p>

<p>The following should be displayed indicating the Coherence Pods are being scraped.</p>



<v-card>
<v-card-text class="overflow-y-hidden" style="text-align:center">
<img src="./images/prometheus-targets.png" alt="Prometheus Targets"width="950" />
</v-card-text>
</v-card>

</div>

<h3 id="_6_access_grafana_and_load_the_dashboards">6. Access Grafana and Load the Dashboards</h3>
<div class="section">
<p>Port-forward the Grafana port using the following, replacing the grafana pod</p>

<markup
lang="bash"

>kubectl port-forward $(kubectl get pod -n &lt;namespace&gt; -l app=grafana -o name) -n &lt;namespace&gt; 3000:3000

Forwarding from 127.0.0.1:9090 -&gt; 9090
Forwarding from [::1]:9090 -&gt; 9090</markup>

<p>Access Grafana via the following URL: <a id="" title="" target="_blank" href="http://127.0.0.1:3000/">http://127.0.0.1:3000/</a></p>

<div class="admonition note">
<p class="admonition-inline">The Grafana credentials are username <code>admin</code> password <code>prom-operator</code></p>
</div>
<p>Once logged in, highlight the <code>+</code> icon and select <code>Import</code>.</p>

<p>Import all of the Grafana dashboards from the following location: <a id="" title="" target="_blank" href="https://github.com/oracle/coherence-operator/helm-charts/coherence-operator/dashboards">https://github.com/oracle/coherence-operator/helm-charts/coherence-operator/dashboards</a></p>

<p>Once all the dashboards have been loaded, you can access the Main dasbohard via: <a id="" title="" target="_blank" href="http://127.0.0.1:3000/d/coh-main/coherence-dashboard-main">http://127.0.0.1:3000/d/coh-main/coherence-dashboard-main</a></p>


<h4 id="_7_uninstall_the_coherence_cluster_prometheus">7. Uninstall the Coherence Cluster &amp; Prometheus</h4>
<div class="section">
<markup
lang="bash"

>kubectl delete -n &lt;namespace&gt; -f metrics-cluster.yaml

coherencecluster.coherence.oracle.com "metrics-cluster" deleted

helm delete prometheus --purge

release "prometheus" deleted</markup>

</div>
</div>
</div>
</doc-view>
