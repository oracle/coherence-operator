<doc-view>

<v-layout row wrap>
<v-flex xs12 sm10 lg10>
<v-card class="section-def" v-bind:color="$store.state.currentColor">
<v-card-text class="pa-3">
<v-card class="section-def__card">
<v-card-text>
<dl>
<dt slot=title>Prerequisites</dt>
<dd slot="desc"><p>Everything needed to install and run the Coherence Operator is listed below:</p>
</dd>
</dl>
</v-card-text>
</v-card>
</v-card-text>
</v-card>
</v-flex>
</v-layout>

<h2 id="_prerequisites">Prerequisites</h2>
<div class="section">
<ul class="ulist">
<li>
<p>Access to a Kubernetes v1.11.3+ cluster.</p>

</li>
<li>
<p>Access to Oracle Coherence Docker images.</p>

</li>
</ul>

<h3 id="_image_pull_secrets">Image Pull Secrets</h3>
<div class="section">
<p>In order for the Coherence Operator to be able to install Coherence clusters it needs to be able to pull Coherence
Docker images. These images are not available in public Docker repositories and will typically Kubernetes will need
authentication to be able to pull them. This is achived by creating pull secrets.
Pull secrets are not global and hence secrets will be required in the namespace(s) that Coherence
clusters will be installed into.</p>

</div>
</div>
</doc-view>
