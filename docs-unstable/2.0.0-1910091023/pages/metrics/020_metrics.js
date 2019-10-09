<doc-view>

<v-layout row wrap>
<v-flex xs12 sm10 lg10>
<v-card class="section-def" v-bind:color="$store.state.currentColor">
<v-card-text class="pa-3">
<v-card class="section-def__card">
<v-card-text>
<dl>
<dt slot=title>Enabling Metrics</dt>
<dd slot="desc"><p>Coherence clusters can be deployed with a metrics endpoint enabled that can be scraped by common metrics applications
such as Prometheus.</p>
</dd>
</dl>
</v-card-text>
</v-card>
</v-card-text>
</v-card>
</v-flex>
</v-layout>

<h2 id="_deploying_coherence_clusters_with_metrics_enabled">Deploying Coherence Clusters with Metrics Enabled</h2>
<div class="section">
<div class="admonition note">
<p class="admonition-inline">Note: Use of metrics is available only when using the operator with clusters running
Coherence 12.2.1.4 or later version.</p>
</div>
<p>The Coherence Operator can be installed with a demo Prometheus installation using embedded Prometheus Operator and
Grafana Helm charts. This Prometheus deployment is not intended for production use but is useful for development,
testing and demo purposes.</p>


<h3 id="_1_install_the_coherence_operator_with_prometheus">1. Install the Coherence Operator with Prometheus</h3>
<div class="section">
<p>To enable Prometheus, add the following options to the Operator Helm install command:</p>

<markup
lang="bash"

>--set prometheusoperator.enabled=true
--set prometheusoperator.prometheusOperator.createCustomResource=false</markup>

<p>A more complete helm install command to enable Prometheus is as follows:</p>

<markup
lang="bash"

>helm install \
    --namespace &lt;namespace&gt; \                                                  <span class="conum" data-value="1" />
    --name coherence-operator \
    --set prometheusoperator.enabled=true \
    --set prometheusoperator.prometheusOperator.createCustomResource=false \
    coherence/coherence-operator</markup>

<ul class="colist">
<li data-value="1">Set <code>&lt;namespace&gt;</code> to the Kubernetes namespace that the Coherence Operator should be installed into.</li>
</ul>
<p>After the installation completes, list the pods in the namespace that the Operator was installed into:</p>

<markup
lang="bash"

>kubectl -n &lt;namespace&gt; get pods</markup>

<p>The results returned should look something like the following:</p>

<markup
lang="bash"

>NAME                                                   READY   STATUS    RESTARTS   AGE
operator-coherence-operator-5d779ffc7-7xz7j            1/1     Running   0          53s  <span class="conum" data-value="1" />
operator-grafana-9d7fc9486-46zb7                       2/2     Running   0          53s  <span class="conum" data-value="2" />
operator-kube-state-metrics-7b4fcc5b74-ljdf8           1/1     Running   0          53s
operator-prometheus-node-exporter-kwdr7                1/1     Running   0          53s
operator-prometheusoperato-operator-77c784b8c5-v4bfz   1/1     Running   0          53s
prometheus-operator-prometheusoperato-prometheus-0     3/3     Running   2          38s  <span class="conum" data-value="3" /></markup>

<ul class="colist">
<li data-value="1">The Coherence Operator <code>Pod</code></li>
<li data-value="2">The Grafana <code>Pod</code></li>
<li data-value="3">The Prometheus <code>Pod</code></li>
</ul>
<p>The demo install of Prometheus in the Operator configures Prometheus to use service monitors to work out which Pods
to scrape metrics from. A <code>ServiceMonitor</code> in Prometheus will scrape from a port defined in a Kubernetes <code>Service</code> from
all <code>Pods</code> that match that service&#8217;s selector.</p>

</div>

<h3 id="_2_install_a_coherence_cluster_with_metrics_enabled">2. Install a Coherence Cluster with Metrics Enabled</h3>
<div class="section">
<p>Now that Prometheus is running Coherence clusters can be created that expose metrics on a port on each <code>Pod</code> and also
deploy a <code>Service</code> to expose the metrics that Prometheus can use.</p>

