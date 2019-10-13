<doc-view>

<v-layout row wrap>
<v-flex xs12 sm10 lg10>
<v-card class="section-def" v-bind:color="$store.state.currentColor">
<v-card-text class="pa-3">
<v-card class="section-def__card">
<v-card-text>
<dl>
<dt slot=title>Coherence Operator Installation</dt>
<dd slot="desc"><p>The Coherence Operator is available as a Docker image <code>oracle/coherence-operator:2.0.0-1910131023</code> that can
easily be installed into a Kubernetes cluster.</p>
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
clusters will be installed into.
see <router-link to="/about/04_obtain_coherence_images">Obtain Coherence Images</router-link></p>

</div>
</div>

<h2 id="_installing_with_helm">Installing With Helm</h2>
<div class="section">
<p>The simplest way to install the Coherence Operator is to use the Helm chart.
This will ensure that all of the correct resources are created in Kubernetes.</p>


<h3 id="_add_the_coherence_helm_repository">Add the Coherence Helm Repository</h3>
<div class="section">
<p>Add the <code>coherence</code> helm repository using the following commands:</p>

<markup
lang="bash"

>helm repo add coherence https://oracle.github.io/coherence-operator/charts

helm repo update</markup>

</div>

<h3 id="_install_the_coherence_operator_helm_chart">Install the Coherence Operator Helm chart</h3>
<div class="section">
<p>Once the Coherence Helm repo is configured the Coherence Operator can be installed using a normal Helm install command:</p>

<markup
lang="bash"

>helm install  \
    --namespace &lt;namespace&gt; \
    --name coherence-operator \
    coherence/coherence-operator</markup>

<p>where <code>&lt;namespace&gt;</code> is the namespace that the Coherence Operator will be installed into and the namespace where it will
manage <code>CoherenceClusters</code></p>


<h4 id="_uninstall_the_coherence_operator_helm_chart">Uninstall the Coherence Operator Helm chart</h4>
<div class="section">
<p>To uninstall the operator:</p>

<markup
lang="bash"

>helm delete --purge coherence-operator</markup>

</div>
</div>
</div>
</doc-view>
