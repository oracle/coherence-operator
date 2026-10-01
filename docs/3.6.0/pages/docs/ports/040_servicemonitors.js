<doc-view>

<h2 id="_prometheus_servicemonitors">Prometheus ServiceMonitors</h2>
<div class="section">
<p>When a port exposed on a container is to be used to serve metrics to Prometheus this often requires the addition of
a Prometheus <code>ServiceMonitor</code> resource. The Coherence Operator makes it simple to add a <code>ServiceMonitor</code> for an exposed
port. The advantage of specifying the <code>ServiceMonitor</code> configuration in the <code>Coherence</code> CRD spec is that the
<code>ServiceMonitor</code> resource will be created, updated and deleted as part of the lifecycle of the <code>Coherence</code> resource,
and does not need to be managed separately.</p>

<p>A <code>ServiceMonitor</code> is created for an exposed port by setting the <code>serviceMonitor.enabled</code> field to <code>true</code>.
The Operator will create a <code>ServiceMonitor</code> with the same name as the <code>Service</code>.
The <code>ServiceMonitor</code> created will have a single endpoint for the port being exposed.</p>

<p>For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: test-cluster
spec:
  ports:
    - name: rest
      port: 8080
      serviceMonitor:
        enabled: true  <span class="conum" data-value="1" /></markup>

<ul class="colist">
<li data-value="1">With the <code>serviceMonitor.enabled</code> field set to <code>true</code> a <code>ServiceMonitor</code> resource will be created.</li>
</ul>
<p>The <code>ServiceMonitor</code> created from the spec above will look like this:
For example:</p>

<markup
lang="yaml"

>apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: test-cluster-rest
  labels:
    coherenceCluster: test-cluster
    coherenceComponent: coherence-service-monitor
    coherenceDeployment: test-cluster
    coherenceRole: test-cluster
spec:
  endpoints:
    - port: rest
      relabelings:
        - action: labeldrop
          regex: (endpoint|instance|job|service)
  selector:
    matchLabels:
      coherenceCluster: test-cluster
      coherenceComponent: coherence-service
      coherenceDeployment: test-cluster
      coherencePort: rest
      coherenceRole: test-cluster</markup>


<h3 id="_configure_the_servicemonitor">Configure the ServiceMonitor</h3>
<div class="section">
<p>The <code>Coherence</code> CRD <router-link :to="{path: '/docs/about/04_coherence_spec', hash: '#_servicemonitorspec'}">ServiceMonitorSpec</router-link>
contains many of the fields from the
<a id="" title="" target="_blank" href="https://prometheus-operator.dev/docs/api-reference/api/#monitoring.coreos.com/v1.ServiceMonitorSpec">Prometheus <code>ServiceMonitorSpec</code></a>
and <a id="" title="" target="_blank" href="https://prometheus-operator.dev/docs/api-reference/api/#monitoring.coreos.com/v1.Endpoint">Prometheus Endpoint</a>
to allow the <code>ServiceMonitor</code> to be configured for most use-cases.</p>

<p>In situations where the <code>Coherence</code> CRD does not have the required fields, for example when a different version
of Prometheus has been installed to that used to build the Coherence Operator, then the solution would be to
manually create <code>ServiceMonitors</code> instead of letting them be created by the Coherence Operator.</p>

</div>

<h3 id="_preserve_target_labels">Preserve target labels</h3>
<div class="section">
<p>Set <code>spec.ports[].serviceMonitor.useDefaultLabelDrop: false</code> to stop the built-in
rule from removing <code>endpoint</code>, <code>instance</code>, <code>job</code>, and <code>service</code>. Omitted or <code>true</code>
retains the existing <code>labeldrop (endpoint|instance|job|service)</code> rule, appended
after user relabelings.
This optional setting applies per port to both <code>Coherence</code> and <code>CoherenceJob</code>.
Install the matching Operator and CRDs before using it.</p>

<p>To drop only <code>endpoint</code> and <code>service</code>, configure a custom rule:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: example
spec:
  coherence:
    metrics:
      enabled: true
  ports:
    - name: metrics
      serviceMonitor:
        enabled: true
        useDefaultLabelDrop: false
        relabelings:
          - action: labeldrop
            regex: (endpoint|service)</markup>

<p>Omit the custom rule to retain all four labels. Merge these settings into each
existing affected port entry; do not replace the deployment&#8217;s complete ports list.
<code>coherence.metrics.enabled</code> enables the metrics endpoint, while
<code>ports[].serviceMonitor.enabled</code> separately requests a ServiceMonitor.</p>

<p>User <code>relabelings</code> retain their order and behavior, including rules that remove
<code>job</code> or <code>instance</code>. The option does not change <code>metricRelabelings</code> or <code>honorLabels</code>.
Target relabelings run before scraping; changing metric relabelings cannot repair
target labels needed for scrape reports.</p>

<p>Retaining target labels changes series label sets. With the existing default
<code>honorLabels: false</code>, conflicting application labels may become <code>exported_job</code>
or <code>exported_instance</code>. Check affected queries and dashboards; see the
<a id="" title="" target="_blank" href="https://prometheus.io/docs/prometheus/latest/configuration/configuration/#scrape_config">Prometheus scrape configuration</a>.
Restore <code>true</code> or remove the field to restore the built-in broad rule.</p>

</div>
</div>
</doc-view>
