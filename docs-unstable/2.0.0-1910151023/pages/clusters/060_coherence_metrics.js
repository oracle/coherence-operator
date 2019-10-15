<doc-view>

<h2 id="_coherence_metrics">Coherence Metrics</h2>
<div class="section">
<p>Since version 12.2.1.4 Coherence has had functionality to expose metrics via a http endpoint.
This endpoint is disabled by default in Coherence clusters but can be enabled and configured by setting the relevant
fields in the <code>CoherenceCluster</code> resource.</p>

</div>

<h2 id="_enabling_coherence_metrics">Enabling Coherence Metrics</h2>
<div class="section">
<p>Coherence metrics can be enabled or disabled by setting the <code>coherence.metrics.enabled</code> field.</p>

<div class="admonition note">
<p class="admonition-textlabel">Note</p>
<p ><p>Enabling metrics will add a number of <code>.jar</code> files to the classpath of the Coherence JVM.
In Coherence 12.2.1.4 those <code>.jar</code> file are:</p>

<markup


>org.glassfish.hk2.external:aopalliance-repackaged:jar:2.4.0-b34
org.glassfish.hk2:hk2-api:jar:2.4.0-b34
org.glassfish.hk2:hk2-locator:jar:2.4.0-b34
org.glassfish.hk2:hk2-utils:jar:2.4.0-b34
org.glassfish.hk2.external:javax.inject:jar:2.4.0-b34
com.fasterxml.jackson.core:jackson-annotations:jar:2.9.9
com.fasterxml.jackson.core:jackson-core:jar:2.9.9
com.fasterxml.jackson.core:jackson-databind:jar:2.9.9.2
com.fasterxml.jackson.jaxrs:jackson-jaxrs-base:jar:2.9.9
com.fasterxml.jackson.jaxrs:jackson-jaxrs-json-provider:jar:2.9.9
com.fasterxml.jackson.module:jackson-module-jaxb-annotations:jar:2.9.9
javax.annotation:javax.annotation-api:jar:1.2
javax.validation:validation-api:jar:1.1.0.Final
javax.ws.rs:javax.ws.rs-api:jar:2.0.1
org.glassfish.jersey.core:jersey-client:jar:2.22.4
org.glassfish.jersey.core:jersey-common:jar:2.22.4
org.glassfish.jersey.ext:jersey-entity-filtering:jar:2.22.4
org.glassfish.jersey.bundles.repackaged:jersey-guava:jar:2.22.4
org.glassfish.jersey.media:jersey-media-json-jackson:jar:2.22.4
org.glassfish.jersey.core:jersey-server:jar:2.22.4
org.glassfish.hk2:osgi-resource-locator:jar:1.0.1</markup>

<p>If adding additional application <code>.jar</code> files care should be taken that there are no version conflicts.</p>

<p>If conflicts are an issue there are alternative approaches available to exposing the metrics.</p>

<p>The list above is subject to change in later Coherence patches and version.</p>
</p>
</div>

<h3 id="_enabling_metrics_for_the_implicit_role">Enabling Metrics for the Implicit Role</h3>
<div class="section">
<p>When configuring a single implicit role in a <code>CoherenceCluster</code> metrics can be enabled by setting
the <code>coherence.metrics.enabled</code> to <code>true</code> in the <code>CoherenceCluster</code> <code>spec</code> section.
For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  coherence:
    metrics:
      enabled: true  <span class="conum" data-value="1" /></markup>

<ul class="colist">
<li data-value="1">Metrics will be enabled and the http endpoint will bind to port <code>9612</code> in the container.
The port is not exposed in a <code>Service</code>.</li>
</ul>
</div>

<h3 id="_enabling_metrics_for_explicit_roles">Enabling Metrics for Explicit Roles</h3>
<div class="section">
<p>When configuring a explicit roles in the <code>roles</code> list of a <code>CoherenceCluster</code> metrics can be
enabled or disabled by setting the <code>coherence.metrics.enabled</code> for each role.
For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  roles:
    - role: data
      coherence:
        metrics:
          enabled: true   <span class="conum" data-value="1" />
    - role: proxy
      coherence:
        metrics:
          enabled: false  <span class="conum" data-value="2" /></markup>

<ul class="colist">
<li data-value="1">The <code>data</code> role has the metrics enabled.</li>
<li data-value="2">The <code>proxy</code> role has the metrics disabled.</li>
</ul>
</div>

<h3 id="_enabling_metrics_for_explicit_roles_with_a_default">Enabling Metrics for Explicit Roles with a Default</h3>
<div class="section">
<p>When configuring a explicit roles in the <code>roles</code> list of a <code>CoherenceCluster</code> a default value for the
<code>coherence.metrics.enabled</code> field can be set in the <code>CoherenceCluster</code> <code>spec</code> section that will apply to
all roles in the <code>roles</code> list unless overridden for a specific role.
For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  coherence:
    metrics:
      enabled: true       <span class="conum" data-value="1" />
  roles:
    - role: data          <span class="conum" data-value="2" />
    - role: proxy
      coherence:
        metrics:
          enabled: false  <span class="conum" data-value="3" /></markup>

<ul class="colist">
<li data-value="1">The default value for enabling metrics is <code>true</code> which will apply to all roles in the <code>roles</code> list
unless the field is specifically overridden.</li>
<li data-value="2">The <code>data</code> role does not specify a value for the <code>coherence.metrics.enabled</code> field so it will use the default
value of <code>true</code> so metrics will be enabled.</li>
<li data-value="3">The <code>proxy</code> role overrides the default value for the <code>coherence.metrics.enabled</code> field and sets it to <code>false</code>
so metrics will be disabled.</li>
</ul>
</div>

