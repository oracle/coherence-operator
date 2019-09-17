<doc-view>

<v-layout row wrap>
<v-flex xs12 sm10 lg10>
<v-card class="section-def" v-bind:color="$store.state.currentColor">
<v-card-text class="pa-3">
<v-card class="section-def__card">
<v-card-text>
<dl>
<dt slot=title>Define Coherence Roles</dt>
<dd slot="desc"><p>A <code>CoherenceCluster</code> is made up of one or more roles defined in its <code>spec</code>.</p>
</dd>
</dl>
</v-card-text>
</v-card>
</v-card-text>
</v-card>
</v-flex>
</v-layout>

<h2 id="_define_coherence_roles">Define Coherence Roles</h2>
<div class="section">
<p>A role is what is actually configured in the <code>CoherenceCluster</code> spec. In a traditional Coherence application that may have
had a number of storage enabled members and a number of storage disable Coherence*Extend proxy members this cluster would
have effectively had two roles, "storage" and "proxy".
Some clusters may simply have just a storage role and some complex Coherence applications and clusters may have many roles
and even different roles storage enabled for different caches/services within the same cluster.</p>

<p>The Coherence Operator uses an internal crd named <code>CoherenceRole</code> to represent a role in a Coherence Cluster.
A <code>CoherenceRole</code> would not typically be modified directly outside of a handful of specialized operations, such as scaling.
Any modification to a role would normally be done by modifying that role in the corresponding <code>CoherenceCluster</code> and leaving
the COherence Operator to update the <code>CoherenceRole</code>.</p>

</div>

<h2 id="_defining_a_coherence_role">Defining a Coherence Role</h2>
<div class="section">
<p>All <code>CoherenceCluster</code> resources will have at lest one role defined. This could be the implicit default role or it could
be one more explicit roles.</p>


<h3 id="_implicit_default_role">Implicit Default Role</h3>
<div class="section">
<p>As mentioned previously, all of the fields in a <code>CoherenceCluster</code> <code>spec</code> are optional meaning that the yaml below is
perfectly valid.</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster</markup>

<p>This yaml will create an implicit single role with a default role name of <code>storage</code> and a default replica count of three.</p>

<p>The implicit role can be modified by specifying role related fields in the <code>CoherenceCluster</code> <code>spec</code>.
The role name and replica count of the implicit role can be overridden using the corresponding fields</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  role: data   <span class="conum" data-value="1" />
  replicas: 6  <span class="conum" data-value="2" /></markup>

<ul class="colist">
<li data-value="1">The role name is set with the <code>role</code> field, in this case the role name of the implicit role is now <code>data</code></li>
<li data-value="2">The replica count is set using the <code>replicas</code> field, in this case the implicit role will now have six replicas.</li>
</ul>
<p>Other role fields can also be used, for example, to set the cache configuration file use by the implicit roles:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  role: data
  replicas: 6
  cacheConfig: test-config.xml  <span class="conum" data-value="1" /></markup>

<ul class="colist">
<li data-value="1">The <code>cacheConfig</code> field is used to set the cache configuration to <code>test-config.xml</code>.</li>
</ul>
</div>

<h3 id="_explicit_roles">Explicit Roles</h3>
<div class="section">
<p>It is possible to also create roles explicitly in the <code>roles</code> list of the <code>CoherenceCluster</code> <code>spec</code>.
If creating a Coherence cluster with more than one role then all roles must be defined in the <code>roles</code> list.
If creating a Coherence cluster with a single role it is optional whether the specification of that role is put into
the <code>CoherenceCluster``spec</code> directly as shown above or whether the single role is added to the <code>roles</code> list.</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  roles:           <span class="conum" data-value="1" />
    - role: data   <span class="conum" data-value="2" />
      replicas: 6  <span class="conum" data-value="3" /></markup>

<ul class="colist">
<li data-value="1">The yaml above defines a single explicit role in the <code>roles</code> list</li>
<li data-value="2">When defining explict roles the role name is mandatory. The role name is set with the <code>role</code> field, in this case
the role name of the role is <code>data</code></li>
<li data-value="3">The replica count is set using the <code>replicas</code> field, in this case the role will have six replicas.</li>
</ul>
<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  roles:           <span class="conum" data-value="1" />
    - role: data   <span class="conum" data-value="2" />
      replicas: 6  <span class="conum" data-value="3" />
    - role: proxy  <span class="conum" data-value="4" />
      replicas: 3  <span class="conum" data-value="5" /></markup>

<ul class="colist">
<li data-value="1">The yaml above defines a two explicit roles in the <code>roles</code> list</li>
<li data-value="2">The first role has a role name of <code>data</code></li>
<li data-value="3">and a replica count of six.</li>
<li data-value="4">The second role has a role name of <code>proxy</code></li>
<li data-value="5">and a replica count of three.</li>
</ul>
</div>
</div>
</doc-view>
