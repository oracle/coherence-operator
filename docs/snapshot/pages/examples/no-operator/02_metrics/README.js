<doc-view>

<h2 id="_enabling_coherence_metrics">Enabling Coherence Metrics</h2>
<div class="section">
<p>This example shows how to deploy a simple Coherence cluster in Kubernetes manually, and enabling the Pods in that cluster to expose a http endpoint to allow access to Coherence metrics.
This example expands on the <code>StatefulSet</code> used in the first simple deployment example.</p>

<div class="admonition tip">
<p class="admonition-textlabel">Tip</p>
<p ><p><img src="./images/GitHub-Mark-32px.png" alt="GitHub Mark 32px" />
 The complete source code for this example is in the <a id="" title="" target="_blank" href="https://github.com/oracle/coherence-operator/tree/main/examples/no-operator/02_metrics">Coherence Operator GitHub</a> repository.</p>
</p>
</div>
<p><strong>Prerequisites</strong></p>

<p>This example assumes that you have already built the example server image.</p>

</div>

<h2 id="_create_the_kubernetes_resources">Create the Kubernetes Resources</h2>
<div class="section">
<p>In the simple server example we created some <code>Services</code> and a <code>StatefulSet</code> that ran a Coherence cluster in Kubernetes.
In this example we will just cover the additional configurations we need to make to expose Coherence metrics.
We will not bother repeating the configuration for the <code>Services</code> for the <code>StatefulSet</code> and well known addressing or the Service for exposing Extend. We will assume they are already part of our yaml file.</p>

<p>The <code>coherence-metrics.yaml</code> file that is part of the source for this example contains all those resources.</p>


<h3 id="_the_statefulset">The StatefulSet</h3>
<div class="section">
<p>To expose Coherence metrics we just need to change the <code>StatefulSet</code> to set either the system properties or environment variables to enable metrics. We will also add a container port to the expose metrics endpoint.</p>

<markup
lang="yaml"
title="coherence.yaml"
>apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: storage
spec:
  selector:
    matchLabels:
      coherence.oracle.com/cluster: test-cluster
      coherence.oracle.com/deployment: storage
  serviceName: storage-sts
  replicas: 3
  template:
    metadata:
      labels:
        coherence.oracle.com/cluster: test-cluster
        coherence.oracle.com/deployment: storage
    spec:
      containers:
        - name: coherence
          image: simple-coherence:1.0.0
          command:
            - java
          args:
            - -cp
            - "@/app/jib-classpath-file"
            - -Xms1800m
            - -Xmx1800m
            - "@/app/jib-main-class-file"
          env:
            - name: COHERENCE_CLUSTER
              value: storage
            - name: COHERENCE_WKA
              value: storage-wka
            - name: COHERENCE_CACHECONFIG
              value: "test-cache-config.xml"
            - name: COHERENCE_METRICS_HTTP_ENABLED
              value: "true"
          ports:
            - name: extend
              containerPort: 20000
            - name: metrics
              containerPort: 9612</markup>

<p>The yaml above is identical to that used in the simple server example apart from:
* We added the <code>COHERENCE_METRICS_HTTP_ENABLED</code> environment variable with a value of <code>"true"</code>. Instead of this we could have added <code>-Dcoherence.metrics.http.enabled=true</code> to the <code>args:</code> list to set the <code>coherence.metrics.http.enabled</code> system property to true. Recent versions of Coherence work with both system properties or environment variables, and we just chose to use environment variables in this example.
* We added a port named <code>metrics</code> with a port value of <code>9612</code>, which is the default port that the Coherence metrics endpoint binds to.</p>

</div>
</div>

<h2 id="_deploy_the_cluster">Deploy the Cluster</h2>
<div class="section">
<p>We can combine all the snippets of yaml above into a single file and deploy it to Kubernetes.
The source code for this example contains a file named <code>coherence-metrics.yaml</code> containing all the configuration above.
We can deploy it with the following command:</p>

<markup
lang="bash"

>kubectl apply -f coherence-metrics.yaml</markup>

