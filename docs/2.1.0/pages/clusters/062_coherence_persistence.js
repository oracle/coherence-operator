<doc-view>

<h2 id="_coherence_persistence">Coherence Persistence</h2>
<div class="section">
<p>Coherence persistence is a set of tools and technologies that manage the persistence and recovery of Coherence
distributed caches. Cached data is persisted so that it can be quickly recovered after a catastrophic failure or
after a cluster restart due to planned maintenance. See the main
<a id="" title="" target="_blank" href="https://docs.oracle.com/en/middleware/fusion-middleware/coherence/12.2.1.4/administer/persisting-caches.html#GUID-3DC46E44-21E4-4DC4-9D12-231DE57FE7A1">Coherence documentation</a></p>

<p>The Coherence Operator supports configuring Coherence Persistence in two parts, snapshots and continuous persistence.
Snapshots is the process of saving the state of caches to a named snapshot as a set of files on disc.
This cache state can later be restored by reloading a named snapshot from disc - like a backup/restore operation.
Continuous persistence is where  Coherence continually writes the sate of caches to disc. When a Coherence cluster is
stopped and restarted, either on purpose or due to a failue, the data on disc is automatically reloaded and the cache
state is restored.</p>

<p>Ideally, the storage used for persistence and snapshots is fast local storage such as SSD. When using stand-alone
Coherence it is a simple process to manage local storage but when using Coherence in containers, and especially inside
Kubernetes, managing storage is a more complex task when that storage needs to be persisted longer than the lifetime of
the containers and re-attached to the containers if they are restarted.
The Coherence Operator aims to make using Coherence persistence in Kubernetes simpler by allowing the more common
use-cases to be easily configured.</p>

<p>Each role in a <code>CoherenceCluster</code> resource maps to a Kubernetes <code>StatefulSet</code>. One of the advantages of <code>StatefulSets</code>
is that they allow easy management of <code>PersistentVolumeClaims</code> which are ideal for use as storage for Coherence
persistence as they have a lifetime outside of the <code>Pods</code> are are reattached to the <code>StatefulSet</code> <code>Pods</code> when they are
restarted.</p>

</div>

<h2 id="pvc">Managing Coherence Snapshots using Persistent Volumes</h2>
<div class="section">
<p>When managing Coherence clusters using the Coherence Operator the simples configuration is to write snapshots to a
volume mapped to a <code>PersistentVolumeClaim</code> and to let the <code>StatefulSet</code> manage the <code>PVCs</code>.</p>

<p>Snapshots are configured in the <code>coherence.snapshots</code> section of the role specification in the <code>CoherenceCluster</code> CRD.</p>

<markup
lang="yaml"

>coherence:
  snapshots:
    enabled: true            <span class="conum" data-value="1" />
      persistentVolumeClaim: <span class="conum" data-value="2" />
        # PVC spec...</markup>

<ul class="colist">
<li data-value="1">Snapshots should be enabled by setting the <code>coherence.snapshots.enabled</code> field to true.</li>
<li data-value="2">The <code>persistentVolumeClaim</code> section allows the <code>PVC</code> used for snapshot files to be configured.</li>
</ul>
<p>The default value for <code>coherence.snapshots.enabled</code> is <code>false</code> so no snapshot location will be configured for Coherence
caches to use.</p>

<p>If <code>snapshots.enabled</code> is either undefined or false it is still possible to use Coherence snapshot functionality but in
this case snapshot files will be written to storage inside the Coherence container and will be lost if the container is
shutdown.</p>

<p>For example, if a Kubernetes cluster has a custom <code>StorageClass</code> named <code>fast</code> defined below:</p>

<markup
lang="yaml"

>kind: StorageClass
apiVersion: storage.k8s.io/v1
metadata:
  name: fast
provisioner: k8s.io/minikube-hostpath
parameters:
  type: pd-ssd</markup>

<p>Then a <code>CoherenceCluster</code> can be created with snapshots enabled and configured to use the <code>fast</code> storage class:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  coherence:
    snapshots:
      enabled: true
        persistentVolumeClaim:
          accessModes: [ "ReadWriteOnce" ]
          storageClassName: fast
          resources:
            requests:
              storage: 1Gi</markup>