<h3 id="_exposing_metrics_via_a_service">Exposing Metrics via a Service</h3>
<div class="section">
<p>Enabling metrics only enables the http server so that the endpoint is available in the container.
If external access to the API is required via a service then the port needs to be exposed just like any other
additional ports as described in <router-link to="/clusters/090_ports_and_services">Expose Ports and Services</router-link>.</p>

<p>For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  coherence:
    metrics:
      enabled: true   <span class="conum" data-value="1" />
  ports:
    - name: metrics   <span class="conum" data-value="2" />
      port: 9612</markup>

<ul class="colist">
<li data-value="1">metrics will be enabled and the default port value will be used so that the http endpoint will bind
to port <code>9612</code> in the container.</li>
<li data-value="2">An additional port named <code>metrics</code> is added to the <code>ports</code> array which will cause the metrics port to be
exposed on a service. The port specified is <code>9612</code> as that is the default port that metrics will bind to.</li>
</ul>
</div>

<h3 id="_expose_metrics_on_a_different_port">Expose metrics on a Different Port</h3>
<div class="section">
<p>The default port in the container that metrics uses is 9612. It is possible to change ths port using the
<code>coherence.metrics.port</code> field.</p>

<p>For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  coherence:
    metrics:
      enabled: true
      port: 9000        <span class="conum" data-value="1" />
  ports:
    - name: metrics
      port: 9000        <span class="conum" data-value="2" /></markup>

<ul class="colist">
<li data-value="1">metrics is enabled and configured to bind to port <code>9000</code> in the container.</li>
<li data-value="2">The corresponding <code>port</code> value of <code>9000</code> must be used when exposing the port on a <code>Service</code>.</li>
</ul>
</div>

<h3 id="_configuring_metrics_with_ssl">Configuring metrics With SSL</h3>
<div class="section">
<p>It is possible to configure metrics endpoint to use SSL to secure the communication between server and
client. The SSL configuration is in the <code>coherence.metrics.ssl</code> section of the spec.
See <router-link to="#metrics/020_manegement_over_rest.adoc" @click.native="this.scrollFix('#metrics/020_manegement_over_rest.adoc')">metrics</router-link> for a more in depth guide to configuring SSL.</p>

<p>For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  coherence:
    metrics:
      enabled: true
      ssl:
        enabled: true                            <span class="conum" data-value="1" />
        keyStore: metrics-keys.jks            <span class="conum" data-value="2" />
        keyStoreType: JKS                        <span class="conum" data-value="3" />
        keyStorePasswordFile: store-pass.txt     <span class="conum" data-value="4" />
        keyPasswordFile: key-pass.txt            <span class="conum" data-value="5" />
        keyStoreProvider:                        <span class="conum" data-value="6" />
        keyStoreAlgorithm: SunX509               <span class="conum" data-value="7" />
        trustStore: metrics-trust.jks         <span class="conum" data-value="8" />
        trustStoreType: JKS                      <span class="conum" data-value="9" />
        trustStorePasswordFile: trust-pass.txt   <span class="conum" data-value="10" />
        trustStoreProvider:                      <span class="conum" data-value="11" />
        trustStoreAlgorithm: SunX509             <span class="conum" data-value="12" />
        requireClientCert: true                  <span class="conum" data-value="13" />
        secrets: metrics-secret               <span class="conum" data-value="14" /></markup>

<ul class="colist">
<li data-value="1">The <code>enabled</code> field when set to <code>true</code> enables SSL for metrics or when set to <code>false</code> disables SSL</li>
<li data-value="2">The <code>keyStore</code> field sets the name of the Java key store file that should be used to obtain the server&#8217;s key</li>
<li data-value="3">The optional <code>keyStoreType</code> field sets the type of the key store file, the default value is <code>JKS</code></li>
<li data-value="4">The optional <code>keyStorePasswordFile</code> sets the name of the text file containing the key store password</li>
<li data-value="5">The optional <code>keyPasswordFile</code> sets the name of the text file containing the password of the key in the key store</li>
<li data-value="6">The optional <code>keyStoreProvider</code> sets the provider name for the key store</li>
<li data-value="7">The optional <code>keyStoreAlgorithm</code> sets the algorithm name for the key store, the default value is <code>SunX509</code></li>
<li data-value="8">The <code>trustStore</code> field sets the name of the Java trust store file that should be used to obtain the server&#8217;s key</li>
<li data-value="9">The optional <code>trustStoreType</code> field sets the type of the trust store file, the default value is <code>JKS</code></li>
<li data-value="10">The optional <code>trustStorePasswordFile</code> sets the name of the text file containing the trust store password</li>
<li data-value="11">The optional <code>trustStoreProvider</code> sets the provider name for the trust store</li>
<li data-value="12">The optional <code>trustStoreAlgorithm</code> sets the algorithm name for the trust store, the default value is <code>SunX509</code></li>
<li data-value="13">The optional <code>requireClientCert</code> field if set to <code>true</code> enables two-way SSL where the client must also provide
a valid certificate</li>
<li data-value="14">The optional <code>secrets</code> field sets the name of the Kubernetes <code>Secret</code> to use to obtain the key store, truct store
and password files from.</li>
</ul>
</div>
</div>
</doc-view>
