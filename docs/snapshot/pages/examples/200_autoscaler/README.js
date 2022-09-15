<doc-view>

<h2 id="_kubernetes_horizontal_pod_autoscaler_example">Kubernetes Horizontal Pod autoscaler Example</h2>
<div class="section">
<p>This example shows how to use the Kubernetes Horizontal Pod Autoscaler to scale Coherence clusters.</p>

<div class="admonition tip">
<p class="admonition-textlabel">Tip</p>
<p ><p><img src="./images/GitHub-Mark-32px.png" alt="GitHub Mark 32px" />
 The complete source code for this example is in the <a id="" title="" target="_blank" href="https://github.com/oracle/coherence-operator/tree/master/examples/200_autoscaler">Coherence Operator GitHub</a> repository.</p>
</p>
</div>
</div>

<h2 id="_how_does_the_horizontal_pod_autoscaler_work">How Does the Horizontal Pod autoscaler Work</h2>
<div class="section">
<p>There is a lot of good documentation on the HPA, particularly the <a id="" title="" target="_blank" href="https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/">Kubernetes documentation</a>.</p>

<p>The HPA uses metrics, which it obtains from one of the Kubernetes metrics APIs.
Many cloud providers and custom Kubernetes installations have metrics features that may be able to expose those metrics to
the <code>custom/metrics.k8s.io</code> API.
It is possible to even do everything yourself and build a custom REST endpoint that serves custom metrics to the HPA.
Those alternatives are beyond the scope of this example though so to keep things simple we will use Prometheus.
The diagram below shows, at a high level, how this works.</p>



<v-card>
<v-card-text class="overflow-y-hidden" >
<img src="./images/images/autoscaler.png" alt="autoscaler" />
</v-card-text>
</v-card>

<p>Prometheus will obtain metrics from the Coherence Pod&#8217;s metrics endpoints.
The Prometheus Adapter exposes certain configured metrics polled from Prometheus as custom Kubernetes metrics.
The HPA is configured to poll the custom metrics and use those to scale the <code>Coherence</code> resource (which will in turn cause
the Coherence Operator to scale the <code>StatefulSet</code>).</p>

</div>

<h2 id="_autoscaling_coherence_clusters">Autoscaling Coherence Clusters</h2>
<div class="section">
<p>This example will show how to configure the HPA to scale Coherence clusters based on heap usage metrics.
As Coherence stores data in memory, monitoring heap usage and using it to scale seems a sensible approach.</p>

<p>The <code>Coherence</code> CRD supports the <code>scale</code> sub-resource, which means that the Kubernetes HPA can be
used to scale a <code>Coherence</code> deployment.
In this example we are going to use heap usage as the metric - or to be more specific the amount of heap in use <em>after</em> the
last garbage collection.
This is an important point, plain heap usage is a poor metric to use for scaling decisions because the heap may be very
full at a given point in time, but most of that memory may be garbage so scaling on the plain heap usage figure may cause the
cluster to scale up needlessly as a milli-second later a GC could run, and the heap use shrinks down to acceptable levels.</p>

<p>The problem is that there is no single metric in a JVM that gives heap usage after garbage collection.
Coherence has some metrics that report this value, but they are taken from the <code>MemoryPool</code> MBeans and this is not reliable
for scaling.
For example, if the JVM is using the G1 collector the <code>G1 Old Gen</code> memory pool value for heap use after garbage collection
will be zero unless a full GC has run.
It is quite possible to almost fill the heap without running a full GC so this figure could remain zero or be wildly inaccurate.</p>

<p>A more reliable way to work out the heap usage is to obtain the values for the different heap memory pools from the
Garbage Collector MBeans. There could be multiple of these MBeans with different names depending on which collector
has been configured for the JVM.
The Garbage Collector Mbeans have a <code>LastGCcInfo</code> attribute, which is a composite attribute containing information about the last
garbage collection that ran on this collector. One of the attributes is the <code>endTime</code>, which we can use to determine which
collector&#8217;s <code>LastGCcInfo</code> is the most recent. Once we have this we can obtain the <code>memoryUsageAfterGc</code> attribute for the last gc,
which is a map of memory pool name to heap use data after the GC.
We can use this to then sum up the usages for the different heap memory pools.</p>

<p>The Java code in this example contains a simple MBean class <code>HeapUsage</code> and corresponding MBean interface <code>HeapUsageMBean</code>
that obtain heap use metrics in the way detailed above. There is also a configuration file <code>custom-mbeans.xml</code> that
Coherence will use to automatically add the custom MBean to Coherence management and metrics.
There is Coherence documentation on
<a id="" title="" target="_blank" href="https://docs.oracle.com/en/middleware/standalone/coherence/14.1.1.0/manage/using-coherence-metrics.html#GUID-CFC31D23-06B8-49AF-8996-ADBA806E0DD9">how to add custom metrics</a>
and
<a id="" title="" target="_blank" href="https://docs.oracle.com/en/middleware/standalone/coherence/14.1.1.0/manage/registering-custom-mbeans.html#GUID-1EE749C5-BC0D-4353-B5FE-1C5DCDEAE48C">how to register custom MBeans</a>.</p>

<p>The custom heap use MBean will be added with an ObjectName of <code>Coherence:type=HeapUsage,nodeId=1</code> where <code>nodeId</code> will change to
match the Coherence member id for the specific JVM. There will be one heap usage MBean for each cluster member.</p>

