<doc-view>

<h2 id="_expose_ports_and_services">Expose Ports and Services</h2>
<div class="section">
<p>Any ports that are used by Coherence or by application code that need to be exposed outside of the <code>Pods</code> for a role
need to be declared in the <code>CoherenceCluster</code> spec for the role.</p>

</div>

<h2 id="_default_ports">Default Ports</h2>
<div class="section">
<p>The Coherence container in <code>Pods</code> in a role in a <code>CoherenceCluster</code> has two ports declared by default, none of the ports
are exposed on services.</p>


<div class="table__overflow elevation-1 ">
<table class="datatable table">
<colgroup>
<col style="width: 33.333%;">
<col style="width: 33.333%;">
<col style="width: 33.333%;">
</colgroup>
<thead>
<tr>
<th>Port</th>
<th>Name</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td><code>7</code></td>
<td><code>coherence</code></td>
<td>This is the standard echo port. Nothing in the container binds to this port it is only declared on the container so
that the headless <code>Service</code> used for Coherence WKA can declare a port.</td>
</tr>
<tr>
<td><code>6676</code></td>
<td><code>health</code></td>
<td>This is the port used to expose the default readiness, liveness and StatusHA ReST endpoints on.</td>
</tr>
</tbody>
</table>
</div>
<div class="admonition note">
<p class="admonition-inline">When exposing additional ports as described in the sections below the names for the additional ports cannot be
either <code>coherence</code> or <code>health</code> that are the names used for the default ports above or the <code>Pods</code> may fail to start.</p>
</div>
</div>

<h2 id="_configure_additional_ports">Configure Additional Ports</h2>
<div class="section">
<p>Additional ports can be declared for a role by adding them to the <code>ports</code> array that is part of the <code>role</code> spec.
A <code>port</code> has the following fields:</p>

<markup
lang="yaml"

>ports:
  - name: extend   <span class="conum" data-value="1" />
    port: 20000    <span class="conum" data-value="2" />
    protocol: TCP  <span class="conum" data-value="3" />
    service: {}    <span class="conum" data-value="4" /></markup>

<ul class="colist">
<li data-value="1">The port must have a <code>name</code> that is unique within the role</li>
<li data-value="2">The <code>port</code> value must be specified</li>
<li data-value="3">The <code>protocol</code> is optional and defaults to <code>TCP</code>. The valid values are the same as when declaring ports for
<a id="" title="" target="_blank" href="https://kubernetes.io/docs/concepts/services-networking/service/">Kubernetes services</a> and <code>Pods</code>.</li>
<li data-value="4">The <code>service</code> section is optional and is used to configure the <code>Service</code> that will be used to expose the port.
see <router-link to="#services" @click.native="this.scrollFix('#services')">Configure Services for Additional Ports</router-link></li>
</ul>
<p>By default a Kubernetes <code>Service</code> of type <code>ClusterIP</code> will be created for each additional port. The <code>Service</code> <code>port</code>
and <code>targetPort</code> will both default to the specified <code>port</code> value. The <code>port</code> value for the <code>Service</code> can be overridden
in the <code>service</code> spec.</p>

<p>The name of the <code>Service</code> created will default to a name made up from the cluster name, the role name and the port name
in the format <code>&lt;cluster-name&gt;-&lt;role-name&gt;-&lt;port-name&gt;</code>. This can be overriden by specifying a different name in the
<code>service</code> section of the additional port configuration.</p>

<p>For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster  <span class="conum" data-value="1" />
spec:
  role: data          <span class="conum" data-value="2" />
  ports:
    - name: extend    <span class="conum" data-value="3" />
      port: 20000</markup>

<ul class="colist">
<li data-value="1">The cluster name is <code>test-cluster</code></li>
<li data-value="2">The role name is <code>data</code></li>
<li data-value="3">The port name is <code>extend</code></li>
</ul>
<p>The <code>Service</code> created for the <code>extend</code> port would be:</p>

<markup
lang="yaml"

