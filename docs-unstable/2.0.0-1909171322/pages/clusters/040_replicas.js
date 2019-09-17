<doc-view>

<h2 id="_setting_the_replica_count_for_a_role">Setting the Replica Count for a Role</h2>
<div class="section">
<p>The replica count for a role in a <code>CoherenceCluster</code> is set using the <code>replicas</code> field of a role spec.</p>


<h3 id="_implicit_role_replicas">Implicit Role Replicas</h3>
<div class="section">
<p>When using the implicit role configuration the <code>replicas</code> count is set directly in the <code>CoherenceCluster</code> <code>spec</code> section.
For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  replicas: 6  <span class="conum" data-value="1" /></markup>

<ul class="colist">
<li data-value="1">In this case a cluster will be created with a single implicit role names <code>storage</code> with a replica count of six.
This will result in a <code>StatefulSet</code> with six <code>Pods</code>.</li>
</ul>
</div>

<h3 id="_explicit_role_replicas">Explicit Role Replicas</h3>
<div class="section">
<p>When using the explicit role configuration the <code>replicas</code> count is set for each role in the <code>CoherenceCluster</code> <code>spec</code>
<code>roles</code> list.</p>

<p>For example to create cluster with two explicit roles, <code>data</code> and <code>proxy</code>:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  roles:
    - role: data
      replicas: 6  <span class="conum" data-value="1" />
    - role: proxy
      replicas: 3  <span class="conum" data-value="2" /></markup>

<ul class="colist">
<li data-value="1">The <code>data</code> role has a replica count of six</li>
<li data-value="2">The <code>proxy</code> role has a replic count of three</li>
</ul>
</div>

<h3 id="_explicit_role_with_default_replicas">Explicit Role with Default Replicas</h3>
<div class="section">
<p>When using the explicit role configuration a value for <code>replicas</code> count can be set in the <code>CoherenceCluster</code> <code>spec</code>
section that will be used as the default <code>replicas</code> value for any <code>role</code> in the <code>roles</code> list that does not explicitly
specify a value.</p>

<p>For example to create cluster with three explicit roles, <code>data</code> and <code>proxy</code> and <code>web</code>:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  replicas: 6  <span class="conum" data-value="1" />
  roles:
    - role: data
    - role: proxy
    - role: web
      replicas: 3  <span class="conum" data-value="2" /></markup>

<ul class="colist">
<li data-value="1">The default <code>replicas</code> value is set to six. The <code>data</code> and <code>proxy</code> roles do not have a <code>replicas</code> value so will use
this default value and so will each have a <code>StatefulSet</code> with a replica count of six</li>
<li data-value="2">The <code>web</code> role has an explicit <code>replicas</code> value of three so will have three replicas in its <code>StatefulSet</code></li>
</ul>
</div>
</div>
</doc-view>
