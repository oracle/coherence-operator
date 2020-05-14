<doc-view>

<v-layout row wrap>
<v-flex xs12 sm10 lg10>
<v-card class="section-def" v-bind:color="$store.state.currentColor">
<v-card-text class="pa-3">
<v-card class="section-def__card">
<v-card-text>
<dl>
<dt slot=title>Accessing the Console</dt>
<dd slot="desc"><p>The Coherence Console is a useful debugging and diagnosis tool usually used by administrators.</p>
</dd>
</dl>
</v-card-text>
</v-card>
</v-card-text>
</v-card>
</v-flex>
</v-layout>

<h2 id="_accessing_the_coherence_console">Accessing the Coherence Console</h2>
<div class="section">
<p>The example shows how to access the Coherence Console in a running cluster.</p>

<div class="admonition note">
<p class="admonition-inline">The Coherence Console is for advanced Coherence users and use-cases and care should be taken when using it.</p>
</div>

<h3 id="_1_install_a_coherence_cluster">1. Install a Coherence Cluster</h3>
<div class="section">
<p>Deploy a simple <code>CoherenceCluster</code> resource with a single role like this:</p>

<markup
lang="yaml"
title="example-cluster.yaml"
>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: example-cluster
spec:
  role: storage
  replicas: 3</markup>

<div class="admonition note">
<p class="admonition-inline">Add an <code>imagePullSecrets</code> entry if required to pull images from a private repository.</p>
</div>
<markup
lang="bash"

>kubectl create -n &lt;namespace&gt; -f  example-cluster.yaml

coherencecluster.coherence.oracle.com/example-cluster created

kubectl -n &lt;namespace&gt; get pod -l coherenceCluster=example-cluster

NAME                        READY   STATUS    RESTARTS   AGE
example-cluster-storage-0   1/1     Running   0          59s
example-cluster-storage-1   1/1     Running   0          59s
example-cluster-storage-2   1/1     Running   0          59s</markup>

</div>

<h3 id="_2_connect_to_the_coherence_console_to_add_data">2. Connect to the Coherence Console to add data</h3>
<div class="section">
<markup
lang="bash"

>kubectl exec -it -n &lt;namespace&gt; example-cluster-storage-0 bash /scripts/startCoherence.sh console</markup>

<p>At the prompt type the following to create a cache called <code>test</code>:</p>

<markup
lang="bash"

>cache test</markup>

<p>Use the following to create 10,000 entries of 100 bytes:</p>

<markup
lang="bash"

>bulkput 10000 100 0 100</markup>

<p>Issue the command <code>size</code> to verify the cache entry count.</p>

<p>Lastly issue the <code>help</code> command to show all available commands.</p>

<p>Type <code>bye</code> to exit the console.</p>

</div>

<h3 id="_3_clean_up">3. Clean Up</h3>
<div class="section">
<p>After running the above the Coherence cluster can be removed using <code>kubectl</code>:</p>

<markup
lang="bash"

>kubectl -n &lt;namespace&gt; delete -f example-cluster.yaml</markup>

</div>
</div>
</doc-view>
