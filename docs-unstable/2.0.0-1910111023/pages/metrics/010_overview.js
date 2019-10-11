<doc-view>

<h2 id="_metrics">Metrics</h2>
<div class="section">
<v-layout row wrap class="mb-5">
<v-flex xs12>
<v-container fluid grid-list-md class="pa-0">
<v-layout row wrap class="pillars">
<v-flex xs12 sm4 lg3>
<v-card>
<router-link to="/metrics/020_metrics"><div class="card__link-hover"/>
</router-link>
<v-card-title primary class="headline layout justify-center">
<span style="text-align:center">Enabling Metrics</span>
</v-card-title>
<v-card-text class="caption">
<p></p>
<p>Deploying Coherence clusters with metrics enabled.</p>
</v-card-text>
</v-card>
</v-flex>
<v-flex xs12 sm4 lg3>
<v-card>
<router-link to="/metrics/030_ssl"><div class="card__link-hover"/>
</router-link>
<v-card-title primary class="headline layout justify-center">
<span style="text-align:center">Enabling SSL</span>
</v-card-title>
<v-card-text class="caption">
<p></p>
<p>Enabling SSL for metrics capture.</p>
</v-card-text>
</v-card>
</v-flex>
<v-flex xs12 sm4 lg3>
<v-card>
<router-link to="#metrics/040_scraping.adoc" @click.native="this.scrollFix('#metrics/040_scraping.adoc')"><div class="card__link-hover"/>
</router-link>
<v-card-title primary class="headline layout justify-center">
<span style="text-align:center">Using Your Own Prometheus</span>
</v-card-title>
<v-card-text class="caption">
<p></p>
<p>Scraping metrics from your own Prometheus instance.</p>
</v-card-text>
</v-card>
</v-flex>
</v-layout>
</v-container>
</v-flex>
</v-layout>
</div>
</doc-view>