<p>The Coherence metrics framework will expose the custom metrics with metric names made up from the MBean domain name,
type, and the attribute name. The MBean has attribute names <code>Used</code> and <code>PercentageUsed</code>, so the metric names will be:</p>

<ul class="ulist">
<li>
<p><code>Coherence.HeapUsage.Used</code></p>

</li>
<li>
<p><code>Coherence.HeapUsage.PercentageUsed</code></p>

</li>
</ul>
<p>These metrics will be scoped as application metrics, as opposed to Coherence standard metrics that are vendor scoped.
This means that in Prometheus the names will be converted to:</p>

<ul class="ulist">
<li>
<p><code>application:coherence_heap_usage_used</code></p>

</li>
<li>
<p><code>application:coherence_heap_usage_percentage_used</code></p>

</li>
</ul>
<p>The metrics will have corresponding tags to identify which cluster member (<code>Pod</code>) they relate to.</p>

</div>

<h2 id="_building_the_example">Building the Example</h2>
<div class="section">

<h3 id="_clone_the_coherence_operator_repository">Clone the Coherence Operator Repository:</h3>
<div class="section">
<p>To build the examples, you first need to clone the Operator GitHub repository to your development machine.</p>

<markup
lang="bash"

>git clone https://github.com/oracle/coherence-operator

cd coherence-operator/examples</markup>

</div>

<h3 id="_build_the_examples">Build the Examples</h3>
<div class="section">

