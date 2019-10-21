<doc-view>

<v-layout row wrap>
<v-flex xs12 sm10 lg10>
<v-card class="section-def" v-bind:color="$store.state.currentColor">
<v-card-text class="pa-3">
<v-card class="section-def__card">
<v-card-text>
<dl>
<dt slot=title>Enabling Log Capture</dt>
<dd slot="desc"><p>The Coherence Operator manages data logging through the Elasticsearch, Fluentd and Kibana (EFK) stack.</p>
</dd>
</dl>
</v-card-text>
</v-card>
</v-card-text>
</v-card>
</v-flex>
</v-layout>

<h2 id="_capturing_and_viewing_coherence_cluster_logs">Capturing and viewing Coherence cluster Logs</h2>
<div class="section">
<p>This example shows how to enable log capture and access the Kibana user interface (UI) to view the captured logs.</p>

<p>Logs are scraped via a Fluentd sidecar image, parsed and sent to Elasticsearch. A default
index pattern called <code>coherence-cluster-*</code> is created which holds all captured logs.</p>


<h3 id="install">1. Install the Coherence Operator with Fluentd logging enabled</h3>
<div class="section">
<p>To enable the EFK stack, add the following options to the Operator Helm install command:</p>

<markup
lang="bash"

>--set installEFK=true</markup>

<p>A more complete helm install command to enable Prometheus is as follows:</p>

<markup
lang="bash"

>helm install \
    --namespace &lt;namespace&gt; \
    --name coherence-operator \
    --set installEFK=true \
    coherence/coherence-operator</markup>

<p>After the installation completes, list the pods in the namespace that the Operator was installed into:</p>

<markup
lang="bash"

>kubectl -n &lt;namespace&gt; get pods</markup>

<p>The results returned should look something like the following:</p>

<markup
lang="bash"

>NAME                                                     READY   STATUS    RESTARTS   AGE
coherence-operator-66c6d868b9-rd429                      1/1     Running   0          8m
coherence-operator-grafana-8454698bcf-v5kxw              2/2     Running   0          8m
coherence-operator-kube-state-metrics-6dc8675d87-qnfdw   1/1     Running   0          8m
coherence-operator-prometh-operator-58d94ffbb8-94d4m     1/1     Running   0          8m
coherence-operator-prometheus-node-exporter-vpjjt        1/1     Running   0          8m
elasticsearch-f978d6fdd-dw7qg                            1/1     Running   0          8m   <span class="conum" data-value="1" />
kibana-9964496fd-5tpv9                                   1/1     Running   0          8m   <span class="conum" data-value="2" />
prometheus-coherence-operator-prometh-prometheus-0       3/3     Running   0          8m</markup>

<ul class="colist">
<li data-value="1">The Elasticsearch <code>Pod</code></li>
<li data-value="2">The Kibana <code>Pod</code></li>
</ul>
</div>

<h3 id="install-coh">2. Install a Coherence Cluster with Logging Enabled</h3>
<div class="section">
<p>Deploy a simple logging enabled <code>CoherenceCluster</code> resource with a single role like this:</p>

<markup
lang="yaml"
title="logging-cluster.yaml"
>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: logging-cluster
spec:
  role: storage
  replicas: 3
  logging:
    fluentd:
      enabled: true  <span class="conum" data-value="1" /></markup>

<ul class="colist">
<li data-value="1">Enables log capture via Fluentd</li>
</ul>
<p>The yaml above can be installed into Kubernetes using <code>kubectl</code>:</p>

<markup
lang="bash"

>kubectl -n &lt;namespace&gt; create -f logging-cluster.yaml

coherencecluster.coherence.oracle.com/logging-cluster created

kubectl -n &lt;namespace&gt; get pod -l coherenceCluster=logging-cluster

NAME                        READY   STATUS    RESTARTS   AGE
logging-cluster-storage-0   2/2     Running   0          86s
logging-cluster-storage-1   2/2     Running   0          86s
logging-cluster-storage-2   2/2     Running   0          86s</markup>

<div class="admonition note">
<p class="admonition-inline">Notice that under the <code>Ready</code> column it shows <code>2/2</code>. This means that there are two containers for this
Pod, Coherence and Fluentd, and they are both ready.  The Fluentd container will capture the logs, parse them
and send them to Elasticsearch. Kibana can then be used to view the logs.</p>
</div>
</div>

<h3 id="_3_port_forward_the_kibana_pod">3. Port-forward the Kibana pod</h3>
<div class="section">
<p>First find the Kibana <code>Pod</code>:</p>

<markup
lang="bash"

>kubectl -n coherence-example get pod -l component=kibana -o name</markup>

<p>Using the <code>Pod</code> name use <code>kubectl</code> to create a port forward session to the Kibana <code>Pod</code> so that the
Kibana API on port <code>5601</code> in the <code>Pod</code> can be accessed from the local host.</p>

<markup
lang="bash"

>kubectl -n &lt;namespace&gt; port-forward \
    $(kubectl -n &lt;namespace&gt; get pod -l component=kibana -o name) \
    5601:5601

Forwarding from [::1]:5601 -&gt; 5601
Forwarding from 127.0.0.1:5601 -&gt; 5601</markup>

</div>

<h3 id="_4_access_the_kibana_application_ui">4. Access the Kibana Application UI</h3>
<div class="section">
<p>Access Kibana using the following URL: <a id="" title="" target="_blank" href="http://127.0.0.1:5601/">http://127.0.0.1:5601/</a></p>

<div class="admonition note">
<p class="admonition-inline">It may take approximately 2-3 minutes for the first logs to reach the Elasticsearch instance.</p>
</div>

<h4 id="_default_dashboards">Default Dashboards</h4>
<div class="section">
<ul class="ulist">
<li>
<p>Coherence Cluster - All Messages : Shows all messages</p>

</li>
<li>
<p>Coherence Cluster - Errors and Warnings : Shows errors and warning messages only</p>

</li>
<li>
<p>Coherence Cluster - Persistence : Shows Persistence related messages</p>

</li>
<li>
<p>Coherence Cluster - Configuration Messages: Shows configuration related messages</p>

</li>
<li>
<p>Coherence Cluster - Network : Shows network related messages, such as communication delays and TCP ring disconnects</p>

</li>
<li>
<p>Coherence Cluster - Partitions : Shows partition transfer and loss messages</p>

</li>
<li>
<p>Coherence Cluster - Message Sources : Shows the source (thread) for messages</p>

</li>
</ul>
</div>

<h4 id="_default_queries">Default Queries</h4>
<div class="section">
<p>There are many searches related to common Coherence messages, warnings,
and errors that are loaded and can be accessed via the <code>Discover</code> <code>side-bar
and selecting `Open</code>.</p>

<p>See <router-link to="/logging/040_dashboards">here</router-link> for more information on the default dashboards and searches.</p>

</div>
</div>

<h3 id="_4_clean_up">4. Clean Up</h3>
<div class="section">
<p>After running the above the Coherence cluster can be removed using <code>kubectl</code>:</p>

<markup
lang="bash"

>kubectl -n &lt;namespace&gt; delete -f logging-cluster.yaml</markup>

</div>
</div>
</doc-view>
