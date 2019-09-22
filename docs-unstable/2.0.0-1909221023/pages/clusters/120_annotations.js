<doc-view>

<v-layout row wrap>
<v-flex xs12 sm10 lg10>
<v-card class="section-def" v-bind:color="$store.state.currentColor">
<v-card-text class="pa-3">
<v-card class="section-def__card">
<v-card-text>
<dl>
<dt slot=title>Pod Annotations</dt>
<dd slot="desc"><p>Annotations can be added to the <code>Pods</code> of a role in a <code>CoherenceCluster</code>.</p>
</dd>
</dl>
</v-card-text>
</v-card>
</v-card-text>
</v-card>
</v-flex>
</v-layout>

<h2 id="_pod_annotations">Pod Annotations</h2>
<div class="section">
<p>Custom <a id="" title="" target="_blank" href="https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/">annotations</a>
can be added to the spec of a role which will then be added to all <code>Pods</code> for that role created by
the Coherence Operator.</p>


<h3 id="_configure_pod_annotations_for_the_implicit_role">Configure Pod Annotations for the Implicit Role</h3>
<div class="section">
<p>When creating a <code>CoherenceCluster</code> with a single implicit role annotations can be defined at the <code>spec</code> level.
Annotations are defined as a map of string key value pairs, for example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  annotations:
    acme.com/layer: back   <span class="conum" data-value="1" /></markup>

<ul class="colist">
<li data-value="1">The implicit role will have the annotation <code>acme.com/layer : back</code></li>
</ul>
<p>This will result in a <code>StatefulSet</code> for the role with the annotation added to the <code>PodSpec</code>.</p>

<markup
lang="yaml"

>apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: test-cluster-storage
spec:
  replicas: 3
  selector:
    matchLabels:
      coherenceDeployment: test-cluster-storage
      component: coherencePod
  serviceName: test-cluster-storage
  template:
    metadata:
      annotations:
        acme.com/layer: back   <span class="conum" data-value="1" /></markup>

<ul class="colist">
<li data-value="1">The annotation <code>acme.com/layer: back</code> has been applied to the <code>StatefulSet</code> <code>Pod</code> template.</li>
</ul>
</div>

<h3 id="_configure_pod_annotations_for_explicit_roles">Configure Pod Annotations for Explicit Roles</h3>
<div class="section">
<p>When creating a <code>CoherenceCluster</code> with explicit roles in the <code>roles</code> list annotations can be defined for each role, for example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  roles:
    - role: data
      annotations:
        acme.com/layer: back   <span class="conum" data-value="1" />
    - role: proxy
      annotations:
        acme.com/layer: front   <span class="conum" data-value="2" /></markup>

<ul class="colist">
<li data-value="1">The <code>data</code> role will have the annotation <code>acme.com/layer: back</code></li>
<li data-value="2">The <code>proxy</code> role will have the annotation <code>acme.com/layer: front</code></li>
</ul>
</div>

<h3 id="_configure_pod_annotations_for_explicit_roles_with_defaults">Configure Pod Annotations for Explicit Roles With Defaults</h3>
<div class="section">
<p>When creating a <code>CoherenceCluster</code> with explicit roles in the <code>roles</code> list annotations can be defined as defaults
applied to all roles and also for each role, for example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  annotations:
    acme.com/layer: back          <span class="conum" data-value="1" />
    acme.com/app:   orders
  roles:
    - role: data                  <span class="conum" data-value="2" />
    - role: proxy                 <span class="conum" data-value="3" />
      annotations:
        acme.com/state: none
    - role: web                   <span class="conum" data-value="4" />
      annotations:
        acme.com/three: none
        acme.com/layer: front</markup>

<ul class="colist">
<li data-value="1">There are two default annotations <code>acme.com/layer : back</code> and <code>acme.com/app : orders</code> that will apply to all <code>Pods</code>
in all roles unless specifically overridden.</li>
<li data-value="2">The <code>data</code> role has no other annotations defined so will just have the default annotations <code>acme.com/layer : back</code>
and <code>acme.com/app : orders</code></li>
<li data-value="3">The <code>proxy</code> role specified an annotation <code>acme.com/state : none</code> so will have this annotation as well as the
default annotations <code>acme.com/layer : back</code> and <code>acme.com/app : orders</code></li>
<li data-value="4">The <code>web</code> role specifies the <code>acme.com/three: none</code> annotation and also overrides the <code>acme.com/layer</code> annotation
with the value <code>front</code> so it will have three annotations, <code>acme.com/three: none</code> , <code>acme.com/layer : front</code>
and <code>acme.com/app : orders</code></li>
</ul>
</div>
</div>
</doc-view>
