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

<h2 id="_create_coherencecluster_resources">Create CoherenceCluster Resources</h2>
<div class="section">
<p>The Coherence Operator uses a Kubernetes <code>CustomResourceDefinition</code> named <code>CoherenceCluster</code> to define the <code>spec</code> for a
Coherence cluster.
All of the fields of the <code>CoherenceCluster</code> crd are optional so a Coherence cluster can be created with yaml as
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
so the cluster is inaccessible. It could be accessed by other <code>Pods</code> in the same</p>

</div>

<h2 id="_coherence_roles">Coherence Roles</h2>
<div class="section">
<p>A role is what is actually configured in the <code>CoherenceCluster</code> spec. In a traditional Coherence application that may have
had a number of storage enabled members and a number of storage disable Coherence*Extend proxy members this cluster would
have effectively had two roles, "storage" and "proxy".
Some clusters may simply have just a storage role and some complex Coherence applications and clusters may have many roles
and even different roles storage enabled for different caches/services within the same cluster.</p>

</div>
</doc-view>
