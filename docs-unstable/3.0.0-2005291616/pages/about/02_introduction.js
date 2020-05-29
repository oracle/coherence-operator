<doc-view>

<h2 id="_what_is_the_coherence_operator">What is the Coherence Operator?</h2>
<div class="section">
<p>The Coherence Operator is a <a id="" title="" target="_blank" href="https://kubernetes.io/docs/concepts/extend-kubernetes/operator/">Kubernetes Operator</a> that
is used to manage <a id="" title="" target="_blank" href="https://coherence.java.net">Oracle Coherence</a> clusters in Kubernetes.
The Coherence Operator takes on the tasks of that human Dev Ops resource might carry out when managing Coherence clusters,
such as configuration, installation, safe scaling, management and metrics.</p>

<p>The Coherence Operator is a Go based application built using the <a id="" title="" target="_blank" href="https://github.com/operator-framework/operator-sdk">Operator SDK</a>.
It is distributed as a Docker image and Helm chart for easy installation and configuration.</p>

</div>

<h2 id="_coherence_clusters">Coherence Clusters</h2>
<div class="section">
<p>A Coherence cluster is a number of distributed Java Virtual Machines (JVMs) that communicate to form a single coherent cluster.
In Kubernetes, this concept can be related to a number of Pods that form a single cluster.
In each <code>Pod</code> is a JVM running a Coherence <code>DefaultCacheServer</code>, or a custom application using Coherence.</p>

<p>The operator uses a Kubernetes Custom Resource Definition (CRD) to represent a group of members in a Coherence cluster.
Typically, a deployment would be used to configure one or more members of a specific role in a cluster.
Every field in the <code>Coherence</code> CRD <code>Spec</code> is optional, so a simple cluster can be defined in  yaml as:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: my-cluster <span class="conum" data-value="1" /></markup>

<ul class="colist">
<li data-value="1">In this case the <code>metadata.name</code> field in the <code>Coherence</code> resource yaml will be used as the Coherence cluster name.</li>
</ul>
<p>The operator will use default values for fields that have not been entered, so the above yaml will create
a Coherence deployment using a <code>StatefulSet</code> with a replica count of three, which means that will be three storage
enabled Coherence <code>Pods</code>.</p>

<p>See the <router-link to="/about/04_coherence_spec">Coherence CRD spec</router-link> page for details of all the fields in the CRD.</p>

</div>
</doc-view>
