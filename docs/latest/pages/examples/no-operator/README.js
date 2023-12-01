<doc-view>

<h2 id="_coherence_in_kubernetes_without_the_operator">Coherence in Kubernetes Without the Operator</h2>
<div class="section">
<p>Although this project is all about the Coherence Kubernetes Operator, there are occasions where using an Operator is not possible.
For example, some corporate or cloud security policies ban the use of CRDs, or have very restrictive RBAC policies that ultimately make it impossible to run Operators that uses their own CRDs or require cluster roles (or even just namespace roles).
These example shows how to run a Coherence clusters in Kubernetes manually.
Obviously the features of the Operator such as safe scaling, safe rolling upgrades, etc. will not be available.</p>

<div class="admonition note">
<p class="admonition-textlabel">Note</p>
<p ><p>We really recommend that you try and use the Coherence Operator for managing Coherence clusters in Kubernetes.
It is possible to run the Operator with fewer RBAC permissions, for example without <code>ClusterRoles</code> and only using <code>Roles</code> restricted to a single namespace. The Operator can also run without installing its web-hooks. Ultimately though it requires the CRD to be installed, which could be done manually instead of allowing the Operator to install it.
If you really cannot change the minds of those dictating policies that mean you cannot use the Operator then these examples may be useful.</p>
</p>
</div>
<div class="admonition tip">
<p class="admonition-textlabel">Tip</p>
<p ><p><img src="./images/GitHub-Mark-32px.png" alt="GitHub Mark 32px" />
 The complete source code for the examples is in the <a id="" title="" target="_blank" href="https://github.com/oracle/coherence-operator/tree/main/examples/no-operator/">Coherence Operator GitHub</a> repository.</p>
</p>
</div>

<h3 id="_prerequisites">Prerequisites</h3>
<div class="section">
<p>There are some common prerequisites used by all the examples.</p>

<ul class="ulist">
<li>
<p><strong>The Server Image</strong></p>

</li>
</ul>
<p>These examples use the image built in the <router-link to="/examples/015_simple_image/README">Build a Coherence Server Image</router-link> example.
The image is nothing more than a cache configuration file that has an Extend proxy along with Coherence metrics and management over REST.
We will use this image in the various examples we cover here. When we run the image it will start a simple storage enabled Coherence server.</p>

<ul class="ulist">
<li>
<p><strong>The Test Client</strong>
In the <router-link to="/examples/no-operator/test-client/README"><code>test-client/</code></router-link> directory is a simple Maven project that we will use to run a simple Extend client.</p>

</li>
<li>
<p><strong>Network Policies</strong>
When running in Kubernetes cluster where <code>NetworkPolicy</code> rules are applied there are certain ingress and egress policies required to allow Coherence to work. These are covered in the <router-link to="/examples/095_network_policies/README">Network Policies Example</router-link></p>

</li>
</ul>
</div>
</div>

<h2 id="_the_examples">The Examples</h2>
<div class="section">
<v-layout row wrap class="mb-5">
<v-flex xs12>
<v-container fluid grid-list-md class="pa-0">
<v-layout row wrap class="pillars">
<v-flex xs12 sm4 lg3>
<v-card>
<router-link to="/examples/no-operator/01_simple_server/README"><div class="card__link-hover"/>
</router-link>
<v-card-title primary class="headline layout justify-center">
<span style="text-align:center">Simple Server</span>
</v-card-title>
<v-card-text class="caption">
<p></p>
<p>Run a simple Coherence storage enabled cluster as a StatefulSet and connect an Extend client to it.</p>
</v-card-text>
</v-card>
</v-flex>
<v-flex xs12 sm4 lg3>
<v-card>
<router-link to="/examples/no-operator/02_metrics/README"><div class="card__link-hover"/>
</router-link>
<v-card-title primary class="headline layout justify-center">
<span style="text-align:center">Simple Server with Metrics</span>
</v-card-title>
<v-card-text class="caption">
<p></p>
<p>Expands the simple storage enabled server to expose metrics that can be scraped by Prometheus.</p>
</v-card-text>
</v-card>
</v-flex>
<v-flex xs12 sm4 lg3>
<v-card>
<router-link to="/examples/no-operator/03_extend_tls/README"><div class="card__link-hover"/>
</router-link>
<v-card-title primary class="headline layout justify-center">
<span style="text-align:center">Securing Extend with TLS</span>
</v-card-title>
<v-card-text class="caption">
<p></p>
<p>Expands the simple storage enabled server to secure Extend using TLS.</p>
</v-card-text>
</v-card>
</v-flex>
<v-flex xs12 sm4 lg3>
<v-card>
<router-link to="/examples/no-operator/04_istio/README"><div class="card__link-hover"/>
</router-link>
<v-card-title primary class="headline layout justify-center">
<span style="text-align:center">Running Coherence with Istio</span>
</v-card-title>
<v-card-text class="caption">
<p></p>
<p>Expands the simple storage enabled server to secure Extend using TLS.</p>
</v-card-text>
</v-card>
</v-flex>
</v-layout>
</v-container>
</v-flex>
</v-layout>
</div>
</doc-view>
