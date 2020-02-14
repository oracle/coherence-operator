<doc-view>

<h2 id="_configure_coherence">Configure Coherence</h2>
<div class="section">
<p>The <code>coherence</code> section of the <code>role</code> spec in a <code>CoherenceCluster</code> contains the following fields and sections that may
be configured:</p>

<markup
lang="yaml"

>coherence:
  cacheConfig: coherence-cache-config.xml          <span class="conum" data-value="1" />
  overrideConfig: tangosol-coherence-override.xml  <span class="conum" data-value="2" />
  logLevel: 5                                      <span class="conum" data-value="3" />
  storageEnabled: true                             <span class="conum" data-value="4" />
  imageSpec: {}                                    <span class="conum" data-value="5" />
  management: {}                                   <span class="conum" data-value="6" />
  metrics: {}                                      <span class="conum" data-value="7" />
  persistence: {}                                  <span class="conum" data-value="8" />
  snapshot: {}                                     <span class="conum" data-value="9" />
  excludeFromWKA: false                            <span class="conum" data-value="10" /></markup>

<ul class="colist">
<li data-value="1">The <code>cacheConfig</code> field sets the name of the Coherence cache configuration file to use.
See <router-link to="/clusters/052_coherence_config_files">Coherence Config Files</router-link> for more details.</li>
<li data-value="2">The <code>overrideConfig</code> field sets the name of the Coherence operational override configuration file to use.
See <router-link to="/clusters/052_coherence_config_files">Coherence Config Files</router-link> for more details.</li>
<li data-value="3">The <code>logLevel</code> field sets the log level that Coherence should use.
See <router-link to="/clusters/100_logging">Logging Configuration</router-link> for more details.</li>
<li data-value="4">The <code>storageEnabled</code> field sets whether the role is storage enabled or not.
See <router-link to="/clusters/054_coherence_storage_enabled">Storage Enabled or Disabled Roles</router-link> for more details.</li>
<li data-value="5">The <code>imageSpec</code> section configures the Coherence image details such as image name, pull policy etc.
See <router-link to="/clusters/056_coherence_image">Setting the Coherence Image</router-link> for more details.</li>
<li data-value="6">The <code>management</code> configures how Coherence management over REST behaves, whether it is enabled, etc.
See <router-link to="/clusters/058_coherence_management">Coherence Management Over REST</router-link> for more details.</li>
<li data-value="7">The <code>metrics</code> configures how Coherence metrics behaves, whether it is enabled, etc.
See <router-link to="/clusters/060_coherence_metrics">Coherence Metrics</router-link> for more details.</li>
<li data-value="8">The <code>persistence</code> configures how Coherence management over REST behaves, whether it is enabled, etc.
See <router-link to="/clusters/062_coherence_persistence">Coherence Persistence</router-link> for more details.</li>
<li data-value="9">The <code>snapshot</code> configures how Coherence management over REST behaves, whether it is enabled, etc.
See <router-link to="#clusters/064_coherence_snapshots.adoc" @click.native="this.scrollFix('#clusters/064_coherence_snapshots.adoc')">Coherence Snapshots</router-link> for more details.</li>
<li data-value="10">The <code>excludeFromWKA</code> field configures whether the <code>Pods</code> for this role form part of the Coherence WKA list for the cluster,
see the <router-link to="/about/05_cluster_discovery">Cluster Discovery</router-link> documentation page for more details.</li>
</ul>
</div>
</doc-view>
