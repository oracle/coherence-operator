<doc-view>

<v-layout row wrap>
<v-flex xs12 sm10 lg10>
<v-card class="section-def" v-bind:color="$store.state.currentColor">
<v-card-text class="pa-3">
<v-card class="section-def__card">
<v-card-text>
<dl>
<dt slot=title>Create Coherence Clusters</dt>
<dd slot="desc"><p>Creating a Coherence cluster using the Coherence Operator is as simple as creating any other Kubernetes resource.</p>
</dd>
</dl>
</v-card-text>
</v-card>
</v-card-text>
</v-card>
</v-flex>
</v-layout>

<h2 id="_create_coherence_clusters">Create Coherence Clusters</h2>
<div class="section">
<p>The Coherence Operator uses a Kubernetes <code>CustomResourceDefinition</code> to define the spec for a Coherence cluster.</p>

<p>All of the fields of the <code>CoherenceCluster</code> crd are optional so a Coherence cluster can be created with yaml as
simple as the following:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: my-cluster  <span class="conum" data-value="1" /></markup>

<ul class="colist">
<li data-value="1">The <code>metadata.name</code> field will be used as the Coherence cluster name.</li>
</ul>
<p>The yaml above will create a Coherence cluster with three storage enabled members.
There is not much that can actually be achived with this cluster because no ports are exposed outside of Kubernetes
so the cluster is inaccessible.</p>

</div>

<h2 id="_coherence_roles">Coherence Roles</h2>
<div class="section">

</div>
</doc-view>
