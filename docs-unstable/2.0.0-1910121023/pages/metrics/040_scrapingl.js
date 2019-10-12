<doc-view>

<v-layout row wrap>
<v-flex xs12 sm10 lg10>
<v-card class="section-def" v-bind:color="$store.state.currentColor">
<v-card-text class="pa-3">
<v-card class="section-def__card">
<v-card-text>
<dl>
<dt slot=title>Using Your Own Prometheus</dt>
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

<h2 id="_scraping_metrics_from_your_own_prometheus_instance">Scraping metrics from your own Prometheus instance</h2>
<div class="section">
<div class="admonition note">
<p class="admonition-inline">Note: Use of metrics is available only when using the operator with clusters running
Coherence 12.2.1.4 or later version.</p>
</div>
</div>
</doc-view>
