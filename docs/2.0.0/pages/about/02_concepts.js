<doc-view>

<h2 id="_what_is_the_coherence_operator">What is the Coherence Operator?</h2>
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
<p>A Coherence cluster is a number of distributed Java Virtual Machines (JVMs) that communicate to form a single coherent cluster.
In Kubernetes, this concept can be related to a number of Pods that form a single cluster.
In each <code>Pod</code> is a JVM running a Coherence <code>DefaultCacheServer</code>, or a custom application using Coherence.</p>

<p>The operator uses a Kubernetes Custom Resource Definition (CRD) to represent a Coherence cluster
and the roles within it. Every field in the <code>CoherenceCluster</code> CRD <code>spec</code> is optional so a simple cluster
can be defined in  yaml as:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: my-cluster <span class="conum" data-value="1" /></markup>

<ul class="colist">
<li data-value="1">The <code>metadata.name</code> field in the <code>CoherenceCluster</code> yaml will be used as the Coherence cluster name and must
be unique in a given Kubernetes namespace.</li>
</ul>
<p>The operator will use default values for fields that have not been entered, so the above yaml will create
a Coherence cluster using a <code>StatefulSet</code> with a replica count of three, which means that will be three storage
enabled Coherence <code>Pods</code>.</p>

</div>

<h2 id="_coherence_roles">Coherence Roles</h2>
<div class="section">
<p>A Coherence cluster can be made up of a number of Pods that perform different roles. All of the Pods in a given role
share the same configuration. A cluster usually has at least one role where <code>Pods</code> are storage enabled.</p>

<p>Each role in a Coherence cluster has a name and configuration. A cluster can have zero or many roles defined in the
<code>CoherenceCluster</code> CRD <code>Spec</code>. You can define common configuration shared by all roles to save duplicating
configuration multiple times in the yaml.</p>

<p>The Coherence Operator will create a <code>StatefulSet</code> for each role defined in the <code>CoherenceCluster</code> CRD yaml.
This separation allows roles to be managed and scaled independently from each other. All of the <code>Pods</code> in the
different <code>StatefulSets</code> will form a single Coherence cluster.</p>

<p>There are two ways to describe the specification of a role in a <code>CoherenceCluster</code> CRD depending on whether the cluster
has a single implied role or has one or more explicit roles.</p>


<h3 id="_a_single_implied_role">A Single Implied Role</h3>
<div class="section">
<p>The operator implies that a single role is required when the <code>roles</code> list in the <code>CoherenceCluster</code> CRD yaml is either
empty or missing. For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: my-cluster
spec:
  role: data
  replicas: 6</markup>

<p>The yaml above does not include any roles defined in the <code>roles</code> list of the <code>spec</code> section. When all of the role
configuration is in fields directly in the <code>spec</code> section like this the operator implies that a single role is required
and will use the values defined in the <code>spec</code> section to create a single <code>StatefulSet</code>.</p>

</div>

<h3 id="_a_single_explicit_role">A Single Explicit Role</h3>
<div class="section">
<p>Roles can be defined explicitly by adding the configuration of each role to the <code>roles</code> list in the <code>CoherenceCluster</code>
CRD <code>spec</code> section. For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: my-cluster
spec:
  roles:
    - role: data
      replicas: 6</markup>

<p>In the example above there is explicitly one role defined in the <code>roles</code> list.</p>

</div>

<h3 id="_multiple_explicit_role">Multiple Explicit Role</h3>
<div class="section">
<p>To define a Coherence cluster with multiple roles each role is configured as a separate entry in the <code>roles</code> list.</p>

<p>For example, if a cluster requires two roles, one named <code>storage</code> and another named <code>web</code> the configuration may
look like this:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: my-cluster
spec:
  roles:
    - role: storage  <span class="conum" data-value="1" />
      replicas: 3
    - role: web      <span class="conum" data-value="2" />
      replicas: 2</markup>

<ul class="colist">
<li data-value="1">The <code>storage</code> role is explicitly defined in the <code>roles</code> list</li>
<li data-value="2">The <code>web</code> role is explicitly defined in the <code>roles</code> list</li>
</ul>
<p>This will result in a Coherence cluster made up of two <code>StatefulSets</code>. The <code>storage</code> role will have a <code>StatefulSet</code> with
three <code>Pods</code> and the <code>web</code> role will have a <code>StatefulSet</code> with two <code>Pods</code>. The Coherence cluster will have a total of
five <code>Pods</code>.</p>

</div>

<h3 id="_explicit_roles_with_default_values">Explicit Roles with Default Values</h3>
<div class="section">
<p>When defining explicit roles in the <code>roles</code> list and field added directly to the <code>CoherenceCluster</code> <code>spec</code> section
becomes a default value that is applied to all of the roles in the <code>roles</code> list unless the value is overridden in
the configuration for a specific role. This allows common configuration shared by multiple roles to be maintained in
a single place instead of being duplicated for every role.</p>

<p>For example, if a cluster requires three roles, one named <code>storage</code> with a <code>5g</code> JVM heap, one named <code>proxy</code> with a <code>5g</code>
JVM heap and another named <code>web</code> with a <code>1g</code> JVM heap the configuration may look like this:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: my-cluster
spec:
  jvm:
    memory:
      heapSize: 5g      <span class="conum" data-value="1" />
  roles:
    - role: storage
      replicas: 6
    - role: storage
      replicas: 3
    - role: web
      replicas: 2
      jvm:
        memory:
          heapSize: 1g  <span class="conum" data-value="2" /></markup>

<ul class="colist">
<li data-value="1">The <code>jvm.memory.heapSize</code> value of <code>5g</code> is added directly under the <code>spec</code> section so this value will apply to
all roles meaning all roles will have the JVM options <code>-Xms5g -Xmx5g</code> unless overridden. In this case the <code>storage</code> and
the <code>proxy</code> roles do not set the <code>jvm.memory.heapSize</code> field so they will have a <code>5g</code> JVM heap.</li>
<li data-value="2">The <code>web</code> role overrides the <code>jvm.memory.heapSize</code> field with a value of <code>1g</code> so the JVMs in the <code>web</code> role will
have the JVM options <code>-Xms1g -Xmx1g</code></li>
</ul>
<div class="admonition note">
<p class="admonition-inline">When using default values some default values are overridden by values in a role and sometimes the default and
role values are merged. When the field is a single intrinsic value, for example a number or a string the role value
overrides the default. Where the field is an array/slice or a map it may be merged.
The <router-link to="/clusters/010_introduction">CoherenceCluster CRD section</router-link> documents how fields are overridden or merged.</p>
</div>
</div>
</div>
</doc-view>
