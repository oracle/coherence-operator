<doc-view>

<v-layout row wrap>
<v-flex xs12 sm10 lg10>
<v-card class="section-def" v-bind:color="$store.state.currentColor">
<v-card-text class="pa-3">
<v-card class="section-def__card">
<v-card-text>
<dl>
<dt slot=title>Importing Dashboards</dt>
<dd slot="desc"><p>If required, you can import the Grafana dashboards into your own Grafana instance.</p>
</dd>
</dl>
</v-card-text>
</v-card>
</v-card-text>
</v-card>
</v-flex>
</v-layout>

<h2 id="_importing_grafana_dashboards_into_your_own_instance">Importing Grafana dashboards into your own instance.</h2>
<div class="section">
<div class="admonition note">
<p class="admonition-inline">Note: Use of metrics is available only when using the operator with clusters running
Coherence 12.2.1.4 or later version.</p>
</div>
<p>This example shows you how to import the Grafana dashboards into your own Grafana instance.</p>

<p>By default the Coherence dashboards require a datasource named <code>Prometheus</code> which
should also be the default datasource.</p>

<p>If this datasource is already used and you cannot add another datasource as the default,
then please go to <router-link to="#different" @click.native="this.scrollFix('#different')">Importing with a different datasource</router-link>.</p>


<h3 id="importing">1. Importing using the defaults</h3>
<div class="section">
<p>In your Grafana environment, ensure you either:</p>

<ul class="ulist">
<li>
<p>have a Prometheus datasource named <code>Prometheus</code> which is also marked as the default datasource</p>

</li>
<li>
<p>have added a new Prometheus datasource which you have set as the default</p>

</li>
</ul>
<p>Then continue below.</p>

<ul class="ulist">
<li>
<p>Clone the git repository using</p>

</li>
</ul>
<div class="listing">
<pre>git clone https://github.com/oracle/coherence-operator.git</pre>
</div>

<ul class="ulist">
<li>
<p>Login to Grafana and for each dashboard in the cloned directory <code>&lt;DIR&gt;/helm-charts/coherence-operator/dashboards</code> carry out the
following to upload to Grafana:</p>
<ul class="ulist">
<li>
<p>Highlight the '+' (plus) icons and click on import</p>

</li>
<li>
<p>Click `Upload Json file' button to select a dashboard</p>

</li>
<li>
<p>Leave all the default values and click on <code>Import</code></p>

</li>
</ul>
</li>
</ul>
</div>

<h3 id="different">2. Importing with a different datasource</h3>
<div class="section">
<p>If your Grafana environment has a default datasource that you cannot change or already has a
datasource named <code>Prometheus</code>, follow these steps to import the dashboards:</p>

<ul class="ulist">
<li>
<p>Login to Grafana</p>

</li>
<li>
<p>Create a new datasource named <code>Coherence-Prometheus</code> and point to your Prometheus endpoint</p>

</li>
<li>
<p>Create a temporary directory and copy all the dashboards from the cloned directory
<code>&lt;DIR&gt;/helm-charts/coherence-operator/dashboards</code> to this temporary directory</p>

</li>
<li>
<p>Change to this temporary directory and run the following to update the datasource to <code>Coherence-Prometheus</code> or the dataousrce of your own choice:</p>

</li>
</ul>
<div class="listing">
<pre>for file in *.json
do
    sed -i '' -e 's/"datasource": "Prometheus"/"datasource": "Coherence-Prometheus"/g' \
              -e 's/"datasource": null/"datasource": "Coherence-Prometheus"/g' \
              -e 's/"datasource": "-- Grafana --"/"datasource": "Coherence-Prometheus"/g' $file;
done</pre>
</div>

<ul class="ulist">
<li>
<p>Once you have completed the script, proceed to upload the dashboards as described above.</p>

</li>
</ul>
</div>
</div>
</doc-view>
