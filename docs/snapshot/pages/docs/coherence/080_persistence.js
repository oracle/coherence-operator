<doc-view>

<v-layout row wrap>
<v-flex xs12 sm10 lg10>
<v-card class="section-def" v-bind:color="$store.state.currentColor">
<v-card-text class="pa-3">
<v-card class="section-def__card">
<v-card-text>
<dl>
<dt slot=title>Coherence Persistence</dt>
<dd slot="desc"><p>Coherence persistence is a set of tools and technologies that manage the persistence and recovery of Coherence
distributed caches. Cached data can be persisted so that it can be quickly recovered after a catastrophic failure
or after a cluster restart due to planned maintenance. Persistence and federated caching can be used together
as required.</p>
</dd>
</dl>
</v-card-text>
</v-card>
</v-card-text>
</v-card>
</v-flex>
</v-layout>

<h2 id="_configure_coherence_persistence">Configure Coherence Persistence</h2>
<div class="section">
<p>The <code>Coherence</code> CRD allows the default persistence mode, and the storage location of persistence data to be
configured. Persistence can be configured in the <code>spec.coherence.persistence</code> section of the CRD.
See the <a id="" title="" target="_blank" href="https://docs.oracle.com/en/middleware/standalone/coherence/14.1.1.0/administer/persisting-caches.html#GUID-3DC46E44-21E4-4DC4-9D12-231DE57FE7A1">Coherence Persistence</a>
documentation for more details of how persistence works and its configuration.</p>

</div>

<h2 id="_persistence_mode">Persistence Mode</h2>
<div class="section">
<p>There are three default persistence modes available, <code>active</code>, <code>active-async</code> and <code>on-demand</code>; the default mode is <code>on-demand</code>.
The persistence mode will be set using the <code>spec.coherence.persistence,mode</code> field in the CRD. The value of this field will be
used to set the <code>coherence.distributed.persistence-mode</code> system property in the Coherence JVM.</p>

<p>For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: storage
spec:
  coherence:
    persistence:
      mode: active  <span class="conum" data-value="1" /></markup>

<ul class="colist">
<li data-value="1">The example above sets the persistence mode to <code>active</code> which will effectively pass
<code>-Dcoherence.distributed.persistence-mode=active</code> to the Coherence JVM&#8217;s command line.</li>
</ul>
</div>

<h2 id="_persistence_storage">Persistence Storage</h2>
<div class="section">
<p>The purpose of persistence in Coherence is to store data on disc so that it is available outside of the lifetime of the
JVMs that make up the cluster. In a containerised environment like Kubernetes this means storing that data in storage that
also lives outside of the containers.</p>

<p>When persistence storage has been configured a <code>VolumeMount</code> will be added to the Coherence container mounted at <code>/persistence</code>,
and the <code>coherence.distributed.persistence.base.dir</code> system property will be configured to point to the storage location.</p>


<h3 id="_using_a_persistentvolumeclaim">Using a PersistentVolumeClaim</h3>
<div class="section">
<p>The Coherence Operator creates a <code>StatefulSet</code> for each <code>Coherence</code> resource, so the
logical place to store persistence data is in a <code>PersistentVolumeClaim</code>.</p>

<p>The PVC used for persistence can be configured in the <code>spec.coherence.persistence.persistentVolumeClaim</code> section
of the CRD.</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: storage
spec:
  coherence:
    persistence:
      persistentVolumeClaim:     <span class="conum" data-value="1" />
        storageClassName: "SSD"
        accessModes:
          - ReadWriteOnce
        resources:
          requests:
            storage: 50Gi</markup>

<ul class="colist">
<li data-value="1">The example above configures a 50GB PVC with a storage class name of "SSD"
(assuming the Kubernetes cluster has a storage class of that name configured).</li>
</ul>
<p>The configuration under the <code>spec.coherence.persistence.persistentVolumeClaim</code> section is exactly the same as
configuring a PVC for a normal Kubernetes Pod and all the possible options are beyond the scope of this document.
For more details on configuring PVC, see the Kubernetes
<a id="" title="" target="_blank" href="https://kubernetes.io/docs/concepts/storage/persistent-volumes/">Persistent Volumes</a> documentation.</p>

