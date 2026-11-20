<doc-view>

<v-layout row wrap>
<v-flex xs12 sm10 lg10>
<v-card class="section-def" v-bind:color="$store.state.currentColor">
<v-card-text class="pa-3">
<v-card class="section-def__card">
<v-card-text>
<dl>
<dt slot=title>Adding Annotations</dt>
<dd slot="desc"><p>Annotations can be added to the Coherence cluster&#8217;s <code>StatefulSet</code> and the <code>Pods</code>.
See the official
<a id="" title="" target="_blank" href="https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/">Kubernetes Annotations</a>
documentation for more details on applying annotations to resources.</p>
</dd>
</dl>
</v-card-text>
</v-card>
</v-card-text>
</v-card>
</v-flex>
</v-layout>

<h2 id="_statefulset_annotations">StatefulSet Annotations</h2>
<div class="section">
<p>The default behaviour of the Operator is to copy any annotations added to the <code>Coherence</code> resource to the <code>StatefulSet</code>.
For example:</p>

<markup
lang="yaml"
title="coherence-cluster.yaml"
>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: storage
  annotations:
    key1: value1
    key2: value2</markup>

<p>This will result in a <code>StatefulSet</code> with the following annotations:</p>

<markup
lang="yaml"

>apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: storage
  annotations:
    key1: value1
    key2: value2</markup>

<p>Alternatively, if the <code>StatefulSet</code> should have different annotations to the <code>Coherence</code> resource, the annotations
for the <code>StatefulSet</code> can be specified in the <code>spec.statefulSetAnnotations</code> field of the <code>Coherence</code> resource.
For example:</p>

<markup
lang="yaml"
title="coherence-cluster.yaml"
>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: storage
  annotations:
    key1: value1
    key2: value2
spec:
  replicas: 3
  statefulSetAnnotations:
    key3: value3
    key4: value4</markup>

<p>This will result in a <code>StatefulSet</code> with the following annotations:</p>

<markup
lang="yaml"

>apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: storage
  annotations:
    key3: value3
    key4: value4</markup>

</div>

<h2 id="_pod_annotations">Pod Annotations</h2>
<div class="section">
<p>Additional annotations can be added to the <code>Pods</code> managed by the Operator.
Annotations should be added to the <code>annotations</code> map in the <code>Coherence</code> CRD spec.</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: test
spec:
  annotations:
    key1: value1
    key2: value2</markup>

<p>The annotations will be added the <code>Pods</code>:</p>

<markup
lang="yaml"

>apiVersion: v1
kind: Pod
metadata:
  name: storage-0
  annotations:
    key1: value1
    key2: value2</markup>

</div>
</doc-view>