<h4 id="_prerequisites">Prerequisites</h4>
<div class="section">
<ul class="ulist">
<li>
<p>Java 11+ JDK either [OpenJDK](<a id="" title="" target="_blank" href="https://adoptopenjdk.net/">https://adoptopenjdk.net/</a>) or [Oracle JDK](<a id="" title="" target="_blank" href="https://www.oracle.com/java/technologies/javase-downloads.html">https://www.oracle.com/java/technologies/javase-downloads.html</a>)</p>

</li>
<li>
<p>[Docker](<a id="" title="" target="_blank" href="https://docs.docker.com/install/">https://docs.docker.com/install/</a>) version 17.03+.</p>

</li>
<li>
<p>[kubectl](<a id="" title="" target="_blank" href="https://kubernetes.io/docs/tasks/tools/install-kubectl/">https://kubernetes.io/docs/tasks/tools/install-kubectl/</a>) version v1.13.0+ .</p>

</li>
<li>
<p>Access to a Kubernetes v1.14.0+ cluster.</p>

</li>
<li>
<p>[Helm](<a id="" title="" target="_blank" href="https://helm.sh/docs/intro/install/">https://helm.sh/docs/intro/install/</a>) version 3.2.4+</p>

</li>
</ul>
<p>Building the project requires [Maven](<a id="" title="" target="_blank" href="https://maven.apache.org">https://maven.apache.org</a>) version 3.6.0+.
The commands below use the Maven Wrapper to run the commands, which will install Maven if it is not
already on the development machine. If you already have a suitable version of Maven installed feel free to replace
the use of <code>./mvnw</code> in the examples with your normal Maven command (typically just <code>mvn</code>).</p>


<h5 id="_corporate_proxies">Corporate Proxies</h5>
<div class="section">
<p>If building inside a corporate proxy (or any machine that requires http and https proxies to be configured) then
the build will require the <code>MAVEN_OPTS</code> environment variable to be properly set, for example:</p>

<markup
lang="bash"

>export MAVEN_OPTS="-Dhttps.proxyHost=host -Dhttps.proxyPort=80 -Dhttp.proxyHost=host -Dhttp.proxyPort=80"</markup>

<p>replacing <code>host</code> with the required proxy hostname and <code>80</code> with the proxy&#8217;s port.</p>

</div>
</div>

<h4 id="_build_instructions">Build Instructions</h4>
<div class="section">
<p>The autoscaler example uses the <a id="" title="" target="_blank" href="https://github.com/GoogleContainerTools/jib/tree/master/jib-maven-plugin#build-your-image">JIB Maven plugin</a> to build the example image.
To build the image run the following command from the <code>examples/autoscaler</code> directory:</p>

<markup
lang="bash"

>./mvnw package jib:dockerBuild</markup>

<p>The build will produce various example images, for the autoscaler example we will be using the <code>autoscaler-example:latest</code> image.</p>

</div>
</div>
</div>

<h2 id="_run_the_example">Run the Example</h2>
<div class="section">
<p>Running the example requires a number of components to be installed.
The example will use Prometheus as a custom metrics source, which requires installation of Prometheus and the
Prometheus Adapter custom metrics source.</p>

<div class="admonition note">
<p class="admonition-inline">To simplify the example commands none of the examples below use a Kubernetes namespace.
If you wish to install the components below into a namespace other than <code>default</code>, then use the required
kubectl and Helm namespace options.</p>
</div>

<h3 id="_install_the_coherence_operator">Install the Coherence Operator</h3>
<div class="section">
<p>First install the Coherence Operator, TBD&#8230;&#8203;</p>

</div>

<h3 id="_install_coherence_cluster">Install Coherence cluster</h3>
<div class="section">
<p>With the Coherence Operator running we can now install a simple Coherence cluster.
An example of the yaml required is below:</p>

<markup
lang="yaml"
title="cluster.yaml"
>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: test-cluster
spec:
  image: autoscaler-example:latest  <span class="conum" data-value="1" />
  imagePullPolicy: IfNotPresent
  replicas: 2                       <span class="conum" data-value="2" />
  coherence:
    metrics:
      enabled: true                 <span class="conum" data-value="3" />
  jvm:
    memory:
      heapSize: 500m                <span class="conum" data-value="4" />
  ports:
    - name: metrics                 <span class="conum" data-value="5" />
      serviceMonitor:
        enabled: true               <span class="conum" data-value="6" />
    - name: extend                  <span class="conum" data-value="7" />
      port: 20000</markup>

<ul class="colist">
<li data-value="1">The image used for the application will be the <code>autoscaler-example:latest</code> image we built above.</li>
<li data-value="2">The deployment will initially have 2 replicas.</li>
<li data-value="3">Coherence metrics must be enabled to publish the metrics we require for scaling.</li>
<li data-value="4">In this example the JVM heap has been fixed to <code>500m</code>, which is quite small but this means we do not need to add a lot of data
to cause excessive heap usage when we run the example.</li>
<li data-value="5">The metrics port must also be exposed on a <code>Service</code>.</li>
<li data-value="6">A Prometheus <code>ServiceMonitor</code> must also be enabled for the metrics service so that Prometheus can find the Coherence <code>Pods</code>
and poll metrics from them.</li>
<li data-value="7">This example also exposes a Coherence Extend port so that test data can easily be loaded into the caches.</li>
</ul>
<p>The autoscaler example includes a suitable yaml file named <code>cluster.yaml</code> in the <code>manifests/</code> directory that can be used
to create a Coherence deployment.</p>

<markup
lang="bash"

>kubectl create -f manifests/cluster.yaml</markup>

<p>The <code>Pods</code> that are part of the Coherence cluster can be listed with <code>kubectl</code>.
All the <code>Pods</code> have a label <code>coherenceCluster</code> set by the Coherence Operator to match the name of the
<code>Coherence</code> resource that they belong to, which makes it easier to list <code>Pods</code> for a specific deployment
using <code>kubectl</code>:</p>

<markup
lang="bash"

>kubectl get pod -l coherenceCluster=test-cluster</markup>

<p>In a short time the <code>Pods</code> should both be ready.</p>

<markup
lang="bash"

>NAME             READY   STATUS    RESTARTS   AGE
test-cluster-0   1/1     Running   0          2m52s
test-cluster-1   1/1     Running   0          2m52s</markup>


<h4 id="_test_the_custom_heap_metrics">Test the Custom Heap Metrics</h4>
<div class="section">
<p>The Metrics endpoint will be exposed on port 9612 on each <code>Pod</code>, so it is possible to query the metrics endpoints
for the custom heap metrics. The simplest way to test the metrics is to use the <code>kubectl</code> <code>port-forward</code> command and <code>curl</code>.</p>

<p>In one terminal session start the port forwarder to the first <code>Pod</code>, <code>test-cluster-0</code>:</p>

<markup
lang="bash"

>kubectl port-forward pod/test-cluster-0 9612:9612</markup>

<p>metrics from <code>Pod</code>, <code>test-cluster-0</code> can be queried on <code><a id="" title="" target="_blank" href="http://127.0.0.1:9612/metrics">http://127.0.0.1:9612/metrics</a></code></p>

<p>In a second terminal we can use curl to query the metrics.
The Coherence metrics endpoint serves metrics in two formats, plain text compatible with Prometheus and JSON.
If the required content type has not been specified in the curl command it could be either that is returned.
To specify a content type set the accepted type in the header, for example <code>--header "Accept: text/plain"</code> or
<code>--header "Accept: application/json"</code>.</p>

<p>This command will retrieve metrics from <code>test-cluster-0</code> in the same format that Prometheus would.</p>

<markup
lang="bash"

>curl -s --header "Accept: text/plain" -X GET http://127.0.0.1:9612/metrics</markup>

<p>This will return quite a lot of metrics, somewhere in that output is the custom application metrics for heap usage.
The simplest way to isolate them would be to use <code>grep</code>, for example:</p>

<markup
lang="bash"

>curl -s --header "Accept: text/plain" -X GET http://127.0.0.1:9612/metrics | grep application</markup>

<p>which should show something like:</p>

<markup
lang="bash"

>application:coherence_heap_usage_percentage_used{cluster="test-cluster", machine="docker-desktop", member="test-cluster-0", node_id="2", role="test-cluster", site="test-cluster-sts.operator-test.svc.cluster.local"} 3.09
application:coherence_heap_usage_used{cluster="test-cluster", machine="docker-desktop", member="test-cluster-0", node_id="2", role="test-cluster", site="test-cluster-sts.operator-test.svc.cluster.local"} 16177976</markup>

<p>The first metric <code>application:coherence_heap_usage_percentage_used</code> shows the heap was <code>3.09%</code> full after the last gc.
The second metric <code>application:coherence_heap_usage_used</code> shows that the in-use heap after the last gc was 16177976 bytes,
or around 16 MB.</p>

<p>The port forwarder can be changed to connect to the second <code>Pod</code> <code>test-cluster-1</code>, and the same curl command
will retrieve metrics from the second <code>Pod</code>, which should show different heap use values.</p>

</div>
</div>

<h3 id="_install_prometheus">Install Prometheus</h3>
<div class="section">
<p>The simplest way to install Prometheus as part of an example or demo is to use the
<a id="" title="" target="_blank" href="https://github.com/prometheus-operator/prometheus-operator">Prometheus Operator</a>, which can be
installed using a Helm chart.</p>


<h4 id="_setup_the_helm_repo">Setup the Helm Repo</h4>
<div class="section">
<p>Make sure the <code>stable</code> helm repository has been added to Helm if it isn&#8217;t already present in your local Helm repositories.</p>

<markup
lang="bash"

>helm repo add stable https://kubernetes-charts.storage.googleapis.com/</markup>

<p>Make sure the local Helm repository is up to date.</p>

<markup
lang="bash"

>helm repo update</markup>

</div>

<h4 id="_configure_prometheus_rbac">Configure Prometheus RBAC</h4>
<div class="section">
<p>If you are using a Kubernetes cluster with RBAC enabled then the rules required by Prometheus need to be added.
The autoscale example contains a yaml file with the required RBAC rules in it in the <code>manifests/</code> directory.</p>

<p>The <code>manifests/prometheus-rbac.yaml</code> uses a namespace <code>coherence-example</code> which may need to be changed
if you are installing into a different namespace.</p>

<p>The following commands use <code>sed</code> to replace <code>coherence-example</code> with <code>default</code> and pipe the result to <code>kubectl</code>
to create the RBAC rules in the <code>default</code> Kubernetes namespace.</p>

<markup
lang="bash"

>sed "s/coherence-example/default/g"  manifests/prometheus-rbac.yaml | kubectl create -f -</markup>

</div>

<h4 id="_install_the_prometheus_operator">Install the Prometheus Operator</h4>
<div class="section">
<p>The Prometheus Operator can now be installed using Helm. The autoscaler example contains a simple values files
that can be used when installing the chart in the <code>manifests/</code> directory.</p>

<markup
lang="bash"

>helm install --atomic --version 8.13.9 --wait \
    --set prometheus.service.type=NodePort \
    --values manifests/prometheus-values.yaml prometheus stable/prometheus-operator</markup>

<p>The <code>--wait</code> parameter makes Helm block until all the installed resources are ready.</p>

<p>The command above sets the <code>prometheus.service.type</code> value to <code>NodePort</code> so that the Prometheus UI will be exposed
on a port on the Kubernetes node. This is particularly useful when testing with a local Kubernetes cluster, such as in Docker
on a laptop because the UI can be reached on localhost at that port. The default node port is <code>30090</code>, this can be
changed by setting a different port, e.g: <code>--set prometheus.service.nodePort=9090</code>.</p>

<p>Assuming the default port of <code>30090</code> is used the UI can be reached on <a id="" title="" target="_blank" href="http://127.0.0.1:30090">http://127.0.0.1:30090</a>.</p>



<v-card>
<v-card-text class="overflow-y-hidden" >
<img src="./images/images/prometheus-ui-empty.png" alt="prometheus ui empty" />
</v-card-text>
</v-card>

<p>After Prometheus has started up and is scraping metrics we should be able to see our custom metrics in the UI.
Type the metric name <code>application:coherence_heap_usage_percentage_used</code> in the expression box and click <code>Execute</code>
and Prometheus should show two values for the metric, one for each <code>Pod</code>.</p>



<v-card>
<v-card-text class="overflow-y-hidden" >
<img src="./images/images/prometheus-ui-metrics.png" alt="prometheus ui metrics" />
</v-card-text>
</v-card>

<p>Prometheus is scraping many more Coherence metrics that can also be queried in the UI.</p>

</div>
</div>

<h3 id="_install_prometheus_adapter">Install Prometheus Adapter</h3>
<div class="section">
<p>The next step in the example is to install the Prometheus Adapter. This is a custom metrics server that published metrics
using the Kubernetes <code>custom/metrics.k8s.io</code> API. This is required because the HPA cannot query metrics directly from
Prometheus, only from standard Kubernetes metrics APIs.
As with Prometheus the simplest way to install the adapter is by using the Helm chart.
Before installing though we need to create the adapter configuration so that it can publish our custom metrics.</p>

<p>The documentation for the adapter configuration is not the simplest to understand quickly.
On top of that the adapter documentation shows how to configure the adapter using a <code>ConfigMap</code> whereas the Helm chart
adds the configuration to the Helm values file.</p>

<p>The basic format for configuring a metric in the adapter is as follows:</p>

<markup
lang="yaml"

>- seriesQuery: 'application:coherence_heap_usage_percentage_used'   <span class="conum" data-value="1" />
  resources:
    overrides:   <span class="conum" data-value="2" />
      namespace: <span class="conum" data-value="3" />
        resource: "namespace"
      pod:   <span class="conum" data-value="4" />
        resource: "pod"
      role:  <span class="conum" data-value="5" />
        group: "coherence.oracle.com"
        resource: "coherence"
  name:
    matches: ""
    as: "heap_memory_usage_after_gc_pct"  <span class="conum" data-value="6" />
  metricsQuery: sum(&lt;&lt;.Series&gt;&gt;{&lt;&lt;.LabelMatchers&gt;&gt;}) by (&lt;&lt;.GroupBy&gt;&gt;)  <span class="conum" data-value="7" /></markup>

<ul class="colist">
<li data-value="1">The <code>seriesQuery</code> is the name of the metric to be retrieved from Prometheus.
This is the same name used when querying in the UI.
The name can be qualified further with tags/labels but in our case just the metric name is sufficient.</li>
<li data-value="2">The <code>overrides</code> section matches metric labels to Kubernetes resources, which can be used in queries (more about this below).</li>
<li data-value="3">The metrics have a <code>namespace</code> label (as can be seen in the UI above) and this maps to a Kubernetes <code>Namespace</code> resource.</li>
<li data-value="4">The metrics have a <code>pod</code> label (as can be seen in the UI above) and this maps to a Kubernetes <code>Pod</code> resource.</li>
<li data-value="5">The metrics have a <code>role</code> label (as can be seen in the UI above) and this maps to a Kubernetes
<code>coherence.coherence.oracle.com</code> resource.</li>
<li data-value="6">The <code>name.as</code> field gives the name of the metric in the metrics API.</li>
<li data-value="7">The <code>metricsQuery</code> determines how a specific metric will be fetched, in this case we are summing the values.</li>
</ul>
<p>The configuration above will create a metric in the <code>custom/metrics.k8s.io</code> API named heap_memory_usage_after_gc_pct.
This metric can be retrieved from the API for a namespace, for a Pod or for a Coherence deployment
(the <code>coherence.coherence.oracle.com</code> resource). This is why the <code>metricsQuery</code> uses <code>sum</code>, so that when querying for
a metric at the namespace level we see the total summed up for the namespace.</p>

<p>Summing up the metric might not be the best approach. Imagine that we want to scale when the heap after gc usage exceeds 80%.
Ideally this is when any JVM heap in use after garbage collection exceeds 80%.
Whilst Coherence will distribute data evenly across the cluster so that each member holds a similar amount of data and has
similar heap usage, there could be an occasion where one member for whatever reason is processing extra load and exceeds 80%
before other members.</p>

<p>One way to approach this issue is instead of summing the metric value for a namespace or <code>coherence.coherence.oracle.com</code>
resource we can fetch the maximum value. We do this by changing the <code>metricsQuery</code> to use <code>max</code> as shown below:</p>

<markup
lang="yaml"

>- seriesQuery: 'application:coherence_heap_usage_percentage_used'
  resources:
    overrides:
      namespace:
        resource: "namespace"
      pod:
        resource: "pod"
      role:
        group: "coherence.oracle.com"
        resource: "coherence"
  name:
    matches: ""
    as: "heap_memory_usage_after_gc_max_pct"
  metricsQuery: max(&lt;&lt;.Series&gt;&gt;{&lt;&lt;.LabelMatchers&gt;&gt;}) by (&lt;&lt;.GroupBy&gt;&gt;)</markup>

<p>This is the same configuration as previously but now the <code>metricsQuery</code> uses the <code>max</code> function, and the
metric name has been changed to <code>heap_memory_usage_after_gc_max_pct</code> so that it is obvious it is a maximum value.</p>

<p>We can repeat the configuration above for the <code>application:coherence_heap_usage_used</code> metric too so that we will end up with
four metrics in the <code>custom/metrics.k8s.io</code> API:</p>

<ul class="ulist">
<li>
<p><code>heap_memory_usage_after_gc_max_pct</code></p>

</li>
<li>
<p><code>heap_memory_usage_after_gc_pct</code></p>

</li>
<li>
<p><code>heap_memory_usage_after_gc</code></p>

</li>
<li>
<p><code>heap_memory_usage_after_gc_max</code></p>

</li>
</ul>
<p>The autoscaler example has a Prometheus Adapter Helm chart values file that contains the configuration for the
four metrics. This can be used to install the adapter
<a id="" title="" target="_blank" href="https://hub.helm.sh/charts/prometheus-com/prometheus-adapter">Helm chart</a>:</p>

<div class="admonition note">
<p class="admonition-inline">In the command below the <code>--set prometheus.url=http://prometheus-prometheus-oper-prometheus.default.svc</code>
parameter tells the adapter how to connect to Prometheus.
The Prometheus Operator creates a <code>Service</code> named <code>prometheus-prometheus-oper-prometheus</code> to expose Prometheus.
In this case the command assumes Prometheus is installed in the <code>default</code> namespace.
If you installed Prometheus into a different namespace change the <code>default</code> part of
<code>prometheus-prometheus-oper-prometheus.<strong>default</strong>.svc</code> to the actual namespace name.</p>
</div>
<div class="admonition note">
<p class="admonition-inline">The <code>manifests/prometheus-adapter-values.yaml</code> contains the configurations for metrics that the adapter
will publish. These work with Coherence Operator 3.1.0 and above. If using an earlier 3.0.x version the values
file must first be edited to change all occurrences of <code>resource: "coherence"</code> to <code>resource: "coherence"</code> (to
make the resource name singular).</p>
</div>
<markup
lang="bash"

>helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update

helm install --atomic --wait \
    --set prometheus.url=http://prometheus-prometheus-oper-prometheus.default.svc \
    --values manifests/prometheus-adapter-values.yaml \
    prometheus-adapter prometheus-community/prometheus-adapter</markup>


<h4 id="_query_custom_metrics">Query Custom Metrics</h4>
<div class="section">
<p>Now the Prometheus adapter is running we can query metrics from the <code>custom/metrics.k8s.io</code> API using <code>kubectl</code> raw API access.
This is the same API that the HPA will use to obtain metrics.</p>

<p>If a Coherence cluster had been installed into the <code>default</code> namespace, then metrics could be fetched for all <code>Pods</code> in
that specific namespace, for example to obtain the <code>heap_memory_usage_after_gc_pct</code> metric:</p>

<markup
lang="bash"

>kubectl get --raw /apis/custom.metrics.k8s.io/v1beta1/namespaces/default/pods/*/heap_memory_usage_after_gc_pct</markup>

<p>The <code>*</code> after <code>pods/</code> tells the adapter to fetch metrics for all <code>Pods</code> in the namespace.
To fetch the metric for pods in another namespace change the <code>default</code> part of the URL to the namespace name.</p>

<p>If you have the <code>jq</code> utility installed that formats json then piping the output to <code>jq</code> will make it prettier.</p>

<markup
lang="bash"

>kubectl get --raw /apis/custom.metrics.k8s.io/v1beta1/namespaces/default/pods/*/heap_memory_usage_after_gc_pct | jq</markup>

<p>We could fetch a metric for a specific <code>Pod</code> in the <code>default</code> namespace, for example a <code>Pod</code> named <code>test-cluster-1</code> as follows:</p>

<markup
lang="bash"

>kubectl get --raw /apis/custom.metrics.k8s.io/v1beta1/namespaces/default/pods/test-cluster-1/heap_memory_usage_after_gc_pct</markup>

<p>which might display something like:</p>

<markup
lang="json"

>{
  "kind": "MetricValueList",
  "apiVersion": "custom.metrics.k8s.io/v1beta1",
  "metadata": {
    "selfLink": "/apis/custom.metrics.k8s.io/v1beta1/namespaces/coherence-test/pods/test-cluster-1/heap_memory_usage_after_gc_pct"
  },
  "items": [
    {
      "describedObject": {
        "kind": "Pod",
        "namespace": "operator-test",
        "name": "test-cluster-1",
        "apiVersion": "/v1"
      },
      "metricName": "heap_memory_usage_after_gc_pct",
      "timestamp": "2020-09-02T12:12:01Z",
      "value": "1300m",
      "selector": null
    }
  ]
}</markup>

<div class="admonition note">
<p class="admonition-inline">The format of the <code>value</code> field above might look a little strange. This is because it is a Kubernetes <code>Quantity</code>
format, in this case it is <code>1300m</code> where the <code>m</code> stand for millis. So in this case 1300 millis is 1.3% heap usage.
This is to get around the poor support in yaml and json for accurate floating-point numbers.</p>
</div>
<p>In our case for auto-scaling we are interested in the maximum heap for a specific <code>Coherence</code> resource.
Remember in the Prometheus Adapter configuration we configured the <code>role</code> metric tag to map to
<code>coherence.coherence.oracle.com</code> resources.
We also configured a query that will give back the maximum heap usage value for a query.</p>

<p>The example yaml used to deploy the <code>Coherence</code> resource above will create a resource named <code>test-cluster</code>.
If we installed this into the <code>default</code> Kubernetes namespace then we can fetch the maximum heap use after gc
for the <code>Pods</code> in that <code>Coherence</code> deployment as follows:</p>

<markup
lang="bash"

>kubectl get --raw /apis/custom.metrics.k8s.io/v1beta1/namespaces/default/coherence.coherence.oracle.com/test-cluster/heap_memory_usage_after_gc_max_pct</markup>

<p>which might display something like:</p>

<markup
lang="json"

>{
  "kind": "MetricValueList",
  "apiVersion": "custom.metrics.k8s.io/v1beta1",
  "metadata": {
    "selfLink": "/apis/custom.metrics.k8s.io/v1beta1/namespaces/operator-test/coherence.coherence.oracle.com/test-cluster/heap_memory_usage_after_gc_max_pct"
  },
  "items": [
    {
      "describedObject": {
        "kind": "Coherence",
        "namespace": "operator-test",
        "name": "test-cluster",
        "apiVersion": "coherence.oracle.com/v1"
      },
      "metricName": "heap_memory_usage_after_gc_max_pct",
      "timestamp": "2020-09-02T12:21:02Z",
      "value": "3300m",
      "selector": null
    }
  ]
}</markup>

</div>
</div>

<h3 id="_configure_the_horizontal_pod_autoscaler">Configure The Horizontal Pod autoscaler</h3>
<div class="section">
<p>Now that we have custom metrics in the Kubernets <code>custom.metrics.k8s.io</code> API, the final piece is to add the HPA
configuration for the Coherence deployment that we want to scale.
To configure the HPA we need to create a <code>HorizontalPodautoscaler</code> resource for each Coherence deployment in the same namespace
as we deployed the Coherence deployment to.</p>

<p>Below is an example <code>HorizontalPodautoscaler</code> resource that will scale our example Coherence deployment:</p>

<markup
lang="yaml"
title="hpa.yaml"
>apiVersion: autoscaling/v2beta2
kind: HorizontalPodautoscaler
metadata:
  name: test-cluster-hpa
spec:
  scaleTargetRef:                         <span class="conum" data-value="1" />
    apiVersion: coherence.oracle.com/v1
    kind: Coherence
    name: test-cluster
  minReplicas: 2         <span class="conum" data-value="2" />
  maxReplicas: 5
  metrics:               <span class="conum" data-value="3" />
  - type: Object
    object:
      describedObject:
        apiVersion: coherence.oracle.com/v1
        kind: Coherence
        name: test-cluster
      metric:
        name: heap_memory_usage_after_gc_max_pct  <span class="conum" data-value="4" />
      target:
        type: Value       <span class="conum" data-value="5" />
        value: 80
  behavior:                             <span class="conum" data-value="6" />
    scaleUp:
      stabilizationWindowSeconds: 120
    scaleDown:
      stabilizationWindowSeconds: 120</markup>

<ul class="colist">
<li data-value="1">The <code>scaleTargetRef</code> points to the resource that the HPA will scale. In this case it is our <code>Coherence</code> deployment
which is named <code>test-cluster</code>. The <code>apiVersion</code> and <code>kind</code> fields match those in the <code>Coherence</code> resource.</li>
<li data-value="2">For this example, the Coherence deployment will have a minimum of 2 replicas and a maximum of 5, so the HPA will not scale up too much.</li>
<li data-value="3">The <code>metrics</code> section in the yaml above tells the HPA how to query our custom metric.
In this case we want to query the single max usage value metric for the <code>Coherence</code> deployment (like we did manually when using
kubectl above). To do this we add a metric with a <code>type</code> of <code>Object</code>.
The <code>describedObject</code> section describes the resource to query, in this case kind <code>Coherence</code> in resource group <code>coherence.oracle.com</code> with the name <code>test-cluster</code>.</li>
<li data-value="4">The metric name to query is our custom max heap usage percentage metric <code>heap_memory_usage_after_gc_max_pct</code>.</li>
<li data-value="5">The <code>target</code> section describes the target value for the metric, in this case 80 thousand millis - which is 80%.</li>
<li data-value="6">The <code>behavior</code> section sets a window of 120 seconds so that the HAP will wait at least 120 seconds after scaling up
or down before re-evaluating the metric. This gives Coherence enough time to scale the deployment and for the data to redistribute
and gc to occur. In real life this value would need to be adjusted to work correctly on your actual cluster.</li>
</ul>
<p>The autoscaler example contains yaml to create the <code>HorizontalPodautoscaler</code> resource in the <code>manifests/</code> directory.</p>

<markup
lang="bash"

>kubectl create -f manifests/hpa.yaml</markup>

<p>The <code>hpa.yaml</code> file will create a <code>HorizontalPodautoscaler</code> resource named <code>test-cluster-hpa</code>.
After waiting a minute or two for the HPA to get around to polling our new <code>HorizontalPodautoscaler</code> resource
we can check its status.</p>

<markup
lang="bash"

>kubectl describe horizontalpodautoscaler.autoscaling/test-cluster-hpa</markup>

<p>Which should show something like:</p>

<markup
lang="bash"

>Name:                                                                             test-cluster-hpa
Namespace:                                                                        operator-test
Labels:                                                                           &lt;none&gt;
Annotations:                                                                      &lt;none&gt;
CreationTimestamp:                                                                Wed, 02 Sep 2020 15:58:26 +0300
Reference:                                                                        Coherence/test-cluster
Metrics:                                                                          ( current / target )
  "heap_memory_usage_after_gc_max_pct" on Coherence/test-cluster (target value):  3300m / 80
Min replicas:                                                                     2
Max replicas:                                                                     10
Coherence pods:                                                                   2 current / 2 desired
Conditions:
  Type            Status  Reason               Message
  ----            ------  ------               -------
  AbleToScale     True    ScaleDownStabilized  recent recommendations were higher than current one, applying the highest recent recommendation
  ScalingActive   True    ValidMetricFound     the HPA was able to successfully calculate a replica count from Coherence metric heap_memory_usage_after_gc_max_pct
  ScalingLimited  False   DesiredWithinRange   the desired count is within the acceptable range
Events:           &lt;none&gt;</markup>

<p>We can see that the HPA has successfully polled the metric and obtained a value of <code>3300m</code> (so 3.3%) and has
decided that it does not need to scale.</p>

</div>

<h3 id="_add_data_scale_up">Add Data - Scale Up!</h3>
<div class="section">
<p>The HPA is now monitoring our Coherence deployment so we can now add data to the cluster and see the HPA scale up when
heap use grows.
The autoscaler example Maven pom file has been configured to use the Maven exec plugin to execute a Coherence command line
client that will connect over Coherence Extend to the demo cluster that we have deployed.</p>

<p>First we need to create a port forwarder to expose the Coherence Extend port locally.
Extend is bound to port 20000 in the <code>Pods</code> in our example.</p>

<markup
lang="bash"

>kubectl port-forward pod/test-cluster-0 20000:20000</markup>

<p>The command above forwards port 20000 in the <code>Pod</code> <code>test-cluster-0</code> to the local port 20000.</p>

<p>To start the client, run the following command in a terminal:</p>

<markup
lang="bash"

>./mvnw exec:java -pl autoscaler/</markup>

<p>The command above will start the console client and eventually display a <code>Map (?):</code> prompt.</p>

<p>At the map prompt, first create a cache named <code>test</code> with the <code>cache</code> command, type <code>cache test</code> and hit enter:</p>

<markup
lang="bash"

>Map (?): cache test</markup>

<p>There will now be a cache created in the cluster named <code>test</code>, and the map prompt will change to <code>Map (test):</code>.
We can add random data to this with the <code>bulkput</code> command. The format of the <code>bulkput</code> command is:</p>

<markup
lang="bash"

>bulkput &lt;# of iterations&gt; &lt;block size&gt; &lt;start key&gt; [&lt;batch size&gt; | all]</markup>

<p>So to add 20,000 entries of 10k bytes each starting at key <code>1</code> adding in batches of 1000 we can run
the <code>bulkput 20000 10000 1 1000</code> command at the map prompt:</p>

<markup
lang="bash"

>Map (test): bulkput 20000 10000 1 1000</markup>

<p>We can now look at the <code>HorizontalPodautoscaler</code> resource we create earlier with the command:</p>

<markup
lang="bash"

>kubectl get horizontalpodautoscaler.autoscaling/test-cluster-hpa</markup>

<p>Which will display something like:</p>

<markup
lang="bash"

>NAME               REFERENCE                TARGETS     MINPODS   MAXPODS   REPLICAS   AGE
test-cluster-hpa   Coherence/test-cluster   43700m/80   2         10        2          41m</markup>

<p>The HPA is now saying that the value of our heap use metric is 43.7%, so we can add a bit more data.
It may take a minute or two for the heap to increase and stabilise as different garbage collections happen across the Pods.
We should be able to safely add another 20000 entries putting the heap above 80% and hopefully scaling our deployment.</p>

<p>We need to change the third parameter to bulk put to 20000 otherwise the put will start again at key <code>1</code> and just overwrite the
previous entries, not really adding to the heap.</p>

<markup
lang="bash"

>Map (test): bulkput 20000 10000 20000 1000</markup>

<p>Now run the <code>kubectl describe</code> command on the <code>HorizontalPodautoscaler</code> resource again, and we should see that it has scaled
our cluster. If another 20,000 entries does not cause the heap to exceed 80% then you may need to run the <code>bulkput</code> command
once or twice more with a smaller number of entries to push the heap over 80%.</p>

<div class="admonition note">
<p class="admonition-inline">As previously mentioned, everything with HPA is slightly delayed due to the different components polling, and
stabilization times. It could take a few minutes for the HPA to actually scale the cluster.</p>
</div>
<markup
lang="bash"

>kubectl describe horizontalpodautoscaler.autoscaling/test-cluster-hpa</markup>

<p>The output of the <code>kubectl describe</code> command should now be something like this:</p>

<markup
lang="bash"

>Name:                                                                             test-cluster-hpa
Namespace:                                                                        operator-test
Labels:                                                                           &lt;none&gt;
Annotations:                                                                      &lt;none&gt;
CreationTimestamp:                                                                Wed, 02 Sep 2020 15:58:26 +0300
Reference:                                                                        Coherence/test-cluster
Metrics:                                                                          ( current / target )
  "heap_memory_usage_after_gc_max_pct" on Coherence/test-cluster (target value):  88300m / 80
Min replicas:                                                                     2
Max replicas:                                                                     10
Coherence pods:                                                                   2 current / 3 desired
Conditions:
  Type            Status  Reason              Message
  ----            ------  ------              -------
  AbleToScale     True    SucceededRescale    the HPA controller was able to update the target scale to 3
  ScalingActive   True    ValidMetricFound    the HPA was able to successfully calculate a replica count from Coherence metric heap_memory_usage_after_gc_max_pct
  ScalingLimited  False   DesiredWithinRange  the desired count is within the acceptable range
Events:
  Type    Reason             Age   From                       Message
  ----    ------             ----  ----                       -------
  Normal  SuccessfulRescale  1s    horizontal-pod-autoscaler  New size: 3; reason: Coherence metric heap_memory_usage_after_gc_max_pct above target</markup>

<p>We can see that the heap use value is now <code>88300m</code> or 88.3% and the events section shows that the HPA has scaled the <code>Coherence</code>
deployment to <code>3</code>. We can list the <code>Pods</code> and there should be three:</p>

<markup
lang="bash"

>kubectl get pod -l coherenceCluster=test-cluster</markup>

<markup
lang="bash"

>NAME             READY   STATUS    RESTARTS   AGE
test-cluster-0   1/1     Running   0          3h14m
test-cluster-1   1/1     Running   0          3h14m
test-cluster-2   1/1     Running   0          1m10s</markup>

<div class="admonition note">
<p class="admonition-inline">At this point Coherence will redistribute data to balance it over the three members of the cluster.
It may be that it takes considerable time for this to affect the heap usage as a lot of the cache data will be in the old generation of
the heap and not be immediately collected. This may then trigger another scale after the 120 second stabilization period that
we configured in the <code>HorizontalPodautoscaler</code>.</p>
</div>
</div>

<h3 id="_clean_up">Clean-Up</h3>
<div class="section">
<p>To clean-up after running the example just uninstall everything in the reverse order:</p>

<markup
lang="bash"

>kubectl delete -f manifests/hpa.yaml
helm delete prometheus-adapter
helm delete prometheus
kubectl delete -f manifests/cluster.yaml</markup>

<p>Remove the Prometheus RBAC rules, remembering to change the namespace name.</p>

<markup
lang="bash"

>sed "s/coherence-example/default/g"  manifests/prometheus-rbac.yaml | kubectl delete -f -</markup>

<p>Delete the Coherence deployment.</p>

<markup
lang="bash"

>kubectl delete manifests/cluster.yaml</markup>

<p>Undeploy the Operator.
TBD&#8230;&#8203;</p>

</div>
</div>

<h2 id="_conclusions">Conclusions</h2>
<div class="section">
<p>As we&#8217;ve shown, it is possible to use the HPA to scale a Coherence cluster based on metrics published by Coherence or
custom metrics, but there are some obvious caveats due to how HPA works.
There are inherent delays in the scaling process, the HPA only polls metrics periodically,
which themselves have been polled by Prometheus periodically and hence there can be some delay after
reaching a given heap size before the scale command actually reaches the Coherence Operator.
This will be obvious when running the example below.
Given a suitable configuration the HPA can be useful to scale as load increases but in no way can it
guarantee that an out of memory exception will never happen.</p>

<p>Using the HPA to scale as Coherence Pod&#8217;s heaps become filled is in no way an excuse not to do proper capacity planning
and size your Coherence clusters appropriately.</p>

</div>
</doc-view>
