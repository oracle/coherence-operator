<doc-view>

<v-layout row wrap>
<v-flex xs12 sm10 lg10>
<v-card class="section-def" v-bind:color="$store.state.currentColor">
<v-card-text class="pa-3">
<v-card class="section-def__card">
<v-card-text>
<dl>
<dt slot=title>Coherence Clusters K8s Resources</dt>
<dd slot="desc"><p>When a <code>CoherenceCluster</code> is deployed into Kubernetes the Coherence Operator will create a number of other resources in Kubernetes.</p>
</dd>
</dl>
</v-card-text>
</v-card>
</v-card-text>
</v-card>
</v-flex>
</v-layout>

<h2 id="_kubernetes_resource_relationships_when_creating_coherence_clusters">Kubernetes Resource Relationships When Creating Coherence Clusters</h2>
<div class="section">
<p>A <code>CoherenceCluster</code> is made up of one or more roles.
In theory a <code>CoherenceCluster</code> could have zero roles but this would not by typical.
A role maps to zero or more <code>Pods</code> that will all share the same specification and hence typically take on the same
business role within an application.</p>

<p>In Kubernetes a Coherence role is represented by a <code>CoherenceRole</code> and a <code>CoherenceInternal</code> crd although it is not expected that
these crds are modified directly, they are purely used to allow roles in the same cluster to be managed as independent
entities by the Coherence Operator.</p>

<p>When a resource of type <code>CoherenceCluster</code> is created in Kubernetes the Coherence Operator will create the other resources.</p>

<p>A <code>Service</code> will be created for every <code>CoherenceCluster</code> to be used for Coherence WKA (cluster membership discovery).
Every <code>Pod</code> that is created as part of this cluster will have a label <code>coherenceCluster=&lt;cluster-name&gt;</code> and the WKA <code>Service</code>
uses this label to identify all of the <code>Pods</code> in the same Coherence cluster. The <code>Pods</code> then use the <code>Service</code> as their WKA address.</p>

<p>A <code>CoherenceRole</code> resource will be created for each role in the <code>CoherenceCluster</code> spec that has a replica count greater than zero.
The name of the <code>CoherenceRole</code> will be in the form <code>&lt;cluster-name&gt;-&lt;role-name&gt;</code></p>

<p>Each <code>CoherenceRole</code> will have a related <code>CoherenceInternal</code> resource. The name of the <code>CoherenceInternal</code> will be the same
as the <code>CoherenceRole</code> resource.</p>

<p>Each <code>CoherenceRole</code> will have a related <code>StatefulSet</code> with corresponding <code>Pods</code> and headless <code>Service</code> required by
the <code>StatefulSet</code>. The name of the <code>StatefulSet</code> will be the same as the <code>CoherenceRole</code></p>

<p>For each port that a role in a <code>CoherenceCluster</code> is configured to expose a corresponding <code>Service</code> will be created for that port.
The name of the <code>Service</code> will be <code>&lt;cluster-name&gt;-&lt;role-name&gt;-&lt;port-name&gt;</code> (although this can be overridden when specifying the port
in the <code>CoherenceCLuster</code> spec for that role and port).</p>

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
