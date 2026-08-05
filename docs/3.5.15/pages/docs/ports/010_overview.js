<doc-view>

<h2 id="_overview">Overview</h2>
<div class="section">
<p>Almost every application deployed into a Kubernetes cluster needs to communicate with other processes to provide services
to other processes or consume services to other processes. This is achieved by exposing ports on containers in <code>Pods</code> and
optionally exposing those same ports using <code>Services</code> and ingress.
The <code>Coherence</code> CRD spec makes it simple to add ports to the Coherence container and configure <code>Services</code> to
expose those ports.</p>

<p>Each additional port configured is exposed via its own <code>Service</code>.</p>

<p>If the configuration of <code>Services</code> for ports provided by the <code>Coherence</code> CRD spec is not sufficient or cannot
provide the required <code>Service</code> configuration then it is always possible to just create your own <code>Services</code> in Kubernetes.</p>


<h3 id="_guides_to_adding_and_exposing_ports">Guides to Adding and Exposing Ports</h3>
<div class="section">
<v-layout row wrap class="mb-5">
<v-flex xs12>
<v-container fluid grid-list-md class="pa-0">
<v-layout row wrap class="pillars">
<v-flex xs12 sm4 lg3>
<v-card>
<router-link to="/docs/ports/020_container_ports"><div class="card__link-hover"/>
</router-link>
<v-card-title primary class="headline layout justify-center">
<span style="text-align:center">Adding Ports</span>
</v-card-title>
<v-card-text class="caption">
<p></p>
<p>Adding additional container ports to the Coherence container.</p>
</v-card-text>
</v-card>
</v-flex>
<v-flex xs12 sm4 lg3>
<v-card>
<router-link to="/docs/ports/030_services"><div class="card__link-hover"/>
</router-link>
<v-card-title primary class="headline layout justify-center">
<span style="text-align:center">Expose Ports via Services</span>
</v-card-title>
<v-card-text class="caption">
<p></p>
<p>Configuring Services used to expose ports.</p>
</v-card-text>
</v-card>
</v-flex>
<v-flex xs12 sm4 lg3>
<v-card>
<router-link to="/docs/ports/040_servicemonitors"><div class="card__link-hover"/>
</router-link>
<v-card-title primary class="headline layout justify-center">
<span style="text-align:center">Prometheus ServiceMonitors</span>
</v-card-title>
<v-card-text class="caption">
<p></p>
<p>Adding Prometheus ServiceMonitors to expose ports to be scraped for metrics.</p>
</v-card-text>
</v-card>
</v-flex>
</v-layout>
</v-container>
</v-flex>
</v-layout>
</div>
</div>
</doc-view>
