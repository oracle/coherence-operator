<doc-view>

<h2 id="_set_coherence_cluster_name">Set Coherence Cluster Name</h2>
<div class="section">
<p>The name of the Coherence cluster that a <code>Coherence</code> resource is part of can be set with the <code>cluster</code> field
in the <code>Coherence.Spec</code>. The cluster name is used to set the <code>coherence.cluster</code> system property in the JVM in the Coherence container.</p>


<h3 id="_default_cluster_name">Default Cluster Name</h3>
<div class="section">
<p>The default Coherence cluster name, used when the <code>cluster</code> field is empty, will be the same as the name of the <code>Coherence</code> resource, for example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: test <span class="conum" data-value="1" /></markup>

<ul class="colist">
<li data-value="1">The name of this <code>Coherence</code> resource is <code>test</code>, which will also be used as the Coherence cluster name, effectively passing <code>-Dcoherence.cluster=test</code> to the JVM in the Coherence container.</li>
</ul>
</div>

<h3 id="_specify_a_cluster_name">Specify a Cluster Name</h3>
<div class="section">
<p>In a use case where multiple <code>Coherence</code> resources will be created to form a single Coherence cluster, the <code>cluster</code>
field in all the <code>Coherence</code> resources needs to be set to the same value.</p>

<markup
lang="yaml"
title="cluster.yaml"
>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: storage
spec:
  cluster: test-cluster
---
apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: front-end
spec:
  cluster: test-cluster</markup>

<p>The yaml above contains two <code>Coherence</code> resources, one named <code>storage</code> and one named <code>front-end</code>.
Both of these <code>Coherence</code> resources have the same value for the <code>cluster</code> field, <code>test-cluster</code>,
so the Pods in both deployments will form a single Coherence cluster named <code>test</code>.</p>

</div>
</div>
</doc-view>