<p>We can see all the resources created in Kubernetes are the same as the simple server example:</p>

<markup
lang="bash"

>kubectl get all</markup>

<p>Which will display something like the following:</p>

<markup


>pod/storage-0   1/1     Running   0          10s
pod/storage-1   1/1     Running   0          7s
pod/storage-2   1/1     Running   0          6s

NAME                     TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)     AGE
service/kubernetes       ClusterIP   10.96.0.1        &lt;none&gt;        443/TCP     158d
service/storage-extend   ClusterIP   10.102.198.218   &lt;none&gt;        20000/TCP   10s
service/storage-sts      ClusterIP   None             &lt;none&gt;        7/TCP       10s
service/storage-wka      ClusterIP   None             &lt;none&gt;        7/TCP       10s

NAME                       READY   AGE
statefulset.apps/storage   3/3     10s</markup>

</div>

<h2 id="_retrieve_metrics">Retrieve Metrics</h2>
<div class="section">
<p>To test that we can access metrics we will port-forward to one of the <code>Pods</code> and use <code>curl</code> to get the metrics.
We can choose any of the three <code>Pods</code> to test, or repeat the test for each <code>Pod</code>.
In this example, we&#8217;ll just port-forward local port 9612 to port 9612 in <code>pod/storage-0</code>.</p>

<markup
lang="bash"

>kubectl port-forward pod/storage-0 9612:9612</markup>

<p>Now in another terminal we can run the <code>curl</code> command to get metrics. As we are using port forwarding the host will be <code>127.0.0.1</code> and the port will be <code>9612</code>.</p>

<markup
lang="bash"

>curl -X GET http://127.0.0.1:9612/metrics</markup>

<p>This should then bring back all the Coherence metrics for <code>pod/storage-0</code>. The default format of the response is Prometheus text format.</p>

<p>We can also retrieve individual metrics by name. For example, we can get the <code>Coherence.Cluster.Size</code> metric:</p>

<markup
lang="bash"

>curl -X GET http://127.0.0.1:9612/metrics/Coherence.Cluster.Size</markup>

<p>which will display something like this:</p>

<markup
lang="bash"

>vendor:coherence_cluster_size{cluster="storage", version="21.12.1"} 3</markup>

<p>This displays the metric name in Prometheus format <code>vendor:coherence_cluster_size</code>, the metric labels <code>cluster="storage", version="21.12.1"</code> and the metric value, in this case <code>3</code> as there are three cluster members because we specified a replicas value of 3 in the <code>StatefulSet</code>.</p>

<p>We can also receive the same response as <code>json</code> by using either the accepted media type header <code>"Accept: application/json"</code>:</p>

<markup
lang="bash"

>curl -X GET -H "Accept: application/json" http://127.0.0.1:9612/metrics/Coherence.Cluster.Size</markup>

<p>Or by using the <code>.json</code> suffix on the URL</p>

<markup
lang="bash"

>curl -X GET http://127.0.0.1:9612/metrics/Coherence.Cluster.Size.json</markup>

<p>Both requests will display something like this:</p>

<markup
lang="bash"

>[{"name":"Coherence.Cluster.Size","tags":{"cluster":"storage","version":"21.12.1"},"scope":"VENDOR","value":3}]</markup>

<p>We have now verified that the <code>Pods</code> in the cluster are producing metrics.</p>

</div>

<h2 id="_using_prometheus">Using Prometheus</h2>
<div class="section">
<p>One of the most common ways to analyse metrics in Kubernetes is by using Prometheus.
The recommended way to do this is to deploy Prometheus inside your Kubernetes cluster so that it can scrape metrics directly from <code>Pods</code>. Whilst Prometheus can be installed outside the Kubernetes cluster, this introduces a much more complicated set-up.
If using Prometheus externally to the Kubernetes cluster, the approach recommended by Prometheus is to use federation, which we show in an example below.</p>