<p>The content of the <code>persistentVolumeClaim</code> is any valid yaml for defining a <code>PersistentVolumeClaimSpec</code> that would be
allowed when configuring the <code>spec</code> section of a PVC in the <code>volumeClaimTemplates</code> section of a <code>StatefulSet</code> as
described in the <a id="" title="" target="_blank" href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.16/#persistentvolumeclaimspec-v1-core">Kubernetes API documentation</a>.</p>


<h3 id="_snapshots_using_persistent_volumes_for_a_single_implicit_role">Snapshots using Persistent Volumes for a Single Implicit Role</h3>
<div class="section">
<p>When configuring a <code>CoherenceCluster</code> with a single implicit role the <code>coherence.snapshots</code> configuration is added
directly to the <code>spec</code> section of the <code>CoherenceCluster</code>.
For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  coherence:
    snapshots:                              <span class="conum" data-value="1" />
      enabled: true
        persistentVolumeClaim:
          accessModes: [ "ReadWriteOnce" ]
          storageClassName: fast
          resources:
            requests:
              storage: 1Gi</markup>

<ul class="colist">
<li data-value="1">The implicit <code>storage</code> role has <code>snapshots</code> enabled and configured to use a PVC with custom <code>StorageClass</code>.</li>
</ul>
</div>

<h3 id="_snapshots_using_persistent_volumes_for_explicit_roles">Snapshots using Persistent Volumes for Explicit Roles</h3>
<div class="section">
<p>When configuring a <code>CoherenceCluster</code> with one or more explicit roles the <code>coherence.snapshots</code> configuration is added
directly to the configuration of each role in the <code>roles</code> list of the <code>CoherenceCluster</code>.
For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  roles:
    - role: data                               <span class="conum" data-value="1" />
      coherence:
        snapshots:
          enabled: true
            persistentVolumeClaim:
              accessModes: [ "ReadWriteOnce" ]
              storageClassName: fast
              resources:
                requests:
                  storage: 1Gi
    - role: proxy                              <span class="conum" data-value="2" />
      coherence:
        snapshots:
          enabled: false</markup>

<ul class="colist">
<li data-value="1">The <code>data</code> role has <code>snapshots</code> enabled and configured to use a PVC with custom <code>StorageClass</code>.</li>
<li data-value="2">The <code>proxy</code> role has <code>snapshots</code> explicitly disabled.</li>
</ul>
</div>

<h3 id="_snapshots_using_persistent_volumes_for_explicit_roles_with_defaults">Snapshots using Persistent Volumes for Explicit Roles with Defaults</h3>
<div class="section">
<p>When configuring a explicit roles in the <code>roles</code> list of a <code>CoherenceCluster</code> default values for the
<code>coherence.snapshots</code> configuration can be set in the <code>CoherenceCluster</code> <code>spec</code> section that will apply to
all roles in the <code>roles</code> list unless overridden for a specific role.
For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  coherence:
    snapshots:                                 <span class="conum" data-value="1" />
      enabled: true
        persistentVolumeClaim:
          accessModes: [ "ReadWriteOnce" ]
          storageClassName: fast
          resources:
            requests:
              storage: 1Gi
  roles:
    - role: data                               <span class="conum" data-value="2" />
    - role: proxy                              <span class="conum" data-value="3" />
      coherence:
        snapshots:
          enabled: false</markup>

<ul class="colist">
<li data-value="1">The default <code>snapshots</code> configuration is to enable snapshots using a PVC with custom <code>StorageClass</code>.</li>
<li data-value="2">The <code>data</code> role does not specify an explict <code>snapshots</code> configuration so it will use the defaults.</li>
<li data-value="3">The <code>proxy</code> role has <code>snapshots</code> explicitly disabled.</li>
</ul>
</div>
</div>

<h2 id="_managing_coherence_snapshots_using_standard_volumes">Managing Coherence Snapshots using Standard Volumes</h2>
<div class="section">
<p>Although <code>PersistentVolumeClaims</code> are the recommended way to manage storage for Coherence snapshots the Coherence
Operator also supports using standard Kubernetes <code>Volumes</code> as a storage mechanism.</p>

