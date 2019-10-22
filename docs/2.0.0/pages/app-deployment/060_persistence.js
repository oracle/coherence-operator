<doc-view>

<v-layout row wrap>
<v-flex xs12 sm10 lg10>
<v-card class="section-def" v-bind:color="$store.state.currentColor">
<v-card-text class="pa-3">
<v-card class="section-def__card">
<v-card-text>
<dl>
<dt slot=title>Persistence</dt>
<dd slot="desc"><p>The Coherence persistence feature is used to save a cache to disk and ensures that cache
data can always be recovered especially in the case of a full cluster restart or re-creation.</p>
</dd>
</dl>
</v-card-text>
</v-card>
</v-card-text>
</v-card>
</v-flex>
</v-layout>

<h2 id="_using_coherence_persistence">Using Coherence Persistence</h2>
<div class="section">
<p>When enabling Persistence in the Coherence Operator, you have two options:</p>

<ul class="ulist">
<li>
<p>Use the default Persistent Volume Claim (PVC) - PVC&#8217;s will be automatically be created and bound to pods on startup</p>

</li>
<li>
<p>Specify existing persistent volumes - allows full control of the underlying allocated volumes</p>

</li>
</ul>
<p>For more information on Coherence Persistence, please see the
<a id="" title="" target="_blank" href="https://docs.oracle.com/en/middleware/fusion-middleware/coherence/12.2.1.4/administer/persisting-caches.html">Coherence Documentation</a>.</p>

</div>

<h2 id="_table_of_contents">Table of Contents</h2>
<div class="section">
<ol style="margin-left: 15px;">
<li>
<router-link to="#prereqs" @click.native="this.scrollFix('#prereqs')">Prerequisites</router-link>

</li>
<li>
<router-link to="#default" @click.native="this.scrollFix('#default')">Use Default Persistent Volume Claim</router-link>

</li>
</ol>
</div>

<h2 id="prereqs">Prerequisites</h2>
<div class="section">
<ol style="margin-left: 15px;">
<li>
Install the Coherence Operator

</li>
<li>
Create any secrets required to pull Docker images

</li>
<li>
Create a new <code>working directory</code> and change to that directory

</li>
</ol>
</div>

<h2 id="default">Use Default Persistent Volume Claim</h2>
<div class="section">
<p>By default, when you enable Coherence Persistence, the required infrastructure in
terms of persistent volumes (PV) and persistent volume claims (PVC) is set up automatically. Also, the persistence-mode
is set to <code>active</code>. This allows the Coherence cluster to be restarted and the data to be retained.</p>

<p>This example shows how to enable Persistence with all the defaults.</p>


<h3 id="_1_create_the_coherence_cluster_yaml">1. Create the Coherence cluster yaml</h3>
<div class="section">
<p>In your working directory directory create a file called <code>persistence-cluster.yaml</code> with the following contents:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: persistence-cluster
spec:
  jvm:
    memory:
      heapSize: 512m
  replicas: 3
  coherence:
    persistence:
      enabled: true                    <span class="conum" data-value="1" />
      persistentVolumeClaim:
        accessModes:
        - ReadWriteOnce                <span class="conum" data-value="2" />
        resources:
          requests:
            storage: 1Gi               <span class="conum" data-value="3" /></markup>

<ul class="colist">
<li data-value="1">Enables <code>Active</code> Persistence</li>
<li data-value="2">Specifies that the volume can be mounted as read-write by a single node</li>
<li data-value="3">Sets the size of the Persistent Volume</li>
</ul>
<div class="admonition note">
<p class="admonition-inline">Add an <code>imagePullSecrets</code> entry if required to pull images from a private repository.</p>
</div>
</div>

<h3 id="_2_install_the_coherence_cluster">2. Install the Coherence Cluster</h3>
<div class="section">
<p>Issue the following to install the cluster:</p>

<markup
lang="bash"

>kubectl create -n &lt;namespace&gt; -f persistence-cluster.yaml

coherencecluster.coherence.oracle.com/persistence-cluster created

kubectl -n &lt;namespace&gt; get pod -l coherenceCluster=persistence-cluster

NAME                            READY   STATUS    RESTARTS   AGE
persistence-cluster-storage-0   1/1     Running   0          79s
persistence-cluster-storage-1   0/1     Running   0          79s
persistence-cluster-storage-2   0/1     Running   0          79s</markup>