<h3 id="_install_prometheus">Install Prometheus</h3>
<div class="section">
<p>The simplest way to install Prometheus is to follow the instructions in the Prometheus Operator
<a id="" title="" target="_blank" href="https://prometheus-operator.dev/docs/prologue/quick-start/">Quick Start</a> page.
Prometheus can then be accessed as documented in the
<a id="" title="" target="_blank" href="https://prometheus-operator.dev/docs/prologue/quick-start/#access-prometheus">Access Prometheus section of the Quick Start</a> page.</p>

<p>As described in the Prometheus docs we can create a port-forward process to the Prometheus <code>Service</code>.</p>

<markup
lang="bash"

>kubectl --namespace monitoring port-forward svc/prometheus-k8s 9090</markup>

<p>Then point our browser to <a id="" title="" target="_blank" href="http://localhost:9090">http://localhost:9090</a> to access the Prometheus UI.</p>



<v-card>
<v-card-text class="overflow-y-hidden" >
<img src="./images/img/prom.png" alt="Prometheus UI" />
</v-card-text>
</v-card>

<p>At this stage there will be no Coherence metrics, but we&#8217;ll change that in the next section.</p>

</div>

<h3 id="_create_a_servicemonitor">Create a ServiceMonitor</h3>
<div class="section">
<p>The out of the box Prometheus install uses <code>ServiceMonitor</code> resources to determine which Pods to scrape metrics from.
We can therefore configure Prometheus to scrape our Coherence cluster metrics by adding a <code>Service</code> and <code>ServiceMonitor</code>.</p>

<p>A Prometheus <code>ServiceMonitor</code>, as the name suggests, monitors a <code>Service</code> so we need to create a <code>Service</code> to expose the metrics port.
We are not going to access this <code>Service</code> ourselves, so it does not need to be a load balancer, in fact it can just be a headless service.
Prometheus uses the <code>Service</code> to locate the Pods that it should scrape.</p>

<p>The yaml below is a simple headless service that has a selector that matches labels in our Coherence cluster <code>Pods</code>.</p>

<markup
lang="yaml"
title="prometheus-metrics.yaml"
>apiVersion: v1
kind: Service
metadata:
  name: storage-metrics
  labels:
    coherence.oracle.com/cluster: test-cluster
    coherence.oracle.com/deployment: storage
    coherence.oracle.com/component: metrics-service
spec:
  type: ClusterIP
  ports:
  - name: metrics
    port: 9612
    targetPort: metrics
  selector:
    coherence.oracle.com/cluster: test-cluster
    coherence.oracle.com/deployment: storage</markup>

<p>We can now create a Prometheus <code>ServiceMonitor</code> that tells Prometheus about the <code>Service</code> to use.</p>

<markup
lang="yaml"
title="prometheus-metrics.yaml"
>apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: storage-metrics
  labels:
    coherence.oracle.com/cluster: test-cluster
    coherence.oracle.com/deployment: storage
    coherence.oracle.com/component: service-monitor
spec:
  endpoints:
  - port: metrics
  selector:
    matchLabels:
        coherence.oracle.com/cluster: test-cluster
        coherence.oracle.com/deployment: storage
        coherence.oracle.com/component: metrics-service</markup>

<p>The <code>ServiceMonitor</code> above contains a single endpoint that scrapes the port named <code>metrics</code> in any <code>Service</code> with labels matching those in the <code>matchLabels</code> array, which in this case are the labels we applied to the <code>storage-metrics</code> service above.</p>

<p>The full specification of what can be in a <code>ServiceMonitor</code> can be found in the Prometheus
<a id="" title="" target="_blank" href="https://github.com/prometheus-operator/prometheus-operator/blob/master/Documentation/api.md#servicemonitorspec">ServiceMonitorSpec</a>
documentation.</p>

<p>We can combine both of the above pieces of yaml into a single file and deploy them.
The example source code contains a file named <code>prometheus-metrics.yaml</code> that contains the yaml above.
Create the <code>Service</code> and <code>ServiceMonitor</code> in the same Kubernetes namespace as the Coherence cluster.</p>

<markup
lang="bash"

>kubectl apply -f prometheus-metrics.yaml</markup>

