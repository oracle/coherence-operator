<doc-view>

<v-layout row wrap>
<v-flex xs12 sm10 lg10>
<v-card class="section-def" v-bind:color="$store.state.currentColor">
<v-card-text class="pa-3">
<v-card class="section-def__card">
<v-card-text>
<dl>
<dt slot=title>Overview</dt>
<dd slot="desc"><p>The <code>Coherence</code> resource has a number of fields to configure the behaviour of <code>Coherence</code>,
these fields are in the <code>spec.coherence</code> section of the CRD.</p>
</dd>
</dl>
</v-card-text>
</v-card>
</v-card-text>
</v-card>
</v-flex>
</v-layout>

<h2 id="_configuring_coherence">Configuring Coherence</h2>
<div class="section">
<p>The <code>Coherence</code> CRD has specific fields to configure the most common Coherence settings.
Any other settings can be configured by adding system properties to the <router-link to="/docs/jvm/010_overview">JVM Settings</router-link>.</p>

<p>The following Coherence features can be directly specified in the <code>Coherence</code> spec.</p>

<ul class="ulist">
<li>
<p><router-link to="/docs/coherence/020_cluster_name">Cluster Name</router-link></p>

</li>
<li>
<p><router-link to="/docs/coherence/030_cache_config">Cache Configuration File</router-link></p>

</li>
<li>
<p><router-link to="/docs/coherence/040_override_file">Operational Configuration File</router-link> (aka, the override file)</p>

</li>
<li>
<p><router-link to="/docs/coherence/050_storage_enabled">Storage Enabled</router-link> or disabled deployments</p>

</li>
<li>
<p><router-link to="/docs/coherence/060_log_level">Log Level</router-link></p>

</li>
<li>
<p><router-link to="/docs/coherence/070_wka">Well Known Addressing</router-link> and cluster discovery</p>

</li>
<li>
<p><router-link to="/docs/coherence/080_persistence">Persistence</router-link></p>

</li>
<li>
<p><router-link to="/docs/management/010_overview">Management over REST</router-link></p>

</li>
<li>
<p><router-link to="/docs/metrics/010_overview">Metrics</router-link></p>

</li>
</ul>
<div class="admonition note">
<p class="admonition-inline">The Coherence settings in the <code>Coherence</code> CRD spec typically set system property values that will
be passed through to the Coherence JVM command line, which in turn configure Coherence.
This is the same behaviour that would occur when running Coherence outside of containers.
Whether these system properties actually apply or not depends on the application code. For example,
it is simple to override the Coherence operational configuration file in a jar file deployed as part of an
application&#8217;s image in such a way that will cause all the normal Coherence system properties to be ignored.
If that is done then the Coherence settings discussed in this documentation will not apply.<br>
For example, adding a <code>tangosol-coherence-override.xml</code> file to a jar on the application&#8217;s classpath that contains
an overridden <code>&lt;configurable-cache-factory-config&gt;</code> section with a hard coded cache configuration file name would
mean that the <code>Coherence</code> CRD <code>spec.coherence.cacheConfig</code> field, that sets the <code>coherence.cacheconfig</code> system
property, would be ignored.<br>
It is, therefore, entirely at the application developer&#8217;s discretion whether they use the fields of the <code>Coherence</code> CRD
to configure Coherence, or they put those settings into configuration files, either hard coded into jar files or
picked up at runtime from files mapped from Kubernetes volumes, config maps, secrets, etc.</p>
</div>
</div>
</doc-view>
