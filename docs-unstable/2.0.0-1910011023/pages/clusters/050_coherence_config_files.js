<doc-view>

<v-layout row wrap>
<v-flex xs12 sm10 lg10>
<v-card class="section-def" v-bind:color="$store.state.currentColor">
<v-card-text class="pa-3">
<v-card class="section-def__card">
<v-card-text>
<dl>
<dt slot=title>Coherence Config Files</dt>
<dd slot="desc"><p>The different configuration files commonly used by Coherence can be specified for a role in the role&#8217;s spec.</p>
</dd>
</dl>
</v-card-text>
</v-card>
</v-card-text>
</v-card>
</v-flex>
</v-layout>

<h2 id="_coherence_config_files">Coherence Config Files</h2>
<div class="section">
<p>There are three Coherence configuration files that can be set in a role&#8217;s specification:</p>

<ul class="ulist">
<li>
<p>The <router-link to="#cache-config" @click.native="this.scrollFix('#cache-config')">Coherence Cache Configuration</router-link> file</p>

</li>
<li>
<p>The <router-link to="#override-file" @click.native="this.scrollFix('#override-file')">Coherence Operational Override</router-link> file</p>

</li>
</ul>
<div class="admonition note">
<p class="admonition-inline">Oracle Coherence provides a number of ways to specify the configuration files to be used. If deploying application
code alongside the Coherence JVM it is entirely possible to use the default configuration file names in the application
or to hard code one or more of the configuration file names in the application itself. This section of the documentation
describes <code>CoherenceCluster</code> role configuration that will result in passing the <code>-Dcoherence.cache.config</code> and the
<code>coherence.override</code> system properties to the Coherence JVM. If the deployed application overrides or ignores these
properties then setting the configurations described below will have no effect.</p>
</div>
</div>

<h2 id="cache-config">Setting the Coherence Cache Configuration File</h2>
<div class="section">
<p>The Coherence cache configuration file for a role in a <code>CoherenceCluster</code> is set using the <code>cacheConfig</code> field of a role spec.
The value of this field will end up being passed to the Coherence JVM as the <code>coherence.cache.config</code> system property and
will hence set the value of the cache configuration file used as described in the Coherence documentation.</p>


<h3 id="_set_the_cache_configuration_for_an_implicit_role">Set the Cache Configuration for an Implicit Role</h3>
<div class="section">
<p>When using the implicit role configuration the <code>cacheConfig</code> value is set directly in the <code>CoherenceCluster</code>
<code>spec.coherence</code> section.
For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  coherence:
    cacheConfig: application-cache-config.xml  <span class="conum" data-value="1" /></markup>

<ul class="colist">
<li data-value="1">In this case a cluster will be created with a single implicit role named <code>storage</code> where the <code>coherence.cache.config</code>
system property and hence the cache configuration file used will be <code>application-cache-config.xml</code></li>
</ul>
</div>

<h3 id="_set_the_cache_configuration_for_explicit_role">Set the Cache Configuration for Explicit Role</h3>
<div class="section">
<p>When using the explicit role configuration the <code>cacheConfig</code> value is set for each role in the <code>CoherenceCluster</code> <code>spec</code>
<code>roles</code> list.</p>

<p>For example to create cluster with two explicit roles, <code>data</code> and <code>proxy</code>:</p>

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
        cacheConfig: data-cache-config.xml  <span class="conum" data-value="1" />
    - role: proxy
      coherence:
        cacheConfig: proxy-cache-config.xml  <span class="conum" data-value="2" /></markup>

<ul class="colist">
<li data-value="1">The <code>data</code> role will use the <code>data-cache-config.xml</code> cache configuration file</li>
<li data-value="2">The <code>proxy</code> role will use the <code>proxy-cache-config.xml</code> cache configuration file</li>
</ul>
</div>

<h3 id="_set_the_cache_configuration_for_explicit_roles_with_a_default">Set the Cache Configuration for Explicit Roles with a Default</h3>
<div class="section">
<p>When using the explicit role configuration a value for <code>cacheConfig</code> value can be set in the <code>CoherenceCluster</code> <code>spec</code>
section that will be used as the default <code>cacheConfig</code> value for any <code>role</code> in the <code>roles</code> list that does not explicitly
specify a value.</p>

