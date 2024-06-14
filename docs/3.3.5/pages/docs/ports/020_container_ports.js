<doc-view>

<h2 id="_additional_container_ports">Additional Container Ports</h2>
<div class="section">
<p>Except for rare cases most applications deployed into a Kubernetes cluster will need to expose ports that
they provide services on to other applications.
This is covered in the Kubernetes documentation,
<a id="" title="" target="_blank" href="https://kubernetes.io/docs/concepts/services-networking/connect-applications-service/">Connect Applications with Services</a></p>

<p>The <code>Coherence</code> CRD makes it simple to expose ports and configure their services.
The CRD contains a field named <code>ports</code>, which is an array of named ports.
In the most basic configuration the only required values are the name and port to expose, for example:</p>

<markup
lang="yaml"
title="test-cluster.yaml"
>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: test-cluster
spec:
  ports:
    - name: rest  <span class="conum" data-value="1" />
      port: 8080</markup>

<ul class="colist">
<li data-value="1">This example exposes a single port named <code>rest</code> on port <code>8080</code>.</li>
</ul>
<p>When the example above is deployed the Coherence Operator will add configure the ports for the
Coherence container in the <code>Pods</code> to expose that port and will also create a <code>Service</code> for the port.</p>

<p>For example, the relevant snippet of the <code>StatefulSet</code> configuration would be:</p>

<markup
lang="yaml"

>apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: test-cluster
spec:
  template:
    spec:
      containers:
      - name: coherence
        ports:
          - name: rest           <span class="conum" data-value="1" />
            containerPort: 8080  <span class="conum" data-value="2" /></markup>

<ul class="colist">
<li data-value="1">The Operator has added the <code>rest</code> port to the <code>coherence</code> containers port list.
The <code>name</code> field in the <code>Coherence</code> CRD&#8217;s port spec maps to the <code>name</code> field in the Container port spec.</li>
<li data-value="2">The <code>port</code> field in the <code>Coherence</code> CRD&#8217;s port spec maps to the <code>containerPort</code> in the Container port spec.</li>
</ul>
<p>For each additional port the Operator will create a <code>Service</code> of type <code>ClusterIP</code> with a default configuration.
The name of the service will be the <code>Coherence</code> resource&#8217;s name with the port name appended to it,
so in this case it will be <code>test-cluster-rest</code>. The <code>Service</code> might look like this:</p>

<markup
lang="yaml"

>apiVersion: v1
kind: Service
metadata:
  name: test-cluster-rest                 <span class="conum" data-value="1" />
spec:
  ports:
    - name: rest                          <span class="conum" data-value="2" />
      port: 8080                          <span class="conum" data-value="3" />
      targetPort: rest                    <span class="conum" data-value="4" />
  type: ClusterIP                         <span class="conum" data-value="5" />
  selector:
    coherenceDeployment: test-cluster     <span class="conum" data-value="6" />
    coherenceCluster: test-cluster
    coherenceRole: storage
    coherenceComponent: coherencePod</markup>

<ul class="colist">
<li data-value="1">The <code>Service</code> name will be automatically generated (this can be overridden).</li>
<li data-value="2">The <code>ports</code> section will have just the single port being exposed by this service with the same name as the port.</li>
<li data-value="3">The <code>port</code> exposed by the <code>Service</code> will be the same as the container port value (this can be overridden).</li>
<li data-value="4">The target port will be set to the port being exposed from the container.</li>
<li data-value="5">The default <code>Service</code> type is <code>ClusterIP</code> (this can be overridden).</li>
<li data-value="6">A selector will be created to match the <code>Pods</code> in the <code>Coherence</code> resource.</li>
</ul>
<p>The <code>Coherence</code> CRD spec allows port and service to be further configured and allows a
Prometheus <code>ServiceMonitor</code> to be created for the port if that port is to expose metrics.</p>

<p>See also:</p>

<ul class="ulist">
<li>
<p><router-link to="#ports/030_services.adoc" @click.native="this.scrollFix('#ports/030_services.adoc')">Configure Services for Ports</router-link></p>