</div>

<h3 id="_using_a_normal_volume">Using a Normal Volume</h3>
<div class="section">
<p>An alternative to a PVC is to use a normal Kubernetes Volume to store the persistence data.
An example of this use-case could be when the Kubernetes Nodes that the Coherence Pods are scheduled onto have locally
attached fast SSD drives, which is ideal storage for persistence.
In this case a normal Volume can be configured in the <code>spec.coherence.persistence.volume</code> section of the CRD.</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: storage
spec:
  coherence:
    persistence:                                 <span class="conum" data-value="1" />
      volume:
        hostPath:
          path: /mnt/ssd/coherence/persistence</markup>

<ul class="colist">
<li data-value="1">In the example above a Volume has been configured for persistence, in this case a <code>HostPath</code> volume pointing to
the <code>/mnt/ssd/coherence/persistence</code> directory on the Node.</li>
</ul>
<p>The configuration under the <code>spec.coherence.persistence.volume</code> section is a normal Kubernetes
<a id="" title="" target="_blank" href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#volume-v1-core">VolumeSource</a>
so any valid <code>VolumeSource</code> configuration can be used.
See the Kubernetes <a id="" title="" target="_blank" href="https://kubernetes.io/docs/concepts/storage/volumes/">Volumes</a> documentation for more details.</p>

</div>
</div>

<h2 id="_snapshot_storage">Snapshot Storage</h2>
<div class="section">
<p>Coherence allows on-demand snapshots to be taken of cache data. With the default configuration the snapshot files will
be stored under the same persistence root location as active persistence data.
The <code>Coherence</code> spec allows a different location to be specified for storage of snapshot files so that active data
and snapshot data can be stored in different locations and/or on different storage types in Kubernetes.</p>

<p>The same two options are available for snapshot storage that are available for persistence storage, namely PVCs and
normal Volumes. The <code>spec.coherence.persistence.snapshots</code> section is used to configure snapshot storage.
When this is used a <code>VolumeMount</code> will be added to the Coherence container with a mount path of <code>/snapshots</code>,
and the <code>coherence.distributed.persistence.snapshot.dir</code> system property will be set to point to this location.</p>


<h3 id="_snapshots_using_a_persistentvolumeclaim">Snapshots Using a PersistentVolumeClaim</h3>
<div class="section">
<p>A PVC can be configured for persistence snapshot data as shown below.</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: storage
spec:
  coherence:
    persistence:                                 <span class="conum" data-value="1" />
      volume:
        hostPath:
          path: /mnt/ssd/coherence/persistence
      snapshots:
        persistentVolumeClaim:                   <span class="conum" data-value="2" />
          resources:
            requests:
              storage: 50Gi</markup>

<ul class="colist">
<li data-value="1">Active persistence data will be stored on a normal Volume using a HostPath volume source.</li>
<li data-value="2">Snapshot data will be stored in a 50GB PVC.</li>
</ul>
</div>

<h3 id="_snapshots_using_a_normal_volumes">Snapshots Using a Normal Volumes</h3>
<div class="section">
<p>A normal volume can be configured for snapshot data as shown below.</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: storage
spec:
  coherence:
    persistence:                                 <span class="conum" data-value="1" />
      volume:
        hostPath:
          path: /mnt/ssd/coherence/persistence
      snapshots:
        volume:
          hostPath:
            path: /mnt/ssd/coherence/snapshots   <span class="conum" data-value="2" /></markup>

<ul class="colist">
<li data-value="1">Active persistence data will be stored on a normal Volume using a HostPath volume source.</li>
<li data-value="2">Snapshot data will be stored on a normal Volume using a different HostPath volume source.</li>
</ul>
</div>
</div>
</doc-view>