<p>It can sometimes take a minute or two for Prometheus to discover the <code>ServiceMonitor</code> and start to scrape metrics from the Pods. Once this happens it should be possible to see Coherence metrics for the cluster in Prometheus.</p>



<v-card>
<v-card-text class="overflow-y-hidden" >
<img src="./images/img/prom-coh.png" alt="Prometheus UI" />
</v-card-text>
</v-card>

<p>As shown above, the <code>vendor:coherence_cluster_size</code> metric has been scraped from all three <code>Pods</code> and as expected all <code>Pods</code> have a cluster size value of <code>3</code>.</p>

</div>

<h3 id="_federated_prometheus_metrics">Federated Prometheus Metrics</h3>
<div class="section">
<p>Prometheus Federation is the recommended way to scale Prometheus and to make metrics from inside Kubernetes available in a Prometheus instance outside of Kubernetes. Instead of the external Prometheus instance needing to be configured to locate and connect to <code>Pods</code> inside Kubernetes, it only needs an ingress into Prometheus running inside Kubernetes and can scrape all the metrics from there.
More details can be found in the <a id="" title="" target="_blank" href="https://prometheus.io/docs/prometheus/latest/federation/">Prometheus Federation</a> documentation.</p>

<p>We can install a local Prometheus instance as described in the <a id="" title="" target="_blank" href="https://prometheus.io/docs/prometheus/latest/getting_started/">Prometheus Getting Started</a> guide.</p>

<p>In the Prometheus installation directory we can edit the <code>prometheus.yml</code> file to configure Prometheus to scrape the federation endpoint of Prometheus inside Kubernetes. We need to <strong>add</strong> the federation configuration to the <code>scrape_configs:</code> section as shown below:</p>

<markup
lang="yaml"
title="prometheus.yml"
>scrape_configs:
  - job_name: 'federate'
    scrape_interval: 15s
    honor_labels: true
    metrics_path: '/federate'
    params:
      'match[]':
        - '{__name__=~"vendor:coherence_.*"}'
    static_configs:
      - targets:
        - '127.0.0.1:9091'</markup>

<p>You will notice that we have used <code>127.0.0.1:9091</code> as the target address. This is because when we run our local Prometheus instance it will bind to port 9090 so when we run the port-forward process to allow connections into Prometheus in the cluster we cannot use port <code>9090</code>, so we will forward local port <code>9091</code> to the Prometheus service port <code>9090</code> in Kubernetes.</p>

<p>In the <code>params:</code> section we have specified that the <code>'match[]':</code> field only federates metrics that have a name that starts with <code>vendor:coherence_</code> so in this example we only federate Coherence metrics.</p>

<p>Run the port-forward process so that when we start our local Prometheus instance it can connect to Prometheus in Kubernetes.</p>

<markup
lang="bash"

>kubectl --namespace monitoring port-forward svc/prometheus-k8s 9091:9090</markup>

<p>We&#8217;re now forwarding local port 9091 to Prometheus service port 9090 so we can run the local Prometheus instance.
As described in the Prometheus documentation, from the Prometheus installation directory run the command:</p>

<markup
lang="bash"

>./prometheus --config.file=prometheus.yml</markup>

<p>Once Prometheus starts we can point our browser to <a id="" title="" target="_blank" href="http://localhost:9090">http://localhost:9090</a> to access the prometheus UI.
After a short pause, Prometheus should start to scrap emetrics from inside Kubernetes and we should see them in the UI</p>



<v-card>
<v-card-text class="overflow-y-hidden" >
<img src="./images/img/prom-federate.png" alt="Prometheus UI" />
</v-card-text>
</v-card>

</div>
</div>

<h2 id="_grafana">Grafana</h2>
<div class="section">
<p>We could now install Grafana and configure it to connect to Prometheus, either the local instance or the instance inside Kubernetes. The Coherence Operator provides a number of dashboards that can imported into Grafana. See the Operator
<a id="" title="" target="_blank" href="https://oracle.github.io/coherence-operator/docs/latest/#/docs/metrics/030_importing">Import Grafana Dashboards</a> documentation.</p>

</div>
</doc-view>
