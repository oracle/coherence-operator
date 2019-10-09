<doc-view>

<h2 id="_what_is_the_coherence_operator">What is the Coherence Operator</h2>
<div class="section">
<p>The Coherence Operator is a <a id="" title="" target="_blank" href="https://kubernetes.io/docs/concepts/extend-kubernetes/operator/">Kubernetes Operator</a> that
is used to manage <a id="" title="" target="_blank" href="https://docs.oracle.com/middleware/12213/coherence/">Oracle Coherence</a> clusters in Kubernetes.
The Coherence Operator takes on the tasks of that human Dev Ops resource might carry out when managing Coherence clusters,
such as configuration, installation, safe scaling, management and metrics.</p>

<p>The Coherence Operator is a Go based application built using the <a id="" title="" target="_blank" href="https://github.com/operator-framework/operator-sdk">Operator SDK</a>.
It is distributed as a Docker image and Helm chart for easy installation and configuration.</p>

</div>

<h2 id="_coherence_clusters">Coherence Clusters</h2>
<div class="section">
<p>Traditionally a Coherence cluster is a number of distributed JVMs that communicate to form a single coherent cluster.
In Kubernetes this concept still applies but can now be though of as a number of Pods that form a single cluster.
Inside each <code>Pod</code> is a JVM running Coherence, or some custom application using Coherence.</p>

<p>The Coherence Operator uses a Kubernetes Custom Resource Definition to represent a Coherence cluster
(and the roles withing it, see below). Every field in the <code>CoherenceCluster</code> crd <code>Spec</code> is optional so a cluster
can be defined by yaml as simple as this:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: my-cluster <span class="conum" data-value="1" /></markup>

<ul class="colist">
<li data-value="1">The <code>metadata.name</code> field in the <code>CoherenceCluster</code> will be used as the Coherence cluster name and would obviously
be unique in a given k8s namespace.</li>
</ul>
<p>The Coherence Operator will use default values for fields that have not been entered, so the above yaml will create
a Coherence cluster made up of a <code>StatefulSet</code> with a replica count of 3, so there will be three storage enabled
Coherence <code>Pods</code>.</p>

</div>

<h2 id="_coherence_roles">Coherence Roles</h2>
<div class="section">
<p>A Coherence cluster can be made up of a number of Pods that perform different roles. All of the Pods in a given role
share the same configuration. At a bare minimum a cluster would have at least one role where Pods are storage enabled.</p>

<p>Each role in a Coherence cluster has a name and configuration. A cluster can have zero or many roles defined in the
<code>CoherenceCluster</code> crd <code>Spec</code>. It is possible to define common configuration shared by all roles to save duplicating
configuration multiple times in the yaml.</p>

<p>The Coherence Operator will create a <code>StatefulSet</code> for each role defined in the <code>CoherenceCluster</code> crd yaml.
This separation allows roles to be managed and scaled independently from each other.</p>

<p>As described above the minimal yaml to define a CoherenceCluster is:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: my-cluster</markup>

<p>Although there are no roles described in this yaml the Coherence Operator will create a default role with the name
<code>storage</code> and give it a replica count of three.</p>

<p>There are two ways to describe the specification of a role in a <code>CoherenceCluster</code> crd depending on whether the cluster
created will have a single role or multiple roles.</p>

<p>The same configuration to create a single role three member cluster as the minimal yaml could be specified more fully
as follows:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: my-cluster
spec:
  role: storage <span class="conum" data-value="1" />
  replicas: 3   <span class="conum" data-value="2" /></markup>

<ul class="colist">
<li data-value="1">The <code>role</code> field specifies the name of the role, in this case <code>stroage</code>.</li>
<li data-value="2">The <code>replicas</code> field defines the number of Pods that will be started for this role, in this case three.</li>
</ul>
<p>If a cluster will have multiple roles they are defined in the <code>spec.roles</code> list; so again the same cluster could be
defined more fully as:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: my-cluster
spec:
  roles:             <span class="conum" data-value="1" />
    - role: storage
      replicas: 3</markup>

<ul class="colist">
<li data-value="1">This time the role is defined in the <code>roles</code> section of the yaml. The <code>roles</code> section is a list of one or more role
specifications.</li>
</ul>
<p>Multiple roles can be defined by adding more roles with distinct names to the <code>roles</code> list; for example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: my-cluster
spec:
  roles:
    - role: storage  <span class="conum" data-value="1" />
      replicas: 3    <span class="conum" data-value="2" />
    - role: web      <span class="conum" data-value="3" />
      replicas: 2    <span class="conum" data-value="4" /></markup>

<ul class="colist">
<li data-value="1">In this case there are two roles defined, the first names <code>storage</code>,</li>
<li data-value="2">with three replicas</li>
<li data-value="3">and the second named <code>web</code></li>
<li data-value="4">with two replicas.</li>
</ul>
<p>This will result in a Coherence cluster with a total of five members.
The Coherence Operator would create two <code>StatefulSets</code>, one for <code>storage</code> with three <code>Pods</code> and one for <code>web</code> with two <code>Pods</code>.</p>

</div>
</doc-view>
