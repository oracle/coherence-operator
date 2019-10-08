<doc-view>

<v-layout row wrap>
<v-flex xs12 sm10 lg10>
<v-card class="section-def" v-bind:color="$store.state.currentColor">
<v-card-text class="pa-3">
<v-card class="section-def__card">
<v-card-text>
<dl>
<dt slot=title>Quick Start</dt>
<dd slot="desc"><p>This guide is a simple set of steps to install the Coherence Operator and then use that
to install a simple Coherence cluster.</p>
</dd>
</dl>
</v-card-text>
</v-card>
</v-card-text>
</v-card>
</v-flex>
</v-layout>

<h2 id="_prerequisites">Prerequisites</h2>
<div class="section">
<p>Ensure that the <router-link to="/install/02_prerequisites">Coherence Operator prerequisites</router-link> are available.</p>

</div>

<h2 id="_1_install_the_coherence_operator">1. Install the Coherence Operator</h2>
<div class="section">

<h3 id="_1_1_add_the_coherence_operator_helm_repository">1.1 Add the Coherence Operator Helm repository</h3>
<div class="section">
<markup
lang="bash"

>helm repo add coherence https://oracle.github.io/coherence-operator/charts

helm repo update</markup>

</div>

<h3 id="_1_2_install_the_coherence_operator_helm_chart">1.2. Install the Coherence Operator Helm chart</h3>
<div class="section">
<markup
lang="bash"

>helm install  \
    --namespace &lt;namespace&gt; \
    --name &lt;release-name&gt; \
    coherence/coherence-operator</markup>

<div class="admonition note">
<p class="admonition-inline">Use the same namespace that the operator was installed into,
e.g. if the namespace is <code>coherence</code> the command would be
<code>helm install --namespace coherence --name operator coherence/coherence-operator</code></p>
</div>
<p>See the <router-link to="/install/01_introduction">full install guide</router-link> for more details.</p>

</div>
</div>

<h2 id="_2_install_a_coherence_cluster">2. Install a Coherence Cluster</h2>
<div class="section">

<h3 id="_2_1_install_a_coherence_cluster_using_the_minimal_required_configuration">2.1 Install a Coherence cluster using the minimal required configuration.</h3>
<div class="section">
<p>The minimal required yaml to create a <code>CoherenceCluster</code> resource is shown below.</p>

<markup
lang="yaml"
title="my-cluster.yaml"
>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: my-cluster  <span class="conum" data-value="1" /></markup>

<p>The only required field is <code>metadata.name</code> which will be used as the Coherence cluster name, in this case <code>my-cluster</code></p>

<markup


>kubectl -n &lt;namespace&gt; apply -f my-cluster.yaml</markup>

<div class="admonition note">
<p class="admonition-inline">Use the same namespace that the operator was installed into,
e.g. if the namespace is <code>coherence</code> the command would be
<code>kubectl -n coherence create -f my-cluster.yaml</code></p>
</div>
</div>

<h3 id="_2_2_list_the_coherence_resources">2.2 List the Coherence Resources</h3>
<div class="section">
<p>After installing the <code>my-cluster.yaml</code> above here should be a single <code>coherencecluster</code> resource  named <code>my-cluster</code>
and a single <code>coherencerole</code> resource named <code>my-cluster-storage</code> created in the Coherence Operator namespace.</p>

<markup


>kubectl -n &lt;namespace&gt; get coherencecluster</markup>

<div class="admonition note">
<p class="admonition-inline">Use the same namespace that the operator was installed into, e.g. if the namespace is <code>coherence</code> the command
would be <code>kubectl -n coherence get coherence</code></p>
</div>
<markup


>NAME                                                    AGE
coherencerole.coherence.oracle.com/my-cluster-storage   19s

NAME                                               AGE
coherencecluster.coherence.oracle.com/my-cluster   19s</markup>

<p>See the <router-link to="/clusters/020_k8s_resources">in-depth documentation</router-link> on the Kubernetes resources created by the
Coherence Operator.</p>

</div>

<h3 id="_2_3_list_all_of_the_pods_for_the_coherence_cluster">2.3 List all of the <code>Pods</code> for the Coherence cluster.</h3>
<div class="section">
<p>The Coherence Operator applies a <code>coherenceCluster</code> label to all
of the <code>Pods</code> so this label can be used with the <code>kubectl</code> command to find <code>Pods</code> for a Coherence cluster.</p>

<markup


>kubectl -n &lt;namespace&gt; get pod -l coherenceCluster=my-cluster</markup>

<div class="admonition note">
<p class="admonition-inline">Use the same namespace that the operator was installed into,
e.g. if the namespace is <code>coherence</code> the command would be
<code>kubectl -n coherence get pod -l coherenceCluster=my-cluster</code></p>
</div>
<markup


>NAME                   READY   STATUS    RESTARTS   AGE
my-cluster-storage-0   1/1     Running   0          2m58s
my-cluster-storage-1   1/1     Running   0          2m58s
my-cluster-storage-2   1/1     Running   0          2m58s</markup>

<p>The default cluster size is three so there should be three <code>Pods</code></p>

</div>
</div>

<h2 id="_3_scale_the_coherence_cluster">3. Scale the Coherence Cluster</h2>
<div class="section">

<h3 id="_3_1_use_kubectl_to_scale_up">3.1 Use kubectl to Scale Up</h3>
<div class="section">
<p>Using the <code>kubectl scale</code> command a specific <code>CoherenceRole</code> can be scaled up or down.</p>

<markup


>kubectl -n &lt;namespace&gt; scale coherencerole/storage --replicas=6</markup>

<div class="admonition note">
<p class="admonition-inline">Use the same namespace that the operator was installed into,
e.g. if the namespace is <code>coherence</code> the command would be
<code>kubectl -n coherence scale coherencerole/my-cluster-storage --replicas=6</code></p>
</div>
</div>

<h3 id="_3_2_list_all_of_the_pods_fo_the_coherence_cluster">3.2 List all of the <code>Pods</code> fo the Coherence cluster</h3>
<div class="section">
<markup


>kubectl -n &lt;namespace&gt; get pod -l coherenceCluster=my-cluster</markup>

<div class="admonition note">
<p class="admonition-inline">Use the same namespace that the operator was installed into,
e.g. if the namespace is <code>coherence</code> the command would be
<code>kubectl -n coherence get pod -l coherenceCluster=my-cluster</code></p>
</div>
<markup


>NAME                   READY   STATUS    RESTARTS   AGE
my-cluster-storage-0   1/1     Running   0          4m23s
my-cluster-storage-1   1/1     Running   0          4m23s
my-cluster-storage-2   1/1     Running   0          4m23s
my-cluster-storage-3   1/1     Running   0          1m19s
my-cluster-storage-4   1/1     Running   0          1m19s
my-cluster-storage-5   1/1     Running   0          1m19s</markup>

<p>There should eventually be six running <code>Pods</code>.</p>

</div>
</div>
</doc-view>