<div class="admonition warning">
<p class="admonition-inline">When using standard Kubernetes <code>Volumes</code> for snapshot storage it is important to ensure that
<code>CoherenceClusters</code> are configured and managed in such a way that the same <code>Volumes</code> are reattached to <code>Pods</code> if
clusters are restarted or if individual <code>Pods</code> are restarted or rescheduled by Kubernetes. If this is not done
then snapshot data can be lost. There are many ways to accomplish this using particular <code>Volume</code> types or controlling
<code>Pod</code> scheduling but this configuration is beyond the scope of this document and the relevant Kubernetes or storage
provider documentation should be consulted.</p>
</div>
<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  coherence:
    snapshots:
      enabled: true
        volume:          <span class="conum" data-value="1" />
          hostPath:
            path: /data</markup>

<ul class="colist">
<li data-value="1">Snapshots storage is configured to use a <code>hostPath</code> volume mapped to the <code>/data</code> directory on the host</li>
</ul>
<p>As with configuring snapshots to use <code>PersistentVolumeClaims</code> configuring them to use <code>Volumes</code> can be done at
different levels in the <code>CoherenceCluster</code> spec depending on whether there is a single implicit role, multiple
explicit roles and default values to apply to explicit roles.</p>

</div>

<h2 id="_managing_coherence_persistence_using_persistent_volumes">Managing Coherence Persistence using Persistent Volumes</h2>
<div class="section">
<p>When managing Coherence clusters using the Coherence Operator the simples configuration is to write persistence files
to a volume mapped to a <code>PersistentVolumeClaim</code> and to let the <code>StatefulSet</code> manage the <code>PVCs</code>.</p>

<p>Persistence is configured in the <code>coherence.persistence</code> section of the role specification in the <code>CoherenceCluster</code> CRD.</p>

<markup
lang="yaml"

>coherence:
  persistence:
    enabled: true            <span class="conum" data-value="1" />
      persistentVolumeClaim: <span class="conum" data-value="2" />
        # PVC spec...</markup>

<ul class="colist">
<li data-value="1">Persistence should be enabled by setting the <code>coherence.persistence.enabled</code> field to true.</li>
<li data-value="2">The <code>persistentVolumeClaim</code> section allows the <code>PVC</code> used for snapshot files to be configured.</li>
</ul>
<p>The default value for <code>coherence.persistence.enabled</code> is <code>false</code> so no snapshot location will be configured for Coherence
caches to use.</p>

<p>For example, if a Kubernetes cluster has a custom <code>StorageClass</code> named <code>fast</code> defined below:</p>

<markup
lang="yaml"

>kind: StorageClass
apiVersion: storage.k8s.io/v1
metadata:
  name: fast
provisioner: k8s.io/minikube-hostpath
parameters:
  type: pd-ssd</markup>

<p>Then a <code>CoherenceCluster</code> can be created with persistence enabled and configured to use the <code>fast</code> storage class:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  coherence:
    persistence:
      enabled: true
      persistentVolumeClaim:
        accessModes: [ "ReadWriteOnce" ]
        storageClassName: fast
        resources:
          requests:
            storage: 1Gi</markup>

<p>The content of the <code>persistentVolumeClaim</code> is any valid yaml for defining a <code>PersistentVolumeClaimSpec</code> that would be
allowed when configuring the <code>spec</code> section of a PVC in the <code>volumeClaimTemplates</code> section of a <code>StatefulSet</code> as
described in the <a id="" title="" target="_blank" href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.16/#persistentvolumeclaimspec-v1-core">Kubernetes API documentation</a>.</p>


<h3 id="_persistence_using_persistent_volumes_for_a_single_implicit_role">Persistence using Persistent Volumes for a Single Implicit Role</h3>
<div class="section">
<p>When configuring a <code>CoherenceCluster</code> with a single implicit role the <code>coherence.persistence</code> configuration is added
directly to the <code>spec</code> section of the <code>CoherenceCluster</code>.
For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  coherence:
    persistence:                              <span class="conum" data-value="1" />
      enabled: true
      persistentVolumeClaim:
        accessModes: [ "ReadWriteOnce" ]
        storageClassName: fast
        resources:
          requests:
            storage: 1Gi</markup>