>apiVersion: v1
kind: Service
metadata:
  name: test-cluster-data-extend  <span class="conum" data-value="1" /></markup>

<ul class="colist">
<li data-value="1">the <code>Service</code> name is <code>test-cluster-data-extend</code> made up of the cluster name <code>test-cluster</code> the role name <code>data</code>
and the port name <code>extend</code>.</li>
</ul>

<h3 id="_configure_additional_ports_for_the_implicit_role">Configure Additional Ports for the Implicit Role</h3>
<div class="section">
<p>When configuring a <code>CoherenceCluster</code> with a single implicit role the additional ports are added to the <code>ports</code> array
in the <code>CoherenceCluster</code> <code>spec</code> section.
For Example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  ports:            <span class="conum" data-value="1" />
    - name: extend
      port: 20000
    - name: rest
      port: 8080</markup>

<ul class="colist">
<li data-value="1">The <code>ports</code> array for the single implicit role contains two additional ports. The first named <code>extend</code> on port
<code>20000</code> and the second named <code>rest</code> on port <code>8080</code>. Both of the ports in the above example will be exposed on separate
<code>Services</code> using the default service configuration.</li>
</ul>
</div>

<h3 id="_configure_additional_ports_for_explicit_roles">Configure Additional Ports for Explicit Roles</h3>
<div class="section">
<p>When configuring a <code>CoherenceCluster</code> with explicit roles in the <code>roles</code> list the additional ports are added to
the <code>ports</code> array for each role.
For Example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  roles:
    - role: data             <span class="conum" data-value="1" />
      ports:
        - name: management
          port: 30000
    - role: proxy            <span class="conum" data-value="2" />
      ports:
        - name: extend
          port: 20000</markup>

<ul class="colist">
<li data-value="1">The <code>data</code> role adds an additional port named <code>management</code> with a port value of <code>30000</code> that will be exposed on
a service named <code>test-cluster-data-management</code>.</li>
<li data-value="2">The <code>proxy</code> role adds an additional port named <code>extend</code> with a port value of <code>20000</code> that will be exposed on
a service named <code>test-cluster-data-extend</code>.</li>
</ul>
</div>

<h3 id="_configure_additional_ports_for_explicit_roles_with_defaults">Configure Additional Ports for Explicit Roles with Defaults</h3>
<div class="section">
<p>When configuring a <code>CoherenceCluster</code> with explicit roles default additional ports can be added to the
<code>CoherenceCluster</code> <code>spec.ports</code> array that will apply to all roles in the <code>roles</code> list.
Additional ports can then also be specified for individual roles in the <code>roles</code> list.
The <code>ports</code> array for an individual role will then be a <strong>merge</strong> of the default ports and the role&#8217;s ports.
If a port in a role has the same name as a default port then the role&#8217;s port will override the default port.</p>

<p>For Example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  ports:                   <span class="conum" data-value="1" />
    - name: management
      port: 30000
  roles:
    - role: data           <span class="conum" data-value="2" />
    - role: proxy
      ports:
        - name: extend     <span class="conum" data-value="3" />
          port: 20000
    - role: web
      ports:
        - name: http       <span class="conum" data-value="4" />
          port: 8080
        - name: management <span class="conum" data-value="5" />
          port: 9000</markup>

<ul class="colist">
<li data-value="1">The default additional ports section specifies a single additional port named <code>management</code> on port <code>30000</code>.</li>
<li data-value="2">The <code>data</code> role does not specify any additional ports so will just have the default additional <code>management</code> port
that will be exposed on a service named <code>test-cluster-data-management</code>.</li>
<li data-value="3">The <code>proxy</code> role adds an additional port named <code>extend</code> with a port value of <code>20000</code> that will be exposed on
a service named <code>test-cluster-data-extend</code>. The <code>proxy</code> role will also have the default additional <code>management</code> port
exposed on a service named <code>test-cluster-proxy-management</code>.</li>
<li data-value="4">The <code>web</code> role specified an additional port named <code>http</code> on port <code>8080</code> that will be exposed on a service named
<code>test-cluster-web-http</code>.</li>
<li data-value="5">The <code>web</code> role also overrides the default <code>management</code> port changing the <code>port</code> value from <code>30000</code> to <code>9000</code>
that will be exposed on a service named <code>test-cluster-web-management</code>.</li>
</ul>
</div>
</div>

