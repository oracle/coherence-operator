<doc-view>

<h2 id="_configure_services_for_ports">Configure Services for Ports</h2>
<div class="section">
<p>As described in the <router-link to="/docs/ports/020_container_ports">Additional Container Ports</router-link> documentation,
it is possible to expose additional ports on the Coherence container in the Pods of a <code>Coherence</code> resource.
The Coherence Operator will create a <code>Service</code> to expose each additional port.
By default, the name of the service is the combination of the <code>Coherence</code> resource name and the port name
(this can default behaviour can be overridden as shown below in the <router-link to="#_override_the_service_name" @click.native="this.scrollFix('#_override_the_service_name')"></router-link> section).
The configuration of the <code>Service</code> can be altered using fields in the port spec&#8217;s <code>service</code> section.</p>

<p>For example:</p>

<markup
lang="yaml"
title="test-cluster.yaml"
>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: test-cluster
spec:
  ports:
    - name: rest   <span class="conum" data-value="1" />
      port: 8080
      service: {}  <span class="conum" data-value="2" /></markup>

<ul class="colist">
<li data-value="1">This example exposes a single port named <code>rest</code> on port <code>8080</code>.</li>
<li data-value="2">The <code>service</code> section of the port spec is empty so the Operator will use its default behaviour
to create a <code>Service</code> in the same namespace with the name <code>test-cluster-rest</code>.</li>
</ul>
</div>

<h2 id="_override_the_service_name">Override the Service Name</h2>
<div class="section">
<p>Sometimes it is useful to use a different name than the default for a <code>Service</code> for a port,
for example, when the port is exposing functionality that other applications want to consume on a fixed well know endpoint.
To override the generated service name with another name the <code>service.name</code> field can be set.</p>

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
      service:
        name: payments  <span class="conum" data-value="1" /></markup>

<ul class="colist">
<li data-value="1">By setting the <code>service.name</code> field the <code>Service</code> for this port will be named <code>payments</code>.</li>
</ul>
<p>The service for the above example would look like this:</p>

<markup
lang="yaml"

>apiVersion: v1
kind: Service
metadata:
  name: payments  <span class="conum" data-value="1" />
spec:
  ports:
    - name: rest
      port: 8080
      targetPort: rest
  type: ClusterIP
  selector:
    coherenceDeployment: test-cluster
    coherenceCluster: test-cluster
    coherenceRole: storage
    coherenceComponent: coherencePod</markup>

<ul class="colist">
<li data-value="1">The <code>Service</code> name is <code>payments</code> instead of <code>test-cluster-rest</code></li>
</ul>
</div>

<h2 id="_override_the_service_port">Override the Service Port</h2>
<div class="section">
<p>It is sometimes useful to be able to expose a service on a different port on the <code>Service</code> to that being used by the container.
One use-case for this would be where the <code>Coherence</code> deployment is providing a http service where the container
exposes the service on port <code>8080</code> whereas the <code>Service</code> can use port <code>80</code>.</p>

<p>For example, using the same example payemnts service above:</p>

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
      service:
        name: payments  <span class="conum" data-value="1" />
        port: 80        <span class="conum" data-value="2" /></markup>

<ul class="colist">
<li data-value="1">The <code>Service</code> name will be <code>payments</code></li>
<li data-value="2">The <code>Service</code> port will be <code>80</code></li>
</ul>
<p>This then allows the payments service to be accessed on a simple url of <code><a id="" title="" target="_blank" href="http://payments">http://payments</a></code></p>

</div>

<h2 id="_disable_service_creation">Disable Service Creation</h2>
<div class="section">
<p>Sometimes it may be desirable to expose a port on the Coherence container but not have the Operator automatically
create a <code>Service</code> to expose the port. For example, maybe the port is to be exposed via some other load balancer
service controlled by another system.
To disable automatic service creation set the <code>service.enabled</code> field to <code>false</code>.</p>

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
      service:
        enabled: false  <span class="conum" data-value="1" /></markup>

<ul class="colist">
<li data-value="1">With the <code>service.enabled</code> field set to <code>false</code> no <code>Service</code> will be created.</li>
</ul>
</div>

<h2 id="_other_service_configuration">Other Service Configuration</h2>
<div class="section">
<p>The <code>Coherence</code> resource CRD allows many other settings to be configured on the <code>Service</code>.
These fields are identical to the corresponding fields in the Kubernetes <code>Service</code> spec.</p>

<p>See the <code>Coherence</code> CRD <router-link to="#_servicespec" @click.native="this.scrollFix('#_servicespec')">Service Spec</router-link> documentation
and the Kubernetes
<a id="" title="" target="_blank" href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#servicespec-v1-core">Service API reference</a>.</p>

</div>
</doc-view>