<ul class="colist">
<li data-value="1">The implicit <code>storage</code> role has <code>persistence</code> enabled and configured to use a PVC with custom <code>StorageClass</code>.</li>
</ul>
</div>

<h3 id="_persistence_using_persistent_volumes_for_explicit_roles">Persistence using Persistent Volumes for Explicit Roles</h3>
<div class="section">
<p>When configuring a <code>CoherenceCluster</code> with one or more explicit roles the <code>coherence.persistence</code> configuration is added
directly to the configuration of each role in the <code>roles</code> list of the <code>CoherenceCluster</code>.
For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  roles:
    - role: data                               <span class="conum" data-value="1" />
      coherence:
        persistence:
          enabled: true
          persistentVolumeClaim:
            accessModes: [ "ReadWriteOnce" ]
            storageClassName: fast
            resources:
              requests:
                storage: 1Gi
    - role: proxy                              <span class="conum" data-value="2" />
      coherence:
        persistence:
          enabled: false</markup>

<ul class="colist">
<li data-value="1">The <code>data</code> role has <code>persistence</code> enabled and configured to use a PVC with custom <code>StorageClass</code>.</li>
<li data-value="2">The <code>proxy</code> role has <code>persistence</code> explicitly disabled.</li>
</ul>
</div>

<h3 id="_persistence_using_persistent_volumes_for_explicit_roles_with_defaults">Persistence using Persistent Volumes for Explicit Roles with Defaults</h3>
<div class="section">
<p>When configuring a explicit roles in the <code>roles</code> list of a <code>CoherenceCluster</code> default values for the
<code>coherence.persistence</code> configuration can be set in the <code>CoherenceCluster</code> <code>spec</code> section that will apply to
all roles in the <code>roles</code> list unless overridden for a specific role.
For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  coherence:
    persistence:                               <span class="conum" data-value="1" />
      enabled: true
      persistentVolumeClaim:
        accessModes: [ "ReadWriteOnce" ]
        storageClassName: fast
        resources:
          requests:
            storage: 1Gi
  roles:
    - role: data                               <span class="conum" data-value="2" />
    - role: proxy                              <span class="conum" data-value="3" />
      coherence:
        persistence:
          enabled: false</markup>

<ul class="colist">
<li data-value="1">The default <code>persistence</code> configuration is to enable persistence using a PVC with custom <code>StorageClass</code>.</li>
<li data-value="2">The <code>data</code> role does not specify an explict <code>persistence</code> configuration so it will use the defaults.</li>
<li data-value="3">The <code>proxy</code> role has <code>persistence</code> explicitly disabled.</li>
</ul>
</div>
</div>

<h2 id="_managing_coherence_persistence_using_standard_volumes">Managing Coherence Persistence using Standard Volumes</h2>
<div class="section">
<p>Although <code>PersistentVolumeClaims</code> are the recommended way to manage storage for Coherence persistence the Coherence
Operator also supports using standard Kubernetes <code>Volumes</code> as a storage mechanism.</p>

<div class="admonition warning">
<p class="admonition-inline">When using standard Kubernetes <code>Volumes</code> for snapshot storage it is important to ensure that
<code>CoherenceClusters</code> are configured and managed in such a way that the same <code>Volumes</code> are reattached to <code>Pods</code> if
clusters are restarted or if individual <code>Pods</code> are restarted or rescheduled by Kubernetes. If this is not done
then snapshot data can be lost. There are many ways to accomplish this using particular <code>Volume</code> types or controlling
<code>Pod</code> scheduling but this configuration is beyond the scope of this document and the relevant Kubernetes or storage
provider documentation should be consulted.</p>
</div>
<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  coherence:
    persistence:
      enabled: true
      volume:          <span class="conum" data-value="1" />
        hostPath:
          path: /data</markup>

<ul class="colist">
<li data-value="1">Snapshots storage is configured to use a <code>hostPath</code> volume mapped to the <code>/data</code> directory on the host</li>
</ul>
<p>As with configuring persistence to use <code>PersistentVolumeClaims</code> configuring them to use <code>Volumes</code> can be done at
different levels in the <code>CoherenceCluster</code> spec depending on whether there is a single implicit role, multiple
explicit roles and default values to apply to explicit roles.</p>

</div>
</doc-view>
