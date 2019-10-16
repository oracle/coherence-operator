<doc-view>

<v-layout row wrap>
<v-flex xs12 sm10 lg10>
<v-card class="section-def" v-bind:color="$store.state.currentColor">
<v-card-text class="pa-3">
<v-card class="section-def__card">
<v-card-text>
<dl>
<dt slot=title>Using VisualVM</dt>
<dd slot="desc"><p><a id="" title="" target="_blank" href="https://visualvm.github.io/">VisualVM</a> is a visual tool integrating commandline JDK tools and lightweight profiling capabilities.
Designed for both development and production time use.</p>
</dd>
</dl>
</v-card-text>
</v-card>
</v-card-text>
</v-card>
</v-flex>
</v-layout>

<h2 id="_access_the_coherence_cluster_via_visualvm">Access the Coherence cluster via VisualVM</h2>
<div class="section">
<p>Coherence management is implemented using Java Management Extensions (JMX). JMX is a Java standard
for managing and monitoring Java applications and services. VisualVM and other JMX tools can be used to
manage and monitor Coherence Clusters via JMX.</p>

<p>This example shows how to connect to a cluster via VisualVM over JMXMP.</p>

<p>Please see <router-link to="#020_manegement_over_rest.adoc" @click.native="this.scrollFix('#020_manegement_over_rest.adoc')">Management over ReST</router-link> for how to connect
to a cluster via the VisualVM plugin using ReST.</p>

<div class="admonition note">
<p class="admonition-inline">See the <a id="" title="" target="_blank" href="https://docs.oracle.com/en/middleware/fusion-middleware/coherence/12.2.1.4/manage/introduction-oracle-coherence-management.html">Coherence Management Documentation</a>
for more information on JMX and Management.</p>
</div>
</div>
</doc-view>
