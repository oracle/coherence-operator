<doc-view>

<h2 id="_what_is_the_coherence_operator">What is the Coherence Operator?</h2>
<div class="section">
<p>The Coherence Operator is a <a id="" title="" target="_blank" href="https://kubernetes.io/docs/concepts/extend-kubernetes/operator/">Kubernetes Operator</a> that
is used to manage <a id="" title="" target="_blank" href="https://oracle.github.io/coherence">Oracle Coherence</a> clusters in Kubernetes.
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

<p>In the above example no <code>spec.image</code> field has been set, so the Operator will use a publicly pullable Coherence CE
image as its default. These images are meant for demos, POCs and experimentation, but for a production application you
should build your own image.</p>

</div>

<h2 id="_using_commercial_coherence_versions">Using Commercial Coherence Versions</h2>
<div class="section">
<div class="admonition note">
<p class="admonition-inline">Whilst the Coherence CE version can be freely deployed anywhere, if your application image uses a commercial
version of Oracle Coherence then you are responsible for making sure your deployment has been properly licensed.</p>
</div>
<p>Oracle&#8217;s current policy is that a license will be required for each Kubernetes Node that images are to be pulled to.
While an image exists on a node it is effectively the same as having installed the software on that node.</p>

<p>One way to ensure that the Pods of a Coherence deployment only get scheduled onto nodes that meet the
license requirement is to configure Pod scheduling, for example a node selector. Node selectors, and other scheduling,
is simple to configure in the <code>Coherence</code> CRD, see the <router-link to="/other/090_pod_scheduling">scheduling documentation</router-link></p>

<p>For example, if a commercial Coherence license exists such that a sub-set of nodes in a Kubernetes cluster
have been covered by the license then those nodes could all be given a label, e.g. <code>coherenceLicense=true</code></p>

<p>When creating a <code>Coherence</code> deployment specify a node selector to match the label:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: test
spec:
  image: my-app:1.0.0         <span class="conum" data-value="1" />
  nodeSelector:
    coherenceLicense: 'true'  <span class="conum" data-value="2" /></markup>

<ul class="colist">
<li data-value="1">The <code>my-app:1.0.0</code> image contains a commercial Coherence version.</li>
<li data-value="2">The <code>nodeSelector</code> will ensure Pods only get scheduled to nodes with the <code>coherenceLicense=true</code> label.</li>
</ul>
<p>There are other ways to configure Pod scheduling supported by the Coherence Operator (such as taints and tolerations)
and there are alternative ways to restrict nodes that Pods can be schedule to, for example a namespace in kubernetes
can be restricted to a sub-set of the cluster&#8217;s nodes. Using a node selector as described above is probably the
simplest approach.</p>

</div>
</doc-view>
