<doc-view>

<v-layout row wrap>
<v-flex xs12 sm10 lg10>
<v-card class="section-def" v-bind:color="$store.state.currentColor">
<v-card-text class="pa-3">
<v-card class="section-def__card">
<v-card-text>
<dl>
<dt slot=title>Coherence Operator Installation</dt>
<dd slot="desc"><p>The Coherence Operator is available as a Docker image <code>oracle/coherence-operator:2.0.0-1909171023</code> that can
easily be installed into a Kubernetes cluster.</p>
</dd>
</dl>
</v-card-text>
</v-card>
</v-card-text>
</v-card>
</v-flex>
</v-layout>

<h2 id="_installation_options">Installation Options</h2>
<div class="section">
<p>There are two ways to install the Coherence Operator</p>

<ul class="ulist">
<li>
<p><router-link to="/install/03_helm_install">Using Helm</router-link> using the Coherence Operator Helm chart</p>

</li>
<li>
<p><router-link to="/install/04_manual_install">Manually</router-link> using Kubernetes APIs (e.g. <code>kubectl</code>)</p>

</li>
</ul>
</div>
</doc-view>
