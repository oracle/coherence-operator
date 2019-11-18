<doc-view>

<v-layout row wrap>
<v-flex xs12 sm10 lg10>
<v-card class="section-def" v-bind:color="$store.state.currentColor">
<v-card-text class="pa-3">
<v-card class="section-def__card">
<v-card-text>
<dl>
<dt slot=title>Enabling SSL</dt>
<dd slot="desc"><p>Coherence clusters can be deployed with a metrics endpoint enabled that can be scraped by common metrics applications
such as Prometheus.</p>
</dd>
</dl>
</v-card-text>
</v-card>
</v-card-text>
</v-card>
</v-flex>
</v-layout>

<h2 id="_enabling_ssl_for_metrics_capture">Enabling SSL for metrics capture</h2>
<div class="section">
<div class="admonition note">
<p class="admonition-inline">Note: Use of metrics is available only when using the operator with clusters running
Coherence 12.2.1.4 or later version.</p>
</div>
<p>Please see <router-link :to="{path: '/clusters/060_coherence_metrics', hash: '#ssl'}">Coherence Metrics Documentation</router-link> for information on how to enable SSL.</p>

</div>
</doc-view>