<p>Deploy a simple metrics enabled <code>CoherenceCluster</code> resource with a single role like this:</p>

<markup
lang="yaml"
title="metrics-cluster.yaml"
>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  role: storage        <span class="conum" data-value="1" />
  replicas: 2
  coherence:
    metrics:
      enabled: true    <span class="conum" data-value="2" />
  ports:
    - name: metrics    <span class="conum" data-value="3" />
      port: 9612       <span class="conum" data-value="4" /></markup>

<ul class="colist">
<li data-value="1">This cluster will have a single role called <code>storage</code></li>
<li data-value="2">The cluster will have two replicas (<code>Pods</code>)</li>
<li data-value="3">The Coherence <code>Pod</code> spec contains a port spec for metrics named <code>metric</code> so this needs to be exposed as a
service by specifying the <code>metrics</code> port in the role&#8217;s <code>ports</code> list</li>
<li data-value="4">The port must be set as <code>9612</code> which is the port that Coherence will expose metrics on</li>
</ul>
<p>The yaml above can be installed into Kubernetes using <code>kubectl</code>:</p>

<markup
lang="bash"

>kubectl -n &lt;namespace&gt; create -f metrics-cluster.yaml</markup>

<p>The Coherence Operator will see the new <code>CoherenceCluster</code> resource and create the cluster with two <code>Pods</code>.
If <code>kubectl get pods -n &lt;namespace&gt;</code> is run again it should now look something like this:</p>

<markup
lang="bash"

>NAME                                                   READY   STATUS    RESTARTS   AGE
operator-coherence-operator-5d779ffc7-7xz7j            1/1     Running   0          53s
operator-grafana-9d7fc9486-46zb7                       2/2     Running   0          53s
operator-kube-state-metrics-7b4fcc5b74-ljdf8           1/1     Running   0          53s
operator-prometheus-node-exporter-kwdr7                1/1     Running   0          53s
operator-prometheusoperato-operator-77c784b8c5-v4bfz   1/1     Running   0          53s
prometheus-operator-prometheusoperato-prometheus-0     3/3     Running   2          38s
test-cluster-storage-0                                 1/1     Running   0          70s  <span class="conum" data-value="1" />
test-cluster-storage-1                                 1/1     Running   0          70s  <span class="conum" data-value="2" /></markup>

<ul class="colist">
<li data-value="1"><code>Pod</code> one of the Coherence cluster.</li>
<li data-value="2"><code>Pod</code> two of the Coherence cluster.</li>
</ul>
<p>If the services are listed for the namespace:</p>

<markup
lang="bash"

>kubectl -n &lt;namespace&gt; get svc</markup>

<p>The list of services will look something like this.</p>

<markup
lang="bash"

>NAME                                    TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)     AGE
operator-grafana                        ClusterIP   10.104.251.51    &lt;none&gt;        80/TCP      31m
operator-kube-state-metrics             ClusterIP   10.110.18.78     &lt;none&gt;        8080/TCP    31m
operator-prometheus-node-exporter       ClusterIP   10.102.181.6     &lt;none&gt;        9100/TCP    31m
operator-prometheusoperato-operator     ClusterIP   10.107.59.229    &lt;none&gt;        8080/TCP    31m
operator-prometheusoperato-prometheus   ClusterIP   10.99.208.18     &lt;none&gt;        9090/TCP    31m
prometheus-operated                     ClusterIP   None             &lt;none&gt;        9090/TCP    31m
test-cluster-storage-headless           ClusterIP   None             &lt;none&gt;        30000/TCP   16m
test-cluster-storage-metrics            ClusterIP   10.109.201.211   &lt;none&gt;        9612/TCP    16m  <span class="conum" data-value="1" />
test-cluster-wka                        ClusterIP   None             &lt;none&gt;        30000/TCP   16m</markup>

<ul class="colist">
<li data-value="1">One of the services will be the service exposing the Coherence metrics.
The service name is typically in the format <code>&lt;cluster-name&gt;-&lt;role-name&gt;-&lt;port-name&gt;</code></li>
</ul>
<p>The Prometheus <code>ServiceMonitor</code> installed by the Coherence Operator is configured to look for services with the
label <code>component=coherence-service-metrics</code>. When ports are exposed in a <code>CoherenceCluster</code>, as has been done here
for metrics, the service created will have a label of the format <code>component=coherence-service-&lt;port-name&gt;</code>, so in
this case the <code>test-cluster-storage-metrics</code> service above will have the label <code>component=coherence-service-metrics</code>.</p>

