<doc-view>

<h2 id="_examples_overview">Examples Overview</h2>
<div class="section">
<p>There are a number of examples which show you how to build and deploy applications for the Coherence Operator.</p>

<div class="admonition tip">
<p class="admonition-textlabel">Tip</p>
<p ><p><img src="./images/GitHub-Mark-32px.png" alt="GitHub Mark 32px" />
 The complete source code for the examples is in the <a id="" title="" target="_blank" href="https://github.com/oracle/coherence-operator/tree/main/examples/">Coherence Operator GitHub</a> repository.</p>
</p>
</div>
<v-layout row wrap class="mb-5">
<v-flex xs12>
<v-container fluid grid-list-md class="pa-0">
<v-layout row wrap class="pillars">
<v-flex xs12 sm4 lg3>
<v-card>
<router-link to="/examples/015_simple_image/README"><div class="card__link-hover"/>
</router-link>
<v-card-title primary class="headline layout justify-center">
<span style="text-align:center">Simple Coherence Image using JIB</span>
</v-card-title>
<v-card-text class="caption">
<p></p>
<p>Building a simple Coherence server image with <a id="" title="" target="_blank" href="https://github.com/GoogleContainerTools/jib/blob/master/README.md">JIB</a> using Maven or Gradle.</p>
</v-card-text>
</v-card>
</v-flex>
<v-flex xs12 sm4 lg3>
<v-card>
<router-link to="/examples/016_simple_docker_image/README"><div class="card__link-hover"/>
</router-link>
<v-card-title primary class="headline layout justify-center">
<span style="text-align:center">Simple Coherence Image using a Dockerfile</span>
</v-card-title>
<v-card-text class="caption">
<p></p>
<p>Building a simple Coherence image with a Dockerfile, that works out of the box with the Operator.</p>
</v-card-text>
</v-card>
</v-flex>
<v-flex xs12 sm4 lg3>
<v-card>
<router-link to="/examples/020_hello_world/README"><div class="card__link-hover"/>
</router-link>
<v-card-title primary class="headline layout justify-center">
<span style="text-align:center">Hello World</span>
</v-card-title>
<v-card-text class="caption">
<p></p>
<p>Deploying the most basic Coherence cluster using the Operator.</p>
</v-card-text>
</v-card>
</v-flex>
</v-layout>
</v-container>
</v-flex>
</v-layout>
<v-layout row wrap class="mb-5">
<v-flex xs12>
<v-container fluid grid-list-md class="pa-0">
<v-layout row wrap class="pillars">
<v-flex xs12 sm4 lg3>
<v-card>
<router-link to="/examples/025_extend_client/README"><div class="card__link-hover"/>
</router-link>
<v-card-title primary class="headline layout justify-center">
<span style="text-align:center">Coherence*Extend Clients</span>
</v-card-title>
<v-card-text class="caption">
<p></p>
<p>An example demonstrating various ways to configure and use Coherence*Extend with Kubernetes.</p>
</v-card-text>
</v-card>
</v-flex>
<v-flex xs12 sm4 lg3>
<v-card>
<router-link to="#examples/020_deployment/README.adoc" @click.native="this.scrollFix('#examples/020_deployment/README.adoc')"><div class="card__link-hover"/>
</router-link>
<v-card-title primary class="headline layout justify-center">
<span style="text-align:center">Deployment</span>
</v-card-title>
<v-card-text class="caption">
<p></p>
<p>This example shows how to deploy Coherence applications using the Coherence Operator.</p>
</v-card-text>
</v-card>
</v-flex>
<v-flex xs12 sm4 lg3>
<v-card>
<router-link to="/examples/090_tls/README"><div class="card__link-hover"/>
</router-link>
<v-card-title primary class="headline layout justify-center">
<span style="text-align:center">TLS</span>
</v-card-title>
<v-card-text class="caption">
<p></p>
<p>Securing Coherence clusters using TLS.</p>
</v-card-text>
</v-card>
</v-flex>
<v-flex xs12 sm4 lg3>
<v-card>
<router-link to="/examples/095_network_policies/README"><div class="card__link-hover"/>
</router-link>
<v-card-title primary class="headline layout justify-center">
<span style="text-align:center">Network Policies</span>
</v-card-title>
<v-card-text class="caption">
<p></p>
<p>An example covering the use of Kubernetes <code>NetworkPolicy</code> rules with the Operator and Coherence clusters.</p>
</v-card-text>
</v-card>
</v-flex>
<v-flex xs12 sm4 lg3>
<v-card>
<router-link to="/examples/100_federation/README"><div class="card__link-hover"/>
</router-link>
<v-card-title primary class="headline layout justify-center">
<span style="text-align:center">Federation</span>
</v-card-title>
<v-card-text class="caption">
<p></p>
<p>This is a simple Coherence federation example. The federation feature requires Coherence Grid Edition.</p>
</v-card-text>
</v-card>
</v-flex>
<v-flex xs12 sm4 lg3>
<v-card>
<router-link to="/examples/200_autoscaler/README"><div class="card__link-hover"/>
</router-link>
<v-card-title primary class="headline layout justify-center">
<span style="text-align:center">Autoscaling</span>
</v-card-title>
<v-card-text class="caption">
<p></p>
<p>Scaling Coherence clusters using the horizontal Pod Autoscaler.</p>
</v-card-text>
</v-card>
</v-flex>
<v-flex xs12 sm4 lg3>
<v-card>
<router-link to="/examples/300_helm/README"><div class="card__link-hover"/>
</router-link>
<v-card-title primary class="headline layout justify-center">
<span style="text-align:center">Helm</span>
</v-card-title>
<v-card-text class="caption">
<p></p>
<p>Manage Coherence resources using Helm.</p>
</v-card-text>
</v-card>
</v-flex>
<v-flex xs12 sm4 lg3>
<v-card>
<router-link to="/examples/400_Istio/README"><div class="card__link-hover"/>
</router-link>
<v-card-title primary class="headline layout justify-center">
<span style="text-align:center">Istio</span>
</v-card-title>
<v-card-text class="caption">
<p></p>
<p>Istio Support</p>
</v-card-text>
</v-card>
</v-flex>
</v-layout>
</v-container>
</v-flex>
</v-layout>
<v-layout row wrap class="mb-5">
<v-flex xs12>
<v-container fluid grid-list-md class="pa-0">
<v-layout row wrap class="pillars">
<v-flex xs12 sm4 lg3>
<v-card>
<router-link to="/examples/900_demo/README"><div class="card__link-hover"/>
</router-link>
<v-card-title primary class="headline layout justify-center">
<span style="text-align:center">Coherence Demo App</span>
</v-card-title>
<v-card-text class="caption">
<p></p>
<p>Deploying the Coherence demo application.</p>
</v-card-text>
</v-card>
</v-flex>
</v-layout>
</v-container>
</v-flex>
</v-layout>
</div>
</doc-view>