<h2 id="services">Configure Services for Additional Ports</h2>
<div class="section">
<p>A number of fields may be specified to configure the <code>Service</code> that will be created to expose the port.</p>

<markup
lang="yaml"

>  ports:
    - name: extend
      port: 20000
      protocol: TCP
      service:
        enabled: true                     <span class="conum" data-value="1" />
        name: test-cluster-data-extend    <span class="conum" data-value="2" />
        port: 20000                       <span class="conum" data-value="3" />
        type:
        externalName:
        sessionAffinity:
        publishNotReadyAddresses:
        externalTrafficPolicy:
        loadBalancerIP:
        healthCheckNodePort:
        loadBalancerSourceRanges: []
        annotations: {}
        sessionAffinityConfig: {}</markup>

<ul class="colist">
<li data-value="1">Optionally enable or disable creation of a <code>Service</code> for the port, the defautl value is <code>true</code>.</li>
<li data-value="2">Optionally override the default generated <code>Service</code> name.</li>
<li data-value="3">Optionally use a different port in the <code>Service</code> to that used by the <code>Container</code>.</li>
</ul>
<p>Apart from the <code>enabled</code> and <code>name</code> fields, all of the fields shown above have exactly the same meaning and default
behaviour that they do for a normal <a id="" title="" target="_blank" href="https://kubernetes.io/docs/concepts/services-networking/service/">Kubernetes Service</a></p>


<h3 id="_enabling_or_disabling_service_creation">Enabling or Disabling Service Creation</h3>
<div class="section">
<p>By default a <code>Service</code> will be created for all additional ports in the <code>ports</code> array. If for some reason this is not
required <code>Service</code> creation can be disabled by setting the <code>service.enabled</code> field to <code>false</code>. The additional port
will still be added as a named port to the Coherence <code>Container</code> spec in the <code>Pod</code>.
For example:</p>

<markup
lang="yaml"

>  ports:
    - name: extend
      port: 20000
      protocol: TCP
      service:
        enabled: false</markup>

</div>

<h3 id="_changing_a_service_name">Changing a Service Name</h3>
<div class="section">
<p>As already described above the name of a <code>Service</code> created for an additional port is a combination of cluster name, role
name and port name. This can be overridden by setting the <code>service.name</code> field to the required name of the <code>Service</code>.</p>

<div class="admonition note">
<p class="admonition-inline">Bear in mind when overriding <code>Service</code> names that they must be unique within the Kubernetes namespace that the
<code>CoherenceCluster</code> is being installed into.</p>
</div>
<p>For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  ports:
    - name: http
      port: 8080
      service:
        name: front-end
        port: 80</markup>

<p>In the example above the service name has been overridden to be <code>front-end</code> and the service port overridden to <code>80</code>,
which will generate a <code>Service</code> like the following:</p>

<markup
lang="yaml"

>apiVersion: v1
kind: Service
metadata:
  name: front-end        <span class="conum" data-value="1" />
spec:
  ports:
    - name: http         <span class="conum" data-value="2" />
      port: 80           <span class="conum" data-value="3" />
      targetPort: 8080   <span class="conum" data-value="4" /></markup>

<ul class="colist">
<li data-value="1">The <code>Service</code> name has been overridden to <code>front-end</code></li>
<li data-value="2">The port name is <code>http</code> the same as the name of the additional port in the role spec.</li>
<li data-value="3">The <code>port</code> is <code>80</code> which is the value from the additional port&#8217;s <code>service.port</code> field.</li>
<li data-value="4">The <code>targetPort</code> is <code>8080</code> which is the port that the container will use from the <code>port</code> value of the additional
port in the role spec.</li>
</ul>
</div>
</div>
</doc-view>
