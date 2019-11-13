<doc-view>

<v-layout row wrap>
<v-flex xs12 sm10 lg10>
<v-card class="section-def" v-bind:color="$store.state.currentColor">
<v-card-text class="pa-3">
<v-card class="section-def__card">
<v-card-text>
<dl>
<dt slot=title>Configure Pod Labels</dt>
<dd slot="desc"><p>Labels can be added to the <code>Pods</code> of a role in a <code>CoherenceCluster</code>.</p>
</dd>
</dl>
</v-card-text>
</v-card>
</v-card-text>
</v-card>
</v-flex>
</v-layout>

<h2 id="_configure_pod_labels">Configure Pod Labels</h2>
<div class="section">
<p>Custom Pod <a id="" title="" target="_blank" href="https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/">labels</a>
can be added to the spec of a role which will then be added to all <code>Pods</code> for that role created by
the Coherence Operator.</p>


<h3 id="_default_labels">Default Labels</h3>
<div class="section">
<p>The Coherence Operator applies the following labels to a role. These labels should not be overridden as they
are used by the Coherence Operator.</p>


<div class="table__overflow elevation-1 ">
<table class="datatable table">
<colgroup>
<col style="width: 50%;">
<col style="width: 50%;">
</colgroup>
<thead>
<tr>
<th>Label</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>coherenceCluster</td>
<td>This label will be set to the owning <code>CoherenceCluster</code> name</td>
</tr>
<tr>
<td>coherenceRole</td>
<td>This label will be set to the role name</td>
</tr>
<tr>
<td>coherenceDeployment</td>
<td>This label will be the concatenated cluster name and role name in the format of the format <code>ClusterName-RoleName</code></td>
</tr>
<tr>
<td>component</td>
<td>This label is always <code>coherencePod</code></td>
</tr>
</tbody>
</table>
</div>
<p>The default labels above make it simple to find all <code>Pods</code> for a Coherence cluster or for a role when querying
Kubernetes (for example with <code>kubectl get</code>).</p>

</div>

<h3 id="_configure_pod_labels_for_the_implicit_role">Configure Pod Labels for the Implicit Role</h3>
<div class="section">
<p>When creating a <code>CoherenceCluster</code> with a single implicit role labels can be defined at the <code>spec</code> level.
Labels are defined as a map of string key value pairs, for example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  labels:            <span class="conum" data-value="1" />
    key1 : value1
    key2 : value2</markup>

<ul class="colist">
<li data-value="1">The implicit role will have the labels <code>key1=value1</code> and <code>key2=value2</code> which will result in all <code>Pods</code>
for the role also having those same labels.</li>
</ul>
</div>

<h3 id="_configure_pod_labels_for_explicit_roles">Configure Pod Labels for Explicit Roles</h3>
<div class="section">
<p>When creating a <code>CoherenceCluster</code> with explicit roles in the <code>roles</code> list labels can be defined for each role,
for example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  roles:
    - role: data
      labels:
        key1 : value1   <span class="conum" data-value="1" />
    - role: proxy
      labels:
        key2 : value2   <span class="conum" data-value="2" /></markup>

<ul class="colist">
<li data-value="1">The <code>data</code> role will have the label <code>key1=value1</code></li>
<li data-value="2">The <code>proxy</code> role will have the labels <code>key2=value2</code></li>
</ul>
</div>

<h3 id="_configure_pod_labels_for_explicit_roles_with_defaults">Configure Pod Labels for Explicit Roles With Defaults</h3>
<div class="section">
<p>When creating a <code>CoherenceCluster</code> with explicit roles in the <code>roles</code> list labels can be defined as defaults
applied to all roles and also for each role. The default labels will be merged with the role labels.
Where labels exist with the same key in both the defaults and the role then the labels in the role will take precedence.</p>

<p>For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  labels:                         <span class="conum" data-value="1" />
    key1 : value1
    key2 : value2
  roles:
    - role: data                  <span class="conum" data-value="2" />
    - role: proxy                 <span class="conum" data-value="3" />
      labels:
        key3 : value3
    - role: web                   <span class="conum" data-value="4" />
      labels:
        key2 : value-two
        key3 : value3</markup>

<ul class="colist">
<li data-value="1">There are two default labels <code>key1=value1</code> and <code>key2=value2</code> that will apply to all <code>Pods</code>
in all roles unless specifically overridden.</li>
<li data-value="2">The <code>data</code> role has no other labels defined so will just have the default labels <code>key1=value1</code> and <code>key2=value2</code></li>
<li data-value="3">The <code>proxy</code> role specified an labels <code>key3=value3</code> so will have this labels as well as the default labels
<code>key1=value1</code> and <code>key2=value2</code></li>
<li data-value="4">The <code>web</code> role specifies the <code>key3=value3</code> labels and also overrides the <code>key2</code> label with the value <code>value-two</code>
so it will have three labels, <code>key1=value1</code> , <code>key2=value-two</code> and <code>key3=value3</code></li>
</ul>
</div>
</div>
</doc-view>
