<doc-view>

<v-layout row wrap>
<v-flex xs12 sm10 lg10>
<v-card class="section-def" v-bind:color="$store.state.currentColor">
<v-card-text class="pa-3">
<v-card class="section-def__card">
<v-card-text>
<dl>
<dt slot=title>Generating Heap Dumps</dt>
<dd slot="desc"><p>Some of the debugging techniques described in <a id="" title="" target="_blank" href="https://docs.oracle.com/en/middleware/fusion-middleware/coherence/12.2.1.4/develop-applications/debugging-coherence.html">Debugging in Coherence</a>
require the creation of files, such as log files and JVM heap dumps, for analysis. You can also create and extract these files in the Coherence Operator.</p>
</dd>
</dl>
</v-card-text>
</v-card>
</v-card-text>
</v-card>
</v-flex>
</v-layout>

<h2 id="_produce_and_extract_a_heap_dump">Produce and extract a heap dump</h2>
<div class="section">
<p>This example shows how to collect a .hprof file for a heap dump.</p>

<p>A single-command technique is also included at the end of this sample.</p>

<div class="admonition note">
<p class="admonition-inline">Coherence Pods are configured to  produce a heap dump on OOM error by default. See
<router-link to="/clusters/080_jvm">Configure The JVM</router-link> for more information.</p>
</div>
<div class="admonition note">
<p class="admonition-inline">You cal also trigger a heap dump via the <a id="" title="" target="_blank" href="https://docs.oracle.com/en/middleware/fusion-middleware/coherence/12.2.1.4/rest-reference/op-management-coherence-cluster-members-memberidentifier-dumpheap-post.html">Management over REST API</a>.</p>
</div>

<h3 id="_1_install_a_coherence_cluster">1. Install a Coherence Cluster</h3>
<div class="section">
<p>Deploy a simple <code>CoherenceCluster</code> resource with a single role like this:</p>

<markup
lang="yaml"
title="heapdump-cluster.yaml"
>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: heapdump-cluster
spec:
  role: storage
  replicas: 3</markup>

<div class="admonition note">
<p class="admonition-inline">Add an <code>imagePullSecrets</code> entry if required to pull images from a private repository.</p>
</div>
<markup
lang="bash"

>kubectl create -n &lt;namespace&gt; -f  heapdump-cluster.yaml

coherencecluster.coherence.oracle.com/heapdump-cluster created

kubectl -n &lt;namespace&gt; get pod -l coherenceCluster=heapdump-cluster

NAME                         READY   STATUS    RESTARTS   AGE
heapdump-cluster-storage-0   1/1     Running   0          59s
heapdump-cluster-storage-1   1/1     Running   0          59s
heapdump-cluster-storage-2   1/1     Running   0          59s</markup>

</div>

<h3 id="_2_obtain_the_pid_of_the_coherence_process">2. Obtain the PID of the Coherence process</h3>
<div class="section">
<p>Obtain the PID of the Coherence process. Generally, the PID is 1. You can also use jps to get the actual PID.</p>

<markup
lang="bash"

>kubectl exec -it -n coherence-example heapdump-cluster-storage-0 -- bash

$  jps
1 Main
153 Jps</markup>

<div class="admonition note">
<p class="admonition-inline">The process with <code>Main</code> is the main process that calls <code>DefaultCacheServer</code> to start a cluster node.</p>
</div>
</div>

<h3 id="_3_use_the_jcmd_command_to_extract_the_heap_dump">3. Use the jcmd command to extract the heap dump</h3>
<div class="section">
<markup
lang="bash"

>$ rm -f /tmp/heap.hprof
$ /usr/java/default/bin/jcmd 1 GC.heap_dump /tmp/heap.hprof
$ exit</markup>

</div>

<h3 id="_4_copy_the_heap_dump_to_local_machine">4. Copy the heap dump to local machine</h3>
<div class="section">
<markup
lang="bash"

>kubectl cp &lt;namespace&gt;/heapdump-cluster-storage-0:/tmp/heap.hprof heap.hprof

tar: Removing leading `/' from member names

ls -l heap.hprof

-rw-r--r--  1 user  staff  21113314 15 Oct 08:50 heap.hprof</markup>

<div class="admonition note">
<p class="admonition-inline">Depending upon whether the Kubernetes cluster is local or remote, this might take some time.</p>
</div>
</div>

<h3 id="_5_single_command_usage">5. Single command usage</h3>
<div class="section">
<p>Assuming that the Coherence PID is 1, you can use this repeatable single-command technique to extract the heap dump:</p>

<markup
lang="bash"

>(kubectl exec heapdump-cluster-storage-0 -n &lt;namespace&gt;  -- /bin/bash -c \
  "rm -f /tmp/heap.hprof; /usr/java/default/bin/jcmd 1 GC.heap_dump /tmp/heap.hprof; cat /tmp/heap.hprof &gt; /dev/stderr" ) 2&gt; heap.hprof</markup>

</div>

<h3 id="_6_clean_up">6. Clean Up</h3>
<div class="section">
<p>After running the above the Coherence cluster can be removed using <code>kubectl</code>:</p>

<markup
lang="bash"

>kubectl -n &lt;namespace&gt; delete -f heapdump-cluster.yaml</markup>

</div>
</div>
</doc-view>
