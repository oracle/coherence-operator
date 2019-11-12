<doc-view>

<v-layout row wrap>
<v-flex xs12 sm10 lg10>
<v-card class="section-def" v-bind:color="$store.state.currentColor">
<v-card-text class="pa-3">
<v-card class="section-def__card">
<v-card-text>
<dl>
<dt slot=title>Using Your Own Elasticsearch</dt>
<dd slot="desc"><p>The Coherence Operator can be configured to instruct Fluentd to push logs to a separate Elasticsearch instance rather thatn the in-built one.</p>
</dd>
</dl>
</v-card-text>
</v-card>
</v-card-text>
</v-card>
</v-flex>
</v-layout>

<h2 id="_pushing_logs_to_your_own_elasticsearch_instance">Pushing logs to your own Elasticsearch instance</h2>
<div class="section">
<p>This example shows how to instruct Fluentd to push data to your own Elasticsearch instance.</p>


<h3 id="install">1. Install the Coherence Operator with custom Elasticsearch endpoint</h3>
<div class="section">
<p>To enable an different Elasticsearch endpoint, add the following options to the Operator Helm install command:</p>

<markup
lang="bash"

>--set elasticsearchEndpoint.host=your-es-host
--set elasticsearchEndpoint.port=your-es-host</markup>

<p>You can also set the user and password if you Elasticsearch instance requires it:</p>

<markup
lang="bash"

>--set elasticsearchEndpoint.user=user
--set elasticsearchEndpoint.password=password</markup>

<div class="admonition note">
<p class="admonition-inline">For this example we have used the Stable ELK Stack at <a id="" title="" target="_blank" href="https://github.com/helm/charts/tree/master/stable/elastic-stack">https://github.com/helm/charts/tree/master/stable/elastic-stack</a>
to install the required components and have Elasticsearch runing on coherence-example-elastic-stack.default.svc.cluster.local:9200</p>
</div>
<p>A more complete helm install command to enable Prometheus is as follows:</p>

<markup
lang="bash"

>helm install \
    --namespace &lt;namespace&gt; \
    --name coherence-operator \
    --set elasticsearchEndpoint.host=coherence-example-elastic-stack.default.svc.cluster.local \
    --set elasticsearchEndpoint.port=9200 \
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
prometheus-coherence-operator-prometh-prometheus-0       3/3     Running   0          8m</markup>

<div class="admonition note">
<p class="admonition-inline">You will notice that there are no Kibana and Elasticsearch pods.</p>
</div>
</div>

<h3 id="install-coh">2. Install a Coherence Cluster with Logging Enabled</h3>
<div class="section">
<div class="admonition note">
<p class="admonition-inline">From this point on there is no difference in installation from when EFK is installed by the Coherence Operator.
This is because when Coherence is installed it will querying the Coherence Operator to receive the new Elasticsearch endpoint.
.</p>
</div>
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

</div>

<h3 id="_3_inspect_the_fluentd_container_logs">3. Inspect the Fluentd container logs</h3>
<div class="section">
<p>Issue the following to view the logs for the Fluentd container on the first Pod:</p>

<markup
lang="bash"

>kubectl logs -n &lt;namespace&gt; logging-cluster-storage-0 -c fluentd</markup>

<p>In the output you will see something similar to the following indicating your Fluentd container
will send data to your own Elasticsearch.</p>

<markup
lang="bash"

> &lt;match coherence-cluster&gt;
    @type elasticsearch
    host "coherence-example-elastic-stack.default.svc.cluster.local"
    port 9020
    user ""
    password xxxxxx
    logstash_format true
    logstash_prefix "coherence-cluster"
  &lt;/match&gt;</markup>

</div>

<h3 id="_4_connect_to_your_kibana_ui">4. Connect to your Kibana UI</h3>
<div class="section">
<p>Connect to your Kibana UI and create an index pattern called <code>coherence-cluster-*</code> to view the
incoming logs.</p>

</div>

<h3 id="_5_clean_up">5. Clean Up</h3>
<div class="section">
<p>After running the above the Coherence cluster can be removed using <code>kubectl</code>:</p>

<markup
lang="bash"

>kubectl -n &lt;namespace&gt; delete -f logging-cluster.yaml</markup>

</div>
</div>
</doc-view>
