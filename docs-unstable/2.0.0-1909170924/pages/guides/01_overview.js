<doc-view>

<h2 id="_getting_started">Getting Started</h2>
<div class="section">
<v-layout row wrap class="mb-5">
<v-flex xs12>
<v-container fluid grid-list-md class="pa-0">
<v-layout row wrap class="pillars">
<v-flex xs12 sm4 lg3>
<v-card>
<router-link to="/install/01_introduction"><div class="card__link-hover"/>
</router-link>
<v-card-title primary class="headline layout justify-center">
<span style="text-align:center">Quickstart</span>
</v-card-title>
<v-card-text class="caption">
<p></p>
<p>Installing and running the Coherence Operator.</p>
</v-card-text>
</v-card>
</v-flex>
<v-flex xs12 sm4 lg3>
<v-card>
<router-link to="#guides/01_overview.adoc" @click.native="this.scrollFix('#guides/01_overview.adoc')"><div class="card__link-hover"/>
</router-link>
<v-card-title primary class="headline layout justify-center">
<span style="text-align:center">Install</span>
</v-card-title>
<v-card-text class="caption">
<p></p>
<p>Follow step-by-step guides to using the Coherence Operator to manage Coherence clusters.</p>
</v-card-text>
</v-card>
</v-flex>
</v-layout>
</v-container>
</v-flex>
</v-layout>
</div>

<h2 id="_more_guides">More Guides</h2>
<div class="section">
<v-layout row wrap class="mb-5">
<v-flex xs12>
<v-container fluid grid-list-md class="pa-0">
<v-layout row wrap class="pillars">
<v-flex xs12 sm4 lg3>
<v-card>
<router-link to="/guides/03_management"><div class="card__link-hover"/>
</router-link>
<v-card-title primary class="headline layout justify-center">
<span style="text-align:center">ReST Management API</span>
</v-card-title>
<v-card-text class="caption">
<p></p>
<p>Managing Coherence clusters with management over ReST.</p>
</v-card-text>
</v-card>
</v-flex>
<v-flex xs12 sm4 lg3>
<v-card>
<router-link to="/guides/04_metrics"><div class="card__link-hover"/>
</router-link>
<v-card-title primary class="headline layout justify-center">
<span style="text-align:center">Coherence Metrics</span>
</v-card-title>
<v-card-text class="caption">
<p></p>
<p>Publishing metrics from Coherence clusters.</p>
</v-card-text>
</v-card>
</v-flex>
<v-flex xs12 sm4 lg3>
<v-card>
<router-link to="#clusters/01_introduction.adoc" @click.native="this.scrollFix('#clusters/01_introduction.adoc')"><div class="card__link-hover"/>
</router-link>
<v-card-title primary class="headline layout justify-center">
<span style="text-align:center">Clusters</span>
</v-card-title>
<v-card-text class="caption">
<p></p>
<p>Managing Coherence clusters.</p>
</v-card-text>
</v-card>
</v-flex>
<v-flex xs12 sm4 lg3>
<v-card>
<router-link to="/developer/01_introduction"><div class="card__link-hover"/>
</router-link>
<v-card-title primary class="headline layout justify-center">
<span style="text-align:center">Developer</span>
</v-card-title>
<v-card-text class="caption">
<p></p>
<p>Developer guide for building the Coherence Operator.</p>
</v-card-text>
</v-card>
</v-flex>
</v-layout>
</v-container>
</v-flex>
</v-layout>
</div>
</doc-view>
