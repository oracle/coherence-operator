<doc-view>

<v-layout row wrap>
<v-flex xs12 sm10 lg10>
<v-card class="section-def" v-bind:color="$store.state.currentColor">
<v-card-text class="pa-3">
<v-card class="section-def__card">
<v-card-text>
<dl>
<dt slot=title>Installing With Helm</dt>
<dd slot="desc"><p>The simplest way to install the Coherence Operator is to use the Helm chart.
This will ensure that all of the correct resources are created in Kubernetes.</p>
</dd>
</dl>
</v-card-text>
</v-card>
</v-card-text>
</v-card>
</v-flex>
</v-layout>

<h2 id="_add_the_coherence_helm_repository">Add the Coherence Helm Repository</h2>
<div class="section">
<p>Add the <code>coherence</code> helm repository using the following commands:</p>

<markup
lang="bash"

>$ helm repo add coherence https://oracle.github.io/coherence-operator/charts

$ helm repo update</markup>

</div>

<h2 id="_install_the_coherence_operator_helm_chart">Install the Coherence Operator Helm chart</h2>
<div class="section">
<p>Once the Coherence Helm repo is configured the Coherence Operator can be installed using a normal Helm install command:</p>

<markup
lang="bash"

>$ helm install  \
    --namespace &lt;namespace&gt; \
    --name coherence-operator \
    coherence/coherence-operator</markup>

<p>where <code>&lt;namespace&gt;</code> is the namespace that the Coherence Operator will be installed into and the namespace where it will
manage <code>CoherenceClusters</code></p>


<h3 id="_uninstall_the_coherence_operator_helm_chart">Uninstall the Coherence Operator Helm chart</h3>
<div class="section">
<p>To uninstall the operator:</p>

<markup
lang="bash"

>$ helm delete --purge coherence-operator</markup>

</div>
</div>
</doc-view>