<p>For example to create cluster with three explicit roles, <code>data</code> and <code>proxy</code> and <code>web</code>:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  coherence:
    cacheConfig: application-cache-config.xml  <span class="conum" data-value="1" />
  roles:
    - role: data
    - role: proxy
    - role: web
      coherence:
        cacheConfig: web-cache-config.xml  <span class="conum" data-value="2" /></markup>

<ul class="colist">
<li data-value="1">The default <code>cacheConfig</code> value is set to <code>application-cache-config.xml</code>. The <code>data</code> and <code>proxy</code> roles do not have
a <code>cacheConfig</code> value so will use this default value and will each have use the <code>application-cache-config.xml</code> file</li>
<li data-value="2">The <code>web</code> role has an explicit <code>cacheConfig</code> value of <code>web-cache-config.xml</code> so will use the <code>web-cache-config.xml</code>
cache configuration file</li>
</ul>
</div>
</div>

<h2 id="override-file">Setting the Coherence Operational Override File</h2>
<div class="section">
<p>The Coherence operational override file for a role in a <code>CoherenceCluster</code> is set using the <code>overrideConfig</code> field of a role spec.
The value of this field will end up being passed to the Coherence JVM as the <code>coherence.override</code> system property and
will hence set the value of the operational override file used as described in the Coherence documentation.</p>


<h3 id="_set_the_operational_override_for_an_implicit_role">Set the Operational Override for an Implicit Role</h3>
<div class="section">
<p>When using the implicit role configuration the <code>overrideConfig</code> value is set directly in the <code>CoherenceCluster</code> <code>spec</code> section.
For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  coherence:
    overrideConfig: application-override.xml  <span class="conum" data-value="1" /></markup>

<ul class="colist">
<li data-value="1">In this case a cluster will be created with a single implicit role named <code>storage</code> where the <code>coherence.override</code>
system property and hence the operational override file used will be <code>application-override.xml</code></li>
</ul>
</div>

<h3 id="_set_the_operational_override_for_explicit_role">Set the Operational Override for Explicit Role</h3>
<div class="section">
<p>When using the explicit role configuration the <code>overrideConfig</code> value is set for each role in the <code>CoherenceCluster</code> <code>spec</code>
<code>roles</code> list.</p>

<p>For example to create cluster with two explicit roles, <code>data</code> and <code>proxy</code>:</p>

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
        overrideConfig: data-override.xml  <span class="conum" data-value="1" />
    - role: proxy
      coherence:
        overrideConfig: proxy-override.xml  <span class="conum" data-value="2" /></markup>

<ul class="colist">
<li data-value="1">The <code>data</code> role will use the <code>data-override.xml</code> operational override file</li>
<li data-value="2">The <code>proxy</code> role will use the <code>proxy-override.xml</code> operational override file</li>
</ul>
</div>

<h3 id="_set_the_operational_override_for_explicit_roles_with_a_default">Set the Operational Override for Explicit Roles with a Default</h3>
<div class="section">
<p>When using the explicit role configuration a value for <code>overrideConfig</code> value can be set in the <code>CoherenceCluster</code> <code>spec</code>
section that will be used as the default <code>overrideConfig</code> value for any <code>role</code> in the <code>roles</code> list that does not explicitly
specify a value.</p>

<p>For example to create cluster with three explicit roles, <code>data</code> and <code>proxy</code> and <code>web</code>:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  coherence:
    overrideConfig: application-override.xml  <span class="conum" data-value="1" />
  roles:
    - role: data
    - role: proxy
    - role: web
      coherence:
        overrideConfig: web-override.xml  <span class="conum" data-value="2" /></markup>

<ul class="colist">
<li data-value="1">The default <code>overrideConfig</code> value is set to <code>application-override.xml</code>. The <code>data</code> and <code>proxy</code> roles do not have
an <code>overrideConfig</code> value so will use this default value and will each have use the <code>application-override.xml</code> file</li>
<li data-value="2">The <code>web</code> role has an explicit <code>overrideConfig</code> value of <code>web-override.xml</code> so will use the <code>web-override.xml</code>
operational override file</li>
</ul>
</div>
</div>
</doc-view>
