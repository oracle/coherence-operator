<doc-view>

<v-layout row wrap>
<v-flex xs12 sm10 lg10>
<v-card class="section-def" v-bind:color="$store.state.currentColor">
<v-card-text class="pa-3">
<v-card class="section-def__card">
<v-card-text>
<dl>
<dt slot=title>Overview</dt>
<dd slot="desc"><p>This section includes links to a number of examples which uses the Coherence Operator.</p>
</dd>
</dl>
</v-card-text>
</v-card>
</v-card-text>
</v-card>
</v-flex>
</v-layout>

<h2 id="_examples_overview">Examples Overview</h2>
<div class="section">
<p>There are a number of examples which show you how to build and deploy applications for the Coherence Operator.</p>


<h3 id="_1_deployment_example">1. Deployment Example</h3>
<div class="section">
<p>This example showcases how to deploy Coherence applications using the Coherence Operator.</p>

<p>The following scenarios are covered:</p>

<ol style="margin-left: 15px;">
<li>
Installing the Coherence Operator

</li>
<li>
Installing a Coherence cluster

</li>
<li>
Deploying a Proxy tier

</li>
<li>
Deploying an storage-disabled application

</li>
<li>
Enabling Active Persistence

</li>
</ol>
<p>After the initial install of the Coherence cluster, the following examples build on the previous ones by issuing a kubectl apply to modify the install adding additional roles.</p>

<p>Please see <a id="" title="" target="_blank" href="https://github.com/oracle/coherence-operator/tree/master/examples/deployment">GitHub</a> for full instructions.</p>

</div>

<h3 id="_2_coherence_demo">2. Coherence Demo</h3>
<div class="section">
<p>The Coherence Demonstration application is an application which demonstrates various Coherence
related features such include Persistence, Federation and Lambda support.  This demonstration
can run stand alone but can also be installed on the Coherence Operator.</p>

<p>When installed using the Coherence Operator, the setup includes two Coherence Clusters, in the same Kubernetes cluster,
which are configured with Active/Active Federation.</p>



<v-card>
<v-card-text class="overflow-y-hidden" style="text-align:center">
<img src="./images/coherence-demo.png" alt="Service Details"width="950" />
</v-card-text>
</v-card>

<p>Please see <a id="" title="" target="_blank" href="https://github.com/coherence-community/coherence-demo">The Coherence Demo GitHub project</a> for full instructions.</p>

</div>
</div>
</doc-view>