</li>
<li>
<p><router-link to="#ports/040_servicemonitors.adoc" @click.native="this.scrollFix('#ports/040_servicemonitors.adoc')">Prometheus ServiceMonitors</router-link></p>

</li>
</ul>

<h3 id="_metrics_management_ports">Metrics &amp; Management Ports</h3>
<div class="section">
<p>Exposing the Coherence metrics port or Coherence Management over REST port are treated as a special case in the
configuration. Normally both the port&#8217;s <code>name</code> and <code>port</code> value are required fields. If the port name is <code>metrics</code>
or <code>management</code> the Operator already knows the <code>port</code> values (either from the defaults or from the metrics or
management configuration) so these do not need to be specified again.</p>

<p>For example, if the <code>Coherence</code> resource above also exposed Coherence metrics and management it might look like this:</p>

<markup
lang="yaml"
title="test-cluster.yaml"
>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: test-cluster
spec:
  coherence:
    metrics:
      enabled: true
      port: 9876
    management:
      enabled: true
      port: 1234
  ports:
    - name: rest         <span class="conum" data-value="1" />
      port: 8080
    - name: metrics      <span class="conum" data-value="2" />
    - name: management   <span class="conum" data-value="3" /></markup>

<ul class="colist">
<li data-value="1">The <code>rest</code> port is not a special case and must have a port defined, in this case <code>8080</code>.</li>
<li data-value="2">The <code>metrics</code> port is exposed, but the port is not required as the Operator already knows the port value,
which is configured in the <code>coherence.metrics</code> section to be 9876.</li>
<li data-value="3">The <code>management</code> port is exposed, but the port is not required as the Operator already knows the port value,
which is configured in the <code>coherence.management</code> section to be 1234.</li>
</ul>
<p>If the port value is not set in <code>coherence.metrics.port</code> or in <code>coherence.management.port</code> then the Operator will
use the defaults for these values, 9612 for metrics and 30000 for management.</p>

</div>
</div>

<h2 id="_configuring_the_port">Configuring the Port</h2>
<div class="section">
<p>The only mandatory fields when adding a port to a <code>Coherence</code> resource are the name and port number.
There are a number of optional fields, which when not specified use the Kubernetes default values.</p>

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
      protocol: TCP
      hostIP: 10.10.1.19
      hostPort: 1000
      nodePort: 5000</markup>

<p>The additional fields, <code>protocol</code>, <code>hostIP</code>, <code>hostPort</code> have the same meaning and same defaults in the
<code>Coherence</code> CRD port spec as they have in a Kubernetes container port
(see the Kubernetes <a id="" title="" target="_blank" href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.29/#containerport-v1-core">ContainerPort</a> API reference).
These fields map directly from the <code>Coherence</code> CRD port spec to the container port spec.</p>

<p>The example above would create a container port shown below:</p>

<markup
lang="yaml"

>apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: test-cluster
spec:
  template:
    spec:
      containers:
      - name: coherence
        ports:
          - name: rest
            containerPort: 8080
            protocol: TCP
            hostIP: 10.10.1.19
            hostPort: 1000</markup>

<p>The <code>nodePort</code> field in the <code>Coherence</code> CRD port spec maps to the <code>nodePort</code> field in the <code>Service</code> port spec.
The <code>nodePort</code> is described in the Kubernetes
<a id="" title="" target="_blank" href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.29/#serviceport-v1-core">ServicePort</a> API reference.</p>

<p>The <code>Coherence</code> CRD example above with <code>nodePort</code> set would create a <code>Service</code> with the same <code>nodePort</code> value:</p>

<markup
lang="yaml"

>apiVersion: v1
kind: Service
metadata:
  name: test-cluster-rest
spec:
  ports:
    - name: rest
      port: 8080
      targetPort: rest
      nodePort: 5000
  type: ClusterIP
  selector:
    coherenceDeployment: test-cluster
    coherenceCluster: test-cluster
    coherenceRole: storage
    coherenceComponent: coherencePod</markup>

</div>
</doc-view>
