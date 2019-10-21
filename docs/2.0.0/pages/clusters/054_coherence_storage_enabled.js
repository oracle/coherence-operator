<doc-view>

<v-layout row wrap>
<v-flex xs12 sm10 lg10>
<v-card class="section-def" v-bind:color="$store.state.currentColor">
<v-card-text class="pa-3">
<v-card class="section-def__card">
<v-card-text>
<dl>
<dt slot=title>Storage Enabled or Disabled Roles</dt>
<dd slot="desc"><p>A Coherence cluster member can be storage enabled or storage disabled and hence a <code>CoherenceCluster</code> role
can be configured to be storage enabled or disabled.</p>
</dd>
</dl>
</v-card-text>
</v-card>
</v-card-text>
</v-card>
</v-flex>
</v-layout>

<h2 id="_storage_enabled_or_disabled_roles">Storage Enabled or Disabled Roles</h2>
<div class="section">
<p>Coherence has a default System property that configures cache services to be storage enabled (i.e. that JVM will be
manage data for caches) or storage disabled (i.e. that member will be not manage data for caches).
A role in a <code>CoherenceCluster</code> can be set as storage enabled or disabled using the <code>storageEnabled</code> field; the value
is a boolean true or false. Setting this property sets the Coherence JVM system property <code>coherence.distributed.localstorage</code>
to true or false.</p>

<p>If the <code>storageEnabled</code> field is not specifically set for a role then the <code>coherence.distributed.localstorage</code> property
will not be set in the JVMs for that role and Coherence&#8217;s default behaviour will apply.</p>

<div class="admonition note">
<p class="admonition-inline">If a custom application is deployed into the Coherence container that specifies a custom cache configuration file
or custom operational configuration file it is entirely possible for the <code>coherence.distributed.localstorage</code> system
property to be ignored if the application configuration files override this value. If this is the case then the settings
described below will have no effect.</p>
</div>
</div>

<h2 id="_storage_enabled_or_disabled_implicit_role">Storage Enabled or Disabled Implicit Role</h2>
<div class="section">
<p>When creating a <code>CoherenceCluster</code> with the single implicit role the <code>storageEnabled</code> field is set in the <code>CoherenceCluster</code>
<code>spec.coherence</code> field. For example</p>

<markup
lang="yaml"
title="Storage Enabled Role"
>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  coherence:
    storageEnabled: true <span class="conum" data-value="1" /></markup>

<ul class="colist">
<li data-value="1">The implicit role will be storage enabled</li>
</ul>
<markup
lang="yaml"
title="Storage Disabled Role"
>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  coherence:
    storageEnabled: false <span class="conum" data-value="1" /></markup>

<ul class="colist">
<li data-value="1">The implicit role will be storage disabled</li>
</ul>
</div>

<h2 id="_storage_enabled_or_disabled_explicit_roles">Storage Enabled or Disabled Explicit Roles</h2>
<div class="section">
<p>When creating a <code>CoherenceCluster</code> with the explicit roles the <code>storageEnabled</code> field is set for each role in
the <code>CoherenceCluster</code> <code>roles</code> list.</p>

<markup
lang="yaml"
title="Storage Enabled Role"
>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  roles:
    - role: data
      coherence:
        storageEnabled: true <span class="conum" data-value="1" />
    - role: proxy
      coherence:
        storageEnabled: false <span class="conum" data-value="2" /></markup>

<ul class="colist">
<li data-value="1">The <code>data</code> role will be storage enabled</li>
<li data-value="2">The <code>proxy</code> role will be storage disabled</li>
</ul>
</div>

<h2 id="_storage_enabled_or_disabled_explicit_roles_with_defaults">Storage Enabled or Disabled Explicit Roles With Defaults</h2>
<div class="section">
<p>When creating a <code>CoherenceCluster</code> with the explicit roles the <code>storageEnabled</code> field is set for each role in
the <code>CoherenceCluster</code> <code>roles</code> list and a default can be set in the <code>CoherenceCluster</code> <code>spec</code>.</p>

<markup
lang="yaml"
title="Storage Enabled Role"
>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  coherence:
    storageEnabled: false     <span class="conum" data-value="1" />
  roles:
    - role: data              <span class="conum" data-value="2" />
      coherence:
        storageEnabled: true
    - role: proxy             <span class="conum" data-value="3" />
    - role: web               <span class="conum" data-value="4" /></markup>

<ul class="colist">
<li data-value="1">The default value will be storage disabled</li>
<li data-value="2">The <code>data</code> role overrides the default and will be storage enabled</li>
<li data-value="3">The <code>proxy</code> role does not have a specific <code>storageEnabled</code> so will be storage disabled</li>
<li data-value="4">The <code>web</code> roles does not have a specific <code>storageEnabled</code> so will be storage disabled</li>
</ul>
</div>
</doc-view>