<p>The labels for the service can be displayed:</p>

<markup
lang="bash"

>kubectl -n &lt;namespace&gt;&gt; get svc/test-cluster-storage-metrics --label-columns=component</markup>

<markup
lang="bash"

>NAME                           TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)    AGE   COMPONENT
test-cluster-storage-metrics   ClusterIP   10.109.201.211   &lt;none&gt;        9612/TCP   26m   coherence-service-metrics</markup>

<p>Which shows that the service does indeed have the required label.</p>

</div>

<h3 id="_3_access_prometheus">3. Access Prometheus</h3>
<div class="section">
<p>Now that Prometheus is running and is able to scrape metrics from the Coherence cluster it should be possible to access
those metrics in Prometheus.</p>

<p>First find the Prometheus <code>Pod</code> name using <code>kubectl</code></p>

<markup
lang="bash"

>kubectl -n &lt;namespace&gt; get pod -l app=prometheus -o name</markup>

<p>Using the <code>Pod</code> name use <code>kubectl</code> to create a port forward session to the Prometheus <code>Pod</code> so that the
Prometheus API on port <code>9090</code> in the <code>Pod</code> can be accessed from the local host.</p>

<markup
lang="bash"

>kubectl -n &lt;namespace&gt; port-forward \
    $(kubectl -n &lt;namespace&gt; get pod -l app=prometheus -o name) \
    9090:9090</markup>

<p>It is now possible to access the Prometheus API on localhost port 9090. This can be used to directly retrieve
Coherence metrics using <code>curl</code>, for example to obtain the cluster size metric:</p>

<markup
lang="bash"

>curl -w '\n' -X GET http://127.0.0.1:9090/api/v1/query?query=vendor:coherence_cluster_size</markup>

<p>It is also possible to browse directly to the Prometheus web UI at <a id="" title="" target="_blank" href="http://127.0.0.1:9090">http://127.0.0.1:9090</a></p>

</div>

<h3 id="_3_access_grafana">3. Access Grafana</h3>
<div class="section">
<p>By default when the Coherence Operator configured to install Prometheus the Prometheus Operator also install a
Grafana <code>Pod</code> and the Coherence Operator imports into Grafana a number of custom dashboards for displaying Coherence
metrics. Grafana can be accessed by using port forwarding in the same way that was done for Prometheus</p>

<p>First find the Grafana <code>Pod</code>:</p>

<markup
lang="bash"

>kubectl -n &lt;namespace&gt; get pod -l app=grafana -o name</markup>

<p>Using the <code>Pod</code> name use <code>kubectl</code> to create a port forward session to the Grafana <code>Pod</code> so that the
Grafana API on port <code>3000</code> in the <code>Pod</code> can be accessed from the local host.</p>

<markup
lang="bash"

>kubectl -n &lt;namespace&gt; port-forward \
    $(kubectl -n &lt;namespace&gt; get pod -l app=grafana -o name) \
    3000:3000</markup>

<p>The custom Coherence dashboards can be accessed by pointing a browser to
<a id="" title="" target="_blank" href="http://127.0.0.1:3000/d/coh-main/coherence-dashboard-main">http://127.0.0.1:3000/d/coh-main/coherence-dashboard-main</a></p>

<p>The Grafana credentials are username <code>admin</code> password <code>prom-operator</code></p>

</div>

<h3 id="_4_cleaning_up">4. Cleaning Up</h3>
<div class="section">
<p>After running the demo above the Coherence cluster can be removed using <code>kubectl</code>:</p>

<markup
lang="bash"

>kubectl -n &lt;namespace&gt; delete -f metrics-cluster.yaml</markup>

<p>The Coherence Operator, along with Prometheus and Grafana can be removed using Helm:</p>

<markup
lang="bash"

>helm delete --purge coherence-operator</markup>

</div>
</div>
</doc-view>
