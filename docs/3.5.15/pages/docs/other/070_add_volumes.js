<doc-view>

<h2 id="_add_pod_volumes">Add Pod Volumes</h2>
<div class="section">
<p>Volumes and volume mappings can easily be added to a <code>Coherence</code> resource Pod to allow application code
deployed in the Pods to access additional storage.</p>

<p>Volumes are added by adding configuration to the <code>volumes</code> list in the <code>Coherence</code> CRD spec.
The configuration of the volume can be any valid yaml that would be used when adding a <code>Volume</code> to a <code>Pod</code> spec.</p>

<p>Volume mounts are added by adding configuration to the <code>volumeMounts</code> list in the <code>Coherence</code> CRD spec.
The configuration of the volume mount can be any valid yaml that would be used when adding a volume mount to a
container in a <code>Pod</code> spec.</p>

<div class="admonition note">
<p class="admonition-inline">Additional volumes added in this way will be added to all containers in the <code>Pod</code>.</p>
</div>
<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: storage
spec:
  volumes:
    - name: data-volume      <span class="conum" data-value="1" />
      nfs:
        path: /shared/data
        server: nfs-server
  volumeMounts:
    - name: data-volume      <span class="conum" data-value="2" />
      mountPath: /data</markup>

<ul class="colist">
<li data-value="1">An additional <code>Volume</code> named <code>data-volume</code> has been added (in this case the volume is an NFS volume).</li>
<li data-value="2">An additional volume mount has been added tthat will mount the <code>data-volume</code> at the <code>/data</code> mount point.</li>
</ul>
<p>The yaml above would result in a <code>Pod</code> spec similar to the following (a lot of the <code>Pod</code> spec has been omitted to just
show the relevant volume information):</p>

<markup
lang="yaml"

>apiVersion: v1
kind: Pod
metadata:
  name: storage-0
spec:
  containers:
    - name: coherence
      volumeMounts:
        - name: data-volume
          mountPath: /data
  volumes:
    - name: data-volume
      nfs:
        path: /shared/data
        server: nfs-server</markup>

</div>
</doc-view>
