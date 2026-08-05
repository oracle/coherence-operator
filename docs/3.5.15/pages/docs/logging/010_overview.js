<doc-view>

<h2 id="_overview">Overview</h2>
<div class="section">
<p>In a container environment like Kubernetes, or any cloud, it is often a requirement to centralize log files
to allow easier analysis and debugging. There are many ways to do this, including collecting container logs,
parsing and shipping log files with something like Fluentd, or using a specialized log appender specific to
your logging framework.</p>

<p>The Coherence Operator does not proscribe any particular method of log capture. The <code>Coherence</code> CRD is flexible
enough to allow any method of log capture that an application or specific cloud environment requires.
This could be as simple as adding JVM arguments to configure the Java logger, or it could be injecting a whole
side-car container to run something like Fluentd. Different approaches have their own pros and cons that need
to be weighed up on a case by case basis.</p>


<h3 id="_logging_guides">Logging Guides</h3>
<div class="section">
<p>The use of Elasticsearch, Fluentd and Kibana is a common approach. For this reason the Coherence Operator
has a set of Kibana dashboards that support the common Coherence logging format.
The logging guides below show one approach to shipping Coherence logs to Elasticsearch and importing the Coherence
dashboards into Kibana.
If this approach does not meet your needs you are obviously free to configure an alternative.</p>

<v-layout row wrap class="mb-5">
<v-flex xs12>
<v-container fluid grid-list-md class="pa-0">
<v-layout row wrap class="pillars">
<v-flex xs12 sm4 lg3>
<v-card>
<router-link to="/docs/logging/020_logging"><div class="card__link-hover"/>
</router-link>
<v-card-title primary class="headline layout justify-center">
<span style="text-align:center">Enabling Log Capture</span>
</v-card-title>
<v-card-text class="caption">
<p></p>
<p>Capturing and viewing Coherence cluster Logs in Elasticsearch using a Fluentd side-car.</p>
</v-card-text>
</v-card>
</v-flex>
<v-flex xs12 sm4 lg3>
<v-card>
<router-link to="/docs/logging/030_kibana"><div class="card__link-hover"/>
</router-link>
<v-card-title primary class="headline layout justify-center">
<span style="text-align:center">Kibana Dashboards</span>
</v-card-title>
<v-card-text class="caption">
<p></p>
<p>Importing and using the Kibana Dashboards available.</p>
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