<p>Check the Persistent Volumes and PVC are automatically created.</p>

<markup
lang="bash"

>kubectl get pvc -n &lt;namespace&gt;

NAME                                               STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS   AGE
persistence-volume-persistence-cluster-storage-0   Bound    pvc-060c61d6-ee2d-11e9-aa71-025000000001   1Gi        RWO            hostpath       2m32s
persistence-volume-persistence-cluster-storage-1   Bound    pvc-061204e8-ee2d-11e9-aa71-025000000001   1Gi        RWO            hostpath       2m32s
persistence-volume-persistence-cluster-storage-2   Bound    pvc-06205b32-ee2d-11e9-aa71-025000000001   1Gi        RWO            hostpath       2m32s</markup>

<p>Wait until all nodes are Running and READY before continuing.</p>


<h4 id="_3_connect_to_the_coherence_console_to_add_data">3. Connect to the Coherence Console to add data</h4>
<div class="section">
<markup
lang="bash"

>kubectl exec -it -n &lt;namespace&gt; persistence-cluster-storage-0 bash /scripts/startCoherence.sh console</markup>

<p>At the prompt type the following to create a cache called <code>test</code>:</p>

<markup
lang="bash"

>cache test</markup>

<p>Use the following to create 10,000 entries of 100 bytes:</p>

<markup
lang="bash"

>bulkput 10000 100 0 100</markup>

<p>Lastly issue the command <code>size</code> to verify the cache entry count.</p>

<p>Type <code>bye</code> to exit the console.</p>

</div>

<h4 id="_4_delete_the_cluster">4. Delete the cluster</h4>
<div class="section">
<div class="admonition note">
<p class="admonition-inline">This will not delete the PVC&#8217;s.</p>
</div>
<markup
lang="bash"

>kubectl -n &lt;namespace&gt; delete -f persistence-cluster.yaml</markup>

<p>Use <code>kubectl get pods -n &lt;namespace&gt;</code> to confirm the pods have terminated.</p>

</div>

<h4 id="_5_confirm_the_pvcs_are_still_present">5. Confirm the PVC&#8217;s are still present</h4>
<div class="section">
<markup
lang="bash"

>kubectl get pvc -n &lt;namespace&gt;

NAME                                               STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS   AGE
persistence-volume-persistence-cluster-storage-0   Bound    pvc-060c61d6-ee2d-11e9-aa71-025000000001   1Gi        RWO            hostpath       2m32s
persistence-volume-persistence-cluster-storage-1   Bound    pvc-061204e8-ee2d-11e9-aa71-025000000001   1Gi        RWO            hostpath       2m32s
persistence-volume-persistence-cluster-storage-2   Bound    pvc-06205b32-ee2d-11e9-aa71-025000000001   1Gi        RWO            hostpath       2m32s</markup>

</div>

<h4 id="_6_re_install_the_coherence_cluster">6. Re-install the Coherence cluster</h4>
<div class="section">
<markup
lang="bash"

>kubectl create -n &lt;namespace&gt; -f persistence-cluster.yaml

coherencecluster.coherence.oracle.com/persistence-cluster created

kubectl -n &lt;namespace&gt; get pod -l coherenceCluster=persistence-cluster

NAME                            READY   STATUS    RESTARTS   AGE
persistence-cluster-storage-0   1/1     Running   0          79s
persistence-cluster-storage-1   0/1     Running   0          79s
persistence-cluster-storage-2   0/1     Running   0          79s</markup>

<p>Wait until the pods are Running and Ready, then confirm the data is still present by using the
<code>cache test</code> and <code>size</code> commands only as in step 3 above.</p>

</div>

<h4 id="_7_uninstall_the_cluster_and_pvcs">7. Uninstall the Cluster and PVC&#8217;s</h4>
<div class="section">
<p>Issue the following to delete the Coherence cluster.</p>

<markup
lang="bash"

>kubectl -n &lt;namespace&gt; delete -f persistence-cluster.yaml</markup>

<p>Ensure all the pods have all terminated before you delete the PVC&#8217;s.</p>

<markup
lang="bash"

>kubectl get pvc -n &lt;namespace&gt; | sed 1d | awk '{print $1}' | xargs kubectl delete pvc -n &lt;namespace&gt;</markup>

</div>
</div>
</div>
</doc-view>
